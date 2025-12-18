package handler

import (
	"net/http"
	"strings"

	"github.com/tango/explore/internal/logic"
	"github.com/tango/explore/internal/svc"
	"github.com/tango/explore/internal/utils"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetShareHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从路径中提取shareId: /api/share/:shareId
		path := r.URL.Path
		parts := strings.Split(path, "/")
		var shareId string
		for i, part := range parts {
			if part == "share" && i+1 < len(parts) {
				shareId = parts[i+1]
				break
			}
		}
		if shareId == "" {
			httpx.Error(w, utils.NewAPIError(400, "分享链接ID不能为空"))
			return
		}

		l := logic.NewGetShareLogic(r.Context(), svcCtx)
		resp, err := l.GetShare(shareId)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
