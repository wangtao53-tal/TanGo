package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
	"github.com/tango/explore/internal/svc"
	"github.com/tango/explore/internal/types"
	"github.com/tango/explore/internal/utils"
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

// StreamConversation 流式对话，集成Eino流式输出和SSE发送（兼容旧版本）
// 注意：新代码应使用 StreamConversationUnified
func (l *StreamLogic) StreamConversation(
	w http.ResponseWriter,
	req types.StreamConversationRequest,
) error {
	// 转换为统一请求类型
	unifiedReq := types.UnifiedStreamConversationRequest{
		MessageType:           req.MessageType,
		Message:               req.Message,
		Audio:                 req.Voice, // Voice字段映射到Audio
		Image:                 req.Image,
		SessionId:             req.SessionId,
		IdentificationContext: req.IdentificationContext,
		UserAge:               req.UserAge,
		MaxContextRounds:      req.MaxContextRounds,
	}
	// 设置默认值
	if unifiedReq.MessageType == "" {
		unifiedReq.MessageType = "text"
	}
	return l.StreamConversationUnified(w, unifiedReq)
}

// StreamConversationUnified 统一流式对话，支持文本、语音、图片三种输入方式
func (l *StreamLogic) StreamConversationUnified(
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

	// 验证messageType字段
	if req.MessageType == "" {
		logger.Errorw("messageType字段必填",
			logx.Field("sessionId", sessionId),
		)
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
	var imageURL string
	var messageType string

	switch req.MessageType {
	case "text":
		// 验证message字段
		if req.Message == "" {
			logger.Errorw("messageType为text时，message字段必填",
				logx.Field("sessionId", sessionId),
			)
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
		// 验证audio字段
		if req.Audio == "" {
			logger.Errorw("messageType为voice时，audio字段必填",
				logx.Field("sessionId", sessionId),
			)
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
			logger.Errorw("语音识别失败",
				logx.Field("error", err),
				logx.Field("sessionId", sessionId),
			)
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
		// 发送语音识别完成事件
		recognizedEvent := types.StreamEvent{
			Type:      "voice_recognized",
			Content:   map[string]interface{}{"text": voiceResp.Text},
			SessionId: sessionId,
		}
		recognizedJSON, _ := json.Marshal(recognizedEvent)
		fmt.Fprintf(w, "event: voice_recognized\ndata: %s\n\n", string(recognizedJSON))
		w.(http.Flusher).Flush()

	case "image":
		// 验证image字段
		if req.Image == "" {
			logger.Errorw("messageType为image时，image字段必填",
				logx.Field("sessionId", sessionId),
			)
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
		// 如果image是base64，需要先上传
		if strings.HasPrefix(req.Image, "data:") || !strings.HasPrefix(req.Image, "http") {
			// 提取base64数据（移除data URL前缀）
			imageData := req.Image
			if strings.HasPrefix(imageData, "data:image/") {
				// 移除 data:image/xxx;base64, 前缀
				parts := strings.SplitN(imageData, ",", 2)
				if len(parts) == 2 {
					imageData = parts[1]
				}
			}
			uploadLogic := NewUploadLogic(l.ctx, l.svcCtx)
			uploadReq := &types.UploadRequest{
				ImageData: imageData,
			}
			uploadResp, err := uploadLogic.Upload(uploadReq)
			if err != nil {
				logger.Errorw("图片上传失败",
					logx.Field("error", err),
					logx.Field("sessionId", sessionId),
				)
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
			imageURL = uploadResp.Url
			// 发送图片上传完成事件
			uploadedEvent := types.StreamEvent{
				Type:      "image_uploaded",
				Content:   map[string]interface{}{"url": uploadResp.Url},
				SessionId: sessionId,
			}
			uploadedJSON, _ := json.Marshal(uploadedEvent)
			fmt.Fprintf(w, "event: image_uploaded\ndata: %s\n\n", string(uploadedJSON))
			w.(http.Flusher).Flush()
		} else {
			imageURL = req.Image
		}
		messageText = req.Message // 图片输入时，message是可选的文本描述
		messageType = "image"
	default:
		logger.Errorw("不支持的messageType",
			logx.Field("messageType", req.MessageType),
			logx.Field("sessionId", sessionId),
		)
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

	// 保存用户消息到存储
	userMessage := types.ConversationMessage{
		Id:        uuid.New().String(),
		Type:      messageType,
		Sender:    "user",
		Content:   messageText,
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

	// 检查Agent是否可用（根据配置决定是否允许Mock降级）
	useAIModel := l.svcCtx.Config.AI.UseAIModel
	if !useAIModel {
		// 如果配置为false，可以使用Mock数据
		logger.Infow("USE_AI_MODEL=false，使用Mock数据",
			logx.Field("sessionId", sessionId),
		)
		return l.streamTextMock(w, sessionId, messageText)
	}

	// 当UseAIModel=true时，必须使用AI模型，禁止Mock降级
	if l.svcCtx.Agent == nil {
		logger.Errorw("Agent未初始化，无法进行流式对话",
			logx.Field("sessionId", sessionId),
			logx.Field("useAIModel", useAIModel),
		)
		// 发送错误事件，不允许降级到Mock数据
		errorEvent := types.StreamEvent{
			Type:      "error",
			Content:   map[string]interface{}{"message": "Agent未初始化，无法进行流式对话"},
			SessionId: sessionId,
		}
		errorJSON, _ := json.Marshal(errorEvent)
		fmt.Fprintf(w, "event: error\ndata: %s\n\n", string(errorJSON))
		w.(http.Flusher).Flush()
		return fmt.Errorf("Agent未初始化")
	}

	graph := l.svcCtx.Agent.GetGraph()
	if graph == nil {
		logger.Errorw("Graph未初始化，无法进行流式对话",
			logx.Field("sessionId", sessionId),
			logx.Field("useAIModel", useAIModel),
		)
		// 发送错误事件，不允许降级到Mock数据
		errorEvent := types.StreamEvent{
			Type:      "error",
			Content:   map[string]interface{}{"message": "Graph未初始化，无法进行流式对话"},
			SessionId: sessionId,
		}
		errorJSON, _ := json.Marshal(errorEvent)
		fmt.Fprintf(w, "event: error\ndata: %s\n\n", string(errorJSON))
		w.(http.Flusher).Flush()
		return fmt.Errorf("Graph未初始化")
	}

	// 获取上下文消息
	contextMessages := l.getContextMessages(sessionId, maxContextRounds)

	// 获取对话节点
	conversationNode := graph.GetConversationNode()
	if conversationNode == nil {
		logger.Errorw("ConversationNode未初始化，无法进行流式对话",
			logx.Field("sessionId", sessionId),
			logx.Field("useAIModel", useAIModel),
		)
		// 发送错误事件，不允许降级到Mock数据
		errorEvent := types.StreamEvent{
			Type:      "error",
			Content:   map[string]interface{}{"message": "ConversationNode未初始化，无法进行流式对话"},
			SessionId: sessionId,
		}
		errorJSON, _ := json.Marshal(errorEvent)
		fmt.Fprintf(w, "event: error\ndata: %s\n\n", string(errorJSON))
		w.(http.Flusher).Flush()
		return fmt.Errorf("ConversationNode未初始化")
	}

	logger.Infow("开始流式对话",
		logx.Field("sessionId", sessionId),
		logx.Field("userAge", userAge),
		logx.Field("objectName", objectName),
		logx.Field("contextRounds", len(contextMessages)/2),
		logx.Field("hasImage", imageURL != ""),
		logx.Field("messageType", req.MessageType),
	)

	// 调用真实的Eino流式接口（传入图片URL）
	streamReader, err := conversationNode.StreamConversation(
		l.ctx,
		messageText,
		contextMessages,
		userAge,
		objectName,
		objectCategory,
		imageURL, // 传入图片URL（如果提供）
	)
	if err != nil {
		logger.Errorw("调用Eino流式接口失败",
			logx.Field("error", err),
			logx.Field("errorType", "model_error"),
			logx.Field("sessionId", sessionId),
			logx.Field("userAge", userAge),
			logx.Field("hasImage", imageURL != ""),
			logx.Field("useAIModel", useAIModel),
		)
		// 发送错误事件，根据配置决定是否允许降级到Mock数据
		if !useAIModel {
			// 如果配置为false，可以使用Mock数据作为降级方案
			logger.Infow("USE_AI_MODEL=false，降级到Mock数据",
				logx.Field("sessionId", sessionId),
				logx.Field("error", err),
			)
			return l.streamTextMock(w, sessionId, messageText)
		}
		// 当UseAIModel=true时，不允许降级到Mock数据
		errorEvent := types.StreamEvent{
			Type:      "error",
			Content:   map[string]interface{}{"message": "Agent模型调用失败: " + err.Error()},
			SessionId: sessionId,
		}
		errorJSON, _ := json.Marshal(errorEvent)
		fmt.Fprintf(w, "event: error\ndata: %s\n\n", string(errorJSON))
		w.(http.Flusher).Flush()
		return fmt.Errorf("调用AI模型失败: %w", err)
	}

	// 实现真实的Eino StreamReader读取逻辑
	logger.Infow("开始读取Eino流式数据",
		logx.Field("streamReader", streamReader != nil),
	)

	// 创建助手消息ID
	assistantMessageId := uuid.New().String()
	fullText := ""
	index := 0
	isMarkdown := false // 跟踪是否为Markdown格式

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
			// 检测Markdown格式（在累积足够文本后检测）
			fullText += msg.Content
			if !isMarkdown && len(fullText) > 10 {
				isMarkdown = utils.DetectMarkdown(fullText)
			}

			// 逐字符发送（用于打字机效果）
			textRunes := []rune(msg.Content)
			for _, char := range textRunes {
				event := types.StreamEvent{
					Type:      "message",
					Content:   string(char),
					Index:     index,
					SessionId: sessionId,
					MessageId: assistantMessageId,
					Markdown:  isMarkdown,
				}
				eventJSON, _ := json.Marshal(event)
				fmt.Fprintf(w, "event: message\ndata: %s\n\n", string(eventJSON))
				w.(http.Flusher).Flush()
				index++
			}
		}
	}

	// 最终检测Markdown格式
	if !isMarkdown {
		isMarkdown = utils.DetectMarkdown(fullText)
	}

	// 保存助手消息到存储
	markdownPtr := &isMarkdown
	assistantMessage := types.ConversationMessage{
		Id:            assistantMessageId,
		Type:          "text",
		Sender:        "assistant",
		Content:       fullText,
		Timestamp:     time.Now().Format(time.RFC3339),
		SessionId:     sessionId,
		StreamingText: fullText, // 保存流式文本
		Markdown:      markdownPtr,
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
	text := fmt.Sprintf("这是一个Mock流式响应 🌟。您的问题是：%s。待接入真实AI模型后，将实现真实的流式文本生成 ✨。", message)
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
