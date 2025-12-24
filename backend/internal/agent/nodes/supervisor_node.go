package nodes

import (
	"context"

	"github.com/cloudwego/eino/schema"
	"github.com/tango/explore/internal/config"
	"github.com/tango/explore/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

// SupervisorNode Supervisor节点（多Agent系统）
type SupervisorNode struct {
	ctx                context.Context
	config             config.AIConfig
	logger             logx.Logger
	intentAgent        *IntentAgentNode
	cognitiveLoadAgent *CognitiveLoadNode
	learningPlannerAgent *LearningPlannerNode
}

// NewSupervisorNode 创建Supervisor节点
func NewSupervisorNode(ctx context.Context, cfg config.AIConfig, logger logx.Logger, intentAgent *IntentAgentNode, cognitiveLoadAgent *CognitiveLoadNode, learningPlannerAgent *LearningPlannerNode) (*SupervisorNode, error) {
	node := &SupervisorNode{
		ctx:                ctx,
		config:             cfg,
		logger:             logger,
		intentAgent:        intentAgent,
		cognitiveLoadAgent: cognitiveLoadAgent,
		learningPlannerAgent: learningPlannerAgent,
	}

	logger.Info("✅ Supervisor节点已初始化")
	return node, nil
}

// Coordinate 协调多Agent协作
func (n *SupervisorNode) Coordinate(ctx context.Context, state *types.SupervisorState, message string, chatHistory []*schema.Message) (*types.LearningPlanDecision, error) {
	n.logger.Infow("Supervisor开始协调多Agent协作",
		logx.Field("objectName", state.ObjectName),
		logx.Field("userAge", state.UserAge),
		logx.Field("conversationRounds", state.ConversationRounds),
	)

	// 确保AgentResults已初始化
	if state.AgentResults == nil {
		state.AgentResults = make(map[string]interface{})
	}

	// 1. 调用Intent Agent识别意图
	intentResult, err := n.intentAgent.RecognizeIntent(ctx, message, chatHistory)
	if err != nil {
		n.logger.Errorw("Intent Agent调用失败", logx.Field("error", err))
		// 降级处理：使用默认意图
		intentResult = &types.FollowUpIntentResult{
			Intent:     "认知型",
			Confidence: 0.5,
			Reason:     "Intent Agent调用失败，使用默认意图",
		}
	}
	state.AgentResults["intent"] = intentResult

	// 2. 调用Cognitive Load Agent判断认知负载
	cognitiveLoadAdvice, err := n.cognitiveLoadAgent.AssessCognitiveLoad(ctx, state.UserAge, state.ConversationRounds, state.RecentOutputLength)
	if err != nil {
		n.logger.Errorw("Cognitive Load Agent调用失败", logx.Field("error", err))
		// 降级处理：使用默认策略
		cognitiveLoadAdvice = &types.CognitiveLoadAdvice{
			Strategy:     "类比讲解",
			Reason:       "Cognitive Load Agent调用失败，使用默认策略",
			MaxSentences: 5,
		}
	}
	state.AgentResults["cognitiveLoad"] = cognitiveLoadAdvice

	// 3. 调用Learning Planner Agent制定学习计划
	decision, err := n.learningPlannerAgent.PlanLearning(ctx, intentResult, cognitiveLoadAdvice, state.ObjectName, state.ObjectCategory, state.UserAge)
	if err != nil {
		n.logger.Errorw("Learning Planner Agent调用失败", logx.Field("error", err))
		return nil, err
	}
	state.AgentResults["learningPlan"] = decision

	// 4. 智能工具选择：根据意图、问题内容和领域Agent选择工具
	selectedTools, toolStrategy := SelectTools(
		intentResult.Intent,
		message,
		decision.DomainAgent,
		intentResult.Confidence,
		n.logger,
	)

	// 将工具信息添加到决策中
	decision.Tools = selectedTools
	decision.ToolStrategy = string(toolStrategy)

	n.logger.Infow("Supervisor协调完成",
		logx.Field("intent", intentResult.Intent),
		logx.Field("strategy", cognitiveLoadAdvice.Strategy),
		logx.Field("domainAgent", decision.DomainAgent),
		logx.Field("action", decision.Action),
		logx.Field("tools", selectedTools),
		logx.Field("toolStrategy", toolStrategy),
	)

	return decision, nil
}

