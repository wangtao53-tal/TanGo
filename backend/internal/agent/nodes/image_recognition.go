package nodes

import (
	"context"
	"encoding/json"
	"math/rand"
	"strings"
	"time"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
	"github.com/tango/explore/internal/config"
	configpkg "github.com/tango/explore/internal/config"
	"github.com/zeromicro/go-zero/core/logx"
)

// ImageRecognitionNode 图片识别节点
type ImageRecognitionNode struct {
	ctx         context.Context
	config      config.AIConfig
	logger      logx.Logger
	chatModel   model.ChatModel     // eino ChatModel 实例（支持 Vision）
	template    prompt.ChatTemplate // 消息模板
	initialized bool
}

// ImageRecognitionResult 图片识别结果
type ImageRecognitionResult struct {
	ObjectName     string
	ObjectCategory string
	Keywords       []string
	Confidence     float64
}

// NewImageRecognitionNode 创建图片识别节点
func NewImageRecognitionNode(ctx context.Context, cfg config.AIConfig, logger logx.Logger) (*ImageRecognitionNode, error) {
	node := &ImageRecognitionNode{
		ctx:    ctx,
		config: cfg,
		logger: logger,
	}

	// 如果配置了 eino 相关参数，初始化 ChatModel（Vision 模型）
	if cfg.EinoBaseURL != "" && cfg.AppID != "" && cfg.AppKey != "" {
		if err := node.initChatModel(ctx); err != nil {
			logger.Errorw("初始化Vision ChatModel失败，将使用Mock模式", logx.Field("error", err))
		} else {
			node.initialized = true
			logger.Info("图片识别节点已初始化Vision ChatModel")
		}
	} else {
		logger.Info("未配置eino参数，图片识别节点将使用Mock模式")
	}

	// 创建消息模板
	node.initTemplate()

	return node, nil
}

// initChatModel 初始化 Vision ChatModel
func (n *ImageRecognitionNode) initChatModel(ctx context.Context) error {
	// 从配置中选择一个图片识别模型
	modelName := ""
	if len(n.config.ImageRecognitionModels) > 0 {
		modelName = n.config.ImageRecognitionModels[0] // 使用第一个模型
	}
	if modelName == "" {
		modelName = configpkg.DefaultImageRecognitionModel1
	}

	cfg := &ark.ChatModelConfig{
		Model: modelName,
	}

	if n.config.EinoBaseURL != "" {
		cfg.BaseURL = n.config.EinoBaseURL
	}

	apiKey := n.config.AppKey
	if apiKey == "" {
		apiKey = n.config.AppID
	}
	if apiKey != "" {
		cfg.APIKey = apiKey
	} else {
		return nil // 返回 nil，使用 Mock 模式
	}

	chatModel, err := ark.NewChatModel(ctx, cfg)
	if err != nil {
		return err
	}

	n.chatModel = chatModel
	return nil
}

// initTemplate 初始化消息模板
func (n *ImageRecognitionNode) initTemplate() {
	n.template = prompt.FromMessages(schema.FString,
		schema.SystemMessage(`你是一个图片识别助手，专门识别图片中的对象。

请分析用户提供的图片，识别图片中的主要对象，并返回JSON格式的结果。

要求：
1. 识别图片中的主要对象名称（中文）
2. 判断对象类别：自然类、生活类、人文类
3. 提取3-5个相关关键词
4. 评估识别置信度（0.0-1.0）

请严格按照以下JSON格式返回：
{
  "objectName": "对象名称（中文）",
  "objectCategory": "自然类/生活类/人文类",
  "keywords": ["关键词1", "关键词2", "关键词3"],
  "confidence": 0.0-1.0之间的浮点数
}`),
		schema.UserMessage("请识别这张图片中的对象。"),
	)
}

// Execute 执行图片识别
func (n *ImageRecognitionNode) Execute(data *GraphData) (*ImageRecognitionResult, error) {
	n.logger.Infow("执行图片识别",
		logx.Field("imageLength", len(data.Image)),
		logx.Field("age", data.Age),
		logx.Field("useRealModel", n.initialized),
	)

	// 如果 ChatModel 已初始化，使用真实模型
	if n.initialized && n.chatModel != nil {
		return n.executeReal(data)
	}

	// 否则使用 Mock 实现
	return n.executeMock(data)
}

// executeMock Mock实现（待替换为真实eino调用）
func (n *ImageRecognitionNode) executeMock(data *GraphData) (*ImageRecognitionResult, error) {
	// Mock识别结果 - 随机返回一个常见对象
	mockObjects := []struct {
		name     string
		category string
		keywords []string
	}{
		{"银杏", "自然类", []string{"植物", "树木", "秋天", "叶子"}},
		{"苹果", "生活类", []string{"水果", "食物", "红色", "健康"}},
		{"蝴蝶", "自然类", []string{"昆虫", "飞行", "美丽", "春天"}},
		{"书本", "人文类", []string{"学习", "知识", "阅读", "教育"}},
		{"汽车", "生活类", []string{"交通工具", "速度", "现代", "出行"}},
		{"月亮", "自然类", []string{"天体", "夜晚", "圆形", "美丽"}},
		{"钢琴", "人文类", []string{"乐器", "音乐", "艺术", "优雅"}},
		{"太阳", "自然类", []string{"恒星", "光明", "温暖", "能量"}},
	}

	rand.Seed(time.Now().UnixNano())
	selected := mockObjects[rand.Intn(len(mockObjects))]

	// 生成随机置信度（0.85-0.99）
	confidence := 0.85 + rand.Float64()*0.14

	result := &ImageRecognitionResult{
		ObjectName:     selected.name,
		ObjectCategory: selected.category,
		Keywords:       selected.keywords,
		Confidence:     confidence,
	}

	n.logger.Infow("图片识别完成（Mock）",
		logx.Field("objectName", result.ObjectName),
		logx.Field("category", result.ObjectCategory),
		logx.Field("confidence", result.Confidence),
	)

	return result, nil
}

// executeReal 真实eino实现
func (n *ImageRecognitionNode) executeReal(data *GraphData) (*ImageRecognitionResult, error) {
	// 使用模板生成基础消息
	messages, err := n.template.Format(n.ctx, map[string]any{})
	if err != nil {
		n.logger.Errorw("模板格式化失败", logx.Field("error", err))
		return n.executeMock(data)
	}

	// 构建多模态消息：添加图片
	// 图片数据可能是 base64 编码的字符串，需要处理
	var imageBase64 string
	var mimeType string = "image/jpeg" // 默认类型

	if strings.HasPrefix(data.Image, "data:") {
		// 已经是 data URL 格式，提取 base64 部分
		parts := strings.SplitN(data.Image, ",", 2)
		if len(parts) == 2 {
			// 提取 MIME 类型
			mimePart := strings.TrimSuffix(strings.SplitN(parts[0], ";", 2)[0], "data:")
			if mimePart != "" {
				mimeType = mimePart
			}
			imageBase64 = parts[1]
		} else {
			imageBase64 = data.Image
		}
	} else {
		// 假设是纯 base64 数据
		imageBase64 = data.Image
	}

	// 修改用户消息，添加图片内容
	userMsg := messages[len(messages)-1] // 获取用户消息
	userMsg.UserInputMultiContent = []schema.MessageInputPart{
		{
			Type: schema.ChatMessagePartTypeImageURL,
			Image: &schema.MessageInputImage{
				MessagePartCommon: schema.MessagePartCommon{
					Base64Data: &imageBase64,
					MIMEType:   mimeType,
				},
				Detail: schema.ImageURLDetailAuto,
			},
		},
		{
			Type: schema.ChatMessagePartTypeText,
			Text: "请识别这张图片中的对象。",
		},
	}
	userMsg.Content = "" // 清空 Content，使用 MultiContent

	// 调用 ChatModel
	result, err := n.chatModel.Generate(n.ctx, messages)
	if err != nil {
		n.logger.Errorw("ChatModel调用失败", logx.Field("error", err))
		return n.executeMock(data)
	}

	// 解析 JSON 结果
	var recognitionResult ImageRecognitionResult
	text := result.Content

	// 尝试提取 JSON
	jsonStart := strings.Index(text, "{")
	jsonEnd := strings.LastIndex(text, "}")
	if jsonStart >= 0 && jsonEnd > jsonStart {
		jsonStr := text[jsonStart : jsonEnd+1]
		if err := json.Unmarshal([]byte(jsonStr), &recognitionResult); err != nil {
			n.logger.Errorw("解析JSON失败", logx.Field("error", err), logx.Field("text", text))
			return n.executeMock(data)
		}
	} else {
		// 无法解析 JSON，降级到 Mock
		n.logger.Errorw("无法从模型响应中提取JSON", logx.Field("text", text))
		return n.executeMock(data)
	}

	n.logger.Infow("图片识别完成（真实模型）",
		logx.Field("objectName", recognitionResult.ObjectName),
		logx.Field("category", recognitionResult.ObjectCategory),
		logx.Field("confidence", recognitionResult.Confidence),
	)

	return &recognitionResult, nil
}
