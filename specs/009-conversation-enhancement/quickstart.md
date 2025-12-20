# 快速开始指南

**创建日期**: 2025-12-21  
**功能**: 对话页面完善 - Agent模型流式返回

## 概述

本指南介绍如何实现统一流式接口，支持文本输入、语音输入、图片输入三种输入方式。统一接口通过 `messageType` 字段明确指定输入类型，根据输入类型自动处理（语音识别、图片上传等），然后通过Agent模型流式返回回答。

## 实现步骤

### 1. 后端修改

#### 1.1 定义统一流式请求类型

**文件**: `backend/internal/types/types.go`

**修改内容**:
- 定义 `UnifiedStreamConversationRequest` 类型，包含 `messageType` 字段（必填）
- 根据 `messageType` 的值，请求包含对应的字段：`message`、`audio` 或 `image`

**示例代码**:

```go
// UnifiedStreamConversationRequest 统一流式对话请求
type UnifiedStreamConversationRequest struct {
	MessageType           string                 `json:"messageType"`                       // 消息类型（必填）：text/voice/image
	Message               string                 `json:"message,optional"`                  // 文本消息，当messageType为text时必填
	Audio                 string                 `json:"audio,optional"`                    // 语音数据（base64），当messageType为voice时必填
	Image                 string                 `json:"image,optional"`                    // 图片数据（base64或URL），当messageType为image时必填
	SessionId             string                 `json:"sessionId,optional"`                // 会话ID，如果为空则创建新会话
	IdentificationContext *IdentificationContext `json:"identificationContext,optional"`    // 识别结果上下文（可选）
	UserAge               int                    `json:"userAge,optional"`                 // 用户年龄（3-18岁），用于内容适配
	MaxContextRounds      int                    `json:"maxContextRounds,optional"`         // 最大上下文轮次，默认20轮
}
```

#### 1.2 扩展 ConversationNode 支持多模态输入

**文件**: `backend/internal/agent/nodes/conversation_node.go`

**修改内容**:
- 扩展 `StreamConversation` 方法签名，添加 `imageURL` 参数
- 如果提供了图片URL，构建多模态消息

**示例代码**:

```go
// StreamConversation 流式对话，支持多模态输入
func (n *ConversationNode) StreamConversation(
	ctx context.Context,
	message string,
	contextMessages []*schema.Message,
	userAge int,
	objectName string,
	objectCategory string,
	imageURL string, // 新增：图片URL参数
) (*schema.StreamReader[*schema.Message], error) {
	// ... 系统prompt构建 ...

	// 构建用户消息
	var userMsg *schema.Message
	if imageURL != "" {
		// 多模态消息（图片+文本）
		userMsg = &schema.Message{
			Role: schema.User,
			UserInputMultiContent: []schema.MessageInputPart{
				{
					Type: schema.ChatMessagePartTypeImageURL,
					Image: &schema.MessageInputImage{
						MessagePartCommon: schema.MessagePartCommon{
							URL: &imageURL,
						},
						Detail: schema.ImageURLDetailAuto,
					},
				},
				{
					Type: schema.ChatMessagePartTypeText,
					Text: message,
				},
			},
		}
	} else {
		// 文本消息
		userMsg = schema.UserMessage(message)
	}
	messages = append(messages, userMsg)

	// 调用Eino ChatModel的Stream接口
	streamReader, err := n.chatModel.Stream(ctx, messages)
	// ...
}
```

#### 1.3 扩展 StreamLogic 实现统一接口处理逻辑

**文件**: `backend/internal/logic/streamlogic.go`

**修改内容**:
- 扩展 `StreamConversation` 方法，接收 `UnifiedStreamConversationRequest`
- 根据 `messageType` 字段处理不同输入类型
- 集成语音识别和图片上传逻辑

**示例代码**:

```go
func (l *StreamLogic) StreamConversation(
	w http.ResponseWriter,
	req types.UnifiedStreamConversationRequest,
) error {
	// 验证 messageType 字段
	if req.MessageType == "" {
		// 发送错误事件
		return fmt.Errorf("messageType字段必填")
	}

	var messageText string
	var imageURL string

	// 根据 messageType 处理不同输入类型
	switch req.MessageType {
	case "text":
		if req.Message == "" {
			return fmt.Errorf("messageType为text时，message字段必填")
		}
		messageText = req.Message

	case "voice":
		if req.Audio == "" {
			return fmt.Errorf("messageType为voice时，audio字段必填")
		}
		// 语音识别
		voiceLogic := NewVoiceLogic(l.ctx, l.svcCtx)
		voiceReq := &types.VoiceRequest{
			Audio:     req.Audio,
			SessionId: req.SessionId,
		}
		voiceResp, err := voiceLogic.RecognizeVoice(voiceReq)
		if err != nil {
			return fmt.Errorf("语音识别失败: %w", err)
		}
		messageText = voiceResp.Text

	case "image":
		if req.Image == "" {
			return fmt.Errorf("messageType为image时，image字段必填")
		}
		// 如果image是base64，需要先上传
		if strings.HasPrefix(req.Image, "data:") || !strings.HasPrefix(req.Image, "http") {
			// 上传图片
			uploadLogic := NewUploadLogic(l.ctx, l.svcCtx)
			uploadReq := &types.UploadRequest{
				ImageData: req.Image,
			}
			uploadResp, err := uploadLogic.Upload(uploadReq)
			if err != nil {
				return fmt.Errorf("图片上传失败: %w", err)
			}
			imageURL = uploadResp.Url
		} else {
			imageURL = req.Image
		}
		messageText = req.Message // 图片输入时，message是可选的文本描述

	default:
		return fmt.Errorf("不支持的messageType: %s", req.MessageType)
	}

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
	// ...
}
```

#### 1.4 更新 StreamConversationHandler 支持统一接口

**文件**: `backend/internal/handler/streamhandler.go`

**修改内容**:
- 更新 `StreamConversationHandler`，接收 `UnifiedStreamConversationRequest`
- 调用 `StreamLogic.StreamConversation` 处理统一接口逻辑

**示例代码**:

```go
func StreamConversationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 设置SSE响应头
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		var req types.UnifiedStreamConversationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			// 发送错误事件
			return
		}

		// 调用统一流式逻辑
		streamLogic := logic.NewStreamLogic(r.Context(), svcCtx)
		if err := streamLogic.StreamConversation(w, req); err != nil {
			// 错误已在StreamConversation中处理
			return
		}
	}
}
```

#### 1.5 修改错误处理，禁止Mock降级

**文件**: `backend/internal/logic/streamlogic.go`

**修改内容**:
- 移除Mock降级逻辑
- 记录详细错误日志
- 向用户发送错误事件

**示例代码**:

```go
// 检查Agent是否可用
if l.svcCtx.Agent == nil {
	logger.Error("Agent未初始化")
	// 发送错误事件，不允许降级到Mock数据
	errorEvent := types.StreamEvent{
		Type:      "error",
		Content:   map[string]interface{}{
			"message": "Agent未初始化，无法进行流式对话",
		},
		SessionId: sessionId,
	}
	errorJSON, _ := json.Marshal(errorEvent)
	fmt.Fprintf(w, "event: error\ndata: %s\n\n", string(errorJSON))
	w.(http.Flusher).Flush()
	return fmt.Errorf("Agent未初始化")
}
```

### 2. API路由配置

**文件**: `backend/api/explore.api`

**修改内容**:
- 更新统一流式接口路由

```api
service explore {
	// 统一流式对话接口（支持文本、语音、图片输入）
	@handler StreamConversationHandler
	post /api/conversation/stream (UnifiedStreamConversationRequest) returns (stream)
}
```

### 3. 前端修改

前端统一调用 `/api/conversation/stream` 接口，根据输入类型设置 `messageType` 字段：

**示例代码**:

```typescript
// 统一流式对话接口
async function sendStreamConversation(
  messageType: "text" | "voice" | "image",
  data: {
    message?: string;
    audio?: string;
    image?: string;
  },
  sessionId?: string
) {
  const requestBody = {
    messageType,
    sessionId,
    ...data,
  };

  const response = await fetch('/api/conversation/stream', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(requestBody),
  });

  // 处理SSE流式响应
  const reader = response.body?.getReader();
  const decoder = new TextDecoder();

  while (true) {
    const { done, value } = await reader.read();
    if (done) break;

    const chunk = decoder.decode(value);
    // 解析SSE事件并更新UI
    parseSSEEvent(chunk);
  }
}

// 文本输入
sendStreamConversation("text", { message: "用户问题" });

// 语音输入
sendStreamConversation("voice", { audio: "base64音频数据" });

// 图片输入
sendStreamConversation("image", { image: "base64图片数据或URL", message: "可选文本描述" });
```

## 测试步骤

### 1. 测试统一接口的文本输入流式返回

1. 调用统一接口 `/api/conversation/stream`，设置 `messageType: "text"`，发送 `message` 字段
2. 验证是否开始接收流式响应
3. 验证流式文本是否正确显示

### 2. 测试统一接口的语音输入流式返回

1. 调用统一接口 `/api/conversation/stream`，设置 `messageType: "voice"`，发送 `audio` 字段
2. 验证语音识别是否成功
3. 验证是否开始接收流式响应
4. 验证流式文本是否正确显示

### 3. 测试统一接口的图片输入流式返回

1. 调用统一接口 `/api/conversation/stream`，设置 `messageType: "image"`，发送 `image` 字段
2. 验证图片上传是否成功（如果需要）
3. 验证是否开始接收流式响应
4. 验证流式文本是否正确显示

### 4. 测试 messageType 字段验证

1. 测试缺少 `messageType` 字段的请求
2. 测试 `messageType` 为无效值的请求
3. 测试 `messageType` 与对应字段不匹配的请求（如 `messageType: "voice"` 但缺少 `audio` 字段）
4. 验证错误提示是否正确

### 5. 测试错误处理

1. 模拟Agent模型调用失败
2. 验证是否记录详细错误日志
3. 验证是否向用户发送错误事件
4. 验证是否没有降级到Mock数据

## 注意事项

1. **统一接口设计**: 对话页面使用一个统一的流式接口，通过 `messageType` 字段区分输入类型
2. **禁止Mock数据**: 统一接口必须调用真实的Agent模型，禁止使用Mock数据
3. **错误处理**: Agent模型调用失败时，必须记录详细错误日志，向用户发送错误事件
4. **messageType字段验证**: 必须验证 `messageType` 字段，确保与对应字段匹配
5. **多模态支持**: 图片输入支持多模态消息（图片+文本），语音输入先转换为文本
6. **独立卡片接口**: 生成知识卡片使用独立的接口，不是统一对话接口

## 总结

通过实现统一流式接口，通过 `messageType` 字段明确指定输入类型，集成语音识别和图片上传逻辑，实现文本、语音、图片三种输入方式的Agent模型流式返回功能。统一接口简化了前端调用逻辑，统一了后端处理流程，并禁止使用Mock数据。

