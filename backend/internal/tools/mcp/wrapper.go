package mcp

import (
	"context"
	"fmt"

	"github.com/tango/explore/internal/tools"
	"github.com/zeromicro/go-zero/core/logx"
)

// MCPToolWrapper MCP工具包装器
// 将MCP资源包装为Tool接口实现
type MCPToolWrapper struct {
	name        string
	description string
	client      *MCPClient
	logger      logx.Logger
	parameters  map[string]interface{}
}

// NewMCPToolWrapper 创建MCP工具包装器
func NewMCPToolWrapper(name string, description string, client *MCPClient, logger logx.Logger, parameters map[string]interface{}) tools.Tool {
	return &MCPToolWrapper{
		name:        name,
		description: description,
		client:      client,
		logger:      logger,
		parameters:  parameters,
	}
}

// Name 返回工具名称
func (t *MCPToolWrapper) Name() string {
	return t.name
}

// Description 返回工具描述
func (t *MCPToolWrapper) Description() string {
	return t.description
}

// Parameters 返回工具参数定义
func (t *MCPToolWrapper) Parameters() map[string]interface{} {
	if t.parameters != nil {
		return t.parameters
	}
	// 默认参数定义
	return map[string]interface{}{
		"type":       "object",
		"properties": map[string]interface{}{},
		"required":   []string{},
	}
}

// Execute 执行工具
func (t *MCPToolWrapper) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	if t.client == nil {
		return nil, fmt.Errorf("MCP客户端未初始化")
	}

	t.logger.Infow("执行MCP工具",
		logx.Field("tool", t.name),
		logx.Field("params", params),
	)

	// 调用MCP资源
	result, err := t.client.CallResource(ctx, t.name, params)
	if err != nil {
		t.logger.Errorw("MCP资源调用失败",
			logx.Field("tool", t.name),
			logx.Field("error", err),
		)
		return nil, err
	}

	return result, nil
}

