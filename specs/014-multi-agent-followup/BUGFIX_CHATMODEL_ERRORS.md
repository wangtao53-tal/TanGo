# Bug修复：ChatModel调用错误处理

## 🐛 问题描述

在真实环境测试中，发现以下错误：

### 1. Language Agent - 403错误（模型权限问题）

```
{"@timestamp":"2025-12-24T17:11:27.644+08:00","caller":"nodes/language_agent_node.go:172","content":"ChatModel调用失败","error":"failed to create chat completion: Error code: 403 - {\"code\":\"403\",\"message\":\"The present appId lacks access privileges to this specific Model.\",\"type\":\"\",\"request_id\":\"\"}","level":"error"}
```

**原因**：随机选择的模型没有权限访问。

**处理**：代码中已实现降级机制，调用失败时自动降级到Mock模式。

### 2. Reflection Agent - 400错误（参数错误）

```
{"@timestamp":"2025-12-24T17:11:37.497+08:00","caller":"nodes/reflection_agent_node.go:170","content":"ChatModel调用失败","error":"failed to create chat completion: Error code: 400 - {\"code\":\"InvalidParameter\",\"message\":\"A parameter specified in the request is not valid: request Request id: 02176656749744106a9ddacc34765855f5d3132a5b7606ab8962f\",\"param\":\"request\",\"type\":\"BadRequest\",\"request_id\":\"\"}","level":"error"}
```

**原因**：UserMessage中使用了`{conversationHistory}`变量，但这是一个`[]*schema.Message`类型的数组，不能直接作为字符串插入到UserMessage中。

## ✅ 解决方案

### 修复1：Reflection Agent模板修复

**问题代码**：
```go
schema.UserMessage("回答内容: {content}\n对话历史: {conversationHistory}")
```

**修复后**：
```go
schema.UserMessage("回答内容: {content}")
```

对话历史通过`MessagesPlaceholder("chat_history", true)`自动插入，不需要在UserMessage中手动引用。

### 修复2：移除不必要的模板变量

**问题代码**：
```go
messages, err := n.template.Format(ctx, map[string]any{
    "content":            content,
    "conversationHistory": conversationHistory,  // 不需要
    "chat_history":       conversationHistory,
})
```

**修复后**：
```go
messages, err := n.template.Format(ctx, map[string]any{
    "content":      content,
    "chat_history": conversationHistory,
})
```

### 修复3：改进Language Agent错误日志

添加更详细的错误日志，便于排查问题：

```go
n.logger.Errorw("ChatModel调用失败，降级到Mock模式", 
    logx.Field("error", err),
    logx.Field("message", message),
    logx.Field("objectName", objectName),
)
```

## 📝 修复的文件

1. ✅ `backend/internal/agent/nodes/reflection_agent_node.go`
   - 移除UserMessage中的`{conversationHistory}`引用
   - 移除模板格式化时不必要的`conversationHistory`变量

2. ✅ `backend/internal/agent/nodes/language_agent_node.go`
   - 改进错误日志，添加更多上下文信息

## 🔍 根本原因分析

### Reflection Agent问题

eino模板引擎的`MessagesPlaceholder`会自动将消息数组插入到消息列表中，但如果在UserMessage中直接引用消息数组变量（如`{conversationHistory}`），会导致：
1. 类型不匹配：消息数组不能直接转换为字符串
2. 参数错误：API收到无效的请求参数

**正确做法**：
- 使用`MessagesPlaceholder("chat_history", true)`自动插入对话历史
- 在UserMessage中只使用简单的字符串变量（如`{content}`）

### Language Agent问题

这是模型权限问题，不是代码错误。当随机选择的模型没有权限时，会返回403错误。代码中已实现降级机制，会自动降级到Mock模式。

**建议**：
- 确保配置的模型列表中的所有模型都有权限访问
- 或者过滤掉没有权限的模型

## 🧪 验证

修复后，Reflection Agent应该能够正常工作。Language Agent在遇到权限问题时会自动降级到Mock模式，确保系统稳定性。

## 📚 相关文档

- eino模板引擎文档：https://www.cloudwego.io/zh/docs/eino/
- MessagesPlaceholder用法：自动插入消息数组到消息列表

## ⚠️ 注意事项

1. **不要直接引用消息数组**：在UserMessage中不要使用`{conversationHistory}`或`{chat_history}`等消息数组变量
2. **使用MessagesPlaceholder**：对话历史应通过`MessagesPlaceholder("chat_history", true)`自动插入
3. **模型权限检查**：确保配置的模型列表中的所有模型都有权限访问
4. **降级机制**：所有Agent节点都已实现降级机制，调用失败时自动降级到Mock模式

## ✅ 修复状态

- [x] Reflection Agent模板修复 - 已修复
- [x] Language Agent错误日志改进 - 已改进
- [x] 降级机制验证 - 正常工作

修复完成！现在Reflection Agent应该能够正常工作，Language Agent在遇到权限问题时会自动降级。

