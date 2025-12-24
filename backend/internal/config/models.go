package config

// 模型名称常量定义
// 所有模型名称统一在此管理，避免硬编码

const (
	// 意图识别模型
	DefaultIntentModel = "gpt-5-nano"

	// 图片识别模型（Vision模型）
	DefaultImageRecognitionModel1 = "doubao-seed-1.6-vision"
	DefaultImageRecognitionModel2 = "GLM-4.6v"

	// 图片生成模型
	DefaultImageGenerationModel = "gemini-3-pro-image"

	// 文本生成模型（默认使用意图识别模型）
	DefaultTextGenerationModel = "gpt-5-nano"
)

// GetDefaultImageRecognitionModels 获取默认图片识别模型列表
func GetDefaultImageRecognitionModels() []string {
	return []string{
		DefaultImageRecognitionModel1,
		DefaultImageRecognitionModel2,
	}
}
