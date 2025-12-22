# TanGo 部署指南

本文档提供两种部署方案，适用于服务器 Node.js 和 Go 版本较低的情况。

## 方案一：Docker 部署（推荐）

### 优点
- 环境隔离，不依赖服务器 Node.js/Go 版本
- 一键部署，易于管理
- 支持容器编排

### 前置要求
- 服务器已安装 Docker 和 Docker Compose

### 部署步骤

#### 1. 准备配置文件

确保项目根目录有 `.env` 文件：

```bash
# 复制示例文件
cp .env.example .env

# 编辑配置文件
vim .env
```

#### 2. 构建和启动

```bash
# 使用 docker-compose（推荐）
docker-compose up -d

# 或使用 docker 命令
docker build -t tango:latest .
docker run -d \
  --name tango \
  -p 8877:8877 \
  --env-file .env \
  tango:latest
```

#### 3. 配置 Nginx（可选）

如果使用 Nginx 提供前端静态文件服务，参考 `nginx.conf.example` 配置：

```bash
# 复制配置示例
cp nginx.conf.example /etc/nginx/sites-available/tango

# 编辑配置（修改域名和路径）
vim /etc/nginx/sites-available/tango

# 创建软链接
ln -s /etc/nginx/sites-available/tango /etc/nginx/sites-enabled/

# 测试配置
nginx -t

# 重载 Nginx
nginx -s reload
```

**注意**：如果使用 Nginx，需要在 Docker 容器中禁用静态文件服务：

```bash
# 在 .env 文件中添加
ENABLE_STATIC_SERVER=false
```

#### 4. 查看日志

```bash
# Docker 日志
docker logs -f tango

# 或使用 docker-compose
docker-compose logs -f
```

#### 5. 停止服务

```bash
# 使用 docker-compose
docker-compose down

# 或使用 docker 命令
docker stop tango
docker rm tango
```

## 方案二：静态编译部署

### 优点
- 无需 Docker，轻量级
- 可执行文件，直接运行
- 适合资源受限的服务器

### 前置要求
- 本地开发环境需要 Node.js 和 Go（用于构建）
- 服务器只需要基本的 Linux 环境

### 部署步骤

#### 1. 本地构建

在本地开发环境执行：

```bash
# 运行构建脚本
./build.sh
```

构建完成后，会在 `build/` 目录生成以下文件：
- `explore` - 后端可执行文件（Linux amd64）
- `frontend/` - 前端静态文件
- `etc/` - 后端配置文件
- `deploy.sh` - 部署脚本
- `stop.sh` - 停止脚本
- `.env.example` - 环境变量示例

#### 2. 上传到服务器

```bash
# 使用 scp 上传
scp -r build/ user@server:/path/to/tango/

# 或使用 rsync
rsync -avz build/ user@server:/path/to/tango/
```

#### 3. 服务器配置

```bash
# SSH 登录服务器
ssh user@server

# 进入部署目录
cd /path/to/tango

# 创建环境变量文件
cp .env.example .env
vim .env  # 编辑配置

# 设置执行权限
chmod +x explore deploy.sh stop.sh
```

#### 4. 启动服务

```bash
# 运行部署脚本
./deploy.sh
```

#### 5. 配置 Nginx

参考 `nginx.conf.example` 配置 Nginx：

```bash
# 复制配置示例
sudo cp nginx.conf.example /etc/nginx/sites-available/tango

# 编辑配置
sudo vim /etc/nginx/sites-available/tango
# 修改以下内容：
# - server_name: 你的域名
# - root: /path/to/tango/frontend（前端静态文件路径）
# - proxy_pass: http://localhost:8877（后端服务地址）

# 启用配置
sudo ln -s /etc/nginx/sites-available/tango /etc/nginx/sites-enabled/

# 测试并重载
sudo nginx -t
sudo nginx -s reload
```

#### 6. 查看日志

```bash
# 查看后端日志
tail -f logs/explore.log
```

#### 7. 停止服务

```bash
./stop.sh
```

## 环境变量配置

两种方案都需要配置 `.env` 文件，主要配置项：

```bash
# 后端服务配置
BACKEND_HOST=0.0.0.0
BACKEND_PORT=8877

# AI 模型配置
EINO_BASE_URL=
APP_ID=
APP_KEY=
USE_AI_MODEL=true

# GitHub 图片上传配置
GITHUB_TOKEN=
GITHUB_OWNER=
GITHUB_REPO=
GITHUB_BRANCH=main
GITHUB_PATH=images/

# Docker 方案专用
ENABLE_STATIC_SERVER=true  # 是否启用后端静态文件服务（使用 Nginx 时设为 false）
```

## 端口说明

- **后端服务端口**：默认 8877（可在 .env 中修改）
- **前端访问端口**：
  - Docker 方案：如果使用后端静态文件服务，直接访问 8877
  - Nginx 方案：通过 Nginx 配置的端口（通常是 80 或 443）

## 常见问题

### 1. Docker 容器无法访问

检查端口映射和防火墙设置：

```bash
# 检查端口是否被占用
lsof -i :8877

# 检查防火墙
sudo ufw status
```

### 2. 静态文件 404

- **Docker 方案**：检查 `ENABLE_STATIC_SERVER` 环境变量
- **静态编译方案**：检查 Nginx 配置中的 `root` 路径是否正确

### 3. API 请求失败

检查后端服务是否正常运行：

```bash
# 测试后端 API
curl http://localhost:8877/api/explore/identify

# 查看日志
tail -f logs/explore.log
```

### 4. 跨域问题

后端已配置 CORS，如果仍有问题，检查：
- Nginx 配置中的 `proxy_set_header` 设置
- 后端服务的 CORS 配置

## 性能优化建议

1. **使用 Nginx 提供静态文件**：比后端提供静态文件性能更好
2. **启用 Gzip 压缩**：在 Nginx 配置中已包含
3. **静态资源缓存**：Nginx 配置中已设置长期缓存
4. **使用 HTTPS**：生产环境建议配置 SSL 证书

## 监控和维护

### 日志管理

```bash
# 查看实时日志
tail -f logs/explore.log

# 日志轮转（建议配置 logrotate）
sudo vim /etc/logrotate.d/tango
```

### 服务管理

```bash
# 使用 systemd 管理服务（静态编译方案）
sudo vim /etc/systemd/system/tango.service
```

systemd 服务配置示例：

```ini
[Unit]
Description=TanGo Backend Service
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/path/to/tango
ExecStart=/path/to/tango/explore -f /path/to/tango/etc/explore.yaml
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

启用服务：

```bash
sudo systemctl enable tango
sudo systemctl start tango
sudo systemctl status tango
```

## 更新部署

### Docker 方案

```bash
# 拉取最新代码
git pull

# 重新构建和部署
docker-compose down
docker-compose build
docker-compose up -d
```

### 静态编译方案

```bash
# 本地重新构建
./build.sh

# 上传新版本到服务器
scp -r build/* user@server:/path/to/tango/

# 重启服务
ssh user@server
cd /path/to/tango
./stop.sh
./deploy.sh
```

## 安全建议

1. **生产环境禁用 CORS 通配符**：修改后端代码，限制允许的域名
2. **使用 HTTPS**：配置 SSL 证书
3. **限制 API 访问**：使用防火墙或 Nginx 访问控制
4. **定期更新依赖**：保持 Docker 镜像和依赖包最新
5. **保护敏感信息**：不要将 `.env` 文件提交到版本控制

