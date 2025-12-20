# Tasks: H5对话落地页 - 追问能力实现

**Input**: Design documents from `/specs/007-conversation-landing-page/`
**Prerequisites**: plan.md ✅, spec.md ✅, research.md ✅, data-model.md ✅, contracts/ ✅

**核心目标**: 实现对话页的追问能力，支持流式输出、打字机效果、图片loading占位

**MVP策略**: 快速迭代，功能快速调通，优先实现追问对话核心功能

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: 项目初始化和基础结构

- [ ] T001 检查并确认Eino框架依赖已安装，验证backend/go.mod中的eino和eino-ext版本
- [ ] T002 [P] 检查前端依赖，确认React、Tailwind CSS、Axios已安装在frontend/package.json
- [ ] T003 [P] 验证现有API路由配置，检查backend/internal/handler/routes.go中的对话相关路由

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: 核心基础设施，必须在所有用户故事之前完成

**⚠️ CRITICAL**: 这些任务完成后才能开始用户故事实现

- [ ] T004 扩展types定义，在backend/internal/types/types.go中添加StreamConversationRequest和StreamEvent类型
- [ ] T005 [P] 扩展存储接口，在backend/internal/storage/memory.go中添加20轮消息限制逻辑
- [ ] T006 [P] 创建对话节点基础结构，在backend/internal/agent/nodes/conversation_node.go中创建ConversationNode结构体

**Checkpoint**: Foundation ready - 用户故事实现可以开始

---

## Phase 3: User Story 2 - 追问对话和流式输出 (Priority: P2) 🎯 MVP核心

**Goal**: 用户在对话页面可以发送追问消息，系统通过流式接口返回回答，支持打字机效果和图片生成loading占位

**Independent Test**: 用户在对话页面输入问题并发送，系统通过流式接口返回回答，文本逐字显示（打字机效果），如果包含图片生成，显示loading占位符。可以独立工作，即使没有历史消息保存，也能完成单轮对话。

### 后端实现 - 流式对话核心

- [X] T007 [US2] 实现基于年级的prompt生成函数，在backend/internal/agent/nodes/conversation_node.go中添加generateSystemPrompt方法
- [X] T008 [US2] 实现上下文消息转换函数，在backend/internal/logic/streamlogic.go中添加convertToEinoMessages方法，将内部消息转换为Eino Message格式
- [X] T009 [US2] 实现Eino流式对话节点，在backend/internal/agent/nodes/conversation_node.go中实现StreamConversation方法，调用Eino ChatModel.Stream接口
- [X] T010 [US2] 扩展流式逻辑，在backend/internal/logic/streamlogic.go中实现StreamConversation方法，集成Eino流式输出和SSE发送（使用Recv()方法读取）
- [X] T011 [US2] 更新流式Handler，在backend/internal/handler/streamhandler.go中实现StreamConversationHandler，处理SSE连接和流式事件发送
- [X] T012 [US2] 添加流式对话路由，在backend/internal/handler/routes.go中注册GET /api/conversation/stream路由
- [X] T013 [US2] 扩展API定义，在backend/api/explore.api中添加流式对话接口定义

### 前端实现 - 流式对话和打字机效果

- [X] T014 [P] [US2] 创建流式对话Hook，在frontend/src/hooks/useStreamConversation.ts中实现useStreamConversation Hook
- [X] T015 [P] [US2] 创建打字机效果Hook，在frontend/src/hooks/useTypingEffect.ts中实现useTypingEffect Hook
- [X] T016 [P] [US2] 创建图片loading占位组件，在frontend/src/components/common/ImagePlaceholder.tsx中实现ImagePlaceholder组件
- [X] T017 [US2] 扩展SSE服务，在frontend/src/services/sse.ts中优化createSSEConnection函数，支持流式事件处理
- [X] T018 [US2] 扩展对话服务，在frontend/src/services/conversation.ts中添加streamConversation函数，封装流式对话调用
- [X] T019 [US2] 更新对话消息组件，在frontend/src/components/conversation/ConversationMessage.tsx中添加打字机效果支持
- [X] T020 [US2] 更新对话列表组件，在frontend/src/components/conversation/ConversationList.tsx中支持流式消息实时更新
- [X] T021 [US2] 更新对话页面，在frontend/src/pages/Result.tsx中集成流式对话功能，实现用户消息发送和AI流式回答显示

### 前后端集成测试

- [ ] T022 [US2] 测试流式对话端到端流程：用户发送消息 → 后端流式返回 → 前端打字机效果显示
- [ ] T023 [US2] 测试图片生成loading占位：图片生成进度显示 → 图片完成后替换占位符

**Checkpoint**: User Story 2应该完全功能正常，可以独立测试。用户可以发送追问消息，看到流式回答和打字机效果。

---

## Phase 4: User Story 1 - 首次拍照后知识卡片生成 (Priority: P1) - 前置支持

**Goal**: 用户完成拍照识别后，系统自动进行意图识别，并生成三张知识卡片，作为对话页面的入口

**Independent Test**: 用户拍照识别后，系统自动生成三张知识卡片并展示在对话页面。用户可以看到卡片内容，了解识别对象的相关知识。

**Note**: 此功能已部分实现，主要需要确保与追问功能的集成

- [ ] T024 [US1] 验证意图识别功能，确保backend/internal/logic/intentlogic.go中的RecognizeIntent方法正常工作
- [ ] T025 [US1] 验证卡片生成功能，确保backend/internal/logic/generatecardslogic.go中的GenerateCards方法正常工作
- [ ] T026 [US1] 确保对话页面自动生成卡片逻辑，在frontend/src/pages/Result.tsx中验证generateCardsAutomatically函数正常工作
- [ ] T027 [US1] 测试知识卡片生成到追问的流程：卡片生成 → 用户发送追问 → 流式回答

**Checkpoint**: User Story 1和User Story 2可以协同工作，用户可以看到卡片并继续追问

---

## Phase 5: User Story 3 - 历史消息保存和上下文关联 (Priority: P3) - 增强功能

**Goal**: 系统保存最近20轮对话历史，在生成回答时使用这些历史消息作为上下文，确保对话的连贯性

**Independent Test**: 用户进行多轮对话后，系统保存最近20轮消息。当用户继续提问时，系统使用这些历史消息作为上下文生成回答，回答内容体现对之前对话的理解。

### 后端实现 - 上下文管理

- [ ] T028 [US3] 实现20轮消息限制逻辑，在backend/internal/storage/memory.go中添加消息数量检查和自动删除最早消息的逻辑
- [ ] T029 [US3] 实现上下文消息获取函数，在backend/internal/logic/streamlogic.go中实现getContextMessages方法，限制为最近20轮
- [ ] T030 [US3] 集成上下文到流式对话，在backend/internal/logic/streamlogic.go的StreamConversation方法中使用getContextMessages获取上下文
- [ ] T031 [US3] 确保消息保存逻辑，在backend/internal/logic/conversationlogic.go中验证消息保存到存储的逻辑

### 前端实现 - 历史消息显示

- [ ] T032 [US3] 实现历史消息恢复，在frontend/src/pages/Result.tsx中添加从localStorage恢复历史消息的逻辑
- [ ] T033 [US3] 实现历史消息持久化，在frontend/src/services/storage.ts中添加对话历史保存到localStorage的逻辑
- [ ] T034 [US3] 测试多轮对话上下文关联：用户问"这是什么？" → AI回答 → 用户问"它有什么特点？" → AI理解"它"指代

**Checkpoint**: User Story 3完成，多轮对话可以正确使用上下文，回答具有连贯性

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: 影响多个用户故事的改进和优化

- [ ] T035 [P] 优化移动端响应式设计，在frontend/src/styles/responsive.ts中添加移动端优先的Tailwind配置
- [ ] T036 [P] 优化打字机效果性能，在frontend/src/hooks/useTypingEffect.ts中使用requestAnimationFrame优化渲染
- [ ] T037 [P] 添加错误处理和重连机制，在frontend/src/services/sse.ts中实现SSE连接错误处理和自动重连
- [ ] T038 [P] 添加加载状态反馈，在frontend/src/pages/Result.tsx中显示流式输出过程中的加载状态
- [ ] T039 [P] 优化图片loading占位动画，在frontend/src/components/common/ImagePlaceholder.tsx中优化进度显示和动画效果
- [ ] T040 添加日志记录，在backend/internal/logic/streamlogic.go中添加流式对话的关键日志
- [ ] T041 验证性能指标：流式回答启动时间<1秒，打字机效果流畅度60fps
- [ ] T042 运行quickstart.md中的端到端测试，验证完整流程

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: 无依赖 - 可以立即开始
- **Foundational (Phase 2)**: 依赖Setup完成 - **阻塞所有用户故事**
- **User Story 2 (Phase 3)**: 依赖Foundational完成 - **MVP核心功能**
- **User Story 1 (Phase 4)**: 依赖Foundational完成 - 前置支持功能
- **User Story 3 (Phase 5)**: 依赖User Story 2完成 - 增强功能
- **Polish (Phase 6)**: 依赖所有用户故事完成

### User Story Dependencies

- **User Story 2 (P2)**: MVP核心，可以独立实现和测试
- **User Story 1 (P1)**: 前置支持，确保与User Story 2集成
- **User Story 3 (P3)**: 增强功能，依赖User Story 2的基础实现

### Within Each User Story

- 后端实现 → 前端实现 → 集成测试
- 核心功能 → 增强功能 → 优化

### Parallel Opportunities

**Phase 2 (Foundational)**:
- T005和T006可以并行（不同文件）

**Phase 3 (User Story 2)**:
- T014, T015, T016可以并行（不同Hook和组件）
- T007, T008可以并行（不同方法）
- T017, T018可以并行（不同服务文件）

**Phase 5 (User Story 3)**:
- T032, T033可以并行（不同功能）

**Phase 6 (Polish)**:
- 所有标记[P]的任务可以并行

---

## Parallel Example: User Story 2

```bash
# 前端Hook和组件可以并行开发：
Task: "创建流式对话Hook，在frontend/src/hooks/useStreamConversation.ts中实现"
Task: "创建打字机效果Hook，在frontend/src/hooks/useTypingEffect.ts中实现"
Task: "创建图片loading占位组件，在frontend/src/components/common/ImagePlaceholder.tsx中实现"

# 后端方法可以并行开发：
Task: "实现基于年级的prompt生成函数，在backend/internal/agent/nodes/conversation_node.go中"
Task: "实现上下文消息转换函数，在backend/internal/logic/streamlogic.go中"
```

---

## Implementation Strategy

### MVP First (快速调通追问功能)

1. ✅ 完成Phase 1: Setup
2. ✅ 完成Phase 2: Foundational
3. ✅ **完成Phase 3: User Story 2 (追问对话和流式输出)** - **MVP核心**
4. **STOP and VALIDATE**: 测试追问功能独立工作
5. 快速调通：用户可以发送追问，看到流式回答和打字机效果

### Incremental Delivery

1. **MVP**: Setup + Foundational + User Story 2 → 追问功能可用
2. **增强**: User Story 1 → 确保卡片生成与追问集成
3. **优化**: User Story 3 → 添加上下文关联
4. **完善**: Polish → 性能优化和错误处理

### 前后端并行开发策略

**后端开发者**:
- Phase 2: T004-T006 (Foundational)
- Phase 3: T007-T013 (后端流式对话实现)
- Phase 5: T028-T031 (上下文管理)

**前端开发者**:
- Phase 3: T014-T021 (前端流式对话和打字机效果)
- Phase 5: T032-T033 (历史消息显示)
- Phase 6: T035-T039 (优化和错误处理)

**集成测试**:
- Phase 3: T022-T023 (端到端测试)
- Phase 4: T027 (流程测试)
- Phase 5: T034 (上下文测试)

---

## Notes

- **[P]标记**: 不同文件，无依赖，可以并行开发
- **[Story]标记**: 映射到特定用户故事，便于追踪
- **MVP优先**: 优先实现User Story 2（追问功能），快速调通
- **前后端分离**: 后端和前端任务可以并行开发
- **快速迭代**: 每个任务完成后立即测试，确保功能可用
- **避免**: 模糊任务、同一文件冲突、跨故事依赖破坏独立性

## Task Summary

- **总任务数**: 42个任务
- **Phase 1 (Setup)**: 3个任务
- **Phase 2 (Foundational)**: 3个任务
- **Phase 3 (User Story 2 - MVP核心)**: 12个任务（后端7个，前端5个）
- **Phase 4 (User Story 1 - 前置支持)**: 4个任务
- **Phase 5 (User Story 3 - 增强功能)**: 7个任务
- **Phase 6 (Polish)**: 8个任务

**MVP范围**: Phase 1 + Phase 2 + Phase 3（18个任务）

**并行机会**: 
- Phase 2: 2个并行任务
- Phase 3: 5个并行任务
- Phase 5: 2个并行任务
- Phase 6: 6个并行任务

