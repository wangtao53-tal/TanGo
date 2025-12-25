package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/tango/explore/internal/logic"
	"github.com/tango/explore/internal/svc"
	"github.com/tango/explore/internal/types"
)

func StreamHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 设置SSE响应头
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Cache-Control")

		// 获取sessionId
		sessionId := r.URL.Query().Get("sessionId")
		if sessionId == "" {
			// 发送错误事件
			fmt.Fprintf(w, "event: error\ndata: {\"message\":\"sessionId is required\"}\n\n")
			w.(http.Flusher).Flush()
			return
		}

		// 获取消息（从查询参数或请求体）
		message := r.URL.Query().Get("message")
		if message == "" {
			message = "开始对话"
		}

		// 创建流式逻辑（TODO: 待实现真实流式逻辑时使用）
		_ = logic.NewStreamLogic(r.Context(), svcCtx)

		// 发送初始连接事件
		fmt.Fprintf(w, "event: connected\ndata: {\"sessionId\":\"%s\"}\n\n", sessionId)
		w.(http.Flusher).Flush()

		// 模拟流式输出（TODO: 接入真实AI模型）
		text := "这是一个Mock流式响应。待接入真实AI模型后，将实现真实的流式文本生成。"
		words := []rune(text)

		for i, word := range words {
			eventData := fmt.Sprintf(`{"type":"text","content":"%s","index":%d}`, string(word), i)
			fmt.Fprintf(w, "event: message\ndata: %s\n\n", eventData)
			w.(http.Flusher).Flush()
			time.Sleep(50 * time.Millisecond) // 模拟延迟
		}

		// 发送完成事件
		fmt.Fprintf(w, "event: done\ndata: {\"type\":\"done\"}\n\n")
		w.(http.Flusher).Flush()
	}
}

// StreamConversationHandler 流式对话Handler（统一接口）
// 支持POST请求，从请求体解析参数，支持文本、语音、图片三种输入方式
func StreamConversationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 设置SSE响应头
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

		// 处理OPTIONS预检请求
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// 从请求体解析JSON参数（POST请求）
		var req types.UnifiedStreamConversationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			// 发送错误事件
			errorEvent := types.StreamEvent{
				Type:    "error",
				Content: map[string]interface{}{"message": "请求参数解析失败: " + err.Error()},
			}
			errorJSON, _ := json.Marshal(errorEvent)
			fmt.Fprintf(w, "event: error\ndata: %s\n\n", string(errorJSON))
			w.(http.Flusher).Flush()
			return
		}

		// 设置默认值
		if req.MaxContextRounds <= 0 {
			req.MaxContextRounds = 20
		}

		// 调用统一流式逻辑
		streamLogic := logic.NewStreamLogic(r.Context(), svcCtx)
		if err := streamLogic.StreamConversationUnified(w, req); err != nil {
			// 错误已在StreamConversationUnified中处理，这里不需要再次发送错误事件
			// 但如果需要，可以在这里添加额外的错误处理
		}
	}
}
