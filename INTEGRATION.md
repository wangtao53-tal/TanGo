# 前后端联调指南

## 项目状态

✅ **前后端基础框架已完成**
- 前端：React 18 + Vite + Tailwind CSS
- 后端：go-zero框架（Go 1.25.3）
- API接口：已实现Mock数据返回

✅ **前后端API联调已完成**
- 前端API服务已配置调用真实后端
- 拍照功能已实现，可调用识别API
- 结果页面已实现，可展示API返回的数据
- 支持降级到Mock数据（后端不可用时）

## 启动步骤

### 1. 启动后端服务

```bash
cd backend
go run explore.go -f etc/explore.yaml
```

后端服务将在 `http://localhost:8877` 启动。

### 2. 启动前端服务

```bash
cd frontend
npm install  # 首次运行需要
npm run dev
```

前端服务将在 `http://localhost:3000` 启动。

### 3. 测试流程

1. 打开浏览器访问 `http://localhost:3000`
2. 点击首页的"Photo Explore"按钮
3. 在拍照页面点击快门按钮，选择一张图片
4. 系统会自动调用后端API进行识别和生成卡片
5. 跳转到结果页面，展示三张知识卡片

## API接口说明

### 1. 图像识别接口

**POST** `/api/explore/identify`

**请求体**:
```json
{
  "image": "base64编码的图片数据",
  "age": 8  // 可选，孩子年龄
}
```

**响应**:
```json
{
  "objectName": "银杏",
  "objectCategory": "自然类",
  "confidence": 0.95,
  "keywords": ["植物", "古老", "叶子"]
}
```

### 2. 知识卡片生成接口

**POST** `/api/explore/generate-cards`

**请求体**:
```json
{
  "objectName": "银杏",
  "objectCategory": "自然类",
  "age": 8,
  "keywords": ["植物", "古老", "叶子"]
}
```

**响应**:
```json
{
  "cards": [
    {
      "type": "science",
      "title": "银杏的科学知识",
      "content": {
        "name": "银杏",
        "explanation": "...",
        "facts": [...],
        "funFact": "..."
      }
    },
    {
      "type": "poetry",
      "title": "古人怎么看银杏",
      "content": {
        "poem": "...",
        "explanation": "...",
        "context": "..."
      }
    },
    {
      "type": "english",
      "title": "用英语说银杏",
      "content": {
        "words": [...],
        "expressions": [...]
      }
    }
  ]
}
```

## 配置说明

### 前端配置

- **API代理**: `frontend/vite.config.ts` 中配置了代理到 `http://localhost:8877`
- **API基础URL**: `frontend/src/services/api.ts` 中配置为 `http://localhost:8877`
- **降级策略**: 如果后端不可用，前端会自动降级到Mock数据

### 后端配置

- **端口**: `backend/etc/explore.yaml` 中配置为 `8877`
- **CORS**: 已配置允许跨域请求
- **Mock数据**: 当前所有接口返回Mock数据，待AI模型接入后替换

## 故障排查

### 后端无法启动

1. 检查Go环境：`go version`
2. 检查依赖：`cd backend && go mod tidy`
3. 检查端口占用：`lsof -i :8877`

### 前端无法连接后端

1. 检查后端是否启动：访问 `http://localhost:8877/api/explore/identify`
2. 检查CORS配置：后端已配置允许所有来源
3. 检查浏览器控制台：查看是否有CORS错误

### API调用失败

1. 前端会自动降级到Mock数据
2. 检查浏览器控制台的错误信息
3. 检查后端日志：`backend/logs/` 目录

## 下一步工作

- [ ] 实现年龄/年级选择组件（首次使用必选）
- [ ] 实现收藏功能的数据持久化
- [ ] 接入真实AI模型（待APP ID提供后）
- [ ] 实现分享功能
- [ ] 实现学习报告功能

