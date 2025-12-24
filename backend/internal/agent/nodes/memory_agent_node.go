package nodes

import (
	"context"
	"time"

	"github.com/tango/explore/internal/config"
	"github.com/tango/explore/internal/storage"
	"github.com/tango/explore/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

// MemoryAgentNode Memory Agent节点
type MemoryAgentNode struct {
	ctx         context.Context
	config      config.AIConfig
	logger      logx.Logger
	memoryStorage *storage.MemoryAgentStorage
}

// NewMemoryAgentNode 创建Memory Agent节点
func NewMemoryAgentNode(ctx context.Context, cfg config.AIConfig, logger logx.Logger, memoryStorage *storage.MemoryAgentStorage) (*MemoryAgentNode, error) {
	node := &MemoryAgentNode{
		ctx:          ctx,
		config:       cfg,
		logger:       logger,
		memoryStorage: memoryStorage,
	}

	logger.Info("✅ Memory Agent节点已初始化")
	return node, nil
}

// RecordMemory 记录学习状态
func (n *MemoryAgentNode) RecordMemory(ctx context.Context, sessionId string, reflectionResult *types.ReflectionResult, content string, objectName string) error {
	n.logger.Infow("执行Memory Agent记忆记录",
		logx.Field("sessionId", sessionId),
		logx.Field("interest", reflectionResult.Interest),
		logx.Field("confusion", reflectionResult.Confusion),
		logx.Field("objectName", objectName),
	)

	// 获取现有记忆记录
	record, exists := n.memoryStorage.GetMemoryRecord(sessionId)
	if !exists {
		// 创建新记录
		record = &types.MemoryRecord{
			SessionId:         sessionId,
			InterestedTopics:  []string{},
			UnderstoodPoints:  []string{},
			UnunderstoodPoints: []string{},
			UpdatedAt:         time.Now(),
		}
	}

	// 根据反思结果更新记忆
	if reflectionResult.Interest && objectName != "" {
		// 添加感兴趣的主题
		n.memoryStorage.AddInterestedTopic(sessionId, objectName)
	}

	if reflectionResult.Confusion {
		// 添加未理解的点
		n.memoryStorage.AddUnunderstoodPoint(sessionId, content)
	} else {
		// 添加已理解的点
		n.memoryStorage.AddUnderstoodPoint(sessionId, content)
	}

	n.logger.Infow("记忆记录完成",
		logx.Field("sessionId", sessionId),
		logx.Field("interestedTopics", record.InterestedTopics),
		logx.Field("understoodPoints", record.UnderstoodPoints),
		logx.Field("ununderstoodPoints", record.UnunderstoodPoints),
	)

	return nil
}

// GetMemory 获取记忆记录
func (n *MemoryAgentNode) GetMemory(ctx context.Context, sessionId string) (*types.MemoryRecord, bool) {
	return n.memoryStorage.GetMemoryRecord(sessionId)
}

