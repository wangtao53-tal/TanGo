package nodes

import (
	"context"
	"math/rand"
	"time"

	"github.com/tango/explore/internal/config"
	"github.com/zeromicro/go-zero/core/logx"
)

// ImageRecognitionNode 图片识别节点
type ImageRecognitionNode struct {
	ctx    context.Context
	config config.AIConfig
	logger logx.Logger
}

// ImageRecognitionResult 图片识别结果
type ImageRecognitionResult struct {
	ObjectName     string
	ObjectCategory string
	Keywords       []string
	Confidence     float64
}

// NewImageRecognitionNode 创建图片识别节点
func NewImageRecognitionNode(ctx context.Context, cfg config.AIConfig, logger logx.Logger) (*ImageRecognitionNode, error) {
	return &ImageRecognitionNode{
		ctx:    ctx,
		config: cfg,
		logger: logger,
	}, nil
}

// Execute 执行图片识别
func (n *ImageRecognitionNode) Execute(data *GraphData) (*ImageRecognitionResult, error) {
	n.logger.Infow("执行图片识别",
		logx.Field("imageLength", len(data.Image)),
		logx.Field("age", data.Age),
	)

	// TODO: 待APP ID提供后，接入真实eino框架调用图片识别模型
	// 当前使用Mock数据
	return n.executeMock(data)
}

// executeMock Mock实现（待替换为真实eino调用）
func (n *ImageRecognitionNode) executeMock(data *GraphData) (*ImageRecognitionResult, error) {
	// Mock识别结果 - 随机返回一个常见对象
	mockObjects := []struct {
		name     string
		category string
		keywords []string
	}{
		{"银杏", "自然类", []string{"植物", "树木", "秋天", "叶子"}},
		{"苹果", "生活类", []string{"水果", "食物", "红色", "健康"}},
		{"蝴蝶", "自然类", []string{"昆虫", "飞行", "美丽", "春天"}},
		{"书本", "人文类", []string{"学习", "知识", "阅读", "教育"}},
		{"汽车", "生活类", []string{"交通工具", "速度", "现代", "出行"}},
		{"月亮", "自然类", []string{"天体", "夜晚", "圆形", "美丽"}},
		{"钢琴", "人文类", []string{"乐器", "音乐", "艺术", "优雅"}},
		{"太阳", "自然类", []string{"恒星", "光明", "温暖", "能量"}},
	}

	rand.Seed(time.Now().UnixNano())
	selected := mockObjects[rand.Intn(len(mockObjects))]

	// 生成随机置信度（0.85-0.99）
	confidence := 0.85 + rand.Float64()*0.14

	result := &ImageRecognitionResult{
		ObjectName:     selected.name,
		ObjectCategory: selected.category,
		Keywords:       selected.keywords,
		Confidence:     confidence,
	}

	n.logger.Infow("图片识别完成（Mock）",
		logx.Field("objectName", result.ObjectName),
		logx.Field("category", result.ObjectCategory),
		logx.Field("confidence", result.Confidence),
	)

	return result, nil
}

// executeReal 真实eino实现（待APP ID提供后实现）
func (n *ImageRecognitionNode) executeReal(data *GraphData) (*ImageRecognitionResult, error) {
	// TODO: 使用eino框架调用图片识别模型
	// 1. 从config中获取模型配置
	// 2. 使用eino的ChatModel或VisionModel接口
	// 3. 调用模型API（通过EinoBaseURL + AppID/AppKey认证）
	// 4. 解析返回结果

	// 示例代码结构（待实现）：
	// model := eino.NewChatModel(...)
	// response, err := model.Invoke(ctx, messages)
	// ...

	return nil, nil
}
