package nodes

import (
	"context"

	"github.com/tango/explore/internal/config"
	"github.com/zeromicro/go-zero/core/logx"
)

// ImageGenerationNode 图片生成节点
type ImageGenerationNode struct {
	ctx    context.Context
	config config.AIConfig
	logger logx.Logger
}

// NewImageGenerationNode 创建图片生成节点
func NewImageGenerationNode(ctx context.Context, cfg config.AIConfig, logger logx.Logger) (*ImageGenerationNode, error) {
	return &ImageGenerationNode{
		ctx:    ctx,
		config: cfg,
		logger: logger,
	}, nil
}

// GenerateCardImage 为卡片生成配图
func (n *ImageGenerationNode) GenerateCardImage(data *GraphData, card interface{}) (string, error) {
	n.logger.Infow("执行图片生成",
		logx.Field("objectName", data.ObjectName),
		logx.Field("cardType", n.getCardType(card)),
	)

	// TODO: 待APP ID提供后，接入真实eino框架调用图片生成模型
	// 当前使用Mock数据
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

// generateImageReal 真实eino实现（待APP ID提供后实现）
func (n *ImageGenerationNode) generateImageReal(data *GraphData, card interface{}) (string, error) {
	// TODO: 使用eino框架调用图片生成模型
	// 1. 从config中获取ImageGenerationModel配置
	// 2. 使用eino的ImageGenerationModel接口
	// 3. 根据卡片内容构建prompt
	// 4. 调用模型API生成图片
	// 5. 返回图片URL或base64数据

	return "", nil
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
