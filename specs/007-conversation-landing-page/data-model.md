# 数据模型：H5对话落地页

**功能**: 007-conversation-landing-page  
**创建日期**: 2025-12-19

## 概述

本文档定义H5对话落地页功能涉及的关键数据实体及其关系，包括对话消息、会话管理、上下文窗口等。

## 核心实体

### 1. 对话会话（Conversation Session）

**用途**: 维护对话会话状态，关联识别结果和对话消息，管理20轮上下文窗口

**属性**:
- `id` (string, required): 会话唯一标识符（UUID格式，如 `session-{timestamp}-{random}`）
- `messages` (ConversationMessage[], required): 消息列表（最多20轮，40条消息）
- `identificationContext` (IdentificationContext, optional): 关联的识别结果上下文
- `userAge` (number, required): 用户年龄（3-18岁，用于内容适配）
- `createdAt` (string, required): 会话创建时间（ISO 8601格式）
- `updatedAt` (string, required): 会话最后更新时间（ISO 8601格式）
- `messageCount` (number, required): 当前消息数量（用于判断是否超过20轮）

**关系**:
- 与 `IdentificationContext` 一对一关联（一个会话对应一个识别结果）
- 与 `ConversationMessage` 一对多关联（一个会话包含多条消息，最多40条）

**验证规则**:
- `id` 必须是有效的UUID格式或自定义格式
- `userAge` 必须在 3-18 范围内
- `messageCount` 不能超过 40（20轮对话）

**状态转换**:
- 创建 → active → closed（可选）

**特殊说明**:
- 消息列表自动维护最近20轮（40条消息），超过时删除最早的消息
- 识别结果上下文作为对话的初始上下文，影响AI生成内容

### 2. 对话消息（Conversation Message）

**用途**: 存储对话中的用户消息和AI助手响应，支持多种消息类型

**属性**:
- `id` (string, required): 消息唯一标识符（UUID格式）
- `sessionId` (string, required): 所属会话ID
- `type` (string, required): 消息类型
  - `"text"`: 文本消息
  - `"image"`: 图片消息
  - `"voice"`: 语音消息（转录文本）
  - `"card"`: 知识卡片消息
- `sender` (string, required): 发送者
  - `"user"`: 用户消息（显示在右侧）
  - `"assistant"`: AI助手消息（显示在左侧）
- `content` (string | object, required): 消息内容
  - 文本消息: `string` - 文本内容
  - 图片消息: `string` - 图片URL或base64
  - 语音消息: `string` - 转录文本
  - 卡片消息: `KnowledgeCard` 对象
- `timestamp` (string, required): 消息时间戳（ISO 8601格式）
- `isStreaming` (boolean, optional): 是否正在流式返回（仅AI消息）
- `streamingText` (string, optional): 流式文本内容（用于打字机效果）

**关系**:
- 与 `ConversationSession` 多对一关联（多条消息属于一个会话）
- 与 `KnowledgeCard` 可选关联（如果消息类型为"card"）

**验证规则**:
- `id` 必须是有效的UUID格式
- `type` 必须是预定义值之一
- `sender` 必须是预定义值之一
- `content` 不能为空
- `sessionId` 必须与所属会话的ID匹配

**状态转换**:
- 创建 → 已发送 → 已接收 → 已展示
- 流式消息: 创建 → 流式更新中 → 流式完成 → 已展示

**特殊说明**:
- 用户消息立即显示（乐观更新）
- AI消息支持流式更新，通过`isStreaming`和`streamingText`字段管理
- 卡片消息在对话中只显示文本内容，不显示图片

### 3. 识别结果上下文（Identification Context）

**用途**: 存储图像识别后的结果，作为对话的初始上下文

**属性**:
- `objectName` (string, required): 识别出的对象名称
- `objectCategory` (string, required): 对象类别（"自然类" | "生活类" | "人文类"）
- `confidence` (number, required): 识别置信度（0-1）
- `keywords` (string[], optional): 相关关键词列表
- `age` (number, optional): 用户年龄（用于内容适配，3-18岁）

**关系**:
- 与 `ConversationSession` 一对一关联（一个识别结果对应一个会话）

**验证规则**:
- `objectName` 不能为空
- `objectCategory` 必须是预定义值之一
- `confidence` 必须在 0-1 范围内
- `age` 必须在 3-18 范围内（如果提供）

**特殊说明**:
- 识别结果上下文影响AI生成内容的主题和风格
- 结合用户年龄，生成适配的内容难度

### 4. 知识卡片（Knowledge Card）

**用途**: 存储生成的知识卡片内容，在对话中展示

**属性**:
- `id` (string, required): 卡片唯一标识符
- `explorationId` (string, optional): 关联的探索记录ID
- `type` (string, required): 卡片类型
  - `"science"`: 主体卡片（科学认知）
  - `"poetry"`: 古诗词卡片
  - `"english"`: 英文知识卡片
- `title` (string, required): 卡片标题
- `content` (object, required): 卡片内容对象（根据类型不同结构不同）
  - 科学认知卡: `{ name, explanation, facts[], funFact }`
  - 古诗词卡: `{ poem, poemSource, explanation, context }`
  - 英文知识卡: `{ keywords[], expressions[], pronunciation }`
- `createdAt` (string, required): 创建时间（ISO 8601格式）

**关系**:
- 与 `ConversationMessage` 可选关联（如果消息类型为"card"）

**验证规则**:
- `type` 必须是预定义值之一
- `title` 不能为空
- `content` 必须符合对应类型的结构要求

**特殊说明**:
- 在对话消息列表中展示时，只渲染文本内容，不渲染图片
- 图片数据保留在数据结构中，但不显示（符合原则七）

### 5. 流式事件（Stream Event）

**用途**: SSE流式传输中的事件数据结构

**属性**:
- `type` (string, required): 事件类型
  - `"connected"`: 连接建立
  - `"message"`: 消息内容（文本、图片、卡片）
  - `"image_progress"`: 图片生成进度
  - `"image_done"`: 图片生成完成
  - `"error"`: 错误事件
  - `"done"`: 流式传输完成
- `content` (string | object, required): 事件内容
  - 文本消息: `string` - 文本片段
  - 图片进度: `{ progress: number }` - 进度百分比（0-100）
  - 图片完成: `{ url: string }` - 图片URL
  - 卡片消息: `KnowledgeCard` 对象
- `index` (number, optional): 文本消息的字符索引（用于打字机效果）
- `sessionId` (string, optional): 会话ID

**关系**:
- 与 `ConversationMessage` 关联（用于更新消息内容）

**验证规则**:
- `type` 必须是预定义值之一
- `content` 不能为空

**特殊说明**:
- SSE事件格式: `event: {type}\ndata: {JSON}\n\n`
- 文本消息逐字符发送，前端实时更新显示

## 数据流

### 首次拍照后知识卡片生成流程

```
用户拍照识别
  ↓
前端接收识别结果 (IdentificationResponse)
  ↓
前端跳转到对话页面，传递识别结果
  ↓
前端创建会话 (ConversationSession)，包含识别结果上下文
  ↓
前端自动进行意图识别（识别为"生成卡片"）
  ↓
前端调用卡片生成API
  ↓
后端并行生成三张知识卡片
  ↓
后端返回卡片数据 (GenerateCardsResponse)
  ↓
前端将卡片转换为消息 (ConversationMessage)，添加到消息列表
  ↓
前端展示三张卡片（只显示文本内容）
```

### 追问对话流程（流式输出）

```
用户在对话页面输入问题
  ↓
前端立即创建用户消息，添加到消息列表（乐观更新）
  ↓
前端调用流式对话API (SSE)
  ↓
后端获取最近20轮对话历史，转换为Eino Message格式
  ↓
后端根据用户年级生成系统prompt
  ↓
后端调用Eino ChatModel.Stream接口
  ↓
后端逐块读取流式数据，通过SSE发送到前端
  ↓
前端接收SSE事件，实时更新AI消息内容（打字机效果）
  ↓
如果包含图片生成，显示loading占位符
  ↓
图片生成完成后，替换占位符为实际图片
  ↓
流式传输完成，保存最终消息到历史记录
```

### 上下文窗口管理流程

```
用户发送新消息
  ↓
后端获取会话的所有消息
  ↓
检查消息数量是否超过40条（20轮）
  ↓
如果超过，只保留最近40条消息
  ↓
将消息转换为Eino Message格式（最多20轮）
  ↓
添加到系统prompt中，作为对话历史
  ↓
调用AI模型生成回答
```

## 前后端数据格式一致性

### 对话消息格式

**前端类型定义**:
```typescript
interface ConversationMessage {
  id: string;
  sessionId: string;
  type: 'text' | 'image' | 'voice' | 'card';
  sender: 'user' | 'assistant';
  content: string | KnowledgeCard;
  timestamp: string;
  isStreaming?: boolean;
  streamingText?: string; // 流式文本内容
}
```

**后端类型定义**:
```go
type ConversationMessage struct {
    ID            string      `json:"id"`
    SessionID     string      `json:"sessionId"`
    Type          string      `json:"type"`
    Sender        string      `json:"sender"`
    Content       interface{} `json:"content"`
    Timestamp     string      `json:"timestamp"`
    IsStreaming   *bool       `json:"isStreaming,omitempty"`
    StreamingText string      `json:"streamingText,omitempty"`
}
```

### SSE事件格式

**前端类型定义**:
```typescript
interface ConversationStreamEvent {
  type: 'connected' | 'message' | 'image_progress' | 'image_done' | 'error' | 'done';
  content: string | { progress?: number; url?: string } | KnowledgeCard;
  index?: number;
  sessionId?: string;
}
```

**后端SSE格式**:
```
event: message
data: {"type":"text","content":"文本片段","index":0}

event: image_progress
data: {"type":"image_progress","content":{"progress":50}}

event: image_done
data: {"type":"image_done","content":{"url":"https://..."}}

event: done
data: {"type":"done"}
```

### 上下文消息格式（Eino）

**后端Eino Message格式**:
```go
// 系统消息
schema.SystemMessage("系统prompt，包含年级信息")

// 用户消息
schema.UserMessage("用户问题")

// 助手消息
schema.AssistantMessage("AI回答", nil)

// 消息列表（最多20轮）
messages := []*schema.Message{
    schema.SystemMessage(...),
    schema.UserMessage(...),
    schema.AssistantMessage(...),
    // ... 最多20轮
}
```

## 状态管理

### 前端状态管理

- **会话状态**: 使用React state管理，包含sessionId、消息列表、识别结果上下文
- **消息列表状态**: 使用React state数组，支持实时更新和流式追加
- **流式状态**: 使用React state管理当前流式消息的文本内容
- **图片loading状态**: 使用React state管理图片生成进度和占位符显示

### 后端状态管理

- **会话状态**: 使用内存map存储（`storage/memory.go`），key为sessionId
- **消息列表**: 存储在会话对象中，自动维护最近20轮限制
- **并发控制**: 使用sync.RWMutex保护并发读写
- **流式连接**: 每个SSE连接独立管理，支持并发流式传输

## 数据持久化

### 前端持久化

- **对话历史**: 使用localStorage保存最近20轮对话（可选，MVP版本不必须）
- **识别结果**: 通过路由state传递，刷新后可能丢失（可接受）
- **会话ID**: 使用localStorage保存，刷新后恢复会话

### 后端持久化

- **会话状态**: 内存存储，服务重启后丢失（MVP版本可接受）
- **识别结果**: 不持久化，仅作为对话上下文使用
- **消息历史**: 存储在内存中，最多20轮，超过自动删除

## 数据验证

### 前端验证

- 使用TypeScript类型检查
- 运行时验证API响应格式
- SSE事件格式验证
- 错误处理和降级方案

### 后端验证

- 使用go-zero的验证机制
- 类型转换和格式检查
- 上下文窗口大小验证（最多20轮）
- 错误返回和日志记录

## 性能优化

### 前端优化

- 消息列表虚拟滚动（如果消息过多）
- 流式文本更新使用防抖/节流
- 图片懒加载和占位符优化
- React.memo优化组件渲染

### 后端优化

- 上下文消息转换缓存（如果可能）
- 流式输出缓冲区管理
- 并发连接数限制
- 内存使用监控和清理

