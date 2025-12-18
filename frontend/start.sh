#!/bin/bash

# TanGo 前端启动脚本

echo "🚀 TanGo 前端启动脚本"
echo "===================="

# 检查 Node.js 是否安装
if ! command -v node &> /dev/null; then
    echo "❌ 未检测到 Node.js"
    echo ""
    echo "请先安装 Node.js:"
    echo "  1. 访问 https://nodejs.org/ 下载安装（推荐 LTS 版本）"
    echo "  2. 或使用 Homebrew: brew install node"
    echo ""
    exit 1
fi

# 显示 Node.js 版本
echo "✅ Node.js 版本: $(node --version)"
echo "✅ npm 版本: $(npm --version)"
echo ""

# 进入前端目录
cd "$(dirname "$0")"

# 检查 node_modules 是否存在
if [ ! -d "node_modules" ]; then
    echo "📦 首次运行，正在安装依赖..."
    echo "   这可能需要几分钟时间，请耐心等待..."
    echo ""
    npm install
    
    if [ $? -ne 0 ]; then
        echo "❌ 依赖安装失败，请检查网络连接或重试"
        exit 1
    fi
    
    echo ""
    echo "✅ 依赖安装完成！"
    echo ""
else
    echo "✅ 依赖已存在，跳过安装"
    echo ""
fi

# 启动开发服务器
echo "🎉 正在启动开发服务器..."
echo "   访问地址: http://localhost:5173"
echo "   按 Ctrl+C 停止服务器"
echo ""
echo "===================="
echo ""

npm run dev
