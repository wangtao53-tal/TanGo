package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/tango/explore/internal/svc"
	"github.com/tango/explore/internal/types"
	"github.com/tango/explore/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type GenerateCardsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGenerateCardsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateCardsLogic {
	return &GenerateCardsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GenerateCardsLogic) GenerateCards(req *types.GenerateCardsRequest) (resp *types.GenerateCardsResponse, err error) {
	// 参数验证
	if req.ObjectName == "" {
		return nil, utils.ErrObjectNameRequired
	}
	if req.ObjectCategory == "" {
		return nil, utils.ErrCategoryRequired
	}
	if req.Age < 3 || req.Age > 18 {
		return nil, utils.ErrInvalidAge
	}

	l.Infow("生成知识卡片", logx.Field("objectName", req.ObjectName), logx.Field("category", req.ObjectCategory), logx.Field("age", req.Age))

	// 使用Agent系统生成卡片
	if l.svcCtx.Agent != nil {
		graph := l.svcCtx.Agent.GetGraph()
		data, err := graph.ExecuteCardGeneration(l.ctx, req.ObjectName, req.ObjectCategory, req.Age, req.Keywords)
		if err != nil {
			l.Errorw("Agent卡片生成失败，返回错误",
				logx.Field("error", err),
				logx.Field("errorDetail", err.Error()),
			)
			// 不再回退到Mock，直接返回错误
			return nil, fmt.Errorf("卡片生成失败: %w", err)
		}

		// 转换Agent返回的卡片数据为types.CardContent
		cards := make([]types.CardContent, 0, len(data.Cards))
		for _, cardData := range data.Cards {
			if cardMap, ok := cardData.(map[string]interface{}); ok {
				card := types.CardContent{
					Type:    getString(cardMap, "type"),
					Title:   getString(cardMap, "title"),
					Content: cardMap["content"],
				}
				cards = append(cards, card)
			}
		}

		resp = &types.GenerateCardsResponse{
			Cards: cards,
		}

		l.Infow("卡片生成完成（Agent）", logx.Field("cardCount", len(cards)))
		return resp, nil
	}

	// 如果Agent未初始化，使用Mock数据
	return l.generateCardsMock(req)
}

// GenerateCardsStream 流式生成知识卡片（每生成完一张立即返回）
func (l *GenerateCardsLogic) GenerateCardsStream(w http.ResponseWriter, req *types.GenerateCardsRequest) error {
	// 参数验证
	if req.ObjectName == "" {
		return utils.ErrObjectNameRequired
	}
	if req.ObjectCategory == "" {
		return utils.ErrCategoryRequired
	}
	if req.Age < 3 || req.Age > 18 {
		return utils.ErrInvalidAge
	}

	l.Infow("开始流式生成知识卡片",
		logx.Field("objectName", req.ObjectName),
		logx.Field("category", req.ObjectCategory),
		logx.Field("age", req.Age),
	)

	// 使用Agent系统生成卡片
	if l.svcCtx.Agent != nil {
		graph := l.svcCtx.Agent.GetGraph()

		l.Infow("使用Agent系统生成卡片",
			logx.Field("objectName", req.ObjectName),
			logx.Field("category", req.ObjectCategory),
			logx.Field("age", req.Age),
			logx.Field("graphNil", graph == nil),
		)

		// 调用ExecuteCardGeneration（并行生成，等待模型返回，不设置超时）
		// 超时控制由HTTP请求层面的Timeout配置控制（在explore.yaml中配置为180秒）
		data, err := graph.ExecuteCardGeneration(l.ctx, req.ObjectName, req.ObjectCategory, req.Age, req.Keywords)
		if err != nil {
			l.Errorw("卡片生成失败，使用Mock",
				logx.Field("error", err),
				logx.Field("errorDetail", err.Error()),
			)
			return l.generateCardsStreamMock(w, req)
		}

		l.Infow("Agent卡片生成成功",
			logx.Field("cardCount", len(data.Cards)),
		)

		// 转换并立即发送每张卡片
		// 由于ExecuteCardGeneration已经并行生成，这里按顺序发送
		// 未来可以优化为真正的流式返回（每生成完一张立即发送）
		cardCount := 0
		for i, cardData := range data.Cards {
			if cardMap, ok := cardData.(map[string]interface{}); ok {
				card := types.CardContent{
					Type:    getString(cardMap, "type"),
					Title:   getString(cardMap, "title"),
					Content: cardMap["content"],
				}
				// 立即发送卡片事件
				cardEvent := map[string]interface{}{
					"type":    "card",
					"content": card,
					"index":   i,
				}
				cardJSON, _ := json.Marshal(cardEvent)
				fmt.Fprintf(w, "event: card\ndata: %s\n\n", string(cardJSON))
				w.(http.Flusher).Flush()
				cardCount++
			}
		}

		// 发送完成事件
		doneEvent := map[string]interface{}{
			"type": "done",
		}
		doneJSON, _ := json.Marshal(doneEvent)
		fmt.Fprintf(w, "event: done\ndata: %s\n\n", string(doneJSON))
		w.(http.Flusher).Flush()

		l.Infow("流式卡片生成完成",
			logx.Field("cardCount", cardCount),
		)
		return nil
	}

	// 如果Agent未初始化，使用Mock数据流式返回
	l.Errorw("Agent未初始化，使用Mock数据流式返回",
		logx.Field("agentNil", l.svcCtx.Agent == nil),
	)
	return l.generateCardsStreamMock(w, req)
}

// generateCardsStreamMock Mock流式返回
func (l *GenerateCardsLogic) generateCardsStreamMock(w http.ResponseWriter, req *types.GenerateCardsRequest) error {
	cards := []types.CardContent{
		l.getMockCardByIndex(0, req.ObjectName, req.Age),
		l.getMockCardByIndex(1, req.ObjectName, req.Age),
		l.getMockCardByIndex(2, req.ObjectName, req.Age),
	}

	// 模拟流式返回，每张卡片间隔100ms
	for i, card := range cards {
		cardEvent := map[string]interface{}{
			"type":    "card",
			"content": card,
			"index":   i,
		}
		cardJSON, _ := json.Marshal(cardEvent)
		fmt.Fprintf(w, "event: card\ndata: %s\n\n", string(cardJSON))
		w.(http.Flusher).Flush()
		time.Sleep(100 * time.Millisecond) // 模拟生成延迟
	}

	// 发送完成事件
	doneEvent := map[string]interface{}{
		"type": "done",
	}
	doneJSON, _ := json.Marshal(doneEvent)
	fmt.Fprintf(w, "event: done\ndata: %s\n\n", string(doneJSON))
	w.(http.Flusher).Flush()

	return nil
}

// getMockCardByIndex 根据索引获取Mock卡片
func (l *GenerateCardsLogic) getMockCardByIndex(idx int, objectName string, age int) types.CardContent {
	switch idx {
	case 0: // 科学卡
		return types.CardContent{
			Type:  "science",
			Title: objectName + "的科学知识",
			Content: map[string]interface{}{
				"name":        objectName,
				"explanation": l.getScienceExplanation(objectName, age),
				"facts":       l.getScienceFacts(objectName, age),
				"funFact":     l.getFunFact(objectName, age),
			},
		}
	case 1: // 诗词卡
		return types.CardContent{
			Type:  "poetry",
			Title: "古人怎么看" + objectName,
			Content: map[string]interface{}{
				"poem":        l.getPoem(objectName),
				"poemSource":  l.getPoemSource(objectName),
				"explanation": l.getPoemExplanation(objectName, age),
				"context":     l.getContext(objectName, age),
			},
		}
	case 2: // 英语卡
		return types.CardContent{
			Type:  "english",
			Title: "用英语说" + objectName,
			Content: map[string]interface{}{
				"keywords":      l.getEnglishKeywords(objectName),
				"expressions":   l.getEnglishExpressions(objectName, age),
				"pronunciation": l.getPronunciation(objectName),
			},
		}
	default:
		return types.CardContent{}
	}
}

// generateCardsMock Mock实现（保留作为回退方案）
func (l *GenerateCardsLogic) generateCardsMock(req *types.GenerateCardsRequest) (*types.GenerateCardsResponse, error) {
	// Mock数据：根据对象名称和年龄生成三张卡片
	cards := []types.CardContent{
		// 科学认知卡
		{
			Type:  "science",
			Title: req.ObjectName + "的科学知识",
			Content: map[string]interface{}{
				"name":        req.ObjectName,
				"explanation": l.getScienceExplanation(req.ObjectName, req.Age),
				"facts":       l.getScienceFacts(req.ObjectName, req.Age),
				"funFact":     l.getFunFact(req.ObjectName, req.Age),
			},
		},
		// 古诗词/人文卡
		{
			Type:  "poetry",
			Title: "古人怎么看" + req.ObjectName,
			Content: map[string]interface{}{
				"poem":        l.getPoem(req.ObjectName),
				"poemSource":  l.getPoemSource(req.ObjectName),
				"explanation": l.getPoemExplanation(req.ObjectName, req.Age),
				"context":     l.getContext(req.ObjectName, req.Age),
			},
		},
		// 英语表达卡
		{
			Type:  "english",
			Title: "用英语说" + req.ObjectName,
			Content: map[string]interface{}{
				"keywords":      l.getEnglishKeywords(req.ObjectName),
				"expressions":   l.getEnglishExpressions(req.ObjectName, req.Age),
				"pronunciation": l.getPronunciation(req.ObjectName),
			},
		},
	}

	resp := &types.GenerateCardsResponse{
		Cards: cards,
	}

	l.Infow("卡片生成完成（Mock）", logx.Field("cardCount", len(cards)))
	return resp, nil
}

// getString 辅助函数：从map中安全获取string值
func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// Mock辅助函数：生成科学认知内容
func (l *GenerateCardsLogic) getScienceExplanation(name string, age int) string {
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
	if exp, ok := explanations[name]; ok {
		return exp
	}
	return name + "是一个有趣的对象，值得我们探索和学习。"
}

func (l *GenerateCardsLogic) getScienceFacts(name string, age int) []string {
	facts := map[string][]string{
		"银杏": {"银杏是现存最古老的树种之一", "银杏的叶子在秋天会变成金黄色", "银杏的果实可以食用，但需要处理"},
		"苹果": {"苹果含有丰富的维生素C", "每天一个苹果有助于健康", "苹果有很多品种，颜色和味道不同"},
		"蝴蝶": {"蝴蝶有美丽的翅膀", "蝴蝶可以帮助传播花粉", "不同种类的蝴蝶有不同的颜色"},
	}
	if f, ok := facts[name]; ok {
		return f
	}
	return []string{"这是一个有趣的事实", "还有更多知识等待探索"}
}

func (l *GenerateCardsLogic) getFunFact(name string, age int) string {
	facts := map[string]string{
		"银杏": "银杏被称为'活化石'，因为它在恐龙时代就已经存在了！",
		"苹果": "苹果的种子含有少量氰化物，但吃几个苹果不会中毒！",
		"蝴蝶": "蝴蝶的翅膀上有细小的鳞片，这些鳞片创造了美丽的颜色！",
	}
	if f, ok := facts[name]; ok {
		return f
	}
	return "关于" + name + "还有很多有趣的知识等待发现！"
}

// Mock辅助函数：生成古诗词内容
func (l *GenerateCardsLogic) getPoem(name string) string {
	poems := map[string]string{
		"银杏": "满地翻黄银杏叶，忽惊天地告成功。",
		"苹果": "苹果红时秋已深，满园香气醉人心。",
		"蝴蝶": "穿花蛱蝶深深见，点水蜻蜓款款飞。",
		"月亮": "床前明月光，疑是地上霜。",
		"太阳": "日出江花红胜火，春来江水绿如蓝。",
	}
	if p, ok := poems[name]; ok {
		return p
	}
	return "关于" + name + "的古诗词，等待我们去发现。"
}

func (l *GenerateCardsLogic) getPoemSource(name string) string {
	sources := map[string]string{
		"银杏": "《夜坐》- 李清照",
		"苹果": "现代诗歌",
		"蝴蝶": "《曲江二首》- 杜甫",
		"月亮": "《静夜思》- 李白",
		"太阳": "《忆江南》- 白居易",
	}
	if s, ok := sources[name]; ok {
		return s
	}
	return "古诗词"
}

func (l *GenerateCardsLogic) getPoemExplanation(name string, age int) string {
	return "这句诗描写了" + name + "的美丽景象，让我们感受到古人的智慧和情感。"
}

func (l *GenerateCardsLogic) getContext(name string, age int) string {
	return "看到" + name + "，我们可以联想到相关的文化和历史，丰富我们的认知。"
}

// Mock辅助函数：生成英语表达内容
func (l *GenerateCardsLogic) getEnglishKeywords(name string) []string {
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
	if k, ok := keywords[name]; ok {
		return k
	}
	return []string{name, "object", "interesting"}
}

func (l *GenerateCardsLogic) getEnglishExpressions(name string, age int) []string {
	expressions := map[string][]string{
		"银杏": {"This is a ginkgo tree.", "The ginkgo leaves are golden in autumn."},
		"苹果": {"This is an apple.", "I like to eat apples."},
		"蝴蝶": {"Look at the beautiful butterfly!", "Butterflies fly in the garden."},
	}
	if e, ok := expressions[name]; ok {
		return e
	}
	return []string{"This is " + name + ".", "It's very interesting."}
}

func (l *GenerateCardsLogic) getPronunciation(name string) string {
	pronunciations := map[string]string{
		"银杏": "ginkgo: /ˈɡɪŋkoʊ/",
		"苹果": "apple: /ˈæpl/",
		"蝴蝶": "butterfly: /ˈbʌtərflaɪ/",
	}
	if p, ok := pronunciations[name]; ok {
		return p
	}
	return name + ": pronunciation"
}
