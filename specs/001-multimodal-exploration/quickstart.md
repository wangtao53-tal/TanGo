# TanGo 快速开始指南

**创建日期**: 2025-12-18  
**功能**: TanGo 多模态探索核心功能

## 项目概述

TanGo（小探号）是一个面向K12学生的探索学习应用，核心功能是"拍一得三"：孩子通过拍照识别真实世界对象，系统使用AI实时生成三张知识卡片（科学认知卡、古诗词/人文卡、英语表达卡）。

## 技术栈

### 前端
- React 18
- Vite
- Tailwind CSS
- TypeScript

### 后端
- Go 1.25.3
- go-zero v1.9.3
- eino（字节云原生AI框架）

## 快速开始

### 1. 环境准备

**前端环境**:
```bash
# 需要 Node.js 18+ 和 npm/yarn
node --version  # 应该 >= 18.0.0
npm --version   # 或 yarn --version
```

**后端环境**:
```bash
# 需要 Go 1.25.3
go version  # 应该显示 go1.25.3
```

### 2. 克隆项目

```bash
git clone <repository-url>
cd TanGo
git checkout 001-multimodal-exploration
```

### 3. 前端启动

```bash
cd frontend
npm install  # 或 yarn install
npm run dev  # 启动开发服务器
```

前端将在 `http://localhost:5173` 启动（Vite默认端口）

### 4. 后端启动

```bash
cd backend
go mod download  # 下载依赖
go run explore.go -f etc/explore.yaml  # 启动服务
```

后端将在配置文件中指定的端口启动（默认可能是8080）

### 5. 配置AI模型

在 `backend/eino/config.yaml` 中配置AI模型：

```yaml
models:
  vision:
    app_id: "your-vision-app-id"
    model: "claude-3-5-sonnet-vision"
  
  llm:
    app_id: "your-llm-app-id"
    model: "gpt-4"
```

**注意**: 需要申请AI模型的APP ID，具体流程请联系项目负责人。

## 项目结构

```
TanGo/
├── frontend/          # 前端工程
│   ├── src/
│   │   ├── components/  # React组件
│   │   ├── pages/       # 页面组件
│   │   ├── services/    # API服务
│   │   └── ...
│   └── package.json
│
├── backend/           # 后端工程
│   ├── api/           # API定义
│   ├── internal/      # 内部代码
│   │   ├── handler/   # HTTP处理器
│   │   ├── logic/     # 业务逻辑
│   │   └── agent/     # AI Agent
│   ├── eino/          # eino配置
│   └── go.mod
│
└── specs/             # 规范文档
    └── 001-multimodal-exploration/
        ├── spec.md    # 功能规范
        ├── plan.md    # 实现计划
        └── ...
```

## 核心功能流程

### 1. 用户首次使用

1. 打开应用
2. 选择年龄/年级（必填）
3. 进入首页

### 2. 探索流程

1. 点击拍照按钮
2. 拍摄真实世界对象
3. 系统识别对象（调用 `/api/explore/identify`）
4. 系统生成三张知识卡片（调用 `/api/explore/generate-cards`）
5. 展示结果页面
6. 用户可以收藏卡片到"我的探索图鉴"
7. 用户可以一键分享给家长

### 3. 家长端查看

1. 孩子分享探索结果，生成分享链接
2. 家长打开分享链接
3. 查看探索记录和收藏的卡片
4. 一键生成学习报告

## 开发指南

### 前端开发

**UI设计稿参考**:
- 所有UI设计稿位于 `stitch_ui/` 文件夹
- 每个页面包含 `code.html`（HTML代码）和 `screen.png`（设计截图）
- 实现时必须完全遵循设计稿的样式和交互

**添加新组件**:
```bash
cd frontend/src/components
# 创建新组件文件
# 参考 stitch_ui/ 中对应页面的HTML结构
```

**添加新页面**:
```bash
cd frontend/src/pages
# 创建新页面文件
# 参考 stitch_ui/ 中对应页面的HTML结构
# 在 App.tsx 中添加路由
```

**Tailwind配置**:
```javascript
// tailwind.config.js 必须包含设计稿中的所有颜色
// 参考 plan.md 中的"UI设计分析"章节
```

**API调用**:
```typescript
import { api } from '@/services/api'

// 调用图像识别接口
const result = await api.identify(imageData, age)
```

### 后端开发

**添加新接口**:
1. 在 `backend/api/explore.api` 中定义接口
2. 运行 `goctl api go -api explore.api -dir .` 生成代码
3. 在 `internal/logic/` 中实现业务逻辑

**AI模型调用**:
```go
// 在 logic 中使用 eino 调用AI模型
result, err := eino.CallVisionModel(imageData)
```

## 测试

### 前端测试

```bash
cd frontend
npm run test  # 运行单元测试
```

### 后端测试

```bash
cd backend
go test ./...  # 运行所有测试
```

## 部署

### 前端部署

```bash
cd frontend
npm run build  # 构建生产版本
# 将 dist/ 目录部署到静态文件服务器
```

### 后端部署

```bash
cd backend
go build -o explore explore.go  # 构建可执行文件
./explore -f etc/explore.yaml    # 运行服务
```

## 常见问题

### Q: AI模型调用失败怎么办？

A: 检查 `backend/eino/config.yaml` 中的APP ID配置是否正确，确认网络连接正常。

### Q: 前端无法连接后端？

A: 检查后端服务是否启动，确认 `frontend/src/services/api.ts` 中的API地址配置正确。

### Q: 图片上传失败？

A: 检查图片格式（仅支持JPEG和PNG），确认图片大小不超过10MB。

## 下一步

1. 阅读 [功能规范](./spec.md) 了解详细需求
2. 阅读 [实现计划](./plan.md) 了解技术架构
3. 阅读 [数据模型](./data-model.md) 了解数据结构
4. 阅读 [API合约](./contracts/README.md) 了解接口定义

## 获取帮助

- 查看项目文档：`specs/001-multimodal-exploration/`
- 联系项目负责人获取AI模型APP ID
- 提交Issue到项目仓库

