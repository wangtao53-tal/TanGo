# TanGo 后端服务

TanGo（小探号）多模态探索核心功能后端服务，基于 go-zero 框架实现，为 4-18 岁孩子提供图像识别、知识卡片生成、智能对话等 AI 能力。

## 📋 目录

- [技术栈](#技术栈)
- [快速开始](#快速开始)
- [项目架构](#项目架构)
- [核心功能](#核心功能)
- [API 接口](#api-接口)
- [配置说明](#配置说明)
- [开发指南](#开发指南)
- [部署说明](#部署说明)

## 🛠 技术栈

- **框架**: Go 1.21+ / go-zero v1.9.3
- **AI 框架**: eino（字节云原生 AI 框架）
- **存储**: 内存存储（MemoryStorage）+ GitHub 存储（GitHubStorage）
- **架构模式**: ReAct Agent（推理-行动循环）

## 🚀 快速开始

### 环境要求

- Go 1.21 或更高版本
- goctl 工具（go-zero 代码生成工具，可选）

### 安装依赖

```bash
cd backend
go mod download
```

### 配置

#### 方式一：使用环境变量（推荐）

创建 `.env` 文件（在项目根目录）：

```bash
# 后端服务配置
BACKEND_HOST=0.0.0.0
BACKEND_PORT=8877

# eino AI 框架配置
EINO_BASE_URL=https://your-eino-base-url
TAL_MLOPS_APP_ID=your-app-id
TAL_MLOPS_APP_KEY=your-app-key

# AI 模型配置（可选，有默认值）
INTENT_MODEL=your-intent-model
IMAGE_RECOGNITION_MODELS=model1,model2
IMAGE_GENERATION_MODEL=your-image-generation-model
TEXT_GENERATION_MODEL=your-text-generation-model
USE_AI_MODEL=true  # true=使用AI模型，false=使用Mock数据

# GitHub 图片上传配置（可选）
GITHUB_TOKEN=your-github-token
GITHUB_OWNER=your-github-owner
GITHUB_REPO=your-repo-name
GITHUB_BRANCH=main
GITHUB_PATH=images/
MAX_IMAGE_SIZE=10485760  # 10MB
```

#### 方式二：使用配置文件

编辑 `etc/explore.yaml`：

```yaml
Name: explore
Host: 0.0.0.0
Port: 8877
Timeout: 180000  # 180秒，确保有足够时间处理3张卡片生成

AI:
  EinoBaseURL: ""
  AppID: ""
  AppKey: ""
  UseAIModel: true  # 是否使用AI模型，false表示使用Mock数据

Upload:
  GitHubToken: ""
  GitHubOwner: ""
  GitHubRepo: ""
  GitHubBranch: "main"
  GitHubPath: "images/"
  MaxImageSize: 10485760
```

**注意**: 环境变量优先级高于配置文件。

### 运行服务

```bash
# 开发模式
go run explore.go -f etc/explore.yaml

# 或使用构建后的二进制文件
go build -o explore explore.go
./explore -f etc/explore.yaml
```

服务将在 `http://0.0.0.0:8877` 启动。

### 测试 API

```bash
# 图像识别
curl -X POST http://localhost:8877/api/explore/identify \
  -H "Content-Type: application/json" \
  -d '{"image": "data:image/jpeg;base64,/9j/4AAQSkZJRg==", "age": 8}'

# 生成知识卡片
curl -X POST http://localhost:8877/api/explore/generate-cards \
  -H "Content-Type: application/json" \
  -d '{"objectName": "银杏", "objectCategory": "自然类", "age": 8}'
```

## 🏗 项目架构

### 目录结构

```
backend/
├── api/                    # API 定义文件（go-zero API 格式）
│   └── explore.api         # API 接口定义
├── internal/
│   ├── handler/            # HTTP 处理器层
│   │   ├── identifyhandler.go
│   │   ├── generatecardshandler.go
│   │   ├── conversationhandler.go
│   │   ├── streamhandler.go      # 流式对话处理器
│   │   └── ...
│   ├── logic/              # 业务逻辑层
│   │   ├── identifylogic.go
│   │   ├── generatecardslogic.go
│   │   ├── conversationlogic.go
│   │   └── ...
│   ├── agent/              # AI Agent 系统
│   │   ├── agent.go        # Agent 主入口
│   │   ├── graph.go        # 调用流程图
│   │   └── nodes/          # Agent 节点
│   │       ├── image_recognition.go    # 图像识别节点
│   │       ├── text_generation.go      # 文本生成节点
│   │       ├── image_generation.go    # 图像生成节点
│   │       ├── intent_recognition.go  # 意图识别节点
│   │       └── conversation_node.go   # 对话节点
│   ├── storage/            # 存储层
│   │   ├── memory.go       # 内存存储（会话、分享链接等）
│   │   └── github.go       # GitHub 存储（图片上传）
│   ├── config/             # 配置管理
│   │   ├── config.go       # 配置结构定义
│   │   └── models.go       # 默认模型配置
│   ├── svc/                # 服务上下文
│   │   └── servicecontext.go
│   ├── types/              # 类型定义
│   │   └── types.go
│   └── utils/              # 工具函数
├── eino/                   # eino 框架配置
│   └── models/
├── etc/                    # 配置文件
│   └── explore.yaml
├── logs/                   # 日志文件
├── explore.go              # 主程序入口
├── go.mod
└── go.sum
```

### 架构分层

```
┌─────────────────────────────────────┐
│         HTTP Handler 层              │  ← 处理 HTTP 请求/响应
├─────────────────────────────────────┤
│         Business Logic 层            │  ← 业务逻辑处理
├─────────────────────────────────────┤
│         AI Agent 层                  │  ← ReAct Agent 系统
│  ┌───────────────────────────────┐  │
│  │  Graph (调用流程图)             │  │
│  │  ├─ ImageRecognitionNode      │  │
│  │  ├─ TextGenerationNode         │  │
│  │  ├─ ImageGenerationNode        │  │
│  │  ├─ IntentRecognitionNode      │  │
│  │  └─ ConversationNode           │  │
│  └───────────────────────────────┘  │
├─────────────────────────────────────┤
│         Storage 层                   │  ← 数据存储
│  ├─ MemoryStorage (内存)             │
│  └─ GitHubStorage (GitHub)           │
└─────────────────────────────────────┘
```

### 核心组件

#### 1. Agent 系统

基于 **ReAct（Reasoning + Acting）** 模式的 AI Agent 系统，通过推理-行动循环处理复杂任务：

- **Graph**: 管理 AI 调用流程，协调各个节点
- **Nodes**: 独立的 AI 能力节点，每个节点负责特定任务
  - `ImageRecognitionNode`: 图像识别
  - `TextGenerationNode`: 文本生成（知识卡片内容）
  - `ImageGenerationNode`: 图像生成（卡片配图）
  - `IntentRecognitionNode`: 意图识别（理解用户意图）
  - `ConversationNode`: 对话生成（智能回复）

#### 2. 存储系统

- **MemoryStorage**: 内存存储，用于会话管理、分享链接等临时数据
  - 自动清理过期会话（默认 24 小时未活跃）
  - 线程安全（使用 `sync.Map`）
- **GitHubStorage**: GitHub 存储，用于图片上传
  - 支持通过 GitHub API 上传图片到仓库
  - 降级方案：如果未配置 GitHub，使用 base64 编码返回

## 🎯 核心功能

### 1. 图像识别

识别图片中的对象，返回对象名称、类别、置信度和关键词。

**流程**:
```
图片输入 → ImageRecognitionNode → 识别结果
```

### 2. 知识卡片生成

根据识别结果生成三张知识卡片：
- **科学认知卡** (science): 科学知识、原理
- **人文认知卡** (poetry): 古诗词、文化知识
- **语言认知卡** (english): 英语表达、词汇

**流程**:
```
识别结果 + 年龄 → TextGenerationNode (3次) → 3张卡片内容
                → ImageGenerationNode (3次) → 3张卡片配图
```

### 3. 智能对话

支持文本、语音、图片三种输入方式的智能对话，使用流式响应（SSE）实现打字机效果。

**流程**:
```
用户输入 → IntentRecognitionNode → 识别意图
        → ConversationNode → 生成回复（流式）
```

### 4. 分享功能

- 创建分享链接：将探索记录和收藏的卡片生成分享链接
- 获取分享数据：通过分享 ID 获取分享内容
- 生成学习报告：统计探索次数、收藏卡片数、类别分布等

### 5. 图片上传

支持将图片上传到 GitHub 仓库，返回可访问的 URL。

## 📡 API 接口

### 探索相关

#### 1. 图像识别

**POST** `/api/explore/identify`

识别图片中的对象。

**请求**:
```json
{
  "image": "data:image/jpeg;base64,...",
  "age": 8  // 可选，用于优化识别
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

#### 2. 生成知识卡片

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

### 对话相关

#### 3. 意图识别

**POST** `/api/conversation/intent`

识别用户消息的意图。

**请求**:
```json
{
  "message": "这是什么？",
  "sessionId": "session-123",  // 可选
  "context": []  // 可选，上下文消息
}
```

**响应**:
```json
{
  "intent": "generate_cards",  // 或 "text_response"
  "confidence": 0.95,
  "reason": "用户询问对象信息，需要生成卡片"
}
```

#### 4. 对话消息（非流式）

**POST** `/api/conversation/message`

发送对话消息，获取回复。

**请求**:
```json
{
  "message": "这是什么？",
  "image": "data:image/jpeg;base64,...",  // 可选
  "voice": "base64...",  // 可选
  "sessionId": "session-123",  // 可选
  "identificationContext": {...}  // 可选，识别结果上下文
}
```

**响应**:
```json
{
  "message": {
    "id": "msg-123",
    "type": "text",
    "sender": "assistant",
    "content": "这是银杏...",
    "timestamp": "2025-01-01T00:00:00Z",
    "sessionId": "session-123"
  },
  "sessionId": "session-123",
  "type": "text"  // 或 "cards"
}
```

#### 5. 流式对话（SSE）

**POST** `/api/conversation/stream`

发送对话消息，通过 Server-Sent Events (SSE) 流式返回回复。

**请求**:
```json
{
  "messageType": "text",  // "text" | "voice" | "image"
  "message": "这是什么？",  // 当 messageType 为 text 时必填
  "audio": "base64...",  // 当 messageType 为 voice 时必填
  "image": "base64...",  // 当 messageType 为 image 时必填
  "sessionId": "session-123",  // 可选
  "userAge": 8,  // 可选，3-18岁
  "maxContextRounds": 20  // 可选，最大上下文轮次
}
```

**响应** (SSE 流):
```
event: connected
data: {"type":"connected","sessionId":"session-123"}

event: message
data: {"type":"message","content":"这是","index":0,"sessionId":"session-123"}

event: message
data: {"type":"message","content":"这是银杏","index":1,"sessionId":"session-123"}

...

event: done
data: {"type":"done","sessionId":"session-123"}
```

### 分享相关

#### 6. 创建分享链接

**POST** `/api/share/create`

创建分享链接。

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
  "shareId": "share-123",
  "shareUrl": "https://tango.example.com/share/share-123",
  "expiresAt": "2025-01-08T00:00:00Z"
}
```

#### 7. 获取分享数据

**GET** `/api/share/:shareId`

获取分享数据。

**响应**:
```json
{
  "explorationRecords": [...],
  "collectedCards": [...],
  "createdAt": "2025-01-01T00:00:00Z",
  "expiresAt": "2025-01-08T00:00:00Z"
}
```

#### 8. 生成学习报告

**POST** `/api/share/report`

生成学习报告。

**请求**:
```json
{
  "shareId": "share-123"
}
```

**响应**:
```json
{
  "totalExplorations": 10,
  "totalCollectedCards": 25,
  "categoryDistribution": {
    "自然类": 5,
    "生活类": 3,
    "人文类": 2
  },
  "recentCards": [...],
  "generatedAt": "2025-01-01T00:00:00Z"
}
```

### 上传相关

#### 9. 图片上传

**POST** `/api/upload/image`

上传图片到 GitHub 仓库或返回 base64 编码。

**请求**:
```json
{
  "imageData": "base64编码的图片数据（不含data URL前缀）",
  "filename": "image.jpg"  // 可选
}
```

**响应**:
```json
{
  "url": "https://raw.githubusercontent.com/...",
  "filename": "image_1234567890.jpg",
  "size": 102400,
  "uploadMethod": "github"  // 或 "base64"
}
```

## ⚙️ 配置说明

### 环境变量配置

所有配置项都支持通过环境变量设置，优先级高于配置文件。

#### 服务配置

- `BACKEND_HOST`: 服务监听地址（默认: `0.0.0.0`）
- `BACKEND_PORT`: 服务端口（默认: `8877`）

#### AI 配置

- `EINO_BASE_URL`: eino 框架基础 URL（必填，如果使用真实模型）
- `TAL_MLOPS_APP_ID`: AI 模型 APP ID（必填，如果使用真实模型）
- `TAL_MLOPS_APP_KEY`: AI 模型 APP Key（必填，如果使用真实模型）
- `INTENT_MODEL`: 意图识别模型（可选，有默认值）
- `IMAGE_RECOGNITION_MODELS`: 图像识别模型列表，逗号分隔（可选，有默认值）
- `IMAGE_GENERATION_MODEL`: 图像生成模型（可选，有默认值）
- `TEXT_GENERATION_MODEL`: 文本生成模型（可选，有默认值）
- `USE_AI_MODEL`: 是否使用 AI 模型（`true`/`false`，默认: `true`）

#### 上传配置

- `GITHUB_TOKEN`: GitHub Personal Access Token（可选）
- `GITHUB_OWNER`: GitHub 用户名或组织名（可选）
- `GITHUB_REPO`: GitHub 仓库名（可选）
- `GITHUB_BRANCH`: GitHub 分支名（默认: `main`）
- `GITHUB_PATH`: 图片存储路径（默认: `images/`）
- `MAX_IMAGE_SIZE`: 图片大小限制，字节（默认: `10485760`，10MB）

### Mock 模式

如果未配置 eino 相关参数（`EINO_BASE_URL` 或 `TAL_MLOPS_APP_ID`），系统会自动使用 Mock 数据：

- 图像识别：随机返回常见对象
- 知识卡片生成：根据对象名称和年龄生成 Mock 卡片内容
- 对话：返回预设的回复

**启用 Mock 模式**:
```bash
USE_AI_MODEL=false
```

## 🔧 开发指南

### 代码生成

使用 goctl 生成代码（如果修改了 `api/explore.api`）：

```bash
# 安装 goctl
go install github.com/zeromicro/go-zero/tools/goctl@latest

# 生成代码
goctl api go -api api/explore.api -dir . -style gozero
```

### 添加新的 API

1. 在 `api/explore.api` 中定义 API
2. 运行 `goctl` 生成代码
3. 在 `internal/logic/` 中实现业务逻辑
4. 在 `internal/handler/` 中处理 HTTP 请求/响应

### 添加新的 Agent 节点

1. 在 `internal/agent/nodes/` 中创建新节点文件
2. 实现 `Node` 接口：
   ```go
   type Node interface {
       Execute(data *GraphData) (*GraphData, error)
   }
   ```
3. 在 `internal/agent/graph.go` 中注册节点
4. 在 `Graph` 中添加执行方法

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./internal/logic/... -v

# 运行测试并查看覆盖率
go test ./... -cover
```

### 日志

日志文件位于 `logs/` 目录：

- `access.log`: 访问日志
- `error.log`: 错误日志
- `severe.log`: 严重错误日志
- `slow.log`: 慢请求日志
- `stat.log`: 统计日志

日志配置在 `etc/explore.yaml` 中：

```yaml
Log:
  ServiceName: explore
  Mode: file
  Path: logs
  Level: info
  Compress: true
  KeepDays: 7
```

## 🚢 部署说明

### Docker 部署

```bash
# 构建镜像
docker build -t tango-backend .

# 运行容器
docker run -d \
  -p 8877:8877 \
  -e EINO_BASE_URL=... \
  -e TAL_MLOPS_APP_ID=... \
  -e TAL_MLOPS_APP_KEY=... \
  tango-backend
```

### 生产环境注意事项

1. **CORS 配置**: 当前允许所有来源，生产环境应限制为特定域名
2. **存储**: 当前使用内存存储，服务重启后数据会丢失，生产环境应使用 Redis
3. **日志**: 配置日志轮转和归档
4. **监控**: 添加健康检查接口和监控指标
5. **安全**: 配置 HTTPS、API 限流、认证等

### 静态文件服务

后端支持同时提供前端静态文件服务（用于 Docker 部署）。可通过环境变量控制：

```bash
ENABLE_STATIC_SERVER=false  # 禁用静态文件服务（使用 Nginx 时）
```

## 📝 开发状态

### 已完成 ✅

- [x] 后端框架搭建（go-zero）
- [x] API 接口定义和实现
- [x] AI Agent 系统（ReAct 模式）
- [x] 图像识别功能
- [x] 知识卡片生成功能
- [x] 智能对话功能（支持流式响应）
- [x] 分享功能
- [x] 图片上传功能（GitHub 存储）
- [x] Mock 数据支持

### 待完善 ⏳

- [ ] 生产环境存储方案（Redis）
- [ ] 性能优化和缓存
- [ ] 完整的单元测试和集成测试
- [ ] API 文档（Swagger/OpenAPI）

## 📚 相关文档

- [go-zero 官方文档](https://go-zero.dev/)
- [eino 框架文档](https://github.com/bytedance/eino)
- [项目根目录 README](../README.md)
- [前端 README](../frontend/README.md)

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

详见项目根目录 LICENSE 文件。
