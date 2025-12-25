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

// LearningPlannerNode Learning Planner Agent节点
type LearningPlannerNode struct {
	ctx         context.Context
	config      config.AIConfig
	logger      logx.Logger
	chatModel   model.ChatModel     // eino ChatModel 实例
	template    prompt.ChatTemplate // 消息模板
	initialized bool
}

// NewLearningPlannerNode 创建Learning Planner Agent节点
func NewLearningPlannerNode(ctx context.Context, cfg config.AIConfig, logger logx.Logger) (*LearningPlannerNode, error) {
	node := &LearningPlannerNode{
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
			logger.Info("✅ Learning Planner Agent节点已初始化ChatModel，将使用真实模型")
		}
	} else {
		logger.Info("未配置eino参数，Learning Planner Agent节点将使用Mock模式")
	}

	// 创建消息模板
	node.initTemplate()

	return node, nil
}

// initChatModel 初始化 ChatModel
func (n *LearningPlannerNode) initChatModel(ctx context.Context) error {
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
func (n *LearningPlannerNode) selectRandomModel(models []string) string {
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
func (n *LearningPlannerNode) initTemplate() {
	n.template = prompt.FromMessages(schema.FString,
		schema.SystemMessage(`你是 Learning Planner Agent（像一位有经验的小学老师）。

输入包括：
- 意图判断（认知型、探因型、表达型、游戏型、情绪型）
- 认知负载建议（简短讲解、类比讲解、深入讲解、反问引导、暂停探索）
- 当前识别对象
- 孩子年龄段

你需要决定：
- 本轮是否继续深入
- 选择哪一个领域 Agent（Science、Language、Humanities）
- 是"讲一点"，还是"问一个问题"

重要规则：
- 你的输出是【下一步教学动作】，而不是知识本身
- 必须严格按照JSON格式返回
- 不要使用任何工具，只返回JSON结果

请严格按照以下JSON格式返回：
{{
  "continue": true或false,
  "domainAgent": "Science|Language|Humanities",
  "action": "讲一点|问一个问题"
}}`),
		schema.UserMessage(`意图判断: {intent}
认知负载建议: {cognitiveLoadAdvice}
识别对象: {objectName}（{objectCategory}）
孩子年龄: {userAge}岁`),
	)
}

// PlanLearning 制定学习计划
func (n *LearningPlannerNode) PlanLearning(ctx context.Context, intentResult *types.FollowUpIntentResult, cognitiveLoadAdvice *types.CognitiveLoadAdvice, objectName string, objectCategory string, userAge int) (*types.LearningPlanDecision, error) {
	n.logger.Infow("执行学习计划制定",
		logx.Field("intent", intentResult.Intent),
		logx.Field("strategy", cognitiveLoadAdvice.Strategy),
		logx.Field("objectName", objectName),
		logx.Field("userAge", userAge),
		logx.Field("useRealModel", n.initialized),
	)

	// 如果 ChatModel 已初始化，使用真实模型
	if n.initialized && n.chatModel != nil {
		return n.executeReal(ctx, intentResult, cognitiveLoadAdvice, objectName, objectCategory, userAge)
	}

	// 否则使用 Mock 实现
	return n.executeMock(intentResult, cognitiveLoadAdvice, objectName, objectCategory, userAge)
}

// executeMock Mock实现
func (n *LearningPlannerNode) executeMock(intentResult *types.FollowUpIntentResult, cognitiveLoadAdvice *types.CognitiveLoadAdvice, objectName string, objectCategory string, userAge int) (*types.LearningPlanDecision, error) {
	// 根据意图选择领域Agent
	var domainAgent string
	switch intentResult.Intent {
	case "探因型":
		domainAgent = "Science"
	case "表达型":
		domainAgent = "Language"
	case "游戏型", "情绪型":
		domainAgent = "Humanities"
	default:
		domainAgent = "Science" // 默认Science
	}

	// 根据认知负载建议决定动作
	action := "讲一点"
	if cognitiveLoadAdvice.Strategy == "反问引导" || cognitiveLoadAdvice.Strategy == "暂停探索" {
		action = "问一个问题"
	}

	return &types.LearningPlanDecision{
		Continue:    true,
		DomainAgent: domainAgent,
		Action:      action,
	}, nil
}

// executeReal 真实eino实现
func (n *LearningPlannerNode) executeReal(ctx context.Context, intentResult *types.FollowUpIntentResult, cognitiveLoadAdvice *types.CognitiveLoadAdvice, objectName string, objectCategory string, userAge int) (*types.LearningPlanDecision, error) {
	messages, err := n.template.Format(ctx, map[string]any{
		"intent":              intentResult.Intent,
		"cognitiveLoadAdvice": cognitiveLoadAdvice.Strategy,
		"objectName":          objectName,
		"objectCategory":      objectCategory,
		"userAge":             userAge,
	})
	if err != nil {
		n.logger.Errorw("模板格式化失败", logx.Field("error", err))
		return n.executeMock(intentResult, cognitiveLoadAdvice, objectName, objectCategory, userAge)
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
		return n.executeMock(intentResult, cognitiveLoadAdvice, objectName, objectCategory, userAge)
	}

	// 解析 JSON 结果
	var decision types.LearningPlanDecision
	text := result.Content

	jsonStart := strings.Index(text, "{")
	jsonEnd := strings.LastIndex(text, "}")
	if jsonStart >= 0 && jsonEnd > jsonStart {
		jsonStr := text[jsonStart : jsonEnd+1]
		if err := json.Unmarshal([]byte(jsonStr), &decision); err != nil {
			n.logger.Errorw("解析JSON失败", logx.Field("error", err), logx.Field("text", text))
			return n.executeMock(intentResult, cognitiveLoadAdvice, objectName, objectCategory, userAge)
		}
	} else {
		return n.executeMock(intentResult, cognitiveLoadAdvice, objectName, objectCategory, userAge)
	}

	// 验证领域Agent类型
	validDomainAgents := []string{"Science", "Language", "Humanities"}
	isValidDomain := false
	for _, validAgent := range validDomainAgents {
		if decision.DomainAgent == validAgent {
			isValidDomain = true
			break
		}
	}
	if !isValidDomain {
		return n.executeMock(intentResult, cognitiveLoadAdvice, objectName, objectCategory, userAge)
	}

	// 验证动作类型
	validActions := []string{"讲一点", "问一个问题"}
	isValidAction := false
	for _, validAction := range validActions {
		if decision.Action == validAction {
			isValidAction = true
			break
		}
	}
	if !isValidAction {
		return n.executeMock(intentResult, cognitiveLoadAdvice, objectName, objectCategory, userAge)
	}

	n.logger.Infow("学习计划制定完成（真实模型）",
		logx.Field("domainAgent", decision.DomainAgent),
		logx.Field("action", decision.Action),
		logx.Field("continue", decision.Continue),
	)

	return &decision, nil
}

