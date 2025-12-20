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

func VoiceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.VoiceRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewVoiceLogic(r.Context(), svcCtx)
		resp, err := l.RecognizeVoice(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}

// VoiceStreamHandler 语音识别并流式返回Handler
func VoiceStreamHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
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
		var req types.VoiceRequest
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
		if req.Audio == "" {
			errorEvent := types.StreamEvent{
				Type:    "error",
				Content: map[string]interface{}{"message": "音频数据不能为空"},
			}
			errorJSON, _ := json.Marshal(errorEvent)
			fmt.Fprintf(w, "event: error\ndata: %s\n\n", string(errorJSON))
			w.(http.Flusher).Flush()
			return
		}

		// 语音识别
		voiceLogic := logic.NewVoiceLogic(r.Context(), svcCtx)
		voiceResp, err := voiceLogic.RecognizeVoice(&req)
		if err != nil {
			// 发送错误事件
			errorEvent := types.StreamEvent{
				Type:    "error",
				Content: map[string]interface{}{"message": "语音识别失败: " + err.Error()},
			}
			errorJSON, _ := json.Marshal(errorEvent)
			fmt.Fprintf(w, "event: error\ndata: %s\n\n", string(errorJSON))
			w.(http.Flusher).Flush()
			return
		}

		// 发送语音识别完成事件
		recognizedEvent := types.StreamEvent{
			Type:    "voice_recognized",
			Content: map[string]interface{}{"text": voiceResp.Text},
		}
		recognizedJSON, _ := json.Marshal(recognizedEvent)
		fmt.Fprintf(w, "event: voice_recognized\ndata: %s\n\n", string(recognizedJSON))
		w.(http.Flusher).Flush()

		// 注意：此handler已废弃，请使用统一接口 /api/conversation/stream
		// 调用统一流式接口
		streamLogic := logic.NewStreamLogic(r.Context(), svcCtx)
		streamReq := types.UnifiedStreamConversationRequest{
			MessageType: "voice",
			Audio:       req.Audio, // 使用原始音频数据
			SessionId:   voiceResp.SessionId,
			// 注意：VoiceRequest中没有这些字段，需要从查询参数或其他地方获取
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
