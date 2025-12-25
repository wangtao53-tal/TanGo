package config

// 模型名称常量定义
// 所有模型名称统一在此管理，避免硬编码

const (
	// 意图识别模型（默认值，兼容旧配置）
	DefaultIntentModel = "gpt-5-nano"

	// 图片识别模型（Vision模型）
	DefaultImageRecognitionModel1 = "doubao-seed-1.6-vision"
	DefaultImageRecognitionModel2 = "GLM-4.6v"

	// 图片生成模型
	DefaultImageGenerationModel = "gemini-3-pro-image"

	// 文本生成模型（默认值，兼容旧配置）
	DefaultTextGenerationModel = "gpt-5-nano"
)

// GetDefaultIntentModels 获取默认意图识别模型列表
func GetDefaultIntentModels() []string {
	return []string{
		"gemini-3-pro-image",
		"gpt-5-nano",
		"doubao-seededit-3-0-i2i",
		"doubao-seed-1.6vision",
		"glm-4.6v",
		"gpt-4o",
		"gemini-2.5-flash-preview",
		"gpt-5-pro",
		"gpt-5.1",
	}
}

// GetDefaultImageRecognitionModels 获取默认图片识别模型列表
func GetDefaultImageRecognitionModels() []string {
	return []string{
		"gemini-3-pro-image",
		"gpt-5-nano",
		"doubao-seededit-3-0-i2i",
		"doubao-seed-1.6vision",
		"glm-4.6v",
		"gpt-4o",
		"gemini-2.5-flash-preview",
		"gpt-5-pro",
		"gpt-5.1",
	}
}

// GetDefaultTextGenerationModels 获取默认文本生成模型列表
func GetDefaultTextGenerationModels() []string {
	return []string{
		"gemini-3-pro-image",
		"gpt-5-nano",
		"doubao-seededit-3-0-i2i",
		"doubao-seed-1.6vision",
		"glm-4.6v",
		"gpt-4o",
		"gemini-2.5-flash-preview",
		"gpt-5-pro",
		"gpt-5.1",
	}
}
