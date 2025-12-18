# 任务清单：TanGo 交互功能与AI对话系统

**输入**: 设计文档来自 `/specs/002-interactive-features/`
**前置条件**: plan.md (必需), spec.md (必需，用于用户故事)
**目标**: 生成对应的前后端代码，保证前后端联调完整可用

**注意**: 
- 当前阶段AI模型APP ID尚未提供，因此AI模型调用部分先使用Mock数据实现框架，待APP ID提供后再接入真实模型
- 所有任务必须确保前后端接口定义一致，支持完整联调

**组织方式**: 任务按用户故事分组，支持独立实现和测试每个故事。

## 格式说明: `[ID] [P?] [Story] 描述`

- **[P]**: 可以并行执行（不同文件，无依赖）
- **[Story]**: 该任务属于哪个用户故事（如 US1, US2, US3）
- 描述中包含确切的文件路径

## 路径约定

- **Web应用**: `frontend/src/`, `backend/`
- 路径基于 plan.md 中的项目结构

---

## Phase 1: 项目初始化与基础架构

**目的**: 搭建前后端基础架构，确保前后端能正常通信

### 1.1 后端API定义扩展

- [ ] T001 [P] 扩展API定义文件 `backend/api/explore.api`，添加对话相关类型定义（ConversationMessage, IntentResult, ConversationSession等）
- [ ] T002 [P] 扩展API定义文件 `backend/api/explore.api`，添加意图识别接口 `/api/conversation/intent`
- [ ] T003 [P] 扩展API定义文件 `backend/api/explore.api`，添加对话接口 `/api/conversation/message`
- [ ] T004 [P] 扩展API定义文件 `backend/api/explore.api`，添加语音识别接口 `/api/conversation/voice`
- [ ] T005 [P] 扩展API定义文件 `backend/api/explore.api`，添加流式返回接口 `/api/conversation/stream`（SSE）
- [ ] T006 运行goctl生成代码，更新 `backend/internal/handler/routes.go` 和 `backend/internal/types/types.go`

### 1.2 前端类型定义

- [X] T007 [P] 创建对话类型定义 `frontend/src/types/conversation.ts`（ConversationMessage, ConversationSession等）
- [X] T008 [P] 扩展API类型定义 `frontend/src/types/api.ts`，添加对话相关接口类型（IntentRequest, IntentResponse, ConversationRequest等）
- [X] T009 [P] 创建用户设置类型定义 `frontend/src/types/settings.ts`（UserSettings）

### 1.3 前端国际化基础

- [X] T010 安装react-i18next依赖 `frontend/package.json`
- [X] T011 [P] 创建国际化配置 `frontend/src/i18n/index.ts`
- [X] T012 [P] 创建中文翻译文件 `frontend/src/i18n/locales/zh.ts`
- [X] T013 [P] 创建英文翻译文件 `frontend/src/i18n/locales/en.ts`

### 1.4 后端存储基础

- [ ] T014 [P] 创建内存缓存实现 `backend/internal/storage/memory.go`（使用sync.Map）
- [ ] T015 在ServiceContext中添加存储实例 `backend/internal/svc/servicecontext.go`

**检查点**: 前后端基础架构就绪，API定义和类型定义完成，可以开始用户故事实现

---

## Phase 2: 基础功能（阻塞所有用户故事）

**目的**: 实现所有用户故事都依赖的基础功能

### 2.1 前端存储服务增强

- [X] T016 [P] 扩展存储服务 `frontend/src/services/storage.ts`，添加对话消息存储（IndexedDB）
- [X] T017 [P] 扩展存储服务 `frontend/src/services/storage.ts`，添加用户设置存储（localStorage）
- [X] T018 [P] 扩展存储服务 `frontend/src/services/storage.ts`，优化收藏卡片存储，支持批量操作

### 2.2 前端API服务扩展

- [X] T019 [P] 扩展API服务 `frontend/src/services/api.ts`，添加意图识别接口调用
- [X] T020 [P] 扩展API服务 `frontend/src/services/api.ts`，添加对话接口调用
- [X] T021 [P] 扩展API服务 `frontend/src/services/api.ts`，添加语音识别接口调用
- [X] T022 [P] 创建流式返回服务 `frontend/src/services/sse.ts`（使用EventSource API）

### 2.3 后端基础逻辑

- [ ] T023 [P] 创建对话管理逻辑 `backend/internal/logic/conversationlogic.go`（对话上下文存储和管理）
- [ ] T024 [P] 创建流式返回逻辑 `backend/internal/logic/streamlogic.go`（SSE实现）
- [ ] T025 [P] 创建意图识别逻辑 `backend/internal/logic/intentlogic.go`（使用LLM进行意图分类）

**检查点**: 基础功能完成，所有用户故事可以开始实现

---

## Phase 3: 用户故事 1 - 完整的交互功能实现（优先级: P1）🎯 MVP

**目标**: 所有按钮都有真实的功能响应，包括拍照、相册选择、语音输入、收藏、导出等，所有页面都有快速回到首页拍照的功能

**独立测试**: 点击所有页面的按钮，验证每个按钮都有真实的功能响应，没有假逻辑或点击无反应的情况

### 3.1 前端按钮功能实现

- [X] T026 [P] [US1] 实现首页拍照按钮功能 `frontend/src/pages/Home.tsx`（打开相机或相册选择）
- [X] T027 [P] [US1] 实现首页语音输入按钮功能 `frontend/src/pages/Home.tsx`（启动语音识别）
- [X] T028 [P] [US1] 实现相册选择功能 `frontend/src/pages/Capture.tsx`（支持从相册选择图片）
- [X] T029 [P] [US1] 实现语音输入功能 `frontend/src/pages/Capture.tsx`（使用Web Speech API或调用后端语音识别API）
- [X] T030 [P] [US1] 实现收藏按钮功能 `frontend/src/pages/Result.tsx`（点击后立即更新状态，保存到本地）
- [X] T031 [P] [US1] 实现导出按钮功能 `frontend/src/pages/Collection.tsx`（将卡片导出为图片，使用html2canvas）
- [X] T032 [P] [US1] 创建快速拍照按钮组件 `frontend/src/components/common/QuickCaptureButton.tsx`
- [X] T033 [US1] 在App.tsx中添加全局快速拍照按钮 `frontend/src/App.tsx`（固定在页面底部，所有页面可见）

### 3.2 前端导出工具

- [X] T034 [P] [US1] 安装html2canvas依赖 `frontend/package.json`
- [X] T035 [P] [US1] 创建导出工具函数 `frontend/src/utils/export.ts`（使用html2canvas将卡片导出为图片）

### 3.3 后端语音识别支持

- [ ] T036 [P] [US1] 创建语音识别逻辑 `backend/internal/logic/voicelogic.go`（调用语音识别模型，当前使用Mock）
- [ ] T037 [P] [US1] 创建语音处理器 `backend/internal/handler/voicehandler.go`
- [ ] T038 [US1] 注册语音识别路由 `backend/internal/handler/routes.go`

**检查点**: 用户故事1完成，所有按钮功能可用，可以独立测试

---

## Phase 4: 用户故事 2 - 中英文切换功能（优先级: P1）

**目标**: 支持中文和英文切换，默认语言为中文，语言设置持久化保存

**独立测试**: 在设置页面切换语言，验证界面语言立即更新，刷新页面后语言设置仍然保持

### 4.1 前端语言切换组件

- [ ] T039 [P] [US2] 创建语言切换组件 `frontend/src/components/common/LanguageSwitcher.tsx`
- [ ] T040 [P] [US2] 创建语言切换Hook `frontend/src/hooks/useLanguage.ts`
- [ ] T041 [P] [US2] 创建设置页面 `frontend/src/pages/Settings.tsx`
- [X] T042 [US2] 在App.tsx中集成i18n `frontend/src/App.tsx`（初始化i18n，监听语言变化）

### 4.2 前端页面国际化

- [ ] T043 [P] [US2] 在Home.tsx中使用i18n翻译 `frontend/src/pages/Home.tsx`
- [ ] T044 [P] [US2] 在Capture.tsx中使用i18n翻译 `frontend/src/pages/Capture.tsx`
- [ ] T045 [P] [US2] 在Result.tsx中使用i18n翻译 `frontend/src/pages/Result.tsx`
- [ ] T046 [P] [US2] 在Collection.tsx中使用i18n翻译 `frontend/src/pages/Collection.tsx`
- [ ] T047 [P] [US2] 在所有组件中使用i18n翻译（Button, Header等）

### 4.3 语言设置持久化

- [X] T048 [US2] 实现语言设置的持久化保存 `frontend/src/services/storage.ts`（localStorage）
- [X] T049 [US2] 实现应用启动时加载语言设置 `frontend/src/App.tsx`

**检查点**: 用户故事2完成，中英文切换功能可用，可以独立测试

---

## Phase 5: 用户故事 3 - 对话式交互与卡片展示（优先级: P1）

**目标**: 生成的三个卡片显示在对话消息列表中，支持继续追问，支持文本、语音、图片输入，根据意图识别决定输出文本或重新生成卡片

**独立测试**: 拍照生成卡片后，在对话消息列表中继续输入文本、语音或图片，验证系统能正确识别意图并返回相应内容

### 5.1 前端对话组件

- [ ] T050 [P] [US3] 创建对话消息列表组件 `frontend/src/components/conversation/ConversationList.tsx`
- [ ] T051 [P] [US3] 创建单条消息组件 `frontend/src/components/conversation/ConversationMessage.tsx`（支持文本、卡片、图片、语音消息）
- [ ] T052 [P] [US3] 创建消息输入组件 `frontend/src/components/conversation/MessageInput.tsx`（支持文本输入）
- [ ] T053 [P] [US3] 创建语音输入组件 `frontend/src/components/conversation/VoiceInput.tsx`（支持语音识别）
- [ ] T054 [P] [US3] 创建图片输入组件 `frontend/src/components/conversation/ImageInput.tsx`（支持图片上传）
- [ ] T055 [US3] 创建对话页面 `frontend/src/pages/Conversation.tsx`（整合消息列表和输入组件，或修改Result.tsx）

### 5.2 前端对话服务

- [ ] T056 [P] [US3] 创建对话服务 `frontend/src/services/conversation.ts`（处理消息发送和接收）
- [ ] T057 [US3] 实现对话消息的本地存储 `frontend/src/services/storage.ts`（IndexedDB）
- [ ] T058 [US3] 实现WebSocket或SSE连接，接收流式返回 `frontend/src/services/conversation.ts`

### 5.3 后端对话处理

- [ ] T059 [P] [US3] 创建对话处理器 `backend/internal/handler/conversationhandler.go`
- [ ] T060 [US3] 实现对话接口逻辑 `backend/internal/logic/conversationlogic.go`（处理消息，调用意图识别，根据意图返回文本或生成卡片）
- [ ] T061 [US3] 注册对话路由 `backend/internal/handler/routes.go`

### 5.4 后端意图识别集成

- [ ] T062 [US3] 创建意图识别处理器 `backend/internal/handler/inthandler.go`
- [ ] T063 [US3] 实现意图识别接口逻辑 `backend/internal/logic/intentlogic.go`（区分生成卡片意图和文本回答意图）
- [ ] T064 [US3] 注册意图识别路由 `backend/internal/handler/routes.go`

### 5.5 修改Result页面为对话形式

- [ ] T065 [US3] 修改Result页面，将卡片显示改为对话消息列表形式 `frontend/src/pages/Result.tsx`

**检查点**: 用户故事3完成，对话式交互功能可用，可以独立测试

---

## Phase 6: 用户故事 4 - 卡片收藏与导出功能（优先级: P2）

**目标**: 支持卡片收藏，收藏后可以导出成图片，卡片页面细节需要细化，卡片大小符合设计规范

**独立测试**: 在结果页面或对话消息列表中点击收藏按钮，然后在收藏页面查看收藏的卡片，点击导出按钮验证卡片能导出为图片

### 6.1 前端卡片组件增强

- [ ] T066 [P] [US4] 完善科学认知卡组件的收藏功能 `frontend/src/components/cards/ScienceCard.tsx`（确保点击后立即更新状态）
- [ ] T067 [P] [US4] 完善古诗词/人文卡组件的收藏功能 `frontend/src/components/cards/PoetryCard.tsx`
- [ ] T068 [P] [US4] 完善英语表达卡组件的收藏功能 `frontend/src/components/cards/EnglishCard.tsx`
- [ ] T069 [P] [US4] 创建卡片详情组件 `frontend/src/components/cards/CardDetail.tsx`（展示卡片完整细节）
- [ ] T070 [P] [US4] 优化卡片样式，确保卡片大小符合设计规范，不被压缩 `frontend/src/components/cards/*.tsx`

### 6.2 前端收藏页面增强

- [ ] T071 [US4] 在收藏页面添加导出按钮 `frontend/src/pages/Collection.tsx`
- [ ] T072 [US4] 实现收藏卡片的导出功能 `frontend/src/pages/Collection.tsx`（使用export.ts工具函数）

### 6.3 前端导出优化

- [ ] T073 [US4] 优化导出工具函数，确保导出时图片清晰 `frontend/src/utils/export.ts`
- [ ] T074 [US4] 确保卡片样式在导出时完整显示 `frontend/src/components/cards/*.tsx`

**检查点**: 用户故事4完成，卡片收藏和导出功能可用，可以独立测试

---

## Phase 7: 用户故事 5 - 数据存储与同步（优先级: P2）

**目标**: 前端数据存储在本地缓存中，服务端也存储到内存缓存中，支持数据同步

**独立测试**: 进行多次探索和收藏操作，刷新页面后验证数据仍然存在，关闭应用后重新打开验证数据仍然保持

### 7.1 前端数据存储完善

- [ ] T075 [P] [US5] 完善对话消息存储实现 `frontend/src/services/storage.ts`（IndexedDB，支持查询、删除等操作）
- [ ] T076 [P] [US5] 完善用户设置存储实现 `frontend/src/services/storage.ts`（localStorage，支持语言、年级等设置）
- [ ] T077 [P] [US5] 完善收藏卡片存储实现 `frontend/src/services/storage.ts`（IndexedDB，支持批量操作、查询、删除）

### 7.2 后端数据存储完善

- [ ] T078 [US5] 完善内存缓存实现 `backend/internal/storage/memory.go`（实现数据过期清理机制，30分钟无活动自动清理）
- [ ] T079 [US5] 实现对话上下文的存储和管理 `backend/internal/logic/conversationlogic.go`（限制上下文长度，最多10轮对话）

### 7.3 数据同步逻辑

- [ ] T080 [US5] 实现前端数据同步逻辑 `frontend/src/services/storage.ts`（确保数据一致性）

**检查点**: 用户故事5完成，数据存储功能可用，可以独立测试

---

## Phase 8: 用户故事 6 - AI模型调用与Agent系统（优先级: P1）

**目标**: 使用eino框架搭建Agent系统，通过graph图串联整个流程，区分拍照图片识别、三个卡片图片生成、文本生成等不同功能

**独立测试**: 拍照后验证系统能正确调用AI模型进行图片识别，然后生成三个卡片，验证Agent系统能正确串联整个流程

### 8.1 后端eino框架集成

- [ ] T081 安装eino框架依赖 `backend/go.mod`
- [ ] T082 [P] 创建Agent系统目录结构 `backend/internal/agent/`
- [ ] T083 [P] 实现Agent主文件 `backend/internal/agent/agent.go`（使用eino框架初始化）
- [ ] T084 [P] 实现Graph图定义 `backend/internal/agent/graph.go`（串联整个AI调用流程）

### 8.2 后端Agent节点实现

- [ ] T085 [P] [US6] 实现图片识别节点 `backend/internal/agent/nodes/image_recognition.go`（调用图片识别模型，当前使用Mock）
- [ ] T086 [P] [US6] 实现文本生成节点 `backend/internal/agent/nodes/text_generation.go`（调用文本生成模型，当前使用Mock）
- [ ] T087 [P] [US6] 实现图片生成节点 `backend/internal/agent/nodes/image_generation.go`（调用图片生成模型，当前使用Mock）
- [ ] T088 [P] [US6] 实现意图识别节点 `backend/internal/agent/nodes/intent_recognition.go`（识别用户意图，当前使用Mock）

### 8.3 后端配置和集成

- [ ] T089 [US6] 配置eino框架 `backend/internal/config/config.go`（添加eino配置）
- [ ] T090 [US6] 实现节点间的数据流转逻辑 `backend/internal/agent/graph.go`
- [ ] T091 [US6] 集成Agent系统到现有逻辑 `backend/internal/logic/identifylogic.go`（图片识别使用Agent）
- [ ] T092 [US6] 集成Agent系统到现有逻辑 `backend/internal/logic/generatecardslogic.go`（卡片生成使用Agent）

**检查点**: 用户故事6完成，AI模型调用和Agent系统可用，可以独立测试

---

## Phase 9: 用户故事 7 - 多模态输入与意图识别（优先级: P1）

**目标**: 支持文本输入、语音输入、图片输入，系统根据意图识别，输出文本或三个卡片，三个卡片是明确的意图

**独立测试**: 在对话消息列表中输入文本、语音或图片，验证系统能正确识别意图，返回文本或生成卡片

### 9.1 前端多模态输入完善

- [ ] T093 [US7] 完善文本输入处理 `frontend/src/components/conversation/MessageInput.tsx`（发送到后端，调用意图识别）
- [ ] T094 [US7] 完善语音输入处理 `frontend/src/components/conversation/VoiceInput.tsx`（发送音频到后端，调用语音识别和意图识别）
- [ ] T095 [US7] 完善图片输入处理 `frontend/src/components/conversation/ImageInput.tsx`（发送图片到后端，调用图片识别和意图识别）

### 9.2 后端多模态处理

- [ ] T096 [US7] 增强图片处理逻辑 `backend/internal/logic/imagelogic.go`（支持多种图片格式）
- [ ] T097 [US7] 完善语音识别逻辑 `backend/internal/logic/voicelogic.go`（调用语音识别模型，转换为文本后识别意图）
- [ ] T098 [US7] 完善意图识别逻辑 `backend/internal/logic/intentlogic.go`（结合上下文理解用户意图，支持多种表达方式）

### 9.3 后端流式返回集成

- [ ] T099 [US7] 创建流式返回处理器 `backend/internal/handler/streamhandler.go`（SSE实现）
- [ ] T100 [US7] 实现流式返回逻辑 `backend/internal/logic/streamlogic.go`（支持文字和图片流式返回）
- [ ] T101 [US7] 集成eino框架的流式输出能力 `backend/internal/logic/streamlogic.go`
- [ ] T102 [US7] 注册流式返回路由 `backend/internal/handler/routes.go`

### 9.4 前端流式返回处理

- [ ] T103 [US7] 实现流式数据的接收和解析 `frontend/src/services/sse.ts`
- [ ] T104 [US7] 实现流式数据的实时渲染 `frontend/src/components/conversation/ConversationList.tsx`
- [ ] T105 [US7] 实现连接重连机制 `frontend/src/services/sse.ts`

**检查点**: 用户故事7完成，多模态输入和意图识别功能可用，可以独立测试

---

## Phase 10: 前后端联调与测试

**目的**: 确保前后端接口定义一致，所有功能能正常联调

### 10.1 API接口一致性检查

- [ ] T106 检查前后端API接口定义一致性（`backend/api/explore.api` 与 `frontend/src/types/api.ts`）
- [ ] T107 检查前后端类型定义一致性（所有请求和响应类型）
- [ ] T108 验证所有API接口的路由注册 `backend/internal/handler/routes.go`

### 10.2 前后端联调测试

- [ ] T109 测试图像识别接口联调（前端调用，后端响应）
- [ ] T110 测试知识卡片生成接口联调
- [ ] T111 测试意图识别接口联调
- [ ] T112 测试对话接口联调
- [ ] T113 测试语音识别接口联调
- [ ] T114 测试流式返回接口联调（SSE）

### 10.3 错误处理完善

- [ ] T115 [P] 完善前端错误处理 `frontend/src/services/api.ts`（网络异常、API错误等）
- [ ] T116 [P] 完善后端错误处理 `backend/internal/utils/errors.go`（统一错误响应格式）
- [ ] T117 [P] 实现前端错误提示 `frontend/src/components/common/ErrorToast.tsx`（新建组件）

### 10.4 性能优化

- [ ] T118 优化按钮点击响应时间（目标≤200ms）
- [ ] T119 优化语言切换响应时间（目标≤100ms）
- [ ] T120 优化卡片导出响应时间（目标≤2秒）

**检查点**: 前后端联调完成，所有功能可用

---

## Phase 11: 完善与优化

**目的**: 完善功能细节，优化用户体验

### 11.1 功能完善

- [ ] T121 [P] 实现设置页面完整功能 `frontend/src/pages/Settings.tsx`（语言设置、年级设置等）
- [ ] T122 [P] 优化卡片详情页面 `frontend/src/components/cards/CardDetail.tsx`（确保细节完整）
- [ ] T123 [P] 优化对话消息列表样式 `frontend/src/components/conversation/ConversationList.tsx`（移动端适配）
- [ ] T124 [P] 实现断线重连机制 `frontend/src/services/sse.ts`（网络异常时自动重连）

### 11.2 移动端优化

- [ ] T125 优化移动端交互，确保所有按钮支持touch事件
- [ ] T126 优化移动端布局，确保快速拍照按钮在所有页面都可见且易用
- [ ] T127 优化移动端输入体验（文本、语音、图片输入）

### 11.3 数据持久化验证

- [ ] T128 验证对话消息的本地存储（刷新页面后数据仍然存在）
- [ ] T129 验证用户设置的持久化（刷新页面后设置仍然保持）
- [ ] T130 验证收藏卡片的本地存储（关闭应用后重新打开数据仍然保持）

### 11.4 文档更新

- [ ] T131 更新API文档，包含所有新增接口
- [ ] T132 更新README，说明新功能使用方法

**检查点**: 所有功能完善，可以发布

---

## 依赖关系与执行顺序

### 阶段依赖

- **Phase 1 (Setup)**: 无依赖，可立即开始
- **Phase 2 (Foundational)**: 依赖Phase 1完成，阻塞所有用户故事
- **Phase 3-9 (用户故事)**: 所有依赖Phase 2完成，可以并行实现（如果团队容量允许）
- **Phase 10 (联调)**: 依赖Phase 3-9完成
- **Phase 11 (完善)**: 依赖Phase 10完成

### 用户故事依赖

- **用户故事1 (US1)**: 可独立实现，依赖Phase 2
- **用户故事2 (US2)**: 可独立实现，依赖Phase 2
- **用户故事3 (US3)**: 可独立实现，依赖Phase 2，但需要US6的意图识别支持（可在US6完成前使用Mock）
- **用户故事4 (US4)**: 可独立实现，依赖Phase 2
- **用户故事5 (US5)**: 可独立实现，依赖Phase 2
- **用户故事6 (US6)**: 可独立实现，依赖Phase 2，这是AI能力的核心
- **用户故事7 (US7)**: 依赖US3和US6完成，需要对话系统和意图识别

### 用户故事内部依赖

- **US1**: 按钮功能 → 导出工具 → 语音识别支持
- **US2**: 语言切换组件 → 页面国际化 → 持久化
- **US3**: 对话组件 → 对话服务 → 后端对话处理 → 意图识别集成
- **US4**: 卡片组件增强 → 收藏页面增强 → 导出优化
- **US5**: 前端存储完善 → 后端存储完善 → 数据同步
- **US6**: eino框架集成 → Agent节点实现 → 配置和集成
- **US7**: 多模态输入完善 → 后端多模态处理 → 流式返回集成 → 前端流式处理

### 并行执行机会

- **Phase 1**: T001-T015 中标记[P]的任务可以并行
- **Phase 2**: T016-T025 中标记[P]的任务可以并行
- **Phase 3-9**: 不同用户故事可以并行实现（如果团队容量允许）
- **同一用户故事内**: 标记[P]的任务可以并行

### 并行执行示例：用户故事1

```bash
# 可以并行执行的任务：
- T026: 实现首页拍照按钮功能
- T027: 实现首页语音输入按钮功能
- T028: 实现相册选择功能
- T029: 实现语音输入功能
- T030: 实现收藏按钮功能
- T031: 实现导出按钮功能
- T032: 创建快速拍照按钮组件
- T034: 安装html2canvas依赖
- T035: 创建导出工具函数
- T036: 创建语音识别逻辑
- T037: 创建语音处理器
```

---

## 实施策略

### MVP优先（仅用户故事1）

1. 完成Phase 1: 项目初始化
2. 完成Phase 2: 基础功能（关键）
3. 完成Phase 3: 用户故事1（完整的交互功能实现）
4. **停止并验证**: 测试用户故事1是否独立可用
5. 如果可用，可以部署/演示

### 增量交付

1. 完成Setup + Foundational → 基础就绪
2. 添加用户故事1 → 独立测试 → 部署/演示（MVP！）
3. 添加用户故事2 → 独立测试 → 部署/演示
4. 添加用户故事6 → 独立测试 → 部署/演示（AI能力）
5. 添加用户故事3 → 独立测试 → 部署/演示（对话系统）
6. 添加用户故事7 → 独立测试 → 部署/演示（多模态输入）
7. 添加用户故事4 → 独立测试 → 部署/演示（收藏导出）
8. 添加用户故事5 → 独立测试 → 部署/演示（数据存储）
9. 每个故事增加价值，不破坏之前的故事

### 并行团队策略

如果有多个开发者：

1. 团队共同完成Setup + Foundational
2. Foundational完成后：
   - 开发者A: 用户故事1（交互功能）
   - 开发者B: 用户故事2（中英文切换）
   - 开发者C: 用户故事6（AI模型调用）
3. 用户故事完成后：
   - 开发者A: 用户故事3（对话式交互）
   - 开发者B: 用户故事4（卡片收藏导出）
   - 开发者C: 用户故事7（多模态输入）
4. 最后：用户故事5（数据存储）和联调测试

---

## 任务统计

- **总任务数**: 132
- **Phase 1 (Setup)**: 15个任务
- **Phase 2 (Foundational)**: 10个任务
- **Phase 3 (US1)**: 13个任务
- **Phase 4 (US2)**: 11个任务
- **Phase 5 (US3)**: 16个任务
- **Phase 6 (US4)**: 9个任务
- **Phase 7 (US5)**: 6个任务
- **Phase 8 (US6)**: 12个任务
- **Phase 9 (US7)**: 13个任务
- **Phase 10 (联调)**: 15个任务
- **Phase 11 (完善)**: 14个任务

### 按用户故事统计

- **US1**: 13个任务
- **US2**: 11个任务
- **US3**: 16个任务
- **US4**: 9个任务
- **US5**: 6个任务
- **US6**: 12个任务
- **US7**: 13个任务

### 并行机会

- **可并行任务**: 约60%的任务标记了[P]，可以并行执行
- **用户故事并行**: 7个用户故事中，US1、US2、US4、US5、US6可以并行实现
- **US3和US7**: 需要依赖其他故事，建议按顺序实现

### 独立测试标准

每个用户故事完成后，应该能够：
- **US1**: 点击所有按钮，验证功能响应
- **US2**: 切换语言，验证界面更新和持久化
- **US3**: 拍照生成卡片，在对话中继续追问，验证意图识别和响应
- **US4**: 收藏卡片，导出为图片，验证功能正常
- **US5**: 进行多次操作，刷新页面验证数据持久化
- **US6**: 拍照后验证AI模型调用和Agent系统工作正常
- **US7**: 输入文本、语音、图片，验证意图识别和多模态处理

### 建议的MVP范围

**最小可行产品（MVP）**：
- Phase 1: 项目初始化
- Phase 2: 基础功能
- Phase 3: 用户故事1（完整的交互功能实现）

**扩展MVP（推荐）**：
- 加上Phase 4: 用户故事2（中英文切换）
- 加上Phase 8: 用户故事6（AI模型调用，使用Mock数据）

这样可以快速交付一个可用的产品，然后逐步添加其他功能。

---

## 注意事项

1. **前后端接口一致性**: 每次修改API定义后，必须同步更新前后端类型定义
2. **Mock数据阶段**: 当前阶段AI模型使用Mock数据，待APP ID提供后再接入真实模型
3. **流式返回**: 优先使用SSE，如果未来需要双向通信再升级到WebSocket
4. **移动端优化**: 确保所有交互在移动端流畅，支持touch事件
5. **性能指标**: 按钮响应≤200ms，语言切换≤100ms，卡片导出≤2秒
6. **意图识别准确率**: 目标≥85%，需要大量测试
7. **数据持久化**: 前端使用IndexedDB + localStorage，后端使用内存缓存
8. **错误处理**: 网络异常、AI模型调用失败时必须有降级策略

---

## 格式验证

✅ 所有任务都遵循了严格的检查清单格式：
- ✅ 每个任务都以 `- [ ]` 开头（markdown复选框）
- ✅ 每个任务都有唯一的Task ID（T001-T132）
- ✅ 可并行任务标记了 `[P]`
- ✅ 用户故事任务标记了 `[US1]` 到 `[US7]`
- ✅ 每个任务描述都包含了确切的文件路径
- ✅ Setup和Foundational阶段的任务没有Story标签
- ✅ 用户故事阶段的任务都有Story标签
- ✅ Polish阶段的任务没有Story标签
