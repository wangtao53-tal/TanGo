# API 合约文档

## 图片上传接口

### POST /api/upload/image

上传图片到 GitHub 仓库，返回图片 URL。

**请求**：
```json
{
  "imageData": "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==",
  "filename": "user-upload-123.jpg"
}
```

**响应（成功）**：
```json
{
  "url": "https://raw.githubusercontent.com/owner/repo/main/images/1234567890-abc.jpg",
  "filename": "1234567890-abc.jpg",
  "size": 245678,
  "uploadMethod": "github"
}
```

**响应（失败，回退到 base64）**：
```json
{
  "url": "data:image/jpeg;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==",
  "filename": "fallback-1234567890.jpg",
  "size": 245678,
  "uploadMethod": "base64"
}
```

**错误响应**：
```json
{
  "code": 400,
  "message": "图片数据无效",
  "detail": "Base64 字符串格式错误"
}
```

## 错误码

- `400`: 请求参数错误（图片数据无效、文件类型不支持、文件大小超限）
- `401`: 认证失败（GitHub token 无效）
- `403`: 权限不足或速率限制（GitHub API rate limit）
- `500`: 服务器内部错误
- `502`: GitHub API 调用失败

## 使用示例

### 前端调用示例

```typescript
// 1. 压缩图片
const compressed = await compressImage(file, 1920, 1920, 0.8);

// 2. 转换为 base64
const base64 = await fileToBase64(compressed);
const imageData = extractBase64Data(base64); // 移除 data URL 前缀

// 3. 上传图片
const response = await fetch('/api/upload/image', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    imageData: imageData,
    filename: 'user-upload.jpg',
  }),
});

const result = await response.json();

// 4. 使用返回的 URL 调用识别 API
const identifyResult = await identifyImage({
  image: result.url, // 使用 URL 替代 base64
  age: 8,
});
```

### CURL 示例

```bash
curl -X POST http://localhost:8877/api/upload/image \
  -H "Content-Type: application/json" \
  -d '{
    "imageData": "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==",
    "filename": "test.jpg"
  }'
```
