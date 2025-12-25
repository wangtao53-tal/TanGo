//
// @Copyright (c) 2025 by Mochi, All Rights Reserved.
// @Title: HealthAvailableHandler
// @Description: 健康检查 - 可用性检查端点，用于 Kubernetes readiness probe
// @FilePath: /backend/internal/handler/healthavailablehandler.go
// @Author: wangtao53-tal
// @Date: 2025-12-19 20:42:54
// @LastEditors: wangtao53-tal
// @LastEditTime: 2025-12-19 20:47:23

package handler

import (
	"net/http"

	"github.com/tango/explore/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// HealthAvailableHandler 健康检查 - 可用性检查端点
// 用于 Kubernetes readiness probe，检查服务是否准备好接受流量
func HealthAvailableHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]string{"status": "available"}
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}
