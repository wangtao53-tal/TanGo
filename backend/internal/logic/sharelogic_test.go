package logic

import (
	"testing"
	"time"

	"github.com/tango/explore/internal/types"
)

func TestShareStore(t *testing.T) {
	store := GetShareStore()

	// 测试保存和获取
	shareId := "test-share-id"
	data := &ShareData{
		ShareId:           shareId,
		ExplorationRecords: []types.ExplorationRecord{},
		CollectedCards:     []types.KnowledgeCard{},
		CreatedAt:          time.Now(),
		ExpiresAt:          time.Now().Add(7 * 24 * time.Hour),
	}

	store.Save(shareId, data)

	retrieved, ok := store.Get(shareId)
	if !ok {
		t.Fatal("Should be able to get saved share data")
	}

	if retrieved.ShareId != shareId {
		t.Errorf("ShareId mismatch: expected %s, got %s", shareId, retrieved.ShareId)
	}

	// 测试不存在的shareId
	_, ok = store.Get("non-existent-id")
	if ok {
		t.Error("Should return false for non-existent shareId")
	}

	// 测试过期数据
	expiredId := "expired-id"
	expiredData := &ShareData{
		ShareId:           expiredId,
		ExplorationRecords: []types.ExplorationRecord{},
		CollectedCards:     []types.KnowledgeCard{},
		CreatedAt:          time.Now().Add(-8 * 24 * time.Hour),
		ExpiresAt:          time.Now().Add(-1 * 24 * time.Hour), // 已过期
	}
	store.Save(expiredId, expiredData)

	_, ok = store.Get(expiredId)
	if ok {
		t.Error("Should return false for expired share data")
	}

	// 测试删除
	store.Delete(shareId)
	_, ok = store.Get(shareId)
	if ok {
		t.Error("Should return false after deletion")
	}
}

