package logic

import (
	"context"
	"math/rand"
	"time"

	"github.com/tango/explore/internal/svc"
	"github.com/tango/explore/internal/types"
	"github.com/tango/explore/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type IdentifyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIdentifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IdentifyLogic {
	return &IdentifyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IdentifyLogic) Identify(req *types.IdentifyRequest) (resp *types.IdentifyResponse, err error) {
	// 参数验证
	if req.Image == "" {
		return nil, utils.ErrImageRequired
	}

	// TODO: 待APP ID提供后，接入真实AI模型
	// 当前使用Mock数据
	l.Infow("识别图片", logx.Field("imageLength", len(req.Image)))

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

	resp = &types.IdentifyResponse{
		ObjectName:     selected.name,
		ObjectCategory: selected.category,
		Confidence:     confidence,
		Keywords:       selected.keywords,
	}

	l.Infow("识别完成", logx.Field("objectName", resp.ObjectName), logx.Field("category", resp.ObjectCategory), logx.Field("confidence", resp.Confidence))
	return resp, nil
}
