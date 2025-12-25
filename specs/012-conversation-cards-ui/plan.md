# Implementation Plan: 对话页三张卡片生成页面美化

**Branch**: `dev-mvp-20251218` | **Date**: 2025-12-20 | **Spec**: [spec.md](./spec.md)

**Note**: MVP版本阶段，所有开发工作统一在 `dev-mvp-20251218` 分支进行，不采用一个功能一个分支的策略。

## Summary

优化对话页面中三张知识卡片的展示和交互体验，在保持页面整体布局不变的前提下，实现：
1. **卡片内容自适应显示**：卡片内容在固定比例（16:9）内完整显示，避免内部滚动
2. **渐进式卡片展示**：先显示第一张科学认知卡，移动端支持滑动切换，PC端可一并展开
3. **儿童友好视觉设计**：优化字体样式（圆润字体、合适字号行高）、色彩搭配（活泼有趣）、交互反馈（及时明确）
4. **卡片收藏和保存功能**：支持点击收藏、长按保存到本地

技术方案：创建卡片容器组件实现滑动切换，优化现有卡片组件的视觉样式和交互体验，使用CSS动画和响应式设计确保移动端优先。

## Technical Context

**Language/Version**: TypeScript 5.x, React 18.x  
**Primary Dependencies**: React, React Router, Tailwind CSS, html2canvas (用于保存图片)  
**Storage**: IndexedDB (通过现有的cardStorage服务)  
**Testing**: Jest + React Testing Library (前端单元测试)  
**Target Platform**: Web (移动端优先，支持PC端)  
**Project Type**: Web application (frontend)  
**Performance Goals**: 卡片滑动响应时间≤100ms，卡片渲染时间≤500ms，交互反馈响应时间≤200ms  
**Constraints**: 页面整体布局结构不变，仅优化卡片组件内部；支持320px-1920px宽度设备；字体可读性符合8-12岁儿童  
**Scale/Scope**: 单页面优化，3个卡片组件，1个卡片容器组件

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**规范检查项**（基于 `.specify/memory/constitution.md`）：

- [x] **原则一：中文优先规范** - 所有文档和生成内容必须使用中文（除非技术限制）
- [x] **原则二：K12 教育游戏化设计规范** - 设计必须符合儿童友好性、游戏化元素、玩中学理念，支持探索世界、学习古诗文、学习英语，知识卡片支持文本转语音
- [x] **原则三：可发布应用规范** - 实现必须达到生产级标准，遵循MVP优先原则，关键接口响应时间≤5秒，流式消息实时渲染
- [x] **原则四：多语言和年级设置规范** - 前端项目中文优先，所有页面默认显示中文，中文是主要语言，支持中英文设置和K12年级设置
- [x] **原则五：AI优先（模型优先）规范** - 模型调用优先，Agent eino框架优先，对话页面必须使用真实模型，Mock数据仅允许用于开发/测试环境，生产环境禁止使用Mock数据
- [x] **原则六：移动端优先规范** - 确保移动端交互完整性，统一拍照入口，支持随时随地探索
- [x] **原则七：用户体验流程规范** - 识别后直接跳转问答页，用户消息必须展示，消息卡片暂不显示图片
- [x] **原则八：对话Agent技术规范** - 对话Agent必须基于Eino Graph实现，支持联网获取信息、图文混排输出、SSE流式输出、打字机效果、实时渲染和Markdown格式支持，语音输入和图片上传必须支持Agent模型流式返回，禁止使用Mock数据

**合规性说明**：所有设计均符合项目规范要求，无违反项。特别符合原则二的K12教育游戏化设计规范（儿童友好性）和原则六的移动端优先规范。

## Project Structure

### Documentation (this feature)

```text
specs/012-conversation-cards-ui/
├── plan.md              # This file (/speckit.plan command output)
├── spec.md              # Feature specification
├── checklists/
│   └── requirements.md  # Quality checklist
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
frontend/
├── src/
│   ├── components/
│   │   ├── cards/
│   │   │   ├── ScienceCard.tsx          # 科学认知卡（需要优化：字体、色彩、交互）
│   │   │   ├── PoetryCard.tsx            # 古诗词卡（需要优化：字体、色彩、交互）
│   │   │   ├── EnglishCard.tsx          # 英语学习卡（需要优化：字体、色彩、交互）
│   │   │   └── CardCarousel.tsx          # 新建：卡片轮播容器组件（滑动切换、渐进式展示）
│   │   └── conversation/
│   │       ├── ConversationList.tsx     # 对话列表（保持不变）
│   │       └── ConversationMessage.tsx   # 消息组件（需要修改：卡片消息使用CardCarousel）
│   ├── styles/
│   │   ├── theme.ts                      # 主题配置（需要扩展：儿童友好字体、色彩配置）
│   │   └── cardStyles.ts                 # 新建：卡片专用样式配置
│   ├── hooks/
│   │   └── useCardSwipe.ts               # 新建：卡片滑动Hook（触摸事件处理）
│   │   └── useCardExport.ts              # 新建：卡片导出Hook（保存到本地）
│   ├── utils/
│   │   └── cardExport.ts                 # 新建：卡片导出工具（html2canvas封装）
│   └── pages/
│       └── Result.tsx                    # 结果页面（保持不变，仅卡片展示优化）
```

**Structure Decision**: 采用组件化设计，创建独立的CardCarousel组件处理卡片展示逻辑，保持现有卡片组件（ScienceCard、PoetryCard、EnglishCard）的独立性，通过样式和Hook优化视觉和交互体验。页面整体布局（Result.tsx、ConversationList.tsx）保持不变。

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

无违反项，无需填写。

## Phase 0: Research & Analysis

### 现有实现分析

1. **卡片组件现状**：
   - `ScienceCard.tsx`、`PoetryCard.tsx`、`EnglishCard.tsx` 已实现基础功能
   - 支持文本转语音、收藏功能
   - 使用Tailwind CSS样式，有基础的主题色彩（science-green、sunny-orange、sky-blue）
   - 卡片内容区域有 `overflow-y-auto`，可能导致滚动

2. **卡片展示方式**：
   - 当前在 `ConversationMessage.tsx` 中直接渲染卡片组件
   - 三张卡片通过流式API依次生成，每生成一张立即显示
   - 卡片在对话列表中垂直排列，可能导致页面过长

3. **字体和样式系统**：
   - 使用 `index.css` 中的主题配置
   - 字体：`--font-display: Manrope, "Noto Sans SC", sans-serif`
   - 字体：`--font-body: "Noto Sans", Quicksand, sans-serif`
   - 需要引入儿童友好字体（如圆体、幼圆）

### 技术选型

1. **滑动切换实现**：
   - 方案A：使用 `react-swipeable` 库（推荐）
   - 方案B：原生触摸事件处理（更灵活，但需要更多代码）
   - **选择**：方案A，使用 `react-swipeable` 库，简化实现，确保跨平台兼容性

2. **卡片导出功能**：
   - 方案A：使用 `html2canvas` 库（推荐）
   - 方案B：使用 `dom-to-image` 库
   - **选择**：方案A，`html2canvas` 更成熟，支持更好的样式渲染

3. **儿童友好字体**：
   - 方案A：使用Google Fonts的儿童字体（如 `Comfortaa`、`Nunito`）
   - 方案B：使用系统字体（如 `PingFang SC` 的圆体、`Microsoft YaHei UI`）
   - 方案C：引入中文字体文件（如 `幼圆`、`圆体`）
   - **选择**：方案A + 方案B组合，优先使用Google Fonts的 `Comfortaa` 或 `Nunito`，中文字体使用系统字体回退

4. **固定比例显示**：
   - 方案A：使用CSS `aspect-ratio` 属性（推荐）
   - 方案B：使用JavaScript计算高度
   - **选择**：方案A，CSS原生支持，性能更好

### 依赖项

**新增依赖**：
- `react-swipeable`: ^9.0.0 (卡片滑动切换)
- `html2canvas`: ^1.4.1 (卡片导出为图片)

**现有依赖**（无需新增）：
- React 18.x
- Tailwind CSS
- React Router

## Phase 1: Design & Architecture

### 组件设计

#### 1. CardCarousel 组件

**职责**：管理三张卡片的展示和切换逻辑

**Props**：
```typescript
interface CardCarouselProps {
  cards: KnowledgeCard[];  // 三张卡片数据（可能部分未生成）
  currentIndex: number;    // 当前显示的卡片索引
  onIndexChange: (index: number) => void;  // 索引变化回调
  onCollect?: (cardId: string) => void;   // 收藏回调
  onExport?: (cardId: string) => void;     // 导出回调
  isMobile: boolean;       // 是否为移动端
}
```

**功能**：
- 渐进式展示：只显示已生成的卡片，第一张优先显示
- 滑动切换：移动端支持左右滑动，PC端支持点击切换或键盘导航
- 固定比例：使用 `aspect-ratio: 16/9` 确保卡片在固定比例内显示
- 响应式：移动端单张显示，PC端可三张并排（可选）

**状态管理**：
- `currentIndex`: 当前显示的卡片索引
- `isTransitioning`: 是否正在切换动画
- `cardsReady`: 已生成的卡片列表

#### 2. 卡片组件优化

**ScienceCard、PoetryCard、EnglishCard**：
- 移除内容区域的 `overflow-y-auto`
- 添加固定比例容器：`aspect-ratio: 16/9`
- 优化字体样式：使用儿童友好字体，字号≥14px，行高≥1.6倍
- 优化色彩：保持主题色，增加对比度，确保可读性
- 优化交互反馈：按钮点击动画、收藏状态变化动画

#### 3. 样式系统扩展

**cardStyles.ts**：
```typescript
export const cardStyles = {
  // 儿童友好字体配置
  fonts: {
    childFriendly: {
      chinese: '"Comfortaa", "PingFang SC", "Microsoft YaHei UI", sans-serif',
      english: '"Comfortaa", "Nunito", sans-serif',
    },
    sizes: {
      title: 'clamp(18px, 4vw, 24px)',      // 标题：18-24px
      body: 'clamp(14px, 3vw, 16px)',      // 正文：14-16px
      small: 'clamp(12px, 2.5vw, 14px)',   // 小字：12-14px
    },
    lineHeight: {
      title: 1.4,
      body: 1.6,
      small: 1.5,
    },
  },
  // 固定比例配置
  aspectRatio: '16/9',
  // 色彩配置（增强对比度）
  colors: {
    // 保持现有主题色，但增强对比度
    scienceGreen: '#76FF7A',
    sunnyOrange: '#FF9E64',
    skyBlue: '#40C4FF',
    // 文本颜色（确保对比度≥4.5:1）
    textDark: '#1F2937',
    textLight: '#FFFFFF',
  },
};
```

### 数据流设计

1. **卡片生成流程**（保持不变）：
   - `Result.tsx` 调用 `generateCardsStream` API
   - 每生成一张卡片，添加到 `messages` 状态
   - `ConversationMessage` 检测到卡片消息，使用 `CardCarousel` 渲染

2. **卡片展示流程**（新增）：
   - `CardCarousel` 接收卡片数组，按生成顺序排序
   - 初始状态：`currentIndex = 0`（显示第一张）
   - 用户滑动：更新 `currentIndex`，触发切换动画
   - PC端：检测屏幕宽度，决定单张显示或三张并排

3. **收藏和导出流程**：
   - 收藏：点击收藏按钮 → `onCollect` 回调 → `Result.tsx` 的 `handleCollect`
   - 导出：长按卡片 → `useCardExport` Hook → `html2canvas` 转换 → 下载图片

### 交互设计

1. **移动端滑动**：
   - 左滑：切换到下一张（文言文 → 英语学习）
   - 右滑：切换到上一张（英语学习 → 文言文）
   - 滑动阈值：30% 屏幕宽度
   - 动画：CSS `transform: translateX()` + `transition`

2. **PC端切换**：
   - 方案A：三张并排显示（屏幕宽度≥1024px）
   - 方案B：单张显示 + 左右箭头切换
   - **选择**：方案B，保持一致性，避免布局变化

3. **交互反馈**：
   - 按钮点击：`scale(0.95)` 动画 + 200ms
   - 收藏状态变化：星星图标填充动画 + 300ms
   - 滑动切换：卡片滑动动画 + 300ms ease-out
   - 加载状态：骨架屏或加载动画

## Phase 2: Implementation Tasks

### Task 1: 安装依赖和配置字体

**文件**：`frontend/package.json`, `frontend/src/index.css`

**任务**：
1. 安装 `react-swipeable` 和 `html2canvas`
2. 在 `index.css` 中引入儿童友好字体（Google Fonts）
3. 扩展主题配置，添加儿童友好字体和色彩配置

**验收标准**：
- 依赖安装成功
- 字体正确加载
- 主题配置可用

### Task 2: 创建卡片样式配置

**文件**：`frontend/src/styles/cardStyles.ts`

**任务**：
1. 定义儿童友好字体配置
2. 定义固定比例配置
3. 定义色彩配置（增强对比度）

**验收标准**：
- 配置导出可用
- 字体配置符合规范（字号≥14px，行高≥1.6倍）

### Task 3: 创建卡片滑动Hook

**文件**：`frontend/src/hooks/useCardSwipe.ts`

**任务**：
1. 实现触摸事件处理
2. 实现滑动阈值判断
3. 实现滑动动画控制

**验收标准**：
- Hook可正常使用
- 滑动响应时间≤100ms
- 支持左右滑动

### Task 4: 创建卡片导出Hook

**文件**：`frontend/src/hooks/useCardExport.ts`, `frontend/src/utils/cardExport.ts`

**任务**：
1. 封装 `html2canvas` 调用
2. 实现图片下载功能
3. 处理移动端和PC端的差异

**验收标准**：
- 导出功能正常
- 图片质量清晰
- 移动端和PC端都能使用

### Task 5: 创建CardCarousel组件

**文件**：`frontend/src/components/cards/CardCarousel.tsx`

**任务**：
1. 实现渐进式展示逻辑（只显示已生成的卡片）
2. 集成 `react-swipeable` 实现滑动切换
3. 实现固定比例显示（`aspect-ratio: 16/9`）
4. 实现响应式布局（移动端单张，PC端可选并排）
5. 集成收藏和导出功能

**验收标准**：
- 组件可正常渲染
- 滑动切换流畅（响应时间≤100ms）
- 固定比例显示正确
- 移动端和PC端适配正确

### Task 6: 优化ScienceCard组件

**文件**：`frontend/src/components/cards/ScienceCard.tsx`

**任务**：
1. 移除内容区域的 `overflow-y-auto`
2. 添加固定比例容器
3. 应用儿童友好字体样式
4. 优化色彩对比度
5. 优化交互反馈动画

**验收标准**：
- 卡片内容在固定比例内完整显示
- 字体样式符合儿童友好要求
- 交互反馈及时明确

### Task 7: 优化PoetryCard组件

**文件**：`frontend/src/components/cards/PoetryCard.tsx`

**任务**：同Task 6

**验收标准**：同Task 6

### Task 8: 优化EnglishCard组件

**文件**：`frontend/src/components/cards/EnglishCard.tsx`

**任务**：同Task 6

**验收标准**：同Task 6

### Task 9: 修改ConversationMessage组件

**文件**：`frontend/src/components/conversation/ConversationMessage.tsx`

**任务**：
1. 检测卡片消息，使用 `CardCarousel` 渲染
2. 传递收藏和导出回调
3. 保持其他消息类型不变

**验收标准**：
- 卡片消息使用CardCarousel渲染
- 其他消息类型正常显示
- 收藏和导出功能正常

### Task 10: 测试和优化

**文件**：所有相关文件

**任务**：
1. 测试移动端滑动功能
2. 测试PC端显示和切换
3. 测试字体可读性
4. 测试色彩对比度
5. 测试收藏和导出功能
6. 性能优化（减少重渲染）

**验收标准**：
- 所有功能正常
- 性能指标达标（响应时间、渲染时间）
- 符合成功标准（SC-001至SC-011）

## Risk & Mitigation

### 风险1：固定比例可能导致内容溢出

**风险**：卡片内容过长，在16:9比例内无法完整显示

**缓解措施**：
- 使用响应式字体大小（`clamp()`）
- 优化内容布局，合理使用间距
- 如果内容确实过长，允许轻微滚动（但尽量避免）

### 风险2：滑动动画性能问题

**风险**：滑动动画在低端设备上卡顿

**缓解措施**：
- 使用CSS `transform` 而非 `left/top`（GPU加速）
- 使用 `will-change` 提示浏览器优化
- 限制动画帧率，使用 `requestAnimationFrame`

### 风险3：字体加载延迟

**风险**：Google Fonts加载慢，导致字体闪烁

**缓解措施**：
- 使用 `font-display: swap` 确保文本立即显示
- 提供系统字体回退
- 考虑预加载字体

### 风险4：html2canvas兼容性问题

**风险**：某些浏览器不支持或渲染效果差

**缓解措施**：
- 检测浏览器支持情况
- 提供降级方案（截图API或提示用户手动截图）
- 测试主流浏览器兼容性

## Success Metrics

- **SC-001**: 页面整体布局结构保持不变率100% ✅
- **SC-002**: 100%的卡片内容在固定比例内完整显示 ✅
- **SC-003**: 卡片内容自动撑开到固定比例的成功率达到100% ✅
- **SC-004**: 移动端卡片滑动操作的响应时间≤100毫秒 ✅
- **SC-005**: 卡片字体可读性测试中，8-12岁儿童能够轻松阅读的比例≥95% ✅
- **SC-006**: 卡片色彩搭配评分≥4.5/5.0 ✅
- **SC-007**: 卡片交互反馈响应时间≤200毫秒 ✅
- **SC-008**: 卡片收藏操作的响应时间≤300毫秒 ✅
- **SC-009**: 卡片保存到本地的成功率≥98% ✅
- **SC-010**: PC端三张卡片一并展开时，布局适配成功率100% ✅
- **SC-011**: 渐进式卡片展示功能中，第一张卡片显示时间≤2秒 ✅

