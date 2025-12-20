package nodes

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
	"github.com/tango/explore/internal/config"
	"github.com/zeromicro/go-zero/core/logx"
)

// ConversationNode 对话节点
type ConversationNode struct {
	ctx         context.Context
	config      config.AIConfig
	logger      logx.Logger
	chatModel   model.ChatModel     // eino ChatModel 实例
	template    prompt.ChatTemplate // 对话模板
	initialized bool
}

// NewConversationNode 创建对话节点
func NewConversationNode(ctx context.Context, cfg config.AIConfig, logger logx.Logger) (*ConversationNode, error) {
	node := &ConversationNode{
		ctx:    ctx,
		config: cfg,
		logger: logger,
	}

	// 如果配置了 eino 相关参数，初始化 ChatModel
	hasEinoBaseURL := cfg.EinoBaseURL != ""
	hasAppID := cfg.AppID != ""
	hasAppKey := cfg.AppKey != ""

	if hasEinoBaseURL && hasAppID && hasAppKey {
		logger.Infow("检测到eino配置，尝试初始化对话ChatModel",
			logx.Field("einoBaseURL", cfg.EinoBaseURL),
			logx.Field("appID", hasAppID),
			logx.Field("hasAppKey", hasAppKey),
		)
		if err := node.initChatModel(ctx); err != nil {
			logger.Errorw("初始化对话ChatModel失败，将使用Mock模式",
				logx.Field("error", err),
			)
		} else {
			node.initialized = true
			logger.Info("✅ 对话节点已初始化ChatModel，将使用真实模型")
		}
	} else {
		logger.Errorw("未完整配置eino参数，对话节点将使用Mock模式",
			logx.Field("hasEinoBaseURL", hasEinoBaseURL),
			logx.Field("hasAppID", hasAppID),
			logx.Field("hasAppKey", hasAppKey),
		)
	}

	// 创建对话模板
	node.initTemplate()

	return node, nil
}

// initChatModel 初始化 ChatModel
func (n *ConversationNode) initChatModel(ctx context.Context) error {
	modelName := n.config.TextGenerationModel
	if modelName == "" {
		modelName = config.DefaultTextGenerationModel
	}

	cfg := &ark.ChatModelConfig{
		Model: modelName,
	}

	if n.config.EinoBaseURL != "" {
		cfg.BaseURL = n.config.EinoBaseURL
	}

	// 认证：使用 Bearer Token 格式 ${TAL_MLOPS_APP_ID}:${TAL_MLOPS_APP_KEY}
	if n.config.AppID != "" && n.config.AppKey != "" {
		cfg.APIKey = n.config.AppID + ":" + n.config.AppKey
	} else if n.config.AppKey != "" {
		cfg.APIKey = n.config.AppKey
	} else if n.config.AppID != "" {
		cfg.APIKey = n.config.AppID
	} else {
		return nil // 返回 nil，使用 Mock 模式
	}

	chatModel, err := ark.NewChatModel(ctx, cfg)
	if err != nil {
		return err
	}

	n.chatModel = chatModel
	return nil
}

// initTemplate 初始化对话模板
func (n *ConversationNode) initTemplate() {
	// 对话模板支持动态参数注入
	n.template = prompt.FromMessages(schema.FString,
		schema.MessagesPlaceholder("chat_history", true),
		schema.UserMessage("{message}"),
	)
}

// generateSystemPrompt 根据用户年级生成系统prompt
func (n *ConversationNode) generateSystemPrompt(userAge int, objectName string, objectCategory string) string {
	var difficulty string
	var contentStyle string

	// 根据年龄确定难度和风格
	if userAge <= 6 {
		difficulty = "简单易懂，使用儿童语言"
		contentStyle = "生动有趣，多用比喻和故事"
	} else if userAge <= 12 {
		difficulty = "中等难度，使用日常语言"
		contentStyle = "结合生活实际，激发探索兴趣"
	} else {
		difficulty = "较高难度，可以使用专业术语"
		contentStyle = "深入浅出，培养科学思维"
	}

	prompt := fmt.Sprintf(`你是一个面向%d岁学生的AI助手，专门帮助学生学习课外知识。
要求：
1. 使用%s的语言风格
2. 内容%s
3. 结合%s相关的科学知识、古诗词和英语表达
4. 拓展素质教育，培养探索精神
5. 内容贴合K12课程，但以课外拓展为主`, userAge, difficulty, contentStyle, objectName)

	// 如果有识别对象信息，添加到prompt
	if objectName != "" {
		prompt += fmt.Sprintf("\n6. 当前讨论的对象是：%s（%s）", objectName, objectCategory)
	}

	return prompt
}

// StreamConversation 流式对话，返回流式读取器
func (n *ConversationNode) StreamConversation(
	ctx context.Context,
	message string,
	contextMessages []*schema.Message,
	userAge int,
	objectName string,
	objectCategory string,
) (*schema.StreamReader[*schema.Message], error) {
	if !n.initialized || n.chatModel == nil {
		return nil, fmt.Errorf("ChatModel未初始化，无法进行流式对话")
	}

	// 根据用户年级生成系统prompt
	systemPrompt := n.generateSystemPrompt(userAge, objectName, objectCategory)

	// 构建消息列表
	messages := []*schema.Message{
		schema.SystemMessage(systemPrompt),
	}

	// 添加上下文消息（最多20轮）
	if len(contextMessages) > 0 {
		messages = append(messages, contextMessages...)
	}

	// 添加当前用户消息
	messages = append(messages, schema.UserMessage(message))

	n.logger.Infow("开始流式对话",
		logx.Field("userAge", userAge),
		logx.Field("objectName", objectName),
		logx.Field("contextRounds", len(contextMessages)/2),
		logx.Field("messageLength", len(message)),
	)

	// 调用Eino ChatModel的Stream接口
	streamReader, err := n.chatModel.Stream(ctx, messages)
	if err != nil {
		n.logger.Errorw("调用Eino Stream接口失败",
			logx.Field("error", err),
		)
		return nil, fmt.Errorf("调用AI模型失败: %w", err)
	}

	return streamReader, nil
}

// GenerateText 非流式文本生成（兼容性接口）
func (n *ConversationNode) GenerateText(
	ctx context.Context,
	message string,
	contextMessages []*schema.Message,
	userAge int,
	objectName string,
	objectCategory string,
) (string, error) {
	if !n.initialized || n.chatModel == nil {
		// Mock响应
		return fmt.Sprintf("这是一个Mock响应。待接入真实AI模型后，将根据您的问题和识别结果（%s）生成相应的回答。", objectName), nil
	}

	// 根据用户年级生成系统prompt
	systemPrompt := n.generateSystemPrompt(userAge, objectName, objectCategory)

	// 构建消息列表
	messages := []*schema.Message{
		schema.SystemMessage(systemPrompt),
	}

	// 添加上下文消息
	if len(contextMessages) > 0 {
		messages = append(messages, contextMessages...)
	}

	// 添加当前用户消息
	messages = append(messages, schema.UserMessage(message))

	// 调用Eino ChatModel的Generate接口
	result, err := n.chatModel.Generate(ctx, messages)
	if err != nil {
		n.logger.Errorw("调用Eino Generate接口失败",
			logx.Field("error", err),
		)
		return "", fmt.Errorf("调用AI模型失败: %w", err)
	}

	// 提取文本内容（Message.Content 是 string 类型）
	if result != nil && result.Content != "" {
		return result.Content, nil
	}

	return "", fmt.Errorf("无法从模型响应中提取文本内容")
}

// MockStreamConversation Mock流式对话（用于测试或降级）
func (n *ConversationNode) MockStreamConversation(message string) []string {
	text := fmt.Sprintf("这是一个Mock流式响应。您的问题是：%s。待接入真实AI模型后，将实现真实的流式文本生成。", message)
	words := []rune(text)
	result := make([]string, 0, len(words))
	for _, word := range words {
		result = append(result, string(word))
	}
	return result
}

