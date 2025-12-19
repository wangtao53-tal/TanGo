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

	// 验证 base64 图片数据
	if err := utils.ValidateBase64Image(req.ImageData, maxSize); err != nil {
		return nil, err
	}

	// 解码 base64 数据
	imageData, err := base64.StdEncoding.DecodeString(req.ImageData)
	if err != nil {
		l.Errorw("Base64 解码失败",
			logx.Field("error", err),
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

	// 尝试上传到 GitHub
	var imageURL string
	var uploadMethod string

	if l.svcCtx.GitHubStorage != nil {
		url, err := l.svcCtx.GitHubStorage.Upload(imageData, filename)
		if err != nil {
			// GitHub 上传失败，记录日志但继续降级处理
			l.Errorw("GitHub 上传失败，降级到 base64",
				logx.Field("error", err),
				logx.Field("errorDetail", err.Error()),
			)
			// 降级到 base64 data URL
			uploadMethod = "base64"
			imageURL = fmt.Sprintf("data:image/jpeg;base64,%s", req.ImageData)
		} else {
			uploadMethod = "github"
			imageURL = url
		}
	} else {
		// GitHub 存储未初始化，直接使用 base64
		l.Infow("GitHub 存储未初始化，使用 base64",
			logx.Field("filename", filename),
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
