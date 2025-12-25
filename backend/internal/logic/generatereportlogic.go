package logic

import (
	"context"
	"sort"
	"time"

	"github.com/tango/explore/internal/svc"
	"github.com/tango/explore/internal/types"
	"github.com/tango/explore/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type GenerateReportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGenerateReportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateReportLogic {
	return &GenerateReportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GenerateReportLogic) GenerateReport(req *types.GenerateReportRequest) (resp *types.GenerateReportResponse, err error) {
	// 参数验证
	if req.ShareId == "" {
		return nil, utils.NewAPIError(400, "分享链接ID不能为空")
	}

	// 从内存存储获取分享数据
	store := GetShareStore()
	data, ok := store.Get(req.ShareId)
	if !ok {
		return nil, utils.ErrShareNotFound
	}

	// 统计探索次数
	totalExplorations := len(data.ExplorationRecords)

	// 统计收藏卡片数
	totalCollectedCards := len(data.CollectedCards)

	// 计算类别分布
	categoryDistribution := make(map[string]int)
	for _, record := range data.ExplorationRecords {
		category := record.ObjectCategory
		categoryDistribution[category]++
	}

	// 获取最近收藏的卡片（最多10张）
	recentCards := make([]types.KnowledgeCard, 0)
	if len(data.CollectedCards) > 0 {
		// 按收藏时间排序（假设CollectedAt字段已设置）
		sortedCards := make([]types.KnowledgeCard, len(data.CollectedCards))
		copy(sortedCards, data.CollectedCards)
		
		// 简单排序：有CollectedAt的在前
		sort.Slice(sortedCards, func(i, j int) bool {
			if sortedCards[i].CollectedAt == "" {
				return false
			}
			if sortedCards[j].CollectedAt == "" {
				return true
			}
			return sortedCards[i].CollectedAt > sortedCards[j].CollectedAt
		})

		// 取前10张
		maxCards := 10
		if len(sortedCards) < maxCards {
			maxCards = len(sortedCards)
		}
		recentCards = sortedCards[:maxCards]
	}

	resp = &types.GenerateReportResponse{
		TotalExplorations:    totalExplorations,
		TotalCollectedCards:   totalCollectedCards,
		CategoryDistribution: categoryDistribution,
		RecentCards:          recentCards,
		GeneratedAt:          time.Now().Format("2006-01-02T15:04:05Z07:00"),
	}

	l.Infow("生成学习报告", logx.Field("shareId", req.ShareId), logx.Field("totalExplorations", totalExplorations), logx.Field("totalCollectedCards", totalCollectedCards))
	return resp, nil
}
