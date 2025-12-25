# Data Model: 多Agent追问功能优化

**Date**: 2025-01-27  
**Feature**: 多Agent追问功能优化  
**Plan**: [plan.md](./plan.md)

## Entities

### SupervisorState (Follow-up Supervisor状态)

**用途**: 存储Supervisor的上下文状态，用于决策和Agent选择

**字段**:
- `ObjectName` (string): 当前识别对象名称
- `ObjectCategory` (string): 对象类别（自然类/生活类/人文类）
- `Cards` ([]CardContent): 已生成的三张卡片（科学、诗词、英语）
- `UserAge` (int): 孩子年龄/年级（3-18岁）
- `ConversationRounds` (int): 最近对话轮数
- `RecentOutputLength` (int): 最近输出长度（字符数）
- `AgentResults` (map[string]interface{}): 子Agent的返回结果
- `SessionId` (string): 会话ID

**验证规则**:
- UserAge必须在3-18之间
- ConversationRounds必须≥0
- RecentOutputLength必须≥0

### IntentResult (意图识别结果)

**用途**: 存储Intent Agent的识别结果

**字段**:
- `Intent` (string): 意图类型（认知型、探因型、表达型、游戏型、情绪型）
- `Confidence` (float64): 置信度（0.0-1.0）
- `Reason` (string, optional): 识别原因（可选）

**验证规则**:
- Intent必须是以下值之一：认知型、探因型、表达型、游戏型、情绪型
- Confidence必须在0.0-1.0之间

### CognitiveLoadAdvice (认知负载建议)

**用途**: 存储Cognitive Load Agent的建议结果

**字段**:
- `Strategy` (string): 输出策略（简短讲解、类比讲解、反问引导、暂停探索）
- `Reason` (string): 建议理由
- `MaxSentences` (int): 最大句子数（3/5/7）

**验证规则**:
- Strategy必须是以下值之一：简短讲解、类比讲解、反问引导、暂停探索
- MaxSentences必须≥0

### LearningPlanDecision (学习计划决策)

**用途**: 存储Learning Planner Agent的决策结果

**字段**:
- `Continue` (bool): 是否继续深入
- `DomainAgent` (string): 选择的领域Agent（Science、Language、Humanities）
- `Action` (string): 教学动作（讲一点、问一个问题）

**验证规则**:
- DomainAgent必须是以下值之一：Science、Language、Humanities
- Action必须是以下值之一：讲一点、问一个问题

### DomainAgentResponse (领域Agent回答)

**用途**: 存储Domain Agent（Science/Language/Humanities）的回答结果

**字段**:
- `DomainType` (string): 领域类型（Science、Language、Humanities）
- `Content` (string): 回答内容
- `ToolsUsed` ([]string): 使用的工具列表
- `ToolResults` (map[string]interface{}): 工具调用结果

**验证规则**:
- DomainType必须是以下值之一：Science、Language、Humanities
- Content不能为空

### InteractionOptimization (交互优化结果)

**用途**: 存储Interaction Agent的优化结果

**字段**:
- `OptimizedContent` (string): 优化后的回答内容
- `EndingAction` (string): 结尾动作（你想不想试试？、我们下一步看什么？、要不要换个角度？）

**验证规则**:
- OptimizedContent不能为空
- EndingAction可以为空（可选）

### ReflectionResult (反思结果)

**用途**: 存储Reflection Agent的反思结果

**字段**:
- `Interest` (bool): 是否表现出兴趣
- `Confusion` (bool): 是否出现困惑
- `Relax` (bool): 是否需要放松

**验证规则**:
- 无特殊验证规则

### MemoryRecord (记忆记录)

**用途**: 存储Memory Agent记录的学习状态

**字段**:
- `SessionId` (string): 会话ID
- `InterestedTopics` ([]string): 感兴趣的主题列表
- `UnderstoodPoints` ([]string): 已理解的点列表
- `UnunderstoodPoints` ([]string): 未理解的点列表
- `UpdatedAt` (time.Time): 更新时间

**验证规则**:
- SessionId不能为空
- UpdatedAt不能为空

### GraphExecutionState (Graph执行状态)

**用途**: 存储Graph执行过程中的状态信息

**字段**:
- `CurrentNode` (string): 当前执行的Agent节点
- `ExecutionPath` ([]string): 执行路径（节点列表）
- `IntermediateResults` (map[string]interface{}): 中间结果
- `ErrorState` (string, optional): 错误状态（可选）
- `StartTime` (time.Time): 开始时间
- `EndTime` (time.Time, optional): 结束时间（可选）

**验证规则**:
- CurrentNode不能为空
- StartTime不能为空

## Relationships

- SupervisorState → IntentResult: 1:1 (Supervisor调用Intent Agent)
- SupervisorState → CognitiveLoadAdvice: 1:1 (Supervisor调用Cognitive Load Agent)
- SupervisorState → LearningPlanDecision: 1:1 (Supervisor调用Learning Planner Agent)
- LearningPlanDecision → DomainAgentResponse: 1:1 (Learning Planner选择Domain Agent)
- DomainAgentResponse → InteractionOptimization: 1:1 (Domain Agent输出给Interaction Agent)
- InteractionOptimization → ReflectionResult: 1:1 (Interaction Agent输出给Reflection Agent)
- ReflectionResult → MemoryRecord: 1:1 (Reflection Agent输出给Memory Agent)
- MemoryRecord → SupervisorState: N:1 (Memory Agent记录用于后续Supervisor决策)

## State Transitions

### Graph执行流程状态转换

```
Start → Supervisor → Intent Agent → Cognitive Load Agent → Learning Planner Agent → Domain Agent → Interaction Agent → Reflection Agent → Memory Agent → End
```

### MemoryRecord状态转换

```
Created → Updated (每次Reflection Agent输出时更新)
```

## Data Volume Assumptions

- 单次对话最多20轮
- MemoryRecord按sessionId存储，每个session最多1条记录
- GraphExecutionState按请求存储，不持久化
- 所有数据存储在内存中，不持久化到数据库

