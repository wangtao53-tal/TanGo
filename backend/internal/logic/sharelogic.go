package logic

import (
	"sync"
	"time"

	"github.com/tango/explore/internal/types"
)

// ShareStore 分享链接存储（内存实现）
type ShareStore struct {
	mu    sync.RWMutex
	shares map[string]*ShareData
}

// ShareData 分享数据
type ShareData struct {
	ShareId           string
	ExplorationRecords []types.ExplorationRecord
	CollectedCards     []types.KnowledgeCard
	CreatedAt          time.Time
	ExpiresAt          time.Time
}

var (
	shareStore *ShareStore
	once       sync.Once
)

// GetShareStore 获取分享存储单例
func GetShareStore() *ShareStore {
	once.Do(func() {
		shareStore = &ShareStore{
			shares: make(map[string]*ShareData),
		}
		// 启动清理goroutine
		go shareStore.cleanup()
	})
	return shareStore
}

// Save 保存分享数据
func (s *ShareStore) Save(shareId string, data *ShareData) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.shares[shareId] = data
}

// Get 获取分享数据
func (s *ShareStore) Get(shareId string) (*ShareData, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	data, ok := s.shares[shareId]
	if !ok {
		return nil, false
	}
	// 检查是否过期
	if time.Now().After(data.ExpiresAt) {
		return nil, false
	}
	return data, true
}

// Delete 删除分享数据
func (s *ShareStore) Delete(shareId string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.shares, shareId)
}

// cleanup 定期清理过期的分享链接
func (s *ShareStore) cleanup() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for shareId, data := range s.shares {
			if now.After(data.ExpiresAt) {
				delete(s.shares, shareId)
			}
		}
		s.mu.Unlock()
	}
}

