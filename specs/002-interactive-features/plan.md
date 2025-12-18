# 实现计划：TanGo 交互功能与AI对话系统

**分支**: `002-interactive-features` | **日期**: 2025-12-18 | **规范**: [spec.md](./spec.md)
**输入**: 功能规范来自 `/specs/002-interactive-features/spec.md`

## Summary

实现 TanGo 交互功能与AI对话系统，包括：前端所有按钮的真实功能实现、中英文切换、对话式交互系统、卡片收藏导出、AI模型调用与Agent系统（使用eino框架）、多模态输入与意图识别。采用前后端分离架构，前端使用 React 18 + Vite + Tailwind CSS（移动端优先），后端使用 go-zero + eino 框架，通过 graph 图串联整个AI调用流程，支持流式返回。

## Technical Context

**Language/Version**: 
- 前端：TypeScript (ES2020+), React 18
- 后端：Go 1.25.3

**Primary Dependencies**: 
- 前端：React 18, Vite, Tailwind CSS, React Router, Axios, react-i18next (国际化)
- 后端：go-zero v1.9.3, eino (字节云原生AI框架), WebSocket/SSE (流式传输)
- AI模型：图片识别模型、文本生成模型、图片生成模型（通过eino框架调用）

**Storage**: 
- 前端：IndexedDB（探索记录、收藏卡片）+ localStorage（用户设置）
- 后端：内存缓存（sync.Map，用于对话上下文、临时数据）

**Testing**: 
- 前端：Vitest + React Testing Library
- 后端：Go testing package

**Target Platform**: 
- Web H5（移动端优先，兼容PC端）
- 现代浏览器（支持WebSocket、SSE、IndexedDB）

**Project Type**: Web应用（前后端分离）

**Performance Goals**: 
- 按钮点击响应时间≤200ms（90%请求）
- 语言切换响应时间≤100ms
- 意图识别响应时间≤1秒（90%请求）
- 流式返回延迟≤100ms（首字节时间）
- 卡片导出响应时间≤2秒（90%请求）

**Constraints**: 
- 必须使用eino框架搭建Agent系统
- 必须支持流式返回（WebSocket或SSE）
- 移动端交互必须流畅，支持touch事件
- 卡片大小必须符合设计规范，不能被压缩

**Scale/Scope**: 
- 支持多轮对话（上下文长度：10轮）
- 支持并发用户数：100+（MVP版本）
- 意图识别准确率≥85%（目标90%+）

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**规范检查项**（基于 `.specify/memory/constitution.md`）：

- [x] **原则一：中文优先规范** - 所有文档和生成内容必须使用中文（除非技术限制）
- [x] **原则二：K12 教育游戏化设计规范** - 设计必须符合儿童友好性、游戏化元素、玩中学理念
- [x] **原则三：可发布应用规范** - 实现必须达到生产级标准，可正常运行和发布
- [x] **原则四：多语言和年级设置规范** - 支持中英文设置和K12年级设置，默认中文
- [x] **原则五：AI优先（模型优先）规范** - 所有AI功能通过后端统一调用，支持流式返回
- [x] **原则六：移动端优先规范** - 确保移动端交互完整性，统一拍照入口，支持随时随地探索

**合规性说明**：所有设计均符合项目规范要求，无违反项。

## 前后端改动点详细列表

### 前端改动点

#### 1. 按钮功能实现

**文件路径**：
- `frontend/src/pages/Home.tsx` - 首页按钮功能
- `frontend/src/pages/Capture.tsx` - 拍照页面按钮功能
- `frontend/src/pages/Result.tsx` - 结果页面按钮功能
- `frontend/src/pages/Collection.tsx` - 收藏页面按钮功能
- `frontend/src/components/common/Button.tsx` - 通用按钮组件

**改动内容**：
- [ ] 实现首页拍照按钮：打开相机或相册选择
- [ ] 实现首页语音输入按钮：启动语音识别
- [ ] 实现相册选择功能：支持从相册选择图片
- [ ] 实现语音输入功能：使用Web Speech API或调用后端语音识别API
- [ ] 实现所有页面的快速拍照按钮：固定在页面底部/顶部，始终可见
- [ ] 实现收藏按钮：点击后立即更新状态，保存到本地
- [ ] 实现导出按钮：将卡片导出为图片（使用html2canvas或类似库）
- [ ] 实现设置按钮：打开设置页面

#### 2. 中英文切换功能

**文件路径**：
- `frontend/src/i18n/` - 国际化配置（新建目录）
  - `index.ts` - i18n初始化配置
  - `locales/zh.ts` - 中文翻译文件
  - `locales/en.ts` - 英文翻译文件
- `frontend/src/components/common/LanguageSwitcher.tsx` - 语言切换组件（新建）
- `frontend/src/pages/Settings.tsx` - 设置页面（新建）
- `frontend/src/services/storage.ts` - 添加语言设置存储
- `frontend/src/hooks/useLanguage.ts` - 语言切换Hook（新建）

**改动内容**：
- [ ] 安装react-i18next依赖
- [ ] 创建国际化配置文件（i18n/index.ts）
- [ ] 创建中英文翻译文件（locales/zh.ts, locales/en.ts）
- [ ] 实现语言切换组件（LanguageSwitcher.tsx）
- [ ] 实现设置页面（Settings.tsx）
- [ ] 在用户设置存储中添加语言字段
- [ ] 实现语言切换Hook（useLanguage.ts）
- [ ] 在所有页面组件中使用i18n翻译
- [ ] 实现语言设置的持久化保存（localStorage）

#### 3. 对话式交互系统

**文件路径**：
- `frontend/src/components/conversation/` - 对话相关组件（新建目录）
  - `ConversationList.tsx` - 对话消息列表组件（新建）
  - `ConversationMessage.tsx` - 单条消息组件（新建）
  - `MessageInput.tsx` - 消息输入组件（新建）
  - `VoiceInput.tsx` - 语音输入组件（新建）
  - `ImageInput.tsx` - 图片输入组件（新建）
- `frontend/src/pages/Conversation.tsx` - 对话页面（新建，替换Result.tsx或合并）
- `frontend/src/types/conversation.ts` - 对话相关类型定义（新建）
- `frontend/src/services/conversation.ts` - 对话服务（新建）
- `frontend/src/services/storage.ts` - 添加对话消息存储

**改动内容**：
- [ ] 创建对话消息类型定义（ConversationMessage）
- [ ] 实现对话消息列表组件（ConversationList.tsx）
- [ ] 实现单条消息组件（ConversationMessage.tsx），支持文本、卡片、图片、语音消息
- [ ] 实现消息输入组件（MessageInput.tsx），支持文本输入
- [ ] 实现语音输入组件（VoiceInput.tsx），支持语音识别
- [ ] 实现图片输入组件（ImageInput.tsx），支持图片上传
- [ ] 创建对话页面（Conversation.tsx），整合消息列表和输入组件
- [ ] 实现对话服务（conversation.ts），处理消息发送和接收
- [ ] 实现对话消息的本地存储（IndexedDB）
- [ ] 修改Result页面，将卡片显示改为对话消息列表形式
- [ ] 实现WebSocket或SSE连接，接收流式返回

#### 4. 卡片收藏和导出功能

**文件路径**：
- `frontend/src/components/cards/ScienceCard.tsx` - 科学认知卡组件
- `frontend/src/components/cards/PoetryCard.tsx` - 古诗词/人文卡组件
- `frontend/src/components/cards/EnglishCard.tsx` - 英语表达卡组件
- `frontend/src/components/cards/CardDetail.tsx` - 卡片详情组件（新建）
- `frontend/src/utils/export.ts` - 导出工具函数（新建）
- `frontend/src/services/storage.ts` - 添加收藏功能

**改动内容**：
- [ ] 完善卡片组件的收藏功能，确保点击后立即更新状态
- [ ] 实现卡片详情页面（CardDetail.tsx），展示卡片完整细节
- [ ] 确保卡片大小符合设计规范，使用固定尺寸，不被压缩
- [ ] 实现卡片导出功能（export.ts），使用html2canvas将卡片导出为图片
- [ ] 优化卡片样式，确保导出时图片清晰
- [ ] 实现收藏卡片的本地存储（IndexedDB）
- [ ] 在收藏页面添加导出按钮

#### 5. 统一拍照入口

**文件路径**：
- `frontend/src/components/common/QuickCaptureButton.tsx` - 快速拍照按钮组件（新建）
- `frontend/src/App.tsx` - 应用根组件，添加全局快速拍照按钮

**改动内容**：
- [ ] 创建快速拍照按钮组件（QuickCaptureButton.tsx）
- [ ] 在App.tsx中添加全局快速拍照按钮，固定在页面底部
- [ ] 确保快速拍照按钮在所有页面都可见
- [ ] 实现点击后跳转到拍照页面的功能

#### 6. 数据存储增强

**文件路径**：
- `frontend/src/services/storage.ts` - 存储服务
- `frontend/src/types/exploration.ts` - 类型定义

**改动内容**：
- [ ] 添加对话消息存储（IndexedDB）
- [ ] 添加用户设置存储（localStorage），包括语言、年级
- [ ] 优化收藏卡片存储，支持批量操作
- [ ] 实现数据同步逻辑（前端本地存储）

### 后端改动点

#### 1. AI模型调用与Agent系统（eino框架）

**文件路径**：
- `backend/internal/agent/` - Agent系统（新建目录）
  - `agent.go` - Agent主文件（新建）
  - `graph.go` - Graph图定义（新建）
  - `nodes/` - 节点实现（新建目录）
    - `image_recognition.go` - 图片识别节点（新建）
    - `text_generation.go` - 文本生成节点（新建）
    - `image_generation.go` - 图片生成节点（新建）
    - `intent_recognition.go` - 意图识别节点（新建）
- `backend/internal/config/config.go` - 添加eino配置
- `backend/go.mod` - 添加eino依赖

**改动内容**：
- [ ] 安装eino框架依赖
- [ ] 创建Agent系统目录结构
- [ ] 实现Agent主文件（agent.go），使用eino框架初始化
- [ ] 实现Graph图定义（graph.go），串联整个AI调用流程
- [ ] 实现图片识别节点（image_recognition.go），调用图片识别模型
- [ ] 实现文本生成节点（text_generation.go），调用文本生成模型
- [ ] 实现图片生成节点（image_generation.go），调用图片生成模型
- [ ] 实现意图识别节点（intent_recognition.go），识别用户意图
- [ ] 配置eino框架，连接AI模型服务
- [ ] 实现节点间的数据流转逻辑

#### 2. 意图识别功能

**文件路径**：
- `backend/internal/logic/intentlogic.go` - 意图识别逻辑（新建）
- `backend/internal/handler/inthandler.go` - 意图识别处理器（新建）
- `backend/api/explore.api` - 添加意图识别API定义

**改动内容**：
- [ ] 创建意图识别逻辑（intentlogic.go）
- [ ] 实现意图识别算法，区分生成卡片意图和文本回答意图
- [ ] 创建意图识别处理器（inthandler.go）
- [ ] 在API定义中添加意图识别接口
- [ ] 实现意图识别的路由注册

#### 3. 流式返回支持

**文件路径**：
- `backend/internal/handler/streamhandler.go` - 流式返回处理器（新建）
- `backend/internal/logic/streamlogic.go` - 流式返回逻辑（新建）
- `backend/api/explore.api` - 添加流式返回API定义

**改动内容**：
- [ ] 实现WebSocket处理器（streamhandler.go）
- [ ] 或实现SSE处理器（Server-Sent Events）
- [ ] 实现流式返回逻辑（streamlogic.go），支持文字和图片流式返回
- [ ] 在API定义中添加流式返回接口
- [ ] 实现流式返回的路由注册
- [ ] 集成eino框架的流式输出能力

#### 4. 对话管理功能

**文件路径**：
- `backend/internal/logic/conversationlogic.go` - 对话管理逻辑（新建）
- `backend/internal/handler/conversationhandler.go` - 对话处理器（新建）
- `backend/internal/storage/` - 存储层（新建目录）
  - `memory.go` - 内存缓存实现（新建）
- `backend/api/explore.api` - 添加对话API定义

**改动内容**：
- [ ] 创建对话管理逻辑（conversationlogic.go）
- [ ] 实现对话上下文的存储和管理（内存缓存）
- [ ] 创建对话处理器（conversationhandler.go）
- [ ] 实现内存缓存存储（memory.go），使用sync.Map
- [ ] 在API定义中添加对话相关接口
- [ ] 实现对话的路由注册
- [ ] 实现多轮对话的上下文保持

#### 5. 多模态输入处理

**文件路径**：
- `backend/internal/logic/voicelogic.go` - 语音识别逻辑（新建）
- `backend/internal/logic/imagelogic.go` - 图片处理逻辑（增强）
- `backend/internal/handler/voicehandler.go` - 语音处理器（新建）
- `backend/api/explore.api` - 添加语音识别API定义

**改动内容**：
- [ ] 创建语音识别逻辑（voicelogic.go），调用语音识别模型
- [ ] 增强图片处理逻辑（imagelogic.go），支持多种图片格式
- [ ] 创建语音处理器（voicehandler.go）
- [ ] 在API定义中添加语音识别接口
- [ ] 实现语音识别的路由注册
- [ ] 集成eino框架的语音识别能力

#### 6. 数据存储（内存缓存）

**文件路径**：
- `backend/internal/storage/memory.go` - 内存缓存实现（新建）
- `backend/internal/svc/servicecontext.go` - 添加存储上下文

**改动内容**：
- [ ] 实现内存缓存存储（memory.go），使用sync.Map
- [ ] 实现对话上下文的存储
- [ ] 实现临时数据的存储（分享链接等）
- [ ] 在ServiceContext中添加存储实例
- [ ] 实现数据过期清理机制

### 前后端通信改动

#### 1. WebSocket/SSE连接

**文件路径**：
- `frontend/src/services/websocket.ts` - WebSocket服务（新建）
- `frontend/src/services/sse.ts` - SSE服务（新建，如果使用SSE）
- `backend/internal/handler/streamhandler.go` - 流式返回处理器

**改动内容**：
- [ ] 实现WebSocket客户端连接（websocket.ts）
- [ ] 或实现SSE客户端连接（sse.ts）
- [ ] 实现流式数据的接收和解析
- [ ] 实现流式数据的实时渲染
- [ ] 实现连接重连机制
- [ ] 实现错误处理和降级策略

#### 2. API接口扩展

**文件路径**：
- `backend/api/explore.api` - API定义文件
- `frontend/src/services/api.ts` - API服务封装

**改动内容**：
- [ ] 添加意图识别API接口
- [ ] 添加对话管理API接口
- [ ] 添加语音识别API接口
- [ ] 添加流式返回API接口
- [ ] 更新前端API服务，添加新接口调用
- [ ] 实现API错误处理和重试机制

## Project Structure

### Documentation (this feature)

```text
specs/002-interactive-features/
├── plan.md              # This file (/speckit.plan command output)
├── spec.md              # 功能规范
└── tasks.md             # 任务清单（待生成）
```

### Source Code (repository root)

```text
frontend/
├── src/
│   ├── components/
│   │   ├── common/
│   │   │   ├── Button.tsx              # 通用按钮组件（需增强）
│   │   │   ├── LanguageSwitcher.tsx     # 语言切换组件（新建）
│   │   │   └── QuickCaptureButton.tsx   # 快速拍照按钮（新建）
│   │   ├── conversation/                # 对话相关组件（新建目录）
│   │   │   ├── ConversationList.tsx     # 对话消息列表（新建）
│   │   │   ├── ConversationMessage.tsx  # 单条消息（新建）
│   │   │   ├── MessageInput.tsx         # 消息输入（新建）
│   │   │   ├── VoiceInput.tsx           # 语音输入（新建）
│   │   │   └── ImageInput.tsx           # 图片输入（新建）
│   │   ├── cards/
│   │   │   ├── ScienceCard.tsx          # 科学认知卡（需增强收藏功能）
│   │   │   ├── PoetryCard.tsx           # 古诗词/人文卡（需增强收藏功能）
│   │   │   ├── EnglishCard.tsx          # 英语表达卡（需增强收藏功能）
│   │   │   └── CardDetail.tsx           # 卡片详情（新建）
│   │   └── collection/
│   │       └── CollectionCard.tsx      # 收藏卡片（需增强导出功能）
│   ├── pages/
│   │   ├── Home.tsx                     # 首页（需增强按钮功能）
│   │   ├── Capture.tsx                  # 拍照页面（需增强按钮功能）
│   │   ├── Result.tsx                   # 结果页面（需改为对话页面）
│   │   ├── Conversation.tsx             # 对话页面（新建，或合并到Result）
│   │   ├── Collection.tsx               # 收藏页面（需增强导出功能）
│   │   └── Settings.tsx                 # 设置页面（新建）
│   ├── services/
│   │   ├── api.ts                       # API服务（需扩展）
│   │   ├── storage.ts                   # 存储服务（需增强）
│   │   ├── conversation.ts              # 对话服务（新建）
│   │   ├── websocket.ts                 # WebSocket服务（新建）
│   │   └── sse.ts                       # SSE服务（新建，如果使用SSE）
│   ├── hooks/
│   │   └── useLanguage.ts               # 语言切换Hook（新建）
│   ├── i18n/                            # 国际化（新建目录）
│   │   ├── index.ts                     # i18n配置（新建）
│   │   └── locales/
│   │       ├── zh.ts                    # 中文翻译（新建）
│   │       └── en.ts                    # 英文翻译（新建）
│   ├── types/
│   │   ├── conversation.ts              # 对话类型定义（新建）
│   │   └── exploration.ts               # 探索类型（需扩展）
│   └── utils/
│       └── export.ts                    # 导出工具（新建）
│
backend/
├── internal/
│   ├── agent/                           # Agent系统（新建目录）
│   │   ├── agent.go                     # Agent主文件（新建）
│   │   ├── graph.go                     # Graph图定义（新建）
│   │   └── nodes/                       # 节点实现（新建目录）
│   │       ├── image_recognition.go     # 图片识别节点（新建）
│   │       ├── text_generation.go       # 文本生成节点（新建）
│   │       ├── image_generation.go       # 图片生成节点（新建）
│   │       └── intent_recognition.go    # 意图识别节点（新建）
│   ├── handler/
│   │   ├── conversationhandler.go       # 对话处理器（新建）
│   │   ├── streamhandler.go             # 流式返回处理器（新建）
│   │   ├── voicehandler.go              # 语音处理器（新建）
│   │   └── inthandler.go                 # 意图识别处理器（新建）
│   ├── logic/
│   │   ├── conversationlogic.go          # 对话管理逻辑（新建）
│   │   ├── streamlogic.go                # 流式返回逻辑（新建）
│   │   ├── voicelogic.go                 # 语音识别逻辑（新建）
│   │   └── intentlogic.go                # 意图识别逻辑（新建）
│   ├── storage/
│   │   └── memory.go                    # 内存缓存实现（新建）
│   └── svc/
│       └── servicecontext.go             # 添加存储上下文
├── api/
│   └── explore.api                      # API定义（需扩展）
└── go.mod                                # 添加eino依赖
```

## 技术决策与研究

### 1. eino框架集成方案

**决策**：使用字节的eino框架搭建Agent系统，通过graph图串联整个AI调用流程

**理由**：
- eino是字节开源的云原生AI框架，支持多种AI模型
- 提供统一的模型调用接口，简化集成
- 支持graph图模式，可以灵活串联不同的AI能力
- 支持流式输出，满足实时交互需求

**实现方式**：
- 使用单Agent模式，通过graph图定义工作流
- Graph节点包括：图片识别、文本生成、图片生成、意图识别
- 每个节点通过eino框架调用相应的AI模型
- 节点间通过数据流传递结果

**待确认**：
- eino与go-zero的集成方式
- Graph图的定义和配置方法
- 支持的AI模型类型和配置方法
- APP ID申请和配置流程

### 2. 流式返回方案选择

**决策**：优先使用Server-Sent Events (SSE)，备选WebSocket

**理由**：
- SSE更简单，基于HTTP，易于实现和调试
- 单向流式传输（服务器到客户端）满足需求
- 自动重连机制，连接管理简单
- 如果未来需要双向通信，再升级到WebSocket

**实现方式**：
- 后端实现SSE处理器，通过HTTP响应流式返回数据
- 前端使用EventSource API接收流式数据
- 支持文字和图片的流式传输
- 实现断线重连机制

**替代方案**：WebSocket（如果需要双向通信或更复杂的交互）

### 3. 意图识别实现方案

**决策**：使用LLM进行意图识别，结合规则判断

**理由**：
- LLM能理解自然语言，识别准确率高
- 支持多种表达方式（"帮我生成小卡片"、"生成卡片"等）
- 可以结合上下文理解用户意图
- 通过eino框架调用，统一管理

**实现方式**：
- 使用LLM进行意图分类（generate_cards/text_response/image_recognition）
- 结合关键词规则作为快速判断
- 返回意图类型和置信度
- 根据意图类型调用相应的处理逻辑

### 4. 多模态输入处理方案

**决策**：统一通过后端处理，前端只负责采集和展示

**理由**：
- 符合AI优先原则，所有AI能力通过后端统一调用
- 便于统一管理和优化
- 支持流式返回，提升用户体验

**实现方式**：
- **文本输入**：直接发送到后端，通过意图识别处理
- **语音输入**：前端使用Web Speech API或发送音频到后端，后端调用语音识别模型
- **图片输入**：发送图片到后端，后端调用图片识别模型

### 5. 对话上下文管理方案

**决策**：使用内存缓存（sync.Map）存储对话上下文

**理由**：
- MVP版本不需要持久化存储
- 内存缓存速度快，满足实时交互需求
- 使用sync.Map保证并发安全
- 可以设置过期时间，自动清理

**实现方式**：
- 每个对话会话分配唯一ID
- 使用sync.Map存储对话历史（key: sessionId, value: []Message）
- 限制上下文长度（最多10轮对话）
- 实现过期清理机制（30分钟无活动自动清理）

### 6. 卡片导出实现方案

**决策**：使用html2canvas库将卡片DOM转换为图片

**理由**：
- 实现简单，无需后端支持
- 支持自定义图片质量
- 可以导出为PNG或JPEG格式
- 支持下载和分享

**实现方式**：
- 使用html2canvas库捕获卡片DOM
- 设置合适的图片尺寸和质量
- 确保卡片样式完整，不被压缩
- 支持下载到本地或分享

## 数据模型设计

### 核心实体

#### 1. 对话消息 (ConversationMessage)

```typescript
interface ConversationMessage {
  id: string;                    // 消息ID
  type: 'text' | 'card' | 'image' | 'voice';  // 消息类型
  content: any;                  // 消息内容（根据类型不同）
  timestamp: string;             // 消息时间戳
  sender: 'user' | 'assistant';  // 发送者
  sessionId?: string;            // 对话会话ID
}
```

#### 2. 用户设置 (UserSettings)

```typescript
interface UserSettings {
  language: 'zh' | 'en';         // 语言设置
  grade?: string;                 // 年级设置（K1-K12）
  lastUpdated: string;            // 最后更新时间
}
```

#### 3. 意图识别结果 (IntentResult)

```go
type IntentResult struct {
    Intent      string                 `json:"intent"`      // 意图类型
    Confidence  float64                `json:"confidence"`  // 置信度
    Parameters  map[string]interface{} `json:"parameters"`  // 意图参数
}
```

#### 4. 对话会话 (ConversationSession)

```go
type ConversationSession struct {
    SessionId   string                `json:"sessionId"`
    Messages    []ConversationMessage `json:"messages"`
    CreatedAt   time.Time             `json:"createdAt"`
    LastActive  time.Time             `json:"lastActive"`
}
```

## API接口设计

### 新增接口

#### 1. 意图识别接口

**POST** `/api/conversation/intent`

**请求体**:
```json
{
  "text": "帮我生成小卡片",
  "sessionId": "session-uuid",
  "context": []
}
```

**响应**:
```json
{
  "intent": "generate_cards",
  "confidence": 0.95,
  "parameters": {}
}
```

#### 2. 对话接口

**POST** `/api/conversation/message`

**请求体**:
```json
{
  "sessionId": "session-uuid",
  "type": "text",
  "content": "这是什么植物？",
  "inputType": "text"  // text/voice/image
}
```

**响应**（流式）:
```
data: {"type":"text","content":"这是银杏..."}
data: {"type":"text","content":"银杏是..."}
```

#### 3. 语音识别接口

**POST** `/api/conversation/voice`

**请求体**:
```json
{
  "audio": "base64编码的音频数据",
  "sessionId": "session-uuid"
}
```

**响应**:
```json
{
  "text": "识别出的文本",
  "intent": "generate_cards",
  "confidence": 0.9
}
```

#### 4. 流式返回接口

**GET** `/api/conversation/stream?sessionId=xxx`

**响应**（SSE）:
```
event: message
data: {"type":"text","content":"这是..."}

event: card
data: {"type":"science","title":"...","content":{...}}

event: image
data: {"url":"...","alt":"..."}
```

## 实施步骤

### Phase 1: 基础功能实现（优先级P1）

1. **前端按钮功能实现**
   - 实现所有页面的按钮真实功能
   - 实现统一拍照入口
   - 实现相册选择和语音输入

2. **中英文切换功能**
   - 安装react-i18next
   - 创建国际化配置和翻译文件
   - 实现语言切换组件和设置页面

3. **对话式交互系统**
   - 创建对话消息列表组件
   - 实现消息输入组件（文本、语音、图片）
   - 创建对话页面，整合消息列表和输入

### Phase 2: AI能力集成（优先级P1）

4. **eino框架集成**
   - 安装eino框架依赖
   - 创建Agent系统和Graph图
   - 实现各个节点（图片识别、文本生成、图片生成、意图识别）

5. **意图识别功能**
   - 实现意图识别逻辑
   - 创建意图识别处理器和API接口
   - 集成到对话流程中

6. **流式返回支持**
   - 实现SSE处理器
   - 实现流式返回逻辑
   - 前端实现流式数据接收和渲染

### Phase 3: 增强功能（优先级P2）

7. **卡片收藏和导出**
   - 完善卡片收藏功能
   - 实现卡片导出功能
   - 优化卡片样式和大小

8. **数据存储增强**
   - 实现对话消息存储
   - 实现用户设置存储
   - 实现内存缓存存储

## 注意事项

1. **eino框架集成**：需要确认eino框架的具体使用方法和配置方式
2. **AI模型调用**：需要申请APP ID，配置AI模型服务
3. **流式返回**：确保SSE连接稳定，实现断线重连
4. **意图识别**：需要大量测试，确保识别准确率≥85%
5. **移动端优化**：确保所有交互在移动端流畅，支持touch事件
6. **卡片导出**：确保导出的图片清晰，符合设计规范
7. **性能优化**：按钮响应时间≤200ms，语言切换≤100ms

## 下一步

1. **确认eino框架集成方式**：查阅eino文档，确认与go-zero的集成方法
2. **申请AI模型APP ID**：申请图片识别、文本生成、图片生成模型的APP ID
3. **生成任务清单**：使用 `/speckit.tasks` 命令生成详细的任务清单
4. **开始实现**：按照任务清单开始编码，优先实现P1功能