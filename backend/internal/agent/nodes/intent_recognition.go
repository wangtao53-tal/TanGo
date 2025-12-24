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

// IntentRecognitionNode 意图识别节点
type IntentRecognitionNode struct {
	ctx         context.Context
	config      config.AIConfig
	logger      logx.Logger
	chatModel   model.ChatModel     // eino ChatModel 实例
	template    prompt.ChatTemplate // 消息模板
	initialized bool
}

// IntentRecognitionResult 意图识别结果
type IntentRecognitionResult struct {
	Intent     string
	Confidence float64
	Reason     string
}

// NewIntentRecognitionNode 创建意图识别节点
func NewIntentRecognitionNode(ctx context.Context, cfg config.AIConfig, logger logx.Logger) (*IntentRecognitionNode, error) {
	node := &IntentRecognitionNode{
		ctx:    ctx,
		config: cfg,
		logger: logger,
	}

	// 如果配置了 eino 相关参数，初始化 ChatModel
	if cfg.EinoBaseURL != "" && cfg.AppID != "" && cfg.AppKey != "" {
		if err := node.initChatModel(ctx); err != nil {
			logger.Errorw("初始化ChatModel失败，将使用Mock模式", logx.Field("error", err))
			// 继续使用 Mock 模式
		} else {
			node.initialized = true
			logger.Info("意图识别节点已初始化ChatModel")
		}
	} else {
		logger.Info("未配置eino参数，意图识别节点将使用Mock模式")
	}

	// 创建消息模板
	node.initTemplate()

	return node, nil
}

// initChatModel 初始化 ChatModel（使用随机选择的模型）
func (n *IntentRecognitionNode) initChatModel(ctx context.Context) error {
	// 从配置中随机选择一个意图识别模型
	modelName := n.selectRandomModel(n.config.IntentModels)
	if modelName == "" {
		models := configpkg.GetDefaultIntentModels()
		if len(models) > 0 {
			modelName = n.selectRandomModel(models)
		}
		if modelName == "" {
			modelName = configpkg.DefaultIntentModel
		}
	}

	// 构建配置：使用 AppID 作为 APIKey（根据实际认证方式调整）
	cfg := &ark.ChatModelConfig{
		Model: modelName,
	}

	// 如果配置了 BaseURL，使用自定义地址
	if n.config.EinoBaseURL != "" {
		cfg.BaseURL = n.config.EinoBaseURL
	}

	// 认证：使用 Bearer Token 格式 ${TAL_MLOPS_APP_ID}:${TAL_MLOPS_APP_KEY}
	// 注意：eino 框架的 APIKey 字段可能已经处理了 Bearer Token 格式
	// 如果框架不支持，需要手动构造 Bearer Token
	if n.config.AppID != "" && n.config.AppKey != "" {
		// 使用 AppID:AppKey 格式作为 APIKey（eino 框架可能内部处理 Bearer Token）
		cfg.APIKey = n.config.AppID + ":" + n.config.AppKey
	} else if n.config.AppKey != "" {
		// 如果只有 AppKey，使用 AppKey
		cfg.APIKey = n.config.AppKey
	} else if n.config.AppID != "" {
		// 如果只有 AppID，使用 AppID
		cfg.APIKey = n.config.AppID
	} else {
		// 如果没有配置认证信息，无法初始化
		return nil // 返回 nil，使用 Mock 模式
	}

	chatModel, err := ark.NewChatModel(ctx, cfg)
	if err != nil {
		return err
	}

	n.chatModel = chatModel
	n.logger.Infow("意图识别模型已初始化", logx.Field("model", modelName))
	return nil
}

// selectRandomModel 从模型列表中随机选择一个模型
func (n *IntentRecognitionNode) selectRandomModel(models []string) string {
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
func (n *IntentRecognitionNode) initTemplate() {
	n.template = prompt.FromMessages(schema.FString,
		schema.SystemMessage(`你是一个意图识别助手。请识别用户消息的意图，并返回JSON格式的结果。

意图类型：
1. generate_cards: 用户想要生成知识卡片（例如："生成卡片"、"帮我生成小卡片"等）
2. text_response: 用户想要文本回答（其他所有情况）

请严格按照以下JSON格式返回：
{
  "intent": "generate_cards" 或 "text_response",
  "confidence": 0.0-1.0之间的浮点数,
  "reason": "识别原因"
}`),
		schema.MessagesPlaceholder("chat_history", true),
		schema.UserMessage("用户消息: {message}"),
	)
}

// Execute 执行意图识别
func (n *IntentRecognitionNode) Execute(data *GraphData, context []interface{}) (*IntentRecognitionResult, error) {
	n.logger.Infow("执行意图识别",
		logx.Field("message", data.Text),
		logx.Field("contextLength", len(context)),
		logx.Field("useRealModel", n.initialized),
	)

	// 如果 ChatModel 已初始化，使用真实模型
	if n.initialized && n.chatModel != nil {
		return n.executeReal(data, context)
	}

	// 否则使用 Mock 实现
	return n.executeMock(data, context)
}

// executeMock Mock实现（待替换为真实eino调用）
func (n *IntentRecognitionNode) executeMock(data *GraphData, context []interface{}) (*IntentRecognitionResult, error) {
	message := strings.ToLower(data.Text)

	// 规则判断：如果包含生成卡片相关关键词，识别为generate_cards意图
	generateCardKeywords := []string{
		"生成", "卡片", "小卡片", "知识卡片",
		"生成卡片", "帮我生成", "我要卡片",
		"create", "card", "generate", "cards",
	}

	for _, keyword := range generateCardKeywords {
		if strings.Contains(message, keyword) {
			result := &IntentRecognitionResult{
				Intent:     "generate_cards",
				Confidence: 0.9,
				Reason:     "检测到生成卡片关键词: " + keyword,
			}
			n.logger.Infow("意图识别完成（Mock-规则）",
				logx.Field("intent", result.Intent),
				logx.Field("confidence", result.Confidence),
			)
			return result, nil
		}
	}

	// 默认返回文本回答意图
	result := &IntentRecognitionResult{
		Intent:     "text_response",
		Confidence: 0.8,
		Reason:     "未检测到生成卡片意图，默认文本回答",
	}

	n.logger.Infow("意图识别完成（Mock-默认）",
		logx.Field("intent", result.Intent),
		logx.Field("confidence", result.Confidence),
	)
	return result, nil
}

// executeReal 真实eino实现
func (n *IntentRecognitionNode) executeReal(data *GraphData, context []interface{}) (*IntentRecognitionResult, error) {
	// 每次调用时重新初始化 ChatModel，使用随机选择的模型
	if err := n.initChatModel(n.ctx); err != nil {
		n.logger.Errorw("重新初始化ChatModel失败，使用已初始化的模型",
			logx.Field("error", err),
		)
		// 如果重新初始化失败，继续使用已初始化的模型
	}

	// 转换上下文为 eino Message 格式
	chatHistory := make([]*schema.Message, 0)
	for _, ctxItem := range context {
		// 尝试转换上下文项为 Message
		if msg, ok := ctxItem.(*schema.Message); ok {
			chatHistory = append(chatHistory, msg)
		}
		// TODO: 根据实际上下文类型进行更多转换逻辑
	}

	// 使用模板生成消息
	messages, err := n.template.Format(n.ctx, map[string]any{
		"message":      data.Text,
		"chat_history": chatHistory,
	})
	if err != nil {
		n.logger.Errorw("模板格式化失败", logx.Field("error", err))
		// 降级到 Mock
		return n.executeMock(data, context)
	}

	// 调用 ChatModel
	result, err := n.chatModel.Generate(n.ctx, messages)
	if err != nil {
		n.logger.Errorw("ChatModel调用失败", logx.Field("error", err))
		// 降级到 Mock
		return n.executeMock(data, context)
	}

	// 解析 JSON 结果
	var intentResult IntentRecognitionResult
	text := result.Content // Message.Content 是 string 类型

	// 尝试提取 JSON（可能包含 markdown 代码块）
	jsonStart := strings.Index(text, "{")
	jsonEnd := strings.LastIndex(text, "}")
	if jsonStart >= 0 && jsonEnd > jsonStart {
		jsonStr := text[jsonStart : jsonEnd+1]
		if err := json.Unmarshal([]byte(jsonStr), &intentResult); err != nil {
			n.logger.Errorw("解析JSON失败", logx.Field("error", err), logx.Field("text", text))
			// 降级到 Mock
			return n.executeMock(data, context)
		}
	} else {
		// 如果无法解析 JSON，尝试从文本中提取意图
		textLower := strings.ToLower(text)
		if strings.Contains(textLower, "generate_cards") || strings.Contains(textLower, "生成卡片") {
			intentResult.Intent = "generate_cards"
			intentResult.Confidence = 0.8
		} else {
			intentResult.Intent = "text_response"
			intentResult.Confidence = 0.8
		}
		intentResult.Reason = text
	}

	n.logger.Infow("意图识别完成（真实模型）",
		logx.Field("intent", intentResult.Intent),
		logx.Field("confidence", intentResult.Confidence),
	)

	return &intentResult, nil
}
