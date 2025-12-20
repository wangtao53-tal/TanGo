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

func UploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UploadRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewUploadLogic(r.Context(), svcCtx)
		resp, err := l.Upload(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

// UploadStreamHandler 图片上传并流式返回Handler
func UploadStreamHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
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
		var req types.UploadRequest
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

		// 参数验证
		if req.ImageData == "" {
			errorEvent := types.StreamEvent{
				Type:    "error",
				Content: map[string]interface{}{"message": "图片数据不能为空"},
			}
			errorJSON, _ := json.Marshal(errorEvent)
			fmt.Fprintf(w, "event: error\ndata: %s\n\n", string(errorJSON))
			w.(http.Flusher).Flush()
			return
		}

		// 图片上传
		uploadLogic := logic.NewUploadLogic(r.Context(), svcCtx)
		uploadResp, err := uploadLogic.Upload(&req)
		if err != nil {
			// 发送错误事件
			errorEvent := types.StreamEvent{
				Type:    "error",
				Content: map[string]interface{}{"message": "图片上传失败: " + err.Error()},
			}
			errorJSON, _ := json.Marshal(errorEvent)
			fmt.Fprintf(w, "event: error\ndata: %s\n\n", string(errorJSON))
			w.(http.Flusher).Flush()
			return
		}

		// 发送图片上传完成事件
		uploadedEvent := types.StreamEvent{
			Type:    "image_uploaded",
			Content: map[string]interface{}{"url": uploadResp.Url, "filename": uploadResp.Filename},
		}
		uploadedJSON, _ := json.Marshal(uploadedEvent)
		fmt.Fprintf(w, "event: image_uploaded\ndata: %s\n\n", string(uploadedJSON))
		w.(http.Flusher).Flush()

		// 注意：此handler已废弃，请使用统一接口 /api/conversation/stream
		// 调用统一流式接口（传入图片URL）
		streamLogic := logic.NewStreamLogic(r.Context(), svcCtx)
		streamReq := types.UnifiedStreamConversationRequest{
			MessageType: "image",
			Image:       uploadResp.Url,
			// 注意：UploadRequest中没有这些字段，需要从查询参数或其他地方获取
			// SessionId:             "",
			// Message:               "",
			// UserAge:               0,
			// MaxContextRounds:      20,
			// IdentificationContext: nil,
		}
		if err := streamLogic.StreamConversationUnified(w, streamReq); err != nil {
			// 错误已在StreamConversationUnified中处理
			return
		}
	}
}
