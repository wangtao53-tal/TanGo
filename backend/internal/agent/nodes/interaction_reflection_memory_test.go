package nodes

import (
	"context"
	"testing"

	"github.com/tango/explore/internal/config"
	"github.com/tango/explore/internal/storage"
	"github.com/tango/explore/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

func TestInteractionAgentNode_OptimizeInteraction(t *testing.T) {
	ctx := context.Background()
	logger := logx.WithContext(ctx)
	cfg := config.AIConfig{
		EinoBaseURL: "",
		AppID:       "",
		AppKey:      "",
	}

	node, err := NewInteractionAgentNode(ctx, cfg, logger)
	if err != nil {
		t.Fatalf("Failed to create InteractionAgentNode: %v", err)
	}

	originalContent := "这是关于银杏的科学知识。"
	result, err := node.OptimizeInteraction(ctx, originalContent)
	if err != nil {
		t.Errorf("OptimizeInteraction failed: %v", err)
		return
	}

	if result.OptimizedContent == "" {
		t.Error("OptimizedContent should not be empty")
	}
}

func TestReflectionAgentNode_Reflect(t *testing.T) {
	ctx := context.Background()
	logger := logx.WithContext(ctx)
	cfg := config.AIConfig{
		EinoBaseURL: "",
		AppID:       "",
		AppKey:      "",
	}

	node, err := NewReflectionAgentNode(ctx, cfg, logger)
	if err != nil {
		t.Fatalf("Failed to create ReflectionAgentNode: %v", err)
	}

	result, err := node.Reflect(ctx, "这是关于银杏的科学知识。", nil)
	if err != nil {
		t.Errorf("Reflect failed: %v", err)
		return
	}

	// 验证结果不为nil
	if result == nil {
		t.Error("ReflectionResult should not be nil")
	}
}

func TestMemoryAgentNode_RecordMemory(t *testing.T) {
	ctx := context.Background()
	logger := logx.WithContext(ctx)
	cfg := config.AIConfig{
		EinoBaseURL: "",
		AppID:       "",
		AppKey:      "",
	}

	memoryStorage := storage.NewMemoryAgentStorage()
	node, err := NewMemoryAgentNode(ctx, cfg, logger, memoryStorage)
	if err != nil {
		t.Fatalf("Failed to create MemoryAgentNode: %v", err)
	}

	sessionId := "test-session-123"
	reflectionResult := &types.ReflectionResult{
		Interest:  true,
		Confusion: false,
		Relax:     false,
	}

	err = node.RecordMemory(ctx, sessionId, reflectionResult, "这是关于银杏的科学知识。", "银杏")
	if err != nil {
		t.Errorf("RecordMemory failed: %v", err)
		return
	}

	// 验证记忆记录
	record, exists := node.GetMemory(ctx, sessionId)
	if !exists {
		t.Error("Memory record should exist after recording")
		return
	}

	if record.SessionId != sessionId {
		t.Errorf("Expected sessionId %s, got %s", sessionId, record.SessionId)
	}
}

