package logic

import (
	"context"
	"fmt"
	"time"

	"github.com/tango/explore/internal/svc"
	"github.com/tango/explore/internal/types"
	"github.com/tango/explore/internal/utils"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateShareLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateShareLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateShareLogic {
	return &CreateShareLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateShareLogic) CreateShare(req *types.CreateShareRequest) (resp *types.CreateShareResponse, err error) {
	// 参数验证
	if len(req.ExplorationRecords) == 0 && len(req.CollectedCards) == 0 {
		return nil, utils.NewAPIError(400, "探索记录和收藏卡片不能同时为空")
	}

	// 生成分享ID
	shareId := uuid.New().String()
	now := time.Now()
	expiresAt := now.Add(7 * 24 * time.Hour) // 7天后过期

	// 保存到内存存储
	store := GetShareStore()
	store.Save(shareId, &ShareData{
		ShareId:           shareId,
		ExplorationRecords: req.ExplorationRecords,
		CollectedCards:     req.CollectedCards,
		CreatedAt:          now,
		ExpiresAt:          expiresAt,
	})

	// 生成分享URL（这里使用相对路径，实际部署时需要配置完整URL）
	shareUrl := fmt.Sprintf("/api/share/%s", shareId)

	resp = &types.CreateShareResponse{
		ShareId:   shareId,
		ShareUrl:  shareUrl,
		ExpiresAt: expiresAt.Format(time.RFC3339),
	}

	l.Infow("创建分享链接", logx.Field("shareId", shareId), logx.Field("recordCount", len(req.ExplorationRecords)), logx.Field("cardCount", len(req.CollectedCards)))
	return resp, nil
}
