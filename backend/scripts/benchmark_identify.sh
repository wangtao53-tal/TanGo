#!/bin/bash
# 性能测试脚本：测试 /api/explore/identify 接口性能

set -e

# 配置
API_URL="${API_URL:-http://localhost:8888/api/explore/identify}"
TEST_IMAGE_URL="${TEST_IMAGE_URL:-https://example.com/test.jpg}"
ITERATIONS="${ITERATIONS:-10}"
CONCURRENT="${CONCURRENT:-1}"

echo "=========================================="
echo "性能测试：/api/explore/identify"
echo "=========================================="
echo "API URL: $API_URL"
echo "测试图片URL: $TEST_IMAGE_URL"
echo "迭代次数: $ITERATIONS"
echo "并发数: $CONCURRENT"
echo "=========================================="
echo ""

# 检查依赖
if ! command -v curl &> /dev/null; then
    echo "错误: 需要安装 curl"
    exit 1
fi

# 测试函数
test_identify() {
    local start_time=$(date +%s%N)
    local response=$(curl -s -w "\n%{http_code}\n%{time_total}" \
        -X POST "$API_URL" \
        -H "Content-Type: application/json" \
        -d "{\"image\":\"$TEST_IMAGE_URL\",\"age\":8}" \
        2>&1)
    local end_time=$(date +%s%N)
    
    local http_code=$(echo "$response" | tail -n 2 | head -n 1)
    local time_total=$(echo "$response" | tail -n 1)
    local duration_ms=$(( (end_time - start_time) / 1000000 ))
    
    echo "$duration_ms $time_total $http_code"
}

# 执行测试
echo "开始性能测试..."
echo ""

results=()
success_count=0
total_time=0
min_time=999999
max_time=0

for i in $(seq 1 $ITERATIONS); do
    echo -n "测试 $i/$ITERATIONS... "
    result=$(test_identify)
    duration_ms=$(echo $result | awk '{print $1}')
    time_total=$(echo $result | awk '{print $2}')
    http_code=$(echo $result | awk '{print $3}')
    
    if [ "$http_code" = "200" ]; then
        success_count=$((success_count + 1))
        results+=($duration_ms)
        total_time=$((total_time + duration_ms))
        
        if [ $duration_ms -lt $min_time ]; then
            min_time=$duration_ms
        fi
        if [ $duration_ms -gt $max_time ]; then
            max_time=$duration_ms
        fi
        
        echo "✓ ${duration_ms}ms"
    else
        echo "✗ HTTP $http_code"
    fi
    
    # 避免请求过快
    sleep 1
done

echo ""
echo "=========================================="
echo "测试结果汇总"
echo "=========================================="
echo "总请求数: $ITERATIONS"
echo "成功请求: $success_count"
echo "失败请求: $((ITERATIONS - success_count))"
echo "成功率: $(awk "BEGIN {printf \"%.2f%%\", $success_count/$ITERATIONS*100}")"
echo ""

if [ $success_count -gt 0 ]; then
    avg_time=$((total_time / success_count))
    echo "响应时间统计:"
    echo "  最小: ${min_time}ms"
    echo "  最大: ${max_time}ms"
    echo "  平均: ${avg_time}ms"
    echo "  总计: ${total_time}ms"
    echo ""
    
    # 计算百分位数（简单版本）
    IFS=$'\n' sorted=($(sort -n <<<"${results[*]}"))
    p50_index=$((success_count / 2))
    p90_index=$((success_count * 9 / 10))
    p95_index=$((success_count * 95 / 100))
    
    echo "响应时间百分位数:"
    echo "  P50: ${sorted[$p50_index]}ms"
    echo "  P90: ${sorted[$p90_index]}ms"
    echo "  P95: ${sorted[$p95_index]}ms"
fi

echo "=========================================="
