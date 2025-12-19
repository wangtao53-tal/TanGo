# TanGo 快速启动指南

> 面向 3–12 岁儿童及小学阶段学生的探索学习平台

## 📋 目录

- [前置要求](#前置要求)
- [快速开始](#快速开始)
- [访问服务](#访问服务)
- [查看日志](#查看日志)
- [停止服务](#停止服务)
- [常见问题](#常见问题)
- [开发模式](#开发模式)
- [生产部署](#生产部署)
- [更多信息](#更多信息)

## 📦 前置要求

在开始之前，请确保已安装以下工具：

| 工具 | 版本要求 | 检查命令 |
|------|---------|---------|
| Go | 1.21+ | `go version` |
| Node.js | 18+ | `node -v` |
| npm/yarn | 最新版 | `npm -v` 或 `yarn -v` |

> 💡 **提示**: 如果未安装，请访问 [Go 官网](https://go.dev/) 和 [Node.js 官网](https://nodejs.org/) 下载安装。

## 🚀 快速开始

### 步骤 1: 配置环境变量

1. **复制配置文件模板**

   ```bash
   cp .env.example .env
   ```

2. **编辑 `.env` 文件**

   使用文本编辑器打开 `.env` 文件，填入实际的配置值：

   ```bash
   # 必填：AI模型认证信息
   TAL_MLOPS_APP_ID=your_app_id_here
   TAL_MLOPS_APP_KEY=your_app_key_here
   
   # 可选：端口配置（默认值）
   # BACKEND_PORT=8877
   # FRONTEND_PORT=3000
   
   # 其他配置根据需要调整
   ```

   > ⚠️ **注意**: 
   > - 如果未配置 AI 模型认证信息，系统会自动使用 Mock 模式（功能可用但使用模拟数据）
   > - 首次启动时，启动脚本会自动检测并提示创建 `.env` 文件

### 步骤 2: 一键启动（推荐）

这是最简单快捷的启动方式，脚本会自动处理所有依赖和配置。

**Linux/macOS:**
```bash
chmod +x start.sh  # 首次运行需要添加执行权限
./start.sh
```

**Windows:**
```cmd
start.bat
```

**脚本会自动执行：**
- ✅ 检查依赖环境（Go、Node.js）
- ✅ 检查并加载 `.env` 配置
- ✅ 安装前端依赖（如需要）
- ✅ 启动后端服务（端口 8877）
- ✅ 启动前端服务（端口 3000）

> 💡 **提示**: 启动脚本会在后台运行服务，日志会保存到 `backend.log` 和 `frontend.log` 文件中。

### 步骤 3: 验证启动

启动成功后，你应该看到：

```
✅ 后端服务已启动 (PID: xxxx)
✅ 前端服务已启动 (PID: xxxx)
```

如果看到错误信息，请查看 [常见问题](#常见问题) 部分。

### 手动启动（可选）

如果需要分别启动或调试，可以手动启动：

**启动后端：**
```bash
cd backend
go run explore.go -f etc/explore.yaml
```

**启动前端（新终端窗口）：**
```bash
cd frontend
npm install  # 首次运行需要安装依赖
npm run dev
```

> 💡 **提示**: 手动启动时，日志会直接输出到终端，便于调试。

## 🌐 访问服务

启动成功后，可以通过以下地址访问：

| 服务 | 地址 | 说明 |
|------|------|------|
| **前端应用** | http://localhost:3000 | 主应用界面 |
| **后端 API** | http://localhost:8877 | API 服务 |
| **API 文档** | http://localhost:8877/api/explore/identify | API 接口文档 |

> 💡 **提示**: 如果端口被占用，可以在 `.env` 文件中修改 `BACKEND_PORT` 和 `FRONTEND_PORT` 配置。

## 📝 查看日志

### 使用启动脚本时

日志文件位于项目根目录：

```bash
# 查看后端日志（实时）
tail -f backend.log

# 查看前端日志（实时）
tail -f frontend.log

# 查看最近 50 行后端日志
tail -n 50 backend.log

# 查看最近 50 行前端日志
tail -n 50 frontend.log
```

### 手动启动时

- 后端日志直接输出到终端
- 前端日志直接输出到终端（Vite 开发服务器）

> 💡 **提示**: 使用 `Ctrl+C` 可以停止 `tail -f` 命令。

## 🛑 停止服务

### 使用停止脚本（推荐）

**Linux/macOS:**
```bash
./stop.sh
```

**Windows:**
```cmd
stop.bat
```

**脚本会自动：**
- ✅ 停止所有后端服务进程
- ✅ 停止所有前端服务进程
- ✅ 清理端口占用
- ✅ 可选清理日志文件

### 使用启动脚本时

在启动脚本运行的终端中，按 `Ctrl+C` 即可停止所有服务。

### 手动启动时

在各自终端窗口按 `Ctrl+C` 停止对应的服务。

### 强制停止（如果脚本失效）

如果服务无法正常停止，可以手动终止进程：

```bash
# 查找并终止后端进程
lsof -ti:8877 | xargs kill -9

# 查找并终止前端进程
lsof -ti:3000 | xargs kill -9

# 或者使用 pkill（Linux/macOS）
pkill -f "go run explore.go"
pkill -f "vite"
```

## ❓ 常见问题

### 1. 端口被占用

**问题**: 启动时提示端口已被占用

**解决方案**:

1. 修改 `.env` 文件中的端口配置：
   ```bash
   BACKEND_PORT=8878  # 修改后端端口
   FRONTEND_PORT=3001  # 修改前端端口
   ```

2. 或者停止占用端口的进程：
   ```bash
   # 查看占用端口的进程
   lsof -i:8877  # 后端端口
   lsof -i:3000  # 前端端口
   
   # 终止进程（替换 PID 为实际进程 ID）
   kill -9 <PID>
   ```

### 2. 前端无法连接后端

**问题**: 前端页面显示无法连接到后端 API

**解决方案**:

1. 检查后端服务是否正常运行：
   ```bash
   curl http://localhost:8877/api/explore/identify
   ```

2. 检查 `.env` 文件中的配置：
   ```bash
   VITE_API_BASE_URL=http://localhost:8877
   # 或
   VITE_BACKEND_HOST=localhost
   VITE_BACKEND_PORT=8877
   ```

3. 检查防火墙设置，确保端口未被阻止

4. 查看后端日志确认是否有错误信息

### 3. AI 模型调用失败

**问题**: AI 功能无法使用或返回错误

**解决方案**:

1. **检查认证信息**:
   - 确认 `TAL_MLOPS_APP_ID` 和 `TAL_MLOPS_APP_KEY` 已正确配置
   - 验证认证信息是否有效

2. **检查服务地址**:
   - 确认 `EINO_BASE_URL` 配置正确
   - 测试网络连接是否正常

3. **查看日志**:
   ```bash
   tail -f backend.log | grep -i error
   ```

4. **Mock 模式**:
   - 如果配置未设置或无效，系统会自动使用 Mock 模式
   - Mock 模式下功能仍然可用，但使用模拟数据
   - 查看日志确认是否进入 Mock 模式

### 4. 依赖安装失败

**问题**: `npm install` 或 `go mod download` 失败

**解决方案**:

1. **前端依赖**:
   ```bash
   # 清除缓存后重试
   npm cache clean --force
   rm -rf node_modules package-lock.json
   npm install
   ```

2. **后端依赖**:
   ```bash
   cd backend
   go clean -modcache
   go mod download
   ```

3. **网络问题**:
   - 检查网络连接
   - 考虑使用国内镜像源（如淘宝 npm 镜像）

### 5. 权限问题（Linux/macOS）

**问题**: `./start.sh` 提示权限不足

**解决方案**:
```bash
chmod +x start.sh
chmod +x stop.sh
```

## 💻 开发模式

### 后端开发

```bash
cd backend
go run explore.go -f etc/explore.yaml
```

**特点**:
- 修改代码后需要手动重启服务
- 日志直接输出到终端，便于调试
- 支持 Go 的热重载工具（如 `air`）

### 前端开发

```bash
cd frontend
npm run dev
```

**特点**:
- 支持热重载（HMR），修改代码后自动刷新
- 实时显示编译错误和警告
- 开发服务器运行在 http://localhost:3000

### 开发工具推荐

- **后端**: 使用 [air](https://github.com/cosmtrek/air) 实现自动重载
- **前端**: Vite 已内置热重载支持
- **调试**: 使用浏览器开发者工具和 Go 调试器

## 🚢 生产部署

### 构建前端

```bash
cd frontend
npm run build
```

构建产物位于 `frontend/dist` 目录，可以部署到任何静态文件服务器。

### 构建后端

```bash
cd backend
go build -o explore explore.go
```

生成可执行文件 `explore`，然后运行：

```bash
./explore -f etc/explore.yaml
```

### 部署建议

1. **环境变量**: 确保生产环境的 `.env` 文件配置正确
2. **反向代理**: 使用 Nginx 或 Caddy 作为反向代理
3. **进程管理**: 使用 systemd、supervisor 或 PM2 管理进程
4. **日志管理**: 配置日志轮转和集中收集
5. **监控**: 设置健康检查和监控告警

## 📚 更多信息

- **配置说明**: 查看 [CONFIG.md](./CONFIG.md)（如果存在）
- **API 文档**: 查看 `backend/api/explore.api`
- **项目结构**: 查看 [README.md](./README.md)
- **问题反馈**: 提交 Issue 或联系开发团队

---

> 💡 **提示**: 如果遇到其他问题，请查看日志文件或联系技术支持。
