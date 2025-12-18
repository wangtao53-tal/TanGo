package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/tango/explore/internal/logic"
	"github.com/tango/explore/internal/svc"
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

		// 创建流式逻辑
		streamLogic := logic.NewStreamLogic(r.Context(), svcCtx)

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
