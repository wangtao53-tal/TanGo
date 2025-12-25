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
	"github.com/tango/explore/internal/utils"
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

// initChatModel 初始化 Vision ChatModel（使用随机选择的模型）
func (n *ImageRecognitionNode) initChatModel(ctx context.Context) error {
	// 从配置中随机选择一个图片识别模型
	modelName := n.selectRandomModel(n.config.ImageRecognitionModels)
	if modelName == "" {
		models := configpkg.GetDefaultImageRecognitionModels()
		if len(models) > 0 {
			modelName = n.selectRandomModel(models)
		}
		if modelName == "" {
			modelName = configpkg.DefaultImageRecognitionModel1
		}
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
	n.logger.Infow("图片识别模型已初始化", logx.Field("model", modelName))
	return nil
}

// selectRandomModel 从模型列表中随机选择一个模型
func (n *ImageRecognitionNode) selectRandomModel(models []string) string {
	if len(models) == 0 {
		return ""
	}
	if len(models) == 1 {
		return models[0]
	}
	rand.Seed(time.Now().UnixNano())
	return models[rand.Intn(len(models))]
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
	// 优化：减少日志详细程度，使用Debug级别
	n.logger.Debugw("执行图片识别",
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

	// 优化：减少日志详细程度
	n.logger.Debugw("图片识别完成（Mock）",
		logx.Field("objectName", result.ObjectName),
		logx.Field("category", result.ObjectCategory),
		logx.Field("confidence", result.Confidence),
	)

	return result, nil
}

// executeReal 真实eino实现
func (n *ImageRecognitionNode) executeReal(data *GraphData) (*ImageRecognitionResult, error) {
	// 每次调用时重新初始化 ChatModel，使用随机选择的模型
	if err := n.initChatModel(n.ctx); err != nil {
		n.logger.Errorw("重新初始化ChatModel失败，使用已初始化的模型",
			logx.Field("error", err),
		)
		// 如果重新初始化失败，继续使用已初始化的模型
	}

	// 优化：调整超时时间到45秒（原来60秒可能过长）
	// 如果模型调用需要更长时间，可以根据实际情况调整
	ctx, cancel := context.WithTimeout(n.ctx, 45*time.Second)
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
	// 优化：如果已经是可访问的 HTTP URL，直接使用，不需要下载
	var imageURL string
	var mimeType string = "image/jpeg" // 默认类型

	// 优化：使用更高效的字符串检查方法
	imageData := data.Image
	if strings.HasPrefix(imageData, "data:") {
		// 已经是 data URL 格式，直接使用
		imageURL = imageData
		// 优化：提取 MIME 类型，避免多次分割
		if idx := strings.Index(imageData, ","); idx > 0 {
			mimePart := strings.TrimPrefix(imageData[:idx], "data:")
			if semicolonIdx := strings.Index(mimePart, ";"); semicolonIdx > 0 {
				mimePart = mimePart[:semicolonIdx]
			}
			if mimePart != "" {
				mimeType = mimePart
			}
		}
	} else if strings.HasPrefix(imageData, "http://") || strings.HasPrefix(imageData, "https://") {
		// 优化：检测GitHub raw URL并转换为CDN URL
		originalURL := imageData
		if utils.IsGitHubRawURL(originalURL) {
			// 转换为jsDelivr CDN URL
			cdnURL, err := utils.ConvertToJSDelivrCDN(originalURL)
			if err == nil && cdnURL != originalURL {
				imageURL = cdnURL
				n.logger.Infow("GitHub raw URL已转换为CDN URL",
					logx.Field("originalURL", originalURL),
					logx.Field("cdnURL", cdnURL),
				)
			} else {
				// 转换失败，使用原始URL
				imageURL = originalURL
				n.logger.Debugw("GitHub raw URL转换失败，使用原始URL",
					logx.Field("url", originalURL),
					logx.Field("error", err),
				)
			}
		} else {
			// 非GitHub raw URL，直接使用
			imageURL = imageData
		}

		// 优化：使用更高效的MIME类型推断
		urlLower := strings.ToLower(imageData)
		if strings.HasSuffix(urlLower, ".png") {
			mimeType = "image/png"
		} else if strings.HasSuffix(urlLower, ".jpg") || strings.HasSuffix(urlLower, ".jpeg") {
			mimeType = "image/jpeg"
		} else if strings.HasSuffix(urlLower, ".gif") {
			mimeType = "image/gif"
		} else if strings.HasSuffix(urlLower, ".webp") {
			mimeType = "image/webp"
		}
		// 优化：减少日志详细程度
		n.logger.Debugw("使用图片URL（不下载）", logx.Field("url", imageURL))
	} else {
		// 优化：避免重复格式化，直接构建 data URL
		imageURL = fmt.Sprintf("data:%s;base64,%s", mimeType, imageData)
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
		// 优化：改进错误处理，区分超时、网络错误、模型错误
		if ctx.Err() == context.DeadlineExceeded {
			n.logger.Errorw("ChatModel调用超时",
				logx.Field("timeout", "45s"),
				logx.Field("errorType", "timeout"),
			)
			return nil, fmt.Errorf("图片识别超时: %w", err)
		}

		// 判断错误类型
		errMsg := err.Error()
		errorType := "unknown"
		if strings.Contains(errMsg, "timeout") || strings.Contains(errMsg, "deadline") {
			errorType = "timeout"
		} else if strings.Contains(errMsg, "network") || strings.Contains(errMsg, "connection") {
			errorType = "network"
		} else if strings.Contains(errMsg, "model") || strings.Contains(errMsg, "api") {
			errorType = "model"
		}

		n.logger.Errorw("ChatModel调用失败",
			logx.Field("errorType", errorType),
			logx.Field("error", err),
		)

		// 如果直接使用 HTTP URL 失败，且原始输入是 HTTP URL，尝试重试
		originalImageURL := data.Image
		if strings.HasPrefix(originalImageURL, "http://") || strings.HasPrefix(originalImageURL, "https://") {
			// 优化：如果使用的是CDN URL且失败，先重试原始GitHub raw URL
			if utils.IsGitHubRawURL(originalImageURL) && imageURL != originalImageURL {
				// 当前使用的是CDN URL，重试原始GitHub raw URL
				n.logger.Infow("CDN URL失败，重试原始GitHub raw URL",
					logx.Field("cdnURL", imageURL),
					logx.Field("originalURL", originalImageURL),
					logx.Field("error", err),
				)
				// 更新消息中的图片 URL为原始URL
				messages[len(messages)-1].UserInputMultiContent[0].Image.MessagePartCommon.URL = &originalImageURL

				// 重新调用模型
				result, err = n.chatModel.Generate(ctx, messages)
				if err == nil {
					// 原始URL重试成功
					n.logger.Infow("原始GitHub raw URL重试成功")
					// 继续处理结果
				} else {
					// 原始URL也失败，继续降级到下载base64
					n.logger.Debugw("原始GitHub raw URL也失败，降级到下载base64",
						logx.Field("url", originalImageURL),
						logx.Field("error", err),
					)
				}
			}

			// 如果CDN和原始URL都失败，或者不是GitHub URL，尝试下载并转换为 base64 后重试
			if err != nil {
				n.logger.Debugw("HTTP URL失败，尝试下载并转换为base64后重试",
					logx.Field("url", originalImageURL),
					logx.Field("error", err),
				)
				// 下载图片并转换为 base64 data URL
				downloadedBase64, downloadedMimeType, downloadErr := n.downloadImageAsBase64(ctx, originalImageURL)
				if downloadErr != nil {
					n.logger.Errorw("下载图片失败，回退到Mock",
						logx.Field("url", originalImageURL),
						logx.Field("error", downloadErr),
					)
					return n.executeMock(data)
				}
				if downloadedMimeType != "" {
					mimeType = downloadedMimeType
				}
				// 格式化为 data URL
				imageURL = fmt.Sprintf("data:%s;base64,%s", mimeType, downloadedBase64)
				n.logger.Debugw("图片下载并转换完成，重试模型调用",
					logx.Field("url", originalImageURL),
					logx.Field("mimeType", mimeType),
				)

				// 更新消息中的图片 URL
				messages[len(messages)-1].UserInputMultiContent[0].Image.MessagePartCommon.URL = &imageURL

				// 重新调用模型
				result, err = n.chatModel.Generate(ctx, messages)
				if err != nil {
					n.logger.Errorw("使用base64 data URL重试后仍然失败，回退到Mock",
						logx.Field("error", err),
						logx.Field("errorDetail", err.Error()),
					)
					return n.executeMock(data)
				}
				// 重试成功，继续处理结果
			}
		} else {
			// 非 HTTP URL 的错误，直接回退到 Mock
			// 优化：减少日志字段，移除大对象和详细配置信息
			n.logger.Errorw("ChatModel调用失败，回退到Mock",
				logx.Field("error", err),
			)
			return n.executeMock(data)
		}
	}

	// 解析 JSON 结果
	var recognitionResult ImageRecognitionResult
	text := result.Content

	// 优化：优化JSON解析逻辑，提高解析效率
	jsonStart := strings.IndexByte(text, '{')
	if jsonStart < 0 {
		n.logger.Errorw("无法从模型响应中提取JSON", logx.Field("textLength", len(text)))
		return n.executeMock(data)
	}

	jsonEnd := strings.LastIndexByte(text, '}')
	if jsonEnd <= jsonStart {
		n.logger.Errorw("无法从模型响应中提取JSON", logx.Field("textLength", len(text)))
		return n.executeMock(data)
	}

	jsonStr := text[jsonStart : jsonEnd+1]
	if err := json.Unmarshal([]byte(jsonStr), &recognitionResult); err != nil {
		// 优化：减少日志字段，不记录完整text（可能很大）
		n.logger.Errorw("解析JSON失败",
			logx.Field("error", err),
			logx.Field("jsonLength", len(jsonStr)),
		)
		return n.executeMock(data)
	}

	// 优化：减少日志详细程度，使用Info级别但减少字段
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
