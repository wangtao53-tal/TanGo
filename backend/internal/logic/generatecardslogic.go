package logic

import (
	"context"

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

	// TODO: 待APP ID提供后，接入真实AI模型
	// 当前使用Mock数据，根据年龄调整内容难度
	l.Infow("生成知识卡片", logx.Field("objectName", req.ObjectName), logx.Field("category", req.ObjectCategory), logx.Field("age", req.Age))

	// Mock数据：根据对象名称和年龄生成三张卡片
	cards := []types.CardContent{
		// 科学认知卡
		{
			Type:  "science",
			Title: req.ObjectName + "的科学知识",
			Content: map[string]interface{}{
				"name":      req.ObjectName,
				"explanation": l.getScienceExplanation(req.ObjectName, req.Age),
				"facts":     l.getScienceFacts(req.ObjectName, req.Age),
				"funFact":   l.getFunFact(req.ObjectName, req.Age),
			},
		},
		// 古诗词/人文卡
		{
			Type:  "poetry",
			Title: "古人怎么看" + req.ObjectName,
			Content: map[string]interface{}{
				"poem":        l.getPoem(req.ObjectName),
				"poemSource":   l.getPoemSource(req.ObjectName),
				"explanation":  l.getPoemExplanation(req.ObjectName, req.Age),
				"context":     l.getContext(req.ObjectName, req.Age),
			},
		},
		// 英语表达卡
		{
			Type:  "english",
			Title: "用英语说" + req.ObjectName,
			Content: map[string]interface{}{
				"keywords":     l.getEnglishKeywords(req.ObjectName),
				"expressions":  l.getEnglishExpressions(req.ObjectName, req.Age),
				"pronunciation": l.getPronunciation(req.ObjectName),
			},
		},
	}

	resp = &types.GenerateCardsResponse{
		Cards: cards,
	}

	l.Infow("卡片生成完成", logx.Field("cardCount", len(cards)))
	return resp, nil
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
