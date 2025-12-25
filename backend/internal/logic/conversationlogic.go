package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tango/explore/internal/svc"
	"github.com/tango/explore/internal/types"
	"github.com/tango/explore/internal/utils"
	"github.com/zeromicro/go-zero/core/logx"
)

type ConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConversationLogic {
	return &ConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ProcessMessage 处理对话消息
func (l *ConversationLogic) ProcessMessage(req *types.ConversationRequest) (*types.ConversationResponse, error) {
	// 参数验证
	if req.Message == "" && req.Image == "" && req.Voice == "" {
		return nil, fmt.Errorf("消息内容不能为空")
	}
	// 生成或使用现有会话ID
	sessionId := req.SessionId
	if sessionId == "" {
		sessionId = uuid.New().String()
	}

	// 如果提供了识别结果上下文，保存到会话数据中
	if req.IdentificationContext != nil {
		l.svcCtx.Storage.SetData(sessionId, "identificationContext", req.IdentificationContext)
	}

	// 创建用户消息
	userMessage := types.ConversationMessage{
		Id:        uuid.New().String(),
		Type:      "text",
		Sender:    "user",
		Content:   req.Message,
		Timestamp: time.Now().Format(time.RFC3339),
		SessionId: sessionId,
	}

	// 如果有图片，添加到消息内容
	if req.Image != "" {
		userMessage.Type = "image"
		content := map[string]interface{}{
			"text":  req.Message,
			"image": req.Image,
		}
		contentBytes, _ := json.Marshal(content)
		userMessage.Content = string(contentBytes)
	}

	// 如果有语音，添加到消息内容
	if req.Voice != "" {
		userMessage.Type = "voice"
		content := map[string]interface{}{
			"text":  req.Message,
			"voice": req.Voice,
		}
		contentBytes, _ := json.Marshal(content)
		userMessage.Content = string(contentBytes)
	}

	// 保存用户消息到存储
	l.svcCtx.Storage.AddMessage(sessionId, userMessage)

	// 调用意图识别
	intentLogic := NewIntentLogic(l.ctx, l.svcCtx)
	intentReq := &types.IntentRequest{
		Message:   req.Message,
		SessionId: sessionId,
		Context:   l.getContextMessages(sessionId),
	}

	intentResult, err := intentLogic.RecognizeIntent(intentReq)
	if err != nil {
		logx.Errorf("意图识别失败: %v", err)
		return nil, fmt.Errorf("意图识别失败: %w", err)
	}

	// 根据意图生成响应
	var assistantMessage types.ConversationMessage
	var responseType string

	if intentResult.Intent == "generate_cards" {
		// 生成卡片逻辑（这里需要调用卡片生成逻辑）
		// 暂时返回文本响应
		assistantMessage = types.ConversationMessage{
			Id:        uuid.New().String(),
			Type:      "card",
			Sender:    "assistant",
			Content:   "正在为您生成知识卡片...",
			Timestamp: time.Now().Format(time.RFC3339),
			SessionId: sessionId,
		}
		responseType = "cards"
	} else {
		// 文本回答
		textContent := l.generateTextResponse(req.Message, sessionId)
		isMarkdown := utils.DetectMarkdown(textContent)
		markdownPtr := &isMarkdown

		assistantMessage = types.ConversationMessage{
			Id:        uuid.New().String(),
			Type:      "text",
			Sender:    "assistant",
			Content:   textContent,
			Timestamp: time.Now().Format(time.RFC3339),
			SessionId: sessionId,
			Markdown:  markdownPtr,
		}
		responseType = "text"
	}

	// 保存助手消息
	l.svcCtx.Storage.AddMessage(sessionId, assistantMessage)

	logx.Infow("对话消息处理成功",
		logx.Field("sessionId", sessionId),
		logx.Field("intent", intentResult.Intent),
		logx.Field("responseType", responseType),
	)

	return &types.ConversationResponse{
		Message:   assistantMessage,
		SessionId: sessionId,
		Type:      responseType,
	}, nil
}

// getContextMessages 获取上下文消息（最多10轮）
func (l *ConversationLogic) getContextMessages(sessionId string) []types.ConversationMessage {
	messages := l.svcCtx.Storage.GetMessages(sessionId)
	contextMessages := make([]types.ConversationMessage, 0)

	// 只取最后10轮对话（20条消息）
	start := 0
	if len(messages) > 20 {
		start = len(messages) - 20
	}

	for i := start; i < len(messages); i++ {
		if msg, ok := messages[i].(types.ConversationMessage); ok {
			contextMessages = append(contextMessages, msg)
		}
	}

	return contextMessages
}

// generateTextResponse 生成文本响应（Mock实现，待接入真实AI模型）
func (l *ConversationLogic) generateTextResponse(message string, sessionId string) string {
	// 获取识别结果上下文（如果存在）
	var identificationContext *types.IdentificationContext
	if ctxData, ok := l.svcCtx.Storage.GetData(sessionId, "identificationContext"); ok {
		if ctx, ok := ctxData.(*types.IdentificationContext); ok {
			identificationContext = ctx
		}
	}

	// TODO: 接入真实AI模型（通过eino框架）
	// 当前使用Mock数据
	// 如果存在识别结果上下文，可以在提示中包含识别结果信息
	if identificationContext != nil {
		return "这是一个Mock响应。待接入真实AI模型后，将根据您的问题和识别结果（" + identificationContext.ObjectName + "）生成相应的回答。"
	}
	return "这是一个Mock响应。待接入真实AI模型后，将根据您的问题生成相应的回答。"
}
