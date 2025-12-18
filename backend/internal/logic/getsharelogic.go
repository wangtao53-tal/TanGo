package logic

import (
	"context"

	"github.com/tango/explore/internal/svc"
	"github.com/tango/explore/internal/types"
	"github.com/tango/explore/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetShareLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetShareLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetShareLogic {
	return &GetShareLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetShareLogic) GetShare(shareId string) (resp *types.GetShareResponse, err error) {
	// 从内存存储获取分享数据
	store := GetShareStore()
	data, ok := store.Get(shareId)
	if !ok {
		return nil, utils.ErrShareNotFound
	}

	resp = &types.GetShareResponse{
		ExplorationRecords: data.ExplorationRecords,
		CollectedCards:     data.CollectedCards,
		CreatedAt:          data.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		ExpiresAt:          data.ExpiresAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	l.Infow("获取分享数据", logx.Field("shareId", shareId))
	return resp, nil
}
