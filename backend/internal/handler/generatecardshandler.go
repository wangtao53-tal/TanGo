package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/tango/explore/internal/logic"
	"github.com/tango/explore/internal/svc"
	"github.com/tango/explore/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GenerateCardsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GenerateCardsRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		// 检查是否使用流式返回模式
		useStream := r.URL.Query().Get("stream") == "true"

		if useStream {
			// 流式返回模式：每生成完一张卡片立即返回
			generateCardsStream(w, r, &req, svcCtx)
			return
		}

		// 同步返回模式（保持兼容）：优化超时时间到6秒（目标5秒内）
		ctx, cancel := context.WithTimeout(r.Context(), 6*time.Second)
		defer cancel()

		l := logic.NewGenerateCardsLogic(ctx, svcCtx)
		resp, err := l.GenerateCards(&req)
		if err != nil {
			// 检查是否是超时错误
			if ctx.Err() == context.DeadlineExceeded {
				httpx.ErrorCtx(ctx, w, err)
			} else {
				httpx.Error(w, err)
			}
		} else {
			httpx.OkJson(w, resp)
		}
	}
}

// generateCardsStream 流式返回知识卡片
func generateCardsStream(w http.ResponseWriter, r *http.Request, req *types.GenerateCardsRequest, svcCtx *svc.ServiceContext) {
	// 设置SSE响应头
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	ctx, cancel := context.WithTimeout(r.Context(), 6*time.Second)
	defer cancel()

	// 调用逻辑层生成卡片（流式）
	logic := logic.NewGenerateCardsLogic(ctx, svcCtx)
	if err := logic.GenerateCardsStream(w, req); err != nil {
		// 发送错误事件
		errorEvent := map[string]interface{}{
			"type":    "error",
			"content": map[string]interface{}{"message": err.Error()},
		}
		errorJSON, _ := json.Marshal(errorEvent)
		fmt.Fprintf(w, "event: error\ndata: %s\n\n", string(errorJSON))
		w.(http.Flusher).Flush()
	}
}
