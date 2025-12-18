package utils

import (
	"github.com/zeromicro/go-zero/core/logx"
)

// InitLogger 初始化日志配置
func InitLogger() {
	// 可以在这里配置日志格式、输出位置等
	// 当前使用go-zero默认配置
	logx.DisableStat()
}

