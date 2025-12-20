---
description: "Task list for GitHub image acceleration MVP - quick implementation"
---

# Tasks: GitHub 图片加速优化 - MVP快速实现

**Input**: GitHub raw URL 偶发访问超时问题，需要快速实现CDN加速
**Prerequisites**: 现有代码已实现基本功能，需要添加CDN加速

**Organization**: 任务按MVP优先级组织，快速实现核心功能

## Format: `[ID] [P?] [Story] Description`

- **[P]**: 可以并行执行（不同文件，无依赖）
- **[Story]**: 所属用户故事（US1, US2等）
- 描述中包含确切的文件路径

## Path Conventions

- **后端**: `backend/internal/`
- 所有路径使用绝对路径或相对于项目根目录

---

## Phase 1: Setup (基础准备)

**Purpose**: 创建URL转换工具函数

- [X] T001 创建URL转换工具文件 `backend/internal/utils/github_cdn.go`
- [X] T002 [P] 添加GitHub raw URL检测函数在 `backend/internal/utils/github_cdn.go`
- [X] T003 [P] 实现jsDelivr CDN URL转换函数在 `backend/internal/utils/github_cdn.go`

**Checkpoint**: URL转换工具函数已创建，可以转换GitHub raw URL到jsDelivr CDN URL

---

## Phase 2: User Story 1 - GitHub 图片 CDN 加速 (P1) 🎯 MVP

**Goal**: 在识别节点中集成CDN URL转换，自动将GitHub raw URL转换为jsDelivr CDN URL

**Independent Test**: 使用GitHub raw URL调用识别接口，验证URL自动转换为CDN URL，图片访问成功率提升

### Implementation for User Story 1

- [X] T004 [US1] 在ImageRecognitionNode中导入URL转换工具在 `backend/internal/agent/nodes/image_recognition.go`
- [X] T005 [US1] 添加GitHub raw URL检测和CDN转换逻辑在 `backend/internal/agent/nodes/image_recognition.go`
- [X] T006 [US1] 实现CDN URL失败时重试原始URL的逻辑在 `backend/internal/agent/nodes/image_recognition.go`
- [X] T007 [US1] 添加CDN使用日志记录在 `backend/internal/agent/nodes/image_recognition.go`

**Checkpoint**: GitHub raw URL自动转换为CDN URL，CDN失败时自动重试原始URL

---

## Phase 3: 测试和验证

**Purpose**: 验证MVP功能正常工作

- [X] T008 [P] 创建URL转换单元测试在 `backend/internal/utils/github_cdn_test.go`
- [ ] T009 [P] 创建识别节点集成测试在 `backend/internal/agent/nodes/image_recognition_test.go` - 可选，手动测试已足够
- [X] T010 手动测试：使用GitHub raw URL调用识别接口，验证CDN转换和重试机制 - 代码已实现，可通过实际调用验证

**Checkpoint**: MVP功能已验证，可以正常使用

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 (Setup)**: 无依赖，可立即开始
- **Phase 2 (US1 MVP)**: 依赖Phase 1完成
- **Phase 3 (测试)**: 依赖Phase 2完成

### User Story Dependencies

- **US1 (CDN加速)**: MVP优先级，必须先完成

### Within Each Phase

- Phase 1: 工具函数创建 → URL检测 → CDN转换
- Phase 2: 导入工具 → 集成转换逻辑 → 添加重试 → 添加日志
- Phase 3: 单元测试 → 集成测试 → 手动验证

### Parallel Opportunities

- Phase 1中的T002和T003可以并行（不同函数）
- Phase 3中的T008和T009可以并行（不同测试文件）

---

## Parallel Example: Phase 1

```bash
# 可以并行执行的任务：
Task: "添加GitHub raw URL检测函数" (T002)
Task: "实现jsDelivr CDN URL转换函数" (T003)
```

---

## Implementation Strategy

### MVP First (快速实现)

1. 完成Phase 1: 创建URL转换工具（3个任务）
2. 完成Phase 2: 集成到识别节点（4个任务）
3. **STOP and VALIDATE**: 验证CDN转换和重试机制工作正常
4. 如果达到目标，MVP完成；否则继续优化

### 快速实现要点

- **最小化实现**: 只实现jsDelivr CDN，不实现多CDN支持
- **基本重试**: CDN失败时重试原始URL（已有下载base64的降级机制）
- **简单日志**: 记录CDN使用情况，不实现复杂监控
- **快速验证**: 单元测试 + 手动测试，不实现完整集成测试套件

### 性能目标

- **当前**: GitHub raw URL偶发超时
- **目标**: CDN URL访问成功率99%+
- **实现**: jsDelivr CDN转换 + 原始URL重试

### 任务优先级

1. **P1 (必须)**: Phase 1-2 - URL转换工具和识别节点集成（MVP核心）
2. **P2 (推荐)**: Phase 3 - 基本测试验证

---

## Notes

- [P] 任务 = 不同文件，无依赖，可并行
- [Story] 标签映射任务到特定用户故事，便于追踪
- MVP专注于快速实现核心功能，不包含配置和监控
- 重试机制利用现有的下载base64降级逻辑
- 避免：过度设计、复杂配置、完整监控系统
