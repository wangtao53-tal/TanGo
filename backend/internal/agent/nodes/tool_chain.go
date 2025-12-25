package nodes

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"github.com/tango/explore/internal/tools"
	"github.com/zeromicro/go-zero/core/logx"
)

// ToolChain 工具调用链处理器
// 支持多轮工具调用，限制调用深度，控制超时
type ToolChain struct {
	maxDepth     int           // 最大调用深度
	timeout      time.Duration // 超时时间
	toolRegistry *tools.ToolRegistry
	logger       logx.Logger
}

// NewToolChain 创建工具调用链处理器
func NewToolChain(toolRegistry *tools.ToolRegistry, logger logx.Logger) *ToolChain {
	return &ToolChain{
		maxDepth:     3,                // 默认最大深度3层
		timeout:      10 * time.Second, // 默认超时10秒
		toolRegistry: toolRegistry,
		logger:       logger,
	}
}

// ExecuteToolChain 执行工具调用链
// 支持多轮工具调用，直到达到最大深度或没有工具调用请求
func (tc *ToolChain) ExecuteToolChain(
	ctx context.Context,
	messages []*schema.Message,
	chatModel model.ChatModel,
	initialTools []string, // 初始推荐的工具列表（可选）
) ([]*schema.Message, []string, map[string]interface{}, error) {
	// 设置超时
	ctx, cancel := context.WithTimeout(ctx, tc.timeout)
	defer cancel()

	currentMessages := messages
	toolsUsed := []string{}
	toolResults := make(map[string]interface{})
	depth := 0

	for depth < tc.maxDepth {
		// 调用ChatModel
		result, err := chatModel.Generate(ctx, currentMessages)
		if err != nil {
			tc.logger.Errorw("工具调用链中ChatModel调用失败",
				logx.Field("depth", depth),
				logx.Field("error", err),
			)
			return currentMessages, toolsUsed, toolResults, err
		}

		// 记录ChatModel返回的原始结果（用于调试）
		tc.logger.Infow("📨 ChatModel返回结果",
			logx.Field("depth", depth),
			logx.Field("hasContent", result.Content != ""),
			logx.Field("contentLength", len(result.Content)),
			logx.Field("hasToolCalls", len(result.ToolCalls) > 0),
			logx.Field("toolCallsCount", len(result.ToolCalls)),
			logx.Field("toolCalls", func() []string {
				if len(result.ToolCalls) == 0 {
					return []string{}
				}
				names := make([]string, 0, len(result.ToolCalls))
				for _, tc := range result.ToolCalls {
					if len(tc.Function.Name) > 0 {
						names = append(names, tc.Function.Name)
					}
				}
				return names
			}()),
		)

		// 检查是否有工具调用请求
		if len(result.ToolCalls) == 0 {
			// 没有工具调用，结束链
			if depth == 0 {
				// 第一轮就没有工具调用
				tc.logger.Infow("🔚 工具调用链结束（第一轮无工具调用请求）",
					logx.Field("depth", depth),
					logx.Field("toolsUsed", toolsUsed),
					logx.Field("initialTools", initialTools),
					logx.Field("resultContent", result.Content),
				)
			} else {
				// 后续轮次没有工具调用，说明工具调用已完成并整合
				tc.logger.Infow("✅ 工具调用链完成（工具结果已整合）",
					logx.Field("depth", depth),
					logx.Field("toolsUsed", toolsUsed),
					logx.Field("totalToolsUsed", len(toolsUsed)),
					logx.Field("resultContent", result.Content),
				)
			}
			// 将最终结果添加到消息列表
			currentMessages = append(currentMessages, result)
			return currentMessages, toolsUsed, toolResults, nil
		}

		// 记录检测到的工具调用请求
		tc.logger.Infow("🔍 检测到工具调用请求",
			logx.Field("depth", depth),
			logx.Field("tool_call_count", len(result.ToolCalls)),
			logx.Field("tool_calls", func() []string {
				names := make([]string, 0, len(result.ToolCalls))
				for _, tc := range result.ToolCalls {
					if len(tc.Function.Name) > 0 {
						names = append(names, tc.Function.Name)
					}
				}
				return names
			}()),
		)

		// 执行工具调用
		toolMessages, roundTools, roundResults := tc.executeToolRound(ctx, result.ToolCalls)
		if len(toolMessages) == 0 {
			// 没有成功执行的工具，结束链
			tc.logger.Errorw("工具调用链中断：没有成功执行的工具",
				logx.Field("depth", depth),
			)
			currentMessages = append(currentMessages, result)
			return currentMessages, toolsUsed, toolResults, nil
		}

		// 记录工具使用
		toolsUsed = append(toolsUsed, roundTools...)
		for k, v := range roundResults {
			toolResults[k] = v
		}

		// 添加工具结果到消息列表，继续下一轮
		currentMessages = append(currentMessages, result)
		currentMessages = append(currentMessages, toolMessages...)

		depth++
		tc.logger.Infow("🔄 工具调用链继续（等待ChatModel整合工具结果）",
			logx.Field("depth", depth),
			logx.Field("toolsUsed", roundTools),
			logx.Field("totalToolsUsed", len(toolsUsed)),
			logx.Field("toolResults", func() []string {
				keys := make([]string, 0, len(toolResults))
				for k := range toolResults {
					keys = append(keys, k)
				}
				return keys
			}()),
		)
	}

	// 达到最大深度，最后一次调用ChatModel整合结果
	finalResult, err := chatModel.Generate(ctx, currentMessages)
	if err != nil {
		tc.logger.Errorw("工具调用链最终整合失败", logx.Field("error", err))
		return currentMessages, toolsUsed, toolResults, err
	}

	currentMessages = append(currentMessages, finalResult)
	return currentMessages, toolsUsed, toolResults, nil
}

// executeToolRound 执行一轮工具调用
func (tc *ToolChain) executeToolRound(ctx context.Context, toolCalls []schema.ToolCall) ([]*schema.Message, []string, map[string]interface{}) {
	toolMessages := make([]*schema.Message, 0, len(toolCalls))
	toolsUsed := []string{}
	toolResults := make(map[string]interface{})

	for _, toolCall := range toolCalls {
		if len(toolCall.Function.Name) == 0 {
			continue
		}

		toolName := toolCall.Function.Name

		// 记录工具调用开始
		tc.logger.Infow("🔧 开始执行工具调用",
			logx.Field("tool", toolName),
			logx.Field("tool_call_id", toolCall.ID),
		)

		tool, ok := tc.toolRegistry.GetTool(toolName)
		if !ok {
			tc.logger.Errorw("❌ 工具未找到", logx.Field("tool", toolName))
			continue
		}

		// 解析参数
		params := make(map[string]interface{})
		if toolCall.Function.Arguments != "" {
			if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &params); err != nil {
				tc.logger.Errorw("❌ 工具参数解析失败",
					logx.Field("tool", toolName),
					logx.Field("arguments", toolCall.Function.Arguments),
					logx.Field("error", err),
				)
				continue
			}
		}

		// 记录工具参数
		if len(params) > 0 {
			tc.logger.Infow("📥 工具调用参数",
				logx.Field("tool", toolName),
				logx.Field("params", params),
			)
		} else {
			tc.logger.Infow("📥 工具调用参数（无参数）",
				logx.Field("tool", toolName),
			)
		}

		// 执行工具
		toolResult, err := tool.Execute(ctx, params)
		if err != nil {
			tc.logger.Errorw("❌ 工具调用失败",
				logx.Field("tool", toolName),
				logx.Field("error", err),
			)
			continue
		}

		// 记录工具使用
		toolsUsed = append(toolsUsed, toolName)
		toolResults[toolName] = toolResult

		// 特别处理时间工具的结果输出
		if toolName == "get_current_time" {
			if resultMap, ok := toolResult.(map[string]interface{}); ok {
				tc.logger.Infow("✅ 时间工具调用成功",
					logx.Field("tool", toolName),
					logx.Field("datetime", resultMap["datetime"]),
					logx.Field("date", resultMap["date"]),
					logx.Field("time", resultMap["time"]),
					logx.Field("weekday", resultMap["weekday"]),
					logx.Field("result", toolResult),
				)
			} else {
				tc.logger.Infow("✅ 时间工具调用成功",
					logx.Field("tool", toolName),
					logx.Field("result", toolResult),
				)
			}
		} else {
			// 其他工具的结果输出
			tc.logger.Infow("✅ 工具调用成功",
				logx.Field("tool", toolName),
				logx.Field("result_type", fmt.Sprintf("%T", toolResult)),
				logx.Field("result", toolResult),
			)
		}

		// 创建工具消息
		resultJSON, _ := json.Marshal(toolResult)
		toolMessage := schema.ToolMessage(string(resultJSON), toolCall.ID)
		toolMessages = append(toolMessages, toolMessage)
	}

	return toolMessages, toolsUsed, toolResults
}
