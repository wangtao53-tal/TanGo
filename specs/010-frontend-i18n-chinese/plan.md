# Implementation Plan: 前端中文优先国际化

**Branch**: `dev-mvp-20251218` | **Date**: 2025-12-22 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/010-frontend-i18n-chinese/spec.md`

**Note**: MVP版本阶段，所有开发工作统一在 `dev-mvp-20251218` 分支进行，不采用一个功能一个分支的策略。

## Summary

实现前端应用的中文优先国际化功能，确保所有页面默认显示中文，支持用户切换英文。将所有硬编码的英文文本替换为i18n翻译key，完善中文和英文翻译文件，确保语言切换功能在所有页面正常工作。

**技术方案**: 使用现有的react-i18next框架，扩展翻译文件，替换所有硬编码文本，确保默认语言为中文，语言切换后立即生效。

## Technical Context

**Language/Version**: TypeScript 5.x, React 18.x  
**Primary Dependencies**: react-i18next, i18next  
**Storage**: localStorage (语言设置持久化)  
**Testing**: 手动测试 + 代码审查  
**Target Platform**: Web应用（桌面和移动端浏览器）  
**Project Type**: Web application (frontend)  
**Performance Goals**: 语言切换响应时间 < 1秒，无需页面刷新  
**Constraints**: 
- 必须保持现有i18n框架结构
- 默认语言必须为中文
- 所有页面必须支持语言切换
- 翻译缺失时fallback到中文

**Scale/Scope**: 
- 6个主要页面（首页、拍照页、对话页、收藏页、报告页、设置页）
- 多个共享组件（Header、LittleStar等）
- 约100+个需要翻译的文本key

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

**合规性说明**：本功能完全符合所有规范要求，特别是原则四（多语言和年级设置规范）是本功能的核心目标。

## Project Structure

### Documentation (this feature)

```text
specs/010-frontend-i18n-chinese/
├── plan.md              # This file (/speckit.plan command output)
├── spec.md              # Feature specification
└── checklists/
    └── requirements.md  # Specification quality checklist
```

### Source Code (repository root)

```text
frontend/
├── src/
│   ├── i18n/
│   │   ├── index.ts                    # i18n配置（需更新默认语言为中文）
│   │   └── locales/
│   │       ├── zh.ts                   # 中文翻译文件（需扩展）
│   │       └── en.ts                   # 英文翻译文件（需扩展）
│   ├── components/
│   │   ├── common/
│   │   │   ├── Header.tsx              # Header组件（需替换硬编码文本）
│   │   │   └── LittleStar.tsx          # LittleStar组件（需替换硬编码文本）
│   │   └── collection/
│   │       └── CollectionGrid.tsx     # 收藏网格组件（需替换硬编码文本）
│   └── pages/
│       ├── Home.tsx                    # 首页（需替换硬编码文本）
│       ├── Capture.tsx                 # 拍照页（需替换硬编码文本）
│       ├── Result.tsx                  # 对话页（需替换硬编码文本）
│       ├── Collection.tsx              # 收藏页（需替换硬编码文本）
│       ├── LearningReport.tsx          # 报告页（需替换硬编码文本）
│       └── Settings.tsx                 # 设置页（需替换年级标签）
```

**Structure Decision**: 使用现有的前端项目结构，主要修改i18n翻译文件和各个页面/组件文件，替换硬编码文本为翻译key。

## 需要国际化的文本清单

### 1. Header组件 (`frontend/src/components/common/Header.tsx`)

| 位置 | 当前文本 | 翻译Key | 中文 | 英文 |
|------|---------|---------|------|------|
| title默认值 | `'Little Explorer'` | `header.title` | 小小探索家 | Little Explorer |
| 收藏链接 | `'My Favorites'` | `header.favorites` | 我的收藏 | My Favorites |

### 2. 首页 (`frontend/src/pages/Home.tsx`)

| 位置 | 当前文本 | 翻译Key | 中文 | 英文 |
|------|---------|---------|------|------|
| 科学认知卡片标题 | `'科学认知'` | `home.cardScience` | 科学认知 | Science |
| 人文素养卡片标题 | `'人文素养'` | `home.cardHumanities` | 人文素养 | Humanities |
| 语言能力卡片标题 | `'语言能力'` | `home.cardLanguage` | 语言能力 | Language |
| LittleStar消息 | `'拍一拍，发现有趣的知识吧～'` | `home.littleStarMessage` | 拍一拍，发现有趣的知识吧～ | Take a photo and discover interesting knowledge! |

### 3. 拍照页 (`frontend/src/pages/Capture.tsx`)

| 位置 | 当前文本 | 翻译Key | 中文 | 英文 |
|------|---------|---------|------|------|
| Header标签 | `'AI Auto-Detect'` | `capture.aiAutoDetect` | AI自动识别 | AI Auto-Detect |

### 4. 对话页 (`frontend/src/pages/Result.tsx`)

| 位置 | 当前文本 | 翻译Key | 中文 | 英文 |
|------|---------|---------|------|------|
| 发现新朋友提示 | `'You found a new friend!'` | `result.foundNewFriend` | 你发现了一个新朋友！ | You found a new friend! |
| 标题前缀 | `'It's a'` | `result.itsA` | 这是一个 | It's a |
| AI Companion标签 | `'AI Companion says:'` | `result.aiCompanionSays` | AI小伙伴说： | AI Companion says: |
| AI Companion消息fallback | `'"Wow! A ${objectName}! Let's explore its secrets!"'` | `result.aiCompanionMessage` | "哇！这是一个${objectName}！让我们探索它的秘密吧！" | "Wow! A ${objectName}! Let's explore its secrets!" |

### 5. 收藏页 (`frontend/src/pages/Collection.tsx`)

| 位置 | 当前文本 | 翻译Key | 中文 | 英文 |
|------|---------|---------|------|------|
| 页面标题 | `'My Favorites'` | `collection.title` | 我的收藏 | My Favorites |
| 副标题 | `'Keep exploring your collection of wonders!'` | `collection.subtitle` | 继续探索你的收藏吧！ | Keep exploring your collection of wonders! |
| 导出全部按钮 | `'导出全部'` | `collection.exportAll` | 导出全部 | Export All |
| 家长模式标签 | `'Parent Mode'` | `collection.parentMode` | 家长模式 | Parent Mode |
| 清空所有按钮 | `'Clear All'` | `collection.clearAll` | 清空所有 | Clear All |
| 清空所有提示 | `'Only available in Parent Mode'` | `collection.clearAllHint` | 仅在家长模式下可用 | Only available in Parent Mode |
| Little Star Says标签 | `'Little Star Says:'` | `collection.littleStarSays` | 小星星说： | Little Star Says: |
| Little Star消息 | `'Go explore interesting knowledge and collect more favorite cards! I'm waiting for your discoveries! ✨'` | `collection.littleStarMessage` | 去探索有趣的知识，收藏更多喜欢的卡片吧！我在等待你的发现！✨ | Go explore interesting knowledge and collect more favorite cards! I'm waiting for your discoveries! ✨ |
| 导出失败提示 | `'导出失败，请重试'` | `collection.exportError` | 导出失败，请重试 | Export failed, please try again |
| 加载中 | `'加载中...'` | `common.loading` | 加载中... | Loading... |

### 6. 报告页 (`frontend/src/pages/LearningReport.tsx`)

| 位置 | 当前文本 | 翻译Key | 中文 | 英文 |
|------|---------|---------|------|------|
| 报告标签 | `'Weekly Report'` | `report.weeklyReport` | 周报 | Weekly Report |
| 标题问候 | `'Hi, Little Explorer!'` | `report.greeting` | 你好，小小探索家！ | Hi, Little Explorer! |
| 副标题 | `'You're doing great! Look at your growth this week.'` | `report.subtitle` | 你做得很好！看看你这周的成长吧。 | You're doing great! Look at your growth this week. |
| 探索次数标签 | `'Exploration Stars'` | `report.explorationStars` | 探索次数 | Exploration Stars |
| 探索鼓励 | `'Keep exploring!'` | `report.keepExploring` | 继续探索！ | Keep exploring! |
| 收藏总数标签 | `'Total Favorites'` | `report.totalFavorites` | 收藏总数 | Total Favorites |
| 收藏鼓励 | `'Great collection!'` | `report.greatCollection` | 收藏很棒！ | Great collection! |
| 专家等级标签 | `'Little Expert'` | `report.littleExpert` | 小小专家 | Little Expert |
| 专家等级名称 | `'Nature Master'` | `report.natureMaster` | 自然大师 | Nature Master |
| 升级提示 | `'Level Up! 🚀'` | `report.levelUp` | 升级了！🚀 | Level Up! 🚀 |
| 知识地图标题 | `'Knowledge Map'` | `report.knowledgeMap` | 知识地图 | Knowledge Map |
| 总数标签 | `'Total'` | `report.total` | 总数 | Total |
| 自然类标签 | `'Natural'` | `report.categoryNatural` | 自然类 | Natural |
| 生活类标签 | `'Life'` | `report.categoryLife` | 生活类 | Life |
| 人文类标签 | `'Humanities'` | `report.categoryHumanities` | 人文类 | Humanities |
| 项目数标签 | `'items'` | `report.items` | 项 | items |
| 最近收藏标题 | `'Recent Favorites'` | `report.recentFavorites` | 最近收藏 | Recent Favorites |
| 最近收藏消息 | `'最近收藏了 {totalCollectedCards} 张卡片'` | `report.recentFavoritesMessage` | 最近收藏了 {totalCollectedCards} 张卡片 | Recently collected {totalCollectedCards} cards |
| 空状态消息 | `'还没有收藏任何卡片'` | `report.noCards` | 还没有收藏任何卡片 | No cards collected yet |

### 7. 设置页 (`frontend/src/pages/Settings.tsx`)

| 位置 | 当前文本 | 翻译Key | 中文 | 英文 |
|------|---------|---------|------|------|
| K1标签 | `'Kindergarten 1'` | `settings.gradeK1` | 幼儿园小班 | Kindergarten 1 |
| K2标签 | `'Kindergarten 2'` | `settings.gradeK2` | 幼儿园中班 | Kindergarten 2 |
| K3标签 | `'Kindergarten 3'` | `settings.gradeK3` | 幼儿园大班 | Kindergarten 3 |
| G1-G12标签 | `'Grade 1'` - `'Grade 12'` | `settings.gradeG1` - `settings.gradeG12` | 一年级 - 十二年级 | Grade 1 - Grade 12 |
| 关于应用描述 | `'TanGo - 探索世界的知识卡片应用'` | `settings.appDescription` | TanGo - 探索世界的知识卡片应用 | TanGo - Knowledge Card App for Exploring the World |

### 8. LittleStar组件 (`frontend/src/components/common/LittleStar.tsx`)

| 位置 | 当前文本 | 翻译Key | 中文 | 英文 |
|------|---------|---------|------|------|
| 名称标签 | `'Little Star'` | `littleStar.name` | 小星星 | Little Star |

### 9. CollectionGrid组件 (`frontend/src/components/collection/CollectionGrid.tsx`)

| 位置 | 当前文本 | 翻译Key | 中文 | 英文 |
|------|---------|---------|------|------|
| 空状态消息 | `'还没有收藏任何卡片，快去探索吧！'` | `collection.emptyMessage` | 还没有收藏任何卡片，快去探索吧！ | No cards collected yet, go explore! |
| 导出失败提示 | `'导出失败，请重试'` | `collection.exportError` | 导出失败，请重试 | Export failed, please try again |

### 10. 其他通用文本

| 位置 | 当前文本 | 翻译Key | 中文 | 英文 |
|------|---------|---------|------|------|
| 学习报告链接 | `'Learning Report'` | `common.report` | 学习报告 | Learning Report |

## 实现步骤

### Phase 1: 扩展翻译文件

1. **更新中文翻译文件** (`frontend/src/i18n/locales/zh.ts`)
   - 添加所有新发现的翻译key
   - 确保覆盖所有页面和组件

2. **更新英文翻译文件** (`frontend/src/i18n/locales/en.ts`)
   - 添加所有新发现的翻译key
   - 确保与中文翻译文件结构一致

3. **验证i18n配置** (`frontend/src/i18n/index.ts`)
   - 确保默认语言为中文 (`lng: 'zh'`)
   - 确保fallback语言为中文 (`fallbackLng: 'zh'`)

### Phase 2: 替换硬编码文本

按页面顺序替换所有硬编码文本：

1. **Header组件** - 替换title默认值和链接文本
2. **首页** - 替换卡片标题和LittleStar消息
3. **拍照页** - 替换header标签
4. **对话页** - 替换所有英文提示文本
5. **收藏页** - 替换所有英文文本和中文硬编码文本
6. **报告页** - 替换所有英文文本和中文硬编码文本
7. **设置页** - 替换年级标签和应用描述
8. **LittleStar组件** - 替换名称标签
9. **CollectionGrid组件** - 替换空状态和错误消息

### Phase 3: 测试和验证

1. **功能测试**
   - 清除localStorage，验证默认显示中文
   - 切换语言，验证所有页面立即更新
   - 刷新页面，验证语言设置持久化

2. **完整性检查**
   - 检查所有页面无硬编码英文文本
   - 检查所有页面无硬编码中文文本（应使用i18n）
   - 验证翻译文件完整性

3. **边界情况测试**
   - 测试翻译key缺失时的fallback行为
   - 测试快速切换语言的处理
   - 测试localStorage清除后的恢复

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

无违反规范的情况。

