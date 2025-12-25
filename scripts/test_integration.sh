#!/bin/bash
# 前后端集成测试脚本

set -e

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BACKEND_DIR="$ROOT_DIR/backend"
FRONTEND_DIR="$ROOT_DIR/frontend"
ENV_FILE="$ROOT_DIR/.env"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}=== 前后端集成测试 ===${NC}"
echo ""

# 加载环境变量
if [ -f "$ENV_FILE" ]; then
    source "$ENV_FILE"
fi

BACKEND_HOST=${BACKEND_HOST:-localhost}
BACKEND_PORT=${BACKEND_PORT:-8877}
FRONTEND_PORT=${FRONTEND_PORT:-3000}
API_BASE="http://$BACKEND_HOST:$BACKEND_PORT"

# 测试1: 检查后端服务
echo -e "${BLUE}1. 检查后端服务...${NC}"
if curl -s -f "$API_BASE" > /dev/null 2>&1; then
    echo -e "${GREEN}✅ 后端服务正在运行: $API_BASE${NC}"
else
    echo -e "${RED}❌ 后端服务未运行: $API_BASE${NC}"
    echo -e "${YELLOW}   请先启动后端服务: cd $BACKEND_DIR && go run explore.go${NC}"
    exit 1
fi

# 测试2: 检查 CORS
echo ""
echo -e "${BLUE}2. 检查 CORS 配置...${NC}"
cors_headers=$(curl -s -I -X OPTIONS "$API_BASE/api/explore/identify" \
    -H "Origin: http://localhost:$FRONTEND_PORT" \
    -H "Access-Control-Request-Method: POST" 2>&1)

if echo "$cors_headers" | grep -qi "access-control-allow-origin"; then
    echo -e "${GREEN}✅ CORS 配置正确${NC}"
else
    echo -e "${YELLOW}⚠️  CORS 配置可能有问题${NC}"
fi

# 测试3: 测试图片识别接口
echo ""
echo -e "${BLUE}3. 测试图片识别接口...${NC}"
TEST_IMAGE="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg=="

response=$(curl -s -w "\n%{http_code}" -X POST "$API_BASE/api/explore/identify" \
    -H "Content-Type: application/json" \
    -H "Origin: http://localhost:$FRONTEND_PORT" \
    -d "{
        \"image\": \"$TEST_IMAGE\",
        \"age\": 8
    }" 2>&1)

http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

if [ "$http_code" = "200" ]; then
    echo -e "${GREEN}✅ 图片识别接口正常 (HTTP $http_code)${NC}"
    
    # 检查响应字段
    if echo "$body" | grep -q "objectName"; then
        object_name=$(echo "$body" | jq -r '.objectName' 2>/dev/null || echo "")
        object_category=$(echo "$body" | jq -r '.objectCategory' 2>/dev/null || echo "")
        echo -e "${GREEN}   ✅ 识别结果: $object_name ($object_category)${NC}"
    else
        echo -e "${YELLOW}   ⚠️  响应格式可能不正确${NC}"
    fi
else
    echo -e "${RED}❌ 图片识别接口失败 (HTTP $http_code)${NC}"
    echo -e "${YELLOW}   响应: $body${NC}"
fi

# 测试4: 测试卡片生成接口
echo ""
echo -e "${BLUE}4. 测试卡片生成接口...${NC}"

response=$(curl -s -w "\n%{http_code}" -X POST "$API_BASE/api/explore/generate-cards" \
    -H "Content-Type: application/json" \
    -H "Origin: http://localhost:$FRONTEND_PORT" \
    -d '{
        "objectName": "银杏",
        "objectCategory": "自然类",
        "age": 8,
        "keywords": ["植物", "树木", "秋天"]
    }' 2>&1)

http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

if [ "$http_code" = "200" ]; then
    echo -e "${GREEN}✅ 卡片生成接口正常 (HTTP $http_code)${NC}"
    
    # 检查响应字段
    card_count=$(echo "$body" | jq '.cards | length' 2>/dev/null || echo "0")
    if [ "$card_count" -gt 0 ]; then
        echo -e "${GREEN}   ✅ 生成了 $card_count 张卡片${NC}"
        
        # 检查卡片类型
        card_types=$(echo "$body" | jq -r '.cards[].type' 2>/dev/null | sort -u | tr '\n' ' ')
        echo -e "${BLUE}   卡片类型: $card_types${NC}"
    else
        echo -e "${YELLOW}   ⚠️  未生成卡片${NC}"
    fi
else
    echo -e "${RED}❌ 卡片生成接口失败 (HTTP $http_code)${NC}"
    echo -e "${YELLOW}   响应: $body${NC}"
fi

# 测试5: 检查前端服务
echo ""
echo -e "${BLUE}5. 检查前端服务...${NC}"
FRONTEND_URL="http://localhost:$FRONTEND_PORT"

if curl -s -f "$FRONTEND_URL" > /dev/null 2>&1; then
    echo -e "${GREEN}✅ 前端服务正在运行: $FRONTEND_URL${NC}"
else
    echo -e "${YELLOW}⚠️  前端服务未运行: $FRONTEND_URL${NC}"
    echo -e "${YELLOW}   请启动前端服务: cd $FRONTEND_DIR && npm run dev${NC}"
fi

# 测试6: 检查 API 代理
echo ""
echo -e "${BLUE}6. 检查前端 API 代理...${NC}"
if [ -f "$FRONTEND_DIR/vite.config.ts" ]; then
    if grep -q "proxy" "$FRONTEND_DIR/vite.config.ts"; then
        echo -e "${GREEN}✅ Vite 代理配置已设置${NC}"
    else
        echo -e "${YELLOW}⚠️  Vite 代理配置可能缺失${NC}"
    fi
else
    echo -e "${YELLOW}⚠️  未找到 vite.config.ts${NC}"
fi

# 总结
echo ""
echo -e "${BLUE}=== 测试总结 ===${NC}"
echo -e "${GREEN}✅ 后端服务: $API_BASE${NC}"
echo -e "${GREEN}✅ 前端服务: $FRONTEND_URL${NC}"
echo ""
echo -e "${YELLOW}下一步:${NC}"
echo -e "  1. 确保后端服务正在运行"
echo -e "  2. 确保前端服务正在运行"
echo -e "  3. 在浏览器中访问前端页面进行端到端测试"
echo -e "  4. 检查浏览器控制台是否有错误"
echo -e "  5. 检查后端日志确认模型调用状态"
