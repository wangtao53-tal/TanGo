package tools

import (
	"github.com/zeromicro/go-zero/core/logx"
)

// InitDefaultTools 初始化默认工具到全局注册表
// 注意：为了避免import cycle，工具实例化在调用方完成
func InitDefaultTools(logger logx.Logger, toolsToRegister []Tool) {
	registry := GetDefaultRegistry(logger)

	// 注册传入的工具
	for _, tool := range toolsToRegister {
		registry.Register(tool)
	}

	logger.Info("✅ 默认工具已注册到全局注册表")
}

