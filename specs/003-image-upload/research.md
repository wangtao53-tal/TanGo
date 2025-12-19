# 研究文档：图片上传存储功能

**日期**: 2025-12-18  
**功能**: 003-image-upload

## 研究目标

解决前端图片上传时遇到的 413 错误（Request Entity Too Large）问题，实现图片上传到 GitHub 仓库的功能，使用图片 URL 替代 base64 数据。

## 研究结果

### 1. GitHub API 图片上传方案

#### 决策：使用 GitHub Contents API

**方案选择**：
- ✅ **GitHub Contents API** (`PUT /repos/{owner}/{repo}/contents/{path}`)
- ❌ GitHub Releases API（不适合频繁上传）
- ❌ GitHub LFS（需要额外配置，成本较高）

**理由**：
- Contents API 简单直接，适合小文件（<100MB）
- 上传后可通过 `raw.githubusercontent.com` 直接访问，URL 稳定
- 支持 base64 编码上传，与现有流程兼容
- 认证方式简单（Personal Access Token）

**实现方式**：
```go
// 使用 GitHub REST API v3
PUT https://api.github.com/repos/{owner}/{repo}/contents/{path}
Headers:
  Authorization: token {GITHUB_TOKEN}
  Content-Type: application/json
Body:
  {
    "message": "Upload image",
    "content": "{base64_encoded_image}",
    "branch": "main"
  }
```

**限制和注意事项**：
- 文件大小限制：GitHub 建议单个文件 < 50MB，仓库总大小 < 100GB
- API 速率限制：认证请求 5000 次/小时，未认证 60 次/小时
- Base64 编码会增加约 33% 的文件大小
- 需要 Personal Access Token（repo 权限）

**URL 格式**：
- 上传后的图片 URL：`https://raw.githubusercontent.com/{owner}/{repo}/{branch}/{path}`
- URL 永久有效（除非文件被删除）

#### 替代方案考虑

**GitHub Releases API**：
- ❌ 不适合：主要用于发布版本，不适合频繁上传图片
- ❌ 操作复杂：需要创建 release，然后上传 asset

**GitHub LFS (Large File Storage)**：
- ❌ 不适合：需要额外配置 Git LFS，成本较高
- ❌ 操作复杂：需要 Git 命令操作，不适合 API 调用

**外部存储服务（OSS、CDN）**：
- ⚠️ 未来考虑：如果 GitHub 方案不满足需求，可考虑阿里云 OSS、腾讯云 COS 等
- 优势：专业存储服务，CDN 加速，更好的性能
- 劣势：需要额外成本，配置更复杂

### 2. 前端图片压缩策略

#### 决策：使用 Canvas API 进行客户端压缩

**压缩策略**：
1. **尺寸压缩**：限制最大宽度/高度（如 1920px）
2. **质量压缩**：JPEG 质量设置为 0.8（80%）
3. **格式转换**：统一转换为 JPEG 格式（体积更小）

**实现方式**：
```typescript
// 使用 Canvas API 压缩图片
function compressImage(file: File, maxWidth: number, maxHeight: number, quality: number): Promise<Blob> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = (e) => {
      const img = new Image();
      img.onload = () => {
        const canvas = document.createElement('canvas');
        // 计算新尺寸
        let width = img.width;
        let height = img.height;
        if (width > maxWidth || height > maxHeight) {
          const ratio = Math.min(maxWidth / width, maxHeight / height);
          width = width * ratio;
          height = height * ratio;
        }
        canvas.width = width;
        canvas.height = height;
        const ctx = canvas.getContext('2d');
        ctx?.drawImage(img, 0, 0, width, height);
        canvas.toBlob((blob) => {
          if (blob) resolve(blob);
          else reject(new Error('压缩失败'));
        }, 'image/jpeg', quality);
      };
      img.src = e.target?.result as string;
    };
    reader.onerror = reject;
    reader.readAsDataURL(file);
  });
}
```

**性能优化**：
- ✅ 移动端优先：先压缩尺寸再压缩质量，减少处理时间
- ✅ 异步处理：使用 Promise，不阻塞 UI
- ✅ 内存管理：及时释放 Canvas 和 Image 对象
- ✅ 进度提示：显示压缩进度（可选）

**压缩参数建议**：
- 最大尺寸：1920x1920px（移动端可降至 1280x1280px）
- 质量：0.8（80%），平衡文件大小和图片质量
- 目标文件大小：< 2MB（压缩后）

**移动端性能考虑**：
- ⚠️ 大图片处理可能较慢（>5MB 原始图片）
- ⚠️ 内存限制：移动浏览器可能限制 Canvas 大小
- ✅ 解决方案：先检查文件大小，超过 10MB 提示用户选择更小的图片

### 3. 降级方案

#### 决策：多层降级策略

**降级层级**：
1. **第一层**：尝试上传到 GitHub
2. **第二层**：GitHub 上传失败时，回退到 base64（现有方案）
3. **第三层**：base64 也失败时，提示用户重试或选择更小的图片

**错误处理**：
```typescript
async function uploadImage(file: File): Promise<string> {
  try {
    // 1. 压缩图片
    const compressed = await compressImage(file, 1920, 1920, 0.8);
    
    // 2. 尝试上传到 GitHub
    try {
      const url = await uploadToGitHub(compressed);
      return url;
    } catch (githubError) {
      console.warn('GitHub 上传失败，回退到 base64:', githubError);
      // 3. 回退到 base64
      const base64 = await fileToBase64(compressed);
      return base64;
    }
  } catch (error) {
    throw new Error('图片处理失败，请重试或选择更小的图片');
  }
}
```

**GitHub API 限流处理**：
- 检测 403 错误（Rate Limit）
- 返回友好的错误提示
- 自动回退到 base64 方案

### 4. 安全性

#### 决策：后端统一处理 GitHub API 调用

**安全措施**：
1. **Token 安全**：
   - GitHub token 仅在后端使用，不暴露给前端
   - Token 存储在环境变量中（`.env` 文件）
   - Token 权限最小化（仅 repo 权限）

2. **文件验证**：
   - 后端验证文件类型（仅允许图片格式）
   - 后端验证文件大小（< 10MB）
   - 文件名安全处理（防止路径遍历）

3. **内容安全**：
   - 不存储用户敏感信息
   - 图片仅用于识别，不用于其他用途
   - 考虑添加图片内容检查（可选，未来实现）

**实现方式**：
```go
// 后端统一处理 GitHub 上传
func (l *UploadLogic) UploadImage(req *types.UploadRequest) (*types.UploadResponse, error) {
    // 1. 验证文件类型和大小
    if err := validateImage(req.ImageData); err != nil {
        return nil, err
    }
    
    // 2. 调用 GitHub API
    url, err := l.githubStorage.Upload(req.ImageData, req.Filename)
    if err != nil {
        return nil, err
    }
    
    return &types.UploadResponse{Url: url}, nil
}
```

### 5. go-zero 文件上传处理

#### 决策：使用 multipart/form-data 或 base64 JSON

**方案选择**：
- ✅ **Base64 JSON**（推荐）：与现有 API 风格一致
- ⚠️ Multipart/form-data：需要修改现有 API 结构

**理由**：
- 现有 API 使用 JSON 格式，保持一致性
- Base64 编码在后端处理，前端只需发送字符串
- 简化错误处理和降级逻辑

**实现方式**：
```go
// API 定义
type UploadRequest {
    ImageData string `json:"imageData"` // base64 编码的图片
    Filename  string `json:"filename,optional"` // 可选文件名
}

// Handler 处理
func UploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req types.UploadRequest
        if err := httpx.Parse(r, &req); err != nil {
            httpx.Error(w, err)
            return
        }
        // ... 处理逻辑
    }
}
```

## 技术决策总结

| 决策项 | 选择 | 理由 |
|--------|------|------|
| 存储方案 | GitHub Contents API | 简单、免费、URL 稳定 |
| 压缩方案 | Canvas API 客户端压缩 | 减少上传大小，提升性能 |
| 上传方式 | Base64 JSON | 与现有 API 风格一致 |
| 降级策略 | GitHub → Base64 | 确保功能可用性 |
| Token 管理 | 后端环境变量 | 安全性高，不暴露给前端 |

## 待解决问题

1. **GitHub 仓库选择**：
   - 使用现有仓库还是创建新仓库？
   - 建议：创建专门的图片存储仓库（如 `tango-images`）

2. **文件命名策略**：
   - 如何生成唯一文件名？
   - 建议：`{timestamp}-{random}.jpg` 或 `{hash}.jpg`

3. **图片组织方式**：
   - 是否需要按日期/用户组织？
   - 建议：MVP 版本简单扁平结构，未来可按需优化

4. **CDN 加速**：
   - `raw.githubusercontent.com` 是否有 CDN？
   - 研究：GitHub 的 raw 文件已通过 CDN 加速，性能可接受

## 下一步行动

1. ✅ 研究完成，生成数据模型和 API 合约
2. 实现前端图片压缩功能
3. 实现后端 GitHub 上传功能
4. 实现降级方案和错误处理
5. 测试和优化
