package nodes

import (
	"context"
	"testing"

	"github.com/tango/explore/internal/config"
	"github.com/zeromicro/go-zero/core/logx"
)

func TestCognitiveLoadNode_AssessCognitiveLoad(t *testing.T) {
	ctx := context.Background()
	logger := logx.WithContext(ctx)
	cfg := config.AIConfig{
		EinoBaseURL: "",
		AppID:       "",
		AppKey:      "",
	}

	node, err := NewCognitiveLoadNode(ctx, cfg, logger)
	if err != nil {
		t.Fatalf("Failed to create CognitiveLoadNode: %v", err)
	}

	testCases := []struct {
		name                string
		userAge             int
		conversationRounds  int
		recentOutputLength   int
		expectedStrategy     string
		expectedMaxSentences int
	}{
		{"3-6岁简短讲解", 5, 1, 100, "简短讲解", 3},
		{"7-12岁类比讲解", 10, 1, 100, "类比讲解", 5},
		{"13-18岁深入讲解", 15, 1, 100, "深入讲解", 7},
		{"连续追问>5轮", 10, 6, 100, "反问引导", 2},
		{"输出>500字", 10, 1, 600, "暂停探索", 1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			advice, err := node.AssessCognitiveLoad(ctx, tc.userAge, tc.conversationRounds, tc.recentOutputLength)
			if err != nil {
				t.Errorf("AssessCognitiveLoad failed: %v", err)
				return
			}

			if advice.Strategy != tc.expectedStrategy {
				t.Errorf("Expected strategy %s, got %s", tc.expectedStrategy, advice.Strategy)
			}

			if advice.MaxSentences != tc.expectedMaxSentences {
				t.Errorf("Expected maxSentences %d, got %d", tc.expectedMaxSentences, advice.MaxSentences)
			}

			if advice.Reason == "" {
				t.Error("Reason should not be empty")
			}

			// 验证策略类型是否在有效范围内
			validStrategies := []string{"简短讲解", "类比讲解", "深入讲解", "反问引导", "暂停探索"}
			isValid := false
			for _, validStrategy := range validStrategies {
				if advice.Strategy == validStrategy {
					isValid = true
					break
				}
			}
			if !isValid {
				t.Errorf("Invalid strategy type: %s", advice.Strategy)
			}
		})
	}
}

