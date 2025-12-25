package nodes

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/schema"
	"github.com/tango/explore/internal/config"
	configpkg "github.com/tango/explore/internal/config"
	"github.com/zeromicro/go-zero/core/logx"
)

// ImageGenerationNode 图片生成节点
type ImageGenerationNode struct {
	ctx         context.Context
	config      config.AIConfig
	logger      logx.Logger
	imageModel  *ark.ImageGenerationModel // eino ImageGenerationModel 实例
	initialized bool
}

// NewImageGenerationNode 创建图片生成节点
func NewImageGenerationNode(ctx context.Context, cfg config.AIConfig, logger logx.Logger) (*ImageGenerationNode, error) {
	node := &ImageGenerationNode{
		ctx:    ctx,
		config: cfg,
		logger: logger,
	}

	// 如果配置了 eino 相关参数，初始化 ImageGenerationModel
	if cfg.EinoBaseURL != "" && cfg.AppID != "" && cfg.AppKey != "" {
		if err := node.initImageModel(ctx); err != nil {
			logger.Errorw("初始化ImageGenerationModel失败，将使用Mock模式", logx.Field("error", err))
		} else {
			node.initialized = true
			logger.Info("图片生成节点已初始化ImageGenerationModel")
		}
	} else {
		logger.Info("未配置eino参数，图片生成节点将使用Mock模式")
	}

	return node, nil
}

// initImageModel 初始化 ImageGenerationModel
func (n *ImageGenerationNode) initImageModel(ctx context.Context) error {
	modelName := n.config.ImageGenerationModel
	if modelName == "" {
		modelName = configpkg.DefaultImageGenerationModel
	}

	cfg := &ark.ImageGenerationConfig{
		Model: modelName,
	}

	if n.config.EinoBaseURL != "" {
		cfg.BaseURL = n.config.EinoBaseURL
	}

	// 认证：使用 Bearer Token 格式 ${TAL_MLOPS_APP_ID}:${TAL_MLOPS_APP_KEY}
	if n.config.AppID != "" && n.config.AppKey != "" {
		cfg.APIKey = n.config.AppID + ":" + n.config.AppKey
	} else if n.config.AppKey != "" {
		cfg.APIKey = n.config.AppKey
	} else if n.config.AppID != "" {
		cfg.APIKey = n.config.AppID
	} else {
		return nil // 返回 nil，使用 Mock 模式
	}

	imageModel, err := ark.NewImageGenerationModel(ctx, cfg)
	if err != nil {
		return err
	}

	n.imageModel = imageModel
	return nil
}

// GenerateCardImage 为卡片生成配图
func (n *ImageGenerationNode) GenerateCardImage(data *GraphData, card interface{}) (string, error) {
	n.logger.Infow("执行图片生成",
		logx.Field("objectName", data.ObjectName),
		logx.Field("cardType", n.getCardType(card)),
		logx.Field("useRealModel", n.initialized),
	)

	if n.initialized && n.imageModel != nil {
		return n.generateImageReal(data, card)
	}

	return n.generateImageMock(data, card)
}

// generateImageMock Mock实现（待替换为真实eino调用）
func (n *ImageGenerationNode) generateImageMock(data *GraphData, card interface{}) (string, error) {
	// Mock返回一个占位图片URL
	// 实际实现中，这里应该调用图片生成模型，生成符合卡片内容的配图
	imageURL := "https://via.placeholder.com/400x300?text=" + data.ObjectName

	n.logger.Infow("图片生成完成（Mock）", logx.Field("imageURL", imageURL))
	return imageURL, nil
}

// generateImageReal 真实eino实现
func (n *ImageGenerationNode) generateImageReal(data *GraphData, card interface{}) (string, error) {
	// 根据卡片类型生成不同的 prompt
	prompt := n.buildImagePrompt(data, card)

	// 构建消息
	messages := []*schema.Message{
		{
			Role:    schema.User,
			Content: prompt,
		},
	}

	// 调用 ImageGenerationModel
	result, err := n.imageModel.Generate(n.ctx, messages)
	if err != nil {
		n.logger.Errorw("ImageGenerationModel调用失败", logx.Field("error", err))
		return n.generateImageMock(data, card)
	}

	// 解析返回结果
	// ImageGenerationModel 返回的 Message 可能包含图片 URL 或 base64 数据
	// 需要根据实际返回格式解析
	imageURL := n.extractImageURL(result)
	if imageURL == "" {
		n.logger.Errorw("无法从模型响应中提取图片URL", logx.Field("result", result))
		return n.generateImageMock(data, card)
	}

	n.logger.Infow("图片生成完成（真实模型）", logx.Field("imageURL", imageURL))
	return imageURL, nil
}

// buildImagePrompt 根据卡片类型构建图片生成 prompt
func (n *ImageGenerationNode) buildImagePrompt(data *GraphData, card interface{}) string {
	cardType := n.getCardType(card)
	objectName := data.ObjectName

	switch cardType {
	case "science":
		return fmt.Sprintf("生成一张关于%s的科学知识卡片配图，风格：简洁、教育性、适合K12学生，色彩明亮", objectName)
	case "poetry":
		return fmt.Sprintf("生成一张关于%s的古诗词卡片配图，风格：中国风、诗意、古典、适合K12学生", objectName)
	case "english":
		return fmt.Sprintf("生成一张关于%s的英语学习卡片配图，风格：现代、简洁、国际化、适合K12学生", objectName)
	default:
		return fmt.Sprintf("生成一张关于%s的卡片配图，风格：简洁、适合K12学生", objectName)
	}
}

// extractImageURL 从模型响应中提取图片 URL
func (n *ImageGenerationNode) extractImageURL(result *schema.Message) string {
	// 检查 AssistantGenMultiContent
	if len(result.AssistantGenMultiContent) > 0 {
		for _, part := range result.AssistantGenMultiContent {
			if part.Image != nil {
				if part.Image.URL != nil {
					return *part.Image.URL
				}
				if part.Image.Base64Data != nil {
					// 如果是 base64 数据，转换为 data URL
					mimeType := part.Image.MIMEType
					if mimeType == "" {
						mimeType = "image/png"
					}
					return fmt.Sprintf("data:%s;base64,%s", mimeType, *part.Image.Base64Data)
				}
			}
		}
	}

	// 检查 Content 字段（可能包含 JSON）
	if result.Content != "" {
		var response map[string]interface{}
		if err := json.Unmarshal([]byte(result.Content), &response); err == nil {
			if url, ok := response["url"].(string); ok {
				return url
			}
			if imageURL, ok := response["image_url"].(string); ok {
				return imageURL
			}
		}
	}

	return ""
}

// getCardType 获取卡片类型（辅助函数）
func (n *ImageGenerationNode) getCardType(card interface{}) string {
	if cardMap, ok := card.(map[string]interface{}); ok {
		if cardType, ok := cardMap["type"].(string); ok {
			return cardType
		}
	}
	return "unknown"
}
