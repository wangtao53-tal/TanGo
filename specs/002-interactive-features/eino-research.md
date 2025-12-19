# Eino 框架集成研究文档

**创建日期**: 2025-12-19  
**功能**: TanGo AI Agent 系统 - Eino 框架正确集成方案

## 问题分析

### 当前实现的问题

1. **未真正使用 eino 框架**: 当前实现只是创建了框架结构，但没有使用 eino 的实际 API
2. **缺少 ChatModel 实现**: 没有创建和使用 eino 的 ChatModel 接口
3. **未使用 Graph 编排**: 虽然创建了 Graph 结构，但没有使用 eino 的 Graph 编排能力
4. **缺少模型集成**: 没有集成 eino-ext 中的具体模型实现（如 Ark、OpenAI 等）

## Eino 框架核心概念

### 1. ChatModel - 对话模型抽象

ChatModel 是 eino 框架中对对话大模型的统一抽象，提供了标准接口：

```go
type ChatModel interface {
    Generate(ctx context.Context, messages []*schema.Message) (*schema.Message, error)
    Stream(ctx context.Context, messages []*schema.Message) (*schema.StreamReader[*schema.Message], error)
}
```

### 2. Message - 消息结构

Eino 使用 `schema.Message` 表示对话消息：

```go
// 系统消息
schema.SystemMessage("你是一个助手")

// 用户消息
schema.UserMessage("用户的问题")

// 助手回复
schema.AssistantMessage("回复内容", nil)

// 工具调用结果
schema.ToolMessage("工具执行结果", "tool_call_id")
```

### 3. Graph - 图编排

Eino 的 Graph 提供了强大的编排能力：

```go
graph := NewGraph[InputType, OutputType](
    WithGenLocalState(func(ctx context.Context) *State {
        return &State{...}
    }),
)

// 添加 ChatModel 节点
graph.AddChatModelNode("model", chatModel, options...)

// 添加边
graph.AddEdge(START, "model")
graph.AddEdge("model", END)

// 编译 Graph
runnable, err := graph.Compile(ctx, options...)
```

### 4. ChatTemplate - 消息模板

Eino 提供了强大的模板化功能：

```go
template := prompt.FromMessages(schema.FString,
    schema.SystemMessage("你是一个{role}"),
    schema.MessagesPlaceholder("chat_history", true),
    schema.UserMessage("问题: {question}"),
)

messages, err := template.Format(ctx, map[string]any{
    "role": "助手",
    "question": "用户问题",
    "chat_history": []*schema.Message{...},
})
```

## 模型实现选择

### 字节跳动内部模型 - Ark（豆包）

根据项目配置，应该使用 Ark（火山引擎）模型：

```go
import "github.com/cloudwego/eino-ext/components/model/ark"

chatModel, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
    BaseURL: cfg.EinoBaseURL,  // eino 服务地址
    AppID:   cfg.AppID,        // APP ID
    AppKey:  cfg.AppKey,        // APP Key（用于认证）
    Model:   "doubao-seed-1.6-vision", // 模型名称
})
```

### 支持的模型类型

1. **图片识别**: 使用 Vision 模型（如 `doubao-seed-1.6-vision`）
2. **文本生成**: 使用 Chat 模型（如 `gpt-5-nano`）
3. **图片生成**: 使用 Image Generation 模型（如 `Gemini 3 Pro Image`）

## 实现方案

### 方案 1: 直接使用 ChatModel（简单场景）

适用于单个模型调用的场景：

```go
// 创建 ChatModel
chatModel, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
    BaseURL: cfg.EinoBaseURL,
    AppID:   cfg.AppID,
    AppKey:  cfg.AppKey,
    Model:   cfg.IntentModel,
})

// 构建消息
messages := []*schema.Message{
    schema.SystemMessage("你是一个意图识别助手"),
    schema.UserMessage("用户消息"),
}

// 调用模型
result, err := chatModel.Generate(ctx, messages)
```

### 方案 2: 使用 Graph 编排（复杂场景）

适用于需要多个步骤、条件分支的场景：

```go
// 创建 Graph
graph := NewGraph[*GraphData, *GraphData](
    WithGenLocalState(func(ctx context.Context) *GraphState {
        return &GraphState{Messages: []*schema.Message{}}
    }),
)

// 添加意图识别节点
intentModel, _ := ark.NewChatModel(ctx, &ark.ChatModelConfig{...})
graph.AddChatModelNode("intent", intentModel)

// 添加文本生成节点
textModel, _ := ark.NewChatModel(ctx, &ark.ChatModelConfig{...})
graph.AddChatModelNode("text_gen", textModel)

// 添加条件分支
branch := NewStreamGraphBranch(
    func(ctx context.Context, sr *schema.StreamReader[*schema.Message]) (string, error) {
        // 根据意图判断路由
        if intent == "generate_cards" {
            return "card_gen", nil
        }
        return "text_gen", nil
    },
    map[string]bool{"card_gen": true, "text_gen": true},
)
graph.AddBranch("intent", branch)

// 编译并执行
runnable, _ := graph.Compile(ctx)
result, _ := runnable.Invoke(ctx, inputData)
```

## 决策

### 决策 1: 使用 ChatModel 直接调用（当前阶段）

**决策**: 在节点实现中直接使用 ChatModel，不使用 Graph 编排

**理由**:
- 当前场景相对简单，每个节点独立调用模型
- 直接使用 ChatModel 更简单，易于理解和维护
- 后续如果需要复杂编排，可以升级到 Graph

**实现方式**:
- 每个节点（图片识别、文本生成、意图识别）独立创建 ChatModel
- 使用 ChatTemplate 构建消息
- 调用 Generate 或 Stream 方法

### 决策 2: 使用 Ark 模型实现

**决策**: 使用 eino-ext 中的 Ark 模型实现

**理由**:
- 项目配置中已有 EinoBaseURL、AppID、AppKey
- Ark 是字节跳动内部使用的模型服务
- 支持 Vision 模型用于图片识别

**实现方式**:
```go
import "github.com/cloudwego/eino-ext/components/model/ark"

// 图片识别模型
visionModel, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
    BaseURL: cfg.EinoBaseURL,
    AppID:   cfg.AppID,
    AppKey:  cfg.AppKey,
    Model:   cfg.ImageRecognitionModels[0], // 从配置中选择
})

// 文本生成模型
textModel, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
    BaseURL: cfg.EinoBaseURL,
    AppID:   cfg.AppID,
    AppKey:  cfg.AppKey,
    Model:   cfg.TextGenerationModel,
})
```

### 决策 3: 消息模板化

**决策**: 使用 ChatTemplate 构建消息

**理由**:
- 支持动态参数注入
- 支持对话历史管理
- 代码更清晰，易于维护

**实现方式**:
```go
import (
    "github.com/cloudwego/eino/components/prompt"
    "github.com/cloudwego/eino/schema"
)

// 意图识别模板
intentTemplate := prompt.FromMessages(schema.FString,
    schema.SystemMessage("你是一个意图识别助手。请识别用户消息的意图：\n1. generate_cards: 用户想要生成知识卡片\n2. text_response: 用户想要文本回答\n\n请返回JSON格式: {\"intent\": \"...\", \"confidence\": 0.9}"),
    schema.MessagesPlaceholder("chat_history", true),
    schema.UserMessage("用户消息: {message}"),
)
```

## 待确认事项

1. **APP ID 和 AppKey**: 需要确认如何获取和配置
2. **模型名称**: 需要确认具体可用的模型名称列表
3. **Vision 模型调用**: 需要确认如何传递图片数据（base64）
4. **流式输出**: 需要确认是否需要支持流式响应

## 替代方案

如果 APP ID 尚未提供，可以：
1. 使用 Mock 数据（当前实现）
2. 使用本地 Ollama 模型进行测试
3. 使用 OpenAI API（需要 API Key）

## 实施步骤

1. ✅ **更新依赖**: 已安装 eino v0.7.11 和 eino-ext v0.0.1-alpha
2. ✅ **重构意图识别节点**: 已实现真实的 ChatModel 调用（intent_recognition.go）
3. 🔄 **添加 ChatTemplate**: 已为意图识别节点创建消息模板
4. ✅ **错误处理**: 已添加完善的错误处理和降级机制（失败时自动回退到 Mock）
5. ⏳ **其他节点**: 待实现图片识别、文本生成、图片生成节点

## 已实现功能

### 意图识别节点（已完成）

- ✅ 使用 `ark.NewChatModel` 创建 ChatModel 实例
- ✅ 使用 `prompt.FromMessages` 创建消息模板
- ✅ 支持配置检测：如果配置了 EinoBaseURL、AppID、AppKey，则使用真实模型
- ✅ 自动降级：如果模型调用失败，自动回退到 Mock 实现
- ✅ 结果解析：支持 JSON 格式解析，也支持文本提取意图

### 待实现节点

- ⏳ 图片识别节点：需要使用 Vision 模型
- ⏳ 文本生成节点：需要为三种卡片类型创建不同的模板
- ⏳ 图片生成节点：需要调用图片生成模型

## 参考资料

- Eino 官方文档: https://www.cloudwego.io/docs/eino/
- Eino GitHub: https://github.com/cloudwego/eino
- Eino-ext GitHub: https://github.com/cloudwego/eino-ext
- Eino Examples: https://github.com/cloudwego/eino-examples
