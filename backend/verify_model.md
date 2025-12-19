# 后端模型调用验证指南

## 当前实现状态

### ✅ 已实现的功能

1. **配置加载**：从 `.env` 文件正确加载配置
2. **Agent 初始化**：根据配置自动初始化 Agent 系统
3. **模型调用**：支持四种模型调用
   - 意图识别模型（Intent Recognition）
   - 图片识别模型（Image Recognition）
   - 文本生成模型（Text Generation）
   - 图片生成模型（Image Generation）

### ⚠️ 需要注意的问题

1. **认证方式**：当前使用 `AppKey` 或 `AppID` 作为 `APIKey`，但根据规范应该使用 `Bearer ${TAL_MLOPS_APP_ID}:${TAL_MLOPS_APP_KEY}` 格式。eino 框架可能已经处理了这个，需要实际测试验证。

2. **Mock 模式**：如果配置不完整或模型调用失败，会自动回退到 Mock 模式。

## 验证步骤

### 1. 检查配置是否正确加载

#### 方法一：查看启动日志

启动后端服务，查看日志输出：

```bash
cd backend
go run explore.go
```

**期望看到的日志**：
- ✅ `Agent初始化完成` - 表示 Agent 成功初始化
- ✅ `意图识别节点已初始化ChatModel` - 表示模型初始化成功
- ✅ `图片识别节点已初始化Vision ChatModel` - 表示 Vision 模型初始化成功
- ✅ `文本生成节点已初始化ChatModel` - 表示文本生成模型初始化成功
- ✅ `图片生成节点已初始化ImageGenerationModel` - 表示图片生成模型初始化成功

**如果看到以下日志，说明使用 Mock 模式**：
- ⚠️ `未配置eino，将使用Mock数据` - 配置未加载
- ⚠️ `未配置eino参数，意图识别节点将使用Mock模式` - 配置不完整
- ⚠️ `初始化ChatModel失败，将使用Mock模式` - 模型初始化失败

#### 方法二：检查环境变量

```bash
# 检查 .env 文件是否存在
ls -la .env

# 检查环境变量是否已加载（在启动脚本中）
echo $EINO_BASE_URL
echo $TAL_MLOPS_APP_ID
echo $TAL_MLOPS_APP_KEY
```

### 2. 测试 API 接口

#### 测试图片识别接口

```bash
# 准备一张测试图片（base64编码）
# 假设图片已编码为 base64 字符串

curl -X POST http://localhost:8877/api/explore/identify \
  -H "Content-Type: application/json" \
  -d '{
    "image": "data:image/jpeg;base64,/9j/4AAQSkZJRg...",
    "age": 8
  }'
```

**期望响应**：
- 如果使用真实模型：返回真实的识别结果
- 如果使用 Mock 模式：返回随机对象（如"银杏"、"苹果"等）

#### 测试知识卡片生成接口

```bash
curl -X POST http://localhost:8877/api/explore/generate-cards \
  -H "Content-Type: application/json" \
  -d '{
    "objectName": "银杏",
    "objectCategory": "自然类",
    "age": 8,
    "keywords": ["植物", "树木", "秋天"]
  }'
```

**期望响应**：
- 如果使用真实模型：返回 AI 生成的三张卡片内容
- 如果使用 Mock 模式：返回预定义的 Mock 卡片内容

### 3. 检查日志输出

查看后端日志文件（如果配置了文件日志）：

```bash
tail -f backend/logs/explore.log
```

**关键日志信息**：
- `执行图片识别` - 开始识别
- `图片识别完成（真实模型）` - 使用真实模型
- `图片识别完成（Mock）` - 使用 Mock 模式
- `ChatModel调用失败` - 模型调用失败

### 4. 验证配置完整性

检查 `.env` 文件是否包含所有必需配置：

```bash
cat .env | grep -E "EINO_BASE_URL|TAL_MLOPS_APP_ID|TAL_MLOPS_APP_KEY|INTENT_MODEL|IMAGE_RECOGNITION_MODELS|IMAGE_GENERATION_MODEL|TEXT_GENERATION_MODEL"
```

**必需配置项**：
- `EINO_BASE_URL` - eino 服务地址
- `TAL_MLOPS_APP_ID` - APP ID
- `TAL_MLOPS_APP_KEY` - APP Key
- `INTENT_MODEL` - 意图识别模型（可选，有默认值）
- `IMAGE_RECOGNITION_MODELS` - 图片识别模型列表（可选，有默认值）
- `IMAGE_GENERATION_MODEL` - 图片生成模型（可选，有默认值）
- `TEXT_GENERATION_MODEL` - 文本生成模型（可选，有默认值）

## 常见问题排查

### 问题1：Agent 未初始化

**症状**：日志显示 `未配置eino，将使用Mock数据`

**可能原因**：
1. `.env` 文件不存在或路径不正确
2. 环境变量未正确加载
3. `EINO_BASE_URL` 和 `TAL_MLOPS_APP_ID` 都为空

**解决方法**：
1. 检查 `.env` 文件是否存在
2. 确认配置项已正确填写
3. 检查 `explore.go` 中的 `.env` 文件查找逻辑

### 问题2：模型初始化失败

**症状**：日志显示 `初始化ChatModel失败，将使用Mock模式`

**可能原因**：
1. API 认证失败（AppID/AppKey 错误）
2. 模型名称不正确
3. 网络连接问题
4. eino 服务地址不正确

**解决方法**：
1. 检查 `TAL_MLOPS_APP_ID` 和 `TAL_MLOPS_APP_KEY` 是否正确
2. 检查模型名称是否在 eino 服务中可用
3. 检查网络连接和防火墙设置
4. 验证 `EINO_BASE_URL` 是否正确

### 问题3：模型调用失败但无错误日志

**症状**：API 返回 Mock 数据，但日志中没有错误信息

**可能原因**：
1. 模型调用超时
2. 返回格式不正确
3. JSON 解析失败

**解决方法**：
1. 检查日志级别设置
2. 增加详细的错误日志
3. 检查模型返回格式是否符合预期

## 验证脚本

创建一个简单的验证脚本：

```bash
#!/bin/bash
# verify_model.sh

echo "=== 验证后端模型配置 ==="

# 检查 .env 文件
if [ ! -f .env ]; then
    echo "❌ .env 文件不存在"
    exit 1
fi

echo "✅ .env 文件存在"

# 检查必需配置
source .env

if [ -z "$EINO_BASE_URL" ]; then
    echo "⚠️  EINO_BASE_URL 未配置"
else
    echo "✅ EINO_BASE_URL: $EINO_BASE_URL"
fi

if [ -z "$TAL_MLOPS_APP_ID" ]; then
    echo "⚠️  TAL_MLOPS_APP_ID 未配置"
else
    echo "✅ TAL_MLOPS_APP_ID: ${TAL_MLOPS_APP_ID:0:10}..."
fi

if [ -z "$TAL_MLOPS_APP_KEY" ]; then
    echo "⚠️  TAL_MLOPS_APP_KEY 未配置"
else
    echo "✅ TAL_MLOPS_APP_KEY: ${TAL_MLOPS_APP_KEY:0:10}..."
fi

echo ""
echo "=== 测试 API 接口 ==="

# 测试健康检查（如果有）
# curl -s http://localhost:8877/health

echo "请启动后端服务并查看日志以确认模型是否成功初始化"
```

## 下一步

1. **实际测试**：使用真实的 AppID 和 AppKey 进行测试
2. **验证认证**：确认 eino 框架是否正确处理 Bearer Token 认证
3. **性能测试**：测试模型调用的响应时间和成功率
4. **错误处理**：完善错误处理和日志记录
