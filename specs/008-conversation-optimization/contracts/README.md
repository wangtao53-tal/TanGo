# API契约：对话体验与性能优化

**功能**: 008-conversation-optimization  
**创建日期**: 2025-12-20

## 概述

本文档定义对话体验与性能优化功能的API契约。本功能主要优化现有API，新增流式返回支持，确保前后端数据格式一致。

## API变更说明

### 1. 流式消息实时渲染

**变更内容**：
- 扩展 `ConversationMessage` 类型，新增 `streamingText` 和 `markdown` 字段
- 优化SSE事件格式，确保实时渲染

**影响范围**：
- `/api/conversation/stream` - 流式对话接口
- `/api/conversation/message` - 对话消息接口（如果使用流式返回）

### 2. Markdown格式支持

**变更内容**：
- `ConversationMessage` 新增 `markdown` 字段，标识内容是否为Markdown格式
- 流式消息内容支持Markdown格式

**影响范围**：
- 所有返回文本消息的接口

### 3. 知识卡片生成性能优化

**变更内容**：
- `/api/explore/generate-cards` 支持流式返回
- 每生成完一张卡片立即返回，不等待所有卡片完成

**影响范围**：
- `/api/explore/generate-cards` - 知识卡片生成接口

### 4. 文本转语音

**变更内容**：
- 前端功能，无需后端API变更
- `KnowledgeCard` 类型新增 `audioText` 字段（可选）

**影响范围**：
- 无后端API变更

---

## API详细定义

### 扩展的对话消息类型

```typescript
interface ConversationMessage {
  id: string;                    // 消息ID
  type: 'text' | 'image' | 'voice' | 'card';  // 消息类型
  sender: 'user' | 'assistant';  // 发送者
  content: string | object;       // 消息内容
  timestamp: string;              // 时间戳（ISO 8601）
  sessionId?: string;             // 会话ID
  isStreaming?: boolean;          // 是否正在流式返回（新增）
  streamingText?: string;         // 流式传输中的累积文本（新增）
  markdown?: boolean;              // 内容是否包含Markdown格式（新增）
}
```

### 扩展的流式事件类型

```typescript
interface StreamEvent {
  type: 'connected' | 'message' | 'card' | 'image' | 'done' | 'error';
  content?: string | object;     // 事件内容
  index?: number;                 // 消息索引（用于message类型）
  sessionId: string;              // 会话ID
  messageId?: string;             // 消息ID
  markdown?: boolean;             // 内容是否为Markdown格式（新增）
}
```

### 扩展的知识卡片类型

```typescript
interface KnowledgeCard {
  id: string;
  type: 'science' | 'poetry' | 'english';
  title: string;
  content: {
    text: string;                 // 文本内容（用于文本转语音）
    image?: string;               // 图片URL（可选）
  };
  audioText?: string;              // 用于文本转语音的文本（新增，可选）
  explorationId?: string;
  createdAt: string;
}
```

---

## API端点详细说明

### 1. 流式对话接口（优化）

**端点**: `POST /api/conversation/stream`

**请求格式**:
```json
{
  "sessionId": "session-xxx",
  "message": "用户消息",
  "messageType": "text"
}
```

**响应格式（SSE事件流）**:

**连接建立事件**:
```
event: connected
data: {"type":"connected","sessionId":"session-xxx","messageId":"msg-xxx"}
```

**文本片段事件**:
```
event: message
data: {"type":"message","content":"文本片段","index":0,"sessionId":"session-xxx","messageId":"msg-xxx","markdown":true}
```

**卡片事件**（流式返回时）:
```
event: card
data: {"type":"card","content":{CardContent对象},"sessionId":"session-xxx","messageId":"msg-xxx"}
```

**完成事件**:
```
event: done
data: {"type":"done","sessionId":"session-xxx","messageId":"msg-xxx"}
```

**错误事件**:
```
event: error
data: {"type":"error","content":{"message":"错误信息"},"sessionId":"session-xxx"}
```

**变更说明**:
- 新增 `markdown` 字段，标识文本内容是否为Markdown格式
- 确保每个文本片段事件立即发送，不累积

---

### 2. 知识卡片生成接口（优化）

**端点**: `POST /api/explore/generate-cards`

**请求格式**（不变）:
```json
{
  "objectName": "对象名称",
  "objectCategory": "自然类",
  "age": 8,
  "keywords": ["关键词1", "关键词2"]
}
```

**响应格式（两种模式）**:

**模式1：同步返回（保持兼容）**
```json
{
  "cards": [
    {
      "type": "science",
      "title": "标题",
      "content": {...}
    }
  ]
}
```

**模式2：流式返回（新增，推荐）**

通过SSE流式返回，每个卡片作为一个事件：

**首张卡片事件**:
```
event: card
data: {"type":"card","content":{"type":"science","title":"标题","content":{...}},"sessionId":"session-xxx","index":0}
```

**第二张卡片事件**:
```
event: card
data: {"type":"card","content":{"type":"poetry","title":"标题","content":{...}},"sessionId":"session-xxx","index":1}
```

**第三张卡片事件**:
```
event: card
data: {"type":"card","content":{"type":"english","title":"标题","content":{...}},"sessionId":"session-xxx","index":2}
```

**完成事件**:
```
event: done
data: {"type":"done","sessionId":"session-xxx"}
```

**变更说明**:
- 新增流式返回模式，通过查询参数 `stream=true` 启用
- 流式返回时，每生成完一张卡片立即返回
- 保持同步返回模式，确保向后兼容

**使用示例**:
```typescript
// 同步返回（默认）
POST /api/explore/generate-cards

// 流式返回（新增）
POST /api/explore/generate-cards?stream=true
```

---

### 3. 对话消息接口（扩展）

**端点**: `POST /api/conversation/message`

**请求格式**（不变）:
```json
{
  "sessionId": "session-xxx",
  "message": "用户消息",
  "messageType": "text"
}
```

**响应格式**（扩展）:
```json
{
  "sessionId": "session-xxx",
  "message": {
    "id": "msg-xxx",
    "type": "text",
    "sender": "assistant",
    "content": "消息内容",
    "timestamp": "2025-12-20T10:00:00Z",
    "sessionId": "session-xxx",
    "isStreaming": false,
    "markdown": true
  },
  "useStreaming": true
}
```

**变更说明**:
- 响应消息新增 `markdown` 字段
- `useStreaming` 字段建议使用流式返回

---

## 前后端数据一致性

### 类型定义一致性

**后端（Go）**:
```go
type ConversationMessage struct {
    Id          string      `json:"id"`
    Type        string      `json:"type"`
    Sender      string      `json:"sender"`
    Content     interface{} `json:"content"`
    Timestamp   string      `json:"timestamp"`
    SessionId   string      `json:"sessionId,optional"`
    IsStreaming *bool       `json:"isStreaming,optional"`
    StreamingText string    `json:"streamingText,optional"`  // 新增
    Markdown    *bool       `json:"markdown,optional"`        // 新增
}
```

**前端（TypeScript）**:
```typescript
interface ConversationMessage {
  id: string;
  type: 'text' | 'image' | 'voice' | 'card';
  sender: 'user' | 'assistant';
  content: string | object;
  timestamp: string;
  sessionId?: string;
  isStreaming?: boolean;
  streamingText?: string;  // 新增
  markdown?: boolean;      // 新增
}
```

### 字段映射表

| 字段名 | 后端类型 | 前端类型 | 说明 |
|--------|---------|---------|------|
| id | string | string | 消息ID |
| type | string | 'text'\|'image'\|'voice'\|'card' | 消息类型 |
| sender | string | 'user'\|'assistant' | 发送者 |
| content | interface{} | string\|object | 消息内容 |
| timestamp | string | string | 时间戳 |
| sessionId | string (optional) | string? | 会话ID |
| isStreaming | *bool (optional) | boolean? | 是否流式返回 |
| streamingText | string (optional) | string? | 流式文本（新增） |
| markdown | *bool (optional) | boolean? | Markdown标识（新增） |

---

## 错误处理

### 错误响应格式

```json
{
  "code": 400,
  "message": "错误信息",
  "detail": "错误详情（可选）"
}
```

### 常见错误码

- `400`: 请求参数错误
- `500`: 服务器内部错误
- `503`: 服务不可用（如AI模型服务不可用）
- `504`: 网关超时（如AI模型调用超时）

### 流式传输错误处理

**网络中断**:
- 前端保存已接收的内容
- 显示错误提示
- 提供重连机制

**超时错误**:
- 后端设置合理的超时时间（如每张卡片10秒）
- 超时后返回已生成的卡片
- 前端显示部分结果和错误提示

---

## 性能要求

### 响应时间目标

- **流式消息首字符延迟**: ≤1秒（90%请求）
- **流式消息实时渲染**: 接收到数据后100毫秒内更新UI
- **知识卡片生成**: ≤5秒（95%请求）
- **流式返回首张卡片**: ≤2秒（90%请求）

### 并发要求

- 支持200+并发对话会话
- 支持50+并发知识卡片生成请求

---

## 向后兼容性

### 兼容策略

1. **新增字段为可选**：所有新增字段都标记为 `optional`，确保向后兼容
2. **保持现有接口**：不删除或修改现有接口，只扩展
3. **渐进式迁移**：前端可以逐步迁移到新字段，不影响现有功能

### 迁移路径

1. **阶段1**：后端支持新字段（可选），前端开始使用
2. **阶段2**：前端完全迁移到新字段
3. **阶段3**：移除旧字段（如果需要，未来版本）

---

## 测试建议

### 接口测试

1. **流式消息实时渲染测试**：
   - 验证每个文本片段事件是否立即发送
   - 验证前端是否实时更新UI
   - 验证网络中断时的处理

2. **Markdown格式测试**：
   - 验证Markdown内容是否正确渲染
   - 验证流式Markdown渲染是否同步

3. **性能测试**：
   - 验证知识卡片生成是否在5秒内完成
   - 验证并发请求的处理能力
   - 验证流式返回的首张卡片延迟

4. **数据一致性测试**：
   - 验证前后端类型定义一致性
   - 验证字段映射正确性
   - 验证错误处理一致性

---

## 总结

本功能主要扩展了现有API，新增了流式消息实时渲染、Markdown格式支持、知识卡片流式返回等功能。所有变更都保持向后兼容，确保现有功能不受影响。前后端通过统一的类型定义和API契约保证数据一致性。
