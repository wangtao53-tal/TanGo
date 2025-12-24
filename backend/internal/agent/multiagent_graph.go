package agent

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/schema"
	"github.com/tango/explore/internal/agent/nodes"
	"github.com/tango/explore/internal/config"
	"github.com/tango/explore/internal/storage"
	"github.com/tango/explore/internal/tools"
	"github.com/tango/explore/internal/tools/base"
	"github.com/tango/explore/internal/tools/mcp"
	"github.com/tango/explore/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

// MultiAgentGraph 多Agent Graph结构
type MultiAgentGraph struct {
	ctx    context.Context
	config config.AIConfig
	logger logx.Logger

	// Agent节点实例
	supervisorNode      *nodes.SupervisorNode
	intentAgentNode     *nodes.IntentAgentNode
	cognitiveLoadNode   *nodes.CognitiveLoadNode
	learningPlannerNode *nodes.LearningPlannerNode
	scienceAgentNode    *nodes.ScienceAgentNode
	languageAgentNode   *nodes.LanguageAgentNode
	humanitiesAgentNode *nodes.HumanitiesAgentNode
	interactionAgentNode *nodes.InteractionAgentNode
	reflectionAgentNode *nodes.ReflectionAgentNode
	memoryAgentNode    *nodes.MemoryAgentNode

	// 存储
	memoryStorage *storage.MemoryAgentStorage
}

// NewMultiAgentGraph 创建MultiAgentGraph实例
func NewMultiAgentGraph(ctx context.Context, cfg config.AIConfig, logger logx.Logger) (*MultiAgentGraph, error) {
	graph := &MultiAgentGraph{
		ctx:    ctx,
		config: cfg,
		logger: logger,
	}

	// 初始化Memory存储
	graph.memoryStorage = storage.NewMemoryAgentStorage()

	// 初始化各个Agent节点
	var err error

	// 1. Intent Agent
	graph.intentAgentNode, err = nodes.NewIntentAgentNode(ctx, cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("初始化Intent Agent失败: %w", err)
	}

	// 2. Cognitive Load Agent
	graph.cognitiveLoadNode, err = nodes.NewCognitiveLoadNode(ctx, cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("初始化Cognitive Load Agent失败: %w", err)
	}

	// 3. Learning Planner Agent
	graph.learningPlannerNode, err = nodes.NewLearningPlannerNode(ctx, cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("初始化Learning Planner Agent失败: %w", err)
	}

	// 4. Supervisor Node（依赖Intent、Cognitive Load、Learning Planner）
	graph.supervisorNode, err = nodes.NewSupervisorNode(ctx, cfg, logger, graph.intentAgentNode, graph.cognitiveLoadNode, graph.learningPlannerNode)
	if err != nil {
		return nil, fmt.Errorf("初始化Supervisor节点失败: %w", err)
	}

	// 初始化工具注册表
	toolRegistry := tools.GetDefaultRegistry(logger)

	// 注册基础工具
	baseTools := []tools.Tool{
		base.NewSimpleFactLookupTool(logger),
		base.NewSimpleDictionaryTool(logger),
		base.NewPronunciationHintTool(logger),
		base.NewGetCurrentTimeTool(logger),
		base.NewImageGenerateSimpleTool(logger),
	}
	tools.InitDefaultTools(logger, baseTools)

	// 注册MCP工具（如果启用）
	mcpConfig := config.GetMCPConfig(logger)
	if mcpConfig != nil && mcpConfig.Enabled {
		if err := mcp.DiscoverAndRegisterMCPTools(mcpConfig, toolRegistry, logger); err != nil {
			logger.Errorw("注册MCP工具失败", logx.Field("error", err))
		}
	}

	// 5. Domain Agents（传递工具注册表）
	graph.scienceAgentNode, err = nodes.NewScienceAgentNode(ctx, cfg, logger, toolRegistry)
	if err != nil {
		return nil, fmt.Errorf("初始化Science Agent失败: %w", err)
	}

	graph.languageAgentNode, err = nodes.NewLanguageAgentNode(ctx, cfg, logger, toolRegistry)
	if err != nil {
		return nil, fmt.Errorf("初始化Language Agent失败: %w", err)
	}

	graph.humanitiesAgentNode, err = nodes.NewHumanitiesAgentNode(ctx, cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("初始化Humanities Agent失败: %w", err)
	}

	// 6. Interaction Agent
	graph.interactionAgentNode, err = nodes.NewInteractionAgentNode(ctx, cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("初始化Interaction Agent失败: %w", err)
	}

	// 7. Reflection Agent
	graph.reflectionAgentNode, err = nodes.NewReflectionAgentNode(ctx, cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("初始化Reflection Agent失败: %w", err)
	}

	// 8. Memory Agent
	graph.memoryAgentNode, err = nodes.NewMemoryAgentNode(ctx, cfg, logger, graph.memoryStorage)
	if err != nil {
		return nil, fmt.Errorf("初始化Memory Agent失败: %w", err)
	}

	logger.Info("✅ MultiAgentGraph初始化完成")
	return graph, nil
}

// ExecuteMultiAgentConversation 执行多Agent对话流程
func (g *MultiAgentGraph) ExecuteMultiAgentConversation(
	ctx context.Context,
	req *types.UnifiedStreamConversationRequest,
	chatHistory []*schema.Message,
) (string, error) {
	g.logger.Infow("开始执行多Agent对话流程",
		logx.Field("sessionId", req.SessionId),
		logx.Field("messageType", req.MessageType),
		logx.Field("userAge", req.UserAge),
	)

	// 1. 构建SupervisorState
	state := &types.SupervisorState{
		ObjectName:         "",
		ObjectCategory:     "",
		Cards:              []types.CardContent{},
		UserAge:            req.UserAge,
		ConversationRounds: len(chatHistory) / 2, // 简单估算对话轮数
		RecentOutputLength: 0,                      // TODO: 从chatHistory计算
		AgentResults:       make(map[string]interface{}),
		SessionId:          req.SessionId,
	}
	
	// 确保AgentResults已初始化
	if state.AgentResults == nil {
		state.AgentResults = make(map[string]interface{})
	}

	// 从IdentificationContext获取对象信息
	if req.IdentificationContext != nil {
		state.ObjectName = req.IdentificationContext.ObjectName
		state.ObjectCategory = req.IdentificationContext.ObjectCategory
	}

	// 获取用户消息
	message := req.Message
	if req.MessageType == "voice" && req.Audio != "" {
		// TODO: 语音识别
		message = "语音消息（待识别）"
	}

	// 2. Supervisor协调：调用Intent、Cognitive Load、Learning Planner
	decision, err := g.supervisorNode.Coordinate(ctx, state, message, chatHistory)
	if err != nil {
		return "", fmt.Errorf("Supervisor协调失败: %w", err)
	}

	// 3. 根据决策选择Domain Agent
	var domainResponse *types.DomainAgentResponse
	maxSentences := 5 // 默认值
	if cognitiveLoadAdvice, ok := state.AgentResults["cognitiveLoad"].(*types.CognitiveLoadAdvice); ok && cognitiveLoadAdvice != nil {
		maxSentences = cognitiveLoadAdvice.MaxSentences
	}
	
	switch decision.DomainAgent {
	case "Science":
		domainResponse, err = g.scienceAgentNode.GenerateScienceAnswer(ctx, message, state.ObjectName, state.ObjectCategory, state.UserAge, chatHistory, maxSentences, decision.Tools)
	case "Language":
		domainResponse, err = g.languageAgentNode.GenerateLanguageAnswer(ctx, message, state.ObjectName, state.ObjectCategory, state.UserAge, chatHistory, decision.Tools)
	case "Humanities":
		domainResponse, err = g.humanitiesAgentNode.GenerateHumanitiesAnswer(ctx, message, state.ObjectName, state.ObjectCategory, state.UserAge, chatHistory)
	default:
		return "", fmt.Errorf("未知的领域Agent: %s", decision.DomainAgent)
	}
	if err != nil {
		return "", fmt.Errorf("Domain Agent生成回答失败: %w", err)
	}

	// 4. Interaction Agent优化交互
	interactionResult, err := g.interactionAgentNode.OptimizeInteraction(ctx, domainResponse.Content)
	if err != nil {
		g.logger.Errorw("Interaction Agent优化失败，使用原始回答", logx.Field("error", err))
		interactionResult = &types.InteractionOptimization{
			OptimizedContent: domainResponse.Content,
			EndingAction:     "",
		}
	}

	// 5. Reflection Agent反思判断
	reflectionResult, err := g.reflectionAgentNode.Reflect(ctx, interactionResult.OptimizedContent, chatHistory)
	if err != nil {
		g.logger.Errorw("Reflection Agent反思失败", logx.Field("error", err))
		reflectionResult = &types.ReflectionResult{
			Interest:  true,
			Confusion: false,
			Relax:     false,
		}
	}

	// 6. Memory Agent记录学习状态
	err = g.memoryAgentNode.RecordMemory(ctx, req.SessionId, reflectionResult, interactionResult.OptimizedContent, state.ObjectName)
	if err != nil {
		g.logger.Errorw("Memory Agent记录失败", logx.Field("error", err))
	}

	g.logger.Infow("多Agent对话流程完成",
		logx.Field("domainAgent", decision.DomainAgent),
		logx.Field("contentLength", len(interactionResult.OptimizedContent)),
	)

	return interactionResult.OptimizedContent, nil
}

