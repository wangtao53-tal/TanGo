package nodes

import (
	"context"

	"github.com/tango/explore/internal/config"
	"github.com/zeromicro/go-zero/core/logx"
)

// TextGenerationNode 文本生成节点
type TextGenerationNode struct {
	ctx    context.Context
	config config.AIConfig
	logger logx.Logger
}

// NewTextGenerationNode 创建文本生成节点
func NewTextGenerationNode(ctx context.Context, cfg config.AIConfig, logger logx.Logger) (*TextGenerationNode, error) {
	return &TextGenerationNode{
		ctx:    ctx,
		config: cfg,
		logger: logger,
	}, nil
}

// GenerateText 生成文本回答
func (n *TextGenerationNode) GenerateText(data *GraphData, context []interface{}) (string, error) {
	n.logger.Infow("执行文本生成",
		logx.Field("message", data.Text),
		logx.Field("contextLength", len(context)),
	)

	// TODO: 待APP ID提供后，接入真实eino框架调用文本生成模型
	// 当前使用Mock数据
	return n.generateTextMock(data, context)
}

// GenerateScienceCard 生成科学认知卡内容
func (n *TextGenerationNode) GenerateScienceCard(data *GraphData) (map[string]interface{}, error) {
	n.logger.Infow("生成科学认知卡",
		logx.Field("objectName", data.ObjectName),
		logx.Field("age", data.Age),
	)

	// TODO: 待APP ID提供后，接入真实eino框架
	// 当前使用Mock数据
	return n.generateScienceCardMock(data)
}

// GeneratePoetryCard 生成古诗词/人文卡内容
func (n *TextGenerationNode) GeneratePoetryCard(data *GraphData) (map[string]interface{}, error) {
	n.logger.Infow("生成古诗词卡",
		logx.Field("objectName", data.ObjectName),
		logx.Field("age", data.Age),
	)

	// TODO: 待APP ID提供后，接入真实eino框架
	// 当前使用Mock数据
	return n.generatePoetryCardMock(data)
}

// GenerateEnglishCard 生成英语表达卡内容
func (n *TextGenerationNode) GenerateEnglishCard(data *GraphData) (map[string]interface{}, error) {
	n.logger.Infow("生成英语表达卡",
		logx.Field("objectName", data.ObjectName),
		logx.Field("age", data.Age),
	)

	// TODO: 待APP ID提供后，接入真实eino框架
	// 当前使用Mock数据
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

// generateTextReal 真实eino实现（待APP ID提供后实现）
func (n *TextGenerationNode) generateTextReal(data *GraphData, context []interface{}) (string, error) {
	// TODO: 使用eino框架调用文本生成模型
	// 1. 从config中获取TextGenerationModel配置
	// 2. 使用eino的ChatModel接口
	// 3. 构建prompt，包含上下文
	// 4. 调用模型API
	// 5. 解析返回结果

	return "", nil
}
