# 实现计划：优化识别流程与对话体验

**分支**: `004-optimize-identify-flow` | **日期**: 2025-12-18 | **规范**: [spec.md](./spec.md)
**输入**: 功能规范来自 `/specs/004-optimize-identify-flow/spec.md`

## Summary

优化识别流程与对话体验：将拍照识别后的流程从"识别→生成卡片→跳转结果页"改为"识别→跳转问答页→展示识别结果→通过对话触发卡片生成"。同时优化消息展示机制，确保用户问题立即显示，系统响应通过流式返回实时展示。注重并发性能优化和前后端数据一致性，保证用户体验流畅。

## Technical Context

**Language/Version**: 
- 前端：TypeScript (ES2020+), React 18
- 后端：Go 1.25.3 (darwin/arm64)

**Primary Dependencies**: 
- 前端：React 18, Vite, Tailwind CSS, React Router, Axios, Server-Sent Events (SSE)
- 后端：go-zero v1.9.3, eino (字节云原生AI框架)
- AI模型：图像识别模型、大语言模型（通过eino框架调用）

**Storage**: 
- 前端：浏览器本地存储（localStorage/IndexedDB）用于对话状态和消息历史
- 后端：内存缓存（会话状态、识别结果上下文）用于维护对话上下文

**Testing**: 
- 前端：Vitest + React Testing Library
- 后端：Go testing package + go-zero test tools

**Target Platform**: 
- Web H5（移动端优先，兼容PC端）
- 现代浏览器（Chrome 90+, Safari 14+, Firefox 88+）

**Project Type**: Web应用（前后端分离）

**Performance Goals**: 
- 识别后跳转问答页响应时间≤2秒（90%请求）
- 用户消息展示延迟≤0.5秒（本地渲染）
- 流式返回首字符延迟≤1秒（90%请求）
- 支持并发对话会话：200+（MVP版本）
- 前后端数据一致性：100%（所有消息正确同步）

**Constraints**: 
- 必须保证前后端接口数据格式一致
- 必须处理并发场景（多用户同时使用）
- 必须优化性能，确保用户体验流畅
- 必须保证状态传递和恢复的正确性

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
  
- [x] **原则二：K12 教育游戏化设计规范** - 设计必须符合儿童友好性、游戏化元素、玩中学理念
  - ✅ UI设计：问答页面保持简洁有趣，符合儿童认知
  - ✅ 交互设计：即时反馈，流畅的对话体验
  - ✅ 内容设计：识别结果展示清晰，对话引导自然
  
- [x] **原则三：可发布应用规范** - 实现必须达到生产级标准，可正常运行和发布
  - ✅ 错误处理：完善的错误提示和降级方案
  - ✅ 性能优化：响应时间满足要求（跳转≤2秒，消息展示≤0.5秒）
  - ✅ 并发处理：支持200+并发会话
  
- [x] **原则四：多语言和年级设置规范** - 支持中英文设置和K12年级设置，默认中文
  - ✅ 问答页面支持多语言切换
  - ✅ 识别结果和对话内容根据语言设置显示
  
- [x] **原则五：AI优先（模型优先）规范** - 所有AI功能通过后端统一调用，支持流式返回
  - ✅ 卡片生成通过后端API调用
  - ✅ 系统响应通过SSE流式返回
  - ✅ 前端实现流式数据接收和实时渲染
  
- [x] **原则六：移动端优先规范** - 确保移动端交互完整性，统一拍照入口，支持随时随地探索
  - ✅ 问答页面移动端适配
  - ✅ 消息输入和展示优化移动端体验
  
- [x] **原则七：用户体验流程规范** - 识别后直接跳转问答页，用户消息必须展示，消息卡片暂不显示图片
  - ✅ 识别后直接跳转问答页面（核心要求）
  - ✅ 识别结果首先展示在问答页面
  - ✅ 用户消息立即显示
  - ✅ 知识卡片只显示文本内容，不显示图片

**合规性说明**：所有设计均符合项目规范要求，无违反项。特别符合原则七的用户体验流程规范。

## Project Structure

### Documentation (this feature)

```text
specs/004-optimize-identify-flow/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
frontend/                    # 前端工程（React 18 + Vite + Tailwind CSS）
├── src/
│   ├── pages/
│   │   ├── Capture.tsx      # 拍照页面（需要修改：识别后跳转问答页）
│   │   └── Result.tsx       # 问答页面（需要修改：展示识别结果，支持对话触发卡片生成）
│   ├── components/
│   │   ├── conversation/
│   │   │   ├── ConversationList.tsx      # 对话列表（需要优化：立即显示用户消息）
│   │   │   ├── ConversationMessage.tsx    # 消息组件（需要修改：卡片不显示图片）
│   │   │   └── MessageInput.tsx           # 消息输入（已支持）
│   │   └── cards/
│   │       ├── ScienceCard.tsx            # 科学卡片（需要修改：不显示图片）
│   │       ├── PoetryCard.tsx             # 诗词卡片（需要修改：不显示图片）
│   │       └── EnglishCard.tsx            # 英语卡片（需要修改：不显示图片）
│   ├── services/
│   │   ├── api.ts                         # API服务（需要修改：识别后不调用生成卡片）
│   │   ├── conversation.ts                # 对话服务（需要优化：消息立即显示）
│   │   └── sse.ts                         # SSE流式服务（已支持）
│   └── types/
│       ├── api.ts                         # API类型（需要更新）
│       └── conversation.ts                # 对话类型（需要更新）

backend/                     # 后端工程（go-zero + eino）
├── internal/
│   ├── handler/
│   │   ├── identifyhandler.go            # 识别处理器（已存在）
│   │   └── conversationhandler.go        # 对话处理器（需要优化：支持识别结果上下文）
│   ├── logic/
│   │   ├── identifylogic.go              # 识别逻辑（已存在）
│   │   └── conversationlogic.go          # 对话逻辑（需要修改：支持识别结果上下文，优化并发）
│   ├── types/
│   │   └── types.go                      # 类型定义（需要更新：添加识别结果上下文）
│   └── agent/
│       └── graph.go                       # Agent图（需要修改：支持对话触发卡片生成）
└── etc/
    └── explore.yaml                       # 配置文件（已存在）
```

**Structure Decision**: 采用现有前后端分离架构，主要修改前端页面跳转逻辑、消息展示机制和后端对话处理逻辑，确保前后端数据格式一致。

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

无违反项，所有设计均符合项目规范。

## Phase 0: 研究与设计决策

已完成研究文档 `research.md`，解决了以下关键技术问题：

1. **识别后跳转流程优化**: 修改前端跳转逻辑，识别后直接跳转问答页
2. **用户消息立即展示**: 采用乐观更新策略，确保≤0.5秒显示
3. **知识卡片不显示图片**: 移除图片渲染逻辑，只显示文本
4. **识别结果上下文传递**: 使用React Router state传递，后端维护关联
5. **前后端数据一致性**: 使用TypeScript类型定义，确保格式一致
6. **并发性能优化**: 使用goroutine和内存缓存，支持200+并发会话
7. **流式返回优化**: 使用SSE实现流式返回，确保≤1秒首字符延迟

## Phase 1: 数据模型与API合约

已完成以下设计文档：

1. **数据模型** (`data-model.md`): 
   - 定义了识别结果、对话会话、对话消息、知识卡片等核心实体
   - 明确了前后端数据格式一致性要求
   - 定义了状态管理和数据流

2. **API合约** (`contracts/conversation.api`):
   - 定义了对话消息处理接口
   - 定义了流式返回（SSE）接口
   - 定义了意图识别接口
   - 确保前后端接口数据格式一致

3. **快速开始** (`quickstart.md`):
   - 提供了功能验证步骤
   - 定义了关键检查点
   - 提供了常见问题排查指南

## 规范检查（设计后重新评估）

*重新评估设计后的规范合规性*

**规范检查项**（基于 `.specify/memory/constitution.md`）：

- [x] **原则一：中文优先规范** - ✅ 所有设计文档使用中文，代码注释优先中文
- [x] **原则二：K12 教育游戏化设计规范** - ✅ 问答页面设计符合儿童友好性要求
- [x] **原则三：可发布应用规范** - ✅ 设计考虑了错误处理、性能优化、并发支持
- [x] **原则四：多语言和年级设置规范** - ✅ 问答页面支持多语言切换
- [x] **原则五：AI优先（模型优先）规范** - ✅ 卡片生成通过后端API，支持流式返回
- [x] **原则六：移动端优先规范** - ✅ 问答页面移动端适配，消息输入优化
- [x] **原则七：用户体验流程规范** - ✅ 识别后直接跳转问答页，用户消息立即显示，卡片不显示图片

**合规性说明**：所有设计均符合项目规范要求，无违反项。特别注重了用户体验流程优化和前后端数据一致性。
