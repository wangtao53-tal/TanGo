package svc

import (
	"context"
	"strings"

	"github.com/tango/explore/internal/agent"
	"github.com/tango/explore/internal/config"
	"github.com/tango/explore/internal/storage"
	"github.com/zeromicro/go-zero/core/logx"
)

type ServiceContext struct {
	Config        config.Config
	Storage       *storage.MemoryStorage
	Agent         *agent.Agent
	GitHubStorage *storage.GitHubStorage
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

	// 初始化 GitHub 存储
	var githubStorage *storage.GitHubStorage
	if c.Upload.GitHubToken != "" && c.Upload.GitHubOwner != "" && c.Upload.GitHubRepo != "" {
		// 检查token是否是占位符
		isPlaceholder := strings.Contains(c.Upload.GitHubToken, "xxxxx") || strings.HasPrefix(c.Upload.GitHubToken, "ghp_xxxxxxxx")
		if isPlaceholder {
			logger.Errorw("GitHub token 是占位符，GitHub 存储将无法正常工作",
				logx.Field("hint", "请将 .env 文件中的 GITHUB_TOKEN 替换为真实的 token"),
				logx.Field("owner", c.Upload.GitHubOwner),
				logx.Field("repo", c.Upload.GitHubRepo),
			)
			logger.Info("⚠️  警告：由于 GitHub token 是占位符，图片上传将降级到 base64 方案")
			// 仍然创建实例，但在Upload时会返回错误并降级到base64
		}
		githubStorage = storage.NewGitHubStorage(c.Upload, logger)
		if !isPlaceholder {
			logger.Info("GitHub 存储初始化成功",
				logx.Field("owner", c.Upload.GitHubOwner),
				logx.Field("repo", c.Upload.GitHubRepo),
				logx.Field("branch", c.Upload.GitHubBranch),
				logx.Field("path", c.Upload.GitHubPath),
			)
		}
	} else {
		logger.Infow("未配置 GitHub 上传参数，图片上传将使用 base64 降级方案",
			logx.Field("hasToken", c.Upload.GitHubToken != ""),
			logx.Field("hasOwner", c.Upload.GitHubOwner != ""),
			logx.Field("hasRepo", c.Upload.GitHubRepo != ""),
		)
		logger.Info("如需使用 GitHub 存储，请在.env文件中配置：GITHUB_TOKEN、GITHUB_OWNER、GITHUB_REPO")
	}

	return &ServiceContext{
		Config:        c,
		Storage:       storage.NewMemoryStorage(),
		Agent:         aiAgent,
		GitHubStorage: githubStorage,
	}
}
