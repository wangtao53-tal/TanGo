# 技术研究与决策文档

**创建日期**: 2025-12-21  
**功能**: 对话页面完善 - Agent模型流式返回

## 研究目标

实现统一流式接口，支持文本输入、语音输入、图片输入三种输入方式，通过 `messageType` 字段明确指定输入类型，替换Mock数据，使用真实的Agent模型（基于Eino Graph）。

## 技术研究

### 1. 统一接口设计

**当前状态**:
- 现有多个独立接口：文本对话接口、语音识别接口、图片上传接口
- 前端需要根据输入类型调用不同的接口

**实现方案**:
- **统一接口设计（已确定）**: 实现一个统一的流式接口，通过 `messageType` 字段明确指定输入类型
  - 优点：前端只需调用一个接口，简化调用逻辑；后端统一处理，便于维护
  - 实施：扩展 `StreamConversationHandler`，根据 `messageType` 字段处理不同输入类型

**决策**: 采用统一接口设计，通过 `messageType` 字段区分输入类型。

**技术细节**:
- 统一接口：`POST /api/conversation/stream`
- 请求格式：必须包含 `messageType` 字段（`"text"|"voice"|"image"`）
- 根据 `messageType` 的值：
  - `"text"` → 必须包含 `message` 字段，直接调用Agent模型
  - `"voice"` → 必须包含 `audio` 字段，先调用语音识别，再调用Agent模型
  - `"image"` → 必须包含 `image` 字段，先上传图片获取URL，再调用Agent模型（多模态）

### 2. 统一接口内部处理流程

**处理流程**:
1. 接收统一接口请求，验证 `messageType` 字段
2. 根据 `messageType` 处理：
   - `"text"`: 直接使用 `message` 字段
   - `"voice"`: 调用 `VoiceLogic.RecognizeVoice` 识别语音，获取文本
   - `"image"`: 调用 `UploadLogic.Upload` 上传图片，获取图片URL
3. 构建Agent模型输入消息（文本或 multimodal）
4. 调用 `ConversationNode.StreamConversation` 获取流式响应
5. 通过SSE返回流式响应

**技术细节**:
- 在 `StreamLogic.StreamConversation` 中实现统一处理逻辑
- 根据 `messageType` 调用相应的处理函数
- 需要传递会话ID、用户年龄、识别结果上下文等信息

### 3. Eino Graph多模态输入支持

**研究发现**:
- Eino框架支持多模态输入，在 `image_recognition.go` 中可以看到使用了 `UserInputMultiContent` 和 `MessageInputPart`
- `schema.Message` 支持 `UserInputMultiContent` 字段，可以包含多个 `MessageInputPart`
- `MessageInputPart` 支持 `ChatMessagePartTypeImageURL` 和 `ChatMessagePartTypeText` 类型

**实现方案**:
- **方案A（推荐）**: 扩展 `ConversationNode.StreamConversation` 方法，支持多模态输入
  - 优点：充分利用Eino框架的多模态能力，支持图文混排对话
  - 缺点：需要修改 `ConversationNode` 的实现
  - 实施：修改 `StreamConversation` 方法签名，支持图片URL参数，构建多模态消息

- **方案B**: 图片上传后，先调用图片识别接口，将识别结果作为文本输入
  - 优点：不需要修改 `ConversationNode` 的实现
  - 缺点：需要两次AI调用，响应时间较长
  - 实施：图片上传后，先调用 `ImageRecognitionNode`，然后将识别结果作为文本输入

**决策**: 采用方案A，扩展 `ConversationNode.StreamConversation` 方法，支持多模态输入。

**技术细节**:
- 修改 `ConversationNode.StreamConversation` 方法签名，添加 `imageURL` 参数
- 如果提供了图片URL，构建多模态消息：
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
  ```
- 如果没有图片URL，使用原有的文本消息格式

### 4. 错误处理机制

**当前状态**:
- `streamlogic.go` 中的 `StreamConversation` 方法已经有错误处理
- Agent模型调用失败时，会记录错误日志，但会降级到Mock数据

**实现方案**:
- **禁止降级到Mock数据**: 根据规范要求，禁止使用Mock数据
- **错误处理流程**:
  1. Agent模型调用失败时，记录详细错误日志（包含错误类型、错误信息、会话ID等）
  2. 向用户发送错误事件（通过SSE）
  3. 返回错误响应，不允许降级到Mock数据

**技术细节**:
- 修改 `streamlogic.go` 中的错误处理逻辑，移除Mock降级
- 记录详细的错误日志，包括：
  - 错误类型（超时、网络错误、模型错误等）
  - 错误信息
  - 会话ID
  - 用户年龄
  - 消息内容（脱敏处理）
- 通过SSE发送错误事件：
  ```go
  errorEvent := types.StreamEvent{
      Type:      "error",
      Content:   map[string]interface{}{
          "message": "Agent模型调用失败",
          "error": err.Error(),
      },
      SessionId: sessionId,
  }
  ```

### 5. 前端集成方案

**当前状态**:
- 前端已经有流式对话的SSE连接实现
- 语音输入和图片上传后，需要手动调用流式对话接口

**实现方案**:
- **方案A（推荐）**: 修改后端接口，语音识别和图片上传后直接返回SSE流式响应
  - 优点：前端不需要修改，只需要处理SSE响应
  - 缺点：需要修改后端API接口
  - 实施：修改 `VoiceHandler` 和 `UploadHandler`，返回SSE流式响应

- **方案B**: 保持现有接口不变，前端在收到识别/上传结果后自动调用流式对话接口
  - 优点：不需要修改后端API接口
  - 缺点：需要修改前端代码，需要处理两次请求
  - 实施：前端在收到语音识别或图片上传结果后，自动调用流式对话接口

**决策**: 采用方案A，修改后端接口，直接返回SSE流式响应。

**技术细节**:
- 修改 `VoiceHandler`，设置SSE响应头，调用 `StreamLogic.StreamConversation`
- 修改 `UploadHandler`，设置SSE响应头，调用 `StreamLogic.StreamConversation`
- 前端不需要修改，只需要处理SSE响应

## 关键技术决策

### 决策1: 语音输入后直接返回流式响应
- **决策**: 修改 `VoiceHandler`，语音识别后直接调用流式对话接口，返回SSE流式响应
- **理由**: 用户体验更好，语音识别后立即开始流式返回，减少请求次数
- **实施**: 修改 `VoiceHandler` 和 `VoiceLogic`，识别后调用 `StreamLogic.StreamConversation`

### 决策2: 图片上传后直接返回流式响应
- **决策**: 修改 `UploadHandler`，图片上传后直接调用流式对话接口，返回SSE流式响应
- **理由**: 用户体验更好，图片上传后立即开始流式返回，减少请求次数
- **实施**: 修改 `UploadHandler` 和 `UploadLogic`，上传后调用 `StreamLogic.StreamConversation`

### 决策3: 统一接口内部处理流程
- **决策**: 统一接口根据 `messageType` 自动调用语音识别或图片上传，然后调用Agent模型
- **理由**: 简化前端调用，后端统一处理，提升用户体验
- **实施**: 在 `StreamLogic.StreamConversation` 中实现统一处理逻辑，根据 `messageType` 调用相应处理函数

### 决策4: 扩展ConversationNode支持多模态输入
- **决策**: 扩展 `ConversationNode.StreamConversation` 方法，支持图片URL参数，构建多模态消息
- **理由**: 充分利用Eino框架的多模态能力，支持图文混排对话，提升用户体验
- **实施**: 修改 `ConversationNode.StreamConversation` 方法签名，支持图片URL参数

### 决策5: 禁止降级到Mock数据
- **决策**: Agent模型调用失败时，记录详细错误日志，向用户发送错误事件，不允许降级到Mock数据
- **理由**: 符合规范要求，确保用户获得真实的AI响应
- **实施**: 修改 `streamlogic.go` 中的错误处理逻辑，移除Mock降级

## 实施要点

1. **后端修改**:
   - 扩展 `StreamConversationHandler`，实现统一接口处理逻辑
   - 定义 `UnifiedStreamConversationRequest` 类型，包含 `messageType` 字段（必填）
   - 在 `StreamLogic.StreamConversation` 中根据 `messageType` 处理不同输入类型
   - 集成语音识别：当 `messageType` 为 `"voice"` 时，调用 `VoiceLogic.RecognizeVoice`
   - 集成图片上传：当 `messageType` 为 `"image"` 时，调用 `UploadLogic.Upload`
   - 扩展 `ConversationNode.StreamConversation` 方法，支持多模态输入
   - 修改错误处理逻辑，禁止降级到Mock数据

2. **类型定义扩展**:
   - 定义 `UnifiedStreamConversationRequest`，包含 `messageType` 字段（必填）
   - 根据 `messageType` 的值，请求包含对应的字段：`message`、`audio` 或 `image`
   - 扩展 `StreamEvent`，支持多模态事件类型

3. **前端修改**:
   - 更新对话服务，统一调用 `/api/conversation/stream` 接口
   - 文本输入：设置 `messageType: "text"`，发送 `message` 字段
   - 语音输入：设置 `messageType: "voice"`，发送 `audio` 字段
   - 图片输入：设置 `messageType: "image"`，发送 `image` 字段

4. **测试要点**:
   - 测试统一接口的文本输入流式返回
   - 测试统一接口的语音输入流式返回（包含语音识别）
   - 测试统一接口的图片输入流式返回（包含图片上传）
   - 测试 `messageType` 字段验证和错误处理
   - 测试多模态输入的流式返回
   - 测试错误处理机制

## 风险与缓解

1. **风险**: Eino框架的多模态支持可能不完整
   - **缓解**: 先实现文本+图片的多模态输入，如果Eino不支持，降级到方案B（先识别再对话）

2. **风险**: 修改API接口可能影响现有前端代码
   - **缓解**: 保持向后兼容，新增SSE流式响应接口，保留原有JSON响应接口

3. **风险**: 错误处理可能导致用户体验下降
   - **缓解**: 提供友好的错误提示，记录详细错误日志，便于问题排查

## 总结

通过实现统一流式接口，通过 `messageType` 字段明确指定输入类型，集成语音识别和图片上传逻辑，扩展 `ConversationNode` 支持多模态输入，禁止降级到Mock数据，可以实现文本、语音、图片三种输入方式的Agent模型流式返回功能。统一接口简化了前端调用逻辑，统一了后端处理流程，提升了代码可维护性。

