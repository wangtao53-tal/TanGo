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

// generateSystemPrompt 根据用户年龄生成系统prompt
func (n *ConversationNode) generateSystemPrompt(userAge int, objectName string, objectCategory string) string {
	var difficulty string
	var contentStyle string
	var interactionStyle string
	var knowledgeDepth string

	// 根据年龄段确定难度、风格和交互方式
	// 3-6岁：幼儿阶段
	if userAge <= 6 {
		difficulty = "最简单易懂，使用儿童语言，避免专业术语"
		contentStyle = "生动有趣，多用比喻、拟人和故事，像讲故事一样"
		interactionStyle = "多用提问和互动，如'你见过吗？'、'你觉得呢？'，鼓励孩子观察和表达"
		knowledgeDepth = "基础认知，重点培养观察力和好奇心，内容要贴近日常生活"
	} else if userAge <= 12 {
		// 7-12岁：小学阶段
		difficulty = "简单易懂，使用日常语言，可以适当使用基础科学术语"
		contentStyle = "结合生活实际，激发探索兴趣，可以加入简单的科学原理"
		interactionStyle = "引导式提问，如'为什么？'、'怎么样？'，培养思考习惯"
		knowledgeDepth = "中等深度，结合课本知识但以拓展为主，培养科学思维和探索精神"
	} else {
		// 13-18岁：中学阶段
		difficulty = "准确专业，可以使用科学术语，但要深入浅出地解释"
		contentStyle = "深入浅出，培养科学思维，可以涉及跨学科知识和前沿科学"
		interactionStyle = "引导深度思考，培养批判性思维，可以讨论科学问题和实际应用"
		knowledgeDepth = "较高深度，可以涉及学科知识、科学原理和实际应用，培养科学素养"
	}

	prompt := fmt.Sprintf(`你是一个面向%d岁学生的AI助手，专门帮助学生学习课外知识。

要求：
1. 语言风格：%s
2. 内容风格：%s
3. 交互方式：%s
4. 知识深度：%s
5. 结合%s相关的科学知识、古诗词和英语表达
6. 拓展素质教育，培养探索精神和学习兴趣
7. 内容贴合K12课程，但以课外拓展为主，避免直接讲解课本内容`, 
		userAge, difficulty, contentStyle, interactionStyle, knowledgeDepth, objectName)

	// 如果有识别对象信息，添加到prompt
	if objectName != "" {
		prompt += fmt.Sprintf("\n8. 当前讨论的对象是：%s（%s），可以围绕这个对象展开相关知识的拓展", objectName, objectCategory)
	}

	return prompt
}

// StreamConversation 流式对话，返回流式读取器，支持多模态输入
func (n *ConversationNode) StreamConversation(
	ctx context.Context,
	message string,
	contextMessages []*schema.Message,
	userAge int,
	objectName string,
	objectCategory string,
	imageURL string, // 新增：图片URL参数，支持多模态输入
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

	// 构建用户消息（支持多模态）
	var userMsg *schema.Message
	if imageURL != "" {
		// 多模态消息（图片+文本，如果文本不为空）
		parts := []schema.MessageInputPart{
			{
				Type: schema.ChatMessagePartTypeImageURL,
				Image: &schema.MessageInputImage{
					MessagePartCommon: schema.MessagePartCommon{
						URL: &imageURL,
					},
					Detail: schema.ImageURLDetailAuto,
				},
			},
		}
		// 只有当文本不为空时才添加文本部分
		if message != "" {
			parts = append(parts, schema.MessageInputPart{
				Type: schema.ChatMessagePartTypeText,
				Text: message,
			})
		}
		userMsg = &schema.Message{
			Role:                schema.User,
			UserInputMultiContent: parts,
		}
		n.logger.Infow("构建多模态消息",
			logx.Field("hasImage", true),
			logx.Field("imageURL", imageURL),
			logx.Field("textLength", len(message)),
			logx.Field("hasText", message != ""),
		)
	} else {
		// 文本消息
		userMsg = schema.UserMessage(message)
	}
	messages = append(messages, userMsg)

	n.logger.Infow("开始流式对话",
		logx.Field("userAge", userAge),
		logx.Field("objectName", objectName),
		logx.Field("contextRounds", len(contextMessages)/2),
		logx.Field("messageLength", len(message)),
		logx.Field("hasImage", imageURL != ""),
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

