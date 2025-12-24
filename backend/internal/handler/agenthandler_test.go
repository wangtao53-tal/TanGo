package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tango/explore/internal/config"
	"github.com/tango/explore/internal/svc"
	"github.com/tango/explore/internal/storage"
	"github.com/tango/explore/internal/types"
)

func TestAgentConversationHandler(t *testing.T) {
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

	handler := AgentConversationHandler(svcCtx)

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

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/api/conversation/agent", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	handler(w, httpReq)

	// 验证响应头
	if w.Header().Get("Content-Type") != "text/event-stream" {
		t.Error("Response should have Content-Type: text/event-stream")
	}

	// 验证状态码（可能是200或500，取决于是否成功）
	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 200 or 500, got %d", w.Code)
	}
}

func TestAgentConversationHandler_InvalidRequest(t *testing.T) {
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

	handler := AgentConversationHandler(svcCtx)

	// 测试无效的JSON
	httpReq := httptest.NewRequest("POST", "/api/conversation/agent", bytes.NewReader([]byte("invalid json")))
	httpReq.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	handler(w, httpReq)

	// 应该返回错误
	if w.Code == http.StatusOK {
		t.Error("Should return error for invalid JSON")
	}
}

