# Implementation Plan: H5对话落地页

**Branch**: `007-conversation-landing-page` | **Date**: 2025-12-19 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/007-conversation-landing-page/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

构建一个web H5对话落地页，实现拍照识别后的意图识别、知识卡片生成和追问对话功能。后端基于Eino框架实现流式对话，支持20轮上下文窗口，根据学生年级生成适配内容。前端支持流式输出、打字机效果、图片loading占位，移动端优先设计，兼容PC端。

## Technical Context

**Language/Version**: 
- 后端: Go 1.21+ (go-zero框架)
- 前端: TypeScript 5.9+, React 19.2+, Tailwind CSS 4.1+

**Primary Dependencies**: 
- 后端: Eino框架 (github.com/cloudwego/eino), go-zero, eino-ext (Ark模型)
- 前端: React, React Router, Tailwind CSS, Axios

**Storage**: 
- 内存存储: 对话历史消息（最多20轮，临时存储）
- 本地存储: 前端localStorage保存对话历史

**Testing**: 
- 后端: go-zero测试框架
- 前端: React Testing Library (可选)

**Target Platform**: 
- Web应用: 现代浏览器（Chrome 90+, Safari 14+, Firefox 88+）
- 移动端优先: iOS Safari, Android Chrome
- PC端兼容: Chrome, Edge, Firefox

**Project Type**: Web application (frontend + backend)

**Performance Goals**: 
- 卡片生成响应时间: <3秒
- 流式回答启动时间: <1秒
- 打字机效果流畅度: 60fps
- 图片loading占位替换: <2秒

**Constraints**: 
- 上下文窗口限制: 最多20轮对话
- 移动端性能: 支持低端设备流畅运行
- 网络环境: 支持弱网环境下的流式传输
- 浏览器兼容: 必须支持SSE (Server-Sent Events)

**Scale/Scope**: 
- 单用户对话场景
- 支持并发用户对话（后端无状态设计）
- 前端单页面应用（SPA）

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**规范检查项**（基于 `.specify/memory/constitution.md`）：

- [x] **原则一：中文优先规范** - 所有文档和生成内容必须使用中文（除非技术限制）
- [x] **原则二：K12 教育游戏化设计规范** - 设计必须符合儿童友好性、游戏化元素、玩中学理念，支持探索世界、学习古诗文、学习英语
- [x] **原则三：可发布应用规范** - 实现必须达到生产级标准，遵循MVP优先原则，快速实现功能并确保完整可用性
- [x] **原则四：多语言和年级设置规范** - 支持中英文设置和K12年级设置，默认中文
- [x] **原则五：AI优先（模型优先）规范** - 所有AI功能通过后端统一调用，使用Eino框架实现，支持流式返回
- [x] **原则六：移动端优先规范** - 确保移动端交互完整性，统一拍照入口，支持随时随地探索
- [x] **原则七：用户体验流程规范** - 识别后直接跳转问答页，用户消息必须展示，消息卡片暂不显示图片
- [x] **原则八：对话Agent技术规范** - 对话Agent必须基于Eino Graph实现，支持联网获取信息、图文混排输出、SSE流式输出和打字机效果

**合规性说明**：所有设计均符合项目规范要求，无违反项。

## Project Structure

### Documentation (this feature)

```text
specs/007-conversation-landing-page/
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
│   ├── agent/
│   │   ├── agent.go                    # Agent系统（已有）
│   │   ├── graph.go                    # Graph编排（已有，需扩展）
│   │   └── nodes/
│   │       ├── image_recognition.go     # 图片识别节点（已有）
│   │       ├── text_generation.go       # 文本生成节点（已有，需扩展流式支持）
│   │       ├── image_generation.go      # 图片生成节点（已有）
│   │       ├── intent_recognition.go   # 意图识别节点（已有）
│   │       └── conversation_node.go    # 对话节点（新建，基于Eino Graph实现流式对话）
│   ├── handler/
│   │   ├── conversationhandler.go       # 对话处理器（已有，需扩展）
│   │   ├── streamhandler.go             # 流式返回处理器（已有，需接入Eino）
│   │   └── routes.go                    # 路由配置（已有）
│   ├── logic/
│   │   ├── conversationlogic.go         # 对话管理逻辑（已有，需扩展）
│   │   ├── streamlogic.go               # 流式返回逻辑（已有，需接入Eino）
│   │   └── intentlogic.go               # 意图识别逻辑（已有）
│   ├── storage/
│   │   └── memory.go                    # 内存缓存实现（已有，需扩展20轮限制）
│   └── types/
│       └── types.go                      # 类型定义（已有，需扩展）
├── api/
│   └── explore.api                      # API定义（需扩展流式对话接口）
└── go.mod                                # 依赖管理（已有Eino依赖）

frontend/
├── src/
│   ├── pages/
│   │   └── Result.tsx                   # 对话落地页（已有，需优化）
│   ├── components/
│   │   ├── conversation/
│   │   │   ├── ConversationList.tsx     # 对话列表组件（已有，需优化）
│   │   │   ├── ConversationMessage.tsx  # 消息组件（已有，需支持打字机效果）
│   │   │   ├── MessageInput.tsx         # 消息输入组件（已有）
│   │   │   └── TypingIndicator.tsx      # 打字机效果组件（新建）
│   │   └── common/
│   │       └── ImagePlaceholder.tsx     # 图片loading占位组件（新建）
│   ├── services/
│   │   ├── conversation.ts               # 对话服务（已有，需扩展流式支持）
│   │   ├── sse.ts                       # SSE服务（已有，需优化）
│   │   └── api.ts                       # API服务（已有）
│   ├── hooks/
│   │   ├── useStreamConversation.ts     # 流式对话Hook（新建）
│   │   └── useTypingEffect.ts          # 打字机效果Hook（新建）
│   ├── styles/
│   │   └── responsive.ts                # 响应式样式配置（新建，Tailwind多端适配）
│   └── types/
│       └── conversation.ts              # 对话类型定义（已有，需扩展）
├── tailwind.config.js                    # Tailwind配置（需扩展多端适配）
└── package.json                          # 依赖管理（已有）
```

**Structure Decision**: 采用Web应用架构（frontend + backend），后端基于go-zero和Eino框架，前端基于React和Tailwind CSS。对话功能通过新增conversation_node.go实现基于Eino Graph的流式对话，前端通过新增hooks和组件支持流式输出和打字机效果。

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

无违反项。

## Phase 0: Outline & Research ✅

**状态**: 已完成

**输出文件**: `research.md`

**研究内容**:
1. ✅ Eino框架流式输出实现方案
2. ✅ 上下文窗口管理（20轮对话历史）
3. ✅ 基于年级的prompt生成策略
4. ✅ 前端打字机效果实现
5. ✅ 图片loading占位实现
6. ✅ 移动端优先的响应式设计

**关键决策**:
- 使用Eino ChatModel.Stream接口实现流式输出
- 在内存存储中维护最近20轮对话，转换为Eino Message格式
- 根据用户年级（3-18岁）动态生成不同难度的prompt
- 使用React Hook实现逐字显示效果，支持流式数据更新
- 使用React组件实现图片loading占位，支持进度显示
- 使用Tailwind CSS实现移动端优先的响应式设计

## Phase 1: Design & Contracts ✅

**状态**: 已完成

**输出文件**:
- ✅ `data-model.md` - 数据模型定义
- ✅ `contracts/conversation-stream.api` - API合约定义
- ✅ `contracts/README.md` - API文档
- ✅ `quickstart.md` - 快速开始指南

**数据模型**:
- 对话会话（ConversationSession）：管理20轮上下文窗口
- 对话消息（ConversationMessage）：支持多种消息类型和流式更新
- 识别结果上下文（IdentificationContext）：作为对话初始上下文
- 知识卡片（KnowledgeCard）：在对话中展示
- 流式事件（StreamEvent）：SSE事件数据结构

**API合约**:
- `GET /api/conversation/stream` - SSE流式对话接口
- `POST /api/conversation/message` - 非流式对话接口（兼容性）

**Agent Context更新**: ✅ 已完成
