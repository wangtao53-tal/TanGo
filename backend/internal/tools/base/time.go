package base

import (
	"context"
	"time"

	"github.com/tango/explore/internal/tools"
	"github.com/zeromicro/go-zero/core/logx"
)

// GetCurrentTimeTool get_current_time工具实现
// 用于获取当前时间，主要用于Science Agent
type GetCurrentTimeTool struct {
	logger logx.Logger
}

// NewGetCurrentTimeTool 创建get_current_time工具实例
func NewGetCurrentTimeTool(logger logx.Logger) tools.Tool {
	return &GetCurrentTimeTool{
		logger: logger,
	}
}

// Name 返回工具名称
func (t *GetCurrentTimeTool) Name() string {
	return "get_current_time"
}

// Description 返回工具描述
func (t *GetCurrentTimeTool) Description() string {
	return "获取当前时间信息，包括日期、时间、星期等。不需要参数。"
}

// Parameters 返回工具参数定义（JSON Schema格式）
func (t *GetCurrentTimeTool) Parameters() map[string]interface{} {
	return map[string]interface{}{
		"type":       "object",
		"properties": map[string]interface{}{},
		"required":   []string{},
	}
}

// Execute 执行工具
func (t *GetCurrentTimeTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	// 设置超时
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	t.logger.Infow("⏰ 执行get_current_time工具",
		logx.Field("params", params),
	)

	// 获取当前时间
	now := time.Now()

	// 格式化时间信息
	result := map[string]interface{}{
		"datetime":    now.Format("2006-01-02 15:04:05"),
		"date":        now.Format("2006年1月2日"),
		"time":        now.Format("15:04:05"),
		"weekday":     getWeekdayCN(now.Weekday()),
		"year":        now.Year(),
		"month":       int(now.Month()),
		"day":         now.Day(),
		"hour":        now.Hour(),
		"minute":      now.Minute(),
		"second":      now.Second(),
		"timestamp":  now.Unix(),
	}

	t.logger.Infow("⏰ get_current_time工具执行完成",
		logx.Field("datetime", result["datetime"]),
		logx.Field("date", result["date"]),
		logx.Field("time", result["time"]),
		logx.Field("weekday", result["weekday"]),
		logx.Field("result", result),
	)

	return result, nil
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

