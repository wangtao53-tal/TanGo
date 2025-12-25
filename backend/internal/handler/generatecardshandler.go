package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

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

		// 同步返回模式：等待模型返回，不设置超时
		// 超时控制由HTTP请求层面的Timeout配置控制（在explore.yaml中配置为180秒）
		l := logic.NewGenerateCardsLogic(r.Context(), svcCtx)
		resp, err := l.GenerateCards(&req)
		if err != nil {
			httpx.Error(w, err)
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

	// 调用逻辑层生成卡片（流式）
	// 等待模型返回，不设置超时
	// 超时控制由HTTP请求层面的Timeout配置控制（在explore.yaml中配置为180秒）
	logic := logic.NewGenerateCardsLogic(r.Context(), svcCtx)
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
