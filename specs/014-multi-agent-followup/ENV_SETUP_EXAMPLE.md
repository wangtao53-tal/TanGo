# 环境变量配置示例

## 📋 完整配置示例

创建 `.env` 文件（在项目根目录或 `backend/` 目录）：

```bash
# ============================================
# eino框架基础配置（必需）
# ============================================
EINO_BASE_URL=https://your-eino-api-endpoint.com
TAL_MLOPS_APP_ID=your_app_id
TAL_MLOPS_APP_KEY=your_app_key

# ============================================
# 模型列表配置（推荐）
# ============================================
# 文本生成模型列表（用于多Agent节点：Intent, Cognitive Load, Learning Planner, 
# Science, Language, Humanities, Interaction, Reflection）
# 多个模型用逗号分隔，系统会随机选择一个
TEXT_GENERATION_MODELS=gpt-5-nano,gemini-2.5-flash-preview,gpt-4o,doubao-seed-1.6vision

# 意图识别模型列表（可选，如果未设置则使用TEXT_GENERATION_MODELS）
INTENT_MODELS=gpt-5-nano,gemini-2.5-flash-preview

# 图片识别模型列表（可选，用于Vision模型）
IMAGE_RECOGNITION_MODELS=doubao-seed-1.6-vision,GLM-4.6v,gemini-3-pro-image

# 图片生成模型（单个模型）
IMAGE_GENERATION_MODEL=gemini-3-pro-image

# ============================================
# 其他配置
# ============================================
# 是否使用AI模型（默认true，false表示使用Mock数据）
USE_AI_MODEL=true
```

## 🎯 最小配置（仅多Agent功能）

如果只需要测试多Agent功能，最小配置如下：

```bash
EINO_BASE_URL=https://your-eino-api-endpoint.com
TAL_MLOPS_APP_ID=your_app_id
TAL_MLOPS_APP_KEY=your_app_key
TEXT_GENERATION_MODELS=gpt-5-nano
```

系统会使用默认模型列表作为备选。

## 🔍 验证配置

### 方法1：使用验证脚本

```bash
cd backend
bash scripts/verify_model.sh
```

### 方法2：检查启动日志

启动服务后，查看日志确认配置：

```bash
cd backend
go run explore.go
```

查找以下日志：
- ✅ `Intent Agent节点已初始化ChatModel，将使用真实模型`
- ✅ `Science Agent节点已初始化ChatModel`
- ✅ `Intent Agent模型已初始化 model=gpt-5-nano`

如果看到：
- ⚠️ `未配置eino参数，XXX Agent节点将使用Mock模式`

说明配置未生效，检查环境变量是否正确设置。

## 📝 环境变量解析逻辑

环境变量在 `backend/explore.go` 中解析：

```go
// 解析 TEXT_GENERATION_MODELS
if modelsStr := os.Getenv("TEXT_GENERATION_MODELS"); modelsStr != "" {
    cfg.AI.TextGenerationModels = strings.Split(modelsStr, ",")
    // 去除空格
    for i, model := range cfg.AI.TextGenerationModels {
        cfg.AI.TextGenerationModels[i] = strings.TrimSpace(model)
    }
}

// 解析 INTENT_MODELS
if modelsStr := os.Getenv("INTENT_MODELS"); modelsStr != "" {
    cfg.AI.IntentModels = strings.Split(modelsStr, ",")
    // 去除空格
    for i, model := range cfg.AI.IntentModels {
        cfg.AI.IntentModels[i] = strings.TrimSpace(model)
    }
}
```

## 🎲 模型选择示例

### 场景1：配置单个模型

```bash
TEXT_GENERATION_MODELS=gpt-5-nano
```

**结果**：所有Agent节点都使用 `gpt-5-nano`

### 场景2：配置多个模型

```bash
TEXT_GENERATION_MODELS=gpt-5-nano,gemini-2.5-flash-preview,gpt-4o
```

**结果**：每个Agent节点随机选择一个模型：
- Intent Agent → `gpt-5-nano`
- Cognitive Load Agent → `gemini-2.5-flash-preview`
- Learning Planner Agent → `gpt-4o`
- Science Agent → `gpt-5-nano`（可能重复）
- ...

### 场景3：未配置模型列表

```bash
# 不设置 TEXT_GENERATION_MODELS
```

**结果**：系统使用默认模型列表：
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

## 🔐 认证格式

认证使用 Bearer Token 格式：

```
APIKey = AppID + ":" + AppKey
```

例如：
- `TAL_MLOPS_APP_ID=app123`
- `TAL_MLOPS_APP_KEY=key456`
- 实际APIKey = `app123:key456`

## ⚠️ 注意事项

1. **环境变量优先级**：环境变量 > YAML配置
2. **模型列表格式**：多个模型用**逗号分隔**，不要有空格（系统会自动去除）
3. **模型名称**：确保模型名称在eino平台中可用
4. **配置位置**：`.env` 文件应放在项目根目录或 `backend/` 目录
5. **重启服务**：修改环境变量后需要重启服务才能生效

## 🐛 常见问题

### Q1: 环境变量未生效

**A**: 检查：
1. `.env` 文件位置是否正确
2. 环境变量名称是否正确（区分大小写）
3. 是否重启了服务

### Q2: 模型初始化失败

**A**: 检查：
1. `EINO_BASE_URL` 是否正确
2. `TAL_MLOPS_APP_ID` 和 `TAL_MLOPS_APP_KEY` 是否正确
3. 模型名称是否在eino平台中可用
4. 网络连接是否正常

### Q3: 如何固定使用某个模型

**A**: 配置单个模型：
```bash
TEXT_GENERATION_MODELS=gpt-5-nano
```

### Q4: 如何为不同Agent配置不同模型

**A**: 当前实现中，所有Agent共享 `TEXT_GENERATION_MODELS`。如需独立配置，需要修改代码添加独立的配置字段。

