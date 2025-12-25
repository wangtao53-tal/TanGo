# 提示词优化：直接与孩子对话

## 🎯 问题描述

在真实环境测试中，发现Agent生成的回答语气不对，出现了给家长看的指导性内容，例如：
- "跟小朋友可以这样聊"
- "你:" 这样的对话示例格式
- "先说一句可以直接用的话"
- "跟小朋友可以这样聊: 1. 回答来源:"

这些表述表明Agent把自己定位成了给家长看的指导手册，而不是直接和孩子对话的AI伙伴。

## ✅ 优化方案

### 核心原则

1. **直接对话**：Agent是直接和孩子对话的AI伙伴，不是给家长看的指导手册
2. **第一人称**：使用"我"或直接称呼"你"（孩子）来对话
3. **禁止指导性语言**：不要出现"跟小朋友可以这样聊"、"你可以说"等指导性语言
4. **禁止示例格式**：不要出现"你:"这样的对话示例格式
5. **朋友语气**：像朋友一样聊天，而不是老师或家长

### 优化内容

#### 1. Science Agent

**优化前**：
```
你是 Science Agent，负责用孩子能理解的方式解释自然现象。
...
schema.UserMessage("孩子的问题: {message}\n识别对象: {objectName}（{objectCategory}）\n孩子年龄: {userAge}岁")
```

**优化后**：
```
你是 Science Agent，一个直接和孩子对话的AI伙伴，用简单有趣的方式解释科学知识。

重要规则：
- 直接回答孩子的问题，就像朋友聊天一样
- 不要出现"跟小朋友可以这样聊"、"你可以说"等指导性语言
- 不要出现"你:"这样的对话示例格式
- 用"我"或直接称呼"你"（孩子）来对话
...
schema.UserMessage("{message}")
```

#### 2. Language Agent

**优化前**：
```
你是 Language Agent。
目标：
- 让孩子"说得出口"
- 不讲语法规则
...
schema.UserMessage("孩子的问题: {message}\n识别对象: {objectName}（{objectCategory}）\n孩子年龄: {userAge}岁")
```

**优化后**：
```
你是 Language Agent，一个直接和孩子对话的AI伙伴，帮助孩子用语言表达自己的想法。

重要规则：
- 直接回答孩子的问题，就像朋友聊天一样
- 不要出现"跟小朋友可以这样聊"、"你可以说"等指导性语言
- 不要出现"你:"这样的对话示例格式
- 用"我"或直接称呼"你"（孩子）来对话
...
schema.UserMessage("{message}")
```

#### 3. Humanities Agent

**优化前**：
```
你是 Humanities Agent。
你负责把自然与文化连接起来：
...
schema.UserMessage("孩子的问题: {message}\n识别对象: {objectName}（{objectCategory}）\n孩子年龄: {userAge}岁")
```

**优化后**：
```
你是 Humanities Agent，一个直接和孩子对话的AI伙伴，把自然与文化连接起来。

重要规则：
- 直接回答孩子的问题，就像朋友聊天一样
- 不要出现"跟小朋友可以这样聊"、"你可以说"等指导性语言
- 不要出现"你:"这样的对话示例格式
- 用"我"或直接称呼"你"（孩子）来对话
...
schema.UserMessage("{message}")
```

#### 4. Interaction Agent

**优化前**：
```
你是 Interaction Agent。
你的职责是：
- 把内容说"轻"
- 给孩子一个可选动作
...
```

**优化后**：
```
你是 Interaction Agent，负责优化回答，让它更轻松友好。

重要规则：
- 直接优化回答内容，不要出现"跟小朋友可以这样聊"等指导性语言
- 不要出现"你:"这样的对话示例格式
- 把内容说"轻"，添加轻松友好的结尾
...
```

#### 5. Intent Agent

**优化前**：
```
schema.UserMessage("孩子的问题: {message}")
```

**优化后**：
```
schema.UserMessage("{message}")
```

## 📝 优化的文件

1. ✅ `backend/internal/agent/nodes/science_agent_node.go`
   - 优化SystemMessage，强调直接对话
   - 简化UserMessage，只保留message

2. ✅ `backend/internal/agent/nodes/language_agent_node.go`
   - 优化SystemMessage，强调直接对话
   - 简化UserMessage，只保留message

3. ✅ `backend/internal/agent/nodes/humanities_agent_node.go`
   - 优化SystemMessage，强调直接对话
   - 简化UserMessage，只保留message

4. ✅ `backend/internal/agent/nodes/interaction_agent_node.go`
   - 优化SystemMessage，强调直接优化回答

5. ✅ `backend/internal/agent/nodes/intent_agent_node.go`
   - 简化UserMessage，只保留message

## 🎨 优化效果对比

### 优化前（错误示例）
```
跟小朋友可以这样聊: 1. 回答来源:
"奥特曼不是美国的,也不是日本的,他是宇宙来的。"
"他住在M78星云那个地方,不是在地球上的国家。"

2. 顺着孩子的话接:
孩子:奥特曼是中国的吗?
你:"不是哦,奥特曼是宇宙来的,不是哪个国家的。"
```

### 优化后（正确示例）
```
奥特曼不是哪个国家的，他是宇宙里的英雄。他住在M78星云，不是在地球上的国家。

奥特曼是日本拍的电视剧，但故事里的奥特曼是宇宙英雄，不是日本人。你想不想了解更多关于奥特曼的故事？
```

## 🔍 关键改进点

1. **移除元信息**：不再在UserMessage中包含"孩子的问题:"、"识别对象:"等元信息
2. **强调直接对话**：在SystemMessage中明确强调"直接和孩子对话"
3. **禁止指导性语言**：明确禁止"跟小朋友可以这样聊"等表述
4. **禁止示例格式**：明确禁止"你:"这样的对话示例格式
5. **简化输入**：UserMessage只包含孩子的实际问题

## ⚠️ 注意事项

1. **保持对话自然**：优化后的回答应该像朋友聊天一样自然
2. **适龄性**：根据孩子年龄调整语言难度，但始终保持直接对话的语气
3. **探索导向**：让孩子感受到探索的乐趣，而不是学习的压力
4. **一致性**：所有Domain Agent都应该遵循相同的对话原则

## ✅ 优化状态

- [x] Science Agent提示词优化 - 已完成
- [x] Language Agent提示词优化 - 已完成
- [x] Humanities Agent提示词优化 - 已完成
- [x] Interaction Agent提示词优化 - 已完成
- [x] Intent Agent UserMessage简化 - 已完成

优化完成！现在Agent会直接和孩子对话，而不是生成给家长看的指导手册。

