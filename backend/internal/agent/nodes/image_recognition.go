package nodes

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
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
	hasEinoBaseURL := cfg.EinoBaseURL != ""
	hasAppID := cfg.AppID != ""
	hasAppKey := cfg.AppKey != ""

	if hasEinoBaseURL && hasAppID && hasAppKey {
		logger.Infow("检测到eino配置，尝试初始化Vision ChatModel",
			logx.Field("einoBaseURL", cfg.EinoBaseURL),
			logx.Field("appID", cfg.AppID),
			logx.Field("hasAppKey", hasAppKey),
		)
		if err := node.initChatModel(ctx); err != nil {
			logger.Errorw("初始化Vision ChatModel失败，将使用Mock模式",
				logx.Field("error", err),
				logx.Field("errorDetail", err.Error()),
			)
		} else {
			node.initialized = true
			logger.Info("✅ 图片识别节点已初始化Vision ChatModel，将使用真实模型")
		}
	} else {
		logger.Errorw("未完整配置eino参数，图片识别节点将使用Mock模式",
			logx.Field("hasEinoBaseURL", hasEinoBaseURL),
			logx.Field("hasAppID", hasAppID),
			logx.Field("hasAppKey", hasAppKey),
		)
		logger.Info("提示：需要同时配置 EINO_BASE_URL、TAL_MLOPS_APP_ID、TAL_MLOPS_APP_KEY 才能使用真实模型")
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

	// 认证：使用与其他节点一致的格式
	// eino 框架的 APIKey 字段会自动处理 Bearer Token 格式
	// 使用 AppID:AppKey 格式，框架会在内部添加 Bearer 前缀
	if n.config.AppID != "" && n.config.AppKey != "" {
		// 使用 AppID:AppKey 格式（与其他节点一致）
		cfg.APIKey = n.config.AppID + ":" + n.config.AppKey
		n.logger.Infow("使用 AppID:AppKey 作为 APIKey",
			logx.Field("appIDLength", len(n.config.AppID)),
			logx.Field("appKeyLength", len(n.config.AppKey)),
		)
	} else if n.config.AppKey != "" {
		// 如果只有 AppKey，直接使用（可能是完整的认证 token）
		cfg.APIKey = n.config.AppKey
		n.logger.Infow("使用 AppKey 作为 APIKey")
	} else if n.config.AppID != "" {
		// 如果只有 AppID，使用 AppID
		cfg.APIKey = n.config.AppID
		n.logger.Infow("使用 AppID 作为 APIKey")
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
	// 创建带超时的 context（60秒超时）
	ctx, cancel := context.WithTimeout(n.ctx, 60*time.Second)
	defer cancel()

	// 图片识别模板是静态的，不需要变量，直接构建消息
	// 避免使用模板格式化，因为模板中的JSON示例可能被误解析为变量
	messages := []*schema.Message{
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
	}

	// 构建多模态消息：添加图片
	// 图片数据可能是 base64 编码的字符串、data URL 或 HTTP URL，需要处理
	var imageURL string
	var mimeType string = "image/jpeg" // 默认类型

	if strings.HasPrefix(data.Image, "data:") {
		// 已经是 data URL 格式，直接使用
		imageURL = data.Image
		// 提取 MIME 类型
		parts := strings.SplitN(data.Image, ",", 2)
		if len(parts) == 2 {
			mimePart := strings.TrimSuffix(strings.SplitN(parts[0], ";", 2)[0], "data:")
			if mimePart != "" {
				mimeType = mimePart
			}
		}
	} else if strings.HasPrefix(data.Image, "http://") || strings.HasPrefix(data.Image, "https://") {
		// 如果是 HTTP/HTTPS URL，下载图片并转换为 base64 data URL
		n.logger.Infow("检测到图片URL，开始下载",
			logx.Field("url", data.Image),
		)
		downloadedBase64, downloadedMimeType, err := n.downloadImageAsBase64(ctx, data.Image)
		if err != nil {
			n.logger.Errorw("下载图片失败",
				logx.Field("url", data.Image),
				logx.Field("error", err),
				logx.Field("errorDetail", err.Error()),
			)
			return nil, fmt.Errorf("下载图片失败: %w", err)
		}
		if downloadedMimeType != "" {
			mimeType = downloadedMimeType
		}
		// 格式化为 data URL
		imageURL = fmt.Sprintf("data:%s;base64,%s", mimeType, downloadedBase64)
		n.logger.Infow("图片下载并转换完成",
			logx.Field("url", data.Image),
			logx.Field("mimeType", mimeType),
			logx.Field("base64Length", len(downloadedBase64)),
		)
	} else {
		// 假设是纯 base64 数据，格式化为 data URL
		imageURL = fmt.Sprintf("data:%s;base64,%s", mimeType, data.Image)
	}

	// 创建用户消息，包含图片和文本
	// 使用 URL 字段而不是 Base64Data，因为 eino 期望 data URL 格式
	userMsg := &schema.Message{
		Role: schema.User,
		UserInputMultiContent: []schema.MessageInputPart{
			{
				Type: schema.ChatMessagePartTypeImageURL,
				Image: &schema.MessageInputImage{
					MessagePartCommon: schema.MessagePartCommon{
						URL: &imageURL,
					},
					Detail: schema.ImageURLDetailAuto,
				},
			},
			{
				Type: schema.ChatMessagePartTypeText,
				Text: "请识别这张图片中的对象。",
			},
		},
	}
	messages = append(messages, userMsg)

	// 调用 ChatModel（使用带超时的 context）
	result, err := n.chatModel.Generate(ctx, messages)
	if err != nil {
		// 检查是否是超时错误
		if ctx.Err() == context.DeadlineExceeded {
			n.logger.Errorw("ChatModel调用超时",
				logx.Field("timeout", "60s"),
			)
			return nil, fmt.Errorf("图片识别超时: %w", err)
		}
		n.logger.Errorw("ChatModel调用失败",
			logx.Field("error", err),
			logx.Field("errorDetail", err.Error()),
			logx.Field("baseURL", n.config.EinoBaseURL),
			logx.Field("hasAppID", n.config.AppID != ""),
			logx.Field("hasAppKey", n.config.AppKey != ""),
			logx.Field("appIDLength", len(n.config.AppID)),
			logx.Field("appKeyLength", len(n.config.AppKey)),
		)
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

// downloadImageAsBase64 从 URL 下载图片并转换为 base64
func (n *ImageRecognitionNode) downloadImageAsBase64(ctx context.Context, url string) (base64Data string, mimeType string, err error) {
	// 创建 HTTP 请求，使用传入的 context
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", "", fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置 User-Agent，避免某些服务器拒绝请求
	req.Header.Set("User-Agent", "TanGo-ImageRecognition/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("下载图片失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("下载图片失败，状态码: %d", resp.StatusCode)
	}

	// 读取图片数据
	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("读取图片数据失败: %w", err)
	}

	// 获取 MIME 类型
	mimeType = resp.Header.Get("Content-Type")
	if mimeType == "" {
		// 根据 URL 扩展名推断 MIME 类型
		urlLower := strings.ToLower(url)
		if strings.HasSuffix(urlLower, ".png") {
			mimeType = "image/png"
		} else if strings.HasSuffix(urlLower, ".jpg") || strings.HasSuffix(urlLower, ".jpeg") {
			mimeType = "image/jpeg"
		} else if strings.HasSuffix(urlLower, ".gif") {
			mimeType = "image/gif"
		} else if strings.HasSuffix(urlLower, ".webp") {
			mimeType = "image/webp"
		} else {
			mimeType = "image/jpeg" // 默认类型
		}
	}

	// 转换为 base64
	base64Data = base64.StdEncoding.EncodeToString(imageData)

	return base64Data, mimeType, nil
}
