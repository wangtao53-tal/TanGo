package logic

import (
	"context"
	"testing"

	"github.com/tango/explore/internal/svc"
	"github.com/tango/explore/internal/types"
)

func TestGenerateCardsLogic_GenerateCards(t *testing.T) {
	ctx := context.Background()
	svcCtx := &svc.ServiceContext{}
	logic := NewGenerateCardsLogic(ctx, svcCtx)

	// 测试正常情况
	req := &types.GenerateCardsRequest{
		ObjectName:     "银杏",
		ObjectCategory: "自然类",
		Age:            8,
	}

	resp, err := logic.GenerateCards(req)
	if err != nil {
		t.Fatalf("GenerateCards failed: %v", err)
	}

	if len(resp.Cards) != 3 {
		t.Errorf("Should generate 3 cards, got %d", len(resp.Cards))
	}

	// 验证卡片类型
	cardTypes := make(map[string]bool)
	for _, card := range resp.Cards {
		cardTypes[card.Type] = true
		if card.Title == "" {
			t.Error("Card title should not be empty")
		}
		if card.Content == nil {
			t.Error("Card content should not be nil")
		}
	}

	expectedTypes := []string{"science", "poetry", "english"}
	for _, expectedType := range expectedTypes {
		if !cardTypes[expectedType] {
			t.Errorf("Should have card type: %s", expectedType)
		}
	}

	// 测试参数验证
	req2 := &types.GenerateCardsRequest{
		ObjectName:     "",
		ObjectCategory: "自然类",
		Age:            8,
	}
	_, err2 := logic.GenerateCards(req2)
	if err2 == nil {
		t.Error("Should return error when ObjectName is empty")
	}

	req3 := &types.GenerateCardsRequest{
		ObjectName:     "银杏",
		ObjectCategory: "",
		Age:            8,
	}
	_, err3 := logic.GenerateCards(req3)
	if err3 == nil {
		t.Error("Should return error when ObjectCategory is empty")
	}

	req4 := &types.GenerateCardsRequest{
		ObjectName:     "银杏",
		ObjectCategory: "自然类",
		Age:            2, // 无效年龄
	}
	_, err4 := logic.GenerateCards(req4)
	if err4 == nil {
		t.Error("Should return error when age is invalid")
	}
}

