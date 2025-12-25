# 快速开始：图片上传存储功能

**日期**: 2025-12-18  
**功能**: 003-image-upload

## 功能概述

实现图片上传到 GitHub 仓库的功能，解决前端 base64 图片过大导致的 413 错误问题。上传后的图片 URL 可以替代 base64 数据，减少请求体大小。

## 前置要求

### 1. GitHub 仓库准备

1. **创建 GitHub 仓库**（如果还没有）：
   - 仓库名：建议使用 `tango-images` 或类似名称
   - 仓库类型：Public 或 Private（Private 需要 token 有相应权限）
   - 初始化：可以创建一个空的 `main` 分支

2. **创建 GitHub Personal Access Token**：
   - 访问：https://github.com/settings/tokens
   - 点击 "Generate new token (classic)"
   - 权限选择：至少需要 `repo` 权限（如果是私有仓库）
   - 复制生成的 token（只显示一次，请妥善保存）

### 2. 环境配置

**后端环境变量**（`.env` 文件）：
```bash
# GitHub 配置
GITHUB_TOKEN=ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
GITHUB_OWNER=your-username
GITHUB_REPO=tango-images
GITHUB_BRANCH=main
GITHUB_PATH=images/

# 可选配置
MAX_IMAGE_SIZE=10485760  # 10MB
```

**配置文件**（`backend/etc/explore.yaml`）：
```yaml
Upload:
  GitHubToken: ""  # 从环境变量 GITHUB_TOKEN 读取
  GitHubOwner: ""  # 从环境变量 GITHUB_OWNER 读取
  GitHubRepo: ""   # 从环境变量 GITHUB_REPO 读取
  GitHubBranch: "main"
  GitHubPath: "images/"
  MaxImageSize: 10485760
```

## 安装步骤

### 1. 后端实现

1. **生成 API 代码**：
```bash
cd backend
goctl api go -api api/explore.api -dir . -style gozero
```

2. **实现上传逻辑**：
   - 创建 `internal/handler/uploadhandler.go`
   - 创建 `internal/logic/uploadlogic.go`
   - 创建 `internal/storage/github.go`

3. **安装依赖**（如果需要）：
```bash
go get github.com/google/go-github/v50/github
go get golang.org/x/oauth2
```

### 2. 前端实现

1. **更新图片处理工具**（`frontend/src/utils/image.ts`）：
   - 添加图片压缩函数
   - 添加 base64 提取函数

2. **更新 API 服务**（`frontend/src/services/api.ts`）：
   - 添加图片上传接口

3. **更新页面组件**：
   - `Capture.tsx`: 上传图片后再调用识别 API
   - `Result.tsx`: 上传图片后再发送消息

## 使用流程

### 用户操作流程

1. **用户选择图片**：
   - 点击拍照按钮或选择文件
   - 前端自动压缩图片（Canvas API）

2. **图片上传**：
   - 前端调用 `/api/upload/image` 上传图片
   - 后端尝试上传到 GitHub
   - 如果失败，自动回退到 base64

3. **使用图片 URL**：
   - 前端获取图片 URL
   - 使用 URL 调用识别 API（替代 base64）
   - 识别 API 支持 URL 输入（已在之前实现）

### 开发者测试流程

1. **测试图片上传**：
```bash
# 准备测试图片（base64）
IMAGE_BASE64="iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg=="

# 调用上传接口
curl -X POST http://localhost:8877/api/upload/image \
  -H "Content-Type: application/json" \
  -d "{\"imageData\": \"$IMAGE_BASE64\", \"filename\": \"test.jpg\"}"
```

2. **验证上传结果**：
   - 检查返回的 URL 是否可访问
   - 检查 GitHub 仓库中是否有新文件

3. **测试识别接口**（使用上传后的 URL）：
```bash
curl -X POST http://localhost:8877/api/explore/identify \
  -H "Content-Type: application/json" \
  -d "{
    \"image\": \"https://raw.githubusercontent.com/owner/repo/main/images/test.jpg\",
    \"age\": 8
  }"
```

## 故障排查

### 常见问题

1. **413 错误仍然出现**：
   - 检查图片是否已压缩
   - 检查压缩后的 base64 大小
   - 检查后端请求体大小限制配置

2. **GitHub 上传失败**：
   - 检查 GitHub token 是否有效
   - 检查 token 权限是否足够
   - 检查仓库是否存在且可访问
   - 检查 GitHub API 速率限制

3. **图片 URL 无法访问**：
   - 检查 GitHub 仓库是否为 Public
   - 检查文件路径是否正确
   - 检查分支名称是否正确

4. **识别接口不支持 URL**：
   - 确认已实现 URL 下载功能（已在之前实现）
   - 检查 URL 格式是否正确

## 性能优化建议

1. **前端压缩**：
   - 压缩质量：0.8（80%）
   - 最大尺寸：1920x1920px
   - 目标大小：< 2MB

2. **后端优化**：
   - 使用并发上传（如果需要）
   - 实现上传缓存（避免重复上传）
   - 监控 GitHub API 速率限制

3. **CDN 加速**：
   - GitHub raw 文件已通过 CDN 加速
   - 如需更快的访问速度，可考虑使用 jsDelivr CDN

## 下一步

1. ✅ 实现前端图片压缩功能
2. ✅ 实现后端 GitHub 上传功能
3. ✅ 实现降级方案和错误处理
4. ✅ 测试和优化
5. ⚠️ 监控 GitHub API 使用情况
6. ⚠️ 考虑未来迁移到专业存储服务（如 OSS、COS）
