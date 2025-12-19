package nodes

import (
	"context"
	"encoding/json"
	"strings"

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
	if cfg.EinoBaseURL != "" && cfg.AppID != "" && cfg.AppKey != "" {
		if err := node.initChatModel(ctx); err != nil {
			logger.Errorw("初始化ChatModel失败，将使用Mock模式", logx.Field("error", err))
		} else {
			node.initialized = true
			logger.Info("文本生成节点已初始化ChatModel")
		}
	} else {
		logger.Info("未配置eino参数，文本生成节点将使用Mock模式")
	}

	// 创建所有模板
	node.initTemplates()

	return node, nil
}

// initChatModel 初始化 ChatModel
func (n *TextGenerationNode) initChatModel(ctx context.Context) error {
	modelName := n.config.TextGenerationModel
	if modelName == "" {
		modelName = "gpt-5-nano" // 默认模型
	}

	cfg := &ark.ChatModelConfig{
		Model: modelName,
	}

	if n.config.EinoBaseURL != "" {
		cfg.BaseURL = n.config.EinoBaseURL
	}

	apiKey := n.config.AppKey
	if apiKey == "" {
		apiKey = n.config.AppID
	}
	if apiKey != "" {
		cfg.APIKey = apiKey
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

// initTemplates 初始化所有消息模板
func (n *TextGenerationNode) initTemplates() {
	// 科学认知卡模板
	n.scienceTemplate = prompt.FromMessages(schema.FString,
		schema.SystemMessage(`你是一个K12教育内容生成助手，专门为{age}岁的孩子生成科学认知卡片内容。

要求：
1. 用简单易懂的语言解释{objectName}的科学知识
2. 提供2-3个有趣的事实
3. 添加一个趣味知识
4. 内容要符合{age}岁孩子的认知水平

请返回JSON格式：
{
  "name": "{objectName}",
  "explanation": "科学解释",
  "facts": ["事实1", "事实2", "事实3"],
  "funFact": "趣味知识"
}`),
		schema.UserMessage("请为{objectName}生成科学认知卡内容，适合{age}岁孩子。"),
	)

	// 古诗词卡模板
	n.poetryTemplate = prompt.FromMessages(schema.FString,
		schema.SystemMessage(`你是一个古诗词专家，专门为K12教育生成古诗词卡片内容。

要求：
1. 找到与{objectName}相关的古诗词（优先选择经典名句）
2. 标注诗词来源（作者和诗名）
3. 用{age}岁孩子能理解的语言解释诗词含义
4. 提供文化背景说明

请返回JSON格式：
{
  "poem": "古诗词内容",
  "poemSource": "作者 - 诗名",
  "explanation": "诗词解释",
  "context": "文化背景"
}`),
		schema.UserMessage("请为{objectName}生成古诗词卡片内容，适合{age}岁孩子。"),
	)

	// 英语表达卡模板
	n.englishTemplate = prompt.FromMessages(schema.FString,
		schema.SystemMessage(`你是一个英语教学专家，专门为K12教育生成英语表达卡片内容。

要求：
1. 提供{objectName}的英语关键词（3-5个）
2. 提供2-3个适合{age}岁孩子的英语表达句子
3. 提供发音指导

请返回JSON格式：
{
  "keywords": ["关键词1", "关键词2", "关键词3"],
  "expressions": ["句子1", "句子2", "句子3"],
  "pronunciation": "发音指导"
}`),
		schema.UserMessage("请为{objectName}生成英语表达卡片内容，适合{age}岁孩子。"),
	)

	// 文本回答模板
	n.textTemplate = prompt.FromMessages(schema.FString,
		schema.SystemMessage("你是一个友好的K12教育助手，用简单易懂的语言回答孩子的问题。"),
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

	if n.initialized && n.chatModel != nil {
		return n.generateTextReal(data, context)
	}

	return n.generateTextMock(data, context)
}

// GenerateScienceCard 生成科学认知卡内容
func (n *TextGenerationNode) GenerateScienceCard(data *GraphData) (map[string]interface{}, error) {
	n.logger.Infow("生成科学认知卡",
		logx.Field("objectName", data.ObjectName),
		logx.Field("age", data.Age),
		logx.Field("useRealModel", n.initialized),
	)

	if n.initialized && n.chatModel != nil {
		return n.generateScienceCardReal(data)
	}

	return n.generateScienceCardMock(data)
}

// GeneratePoetryCard 生成古诗词/人文卡内容
func (n *TextGenerationNode) GeneratePoetryCard(data *GraphData) (map[string]interface{}, error) {
	n.logger.Infow("生成古诗词卡",
		logx.Field("objectName", data.ObjectName),
		logx.Field("age", data.Age),
		logx.Field("useRealModel", n.initialized),
	)

	if n.initialized && n.chatModel != nil {
		return n.generatePoetryCardReal(data)
	}

	return n.generatePoetryCardMock(data)
}

// GenerateEnglishCard 生成英语表达卡内容
func (n *TextGenerationNode) GenerateEnglishCard(data *GraphData) (map[string]interface{}, error) {
	n.logger.Infow("生成英语表达卡",
		logx.Field("objectName", data.ObjectName),
		logx.Field("age", data.Age),
		logx.Field("useRealModel", n.initialized),
	)

	if n.initialized && n.chatModel != nil {
		return n.generateEnglishCardReal(data)
	}

	return n.generateEnglishCardMock(data)
}

// generateTextMock Mock实现（待替换为真实eino调用）
func (n *TextGenerationNode) generateTextMock(data *GraphData, context []interface{}) (string, error) {
	// Mock文本响应
	text := "这是一个Mock文本响应。待接入真实AI模型后，将根据您的问题生成相应的回答。"
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
		explanation = data.ObjectName + "是一个有趣的对象，值得我们探索和学习。"
	}

	card := map[string]interface{}{
		"type":  "science",
		"title": data.ObjectName + "的科学知识",
		"content": map[string]interface{}{
			"name":        data.ObjectName,
			"explanation": explanation,
			"facts": []string{
				"关于" + data.ObjectName + "的有趣事实1",
				"关于" + data.ObjectName + "的有趣事实2",
			},
			"funFact": "关于" + data.ObjectName + "的趣味知识！",
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
	if poem == "" {
		poem = "关于" + data.ObjectName + "的古诗词，等待我们去发现。"
	}

	card := map[string]interface{}{
		"type":  "poetry",
		"title": "古人怎么看" + data.ObjectName,
		"content": map[string]interface{}{
			"poem":        poem,
			"poemSource":  "古诗词",
			"explanation": "这句诗描写了" + data.ObjectName + "的美丽景象，让我们感受到古人的智慧和情感。",
			"context":     "看到" + data.ObjectName + "，我们可以联想到相关的文化和历史，丰富我们的认知。",
		},
	}

	n.logger.Info("古诗词卡生成完成（Mock）")
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

// generateScienceCardReal 真实eino实现科学认知卡
func (n *TextGenerationNode) generateScienceCardReal(data *GraphData) (map[string]interface{}, error) {
	messages, err := n.scienceTemplate.Format(n.ctx, map[string]any{
		"objectName": data.ObjectName,
		"age":        data.Age,
	})
	if err != nil {
		n.logger.Errorw("模板格式化失败", logx.Field("error", err))
		return n.generateScienceCardMock(data)
	}

	result, err := n.chatModel.Generate(n.ctx, messages)
	if err != nil {
		n.logger.Errorw("ChatModel调用失败", logx.Field("error", err))
		return n.generateScienceCardMock(data)
	}

	// 解析 JSON 结果
	var cardContent map[string]interface{}
	text := result.Content
	jsonStart := strings.Index(text, "{")
	jsonEnd := strings.LastIndex(text, "}")
	if jsonStart >= 0 && jsonEnd > jsonStart {
		jsonStr := text[jsonStart : jsonEnd+1]
		if err := json.Unmarshal([]byte(jsonStr), &cardContent); err != nil {
			n.logger.Errorw("解析JSON失败", logx.Field("error", err))
			return n.generateScienceCardMock(data)
		}
	} else {
		return n.generateScienceCardMock(data)
	}

	card := map[string]interface{}{
		"type":    "science",
		"title":   data.ObjectName + "的科学知识",
		"content": cardContent,
	}

	n.logger.Info("科学认知卡生成完成（真实模型）")
	return card, nil
}

// generatePoetryCardReal 真实eino实现古诗词卡
func (n *TextGenerationNode) generatePoetryCardReal(data *GraphData) (map[string]interface{}, error) {
	messages, err := n.poetryTemplate.Format(n.ctx, map[string]any{
		"objectName": data.ObjectName,
		"age":        data.Age,
	})
	if err != nil {
		n.logger.Errorw("模板格式化失败", logx.Field("error", err))
		return n.generatePoetryCardMock(data)
	}

	result, err := n.chatModel.Generate(n.ctx, messages)
	if err != nil {
		n.logger.Errorw("ChatModel调用失败", logx.Field("error", err))
		return n.generatePoetryCardMock(data)
	}

	// 解析 JSON 结果
	var cardContent map[string]interface{}
	text := result.Content
	jsonStart := strings.Index(text, "{")
	jsonEnd := strings.LastIndex(text, "}")
	if jsonStart >= 0 && jsonEnd > jsonStart {
		jsonStr := text[jsonStart : jsonEnd+1]
		if err := json.Unmarshal([]byte(jsonStr), &cardContent); err != nil {
			n.logger.Errorw("解析JSON失败", logx.Field("error", err))
			return n.generatePoetryCardMock(data)
		}
	} else {
		return n.generatePoetryCardMock(data)
	}

	card := map[string]interface{}{
		"type":    "poetry",
		"title":   "古人怎么看" + data.ObjectName,
		"content": cardContent,
	}

	n.logger.Info("古诗词卡生成完成（真实模型）")
	return card, nil
}

// generateEnglishCardReal 真实eino实现英语表达卡
func (n *TextGenerationNode) generateEnglishCardReal(data *GraphData) (map[string]interface{}, error) {
	messages, err := n.englishTemplate.Format(n.ctx, map[string]any{
		"objectName": data.ObjectName,
		"age":        data.Age,
	})
	if err != nil {
		n.logger.Errorw("模板格式化失败", logx.Field("error", err))
		return n.generateEnglishCardMock(data)
	}

	result, err := n.chatModel.Generate(n.ctx, messages)
	if err != nil {
		n.logger.Errorw("ChatModel调用失败", logx.Field("error", err))
		return n.generateEnglishCardMock(data)
	}

	// 解析 JSON 结果
	var cardContent map[string]interface{}
	text := result.Content
	jsonStart := strings.Index(text, "{")
	jsonEnd := strings.LastIndex(text, "}")
	if jsonStart >= 0 && jsonEnd > jsonStart {
		jsonStr := text[jsonStart : jsonEnd+1]
		if err := json.Unmarshal([]byte(jsonStr), &cardContent); err != nil {
			n.logger.Errorw("解析JSON失败", logx.Field("error", err))
			return n.generateEnglishCardMock(data)
		}
	} else {
		return n.generateEnglishCardMock(data)
	}

	card := map[string]interface{}{
		"type":    "english",
		"title":   "用英语说" + data.ObjectName,
		"content": cardContent,
	}

	n.logger.Info("英语表达卡生成完成（真实模型）")
	return card, nil
}
