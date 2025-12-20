# 实现计划：对话体验与性能优化

**分支**: `008-conversation-optimization` | **日期**: 2025-12-20 | **规范**: [spec.md](./spec.md)
**输入**: 功能规范来自 `/specs/008-conversation-optimization/spec.md`

## Summary

优化对话体验与性能：修复流式消息实时渲染问题，实现Markdown格式支持，将知识卡片生成接口从40秒优化到5秒内，并添加知识卡片文本转语音功能。注重并发性能优化和前后端数据一致性，确保用户体验流畅，保证前后端正确执行。

## Technical Context

**Language/Version**: 
- 前端：TypeScript (ES2020+), React 19.2.0
- 后端：Go 1.21.4

**Primary Dependencies**: 
- 前端：React 19, Vite, Tailwind CSS, React Router, Axios, Server-Sent Events (SSE)
- 后端：go-zero v1.9.3, eino v0.7.11 (字节云原生AI框架)
- AI模型：大语言模型（通过eino框架调用Ark模型）
- Markdown渲染：react-markdown（需要新增）
- 文本转语音：Web Speech API 或第三方TTS服务（需要评估）

**Storage**: 
- 前端：浏览器本地存储（localStorage/IndexedDB）用于对话状态和消息历史
- 后端：内存缓存（会话状态、对话上下文）用于维护对话状态

**Testing**: 
- 前端：Vitest + React Testing Library
- 后端：Go testing package + go-zero test tools

**Target Platform**: 
- Web H5（移动端优先，兼容PC端）
- 现代浏览器（Chrome 90+, Safari 14+, Firefox 88+）

**Project Type**: Web应用（前后端分离）

**Performance Goals**: 
- 流式消息实时渲染：接收到第一个数据片段后100毫秒内更新UI
- 知识卡片生成响应时间：从40秒优化到5秒内（95%请求）
- Markdown渲染同步延迟：不超过200毫秒
- 支持并发对话会话：200+（保持现有能力）
- 前后端数据一致性：100%（所有消息正确同步）
- 文本转语音启动时间：≤1秒（90%请求）

**Constraints**: 
- 必须保证前后端接口数据格式一致
- 必须处理并发场景（多用户同时使用）
- 必须优化性能，确保用户体验流畅
- 必须保证流式消息实时渲染，不能延迟到最后统一渲染
- 必须支持Markdown格式实时渲染，与流式输出同步
- 必须保证性能优化后内容质量和完整性
- 必须支持文本转语音的多语言和播放控制

**Scale/Scope**: 
- MVP版本：支持200+并发对话会话
- 目标用户：K12学生（3-12岁）
- 预计并发：50-200用户（演示阶段）

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**规范检查项**（基于 `.specify/memory/constitution.md`）：

- [x] **原则一：中文优先规范** - 所有文档和生成内容必须使用中文（除非技术限制）
  - ✅ 前端界面文本全部中文
  - ✅ AI生成的对话内容使用中文
  - ✅ 代码注释优先使用中文
  
- [x] **原则二：K12 教育游戏化设计规范** - 设计必须符合儿童友好性、游戏化元素、玩中学理念，知识卡片支持文本转语音
  - ✅ UI设计：对话页面保持简洁有趣，符合儿童认知
  - ✅ 交互设计：实时反馈，流畅的对话体验
  - ✅ 文本转语音：支持多语言，语音清晰自然，适合儿童听力
  - ✅ 播放控制：支持播放、暂停、停止，操作简单直观
  
- [x] **原则三：可发布应用规范** - 实现必须达到生产级标准，关键接口响应时间≤5秒，流式消息实时渲染
  - ✅ 性能优化：知识卡片生成从40秒优化到5秒内
  - ✅ 实时渲染：流式消息必须实时渲染，不能延迟到最后统一渲染
  - ✅ 错误处理：完善的错误提示和降级方案
  - ✅ 并发处理：支持200+并发会话
  
- [x] **原则四：多语言和年级设置规范** - 支持中英文设置和K12年级设置，默认中文
  - ✅ 对话页面支持多语言切换
  - ✅ 文本转语音支持中英文等多种语言
  
- [x] **原则五：AI优先（模型优先）规范** - 所有AI功能通过后端统一调用，使用Eino框架实现，支持流式返回
  - ✅ 知识卡片生成通过后端API调用
  - ✅ 系统响应通过SSE流式返回
  - ✅ 前端实现流式数据接收和实时渲染
  
- [x] **原则六：移动端优先规范** - 确保移动端交互完整性，统一拍照入口，支持随时随地探索
  - ✅ 对话页面移动端适配
  - ✅ 消息输入和展示优化移动端体验
  - ✅ 文本转语音按钮移动端友好
  
- [x] **原则七：用户体验流程规范** - 识别后直接跳转问答页，用户消息必须展示，消息卡片暂不显示图片
  - ✅ 用户消息立即显示
  - ✅ 系统响应通过流式返回实时展示
  
- [x] **原则八：对话Agent技术规范** - 对话Agent必须基于Eino Graph实现，支持SSE流式输出、打字机效果、实时渲染和Markdown格式支持
  - ✅ 流式消息实时渲染：前端接收到SSE数据后立即更新UI
  - ✅ Markdown格式支持：流式消息内容支持Markdown格式渲染
  - ✅ Markdown渲染与流式输出同步：实时更新渲染结果
  - ✅ 打字机效果：文本内容逐字显示，与Markdown渲染协调工作

**合规性说明**：所有设计均符合项目规范要求，无违反项。特别符合原则三的性能要求和原则八的流式消息实时渲染、Markdown格式支持要求。

## Project Structure

### Documentation (this feature)

```text
specs/008-conversation-optimization/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
backend/
├── internal/
│   ├── handler/
│   │   ├── conversationhandler.go    # 对话接口处理器
│   │   ├── generatecardshandler.go   # 知识卡片生成接口处理器（需要优化）
│   │   └── streamhandler.go          # 流式输出处理器
│   ├── logic/
│   │   ├── conversationlogic.go      # 对话业务逻辑
│   │   ├── generatecardslogic.go    # 知识卡片生成逻辑（需要优化）
│   │   └── streamlogic.go           # 流式输出逻辑（需要优化实时渲染）
│   ├── agent/
│   │   ├── graph.go                 # Eino Graph实现（需要优化性能）
│   │   └── nodes/
│   │       └── text_generation.go   # 文本生成节点（需要优化）
│   └── types/
│       └── types.go                 # 类型定义（需要扩展）
└── api/
    └── explore.api                  # API定义（需要更新）

frontend/
├── src/
│   ├── components/
│   │   └── conversation/
│   │       ├── ConversationMessage.tsx    # 消息组件（需要支持Markdown渲染）
│   │       └── ConversationList.tsx       # 消息列表组件
│   ├── pages/
│   │   └── Result.tsx                    # 对话页面（需要优化流式渲染）
│   ├── hooks/
│   │   ├── useStreamConversation.ts      # 流式对话Hook（需要优化实时渲染）
│   │   └── useTypingEffect.ts           # 打字机效果Hook（需要与Markdown协调）
│   ├── services/
│   │   ├── sse.ts                       # SSE连接服务（需要优化）
│   │   ├── sse-post.ts                  # POST+SSE连接服务（需要优化）
│   │   └── conversation.ts              # 对话服务
│   ├── components/
│   │   └── cards/
│   │       ├── ScienceCard.tsx          # 科学卡片组件（需要添加文本转语音）
│   │       ├── PoetryCard.tsx           # 诗词卡片组件（需要添加文本转语音）
│   │       └── EnglishCard.tsx          # 英语卡片组件（需要添加文本转语音）
│   └── types/
│       ├── conversation.ts              # 对话类型定义
│       └── api.ts                       # API类型定义
└── package.json                         # 需要添加react-markdown依赖
```

**Structure Decision**: Web应用（前后端分离），前端使用React + TypeScript，后端使用Go + go-zero + eino框架。前后端通过REST API和SSE进行通信，确保数据格式一致。

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

无违反项，所有设计均符合项目规范。

---

## Phase 0: 研究完成

**状态**: ✅ 完成  
**输出**: [research.md](./research.md)

### 研究成果

1. **流式消息实时渲染**：使用增量更新机制，确保每次接收到数据后立即更新UI
2. **Markdown格式支持**：使用`react-markdown`库，与流式输出协调工作
3. **性能优化**：多策略组合，重点优化AI模型调用和实现流式返回
4. **文本转语音**：使用Web Speech API，浏览器原生方案

### 关键技术决策

- 流式消息实时渲染：使用`flushSync`强制同步更新
- Markdown渲染：使用`react-markdown`库
- 性能优化：并行生成+流式返回+超时控制
- 文本转语音：Web Speech API

---

## Phase 1: 设计完成

**状态**: ✅ 完成  
**输出**: 
- [data-model.md](./data-model.md)
- [contracts/README.md](./contracts/README.md)
- [quickstart.md](./quickstart.md)

### 设计成果

1. **数据模型**：扩展了流式消息、知识卡片等实体，新增文本转语音状态
2. **API契约**：定义了流式消息实时渲染、Markdown支持、知识卡片流式返回的API格式
3. **前后端一致性**：确保类型定义和字段映射一致

### 关键设计决策

- 所有新增字段为可选，确保向后兼容
- 支持同步和流式两种返回模式
- 前后端使用统一的类型定义

---

## Phase 2: 任务分解

**状态**: ⏳ 待执行  
**下一步**: 使用 `/speckit.tasks` 命令创建任务清单

### 预计任务范围

1. **前端任务**：
   - 优化流式消息实时渲染
   - 集成Markdown渲染
   - 实现文本转语音功能
   - 优化性能

2. **后端任务**：
   - 优化知识卡片生成性能
   - 实现流式返回
   - 扩展类型定义

3. **测试任务**：
   - 性能测试
   - 用户体验测试
   - 数据一致性测试
