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

// ReflectionAgentNode Reflection Agent节点
type ReflectionAgentNode struct {
	ctx         context.Context
	config      config.AIConfig
	logger      logx.Logger
	chatModel   model.ChatModel     // eino ChatModel 实例
	template    prompt.ChatTemplate // 消息模板
	initialized bool
}

// NewReflectionAgentNode 创建Reflection Agent节点
func NewReflectionAgentNode(ctx context.Context, cfg config.AIConfig, logger logx.Logger) (*ReflectionAgentNode, error) {
	node := &ReflectionAgentNode{
		ctx:    ctx,
		config: cfg,
		logger: logger,
	}

	if cfg.EinoBaseURL != "" && cfg.AppID != "" && cfg.AppKey != "" {
		if err := node.initChatModel(ctx); err != nil {
			logger.Errorw("初始化ChatModel失败，将使用Mock模式", logx.Field("error", err))
		} else {
			node.initialized = true
			logger.Info("✅ Reflection Agent节点已初始化ChatModel")
		}
	} else {
		logger.Info("未配置eino参数，Reflection Agent节点将使用Mock模式")
	}

	node.initTemplate()
	return node, nil
}

// initChatModel 初始化 ChatModel
func (n *ReflectionAgentNode) initChatModel(ctx context.Context) error {
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

	if n.config.AppID != "" && n.config.AppKey != "" {
		cfg.APIKey = n.config.AppID + ":" + n.config.AppKey
	} else if n.config.AppKey != "" {
		cfg.APIKey = n.config.AppKey
	} else if n.config.AppID != "" {
		cfg.APIKey = n.config.AppID
	} else {
		return nil
	}

	chatModel, err := ark.NewChatModel(ctx, cfg)
	if err != nil {
		return err
	}

	n.chatModel = chatModel
	return nil
}

// selectRandomModel 从模型列表中随机选择一个模型
func (n *ReflectionAgentNode) selectRandomModel(models []string) string {
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
func (n *ReflectionAgentNode) initTemplate() {
	n.template = prompt.FromMessages(schema.FString,
		schema.SystemMessage(`你是 Reflection Agent。

判断孩子是否：
- 表现出兴趣
- 出现困惑
- 需要放松

根据对话历史和回答内容进行判断。

重要规则：
- 不要使用任何工具，只返回JSON结果
- 必须严格按照JSON格式返回

请严格按照以下JSON格式返回：
{{
  "interest": true或false,
  "confusion": true或false,
  "relax": true或false
}}`),
		schema.MessagesPlaceholder("chat_history", true),
		schema.UserMessage("回答内容: {content}"),
	)
}

// Reflect 反思判断
func (n *ReflectionAgentNode) Reflect(ctx context.Context, content string, conversationHistory []*schema.Message) (*types.ReflectionResult, error) {
	n.logger.Infow("执行Reflection Agent反思判断",
		logx.Field("contentLength", len(content)),
		logx.Field("conversationHistoryLength", len(conversationHistory)),
		logx.Field("useRealModel", n.initialized),
	)

	if n.initialized && n.chatModel != nil {
		return n.executeReal(ctx, content, conversationHistory)
	}

	return n.executeMock(content)
}

// executeMock Mock实现
func (n *ReflectionAgentNode) executeMock(content string) (*types.ReflectionResult, error) {
	// 简单规则判断
	contentLower := strings.ToLower(content)
	hasConfusion := strings.Contains(contentLower, "不懂") || strings.Contains(contentLower, "太难") || strings.Contains(contentLower, "不明白")
	hasInterest := !hasConfusion // 如果没有困惑，假设有兴趣

	return &types.ReflectionResult{
		Interest:  hasInterest,
		Confusion: hasConfusion,
		Relax:     false,
	}, nil
}

// executeReal 真实eino实现
func (n *ReflectionAgentNode) executeReal(ctx context.Context, content string, conversationHistory []*schema.Message) (*types.ReflectionResult, error) {
	messages, err := n.template.Format(ctx, map[string]any{
		"content":      content,
		"chat_history": conversationHistory,
	})
	if err != nil {
		n.logger.Errorw("模板格式化失败", logx.Field("error", err))
		return n.executeMock(content)
	}

	// 确保消息格式正确，移除任何可能导致工具调用错误的字段
	// 过滤消息，确保只包含有效的文本消息
	cleanMessages := make([]*schema.Message, 0, len(messages))
	for _, msg := range messages {
		if msg != nil && msg.Role != "" {
			// 创建干净的消息副本，只保留必要的字段
			cleanMsg := &schema.Message{
				Role:    msg.Role,
				Content: msg.Content,
			}
			cleanMessages = append(cleanMessages, cleanMsg)
		}
	}

	result, err := n.chatModel.Generate(ctx, cleanMessages)
	if err != nil {
		n.logger.Errorw("ChatModel调用失败", logx.Field("error", err))
		return n.executeMock(content)
	}

	// 解析 JSON 结果
	var reflectionResult types.ReflectionResult
	text := result.Content

	jsonStart := strings.Index(text, "{")
	jsonEnd := strings.LastIndex(text, "}")
	if jsonStart >= 0 && jsonEnd > jsonStart {
		jsonStr := text[jsonStart : jsonEnd+1]
		if err := json.Unmarshal([]byte(jsonStr), &reflectionResult); err != nil {
			n.logger.Errorw("解析JSON失败", logx.Field("error", err), logx.Field("text", text))
			return n.executeMock(content)
		}
	} else {
		return n.executeMock(content)
	}

	return &reflectionResult, nil
}

