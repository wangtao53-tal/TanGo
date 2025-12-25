package storage

import (
	"sync"
	"time"

	"github.com/tango/explore/internal/types"
)

// MemoryAgentStorage Memory Agent存储实现
type MemoryAgentStorage struct {
	records sync.Map // key: sessionId, value: *types.MemoryRecord
	mu      sync.RWMutex
}

// NewMemoryAgentStorage 创建新的Memory Agent存储实例
func NewMemoryAgentStorage() *MemoryAgentStorage {
	return &MemoryAgentStorage{}
}

// GetMemoryRecord 获取记忆记录
func (m *MemoryAgentStorage) GetMemoryRecord(sessionId string) (*types.MemoryRecord, bool) {
	value, ok := m.records.Load(sessionId)
	if !ok {
		return nil, false
	}
	return value.(*types.MemoryRecord), true
}

// SetMemoryRecord 设置记忆记录
func (m *MemoryAgentStorage) SetMemoryRecord(sessionId string, record *types.MemoryRecord) {
	record.UpdatedAt = time.Now()
	m.records.Store(sessionId, record)
}

// UpdateMemoryRecord 更新记忆记录
func (m *MemoryAgentStorage) UpdateMemoryRecord(sessionId string, record *types.MemoryRecord) {
	record.UpdatedAt = time.Now()
	m.records.Store(sessionId, record)
}

// AddInterestedTopic 添加感兴趣的主题
func (m *MemoryAgentStorage) AddInterestedTopic(sessionId string, topic string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	record, ok := m.GetMemoryRecord(sessionId)
	if !ok {
		record = &types.MemoryRecord{
			SessionId:        sessionId,
			InterestedTopics: []string{},
			UnderstoodPoints: []string{},
			UnunderstoodPoints: []string{},
			UpdatedAt:        time.Now(),
		}
	}

	// 检查是否已存在
	for _, t := range record.InterestedTopics {
		if t == topic {
			return
		}
	}

	record.InterestedTopics = append(record.InterestedTopics, topic)
	m.SetMemoryRecord(sessionId, record)
}

// AddUnderstoodPoint 添加已理解的点
func (m *MemoryAgentStorage) AddUnderstoodPoint(sessionId string, point string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	record, ok := m.GetMemoryRecord(sessionId)
	if !ok {
		record = &types.MemoryRecord{
			SessionId:        sessionId,
			InterestedTopics: []string{},
			UnderstoodPoints: []string{},
			UnunderstoodPoints: []string{},
			UpdatedAt:        time.Now(),
		}
	}

	// 检查是否已存在
	for _, p := range record.UnderstoodPoints {
		if p == point {
			return
		}
	}

	record.UnderstoodPoints = append(record.UnderstoodPoints, point)
	m.SetMemoryRecord(sessionId, record)
}

// AddUnunderstoodPoint 添加未理解的点
func (m *MemoryAgentStorage) AddUnunderstoodPoint(sessionId string, point string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	record, ok := m.GetMemoryRecord(sessionId)
	if !ok {
		record = &types.MemoryRecord{
			SessionId:        sessionId,
			InterestedTopics: []string{},
			UnderstoodPoints: []string{},
			UnunderstoodPoints: []string{},
			UpdatedAt:        time.Now(),
		}
	}

	// 检查是否已存在
	for _, p := range record.UnunderstoodPoints {
		if p == point {
			return
		}
	}

	record.UnunderstoodPoints = append(record.UnunderstoodPoints, point)
	m.SetMemoryRecord(sessionId, record)
}

// DeleteMemoryRecord 删除记忆记录
func (m *MemoryAgentStorage) DeleteMemoryRecord(sessionId string) {
	m.records.Delete(sessionId)
}

