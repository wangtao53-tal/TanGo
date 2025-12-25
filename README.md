# TanGo（小探号）

> 《Hackathon 2025》专为3-12岁儿童打造，跟随小探号一起探索世界的冒险！

TanGo（小探号）是一款多模态探索学习应用，通过拍照识别真实世界对象，为孩子们生成三张定制化的知识卡片（科学认知卡、人文认知卡、语言认知卡），实现"一次观察，三重收获"的学习体验。

## ✨ 核心特性

### 🎯 拍一得三知识卡片
- **科学认知卡**: 科学知识、原理和事实
- **人文认知卡**: 古诗词、文化知识和历史背景
- **语言认知卡**: 英语表达、核心词汇和口语练习

### 🤖 AI 驱动的智能识别
- 图像识别：准确识别自然类、生活类、人文类对象
- 年龄适配：根据孩子年龄（3-18岁）自动调整内容难度
- 智能对话：支持文本、语音、图片多模态输入

### 📚 探索与收藏
- 探索图鉴：收藏喜欢的知识卡片
- 一键分享：分享探索结果给家长
- 学习报告：生成学习统计报告

### 🌍 国际化支持
- 中文/英文切换
- 符合 K12 教育分级标准

## 🏗 技术架构

### 前端技术栈
- **框架**: React 19.2.0 + TypeScript 5.9.3
- **构建工具**: Vite 7.2.4
- **路由**: React Router DOM 7.11.0
- **样式**: Tailwind CSS 4.1.18
- **国际化**: react-i18next 15.7.4
- **HTTP 客户端**: Axios 1.13.2
- **其他**: react-markdown, html2canvas, react-swipeable

### 后端技术栈
- **框架**: Go 1.21+ / go-zero v1.9.3
- **AI 框架**: eino（字节云原生 AI 框架）
- **架构模式**: ReAct Agent（推理-行动循环）
- **存储**: 内存存储（MemoryStorage）+ GitHub 存储（GitHubStorage）

### 项目结构
```
TanGo/
├── frontend/              # 前端应用
│   ├── src/
│   │   ├── components/   # React 组件
│   │   ├── pages/        # 页面组件
│   │   ├── services/     # API 服务
│   │   ├── hooks/        # React Hooks
│   │   ├── i18n/         # 国际化
│   │   └── utils/        # 工具函数
│   └── package.json
├── backend/              # 后端服务
│   ├── api/              # API 定义
│   ├── internal/
│   │   ├── handler/      # HTTP 处理器
│   │   ├── logic/        # 业务逻辑
│   │   ├── agent/        # AI Agent 系统
│   │   └── storage/      # 存储层
│   └── go.mod
├── specs/                # 功能规范文档
├── stitch_ui/            # UI 设计稿
├── build/                # 构建产物
├── docker-compose.yml     # Docker 编排
├── Dockerfile            # Docker 镜像
├── build.sh              # 构建脚本
└── start.sh              # 启动脚本
```

## 🚀 快速开始

### 环境要求

- **Go**: 1.21 或更高版本
- **Node.js**: 18+ (推荐 LTS 版本)
- **npm**: 或 yarn

### 安装步骤

1. **克隆项目**
```bash
git clone <repository-url>
cd TanGo
```

2. **配置环境变量**

创建 `.env` 文件（在项目根目录）：
```bash
# 后端服务配置
BACKEND_HOST=0.0.0.0
BACKEND_PORT=8877

# 前端服务配置
FRONTEND_PORT=3000

# eino AI 框架配置（可选，未配置将使用 Mock 数据）
EINO_BASE_URL=https://your-eino-base-url
TAL_MLOPS_APP_ID=your-app-id
TAL_MLOPS_APP_KEY=your-app-key
USE_AI_MODEL=true  # true=使用AI模型，false=使用Mock数据

# AI 模型配置（可选，有默认值）
INTENT_MODEL=your-intent-model
IMAGE_RECOGNITION_MODELS=model1,model2
IMAGE_GENERATION_MODEL=your-image-generation-model
TEXT_GENERATION_MODEL=your-text-generation-model

# GitHub 图片上传配置（可选）
GITHUB_TOKEN=your-github-token
GITHUB_OWNER=your-github-owner
GITHUB_REPO=your-repo-name
GITHUB_BRANCH=main
GITHUB_PATH=images/

# 前端 API 地址（开发环境）
VITE_API_BASE_URL=http://localhost:8877
```

3. **安装依赖**

```bash
# 安装后端依赖
cd backend
go mod download

# 安装前端依赖
cd ../frontend
npm install
```

4. **启动服务**

#### 方式一：使用启动脚本（推荐）

```bash
# 在项目根目录
chmod +x start.sh
./start.sh
```

脚本会自动：
- 检查环境变量配置
- 检查依赖是否安装
- 启动后端服务（端口 8877）
- 启动前端服务（端口 3000）

#### 方式二：手动启动

```bash
# 启动后端（终端 1）
cd backend
go run explore.go -f etc/explore.yaml

# 启动前端（终端 2）
cd frontend
npm run dev
```

5. **访问应用**

- **前端**: http://localhost:3000
- **后端 API**: http://localhost:8877

## 📖 使用指南

### 开发模式

#### 前端开发
```bash
cd frontend
npm run dev        # 启动开发服务器
npm run build      # 构建生产版本
npm run lint       # 代码检查
npm run preview    # 预览构建结果
```

#### 后端开发
```bash
cd backend
go run explore.go -f etc/explore.yaml  # 启动开发服务器
go test ./...                          # 运行测试
go build -o explore explore.go         # 构建可执行文件
```

### API 接口

#### 图像识别
```bash
POST /api/explore/identify
Content-Type: application/json

{
  "image": "data:image/jpeg;base64,...",
  "age": 8
}
```

#### 生成知识卡片
```bash
POST /api/explore/generate-cards
Content-Type: application/json

{
  "objectName": "银杏",
  "objectCategory": "自然类",
  "age": 8,
  "keywords": ["植物", "树木"]
}
```

#### 流式对话（SSE）
```bash
POST /api/conversation/stream
Content-Type: application/json

{
  "messageType": "text",
  "message": "这是什么？",
  "sessionId": "session-123",
  "userAge": 8
}
```

更多 API 文档请参考 [backend/README.md](./backend/README.md)。

## 🐳 Docker 部署

### 使用 Docker Compose

```bash
# 构建并启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

### 使用 Dockerfile

```bash
# 构建镜像
docker build -t tango:latest .

# 运行容器
docker run -d \
  -p 8877:8877 \
  -e EINO_BASE_URL=... \
  -e TAL_MLOPS_APP_ID=... \
  -e TAL_MLOPS_APP_KEY=... \
  tango:latest
```

## 📦 构建部署

### 构建生产版本

```bash
# 使用构建脚本（推荐）
chmod +x build.sh
./build.sh
```

构建产物位于 `build/` 目录：
- `build/explore` - 后端可执行文件
- `build/frontend/` - 前端静态文件
- `build/etc/` - 配置文件
- `build/deploy.sh` - 部署脚本

### 部署到服务器

1. **上传构建产物**
```bash
scp -r build/* user@server:/path/to/tango/
```

2. **配置环境变量**
```bash
cd /path/to/tango
cp .env.example .env
vim .env  # 编辑配置
```

3. **启动服务**
```bash
chmod +x deploy.sh
./deploy.sh
```

4. **配置 Nginx**（可选）

参考 `nginx.conf.example` 配置文件，配置反向代理和静态文件服务。

## 🧪 测试

### 前端测试
```bash
cd frontend
npm test
```

### 后端测试
```bash
cd backend
go test ./...
go test ./... -cover  # 查看覆盖率
```

### 集成测试
```bash
# 运行集成测试脚本
chmod +x scripts/test_integration.sh
./scripts/test_integration.sh
```

## 📚 文档

- [前端开发文档](./frontend/README.md)
- [后端开发文档](./backend/README.md)
- [功能规范文档](./specs/)
- [API 接口文档](./backend/README.md#api-接口)

## 🎨 UI 设计

UI 设计稿位于 `stitch_ui/` 目录，包含：
- 首页设计
- 拍照界面
- 识别结果页面
- 知识卡片详情页
- 收藏页面
- 学习报告页面

前端实现必须完全遵循设计稿规范。

## 🔧 配置说明

### Mock 模式

如果未配置 eino 相关参数，系统会自动使用 Mock 数据：
- 图像识别：随机返回常见对象
- 知识卡片生成：根据对象名称和年龄生成 Mock 卡片内容
- 对话：返回预设的回复

**启用 Mock 模式**:
```bash
USE_AI_MODEL=false
```

### 环境变量优先级

1. 环境变量（最高优先级）
2. `.env` 文件
3. 配置文件 `backend/etc/explore.yaml`（最低优先级）

## 🤝 贡献指南

### Git 提交规范

使用中文编写 commit message，格式：`<type>(<scope>): <subject>`

**提交类型**:
- `feat`: 新增功能
- `fix`: 修复 bug
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 代码重构
- `perf`: 性能优化
- `test`: 测试相关
- `chore`: 构建工具、依赖更新
- `ui`: UI/UX 相关

**作用域格式**: `<area>:<module>`

**示例**:
```bash
feat(frontend:pages): 实现识别结果页面
fix(backend:handler): 修复API调用错误处理
docs(README): 更新快速开始指南
```

详细规范请参考项目根目录的 Git 提交规范文档。

### 代码规范

- **前端**: 遵循 ESLint 配置，使用 TypeScript 严格模式
- **后端**: 遵循 Go 代码规范，使用 gofmt 格式化
- **注释**: 优先使用中文注释
- **命名**: 前端使用 camelCase，后端使用 Go 命名规范

## 📝 开发状态

### 已完成 ✅

- [x] 项目框架搭建（前后端分离）
- [x] 核心功能实现（拍照识别、知识卡片生成）
- [x] 智能对话功能（支持流式响应）
- [x] 收藏和分享功能
- [x] 国际化支持（中文/英文）
- [x] Docker 部署支持
- [x] Mock 数据支持

### 待完善 ⏳

- [ ] 完整的单元测试和 E2E 测试
- [ ] 性能优化和缓存策略
- [ ] 生产环境存储方案（Redis）
- [ ] API 文档（Swagger/OpenAPI）
- [ ] 监控和日志系统
- [ ] PWA 支持

## 🐛 常见问题

### 端口被占用

```bash
# 查看端口占用
lsof -ti:8877  # 后端端口
lsof -ti:3000  # 前端端口

# 停止服务
./stop.sh  # 或手动 kill 进程
```

### 依赖安装失败

```bash
# 前端
cd frontend
rm -rf node_modules package-lock.json
npm cache clean --force
npm install

# 后端
cd backend
go mod tidy
go mod download
```

### API 请求失败

- 检查后端服务是否启动
- 检查环境变量配置（`VITE_API_BASE_URL`）
- 检查 CORS 配置（后端需要允许前端域名）
- 查看浏览器控制台和服务器日志

## 📄 许可证

详见 [LICENSE](./LICENSE) 文件。

## 🙏 致谢

- [go-zero](https://go-zero.dev/) - 优秀的 Go 微服务框架
- [eino](https://github.com/bytedance/eino) - 字节云原生 AI 框架
- [React](https://react.dev/) - 前端 UI 框架
- [Vite](https://vite.dev/) - 下一代前端构建工具

## 📮 联系方式

如有问题或建议，欢迎提交 Issue 或 Pull Request。

---

**TanGo（小探号）** - 让探索成为最好的学习方式 🌟
