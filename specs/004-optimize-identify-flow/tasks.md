# Tasks: 优化识别流程与对话体验

**Input**: 设计文档来自 `/specs/004-optimize-identify-flow/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Organization**: 任务按用户故事组织，支持独立实现和测试。明确区分前端和后端任务。

## Format: `[ID] [P?] [Story] Description`

- **[P]**: 可并行执行（不同文件，无依赖）
- **[Story]**: 任务所属的用户故事（US1, US2, US3）
- 包含精确的文件路径
- 明确标注前端（frontend）或后端（backend）

## Path Conventions

- **前端**: `frontend/src/`
- **后端**: `backend/internal/`

---

## Phase 1: Setup (共享基础设施)

**Purpose**: 项目初始化和基础结构

- [x] T001 [P] 更新前端类型定义，添加识别结果上下文类型到 `frontend/src/types/api.ts`
- [x] T002 [P] 更新后端类型定义，添加识别结果上下文类型到 `backend/internal/types/types.go`
- [x] T003 [P] 更新对话消息类型，添加识别结果上下文字段到 `frontend/src/types/conversation.ts`

---

## Phase 2: Foundational (阻塞性前置条件)

**Purpose**: 核心基础设施，必须在所有用户故事实现前完成

**⚠️ CRITICAL**: 此阶段完成后才能开始用户故事实现

- [x] T004 更新API合约，确保前后端数据格式一致，参考 `specs/004-optimize-identify-flow/contracts/conversation.api`
- [x] T005 [P] 后端：更新会话状态管理，支持识别结果上下文关联到 `backend/internal/logic/conversationlogic.go`
- [x] T006 [P] 前端：更新API服务类型定义，确保与后端一致到 `frontend/src/services/api.ts`

**Checkpoint**: 基础设施就绪 - 用户故事实现可以开始并行进行

---

## Phase 3: User Story 1 - 识别后直接跳转问答页（优先级: P1）🎯 MVP

**Goal**: 用户拍照识别后，系统直接跳转到问答页面，首先展示识别结果，然后通过对话触发生成知识卡片

**Independent Test**: 拍照识别后验证系统是否直接跳转到问答页面，是否首先展示识别结果，是否支持通过对话触发卡片生成

### 前端实现 - User Story 1

- [x] T007 [US1] 前端：修改拍照页面，识别后直接跳转问答页，移除生成卡片调用到 `frontend/src/pages/Capture.tsx`
  - 移除 `generateCards` API调用（第64-70行）
  - 修改 `navigate` 跳转，只传递识别结果，不传递卡片数据（第82-90行）
  - 确保跳转路径为 `/result`

- [x] T008 [US1] 前端：修改问答页面，接收识别结果并展示为初始消息到 `frontend/src/pages/Result.tsx`
  - 修改 `useEffect`，接收识别结果而非卡片数据（第42-84行）
  - 创建初始系统消息，展示识别结果（对象名称、类别、置信度）
  - 移除卡片相关的初始消息创建逻辑

- [x] T009 [US1] 前端：更新LocationState接口，移除cards字段，添加识别结果字段到 `frontend/src/pages/Result.tsx`
  - 更新 `LocationState` 接口定义（第22-28行）
  - 确保类型安全

- [x] T010 [US1] 前端：更新API服务，移除识别后自动生成卡片的逻辑到 `frontend/src/services/api.ts`
  - 确保 `identifyImage` 函数只返回识别结果
  - 移除识别后自动调用 `generateCards` 的逻辑

### 后端实现 - User Story 1

- [x] T011 [US1] 后端：更新对话处理器，支持接收识别结果上下文到 `backend/internal/handler/conversationhandler.go`
  - 更新请求结构，支持 `IdentificationContext` 字段
  - 将识别结果上下文传递给逻辑层

- [x] T012 [US1] 后端：更新对话逻辑，维护识别结果上下文关联到 `backend/internal/logic/conversationlogic.go`
  - 在创建会话时，保存识别结果上下文
  - 在对话处理时，将识别结果上下文传递给AI模型
  - 确保识别结果与对话会话正确关联

- [x] T013 [US1] 后端：更新类型定义，添加识别结果上下文结构到 `backend/internal/types/types.go`
  - 添加 `IdentificationContext` 结构体
  - 更新 `ConversationRequest` 结构体，添加可选字段

**Checkpoint**: 此时，User Story 1 应该完全功能正常并可独立测试

---

## Phase 4: User Story 2 - 知识卡片简化显示（优先级: P2）

**Goal**: 知识卡片在对话中展示时，只显示文本内容，不显示图片

**Independent Test**: 在对话中生成知识卡片后，验证卡片是否只显示文本内容，不显示图片

### 前端实现 - User Story 2

- [x] T014 [P] [US2] 前端：修改科学卡片组件，移除图片显示逻辑到 `frontend/src/components/cards/ScienceCard.tsx`
  - 移除图片渲染代码
  - 保留文本内容显示
  - 优化布局，确保纯文本展示清晰

- [x] T015 [P] [US2] 前端：修改诗词卡片组件，移除图片显示逻辑到 `frontend/src/components/cards/PoetryCard.tsx`
  - 移除图片渲染代码
  - 保留文本内容显示
  - 优化布局，确保纯文本展示清晰

- [x] T016 [P] [US2] 前端：修改英语卡片组件，移除图片显示逻辑到 `frontend/src/components/cards/EnglishCard.tsx`
  - 移除图片渲染代码
  - 保留文本内容显示
  - 优化布局，确保纯文本展示清晰

- [x] T017 [US2] 前端：更新对话消息组件，确保卡片消息不显示图片到 `frontend/src/components/conversation/ConversationMessage.tsx`
  - 检查卡片消息渲染逻辑
  - 确保调用卡片组件时不传递图片数据
  - 验证所有卡片类型都遵循不显示图片的规则

**Checkpoint**: 此时，User Story 2 应该完全功能正常并可独立测试

---

## Phase 5: User Story 3 - 追问时即时展示用户问题（优先级: P1）

**Goal**: 用户在对话中发送问题时，系统立即在对话界面中展示用户发送的问题，然后通过流式返回展示系统响应

**Independent Test**: 在对话中发送问题（文本、语音、图片），验证问题是否立即显示在对话列表中，系统响应是否通过流式返回实时展示

### 前端实现 - User Story 3

- [x] T018 [US3] 前端：修改对话服务，实现乐观更新，立即显示用户消息到 `frontend/src/services/conversation.ts`
  - 在 `sendMessage` 函数中，发送API请求前先创建用户消息对象
  - 立即将用户消息添加到消息列表（乐观更新）
  - 使用本地时间戳和临时ID
  - API成功后更新消息ID，失败则回滚并显示错误

- [x] T019 [US3] 前端：更新消息发送处理，支持文本、语音、图片的乐观更新到 `frontend/src/pages/Result.tsx`
  - 修改 `handleSendMessage` 函数（第104-138行）
  - 在发送前立即添加用户消息到消息列表
  - 确保消息包含时间戳和发送者标识

- [x] T020 [US3] 前端：更新语音输入处理，实现乐观更新到 `frontend/src/pages/Result.tsx`
  - 修改 `handleVoiceResult` 函数（第140-142行）
  - 语音识别后立即显示识别结果作为用户消息

- [x] T021 [US3] 前端：更新图片输入处理，实现乐观更新到 `frontend/src/pages/Result.tsx`
  - 修改 `handleImageSelect` 函数（第144-179行）
  - 图片上传后立即显示图片作为用户消息

- [x] T022 [US3] 前端：优化对话列表组件，支持立即显示用户消息到 `frontend/src/components/conversation/ConversationList.tsx`
  - 确保消息列表能够实时更新
  - 优化渲染性能，支持快速添加消息

- [x] T023 [US3] 前端：优化流式返回处理，确保系统响应实时展示到 `frontend/src/services/sse.ts`
  - 确保SSE连接正确建立
  - 实时接收并更新系统消息
  - 提供加载状态和进度反馈

### 后端实现 - User Story 3

- [x] T024 [US3] 后端：优化对话逻辑，支持识别结果上下文传递到 `backend/internal/logic/conversationlogic.go`
  - 在处理消息时，如果存在识别结果上下文，将其传递给AI模型
  - 确保多轮对话能够正确关联识别结果

- [x] T025 [US3] 后端：优化并发处理，支持200+并发会话到 `backend/internal/logic/conversationlogic.go`
  - 使用goroutine处理并发请求
  - 使用sync.RWMutex保护会话状态并发读写
  - 实现会话过期清理机制（已在storage中实现）

- [ ] T026 [US3] 后端：优化流式返回，确保首字符延迟≤1秒到 `backend/internal/handler/conversationhandler.go`
  - 优化SSE实现，立即发送第一个chunk
  - 确保流式返回性能满足要求
  - 实现连接重试机制
  - 注意：流式返回优化需要后续实现，当前基础功能已完成

- [ ] T027 [US3] 后端：更新Agent图，支持对话触发卡片生成到 `backend/internal/agent/graph.go`
  - 修改意图识别逻辑，支持"生成卡片"意图
  - 在对话中触发卡片生成时，通过流式返回发送卡片
  - 确保卡片生成流程正确
  - 注意：Agent图更新需要后续实现，当前基础功能已完成

**Checkpoint**: 此时，User Story 3 应该完全功能正常并可独立测试

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: 影响多个用户故事的改进

- [x] T028 [P] 前端：添加错误处理，确保识别失败时友好提示到 `frontend/src/pages/Capture.tsx`
- [x] T029 [P] 前端：添加错误处理，确保对话失败时友好提示到 `frontend/src/pages/Result.tsx`
- [x] T030 [P] 后端：添加错误处理和日志记录到 `backend/internal/logic/conversationlogic.go`
- [x] T031 [P] 前端：优化性能，确保消息展示延迟≤0.5秒到 `frontend/src/components/conversation/ConversationList.tsx`（通过乐观更新实现）
- [x] T032 [P] 后端：实现会话清理机制，移除过期会话到 `backend/internal/logic/conversationlogic.go`（已在storage中实现）
- [x] T033 [P] 前端：验证前后端数据格式一致性，确保所有消息正确同步（类型定义已统一）
- [x] T034 [P] 后端：验证前后端数据格式一致性，确保API响应格式正确（类型定义已统一）
- [ ] T035 运行快速开始验证步骤，参考 `specs/004-optimize-identify-flow/quickstart.md`（需要手动验证）

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: 无依赖 - 可立即开始
- **Foundational (Phase 2)**: 依赖Setup完成 - 阻塞所有用户故事
- **User Stories (Phase 3+)**: 所有依赖Foundational阶段完成
  - User Story 1 和 User Story 3 可以并行（都是P1优先级）
  - User Story 2 可以独立实现（P2优先级）
- **Polish (Final Phase)**: 依赖所有期望的用户故事完成

### User Story Dependencies

- **User Story 1 (P1)**: 可在Foundational完成后开始 - 无依赖其他故事
- **User Story 2 (P2)**: 可在Foundational完成后开始 - 无依赖其他故事
- **User Story 3 (P1)**: 可在Foundational完成后开始 - 无依赖其他故事

### Within Each User Story

- 前端和后端任务可以并行进行
- 类型定义应在实现前完成
- 核心逻辑应在UI组件前完成
- 故事完成后才能移动到下一个优先级

### Parallel Opportunities

- 所有Setup任务标记[P]的可并行运行
- 所有Foundational任务标记[P]的可并行运行（在Phase 2内）
- Foundational完成后，所有用户故事可以并行开始（如果团队容量允许）
- User Story 2中的三个卡片组件任务可以并行（T014, T015, T016）
- 前端和后端任务可以并行进行

---

## Implementation Strategy

### MVP First (User Story 1 + User Story 3)

1. 完成Phase 1: Setup
2. 完成Phase 2: Foundational（关键 - 阻塞所有故事）
3. 完成Phase 3: User Story 1（识别后跳转问答页）
4. 完成Phase 5: User Story 3（即时展示用户问题）
5. **停止并验证**: 测试User Story 1和User Story 3独立工作
6. 如果准备就绪，部署/演示

### Incremental Delivery

1. 完成Setup + Foundational → 基础设施就绪
2. 添加User Story 1 → 独立测试 → 部署/演示（MVP核心功能）
3. 添加User Story 3 → 独立测试 → 部署/演示（完整对话体验）
4. 添加User Story 2 → 独立测试 → 部署/演示（界面优化）
5. 每个故事增加价值而不破坏之前的故事

### Parallel Team Strategy

多开发者协作：

1. 团队一起完成Setup + Foundational
2. Foundational完成后：
   - 开发者A: User Story 1（前端）
   - 开发者B: User Story 1（后端）
   - 开发者C: User Story 3（前端）
   - 开发者D: User Story 3（后端）
3. 故事完成后独立集成

---

## Notes

- [P] 任务 = 不同文件，无依赖
- [Story] 标签将任务映射到特定用户故事以便追溯
- 每个用户故事应该独立完成和测试
- 前端和后端任务明确区分，可以并行进行
- 提交前验证每个任务或逻辑组
- 在任何检查点停止以独立验证故事
- 避免：模糊任务、同一文件冲突、跨故事依赖破坏独立性

---

## Task Summary

**总任务数**: 35

**按用户故事分布**:
- User Story 1: 7个任务（前端4个，后端3个）
- User Story 2: 4个任务（全部前端）
- User Story 3: 10个任务（前端6个，后端4个）
- Setup: 3个任务
- Foundational: 3个任务
- Polish: 8个任务

**并行机会**:
- Setup阶段: 3个并行任务
- Foundational阶段: 2个并行任务
- User Story 2: 3个卡片组件可并行
- 前端和后端任务可并行进行

**建议MVP范围**: User Story 1 + User Story 3（核心功能）
