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

func GenerateCardsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GenerateCardsRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		// 创建带超时的 context（120秒，因为需要生成3张卡片，每张可能需要较长时间）
		ctx, cancel := context.WithTimeout(r.Context(), 120*time.Second)
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
