# 数据模型：图片上传存储功能

**日期**: 2025-12-18  
**功能**: 003-image-upload

## 实体定义

### 1. 图片上传请求 (UploadRequest)

**用途**：前端向后端发送图片上传请求

**字段**：
- `imageData` (string, 必填): Base64 编码的图片数据（不含 data URL 前缀）
- `filename` (string, 可选): 文件名（如果不提供，后端自动生成）

**验证规则**：
- `imageData` 不能为空
- `imageData` 必须是有效的 base64 字符串
- `filename` 如果提供，必须符合安全命名规范（仅字母、数字、连字符、下划线）

**示例**：
```json
{
  "imageData": "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==",
  "filename": "user-upload-123.jpg"
}
```

### 2. 图片上传响应 (UploadResponse)

**用途**：后端返回图片上传结果

**字段**：
- `url` (string, 必填): 图片的访问 URL（GitHub raw URL 或 base64 data URL）
- `filename` (string, 必填): 实际存储的文件名
- `size` (int, 可选): 图片大小（字节）
- `uploadMethod` (string, 可选): 上传方式（"github" 或 "base64"）

**示例**：
```json
{
  "url": "https://raw.githubusercontent.com/owner/repo/main/images/1234567890-abc.jpg",
  "filename": "1234567890-abc.jpg",
  "size": 245678,
  "uploadMethod": "github"
}
```

### 3. 图片元数据 (ImageMetadata)

**用途**：存储图片的元数据信息（可选，用于未来扩展）

**字段**：
- `url` (string): 图片 URL
- `filename` (string): 文件名
- `size` (int): 文件大小（字节）
- `mimeType` (string): MIME 类型（如 "image/jpeg"）
- `width` (int, 可选): 图片宽度（像素）
- `height` (int, 可选): 图片高度（像素）
- `uploadedAt` (string): 上传时间（ISO 8601 格式）
- `uploadMethod` (string): 上传方式

**存储位置**：
- 当前版本：不持久化存储（仅返回给前端）
- 未来版本：可考虑存储在数据库或内存缓存中

## 状态流转

### 图片上传流程

```
前端选择图片
    ↓
前端压缩图片（Canvas API）
    ↓
前端转换为 base64
    ↓
前端发送上传请求（POST /api/upload/image）
    ↓
后端验证图片数据
    ↓
后端尝试上传到 GitHub
    ├─ 成功 → 返回 GitHub URL
    └─ 失败 → 回退到 base64，返回 data URL
    ↓
前端接收 URL
    ↓
前端使用 URL 调用识别 API（替代 base64）
```

## 数据验证规则

### 前端验证

1. **文件类型验证**：
   - 仅允许图片格式：JPEG、PNG、WebP、GIF
   - 通过文件扩展名和 MIME 类型双重验证

2. **文件大小验证**：
   - 原始文件大小：< 10MB
   - 压缩后目标大小：< 2MB

3. **Base64 验证**：
   - 必须是有效的 base64 字符串
   - 长度合理（不超过 15MB base64 字符串）

### 后端验证

1. **Base64 解码验证**：
   - 能够成功解码为二进制数据
   - 解码后的数据是有效的图片格式

2. **图片格式验证**：
   - 通过文件头（Magic Number）验证
   - 支持的格式：JPEG、PNG、WebP

3. **文件大小验证**：
   - 解码后的二进制数据 < 10MB

4. **文件名安全验证**：
   - 防止路径遍历（`../`、`..\\`）
   - 仅允许安全字符（字母、数字、连字符、下划线、点）

## 错误处理

### 错误类型

1. **验证错误** (400 Bad Request)：
   - 图片数据无效
   - 文件类型不支持
   - 文件大小超限

2. **GitHub API 错误** (502 Bad Gateway)：
   - GitHub API 调用失败
   - Token 无效或权限不足
   - 速率限制（Rate Limit）

3. **服务器错误** (500 Internal Server Error)：
   - 图片处理失败
   - 存储服务不可用

### 错误响应格式

```json
{
  "code": 400,
  "message": "图片数据无效",
  "detail": "Base64 字符串格式错误"
}
```

## 配置项

### 后端配置

**环境变量**：
- `GITHUB_TOKEN`: GitHub Personal Access Token（必填）
- `GITHUB_OWNER`: GitHub 用户名或组织名（必填）
- `GITHUB_REPO`: GitHub 仓库名（必填）
- `GITHUB_BRANCH`: GitHub 分支名（可选，默认 "main"）
- `GITHUB_PATH`: 图片存储路径（可选，默认 "images/"）
- `MAX_IMAGE_SIZE`: 最大图片大小，单位字节（可选，默认 10485760，即 10MB）

**配置文件** (`etc/explore.yaml`)：
```yaml
Upload:
  GitHubToken: ""  # 从环境变量读取
  GitHubOwner: ""  # 从环境变量读取
  GitHubRepo: ""  # 从环境变量读取
  GitHubBranch: "main"
  GitHubPath: "images/"
  MaxImageSize: 10485760  # 10MB
```

## 未来扩展

### 可能的扩展方向

1. **图片管理**：
   - 图片列表查询
   - 图片删除功能
   - 图片使用统计

2. **多存储后端**：
   - 支持其他存储服务（OSS、COS 等）
   - 存储后端可配置切换

3. **图片处理**：
   - 自动生成缩略图
   - 图片格式转换
   - 图片水印（可选）

4. **缓存机制**：
   - CDN 缓存
   - 本地缓存（减少重复上传）
