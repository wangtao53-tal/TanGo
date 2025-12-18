package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	// AI模型配置
	AI AIConfig
}

// AIConfig AI模型配置
type AIConfig struct {
	// eino框架配置
	EinoBaseURL string `json:",optional,env=EINO_BASE_URL"`

	// AI模型APP ID和Key（用于Bearer Token认证）
	AppID  string `json:",optional,env=TAL_MLOPS_APP_ID"`
	AppKey string `json:",optional,env=TAL_MLOPS_APP_KEY"`

	// 意图识别模型
	IntentModel string `json:",optional,env=INTENT_MODEL,default=gpt-5-nano"`

	// 图片识别模型（随机选择其中一个）
	ImageRecognitionModels []string `json:",optional,env=IMAGE_RECOGNITION_MODELS"`

	// 图片生成模型
	ImageGenerationModel string `json:",optional,env=IMAGE_GENERATION_MODEL,default=Gemini 3 Pro Image"`

	// 文本生成模型
	TextGenerationModel string `json:",optional,env=TEXT_GENERATION_MODEL"`
}
