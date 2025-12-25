# 数据模型定义

**创建日期**: 2025-12-21  
**功能**: 对话页面完善 - Agent模型流式返回

## 概述

本文档定义了语音输入和图片上传后Agent模型流式返回功能的数据模型。主要扩展了流式对话请求和响应类型，支持多模态输入。

## 核心实体

### 1. UnifiedStreamConversationRequest（统一流式对话请求）

统一流式对话请求，通过 `messageType` 字段明确指定输入类型，支持文本、语音、图片三种输入方式。

```typescript
interface UnifiedStreamConversationRequest {
  // 必填字段
  messageType: "text" | "voice" | "image";  // 消息类型（必填），明确指定输入类型
  
  // 条件字段（根据messageType必填其一）
  message?: string;                  // 用户消息内容（文本），当messageType为text时必填
  audio?: string;                    // 语音数据（base64），当messageType为voice时必填
  image?: string;                    // 图片数据（base64或URL），当messageType为image时必填
  
  // 可选字段
  sessionId?: string;                // 会话ID，如果为空则创建新会话
  identificationContext?: IdentificationContext;  // 识别结果上下文（可选）
  userAge?: number;                  // 用户年龄（3-18岁），用于内容适配
  maxContextRounds?: number;         // 最大上下文轮次，默认20轮
}
```

**设计说明**:
- `messageType` 字段（必填）：明确指定输入类型，避免字段冲突和歧义
- 根据 `messageType` 的值，请求必须包含对应的字段：
  - `messageType: "text"` → 必须包含 `message` 字段
  - `messageType: "voice"` → 必须包含 `audio` 字段
  - `messageType: "image"` → 必须包含 `image` 字段
- 统一接口设计：一个接口支持所有输入类型，简化前端调用

### 2. StreamEvent（扩展）

SSE流式事件类型，扩展支持多模态事件。

```typescript
interface StreamEvent {
  // 基础字段
  type: string;                    // 事件类型：connected/message/image_progress/image_done/card/error/done
  content: any;                     // 事件内容
  index?: number;                   // 文本消息的字符索引（用于打字机效果）
  sessionId?: string;               // 会话ID
  messageId?: string;               // 消息ID
  markdown?: boolean;               // 内容是否为Markdown格式
  
  // 多模态字段
  progress?: number;                // 图片生成进度（0-100）
  imageUrl?: string;                // 图片URL（用于图文混排）
}
```

**扩展说明**:
- 保持现有字段不变，确保向后兼容
- 新增 `imageUrl` 字段：用于图文混排输出

### 3. ConversationMessage（扩展）

对话消息，扩展支持多模态消息类型。

```typescript
interface ConversationMessage {
  // 基础字段
  id: string;                       // 消息ID
  type: string;                     // 消息类型：text/image/voice/card
  sender: string;                    // 发送者：user/assistant
  content: any;                      // 消息内容（文本、图片、语音、卡片等）
  timestamp: string;                 // 消息时间戳
  sessionId?: string;                // 会话ID
  
  // 流式字段
  isStreaming?: boolean;             // 是否正在流式返回
  streamingText?: string;            // 流式传输中的累积文本（仅系统消息）
  markdown?: boolean;                // 内容是否包含Markdown格式（仅文本消息）
  
  // 多模态字段
  imageUrl?: string;                 // 图片URL（用于图片消息）
  voiceUrl?: string;                 // 语音URL（用于语音消息）
}
```

**扩展说明**:
- 保持现有字段不变，确保向后兼容
- 新增 `imageUrl` 和 `voiceUrl` 字段：用于多模态消息

## 数据流

### 1. 统一接口处理流程

#### 文本输入流程

```
用户输入文本
  ↓
前端调用统一接口 /api/conversation/stream
  messageType: "text", message: "用户文本"
  ↓
后端直接使用 message 字段
  ↓
后端调用 StreamLogic.StreamConversation
  ↓
后端通过SSE流式返回Agent模型响应
  ↓
前端接收SSE事件，实时渲染
```

#### 语音输入流程

```
用户录制语音
  ↓
前端调用统一接口 /api/conversation/stream
  messageType: "voice", audio: "base64音频数据"
  ↓
后端识别语音，转换为文本（调用 VoiceLogic.RecognizeVoice）
  ↓
后端使用识别的文本调用 StreamLogic.StreamConversation
  ↓
后端通过SSE流式返回Agent模型响应
  ↓
前端接收SSE事件，实时渲染
```

#### 图片输入流程

```
用户上传图片
  ↓
前端调用统一接口 /api/conversation/stream
  messageType: "image", image: "base64图片数据或URL"
  ↓
后端上传图片，获取图片URL（调用 UploadLogic.Upload，如果需要）
  ↓
后端构建多模态消息（图片+文本），调用 StreamLogic.StreamConversation
  ↓
后端通过SSE流式返回Agent模型响应
  ↓
前端接收SSE事件，实时渲染
```

**数据转换**:
- 文本输入：`message` → 用户消息 → Agent模型输入 → 流式响应
- 语音输入：`audio`（base64） → 识别文本 → 用户消息 → Agent模型输入 → 流式响应
- 图片输入：`image`（base64或URL） → 图片URL → 多模态消息（图片+文本） → Agent模型输入 → 流式响应

## 消息构建规则

### 1. 文本消息

```go
messages = []*schema.Message{
    schema.SystemMessage(systemPrompt),
    // ... 上下文消息 ...
    schema.UserMessage(message),
}
```

### 2. 图片消息（多模态）

```go
userMsg := &schema.Message{
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
            Text: message, // 用户输入的文本（如果有）
        },
    },
}
messages = append(messages, userMsg)
```

### 3. 语音消息

语音消息先转换为文本，然后按照文本消息处理：

```go
// 语音识别后
recognizedText := "识别的文本内容"
messages = []*schema.Message{
    schema.SystemMessage(systemPrompt),
    // ... 上下文消息 ...
    schema.UserMessage(recognizedText),
}
```

## 状态管理

### 1. 会话状态

- **会话ID**: 唯一标识一个对话会话
- **消息列表**: 存储对话历史消息（最多20轮）
- **识别结果上下文**: 存储识别结果信息（对象名称、类别等）

### 2. 流式状态

- **连接状态**: SSE连接是否建立
- **流式文本**: 累积的流式文本内容
- **Markdown标识**: 内容是否为Markdown格式
- **错误状态**: 是否有错误发生

## 验证规则

### 1. 请求验证

- `messageType` 字段：必填，必须是 `"text"`、`"voice"` 或 `"image"` 之一
- `message` 字段：当 `messageType` 为 `"text"` 时必填
- `audio` 字段：当 `messageType` 为 `"voice"` 时必填
- `image` 字段：当 `messageType` 为 `"image"` 时必填
- `userAge` 字段：范围3-18岁（可选）
- `maxContextRounds` 字段：范围1-20，默认20（可选）

### 2. 响应验证

- SSE事件格式必须符合规范
- 事件类型必须为预定义类型之一
- 事件内容必须符合对应类型的格式要求

## 错误处理

### 1. 错误类型

- **网络错误**: SSE连接中断、超时等
- **模型错误**: Agent模型调用失败、返回错误等
- **参数错误**: 请求参数验证失败等
- **业务错误**: 语音识别失败、图片上传失败等

### 2. 错误响应格式

```typescript
interface StreamEvent {
  type: "error";
  content: {
    message: string;        // 错误消息
    error?: string;         // 错误详情（可选）
    code?: number;          // 错误码（可选）
  };
  sessionId?: string;
}
```

## 总结

通过扩展 `StreamConversationRequest` 和 `StreamEvent` 类型，支持多模态输入（文本、语音、图片），实现语音输入和图片上传后的Agent模型流式返回功能。所有扩展字段均为可选，确保向后兼容。

