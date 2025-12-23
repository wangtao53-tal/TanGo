# API 地址配置修复说明

## 问题描述

部署后调用 `/api/upload/image` 等接口时，前端使用了错误的 API 地址（`http://localhost:8877`），导致无法访问后端。

## 解决方案

已修复前端代码，现在生产环境默认使用**相对路径**，通过 Nginx 代理访问后端。

## 两种部署场景

### 场景 1：使用 Nginx 代理（推荐）

**配置方式**：无需额外配置，使用默认的相对路径即可。

1. **构建前端**（无需设置环境变量）：
```bash
./build.sh
```

2. **配置 Nginx**（参考 `nginx.conf.example`）：
```nginx
location /api {
    proxy_pass http://localhost:8877;
    # ... 其他配置
}
```

3. **访问方式**：
- 前端：`http://your-domain.com/`
- API：`http://your-domain.com/api/...`（自动代理到后端）

### 场景 2：不使用 Nginx，直接访问后端

**配置方式**：构建时设置 `VITE_API_BASE_URL` 环境变量。

1. **构建前端**（设置 API 地址）：
```bash
# 方式 1：通过环境变量
export VITE_API_BASE_URL=http://your-server:8877
./build.sh

# 方式 2：一行命令
VITE_API_BASE_URL=http://your-server:8877 ./build.sh
```

2. **访问方式**：
- 前端：`http://your-server:8877/`
- API：`http://your-server:8877/api/...`

## 已修复的文件

以下文件已更新，生产环境默认使用相对路径：

1. `frontend/src/services/api.ts` - 主要 API 客户端
2. `frontend/src/services/sse.ts` - SSE 连接
3. `frontend/src/services/sse-post.ts` - POST + SSE 连接
4. `frontend/src/hooks/useStreamConversation.ts` - 流式对话 Hook

## 重新部署步骤

### 如果使用 Nginx（推荐）

```bash
# 1. 重新构建（无需设置环境变量）
./build.sh

# 2. 上传新的前端文件到服务器
scp -r build/frontend/* user@server:/path/to/tango/frontend/

# 3. 确保 Nginx 配置正确（参考 nginx.conf.example）
# 4. 重启 Nginx（如果需要）
sudo nginx -s reload
```

### 如果不使用 Nginx

```bash
# 1. 设置 API 地址并构建
VITE_API_BASE_URL=http://your-server:8877 ./build.sh

# 2. 上传新的前端文件到服务器
scp -r build/frontend/* user@server:/path/to/tango/frontend/
```

## 验证修复

部署后，打开浏览器开发者工具（F12），查看 Network 标签：

1. **使用 Nginx**：API 请求应该是 `/api/upload/image`（相对路径）
2. **不使用 Nginx**：API 请求应该是 `http://your-server:8877/api/upload/image`（完整 URL）

## 常见问题

### Q: 为什么生产环境使用相对路径？

A: 使用相对路径可以让 Nginx 统一处理前端和 API 请求，避免跨域问题，也更灵活。

### Q: 如何检查当前使用的 API 地址？

A: 在浏览器控制台执行：
```javascript
console.log(import.meta.env.VITE_API_BASE_URL || '使用相对路径')
```

### Q: 如果后端和前端不在同一域名怎么办？

A: 构建时设置 `VITE_API_BASE_URL` 为完整的后端地址：
```bash
VITE_API_BASE_URL=https://api.your-domain.com ./build.sh
```

### Q: Docker 部署需要设置吗？

A: 如果 Docker 容器内后端提供静态文件服务，且使用同一端口，无需设置（使用相对路径）。
如果使用 Nginx 在容器外，也无需设置（使用相对路径）。

