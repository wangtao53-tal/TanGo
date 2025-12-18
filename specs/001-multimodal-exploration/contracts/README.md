# API 合约文档

## 概述

本文档定义了 TanGo 探索服务的所有 API 接口。使用 go-zero 的 API 定义格式（.api 文件）。

## API 列表

### 1. 图像识别

**接口**: `POST /api/explore/identify`

**描述**: 识别拍照对象，返回对象名称、类别和置信度

**请求**:
```json
{
  "image": "data:image/jpeg;base64,...",
  "age": 8
}
```

**响应**:
```json
{
  "objectName": "银杏",
  "objectCategory": "自然类",
  "confidence": 0.95,
  "keywords": ["植物", "树木", "秋天"]
}
```

**错误码**:
- 400: 请求参数错误（图片格式不支持、base64格式错误等）
- 500: 服务器内部错误（AI模型调用失败等）

### 2. 生成知识卡片

**接口**: `POST /api/explore/generate-cards`

**描述**: 根据识别结果和年龄，生成三张知识卡片

**请求**:
```json
{
  "objectName": "银杏",
  "objectCategory": "自然类",
  "age": 8,
  "keywords": ["植物", "树木", "秋天"]
}
```

**响应**:
```json
{
  "cards": [
    {
      "type": "science",
      "title": "银杏",
      "content": {
        "name": "银杏",
        "explanation": "银杏是非常古老的植物...",
        "facts": [...],
        "funFact": "..."
      }
    },
    {
      "type": "poetry",
      "title": "古人怎么看银杏",
      "content": {...}
    },
    {
      "type": "english",
      "title": "用英语说银杏",
      "content": {...}
    }
  ]
}
```

**错误码**:
- 400: 请求参数错误（缺少必填字段等）
- 500: 服务器内部错误（AI模型调用失败等）

### 3. 创建分享链接

**接口**: `POST /api/share/create`

**描述**: 创建分享链接，用于家长端查看

**请求**:
```json
{
  "explorationRecords": [...],
  "collectedCards": [...]
}
```

**响应**:
```json
{
  "shareId": "share-123e4567-e89b-12d3-a456-426614174000",
  "shareUrl": "https://tango.example.com/share/share-123e4567-e89b-12d3-a456-426614174000",
  "expiresAt": "2025-12-25T10:00:00Z"
}
```

**错误码**:
- 400: 请求参数错误（数据格式错误等）
- 500: 服务器内部错误

### 4. 获取分享数据

**接口**: `GET /api/share/:shareId`

**描述**: 通过分享链接ID获取探索数据

**响应**:
```json
{
  "explorationRecords": [...],
  "collectedCards": [...],
  "createdAt": "2025-12-18T10:00:00Z",
  "expiresAt": "2025-12-25T10:00:00Z"
}
```

**错误码**:
- 404: 分享链接不存在或已过期
- 500: 服务器内部错误

### 5. 生成学习报告

**接口**: `POST /api/share/report`

**描述**: 一键生成学习报告

**请求**:
```json
{
  "shareId": "share-123e4567-e89b-12d3-a456-426614174000"
}
```

**响应**:
```json
{
  "totalExplorations": 15,
  "totalCollectedCards": 45,
  "categoryDistribution": {
    "自然类": 8,
    "生活类": 5,
    "人文类": 2
  },
  "recentCards": [...],
  "generatedAt": "2025-12-18T10:00:00Z"
}
```

**错误码**:
- 404: 分享链接不存在或已过期
- 500: 服务器内部错误

## 通用错误响应

所有接口在发生错误时，返回统一的错误格式：

```json
{
  "code": 400,
  "message": "请求参数错误",
  "detail": "图片格式不支持，仅支持 JPEG 和 PNG 格式"
}
```

## 性能要求

- 图像识别接口：响应时间≤3秒（90%请求）
- 知识卡片生成接口：响应时间≤5秒（90%请求）
- 其他接口：响应时间≤1秒（90%请求）

## 安全要求

- 所有接口使用 HTTPS
- 分享链接使用 UUID v4，不可预测
- 分享链接自动过期（7天）
- 不存储用户原始图片数据

## 版本管理

当前版本：v1.0.0

未来版本变更将通过 URL 路径版本控制（如 `/api/v2/...`）

