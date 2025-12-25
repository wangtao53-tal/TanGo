package nodes

import (
	"context"
	"testing"

	"github.com/tango/explore/internal/config"
	"github.com/tango/explore/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

func TestLearningPlannerNode_PlanLearning(t *testing.T) {
	ctx := context.Background()
	logger := logx.WithContext(ctx)
	cfg := config.AIConfig{
		EinoBaseURL: "",
		AppID:       "",
		AppKey:      "",
	}

	node, err := NewLearningPlannerNode(ctx, cfg, logger)
	if err != nil {
		t.Fatalf("Failed to create LearningPlannerNode: %v", err)
	}

	testCases := []struct {
		name                string
		intentResult        *types.FollowUpIntentResult
		cognitiveLoadAdvice *types.CognitiveLoadAdvice
		objectName          string
		objectCategory      string
		userAge             int
	}{
		{
			"探因型选择Science",
			&types.FollowUpIntentResult{Intent: "探因型", Confidence: 0.9},
			&types.CognitiveLoadAdvice{Strategy: "类比讲解", MaxSentences: 5},
			"银杏",
			"自然类",
			10,
		},
		{
			"表达型选择Language",
			&types.FollowUpIntentResult{Intent: "表达型", Confidence: 0.9},
			&types.CognitiveLoadAdvice{Strategy: "类比讲解", MaxSentences: 5},
			"银杏",
			"自然类",
			10,
		},
		{
			"游戏型选择Humanities",
			&types.FollowUpIntentResult{Intent: "游戏型", Confidence: 0.9},
			&types.CognitiveLoadAdvice{Strategy: "类比讲解", MaxSentences: 5},
			"银杏",
			"自然类",
			10,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			decision, err := node.PlanLearning(ctx, tc.intentResult, tc.cognitiveLoadAdvice, tc.objectName, tc.objectCategory, tc.userAge)
			if err != nil {
				t.Errorf("PlanLearning failed: %v", err)
				return
			}

			if decision.DomainAgent == "" {
				t.Error("DomainAgent should not be empty")
			}

			if decision.Action == "" {
				t.Error("Action should not be empty")
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
				t.Errorf("Invalid domain agent type: %s", decision.DomainAgent)
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
				t.Errorf("Invalid action type: %s", decision.Action)
			}
		})
	}
}

