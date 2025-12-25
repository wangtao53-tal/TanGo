package handler

import (
	"encoding/json"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/tango/explore/internal/logic"
	"github.com/tango/explore/internal/svc"
	"github.com/tango/explore/internal/types"
)

func AgentConversationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UnifiedStreamConversationRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 记录请求日志
		logger := logx.WithContext(r.Context())
		logger.Infow("收到多Agent对话请求",
			logx.Field("messageType", req.MessageType),
			logx.Field("sessionId", req.SessionId),
			logx.Field("userAge", req.UserAge),
		)

		// 调用AgentLogic
		l := logic.NewAgentLogic(r.Context(), svcCtx)
		if err := l.StreamAgentConversation(w, req); err != nil {
			logger.Errorw("多Agent对话处理失败", logx.Field("error", err))
			// 发送错误事件
			errorEvent := types.StreamEvent{
				Type:      "error",
				Content:   map[string]interface{}{"message": err.Error()},
				SessionId: req.SessionId,
			}
			errorJSON, _ := json.Marshal(errorEvent)
			w.Header().Set("Content-Type", "text/event-stream")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("event: error\ndata: " + string(errorJSON) + "\n\n"))
			return
		}
	}
}

