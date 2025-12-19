# 数据模型：优化识别流程与对话体验

**功能**: 004-optimize-identify-flow  
**创建日期**: 2025-12-18

## 概述

本文档定义优化识别流程与对话体验功能涉及的关键数据实体及其关系。

## 核心实体

### 1. 识别结果（Identification Result）

**用途**: 存储图像识别后的结果，作为对话的初始上下文

**属性**:
- `objectName` (string, required): 识别出的对象名称
- `objectCategory` (string, required): 对象类别（"自然类" | "生活类" | "人文类"）
- `confidence` (number, required): 识别置信度（0-1）
- `keywords` (string[], optional): 相关关键词列表
- `age` (number, optional): 用户年龄（用于内容适配）
- `timestamp` (string, required): 识别时间戳（ISO 8601格式）

**关系**:
- 与 `ConversationSession` 一对一关联（一个识别结果对应一个会话）

**验证规则**:
- `objectName` 不能为空
- `objectCategory` 必须是预定义值之一
- `confidence` 必须在 0-1 范围内

**状态转换**:
- 创建 → 已展示 → 已关联会话

### 2. 对话会话（Conversation Session）

**用途**: 维护对话会话状态，关联识别结果和对话消息

**属性**:
- `sessionId` (string, required): 会话唯一标识符（UUID）
- `identificationResult` (IdentificationResult, optional): 关联的识别结果
- `createdAt` (string, required): 会话创建时间（ISO 8601格式）
- `updatedAt` (string, required): 会话最后更新时间（ISO 8601格式）
- `messageCount` (number, required): 消息数量
- `status` (string, required): 会话状态（"active" | "closed"）

**关系**:
- 与 `IdentificationResult` 一对一关联
- 与 `ConversationMessage` 一对多关联（一个会话包含多条消息）

**验证规则**:
- `sessionId` 必须是有效的UUID格式
- `status` 必须是预定义值之一

**状态转换**:
- 创建 → active → closed

### 3. 对话消息（Conversation Message）

**用途**: 存储对话中的用户消息和系统响应

**属性**:
- `id` (string, required): 消息唯一标识符（UUID）
- `sessionId` (string, required): 所属会话ID
- `type` (string, required): 消息类型（"text" | "image" | "voice" | "card"）
- `sender` (string, required): 发送者（"user" | "assistant"）
- `content` (string | object, required): 消息内容
  - 文本消息: string
  - 图片消息: string (base64或URL)
  - 语音消息: string (转录文本)
  - 卡片消息: KnowledgeCard对象
- `timestamp` (string, required): 消息时间戳（ISO 8601格式）
- `isStreaming` (boolean, optional): 是否正在流式返回（仅系统消息）

**关系**:
- 与 `ConversationSession` 多对一关联（多条消息属于一个会话）
- 与 `KnowledgeCard` 可选关联（如果消息类型为"card"）

**验证规则**:
- `id` 必须是有效的UUID格式
- `type` 必须是预定义值之一
- `sender` 必须是预定义值之一
- `content` 不能为空

**状态转换**:
- 创建 → 已发送 → 已接收 → 已展示

### 4. 知识卡片（Knowledge Card）

**用途**: 存储生成的知识卡片内容（在对话中展示时只显示文本）

**属性**:
- `id` (string, required): 卡片唯一标识符
- `type` (string, required): 卡片类型（"science" | "poetry" | "english"）
- `title` (string, required): 卡片标题
- `content` (object, required): 卡片内容对象
  - `text` (string, required): 文本内容
  - `image` (string, optional): 图片URL或base64（对话中不显示）
- `explorationId` (string, optional): 关联的探索ID
- `createdAt` (string, required): 创建时间（ISO 8601格式）

**关系**:
- 与 `ConversationMessage` 可选关联（如果消息类型为"card"）

**验证规则**:
- `type` 必须是预定义值之一
- `title` 不能为空
- `content.text` 不能为空

**特殊说明**:
- 在对话消息列表中展示时，只渲染 `content.text`，不渲染 `content.image`
- 图片数据保留在数据结构中，但不显示

## 数据流

### 识别流程数据流

```
用户拍照
  ↓
前端调用识别API
  ↓
后端返回识别结果 (IdentificationResult)
  ↓
前端跳转到问答页面，传递识别结果
  ↓
问答页面创建初始系统消息，展示识别结果
  ↓
创建对话会话 (ConversationSession)，关联识别结果
```

### 对话流程数据流

```
用户发送消息
  ↓
前端立即创建用户消息 (ConversationMessage)，添加到消息列表（乐观更新）
  ↓
前端调用对话API，传递消息和sessionId
  ↓
后端处理消息，更新会话状态
  ↓
后端通过SSE流式返回系统响应
  ↓
前端接收流式数据，实时更新系统消息
```

### 卡片生成流程数据流

```
用户发送"生成卡片"请求
  ↓
前端立即显示用户消息
  ↓
后端识别意图为"generate_cards"
  ↓
后端调用AI模型生成三张卡片
  ↓
后端通过SSE流式返回卡片消息（每张卡片一条消息）
  ↓
前端接收卡片消息，只显示文本内容，不显示图片
```

## 前后端数据格式一致性

### 识别结果格式

**前端类型定义**:
```typescript
interface IdentifyResponse {
  objectName: string;
  objectCategory: '自然类' | '生活类' | '人文类';
  confidence: number;
  keywords?: string[];
}
```

**后端类型定义**:
```go
type IdentifyResponse struct {
    ObjectName     string   `json:"objectName"`
    ObjectCategory string   `json:"objectCategory"`
    Confidence     float64  `json:"confidence"`
    Keywords       []string `json:"keywords,omitempty"`
}
```

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
}
```

**后端类型定义**:
```go
type ConversationMessage struct {
    ID        string      `json:"id"`
    SessionID string      `json:"sessionId"`
    Type      string      `json:"type"`
    Sender    string      `json:"sender"`
    Content   interface{} `json:"content"`
    Timestamp string      `json:"timestamp"`
    IsStreaming *bool     `json:"isStreaming,omitempty"`
}
```

## 状态管理

### 前端状态管理

- **识别结果状态**: 使用React Router state传递，页面加载时读取
- **对话会话状态**: 使用React state管理，包含sessionId和消息列表
- **消息列表状态**: 使用React state数组，支持实时更新

### 后端状态管理

- **会话状态**: 使用内存map存储，key为sessionId，value为会话对象
- **并发控制**: 使用sync.RWMutex保护并发读写
- **会话清理**: 实现定时清理机制，移除过期会话

## 数据持久化

### 前端持久化

- **对话历史**: 可选使用localStorage或IndexedDB存储（MVP版本不必须）
- **识别结果**: 通过路由state传递，刷新后可能丢失（可接受）

### 后端持久化

- **会话状态**: 内存存储，服务重启后丢失（MVP版本可接受）
- **识别结果**: 不持久化，仅作为对话上下文使用

## 数据验证

### 前端验证

- 使用TypeScript类型检查
- 运行时验证API响应格式
- 错误处理和降级方案

### 后端验证

- 使用go-zero的验证机制
- 类型转换和格式检查
- 错误返回和日志记录
