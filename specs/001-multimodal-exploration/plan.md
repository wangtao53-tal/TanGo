# 实现计划：TanGo 多模态探索核心功能

**分支**: `001-multimodal-exploration` | **日期**: 2025-12-18 | **规范**: [spec.md](./spec.md)
**输入**: 功能规范来自 `/specs/001-multimodal-exploration/spec.md`

## Summary

实现 TanGo（小探号）核心功能：孩子通过拍照识别真实世界对象，系统使用AI实时生成三张知识卡片（科学认知卡、古诗词/人文卡、英语表达卡），实现"拍一得三"的核心亮点。采用前后端分离架构，前端使用 React 18 + Vite + Tailwind CSS（移动端优先），后端使用 go-zero 框架，AI部分通过字节的 eino 框架集成，使用 ReAct Agent 模式实现AI模型调用。

## Technical Context

**Language/Version**: 
- 前端：JavaScript/TypeScript (ES2020+)
- 后端：Go 1.25.3 (darwin/arm64)

**Primary Dependencies**: 
- 前端：React 18, Vite, Tailwind CSS, React Router, Axios
- 后端：go-zero v1.9.3, eino (字节云原生AI框架)
- AI模型：图像识别模型（待定APP ID）、大语言模型（GPT-4/Claude等，待定APP ID）

**Storage**: 
- 前端：浏览器本地存储（localStorage/IndexedDB）用于用户档案和收藏记录
- 后端：内存缓存（Redis，可选）用于分享链接临时存储；文件存储用于图片上传

**Testing**: 
- 前端：Vitest + React Testing Library
- 后端：Go testing package + go-zero test tools

**Target Platform**: 
- Web H5（移动端优先，兼容PC端）
- 现代浏览器（Chrome 90+, Safari 14+, Firefox 88+）

**Project Type**: Web应用（前后端分离）

**Performance Goals**: 
- 图像识别响应时间≤3秒（90%请求）
- 知识卡片生成响应时间≤5秒（90%请求）
- 前端首屏加载时间≤2秒
- 支持并发用户数：100+（MVP版本）

**Constraints**: 
- 黑客马拉松时间限制，优先实现核心功能
- 完全匿名使用，无需账户系统
- 数据仅本地存储，分享链接临时有效
- 必须保护儿童隐私数据

**Scale/Scope**: 
- MVP版本：支持80-100种常见对象识别
- 目标用户：K12学生（3-12岁）
- 预计并发：50-100用户（演示阶段）

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**规范检查项**（基于 `.specify/memory/constitution.md`）：

- [x] **原则一：中文优先规范** - 所有文档和生成内容必须使用中文（除非技术限制）
  - ✅ 前端界面文本全部中文
  - ✅ AI生成的知识卡片内容使用中文
  - ✅ 代码注释优先使用中文
  
- [x] **原则二：K12 教育游戏化设计规范** - 设计必须符合儿童友好性、游戏化元素、玩中学理念
  - ✅ UI设计：大按钮、清晰图标、简单操作流程（使用Tailwind CSS实现）
  - ✅ 交互设计：游戏化元素（收藏图鉴、分享成就感）
  - ✅ 内容设计：平衡趣味性和教育性，符合K12分级（通过AI prompt控制）
  
- [x] **原则三：可发布应用规范** - 实现必须达到生产级标准，可正常运行和发布
  - ✅ 错误处理：完善的错误提示和降级方案（API错误响应、前端错误处理）
  - ✅ 性能优化：响应时间满足要求（识别≤3秒，卡片生成≤5秒）
  - ✅ 安全性：数据传输加密（HTTPS），保护儿童隐私（本地存储、分享链接过期）

**合规性说明**：所有设计均符合项目规范要求，无违反项。

## Project Structure

### Documentation (this feature)

```text
specs/001-multimodal-exploration/
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
│   ├── components/         # React组件
│   │   ├── common/         # 通用组件（按钮、卡片等）
│   │   ├── camera/         # 拍照组件
│   │   ├── cards/          # 知识卡片展示组件
│   │   └── collection/     # 探索图鉴组件
│   ├── pages/              # 页面组件
│   │   ├── Home.tsx        # 首页（拍照入口）
│   │   ├── Result.tsx      # 识别结果页（展示三张卡片）
│   │   ├── Collection.tsx   # 我的探索图鉴
│   │   └── Share.tsx       # 分享页面（家长端）
│   ├── services/           # API服务
│   │   ├── api.ts          # API调用封装
│   │   └── storage.ts      # 本地存储服务
│   ├── hooks/              # React Hooks
│   ├── utils/              # 工具函数
│   ├── types/              # TypeScript类型定义
│   ├── App.tsx             # 根组件
│   └── main.tsx            # 入口文件
├── public/                 # 静态资源
├── index.html
├── vite.config.ts          # Vite配置
├── tailwind.config.js      # Tailwind CSS配置
├── tsconfig.json           # TypeScript配置
├── package.json
└── README.md

backend/                     # 后端工程（go-zero）
├── api/                     # API定义
│   └── explore.api         # 探索相关API定义
├── internal/
│   ├── handler/            # HTTP处理器
│   │   └── explorehandler.go
│   ├── logic/              # 业务逻辑
│   │   ├── explorationlogic.go
│   │   └── cardgenerationlogic.go
│   ├── svc/                # 服务上下文
│   │   └── servicecontext.go
│   ├── types/              # 类型定义
│   │   └── types.go
│   └── agent/              # AI Agent实现
│       ├── reactagent.go   # ReAct Agent核心
│       ├── visionagent.go  # 图像识别Agent
│       └── cardagent.go    # 知识卡片生成Agent
├── eino/                   # eino框架集成
│   ├── config.yaml         # eino配置
│   └── models/             # AI模型配置
│       ├── vision.yaml     # 图像识别模型配置
│       └── llm.yaml        # 大语言模型配置
├── go.mod
├── go.sum
├── etc/
│   └── explore.yaml       # 服务配置
└── README.md
```

**Structure Decision**: 采用前后端分离架构，前端和后端分别位于根目录下的 `frontend/` 和 `backend/` 文件夹。前端使用 React 18 + Vite + Tailwind CSS，移动端优先设计，完全基于提供的UI设计稿实现。后端使用 go-zero 框架，AI部分通过 eino 框架集成，使用 ReAct Agent 模式实现AI模型调用链。

**UI设计稿位置**: `stitch_ui/` 文件夹包含所有页面的HTML设计稿和截图，作为前端实现的参考标准。

## UI设计分析

基于提供的UI设计稿（`stitch_ui/` 文件夹），已分析所有页面的设计细节，以下是关键发现：

### 设计系统

**颜色主题**（基于设计稿提取）：
- **Primary Green**: `#4cdf20` / `#76FF7A` (科学认知卡主题色)
- **Sunny Orange**: `#FF9E64` (古诗词/人文卡主题色)
- **Sky Blue**: `#40C4FF` (英语表达卡主题色)
- **Warm Yellow**: `#FFE580` (强调色、按钮)
- **Peach Pink**: `#FFB7C5` (Little Star对话气泡)
- **Cloud White**: `#F8F8F8` (背景色)
- **Text Main**: `#1F2937` / `#2D3748` (主文本)
- **Text Sub**: `#6B7280` / `#718096` (次要文本)

**字体系统**：
- **Display Font**: Manrope, Fredoka, Space Grotesk（标题、按钮）
- **Body Font**: Noto Sans, Quicksand（正文）
- **Material Icons**: Material Symbols Outlined（图标系统）

**设计特点**：
- 大圆角设计（border-radius: 1rem - 3rem）
- 柔和阴影（soft shadows）
- 动画效果（float, pulse-glow, bounce等）
- 卡片式布局（card-based design）
- 游戏化元素（收藏、等级、成就等）

### 页面分析

#### 1. 首页（homepage_/_main_interface/）

**核心元素**：
- 大圆形拍照按钮（size-48 sm:size-64），带脉冲发光动画
- 语音触发按钮（浮动在拍照按钮右上角）
- 三个功能卡片展示区域（科学认知、人文素养、语言能力）
- Little Star对话气泡（底部固定，浮动动画）
- 背景装饰元素（模糊圆形，浮动动画）

**交互设计**：
- 拍照按钮：hover效果、点击缩放效果、脉冲发光动画
- 功能卡片：hover时上浮、边框颜色变化、阴影增强
- 响应式设计：移动端隐藏部分文字，PC端完整显示

#### 2. 拍照页面（capture_/_scan_interface/）

**核心元素**：
- 相机取景框（带圆角边框，金色主题 `#FFD700`）
- 扫描线动画（从顶部到底部，3秒循环）
- 快门按钮（大圆形，金色，带按压效果）
- 语音模式按钮（右侧悬浮）
- 相册按钮（底部左侧）
- 返回按钮（底部右侧）
- AI自动识别提示（"Little Star is identifying..."）

**交互设计**：
- 扫描线动画：模拟识别过程
- 快门按钮：active状态缩放效果
- 取景框：四个角的装饰元素

#### 3. 识别结果页面（recognition_result_page_1/2/3/）

**核心元素**：
- 三张知识卡片横向滑动展示（snap scroll）
- 每张卡片有独特的颜色主题：
  - 科学认知卡：绿色边框 `#76FF7A`
  - 古诗词/人文卡：橙色边框 `#FF9E64`
  - 英语表达卡：蓝色边框 `#40C4FF`
- 卡片内容：
  - 顶部图片区域（45%高度）
  - 底部内容区域（55%高度）
  - 收藏按钮（星星图标）
  - 播放/听按钮
- 底部固定栏：返回按钮、进度指示器、收藏所有卡片按钮

**交互设计**：
- 卡片轮播：snap scroll，左右箭头导航（PC端）
- 卡片hover：上浮效果、阴影增强
- 收藏按钮：点击变色、缩放动画

#### 4. 卡片详情页（science_card_detail_page/等）

**核心元素**：
- 左侧：卡片视觉展示（3:4比例，带光晕效果）
- 右侧：详细信息
  - 标题和标签
  - 音频播放按钮
  - 关键事实网格
  - 趣味知识点气泡
  - 收藏按钮
- 相关卡片推荐（横向滚动）

**交互设计**：
- 卡片hover：光晕效果
- 音频按钮：播放状态切换
- 收藏按钮：点击添加到收藏

#### 5. 收藏页面（favorites_page/）

**核心元素**：
- 侧边栏导航（PC端，移动端隐藏）
- 分类筛选按钮（All, Natural, Life, Humanities）
- 卡片网格布局（响应式：1列/2列/3列）
- 每个卡片：
  - 缩略图
  - 标题
  - 类别标签
  - 时间戳
  - "Re-explore"按钮
- Little Star鼓励消息（底部）

**交互设计**：
- 筛选按钮：激活状态高亮
- 卡片hover：上浮、阴影增强
- 响应式布局：移动端单列，平板2列，PC端3列

#### 6. 学习报告页面（learning_report_page/）

**核心元素**：
- 报告头部（标题、等级徽章）
- 统计卡片：
  - 总探索次数
  - 总收藏卡片数
  - 类别分布（饼图/柱状图）
- 最近收藏卡片列表
- 一键生成报告按钮

**交互设计**：
- 统计卡片：hover效果
- 数据可视化：简单的图表展示

### UI实现要求

1. **完全遵循设计稿**：
   - 所有颜色、字体、间距必须与设计稿一致
   - 所有动画效果必须实现
   - 所有交互状态必须实现

2. **组件化实现**：
   - 将设计稿中的HTML结构转换为React组件
   - 提取可复用的组件（按钮、卡片、对话框等）
   - 使用Tailwind CSS实现样式

3. **响应式设计**：
   - 移动端优先（Mobile First）
   - 使用Tailwind的响应式断点（sm, md, lg）
   - 确保在所有设备上体验良好

4. **动画实现**：
   - 使用CSS动画和Tailwind动画类
   - 关键动画：float、pulse-glow、bounce、scan-line
   - 过渡效果：hover、active、focus状态

5. **无障碍性**：
   - 所有交互元素必须有适当的aria标签
   - 确保键盘导航可用
   - 颜色对比度符合WCAG标准

## AI模型使用场景分析

### 场景1：图像识别（Vision Model）

**用途**：识别拍照对象，判断对象类别和名称

**输入**：
- 图像数据（base64或文件）
- 可选：对象类别提示（自然类/生活类/人文类）

**输出**：
- 对象名称（中文）
- 对象类别（自然类/生活类/人文类）
- 识别置信度（0-1）
- 相关标签/关键词

**推荐模型**：
- 通用图像识别模型（如：GPT-4V, Claude 3.5 Sonnet with Vision）
- 或专用图像分类模型（通过eino集成）

**调用时机**：用户拍照后立即调用

**性能要求**：响应时间≤3秒

### 场景2：知识卡片生成（LLM Model）

**用途**：根据识别结果和年龄，生成三张知识卡片内容

**输入**：
- 对象名称和类别
- 孩子年龄/年级（K1-K12）
- 对象相关关键词

**输出**：
- 科学认知卡内容（名称、解释、关键事实、趣味知识点）
- 古诗词/人文卡内容（相关诗词、白话解释、情境联想）
- 英语表达卡内容（核心单词、口语表达）

**推荐模型**：
- GPT-4 或 Claude 3.5 Sonnet（支持长文本生成，内容质量高）
- 通过eino框架统一调用

**调用时机**：图像识别完成后立即调用

**性能要求**：响应时间≤5秒（包含图像识别总时间）

**Prompt设计**：
- 需要针对不同年龄设计不同的prompt模板
- 确保生成内容符合K12分级要求
- 控制输出格式（JSON结构化输出）

## AI Agent框架建议

### 推荐方案：ReAct Agent（单Agent模式）

**理由**：
1. **简单直接**：我们的场景主要是顺序调用两个AI模型（图像识别 → 知识卡片生成），不需要复杂的多Agent协作
2. **易于实现**：ReAct Agent模式适合"推理-行动-观察"的循环，我们的流程是线性的
3. **快速开发**：适合黑客马拉松快速实现，降低复杂度
4. **易于调试**：单Agent模式更容易追踪和调试

**架构设计**：

```
用户拍照
  ↓
ReAct Agent启动
  ↓
Action 1: 调用Vision Model（图像识别）
  ↓
Observation: 获取识别结果（对象名称、类别、置信度）
  ↓
Reasoning: 判断识别是否成功，是否需要重试
  ↓
Action 2: 调用LLM Model（知识卡片生成）
  ↓
Observation: 获取三张卡片内容
  ↓
Reasoning: 验证内容质量，是否符合年龄要求
  ↓
Final Action: 返回结果给前端
```

**实现方式**：
- 使用eino框架封装AI模型调用
- 在go-zero的logic层实现ReAct Agent逻辑
- Agent负责协调Vision Model和LLM Model的调用顺序
- 实现简单的错误重试机制

### 不推荐：Multi-Agent方案

**原因**：
1. **过度设计**：我们的场景不需要多个Agent协作
2. **复杂度高**：需要Agent间通信、任务分配等机制
3. **开发时间长**：不适合黑客马拉松快速实现
4. **维护成本高**：多Agent系统调试困难

**如果未来需要扩展**：
- 可以考虑将"图像识别"和"知识卡片生成"拆分为两个Agent
- 但MVP阶段建议使用单Agent模式

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

无违反项，所有设计均符合项目规范。

## Phase 0: Research & Technology Decisions

### 技术选型研究

#### 1. eino框架集成研究

**决策**：使用字节的eino框架统一管理AI模型调用

**理由**：
- eino是字节开源的云原生AI框架，支持多种AI模型
- 提供统一的模型调用接口，简化集成
- 支持模型配置管理和监控

**待确认**：
- eino与go-zero的集成方式
- 支持的AI模型类型和配置方法
- APP ID申请和配置流程

**替代方案**：直接调用AI模型API（如OpenAI API、Claude API），但需要自行管理配置和错误处理

#### 2. 图像识别模型选择

**决策**：待定，需要根据APP ID申请情况确定

**候选方案**：
- GPT-4V（如果可用）：通用性强，识别准确率高
- Claude 3.5 Sonnet with Vision：识别准确率高，支持中文
- 专用图像分类模型：通过eino集成

**选择标准**：
- 识别准确率≥90%（目标95%+）
- 支持中文输出
- 响应时间≤3秒
- 支持80-100种常见对象识别

#### 3. 知识卡片生成模型选择

**决策**：使用GPT-4或Claude 3.5 Sonnet

**理由**：
- 内容生成质量高
- 支持长文本生成
- 支持结构化输出（JSON格式）
- 支持中文内容生成

**待确认**：
- 具体使用哪个模型（根据APP ID申请情况）
- Prompt设计优化
- 内容质量控制机制

#### 4. 前端本地存储方案

**决策**：使用IndexedDB + localStorage混合方案

**理由**：
- IndexedDB：存储大量探索记录和知识卡片（支持结构化数据）
- localStorage：存储用户设置（年龄、年级等简单配置）

**数据结构**：
- IndexedDB：探索记录表、知识卡片表
- localStorage：userProfile（年龄、年级）

#### 5. 分享链接实现方案

**决策**：使用临时分享链接，数据存储在内存或Redis

**实现方式**：
- 前端生成分享链接（包含唯一ID）
- 后端创建临时存储（内存或Redis，TTL 7天）
- 家长端通过链接ID获取数据
- 无需账户系统，简化实现

**数据结构**：
```json
{
  "shareId": "uuid",
  "explorationRecords": [...],
  "collectedCards": [...],
  "createdAt": "timestamp",
  "expiresAt": "timestamp"
}
```

## Phase 1: Data Model & API Contracts

### 数据模型设计

详见 `data-model.md`（将在Phase 1生成）

**核心实体**：
1. **UserProfile**（本地存储）
2. **ExplorationRecord**（本地存储 + 分享时上传）
3. **KnowledgeCard**（本地存储 + 分享时上传）
4. **ShareLink**（后端临时存储）

### API接口设计

详见 `contracts/` 目录（将在Phase 1生成）

**核心接口**：
1. `POST /api/explore/identify` - 图像识别
2. `POST /api/explore/generate-cards` - 生成知识卡片
3. `POST /api/share/create` - 创建分享链接
4. `GET /api/share/:shareId` - 获取分享数据
5. `POST /api/share/report` - 生成学习报告

## Next Steps

1. **Phase 0完成**：完成技术选型研究，确定AI模型和eino集成方式
2. **Phase 1完成**：完成数据模型设计和API接口定义
3. **Phase 2**：使用 `/speckit.tasks` 命令生成任务清单
4. **开始实现**：按照任务清单开始编码

## Notes

- 所有AI模型交互通过eino框架实现，需要申请APP ID
- **UI设计稿已完整提供**：`stitch_ui/` 文件夹包含所有页面的HTML和截图
- 前端实现必须完全遵循设计稿的样式、颜色、动画和交互
- 重点关注核心亮点"拍一得三"功能的实现和演示效果
- 优先保证核心功能可用，辅助功能可简化实现

## UI设计整合完成

✅ **已完成UI设计分析**：
- 已分析所有页面设计稿（首页、拍照、结果、详情、收藏、报告）
- 已提取完整的颜色系统、字体系统、动画效果
- 已识别所有可复用组件和交互模式
- 已更新项目结构，包含详细的组件和页面列表
- 已更新research.md，添加UI实现方案
- 已更新quickstart.md，添加UI开发指南

**下一步**：使用 `/speckit.tasks` 命令生成详细的任务清单，开始实现阶段。
