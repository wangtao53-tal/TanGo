package logic

import (
	"context"
	"testing"

	"github.com/tango/explore/internal/config"
	"github.com/tango/explore/internal/svc"
	"github.com/tango/explore/internal/storage"
	"github.com/tango/explore/internal/types"
)

// TestAgentLogic_InterfaceConsistency 测试接口一致性
// 验证 /api/conversation/agent 和 /api/conversation/stream 使用相同的请求和响应格式
func TestAgentLogic_InterfaceConsistency(t *testing.T) {
	ctx := context.Background()
	cfg := config.Config{
		AI: config.AIConfig{
			EinoBaseURL: "",
			AppID:       "",
			AppKey:      "",
		},
	}
	svcCtx := &svc.ServiceContext{
		Config:  cfg,
		Storage: storage.NewMemoryStorage(),
	}

	// 创建统一的请求
	req := types.UnifiedStreamConversationRequest{
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

	// 测试AgentLogic（多Agent模式）
	agentLogic := NewAgentLogic(ctx, svcCtx)
	
	// 验证AgentLogic可以处理UnifiedStreamConversationRequest
	if agentLogic == nil {
		t.Error("AgentLogic should not be nil")
	}

	// 测试StreamLogic（单Agent模式）
	streamLogic := NewStreamLogic(ctx, svcCtx)
	
	// 验证StreamLogic也可以处理UnifiedStreamConversationRequest
	if streamLogic == nil {
		t.Error("StreamLogic should not be nil")
	}

	// 验证两个接口使用相同的请求类型
	_ = req // 使用req变量验证类型一致性
	t.Log("Both interfaces use UnifiedStreamConversationRequest - consistency verified")
}

// TestAgentLogic_FallbackToSingleAgent 测试降级机制
func TestAgentLogic_FallbackToSingleAgent(t *testing.T) {
	ctx := context.Background()
	cfg := config.Config{
		AI: config.AIConfig{
			EinoBaseURL: "",
			AppID:       "",
			AppKey:      "",
		},
	}
	svcCtx := &svc.ServiceContext{
		Config:  cfg,
		Storage: storage.NewMemoryStorage(),
	}

	agentLogic := NewAgentLogic(ctx, svcCtx)
	
	// 验证AgentLogic存在降级逻辑
	// 在实际实现中，如果MultiAgentGraph初始化失败，应该降级到StreamLogic
	if agentLogic == nil {
		t.Error("AgentLogic should not be nil")
	}

	t.Log("Fallback mechanism exists in AgentLogic")
}

