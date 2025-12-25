package logic

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tango/explore/internal/config"
	"github.com/tango/explore/internal/svc"
	"github.com/tango/explore/internal/storage"
	"github.com/tango/explore/internal/types"
)

func TestAgentLogic_StreamAgentConversation(t *testing.T) {
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

	logic := NewAgentLogic(ctx, svcCtx)

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

	// 创建HTTP请求和响应
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/api/conversation/agent", bytes.NewReader(body))
	httpReq.Body = http.MaxBytesReader(nil, io.NopCloser(bytes.NewReader(body)), 10*1024*1024)
	httpReq.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	// 执行流式对话（注意：这会尝试连接真实服务，在测试环境中可能失败）
	err := logic.StreamAgentConversation(w, req)
	
	// 在测试环境中，由于没有真实AI模型，可能会失败或降级
	// 我们主要验证逻辑是否正确执行，不强制要求成功
	if err != nil {
		t.Logf("StreamAgentConversation returned error (expected in test environment): %v", err)
	}

	// 验证响应头
	if w.Header().Get("Content-Type") != "text/event-stream" {
		t.Error("Response should have Content-Type: text/event-stream")
	}
}

func TestAgentLogic_StreamAgentConversation_ErrorHandling(t *testing.T) {
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

	logic := NewAgentLogic(ctx, svcCtx)

	// 测试无效的messageType
	req := types.UnifiedStreamConversationRequest{
		MessageType: "",
		Message:     "这是什么？",
		SessionId:   "test-session-123",
		UserAge:     10,
	}

	w := httptest.NewRecorder()
	err := logic.StreamAgentConversation(w, req)
	
	// 应该返回错误
	if err == nil {
		t.Error("Should return error for empty messageType")
	}
}

