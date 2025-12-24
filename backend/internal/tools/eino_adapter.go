package tools

import (
	"context"
	"encoding/json"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// EinoToolAdapter 将我们的Tool接口适配为eino的tool.BaseTool
type EinoToolAdapter struct {
	tool Tool
}

// NewEinoToolAdapter 创建适配器
func NewEinoToolAdapter(t Tool) *EinoToolAdapter {
	return &EinoToolAdapter{tool: t}
}

// Info 返回工具信息（实现tool.BaseTool接口）
func (a *EinoToolAdapter) Info(ctx context.Context) (*schema.ToolInfo, error) {
	params := a.tool.Parameters()

	// 转换参数定义
	paramInfos := make(map[string]*schema.ParameterInfo)
	if props, ok := params["properties"].(map[string]interface{}); ok {
		for name, prop := range props {
			if propMap, ok := prop.(map[string]interface{}); ok {
				paramType := getString(propMap, "type", "string")
				// 转换类型字符串到schema.DataType
				var dataType schema.DataType
				switch paramType {
				case "string":
					dataType = schema.String
				case "number", "integer":
					dataType = schema.Number
				case "boolean":
					dataType = schema.Boolean
				default:
					dataType = schema.String
				}

				paramInfo := &schema.ParameterInfo{
					Type: dataType,
				}
				if desc, ok := propMap["description"].(string); ok {
					paramInfo.Desc = desc
				}
				paramInfos[name] = paramInfo
			}
		}
	}

	return &schema.ToolInfo{
		Name:        a.tool.Name(),
		Desc:        a.tool.Description(),
		ParamsOneOf: schema.NewParamsOneOfByParams(paramInfos),
	}, nil
}

// InvokableRun 执行工具（实现tool.BaseTool接口）
func (a *EinoToolAdapter) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	// 解析参数
	var params map[string]interface{}
	if argumentsInJSON != "" {
		if err := json.Unmarshal([]byte(argumentsInJSON), &params); err != nil {
			return "", err
		}
	} else {
		params = make(map[string]interface{})
	}

	// 执行工具
	result, err := a.tool.Execute(ctx, params)
	if err != nil {
		return "", err
	}

	// 将结果转换为JSON字符串
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	return string(resultJSON), nil
}

// getString 辅助函数
func getString(m map[string]interface{}, key, defaultValue string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return defaultValue
}

// ConvertToolNamesToEinoTools 根据工具名称列表从注册表中获取工具并转换为eino格式
func ConvertToolNamesToEinoTools(registry *ToolRegistry, toolNames []string, ctx context.Context) ([]*schema.ToolInfo, error) {
	tools := make([]Tool, 0, len(toolNames))

	for _, name := range toolNames {
		if tool, ok := registry.GetTool(name); ok {
			tools = append(tools, tool)
		}
	}

	return ConvertToEinoTools(tools, ctx)
}

// ConvertToEinoTools 将Tool接口列表转换为eino的schema.ToolInfo列表
func ConvertToEinoTools(toolList []Tool, ctx context.Context) ([]*schema.ToolInfo, error) {
	toolInfos := make([]*schema.ToolInfo, 0, len(toolList))

	for _, t := range toolList {
		if t != nil {
			adapter := NewEinoToolAdapter(t)
			info, err := adapter.Info(ctx)
			if err != nil {
				return nil, err
			}
			toolInfos = append(toolInfos, info)
		}
	}

	return toolInfos, nil
}

