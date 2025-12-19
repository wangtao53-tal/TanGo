// package
package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/tango/explore/internal/logic"
	"github.com/tango/explore/internal/svc"
	"github.com/tango/explore/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func IdentifyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.IdentifyRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		// 创建带超时的 context（90秒，与配置文件中的 Timeout 一致）
		ctx, cancel := context.WithTimeout(r.Context(), 90*time.Second)
		defer cancel()

		l := logic.NewIdentifyLogic(ctx, svcCtx)
		resp, err := l.Identify(&req)
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
