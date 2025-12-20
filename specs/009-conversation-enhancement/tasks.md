# Tasks: 对话页面完善 - Agent模型流式返回

**Input**: Design documents from `/specs/009-conversation-enhancement/`
**Prerequisites**: plan.md ✅, spec.md ✅, research.md ✅, data-model.md ✅, contracts/ ✅

**核心目标**: 实现统一流式接口，支持文本输入、语音输入、图片输入三种输入方式，通过 `messageType` 字段明确指定输入类型。通过 `.env` 配置文件中的 `USE_AI_MODEL` 字段控制是否使用AI模型调用，默认值为 `true`（使用AI模型）。

**MVP策略**: 快速迭代，优先实现统一接口的核心功能，支持三种输入类型的流式返回

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US2, US3)
- Include exact file paths in descriptions

## Phase 1: 后端核心实现

**Purpose**: 实现统一流式接口，支持文本、语音、图片三种输入方式

### 1.1 添加配置字段支持

- [X] T001 [P] [US1,US2,US3] 在 `backend/internal/config/config.go` 的 `AIConfig` 结构体中添加 `UseAIModel` 字段（布尔类型），从环境变量 `USE_AI_MODEL` 读取配置，默认值为 `true`（使用AI模型）

### 1.2 定义统一流式请求类型

- [X] T002 [P] [US1,US2,US3] 定义 `UnifiedStreamConversationRequest` 类型，在 `backend/internal/types/stream_types.go` 中添加统一流式请求类型，包含 `messageType` 字段（必填）和对应的条件字段（`message`、`audio`、`image`）

### 1.3 扩展 ConversationNode 支持多模态输入

- [X] T003 [US3] 扩展 `StreamConversation` 方法签名，在 `backend/internal/agent/nodes/conversation_node.go` 中添加 `imageURL` 参数，支持多模态消息构建

### 1.4 扩展 StreamLogic 实现统一接口处理逻辑

- [X] T004 [US1,US2,US3] 扩展 `StreamConversation` 方法，在 `backend/internal/logic/streamlogic.go` 中接收 `UnifiedStreamConversationRequest`，根据 `messageType` 字段处理不同输入类型：
  - `messageType: "text"` → 直接使用 `message` 字段
  - `messageType: "voice"` → 调用 `VoiceLogic.RecognizeVoice` 识别语音，获取文本
  - `messageType: "image"` → 调用 `UploadLogic.Upload` 上传图片，获取图片URL

### 1.5 实现配置字段控制逻辑

- [X] T005 [US1,US2,US3] 在 `StreamLogic.StreamConversationUnified` 方法中，根据 `UseAIModel` 配置决定是否使用AI模型：
  - 当 `UseAIModel=true`（默认值）时，必须使用AI模型，禁止降级到Mock数据
  - 当 `UseAIModel=false` 时，可以使用Mock数据作为降级方案（仅用于开发测试场景）
  - 记录详细错误日志，无论配置如何

### 1.6 修改错误处理逻辑

- [X] T006 [US1,US2,US3] 修改 `StreamConversationUnified` 方法，在 `backend/internal/logic/streamlogic.go` 中根据 `UseAIModel` 配置实现错误处理逻辑：
  - 当 `UseAIModel=true` 时，Agent模型调用失败时记录详细错误日志，向用户发送错误事件，不允许降级到Mock数据
  - 当 `UseAIModel=false` 时，可以使用Mock数据作为降级方案

### 1.7 更新 StreamConversationHandler 支持统一接口

- [X] T007 [US1,US2,US3] 更新 `StreamConversationHandler`，在 `backend/internal/handler/streamhandler.go` 中接收 `UnifiedStreamConversationRequest`，调用统一流式逻辑

### 1.8 添加 messageType 字段验证

- [X] T008 [US1,US2,US3] 在 `StreamLogic.StreamConversationUnified` 中添加 `messageType` 字段验证，确保字段值与对应字段匹配（如 `messageType: "voice"` 时必须包含 `audio` 字段）

### 1.9 更新API定义

- [X] T009 [US1,US2,US3] 更新API定义，在 `backend/api/explore.api` 中更新 `/api/conversation/stream` 接口，使用 `UnifiedStreamConversationRequest` 类型

## Phase 2: 前端集成

**Purpose**: 前端统一调用统一流式接口

- [X] T010 [US1] 更新文本输入处理逻辑，在 `frontend/src/services/conversation.ts` 中调用统一接口 `/api/conversation/stream`，设置 `messageType: "text"`，发送 `message` 字段
- [X] T011 [US2] 更新语音输入处理逻辑，在 `frontend/src/services/conversation.ts` 中调用统一接口 `/api/conversation/stream`，设置 `messageType: "voice"`，发送 `audio` 字段
- [X] T012 [US3] 更新图片输入处理逻辑，在 `frontend/src/services/conversation.ts` 中调用统一接口 `/api/conversation/stream`，设置 `messageType: "image"`，发送 `image` 字段
- [X] T013 [P] [US1,US2,US3] 统一错误处理机制，在 `frontend/src/services/conversation.ts` 中统一处理统一接口的错误响应

## Phase 3: 测试验证

- [ ] T014 [US1,US2,US3] 测试配置字段读取：验证 `USE_AI_MODEL` 配置字段正确读取和应用，默认值为 `true`（使用AI模型）
- [ ] T015 [US1] 测试统一接口的文本输入流式返回（`USE_AI_MODEL=true`）：发送文本消息 → 流式返回
- [ ] T016 [US2] 测试统一接口的语音输入流式返回（`USE_AI_MODEL=true`）：发送语音数据 → 语音识别 → 流式返回
- [ ] T017 [US3] 测试统一接口的图片输入流式返回（`USE_AI_MODEL=true`）：发送图片数据 → 图片上传 → 流式返回
- [ ] T018 [US1,US2,US3] 测试 `messageType` 字段验证：测试缺少字段、字段不匹配等错误情况
- [ ] T019 [US1,US2,US3] 测试错误处理（`USE_AI_MODEL=true`）：模拟Agent模型调用失败，验证错误日志和错误事件，验证禁止降级到Mock数据
- [ ] T020 [US1,US2,US3] 测试错误处理（`USE_AI_MODEL=false`）：模拟Agent模型调用失败，验证可以使用Mock数据作为降级方案
- [ ] T021 [US3] 测试多模态输入：图片+文本的多模态消息构建和流式返回

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1**: 后端核心实现，必须按顺序完成
- **Phase 2**: 前端集成，依赖Phase 1完成，可以并行开发（不同输入类型）
- **Phase 3**: 测试验证，依赖Phase 1和Phase 2完成

### Task Dependencies

- T001 → T005, T006: 配置字段定义后才能使用
- T002 → T004, T007: 类型定义后才能使用
- T003 → T004: ConversationNode扩展后才能使用多模态输入
- T004 → T007: StreamLogic扩展后才能被Handler调用
- T005 → T004, T006: 配置字段控制逻辑实现后才能正确处理
- T006 → T004: 错误处理修改后才能正确调用
- T008 → T004: 字段验证后才能正确处理请求
- T009 → T007: API定义后才能更新Handler

### Parallel Opportunities

- T001可以与其他任务并行（配置字段定义）
- T002可以与其他任务并行（类型定义）
- T003可以与其他任务并行（ConversationNode扩展）
- T010、T011、T012可以并行（不同输入类型的处理）

## Implementation Strategy

### MVP First

1. 完成Phase 1: 后端核心实现（统一接口 + 配置字段）
2. **STOP and VALIDATE**: 测试配置字段读取和应用，测试统一接口的文本、语音、图片输入流式返回功能（`USE_AI_MODEL=true`）
3. 快速调通：三种输入方式都能通过统一接口看到Agent模型流式返回，配置字段正确控制是否使用AI模型

### Incremental Delivery

1. **MVP**: Phase 1 → 统一接口支持三种输入类型的流式返回可用，配置字段正确控制是否使用AI模型（默认使用AI模型）
2. **增强**: Phase 2 → 前端统一调用统一接口
3. **完善**: Phase 3 → 测试验证和错误处理，包括配置字段测试和Mock数据降级测试

## Task Summary

- **总任务数**: 21个任务
- **Phase 1 (后端核心)**: 9个任务
- **Phase 2 (前端集成)**: 4个任务
- **Phase 3 (测试验证)**: 8个任务

**MVP范围**: Phase 1（9个任务）

**并行机会**: 
- Phase 1: T001、T002、T003可以并行
- Phase 2: T010、T011、T012可以并行

