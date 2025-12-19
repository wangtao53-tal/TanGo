package storage

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

	// Base64 编码图片数据
	encoded := base64.StdEncoding.EncodeToString(imageData)

	// 构建请求
	reqBody := UploadRequest{
		Message: fmt.Sprintf("Upload image: %s", filename),
		Content: encoded,
		Branch:  g.config.GitHubBranch,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		g.logger.Errorw("序列化请求失败",
			logx.Field("error", err),
			logx.Field("filename", filename),
		)
		return "", utils.NewAPIError(500, "请求序列化失败", err.Error())
	}

	// 创建 HTTP 请求
	url := fmt.Sprintf("%s/%s", g.baseURL, filePath)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		g.logger.Errorw("创建请求失败",
			logx.Field("error", err),
			logx.Field("url", url),
		)
		return "", utils.NewAPIError(500, "创建请求失败", err.Error())
	}

	// 设置请求头
	req.Header.Set("Authorization", fmt.Sprintf("token %s", g.config.GitHubToken))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	// 发送请求
	resp, err := g.client.Do(req)
	if err != nil {
		g.logger.Errorw("GitHub API 请求失败",
			logx.Field("error", err),
			logx.Field("url", url),
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
