#!/bin/bash

# TanGo 项目停止脚本
# 功能：停止所有运行中的前后端服务

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

# 加载 .env 文件获取端口配置
ENV_FILE="$ROOT_DIR/.env"
if [ -f "$ENV_FILE" ]; then
    set -a
    while IFS= read -r line || [ -n "$line" ]; do
        line_trimmed=$(echo "$line" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
        if [[ "$line_trimmed" =~ ^# ]] || [[ -z "$line_trimmed" ]]; then
            continue
        fi
        export "$line_trimmed" 2>/dev/null || true
    done < "$ENV_FILE"
    set +a
fi

BACKEND_PORT=${BACKEND_PORT:-8877}
FRONTEND_PORT=${FRONTEND_PORT:-3000}

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  TanGo 服务停止脚本${NC}"
echo -e "${BLUE}========================================${NC}\n"

# 停止后端服务
stop_backend() {
    echo -e "${BLUE}停止后端服务...${NC}"
    
    # 方法1: 通过端口查找进程
    BACKEND_PIDS=$(lsof -ti:${BACKEND_PORT} 2>/dev/null || true)
    
    # 方法2: 通过进程名查找
    GO_PIDS=$(pgrep -f "go run explore.go" 2>/dev/null || true)
    EXPLORE_PIDS=$(pgrep -f "explore.go" 2>/dev/null || true)
    
    # 合并所有进程ID
    ALL_PIDS=$(echo "$BACKEND_PIDS $GO_PIDS $EXPLORE_PIDS" | tr ' ' '\n' | sort -u | tr '\n' ' ')
    
    if [ -z "$ALL_PIDS" ] || [ "$ALL_PIDS" = " " ]; then
        echo -e "${YELLOW}未找到运行中的后端服务${NC}"
    else
        for pid in $ALL_PIDS; do
            if [ ! -z "$pid" ] && ps -p $pid > /dev/null 2>&1; then
                echo -e "${YELLOW}正在停止后端进程 (PID: $pid)...${NC}"
                kill $pid 2>/dev/null || true
                # 等待进程结束
                for i in {1..5}; do
                    if ! ps -p $pid > /dev/null 2>&1; then
                        break
                    fi
                    sleep 1
                done
                # 如果还在运行，强制杀死
                if ps -p $pid > /dev/null 2>&1; then
                    echo -e "${YELLOW}强制停止后端进程 (PID: $pid)...${NC}"
                    kill -9 $pid 2>/dev/null || true
                fi
            fi
        done
        echo -e "${GREEN}✓ 后端服务已停止${NC}"
    fi
    
    # 清理端口
    if lsof -ti:${BACKEND_PORT} > /dev/null 2>&1; then
        echo -e "${YELLOW}清理端口 ${BACKEND_PORT}...${NC}"
        lsof -ti:${BACKEND_PORT} | xargs kill -9 2>/dev/null || true
        sleep 1
    fi
}

# 停止前端服务
stop_frontend() {
    echo -e "${BLUE}停止前端服务...${NC}"
    
    # 方法1: 通过端口查找进程
    FRONTEND_PIDS=$(lsof -ti:${FRONTEND_PORT} 2>/dev/null || true)
    
    # 方法2: 通过进程名查找
    VITE_PIDS=$(pgrep -f "vite" 2>/dev/null || true)
    NPM_PIDS=$(pgrep -f "npm.*dev" 2>/dev/null || true)
    NODE_PIDS=$(pgrep -f "node.*vite" 2>/dev/null || true)
    
    # 合并所有进程ID
    ALL_PIDS=$(echo "$FRONTEND_PIDS $VITE_PIDS $NPM_PIDS $NODE_PIDS" | tr ' ' '\n' | sort -u | tr '\n' ' ')
    
    if [ -z "$ALL_PIDS" ] || [ "$ALL_PIDS" = " " ]; then
        echo -e "${YELLOW}未找到运行中的前端服务${NC}"
    else
        for pid in $ALL_PIDS; do
            if [ ! -z "$pid" ] && ps -p $pid > /dev/null 2>&1; then
                # 检查是否是 TanGo 相关进程
                CMD=$(ps -p $pid -o command= 2>/dev/null || echo "")
                if [[ "$CMD" =~ (vite|npm.*dev|node.*vite) ]] && [[ "$CMD" =~ "$FRONTEND_DIR" ]] || [[ "$CMD" =~ "TanGo" ]]; then
                    echo -e "${YELLOW}正在停止前端进程 (PID: $pid)...${NC}"
                    kill $pid 2>/dev/null || true
                    # 等待进程结束
                    for i in {1..5}; do
                        if ! ps -p $pid > /dev/null 2>&1; then
                            break
                        fi
                        sleep 1
                    done
                    # 如果还在运行，强制杀死
                    if ps -p $pid > /dev/null 2>&1; then
                        echo -e "${YELLOW}强制停止前端进程 (PID: $pid)...${NC}"
                        kill -9 $pid 2>/dev/null || true
                    fi
                fi
            fi
        done
        echo -e "${GREEN}✓ 前端服务已停止${NC}"
    fi
    
    # 清理端口
    if lsof -ti:${FRONTEND_PORT} > /dev/null 2>&1; then
        echo -e "${YELLOW}清理端口 ${FRONTEND_PORT}...${NC}"
        lsof -ti:${FRONTEND_PORT} | xargs kill -9 2>/dev/null || true
        sleep 1
    fi
}

# 清理日志文件（可选）
clean_logs() {
    echo -e "\n${BLUE}清理日志文件...${NC}"
    if [ -f "$ROOT_DIR/backend.log" ]; then
        read -p "是否删除 backend.log? (y/n): " answer
        if [ "$answer" = "y" ] || [ "$answer" = "Y" ]; then
            rm "$ROOT_DIR/backend.log"
            echo -e "${GREEN}✓ 已删除 backend.log${NC}"
        fi
    fi
    if [ -f "$ROOT_DIR/frontend.log" ]; then
        read -p "是否删除 frontend.log? (y/n): " answer
        if [ "$answer" = "y" ] || [ "$answer" = "Y" ]; then
            rm "$ROOT_DIR/frontend.log"
            echo -e "${GREEN}✓ 已删除 frontend.log${NC}"
        fi
    fi
}

# 主函数
main() {
    # 停止后端
    stop_backend
    
    # 停止前端
    stop_frontend
    
    # 验证端口是否已释放
    echo -e "\n${BLUE}验证端口状态...${NC}"
    if lsof -ti:${BACKEND_PORT} > /dev/null 2>&1; then
        echo -e "${RED}⚠ 端口 ${BACKEND_PORT} 仍被占用${NC}"
    else
        echo -e "${GREEN}✓ 端口 ${BACKEND_PORT} 已释放${NC}"
    fi
    
    if lsof -ti:${FRONTEND_PORT} > /dev/null 2>&1; then
        echo -e "${RED}⚠ 端口 ${FRONTEND_PORT} 仍被占用${NC}"
    else
        echo -e "${GREEN}✓ 端口 ${FRONTEND_PORT} 已释放${NC}"
    fi
    
    # 询问是否清理日志
    echo ""
    read -p "是否清理日志文件? (y/n): " answer
    if [ "$answer" = "y" ] || [ "$answer" = "Y" ]; then
        clean_logs
    fi
    
    echo -e "\n${GREEN}========================================${NC}"
    echo -e "${GREEN}  所有服务已停止${NC}"
    echo -e "${GREEN}========================================${NC}\n"
}

# 运行主函数
main
