package mcp

import (
	"context"
	"fmt"

	"github.com/tango/explore/internal/config"
	"github.com/tango/explore/internal/tools"
	"github.com/zeromicro/go-zero/core/logx"
)

// DiscoverAndRegisterMCPTools 发现并注册MCP工具
func DiscoverAndRegisterMCPTools(mcpConfig *config.MCPConfig, registry *tools.ToolRegistry, logger logx.Logger) error {
	if !mcpConfig.Enabled {
		logger.Info("MCP功能未启用，跳过MCP工具发现")
		return nil
	}

	if len(mcpConfig.Servers) == 0 {
		logger.Info("未配置MCP服务器，跳过MCP工具发现")
		return nil
	}

	logger.Infow("开始发现MCP工具",
		logx.Field("server_count", len(mcpConfig.Servers)),
	)

	// 遍历所有MCP服务器
	for serverName, serverConfig := range mcpConfig.Servers {
		if !serverConfig.Enabled {
			logger.Infow("跳过未启用的MCP服务器", logx.Field("server", serverName))
			continue
		}

		// 根据服务器名称创建对应的工具
		switch serverName {
		case "tal_time":
			tool, err := NewTalTimeTool(serverConfig, logger)
			if err != nil {
				logger.Errorw("创建tal_time工具失败",
					logx.Field("error", err),
				)
				continue
			}
			registry.Register(tool)
			logger.Infow("已注册MCP工具", logx.Field("tool", "tal_time"))

		default:
			// 对于其他MCP服务器，尝试通用包装
			client, err := NewMCPClient(serverConfig, logger)
			if err != nil {
				logger.Errorw("创建MCP客户端失败",
					logx.Field("server", serverName),
					logx.Field("error", err),
				)
				continue
			}

			// 尝试发现资源
			resources, err := client.DiscoverResources(context.Background())
			if err != nil {
				logger.Errorw("发现MCP资源失败",
					logx.Field("server", serverName),
					logx.Field("error", err),
				)
				// 继续处理，即使发现失败
			}

			// 为每个资源创建工具包装器
			for _, resourceName := range resources {
				wrapper := NewMCPToolWrapper(
					resourceName,
					fmt.Sprintf("MCP资源: %s/%s", serverName, resourceName),
					client,
					logger,
					nil, // 使用默认参数定义
				)
				registry.Register(wrapper)
				logger.Infow("已注册MCP工具",
					logx.Field("tool", resourceName),
					logx.Field("server", serverName),
				)
			}
		}
	}

	logger.Info("MCP工具发现完成")
	return nil
}

