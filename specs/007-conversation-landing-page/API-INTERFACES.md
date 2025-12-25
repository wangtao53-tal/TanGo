# 对话接口说明文档

## 接口概览

对话功能提供两个接口，用于不同的使用场景：

### 1. `/api/conversation/message` (POST) - 非流式接口

**用途**: 发送消息并获取完整响应（一次性返回）

**特点**:
- 非流式：等待AI生成完整回答后一次性返回
- 适用于：不需要实时反馈的场景，或者兼容性要求
- 响应格式：标准JSON响应

**请求示例**:
```json
POST /api/conversation/message
{
  "message": "这是什么？",
  "sessionId": "session-123",
  "identificationContext": {
    "objectName": "银杏",
    "objectCategory": "自然类",
    "confidence": 0.95,
    "age": 8
  }
}
```

**响应示例**:
```json
{
  "message": {
    "id": "msg-123",
    "type": "text",
    "content": "这是银杏，是一种古老的植物...",
    "sender": "assistant",
    "timestamp": "2025-12-19T10:00:00Z"
  },
  "sessionId": "session-123",
  "type": "text"
}
```

---

### 2. `/api/conversation/stream` (POST) - 流式接口 ⭐ 推荐

**用途**: 发送消息并获取流式响应（实时逐字返回）

**特点**:
- 流式：使用SSE (Server-Sent Events) 实时返回
- 支持打字机效果：文本逐字显示
- POST请求：参数通过请求体传递，避免中文在URL中的编码问题
- 适用于：需要实时反馈的场景，提升用户体验
- 响应格式：SSE事件流

**请求示例**:
```json
POST /api/conversation/stream
Content-Type: application/json

{
  "sessionId": "session-123",
  "message": "这是什么？",
  "messageType": "text",
  "userAge": 8,
  "maxContextRounds": 20,
  "identificationContext": {
    "objectName": "银杏",
    "objectCategory": "自然类",
    "confidence": 0.95,
    "age": 8
  }
}
```

**响应格式** (SSE):
```
event: connected
data: {"type":"connected","sessionId":"session-123"}

event: message
data: {"type":"message","content":"这","index":0,"sessionId":"session-123"}

event: message
data: {"type":"message","content":"是","index":1,"sessionId":"session-123"}

event: message
data: {"type":"message","content":"银","index":2,"sessionId":"session-123"}

...

event: done
data: {"type":"done","sessionId":"session-123"}
```

**事件类型**:
- `connected`: 连接建立
- `message`: 文本消息（逐字符）
- `image_progress`: 图片生成进度
- `image_done`: 图片生成完成
- `card`: 知识卡片
- `error`: 错误
- `done`: 流式完成

---

## 为什么之前会调用两个接口？

**问题**: 之前的实现中，追问时会先调用 `/api/conversation/message`，然后再调用 `/api/conversation/stream`，导致：
1. 重复调用：浪费资源
2. 第一个接口的响应被忽略
3. 用户体验不佳：需要等待两次请求

**原因**: 
- 代码中先调用了 `sendMessage()` 函数（调用非流式接口）
- 然后又创建了 EventSource 连接流式接口
- 两个接口都被调用，但只需要流式接口

---

## 修复方案

**现在**: 直接使用 `/api/conversation/stream` 接口（POST方式）

**修改内容**:
1. 移除了对 `/api/conversation/message` 的调用
2. 使用 fetch API 发送 POST 请求，然后手动解析 SSE 流
3. 通过请求体传递所有必要信息（避免中文在URL中的编码问题）

**优势**:
- ✅ 只调用一个接口，减少网络请求
- ✅ 实时流式返回，打字机效果
- ✅ POST请求，参数通过请求体传递，避免中文编码问题
- ✅ 更好的用户体验
- ✅ 减少服务器负载

---

## 使用建议

### 推荐使用流式接口 (`/api/conversation/stream`)

**适用场景**:
- ✅ 所有对话场景（默认）
- ✅ 需要实时反馈
- ✅ 需要打字机效果
- ✅ 需要更好的用户体验

### 非流式接口 (`/api/conversation/message`) 保留用于

**适用场景**:
- ⚠️ 兼容性要求（不支持SSE的环境）
- ⚠️ 不需要实时反馈的场景
- ⚠️ 测试和调试

---

## 参数说明

### 流式接口参数（POST请求体）

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `sessionId` | string | 否 | 会话ID，如果为空则创建新会话 |
| `message` | string | 是 | 用户消息内容 |
| `messageType` | string | 否 | 消息类型（text/image/voice），默认text |
| `image` | string | 否 | 图片数据（base64或URL），当messageType为image时必填 |
| `voice` | string | 否 | 语音数据（base64），当messageType为voice时必填 |
| `userAge` | int | 否 | 用户年龄（3-18岁），用于内容适配 |
| `maxContextRounds` | int | 否 | 最大上下文轮次，默认20轮 |
| `identificationContext` | object | 否 | 识别结果上下文（首次发送时传递） |
| `identificationContext.objectName` | string | 否 | 识别对象名称 |
| `identificationContext.objectCategory` | string | 否 | 识别对象类别 |
| `identificationContext.confidence` | float | 否 | 识别置信度 |
| `identificationContext.age` | int | 否 | 用户年龄（与userAge重复，优先使用userAge） |

---

## 前端实现示例

```typescript
import { createPostSSEConnection, closePostSSEConnection } from '../services/sse-post';
import type { StreamConversationRequest } from '../types/api';

// 构建POST请求参数
const streamRequest: StreamConversationRequest = {
  sessionId,
  message: text,
  messageType: 'text',
  userAge: identificationContext?.age,
  maxContextRounds: 20,
};

// 首次发送时传递识别结果上下文
if (identificationContext && isFirstMessage) {
  streamRequest.identificationContext = identificationContext;
}

// 使用POST + SSE连接
const abortController = createPostSSEConnection(streamRequest, {
  onMessage: (data: ConversationStreamEvent) => {
    if (data.type === 'message' && data.content) {
      // 处理流式消息
      accumulatedText += data.content;
    }
  },
  onError: (error: Error) => {
    // 处理错误
    console.error('流式返回错误:', error);
  },
  onClose: () => {
    // 流式完成
    console.log('流式完成');
  },
});

// 需要时可以取消连接
// closePostSSEConnection(abortController);
```

---

## 总结

- **流式接口** (`/api/conversation/stream` POST) 是推荐使用的接口，提供实时流式返回和打字机效果
- **POST方式**: 参数通过请求体传递，避免中文在URL中的编码问题
- **非流式接口** (`/api/conversation/message`) 保留用于兼容性场景
- **现在实现**: 追问时只调用流式接口（POST方式），不再重复调用非流式接口
- **用户体验**: 实时反馈，打字机效果，更好的交互体验
- **技术优势**: POST请求体支持复杂数据结构，避免URL长度限制和编码问题

