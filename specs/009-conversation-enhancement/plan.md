# Implementation Plan: 对话页面完善 - Agent模型流式返回

**Branch**: `009-conversation-enhancement` | **Date**: 2025-12-21 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/009-conversation-enhancement/spec.md`

**Note**: This plan focuses on implementing a unified streaming interface that supports text, voice, and image inputs, with configurable AI model usage via environment variables.

## Summary

实现对话页面的统一流式接口，支持文本输入、语音输入、图片输入三种输入方式。统一接口通过 `messageType` 字段明确指定输入类型（`"text"|"voice"|"image"`），根据输入类型自动处理（语音识别、图片上传等），然后通过Agent模型流式返回回答。系统支持通过 `.env` 配置文件中的 `USE_AI_MODEL` 字段控制是否使用AI模型调用，默认值为 `true`（使用AI模型）。当 `USE_AI_MODEL=true` 时，禁止使用Mock数据；当 `USE_AI_MODEL=false` 时，可以使用Mock数据作为降级方案（仅用于开发测试场景）。除了生成知识卡片使用另一个独立接口。

## Technical Context

**Language/Version**: 
- 后端: Go 1.21+ (go-zero框架)
- 前端: TypeScript 5.9+, React 19.2+

**Primary Dependencies**: 
- 后端: Eino框架 (github.com/cloudwego/eino), go-zero, eino-ext (Ark模型)
- 前端: React, React Router, Tailwind CSS, Axios, Server-Sent Events (SSE)

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
- 统一接口流式响应启动时间: ≤5秒（从用户发送消息到开始接收流式响应）
- 语音输入处理时间: <2秒（从接收语音数据到完成识别）
- 图片上传处理时间: <2秒（从接收图片数据到完成上传）
- 流式输出过程中，文本逐字显示的延迟: ≤100ms（每个字符到达后立即显示）

**Constraints**: 
- 必须实现统一流式接口，支持文本、语音、图片三种输入方式
- 统一接口必须通过 `messageType` 字段明确指定输入类型
- 必须支持通过 `.env` 配置文件中的 `USE_AI_MODEL` 字段控制是否使用AI模型调用，默认值为 `true`（使用AI模型）
- 当 `USE_AI_MODEL=true`（默认值）时，必须使用真实的Agent模型，禁止使用Mock数据
- 当 `USE_AI_MODEL=false` 时，可以使用Mock数据作为降级方案（仅用于开发测试场景）
- 必须支持多模态输入（文本、语音、图片）
- 必须维护对话上下文（最多20轮）
- 必须支持错误处理和降级机制（但必须记录错误日志）

**Scale/Scope**: 
- 单用户对话场景
- 支持并发用户对话（后端无状态设计）
- 前端单页面应用（SPA）

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**规范检查项**（基于项目规范）：

- [x] **原则一：中文优先规范** - 所有文档和生成内容必须使用中文（除非技术限制）
- [x] **原则二：K12 教育游戏化设计规范** - 设计必须符合儿童友好性、游戏化元素、玩中学理念
- [x] **原则三：可发布应用规范** - 实现必须达到生产级标准，遵循MVP优先原则
- [x] **原则四：多语言和年级设置规范** - 支持中英文设置和K12年级设置，默认中文
- [x] **原则五：AI优先（模型优先）规范** - 所有AI功能通过后端统一调用，使用Eino框架实现，支持流式返回
- [x] **原则六：移动端优先规范** - 确保移动端交互完整性，统一拍照入口，支持随时随地探索
- [x] **原则七：用户体验流程规范** - 识别后直接跳转问答页，用户消息必须展示
- [x] **原则八：对话Agent技术规范** - 对话Agent必须基于Eino Graph实现，支持联网获取信息、图文混排输出、SSE流式输出和打字机效果。**语音输入和图片上传必须支持Agent模型流式返回内容，通过配置字段 `USE_AI_MODEL` 控制是否使用AI模型调用，默认值为 `true`（使用AI模型）**

**合规性说明**：所有设计均符合项目规范要求，无违反项。特别符合原则八的要求：语音输入和图片上传必须支持Agent模型流式返回内容。通过配置字段 `USE_AI_MODEL` 提供灵活性，默认使用AI模型，同时支持开发测试场景使用Mock数据。

## Project Structure

### Documentation (this feature)

```text
specs/009-conversation-enhancement/
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
│   │   ├── graph.go                    # Graph编排（已有）
│   │   └── nodes/
│   │       └── conversation_node.go    # 对话节点（已有，需扩展支持多模态）
│   ├── config/
│   │   └── config.go                   # 配置定义（需添加UseAIModel字段）
│   ├── handler/
│   │   ├── streamhandler.go            # 统一流式接口处理器（已有，需扩展支持messageType）
│   │   └── routes.go                   # 路由配置（已有）
│   ├── logic/
│   │   ├── voicelogic.go               # 语音识别逻辑（已有，统一接口内部调用）
│   │   ├── uploadlogic.go             # 图片上传逻辑（已有，统一接口内部调用）
│   │   ├── streamlogic.go            # 流式返回逻辑（已有，需扩展支持messageType、多模态和配置字段）
│   │   └── conversationlogic.go       # 对话管理逻辑（已有）
│   ├── storage/
│   │   └── memory.go                   # 内存缓存实现（已有）
│   └── types/
│       └── types.go                    # 类型定义（需扩展统一流式请求类型，包含messageType字段）
├── api/
│   └── explore.api                     # API定义（需更新为统一流式接口）
└── go.mod                               # 依赖管理（已有Eino依赖）

frontend/
├── src/
│   ├── pages/
│   │   └── Result.tsx                  # 对话落地页（已有，需使用统一接口）
│   ├── components/
│   │   └── conversation/
│   │       ├── ConversationList.tsx    # 对话列表组件（已有）
│   │       ├── ConversationMessage.tsx  # 消息组件（已有）
│   │       └── MessageInput.tsx       # 消息输入组件（已有，需支持多模态）
│   ├── services/
│   │   ├── conversation.ts             # 对话服务（需更新为统一接口调用）
│   │   └── sse.ts                      # SSE服务（已有）
│   └── types/
│       └── conversation.ts             # 对话类型定义（需扩展统一请求类型）
└── package.json                        # 依赖管理（已有）
```

**Structure Decision**: 采用Web应用架构（frontend + backend），后端基于go-zero和Eino框架，前端基于React和Tailwind CSS。实现统一流式接口，通过 `messageType` 字段支持文本、语音、图片三种输入方式，一个接口搞定所有输入和输出。

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

无违反项。

## Phase 0: Outline & Research ✅

**状态**: 已完成  
**输出**: [research.md](./research.md)

**研究内容**:
1. ✅ 统一接口设计 - 实现一个统一的流式接口，支持文本、语音、图片三种输入方式
2. ✅ 输入类型识别方式 - 通过 `messageType` 字段明确指定输入类型（`"text"|"voice"|"image"`）
3. ✅ 统一接口内部处理流程 - 根据 `messageType` 自动调用语音识别或图片上传，然后调用Agent模型
4. ✅ Eino Graph是否支持多模态输入（文本+图片） - 支持，使用 `UserInputMultiContent` 和 `MessageInputPart`
5. ✅ 配置字段设计 - 通过 `.env` 配置文件中的 `USE_AI_MODEL` 字段控制是否使用AI模型调用，默认值为 `true`
6. ✅ 错误处理机制 - 当 `USE_AI_MODEL=true` 时，禁止降级到Mock数据；当 `USE_AI_MODEL=false` 时，可以使用Mock数据作为降级方案
7. ✅ 前端如何统一处理多模态输入的流式返回 - 前端统一调用一个接口，根据输入类型设置 `messageType` 字段

**关键决策**:
- **统一接口设计**：实现一个统一的流式接口 `POST /api/conversation/stream`，支持文本、语音、图片三种输入方式
- **输入类型识别**：通过 `messageType` 字段（必填）明确指定输入类型：`"text"`、`"voice"` 或 `"image"`
- **内部处理流程**：
  - `messageType: "text"` → 直接使用 `message` 字段调用Agent模型
  - `messageType: "voice"` → 先调用语音识别，将识别的文本调用Agent模型
  - `messageType: "image"` → 先上传图片获取URL，构建多模态消息调用Agent模型
- **多模态支持**：扩展 `ConversationNode.StreamConversation` 方法，支持图片URL参数，构建多模态消息
- **配置字段设计**：在 `backend/internal/config/config.go` 的 `AIConfig` 结构体中添加 `UseAIModel` 字段，从环境变量 `USE_AI_MODEL` 读取配置，默认值为 `true`（使用AI模型）
- **错误处理**：当 `UseAIModel=true`（默认值）时，Agent模型调用失败时记录详细错误日志，向用户发送错误事件，不允许降级到Mock数据；当 `UseAIModel=false` 时，可以使用Mock数据作为降级方案

## Phase 1: Design & Contracts ✅

**状态**: 已完成  
**输出**: 
- ✅ [data-model.md](./data-model.md) - 数据模型定义（扩展多模态消息类型）
- ✅ [contracts/stream-conversation-multimodal.api](./contracts/stream-conversation-multimodal.api) - API合约定义（多模态流式接口）
- ✅ [contracts/README.md](./contracts/README.md) - API文档
- ✅ [quickstart.md](./quickstart.md) - 快速开始指南

**数据模型**:
- ✅ 定义 `UnifiedStreamConversationRequest` 统一流式请求类型，包含 `messageType` 字段（必填）
- ✅ 根据 `messageType` 的值，请求包含对应的字段：`message`（text时）、`audio`（voice时）、`image`（image时）
- ✅ 扩展 `ConversationMessage` 支持多模态消息类型
- ✅ 扩展 `StreamEvent` 支持多模态流式事件

**API合约**:
- ✅ `POST /api/conversation/stream` - 统一流式对话接口，支持文本、语音、图片三种输入方式
- ✅ 请求格式：必须包含 `messageType` 字段（`"text"|"voice"|"image"`），根据 `messageType` 包含对应字段
- ✅ 响应格式：SSE流式响应，统一处理所有输入类型的输出

**设计决策**:
- ✅ 统一接口设计：一个接口支持所有输入类型，通过 `messageType` 字段区分
- ✅ 内部处理流程：统一接口根据 `messageType` 自动调用语音识别或图片上传，然后调用Agent模型
- ✅ 多模态支持：扩展 `ConversationNode.StreamConversation` 方法，支持图片URL参数，构建多模态消息
- ✅ 配置字段设计：在 `AIConfig` 结构体中添加 `UseAIModel` 字段，从环境变量 `USE_AI_MODEL` 读取配置，默认值为 `true`
- ✅ 错误处理：当 `UseAIModel=true`（默认值）时，禁止降级到Mock数据，记录详细错误日志；当 `UseAIModel=false` 时，可以使用Mock数据作为降级方案

**Agent Context更新**: ✅ 已完成（通过research.md和quickstart.md文档说明）

## Phase 2: Task Breakdown

**状态**: ⏳ 待执行  
**下一步**: 使用 `/speckit.tasks` 命令创建任务清单

### 预计任务范围

1. **后端任务**：
   - 在 `config.go` 的 `AIConfig` 结构体中添加 `UseAIModel` 字段，从环境变量 `USE_AI_MODEL` 读取配置，默认值为 `true`
   - 扩展 `streamlogic.go`，实现统一接口逻辑，根据 `messageType` 字段处理不同输入类型
   - 在统一接口中集成语音识别：当 `messageType` 为 `"voice"` 时，调用 `VoiceLogic.RecognizeVoice`
   - 在统一接口中集成图片上传：当 `messageType` 为 `"image"` 时，调用 `UploadLogic.Upload`
   - 在统一接口处理逻辑中，根据 `UseAIModel` 配置决定是否使用AI模型，当 `UseAIModel=true` 时禁止降级到Mock数据
   - 扩展 `conversation_node.go`，支持多模态输入（图片+文本）
   - 扩展类型定义，定义 `UnifiedStreamConversationRequest`，包含 `messageType` 字段
   - 更新 `StreamConversationHandler`，支持统一接口处理逻辑
   - 完善错误处理和日志记录

2. **前端任务**：
   - 更新对话服务，统一调用 `/api/conversation/stream` 接口
   - 文本输入：设置 `messageType: "text"`，发送 `message` 字段
   - 语音输入：设置 `messageType: "voice"`，发送 `audio` 字段
   - 图片输入：设置 `messageType: "image"`，发送 `image` 字段
   - 优化多模态消息的展示
   - 统一错误处理机制

3. **测试任务**：
   - 测试配置字段读取：验证 `USE_AI_MODEL` 配置字段正确读取和应用，默认值为 `true`
   - 测试统一接口的文本输入流式返回（`USE_AI_MODEL=true`）
   - 测试统一接口的语音输入流式返回（包含语音识别，`USE_AI_MODEL=true`）
   - 测试统一接口的图片输入流式返回（包含图片上传，`USE_AI_MODEL=true`）
   - 测试 `messageType` 字段验证和错误处理
   - 测试错误处理机制：当 `USE_AI_MODEL=true` 时，验证禁止降级到Mock数据；当 `USE_AI_MODEL=false` 时，验证可以使用Mock数据
   - 测试多模态输入的并发场景

