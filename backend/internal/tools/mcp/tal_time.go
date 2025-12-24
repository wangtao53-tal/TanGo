package mcp

import (
	"context"
	"time"

	"github.com/tango/explore/internal/config"
	"github.com/tango/explore/internal/tools"
	"github.com/zeromicro/go-zero/core/logx"
)

// TalTimeTool tal_time MCP工具包装
// 用于获取当前时间信息
type TalTimeTool struct {
	client *MCPClient
	logger logx.Logger
}

// NewTalTimeTool 创建tal_time工具实例
func NewTalTimeTool(cfg config.MCPServerConfig, logger logx.Logger) (tools.Tool, error) {
	client, err := NewMCPClient(cfg, logger)
	if err != nil {
		return nil, err
	}

	return &TalTimeTool{
		client: client,
		logger: logger,
	}, nil
}

// Name 返回工具名称
func (t *TalTimeTool) Name() string {
	return "tal_time"
}

// Description 返回工具描述
func (t *TalTimeTool) Description() string {
	return "获取当前时间信息，通过MCP服务器获取准确的时间数据。"
}

// Parameters 返回工具参数定义
func (t *TalTimeTool) Parameters() map[string]interface{} {
	return map[string]interface{}{
		"type":       "object",
		"properties": map[string]interface{}{},
		"required":   []string{},
	}
}

// Execute 执行工具
func (t *TalTimeTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	t.logger.Info("执行tal_time工具")

	// 尝试调用MCP资源
	if t.client != nil {
		result, err := t.client.CallResource(ctx, "time", params)
		if err == nil {
			return result, nil
		}
		// 如果MCP调用失败，降级到本地时间
		t.logger.Errorw("MCP调用失败，使用本地时间", logx.Field("error", err))
	}

	// 降级：返回本地时间
	now := time.Now()
	return map[string]interface{}{
		"datetime":   now.Format("2006-01-02 15:04:05"),
		"date":       now.Format("2006年1月2日"),
		"time":       now.Format("15:04:05"),
		"weekday":    getWeekdayCN(now.Weekday()),
		"year":       now.Year(),
		"month":      int(now.Month()),
		"day":        now.Day(),
		"hour":       now.Hour(),
		"minute":     now.Minute(),
		"second":     now.Second(),
		"timestamp":  now.Unix(),
		"source":     "local_fallback",
	}, nil
}

// getWeekdayCN 获取中文星期
func getWeekdayCN(weekday time.Weekday) string {
	weekdays := map[time.Weekday]string{
		time.Sunday:    "星期日",
		time.Monday:    "星期一",
		time.Tuesday:   "星期二",
		time.Wednesday: "星期三",
		time.Thursday:  "星期四",
		time.Friday:    "星期五",
		time.Saturday:  "星期六",
	}
	return weekdays[weekday]
}

