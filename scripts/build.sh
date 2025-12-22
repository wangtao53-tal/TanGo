#!/bin/bash

# TanGo 本地构建脚本
# 用于在本地编译前端和后端，生成可部署的文件

set -e  # 遇到错误立即退出

echo "🔨 TanGo 本地构建脚本"
echo "===================="
echo ""

# 获取项目根目录
PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$PROJECT_ROOT"

# 构建输出目录
BUILD_DIR="$PROJECT_ROOT/build"
FRONTEND_DIST="$PROJECT_ROOT/frontend/dist"
BACKEND_BINARY="$BUILD_DIR/explore"

# 清理旧的构建文件
echo "🧹 清理旧的构建文件..."
rm -rf "$BUILD_DIR"
mkdir -p "$BUILD_DIR"

# ==================== 前端构建 ====================
echo ""
echo "📦 开始构建前端..."
echo "-------------------"

cd "$PROJECT_ROOT/frontend"

# 检查 Node.js
if ! command -v node &> /dev/null; then
    echo "❌ 未检测到 Node.js，请先安装 Node.js"
    exit 1
fi

echo "✅ Node.js 版本: $(node --version)"
echo "✅ npm 版本: $(npm --version)"

# 安装依赖（如果需要）
if [ ! -d "node_modules" ]; then
    echo "📥 安装前端依赖..."
    npm install
fi

# 构建前端
echo "🔨 构建前端静态文件..."
npm run build

if [ ! -d "dist" ]; then
    echo "❌ 前端构建失败：未找到 dist 目录"
    exit 1
fi

echo "✅ 前端构建完成：$FRONTEND_DIST"
echo ""

# ==================== 后端构建 ====================
echo "📦 开始构建后端..."
echo "-------------------"

cd "$PROJECT_ROOT/backend"

# 检查 Go
if ! command -v go &> /dev/null; then
    echo "❌ 未检测到 Go，请先安装 Go"
    exit 1
fi

echo "✅ Go 版本: $(go version)"

# 下载依赖
echo "📥 下载 Go 依赖..."
go mod download

# 交叉编译（Linux amd64）
echo "🔨 交叉编译后端（Linux amd64）..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o "$BACKEND_BINARY" \
    explore.go

if [ ! -f "$BACKEND_BINARY" ]; then
    echo "❌ 后端构建失败：未找到可执行文件"
    exit 1
fi

echo "✅ 后端构建完成：$BACKEND_BINARY"
echo ""

# ==================== 打包部署文件 ====================
echo "📦 打包部署文件..."
echo "-------------------"

# 创建部署目录结构
DEPLOY_DIR="$BUILD_DIR/deploy"
mkdir -p "$DEPLOY_DIR"

# 复制后端可执行文件
cp "$BACKEND_BINARY" "$DEPLOY_DIR/"

# 复制后端配置文件
echo "📋 复制配置文件..."
mkdir -p "$DEPLOY_DIR/etc"
cp -r "$PROJECT_ROOT/backend/etc"/* "$DEPLOY_DIR/etc/" 2>/dev/null || true

# 复制前端静态文件
echo "📋 复制前端静态文件..."
mkdir -p "$DEPLOY_DIR/frontend"
cp -r "$FRONTEND_DIST"/* "$DEPLOY_DIR/frontend/"

# 创建启动脚本
cat > "$DEPLOY_DIR/start.sh" << 'EOF'
#!/bin/bash

# TanGo 部署启动脚本

# 获取脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$SCRIPT_DIR"

# 创建日志目录
mkdir -p logs

# 设置时区
export TZ=Asia/Shanghai

# 启动服务
# 注意：需要根据实际情况修改配置文件路径
./explore -f etc/explore.yaml
EOF

chmod +x "$DEPLOY_DIR/start.sh"

# 创建部署说明文件
cat > "$DEPLOY_DIR/README.md" << 'EOF'
# TanGo 部署包

## 文件说明

- `explore`: 后端可执行文件（Linux amd64）
- `etc/`: 后端配置文件目录
- `frontend/`: 前端静态文件目录
- `start.sh`: 启动脚本

## 部署步骤

1. **上传文件到服务器**
   ```bash
   # 将整个 deploy 目录上传到服务器
   scp -r deploy/* user@server:/path/to/tango/
   ```

2. **在服务器上设置权限**
   ```bash
   chmod +x explore
   chmod +x start.sh
   ```

3. **配置环境变量**
   ```bash
   # 在服务器上创建 .env 文件（可选）
   # 或直接修改 etc/explore.yaml
   ```

4. **启动服务**
   ```bash
   ./start.sh
   # 或后台运行
   nohup ./start.sh > logs/app.log 2>&1 &
   ```

## 注意事项

- 确保服务器是 Linux amd64 架构
- 确保服务器有执行权限
- 根据实际情况修改配置文件中的端口和路径
- 建议使用 systemd 或 supervisor 管理进程
EOF

# 创建压缩包
echo "📦 创建部署压缩包..."
cd "$BUILD_DIR"
tar -czf "tango-deploy-$(date +%Y%m%d-%H%M%S).tar.gz" deploy/

echo ""
echo "✅ 构建完成！"
echo "===================="
echo ""
echo "📁 构建输出目录: $BUILD_DIR"
echo "📦 部署目录: $DEPLOY_DIR"
echo "📦 压缩包: $BUILD_DIR/tango-deploy-*.tar.gz"
echo ""
echo "🚀 部署步骤:"
echo "   1. 将 deploy 目录上传到服务器"
echo "   2. 在服务器上运行: chmod +x explore start.sh"
echo "   3. 运行: ./start.sh"
echo ""

