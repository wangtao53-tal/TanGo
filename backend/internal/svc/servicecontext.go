package svc

import (
	"context"

	"github.com/tango/explore/internal/agent"
	"github.com/tango/explore/internal/config"
	"github.com/tango/explore/internal/storage"
	"github.com/zeromicro/go-zero/core/logx"
)

type ServiceContext struct {
	Config  config.Config
	Storage *storage.MemoryStorage
	Agent   *agent.Agent
}

func NewServiceContext(c config.Config) *ServiceContext {
	ctx := context.Background()
	logger := logx.WithContext(ctx)

	// 初始化Agent系统
	var aiAgent *agent.Agent
	var err error
	if c.AI.EinoBaseURL != "" || c.AI.AppID != "" {
		// 如果配置了eino相关配置，初始化Agent
		aiAgent, err = agent.NewAgent(ctx, c.AI)
		if err != nil {
			logger.Errorw("Agent初始化失败", logx.Field("error", err))
			// 继续运行，使用Mock数据
		} else {
			logger.Info("Agent系统初始化成功")
		}
	} else {
		logger.Info("未配置eino，将使用Mock数据")
	}

	return &ServiceContext{
		Config:  c,
		Storage: storage.NewMemoryStorage(),
		Agent:   aiAgent,
	}
}
