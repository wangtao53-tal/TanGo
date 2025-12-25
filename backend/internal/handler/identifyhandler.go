// package
package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/tango/explore/internal/logic"
	"github.com/tango/explore/internal/svc"
	"github.com/tango/explore/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func IdentifyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 性能监控：记录请求开始时间
		startTime := time.Now()

		// 优化：限制请求体大小，防止过大图片导致内存问题
		r.Body = http.MaxBytesReader(w, r.Body, 10*1024*1024) // 10MB限制

		var req types.IdentifyRequest
		if err := httpx.Parse(r, &req); err != nil {
			// 优化：提供更友好的错误消息
			if err.Error() == "http: request body too large" {
				httpx.Error(w, fmt.Errorf("图片数据过大，请使用小于10MB的图片"))
			} else {
				httpx.Error(w, err)
			}
			return
		}

		// 优化：提前验证图片数据大小
		if len(req.Image) > 10*1024*1024 {
			httpx.Error(w, fmt.Errorf("图片数据过大，请使用小于10MB的图片"))
			return
		}

		// 优化：调整超时时间到45秒（原来90秒可能过长）
		// 模型调用超时是60秒，handler层设置为45秒可以更快失败
		ctx, cancel := context.WithTimeout(r.Context(), 45*time.Second)
		defer cancel()

		l := logic.NewIdentifyLogic(ctx, svcCtx)
		resp, err := l.Identify(&req)

		// 性能监控：记录请求耗时
		duration := time.Since(startTime)
		if err != nil {
			// 检查是否是超时错误
			if ctx.Err() == context.DeadlineExceeded {
				logx.WithContext(ctx).Errorw("识别请求超时",
					logx.Field("duration_ms", duration.Milliseconds()),
					logx.Field("timeout", true),
				)
				httpx.ErrorCtx(ctx, w, fmt.Errorf("识别请求超时，请稍后重试"))
			} else {
				httpx.Error(w, err)
			}
		} else {
			httpx.OkJson(w, resp)
		}

		// 记录性能指标（即使出错也记录）
		logx.WithContext(ctx).Infow("识别接口请求完成",
			logx.Field("duration_ms", duration.Milliseconds()),
			logx.Field("duration_sec", duration.Seconds()),
			logx.Field("success", err == nil),
			logx.Field("timeout", ctx.Err() == context.DeadlineExceeded),
			logx.Field("imageSize", len(req.Image)),
		)
	}
}
