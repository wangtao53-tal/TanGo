//
// @Copyright (c) 2025 by Mochi, All Rights Reserved.
// @Title: HealthAliveHandler
// @Description: 健康检查 - 存活检查端点，用于 Kubernetes liveness probe
// @FilePath: /backend/internal/handler/healthalivehandler.go
// @Author: wangtao53-tal
// @Date: 2025-12-19 20:42:51
// @LastEditors: wangtao53-tal
// @LastEditTime: 2025-12-19 20:48:20

package handler

import (
	"net/http"

	"github.com/tango/explore/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// HealthAliveHandler 健康检查 - 存活检查端点
// 用于 Kubernetes liveness probe，检查服务是否存活
func HealthAliveHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]string{"status": "alive"}
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}
