package logic

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/tango/explore/internal/svc"
	"github.com/tango/explore/internal/types"
	"github.com/tango/explore/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadLogic {
	return &UploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadLogic) Upload(req *types.UploadRequest) (resp *types.UploadResponse, err error) {
	// 参数验证
	if req.ImageData == "" {
		return nil, utils.ErrImageDataRequired
	}

	// 清理 base64 字符串（在验证前清理，确保一致性）
	originalLength := len(req.ImageData)
	req.ImageData = utils.CleanBase64String(req.ImageData)
	
	// 清理后检查是否为空
	if req.ImageData == "" {
		l.Errorw("Base64 数据清理后为空",
			logx.Field("originalLength", originalLength),
		)
		return nil, utils.ErrImageDataRequired
	}
	
	// 记录清理信息（如果长度发生变化，说明有空白字符被移除）
	if len(req.ImageData) != originalLength {
		l.Infow("Base64 字符串已清理空白字符",
			logx.Field("originalLength", originalLength),
			logx.Field("cleanedLength", len(req.ImageData)),
			logx.Field("removedChars", originalLength-len(req.ImageData)),
		)
	}

	// 验证文件名（如果提供）
	if req.Filename != "" {
		if err := utils.ValidateFilename(req.Filename); err != nil {
			return nil, err
		}
	}

	// 获取配置
	maxSize := int64(10 * 1024 * 1024) // 默认 10MB
	if l.svcCtx.Config.Upload.MaxImageSize > 0 {
		maxSize = l.svcCtx.Config.Upload.MaxImageSize
	}

	// 验证 base64 图片数据（内部会再次清理，但我们已经清理过了，所以这里是安全的）
	if err := utils.ValidateBase64Image(req.ImageData, maxSize); err != nil {
		l.Errorw("Base64 图片验证失败",
			logx.Field("error", err),
			logx.Field("dataLength", len(req.ImageData)),
			logx.Field("dataPreview", getDataPreview(req.ImageData)),
		)
		return nil, err
	}

	// 解码 base64 数据（使用已清理的字符串）
	imageData, err := base64.StdEncoding.DecodeString(req.ImageData)
	if err != nil {
		l.Errorw("Base64 解码失败",
			logx.Field("error", err),
			logx.Field("dataLength", len(req.ImageData)),
			logx.Field("dataPreview", getDataPreview(req.ImageData)),
		)
		return nil, utils.ErrImageDataInvalid
	}

	// 生成文件名
	filename := req.Filename
	if filename == "" {
		// 从图片数据推断扩展名
		ext := ".jpg" // 默认 JPEG
		if len(imageData) >= 4 {
			if imageData[0] == 0x89 && imageData[1] == 0x50 && imageData[2] == 0x4E && imageData[3] == 0x47 {
				ext = ".png"
			} else if len(imageData) >= 12 && imageData[0] == 0x52 && imageData[1] == 0x49 && imageData[2] == 0x46 && imageData[3] == 0x46 {
				if imageData[8] == 0x57 && imageData[9] == 0x45 && imageData[10] == 0x42 && imageData[11] == 0x50 {
					ext = ".webp"
				}
			} else if len(imageData) >= 6 && imageData[0] == 0x47 && imageData[1] == 0x49 && imageData[2] == 0x46 && imageData[3] == 0x38 {
				ext = ".gif"
			}
		}
		filename = utils.GenerateFilename(ext)
	}

	l.Infow("开始上传图片",
		logx.Field("filename", filename),
		logx.Field("size", len(imageData)),
		logx.Field("hasGitHubStorage", l.svcCtx.GitHubStorage != nil),
	)

	// 优先上传到 GitHub（默认方式）
	var imageURL string
	var uploadMethod string

	if l.svcCtx.GitHubStorage != nil {
		// GitHub 存储已初始化，优先使用 GitHub 上传
		url, err := l.svcCtx.GitHubStorage.Upload(imageData, filename)
		if err != nil {
			// GitHub 上传失败，记录详细错误信息
			l.Errorw("GitHub 上传失败",
				logx.Field("error", err),
				logx.Field("errorDetail", err.Error()),
				logx.Field("filename", filename),
				logx.Field("size", len(imageData)),
			)
			// 如果GitHub配置了但上传失败，仍然降级到base64（保证功能可用）
			// 但记录警告，提示检查GitHub配置
			l.Infow("GitHub 上传失败，降级到 base64（请检查GitHub配置）",
				logx.Field("filename", filename),
			)
			uploadMethod = "base64"
			imageURL = fmt.Sprintf("data:image/jpeg;base64,%s", req.ImageData)
		} else {
			// GitHub 上传成功
			uploadMethod = "github"
			imageURL = url
			l.Infow("图片已上传到 GitHub",
				logx.Field("url", url),
				logx.Field("filename", filename),
			)
		}
	} else {
		// GitHub 存储未初始化，使用 base64（降级方案）
		l.Infow("GitHub 存储未配置，使用 base64 降级方案",
			logx.Field("filename", filename),
			logx.Field("hint", "如需使用GitHub存储，请配置GITHUB_TOKEN、GITHUB_OWNER、GITHUB_REPO"),
		)
		uploadMethod = "base64"
		imageURL = fmt.Sprintf("data:image/jpeg;base64,%s", req.ImageData)
	}

	// 构建响应
	resp = &types.UploadResponse{
		Url:          imageURL,
		Filename:     filename,
		Size:         len(imageData),
		UploadMethod: uploadMethod,
	}

	l.Infow("图片上传完成",
		logx.Field("filename", filename),
		logx.Field("url", imageURL),
		logx.Field("method", uploadMethod),
		logx.Field("size", len(imageData)),
	)

	return resp, nil
}

// getDataPreview 获取数据预览（用于日志，避免记录过长的base64字符串）
func getDataPreview(data string) string {
	if len(data) == 0 {
		return ""
	}
	previewLen := 50
	if len(data) < previewLen {
		previewLen = len(data)
	}
	return data[:previewLen] + "..."
}
