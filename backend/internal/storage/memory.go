package storage

import (
	"sync"
	"time"
)

// SessionData 会话数据
type SessionData struct {
	SessionId  string
	Messages   []interface{}
	CreatedAt  time.Time
	LastActive time.Time
	Data       map[string]interface{}
}

// MemoryStorage 内存缓存实现
type MemoryStorage struct {
	sessions sync.Map // key: sessionId, value: *SessionData
	mu       sync.RWMutex
}

// NewMemoryStorage 创建新的内存存储实例
func NewMemoryStorage() *MemoryStorage {
	storage := &MemoryStorage{}
	// 启动清理协程
	go storage.startCleanup()
	return storage
}

// SetSession 设置会话数据
func (m *MemoryStorage) SetSession(sessionId string, data *SessionData) {
	data.LastActive = time.Now()
	m.sessions.Store(sessionId, data)
}

// GetSession 获取会话数据
func (m *MemoryStorage) GetSession(sessionId string) (*SessionData, bool) {
	value, ok := m.sessions.Load(sessionId)
	if !ok {
		return nil, false
	}
	session := value.(*SessionData)
	// 更新最后活动时间
	session.LastActive = time.Now()
	return session, true
}

// DeleteSession 删除会话
func (m *MemoryStorage) DeleteSession(sessionId string) {
	m.sessions.Delete(sessionId)
}

// AddMessage 添加消息到会话
func (m *MemoryStorage) AddMessage(sessionId string, message interface{}) {
	value, ok := m.sessions.Load(sessionId)
	if !ok {
		// 创建新会话
		session := &SessionData{
			SessionId:  sessionId,
			Messages:   []interface{}{message},
			CreatedAt:  time.Now(),
			LastActive: time.Now(),
			Data:       make(map[string]interface{}),
		}
		m.sessions.Store(sessionId, session)
		return
	}

	session := value.(*SessionData)
	session.Messages = append(session.Messages, message)
	// 限制上下文长度，最多保留20轮对话（40条消息）
	if len(session.Messages) > 40 {
		session.Messages = session.Messages[len(session.Messages)-40:]
	}
	session.LastActive = time.Now()
	m.sessions.Store(sessionId, session)
}

// GetMessages 获取会话消息列表
func (m *MemoryStorage) GetMessages(sessionId string) []interface{} {
	value, ok := m.sessions.Load(sessionId)
	if !ok {
		return []interface{}{}
	}
	session := value.(*SessionData)
	return session.Messages
}

// SetData 设置会话的额外数据
func (m *MemoryStorage) SetData(sessionId string, key string, value interface{}) {
	valueData, ok := m.sessions.Load(sessionId)
	if !ok {
		session := &SessionData{
			SessionId:  sessionId,
			Messages:   []interface{}{},
			CreatedAt:  time.Now(),
			LastActive: time.Now(),
			Data:       map[string]interface{}{key: value},
		}
		m.sessions.Store(sessionId, session)
		return
	}

	session := valueData.(*SessionData)
	if session.Data == nil {
		session.Data = make(map[string]interface{})
	}
	session.Data[key] = value
	session.LastActive = time.Now()
	m.sessions.Store(sessionId, session)
}

// GetData 获取会话的额外数据
func (m *MemoryStorage) GetData(sessionId string, key string) (interface{}, bool) {
	value, ok := m.sessions.Load(sessionId)
	if !ok {
		return nil, false
	}
	session := value.(*SessionData)
	if session.Data == nil {
		return nil, false
	}
	val, ok := session.Data[key]
	return val, ok
}

// startCleanup 启动清理协程，定期清理过期会话（30分钟无活动）
func (m *MemoryStorage) startCleanup() {
	ticker := time.NewTicker(5 * time.Minute) // 每5分钟检查一次
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		m.sessions.Range(func(key, value interface{}) bool {
			session := value.(*SessionData)
			// 如果30分钟无活动，删除会话
			if now.Sub(session.LastActive) > 30*time.Minute {
				m.sessions.Delete(key)
			}
			return true
		})
	}
}
