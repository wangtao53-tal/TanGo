# API 合约文档

## 概述

本文档定义了优化识别流程与对话体验功能的API接口规范，确保前后端数据格式一致。

## API 端点

### 1. 对话消息处理

**端点**: `POST /api/conversation/message`

**请求**:
```json
{
  "sessionId": "可选，如果为空则创建新会话",
  "message": "用户消息内容",
  "messageType": "text|image|voice",
  "image": "可选，图片数据（base64或URL）",
  "identificationContext": {
    "objectName": "对象名称",
    "objectCategory": "对象类别",
    "confidence": 0.95,
    "keywords": ["关键词1", "关键词2"],
    "age": 8
  }
}
```

**响应**:
```json
{
  "sessionId": "会话ID",
  "message": {
    "id": "消息ID",
    "sessionId": "会话ID",
    "type": "text|image|voice|card",
    "sender": "user|assistant",
    "content": "消息内容",
    "timestamp": "2025-12-18T10:00:00Z",
    "isStreaming": false
  },
  "useStreaming": true
}
```

**说明**:
- 如果 `useStreaming` 为 `true`，客户端应建立SSE连接接收流式响应
- `identificationContext` 用于关联识别结果，支持多轮对话

### 2. 流式返回（SSE）

**端点**: `GET /api/conversation/stream/:sessionId`

**响应格式** (Server-Sent Events):
```
event: message_start
data: {"type":"message_start","messageId":"msg-123","sessionId":"session-456"}

event: message_delta
data: {"type":"message_delta","data":"这是","messageId":"msg-123","sessionId":"session-456"}

event: message_delta
data: {"type":"message_delta","data":"流式","messageId":"msg-123","sessionId":"session-456"}

event: message_end
data: {"type":"message_end","messageId":"msg-123","sessionId":"session-456"}
```

**事件类型**:
- `message_start`: 消息开始
- `message_delta`: 消息增量内容
- `message_end`: 消息结束
- `card_start`: 卡片开始
- `card_delta`: 卡片增量内容
- `card_end`: 卡片结束
- `error`: 错误事件

### 3. 意图识别

**端点**: `POST /api/conversation/intent`

**请求**:
```json
{
  "message": "帮我生成卡片",
  "sessionId": "session-123"
}
```

**响应**:
```json
{
  "intent": "generate_cards|text_response|other",
  "confidence": 0.95
}
```

### 4. 获取会话历史

**端点**: `GET /api/conversation/history/:sessionId`

**响应**:
```json
[
  {
    "id": "msg-1",
    "sessionId": "session-123",
    "type": "text",
    "sender": "user",
    "content": "这是什么？",
    "timestamp": "2025-12-18T10:00:00Z"
  },
  {
    "id": "msg-2",
    "sessionId": "session-123",
    "type": "text",
    "sender": "assistant",
    "content": "这是一只猫。",
    "timestamp": "2025-12-18T10:00:01Z"
  }
]
```

## 数据格式一致性

### 前端类型定义

参考 `frontend/src/types/conversation.ts` 和 `frontend/src/types/api.ts`

### 后端类型定义

参考 `backend/internal/types/types.go`

### 关键一致性要求

1. **消息ID格式**: 使用UUID格式，前后端保持一致
2. **时间戳格式**: 使用ISO 8601格式（RFC3339）
3. **消息类型**: 使用预定义枚举值，前后端保持一致
4. **卡片内容**: 图片字段保留但不显示，前后端数据结构一致

## 错误处理

所有API错误响应格式：
```json
{
  "code": 400,
  "message": "错误描述",
  "detail": "可选，详细错误信息"
}
```

常见错误码：
- `400`: 请求参数错误
- `401`: 未授权
- `404`: 资源不存在
- `500`: 服务器内部错误
- `503`: 服务不可用

## 性能要求

- 对话消息处理响应时间: ≤2秒（90%请求）
- 流式返回首字符延迟: ≤1秒（90%请求）
- 支持并发会话数: 200+
- 消息展示延迟: ≤0.5秒（前端乐观更新）

## 版本历史

- **v1.1.0** (2025-12-18): 优化识别流程，支持识别结果上下文，优化流式返回
