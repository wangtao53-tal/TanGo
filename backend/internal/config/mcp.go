package config

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/zeromicro/go-zero/core/logx"
)

// MCPConfig MCP服务器配置
type MCPConfig struct {
	// 是否启用MCP功能
	Enabled bool `json:",optional,env=MCP_ENABLED"`

	// MCP服务器配置列表（从环境变量MCP_SERVERS读取JSON配置）
	Servers map[string]MCPServerConfig `json:",optional"`

	// MCP配置文件路径（可选，从文件读取配置）
	ConfigPath string `json:",optional,env=MCP_CONFIG_PATH"`
}

// MCPServerConfig MCP服务器配置
type MCPServerConfig struct {
	// 服务器类型：url, command
	Type string `json:"type,omitempty"`

	// URL（用于SSE或HTTP类型的服务器）
	URL string `json:"url,omitempty"`

	// 命令（用于command类型的服务器）
	Command string `json:"command,omitempty"`

	// 命令参数
	Args []string `json:"args,omitempty"`

	// HTTP请求头（用于需要认证的服务器）
	Headers map[string]string `json:"headers,omitempty"`

	// 是否启用此服务器
	Enabled bool `json:"enabled,omitempty"`
}

var (
	mcpConfigOnce sync.Once
	mcpConfig     *MCPConfig
)

// LoadMCPConfig 加载MCP配置
func LoadMCPConfig(logger logx.Logger) *MCPConfig {
	mcpConfigOnce.Do(func() {
		mcpConfig = &MCPConfig{
			Enabled: false,
			Servers: make(map[string]MCPServerConfig),
		}

		// 从环境变量读取MCP_ENABLED
		if enabled := os.Getenv("MCP_ENABLED"); enabled == "true" {
			mcpConfig.Enabled = true
		}

		// 从环境变量读取MCP_SERVERS（JSON格式）
		if serversJSON := os.Getenv("MCP_SERVERS"); serversJSON != "" {
			if err := json.Unmarshal([]byte(serversJSON), &mcpConfig.Servers); err != nil {
				logger.Errorw("解析MCP_SERVERS环境变量失败", logx.Field("error", err))
			} else {
				logger.Infow("从环境变量加载MCP服务器配置",
					logx.Field("server_count", len(mcpConfig.Servers)),
				)
			}
		}

		// 从配置文件读取（如果指定了配置路径）
		if configPath := os.Getenv("MCP_CONFIG_PATH"); configPath != "" {
			if data, err := os.ReadFile(configPath); err == nil {
				if err := json.Unmarshal(data, mcpConfig); err != nil {
					logger.Errorw("解析MCP配置文件失败", logx.Field("error", err))
				} else {
					logger.Infow("从配置文件加载MCP配置", logx.Field("path", configPath))
				}
			}
		}

		// 如果没有配置，使用默认配置（tal_time）
		if len(mcpConfig.Servers) == 0 && mcpConfig.Enabled {
			mcpConfig.Servers["tal_time"] = MCPServerConfig{
				Type:    "url",
				URL:     "http://mcp.tal.com/time/sse",
				Enabled: true,
			}
			logger.Info("使用默认MCP配置（tal_time）")
		}
	})

	return mcpConfig
}

// GetMCPConfig 获取MCP配置（单例模式）
func GetMCPConfig(logger logx.Logger) *MCPConfig {
	if mcpConfig == nil {
		return LoadMCPConfig(logger)
	}
	return mcpConfig
}

