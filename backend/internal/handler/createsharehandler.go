package handler

import (
	"fmt"
	"net/http"

	"github.com/tango/explore/internal/logic"
	"github.com/tango/explore/internal/svc"
	"github.com/tango/explore/internal/types"
	"github.com/tango/explore/internal/utils"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CreateShareHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 限制请求体大小为20MB（因为可能包含多张图片的base64数据）
		r.Body = http.MaxBytesReader(w, r.Body, 20*1024*1024) // 20MB限制

		var req types.CreateShareRequest
		if err := httpx.Parse(r, &req); err != nil {
			// 提供更友好的错误消息
			if err.Error() == "http: request body too large" {
				httpx.Error(w, utils.NewAPIError(413, "请求数据过大，请减少分享的探索记录数量（最多10条）"))
			} else {
				httpx.Error(w, err)
			}
			return
		}

		// 限制探索记录数量为最多10条
		if len(req.ExplorationRecords) > 10 {
			httpx.Error(w, utils.NewAPIError(400, fmt.Sprintf("探索记录数量过多，最多只能分享10条，当前有%d条", len(req.ExplorationRecords))))
			return
		}

		l := logic.NewCreateShareLogic(r.Context(), svcCtx)
		resp, err := l.CreateShare(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
