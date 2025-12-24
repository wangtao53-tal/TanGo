package nodes

import (
	"context"
	"math/rand"
	"time"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
	"github.com/tango/explore/internal/config"
	"github.com/tango/explore/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

// HumanitiesAgentNode Humanities Agent节点
type HumanitiesAgentNode struct {
	ctx         context.Context
	config      config.AIConfig
	logger      logx.Logger
	chatModel   model.ChatModel     // eino ChatModel 实例
	template    prompt.ChatTemplate // 消息模板
	initialized bool
}

// NewHumanitiesAgentNode 创建Humanities Agent节点
func NewHumanitiesAgentNode(ctx context.Context, cfg config.AIConfig, logger logx.Logger) (*HumanitiesAgentNode, error) {
	node := &HumanitiesAgentNode{
		ctx:    ctx,
		config: cfg,
		logger: logger,
	}

	if cfg.EinoBaseURL != "" && cfg.AppID != "" && cfg.AppKey != "" {
		if err := node.initChatModel(ctx); err != nil {
			logger.Errorw("初始化ChatModel失败，将使用Mock模式", logx.Field("error", err))
		} else {
			node.initialized = true
			logger.Info("✅ Humanities Agent节点已初始化ChatModel")
		}
	} else {
		logger.Info("未配置eino参数，Humanities Agent节点将使用Mock模式")
	}

	node.initTemplate()
	return node, nil
}

// initChatModel 初始化 ChatModel
func (n *HumanitiesAgentNode) initChatModel(ctx context.Context) error {
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
func (n *HumanitiesAgentNode) selectRandomModel(models []string) string {
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
func (n *HumanitiesAgentNode) initTemplate() {
	n.template = prompt.FromMessages(schema.FString,
		schema.SystemMessage(`你是 Humanities Agent，一个直接和孩子对话的AI伙伴，把自然与文化连接起来。

重要规则：
- 直接回答孩子的问题，就像朋友聊天一样
- 不要出现"跟小朋友可以这样聊"、"你可以说"等指导性语言
- 不要出现"你:"这样的对话示例格式
- 用"我"或直接称呼"你"（孩子）来对话
- 把自然与文化连接起来：一句诗、一个故事、一个画面感
- 不要求背诵，必须和眼前看到的事物有关
- 让孩子感受到文化的魅力

记住：你是直接和孩子对话的AI伙伴，不是给家长看的指导手册！`),
		schema.MessagesPlaceholder("chat_history", true),
		schema.UserMessage("{message}"),
	)
}

// GenerateHumanitiesAnswer 生成人文回答
func (n *HumanitiesAgentNode) GenerateHumanitiesAnswer(ctx context.Context, message string, objectName string, objectCategory string, userAge int, chatHistory []*schema.Message) (*types.DomainAgentResponse, error) {
	n.logger.Infow("执行Humanities Agent回答生成",
		logx.Field("message", message),
		logx.Field("objectName", objectName),
		logx.Field("userAge", userAge),
		logx.Field("useRealModel", n.initialized),
	)

	if n.initialized && n.chatModel != nil {
		return n.executeReal(ctx, message, objectName, objectCategory, userAge, chatHistory)
	}

	return n.executeMock(message, objectName, userAge)
}

// executeMock Mock实现
func (n *HumanitiesAgentNode) executeMock(message string, objectName string, userAge int) (*types.DomainAgentResponse, error) {
	content := "关于" + objectName + "，有一句古诗说得很美。"
	if userAge <= 6 {
		content = "看到" + objectName + "，我想起了一个有趣的故事。"
	} else if userAge <= 12 {
		content = objectName + "在古诗词中经常出现，比如..."
	} else {
		content = objectName + "承载着丰富的文化内涵，让我们一起来探索。"
	}

	return &types.DomainAgentResponse{
		DomainType:  "Humanities",
		Content:     content,
		ToolsUsed:   []string{},
		ToolResults: make(map[string]interface{}),
	}, nil
}

// executeReal 真实eino实现
func (n *HumanitiesAgentNode) executeReal(ctx context.Context, message string, objectName string, objectCategory string, userAge int, chatHistory []*schema.Message) (*types.DomainAgentResponse, error) {
	messages, err := n.template.Format(ctx, map[string]any{
		"message":       message,
		"objectName":    objectName,
		"objectCategory": objectCategory,
		"userAge":       userAge,
		"chat_history":  chatHistory,
	})
	if err != nil {
		n.logger.Errorw("模板格式化失败", logx.Field("error", err))
		return n.executeMock(message, objectName, userAge)
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

	result, err := n.chatModel.Generate(ctx, cleanMessages)
	if err != nil {
		n.logger.Errorw("ChatModel调用失败", logx.Field("error", err))
		return n.executeMock(message, objectName, userAge)
	}

	return &types.DomainAgentResponse{
		DomainType:  "Humanities",
		Content:     result.Content,
		ToolsUsed:   []string{},
		ToolResults: make(map[string]interface{}),
	}, nil
}

