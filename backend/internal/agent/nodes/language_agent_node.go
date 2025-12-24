package nodes

import (
	"context"
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

// LanguageAgentNode Language Agent节点
type LanguageAgentNode struct {
	ctx          context.Context
	config       config.AIConfig
	logger       logx.Logger
	chatModel    model.ChatModel     // eino ChatModel 实例
	template     prompt.ChatTemplate // 消息模板
	toolRegistry *tools.ToolRegistry // 工具注册表
	initialized  bool
}

// NewLanguageAgentNode 创建Language Agent节点
func NewLanguageAgentNode(ctx context.Context, cfg config.AIConfig, logger logx.Logger, toolRegistry *tools.ToolRegistry) (*LanguageAgentNode, error) {
	node := &LanguageAgentNode{
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
			logger.Info("✅ Language Agent节点已初始化ChatModel")
		}
	} else {
		logger.Info("未配置eino参数，Language Agent节点将使用Mock模式")
	}

	node.initTemplate()
	return node, nil
}

// initChatModel 初始化 ChatModel（支持工具调用）
func (n *LanguageAgentNode) initChatModel(ctx context.Context) error {
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
		// 获取Language Agent可用的工具
		agentTools := n.toolRegistry.GetToolsForAgent("Language")
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
					n.logger.Infow("✅ 注册工具到Language Agent ChatModel",
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
func (n *LanguageAgentNode) selectRandomModel(models []string) string {
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
func (n *LanguageAgentNode) initTemplate() {
	n.template = prompt.FromMessages(schema.FString,
		schema.SystemMessage(`你是 Language Agent，一个直接和孩子对话的AI伙伴，帮助孩子用语言表达自己的想法。

重要规则：
- 直接回答孩子的问题，就像朋友聊天一样
- 不要出现"跟小朋友可以这样聊"、"你可以说"等指导性语言
- 不要出现"你:"这样的对话示例格式
- 用"我"或直接称呼"你"（孩子）来对话
- 让孩子"说得出口"，不讲语法规则
- 用孩子日常语言，包含可模仿的句子
- 让孩子感受到表达的乐趣

你可以调用：
- simple_dictionary: 查找单词
- pronunciation_hint: 发音提示

如果工具调用失败，不依赖工具也能生成基本回答。

记住：你是直接和孩子对话的AI伙伴，不是给家长看的指导手册！`),
		schema.MessagesPlaceholder("chat_history", true),
		schema.UserMessage("{message}"),
	)
}

// GenerateLanguageAnswer 生成语言回答
func (n *LanguageAgentNode) GenerateLanguageAnswer(ctx context.Context, message string, objectName string, objectCategory string, userAge int, chatHistory []*schema.Message, recommendedTools []string) (*types.DomainAgentResponse, error) {
	n.logger.Infow("执行Language Agent回答生成",
		logx.Field("message", message),
		logx.Field("objectName", objectName),
		logx.Field("userAge", userAge),
		logx.Field("recommendedTools", recommendedTools),
		logx.Field("useRealModel", n.initialized),
	)

	if n.initialized && n.chatModel != nil {
		return n.executeReal(ctx, message, objectName, objectCategory, userAge, chatHistory, recommendedTools)
	}

	return n.executeMock(message, objectName, userAge)
}

// executeMock Mock实现
func (n *LanguageAgentNode) executeMock(message string, objectName string, userAge int) (*types.DomainAgentResponse, error) {
	content := "用英语说" + objectName + "是 \"" + objectName + "\"。你可以说：This is " + objectName + "."
	if userAge <= 6 {
		content = "这个叫" + objectName + "，你可以说：这是" + objectName + "。"
	}

	return &types.DomainAgentResponse{
		DomainType:  "Language",
		Content:     content,
		ToolsUsed:   []string{},
		ToolResults: make(map[string]interface{}),
	}, nil
}

// executeReal 真实eino实现（支持工具调用）
func (n *LanguageAgentNode) executeReal(ctx context.Context, message string, objectName string, objectCategory string, userAge int, chatHistory []*schema.Message, recommendedTools []string) (*types.DomainAgentResponse, error) {
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
	
	// 添加用户消息
	messages = append(messages, schema.UserMessage(message))

	// 如果有关键工具推荐，动态注册（补充到已注册的工具）
	if len(recommendedTools) > 0 && n.toolRegistry != nil {
		// 获取已注册的工具
		existingTools := n.toolRegistry.GetToolsForAgent("Language")

		// 合并推荐的工具（去重）
		toolMap := make(map[string]bool)
		allTools := make([]tools.Tool, 0, len(existingTools))
		for _, t := range existingTools {
			allTools = append(allTools, t)
			toolMap[t.Name()] = true
		}
		for _, name := range recommendedTools {
			if !toolMap[name] {
				if tool, ok := n.toolRegistry.GetTool(name); ok {
					allTools = append(allTools, tool)
					toolMap[name] = true
				}
			}
		}

		// 转换为工具信息并重新绑定
		if len(allTools) > 0 {
			toolInfos, err := tools.ConvertToEinoTools(allTools, ctx)
			if err == nil && len(toolInfos) > 0 {
				if err := n.chatModel.BindTools(toolInfos); err != nil {
					n.logger.Errorw("动态绑定工具失败", logx.Field("error", err))
				} else {
					n.logger.Infow("🔄 动态注册推荐工具",
						logx.Field("recommended_tools", recommendedTools),
						logx.Field("total_tools", len(toolInfos)),
					)
				}
			}
		}
	}

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

	// 使用工具调用链处理工具调用
	toolChain := NewToolChain(n.toolRegistry, n.logger)
	finalMessages, toolsUsed, toolResults, err := toolChain.ExecuteToolChain(ctx, cleanMessages, n.chatModel, recommendedTools)
	if err != nil {
		n.logger.Errorw("工具调用链执行失败", logx.Field("error", err))
		// 降级：直接调用ChatModel
		result, err := n.chatModel.Generate(ctx, cleanMessages)
		if err != nil {
			n.logger.Errorw("ChatModel调用失败，降级到Mock模式",
				logx.Field("error", err),
				logx.Field("message", message),
				logx.Field("objectName", objectName),
			)
			return n.executeMock(message, objectName, userAge)
		}
		return &types.DomainAgentResponse{
			DomainType:  "Language",
			Content:     result.Content,
			ToolsUsed:   []string{},
			ToolResults: make(map[string]interface{}),
		}, nil
	}

	// 获取最终结果（最后一条消息）
	var result *schema.Message
	if len(finalMessages) > 0 {
		result = finalMessages[len(finalMessages)-1]
	} else {
		// 如果没有结果，降级处理
		return n.executeMock(message, objectName, userAge)
	}

	return &types.DomainAgentResponse{
		DomainType:  "Language",
		Content:     result.Content,
		ToolsUsed:   toolsUsed,
		ToolResults: toolResults,
	}, nil
}

// buildSystemMessageWithTools 根据推荐的工具构建SystemMessage
func (n *LanguageAgentNode) buildSystemMessageWithTools(recommendedTools []string) string {
	baseMessage := `你是 Language Agent，一个直接和孩子对话的AI伙伴，帮助孩子用语言表达自己的想法。

重要规则：
- 直接回答孩子的问题，就像朋友聊天一样
- 不要出现"跟小朋友可以这样聊"、"你可以说"等指导性语言
- 不要出现"你:"这样的对话示例格式
- 用"我"或直接称呼"你"（孩子）来对话
- 让孩子"说得出口"，不讲语法规则
- 用孩子日常语言，包含可模仿的句子
- 让孩子感受到表达的乐趣`

	// 如果有推荐的工具，添加到SystemMessage中
	if len(recommendedTools) > 0 {
		toolDescriptions := n.getToolDescriptions(recommendedTools)
		if toolDescriptions != "" {
			baseMessage += "\n\n你可以调用的工具：\n" + toolDescriptions
			baseMessage += "\n\n重要：当问题需要使用工具时，你必须调用相应的工具来获取信息，然后再回答。"
			baseMessage += "\n\n工具调用规则："
			baseMessage += "\n- 如果问时间相关的问题（几点了、现在几点、什么时候、现在几时），必须调用get_current_time工具"
			baseMessage += "\n- 如果问单词的意思或发音，必须调用simple_dictionary或pronunciation_hint工具"
			baseMessage += "\n- 调用工具时，使用function calling格式，不要直接回答"
			baseMessage += "\n\n如果工具调用失败，不依赖工具也能生成基本回答。"
		}
	} else {
		// 默认工具列表
		baseMessage += "\n\n你可以调用：\n- simple_dictionary: 查找单词\n- pronunciation_hint: 发音提示\n\n如果工具调用失败，不依赖工具也能生成基本回答。"
	}

	baseMessage += "\n\n记住：你是直接和孩子对话的AI伙伴，不是给家长看的指导手册！"
	return baseMessage
}

// getToolDescriptions 获取工具描述列表
func (n *LanguageAgentNode) getToolDescriptions(toolNames []string) string {
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

