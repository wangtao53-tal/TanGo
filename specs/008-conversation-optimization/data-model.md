# 数据模型：对话体验与性能优化

**功能**: 008-conversation-optimization  
**创建日期**: 2025-12-20

## 概述

本文档定义对话体验与性能优化功能涉及的关键数据实体及其关系。本功能主要优化现有实体，新增流式消息状态和文本转语音状态。

## 核心实体

### 1. 流式消息（Streaming Message）

**用途**: 扩展对话消息，支持流式传输和实时渲染

**属性**:
- `id` (string, required): 消息唯一标识符（UUID）
- `sessionId` (string, required): 所属会话ID
- `type` (string, required): 消息类型（"text" | "image" | "voice" | "card"）
- `sender` (string, required): 发送者（"user" | "assistant"）
- `content` (string | object, required): 消息内容
  - 文本消息: string（可能包含Markdown格式）
  - 图片消息: string (base64或URL)
  - 语音消息: string (转录文本)
  - 卡片消息: KnowledgeCard对象
- `timestamp` (string, required): 消息时间戳（ISO 8601格式）
- `isStreaming` (boolean, optional): 是否正在流式返回（仅系统消息）
- `streamingText` (string, optional): 流式传输中的累积文本（仅系统消息，用于实时渲染）
- `markdown` (boolean, optional): 内容是否包含Markdown格式（仅文本消息）

**关系**:
- 与 `ConversationSession` 多对一关联（多条消息属于一个会话）

**验证规则**:
- `id` 必须是有效的UUID格式
- `type` 必须是预定义值之一
- `sender` 必须是预定义值之一
- `content` 不能为空
- `streamingText` 仅在 `isStreaming` 为 true 时存在

**状态转换**:
- 创建 → 流式传输中（isStreaming=true） → 流式完成（isStreaming=false） → 已展示

**新增字段说明**:
- `streamingText`: 用于实时渲染，每次接收到新的文本片段时更新
- `markdown`: 标识内容是否为Markdown格式，前端据此决定是否使用Markdown渲染

---

### 2. 流式事件（Stream Event）

**用途**: SSE流式传输中的事件数据结构

**属性**:
- `type` (string, required): 事件类型（"connected" | "message" | "card" | "image" | "done" | "error"）
- `content` (string | object, optional): 事件内容
  - message类型: string（文本片段）
  - card类型: CardContent对象
  - image类型: string（图片URL）
  - error类型: Error对象
- `index` (number, optional): 消息索引（用于message类型，标识字符位置）
- `sessionId` (string, required): 会话ID
- `messageId` (string, optional): 消息ID（用于关联消息）

**关系**:
- 与 `ConversationMessage` 多对一关联（多个事件对应一条消息）

**验证规则**:
- `type` 必须是预定义值之一
- `content` 在message、card、image类型时必须存在
- `index` 在message类型时存在

**事件类型说明**:
- `connected`: 连接建立
- `message`: 文本片段（逐字符或逐段）
- `card`: 知识卡片（流式返回时）
- `image`: 图片（流式返回时）
- `done`: 流式传输完成
- `error`: 错误事件

---

### 3. 知识卡片（Knowledge Card）- 扩展

**用途**: 存储生成的知识卡片内容，新增文本转语音支持

**属性**（继承现有属性，新增）:
- `id` (string, required): 卡片唯一标识符
- `type` (string, required): 卡片类型（"science" | "poetry" | "english"）
- `title` (string, required): 卡片标题
- `content` (object, required): 卡片内容对象
  - `text` (string, required): 文本内容（用于文本转语音）
  - `image` (string, optional): 图片URL或base64（对话中不显示）
- `explorationId` (string, optional): 关联的探索ID
- `createdAt` (string, required): 创建时间（ISO 8601格式）
- `audioText` (string, optional): 用于文本转语音的文本内容（提取自content，支持多语言）

**关系**:
- 与 `ConversationMessage` 可选关联（如果消息类型为"card"）

**验证规则**:
- `type` 必须是预定义值之一
- `title` 不能为空
- `content.text` 不能为空（用于文本转语音）

**新增字段说明**:
- `audioText`: 从卡片内容中提取的纯文本，用于文本转语音，支持中英文混合

---

### 4. 文本转语音状态（Text-to-Speech State）

**用途**: 管理知识卡片的文本转语音播放状态

**属性**:
- `cardId` (string, required): 关联的卡片ID
- `isPlaying` (boolean, required): 是否正在播放
- `isPaused` (boolean, required): 是否暂停
- `language` (string, required): 当前播放语言（"zh-CN" | "en-US"）
- `rate` (number, optional): 语速（0.1-10，默认0.9）
- `pitch` (number, optional): 音调（0-2，默认1.0）
- `currentCardId` (string, optional): 当前正在播放的卡片ID（全局状态）

**关系**:
- 与 `KnowledgeCard` 一对一关联（一个卡片对应一个播放状态）

**验证规则**:
- `cardId` 必须是有效的卡片ID
- `language` 必须是支持的语言之一
- `rate` 必须在 0.1-10 范围内
- `pitch` 必须在 0-2 范围内

**状态转换**:
- 初始 → 播放中（isPlaying=true） → 暂停（isPaused=true） → 继续播放 → 停止（isPlaying=false）

---

### 5. 性能指标（Performance Metrics）

**用途**: 记录性能优化相关的指标

**属性**:
- `requestId` (string, required): 请求唯一标识符
- `endpoint` (string, required): API端点（如"/api/explore/generate-cards"）
- `startTime` (number, required): 请求开始时间（Unix时间戳，毫秒）
- `endTime` (number, optional): 请求结束时间（Unix时间戳，毫秒）
- `duration` (number, optional): 请求持续时间（毫秒）
- `status` (string, required): 请求状态（"success" | "timeout" | "error"）
- `cardCount` (number, optional): 生成的卡片数量（仅generate-cards接口）
- `firstCardTime` (number, optional): 首张卡片生成时间（毫秒，流式返回时）

**关系**:
- 与 `GenerateCardsRequest` 一对一关联（一个请求对应一个指标）

**验证规则**:
- `endpoint` 必须是有效的API端点
- `duration` 必须大于0
- `status` 必须是预定义值之一

**用途**:
- 监控性能优化效果
- 分析性能瓶颈
- 生成性能报告

---

## 数据流

### 流式消息实时渲染流程

```
用户发送消息
  ↓
前端创建消息对象（isStreaming=false）
  ↓
发送到后端
  ↓
后端开始流式返回（SSE）
  ↓
前端接收第一个事件（type="message"）
  ↓
立即更新消息对象（isStreaming=true, streamingText=content）
  ↓
前端实时渲染（每次接收到新片段立即更新）
  ↓
接收完成事件（type="done"）
  ↓
更新消息对象（isStreaming=false, content=streamingText）
```

### Markdown渲染流程

```
接收流式文本片段
  ↓
累积到streamingText
  ↓
检测是否包含Markdown格式（通过markdown字段或内容检测）
  ↓
使用react-markdown渲染
  ↓
实时更新渲染结果
  ↓
流式完成，最终渲染
```

### 知识卡片生成优化流程

```
用户请求生成卡片
  ↓
后端并行生成三张卡片（goroutine）
  ↓
每生成完一张卡片，立即通过SSE返回（流式返回）
  ↓
前端接收到卡片事件，立即显示
  ↓
所有卡片生成完成，返回done事件
```

### 文本转语音流程

```
用户点击"听"按钮
  ↓
提取卡片audioText
  ↓
检测语言（中文/英文）
  ↓
创建SpeechSynthesisUtterance
  ↓
设置语言、语速、音调
  ↓
开始播放
  ↓
用户控制（播放/暂停/停止）
```

---

## 前后端数据一致性

### 流式事件格式（SSE）

**后端发送格式**:
```
event: message
data: {"type":"message","content":"文本片段","index":0,"sessionId":"xxx","messageId":"xxx"}

event: card
data: {"type":"card","content":{CardContent对象},"sessionId":"xxx"}

event: done
data: {"type":"done","sessionId":"xxx","messageId":"xxx"}

event: error
data: {"type":"error","content":{"message":"错误信息"},"sessionId":"xxx"}
```

**前端接收格式**:
- 与后端格式完全一致
- 使用TypeScript类型定义确保类型安全

### 知识卡片生成响应格式

**当前格式（同步）**:
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

**优化后格式（流式）**:
- 通过SSE事件流式返回
- 每个卡片作为一个card事件
- 最后发送done事件

### 对话消息格式

**扩展字段**:
```json
{
  "id": "msg-xxx",
  "type": "text",
  "sender": "assistant",
  "content": "消息内容",
  "timestamp": "2025-12-20T10:00:00Z",
  "sessionId": "session-xxx",
  "isStreaming": false,
  "streamingText": "流式传输中的文本",
  "markdown": true
}
```

---

## 验证规则总结

### 流式消息验证

1. `isStreaming` 为 true 时，`streamingText` 必须存在
2. `isStreaming` 为 false 时，`content` 必须包含完整内容
3. `markdown` 为 true 时，`content` 或 `streamingText` 必须包含有效的Markdown格式

### 知识卡片验证

1. `content.text` 必须存在（用于文本转语音）
2. `audioText` 如果存在，必须是从 `content` 提取的有效文本

### 性能指标验证

1. `duration` 必须大于0
2. `endTime` 必须大于 `startTime`
3. `status` 为 "success" 时，`duration` 必须存在

---

## 数据迁移

### 现有数据兼容性

- 现有 `ConversationMessage` 结构保持不变，新增字段为可选
- 现有 `KnowledgeCard` 结构保持不变，新增字段为可选
- 前端需要处理字段不存在的情况（向后兼容）

### 迁移步骤

1. **阶段1**：后端支持新字段（可选），前端开始使用新字段
2. **阶段2**：前端完全迁移到新字段，后端确保新字段存在
3. **阶段3**：移除旧字段（如果需要）

---

## 总结

本功能主要扩展了现有数据模型，新增了流式消息状态、Markdown标识、文本转语音状态等字段。所有新增字段都是可选的，确保向后兼容。前后端通过统一的类型定义和API契约保证数据一致性。
