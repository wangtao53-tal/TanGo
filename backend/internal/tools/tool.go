package tools

import (
	"context"
)

// Tool 工具接口定义
// 所有工具必须实现此接口，包括基础工具和MCP工具包装器
type Tool interface {
	// Name 返回工具名称，用于注册和调用
	Name() string

	// Description 返回工具描述，用于Eino工具注册
	Description() string

	// Execute 执行工具，接收参数并返回结果
	// params: 工具参数，map格式，key为参数名，value为参数值
	// 返回: 工具执行结果，可以是任意类型
	Execute(ctx context.Context, params map[string]interface{}) (interface{}, error)

	// Parameters 返回工具参数定义，用于Eino工具注册
	// 返回格式应符合JSON Schema规范
	Parameters() map[string]interface{}
}

