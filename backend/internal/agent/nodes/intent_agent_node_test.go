package nodes

import (
	"context"
	"testing"

	"github.com/tango/explore/internal/config"
	"github.com/zeromicro/go-zero/core/logx"
)

func TestIntentAgentNode_RecognizeIntent(t *testing.T) {
	ctx := context.Background()
	logger := logx.WithContext(ctx)
	cfg := config.AIConfig{
		EinoBaseURL: "",
		AppID:       "",
		AppKey:      "",
	}

	node, err := NewIntentAgentNode(ctx, cfg, logger)
	if err != nil {
		t.Fatalf("Failed to create IntentAgentNode: %v", err)
	}

	testCases := []struct {
		name     string
		message  string
		expected string
	}{
		{"认知型意图", "这是什么？", "认知型"},
		{"探因型意图", "为什么会这样？", "探因型"},
		{"表达型意图", "用英语怎么说？", "表达型"},
		{"游戏型意图", "好玩吗？", "游戏型"},
		{"情绪型意图", "我不懂", "情绪型"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := node.RecognizeIntent(ctx, tc.message, nil)
			if err != nil {
				t.Errorf("RecognizeIntent failed: %v", err)
				return
			}

			if result.Intent == "" {
				t.Error("Intent should not be empty")
			}

			if result.Confidence < 0 || result.Confidence > 1 {
				t.Errorf("Confidence should be between 0 and 1, got %f", result.Confidence)
			}

			// 验证意图类型是否在有效范围内
			validIntents := []string{"认知型", "探因型", "表达型", "游戏型", "情绪型"}
			isValid := false
			for _, validIntent := range validIntents {
				if result.Intent == validIntent {
					isValid = true
					break
				}
			}
			if !isValid {
				t.Errorf("Invalid intent type: %s", result.Intent)
			}
		})
	}
}

func TestIntentAgentNode_RecognizeIntent_EdgeCases(t *testing.T) {
	ctx := context.Background()
	logger := logx.WithContext(ctx)
	cfg := config.AIConfig{
		EinoBaseURL: "",
		AppID:       "",
		AppKey:      "",
	}

	node, err := NewIntentAgentNode(ctx, cfg, logger)
	if err != nil {
		t.Fatalf("Failed to create IntentAgentNode: %v", err)
	}

	// 测试模糊意图
	result, err := node.RecognizeIntent(ctx, "嗯...", nil)
	if err != nil {
		t.Errorf("RecognizeIntent failed for ambiguous intent: %v", err)
	}
	if result.Intent == "" {
		t.Error("Intent should not be empty even for ambiguous input")
	}

	// 测试空消息
	result2, err2 := node.RecognizeIntent(ctx, "", nil)
	if err2 != nil {
		t.Errorf("RecognizeIntent should handle empty message gracefully: %v", err2)
	}
	if result2 != nil && result2.Intent == "" {
		t.Error("Intent should have a default value for empty message")
	}
}

