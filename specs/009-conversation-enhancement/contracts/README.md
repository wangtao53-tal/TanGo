# API 合约文档

**创建日期**: 2025-12-21  
**功能**: 对话页面完善 - Agent模型流式返回

## 概述

本文档定义了统一流式对话接口的API规范。统一接口支持文本输入、语音输入、图片输入三种输入方式，通过 `messageType` 字段明确指定输入类型。

## API 接口

### 1. 统一流式对话接口

**接口**: `POST /api/conversation/stream`

**描述**: 统一流式对话接口，支持文本、语音、图片三种输入方式，通过 `messageType` 字段明确指定输入类型，返回Agent模型的流式响应。

**请求参数**:

```json
{
  "messageType": "text|voice|image (必填)",
  "message": "string (当messageType为text时必填)",
  "audio": "string (当messageType为voice时必填，base64编码的音频数据)",
  "image": "string (当messageType为image时必填，base64编码的图片数据或URL)",
  "sessionId": "string (可选)",
  "identificationContext": {
    "objectName": "string",
    "objectCategory": "string",
    "confidence": 0.0,
    "keywords": ["string"],
    "age": 0
  } (可选),
  "userAge": 0 (可选，3-18岁),
  "maxContextRounds": 20 (可选，默认20)
}
```

**字段说明**:
- `messageType`（必填）：明确指定输入类型
  - `"text"`: 文本输入，必须包含 `message` 字段
  - `"voice"`: 语音输入，必须包含 `audio` 字段
  - `"image"`: 图片输入，必须包含 `image` 字段

**响应格式**: SSE流式响应

**事件类型**:
- `connected`: 连接建立
- `voice_recognized`: 语音识别完成（仅当messageType为voice时，包含识别的文本）
- `image_uploaded`: 图片上传完成（仅当messageType为image时，包含图片URL）
- `message`: 文本消息（逐字符）
- `error`: 错误事件
- `done`: 流式完成

**示例请求**:

```bash
# 文本输入
curl -X POST http://localhost:8888/api/conversation/stream \
  -H "Content-Type: application/json" \
  -d '{
    "messageType": "text",
    "message": "这是什么？",
    "sessionId": "session-123",
    "userAge": 8
  }'

# 语音输入
curl -X POST http://localhost:8888/api/conversation/stream \
  -H "Content-Type: application/json" \
  -d '{
    "messageType": "voice",
    "audio": "UklGRiQAAABXQVZFZm10...",
    "sessionId": "session-123",
    "userAge": 8
  }'

# 图片输入
curl -X POST http://localhost:8888/api/conversation/stream \
  -H "Content-Type: application/json" \
  -d '{
    "messageType": "image",
    "image": "data:image/jpeg;base64,/9j/4AAQ...",
    "message": "请介绍一下这张图片",
    "sessionId": "session-123",
    "userAge": 8
  }'
```

## SSE 事件格式

### 连接建立事件

```
event: connected
data: {"type":"connected","sessionId":"session-123"}
```

### 文本消息事件

```
event: message
data: {"type":"message","content":"你","index":0,"sessionId":"session-123","messageId":"msg-123","markdown":false}
```

### 错误事件

```
event: error
data: {"type":"error","content":{"message":"Agent模型调用失败","error":"connection timeout"},"sessionId":"session-123"}
```

### 完成事件

```
event: done
data: {"type":"done","sessionId":"session-123","messageId":"msg-123"}
```

## 错误处理

### 错误码

- `400`: 请求参数错误
- `500`: 服务器内部错误
- `503`: 服务不可用（Agent模型调用失败）

### 错误响应格式

```json
{
  "type": "error",
  "content": {
    "message": "错误消息",
    "error": "错误详情（可选）",
    "code": 500
  },
  "sessionId": "session-123"
}
```

## 注意事项

1. **统一接口设计**: 对话页面使用一个统一的流式接口，通过 `messageType` 字段区分输入类型
2. **messageType字段验证**: 必须验证 `messageType` 字段，确保与对应字段匹配
3. **禁止Mock数据**: 统一接口必须调用真实的Agent模型，禁止使用Mock数据
4. **错误处理**: Agent模型调用失败时，必须记录详细错误日志，向用户发送错误事件
5. **流式响应**: 统一接口返回SSE流式响应，支持实时渲染
6. **多模态支持**: 图片输入支持多模态消息（图片+文本），语音输入先转换为文本
7. **独立卡片接口**: 生成知识卡片使用独立的接口，不是统一对话接口

## 向后兼容

- 保留原有的 `/api/conversation/voice` 和 `/api/upload/image` 接口（返回JSON响应）
- 统一流式接口 `/api/conversation/stream` 支持所有输入类型
- 前端统一调用 `/api/conversation/stream` 接口，根据输入类型设置 `messageType` 字段

