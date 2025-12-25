package nodes

import (
	"context"
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

// InteractionAgentNode Interaction Agent节点
type InteractionAgentNode struct {
	ctx         context.Context
	config      config.AIConfig
	logger      logx.Logger
	chatModel   model.ChatModel     // eino ChatModel 实例
	template    prompt.ChatTemplate // 消息模板
	initialized bool
}

// NewInteractionAgentNode 创建Interaction Agent节点
func NewInteractionAgentNode(ctx context.Context, cfg config.AIConfig, logger logx.Logger) (*InteractionAgentNode, error) {
	node := &InteractionAgentNode{
		ctx:    ctx,
		config: cfg,
		logger: logger,
	}

	if cfg.EinoBaseURL != "" && cfg.AppID != "" && cfg.AppKey != "" {
		if err := node.initChatModel(ctx); err != nil {
			logger.Errorw("初始化ChatModel失败，将使用Mock模式", logx.Field("error", err))
		} else {
			node.initialized = true
			logger.Info("✅ Interaction Agent节点已初始化ChatModel")
		}
	} else {
		logger.Info("未配置eino参数，Interaction Agent节点将使用Mock模式")
	}

	node.initTemplate()
	return node, nil
}

// initChatModel 初始化 ChatModel
func (n *InteractionAgentNode) initChatModel(ctx context.Context) error {
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
func (n *InteractionAgentNode) selectRandomModel(models []string) string {
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
func (n *InteractionAgentNode) initTemplate() {
	n.template = prompt.FromMessages(schema.FString,
		schema.SystemMessage(`你是 Interaction Agent，负责优化回答，让它更轻松友好。

重要规则：
- 直接优化回答内容，不要出现"跟小朋友可以这样聊"等指导性语言
- 不要出现"你:"这样的对话示例格式
- 把内容说"轻"，添加轻松友好的结尾
- 给孩子一个可选动作，不制造学习压力
- 常用的结尾方式：你想不想试试？我们下一步看什么？要不要换个角度？
- 让孩子感受到探索的乐趣
- 不要使用任何工具，只优化文本内容

记住：优化后的回答是直接给孩子看的，不是给家长看的指导手册！`),
		schema.UserMessage("原始回答: {content}"),
	)
}

// OptimizeInteraction 优化交互方式
func (n *InteractionAgentNode) OptimizeInteraction(ctx context.Context, content string) (*types.InteractionOptimization, error) {
	n.logger.Infow("执行Interaction Agent交互优化",
		logx.Field("contentLength", len(content)),
		logx.Field("useRealModel", n.initialized),
	)

	if n.initialized && n.chatModel != nil {
		return n.executeReal(ctx, content)
	}

	return n.executeMock(content)
}

// executeMock Mock实现
func (n *InteractionAgentNode) executeMock(content string) (*types.InteractionOptimization, error) {
	endings := []string{
		"你想不想试试？",
		"我们下一步看什么？",
		"要不要换个角度？",
	}
	ending := endings[rand.Intn(len(endings))]

	optimizedContent := content
	if !strings.HasSuffix(content, "？") && !strings.HasSuffix(content, "?") {
		optimizedContent = content + " " + ending
	}

	return &types.InteractionOptimization{
		OptimizedContent: optimizedContent,
		EndingAction:     ending,
	}, nil
}

// executeReal 真实eino实现
func (n *InteractionAgentNode) executeReal(ctx context.Context, content string) (*types.InteractionOptimization, error) {
	messages, err := n.template.Format(ctx, map[string]any{
		"content": content,
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

	optimizedContent := result.Content
	ending := ""
	if strings.Contains(optimizedContent, "你想不想试试？") {
		ending = "你想不想试试？"
	} else if strings.Contains(optimizedContent, "我们下一步看什么？") {
		ending = "我们下一步看什么？"
	} else if strings.Contains(optimizedContent, "要不要换个角度？") {
		ending = "要不要换个角度？"
	}

	return &types.InteractionOptimization{
		OptimizedContent: optimizedContent,
		EndingAction:     ending,
	}, nil
}

