package agent

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/tango/explore/internal/agent/nodes"
	"github.com/tango/explore/internal/config"
	"github.com/zeromicro/go-zero/core/logx"
)

// Graph AI调用流程图
type Graph struct {
	ctx    context.Context
	config config.AIConfig
	logger logx.Logger

	// 节点实例
	imageRecognitionNode  *nodes.ImageRecognitionNode
	textGenerationNode    *nodes.TextGenerationNode
	imageGenerationNode   *nodes.ImageGenerationNode
	intentRecognitionNode *nodes.IntentRecognitionNode
	conversationNode      *nodes.ConversationNode
}

// NewGraph 创建新的Graph实例
func NewGraph(ctx context.Context, cfg config.AIConfig, logger logx.Logger) (*Graph, error) {
	graph := &Graph{
		ctx:    ctx,
		config: cfg,
		logger: logger,
	}

	// 初始化各个节点
	var err error

	graph.imageRecognitionNode, err = nodes.NewImageRecognitionNode(ctx, cfg, logger)
	if err != nil {
		return nil, err
	}

	graph.textGenerationNode, err = nodes.NewTextGenerationNode(ctx, cfg, logger)
	if err != nil {
		return nil, err
	}

	graph.imageGenerationNode, err = nodes.NewImageGenerationNode(ctx, cfg, logger)
	if err != nil {
		return nil, err
	}

	graph.intentRecognitionNode, err = nodes.NewIntentRecognitionNode(ctx, cfg, logger)
	if err != nil {
		return nil, err
	}

	graph.conversationNode, err = nodes.NewConversationNode(ctx, cfg, logger)
	if err != nil {
		return nil, err
	}

	logger.Info("Graph初始化完成")
	return graph, nil
}

// ExecuteImageRecognition 执行图片识别流程
// 输入: 图片 -> 输出: 对象名称、类别、关键词
func (g *Graph) ExecuteImageRecognition(image string, age int) (*nodes.GraphData, error) {
	data := &nodes.GraphData{
		Image: image,
		Age:   age,
	}

	// 调用图片识别节点
	result, err := g.imageRecognitionNode.Execute(data)
	if err != nil {
		return nil, err
	}

	data.ObjectName = result.ObjectName
	data.ObjectCategory = result.ObjectCategory
	data.Keywords = result.Keywords

	g.logger.Infow("图片识别完成",
		logx.Field("objectName", data.ObjectName),
		logx.Field("category", data.ObjectCategory),
	)

	return data, nil
}

// ExecuteCardGeneration 执行卡片生成流程
// 输入: 对象名称、类别、年龄 -> 输出: 三张卡片（科学、诗词、英语）
// 优化：并行生成三张卡片以减少响应时间，添加超时控制
func (g *Graph) ExecuteCardGeneration(ctx context.Context, objectName, category string, age int, keywords []string) (*nodes.GraphData, error) {
	data := &nodes.GraphData{
		ObjectName:     objectName,
		ObjectCategory: category,
		Age:            age,
		Keywords:       keywords,
	}

	// 使用 WaitGroup 和 channel 并行生成三张卡片
	var wg sync.WaitGroup
	type cardResult struct {
		card interface{}
		err  error
		idx  int // 用于保持顺序：0-科学卡, 1-诗词卡, 2-英语卡
	}
	results := make(chan cardResult, 3)

	// 为每张卡片设置10秒超时
	cardTimeout := 10 * time.Second

	// 1. 并行生成科学认知卡（带超时控制）
	wg.Add(1)
	go func() {
		defer wg.Done()
		cardCtx, cancel := context.WithTimeout(ctx, cardTimeout)
		defer cancel()
		card, err := g.textGenerationNode.GenerateScienceCard(cardCtx, data)
		if cardCtx.Err() == context.DeadlineExceeded {
			err = fmt.Errorf("科学卡生成超时（超过%v）", cardTimeout)
		}
		results <- cardResult{card: card, err: err, idx: 0}
	}()

	// 2. 并行生成古诗词/人文卡（带超时控制）
	wg.Add(1)
	go func() {
		defer wg.Done()
		cardCtx, cancel := context.WithTimeout(ctx, cardTimeout)
		defer cancel()
		card, err := g.textGenerationNode.GeneratePoetryCard(cardCtx, data)
		if cardCtx.Err() == context.DeadlineExceeded {
			err = fmt.Errorf("诗词卡生成超时（超过%v）", cardTimeout)
		}
		results <- cardResult{card: card, err: err, idx: 1}
	}()

	// 3. 并行生成英语表达卡（带超时控制）
	wg.Add(1)
	go func() {
		defer wg.Done()
		cardCtx, cancel := context.WithTimeout(ctx, cardTimeout)
		defer cancel()
		card, err := g.textGenerationNode.GenerateEnglishCard(cardCtx, data)
		if cardCtx.Err() == context.DeadlineExceeded {
			err = fmt.Errorf("英语卡生成超时（超过%v）", cardTimeout)
		}
		results <- cardResult{card: card, err: err, idx: 2}
	}()

	// 等待所有 goroutine 完成
	wg.Wait()
	close(results)

	// 收集结果并保持顺序
	cards := make([]interface{}, 3)
	var firstErr error
	successCount := 0
	for result := range results {
		if result.err != nil {
			if firstErr == nil {
				firstErr = result.err
			}
			g.logger.Errorw("卡片生成失败",
				logx.Field("cardIndex", result.idx),
				logx.Field("error", result.err),
			)
			// 超时或错误时，使用Mock数据作为降级方案
			cards[result.idx] = g.getMockCard(result.idx, objectName, age)
		} else {
			cards[result.idx] = result.card
			successCount++
		}
	}

	// 如果所有卡片都失败，返回错误
	if successCount == 0 && firstErr != nil {
		return nil, firstErr
	}

	// 如果有部分成功，记录警告但继续
	if successCount < 3 {
		g.logger.Infow("部分卡片生成失败，使用降级方案",
			logx.Field("successCount", successCount),
			logx.Field("totalCount", 3),
		)
	}

	// 4. 为每张卡片生成配图（图片生成）
	// TODO: 待APP ID提供后，启用图片生成
	// for i := range cards {
	// 	imageURL, err := g.imageGenerationNode.GenerateCardImage(data, cards[i])
	// 	if err != nil {
	// 		g.logger.Errorw("生成卡片配图失败", logx.Field("error", err))
	// 		continue
	// 	}
	// 	// 将图片URL添加到卡片数据中
	// }

	data.Cards = cards

	g.logger.Infow("卡片生成完成（并行）",
		logx.Field("cardCount", len(cards)),
		logx.Field("successCount", successCount),
	)
	return data, nil
}

// getMockCard 获取Mock卡片作为降级方案
func (g *Graph) getMockCard(idx int, objectName string, age int) map[string]interface{} {
	switch idx {
	case 0: // 科学卡
		return map[string]interface{}{
			"type":  "science",
			"title": objectName + "的科学知识",
			"content": map[string]interface{}{
				"name":        objectName,
				"explanation": objectName + "是一个有趣的对象，值得我们探索和学习。",
				"facts":       []string{"这是一个有趣的事实", "还有更多知识等待探索"},
				"funFact":     "关于" + objectName + "还有很多有趣的知识等待发现！",
			},
		}
	case 1: // 诗词卡
		return map[string]interface{}{
			"type":  "poetry",
			"title": "古人怎么看" + objectName,
			"content": map[string]interface{}{
				"poem":        "关于" + objectName + "的古诗词，等待我们去发现。",
				"poemSource":  "古诗词",
				"explanation": "这句诗描写了" + objectName + "的美丽景象，让我们感受到古人的智慧和情感。",
				"context":     "看到" + objectName + "，我们可以联想到相关的文化和历史，丰富我们的认知。",
			},
		}
	case 2: // 英语卡
		return map[string]interface{}{
			"type":  "english",
			"title": "用英语说" + objectName,
			"content": map[string]interface{}{
				"keywords":      []string{objectName, "object", "interesting"},
				"expressions":   []string{"This is " + objectName + ".", "It's very interesting."},
				"pronunciation": objectName + ": pronunciation",
			},
		}
	default:
		return map[string]interface{}{}
	}
}

// ExecuteIntentRecognition 执行意图识别流程
// 输入: 文本消息、上下文 -> 输出: 意图类型
func (g *Graph) ExecuteIntentRecognition(message string, context []interface{}) (*nodes.GraphData, error) {
	data := &nodes.GraphData{
		Text: message,
	}

	// 调用意图识别节点
	result, err := g.intentRecognitionNode.Execute(data, context)
	if err != nil {
		return nil, err
	}

	data.Intent = result.Intent

	g.logger.Infow("意图识别完成",
		logx.Field("intent", data.Intent),
		logx.Field("confidence", result.Confidence),
	)

	return data, nil
}

// ExecuteTextGeneration 执行文本生成流程
// 输入: 消息、上下文 -> 输出: 文本回答
func (g *Graph) ExecuteTextGeneration(message string, context []interface{}) (*nodes.GraphData, error) {
	data := &nodes.GraphData{
		Text: message,
	}

	// 调用文本生成节点
	result, err := g.textGenerationNode.GenerateText(data, context)
	if err != nil {
		return nil, err
	}

	data.TextResult = result

	g.logger.Infow("文本生成完成", logx.Field("length", len(result)))
	return data, nil
}

// GetConversationNode 获取对话节点
func (g *Graph) GetConversationNode() *nodes.ConversationNode {
	return g.conversationNode
}
