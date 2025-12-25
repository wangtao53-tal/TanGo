# API Contracts: 多Agent追问功能优化

**Date**: 2025-01-27  
**Feature**: 多Agent追问功能优化

## API Endpoints

### POST /api/conversation/agent

多Agent模式流式对话接口，支持文本、语音、图片三种输入方式。

**请求格式**: `UnifiedStreamConversationRequest`

**响应格式**: SSE流式响应（`StreamEvent`）

**请求参数**:
- `messageType` (string, 必填): 消息类型，值为 `"text"`、`"voice"` 或 `"image"`
- `message` (string, 条件必填): 文本消息，当 `messageType` 为 `"text"` 时必填
- `audio` (string, 条件必填): 语音数据（base64），当 `messageType` 为 `"voice"` 时必填
- `image` (string, 条件必填): 图片数据（base64或URL），当 `messageType` 为 `"image"` 时必填
- `sessionId` (string, 可选): 会话ID，如果为空则创建新会话
- `identificationContext` (IdentificationContext, 可选): 识别结果上下文
- `userAge` (int, 可选): 用户年龄（3-18岁），用于内容适配
- `maxContextRounds` (int, 可选): 最大上下文轮次，默认20轮

**响应事件**:
- `connected`: 连接成功事件
- `message`: 消息事件（流式文本内容）
- `done`: 完成事件
- `error`: 错误事件

**多Agent执行流程**:
1. Supervisor接收请求，分析上下文
2. Intent Agent识别意图类型
3. Cognitive Load Agent判断认知负载
4. Learning Planner Agent决定教学动作
5. Domain Agent（Science/Language/Humanities）生成专业回答
6. Interaction Agent优化交互方式
7. Reflection Agent反思学习状态
8. Memory Agent记录学习状态
9. 流式返回最终回答

**接口一致性**:
- 请求格式与 `/api/conversation/stream` 接口完全一致
- 响应格式与 `/api/conversation/stream` 接口完全一致
- 前端可以无缝切换接口调用

## 类型定义

### UnifiedStreamConversationRequest

统一流式对话请求类型，支持文本、语音、图片三种输入方式。

### IdentificationContext

识别结果上下文，包含对象名称、类别、置信度等信息。

### StreamEvent

SSE流式事件类型，包含事件类型、内容、索引等信息。

## 错误处理

- 当多Agent模式执行失败时，可以降级到单Agent模式
- 错误事件包含详细的错误信息
- 前端可以实现自动降级机制

## 示例

### 文本输入请求

```json
{
  "messageType": "text",
  "message": "这是什么？",
  "sessionId": "session-123",
  "userAge": 8
}
```

### 语音输入请求

```json
{
  "messageType": "voice",
  "audio": "base64-encoded-audio-data",
  "sessionId": "session-123",
  "userAge": 8
}
```

### 图片输入请求

```json
{
  "messageType": "image",
  "image": "base64-encoded-image-data",
  "sessionId": "session-123",
  "userAge": 8
}
```

### SSE响应示例

```
event: connected
data: {"type":"connected","sessionId":"session-123"}

event: message
data: {"type":"message","content":"这","index":0,"sessionId":"session-123"}

event: message
data: {"type":"message","content":"是","index":1,"sessionId":"session-123"}

event: done
data: {"type":"done","sessionId":"session-123"}
```

