# 任务清单：图片上传存储功能

**输入**: 设计文档来自 `/specs/003-image-upload/`
**前置条件**: plan.md ✅, research.md ✅, data-model.md ✅, contracts/ ✅

**GitHub 仓库配置**:
- 仓库地址: https://github.com/wangtao53-tal/image
- Owner: wangtao53-tal
- Repo: image
- Branch: main (默认)

## 格式说明

- **[P]**: 可以并行执行（不同文件，无依赖关系）
- **[US]**: 用户故事标识（US1 = 图片上传核心功能）

## Phase 1: Setup（项目初始化）

**目的**: 项目初始化和基础结构

- [x] T001 [P] 更新后端 API 定义文件 `backend/api/explore.api`，添加图片上传接口定义
- [x] T002 [P] 更新后端配置文件 `backend/etc/explore.yaml`，添加 GitHub 上传配置项
- [x] T003 [P] 更新后端配置结构 `backend/internal/config/config.go`，添加 Upload 配置结构体

**检查点**: 配置和 API 定义完成

---

## Phase 2: Foundational（基础设施）

**目的**: 核心基础设施，必须在任何用户故事实现之前完成

**⚠️ 关键**: 在完成此阶段之前，不能开始任何用户故事工作

- [x] T004 创建 GitHub 存储实现 `backend/internal/storage/github.go`
  - 实现 GitHub API 客户端封装
  - 实现图片上传到 GitHub 的功能
  - 实现错误处理和重试逻辑
- [x] T005 [P] 创建图片验证工具 `backend/internal/utils/image.go`
  - 实现 base64 解码验证
  - 实现图片格式验证（通过文件头 Magic Number）
  - 实现文件大小验证
  - 实现文件名安全验证
- [x] T006 [P] 更新错误定义 `backend/internal/utils/errors.go`
  - 添加图片上传相关错误定义

**检查点**: 基础设施就绪，可以开始用户故事实现

---

## Phase 3: User Story 1 - 图片上传核心功能（Priority: P1）🎯 MVP

**目标**: 实现图片上传到 GitHub 仓库的核心功能，解决 413 错误问题

**独立测试**: 可以通过 CURL 或前端调用上传接口，成功上传图片并返回 GitHub URL

### 后端实现

- [x] T007 [US1] 生成 API 代码：运行 `goctl api go -api backend/api/explore.api -dir backend -style gozero`
- [x] T008 [US1] 实现上传处理器 `backend/internal/handler/uploadhandler.go`
  - 解析上传请求
  - 调用业务逻辑
  - 返回响应
- [x] T009 [US1] 实现上传业务逻辑 `backend/internal/logic/uploadlogic.go`
  - 验证图片数据（调用 T005 的工具函数）
  - 生成唯一文件名（时间戳 + 随机字符串）
  - 调用 GitHub 存储上传图片
  - 实现降级逻辑（GitHub 失败时返回 base64 data URL）
  - 错误处理和日志记录
- [x] T010 [US1] 更新路由注册 `backend/internal/handler/routes.go`
  - 注册图片上传路由 `/api/upload/image`
- [x] T011 [US1] 更新服务上下文 `backend/internal/svc/servicecontext.go`
  - 初始化 GitHub 存储客户端

### 前端实现

- [x] T012 [P] [US1] 更新图片处理工具 `frontend/src/utils/image.ts`
  - 实现 `compressImage` 函数（Canvas API 压缩）
    - 参数：file, maxWidth, maxHeight, quality
    - 返回：Promise<Blob>
    - 支持尺寸压缩和质量压缩
    - 统一转换为 JPEG 格式
  - 实现 `extractBase64Data` 函数（如果还没有）
    - 从 data URL 中提取纯 base64 字符串
- [x] T013 [P] [US1] 更新 API 服务 `frontend/src/services/api.ts`
  - 添加 `uploadImage` 函数
    - 调用 `/api/upload/image` 接口
    - 处理响应和错误
- [x] T014 [US1] 更新类型定义 `frontend/src/types/api.ts`
  - 添加 `UploadRequest` 接口
  - 添加 `UploadResponse` 接口

### 前端集成

- [x] T015 [US1] 更新拍照页面 `frontend/src/pages/Capture.tsx`
  - 在 `handleImageSelect` 中：
    - 先压缩图片（调用 `compressImage`）
    - 转换为 base64
    - 调用上传接口（`uploadImage`）
    - 使用返回的 URL 调用识别 API（替代 base64）
    - 添加上传进度提示
    - 处理上传失败（自动降级到 base64）
- [x] T016 [US1] 更新对话页面 `frontend/src/pages/Result.tsx`
  - 在 `handleImageSelect` 中：
    - 先压缩图片
    - 调用上传接口
    - 使用返回的 URL 发送消息
    - 处理上传失败（自动降级）

### 多语言支持

- [ ] T017 [P] [US1] 更新国际化文件 `frontend/src/i18n/locales/zh.ts`
  - 添加图片上传相关文本（上传中、上传成功、上传失败等）
- [ ] T018 [P] [US1] 更新国际化文件 `frontend/src/i18n/locales/en.ts`
  - 添加图片上传相关英文文本

**检查点**: 此时，用户故事 1 应该完全功能化并可独立测试

---

## Phase 4: 错误处理和优化

**目的**: 完善错误处理和用户体验优化

- [ ] T019 [P] 实现 GitHub API 速率限制处理
  - 检测 403 错误（Rate Limit）
  - 返回友好的错误提示
  - 自动降级到 base64
- [ ] T020 [P] 实现上传进度提示组件（可选）
  - 创建 `frontend/src/components/common/UploadProgress.tsx`
  - 显示上传进度百分比
  - 显示上传状态（上传中、成功、失败）
- [ ] T021 优化图片压缩性能
  - 添加压缩超时处理
  - 优化大图片处理（分块处理）
  - 添加压缩进度回调（可选）
- [ ] T022 添加日志和监控
  - 后端：添加详细的上传日志
  - 前端：添加错误日志和性能监控
  - 记录上传成功率、平均上传时间等指标

---

## Phase 5: 测试和验证

**目的**: 确保功能正常工作

- [ ] T023 测试后端上传接口
  - 使用 CURL 测试上传功能
  - 测试各种图片格式（JPEG、PNG、WebP）
  - 测试大图片（接近 10MB）
  - 测试降级逻辑（模拟 GitHub API 失败）
- [ ] T024 测试前端集成
  - 测试拍照页面上传流程
  - 测试对话页面上传流程
  - 测试压缩功能
  - 测试错误处理
- [ ] T025 测试移动端
  - 在移动设备上测试上传功能
  - 验证压缩性能
  - 验证用户体验
- [ ] T026 验证 GitHub 仓库
  - 确认图片已成功上传到 https://github.com/wangtao53-tal/image
  - 确认图片 URL 可访问
  - 验证图片 URL 格式正确

---

## Phase 6: 文档和清理

**目的**: 完善文档和代码清理

- [ ] T027 [P] 更新 README 文档
  - 更新 `backend/README.md`，添加图片上传功能说明
  - 更新 `frontend/README.md`，添加图片压缩说明
- [ ] T028 [P] 更新环境变量示例
  - 更新 `.env.example`，添加 GitHub 配置项
- [ ] T029 代码清理和重构
  - 检查代码规范
  - 优化代码结构
  - 添加必要的注释
- [ ] T030 运行 quickstart.md 验证
  - 按照 `specs/003-image-upload/quickstart.md` 验证所有步骤
  - 确认配置正确
  - 确认功能正常

---

## 依赖关系和执行顺序

### 阶段依赖

- **Setup (Phase 1)**: 无依赖，可立即开始
- **Foundational (Phase 2)**: 依赖 Setup 完成，阻塞所有用户故事
- **User Story 1 (Phase 3)**: 依赖 Foundational 完成
- **错误处理和优化 (Phase 4)**: 依赖 User Story 1 完成
- **测试和验证 (Phase 5)**: 依赖 Phase 3 和 Phase 4 完成
- **文档和清理 (Phase 6)**: 依赖所有功能完成

### 用户故事依赖

- **User Story 1 (P1)**: 可在 Foundational (Phase 2) 完成后开始，不依赖其他故事

### 每个用户故事内部

- 配置和 API 定义 → 基础设施 → 后端实现 → 前端实现 → 前端集成 → 测试

### 并行机会

- Phase 1 中所有标记 [P] 的任务可以并行执行
- Phase 2 中所有标记 [P] 的任务可以并行执行（在 Phase 2 内）
- Phase 3 中：
  - T012, T013, T014 可以并行（前端工具和类型定义）
  - T017, T018 可以并行（多语言文件）
- Phase 4 中所有标记 [P] 的任务可以并行
- Phase 6 中所有标记 [P] 的任务可以并行

---

## 实施策略

### MVP 优先（仅 User Story 1）

1. 完成 Phase 1: Setup
2. 完成 Phase 2: Foundational（关键 - 阻塞所有故事）
3. 完成 Phase 3: User Story 1
4. **停止并验证**: 独立测试 User Story 1
5. 如果就绪，部署/演示

### 增量交付

1. Setup + Foundational → 基础设施就绪
2. 添加 User Story 1 → 独立测试 → 部署/演示（MVP！）
3. 添加错误处理和优化 → 测试 → 部署/演示
4. 添加测试和文档 → 完成

### 并行团队策略

如果有多个开发者：

1. 团队一起完成 Setup + Foundational
2. Foundational 完成后：
   - 开发者 A: 后端实现（T007-T011）
   - 开发者 B: 前端工具和 API（T012-T014）
   - 开发者 C: 前端集成（T015-T016）
3. 完成后集成和测试

---

## 环境配置检查清单

在开始实现前，请确认：

- [ ] GitHub Personal Access Token 已创建（需要 repo 权限）
- [ ] 环境变量已配置：
  - `GITHUB_TOKEN=ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx`
  - `GITHUB_OWNER=wangtao53-tal`
  - `GITHUB_REPO=image`
  - `GITHUB_BRANCH=main`（可选，默认 main）
  - `GITHUB_PATH=images/`（可选，默认 images/）
- [ ] GitHub 仓库 https://github.com/wangtao53-tal/image 已创建且可访问
- [ ] Token 有写入仓库的权限

---

## 注意事项

- [P] 任务 = 不同文件，无依赖关系
- [US1] 标签映射任务到特定用户故事，便于追溯
- 每个用户故事应该可以独立完成和测试
- 提交前验证每个任务
- 在每个检查点停止以验证故事独立性
- 避免：模糊任务、同一文件冲突、破坏独立性的跨故事依赖

---

## 快速开始命令

```bash
# 1. 配置环境变量（.env 文件）
GITHUB_TOKEN=ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
GITHUB_OWNER=wangtao53-tal
GITHUB_REPO=image
GITHUB_BRANCH=main
GITHUB_PATH=images/

# 2. 生成后端 API 代码（T007）
cd backend
goctl api go -api api/explore.api -dir . -style gozero

# 3. 测试上传接口（T023）
curl -X POST http://localhost:8877/api/upload/image \
  -H "Content-Type: application/json" \
  -d '{
    "imageData": "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==",
    "filename": "test.jpg"
  }'
```
