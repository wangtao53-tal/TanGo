package nodes

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
	"github.com/tango/explore/internal/config"
	"github.com/tango/explore/internal/tools"
	"github.com/tango/explore/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

// ScienceAgentNode Science Agent节点
type ScienceAgentNode struct {
	ctx         context.Context
	config      config.AIConfig
	logger      logx.Logger
	chatModel   model.ChatModel     // eino ChatModel 实例
	template    prompt.ChatTemplate // 消息模板
	toolRegistry *tools.ToolRegistry // 工具注册表
	initialized bool
}

// NewScienceAgentNode 创建Science Agent节点
func NewScienceAgentNode(ctx context.Context, cfg config.AIConfig, logger logx.Logger, toolRegistry *tools.ToolRegistry) (*ScienceAgentNode, error) {
	node := &ScienceAgentNode{
		ctx:          ctx,
		config:       cfg,
		logger:       logger,
		toolRegistry: toolRegistry,
	}

	if cfg.EinoBaseURL != "" && cfg.AppID != "" && cfg.AppKey != "" {
		if err := node.initChatModel(ctx); err != nil {
			logger.Errorw("初始化ChatModel失败，将使用Mock模式", logx.Field("error", err))
		} else {
			node.initialized = true
			logger.Info("✅ Science Agent节点已初始化ChatModel")
		}
	} else {
		logger.Info("未配置eino参数，Science Agent节点将使用Mock模式")
	}

	node.initTemplate()
	return node, nil
}

// initChatModel 初始化 ChatModel（支持工具调用）
func (n *ScienceAgentNode) initChatModel(ctx context.Context) error {
	modelName := n.selectRandomModel(n.config.TextGenerationModels)
	if modelName == "" {
		models := config.GetDefaultTextGenerationModels()
		if len(models) > 0 {
			modelName = n.selectRandomModel(models)
		}
		if modelName == "" {
			modelName = config.DefaultTextGenerationModel
		}
	}

	cfg := &ark.ChatModelConfig{
		Model: modelName,
	}

	if n.config.EinoBaseURL != "" {
		cfg.BaseURL = n.config.EinoBaseURL
	}

	if n.config.AppID != "" && n.config.AppKey != "" {
		cfg.APIKey = n.config.AppID + ":" + n.config.AppKey
	} else if n.config.AppKey != "" {
		cfg.APIKey = n.config.AppKey
	} else if n.config.AppID != "" {
		cfg.APIKey = n.config.AppID
	} else {
		return nil
	}

	chatModel, err := ark.NewChatModel(ctx, cfg)
	if err != nil {
		return err
	}

	// 注册工具到ChatModel
	if n.toolRegistry != nil {
		// 获取Science Agent可用的工具
		agentTools := n.toolRegistry.GetToolsForAgent("Science")
		if len(agentTools) > 0 {
			// 转换为eino工具信息
			toolInfos, err := tools.ConvertToEinoTools(agentTools, ctx)
			if err != nil {
				n.logger.Errorw("转换工具信息失败", logx.Field("error", err))
			} else if len(toolInfos) > 0 {
				// 绑定工具到ChatModel
				if err := chatModel.BindTools(toolInfos); err != nil {
					n.logger.Errorw("绑定工具到ChatModel失败", logx.Field("error", err))
				} else {
					n.logger.Infow("✅ 注册工具到Science Agent ChatModel",
						logx.Field("tool_count", len(toolInfos)),
						logx.Field("tools", func() []string {
							names := make([]string, 0, len(agentTools))
							for _, t := range agentTools {
								names = append(names, t.Name())
							}
							return names
						}()),
					)
				}
			}
		}
	}

	n.chatModel = chatModel
	return nil
}

// selectRandomModel 从模型列表中随机选择一个模型
func (n *ScienceAgentNode) selectRandomModel(models []string) string {
	if len(models) == 0 {
		return ""
	}
	if len(models) == 1 {
		return models[0]
	}
	rand.Seed(time.Now().UnixNano())
	return models[rand.Intn(len(models))]
}

// initTemplate 初始化消息模板
func (n *ScienceAgentNode) initTemplate() {
	n.template = prompt.FromMessages(schema.FString,
		schema.SystemMessage(`你是 Science Agent，一个直接和孩子对话的AI伙伴，用简单有趣的方式解释科学知识。

重要规则：
- 直接回答孩子的问题，就像朋友聊天一样
- 不要出现"跟小朋友可以这样聊"、"你可以说"等指导性语言
- 不要出现"你:"这样的对话示例格式
- 用"我"或直接称呼"你"（孩子）来对话
- 只回答一个知识点，用生活类比，不用术语
- 回答简洁，控制在合理长度内
- 让孩子感受到探索的乐趣

你可以调用的工具：
- simple_fact_lookup: 查找简单事实
- get_current_time: 获取当前时间
- image_generate_simple: 生成示意图（仅示意图）

如果工具调用失败，不依赖工具也能生成基本回答。

记住：你是直接和孩子对话的AI伙伴，不是给家长看的指导手册！`),
		schema.MessagesPlaceholder("chat_history", true),
		schema.UserMessage("{message}"),
	)
}

// GenerateScienceAnswer 生成科学回答
func (n *ScienceAgentNode) GenerateScienceAnswer(ctx context.Context, message string, objectName string, objectCategory string, userAge int, chatHistory []*schema.Message, maxSentences int, recommendedTools []string) (*types.DomainAgentResponse, error) {
	n.logger.Infow("执行Science Agent回答生成",
		logx.Field("message", message),
		logx.Field("objectName", objectName),
		logx.Field("userAge", userAge),
		logx.Field("maxSentences", maxSentences),
		logx.Field("recommendedTools", recommendedTools),
		logx.Field("useRealModel", n.initialized),
	)

	if n.initialized && n.chatModel != nil {
		return n.executeReal(ctx, message, objectName, objectCategory, userAge, chatHistory, maxSentences, recommendedTools)
	}

	return n.executeMock(message, objectName, userAge, maxSentences)
}

// executeMock Mock实现
func (n *ScienceAgentNode) executeMock(message string, objectName string, userAge int, maxSentences int) (*types.DomainAgentResponse, error) {
	content := "关于" + objectName + "的科学知识很有趣。"
	if userAge <= 6 {
		content = objectName + "就像我们身边的朋友一样，有很多有趣的特点。"
	} else if userAge <= 12 {
		content = objectName + "的科学原理可以用生活中的例子来解释，比如..."
	} else {
		content = objectName + "涉及的科学知识可以深入探索，让我们一起来了解。"
	}

	return &types.DomainAgentResponse{
		DomainType:  "Science",
		Content:     content,
		ToolsUsed:   []string{},
		ToolResults: make(map[string]interface{}),
	}, nil
}

// executeReal 真实eino实现（支持工具调用）
func (n *ScienceAgentNode) executeReal(ctx context.Context, message string, objectName string, objectCategory string, userAge int, chatHistory []*schema.Message, maxSentences int, recommendedTools []string) (*types.DomainAgentResponse, error) {
	// 根据推荐的工具动态构建SystemMessage
	systemMessage := n.buildSystemMessageWithTools(recommendedTools)
	
	// 构建消息列表
	messages := []*schema.Message{
		schema.SystemMessage(systemMessage),
	}
	
	// 添加对话历史
	if len(chatHistory) > 0 {
		messages = append(messages, chatHistory...)
	}
	
	// 添加用户消息（包含对象信息）
	userMsg := message
	if objectName != "" {
		userMsg = fmt.Sprintf("识别对象：%s（%s）。问题：%s", objectName, objectCategory, message)
	}
	messages = append(messages, schema.UserMessage(userMsg))

	// 确保消息格式正确，移除任何可能导致工具调用错误的字段
	cleanMessages := make([]*schema.Message, 0, len(messages))
	for _, msg := range messages {
		if msg != nil && msg.Role != "" {
			cleanMsg := &schema.Message{
				Role:    msg.Role,
				Content: msg.Content,
			}
			cleanMessages = append(cleanMessages, cleanMsg)
		}
	}

	// 调用ChatModel，可能返回工具调用请求
	result, err := n.chatModel.Generate(ctx, cleanMessages)
	if err != nil {
		n.logger.Errorw("ChatModel调用失败", logx.Field("error", err))
		return n.executeMock(message, objectName, userAge, maxSentences)
	}

	// 检查是否有工具调用请求
	toolsUsed := []string{}
	toolResults := make(map[string]interface{})

	if len(result.ToolCalls) > 0 {
		n.logger.Infow("检测到工具调用请求",
			logx.Field("tool_call_count", len(result.ToolCalls)),
		)

		// 执行工具调用
		toolMessages := make([]*schema.Message, 0, len(result.ToolCalls))
		for _, toolCall := range result.ToolCalls {
			// 检查toolCall是否有效
			if len(toolCall.Function.Name) == 0 {
				continue
			}

			toolName := toolCall.Function.Name
			if n.toolRegistry == nil {
				n.logger.Errorw("工具注册表未初始化", logx.Field("tool", toolName))
				continue
			}

			tool, ok := n.toolRegistry.GetTool(toolName)
			if !ok {
				n.logger.Errorw("工具未找到", logx.Field("tool", toolName))
				continue
			}

			// 解析参数
			params := make(map[string]interface{})
			if toolCall.Function.Arguments != "" {
				if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &params); err != nil {
					n.logger.Errorw("工具参数解析失败",
						logx.Field("tool", toolName),
						logx.Field("error", err),
					)
					continue
				}
			}

			// 执行工具
			toolResult, err := tool.Execute(ctx, params)
			if err != nil {
				n.logger.Errorw("工具调用失败",
					logx.Field("tool", toolName),
					logx.Field("error", err),
				)
				continue
			}

			// 记录工具使用
			toolsUsed = append(toolsUsed, toolName)
			toolResults[toolName] = toolResult

			// 创建工具消息
			resultJSON, _ := json.Marshal(toolResult)
			toolMessage := schema.ToolMessage(string(resultJSON), toolCall.ID)
			toolMessages = append(toolMessages, toolMessage)

			n.logger.Infow("工具调用成功",
				logx.Field("tool", toolName),
			)
		}

		// 如果有工具调用结果，重新调用ChatModel整合结果
		if len(toolMessages) > 0 {
			// 添加工具结果到消息列表
			cleanMessages = append(cleanMessages, result)
			cleanMessages = append(cleanMessages, toolMessages...)

			// 重新调用ChatModel，整合工具结果
			finalResult, err := n.chatModel.Generate(ctx, cleanMessages)
			if err != nil {
				n.logger.Errorw("整合工具结果失败", logx.Field("error", err))
				// 降级：使用原始结果
				content := result.Content
				return &types.DomainAgentResponse{
					DomainType:  "Science",
					Content:     content,
					ToolsUsed:   toolsUsed,
					ToolResults: toolResults,
				}, nil
			}

			result = finalResult
		}
	}

	content := result.Content
	// 限制句子数量（简单实现，按句号分割）
	sentences := splitSentences(content)
	if len(sentences) > maxSentences {
		sentences = sentences[:maxSentences]
		content = joinSentences(sentences)
	}

	return &types.DomainAgentResponse{
		DomainType:  "Science",
		Content:     content,
		ToolsUsed:   toolsUsed,
		ToolResults: toolResults,
	}, nil
}

// splitSentences 分割句子
func splitSentences(text string) []string {
	// 简单实现：按句号、问号、感叹号分割
	// TODO: 更智能的句子分割
	return []string{text}
}

// joinSentences 连接句子
func joinSentences(sentences []string) string {
	result := ""
	for i, s := range sentences {
		if i > 0 {
			result += "。"
		}
		result += s
	}
	return result
}

// buildSystemMessageWithTools 根据推荐的工具构建SystemMessage
func (n *ScienceAgentNode) buildSystemMessageWithTools(recommendedTools []string) string {
	baseMessage := `你是 Science Agent，一个直接和孩子对话的AI伙伴，用简单有趣的方式解释科学知识。

重要规则：
- 直接回答孩子的问题，就像朋友聊天一样
- 不要出现"跟小朋友可以这样聊"、"你可以说"等指导性语言
- 不要出现"你:"这样的对话示例格式
- 用"我"或直接称呼"你"（孩子）来对话
- 只回答一个知识点，用生活类比，不用术语
- 回答简洁，控制在合理长度内
- 让孩子感受到探索的乐趣`

	// 如果有推荐的工具，添加到SystemMessage中
	if len(recommendedTools) > 0 {
		toolDescriptions := n.getToolDescriptions(recommendedTools)
		if toolDescriptions != "" {
			baseMessage += "\n\n你可以调用的工具：\n" + toolDescriptions
			baseMessage += "\n\n重要：当问题需要使用工具时，你必须调用相应的工具来获取信息，然后再回答。"
			baseMessage += "\n\n工具调用规则："
			baseMessage += "\n- 如果问时间相关的问题（几点了、现在几点、什么时候、现在几时），必须调用get_current_time工具"
			baseMessage += "\n- 如果问科学事实或知识，必须调用simple_fact_lookup工具"
			baseMessage += "\n- 如果需要示意图，可以调用image_generate_simple工具"
			baseMessage += "\n- 调用工具时，使用function calling格式，不要直接回答"
			baseMessage += "\n\n如果工具调用失败，不依赖工具也能生成基本回答。"
		}
	} else {
		// 默认工具列表
		baseMessage += "\n\n你可以调用的工具：\n- simple_fact_lookup: 查找简单事实\n- get_current_time: 获取当前时间\n- image_generate_simple: 生成示意图（仅示意图）\n\n如果工具调用失败，不依赖工具也能生成基本回答。"
	}

	baseMessage += "\n\n记住：你是直接和孩子对话的AI伙伴，不是给家长看的指导手册！"
	return baseMessage
}

// getToolDescriptions 获取工具描述列表
func (n *ScienceAgentNode) getToolDescriptions(toolNames []string) string {
	if n.toolRegistry == nil {
		return ""
	}

	descriptions := []string{}
	for _, toolName := range toolNames {
		tool, ok := n.toolRegistry.GetTool(toolName)
		if ok {
			descriptions = append(descriptions, fmt.Sprintf("- %s: %s", tool.Name(), tool.Description()))
		}
	}

	return strings.Join(descriptions, "\n")
}

