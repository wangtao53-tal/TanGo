package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/tango/explore/internal/svc"
	"github.com/tango/explore/internal/types"
)

type StreamLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStreamLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StreamLogic {
	return &StreamLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// StreamResponse 流式返回响应
func (l *StreamLogic) StreamResponse(sessionId string, message string) error {
	// 获取上下文
	messages := l.svcCtx.Storage.GetMessages(sessionId)

	// 调用意图识别
	intentLogic := NewIntentLogic(l.ctx, l.svcCtx)
	intentReq := &types.IntentRequest{
		Message:   message,
		SessionId: sessionId,
		Context:   l.convertToConversationMessages(messages),
	}

	intentResult, err := intentLogic.RecognizeIntent(intentReq)
	if err != nil {
		return err
	}

	// 根据意图生成流式响应
	if intentResult.Intent == "generate_cards" {
		// 生成卡片流式返回
		return l.streamCards(sessionId, message)
	} else {
		// 文本回答流式返回
		return l.streamText(sessionId, message)
	}
}

// streamText 流式返回文本
func (l *StreamLogic) streamText(sessionId string, message string) error {
	// TODO: 实现真实的流式文本生成（通过eino框架）
	// 当前使用Mock数据
	text := "这是一个Mock流式响应。待接入真实AI模型后，将实现真实的流式文本生成。"

	// 模拟流式输出
	words := []rune(text)
	for i, word := range words {
		event := map[string]interface{}{
			"type":    "text",
			"content": string(word),
			"index":   i,
		}
		_ = event                         // 这里应该发送到SSE连接
		time.Sleep(50 * time.Millisecond) // 模拟延迟
	}

	return nil
}

// streamCards 流式返回卡片
func (l *StreamLogic) streamCards(sessionId string, message string) error {
	// TODO: 实现真实的流式卡片生成（通过eino框架）
	// 当前使用Mock数据
	cards := []map[string]interface{}{
		{"type": "science", "title": "科学认知卡", "content": "Mock内容"},
		{"type": "poetry", "title": "古诗词卡", "content": "Mock内容"},
		{"type": "english", "title": "英语表达卡", "content": "Mock内容"},
	}

	for _, card := range cards {
		event := map[string]interface{}{
			"type":    "card",
			"content": card,
		}
		_ = event                          // 这里应该发送到SSE连接
		time.Sleep(200 * time.Millisecond) // 模拟延迟
	}

	return nil
}

// convertToConversationMessages 转换消息列表
func (l *StreamLogic) convertToConversationMessages(messages []interface{}) []types.ConversationMessage {
	result := make([]types.ConversationMessage, 0)
	for _, msg := range messages {
		if convMsg, ok := msg.(types.ConversationMessage); ok {
			result = append(result, convMsg)
		}
	}
	return result
}

// SendSSEEvent 发送SSE事件（辅助函数）
func (l *StreamLogic) SendSSEEvent(eventType string, data interface{}) (string, error) {
	event := map[string]interface{}{
		"type":    eventType,
		"content": data,
	}
	jsonData, err := json.Marshal(event)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("event: %s\ndata: %s\n\n", eventType, string(jsonData)), nil
}
