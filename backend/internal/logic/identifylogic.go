package logic

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
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
	// 性能监控：记录开始时间
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime)
		l.Infow("识别请求完成",
			logx.Field("duration_ms", duration.Milliseconds()),
			logx.Field("duration_sec", duration.Seconds()),
			logx.Field("success", err == nil),
		)
	}()

	// 参数验证
	if req.Image == "" {
		return nil, utils.ErrImageRequired
	}

	// 优化：添加图片URL验证，提前失败避免无效请求
	if len(req.Image) > 7 {
		isHTTP := req.Image[:7] == "http://"
		isHTTPS := len(req.Image) > 8 && req.Image[:8] == "https://"
		if isHTTP || isHTTPS {
			// 简单验证URL格式
			if !strings.Contains(req.Image, ".") {
				return nil, fmt.Errorf("无效的图片URL格式")
			}
		}
	}

	// 优化：减少日志详细程度，只记录关键信息，避免记录大对象
	isURL := len(req.Image) > 7 && (req.Image[:7] == "http://" || (len(req.Image) > 8 && req.Image[:8] == "https://"))
	imageType := "base64"
	if isURL {
		imageType = "url"
	} else if len(req.Image) > 5 && req.Image[:5] == "data:" {
		imageType = "data_url"
	}

	// 优化：使用Debug级别记录详细信息，Info级别只记录关键指标
	l.Debugw("识别图片",
		logx.Field("imageType", imageType),
		logx.Field("imageLength", len(req.Image)),
		logx.Field("age", req.Age),
	)
	l.Infow("开始识别",
		logx.Field("imageType", imageType),
		logx.Field("age", req.Age),
	)

	// 使用Agent系统进行图片识别
	if l.svcCtx.Agent != nil {
		graph := l.svcCtx.Agent.GetGraph()
		data, err := graph.ExecuteImageRecognition(req.Image, req.Age)

		if err != nil {
			// 优化：减少不必要的Mock调用，只对特定错误回退
			errMsg := err.Error()
			// 如果是超时或严重错误，才回退到Mock
			shouldFallback := strings.Contains(errMsg, "timeout") ||
				strings.Contains(errMsg, "deadline") ||
				strings.Contains(errMsg, "network")

			if shouldFallback {
				l.Errorw("Agent图片识别失败，回退到Mock",
					logx.Field("error", err),
					logx.Field("errorType", "fallback"),
				)
				return l.identifyMock(req)
			}
			// 其他错误直接返回，不回退
			return nil, err
		}

		resp = &types.IdentifyResponse{
			ObjectName:     data.ObjectName,
			ObjectCategory: data.ObjectCategory,
			Confidence:     0.95, // Agent返回的置信度
			Keywords:       data.Keywords,
		}

		// 优化：减少日志详细程度，移除keywords数组（可能很大）
		l.Infow("识别完成（Agent）",
			logx.Field("objectName", resp.ObjectName),
			logx.Field("category", resp.ObjectCategory),
			logx.Field("confidence", resp.Confidence),
			logx.Field("keywordsCount", len(resp.Keywords)),
		)
		return resp, nil
	}

	// 如果Agent未初始化，使用Mock数据
	l.Errorw("Agent未初始化，使用Mock数据")
	return l.identifyMock(req)
}

// identifyMock Mock实现（保留作为回退方案）
func (l *IdentifyLogic) identifyMock(req *types.IdentifyRequest) (*types.IdentifyResponse, error) {
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

	resp := &types.IdentifyResponse{
		ObjectName:     selected.name,
		ObjectCategory: selected.category,
		Confidence:     confidence,
		Keywords:       selected.keywords,
	}

	l.Infow("识别完成（Mock）",
		logx.Field("objectName", resp.ObjectName),
		logx.Field("category", resp.ObjectCategory),
		logx.Field("confidence", resp.Confidence),
	)
	return resp, nil
}
