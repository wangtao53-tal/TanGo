#!/bin/bash

# 多Agent系统测试脚本
# 运行所有多Agent相关的测试

set -e

echo "🧪 开始运行多Agent系统测试..."
echo ""

# 进入backend目录
cd "$(dirname "$0")/.." || exit 1

echo "📦 测试Agent节点..."
go test -v ./internal/agent/nodes/... -run "Test.*Agent" || echo "⚠️  Agent节点测试有失败项"

echo ""
echo "📦 测试Supervisor节点..."
go test -v ./internal/agent/nodes -run "TestSupervisor" || echo "⚠️  Supervisor节点测试有失败项"

echo ""
echo "📦 测试MultiAgentGraph..."
go test -v ./internal/agent -run "TestMultiAgentGraph" || echo "⚠️  MultiAgentGraph测试有失败项"

echo ""
echo "📦 测试Memory Agent存储..."
go test -v ./internal/storage -run "TestMemoryAgentStorage" || echo "⚠️  Memory Agent存储测试有失败项"

echo ""
echo "📦 测试AgentLogic..."
go test -v ./internal/logic -run "TestAgentLogic" || echo "⚠️  AgentLogic测试有失败项"

echo ""
echo "📦 测试AgentHandler..."
go test -v ./internal/handler -run "TestAgentConversationHandler" || echo "⚠️  AgentHandler测试有失败项"

echo ""
echo "✅ 多Agent系统测试完成！"
echo ""
echo "💡 提示："
echo "  - 当前测试在Mock模式下运行（未配置eino参数）"
echo "  - 如需测试真实AI模型，请配置EINO_BASE_URL、TAL_MLOPS_APP_ID、TAL_MLOPS_APP_KEY"
echo "  - 前端配置和接口切换功能需要手动测试"

