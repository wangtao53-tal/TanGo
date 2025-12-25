# API 合约文档：H5对话落地页

**功能**: 007-conversation-landing-page  
**创建日期**: 2025-12-19

## 概述

本文档定义H5对话落地页功能的API接口规范，包括流式对话接口、上下文管理、消息格式等。

## API端点

### 1. 流式对话接口（SSE）

**端点**: `GET /api/conversation/stream`

**描述**: 基于Server-Sent Events (SSE)的流式对话接口，支持文本流式输出、图片流式输出、图文混排输出、知识卡片输出。

**请求参数**（Query Parameters）:
- `sessionId` (string, optional): 会话ID，如果为空则创建新会话
- `message` (string, required): 用户消息内容
- `messageType` (string, optional): 消息类型，默认"text"，可选值：text/image/voice
- `image` (string, optional): 图片数据（base64或URL），当messageType为image时必填
- `voice` (string, optional): 语音数据（base64），当messageType为voice时必填
- `userAge` (int, optional): 用户年龄（3-18岁），用于内容适配
- `maxContextRounds` (int, optional): 最大上下文轮次，默认20轮

**请求体**（JSON，可选）:
```json
{
  "identificationContext": {
    "objectName": "银杏",
    "objectCategory": "自然类",
    "confidence": 0.95,
    "keywords": ["植物", "叶子"],
    "age": 8
  }
}
```

**响应格式**（SSE）:
```
event: connected
data: {"type":"connected","sessionId":"session-123"}

event: message
data: {"type":"message","content":"文本片段","index":0}

event: message
data: {"type":"message","content":"文本片段","index":1}

event: image_progress
data: {"type":"image_progress","progress":50}

event: image_done
data: {"type":"image_done","content":{"url":"https://..."}}

event: card
data: {"type":"card","content":{"id":"card-1","type":"science",...}}

event: done
data: {"type":"done"}
```

**事件类型说明**:
- `connected`: 连接建立，返回sessionId
- `message`: 文本消息片段（逐字符发送，用于打字机效果）
- `image_progress`: 图片生成进度（0-100）
- `image_done`: 图片生成完成，返回图片URL
- `card`: 知识卡片消息（完整卡片数据）
- `error`: 错误事件
- `done`: 流式传输完成

**示例请求**:
```bash
curl -N "http://localhost:8877/api/conversation/stream?sessionId=session-123&message=这是什么？&userAge=8"
```

**示例响应**:
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
data: {"type":"done"}
```

### 2. 非流式对话接口（兼容性）

**端点**: `POST /api/conversation/message`

**描述**: 非流式对话接口，用于兼容不支持SSE的场景或简单测试。

**请求体**:
```json
{
  "sessionId": "session-123",
  "message": "这是什么？",
  "messageType": "text",
  "userAge": 8,
  "identificationContext": {
    "objectName": "银杏",
    "objectCategory": "自然类",
    "confidence": 0.95,
    "keywords": ["植物", "叶子"],
    "age": 8
  },
  "maxContextRounds": 20
}
```

**响应**:
```json
{
  "id": "msg-123",
  "sessionId": "session-123",
  "type": "text",
  "sender": "assistant",
  "content": "这是银杏，一种非常古老的植物...",
  "timestamp": "2025-12-19T10:00:00Z",
  "isStreaming": false
}
```

## 数据格式

### StreamEvent 事件格式

**connected事件**:
```json
{
  "type": "connected",
  "sessionId": "session-123"
}
```

**message事件**（文本流式）:
```json
{
  "type": "message",
  "content": "文本片段",
  "index": 0,
  "sessionId": "session-123",
  "messageId": "msg-123"
}
```

**image_progress事件**:
```json
{
  "type": "image_progress",
  "progress": 50,
  "sessionId": "session-123",
  "messageId": "msg-123"
}
```

**image_done事件**:
```json
{
  "type": "image_done",
  "content": {
    "url": "https://example.com/image.jpg"
  },
  "sessionId": "session-123",
  "messageId": "msg-123"
}
```

**card事件**:
```json
{
  "type": "card",
  "content": {
    "id": "card-1",
    "type": "science",
    "title": "银杏的科学知识",
    "content": {
      "name": "银杏",
      "explanation": "银杏是非常古老的植物...",
      "facts": ["事实1", "事实2"],
      "funFact": "趣味知识"
    },
    "explorationId": "exp-123"
  },
  "sessionId": "session-123",
  "messageId": "msg-123"
}
```

**error事件**:
```json
{
  "type": "error",
  "content": {
    "code": 500,
    "message": "内部服务器错误",
    "detail": "详细错误信息"
  },
  "sessionId": "session-123"
}
```

**done事件**:
```json
{
  "type": "done",
  "sessionId": "session-123"
}
```

## 上下文管理

### 上下文窗口限制

- 默认最大上下文轮次：20轮（40条消息：用户+助手）
- 超过20轮时，自动删除最早的消息，只保留最近20轮
- 上下文消息按时间顺序排列，确保对话连贯性

### 上下文消息格式（Eino）

后端将内部消息格式转换为Eino Message格式：

```go
// 系统消息（包含年级信息）
schema.SystemMessage("你是一个面向8岁学生的AI助手...")

// 用户消息
schema.UserMessage("这是什么？")

// 助手消息
schema.AssistantMessage("这是银杏...", nil)

// 消息列表（最多20轮）
messages := []*schema.Message{
    schema.SystemMessage(...),
    schema.UserMessage(...),
    schema.AssistantMessage(...),
    // ... 最多20轮
}
```

## 错误处理

### 错误码

- `400`: 请求参数错误
- `401`: 认证失败
- `404`: 会话不存在
- `500`: 服务器内部错误
- `503`: 服务不可用（AI模型服务不可用）

### 错误响应格式

**SSE错误事件**:
```
event: error
data: {"type":"error","content":{"code":500,"message":"内部服务器错误","detail":"..."}}
```

**HTTP错误响应**:
```json
{
  "code": 500,
  "message": "内部服务器错误",
  "detail": "详细错误信息"
}
```

## 性能要求

- 流式回答启动时间: <1秒
- 文本流式输出延迟: <50ms/字符
- 图片生成进度更新: 每200ms更新一次
- 上下文窗口查询: <10ms

## 安全考虑

- 会话ID验证：确保sessionId有效
- 消息长度限制：单条消息最大10KB
- 上下文大小限制：最多20轮，防止内存溢出
- 并发连接限制：单个sessionId最多1个SSE连接

## 兼容性

- 支持现代浏览器（Chrome 90+, Safari 14+, Firefox 88+）
- 必须支持Server-Sent Events (SSE)
- 不支持SSE的浏览器可以使用非流式接口

