package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
	"github.com/tango/explore/internal/agent"
	"github.com/tango/explore/internal/svc"
	"github.com/tango/explore/internal/types"
	"github.com/tango/explore/internal/utils"
	"github.com/zeromicro/go-zero/core/logx"
)

type AgentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAgentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AgentLogic {
	return &AgentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// StreamAgentConversation 多Agent模式流式对话
func (l *AgentLogic) StreamAgentConversation(
	w http.ResponseWriter,
	req types.UnifiedStreamConversationRequest,
) error {
	logger := logx.WithContext(l.ctx)

	// 设置SSE响应头
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Cache-Control")

	// 生成或使用现有会话ID
	sessionId := req.SessionId
	if sessionId == "" {
		sessionId = uuid.New().String()
	}

	// 发送连接成功事件
	connectedEvent := types.StreamEvent{
		Type:      "connected",
		Content:   map[string]interface{}{"sessionId": sessionId},
		SessionId: sessionId,
	}
	connectedJSON, _ := json.Marshal(connectedEvent)
	fmt.Fprintf(w, "event: connected\ndata: %s\n\n", string(connectedJSON))
	w.(http.Flusher).Flush()

	// 设置最大上下文轮次
	maxContextRounds := req.MaxContextRounds
	if maxContextRounds <= 0 {
		maxContextRounds = 20 // 默认20轮
	}

	// 获取用户年龄和对象信息
	userAge := req.UserAge

	if req.IdentificationContext != nil {
		if userAge == 0 && req.IdentificationContext.Age > 0 {
			userAge = req.IdentificationContext.Age
		}
		l.svcCtx.Storage.SetData(sessionId, "identificationContext", req.IdentificationContext)
	} else {
		if ctxData, ok := l.svcCtx.Storage.GetData(sessionId, "identificationContext"); ok {
			if ctx, ok := ctxData.(*types.IdentificationContext); ok {
				if userAge == 0 && ctx.Age > 0 {
					userAge = ctx.Age
				}
			}
		}
	}

	// 默认年龄
	if userAge == 0 {
		userAge = 8
	}

	// 验证messageType字段
	if req.MessageType == "" {
		logger.Errorw("messageType字段必填", logx.Field("sessionId", sessionId))
		errorEvent := types.StreamEvent{
			Type:      "error",
			Content:   map[string]interface{}{"message": "messageType字段必填"},
			SessionId: sessionId,
		}
		errorJSON, _ := json.Marshal(errorEvent)
		fmt.Fprintf(w, "event: error\ndata: %s\n\n", string(errorJSON))
		w.(http.Flusher).Flush()
		return fmt.Errorf("messageType字段必填")
	}

	// 根据messageType处理不同输入类型
	var messageText string
	var messageType string

	switch req.MessageType {
	case "text":
		if req.Message == "" {
			logger.Errorw("messageType为text时，message字段必填", logx.Field("sessionId", sessionId))
			errorEvent := types.StreamEvent{
				Type:      "error",
				Content:   map[string]interface{}{"message": "messageType为text时，message字段必填"},
				SessionId: sessionId,
			}
			errorJSON, _ := json.Marshal(errorEvent)
			fmt.Fprintf(w, "event: error\ndata: %s\n\n", string(errorJSON))
			w.(http.Flusher).Flush()
			return fmt.Errorf("messageType为text时，message字段必填")
		}
		messageText = req.Message
		messageType = "text"
	case "voice":
		if req.Audio == "" {
			logger.Errorw("messageType为voice时，audio字段必填", logx.Field("sessionId", sessionId))
			errorEvent := types.StreamEvent{
				Type:      "error",
				Content:   map[string]interface{}{"message": "messageType为voice时，audio字段必填"},
				SessionId: sessionId,
			}
			errorJSON, _ := json.Marshal(errorEvent)
			fmt.Fprintf(w, "event: error\ndata: %s\n\n", string(errorJSON))
			w.(http.Flusher).Flush()
			return fmt.Errorf("messageType为voice时，audio字段必填")
		}
		// 语音识别
		voiceLogic := NewVoiceLogic(l.ctx, l.svcCtx)
		voiceReq := &types.VoiceRequest{
			Audio:     req.Audio,
			SessionId: sessionId,
		}
		voiceResp, err := voiceLogic.RecognizeVoice(voiceReq)
		if err != nil {
			logger.Errorw("语音识别失败", logx.Field("error", err), logx.Field("sessionId", sessionId))
			errorEvent := types.StreamEvent{
				Type:      "error",
				Content:   map[string]interface{}{"message": "语音识别失败: " + err.Error()},
				SessionId: sessionId,
			}
			errorJSON, _ := json.Marshal(errorEvent)
			fmt.Fprintf(w, "event: error\ndata: %s\n\n", string(errorJSON))
			w.(http.Flusher).Flush()
			return fmt.Errorf("语音识别失败: %w", err)
		}
		messageText = voiceResp.Text
		messageType = "voice"
		recognizedEvent := types.StreamEvent{
			Type:      "voice_recognized",
			Content:   map[string]interface{}{"text": voiceResp.Text},
			SessionId: sessionId,
		}
		recognizedJSON, _ := json.Marshal(recognizedEvent)
		fmt.Fprintf(w, "event: voice_recognized\ndata: %s\n\n", string(recognizedJSON))
		w.(http.Flusher).Flush()
	case "image":
		if req.Image == "" {
			logger.Errorw("messageType为image时，image字段必填", logx.Field("sessionId", sessionId))
			errorEvent := types.StreamEvent{
				Type:      "error",
				Content:   map[string]interface{}{"message": "messageType为image时，image字段必填"},
				SessionId: sessionId,
			}
			errorJSON, _ := json.Marshal(errorEvent)
			fmt.Fprintf(w, "event: error\ndata: %s\n\n", string(errorJSON))
			w.(http.Flusher).Flush()
			return fmt.Errorf("messageType为image时，image字段必填")
		}
		// 图片上传
		if strings.HasPrefix(req.Image, "data:") || !strings.HasPrefix(req.Image, "http") {
			imageData := req.Image
			if strings.HasPrefix(imageData, "data:image/") {
				parts := strings.SplitN(imageData, ",", 2)
				if len(parts) == 2 {
					imageData = parts[1]
				}
			}
			imageData = utils.CleanBase64String(imageData)
			uploadLogic := NewUploadLogic(l.ctx, l.svcCtx)
			uploadReq := &types.UploadRequest{
				ImageData: imageData,
			}
			uploadResp, err := uploadLogic.Upload(uploadReq)
			if err != nil {
				logger.Errorw("图片上传失败", logx.Field("error", err), logx.Field("sessionId", sessionId))
				errorEvent := types.StreamEvent{
					Type:      "error",
					Content:   map[string]interface{}{"message": "图片上传失败: " + err.Error()},
					SessionId: sessionId,
				}
				errorJSON, _ := json.Marshal(errorEvent)
				fmt.Fprintf(w, "event: error\ndata: %s\n\n", string(errorJSON))
				w.(http.Flusher).Flush()
				return fmt.Errorf("图片上传失败: %w", err)
			}
			uploadedEvent := types.StreamEvent{
				Type:      "image_uploaded",
				Content:   map[string]interface{}{"url": uploadResp.Url},
				SessionId: sessionId,
			}
			uploadedJSON, _ := json.Marshal(uploadedEvent)
			fmt.Fprintf(w, "event: image_uploaded\ndata: %s\n\n", string(uploadedJSON))
			w.(http.Flusher).Flush()
		}
		messageText = req.Message
		messageType = "image"
	default:
		logger.Errorw("不支持的messageType", logx.Field("messageType", req.MessageType), logx.Field("sessionId", sessionId))
		errorEvent := types.StreamEvent{
			Type:      "error",
			Content:   map[string]interface{}{"message": "不支持的messageType: " + req.MessageType},
			SessionId: sessionId,
		}
		errorJSON, _ := json.Marshal(errorEvent)
		fmt.Fprintf(w, "event: error\ndata: %s\n\n", string(errorJSON))
		w.(http.Flusher).Flush()
		return fmt.Errorf("不支持的messageType: %s", req.MessageType)
	}

	// 保存用户消息
	userMessage := types.ConversationMessage{
		Id:        uuid.New().String(),
		Type:      messageType,
		Sender:    "user",
		Content:   messageText,
		Timestamp: time.Now().Format(time.RFC3339),
		SessionId: sessionId,
	}
	l.svcCtx.Storage.AddMessage(sessionId, userMessage)

	// 获取对话历史（转换为eino Message格式）
	messagesRaw := l.svcCtx.Storage.GetMessages(sessionId)
	messages := make([]types.ConversationMessage, 0, len(messagesRaw))
	for _, msgRaw := range messagesRaw {
		if msg, ok := msgRaw.(types.ConversationMessage); ok {
			messages = append(messages, msg)
		}
	}
	chatHistory := l.convertToEinoMessages(messages, maxContextRounds)

	// 构建请求（用于MultiAgentGraph）
	multiAgentReq := &types.UnifiedStreamConversationRequest{
		MessageType:           req.MessageType,
		Message:               messageText,
		Audio:                 req.Audio,
		Image:                 req.Image,
		SessionId:             sessionId,
		IdentificationContext: req.IdentificationContext,
		UserAge:               userAge,
		MaxContextRounds:      maxContextRounds,
	}

	// 尝试调用MultiAgentGraph
	multiAgentGraph, err := agent.NewMultiAgentGraph(l.ctx, l.svcCtx.Config.AI, logger)
	if err != nil {
		logger.Errorw("MultiAgentGraph初始化失败，降级到单Agent模式", logx.Field("error", err))
		// 降级到单Agent模式
		streamLogic := NewStreamLogic(l.ctx, l.svcCtx)
		return streamLogic.StreamConversationUnified(w, req)
	}

	// 调用MultiAgentGraph执行对话
	answer, err := multiAgentGraph.ExecuteMultiAgentConversation(l.ctx, multiAgentReq, chatHistory)
	if err != nil {
		logger.Errorw("MultiAgentGraph执行失败，降级到单Agent模式", logx.Field("error", err))
		// 降级到单Agent模式
		streamLogic := NewStreamLogic(l.ctx, l.svcCtx)
		return streamLogic.StreamConversationUnified(w, req)
	}

	// 流式返回回答
	messageId := uuid.New().String()
	words := []rune(answer)
	for i, word := range words {
		event := types.StreamEvent{
			Type:      "message",
			Content:   string(word),
			Index:     i,
			SessionId: sessionId,
			MessageId: messageId,
			Markdown:  true,
		}
		eventJSON, _ := json.Marshal(event)
		fmt.Fprintf(w, "event: message\ndata: %s\n\n", string(eventJSON))
		w.(http.Flusher).Flush()
		time.Sleep(30 * time.Millisecond) // 打字机效果
	}

	// 保存助手消息
	assistantMessage := types.ConversationMessage{
		Id:        messageId,
		Type:      "text",
		Sender:    "assistant",
		Content:   answer,
		Timestamp: time.Now().Format(time.RFC3339),
		SessionId: sessionId,
		Markdown:  &[]bool{true}[0],
	}
	l.svcCtx.Storage.AddMessage(sessionId, assistantMessage)

	// 发送完成事件
	doneEvent := types.StreamEvent{
		Type:      "done",
		Content:   map[string]interface{}{"messageId": messageId},
		SessionId: sessionId,
		MessageId: messageId,
	}
	doneJSON, _ := json.Marshal(doneEvent)
	fmt.Fprintf(w, "event: done\ndata: %s\n\n", string(doneJSON))
	w.(http.Flusher).Flush()

	return nil
}

// convertToEinoMessages 转换对话消息为eino Message格式
func (l *AgentLogic) convertToEinoMessages(messages []types.ConversationMessage, maxRounds int) []*schema.Message {
	result := make([]*schema.Message, 0)
	
	// 限制消息数量
	startIdx := 0
	if len(messages) > maxRounds*2 {
		startIdx = len(messages) - maxRounds*2
	}

	for i := startIdx; i < len(messages); i++ {
		msg := messages[i]
		var einoMsg *schema.Message
		
		if msg.Sender == "user" {
			einoMsg = schema.UserMessage(fmt.Sprintf("%v", msg.Content))
		} else {
			einoMsg = schema.AssistantMessage(fmt.Sprintf("%v", msg.Content), nil)
		}
		
		result = append(result, einoMsg)
	}
	
	return result
}

