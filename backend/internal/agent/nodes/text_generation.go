package nodes

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
	"github.com/tango/explore/internal/config"
	"github.com/zeromicro/go-zero/core/logx"
)

// TextGenerationNode 文本生成节点
type TextGenerationNode struct {
	ctx             context.Context
	config          config.AIConfig
	logger          logx.Logger
	chatModel       model.ChatModel     // eino ChatModel 实例
	scienceTemplate prompt.ChatTemplate // 科学认知卡模板
	poetryTemplate  prompt.ChatTemplate // 古诗词卡模板
	englishTemplate prompt.ChatTemplate // 英语表达卡模板
	textTemplate    prompt.ChatTemplate // 文本回答模板
	initialized     bool
}

// NewTextGenerationNode 创建文本生成节点
func NewTextGenerationNode(ctx context.Context, cfg config.AIConfig, logger logx.Logger) (*TextGenerationNode, error) {
	node := &TextGenerationNode{
		ctx:    ctx,
		config: cfg,
		logger: logger,
	}

	// 如果配置了 eino 相关参数，初始化 ChatModel
	hasEinoBaseURL := cfg.EinoBaseURL != ""
	hasAppID := cfg.AppID != ""
	hasAppKey := cfg.AppKey != ""

	if hasEinoBaseURL && hasAppID && hasAppKey {
		logger.Infow("检测到eino配置，尝试初始化ChatModel",
			logx.Field("einoBaseURL", cfg.EinoBaseURL),
			logx.Field("appID", hasAppID),
			logx.Field("hasAppKey", hasAppKey),
		)
		if err := node.initChatModel(ctx); err != nil {
			logger.Errorw("初始化ChatModel失败，将使用Mock模式",
				logx.Field("error", err),
				logx.Field("errorDetail", err.Error()),
			)
			// 即使初始化失败，也继续创建节点，但会使用Mock模式
		} else {
			node.initialized = true
			logger.Info("✅ 文本生成节点已初始化ChatModel，将使用真实模型")
		}
	} else {
		logger.Errorw("未完整配置eino参数，文本生成节点将使用Mock模式",
			logx.Field("hasEinoBaseURL", hasEinoBaseURL),
			logx.Field("hasAppID", hasAppID),
			logx.Field("hasAppKey", hasAppKey),
		)
		logger.Info("提示：需要同时配置 EINO_BASE_URL、TAL_MLOPS_APP_ID、TAL_MLOPS_APP_KEY 才能使用真实模型")
	}

	// 检查配置中的UseAIModel设置
	// 如果UseAIModel为true但ChatModel未初始化，记录警告
	if cfg.UseAIModel && !node.initialized {
		logger.Errorw("配置要求使用AI模型，但ChatModel未初始化，将使用Mock模式",
			logx.Field("useAIModel", cfg.UseAIModel),
			logx.Field("initialized", node.initialized),
			logx.Field("chatModelNil", node.chatModel == nil),
		)
	}

	// 创建所有模板
	node.initTemplates()

	return node, nil
}

// initChatModel 初始化 ChatModel（使用随机选择的模型）
func (n *TextGenerationNode) initChatModel(ctx context.Context) error {
	// 从配置中随机选择一个文本生成模型
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

	// 认证：使用 Bearer Token 格式 ${TAL_MLOPS_APP_ID}:${TAL_MLOPS_APP_KEY}
	if n.config.AppID != "" && n.config.AppKey != "" {
		cfg.APIKey = n.config.AppID + ":" + n.config.AppKey
	} else if n.config.AppKey != "" {
		cfg.APIKey = n.config.AppKey
	} else if n.config.AppID != "" {
		cfg.APIKey = n.config.AppID
	} else {
		// 如果缺少认证信息，返回错误而不是nil，让调用方知道初始化失败
		return fmt.Errorf("缺少认证信息：需要配置 AppID 或 AppKey")
	}

	chatModel, err := ark.NewChatModel(ctx, cfg)
	if err != nil {
		n.logger.Errorw("创建ChatModel失败",
			logx.Field("error", err),
			logx.Field("modelName", modelName),
			logx.Field("baseURL", cfg.BaseURL),
		)
		return fmt.Errorf("创建ChatModel失败: %w", err)
	}

	n.chatModel = chatModel
	n.logger.Infow("ChatModel初始化成功",
		logx.Field("modelName", modelName),
		logx.Field("baseURL", cfg.BaseURL),
	)
	return nil
}

// selectRandomModel 从模型列表中随机选择一个模型
func (n *TextGenerationNode) selectRandomModel(models []string) string {
	if len(models) == 0 {
		return ""
	}
	if len(models) == 1 {
		return models[0]
	}
	// 使用当前时间作为随机种子
	rand.Seed(time.Now().UnixNano())
	return models[rand.Intn(len(models))]
}

// getAgePrompt 根据年龄生成对应的prompt要求
func (n *TextGenerationNode) getAgePrompt(age int, cardType string) string {
	var agePrompt string

	// 根据年龄段划分：3-6岁（幼儿）、7-12岁（小学）、13-18岁（中学）
	if age <= 6 {
		// 幼儿阶段（3-6岁）
		switch cardType {
		case "science":
			agePrompt = `要求：
1. 用最简单、最生动的语言解释{objectName}的科学知识，避免专业术语
2. 使用比喻和拟人手法，让内容像故事一样有趣
3. 提供2-3个简单有趣的事实，每个事实不超过一句话，可以适当使用emoji（如 🌟 ✨ 💡 🔍 等）
4. 添加一个趣味知识，用"你知道吗？"开头，可以适当使用emoji（如 🎉 🌈 ⭐ 等）
5. 内容要符合3-6岁孩子的认知水平，使用日常词汇
6. 可以加入互动元素，如"你见过吗？"、"你觉得呢？"等
7. 适当使用emoji让内容更生动，但不要过多，保持可读性`
		case "poetry":
			agePrompt = `要求：
1. 找到与{objectName}相关的古诗词，优先选择简短、朗朗上口的诗句
2. 标注诗词来源（作者和诗名）
3. 用最简单、最形象的语言解释诗词含义，多用比喻，可以适当使用emoji（如 📜 ✨ 🌸 🌙 等）
4. 提供简单的文化背景说明，不超过两句话，可以适当使用emoji（如 🏛️ 📚 等）
5. 解释要符合3-6岁孩子的理解能力，避免复杂概念
6. 适当使用emoji让内容更生动，但不要过多，保持可读性`
		case "english":
			agePrompt = `要求：
1. 提供{objectName}的英语关键词（3-4个），选择最简单、最常用的单词
2. 提供2-3个适合3-6岁孩子的英语表达句子，句子要简短（3-5个单词），可以适当使用emoji（如 🌟 💬 🎯 等）
3. 提供简单的发音指导，用中文拼音或音标标注，可以适当使用emoji（如 🔊 📝 等）
4. 可以加入简单的英语儿歌或韵律，帮助记忆
5. 适当使用emoji让内容更生动，但不要过多，保持可读性`
		}
	} else if age <= 12 {
		// 小学阶段（7-12岁）
		switch cardType {
		case "science":
			agePrompt = `要求：
1. 用简单易懂的语言解释{objectName}的科学知识，可以适当使用基础科学术语
2. 结合生活实际，让孩子能够联系到日常经验
3. 提供2-3个有趣的事实，每个事实可以包含简单的科学原理，可以适当使用emoji（如 🌟 ✨ 💡 🔍 等）
4. 添加一个趣味知识，可以涉及科学小实验或观察方法，可以适当使用emoji（如 🎉 🌈 ⭐ 🔬 等）
5. 内容要符合7-12岁孩子的认知水平，激发探索兴趣
6. 可以加入"为什么"、"怎么样"等引导性问题
7. 适当使用emoji让内容更生动，但不要过多，保持可读性`
		case "poetry":
			agePrompt = `要求：
1. 找到与{objectName}相关的古诗词（优先选择经典名句）
2. 标注诗词来源（作者和诗名）
3. 用7-12岁孩子能理解的语言解释诗词含义，可以适当讲解修辞手法，可以适当使用emoji（如 📜 ✨ 🌸 🌙 等）
4. 提供文化背景说明，包括历史背景和诗人创作意图，可以适当使用emoji（如 🏛️ 📚 🎨 等）
5. 可以引导孩子思考诗词中的情感和意境
6. 适当使用emoji让内容更生动，但不要过多，保持可读性`
		case "english":
			agePrompt = `要求：
1. 提供{objectName}的英语关键词（3-5个），包括基础词汇和相关表达
2. 提供2-3个适合7-12岁孩子的英语表达句子，句子可以稍长（5-8个单词），可以适当使用emoji（如 🌟 💬 🎯 等）
3. 提供发音指导，包括音标和发音技巧，可以适当使用emoji（如 🔊 📝 等）
4. 可以加入简单的语法点或常用搭配，帮助扩展词汇
5. 适当使用emoji让内容更生动，但不要过多，保持可读性`
		}
	} else {
		// 中学阶段（13-18岁）
		switch cardType {
		case "science":
			agePrompt = `要求：
1. 用准确、专业的语言解释{objectName}的科学知识，可以使用科学术语
2. 深入讲解科学原理，可以涉及物理、化学、生物等学科知识
3. 提供2-3个有深度的事实，每个事实可以包含科学原理和实际应用，可以适当使用emoji（如 🌟 ✨ 💡 🔍 等）
4. 添加一个趣味知识，可以涉及前沿科学或跨学科知识，可以适当使用emoji（如 🎉 🌈 ⭐ 🔬 等）
5. 内容要符合13-18岁学生的认知水平，培养科学思维
6. 可以引导思考科学问题，培养批判性思维
7. 适当使用emoji让内容更生动，但不要过多，保持可读性`
		case "poetry":
			agePrompt = `要求：
1. 找到与{objectName}相关的古诗词（优先选择经典名句，可以包含较长的诗句）
2. 标注诗词来源（作者和诗名），可以介绍诗人的生平和创作背景
3. 深入解释诗词含义，分析修辞手法、意象和艺术特色，可以适当使用emoji（如 📜 ✨ 🌸 🌙 等）
4. 提供详细的文化背景说明，包括历史背景、文学流派和艺术价值，可以适当使用emoji（如 🏛️ 📚 🎨 等）
5. 可以引导分析诗词的深层含义和思想情感，培养文学鉴赏能力
6. 适当使用emoji让内容更生动，但不要过多，保持可读性`
		case "english":
			agePrompt = `要求：
1. 提供{objectName}的英语关键词（4-6个），包括高级词汇和相关表达
2. 提供2-3个适合13-18岁学生的英语表达句子，句子可以更复杂（8-12个单词），可以适当使用emoji（如 🌟 💬 🎯 等）
3. 提供详细的发音指导，包括音标、重音和语调，可以适当使用emoji（如 🔊 📝 等）
4. 可以加入语法点、固定搭配和高级表达，帮助提升英语水平
5. 可以介绍相关的英语文化背景或使用场景
6. 适当使用emoji让内容更生动，但不要过多，保持可读性`
		}
	}

	return agePrompt
}

// initTemplates 初始化所有消息模板
func (n *TextGenerationNode) initTemplates() {
	// 科学认知卡模板（使用动态prompt，根据年龄调整）
	n.scienceTemplate = prompt.FromMessages(schema.FString,
		schema.SystemMessage(`你是一个K12教育内容生成助手，专门为{age}岁的孩子生成科学认知卡片内容。

{agePrompt}

请返回JSON格式，包含以下字段：
- name: 对象名称（字符串）
- explanation: 科学解释（字符串，适当使用emoji如 🌟 ✨ 💡 🔍 等让内容更生动）
- facts: 有趣的事实列表（字符串数组，2-3个，每个事实可以适当使用emoji）
- funFact: 趣味知识（字符串，可以适当使用emoji如 🎉 🌈 ⭐ 等）

注意：emoji要适量使用，不要过多，保持内容的可读性。`),
		schema.UserMessage("请为{objectName}生成科学认知卡内容，适合{age}岁孩子。"),
	)

	// 古诗词卡模板（使用动态prompt，根据年龄调整）
	n.poetryTemplate = prompt.FromMessages(schema.FString,
		schema.SystemMessage(`你是一个古诗词专家，专门为K12教育生成古诗词卡片内容。

{agePrompt}

请返回JSON格式，包含以下字段：
- poem: 古诗词内容（字符串）
- poemSource: 作者和诗名（字符串，格式：作者 - 诗名）
- explanation: 诗词解释（字符串，适当使用emoji如 📜 ✨ 🌸 🌙 等让内容更生动）
- context: 文化背景（字符串，可以适当使用emoji如 🏛️ 📚 🎨 等）

注意：emoji要适量使用，不要过多，保持内容的可读性。`),
		schema.UserMessage("请为{objectName}生成古诗词卡片内容，适合{age}岁孩子。"),
	)

	// 英语表达卡模板（使用动态prompt，根据年龄调整）
	n.englishTemplate = prompt.FromMessages(schema.FString,
		schema.SystemMessage(`你是一个英语教学专家，专门为K12教育生成英语表达卡片内容。

{agePrompt}

请返回JSON格式，包含以下字段：
- keywords: 英语关键词列表（字符串数组，3-5个）
- expressions: 英语表达句子列表（字符串数组，2-3个，可以适当使用emoji如 🌟 💬 🎯 等）
- pronunciation: 发音指导（字符串，可以适当使用emoji如 🔊 📝 等）

注意：emoji要适量使用，不要过多，保持内容的可读性。`),
		schema.UserMessage("请为{objectName}生成英语表达卡片内容，适合{age}岁孩子。"),
	)

	// 文本回答模板
	n.textTemplate = prompt.FromMessages(schema.FString,
		schema.SystemMessage("你是一个友好的K12教育助手，用简单易懂的语言回答孩子的问题。适当使用emoji表情符号（如 🌟 ✨ 💡 🔍 📚 🎨 🌈 🦋 🌸 ⭐ 等）让回答更生动有趣，适合小朋友阅读。注意：emoji要适量，不要过多，避免影响阅读体验。"),
		schema.MessagesPlaceholder("chat_history", true),
		schema.UserMessage("{message}"),
	)
}

// GenerateText 生成文本回答
func (n *TextGenerationNode) GenerateText(data *GraphData, context []interface{}) (string, error) {
	n.logger.Infow("执行文本生成",
		logx.Field("message", data.Text),
		logx.Field("contextLength", len(context)),
		logx.Field("useRealModel", n.initialized),
	)

	if n.initialized {
		// 每次调用时重新初始化 ChatModel，使用随机选择的模型
		if err := n.initChatModel(n.ctx); err != nil {
			n.logger.Errorw("重新初始化ChatModel失败，使用已初始化的模型",
				logx.Field("error", err),
			)
			// 如果重新初始化失败，继续使用已初始化的模型
		}
		if n.chatModel != nil {
			return n.generateTextReal(data, context)
		}
	}

	return n.generateTextMock(data, context)
}

// GenerateScienceCard 生成科学认知卡内容
func (n *TextGenerationNode) GenerateScienceCard(ctx context.Context, data *GraphData) (map[string]interface{}, error) {
	n.logger.Infow("生成科学认知卡",
		logx.Field("objectName", data.ObjectName),
		logx.Field("age", data.Age),
		logx.Field("useRealModel", n.initialized),
		logx.Field("chatModelNil", n.chatModel == nil),
	)

	if n.initialized {
		// 每次调用时重新初始化 ChatModel，使用随机选择的模型
		if err := n.initChatModel(ctx); err != nil {
			n.logger.Errorw("重新初始化ChatModel失败，使用已初始化的模型",
				logx.Field("error", err),
			)
			// 如果重新初始化失败，继续使用已初始化的模型
		}
		if n.chatModel != nil {
			n.logger.Infow("尝试使用真实模型生成科学认知卡",
				logx.Field("objectName", data.ObjectName),
			)
			card, err := n.generateScienceCardReal(ctx, data)
			if err != nil {
				n.logger.Errorw("真实模型生成科学认知卡失败，返回错误",
					logx.Field("objectName", data.ObjectName),
					logx.Field("error", err),
					logx.Field("errorDetail", err.Error()),
				)
				return nil, err
			}
			n.logger.Infow("真实模型生成科学认知卡成功",
				logx.Field("objectName", data.ObjectName),
			)
			return card, nil
		}
	}

	n.logger.Errorw("使用Mock模式生成科学认知卡（ChatModel未初始化）",
		logx.Field("objectName", data.ObjectName),
		logx.Field("initialized", n.initialized),
		logx.Field("chatModelNil", n.chatModel == nil),
	)
	return n.generateScienceCardMock(data)
}

// GeneratePoetryCard 生成古诗词/人文卡内容
func (n *TextGenerationNode) GeneratePoetryCard(ctx context.Context, data *GraphData) (map[string]interface{}, error) {
	n.logger.Infow("生成古诗词卡",
		logx.Field("objectName", data.ObjectName),
		logx.Field("age", data.Age),
		logx.Field("useRealModel", n.initialized),
		logx.Field("chatModelNil", n.chatModel == nil),
	)

	if n.initialized {
		// 每次调用时重新初始化 ChatModel，使用随机选择的模型
		if err := n.initChatModel(ctx); err != nil {
			n.logger.Errorw("重新初始化ChatModel失败，使用已初始化的模型",
				logx.Field("error", err),
			)
			// 如果重新初始化失败，继续使用已初始化的模型
		}
		if n.chatModel != nil {
			n.logger.Infow("尝试使用真实模型生成古诗词卡",
				logx.Field("objectName", data.ObjectName),
			)
			card, err := n.generatePoetryCardReal(ctx, data)
			if err != nil {
				n.logger.Errorw("真实模型生成古诗词卡失败，返回错误",
					logx.Field("objectName", data.ObjectName),
					logx.Field("error", err),
					logx.Field("errorDetail", err.Error()),
				)
				return nil, err
			}
			n.logger.Infow("真实模型生成古诗词卡成功",
				logx.Field("objectName", data.ObjectName),
			)
			return card, nil
		}
	}

	n.logger.Errorw("使用Mock模式生成古诗词卡（ChatModel未初始化）",
		logx.Field("objectName", data.ObjectName),
		logx.Field("initialized", n.initialized),
		logx.Field("chatModelNil", n.chatModel == nil),
	)
	return n.generatePoetryCardMock(data)
}

// GenerateEnglishCard 生成英语表达卡内容
func (n *TextGenerationNode) GenerateEnglishCard(ctx context.Context, data *GraphData) (map[string]interface{}, error) {
	n.logger.Infow("生成英语表达卡",
		logx.Field("objectName", data.ObjectName),
		logx.Field("age", data.Age),
		logx.Field("useRealModel", n.initialized),
		logx.Field("chatModelNil", n.chatModel == nil),
	)

	if n.initialized {
		// 每次调用时重新初始化 ChatModel，使用随机选择的模型
		if err := n.initChatModel(ctx); err != nil {
			n.logger.Errorw("重新初始化ChatModel失败，使用已初始化的模型",
				logx.Field("error", err),
			)
			// 如果重新初始化失败，继续使用已初始化的模型
		}
		if n.chatModel != nil {
			n.logger.Infow("尝试使用真实模型生成英语表达卡",
				logx.Field("objectName", data.ObjectName),
			)
			card, err := n.generateEnglishCardReal(ctx, data)
			if err != nil {
				n.logger.Errorw("真实模型生成英语表达卡失败，返回错误",
					logx.Field("objectName", data.ObjectName),
					logx.Field("error", err),
					logx.Field("errorDetail", err.Error()),
				)
				return nil, err
			}
			n.logger.Infow("真实模型生成英语表达卡成功",
				logx.Field("objectName", data.ObjectName),
			)
			return card, nil
		}
	}

	n.logger.Errorw("使用Mock模式生成英语表达卡（ChatModel未初始化）",
		logx.Field("objectName", data.ObjectName),
		logx.Field("initialized", n.initialized),
		logx.Field("chatModelNil", n.chatModel == nil),
	)
	return n.generateEnglishCardMock(data)
}

// generateTextMock Mock实现（待替换为真实eino调用）
func (n *TextGenerationNode) generateTextMock(data *GraphData, context []interface{}) (string, error) {
	// Mock文本响应
	text := "这是一个Mock文本响应 🌟。待接入真实AI模型后，将根据您的问题生成相应的回答 ✨。"
	n.logger.Info("文本生成完成（Mock）")
	return text, nil
}

// generateScienceCardMock Mock实现科学认知卡
func (n *TextGenerationNode) generateScienceCardMock(data *GraphData) (map[string]interface{}, error) {
	explanations := map[string]string{
		"银杏": "银杏是非常古老的植物，已经在地球上生存了2亿多年。",
		"苹果": "苹果是一种营养丰富的水果，含有多种维生素和矿物质。",
		"蝴蝶": "蝴蝶是美丽的昆虫，会经历从卵到幼虫、蛹、成虫的完全变态过程。",
		"书本": "书本是人类知识的载体，记录着历史和智慧。",
		"汽车": "汽车是现代重要的交通工具，使用发动机驱动。",
		"月亮": "月亮是地球的卫星，围绕地球旋转，影响潮汐。",
		"钢琴": "钢琴是一种键盘乐器，可以演奏丰富的音乐。",
		"太阳": "太阳是太阳系的中心，为地球提供光和热。",
	}

	explanation := explanations[data.ObjectName]
	if explanation == "" {
		explanation = data.ObjectName + "是一个有趣的对象 🌟，值得我们探索和学习 ✨。"
	}

	card := map[string]interface{}{
		"type":  "science",
		"title": data.ObjectName + "的科学知识",
		"content": map[string]interface{}{
			"name":        data.ObjectName,
			"explanation": explanation,
			"facts": []string{
				"关于" + data.ObjectName + "的有趣事实1 💡",
				"关于" + data.ObjectName + "的有趣事实2 🔍",
			},
			"funFact": "关于" + data.ObjectName + "的趣味知识 🎉！",
		},
	}

	n.logger.Info("科学认知卡生成完成（Mock）")
	return card, nil
}

// generatePoetryCardMock Mock实现古诗词卡
func (n *TextGenerationNode) generatePoetryCardMock(data *GraphData) (map[string]interface{}, error) {
	poems := map[string]string{
		"银杏": "满地翻黄银杏叶，忽惊天地告成功。",
		"苹果": "苹果红时秋已深，满园香气醉人心。",
		"蝴蝶": "穿花蛱蝶深深见，点水蜻蜓款款飞。",
		"月亮": "床前明月光，疑是地上霜。",
		"太阳": "日出江花红胜火，春来江水绿如蓝。",
	}

	poem := poems[data.ObjectName]
	hasPresetPoem := poem != ""
	if !hasPresetPoem {
		poem = "关于" + data.ObjectName + "的古诗词 📜，等待我们去发现 ✨。"
		n.logger.Infow("古诗词卡Mock：对象名不在预设列表中，使用默认值",
			logx.Field("objectName", data.ObjectName),
			logx.Field("defaultPoem", poem),
		)
	} else {
		n.logger.Infow("古诗词卡Mock：使用预设诗句",
			logx.Field("objectName", data.ObjectName),
			logx.Field("poem", poem),
		)
	}

	card := map[string]interface{}{
		"type":  "poetry",
		"title": "古人怎么看" + data.ObjectName,
		"content": map[string]interface{}{
			"poem":        poem,
			"poemSource":  "古诗词",
			"explanation": "这句诗描写了" + data.ObjectName + "的美丽景象 🌸，让我们感受到古人的智慧和情感 ✨。",
			"context":     "看到" + data.ObjectName + "，我们可以联想到相关的文化和历史 🏛️，丰富我们的认知 📚。",
		},
	}

	n.logger.Infow("古诗词卡生成完成（Mock）",
		logx.Field("objectName", data.ObjectName),
		logx.Field("hasPresetPoem", hasPresetPoem),
		logx.Field("title", card["title"]),
	)
	return card, nil
}

// generateEnglishCardMock Mock实现英语表达卡
func (n *TextGenerationNode) generateEnglishCardMock(data *GraphData) (map[string]interface{}, error) {
	keywords := map[string][]string{
		"银杏": {"ginkgo", "tree", "ancient"},
		"苹果": {"apple", "fruit", "red"},
		"蝴蝶": {"butterfly", "insect", "beautiful"},
		"书本": {"book", "knowledge", "reading"},
		"汽车": {"car", "vehicle", "transport"},
		"月亮": {"moon", "night", "round"},
		"钢琴": {"piano", "music", "instrument"},
		"太阳": {"sun", "bright", "energy"},
	}

	kw := keywords[data.ObjectName]
	if len(kw) == 0 {
		kw = []string{data.ObjectName, "object", "interesting"}
	}

	card := map[string]interface{}{
		"type":  "english",
		"title": "用英语说" + data.ObjectName,
		"content": map[string]interface{}{
			"keywords": kw,
			"expressions": []string{
				"This is " + kw[0] + ".",
				"I like " + kw[0] + ".",
			},
			"pronunciation": kw[0] + ": /pronunciation/",
		},
	}

	n.logger.Info("英语表达卡生成完成（Mock）")
	return card, nil
}

// generateTextReal 真实eino实现
func (n *TextGenerationNode) generateTextReal(data *GraphData, context []interface{}) (string, error) {
	// 转换上下文为 Message 格式
	chatHistory := make([]*schema.Message, 0)
	for _, ctxItem := range context {
		if msg, ok := ctxItem.(*schema.Message); ok {
			chatHistory = append(chatHistory, msg)
		}
	}

	// 使用模板生成消息
	messages, err := n.textTemplate.Format(n.ctx, map[string]any{
		"message":      data.Text,
		"chat_history": chatHistory,
	})
	if err != nil {
		n.logger.Errorw("模板格式化失败", logx.Field("error", err))
		return n.generateTextMock(data, context)
	}

	// 调用 ChatModel
	result, err := n.chatModel.Generate(n.ctx, messages)
	if err != nil {
		n.logger.Errorw("ChatModel调用失败", logx.Field("error", err))
		return n.generateTextMock(data, context)
	}

	return result.Content, nil
}

// parseJSONFromResponse 从模型响应中解析JSON，如果失败则返回错误
func (n *TextGenerationNode) parseJSONFromResponse(text string) (map[string]interface{}, error) {
	jsonStart := strings.Index(text, "{")
	jsonEnd := strings.LastIndex(text, "}")
	if jsonStart >= 0 && jsonEnd > jsonStart {
		jsonStr := text[jsonStart : jsonEnd+1]
		var cardContent map[string]interface{}
		if err := json.Unmarshal([]byte(jsonStr), &cardContent); err != nil {
			return nil, fmt.Errorf("解析JSON失败: %w, JSON字符串: %s", err, jsonStr)
		}
		return cardContent, nil
	}
	return nil, fmt.Errorf("未找到有效的JSON内容")
}

// generateCardWithRetry 生成卡片并支持JSON解析失败时的快速重试
// maxRetries: 最大重试次数（不包括首次调用）
func (n *TextGenerationNode) generateCardWithRetry(
	ctx context.Context,
	cardType string, // "science", "poetry", "english"
	data *GraphData,
	template prompt.ChatTemplate,
	maxRetries int,
) (map[string]interface{}, error) {
	var lastErr error
	var lastContent string

	// 根据卡片类型获取年龄prompt
	var agePrompt string
	switch cardType {
	case "science":
		agePrompt = n.getAgePrompt(data.Age, "science")
	case "poetry":
		agePrompt = n.getAgePrompt(data.Age, "poetry")
	case "english":
		agePrompt = n.getAgePrompt(data.Age, "english")
	default:
		return nil, fmt.Errorf("未知的卡片类型: %s", cardType)
	}

	// 重试逻辑：首次调用 + 最多maxRetries次重试
	for attempt := 0; attempt <= maxRetries; attempt++ {
		// 每次重试都重新构建消息，避免消息累积
		messages, err := template.Format(ctx, map[string]any{
			"objectName": data.ObjectName,
			"age":        strconv.Itoa(data.Age),
			"agePrompt":  agePrompt,
		})
		if err != nil {
			return nil, fmt.Errorf("模板格式化失败: %w", err)
		}

		if attempt > 0 {
			// 重试时，在用户消息中添加强调JSON格式的提示
			// 在最后一条用户消息后添加强调
			if len(messages) > 0 {
				lastMsg := messages[len(messages)-1]
				if lastMsg.Role == schema.User {
					// 创建新的强调消息
					emphasisMsg := &schema.Message{
						Role:    schema.User,
						Content: "重要：请严格按照JSON格式返回，不要包含任何其他文本说明。",
					}
					messages = append(messages, emphasisMsg)
				}
			}
			n.logger.Infow("JSON解析失败，快速重试模型调用",
				logx.Field("cardType", cardType),
				logx.Field("objectName", data.ObjectName),
				logx.Field("attempt", attempt),
				logx.Field("maxRetries", maxRetries),
				logx.Field("lastError", lastErr),
			)
			// 短暂延迟，避免立即重试
			time.Sleep(100 * time.Millisecond)
		}

		// 调用模型
		result, err := n.chatModel.Generate(ctx, messages)
		if err != nil {
			n.logger.Errorw("ChatModel调用失败",
				logx.Field("cardType", cardType),
				logx.Field("objectName", data.ObjectName),
				logx.Field("attempt", attempt),
				logx.Field("error", err),
			)
			// 如果是模型调用错误（非JSON解析错误），直接返回
			if attempt == 0 {
				return nil, fmt.Errorf("ChatModel调用失败: %w", err)
			}
			// 重试时的模型调用错误，继续重试
			lastErr = fmt.Errorf("ChatModel调用失败: %w", err)
			continue
		}

		// 尝试解析JSON
		cardContent, parseErr := n.parseJSONFromResponse(result.Content)
		if parseErr == nil {
			// 解析成功
			if attempt > 0 {
				n.logger.Infow("重试成功，JSON解析通过",
					logx.Field("cardType", cardType),
					logx.Field("objectName", data.ObjectName),
					logx.Field("attempt", attempt),
				)
			}
			return cardContent, nil
		}

		// JSON解析失败，记录错误并准备重试
		lastErr = parseErr
		lastContent = result.Content
		n.logger.Infow("JSON解析失败，准备重试",
			logx.Field("cardType", cardType),
			logx.Field("objectName", data.ObjectName),
			logx.Field("attempt", attempt),
			logx.Field("error", parseErr),
			logx.Field("contentPreview", func() string {
				if len(result.Content) > 200 {
					return result.Content[:200] + "..."
				}
				return result.Content
			}()),
		)
	}

	// 所有重试都失败
	n.logger.Errorw("所有重试均失败，JSON解析失败",
		logx.Field("cardType", cardType),
		logx.Field("objectName", data.ObjectName),
		logx.Field("maxRetries", maxRetries),
		logx.Field("lastError", lastErr),
		logx.Field("lastContentPreview", func() string {
			if len(lastContent) > 200 {
				return lastContent[:200] + "..."
			}
			return lastContent
		}()),
	)
	return nil, fmt.Errorf("JSON解析失败（已重试%d次）: %w, 最后返回内容: %s", maxRetries, lastErr, lastContent)
}

// generateScienceCardReal 真实eino实现科学认知卡
func (n *TextGenerationNode) generateScienceCardReal(ctx context.Context, data *GraphData) (map[string]interface{}, error) {
	n.logger.Infow("开始使用真实模型生成科学认知卡",
		logx.Field("objectName", data.ObjectName),
		logx.Field("age", data.Age),
	)

	// 使用带重试的生成方法（最多重试1次）
	cardContent, err := n.generateCardWithRetry(ctx, "science", data, n.scienceTemplate, 1)
	if err != nil {
		return nil, err
	}

	card := map[string]interface{}{
		"type":    "science",
		"title":   data.ObjectName + "的科学知识",
		"content": cardContent,
	}

	n.logger.Info("✅ 科学认知卡生成完成（真实模型）")
	return card, nil
}

// generatePoetryCardReal 真实eino实现古诗词卡
func (n *TextGenerationNode) generatePoetryCardReal(ctx context.Context, data *GraphData) (map[string]interface{}, error) {
	n.logger.Infow("开始使用真实模型生成古诗词卡",
		logx.Field("objectName", data.ObjectName),
		logx.Field("age", data.Age),
	)

	// 使用带重试的生成方法（最多重试1次）
	cardContent, err := n.generateCardWithRetry(ctx, "poetry", data, n.poetryTemplate, 1)
	if err != nil {
		return nil, err
	}

	card := map[string]interface{}{
		"type":    "poetry",
		"title":   "古人怎么看" + data.ObjectName,
		"content": cardContent,
	}

	n.logger.Infow("✅ 古诗词卡生成完成（真实模型）",
		logx.Field("objectName", data.ObjectName),
		logx.Field("title", card["title"]),
	)
	return card, nil
}

// generateEnglishCardReal 真实eino实现英语表达卡
func (n *TextGenerationNode) generateEnglishCardReal(ctx context.Context, data *GraphData) (map[string]interface{}, error) {
	n.logger.Infow("开始使用真实模型生成英语表达卡",
		logx.Field("objectName", data.ObjectName),
		logx.Field("age", data.Age),
	)

	// 使用带重试的生成方法（最多重试1次）
	cardContent, err := n.generateCardWithRetry(ctx, "english", data, n.englishTemplate, 1)
	if err != nil {
		return nil, err
	}

	card := map[string]interface{}{
		"type":    "english",
		"title":   "用英语说" + data.ObjectName,
		"content": cardContent,
	}

	n.logger.Info("✅ 英语表达卡生成完成（真实模型）")
	return card, nil
}
