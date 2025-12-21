# Tasks: 前端中文优先国际化

**Input**: Design documents from `/specs/010-frontend-i18n-chinese/`
**Prerequisites**: plan.md (required), spec.md (required for user stories)

**Tests**: 手动测试任务已包含在Phase 5中

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Web app**: `frontend/src/`
- Paths shown below use frontend structure

---

## Phase 1: Foundational (Blocking Prerequisites)

**Purpose**: 扩展翻译文件，这是所有用户故事的基础，必须首先完成

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [ ] T001 [P] 扩展中文翻译文件 `frontend/src/i18n/locales/zh.ts`，添加所有新翻译key（header、home、capture、result、collection、report、settings、littleStar等命名空间）
- [ ] T002 [P] 扩展英文翻译文件 `frontend/src/i18n/locales/en.ts`，添加所有新翻译key，确保与中文翻译文件结构一致
- [ ] T003 验证i18n配置 `frontend/src/i18n/index.ts`，确保默认语言为中文 (`lng: 'zh'`)，fallback语言为中文 (`fallbackLng: 'zh'`)

**Checkpoint**: Foundation ready - 翻译文件已扩展，i18n配置已验证，用户故事实现可以开始

---

## Phase 2: User Story 1 - 首次访问应用看到全中文界面 (Priority: P1) 🎯 MVP

**Goal**: 用户首次打开应用时，所有页面默认显示中文，不出现任何英文文本

**Independent Test**: 清除浏览器缓存和localStorage，打开应用，检查所有页面（首页、拍照页、对话页、收藏页、报告页、设置页）是否全部显示中文，没有任何英文文本

### Implementation for User Story 1

- [x] T004 [US1] 替换Header组件 `frontend/src/components/common/Header.tsx` 中的硬编码文本：
  - 将title默认值 `'Little Explorer'` 替换为 `t('header.title')`
  - 将收藏链接文本 `'My Favorites'` 替换为 `t('header.favorites')`
- [x] T005 [US1] 替换首页 `frontend/src/pages/Home.tsx` 中的硬编码文本：
  - 将科学认知卡片标题 `'科学认知'` 替换为 `t('home.cardScience')`
  - 将人文素养卡片标题 `'人文素养'` 替换为 `t('home.cardHumanities')`
  - 将语言能力卡片标题 `'语言能力'` 替换为 `t('home.cardLanguage')`
  - 将LittleStar消息 `'拍一拍，发现有趣的知识吧～'` 替换为 `t('home.littleStarMessage')`
- [x] T006 [US1] 替换拍照页 `frontend/src/pages/Capture.tsx` 中的硬编码文本：
  - 将Header标签 `'AI Auto-Detect'` 替换为 `t('capture.aiAutoDetect')`
- [x] T007 [US1] 替换对话页 `frontend/src/pages/Result.tsx` 中的硬编码文本：
  - 将发现新朋友提示 `'You found a new friend!'` 替换为 `t('result.foundNewFriend')`
  - 将标题前缀 `'It's a'` 替换为 `t('result.itsA')`
  - 将AI Companion标签 `'AI Companion says:'` 替换为 `t('result.aiCompanionSays')`
  - 将AI Companion消息fallback替换为使用 `t('result.aiCompanionMessage', { objectName })`
- [x] T008 [US1] 替换收藏页 `frontend/src/pages/Collection.tsx` 中的硬编码文本：
  - 将页面标题 `'My Favorites'` 替换为 `t('collection.title')`
  - 将副标题 `'Keep exploring your collection of wonders!'` 替换为 `t('collection.subtitle')`
  - 将导出全部按钮 `'导出全部'` 替换为 `t('collection.exportAll')`
  - 将家长模式标签 `'Parent Mode'` 替换为 `t('collection.parentMode')`
  - 将清空所有按钮 `'Clear All'` 替换为 `t('collection.clearAll')`
  - 将清空所有提示 `'Only available in Parent Mode'` 替换为 `t('collection.clearAllHint')`
  - 将Little Star Says标签 `'Little Star Says:'` 替换为 `t('collection.littleStarSays')`
  - 将Little Star消息替换为 `t('collection.littleStarMessage')`
  - 将导出失败提示 `'导出失败，请重试'` 替换为 `t('collection.exportError')`
  - 将加载中 `'加载中...'` 替换为 `t('common.loading')`
- [x] T009 [US1] 替换报告页 `frontend/src/pages/LearningReport.tsx` 中的硬编码文本：
  - 将所有英文文本替换为对应的翻译key（report命名空间下的所有key）
  - 将中文硬编码文本也替换为翻译key（如"最近收藏了"、"还没有收藏任何卡片"等）
- [x] T010 [US1] 替换设置页 `frontend/src/pages/Settings.tsx` 中的硬编码文本：
  - 将所有年级标签（K1-K3, G1-G12）替换为使用翻译key（`settings.gradeK1` 到 `settings.gradeG12`）
  - 将应用描述 `'TanGo - 探索世界的知识卡片应用'` 替换为 `t('settings.appDescription')`
- [x] T011 [US1] 替换LittleStar组件 `frontend/src/components/common/LittleStar.tsx` 中的硬编码文本：
  - 将名称标签 `'Little Star'` 替换为 `t('littleStar.name')`
- [x] T012 [US1] 替换CollectionGrid组件 `frontend/src/components/collection/CollectionGrid.tsx` 中的硬编码文本：
  - 将空状态消息 `'还没有收藏任何卡片，快去探索吧！'` 替换为 `t('collection.emptyMessage')`
  - 将导出失败提示 `'导出失败，请重试'` 替换为 `t('collection.exportError')`

**Checkpoint**: User Story 1完成 - 清除localStorage后，所有页面默认显示中文，无任何英文文本

---

## Phase 3: User Story 2 - 在设置页面切换语言 (Priority: P2)

**Goal**: 用户在设置页面可以通过语言切换器选择中文或英文，切换后立即生效，无需刷新页面

**Independent Test**: 打开设置页面，使用语言切换器从中文切换到英文，验证所有页面立即更新为英文；再切换回中文，验证所有页面立即更新为中文

### Implementation for User Story 2

- [ ] T013 [US2] 验证语言切换器 `frontend/src/components/common/LanguageSwitcher.tsx` 功能正常：
  - 确认切换语言后调用 `i18n.changeLanguage()` 和 `changeLanguage()`
  - 确认语言设置保存到localStorage
  - 确认切换后页面立即更新（无需刷新）
- [ ] T014 [US2] 验证i18n配置 `frontend/src/i18n/index.ts` 支持语言切换：
  - 确认 `i18n.on('languageChanged')` 监听器正确保存语言设置
  - 确认语言切换后所有使用 `useTranslation()` 的组件自动更新

**Checkpoint**: User Story 2完成 - 语言切换功能正常工作，切换后立即生效，设置持久化保存

---

## Phase 4: User Story 3 - 所有页面支持中英文切换 (Priority: P3)

**Goal**: 用户切换语言后，应用的所有页面都能正确显示对应语言的文本，包括header、按钮、提示、标签等所有UI元素

**Independent Test**: 在任意页面切换语言，验证当前页面和所有其他页面的文本都正确更新为对应语言，没有遗漏的硬编码文本

### Implementation for User Story 3

- [ ] T015 [US3] 验证首页语言切换：切换语言后，首页所有文本（header、按钮、卡片标题、LittleStar消息）立即更新
- [ ] T016 [US3] 验证拍照页语言切换：切换语言后，拍照页所有文本（header、按钮、提示）立即更新
- [ ] T017 [US3] 验证对话页语言切换：切换语言后，对话页所有文本（header、消息提示、按钮）立即更新
- [ ] T018 [US3] 验证收藏页语言切换：切换语言后，收藏页所有文本（标题、按钮、提示）立即更新
- [ ] T019 [US3] 验证报告页语言切换：切换语言后，报告页所有文本（标题、统计标签、提示）立即更新
- [ ] T020 [US3] 验证设置页语言切换：切换语言后，设置页所有文本（标题、设置项、年级标签）立即更新

**Checkpoint**: User Story 3完成 - 所有页面都正确支持语言切换，无遗漏的硬编码文本

---

## Phase 5: Testing & Validation

**Purpose**: 功能测试、完整性检查和边界情况测试

### 功能测试

- [ ] T021 [P] 测试默认语言：清除localStorage，打开应用，验证所有页面默认显示中文
- [ ] T022 [P] 测试语言切换：在设置页面切换语言，验证所有页面立即更新为对应语言
- [ ] T023 [P] 测试语言持久化：切换语言后刷新页面，验证语言设置保持不变
- [ ] T024 [P] 测试页面导航：切换语言后导航到其他页面，验证新页面使用选择的语言

### 完整性检查

- [ ] T025 [P] 代码审查：检查所有页面无硬编码英文文本（使用grep搜索常见英文单词）
- [ ] T026 [P] 代码审查：检查所有页面无硬编码中文文本（应使用i18n，除了注释）
- [ ] T027 [P] 翻译文件完整性：验证中文翻译文件覆盖所有使用的key
- [ ] T028 [P] 翻译文件完整性：验证英文翻译文件覆盖所有使用的key
- [ ] T029 [P] 翻译文件一致性：验证中文和英文翻译文件结构一致，无缺失key

### 边界情况测试

- [ ] T030 [P] 测试翻译key缺失：临时删除某个翻译key，验证fallback到中文显示，而不是显示key名称
- [ ] T031 [P] 测试快速切换语言：快速连续切换语言，验证正确处理最后一次选择的语言
- [ ] T032 [P] 测试localStorage清除：清除localStorage后，验证应用恢复为默认中文
- [ ] T033 [P] 测试语言切换过程中页面加载：在页面加载过程中切换语言，验证新加载的内容使用新语言

---

## Dependencies & Execution Order

### Phase Dependencies

- **Foundational (Phase 1)**: No dependencies - can start immediately
- **User Story 1 (Phase 2)**: Depends on Foundational completion - BLOCKS User Stories 2 and 3
- **User Story 2 (Phase 3)**: Depends on Foundational completion - Can run in parallel with US1 after Phase 1
- **User Story 3 (Phase 4)**: Depends on Foundational completion - Can run in parallel with US1/US2 after Phase 1
- **Testing (Phase 5)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 1) - No dependencies on other stories
- **User Story 2 (P2)**: Can start after Foundational (Phase 1) - Independent of US1 and US3
- **User Story 3 (P3)**: Can start after Foundational (Phase 1) - Independent of US1 and US2

### Within Each User Story

- 翻译文件必须在替换硬编码文本之前完成
- 每个页面的替换可以并行进行（不同文件）
- 替换完成后立即测试该页面

### Parallel Opportunities

- Phase 1中的T001和T002可以并行（不同文件）
- Phase 2中的T004-T012可以并行（不同文件）
- Phase 3中的T013和T014可以并行（不同文件）
- Phase 4中的T015-T020可以并行（不同页面）
- Phase 5中的所有测试任务可以并行

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Foundational (扩展翻译文件)
2. Complete Phase 2: User Story 1 (替换所有硬编码文本为中文默认)
3. **STOP and VALIDATE**: 测试User Story 1 - 清除localStorage后所有页面显示中文
4. Deploy/demo if ready

### Incremental Delivery

1. Complete Phase 1 → Foundation ready
2. Add User Story 1 → Test independently → Deploy/Demo (MVP - 全中文界面!)
3. Add User Story 2 → Test independently → Deploy/Demo (语言切换功能)
4. Add User Story 3 → Test independently → Deploy/Demo (完整语言支持)
5. Complete Phase 5 → Final validation → Deploy

### Parallel Team Strategy

With multiple developers:

1. Team completes Phase 1 together (翻译文件)
2. Once Phase 1 is done:
   - Developer A: User Story 1 (替换硬编码文本)
   - Developer B: User Story 2 (验证语言切换)
   - Developer C: User Story 3 (验证所有页面)
3. All complete Phase 5 together (测试和验证)

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Commit after each task or logical group (e.g., after completing one page)
- Stop at any checkpoint to validate story independently
- 翻译文件更新后，建议立即验证i18n配置是否正确加载新key
- 替换硬编码文本时，注意保持原有的样式和格式
- 测试时注意检查动态文本（如包含变量的消息）是否正确处理

