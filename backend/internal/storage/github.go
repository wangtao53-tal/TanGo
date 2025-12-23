package storage

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/tango/explore/internal/config"
	"github.com/tango/explore/internal/utils"
	"github.com/zeromicro/go-zero/core/logx"
)

// GitHubStorage GitHub 存储实现
type GitHubStorage struct {
	config  config.UploadConfig
	client  *http.Client
	logger  logx.Logger
	baseURL string
	rawURL  string
}

// NewGitHubStorage 创建新的 GitHub 存储实例
func NewGitHubStorage(cfg config.UploadConfig, logger logx.Logger) *GitHubStorage {
	// 设置默认值
	if cfg.GitHubBranch == "" {
		cfg.GitHubBranch = "main"
	}
	if cfg.GitHubPath == "" {
		cfg.GitHubPath = "images/"
	}
	if cfg.MaxImageSize == 0 {
		cfg.MaxImageSize = 10 * 1024 * 1024 // 10MB
	}

	// 构建 URL
	baseURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents", cfg.GitHubOwner, cfg.GitHubRepo)
	rawURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s", cfg.GitHubOwner, cfg.GitHubRepo, cfg.GitHubBranch)

	return &GitHubStorage{
		config:  cfg,
		client:  &http.Client{Timeout: 30 * time.Second},
		logger:  logger,
		baseURL: baseURL,
		rawURL:  rawURL,
	}
}

// UploadRequest GitHub API 上传请求
type UploadRequest struct {
	Message string `json:"message"`
	Content string `json:"content"` // base64 编码的文件内容
	Branch  string `json:"branch,omitempty"`
	SHA     string `json:"sha,omitempty"` // 更新已存在文件时需要
}

// FileInfo GitHub API 文件信息响应
type FileInfo struct {
	SHA  string `json:"sha"`
	Path string `json:"path"`
	Size int    `json:"size"`
}

// UploadResponse GitHub API 上传响应
type UploadResponse struct {
	Content struct {
		DownloadURL string `json:"download_url"`
		Path        string `json:"path"`
		SHA         string `json:"sha"`
	} `json:"content"`
	Commit struct {
		SHA string `json:"sha"`
	} `json:"commit"`
}

// getFileSHA 获取文件的 SHA 值（如果文件存在）
func (g *GitHubStorage) getFileSHA(filePath string) (string, error) {
	// 对文件路径进行 URL 编码
	encodedPath := strings.Split(filePath, "/")
	for i, part := range encodedPath {
		encodedPath[i] = url.PathEscape(part)
	}
	encodedFilePath := strings.Join(encodedPath, "/")

	// 构建 URL，添加 ref 参数指定分支
	apiURL := fmt.Sprintf("%s/%s?ref=%s", g.baseURL, encodedFilePath, url.QueryEscape(g.config.GitHubBranch))

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return "", err
	}

	// 设置请求头
	req.Header.Set("Authorization", fmt.Sprintf("token %s", g.config.GitHubToken))
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	// 发送请求
	resp, err := g.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 如果文件不存在，返回空字符串（不是错误）
	if resp.StatusCode == http.StatusNotFound {
		return "", nil
	}

	// 如果请求失败，返回错误
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		g.logger.Infow("获取文件 SHA 失败，继续尝试上传",
			logx.Field("status", resp.StatusCode),
			logx.Field("body", string(body)),
		)
		return "", nil // 如果获取失败，继续尝试上传（可能是新文件）
	}

	// 解析响应
	var fileInfo FileInfo
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil
	}

	if err := json.Unmarshal(body, &fileInfo); err != nil {
		return "", nil
	}

	return fileInfo.SHA, nil
}

// Upload 上传图片到 GitHub
func (g *GitHubStorage) Upload(imageData []byte, filename string) (string, error) {
	// 检查配置
	if g.config.GitHubToken == "" {
		return "", utils.NewAPIError(500, "GitHub token 未配置", "GITHUB_TOKEN 环境变量未设置")
	}
	if g.config.GitHubOwner == "" {
		return "", utils.NewAPIError(500, "GitHub owner 未配置", "GITHUB_OWNER 环境变量未设置")
	}
	if g.config.GitHubRepo == "" {
		return "", utils.NewAPIError(500, "GitHub repo 未配置", "GITHUB_REPO 环境变量未设置")
	}

	// 构建文件路径
	filePath := g.config.GitHubPath
	if filePath != "" && !strings.HasSuffix(filePath, "/") {
		filePath += "/"
	}
	filePath += filename

	// 检查文件是否已存在，如果存在则获取 SHA
	var fileSHA string
	sha, err := g.getFileSHA(filePath)
	if err != nil {
		g.logger.Infow("检查文件是否存在时出错，继续尝试上传",
			logx.Field("error", err),
			logx.Field("filePath", filePath),
		)
	} else {
		fileSHA = sha
		if fileSHA != "" {
			g.logger.Infow("文件已存在，将更新文件",
				logx.Field("filePath", filePath),
				logx.Field("sha", fileSHA),
			)
		}
	}

	// Base64 编码图片数据
	encoded := base64.StdEncoding.EncodeToString(imageData)

	// 构建请求
	reqBody := UploadRequest{
		Message: fmt.Sprintf("Upload image: %s", filename),
		Content: encoded,
		Branch:  g.config.GitHubBranch,
	}
	// 如果文件已存在，添加 SHA
	if fileSHA != "" {
		reqBody.SHA = fileSHA
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		g.logger.Errorw("序列化请求失败",
			logx.Field("error", err),
			logx.Field("filename", filename),
		)
		return "", utils.NewAPIError(500, "请求序列化失败", err.Error())
	}

	// 对文件路径进行 URL 编码
	encodedPath := strings.Split(filePath, "/")
	for i, part := range encodedPath {
		encodedPath[i] = url.PathEscape(part)
	}
	encodedFilePath := strings.Join(encodedPath, "/")

	// 创建 HTTP 请求
	apiURL := fmt.Sprintf("%s/%s", g.baseURL, encodedFilePath)
	req, err := http.NewRequest("PUT", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		g.logger.Errorw("创建请求失败",
			logx.Field("error", err),
			logx.Field("url", apiURL),
		)
		return "", utils.NewAPIError(500, "创建请求失败", err.Error())
	}

	// 设置请求头
	// 检查token是否是占位符
	tokenPreviewLen := 20
	if len(g.config.GitHubToken) < tokenPreviewLen {
		tokenPreviewLen = len(g.config.GitHubToken)
	}
	if strings.Contains(g.config.GitHubToken, "xxxxx") || strings.HasPrefix(g.config.GitHubToken, "ghp_xxxxxxxx") {
		g.logger.Errorw("GitHub token 是占位符，请配置真实的 token",
			logx.Field("tokenPreview", g.config.GitHubToken[:tokenPreviewLen]+"..."),
		)
		return "", utils.NewAPIError(500, "GitHub token 未正确配置", "请将 .env 文件中的 GITHUB_TOKEN 替换为真实的 token")
	}
	
	req.Header.Set("Authorization", fmt.Sprintf("token %s", g.config.GitHubToken))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	// 发送请求
	resp, err := g.client.Do(req)
	if err != nil {
		g.logger.Errorw("GitHub API 请求失败",
			logx.Field("error", err),
			logx.Field("url", apiURL),
		)
		return "", utils.NewAPIError(502, "GitHub API 请求失败", err.Error())
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		g.logger.Errorw("读取响应失败",
			logx.Field("error", err),
			logx.Field("status", resp.StatusCode),
		)
		return "", utils.NewAPIError(502, "读取响应失败", err.Error())
	}

	// 检查状态码
	if resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusOK {
		// 解析响应
		var uploadResp UploadResponse
		if err := json.Unmarshal(body, &uploadResp); err != nil {
			g.logger.Errorw("解析响应失败",
				logx.Field("error", err),
				logx.Field("body", string(body)),
			)
			// 即使解析失败，也可以使用构建的 URL
		}

		// 构建 raw URL
		imageURL := fmt.Sprintf("%s/%s", g.rawURL, filePath)

		g.logger.Infow("图片上传成功",
			logx.Field("filename", filename),
			logx.Field("url", imageURL),
			logx.Field("size", len(imageData)),
		)

		return imageURL, nil
	}

	// 处理错误响应
	if resp.StatusCode == http.StatusUnauthorized {
		g.logger.Errorw("GitHub 认证失败",
			logx.Field("status", resp.StatusCode),
			logx.Field("body", string(body)),
		)
		return "", utils.NewAPIError(401, "GitHub 认证失败", "请检查 GITHUB_TOKEN 是否正确")
	}

	if resp.StatusCode == http.StatusForbidden {
		// 可能是速率限制
		g.logger.Errorw("GitHub API 权限不足或速率限制",
			logx.Field("status", resp.StatusCode),
			logx.Field("body", string(body)),
		)
		return "", utils.ErrGitHubRateLimit
	}

	if resp.StatusCode == http.StatusNotFound {
		g.logger.Errorw("GitHub 仓库或路径不存在",
			logx.Field("status", resp.StatusCode),
			logx.Field("body", string(body)),
		)
		return "", utils.NewAPIError(404, "GitHub 仓库或路径不存在", "请检查仓库配置")
	}

	// 其他错误
	g.logger.Errorw("GitHub 上传失败",
		logx.Field("status", resp.StatusCode),
		logx.Field("body", string(body)),
	)
	return "", utils.NewAPIError(502, "GitHub 上传失败", fmt.Sprintf("状态码: %d, 响应: %s", resp.StatusCode, string(body)))
}
