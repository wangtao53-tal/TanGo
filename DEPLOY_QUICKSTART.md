# TanGo 快速部署指南

## 方案一：Docker 部署（3 步）

```bash
# 1. 配置环境变量
cp .env.example .env
vim .env

# 2. 启动服务
docker-compose up -d

# 3. 配置 Nginx（可选，如果使用 Nginx）
# 参考 nginx.conf.example
```

访问：`http://your-server:8877` 或通过 Nginx 配置的域名

## 方案二：静态编译部署（4 步）

```bash
# 1. 本地构建（需要 Node.js 和 Go）
./build.sh

# 2. 上传到服务器
scp -r build/ user@server:/path/to/tango/

# 3. 服务器配置
ssh user@server
cd /path/to/tango
cp .env.example .env
vim .env
chmod +x explore deploy.sh stop.sh

# 4. 启动服务
./deploy.sh
```

配置 Nginx：参考 `nginx.conf.example`

## 常用命令

### Docker 方案
```bash
# 启动
docker-compose up -d

# 停止
docker-compose down

# 查看日志
docker-compose logs -f

# 重启
docker-compose restart
```

### 静态编译方案
```bash
# 启动
./deploy.sh

# 停止
./stop.sh

# 查看日志
tail -f logs/explore.log
```

## 环境变量最小配置

```bash
BACKEND_HOST=0.0.0.0
BACKEND_PORT=8877
USE_AI_MODEL=true  # 或 false 使用 Mock 数据
```

详细配置请参考 `DEPLOY.md`

