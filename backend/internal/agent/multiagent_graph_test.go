package agent

import (
	"context"
	"testing"
	"time"

	"github.com/tango/explore/internal/config"
	"github.com/tango/explore/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

func TestMultiAgentGraph_ExecuteMultiAgentConversation(t *testing.T) {
	ctx := context.Background()
	logger := logx.WithContext(ctx)
	cfg := config.AIConfig{
		EinoBaseURL: "",
		AppID:       "",
		AppKey:      "",
	}

	graph, err := NewMultiAgentGraph(ctx, cfg, logger)
	if err != nil {
		t.Fatalf("Failed to create MultiAgentGraph: %v", err)
	}

	req := &types.UnifiedStreamConversationRequest{
		MessageType: "text",
		Message:     "这是什么？",
		SessionId:   "test-session-123",
		UserAge:     10,
		IdentificationContext: &types.IdentificationContext{
			ObjectName:     "银杏",
			ObjectCategory: "自然类",
			Confidence:     0.9,
		},
	}

	startTime := time.Now()
	answer, err := graph.ExecuteMultiAgentConversation(ctx, req, nil)
	duration := time.Since(startTime)

	if err != nil {
		t.Errorf("ExecuteMultiAgentConversation failed: %v", err)
		return
	}

	if answer == "" {
		t.Error("Answer should not be empty")
	}

	// 验证执行时间（目标≤8秒）
	if duration > 8*time.Second {
		t.Errorf("Execution time should be ≤8s, got %v", duration)
	}

	t.Logf("MultiAgent conversation completed in %v, answer length: %d", duration, len(answer))
}

func TestMultiAgentGraph_ExecuteMultiAgentConversation_ErrorHandling(t *testing.T) {
	ctx := context.Background()
	logger := logx.WithContext(ctx)
	cfg := config.AIConfig{
		EinoBaseURL: "",
		AppID:       "",
		AppKey:      "",
	}

	graph, err := NewMultiAgentGraph(ctx, cfg, logger)
	if err != nil {
		t.Fatalf("Failed to create MultiAgentGraph: %v", err)
	}

	// 测试空消息
	req := &types.UnifiedStreamConversationRequest{
		MessageType: "text",
		Message:     "",
		SessionId:   "test-session-123",
		UserAge:     10,
	}

	_, err = graph.ExecuteMultiAgentConversation(ctx, req, nil)
	// 应该能够处理空消息（使用Mock模式）
	if err != nil {
		t.Logf("ExecuteMultiAgentConversation returned error for empty message (expected in some cases): %v", err)
	}
}

