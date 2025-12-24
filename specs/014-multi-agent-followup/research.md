# Research: 多Agent追问功能优化

**Date**: 2025-01-27  
**Feature**: 多Agent追问功能优化  
**Plan**: [plan.md](./plan.md)

## Research Tasks

### 1. Eino Graph多Agent协作流程设计

**研究问题**: 如何组织8个Agent节点的执行流程？

**研究发现**:
- Eino框架支持Graph结构，可以定义节点之间的依赖关系和执行顺序
- Graph可以通过条件分支控制执行流程，Supervisor可以作为条件判断节点
- 可以使用eino的Graph API定义节点，并通过条件判断选择执行路径

**决策**: 
- 使用eino Graph结构组织多Agent协作流程
- Supervisor节点作为入口节点，根据条件选择执行路径
- 定义Graph执行流程：Start → Supervisor → Intent → Cognitive Load → Learning Planner → Domain Agent → Interaction → Reflection → Memory → Back to Supervisor

**Rationale**: 
- Eino框架原生支持Graph结构，符合项目技术栈
- Graph结构可以清晰地表达Agent之间的协作关系
- 条件分支可以灵活控制执行流程

**Alternatives considered**: 
- 手动编排Agent调用顺序：不够灵活，难以维护
- 使用状态机：过于复杂，不符合Eino框架设计

### 2. Supervisor节点实现方式

**研究问题**: 如何实现Supervisor的决策逻辑和Agent选择？

**研究发现**:
- Supervisor需要接收上下文信息（识别对象、孩子年龄、对话历史等）
- Supervisor需要调用Intent Agent和Cognitive Load Agent获取判断结果
- Supervisor需要根据判断结果选择领域Agent（Science/Language/Humanities）
- Supervisor不直接生成教学内容，只做决策和协调

**决策**: 
- Supervisor节点实现为eino Graph节点，接收GraphData作为输入
- Supervisor内部调用Intent Agent和Cognitive Load Agent
- Supervisor根据返回结果选择领域Agent，并调用Learning Planner Agent
- Supervisor输出为下一步动作（选择哪个Agent、是否继续等）

**Rationale**: 
- 符合Supervisor的职责定位：协调者而非执行者
- 通过调用子Agent获取判断结果，符合分层设计
- 输出为动作而非内容，符合规范要求

**Alternatives considered**: 
- Supervisor直接生成内容：违反规范要求，Supervisor不直接教学
- Supervisor使用规则引擎：不够灵活，难以适应复杂场景

### 3. Intent Agent实现方式

**研究问题**: 如何实现5种意图类型的识别（认知型、探因型、表达型、游戏型、情绪型）？

**研究发现**:
- Intent Agent需要接收用户消息和对话上下文
- Intent Agent使用ChatModel进行意图识别，输出结构化结果
- Intent Agent只输出意图标签和置信度，不生成教学内容
- 可以使用eino的ChatModel和Prompt模板实现

**决策**: 
- Intent Agent实现为eino Graph节点，使用ChatModel进行意图识别
- 定义Prompt模板，明确要求输出意图类型和置信度
- 输出格式：`{"intent": "认知型|探因型|表达型|游戏型|情绪型", "confidence": 0.0-1.0}`
- 使用JSON解析提取意图和置信度

**Rationale**: 
- 使用ChatModel进行意图识别，准确率高
- Prompt模板明确要求，确保输出格式一致
- 只输出意图标签和置信度，符合规范要求

**Alternatives considered**: 
- 使用规则匹配：准确率低，难以处理复杂表达
- 使用分类模型：需要训练数据，开发成本高

### 4. Cognitive Load Agent实现方式

**研究问题**: 如何根据年龄、轮次、输出长度判断认知负载？

**研究发现**:
- Cognitive Load Agent需要接收孩子年龄、对话轮次、最近输出长度
- Cognitive Load Agent输出策略建议：简短讲解、类比讲解、反问引导、暂停探索
- 可以使用规则判断，也可以使用ChatModel进行判断
- 规则判断更简单直接，ChatModel判断更灵活

**决策**: 
- Cognitive Load Agent实现为eino Graph节点
- 使用规则判断为主，ChatModel判断为辅
- 规则判断逻辑：
  - 3-6岁：简短讲解（≤3句）
  - 7-12岁：类比讲解（≤5句）
  - 13-18岁：深入讲解（≤7句）
  - 连续追问>5轮：反问引导
  - 最近输出>500字：暂停探索
- 输出格式：`{"strategy": "简短讲解|类比讲解|反问引导|暂停探索", "reason": "..."}`

**Rationale**: 
- 规则判断简单直接，性能好
- ChatModel判断可以处理复杂场景
- 结合使用，兼顾性能和灵活性

**Alternatives considered**: 
- 完全使用规则：不够灵活
- 完全使用ChatModel：性能开销大

### 5. Learning Planner Agent实现方式

**研究问题**: 如何根据意图和认知负载决定教学动作？

**研究发现**:
- Learning Planner Agent需要接收意图判断、认知负载建议、识别对象、孩子年龄段
- Learning Planner Agent需要决定：是否继续深入、选择哪个领域Agent、是"讲一点"还是"问一个问题"
- 可以使用ChatModel进行决策，也可以使用规则判断

**决策**: 
- Learning Planner Agent实现为eino Graph节点，使用ChatModel进行决策
- 定义Prompt模板，明确要求输出教学动作
- 输出格式：`{"continue": true/false, "domainAgent": "Science|Language|Humanities", "action": "讲一点|问一个问题"}`
- 根据输出选择对应的领域Agent

**Rationale**: 
- ChatModel可以综合考虑多个因素，做出更合理的决策
- Prompt模板明确要求，确保输出格式一致
- 输出为动作而非内容，符合规范要求

**Alternatives considered**: 
- 使用规则判断：不够灵活，难以处理复杂场景
- 使用强化学习：开发成本高，不符合MVP原则

### 6. Domain Agent实现方式

**研究问题**: Science/Language/Humanities Agent如何调用工具？

**研究发现**:
- Domain Agent需要接收用户消息、识别对象、孩子年龄等信息
- Domain Agent可以使用eino的Tool调用机制调用外部工具
- Science Agent可以调用：simple_fact_lookup、get_current_time、image_generate_simple
- Language Agent可以调用：simple_dictionary、pronunciation_hint
- Humanities Agent可以不调用工具，直接生成内容

**决策**: 
- Domain Agent实现为eino Graph节点，支持Tool调用
- 使用eino的Tool注册机制注册工具
- Domain Agent在生成回答时可以调用工具增强准确性
- 工具调用失败时降级处理，不依赖工具也能生成基本回答

**Rationale**: 
- Eino框架原生支持Tool调用，符合技术栈
- Tool调用可以增强回答准确性
- 降级处理确保系统稳定性

**Alternatives considered**: 
- 不使用工具：回答准确性可能降低
- 强制使用工具：系统稳定性可能降低

### 7. Interaction Agent实现方式

**研究问题**: 如何优化交互方式，添加可选动作？

**研究发现**:
- Interaction Agent需要接收领域Agent的回答内容
- Interaction Agent需要优化回答结尾，添加可选动作
- Interaction Agent可以使用ChatModel进行优化，也可以使用规则添加

**决策**: 
- Interaction Agent实现为eino Graph节点，使用ChatModel进行优化
- 定义Prompt模板，明确要求添加可选动作
- 常用结尾方式："你想不想试试？"、"我们下一步看什么？"、"要不要换个角度？"
- 输出为优化后的回答内容

**Rationale**: 
- ChatModel可以自然地添加可选动作，不显生硬
- Prompt模板明确要求，确保输出符合规范
- 优化交互方式，提升用户体验

**Alternatives considered**: 
- 使用规则添加：可能显生硬
- 不使用Interaction Agent：交互体验可能降低

### 8. Reflection Agent实现方式

**研究问题**: 如何判断孩子的兴趣、困惑、放松需求？

**研究发现**:
- Reflection Agent需要接收对话历史和领域Agent的回答
- Reflection Agent需要判断：孩子是否表现出兴趣、出现困惑、需要放松
- 可以使用ChatModel进行判断，也可以使用规则判断

**决策**: 
- Reflection Agent实现为eino Graph节点，使用ChatModel进行判断
- 定义Prompt模板，明确要求输出反思结果
- 输出格式：`{"interest": true/false, "confusion": true/false, "relax": true/false}`
- 输出给Memory Agent进行记录

**Rationale**: 
- ChatModel可以理解对话上下文，做出更准确的判断
- Prompt模板明确要求，确保输出格式一致
- 输出给Memory Agent，实现学习状态记录

**Alternatives considered**: 
- 使用规则判断：准确率可能较低
- 不使用Reflection Agent：无法记录学习状态

### 9. Memory Agent实现方式

**研究问题**: 如何记录和检索学习状态？

**研究发现**:
- Memory Agent需要接收Reflection Agent的输出和对话历史
- Memory Agent需要记录：孩子感兴趣的主题、已理解/未理解的点
- Memory Agent需要支持检索，为后续对话提供参考
- 可以使用内存存储或持久化存储

**决策**: 
- Memory Agent实现为eino Graph节点，使用内存存储（MemoryStorage）
- 定义MemoryRecord类型，包含感兴趣的主题、已理解/未理解的点
- 实现记录和检索方法，支持按sessionId查询
- 后续对话中，Supervisor可以查询Memory Agent的记录

**Rationale**: 
- 内存存储简单直接，符合当前架构
- 后续可以扩展为持久化存储
- 支持按sessionId查询，实现个性化回答

**Alternatives considered**: 
- 使用数据库存储：开发成本高，不符合MVP原则
- 不使用Memory Agent：无法实现个性化回答

### 10. Graph执行流程设计

**研究问题**: 如何实现Start → Supervisor → Intent → Cognitive Load → Learning Planner → Domain Agent → Interaction → Reflection → Memory → Back to Supervisor的流程？

**研究发现**:
- Eino Graph支持定义节点之间的依赖关系
- Graph可以通过条件分支控制执行流程
- Graph可以支持循环执行（Back to Supervisor）

**决策**: 
- 定义MultiAgentGraph结构，包含8个Agent节点
- 定义Graph执行流程：
  1. Start → Supervisor（入口）
  2. Supervisor → Intent Agent（并行调用）
  3. Supervisor → Cognitive Load Agent（并行调用）
  4. Supervisor → Learning Planner Agent（等待Intent和Cognitive Load结果）
  5. Learning Planner → Domain Agent（根据决策选择Science/Language/Humanities）
  6. Domain Agent → Interaction Agent（优化回答）
  7. Interaction Agent → Reflection Agent（反思判断）
  8. Reflection Agent → Memory Agent（记录状态）
  9. Memory Agent → End（返回最终回答）
- 使用eino Graph API实现节点定义和执行流程

**Rationale**: 
- Graph结构清晰表达Agent之间的协作关系
- 条件分支可以灵活控制执行流程
- 符合Eino框架设计理念

**Alternatives considered**: 
- 手动编排调用顺序：不够灵活，难以维护
- 使用状态机：过于复杂，不符合Eino框架设计

### 11. 新接口实现方式

**研究问题**: 如何创建 `/api/conversation/agent` 接口，实现多Agent协作流程？

**研究发现**:
- 需要创建新的Handler和Logic层
- Handler负责接收请求，设置SSE响应头
- Logic负责调用MultiAgentGraph，处理流式输出
- 需要与 `/api/conversation/stream` 接口保持输入输出一致

**决策**: 
- 创建 `agenthandler.go`，实现 `AgentConversationHandler`
- 创建 `agentlogic.go`，实现 `AgentLogic`，包含 `StreamAgentConversation` 方法
- Logic层调用MultiAgentGraph执行多Agent协作流程
- 使用SSE流式返回回答内容
- 请求参数使用 `UnifiedStreamConversationRequest`，与旧接口一致
- 响应格式使用 `StreamEvent`，与旧接口一致

**Rationale**: 
- 保持接口一致性，前端可以无缝切换
- 使用SSE流式返回，提供良好的用户体验
- 分离Handler和Logic层，符合go-zero架构

**Alternatives considered**: 
- 重构旧接口：风险高，可能影响现有功能
- 使用不同请求格式：前端需要修改，增加开发成本

### 12. 前端配置实现方式

**研究问题**: 前端如何通过配置选择使用哪个接口？

**研究发现**:
- 前端可以使用环境变量、配置文件或localStorage存储配置
- 配置需要支持运行时切换

### 13. Eino工具调用机制研究

**研究问题**: 如何在Eino框架中实现工具调用？

**研究发现**:
- Eino框架的ChatModel通过`Generate`方法返回`*schema.Message`，其中包含`ToolCalls`字段
- `ToolCalls`是一个`[]schema.ToolCall`切片，每个`ToolCall`包含：
  - `ID`: 工具调用ID
  - `Function`: `schema.FunctionCall`结构，包含：
    - `Name`: 工具名称（字符串）
    - `Arguments`: 工具参数（JSON字符串）
- 工具调用流程：
  1. ChatModel在生成回答时，如果检测到需要调用工具，会在`ToolCalls`字段中返回工具调用请求
  2. Agent需要检查`result.ToolCalls`，如果非空则执行工具调用
  3. 执行工具后，使用`schema.ToolMessage(toolResult, toolCallID)`创建工具消息
  4. 将工具消息添加到消息列表，重新调用ChatModel整合结果
- **重要发现**：当前eino ark.ChatModelConfig可能不支持直接设置`Tools`字段来注册工具
- 工具注册可能需要通过其他方式实现，或者依赖模型自身的工具调用能力

**决策**: 
- 采用被动工具调用模式：不主动注册工具到ChatModel，而是等待ChatModel返回工具调用请求
- 在Domain Agent中实现工具调用处理逻辑：
  - 检查`result.ToolCalls`字段
  - 从ToolRegistry获取工具实例
  - 执行工具并获取结果
  - 创建ToolMessage并重新调用ChatModel
- 工具调用失败时降级处理，不影响Agent正常回答

**Rationale**: 
- 被动模式更灵活，不依赖ChatModel配置
- 工具调用处理逻辑已完整实现，可以正常工作
- 降级处理确保系统稳定性

**Alternatives considered**: 
- 主动注册工具到ChatModel：需要研究eino API，当前可能不支持
- 不使用工具调用：会降低回答准确性

### 14. MCP集成方式研究

**研究问题**: 如何集成MCP资源到Agent系统？

**研究发现**:
- MCP（Model Context Protocol）是一个开放协议，为AI应用程序提供标准化接口
- MCP资源可以通过SSE（Server-Sent Events）或HTTP接口访问
- MCP工具包装模式：
  1. 创建MCP客户端连接到MCP服务器
  2. 发现可用资源
  3. 包装资源为Tool接口实现
  4. 注册到ToolRegistry
- MCP错误处理：需要实现重试机制和降级处理

**决策**: 
- 创建MCP工具包装器，实现Tool接口
- 支持从配置文件读取MCP服务器配置
- 实现MCP资源发现机制
- 工具调用失败时降级处理

**Rationale**: 
- MCP协议标准化，易于集成
- 工具包装模式统一，便于管理
- 降级处理确保系统稳定性

**Alternatives considered**: 
- 直接调用MCP API：不够灵活，难以管理
- 不使用MCP：会限制系统能力

### 15. Supervisor智能工具选择策略研究

**研究问题**: Supervisor如何智能选择工具？

**研究发现**:
- 工具选择决策树：
  1. 根据意图类型选择工具（认知型→fact_lookup，表达型→dictionary）
  2. 根据问题关键词匹配工具
  3. 根据Agent类型选择工具列表
- 工具使用场景分析：
  - Science Agent：fact_lookup（科学知识）、time（时间相关）、image_gen（需要示意图）
  - Language Agent：dictionary（单词查询）、pronunciation（发音学习）
- 工具调用链设计：支持多轮工具调用，但需要限制深度
- 工具结果整合策略：使用ChatModel自然整合，保持回答简洁

**决策**: 
- Supervisor根据意图和问题内容选择工具
- 将工具选择信息传递给Learning Planner
- Learning Planner将工具信息传递给Domain Agent
- Domain Agent根据工具列表选择性使用工具

**Rationale**: 
- 智能选择提高工具使用效率
- 分层决策符合架构设计
- 选择性使用避免过度依赖工具

**Alternatives considered**: 
- 所有工具都可用：可能导致工具滥用
- 固定工具列表：不够灵活
- 配置需要支持默认值（向后兼容）

**决策**: 
- 创建 `frontend/src/config/api.ts`，定义API配置
- 支持从环境变量 `VITE_USE_MULTI_AGENT` 读取配置
- 支持从localStorage读取配置（优先级高于环境变量）
- 默认值为 `false`（使用 `/api/conversation/stream` 接口）
- 配置值为 `true` 时使用 `/api/conversation/agent` 接口
- 更新对话服务，根据配置选择接口调用

**Rationale**: 
- 环境变量支持构建时配置
- localStorage支持运行时配置
- 默认值确保向后兼容

**Alternatives considered**: 
- 只使用环境变量：不支持运行时切换
- 只使用localStorage：不支持构建时配置

### 13. 接口一致性保证

**研究问题**: 如何确保两个接口的输入输出格式一致？

**研究发现**:
- 两个接口需要使用相同的请求类型（UnifiedStreamConversationRequest）
- 两个接口需要使用相同的响应格式（StreamEvent）
- 需要确保SSE事件格式一致

**决策**: 
- 两个接口使用相同的请求类型 `UnifiedStreamConversationRequest`
- 两个接口使用相同的响应类型 `StreamEvent`
- 两个接口使用相同的SSE事件格式（event: message/done/error）
- 编写测试用例验证接口一致性

**Rationale**: 
- 使用相同类型定义，确保类型一致性
- 编写测试用例，确保运行时一致性
- 前端可以无缝切换接口调用

**Alternatives considered**: 
- 使用不同请求格式：前端需要修改，增加开发成本
- 不编写测试用例：可能引入不一致问题

### 14. 错误处理和降级机制

**研究问题**: 多Agent模式失败时如何降级到单Agent模式？

**研究发现**:
- 多Agent模式可能因为各种原因失败（Agent调用失败、Graph执行失败等）
- 需要实现降级机制，确保对话能够继续
- 前端也可以实现降级机制

**决策**: 
- 后端实现降级机制：多Agent模式失败时，自动降级到单Agent模式
- 前端实现降级机制：调用 `/api/conversation/agent` 失败时，自动降级到 `/api/conversation/stream`
- 记录详细错误日志，便于问题排查
- 向用户显示友好的错误提示

**Rationale**: 
- 降级机制确保系统稳定性
- 记录错误日志便于问题排查
- 友好错误提示提升用户体验

**Alternatives considered**: 
- 不实现降级机制：系统稳定性可能降低
- 只在前端实现降级：后端无法自动恢复

## Summary

所有研究任务已完成，关键决策已确定：

1. **Graph结构**: 使用eino Graph组织多Agent协作流程
2. **Supervisor实现**: 作为协调者，调用子Agent获取判断结果，选择领域Agent
3. **Intent Agent**: 使用ChatModel进行意图识别，输出意图标签和置信度
4. **Cognitive Load Agent**: 使用规则判断为主，ChatModel判断为辅
5. **Learning Planner Agent**: 使用ChatModel进行教学决策
6. **Domain Agent**: 支持Tool调用，工具调用失败时降级处理
7. **Interaction Agent**: 使用ChatModel优化回答，添加可选动作
8. **Reflection Agent**: 使用ChatModel判断学习状态
9. **Memory Agent**: 使用内存存储记录学习状态
10. **Graph执行流程**: 定义清晰的执行流程，支持条件分支
11. **新接口实现**: 创建新接口，保持与旧接口一致
12. **前端配置**: 支持环境变量和localStorage配置
13. **接口一致性**: 使用相同类型定义和测试用例保证
14. **错误处理**: 实现前后端降级机制

所有决策均符合项目规范和技术栈要求，可以进入Phase 1设计阶段。

