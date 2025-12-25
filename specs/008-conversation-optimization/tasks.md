# Tasks: 对话体验与性能优化

**Input**: Design documents from `/specs/008-conversation-optimization/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story. Tasks are separated by frontend and backend.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3, US4)
- Include exact file paths in descriptions

## Path Conventions

- **Web app**: `backend/`, `frontend/`
- Frontend: `frontend/src/`
- Backend: `backend/internal/`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and dependency setup

### Frontend Setup

- [x] T001 [P] Install react-markdown dependency in frontend/package.json
- [x] T002 [P] Update frontend TypeScript types for streaming message fields in frontend/src/types/conversation.ts

### Backend Setup

- [x] T003 [P] Extend ConversationMessage type with streamingText and markdown fields in backend/internal/types/types.go
- [x] T004 [P] Extend StreamEvent type with markdown field in backend/internal/types/types.go

**Checkpoint**: Dependencies and type definitions ready

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

### Frontend Foundational

- [x] T005 [P] Create useTextToSpeech hook in frontend/src/hooks/useTextToSpeech.ts for text-to-speech functionality
- [x] T006 [P] Update ConversationMessage interface with streamingText and markdown fields in frontend/src/types/conversation.ts

### Backend Foundational

- [x] T007 [P] Update ConversationMessage struct with StreamingText and Markdown fields in backend/internal/types/types.go
- [x] T008 [P] Update StreamEvent struct with Markdown field in backend/internal/types/types.go

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - 流式消息实时渲染（Priority: P1）🎯 MVP

**Goal**: 修复流式消息实时渲染问题，确保每个文本片段到达后立即显示，实现真正的打字机效果

**Independent Test**: 在对话页面发送消息后，观察系统响应是否在接收到流式数据后立即更新UI，而不是等待流结束。可以通过网络节流工具验证在不同网络条件下的实时渲染效果。

### Frontend Implementation for User Story 1

- [x] T009 [US1] Optimize handleSendMessage function to use flushSync for immediate UI updates in frontend/src/pages/Result.tsx
- [x] T010 [US1] Update useStreamConversation hook to immediately update state on each message event in frontend/src/hooks/useStreamConversation.ts
- [x] T011 [US1] Replace accumulatedText pattern with immediate state updates using useState updater function in frontend/src/pages/Result.tsx
- [x] T012 [US1] Add useRef to track streaming text and avoid closure issues in frontend/src/pages/Result.tsx
- [x] T013 [US1] Update ConversationMessage component to display streamingText when isStreaming is true in frontend/src/components/conversation/ConversationMessage.tsx
- [x] T014 [US1] Ensure network interruption handling preserves already received content in frontend/src/pages/Result.tsx

### Backend Implementation for User Story 1

- [x] T015 [US1] Ensure streamlogic.go sends each text fragment immediately without buffering in backend/internal/logic/streamlogic.go
- [x] T016 [US1] Add Flush() call after each message event to ensure immediate transmission in backend/internal/logic/streamlogic.go
- [x] T017 [US1] Verify SSE event format includes index field for character position in backend/internal/logic/streamlogic.go

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently. Stream messages should render in real-time with typing effect.

---

## Phase 4: User Story 2 - 流式消息Markdown格式支持（Priority: P1）

**Goal**: 实现流式消息的Markdown格式渲染，支持标题、列表、代码块、链接等常用Markdown语法

**Independent Test**: 在对话中发送会触发Markdown格式响应的问题（如要求代码示例、列表等），验证响应内容是否正确渲染为格式化的Markdown，而不是纯文本。

### Frontend Implementation for User Story 2

- [x] T018 [P] [US2] Install react-markdown and optional plugins (remark-gfm) in frontend/package.json
- [x] T019 [US2] Integrate react-markdown component in ConversationMessage renderContent function in frontend/src/components/conversation/ConversationMessage.tsx
- [x] T020 [US2] Add markdown detection logic to determine if content should be rendered as Markdown in frontend/src/components/conversation/ConversationMessage.tsx
- [x] T021 [US2] Configure react-markdown with appropriate plugins for code highlighting and GitHub Flavored Markdown in frontend/src/components/conversation/ConversationMessage.tsx
- [x] T022 [US2] Ensure Markdown rendering updates in real-time with streaming text in frontend/src/components/conversation/ConversationMessage.tsx
- [x] T023 [US2] Optimize Markdown rendering performance using React.memo to avoid unnecessary re-renders in frontend/src/components/conversation/ConversationMessage.tsx
- [x] T024 [US2] Add styling for Markdown elements (code blocks, lists, links) in frontend/src/components/conversation/ConversationMessage.tsx

### Backend Implementation for User Story 2

- [x] T025 [US2] Add markdown field detection logic to identify Markdown content in streamlogic.go in backend/internal/logic/streamlogic.go
- [x] T026 [US2] Include markdown field in StreamEvent when content contains Markdown in backend/internal/logic/streamlogic.go
- [x] T027 [US2] Update ConversationMessage response to include markdown field when applicable in backend/internal/logic/conversationlogic.go

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently. Stream messages should render in real-time with Markdown formatting.

---

## Phase 5: User Story 3 - 知识卡片生成性能优化（Priority: P1）

**Goal**: 将知识卡片生成接口响应时间从40秒优化到5秒内，实现流式返回支持

**Independent Test**: 请求生成知识卡片，测量响应时间是否在5秒内完成。可以通过多次请求验证性能优化的稳定性。

### Backend Implementation for User Story 3

- [x] T028 [US3] Optimize ExecuteCardGeneration to ensure true parallel execution of three cards in backend/internal/agent/graph.go
- [x] T029 [US3] Add timeout control for each card generation (10 seconds per card) in backend/internal/agent/graph.go
- [x] T030 [US3] Implement streaming return for generate-cards endpoint to return cards as they complete in backend/internal/logic/generatecardslogic.go
- [x] T031 [US3] Add progress feedback mechanism for card generation in backend/internal/logic/generatecardslogic.go
- [x] T032 [US3] Optimize AI model prompt length and structure to reduce generation time in backend/internal/agent/nodes/text_generation.go
- [x] T033 [US3] Add error handling and fallback for timeout scenarios in backend/internal/logic/generatecardslogic.go
- [x] T034 [US3] Update generate-cards handler to support streaming mode via query parameter in backend/internal/handler/generatecardshandler.go
- [x] T035 [US3] Implement SSE streaming for card generation endpoint in backend/internal/handler/generatecardshandler.go

### Frontend Implementation for User Story 3

- [x] T036 [US3] Update generateCards API call to support streaming mode in frontend/src/services/api.ts
- [x] T037 [US3] Implement SSE connection handler for streaming card generation in frontend/src/services/api.ts
- [x] T038 [US3] Update Result.tsx to handle streaming card events and display cards incrementally in frontend/src/pages/Result.tsx
- [x] T039 [US3] Add progress indicator for card generation when taking longer than 2 seconds in frontend/src/pages/Result.tsx
- [x] T040 [US3] Handle timeout and partial results gracefully in frontend/src/pages/Result.tsx

**Checkpoint**: At this point, User Stories 1, 2, AND 3 should all work independently. Card generation should complete in under 5 seconds with streaming support.

---

## Phase 6: User Story 4 - 知识卡片文本转语音功能（Priority: P2）

**Goal**: 在知识卡片上添加"听"按钮，支持将卡片内容转换为语音播放，支持播放控制

**Independent Test**: 在知识卡片上点击"听"按钮，验证语音是否正确播放卡片内容，支持播放控制。

### Frontend Implementation for User Story 4

- [x] T041 [P] [US4] Implement useTextToSpeech hook with play, pause, stop controls in frontend/src/hooks/useTextToSpeech.ts
- [x] T042 [US4] Add language detection and switching logic (Chinese/English) in frontend/src/hooks/useTextToSpeech.ts
- [x] T043 [US4] Add "听" button to ScienceCard component in frontend/src/components/cards/ScienceCard.tsx
- [x] T044 [US4] Add "听" button to PoetryCard component in frontend/src/components/cards/PoetryCard.tsx
- [x] T045 [US4] Add "听" button to EnglishCard component in frontend/src/components/cards/EnglishCard.tsx
- [x] T046 [US4] Extract audioText from card content for text-to-speech in all card components
- [x] T047 [US4] Implement global state management for current playing card to stop previous playback in frontend/src/pages/Result.tsx
- [x] T048 [US4] Add play/pause/stop button UI and controls in all card components
- [x] T049 [US4] Configure speech synthesis parameters (rate, pitch) for child-friendly audio in frontend/src/hooks/useTextToSpeech.ts
- [x] T050 [US4] Handle Web Speech API errors and browser compatibility in frontend/src/hooks/useTextToSpeech.ts
- [x] T051 [US4] Add audioText field extraction logic to convert card content to plain text in frontend/src/components/cards/*.tsx

**Checkpoint**: At this point, all user stories should be independently functional. Cards should support text-to-speech with full playback controls.

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

### Performance Optimization

- [x] T052 [P] Optimize Markdown rendering performance for long content in frontend/src/components/conversation/ConversationMessage.tsx
- [x] T053 [P] Add request cancellation for streaming connections in frontend/src/services/sse-post.ts
- [x] T054 [P] Implement connection pooling and reuse for SSE connections in frontend/src/services/sse.ts
- [x] T055 [P] Add rate limiting and timeout handling in backend/internal/handler/generatecardshandler.go

### Error Handling & Edge Cases

- [x] T056 [P] Add comprehensive error handling for network interruptions in frontend/src/pages/Result.tsx
- [x] T057 [P] Handle Markdown rendering errors gracefully in frontend/src/components/conversation/ConversationMessage.tsx
- [x] T058 [P] Add fallback for Web Speech API when unavailable in frontend/src/hooks/useTextToSpeech.ts
- [x] T059 [P] Handle concurrent card generation requests in backend/internal/logic/generatecardslogic.go
- [x] T060 [P] Add timeout and error recovery for card generation in backend/internal/agent/graph.go

### Data Consistency

- [x] T061 [P] Verify frontend and backend type definitions match in frontend/src/types/api.ts and backend/internal/types/types.go
- [x] T062 [P] Add validation for streamingText and markdown fields in backend/internal/logic/streamlogic.go
- [x] T063 [P] Ensure SSE event format consistency between frontend and backend

### Testing & Validation

- [ ] T064 [P] Test streaming message real-time rendering with network throttling
- [ ] T065 [P] Test Markdown rendering with various content types (code, lists, links)
- [ ] T066 [P] Performance test card generation to verify 5-second target
- [ ] T067 [P] Test text-to-speech with different languages and content lengths
- [ ] T068 [P] Run quickstart.md validation scenarios

### Documentation

- [x] T069 [P] Update API documentation for streaming endpoints
- [x] T070 [P] Add code comments for streaming and Markdown rendering logic
- [x] T071 [P] Document text-to-speech usage and browser compatibility

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3-6)**: All depend on Foundational phase completion
  - User stories can proceed in parallel (if staffed) or sequentially in priority order (P1 → P2)
- **Polish (Phase 7)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P1)**: Can start after Foundational (Phase 2) - Can work with US1 but independently testable
- **User Story 3 (P1)**: Can start after Foundational (Phase 2) - Independent, but benefits from US1/US2 streaming infrastructure
- **User Story 4 (P2)**: Can start after Foundational (Phase 2) - Independent, pure frontend feature

### Within Each User Story

- Frontend and backend tasks can be worked on in parallel (different files)
- Core implementation before integration
- Story complete before moving to next priority

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel
- All Foundational tasks marked [P] can run in parallel (within Phase 2)
- Once Foundational phase completes, all user stories can start in parallel (if team capacity allows)
- Frontend and backend tasks within a story can run in parallel
- Different user stories can be worked on in parallel by different team members

---

## Parallel Example: User Story 1

```bash
# Frontend and backend can work in parallel:
Frontend: "Optimize handleSendMessage function in frontend/src/pages/Result.tsx"
Backend: "Ensure streamlogic.go sends each text fragment immediately"

# Multiple frontend tasks can run in parallel:
Task: "Update useStreamConversation hook in frontend/src/hooks/useStreamConversation.ts"
Task: "Update ConversationMessage component in frontend/src/components/conversation/ConversationMessage.tsx"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1 (流式消息实时渲染)
4. **STOP and VALIDATE**: Test User Story 1 independently
5. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → Deploy/Demo (MVP!)
3. Add User Story 2 → Test independently → Deploy/Demo
4. Add User Story 3 → Test independently → Deploy/Demo
5. Add User Story 4 → Test independently → Deploy/Demo
6. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A (Frontend): User Story 1 frontend tasks
   - Developer B (Backend): User Story 1 backend tasks + User Story 3 backend tasks
   - Developer C (Frontend): User Story 2 frontend tasks + User Story 4
3. Stories complete and integrate independently

### Frontend/Backend Separation

**Frontend Team Focus**:
- User Story 1: Real-time rendering optimization
- User Story 2: Markdown rendering
- User Story 3: Streaming card reception
- User Story 4: Text-to-speech (complete feature)

**Backend Team Focus**:
- User Story 1: Ensure immediate SSE transmission
- User Story 2: Markdown detection and flagging
- User Story 3: Performance optimization and streaming (main focus)

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Frontend and backend tasks are clearly separated
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence

## Task Summary

- **Total Tasks**: 71
- **Setup Tasks**: 4 (T001-T004)
- **Foundational Tasks**: 4 (T005-T008)
- **User Story 1 Tasks**: 9 (T009-T017) - Frontend: 6, Backend: 3
- **User Story 2 Tasks**: 10 (T018-T027) - Frontend: 7, Backend: 3
- **User Story 3 Tasks**: 13 (T028-T040) - Frontend: 5, Backend: 8
- **User Story 4 Tasks**: 11 (T041-T051) - Frontend: 11, Backend: 0
- **Polish Tasks**: 20 (T052-T071)

**Parallel Opportunities**: 
- Frontend and backend can work in parallel for most user stories
- Multiple frontend tasks can run in parallel (different components)
- Multiple backend tasks can run in parallel (different handlers/logic)

**Suggested MVP Scope**: User Story 1 only (流式消息实时渲染) - 9 tasks total
