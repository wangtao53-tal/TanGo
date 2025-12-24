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

// CognitiveLoadNode Cognitive Load Agent节点
type CognitiveLoadNode struct {
	ctx         context.Context
	config      config.AIConfig
	logger      logx.Logger
	chatModel   model.ChatModel     // eino ChatModel 实例（可选，用于复杂判断）
	template    prompt.ChatTemplate // 消息模板
	initialized bool
}

// NewCognitiveLoadNode 创建Cognitive Load Agent节点
func NewCognitiveLoadNode(ctx context.Context, cfg config.AIConfig, logger logx.Logger) (*CognitiveLoadNode, error) {
	node := &CognitiveLoadNode{
		ctx:    ctx,
		config: cfg,
		logger: logger,
	}

	// Cognitive Load Agent主要使用规则判断，ChatModel作为辅助
	// 如果配置了 eino 相关参数，初始化 ChatModel（用于复杂场景）
	if cfg.EinoBaseURL != "" && cfg.AppID != "" && cfg.AppKey != "" {
		if err := node.initChatModel(ctx); err != nil {
			logger.Errorw("初始化ChatModel失败，将仅使用规则判断", logx.Field("error", err))
		} else {
			node.initialized = true
			logger.Info("✅ Cognitive Load Agent节点已初始化ChatModel，将使用规则+模型判断")
		}
	} else {
		logger.Info("未配置eino参数，Cognitive Load Agent节点将仅使用规则判断")
	}

	// 创建消息模板（用于复杂场景）
	node.initTemplate()

	return node, nil
}

// initChatModel 初始化 ChatModel（可选）
func (n *CognitiveLoadNode) initChatModel(ctx context.Context) error {
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
func (n *CognitiveLoadNode) selectRandomModel(models []string) string {
	if len(models) == 0 {
		return ""
	}
	if len(models) == 1 {
		return models[0]
	}
	rand.Seed(time.Now().UnixNano())
	return models[rand.Intn(len(models))]
}

// initTemplate 初始化消息模板（用于复杂场景）
func (n *CognitiveLoadNode) initTemplate() {
	n.template = prompt.FromMessages(schema.FString,
		schema.SystemMessage(`你是 Cognitive Load Agent。

你的职责是防止信息过量。

根据以下信息判断当前最合适的输出策略：
- 孩子年龄
- 当前对话轮次
- 最近输出长度

输出策略：
1. 简短讲解：适合3-6岁，回答不超过3句话
2. 类比讲解：适合7-12岁，回答不超过5句话
3. 深入讲解：适合13-18岁，回答不超过7句话
4. 反问引导：连续追问超过5轮时使用
5. 暂停探索：最近输出超过500字时使用

重要规则：
- 不要生成知识内容，只给策略建议
- 必须严格按照JSON格式返回

请严格按照以下JSON格式返回：
{{
  "strategy": "简短讲解|类比讲解|深入讲解|反问引导|暂停探索",
  "reason": "建议理由",
  "maxSentences": 3或5或7（根据策略）
}}`),
		schema.UserMessage(`用户年龄: {userAge}岁
当前对话轮次: {conversationRounds}轮
最近输出长度: {recentOutputLength}字`),
	)
}

// AssessCognitiveLoad 评估认知负载
func (n *CognitiveLoadNode) AssessCognitiveLoad(ctx context.Context, userAge int, conversationRounds int, recentOutputLength int) (*types.CognitiveLoadAdvice, error) {
	n.logger.Infow("执行认知负载评估",
		logx.Field("userAge", userAge),
		logx.Field("conversationRounds", conversationRounds),
		logx.Field("recentOutputLength", recentOutputLength),
	)

	// 主要使用规则判断
	advice := n.assessByRules(userAge, conversationRounds, recentOutputLength)

	// 如果ChatModel已初始化，可以用于复杂场景的二次验证
	if n.initialized && n.chatModel != nil && (conversationRounds > 3 || recentOutputLength > 300) {
		// 复杂场景使用ChatModel辅助判断
		modelAdvice, err := n.assessByModel(ctx, userAge, conversationRounds, recentOutputLength)
		if err == nil && modelAdvice != nil {
			// 如果模型判断与规则判断一致，使用模型判断（更灵活）
			if modelAdvice.Strategy == advice.Strategy || n.isStrategyCompatible(modelAdvice.Strategy, advice.Strategy) {
				return modelAdvice, nil
			}
		}
	}

	return advice, nil
}

// assessByRules 使用规则判断认知负载
func (n *CognitiveLoadNode) assessByRules(userAge int, conversationRounds int, recentOutputLength int) *types.CognitiveLoadAdvice {
	// 规则1: 连续追问超过5轮 → 反问引导
	if conversationRounds > 5 {
		return &types.CognitiveLoadAdvice{
			Strategy:     "反问引导",
			Reason:       "连续追问超过5轮，建议反问引导孩子思考",
			MaxSentences: 2,
		}
	}

	// 规则2: 最近输出超过500字 → 暂停探索
	if recentOutputLength > 500 {
		return &types.CognitiveLoadAdvice{
			Strategy:     "暂停探索",
			Reason:       "最近输出超过500字，建议暂停探索，避免信息过载",
			MaxSentences: 1,
		}
	}

	// 规则3: 根据年龄选择策略
	if userAge <= 6 {
		return &types.CognitiveLoadAdvice{
			Strategy:     "简短讲解",
			Reason:       "3-6岁孩子，使用简短讲解策略",
			MaxSentences: 3,
		}
	} else if userAge <= 12 {
		return &types.CognitiveLoadAdvice{
			Strategy:     "类比讲解",
			Reason:       "7-12岁孩子，使用类比讲解策略",
			MaxSentences: 5,
		}
	} else {
		return &types.CognitiveLoadAdvice{
			Strategy:     "深入讲解",
			Reason:       "13-18岁孩子，使用深入讲解策略",
			MaxSentences: 7,
		}
	}
}

// assessByModel 使用ChatModel判断认知负载（复杂场景）
func (n *CognitiveLoadNode) assessByModel(ctx context.Context, userAge int, conversationRounds int, recentOutputLength int) (*types.CognitiveLoadAdvice, error) {
	messages, err := n.template.Format(ctx, map[string]any{
		"userAge":            userAge,
		"conversationRounds": conversationRounds,
		"recentOutputLength": recentOutputLength,
	})
	if err != nil {
		return nil, err
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
		return nil, err
	}

	// 解析 JSON 结果
	var advice types.CognitiveLoadAdvice
	text := result.Content

	jsonStart := strings.Index(text, "{")
	jsonEnd := strings.LastIndex(text, "}")
	if jsonStart >= 0 && jsonEnd > jsonStart {
		jsonStr := text[jsonStart : jsonEnd+1]
		if err := json.Unmarshal([]byte(jsonStr), &advice); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	// 验证策略类型
	validStrategies := []string{"简短讲解", "类比讲解", "深入讲解", "反问引导", "暂停探索"}
	isValid := false
	for _, validStrategy := range validStrategies {
		if advice.Strategy == validStrategy {
			isValid = true
			break
		}
	}
	if !isValid {
		return nil, err
	}

	return &advice, nil
}

// isStrategyCompatible 判断两个策略是否兼容
func (n *CognitiveLoadNode) isStrategyCompatible(strategy1, strategy2 string) bool {
	// 简短讲解、类比讲解、深入讲解可以互相兼容
	compatibleGroups := [][]string{
		{"简短讲解", "类比讲解", "深入讲解"},
		{"反问引导"},
		{"暂停探索"},
	}

	for _, group := range compatibleGroups {
		has1 := false
		has2 := false
		for _, s := range group {
			if s == strategy1 {
				has1 = true
			}
			if s == strategy2 {
				has2 = true
			}
		}
		if has1 && has2 {
			return true
		}
	}
	return false
}

