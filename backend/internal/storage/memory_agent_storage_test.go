package storage

import (
	"testing"
	"time"

	"github.com/tango/explore/internal/types"
)

func TestMemoryAgentStorage_GetSetMemoryRecord(t *testing.T) {
	storage := NewMemoryAgentStorage()
	sessionId := "test-session-123"

	// 测试设置记忆记录
	record := &types.MemoryRecord{
		SessionId:         sessionId,
		InterestedTopics:  []string{"银杏", "自然"},
		UnderstoodPoints:  []string{"银杏是植物"},
		UnunderstoodPoints: []string{},
		UpdatedAt:         time.Now(),
	}

	storage.SetMemoryRecord(sessionId, record)

	// 测试获取记忆记录
	retrievedRecord, exists := storage.GetMemoryRecord(sessionId)
	if !exists {
		t.Error("Memory record should exist")
		return
	}

	if retrievedRecord.SessionId != sessionId {
		t.Errorf("Expected sessionId %s, got %s", sessionId, retrievedRecord.SessionId)
	}

	if len(retrievedRecord.InterestedTopics) != 2 {
		t.Errorf("Expected 2 interested topics, got %d", len(retrievedRecord.InterestedTopics))
	}
}

func TestMemoryAgentStorage_AddInterestedTopic(t *testing.T) {
	storage := NewMemoryAgentStorage()
	sessionId := "test-session-123"

	// 添加感兴趣的主题
	storage.AddInterestedTopic(sessionId, "银杏")
	storage.AddInterestedTopic(sessionId, "自然")

	record, exists := storage.GetMemoryRecord(sessionId)
	if !exists {
		t.Error("Memory record should exist after adding topics")
		return
	}

	if len(record.InterestedTopics) != 2 {
		t.Errorf("Expected 2 interested topics, got %d", len(record.InterestedTopics))
	}

	// 测试重复添加（不应该重复）
	storage.AddInterestedTopic(sessionId, "银杏")
	record2, _ := storage.GetMemoryRecord(sessionId)
	if len(record2.InterestedTopics) != 2 {
		t.Errorf("Should not add duplicate topic, expected 2, got %d", len(record2.InterestedTopics))
	}
}

func TestMemoryAgentStorage_AddUnderstoodPoint(t *testing.T) {
	storage := NewMemoryAgentStorage()
	sessionId := "test-session-123"

	storage.AddUnderstoodPoint(sessionId, "银杏是植物")
	storage.AddUnderstoodPoint(sessionId, "银杏有叶子")

	record, exists := storage.GetMemoryRecord(sessionId)
	if !exists {
		t.Error("Memory record should exist")
		return
	}

	if len(record.UnderstoodPoints) != 2 {
		t.Errorf("Expected 2 understood points, got %d", len(record.UnderstoodPoints))
	}
}

func TestMemoryAgentStorage_AddUnunderstoodPoint(t *testing.T) {
	storage := NewMemoryAgentStorage()
	sessionId := "test-session-123"

	storage.AddUnunderstoodPoint(sessionId, "为什么银杏会变黄？")

	record, exists := storage.GetMemoryRecord(sessionId)
	if !exists {
		t.Error("Memory record should exist")
		return
	}

	if len(record.UnunderstoodPoints) != 1 {
		t.Errorf("Expected 1 ununderstood point, got %d", len(record.UnunderstoodPoints))
	}
}

