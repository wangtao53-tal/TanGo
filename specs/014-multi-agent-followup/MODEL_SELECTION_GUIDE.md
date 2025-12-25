# Agent模型选择机制指南

## 📋 概述

多Agent系统中的每个Agent节点都使用**随机选择模型**的策略，从配置的模型列表中随机选择一个模型进行初始化。这种设计提供了负载均衡和容错能力。

## 🔧 配置方式

### 1. 环境变量配置（推荐）

在 `.env` 文件中配置以下参数：

```bash
# eino框架基础配置
EINO_BASE_URL=https://your-eino-api-endpoint.com

# AI模型认证信息（Bearer Token格式：AppID:AppKey）
TAL_MLOPS_APP_ID=your_app_id
TAL_MLOPS_APP_KEY=your_app_key

# 文本生成模型列表（逗号分隔，用于多Agent节点）
TEXT_GENERATION_MODELS=gpt-5-nano,gemini-2.5-flash-preview,gpt-4o,doubao-seed-1.6vision

# 意图识别模型列表（可选，如果未设置则使用TEXT_GENERATION_MODELS）
INTENT_MODELS=gpt-5-nano,gemini-2.5-flash-preview

# 图片识别模型列表（可选，用于Vision模型）
IMAGE_RECOGNITION_MODELS=doubao-seed-1.6-vision,GLM-4.6v,gemini-3-pro-image

# 是否使用AI模型（默认true）
USE_AI_MODEL=true
```

### 2. YAML配置文件（可选）

在 `backend/etc/explore.yaml` 中配置：

```yaml
AI:
  EinoBaseURL: "https://your-eino-api-endpoint.com"
  AppID: "your_app_id"
  AppKey: "your_app_key"
  UseAIModel: true
```

**注意**：模型列表（`TextGenerationModels`、`IntentModels`等）**不在YAML中配置**，只能通过环境变量配置，避免类型解析问题。

## 🎲 模型选择机制

### 选择流程

每个Agent节点在初始化时，按以下优先级选择模型：

```
1. 从配置的模型列表（TextGenerationModels）中随机选择
   ↓（如果为空）
2. 从默认模型列表（GetDefaultTextGenerationModels()）中随机选择
   ↓（如果为空）
3. 使用默认模型（DefaultTextGenerationModel = "gpt-5-nano"）
```

### 随机选择算法

```go
// selectRandomModel 从模型列表中随机选择一个模型
func (n *IntentAgentNode) selectRandomModel(models []string) string {
    if len(models) == 0 {
        return ""
    }
    if len(models) == 1 {
        return models[0]  // 只有一个模型，直接返回
    }
    rand.Seed(time.Now().UnixNano())
    return models[rand.Intn(len(models))]  // 随机选择
}
```

### 模型选择示例

假设配置了以下模型列表：
```bash
TEXT_GENERATION_MODELS=gpt-5-nano,gemini-2.5-flash-preview,gpt-4o
```

每次初始化Agent节点时，会从这3个模型中**随机选择一个**：
- Intent Agent可能选择 `gpt-5-nano`
- Cognitive Load Agent可能选择 `gemini-2.5-flash-preview`
- Learning Planner Agent可能选择 `gpt-4o`
- Science Agent可能选择 `gpt-5-nano`（可能重复）

## 📊 各Agent节点的模型选择

### 1. Intent Agent（意图识别）
- **模型类型**：文本生成模型
- **模型列表**：`TextGenerationModels` → `GetDefaultTextGenerationModels()`
- **用途**：识别用户意图（认知型、探因型、表达型、游戏型、情绪型）

### 2. Cognitive Load Agent（认知负载）
- **模型类型**：文本生成模型（可选，主要用于复杂场景）
- **模型列表**：`TextGenerationModels` → `GetDefaultTextGenerationModels()`
- **用途**：判断认知负载，主要使用规则判断，ChatModel作为辅助

### 3. Learning Planner Agent（学习计划）
- **模型类型**：文本生成模型
- **模型列表**：`TextGenerationModels` → `GetDefaultTextGenerationModels()`
- **用途**：制定学习计划，选择领域Agent和教学动作

### 4. Science Agent（科学回答）
- **模型类型**：文本生成模型
- **模型列表**：`TextGenerationModels` → `GetDefaultTextGenerationModels()`
- **用途**：生成科学类回答，使用生活类比

### 5. Language Agent（语言回答）
- **模型类型**：文本生成模型
- **模型列表**：`TextGenerationModels` → `GetDefaultTextGenerationModels()`
- **用途**：生成语言类回答，帮助孩子表达

### 6. Humanities Agent（人文回答）
- **模型类型**：文本生成模型
- **模型列表**：`TextGenerationModels` → `GetDefaultTextGenerationModels()`
- **用途**：生成人文类回答，连接自然与文化

### 7. Interaction Agent（交互优化）
- **模型类型**：文本生成模型
- **模型列表**：`TextGenerationModels` → `GetDefaultTextGenerationModels()`
- **用途**：优化回答结尾，添加可选动作

### 8. Reflection Agent（反思判断）
- **模型类型**：文本生成模型
- **模型列表**：`TextGenerationModels` → `GetDefaultTextGenerationModels()`
- **用途**：判断用户兴趣、困惑、放松需求

### 9. Memory Agent（记忆记录）
- **模型类型**：无（不使用ChatModel）
- **用途**：记录学习状态，使用内存存储

## 🔍 模型初始化检查

### 初始化条件

Agent节点只有在满足以下**所有条件**时才会初始化ChatModel：

```go
if cfg.EinoBaseURL != "" && cfg.AppID != "" && cfg.AppKey != "" {
    // 尝试初始化ChatModel
    if err := node.initChatModel(ctx); err != nil {
        // 初始化失败，使用Mock模式
    } else {
        node.initialized = true  // 标记为已初始化
    }
} else {
    // 未配置eino参数，使用Mock模式
}
```

### Mock模式降级

如果以下任一情况发生，Agent节点会使用Mock模式：

1. **未配置eino参数**：`EinoBaseURL`、`AppID`、`AppKey`任一为空
2. **ChatModel初始化失败**：网络错误、认证失败等
3. **模型调用失败**：运行时调用失败时，部分Agent会降级到Mock模式

## 🧪 真实环境测试步骤

### 1. 配置环境变量

创建或编辑 `.env` 文件：

```bash
# 必需配置
EINO_BASE_URL=https://your-eino-api-endpoint.com
TAL_MLOPS_APP_ID=your_app_id
TAL_MLOPS_APP_KEY=your_app_key

# 推荐配置（模型列表）
TEXT_GENERATION_MODELS=gpt-5-nano,gemini-2.5-flash-preview,gpt-4o

# 可选配置
USE_AI_MODEL=true
```

### 2. 启动后端服务

```bash
cd backend
go run explore.go
```

### 3. 检查日志

启动后，查看日志确认各Agent节点是否成功初始化ChatModel：

```
✅ Intent Agent节点已初始化ChatModel，将使用真实模型
✅ Cognitive Load Agent节点已初始化ChatModel，将使用规则+模型判断
✅ Learning Planner Agent节点已初始化ChatModel，将使用真实模型
✅ Science Agent节点已初始化ChatModel
✅ Language Agent节点已初始化ChatModel
✅ Humanities Agent节点已初始化ChatModel
✅ Interaction Agent节点已初始化ChatModel
✅ Reflection Agent节点已初始化ChatModel
✅ Memory Agent节点已初始化（不使用ChatModel）
```

### 4. 测试多Agent对话

使用 `/api/conversation/agent` 接口测试：

```bash
curl -X POST http://localhost:8877/api/conversation/agent \
  -H "Content-Type: application/json" \
  -d '{
    "messageType": "text",
    "message": "这是什么？",
    "sessionId": "test-session-123",
    "userAge": 10,
    "identificationContext": {
      "objectName": "银杏",
      "objectCategory": "自然类",
      "confidence": 0.9
    }
  }'
```

### 5. 验证模型调用

查看日志，确认各Agent节点使用的模型：

```
Intent Agent模型已初始化 model=gpt-5-nano
Science Agent模型已初始化 model=gemini-2.5-flash-preview
Interaction Agent模型已初始化 model=gpt-4o
```

## 📝 默认模型列表

如果未配置 `TEXT_GENERATION_MODELS`，系统会使用以下默认模型列表：

```go
GetDefaultTextGenerationModels() = []string{
    "gemini-3-pro-image",
    "gpt-5-nano",
    "doubao-seededit-3-0-i2i",
    "doubao-seed-1.6vision",
    "glm-4.6v",
    "gpt-4o",
    "gemini-2.5-flash-preview",
    "gpt-5-pro",
    "gpt-5.1",
}
```

## ⚙️ 高级配置

### 为不同Agent配置不同模型列表

当前实现中，所有Agent节点共享 `TextGenerationModels` 配置。如果需要为不同Agent配置不同模型，可以：

1. **修改Agent节点代码**：为每个Agent添加独立的模型配置字段
2. **使用环境变量前缀**：如 `INTENT_MODELS`、`SCIENCE_MODELS` 等（需要修改代码支持）

### 固定模型选择（非随机）

如果需要固定使用某个模型，可以：

1. **配置单个模型**：
   ```bash
   TEXT_GENERATION_MODELS=gpt-5-nano
   ```

2. **修改选择逻辑**：将 `selectRandomModel` 改为 `selectFirstModel` 或 `selectByPriority`

## 🔐 认证方式

Agent节点使用 **Bearer Token** 格式进行认证：

```go
cfg.APIKey = AppID + ":" + AppKey
```

例如：
- `AppID = "app123"`
- `AppKey = "key456"`
- `APIKey = "app123:key456"`

## 🐛 故障排查

### 问题1：Agent节点未初始化ChatModel

**症状**：日志显示"未配置eino参数，XXX Agent节点将使用Mock模式"

**解决方案**：
1. 检查 `.env` 文件中的 `EINO_BASE_URL`、`TAL_MLOPS_APP_ID`、`TAL_MLOPS_APP_KEY` 是否配置
2. 确认环境变量已正确加载（重启服务）
3. 检查YAML配置文件中的配置是否正确

### 问题2：ChatModel初始化失败

**症状**：日志显示"初始化ChatModel失败，将使用Mock模式"

**解决方案**：
1. 检查 `EINO_BASE_URL` 是否正确
2. 检查 `TAL_MLOPS_APP_ID` 和 `TAL_MLOPS_APP_KEY` 是否正确
3. 检查网络连接是否正常
4. 检查eino API服务是否可用

### 问题3：模型调用失败

**症状**：日志显示"ChatModel调用失败"

**解决方案**：
1. 检查模型名称是否正确（在eino平台中可用）
2. 检查API配额是否充足
3. 检查网络连接和超时设置
4. 查看详细错误日志

## 📚 相关文件

- `backend/internal/config/config.go` - 配置结构定义
- `backend/internal/config/models.go` - 默认模型列表
- `backend/internal/agent/nodes/*_agent_node.go` - 各Agent节点实现
- `backend/etc/explore.yaml` - YAML配置文件示例

## 🎯 最佳实践

1. **配置多个模型**：提供负载均衡和容错能力
   ```bash
   TEXT_GENERATION_MODELS=model1,model2,model3
   ```

2. **使用环境变量**：避免在代码中硬编码配置

3. **监控模型调用**：记录每个Agent使用的模型，便于问题排查

4. **测试Mock模式**：确保Mock模式正常工作，作为降级方案

5. **定期检查日志**：确认所有Agent节点正确初始化

