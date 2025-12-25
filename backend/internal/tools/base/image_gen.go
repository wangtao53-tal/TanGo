package base

import (
	"context"
	"fmt"
	"time"

	"github.com/tango/explore/internal/tools"
	"github.com/zeromicro/go-zero/core/logx"
)

// ImageGenerateSimpleTool image_generate_simple工具实现
// 用于生成简单示意图，主要用于Science Agent
type ImageGenerateSimpleTool struct {
	logger logx.Logger
}

// NewImageGenerateSimpleTool 创建image_generate_simple工具实例
func NewImageGenerateSimpleTool(logger logx.Logger) tools.Tool {
	return &ImageGenerateSimpleTool{
		logger: logger,
	}
}

// Name 返回工具名称
func (t *ImageGenerateSimpleTool) Name() string {
	return "image_generate_simple"
}

// Description 返回工具描述
func (t *ImageGenerateSimpleTool) Description() string {
	return "生成简单示意图，用于科学知识讲解。输入描述，返回图片URL。"
}

// Parameters 返回工具参数定义（JSON Schema格式）
func (t *ImageGenerateSimpleTool) Parameters() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"description": map[string]interface{}{
				"type":        "string",
				"description": "图片描述，例如：'太阳系示意图'、'水循环过程图'",
			},
		},
		"required": []string{"description"},
	}
}

// Execute 执行工具
func (t *ImageGenerateSimpleTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	// 设置超时
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 提取参数
	description, ok := params["description"].(string)
	if !ok {
		return nil, fmt.Errorf("参数description必须是字符串类型")
	}

	if description == "" {
		return nil, fmt.Errorf("图片描述不能为空")
	}

	t.logger.Infow("执行image_generate_simple工具",
		logx.Field("description", description),
	)

	// TODO: 后续可以接入真实的图片生成API
	// 当前使用Mock实现
	result := t.generateImage(description)

	return result, nil
}

// generateImage 生成图片（Mock实现）
func (t *ImageGenerateSimpleTool) generateImage(description string) map[string]interface{} {
	// Mock图片URL（实际应该调用图片生成API）
	imageURL := fmt.Sprintf("https://example.com/images/%s.png", description)

	return map[string]interface{}{
		"description": description,
		"image_url":   imageURL,
		"status":      "generated",
		"note":        "这是Mock图片URL，实际使用时需要接入真实图片生成API",
	}
}

