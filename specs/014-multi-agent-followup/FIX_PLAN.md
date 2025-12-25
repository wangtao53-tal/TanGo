# Multi-Agent 系统错误修复计划

## 📋 问题总结

### 1. 模板格式化错误
**错误信息**：
- `could not find key: "continue"` (Learning Planner Agent)
- `could not find key: "interest"` (Reflection Agent)

**原因**：
- SystemMessage中的JSON示例使用了双大括号`{{`和`}}`来转义
- 但模板引擎仍然在尝试解析JSON示例中的键名
- 错误信息显示`could not find key: \n  "continue"`，说明模板引擎在查找一个包含换行符的键名

**影响范围**：
- Learning Planner Agent
- Reflection Agent
- Intent Agent（可能）

### 2. ChatModel调用错误：`Unknown parameter: 'input[0].name'`
**错误信息**：
- `Unknown parameter: 'input[0].name'`

**原因**：
- eino框架在处理消息时可能添加了工具调用相关的字段
- API不支持这些参数
- 消息格式中包含了不应该有的字段（如`UserInputMultiContent`、`ToolCalls`等）

**影响范围**：
- Intent Agent ✅ 已修复
- Learning Planner Agent ✅ 已修复
- Interaction Agent ✅ 已修复
- Reflection Agent ✅ 已修复
- Language Agent ✅ 已修复
- Science Agent ✅ 已修复
- Humanities Agent ✅ 已修复
- Cognitive Load Agent ✅ 已修复

### 3. 模型权限错误（403）
**错误信息**：
- `Error code: 403 - The present appId lacks access privileges to this specific Model.`

**原因**：
- AppID没有访问特定模型的权限
- 这是配置问题，不是代码问题

**影响范围**：
- Language Agent（已实现降级机制）

### 4. 参数验证错误（400）
**错误信息**：
- `Error code: 400 - A parameter specified in the request is not valid`

**原因**：
- 消息格式不正确
- 可能包含无效的参数

**影响范围**：
- Reflection Agent（已修复）

## ✅ 已完成的修复

### 1. 消息清理逻辑
为所有Agent节点添加了消息清理逻辑，确保只包含`Role`和`Content`字段：

```go
// 确保消息格式正确，移除任何可能导致工具调用错误的字段
cleanMessages := make([]*schema.Message, 0, len(messages))
for _, msg := range messages {
    if msg != nil && msg.Role != "" {
        cleanMsg := &schema.Message{
            Role:    msg.Role,
            Content: msg.Content,
        }
        cleanMessages = append(cleanMessages, cleanMsg)
    }
}
```

**已修复的Agent节点**：
- ✅ Intent Agent
- ✅ Learning Planner Agent
- ✅ Interaction Agent
- ✅ Reflection Agent
- ✅ Language Agent
- ✅ Science Agent
- ✅ Humanities Agent
- ✅ Cognitive Load Agent

### 2. SystemMessage优化
在SystemMessage中明确说明不使用工具：

- Learning Planner Agent: "不要使用任何工具，只返回JSON结果"
- Interaction Agent: "不要使用任何工具，只优化文本内容"
- Reflection Agent: "不要使用任何工具，只返回JSON结果"

## 🔧 待修复的问题

### 1. 模板格式化错误（高优先级）

**问题**：JSON示例中的键名被模板引擎误解析为模板变量

**解决方案**：
1. **方案A（推荐）**：将JSON示例移到单独的说明中，不在SystemMessage中直接包含JSON示例
2. **方案B**：使用更严格的转义方式
3. **方案C**：使用代码块格式（markdown）来包裹JSON示例

**实施步骤**：
1. 检查所有Agent节点的SystemMessage
2. 将JSON示例格式改为更安全的方式
3. 测试模板格式化是否正常

**影响的Agent节点**：
- Learning Planner Agent
- Reflection Agent
- Intent Agent（需要检查）

## 📝 修复检查清单

### Phase 1: 消息清理（已完成）
- [x] Intent Agent - 添加消息清理逻辑
- [x] Learning Planner Agent - 添加消息清理逻辑
- [x] Interaction Agent - 添加消息清理逻辑
- [x] Reflection Agent - 添加消息清理逻辑
- [x] Language Agent - 添加消息清理逻辑
- [x] Science Agent - 添加消息清理逻辑
- [x] Humanities Agent - 添加消息清理逻辑
- [x] Cognitive Load Agent - 添加消息清理逻辑

### Phase 2: SystemMessage优化（已完成）
- [x] Learning Planner Agent - 明确说明不使用工具
- [x] Interaction Agent - 明确说明不使用工具
- [x] Reflection Agent - 明确说明不使用工具

### Phase 3: 模板格式化修复（待完成）
- [ ] Learning Planner Agent - 修复JSON示例格式
- [ ] Reflection Agent - 修复JSON示例格式
- [ ] Intent Agent - 检查并修复（如果需要）

## 🎯 下一步行动

1. **立即修复模板格式化错误**
   - 修改Learning Planner Agent的JSON示例格式
   - 修改Reflection Agent的JSON示例格式
   - 测试模板格式化是否正常

2. **测试验证**
   - 运行所有Agent节点的单元测试
   - 进行集成测试
   - 检查错误日志

3. **文档更新**
   - 更新错误修复文档
   - 记录最佳实践

## 📚 参考文档

- [BUGFIX_TEMPLATE_ESCAPE.md](./BUGFIX_TEMPLATE_ESCAPE.md) - 模板转义修复文档
- [BUGFIX_CHATMODEL_ERRORS.md](./BUGFIX_CHATMODEL_ERRORS.md) - ChatModel错误修复文档
- [PROMPT_OPTIMIZATION.md](./PROMPT_OPTIMIZATION.md) - 提示词优化文档

## 🔍 技术细节

### 消息清理逻辑
```go
// 确保消息格式正确，移除任何可能导致工具调用错误的字段
cleanMessages := make([]*schema.Message, 0, len(messages))
for _, msg := range messages {
    if msg != nil && msg.Role != "" {
        cleanMsg := &schema.Message{
            Role:    msg.Role,
            Content: msg.Content,
        }
        cleanMessages = append(cleanMessages, cleanMsg)
    }
}
```

### JSON示例转义问题
当前使用的转义方式（双大括号`{{`和`}}`）在某些情况下可能不够。建议使用以下方式之一：

1. **使用代码块格式**：
```
请严格按照以下JSON格式返回（注意：这是示例，不要解析其中的变量）：
```
json
{
  "continue": true或false,
  "domainAgent": "Science|Language|Humanities",
  "action": "讲一点|问一个问题"
}
```

2. **使用更明确的说明**：
```
请严格按照以下格式返回JSON（不要解析示例中的键名）：
- continue: true或false
- domainAgent: Science|Language|Humanities
- action: 讲一点|问一个问题
```

## ✅ 修复状态总结

- ✅ **消息清理逻辑**：所有Agent节点已完成
- ✅ **SystemMessage优化**：关键Agent节点已完成
- ⏳ **模板格式化修复**：待完成

预计完成时间：1-2小时

