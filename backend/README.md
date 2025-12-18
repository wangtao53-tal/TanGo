# TanGo 后端服务

TanGo（小探号）多模态探索核心功能后端服务，基于 go-zero 框架实现。

## 技术栈

- Go 1.21+
- go-zero v1.9.3
- eino（字节云原生AI框架，待接入）

## 快速开始

### 1. 环境要求

- Go 1.21 或更高版本
- goctl 工具（go-zero代码生成工具）

### 2. 安装依赖

```bash
cd backend
go mod download
```

### 3. 配置

编辑 `etc/explore.yaml` 配置文件：

```yaml
Name: explore
Host: 0.0.0.0
Port: 8888
```

### 4. 运行服务

```bash
go run explore.go -f etc/explore.yaml
```

服务将在 `http://0.0.0.0:8888` 启动

### 5. 测试API

```bash
# 图像识别（Mock）
curl -X POST http://localhost:8888/api/explore/identify \
  -H "Content-Type: application/json" \
  -d '{"image": "data:image/jpeg;base64,/9j/4AAQSkZJRg=="}'

# 生成知识卡片（Mock）
curl -X POST http://localhost:8888/api/explore/generate-cards \
  -H "Content-Type: application/json" \
  -d '{"objectName": "银杏", "objectCategory": "自然类", "age": 8}'
```

## 项目结构

```
backend/
├── api/                    # API定义文件
│   └── explore.api
├── internal/
│   ├── handler/            # HTTP处理器
│   ├── logic/              # 业务逻辑
│   ├── svc/                # 服务上下文
│   ├── types/              # 类型定义
│   ├── utils/              # 工具函数
│   └── agent/              # AI Agent（待实现）
├── eino/                   # eino框架配置（待配置）
│   └── models/
├── etc/                     # 配置文件
│   └── explore.yaml
├── explore.go              # 主程序入口
├── go.mod
└── go.sum
```

## API接口

### 1. 图像识别

**POST** `/api/explore/identify`

识别图片中的对象。

**请求**:
```json
{
  "image": "base64编码的图片数据",
  "age": 8  // 可选
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

### 2. 生成知识卡片

**POST** `/api/explore/generate-cards`

根据识别结果生成三张知识卡片。

**请求**:
```json
{
  "objectName": "银杏",
  "objectCategory": "自然类",
  "age": 8,
  "keywords": ["植物", "树木"]
}
```

**响应**:
```json
{
  "cards": [
    {
      "type": "science",
      "title": "银杏的科学知识",
      "content": {...}
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

### 3. 创建分享链接

**POST** `/api/share/create`

创建分享链接。

### 4. 获取分享数据

**GET** `/api/share/:shareId`

获取分享数据。

### 5. 生成学习报告

**POST** `/api/share/report`

生成学习报告。

## 开发说明

### 当前状态

- ✅ 后端框架搭建完成
- ✅ API接口框架就绪
- ✅ Mock数据实现完成
- ⏳ AI模型集成（待APP ID提供后）

### Mock数据

当前阶段所有AI模型调用使用Mock数据实现，包括：
- 图像识别：随机返回常见对象
- 知识卡片生成：根据对象名称和年龄生成Mock卡片内容

待APP ID提供后，将接入真实AI模型。

## 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./internal/logic/... -v
```

## 注意事项

1. 当前使用Mock数据，所有AI功能待APP ID提供后接入
2. 分享链接存储在内存中，服务重启后数据会丢失（生产环境应使用Redis）
3. CORS已配置为允许所有来源，生产环境应限制为特定域名

