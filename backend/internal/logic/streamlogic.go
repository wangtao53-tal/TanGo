package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
	"github.com/tango/explore/internal/svc"
	"github.com/tango/explore/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
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

// convertToEinoMessages 将内部消息转换为Eino Message格式
func (l *StreamLogic) convertToEinoMessages(messages []interface{}, maxRounds int) []*schema.Message {
	// 只取最后maxRounds轮（maxRounds * 2条消息）
	start := 0
	if len(messages) > maxRounds*2 {
		start = len(messages) - maxRounds*2
	}

	// 转换为Eino Message格式
	einoMessages := make([]*schema.Message, 0)
	for i := start; i < len(messages); i++ {
		if convMsg, ok := messages[i].(types.ConversationMessage); ok {
			if convMsg.Sender == "user" {
				// 提取文本内容
				content := ""
				if str, ok := convMsg.Content.(string); ok {
					content = str
				} else {
					// 如果是对象，转换为JSON字符串
					contentBytes, _ := json.Marshal(convMsg.Content)
					content = string(contentBytes)
				}
				einoMessages = append(einoMessages, schema.UserMessage(content))
			} else if convMsg.Sender == "assistant" {
				// 提取文本内容
				content := ""
				if str, ok := convMsg.Content.(string); ok {
					content = str
				} else {
					contentBytes, _ := json.Marshal(convMsg.Content)
					content = string(contentBytes)
				}
				einoMessages = append(einoMessages, schema.AssistantMessage(content, nil))
			}
		}
	}

	return einoMessages
}

// getContextMessages 获取上下文消息（最多20轮）
func (l *StreamLogic) getContextMessages(sessionId string, maxRounds int) []*schema.Message {
	messages := l.svcCtx.Storage.GetMessages(sessionId)
	return l.convertToEinoMessages(messages, maxRounds)
}

// StreamConversation 流式对话，集成Eino流式输出和SSE发送
func (l *StreamLogic) StreamConversation(
	w http.ResponseWriter,
	req types.StreamConversationRequest,
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

	// 设置最大上下文轮次
	maxContextRounds := req.MaxContextRounds
	if maxContextRounds <= 0 {
		maxContextRounds = 20 // 默认20轮
	}

	// 获取用户年龄（从识别结果上下文或请求参数）
	userAge := req.UserAge
	objectName := ""
	objectCategory := ""

	// 如果有识别结果上下文，提取信息
	if req.IdentificationContext != nil {
		if userAge == 0 && req.IdentificationContext.Age > 0 {
			userAge = req.IdentificationContext.Age
		}
		objectName = req.IdentificationContext.ObjectName
		objectCategory = req.IdentificationContext.ObjectCategory
		// 保存识别结果上下文到会话
		l.svcCtx.Storage.SetData(sessionId, "identificationContext", req.IdentificationContext)
	} else {
		// 尝试从会话数据中获取识别结果上下文
		if ctxData, ok := l.svcCtx.Storage.GetData(sessionId, "identificationContext"); ok {
			if ctx, ok := ctxData.(*types.IdentificationContext); ok {
				if userAge == 0 && ctx.Age > 0 {
					userAge = ctx.Age
				}
				objectName = ctx.ObjectName
				objectCategory = ctx.ObjectCategory
			}
		}
	}

	// 默认年龄
	if userAge == 0 {
		userAge = 8 // 默认8岁
	}

	// 保存用户消息到存储
	userMessage := types.ConversationMessage{
		Id:        uuid.New().String(),
		Type:      "text",
		Sender:    "user",
		Content:   req.Message,
		Timestamp: time.Now().Format(time.RFC3339),
		SessionId: sessionId,
	}
	l.svcCtx.Storage.AddMessage(sessionId, userMessage)

	// 发送连接建立事件
	connectedEvent := types.StreamEvent{
		Type:      "connected",
		SessionId: sessionId,
	}
	connectedJSON, _ := json.Marshal(connectedEvent)
	fmt.Fprintf(w, "event: connected\ndata: %s\n\n", string(connectedJSON))
	w.(http.Flusher).Flush()

	// 检查Agent是否可用
	if l.svcCtx.Agent == nil {
		logger.Error("Agent未初始化，使用Mock流式响应")
		return l.streamTextMock(w, sessionId, req.Message)
	}

	graph := l.svcCtx.Agent.GetGraph()
	if graph == nil {
		logger.Error("Graph未初始化，使用Mock流式响应")
		return l.streamTextMock(w, sessionId, req.Message)
	}

	// 获取对话节点（需要通过Graph访问）
	// 注意：这里需要从Graph中获取conversationNode
	// 由于Graph结构体是私有的，我们需要添加一个方法来获取conversationNode
	// 或者直接在Graph中添加StreamConversation方法
	// 暂时使用Mock实现，后续可以通过Graph方法调用

	// 获取上下文消息
	contextMessages := l.getContextMessages(sessionId, maxContextRounds)

	// 获取对话节点
	conversationNode := graph.GetConversationNode()
	if conversationNode == nil {
		logger.Error("ConversationNode未初始化，使用Mock流式响应")
		return l.streamTextMock(w, sessionId, req.Message)
	}

	logger.Infow("开始流式对话",
		logx.Field("sessionId", sessionId),
		logx.Field("userAge", userAge),
		logx.Field("objectName", objectName),
		logx.Field("contextRounds", len(contextMessages)/2),
	)

	// 调用真实的Eino流式接口
	streamReader, err := conversationNode.StreamConversation(
		l.ctx,
		req.Message,
		contextMessages,
		userAge,
		objectName,
		objectCategory,
	)
	if err != nil {
		logger.Errorw("调用Eino流式接口失败，使用Mock响应",
			logx.Field("error", err),
		)
		return l.streamTextMock(w, sessionId, req.Message)
	}

	// 实现真实的Eino StreamReader读取逻辑
	// StreamReader实现了io.Reader接口，使用Read()方法读取
	logger.Infow("开始读取Eino流式数据",
		logx.Field("streamReader", streamReader != nil),
	)

	// 创建助手消息ID
	assistantMessageId := uuid.New().String()
	fullText := ""
	index := 0

	// 读取流式数据
	// StreamReader使用Recv()方法读取数据
	// 当返回io.EOF时表示流结束
	for {
		msg, err := streamReader.Recv()
		if errors.Is(err, io.EOF) {
			// 流结束，正常退出
			break
		}
		if err != nil {
			// 其他错误
			logger.Errorw("读取流式数据失败",
				logx.Field("error", err),
			)
			// 发送错误事件
			errorEvent := types.StreamEvent{
				Type:      "error",
				Content:   map[string]interface{}{"message": "流式读取失败: " + err.Error()},
				SessionId: sessionId,
			}
			errorJSON, _ := json.Marshal(errorEvent)
			fmt.Fprintf(w, "event: error\ndata: %s\n\n", string(errorJSON))
			w.(http.Flusher).Flush()
			return err
		}

		if msg == nil {
			continue
		}

		// 提取文本内容
		if msg.Content != "" {
			// 逐字符发送（用于打字机效果）
			textRunes := []rune(msg.Content)
			for _, char := range textRunes {
				event := types.StreamEvent{
					Type:      "message",
					Content:   string(char),
					Index:     index,
					SessionId: sessionId,
					MessageId: assistantMessageId,
				}
				eventJSON, _ := json.Marshal(event)
				fmt.Fprintf(w, "event: message\ndata: %s\n\n", string(eventJSON))
				w.(http.Flusher).Flush()
				fullText += string(char)
				index++
			}
		}
	}

	// 保存助手消息到存储
	assistantMessage := types.ConversationMessage{
		Id:        assistantMessageId,
		Type:      "text",
		Sender:    "assistant",
		Content:   fullText,
		Timestamp: time.Now().Format(time.RFC3339),
		SessionId: sessionId,
	}
	l.svcCtx.Storage.AddMessage(sessionId, assistantMessage)

	// 发送完成事件
	doneEvent := types.StreamEvent{
		Type:      "done",
		SessionId: sessionId,
		MessageId: assistantMessageId,
	}
	doneJSON, _ := json.Marshal(doneEvent)
	fmt.Fprintf(w, "event: done\ndata: %s\n\n", string(doneJSON))
	w.(http.Flusher).Flush()

	logger.Infow("流式对话完成",
		logx.Field("sessionId", sessionId),
		logx.Field("messageLength", len(fullText)),
	)

	return nil
}

// streamTextMock Mock流式文本响应
func (l *StreamLogic) streamTextMock(w http.ResponseWriter, sessionId string, message string) error {
	text := fmt.Sprintf("这是一个Mock流式响应。您的问题是：%s。待接入真实AI模型后，将实现真实的流式文本生成。", message)
	words := []rune(text)

	for i, word := range words {
		event := types.StreamEvent{
			Type:      "message",
			Content:   string(word),
			Index:     i,
			SessionId: sessionId,
		}
		eventJSON, _ := json.Marshal(event)
		fmt.Fprintf(w, "event: message\ndata: %s\n\n", string(eventJSON))
		w.(http.Flusher).Flush()
		time.Sleep(50 * time.Millisecond) // 模拟延迟
	}

	// 发送完成事件
	doneEvent := types.StreamEvent{
		Type:      "done",
		SessionId: sessionId,
	}
	doneJSON, _ := json.Marshal(doneEvent)
	fmt.Fprintf(w, "event: done\ndata: %s\n\n", string(doneJSON))
	w.(http.Flusher).Flush()

	return nil
}
