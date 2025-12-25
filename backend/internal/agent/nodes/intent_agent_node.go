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
	"github.com/tango/explore/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

// IntentAgentNode Intent Agent节点（多Agent系统）
type IntentAgentNode struct {
	ctx         context.Context
	config      config.AIConfig
	logger      logx.Logger
	chatModel   model.ChatModel     // eino ChatModel 实例
	template    prompt.ChatTemplate // 消息模板
	initialized bool
}

// NewIntentAgentNode 创建Intent Agent节点
func NewIntentAgentNode(ctx context.Context, cfg config.AIConfig, logger logx.Logger) (*IntentAgentNode, error) {
	node := &IntentAgentNode{
		ctx:    ctx,
		config: cfg,
		logger: logger,
	}

	// 如果配置了 eino 相关参数，初始化 ChatModel
	if cfg.EinoBaseURL != "" && cfg.AppID != "" && cfg.AppKey != "" {
		if err := node.initChatModel(ctx); err != nil {
			logger.Errorw("初始化ChatModel失败，将使用Mock模式", logx.Field("error", err))
		} else {
			node.initialized = true
			logger.Info("✅ Intent Agent节点已初始化ChatModel，将使用真实模型")
		}
	} else {
		logger.Info("未配置eino参数，Intent Agent节点将使用Mock模式")
	}

	// 创建消息模板
	node.initTemplate()

	return node, nil
}

// initChatModel 初始化 ChatModel（使用随机选择的模型）
func (n *IntentAgentNode) initChatModel(ctx context.Context) error {
	// 从配置中随机选择一个文本生成模型
	modelName := n.selectRandomModel(n.config.TextGenerationModels)
	if modelName == "" {
		models := config.GetDefaultTextGenerationModels()
		if len(models) > 0 {
			modelName = n.selectRandomModel(models)
		}
		if modelName == "" {
			modelName = config.DefaultTextGenerationModel
		}
	}

	cfg := &ark.ChatModelConfig{
		Model: modelName,
	}

	if n.config.EinoBaseURL != "" {
		cfg.BaseURL = n.config.EinoBaseURL
	}

	// 认证：使用 Bearer Token 格式 ${TAL_MLOPS_APP_ID}:${TAL_MLOPS_APP_KEY}
	if n.config.AppID != "" && n.config.AppKey != "" {
		cfg.APIKey = n.config.AppID + ":" + n.config.AppKey
	} else if n.config.AppKey != "" {
		cfg.APIKey = n.config.AppKey
	} else if n.config.AppID != "" {
		cfg.APIKey = n.config.AppID
	} else {
		return nil // 返回 nil，使用 Mock 模式
	}

	chatModel, err := ark.NewChatModel(ctx, cfg)
	if err != nil {
		return err
	}

	n.chatModel = chatModel
	n.logger.Infow("Intent Agent模型已初始化", logx.Field("model", modelName))
	return nil
}

// selectRandomModel 从模型列表中随机选择一个模型
func (n *IntentAgentNode) selectRandomModel(models []string) string {
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
func (n *IntentAgentNode) initTemplate() {
	n.template = prompt.FromMessages(schema.FString,
		schema.SystemMessage(`你是 Intent Agent。

你的任务是从孩子的追问中判断主要意图类型。

意图类型：
1. 认知型：孩子想知道"这是什么"（例如："这是什么？"、"它叫什么？"）
2. 探因型：孩子想知道"为什么"或"怎么会"（例如："为什么？"、"怎么会这样？"、"它是怎么形成的？"）
3. 表达型：孩子想知道"怎么说"或"怎么形容"（例如："怎么说？"、"怎么形容？"、"用英语怎么说？"）
4. 游戏型：孩子想知道"好玩吗"或"能不能试试"（例如："好玩吗？"、"能不能试试？"、"我可以玩吗？"）
5. 情绪型：孩子表现出困惑或困难（例如："我不懂"、"太难了"、"我听不明白"）

重要规则：
- 你只输出意图标签和置信度，不生成教学内容
- 必须严格按照JSON格式返回

请严格按照以下JSON格式返回：
{{
  "intent": "认知型|探因型|表达型|游戏型|情绪型",
  "confidence": 0.0-1.0之间的浮点数,
  "reason": "识别原因（可选）"
}}`),
		schema.MessagesPlaceholder("chat_history", true),
		schema.UserMessage("{message}"),
	)
}

// RecognizeIntent 识别意图（多Agent系统）
func (n *IntentAgentNode) RecognizeIntent(ctx context.Context, message string, chatHistory []*schema.Message) (*types.FollowUpIntentResult, error) {
	n.logger.Infow("执行意图识别（多Agent系统）",
		logx.Field("message", message),
		logx.Field("chatHistoryLength", len(chatHistory)),
		logx.Field("useRealModel", n.initialized),
	)

	// 如果 ChatModel 已初始化，使用真实模型
	if n.initialized && n.chatModel != nil {
		return n.executeReal(ctx, message, chatHistory)
	}

	// 否则使用 Mock 实现
	return n.executeMock(message)
}

// executeMock Mock实现
func (n *IntentAgentNode) executeMock(message string) (*types.FollowUpIntentResult, error) {
	messageLower := strings.ToLower(message)

	// 规则判断
	if strings.Contains(messageLower, "为什么") || strings.Contains(messageLower, "怎么会") || strings.Contains(messageLower, "怎么形成") {
		return &types.FollowUpIntentResult{
			Intent:     "探因型",
			Confidence: 0.85,
			Reason:     "检测到探因型关键词",
		}, nil
	}

	if strings.Contains(messageLower, "怎么说") || strings.Contains(messageLower, "怎么形容") || strings.Contains(messageLower, "用英语") {
		return &types.FollowUpIntentResult{
			Intent:     "表达型",
			Confidence: 0.85,
			Reason:     "检测到表达型关键词",
		}, nil
	}

	if strings.Contains(messageLower, "好玩") || strings.Contains(messageLower, "试试") || strings.Contains(messageLower, "可以玩") {
		return &types.FollowUpIntentResult{
			Intent:     "游戏型",
			Confidence: 0.85,
			Reason:     "检测到游戏型关键词",
		}, nil
	}

	if strings.Contains(messageLower, "不懂") || strings.Contains(messageLower, "太难") || strings.Contains(messageLower, "不明白") {
		return &types.FollowUpIntentResult{
			Intent:     "情绪型",
			Confidence: 0.85,
			Reason:     "检测到情绪型关键词",
		}, nil
	}

	// 默认返回认知型
	return &types.FollowUpIntentResult{
		Intent:     "认知型",
		Confidence: 0.8,
		Reason:     "默认认知型意图",
	}, nil
}

// executeReal 真实eino实现
func (n *IntentAgentNode) executeReal(ctx context.Context, message string, chatHistory []*schema.Message) (*types.FollowUpIntentResult, error) {
	// 使用模板生成消息
	messages, err := n.template.Format(ctx, map[string]any{
		"message":      message,
		"chat_history": chatHistory,
	})
	if err != nil {
		n.logger.Errorw("模板格式化失败", logx.Field("error", err))
		return n.executeMock(message)
	}

	// 确保消息格式正确，移除任何可能导致工具调用错误的字段
	cleanMessages := make([]*schema.Message, 0, len(messages))
	for _, msg := range messages {
		if msg != nil && msg.Role != "" {
			cleanMsg := &schema.Message{
				Role:    msg.Role,
				Content: msg.Content,
			}
			cleanMessages = append(cleanMessages, cleanMsg)
		}
	}

	// 调用 ChatModel
	result, err := n.chatModel.Generate(ctx, cleanMessages)
	if err != nil {
		n.logger.Errorw("ChatModel调用失败", logx.Field("error", err))
		return n.executeMock(message)
	}

	// 解析 JSON 结果
	var intentResult types.FollowUpIntentResult
	text := result.Content

	// 尝试提取 JSON（可能包含 markdown 代码块）
	jsonStart := strings.Index(text, "{")
	jsonEnd := strings.LastIndex(text, "}")
	if jsonStart >= 0 && jsonEnd > jsonStart {
		jsonStr := text[jsonStart : jsonEnd+1]
		if err := json.Unmarshal([]byte(jsonStr), &intentResult); err != nil {
			n.logger.Errorw("解析JSON失败", logx.Field("error", err), logx.Field("text", text))
			return n.executeMock(message)
		}
	} else {
		// 如果无法解析 JSON，降级到 Mock
		return n.executeMock(message)
	}

	// 验证意图类型
	validIntents := []string{"认知型", "探因型", "表达型", "游戏型", "情绪型"}
	isValid := false
	for _, validIntent := range validIntents {
		if intentResult.Intent == validIntent {
			isValid = true
			break
		}
	}
	if !isValid {
		n.logger.Errorw("无效的意图类型", logx.Field("intent", intentResult.Intent))
		return n.executeMock(message)
	}

	n.logger.Infow("意图识别完成（真实模型）",
		logx.Field("intent", intentResult.Intent),
		logx.Field("confidence", intentResult.Confidence),
	)

	return &intentResult, nil
}

