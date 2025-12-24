package nodes

import (
	"context"
	"testing"

	"github.com/tango/explore/internal/config"
	"github.com/tango/explore/internal/tools"
	"github.com/zeromicro/go-zero/core/logx"
)

func TestScienceAgentNode_GenerateScienceAnswer(t *testing.T) {
	ctx := context.Background()
	logger := logx.WithContext(ctx)
	cfg := config.AIConfig{
		EinoBaseURL: "",
		AppID:       "",
		AppKey:      "",
	}

	toolRegistry := tools.GetDefaultRegistry(logger)
	node, err := NewScienceAgentNode(ctx, cfg, logger, toolRegistry)
	if err != nil {
		t.Fatalf("Failed to create ScienceAgentNode: %v", err)
	}

	response, err := node.GenerateScienceAnswer(ctx, "这是什么？", "银杏", "自然类", 10, nil, 4, []string{})
	if err != nil {
		t.Errorf("GenerateScienceAnswer failed: %v", err)
		return
	}

	if response.DomainType != "Science" {
		t.Errorf("Expected DomainType Science, got %s", response.DomainType)
	}

	if response.Content == "" {
		t.Error("Content should not be empty")
	}
}

func TestLanguageAgentNode_GenerateLanguageAnswer(t *testing.T) {
	ctx := context.Background()
	logger := logx.WithContext(ctx)
	cfg := config.AIConfig{
		EinoBaseURL: "",
		AppID:       "",
		AppKey:      "",
	}

	toolRegistry := tools.GetDefaultRegistry(logger)
	node, err := NewLanguageAgentNode(ctx, cfg, logger, toolRegistry)
	if err != nil {
		t.Fatalf("Failed to create LanguageAgentNode: %v", err)
	}

	response, err := node.GenerateLanguageAnswer(ctx, "用英语怎么说？", "银杏", "自然类", 10, nil, []string{})
	if err != nil {
		t.Errorf("GenerateLanguageAnswer failed: %v", err)
		return
	}

	if response.DomainType != "Language" {
		t.Errorf("Expected DomainType Language, got %s", response.DomainType)
	}

	if response.Content == "" {
		t.Error("Content should not be empty")
	}
}

func TestHumanitiesAgentNode_GenerateHumanitiesAnswer(t *testing.T) {
	ctx := context.Background()
	logger := logx.WithContext(ctx)
	cfg := config.AIConfig{
		EinoBaseURL: "",
		AppID:       "",
		AppKey:      "",
	}

	node, err := NewHumanitiesAgentNode(ctx, cfg, logger)
	if err != nil {
		t.Fatalf("Failed to create HumanitiesAgentNode: %v", err)
	}

	response, err := node.GenerateHumanitiesAnswer(ctx, "有什么故事吗？", "银杏", "自然类", 10, nil)
	if err != nil {
		t.Errorf("GenerateHumanitiesAnswer failed: %v", err)
		return
	}

	if response.DomainType != "Humanities" {
		t.Errorf("Expected DomainType Humanities, got %s", response.DomainType)
	}

	if response.Content == "" {
		t.Error("Content should not be empty")
	}
}

