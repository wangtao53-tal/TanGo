package logic

import (
	"context"

	"github.com/tango/explore/internal/svc"
	"github.com/tango/explore/internal/types"
	"github.com/tango/explore/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBadgeStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetBadgeStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBadgeStatsLogic {
	return &GetBadgeStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetBadgeStats 获取勋章统计信息
func (l *GetBadgeStatsLogic) GetBadgeStats(req *types.GetBadgeStatsRequest) (resp *types.BadgeDetailResponse, err error) {
	// 参数验证
	if req.ExplorationCount < 0 || req.CollectionCount < 0 || req.ConversationCount < 0 {
		return nil, utils.NewAPIError(400, "统计数据不能为负数")
	}

	// 计算总分
	// 探索次数权重：10分/次
	// 收藏次数权重：5分/次
	// 对话次数权重：3分/次
	totalScore := req.ExplorationCount*10 + req.CollectionCount*5 + req.ConversationCount*3

	// 获取所有等级信息
	allLevels := getAllBadgeLevels()

	// 计算当前等级
	currentLevel := calculateLevel(totalScore, allLevels)
	currentLevelInfo := allLevels[currentLevel-1]

	// 获取下一等级信息
	var nextLevelInfo types.BadgeLevel
	var hasNextLevel bool
	if currentLevel < 10 {
		nextLevelInfo = allLevels[currentLevel]
		hasNextLevel = true
	}

	// 计算进度
	progress := calculateProgress(totalScore, currentLevelInfo, nextLevelInfo, hasNextLevel)

	// 构建响应
	stats := types.UserStats{
		ExplorationCount:  req.ExplorationCount,
		CollectionCount:    req.CollectionCount,
		ConversationCount:  req.ConversationCount,
		TotalScore:         totalScore,
		CurrentLevel:       currentLevel,
		CurrentLevelInfo:   currentLevelInfo,
		Progress:           progress,
	}
	if hasNextLevel {
		stats.NextLevelInfo = nextLevelInfo
	}

	resp = &types.BadgeDetailResponse{
		Stats:     stats,
		AllLevels: allLevels,
	}

	return resp, nil
}

// getAllBadgeLevels 获取所有勋章等级信息
func getAllBadgeLevels() []types.BadgeLevel {
	return []types.BadgeLevel{
		{
			Level:       1,
			Title:       "小小探索家",
			MinScore:    0,
			Icon:        "🌱",
			Color:       "#90EE90", // 浅绿色
			Description: "刚刚开始探索之旅",
		},
		{
			Level:       2,
			Title:       "小小专家",
			MinScore:    50,
			Icon:        "🌿",
			Color:       "#98FB98", // 淡绿色
			Description: "已经掌握了一些知识",
		},
		{
			Level:       3,
			Title:       "自然大师",
			MinScore:    150,
			Icon:        "🌳",
			Color:       "#7CFC00", // 草绿色
			Description: "对自然世界有了深入了解",
		},
		{
			Level:       4,
			Title:       "知识达人",
			MinScore:    300,
			Icon:        "🌟",
			Color:       "#32CD32", // 酸橙绿
			Description: "积累了丰富的知识",
		},
		{
			Level:       5,
			Title:       "探索之星",
			MinScore:    500,
			Icon:        "⭐",
			Color:       "#00FF00", // 纯绿色
			Description: "探索的热情如星星般闪耀",
		},
		{
			Level:       6,
			Title:       "智慧学者",
			MinScore:    750,
			Icon:        "✨",
			Color:       "#00CD00", // 深绿色
			Description: "拥有智慧的学者",
		},
		{
			Level:       7,
			Title:       "博学大师",
			MinScore:    1050,
			Icon:        "🎓",
			Color:       "#228B22", // 森林绿
			Description: "博学多才的大师",
		},
		{
			Level:       8,
			Title:       "知识巨匠",
			MinScore:    1400,
			Icon:        "👑",
			Color:       "#006400", // 深绿色
			Description: "知识的巨匠",
		},
		{
			Level:       9,
			Title:       "探索传奇",
			MinScore:    1800,
			Icon:        "🏆",
			Color:       "#004D00", // 极深绿色
			Description: "探索世界的传奇",
		},
		{
			Level:       10,
			Title:       "终极探索者",
			MinScore:    2250,
			Icon:        "💎",
			Color:       "#003300", // 最深绿色
			Description: "探索世界的终极大师",
		},
	}
}

// calculateLevel 根据总分计算等级
func calculateLevel(totalScore int, levels []types.BadgeLevel) int {
	// 从最高等级开始检查
	for i := len(levels) - 1; i >= 0; i-- {
		if totalScore >= levels[i].MinScore {
			return levels[i].Level
		}
	}
	// 如果都不满足，返回最低等级
	return 1
}

// calculateProgress 计算当前等级进度（0-100）
func calculateProgress(totalScore int, currentLevel types.BadgeLevel, nextLevel types.BadgeLevel, hasNextLevel bool) int {
	if !hasNextLevel {
		// 已经是最高等级
		return 100
	}

	currentScore := totalScore - currentLevel.MinScore
	nextScore := nextLevel.MinScore - currentLevel.MinScore

	if nextScore <= 0 {
		return 100
	}

	progress := (currentScore * 100) / nextScore
	if progress > 100 {
		progress = 100
	}
	if progress < 0 {
		progress = 0
	}

	return progress
}
