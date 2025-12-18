# TanGo 前端启动指南

## 快速开始

### 方法一：使用启动脚本（推荐）

```bash
cd /Users/tal_1/TanGo/frontend
./start.sh
```

### 方法二：手动运行

```bash
# 1. 进入前端目录
cd /Users/tal_1/TanGo/frontend

# 2. 安装依赖（首次运行需要）
npm install

# 3. 启动开发服务器
npm run dev
```

## 前置要求

### 安装 Node.js

如果您的系统未安装 Node.js，请先安装：

**macOS (使用 Homebrew):**
```bash
brew install node
```

**或访问官网下载:**
- 访问 https://nodejs.org/
- 下载并安装 LTS 版本（推荐）

**验证安装:**
```bash
node --version  # 应显示版本号，如 v20.x.x
npm --version   # 应显示版本号，如 10.x.x
```

## 可用命令

- `npm run dev` - 启动开发服务器（热重载）
- `npm run build` - 构建生产版本
- `npm run preview` - 预览生产构建
- `npm run lint` - 运行代码检查

## 访问应用

启动成功后，在浏览器中打开：
- **本地访问**: http://localhost:5173/
- **网络访问**: 终端会显示实际地址

## 常见问题

### 1. 端口被占用
如果 5173 端口被占用，Vite 会自动尝试其他端口（如 5174、5175 等）

### 2. 依赖安装失败
- 检查网络连接
- 尝试清除缓存: `npm cache clean --force`
- 删除 `node_modules` 和 `package-lock.json` 后重新安装

### 3. 权限问题
如果遇到权限问题，可能需要使用 `sudo`（不推荐）或修复 npm 权限

## 开发提示

- 修改代码后会自动热重载，无需手动刷新
- 查看终端输出了解编译状态和错误信息
- 使用浏览器开发者工具调试
