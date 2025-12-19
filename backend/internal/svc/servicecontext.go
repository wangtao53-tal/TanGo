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

	// 检查eino配置
	hasEinoBaseURL := c.AI.EinoBaseURL != ""
	hasAppID := c.AI.AppID != ""
	hasAppKey := c.AI.AppKey != ""

	logger.Infow("检查eino配置",
		logx.Field("hasEinoBaseURL", hasEinoBaseURL),
		logx.Field("hasAppID", hasAppID),
		logx.Field("hasAppKey", hasAppKey),
	)

	if hasEinoBaseURL || hasAppID {
		// 如果配置了eino相关配置，初始化Agent
		aiAgent, err = agent.NewAgent(ctx, c.AI)
		if err != nil {
			logger.Errorw("Agent初始化失败，将使用Mock数据",
				logx.Field("error", err),
				logx.Field("errorDetail", err.Error()),
			)
			// 继续运行，使用Mock数据
		} else {
			logger.Info("Agent系统初始化成功，将使用真实模型")
		}
	} else {
		logger.Errorw("未配置eino参数（EINO_BASE_URL或TAL_MLOPS_APP_ID），将使用Mock数据")
		logger.Info("如需使用真实模型，请在.env文件中配置：EINO_BASE_URL、TAL_MLOPS_APP_ID、TAL_MLOPS_APP_KEY")
	}

	return &ServiceContext{
		Config:  c,
		Storage: storage.NewMemoryStorage(),
		Agent:   aiAgent,
	}
}
