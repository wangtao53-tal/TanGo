# 实现计划：GitHub 图片加速优化

**分支**: `006-github-image-acceleration` | **日期**: 2025-12-20 | **规范**: [spec.md](./spec.md)
**输入**: GitHub raw URL 偶发访问超时问题

## Summary

解决 GitHub raw URL 偶发访问超时问题，通过以下方案优化：
1. 使用 GitHub 图片 CDN 加速服务（jsDelivr、fastgit 等）
2. 添加 URL 转换机制，自动将 raw.githubusercontent.com URL 转换为加速 URL
3. 添加重试机制和备选 URL 支持
4. 优化错误处理，提供更好的降级策略

## Technical Context

**Language/Version**: 
- 后端：Go 1.25.3 (darwin/arm64)

**Primary Dependencies**: 
- 后端：go-zero v1.9.3, eino (字节云原生AI框架)
- AI模型：图像识别模型（通过eino框架调用）
- CDN服务：jsDelivr (主要), FastGit (备选)

**Storage**: 
- GitHub 存储：`backend/internal/storage/github.go`
- 图片识别节点：`backend/internal/agent/nodes/image_recognition.go`

**Testing**: 
- 后端：Go testing package + 网络测试

**Target Platform**: 
- Web API服务

**Project Type**: 后端API服务优化

**Performance Goals**: 
- GitHub 图片访问成功率：从偶发失败提升到 99%+
- 图片加载时间：减少 50%以上（通过CDN加速）
- 超时重试：自动重试失败请求
- CDN响应时间：<2秒（P95）

**Constraints**: 
- 必须保持API接口兼容性
- 必须保证功能正确性
- 不能破坏现有错误处理机制
- CDN服务需要稳定可靠
- 必须支持降级策略（CDN失败时回退到原始URL）

**Scale/Scope**: 
- 影响范围：GitHub 图片URL生成和识别节点URL处理
- 预计影响：所有使用GitHub图片的识别请求

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**规范检查项**（基于 `.specify/memory/constitution.md`）：

- [X] **原则一：中文优先规范** - 所有文档和生成内容使用中文
- [X] **原则二：K12 教育游戏化设计规范** - 不适用（后端优化）
- [X] **原则三：可发布应用规范** - 实现达到生产级标准，可正常运行和发布
- [X] **原则四：多语言和年级设置规范** - 不适用（后端优化）
- [X] **原则五：AI优先（模型优先）规范** - 优化AI模型图片访问，支持流式返回
- [X] **原则六：移动端优先规范** - 不适用（后端优化）
- [X] **原则七：用户体验流程规范** - 提升图片识别成功率，改善用户体验

**合规性说明**：本优化方案完全符合项目规范，提升系统稳定性和用户体验。

## Architecture

### Current Architecture

```
GitHubStorage.Upload()
  ↓
生成 raw.githubusercontent.com URL
  ↓
返回给前端
  ↓
前端调用 /api/explore/identify
  ↓
ImageRecognitionNode.Execute()
  ↓
直接使用 raw URL（可能超时）
  ↓
模型下载图片失败
```

### Optimized Architecture

```
GitHubStorage.Upload()
  ↓
生成 raw.githubusercontent.com URL
  ↓
添加URL转换功能（可选）
  ↓
返回给前端（原始URL或加速URL）
  ↓
前端调用 /api/explore/identify
  ↓
ImageRecognitionNode.Execute()
  ↓
检测GitHub raw URL → 转换为CDN URL
  ↓
模型使用CDN URL（更稳定快速）
  ↓
如果CDN失败 → 重试原始URL
  ↓
如果仍失败 → 下载并转换为base64
```

## Implementation Approach

### Phase 0: Research
- 研究GitHub图片CDN加速方案
- 评估各CDN服务的稳定性和性能
- 确定最佳实践

### Phase 1: URL转换工具
- 创建URL转换工具函数
- 支持多种CDN服务（jsDelivr、fastgit等）
- 添加配置选项

### Phase 2: 集成到识别节点
- 在ImageRecognitionNode中添加URL转换逻辑
- 检测GitHub raw URL并自动转换
- 添加重试机制

### Phase 3: 可选：存储层优化
- 在GitHubStorage中可选生成CDN URL
- 配置选项控制是否使用CDN

### Phase 4: 测试和验证
- 测试CDN URL访问稳定性
- 验证重试机制
- 性能对比测试

## Risk Assessment

**Low Risk**:
- URL转换逻辑（不改变核心功能）
- 添加配置选项（向后兼容）

**Medium Risk**:
- CDN服务稳定性（需要监控）
- 重试机制（可能增加延迟）

**High Risk**:
- CDN服务不可用时的降级策略
- 多个CDN服务的维护成本

## Success Criteria

1. GitHub 图片访问成功率提升到 99%+
2. 图片加载时间减少 50%以上
3. 超时错误显著减少
4. API接口兼容性100%保持
5. 功能正确性100%保持
