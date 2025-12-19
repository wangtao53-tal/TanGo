package logic

import (
	"context"
	"strings"

	"github.com/tango/explore/internal/svc"
	"github.com/tango/explore/internal/types"
)

type IntentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIntentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IntentLogic {
	return &IntentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// RecognizeIntent 识别用户意图
func (l *IntentLogic) RecognizeIntent(req *types.IntentRequest) (*types.IntentResult, error) {
	message := strings.ToLower(strings.TrimSpace(req.Message))

	// 关键词规则判断（快速路径）
	generateCardKeywords := []string{
		"生成", "卡片", "小卡片", "知识卡片",
		"生成卡片", "帮我生成", "制作卡片",
		"create card", "generate card", "make card",
		"card", "cards",
	}

	for _, keyword := range generateCardKeywords {
		if strings.Contains(message, keyword) {
			return &types.IntentResult{
				Intent:     "generate_cards",
				Confidence: 0.9,
				Reason:     "关键词匹配: " + keyword,
			}, nil
		}
	}

	// 使用LLM进行意图识别（通过eino框架调用意图识别模型）
	// 当前使用规则判断作为降级方案
	// 如果消息包含问号或疑问词，倾向于文本回答
	questionWords := []string{"什么", "为什么", "怎么", "如何", "哪里", "哪个", "who", "what", "why", "how", "where", "which", "?"}
	for _, word := range questionWords {
		if strings.Contains(message, word) {
			return &types.IntentResult{
				Intent:     "text_response",
				Confidence: 0.8,
				Reason:     "疑问句识别",
			}, nil
		}
	}

	// 默认返回文本回答
	return &types.IntentResult{
		Intent:     "text_response",
		Confidence: 0.7,
		Reason:     "默认文本回答",
	}, nil
}
