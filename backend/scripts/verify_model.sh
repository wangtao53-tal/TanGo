#!/bin/bash
# 后端模型调用验证脚本

set -e

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
BACKEND_DIR="$ROOT_DIR/backend"
ENV_FILE="$ROOT_DIR/.env"

echo "=== 后端模型调用验证 ==="
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 检查 .env 文件
echo -e "${BLUE}1. 检查配置文件...${NC}"
if [ ! -f "$ENV_FILE" ]; then
    echo -e "${RED}❌ .env 文件不存在: $ENV_FILE${NC}"
    echo -e "${YELLOW}提示: 请从 .env.example 复制并配置${NC}"
    exit 1
fi
echo -e "${GREEN}✅ .env 文件存在${NC}"

# 加载环境变量
source "$ENV_FILE"

# 检查必需配置
echo ""
echo -e "${BLUE}2. 检查配置项...${NC}"

check_config() {
    local key=$1
    local value=${!key}
    
    if [ -z "$value" ] || [ "$value" = "your_app_id_here" ] || [ "$value" = "your_app_key_here" ]; then
        echo -e "${YELLOW}⚠️  $key 未配置或使用默认值${NC}"
        return 1
    else
        # 只显示前10个字符，保护敏感信息
        local display_value="${value:0:10}..."
        echo -e "${GREEN}✅ $key: $display_value${NC}"
        return 0
    fi
}

config_ok=true

if ! check_config "EINO_BASE_URL"; then
    config_ok=false
fi

if ! check_config "TAL_MLOPS_APP_ID"; then
    config_ok=false
fi

if ! check_config "TAL_MLOPS_APP_KEY"; then
    config_ok=false
fi

# 检查可选配置
echo ""
echo -e "${BLUE}3. 检查可选配置（使用默认值）...${NC}"

if [ -z "$INTENT_MODEL" ]; then
    echo -e "${YELLOW}⚠️  INTENT_MODEL 未配置，将使用默认值: gpt-5-nano${NC}"
else
    echo -e "${GREEN}✅ INTENT_MODEL: $INTENT_MODEL${NC}"
fi

if [ -z "$IMAGE_RECOGNITION_MODELS" ]; then
    echo -e "${YELLOW}⚠️  IMAGE_RECOGNITION_MODELS 未配置，将使用默认值${NC}"
else
    echo -e "${GREEN}✅ IMAGE_RECOGNITION_MODELS: $IMAGE_RECOGNITION_MODELS${NC}"
fi

if [ -z "$IMAGE_GENERATION_MODEL" ]; then
    echo -e "${YELLOW}⚠️  IMAGE_GENERATION_MODEL 未配置，将使用默认值: Gemini 3 Pro Image${NC}"
else
    echo -e "${GREEN}✅ IMAGE_GENERATION_MODEL: $IMAGE_GENERATION_MODEL${NC}"
fi

if [ -z "$TEXT_GENERATION_MODEL" ]; then
    echo -e "${YELLOW}⚠️  TEXT_GENERATION_MODEL 未配置，将使用默认值: gpt-5-nano${NC}"
else
    echo -e "${GREEN}✅ TEXT_GENERATION_MODEL: $TEXT_GENERATION_MODEL${NC}"
fi

# 检查后端服务
echo ""
echo -e "${BLUE}4. 检查后端服务...${NC}"

BACKEND_PORT=${BACKEND_PORT:-8877}
BACKEND_HOST=${BACKEND_HOST:-localhost}

if curl -s -f "http://$BACKEND_HOST:$BACKEND_PORT" > /dev/null 2>&1; then
    echo -e "${GREEN}✅ 后端服务正在运行: http://$BACKEND_HOST:$BACKEND_PORT${NC}"
    
    # 测试 API 接口
    echo ""
    echo -e "${BLUE}5. 测试 API 接口...${NC}"
    
    # 测试识别接口（使用 Mock 数据）
    response=$(curl -s -X POST "http://$BACKEND_HOST:$BACKEND_PORT/api/explore/identify" \
        -H "Content-Type: application/json" \
        -d '{"image":"data:image/jpeg;base64,test","age":8}' 2>&1)
    
    if echo "$response" | grep -q "objectName"; then
        echo -e "${GREEN}✅ 识别接口响应正常${NC}"
        echo -e "${BLUE}   响应示例: $(echo "$response" | head -c 100)...${NC}"
    else
        echo -e "${YELLOW}⚠️  识别接口响应异常${NC}"
        echo -e "${YELLOW}   响应: $response${NC}"
    fi
else
    echo -e "${YELLOW}⚠️  后端服务未运行: http://$BACKEND_HOST:$BACKEND_PORT${NC}"
    echo -e "${YELLOW}   请先启动后端服务: cd $BACKEND_DIR && go run explore.go${NC}"
fi

# 总结
echo ""
echo -e "${BLUE}=== 验证总结 ===${NC}"

if [ "$config_ok" = true ]; then
    echo -e "${GREEN}✅ 配置检查通过${NC}"
    echo -e "${GREEN}✅ 可以启动后端服务进行模型调用测试${NC}"
    echo ""
    echo -e "${BLUE}启动命令:${NC}"
    echo -e "  cd $BACKEND_DIR"
    echo -e "  go run explore.go"
    echo ""
    echo -e "${BLUE}查看日志确认:${NC}"
    echo -e "  - 查找 'Agent初始化完成' 表示 Agent 成功初始化"
    echo -e "  - 查找 '已初始化ChatModel' 表示模型初始化成功"
    echo -e "  - 查找 '使用Mock模式' 表示使用 Mock 数据"
else
    echo -e "${RED}❌ 配置检查未通过${NC}"
    echo -e "${YELLOW}请检查 .env 文件中的配置项${NC}"
    exit 1
fi
