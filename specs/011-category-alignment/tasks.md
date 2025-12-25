# 任务清单：知识卡片分类功能对齐与优化

**输入**: 设计文档来自 `/specs/011-category-alignment/`
**前置条件**: plan.md (必需), spec.md (必需，用于用户故事)

**注意**: 本次优化主要在前端进行，不涉及后端API修改。所有分类相关的数据结构已存在，只需确保正确使用和显示。

**组织方式**: 任务按用户故事分组，支持独立实现和测试每个故事。

## 格式说明: `[ID] [P?] [Story] 描述`

- **[P]**: 可以并行执行（不同文件，无依赖）
- **[Story]**: 该任务属于哪个用户故事（如 US1, US2, US3）
- 描述中包含确切的文件路径

## 路径约定

- **Web应用**: `frontend/src/`
- 路径基于 plan.md 中的项目结构

---

## Phase 1: 数据一致性验证与修复（优先级: P1）

**目的**: 确保所有页面的分类信息一致，验证数据保存和筛选逻辑

**独立测试**: 拍照识别一个对象，验证识别结果页显示正确的分类标签，收藏后，在收藏页面和学习报告中都能看到相同的分类信息。

### 1.1 验证探索记录保存逻辑

- [x] T001 [US1] 检查识别结果页保存探索记录的逻辑 `frontend/src/pages/Result.tsx`
  - 验证保存时是否包含objectCategory字段
  - 确认objectCategory值来自识别API响应
  - 确保objectCategory值必须是"自然类"、"生活类"、"人文类"之一

- [x] T002 [US1] 验证探索记录存储服务 `frontend/src/services/storage.ts`
  - 检查explorationStorage.save方法是否正确保存objectCategory字段
  - 验证IndexedDB索引是否支持按objectCategory查询
  - 确认数据格式与类型定义一致

- [x] T003 [US1] 修复探索记录保存逻辑（如需要） `frontend/src/pages/Result.tsx`
  - 如果objectCategory字段缺失，添加保存逻辑
  - 确保从识别API响应中提取objectCategory并保存
  - 添加数据验证，确保objectCategory值有效

### 1.2 验证分类筛选逻辑

- [x] T004 [US1] 检查收藏页面分类筛选组件 `frontend/src/components/collection/CategoryFilter.tsx`
  - 验证筛选选项是否正确（全部、自然类、生活类、人文类）
  - 确认筛选状态管理正常
  - 测试筛选按钮交互

- [x] T005 [US1] 验证收藏网格筛选逻辑 `frontend/src/components/collection/CollectionGrid.tsx`
  - 检查筛选时是否正确使用objectCategory字段
  - 验证"全部"选项显示所有记录
  - 验证特定分类选项只显示对应分类的记录
  - 测试边界情况（空数据、无效分类等）

- [x] T006 [US1] 修复分类筛选逻辑（如需要） `frontend/src/components/collection/CollectionGrid.tsx`
  - 如果筛选逻辑有问题，修复筛选条件
  - 确保筛选结果准确率达到100%
  - 添加错误处理和日志记录

**检查点**: 数据一致性验证完成，探索记录正确保存分类信息，筛选功能正常工作 ✅

---

## Phase 2: UI优化 - 移除无用按钮（优先级: P1）

**目的**: 简化界面，移除无用的UI元素

**独立测试**: 验证学习报告页面没有tip按钮，收藏页面没有"家长模式"和"导出全部"按钮。

### 2.1 移除学习报告tip按钮

- [x] T007 [US2] 移除学习报告知识地图tip按钮 `frontend/src/pages/LearningReport.tsx`
  - 删除第151-153行的tip按钮代码
  - 保持知识地图的其他功能不变
  - 验证页面布局正常，无样式错乱

### 2.2 移除收藏页面无用按钮

- [x] T008 [US3] 移除收藏页面"导出全部"按钮 `frontend/src/pages/Collection.tsx`
  - 删除第108-118行的导出全部按钮代码
  - 删除handleExportAll函数（如果不再使用）
  - 简化页面头部布局

- [x] T009 [US3] 移除收藏页面"家长模式"按钮 `frontend/src/pages/Collection.tsx`
  - 删除第120-136行的家长模式控制代码
  - 删除handleClearAll函数（如果不再使用）
  - 验证页面头部布局正常

**检查点**: UI优化完成，无用按钮已移除，页面布局正常 ✅

---

## Phase 3: 学习报告统计逻辑修复（优先级: P1）

**目的**: 修复统计数据，确保知识地图的三类数量总和与总探索次数完全一致

**独立测试**: 完成3次探索（自然类2次、生活类1次），验证学习报告中总探索次数为3，知识地图中自然类显示2次、生活类显示1次，且知识地图总数与总探索次数一致。

### 3.1 修复知识地图统计逻辑

- [x] T010 [US2] 修复学习报告统计逻辑 `frontend/src/pages/LearningReport.tsx`
  - 修改loadReportData函数中的统计逻辑
  - 确保遍历所有records，统计每个objectCategory的数量
  - 处理边界情况：当objectCategory为空或无效时，使用默认值"自然类"
  - 添加数据验证：确保totalCategories === totalExplorations

- [x] T011 [US2] 添加数据一致性验证 `frontend/src/pages/LearningReport.tsx`
  - 在统计完成后验证数据一致性
  - 如果totalCategories !== totalExplorations，输出警告日志
  - 确保统计逻辑正确，三类数量总和等于总探索次数

- [x] T012 [US2] 确保统计数据实时更新 `frontend/src/pages/LearningReport.tsx`
  - 验证useEffect依赖项正确
  - 确保页面重新加载时统计数据正确
  - 测试完成新探索后统计数据是否正确更新

**检查点**: 统计逻辑修复完成，知识地图三类数量总和等于总探索次数 ✅

---

## Phase 4: 收藏功能优化（优先级: P2）

**目的**: 优化收藏交互体验，支持取消收藏和点击卡片收藏

**独立测试**: 在收藏页面验证点击收藏按钮可以取消收藏，点击卡片可以收藏。

### 4.1 实现取消收藏功能

- [x] T013 [US3] 添加收藏状态到CollectionCard组件 `frontend/src/components/collection/CollectionCard.tsx`
  - 添加isCollected属性到CollectionCardProps接口
  - 添加onToggleCollect回调函数到CollectionCardProps接口
  - 在组件中显示收藏按钮状态（已收藏/未收藏）

- [x] T014 [US3] 实现收藏按钮切换逻辑 `frontend/src/components/collection/CollectionCard.tsx`
  - 添加收藏按钮UI（使用Material Icons：star/star_border）
  - 实现点击收藏按钮切换收藏状态
  - 调用onToggleCollect回调函数

- [x] T015 [US3] 实现收藏状态切换处理 `frontend/src/pages/Collection.tsx`
  - 添加handleToggleCollect函数
  - 更新探索记录的collected字段
  - 保存更新后的记录到IndexedDB
  - 如果取消收藏，从列表中移除该记录
  - 重新加载数据以更新UI

### 4.2 实现点击卡片收藏功能

- [x] T016 [US3] 添加点击卡片收藏功能 `frontend/src/components/collection/CollectionCard.tsx`
  - 在CollectionCard根div上添加onClick事件
  - 如果卡片未收藏，点击卡片时调用onToggleCollect(true)
  - 如果卡片已收藏，点击卡片不执行操作（或执行其他操作）
  - 确保点击收藏按钮时阻止事件冒泡

**检查点**: 收藏功能优化完成，支持取消收藏和点击卡片收藏 ✅

---

## Phase 5: 测试与验证（优先级: P1）

**目的**: 确保所有功能正常工作，数据一致性达到100%

### 5.1 功能测试

- [ ] T017 [US1] 测试数据一致性 `frontend/src/pages/Result.tsx`, `frontend/src/pages/Collection.tsx`, `frontend/src/pages/LearningReport.tsx`
  - 拍照识别后，验证探索记录包含正确的objectCategory
  - 验证识别结果页显示正确的分类标签
  - 验证收藏页面筛选功能正常工作
  - 验证学习报告知识地图统计准确

- [ ] T018 [US2] 测试UI优化 `frontend/src/pages/LearningReport.tsx`, `frontend/src/pages/Collection.tsx`
  - 验证学习报告页面没有tip按钮
  - 验证收藏页面没有"家长模式"和"导出全部"按钮
  - 验证页面布局正常，无样式错乱

- [ ] T019 [US3] 测试收藏功能 `frontend/src/pages/Collection.tsx`, `frontend/src/components/collection/CollectionCard.tsx`
  - 验证点击收藏按钮可以取消收藏
  - 验证点击卡片可以收藏
  - 验证收藏状态正确更新
  - 验证本地存储正确保存收藏状态

- [ ] T020 [US2] 测试数据统计 `frontend/src/pages/LearningReport.tsx`
  - 验证知识地图三类数量总和等于总探索次数
  - 验证完成新探索后，统计数据正确更新
  - 验证边界情况处理正确（空数据、无效分类等）

### 5.2 性能测试

- [ ] T021 [P] 测试分类筛选性能 `frontend/src/components/collection/CollectionGrid.tsx`
  - 验证分类筛选响应时间≤200毫秒
  - 测试大量数据时的筛选性能

- [ ] T022 [P] 测试统计数据更新性能 `frontend/src/pages/LearningReport.tsx`
  - 验证统计数据更新响应时间≤1秒
  - 测试大量数据时的统计性能

- [ ] T023 [P] 测试收藏操作性能 `frontend/src/pages/Collection.tsx`
  - 验证收藏操作响应时间≤500毫秒
  - 测试并发操作处理

### 5.3 兼容性测试

- [ ] T024 [P] 测试移动端兼容性
  - 验证移动端显示正常
  - 验证触摸交互正常

- [ ] T025 [P] 测试PC端兼容性
  - 验证PC端显示正常
  - 验证鼠标交互正常

- [ ] T026 [P] 测试浏览器兼容性
  - 验证Chrome浏览器正常
  - 验证Safari浏览器正常
  - 验证Firefox浏览器正常

**检查点**: 所有测试通过，功能正常工作，性能满足要求 ✅

---

## 依赖关系和执行顺序

### Phase依赖

- **Phase 1 (数据一致性验证)**: 无依赖，可以立即开始
- **Phase 2 (UI优化)**: 无依赖，可以与Phase 1并行进行
- **Phase 3 (统计修复)**: 依赖Phase 1完成（需要确保数据正确保存）
- **Phase 4 (收藏优化)**: 依赖Phase 2完成（需要先移除无用按钮）
- **Phase 5 (测试验证)**: 依赖Phase 1-4完成

### 用户故事依赖

- **用户故事1 (P1)**: 无依赖，可以立即开始
- **用户故事2 (P1)**: 依赖用户故事1完成（需要确保数据一致性）
- **用户故事3 (P2)**: 依赖用户故事1完成（需要确保分类信息正确）

### 并行机会

- Phase 1和Phase 2可以完全并行进行
- Phase 1中的多个验证任务可以并行进行（标记[P]的任务）
- Phase 2中的两个移除按钮任务可以并行进行
- Phase 5中的测试任务可以并行进行（标记[P]的任务）

---

## 实施策略

### MVP优先（用户故事1和2）

1. 完成Phase 1: 数据一致性验证与修复
2. 完成Phase 2: UI优化
3. 完成Phase 3: 学习报告统计修复
4. **停止并验证**: 测试用户故事1和2独立功能
5. 可以演示核心功能

### 增量交付

1. 完成Phase 1 → 测试数据一致性 → 演示
2. 添加Phase 2 → 测试UI优化 → 演示
3. 添加Phase 3 → 测试统计修复 → 演示
4. 添加Phase 4 → 测试收藏优化 → 演示
5. 每个阶段独立交付价值

### 当前阶段重点

**优先完成**:
1. Phase 1: 数据一致性验证与修复（确保核心功能正常）
2. Phase 2: UI优化（简化界面）
3. Phase 3: 统计修复（确保数据准确）

**后续完成**:
- Phase 4: 收藏功能优化（功能增强）
- Phase 5: 全面测试验证

---

## 注意事项

- [P] 任务 = 不同文件，无依赖，可以并行
- [Story] 标签映射任务到特定用户故事，便于追踪
- 每个用户故事应该可以独立完成和测试
- 本次优化主要在前端进行，不涉及后端修改
- 所有分类相关的数据结构已存在，只需确保正确使用
- 提交代码前验证功能可用
- 在每个检查点停止验证故事独立性
- 避免：模糊任务、同一文件冲突、破坏独立性的跨故事依赖

---

## 任务统计

**总任务数**: 26个任务

**按用户故事分布**:
- 用户故事1 (US1): 6个任务 (T001-T006)
- 用户故事2 (US2): 6个任务 (T007, T010-T012, T018, T020)
- 用户故事3 (US3): 5个任务 (T008-T009, T013-T016, T019)
- 测试任务: 9个任务 (T017-T026)

**按优先级分布**:
- P1优先级: 20个任务
- P2优先级: 6个任务

**并行任务**: 12个任务可以并行执行

**预计总时间**: 
- Phase 1: 2小时
- Phase 2: 1小时
- Phase 3: 2小时
- Phase 4: 3小时
- Phase 5: 2小时
- **总计**: 约10小时
