package agent

import (
	"context"

	"github.com/tango/explore/internal/config"
	"github.com/zeromicro/go-zero/core/logx"
)

// Agent AI Agent系统
type Agent struct {
	config config.AIConfig
	graph  *Graph
	logger logx.Logger
	ctx    context.Context
}

// NewAgent 创建新的Agent实例
func NewAgent(ctx context.Context, cfg config.AIConfig) (*Agent, error) {
	logger := logx.WithContext(ctx)

	agent := &Agent{
		config: cfg,
		logger: logger,
		ctx:    ctx,
	}

	// 初始化Graph
	graph, err := NewGraph(ctx, cfg, logger)
	if err != nil {
		return nil, err
	}
	agent.graph = graph

	logger.Infow("Agent初始化完成",
		logx.Field("einoBaseURL", cfg.EinoBaseURL),
		logx.Field("appID", cfg.AppID),
	)

	return agent, nil
}

// GetGraph 获取Graph实例
func (a *Agent) GetGraph() *Graph {
	return a.graph
}

// Close 关闭Agent，清理资源
func (a *Agent) Close() error {
	// TODO: 清理eino相关资源
	return nil
}
