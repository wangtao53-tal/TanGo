package types

import "time"

// SupervisorState Follow-up Supervisor状态
type SupervisorState struct {
	ObjectName         string                 `json:"objectName"`         // 当前识别对象名称
	ObjectCategory     string                 `json:"objectCategory"`     // 对象类别（自然类/生活类/人文类）
	Cards              []CardContent           `json:"cards,optional"`     // 已生成的三张卡片（科学、诗词、英语）
	UserAge            int                    `json:"userAge"`            // 孩子年龄/年级（3-18岁）
	ConversationRounds int                    `json:"conversationRounds"` // 最近对话轮数
	RecentOutputLength int                    `json:"recentOutputLength"` // 最近输出长度（字符数）
	AgentResults       map[string]interface{} `json:"agentResults"`      // 子Agent的返回结果
	SessionId          string                 `json:"sessionId"`         // 会话ID
}

// FollowUpIntentResult 多Agent系统意图识别结果（区别于原有的IntentResult）
type FollowUpIntentResult struct {
	Intent     string  `json:"intent"`               // 意图类型：认知型、探因型、表达型、游戏型、情绪型
	Confidence float64 `json:"confidence"`           // 置信度（0.0-1.0）
	Reason     string  `json:"reason,optional"`      // 识别原因（可选）
}

// CognitiveLoadAdvice 认知负载建议
type CognitiveLoadAdvice struct {
	Strategy     string `json:"strategy"`     // 输出策略：简短讲解、类比讲解、反问引导、暂停探索
	Reason       string `json:"reason"`      // 建议理由
	MaxSentences int    `json:"maxSentences"` // 最大句子数（3/5/7）
}

// LearningPlanDecision 学习计划决策
type LearningPlanDecision struct {
	Continue     bool     `json:"continue"`      // 是否继续深入
	DomainAgent  string   `json:"domainAgent"`   // 选择的领域Agent：Science、Language、Humanities
	Action       string   `json:"action"`        // 教学动作：讲一点、问一个问题
	Tools        []string `json:"tools,optional"` // 推荐的工具列表（可选）
	ToolStrategy string   `json:"toolStrategy,optional"` // 工具使用策略（可选）：direct/enhance/none/multiple
}

// DomainAgentResponse 领域Agent回答
type DomainAgentResponse struct {
	DomainType  string                 `json:"domainType"`  // 领域类型：Science、Language、Humanities
	Content     string                 `json:"content"`     // 回答内容
	ToolsUsed   []string               `json:"toolsUsed"`    // 使用的工具列表
	ToolResults map[string]interface{} `json:"toolResults"` // 工具调用结果
}

// InteractionOptimization 交互优化结果
type InteractionOptimization struct {
	OptimizedContent string `json:"optimizedContent"` // 优化后的回答内容
	EndingAction     string `json:"endingAction"`     // 结尾动作（你想不想试试？、我们下一步看什么？、要不要换个角度？）
}

// ReflectionResult 反思结果
type ReflectionResult struct {
	Interest  bool `json:"interest"`  // 是否表现出兴趣
	Confusion bool `json:"confusion"`  // 是否出现困惑
	Relax     bool `json:"relax"`     // 是否需要放松
}

// MemoryRecord 记忆记录
type MemoryRecord struct {
	SessionId         string    `json:"sessionId"`         // 会话ID
	InterestedTopics  []string  `json:"interestedTopics"`  // 感兴趣的主题列表
	UnderstoodPoints  []string  `json:"understoodPoints"`  // 已理解的点列表
	UnunderstoodPoints []string  `json:"ununderstoodPoints"` // 未理解的点列表
	UpdatedAt         time.Time `json:"updatedAt"`         // 更新时间
}

// GraphExecutionState Graph执行状态
type GraphExecutionState struct {
	CurrentNode        string                 `json:"currentNode"`        // 当前执行的Agent节点
	ExecutionPath      []string               `json:"executionPath"`       // 执行路径（节点列表）
	IntermediateResults map[string]interface{} `json:"intermediateResults"` // 中间结果
	ErrorState         string                 `json:"errorState,optional"` // 错误状态（可选）
	StartTime          time.Time               `json:"startTime"`          // 开始时间
	EndTime            *time.Time              `json:"endTime,optional"`    // 结束时间（可选）
}

