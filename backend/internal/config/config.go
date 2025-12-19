package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	// AI模型配置
	AI AIConfig
	// 图片上传配置
	Upload UploadConfig
}

// AIConfig AI模型配置
type AIConfig struct {
	// eino框架配置
	EinoBaseURL string `json:",optional,env=EINO_BASE_URL"`

	// AI模型APP ID和Key（用于Bearer Token认证）
	AppID  string `json:",optional,env=TAL_MLOPS_APP_ID"`
	AppKey string `json:",optional,env=TAL_MLOPS_APP_KEY"`

	// 意图识别模型（从环境变量 INTENT_MODEL 读取，未设置则使用默认值）
	IntentModel string `json:",optional,env=INTENT_MODEL"`

	// 图片识别模型（从环境变量 IMAGE_RECOGNITION_MODELS 读取，逗号分隔，未设置则使用默认值）
	// 注意：YAML 文件中不配置此字段，避免类型解析问题
	ImageRecognitionModels []string `json:",optional" yaml:",omitempty"`

	// 图片生成模型（从环境变量 IMAGE_GENERATION_MODEL 读取，未设置则使用默认值）
	ImageGenerationModel string `json:",optional,env=IMAGE_GENERATION_MODEL"`

	// 文本生成模型（从环境变量 TEXT_GENERATION_MODEL 读取，未设置则使用默认值）
	TextGenerationModel string `json:",optional,env=TEXT_GENERATION_MODEL"`
}

// UploadConfig 图片上传配置
type UploadConfig struct {
	// GitHub 配置
	GitHubToken  string `json:",optional,env=GITHUB_TOKEN"`  // GitHub Personal Access Token
	GitHubOwner  string `json:",optional,env=GITHUB_OWNER"`  // GitHub 用户名或组织名
	GitHubRepo   string `json:",optional,env=GITHUB_REPO"`   // GitHub 仓库名
	GitHubBranch string `json:",optional,env=GITHUB_BRANCH"` // GitHub 分支名，默认 "main"
	GitHubPath   string `json:",optional,env=GITHUB_PATH"`   // 图片存储路径，默认 "images/"
	// 图片大小限制（字节），默认 10MB
	MaxImageSize int64 `json:",optional,env=MAX_IMAGE_SIZE"`
}
