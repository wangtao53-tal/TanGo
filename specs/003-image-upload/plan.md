# 实现计划：图片上传存储功能

**分支**: `003-image-upload` | **日期**: 2025-12-18 | **规范**: [spec.md](./spec.md)
**输入**: 功能需求来自用户反馈

## Summary

解决前端图片上传时遇到的 413 错误（Request Entity Too Large）问题。当前实现将图片转换为 base64 后直接发送给后端，导致请求体过大。需要实现图片上传存储功能，支持将图片上传到 GitHub 仓库，然后使用图片 URL 替代 base64 数据，减少请求体大小，提升性能和用户体验。

## Technical Context

**Language/Version**: 
- 前端：TypeScript/JavaScript (ES2020+), React 18
- 后端：Go 1.25.3 (darwin/arm64)

**Primary Dependencies**: 
- 前端：React 18, Vite, Axios
- 后端：go-zero v1.9.3
- 图片存储：GitHub API (REST API v3) 或 GitHub Actions + Git 操作
- 图片处理：前端图片压缩（Canvas API）

**Storage**: 
- 图片存储：GitHub 仓库（通过 GitHub API 或 Git 操作）
- 元数据存储：后端内存缓存（可选，用于临时存储上传状态）

**Testing**: 
- 前端：Vitest + React Testing Library
- 后端：Go testing package

**Target Platform**: 
- Web H5（移动端优先）
- 现代浏览器（支持 File API、Canvas API）

**Project Type**: Web应用（前后端分离）

**Performance Goals**: 
- 图片上传响应时间≤5秒（90%请求）
- 支持图片大小：最大 10MB（上传前压缩到 2MB 以内）
- 图片压缩时间≤2秒（前端处理）

**Constraints**: 
- GitHub 仓库大小限制（建议使用 GitHub Releases 或 LFS）
- GitHub API 速率限制（需要认证 token）
- 必须保护用户隐私（图片不包含敏感信息）
- 需要 GitHub token 配置（环境变量）

**Scale/Scope**: 
- MVP版本：支持单张图片上传
- 目标用户：K12学生（3-12岁）
- 预计并发：50-100用户
- 图片存储：GitHub 仓库（建议使用专门的图片存储仓库）

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**规范检查项**（基于 `.specify/memory/constitution.md`）：

- [x] **原则一：中文优先规范** - 所有文档和生成内容必须使用中文（除非技术限制）
  - ✅ 前端界面提示文本使用中文
  - ✅ 错误提示使用中文
  - ✅ 代码注释优先使用中文

- [x] **原则二：K12 教育游戏化设计规范** - 设计必须符合儿童友好性、游戏化元素、玩中学理念
  - ✅ 上传过程提供友好的加载提示
  - ✅ 上传失败时提供清晰的错误提示
  - ✅ 不影响现有的拍照探索流程

- [x] **原则三：可发布应用规范** - 实现必须达到生产级标准，可正常运行和发布
  - ✅ 错误处理：完善的错误提示和降级方案（上传失败时回退到 base64）
  - ✅ 性能优化：前端图片压缩，减少上传时间
  - ✅ 安全性：GitHub token 安全存储（环境变量，不暴露给前端）

- [x] **原则四：多语言和年级设置规范** - 支持中英文设置和K12年级设置，默认中文
  - ✅ 上传功能提示文本支持多语言（通过现有 i18n 系统）

- [x] **原则五：AI优先（模型优先）规范** - 所有AI相关功能统一通过后端调用模型处理
  - ✅ 图片上传通过后端统一处理，前端不直接调用 GitHub API
  - ✅ 后端统一管理 GitHub token 和上传逻辑

- [x] **原则六：移动端优先规范** - 确保移动端交互的完整性
  - ✅ 上传功能在移动端正常工作
  - ✅ 图片压缩在移动端性能可接受
  - ✅ 上传进度提示清晰可见

**合规性说明**：所有设计均符合项目规范要求，无违反项。

## Project Structure

### Documentation (this feature)

```text
specs/003-image-upload/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
frontend/                    # 前端工程
├── src/
│   ├── utils/
│   │   └── image.ts         # 图片处理工具（新增：图片压缩、格式转换）
│   ├── services/
│   │   └── api.ts           # API服务（新增：图片上传接口）
│   └── components/
│       └── common/          # 通用组件（可选：上传进度组件）

backend/                     # 后端工程
├── internal/
│   ├── handler/
│   │   └── uploadhandler.go # 图片上传处理器（新增）
│   ├── logic/
│   │   └── uploadlogic.go   # 图片上传业务逻辑（新增）
│   ├── storage/
│   │   └── github.go        # GitHub 存储实现（新增）
│   └── config/
│       └── config.go        # 配置（新增：GitHub 相关配置）
├── api/
│   └── explore.api          # API定义（新增：图片上传接口）
└── etc/
    └── explore.yaml         # 配置文件（新增：GitHub 配置项）
```

**Structure Decision**: 采用前后端分离架构，前端负责图片压缩和上传请求，后端负责与 GitHub API 交互和图片存储。GitHub token 仅在后端使用，确保安全性。

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

无违反项。

## Phase 0: Research & Outline

### 需要研究的问题

1. **GitHub 图片存储方案**
   - GitHub API 上传文件的最佳实践
   - GitHub Releases vs GitHub LFS vs 直接提交到仓库
   - GitHub API 速率限制和认证方式
   - 图片 URL 的永久性（CDN 加速）

2. **图片压缩策略**
   - 前端图片压缩算法（Canvas API）
   - 压缩质量与文件大小的平衡
   - 支持的图片格式（JPEG、PNG、WebP）
   - 移动端压缩性能

3. **降级方案**
   - 上传失败时的回退策略（base64）
   - 网络错误处理
   - GitHub API 限流处理

4. **安全性**
   - GitHub token 的安全存储
   - 图片内容安全检查（可选）
   - 防止恶意上传

### 研究任务

- [ ] 研究 GitHub API 上传文件的方案（Releases API、Contents API、Git LFS）
- [ ] 研究前端图片压缩的最佳实践和性能优化
- [ ] 研究 GitHub CDN（raw.githubusercontent.com）的访问速度和稳定性
- [ ] 研究 go-zero 框架中文件上传的处理方式
- [ ] 研究图片上传的降级方案和错误处理

## Phase 1: Design & Contracts

### 数据模型

✅ **已完成** - 详见 `data-model.md`

**核心实体**：
- `UploadRequest`: 图片上传请求（base64 数据 + 可选文件名）
- `UploadResponse`: 图片上传响应（URL + 元数据）
- `ImageMetadata`: 图片元数据（未来扩展）

**关键设计决策**：
- 使用 base64 JSON 格式（与现有 API 风格一致）
- 支持自动降级（GitHub → base64）
- 文件名自动生成（时间戳 + 随机字符串）

### API 合约

✅ **已完成** - 详见 `contracts/upload.api` 和 `contracts/README.md`

**接口定义**：
- `POST /api/upload/image`: 图片上传接口
- 请求：`UploadRequest`（base64 图片数据）
- 响应：`UploadResponse`（图片 URL）

**错误处理**：
- 400: 请求参数错误
- 401: 认证失败
- 403: 权限不足或速率限制
- 500: 服务器内部错误
- 502: GitHub API 调用失败

### 快速开始指南

✅ **已完成** - 详见 `quickstart.md`

**包含内容**：
- GitHub 仓库准备步骤
- 环境变量配置
- 安装和实现步骤
- 使用流程和测试方法
- 故障排查指南

## Phase 0 & Phase 1 完成总结

✅ **Phase 0 研究完成**：
- GitHub API 上传方案研究
- 前端图片压缩策略研究
- 降级方案设计
- 安全性考虑

✅ **Phase 1 设计完成**：
- 数据模型定义
- API 合约设计
- 快速开始指南

## Next Steps

1. ✅ Phase 0 研究完成
2. ✅ Phase 1 设计完成
3. ⏭️ 执行 Phase 2：生成任务清单（`/speckit.tasks` 命令）
4. ⏭️ 实现前端图片压缩功能
5. ⏭️ 实现后端 GitHub 上传功能
6. ⏭️ 实现降级方案和错误处理
7. ⏭️ 测试和优化
