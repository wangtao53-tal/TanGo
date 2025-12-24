package nodes

import (
	"context"
	"testing"

	"github.com/tango/explore/internal/config"
	"github.com/tango/explore/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

func TestSupervisorNode_Coordinate(t *testing.T) {
	ctx := context.Background()
	logger := logx.WithContext(ctx)
	cfg := config.AIConfig{
		EinoBaseURL: "",
		AppID:       "",
		AppKey:      "",
	}

	// 创建子Agent节点
	intentAgent, _ := NewIntentAgentNode(ctx, cfg, logger)
	cognitiveLoadAgent, _ := NewCognitiveLoadNode(ctx, cfg, logger)
	learningPlannerAgent, _ := NewLearningPlannerNode(ctx, cfg, logger)

	supervisor, err := NewSupervisorNode(ctx, cfg, logger, intentAgent, cognitiveLoadAgent, learningPlannerAgent)
	if err != nil {
		t.Fatalf("Failed to create SupervisorNode: %v", err)
	}

	state := &types.SupervisorState{
		ObjectName:         "银杏",
		ObjectCategory:     "自然类",
		UserAge:            10,
		ConversationRounds: 1,
		RecentOutputLength: 100,
		AgentResults:       make(map[string]interface{}),
		SessionId:          "test-session-123",
	}

	decision, err := supervisor.Coordinate(ctx, state, "这是什么？", nil)
	if err != nil {
		t.Errorf("Coordinate failed: %v", err)
		return
	}

	if decision == nil {
		t.Error("Decision should not be nil")
		return
	}

	if decision.DomainAgent == "" {
		t.Error("DomainAgent should not be empty")
	}

	if decision.Action == "" {
		t.Error("Action should not be empty")
	}

	// 验证AgentResults已填充
	if state.AgentResults["intent"] == nil {
		t.Error("Intent result should be stored in AgentResults")
	}

	if state.AgentResults["cognitiveLoad"] == nil {
		t.Error("Cognitive load advice should be stored in AgentResults")
	}

	if state.AgentResults["learningPlan"] == nil {
		t.Error("Learning plan decision should be stored in AgentResults")
	}
}

