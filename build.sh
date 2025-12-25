#!/bin/bash

# TanGo 静态编译构建脚本
# 功能：在本地构建前端和后端，生成可部署的文件

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 项目根目录
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BUILD_DIR="$ROOT_DIR/build"
FRONTEND_DIR="$ROOT_DIR/frontend"
BACKEND_DIR="$ROOT_DIR/backend"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  TanGo 静态编译构建脚本${NC}"
echo -e "${BLUE}========================================${NC}\n"

# 清理旧的构建目录
echo -e "${YELLOW}📦 清理构建目录...${NC}"
rm -rf "$BUILD_DIR"
mkdir -p "$BUILD_DIR"

# 构建前端
echo -e "\n${BLUE}🔨 构建前端...${NC}"
cd "$FRONTEND_DIR"

# 检查 node_modules
if [ ! -d "node_modules" ]; then
    echo -e "${YELLOW}前端依赖未安装，正在安装...${NC}"
    npm install
fi

# 检查是否设置了 API 地址环境变量
# 如果使用 Nginx 代理，不需要设置（使用相对路径）
# 如果直接访问后端，需要设置 VITE_API_BASE_URL
if [ -z "$VITE_API_BASE_URL" ]; then
    echo -e "${YELLOW}提示: 未设置 VITE_API_BASE_URL，生产环境将使用相对路径${NC}"
    echo -e "${YELLOW}   - 使用 Nginx 代理: 无需设置（推荐）${NC}"
    echo -e "${YELLOW}   - 直接访问后端: 设置 VITE_API_BASE_URL=http://your-server:8877${NC}"
fi

# 构建前端
npm run build

if [ $? -ne 0 ]; then
    echo -e "${RED}✗ 前端构建失败${NC}"
    exit 1
fi

# 复制前端静态文件
echo -e "${GREEN}✓ 前端构建完成${NC}"
echo -e "${YELLOW}📋 复制前端静态文件...${NC}"
mkdir -p "$BUILD_DIR/frontend"
cp -r "$FRONTEND_DIR/dist"/* "$BUILD_DIR/frontend/"

# 构建后端
echo -e "\n${BLUE}🔨 构建后端...${NC}"
cd "$BACKEND_DIR"

# 检查 Go 环境
if ! command -v go &> /dev/null; then
    echo -e "${RED}✗ 未找到 Go，请先安装 Go 1.21+${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Go 版本: $(go version | awk '{print $3}')${NC}"

# 构建后端（Linux amd64）
echo -e "${YELLOW}正在编译后端（Linux amd64）...${NC}"
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build \
    -ldflags="-w -s" \
    -o "$BUILD_DIR/explore" \
    explore.go

if [ $? -ne 0 ]; then
    echo -e "${RED}✗ 后端构建失败${NC}"
    exit 1
fi

echo -e "${GREEN}✓ 后端构建完成${NC}"

# 复制后端配置文件
echo -e "${YELLOW}📋 复制后端配置文件...${NC}"
mkdir -p "$BUILD_DIR/etc"
cp -r "$BACKEND_DIR/etc"/* "$BUILD_DIR/etc/"

# 创建部署脚本
echo -e "${YELLOW}📝 创建部署脚本...${NC}"
cat > "$BUILD_DIR/deploy.sh" << 'DEPLOY_EOF'
#!/bin/bash
# TanGo 部署脚本

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  TanGo 服务部署脚本${NC}"
echo -e "${BLUE}========================================${NC}\n"

# 检查 .env 文件
if [ ! -f ".env" ]; then
    echo -e "${YELLOW}警告: .env 文件不存在${NC}"
    if [ -f ".env.example" ]; then
        echo -e "${BLUE}提示: 发现 .env.example 文件，是否复制为 .env? (y/n)${NC}"
        read -r answer
        if [ "$answer" = "y" ] || [ "$answer" = "Y" ]; then
            cp .env.example .env
            echo -e "${GREEN}已复制 .env.example 为 .env${NC}"
            echo -e "${YELLOW}请编辑 .env 文件，填入实际的配置值${NC}"
            exit 1
        fi
    fi
    echo -e "${RED}错误: 需要 .env 配置文件${NC}"
    exit 1
fi

# 创建必要的目录
echo -e "${YELLOW}📁 创建必要的目录...${NC}"
mkdir -p logs

# 设置可执行权限
chmod +x explore

# 检查端口是否被占用
BACKEND_PORT=${BACKEND_PORT:-8877}
if lsof -ti:${BACKEND_PORT} > /dev/null 2>&1; then
    echo -e "${YELLOW}警告: 端口 ${BACKEND_PORT} 已被占用${NC}"
    echo -e "${YELLOW}提示: 是否停止现有服务并继续? (y/n)${NC}"
    read -r answer
    if [ "$answer" = "y" ] || [ "$answer" = "Y" ]; then
        lsof -ti:${BACKEND_PORT} | xargs kill -9 2>/dev/null || true
        sleep 1
        echo -e "${GREEN}端口已清理${NC}"
    else
        echo -e "${RED}部署已取消${NC}"
        exit 1
    fi
fi

# 加载环境变量
echo -e "${YELLOW}📋 加载环境变量...${NC}"
set -a
source .env 2>/dev/null || true
set +a

# 启动服务（后台运行）
echo -e "${YELLOW}🚀 启动服务...${NC}"
nohup ./explore -f etc/explore.yaml > logs/explore.log 2>&1 &

# 等待服务启动
sleep 3

# 检查服务是否启动成功
if ps -p $! > /dev/null 2>&1; then
    echo -e "${GREEN}✅ TanGo 服务已启动${NC}"
    echo -e "${BLUE}📋 查看日志: tail -f logs/explore.log${NC}"
    echo -e "${BLUE}🛑 停止服务: pkill -f explore${NC}"
    echo -e "${BLUE}🌐 后端服务地址: http://localhost:${BACKEND_PORT}${NC}"
else
    echo -e "${RED}✗ 服务启动失败，请查看日志: logs/explore.log${NC}"
    if [ -f "logs/explore.log" ]; then
        echo -e "${YELLOW}最后几行日志:${NC}"
        tail -n 10 logs/explore.log
    fi
    exit 1
fi
DEPLOY_EOF

chmod +x "$BUILD_DIR/deploy.sh"

# 创建停止脚本
cat > "$BUILD_DIR/stop.sh" << 'STOP_EOF'
#!/bin/bash
# TanGo 停止脚本

echo "🛑 正在停止 TanGo 服务..."

# 查找并停止 explore 进程
pkill -f explore || true

sleep 1

# 检查是否还有进程在运行
if pgrep -f explore > /dev/null; then
    echo "⚠️  强制停止服务..."
    pkill -9 -f explore || true
fi

echo "✅ TanGo 服务已停止"
STOP_EOF

chmod +x "$BUILD_DIR/stop.sh"

# 创建 .env.example
echo -e "${YELLOW}📝 创建 .env.example...${NC}"
cat > "$BUILD_DIR/.env.example" << 'ENV_EOF'
# 后端服务配置
BACKEND_HOST=0.0.0.0
BACKEND_PORT=8877

# AI 模型配置
EINO_BASE_URL=
APP_ID=
APP_KEY=
USE_AI_MODEL=true
IMAGE_RECOGNITION_MODELS=
INTENT_MODEL=
IMAGE_GENERATION_MODEL=
TEXT_GENERATION_MODEL=

# GitHub 图片上传配置
GITHUB_TOKEN=
GITHUB_OWNER=
GITHUB_REPO=
GITHUB_BRANCH=main
GITHUB_PATH=images/
ENV_EOF

# 创建 README
echo -e "${YELLOW}📝 创建部署说明...${NC}"
cat > "$BUILD_DIR/README.md" << 'README_EOF'
# TanGo 部署说明

## 文件说明

- `explore` - 后端可执行文件
- `etc/` - 后端配置文件目录
- `frontend/` - 前端静态文件目录
- `deploy.sh` - 部署启动脚本
- `stop.sh` - 停止服务脚本
- `.env.example` - 环境变量配置示例

## 部署步骤

### 1. 配置环境变量

```bash
cp .env.example .env
# 编辑 .env 文件，填入实际的配置值
vim .env
```

### 2. 启动后端服务

```bash
./deploy.sh
```

### 3. 配置 Nginx（可选）

如果使用 Nginx 提供前端静态文件服务，请参考 `nginx.conf.example` 配置文件。

### 4. 停止服务

```bash
./stop.sh
```

## 目录结构

```
.
├── explore          # 后端可执行文件
├── etc/             # 后端配置
│   └── explore.yaml
├── frontend/        # 前端静态文件
│   ├── index.html
│   └── assets/
├── logs/            # 日志目录（自动创建）
├── deploy.sh        # 部署脚本
├── stop.sh          # 停止脚本
└── .env             # 环境变量配置
```

## 注意事项

1. 确保服务器有执行权限：`chmod +x explore deploy.sh stop.sh`
2. 确保端口 8877 未被占用（或修改 .env 中的 BACKEND_PORT）
3. 查看日志：`tail -f logs/explore.log`
4. 如果使用 Nginx，需要配置反向代理到后端服务
README_EOF

echo -e "\n${GREEN}========================================${NC}"
echo -e "${GREEN}✅ 构建完成！${NC}"
echo -e "${GREEN}========================================${NC}\n"
echo -e "${BLUE}📦 构建产物位于: ${BUILD_DIR}${NC}\n"
echo -e "${YELLOW}部署步骤：${NC}"
echo -e "1. 将 ${BUILD_DIR} 目录上传到服务器"
echo -e "2. 在服务器上创建 .env 文件（参考 .env.example）"
echo -e "3. 运行 ./deploy.sh 启动服务"
echo -e "4. 配置 Nginx（参考 nginx.conf.example）\n"

