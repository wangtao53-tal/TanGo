#!/bin/bash

# TanGo 项目启动脚本
# 功能：读取 .env 配置文件，同时启动前端和后端服务

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 项目根目录
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$ROOT_DIR/backend"
FRONTEND_DIR="$ROOT_DIR/frontend"
ENV_FILE="$ROOT_DIR/.env"

# 清理函数
cleanup() {
    echo -e "\n${YELLOW}正在关闭服务...${NC}"
    if [ ! -z "$BACKEND_PID" ]; then
        kill $BACKEND_PID 2>/dev/null || true
    fi
    if [ ! -z "$FRONTEND_PID" ]; then
        kill $FRONTEND_PID 2>/dev/null || true
    fi
    # 清理子进程
    pkill -P $$ 2>/dev/null || true
    echo -e "${GREEN}服务已关闭${NC}"
    exit 0
}

# 注册清理函数
trap cleanup SIGINT SIGTERM

# 检查 .env 文件
check_env_file() {
    if [ ! -f "$ENV_FILE" ]; then
        echo -e "${YELLOW}警告: .env 文件不存在${NC}"
        if [ -f "$ROOT_DIR/.env.example" ]; then
            echo -e "${BLUE}提示: 发现 .env.example 文件，是否复制为 .env? (y/n)${NC}"
            read -r answer
            if [ "$answer" = "y" ] || [ "$answer" = "Y" ]; then
                cp "$ROOT_DIR/.env.example" "$ENV_FILE"
                echo -e "${GREEN}已复制 .env.example 为 .env${NC}"
                echo -e "${YELLOW}请编辑 .env 文件，填入实际的配置值${NC}"
                exit 1
            fi
        fi
        echo -e "${RED}错误: 需要 .env 配置文件${NC}"
        exit 1
    fi
    echo -e "${GREEN}✓ 找到 .env 配置文件${NC}"
}

# 加载 .env 文件
load_env() {
    if [ -f "$ENV_FILE" ]; then
        # 使用 set -a 自动导出所有变量
        set -a
        # 读取 .env 文件，忽略注释和空行
        while IFS= read -r line || [ -n "$line" ]; do
            # 跳过注释和空行
            line_trimmed=$(echo "$line" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
            if [[ "$line_trimmed" =~ ^# ]] || [[ -z "$line_trimmed" ]]; then
                continue
            fi
            # 导出环境变量
            export "$line_trimmed" 2>/dev/null || true
        done < "$ENV_FILE"
        set +a
        echo -e "${GREEN}✓ 已加载 .env 配置${NC}"
    fi
}

# 检查依赖
check_dependencies() {
    echo -e "${BLUE}检查依赖...${NC}"
    
    # 检查 Go
    if ! command -v go &> /dev/null; then
        echo -e "${RED}错误: 未找到 Go，请先安装 Go 1.21+${NC}"
        exit 1
    fi
    echo -e "${GREEN}✓ Go: $(go version | awk '{print $3}')${NC}"
    
    # 检查 Node.js
    if ! command -v node &> /dev/null; then
        echo -e "${RED}错误: 未找到 Node.js，请先安装 Node.js${NC}"
        exit 1
    fi
    echo -e "${GREEN}✓ Node.js: $(node --version)${NC}"
    
    # 检查 npm
    if ! command -v npm &> /dev/null; then
        echo -e "${RED}错误: 未找到 npm，请先安装 npm${NC}"
        exit 1
    fi
    echo -e "${GREEN}✓ npm: $(npm --version)${NC}"
}

# 安装前端依赖（如果需要）
install_frontend_deps() {
    echo -e "${BLUE}检查前端依赖...${NC}"
    cd "$FRONTEND_DIR"
    
    # 检查 node_modules 是否存在
    if [ ! -d "node_modules" ]; then
        echo -e "${YELLOW}前端依赖未安装，正在安装...${NC}"
        npm install
        if [ $? -ne 0 ]; then
            echo -e "${RED}✗ 前端依赖安装失败${NC}"
            cd "$ROOT_DIR"
            exit 1
        fi
        echo -e "${GREEN}✓ 前端依赖安装完成${NC}"
    else
        # node_modules 已存在，跳过安装以加快启动速度
        # 如果需要更新依赖，请手动运行: cd frontend && npm install
        echo -e "${GREEN}✓ 前端依赖已存在，跳过安装${NC}"
    fi
    
    cd "$ROOT_DIR"
}

# 启动后端服务
start_backend() {
    echo -e "\n${BLUE}启动后端服务...${NC}"
    cd "$BACKEND_DIR"
    
    # 读取后端端口
    BACKEND_PORT=${BACKEND_PORT:-8877}
    BACKEND_HOST=${BACKEND_HOST:-0.0.0.0}
    
    # 检查端口是否被占用
    if lsof -ti:${BACKEND_PORT} > /dev/null 2>&1; then
        echo -e "${YELLOW}警告: 端口 ${BACKEND_PORT} 已被占用${NC}"
        echo -e "${YELLOW}提示: 运行 ./stop.sh 停止现有服务，或修改 .env 中的 BACKEND_PORT${NC}"
        read -p "是否强制清理端口并继续? (y/n): " answer
        if [ "$answer" = "y" ] || [ "$answer" = "Y" ]; then
            lsof -ti:${BACKEND_PORT} | xargs kill -9 2>/dev/null || true
            sleep 1
            echo -e "${GREEN}端口已清理${NC}"
        else
            echo -e "${RED}启动已取消${NC}"
            exit 1
        fi
    fi
    
    echo -e "${GREEN}后端服务地址: ${BACKEND_HOST}:${BACKEND_PORT}${NC}"
    
    # 启动后端（后台运行）
    nohup go run explore.go -f etc/explore.yaml > "$ROOT_DIR/backend.log" 2>&1 &
    BACKEND_PID=$!
    
    # 等待后端启动
    echo -e "${YELLOW}等待后端服务启动...${NC}"
    sleep 3
    
    # 检查后端是否启动成功
    if ps -p $BACKEND_PID > /dev/null 2>&1; then
        echo -e "${GREEN}✓ 后端服务已启动 (PID: $BACKEND_PID)${NC}"
        echo -e "${BLUE}后端日志: tail -f $ROOT_DIR/backend.log${NC}"
    else
        echo -e "${RED}✗ 后端服务启动失败，请查看日志: $ROOT_DIR/backend.log${NC}"
        if [ -f "$ROOT_DIR/backend.log" ]; then
            echo -e "${YELLOW}最后几行日志:${NC}"
            tail -n 10 "$ROOT_DIR/backend.log"
        fi
        exit 1
    fi
    
    cd "$ROOT_DIR"
}

# 启动前端服务
start_frontend() {
    echo -e "\n${BLUE}启动前端服务...${NC}"
    cd "$FRONTEND_DIR"
    
    # 读取前端端口
    FRONTEND_PORT=${FRONTEND_PORT:-3000}
    
    # 检查端口是否被占用
    if lsof -ti:${FRONTEND_PORT} > /dev/null 2>&1; then
        echo -e "${YELLOW}警告: 端口 ${FRONTEND_PORT} 已被占用${NC}"
        echo -e "${YELLOW}提示: 运行 ./stop.sh 停止现有服务，或修改 .env 中的 FRONTEND_PORT${NC}"
        read -p "是否强制清理端口并继续? (y/n): " answer
        if [ "$answer" = "y" ] || [ "$answer" = "Y" ]; then
            lsof -ti:${FRONTEND_PORT} | xargs kill -9 2>/dev/null || true
            sleep 1
            echo -e "${GREEN}端口已清理${NC}"
        else
            echo -e "${RED}启动已取消${NC}"
            exit 1
        fi
    fi
    
    echo -e "${GREEN}前端服务地址: http://localhost:${FRONTEND_PORT}${NC}"
    
    # 启动前端（后台运行，使用 nohup 确保进程持续运行）
    cd "$FRONTEND_DIR"
    nohup npm run dev > "$ROOT_DIR/frontend.log" 2>&1 &
    FRONTEND_PID=$!
    cd "$ROOT_DIR"
    
    # 等待前端启动
    echo -e "${YELLOW}等待前端服务启动...${NC}"
    sleep 5
    
    # 检查前端是否启动成功（检查进程和端口）
    if ps -p $FRONTEND_PID > /dev/null 2>&1; then
        # 额外检查端口是否被监听
        if lsof -ti:${FRONTEND_PORT} > /dev/null 2>&1; then
            echo -e "${GREEN}✓ 前端服务已启动 (PID: $FRONTEND_PID, 端口: ${FRONTEND_PORT})${NC}"
            echo -e "${BLUE}前端日志: tail -f $ROOT_DIR/frontend.log${NC}"
        else
            echo -e "${YELLOW}⚠ 前端进程已启动但端口未监听，请检查日志${NC}"
            if [ -f "$ROOT_DIR/frontend.log" ]; then
                echo -e "${YELLOW}最后几行日志:${NC}"
                tail -n 10 "$ROOT_DIR/frontend.log"
            fi
        fi
    else
        echo -e "${RED}✗ 前端服务启动失败，请查看日志: $ROOT_DIR/frontend.log${NC}"
        if [ -f "$ROOT_DIR/frontend.log" ]; then
            echo -e "${YELLOW}最后几行日志:${NC}"
            tail -n 10 "$ROOT_DIR/frontend.log"
        fi
        exit 1
    fi
    
    cd "$ROOT_DIR"
}

# 显示服务信息
show_info() {
    local backend_host=${BACKEND_HOST:-localhost}
    local backend_port=${BACKEND_PORT:-8877}
    local frontend_port=${FRONTEND_PORT:-3000}
    
    echo -e "\n${GREEN}========================================${NC}"
    echo -e "${GREEN}  TanGo 服务已启动${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo -e "${BLUE}后端服务:${NC} http://${backend_host}:${backend_port}"
    echo -e "${BLUE}前端服务:${NC} http://localhost:${frontend_port}"
    echo -e "\n${YELLOW}按 Ctrl+C 停止所有服务${NC}"
    echo -e "${BLUE}查看日志:${NC}"
    echo -e "  后端: tail -f $ROOT_DIR/backend.log"
    echo -e "  前端: tail -f $ROOT_DIR/frontend.log"
    echo -e "${GREEN}========================================${NC}\n"
}

# 主函数
main() {
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}  TanGo 项目启动脚本${NC}"
    echo -e "${BLUE}========================================${NC}\n"
    
    # 检查 .env 文件
    check_env_file
    
    # 加载环境变量
    load_env
    
    # 检查依赖
    check_dependencies
    
    # 安装前端依赖（如果需要）
    install_frontend_deps
    
    # 启动后端
    start_backend
    
    # 启动前端
    start_frontend
    
    # 显示服务信息
    show_info
    
    # 持续监控服务状态
    while true; do
        sleep 5
        # 检查后端是否还在运行
        if ! ps -p $BACKEND_PID > /dev/null 2>&1; then
            echo -e "\n${RED}后端服务已停止${NC}"
            cleanup
        fi
        # 检查前端是否还在运行
        if ! ps -p $FRONTEND_PID > /dev/null 2>&1; then
            echo -e "\n${RED}前端服务已停止${NC}"
            cleanup
        fi
    done
}

# 运行主函数
main
