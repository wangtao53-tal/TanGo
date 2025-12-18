package logic

import (
	"context"
	"testing"

	"github.com/tango/explore/internal/svc"
	"github.com/tango/explore/internal/types"
)

func TestIdentifyLogic_Identify(t *testing.T) {
	ctx := context.Background()
	svcCtx := &svc.ServiceContext{}
	logic := NewIdentifyLogic(ctx, svcCtx)

	// 测试正常情况
	req := &types.IdentifyRequest{
		Image: "data:image/jpeg;base64,/9j/4AAQSkZJRg==",
		Age:   8,
	}

	resp, err := logic.Identify(req)
	if err != nil {
		t.Fatalf("Identify failed: %v", err)
	}

	if resp.ObjectName == "" {
		t.Error("ObjectName should not be empty")
	}
	if resp.ObjectCategory == "" {
		t.Error("ObjectCategory should not be empty")
	}
	if resp.Confidence < 0 || resp.Confidence > 1 {
		t.Errorf("Confidence should be between 0 and 1, got %f", resp.Confidence)
	}

	// 测试空图片
	req2 := &types.IdentifyRequest{
		Image: "",
	}
	_, err2 := logic.Identify(req2)
	if err2 == nil {
		t.Error("Should return error when image is empty")
	}
}

