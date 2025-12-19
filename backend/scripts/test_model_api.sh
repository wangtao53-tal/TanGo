#!/bin/bash
# 测试模型 API 调用的脚本

set -e

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
ENV_FILE="$ROOT_DIR/.env"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# 加载环境变量
if [ -f "$ENV_FILE" ]; then
    source "$ENV_FILE"
fi

BACKEND_HOST=${BACKEND_HOST:-localhost}
BACKEND_PORT=${BACKEND_PORT:-8877}
API_BASE="http://$BACKEND_HOST:$BACKEND_PORT"

echo -e "${BLUE}=== 测试后端模型 API ===${NC}"
echo ""

# 测试1: 图片识别接口
echo -e "${BLUE}1. 测试图片识别接口...${NC}"
echo -e "${YELLOW}   请求: POST $API_BASE/api/explore/identify${NC}"

# 使用一个简单的 base64 编码的测试图片（1x1 像素的红色 PNG）
TEST_IMAGE="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg=="

response=$(curl -s -w "\n%{http_code}" -X POST "$API_BASE/api/explore/identify" \
    -H "Content-Type: application/json" \
    -d "{
        \"image\": \"$TEST_IMAGE\",
        \"age\": 8
    }")

http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

if [ "$http_code" = "200" ]; then
    echo -e "${GREEN}✅ 请求成功 (HTTP $http_code)${NC}"
    echo -e "${BLUE}   响应:${NC}"
    echo "$body" | jq '.' 2>/dev/null || echo "$body"
    
    # 检查是否使用真实模型
    if echo "$body" | grep -q "objectName"; then
        object_name=$(echo "$body" | jq -r '.objectName' 2>/dev/null || echo "")
        if [ -n "$object_name" ]; then
            echo -e "${GREEN}   ✅ 识别结果: $object_name${NC}"
        fi
    fi
else
    echo -e "${RED}❌ 请求失败 (HTTP $http_code)${NC}"
    echo -e "${YELLOW}   响应: $body${NC}"
fi

echo ""

# 测试2: 知识卡片生成接口
echo -e "${BLUE}2. 测试知识卡片生成接口...${NC}"
echo -e "${YELLOW}   请求: POST $API_BASE/api/explore/generate-cards${NC}"

response=$(curl -s -w "\n%{http_code}" -X POST "$API_BASE/api/explore/generate-cards" \
    -H "Content-Type: application/json" \
    -d '{
        "objectName": "银杏",
        "objectCategory": "自然类",
        "age": 8,
        "keywords": ["植物", "树木", "秋天"]
    }')

http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

if [ "$http_code" = "200" ]; then
    echo -e "${GREEN}✅ 请求成功 (HTTP $http_code)${NC}"
    echo -e "${BLUE}   响应:${NC}"
    echo "$body" | jq '.' 2>/dev/null || echo "$body"
    
    # 检查卡片数量
    card_count=$(echo "$body" | jq '.cards | length' 2>/dev/null || echo "0")
    if [ "$card_count" -gt 0 ]; then
        echo -e "${GREEN}   ✅ 生成了 $card_count 张卡片${NC}"
    fi
else
    echo -e "${RED}❌ 请求失败 (HTTP $http_code)${NC}"
    echo -e "${YELLOW}   响应: $body${NC}"
fi

echo ""
echo -e "${BLUE}=== 测试完成 ===${NC}"
echo ""
echo -e "${YELLOW}提示:${NC}"
echo -e "  - 如果返回的是随机对象（如'银杏'、'苹果'等），说明使用 Mock 模式"
echo -e "  - 如果返回的是 AI 生成的内容，说明使用真实模型"
echo -e "  - 查看后端日志确认模型初始化状态"
