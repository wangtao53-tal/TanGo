# Feature Specification: 多Agent追问功能优化

**Feature Branch**: `dev-mvp-20251218`  
**Created**: 2025-01-27  
**Status**: Draft  
**Input**: User description: "优化追问功能，不是单一调用模型。 实现多agent功能。一、① Follow-up Supervisor · Prompt Spec（核心中枢）"

**Note**: MVP版本阶段，所有开发工作统一在 `dev-mvp-20251218` 分支进行，不采用一个功能一个分支的策略。

## Clarifications

### Session 2025-01-27

- **Q: 接口重构方式** → A: 核心目标是重构 `/api/conversation/stream` 接口，实现多Agent模式。为了保证可以随时切换，需要添加配置支持，可以选择多Agent模式或单一Agent模型模式。
- **Q: 配置方式** → A: 配置写到env文件中，通过环境变量控制是否使用多Agent模式。当配置为多Agent模式时，使用多Agent协作流程；当配置为单一Agent模式时，使用原有的Conversation Node直接调用模型。
- **Q: 配置字段名称** → A: 使用 `USE_MULTI_AGENT` 环境变量（布尔值），`true` 表示使用多Agent模式，`false` 表示使用单一Agent模型模式。默认值为 `false`（单一Agent模式），确保向后兼容。
- **Q: 接口分离方式** → A: `/api/conversation/stream` 保持为老的单Agent模式接口，新增 `/api/conversation/agent` 接口作为新版本多Agent模式接口。两个接口区分开，前端通过配置区分使用哪个接口来进行追问的流式回答。两个接口的输入和输出必须保证一致。

## User Scenarios & Testing *(mandatory)*

### User Story 1 - 智能追问意图识别 (Priority: P1)

孩子在对话中追问问题时，系统能够准确识别孩子的学习意图（认知型、探因型、表达型、游戏型、情绪型），并根据意图选择合适的响应策略。

**Why this priority**: 意图识别是多Agent系统的基础，只有准确识别孩子的真实意图，才能提供合适的教学内容和交互方式。这是整个追问功能优化的核心前提。

**Independent Test**: 可以独立测试：孩子在对话中提出不同类型的问题（"这是什么？"、"为什么？"、"怎么说？"、"好玩吗？"、"我不懂"），系统能够准确识别意图类型并返回置信度。这个功能可以独立交付价值，即使没有后续的Agent协作，也能帮助系统理解孩子的需求。

**Acceptance Scenarios**:

1. **Given** 孩子在对话中追问问题, **When** 孩子提出认知型问题（如"这是什么？"）, **Then** 系统识别为"认知型"意图，置信度≥80%
2. **Given** 孩子在对话中追问问题, **When** 孩子提出探因型问题（如"为什么？"、"怎么会？"）, **Then** 系统识别为"探因型"意图，置信度≥80%
3. **Given** 孩子在对话中追问问题, **When** 孩子提出表达型问题（如"怎么说？"、"怎么形容？"）, **Then** 系统识别为"表达型"意图，置信度≥80%
4. **Given** 孩子在对话中追问问题, **When** 孩子提出游戏型问题（如"好玩吗？"、"能不能试试？"）, **Then** 系统识别为"游戏型"意图，置信度≥80%
5. **Given** 孩子在对话中追问问题, **When** 孩子表现出困惑情绪（如"我不懂"、"太难了"）, **Then** 系统识别为"情绪型"意图，置信度≥80%

---

### User Story 2 - 认知负载控制 (Priority: P1)

系统根据孩子的年龄、对话轮次和最近输出长度，智能控制每轮回答的信息量，避免信息过载，确保孩子能够理解和消化。

**Why this priority**: 认知负载控制是K12教育场景的关键，不同年龄段的孩子有不同的认知能力。信息过载会导致孩子无法理解，影响学习效果。这是多Agent系统的核心价值之一。

**Independent Test**: 可以独立测试：系统根据孩子年龄（3-6岁、7-12岁、13-18岁）、对话轮次和最近输出长度，自动选择最合适的输出策略（简短讲解、类比讲解、反问引导、暂停探索）。这个功能可以独立工作，确保每轮回答都适合孩子的认知水平。

**Acceptance Scenarios**:

1. **Given** 3-6岁孩子进行对话, **When** 系统生成回答, **Then** 系统选择"简短讲解"策略，回答不超过3句话
2. **Given** 7-12岁孩子进行对话, **When** 系统生成回答, **Then** 系统选择"类比讲解"策略，回答不超过5句话
3. **Given** 13-18岁孩子进行对话, **When** 系统生成回答, **Then** 系统选择"深入讲解"策略，回答不超过7句话
4. **Given** 孩子连续追问超过5轮, **When** 系统生成回答, **Then** 系统选择"反问引导"策略，鼓励孩子思考而不是继续讲解
5. **Given** 最近输出长度超过500字, **When** 系统生成回答, **Then** 系统选择"暂停探索"策略，建议休息或换话题

---

### User Story 3 - Supervisor协调多Agent协作 (Priority: P1)

Follow-up Supervisor作为核心中枢，协调多个子Agent协作，根据孩子的意图和认知负载，选择最合适的领域Agent（Science Agent、Language Agent、Humanities Agent）生成回答。

**Why this priority**: Supervisor是多Agent系统的核心，负责协调各个子Agent的工作，确保回答的质量和一致性。没有Supervisor，各个Agent无法有效协作，系统无法提供连贯的学习体验。

**Independent Test**: 可以独立测试：孩子在对话中追问问题时，Supervisor根据意图识别结果和认知负载建议，选择最合适的领域Agent（Science Agent、Language Agent、Humanities Agent），生成符合孩子年龄和认知水平的回答。这个功能可以独立工作，确保回答的专业性和适龄性。

**Acceptance Scenarios**:

1. **Given** 孩子提出科学相关问题, **When** Supervisor协调Agent协作, **Then** Supervisor选择Science Agent生成回答，回答包含科学知识但用孩子能理解的方式表达
2. **Given** 孩子提出语言表达问题, **When** Supervisor协调Agent协作, **Then** Supervisor选择Language Agent生成回答，回答包含可模仿的句子和发音提示
3. **Given** 孩子提出人文相关问题, **When** Supervisor协调Agent协作, **Then** Supervisor选择Humanities Agent生成回答，回答包含古诗词、故事或文化背景
4. **Given** 孩子的问题涉及多个领域, **When** Supervisor协调Agent协作, **Then** Supervisor选择多个领域Agent协作，生成综合回答
5. **Given** 孩子表现出困惑情绪, **When** Supervisor协调Agent协作, **Then** Supervisor优先选择Interaction Agent，用轻松的方式引导孩子，不制造学习压力

---

### User Story 4 - 领域Agent专业回答 (Priority: P2)

Science Agent、Language Agent、Humanities Agent分别负责不同领域的专业回答，每个Agent专注于自己的领域，提供高质量的内容。

**Why this priority**: 领域Agent提供专业的内容，确保回答的准确性和深度。这是多Agent系统的价值体现，不同领域的专业Agent能够提供更专业、更深入的知识。

**Independent Test**: 可以独立测试：Science Agent回答科学问题时，使用生活类比，不超过4句话；Language Agent回答语言问题时，提供可模仿的句子和发音提示；Humanities Agent回答人文问题时，提供古诗词、故事或文化背景。每个Agent都可以独立工作，提供专业的内容。

**Acceptance Scenarios**:

1. **Given** Science Agent被调用, **When** 回答科学问题, **Then** 回答使用生活类比，避免专业术语，不超过4句话
2. **Given** Language Agent被调用, **When** 回答语言问题, **Then** 回答包含可模仿的句子和发音提示，使用孩子日常语言
3. **Given** Humanities Agent被调用, **When** 回答人文问题, **Then** 回答包含古诗词、故事或文化背景，与当前识别对象相关
4. **Given** 领域Agent需要调用工具, **When** 生成回答, **Then** Agent可以调用相应的工具（如simple_fact_lookup、simple_dictionary、pronunciation_hint），增强回答的准确性
5. **Given** 领域Agent生成回答, **When** 回答内容, **Then** 回答只解决一个认知目标，不一次性讲完完整知识体系

---

### User Story 5 - 交互优化和反思记忆 (Priority: P2)

Interaction Agent负责优化交互方式，用轻松的方式引导孩子；Reflection Agent和Memory Agent负责记录孩子的学习状态和兴趣点，为后续对话提供参考。

**Why this priority**: 交互优化提升用户体验，让孩子感受到学习的乐趣而不是压力。反思和记忆功能帮助系统了解孩子的学习状态，提供个性化的学习体验。

**Independent Test**: 可以独立测试：Interaction Agent在回答结尾提供可选动作（如"你想不想试试？"、"我们下一步看什么？"），不制造学习压力；Reflection Agent判断孩子的兴趣和困惑，Memory Agent记录孩子的学习状态。这个功能可以独立工作，提升用户体验。

**Acceptance Scenarios**:

1. **Given** Interaction Agent被调用, **When** 生成回答结尾, **Then** 回答结尾包含可选动作（如"你想不想试试？"、"我们下一步看什么？"），不制造学习压力
2. **Given** Reflection Agent被调用, **When** 分析孩子反应, **Then** Reflection Agent判断孩子是否表现出兴趣、出现困惑或需要放松
3. **Given** Memory Agent被调用, **When** 记录学习状态, **Then** Memory Agent记录孩子感兴趣的主题、已理解/未理解的点
4. **Given** 后续对话, **When** 系统生成回答, **Then** 系统参考Memory Agent记录的信息，提供个性化的回答
5. **Given** 孩子表现出困惑, **When** Reflection Agent检测到, **Then** Reflection Agent输出给Memory Agent，Memory Agent记录未理解的点，后续对话中避免类似问题

---

### User Story 6 - 新接口创建和前端配置 (Priority: P1)

创建新的 `/api/conversation/agent` 接口作为多Agent模式接口，保留 `/api/conversation/stream` 接口作为单Agent模式接口。前端通过配置选择使用哪个接口，两个接口的输入和输出保持一致。

**Why this priority**: 接口分离是功能实现的基础，通过创建新接口而不是重构旧接口，可以保证向后兼容，降低风险。前端配置切换功能确保可以灵活地在两个接口之间切换，方便测试和逐步迁移。这是整个功能实现的前提条件。

**Independent Test**: 可以独立测试：前端根据配置选择调用 `/api/conversation/agent`（多Agent模式）或 `/api/conversation/stream`（单Agent模式），两个接口接收相同的请求参数，返回相同格式的流式响应。这个功能可以独立工作，不影响现有功能。

**Acceptance Scenarios**:

1. **Given** 前端配置使用多Agent模式, **When** 用户进行追问对话, **Then** 前端调用 `/api/conversation/agent` 接口，接口使用多Agent协作流程（Supervisor → Intent Agent → Cognitive Load Agent → Learning Planner Agent → Domain Agent → Interaction Agent → Reflection Agent → Memory Agent）
2. **Given** 前端配置使用单Agent模式, **When** 用户进行追问对话, **Then** 前端调用 `/api/conversation/stream` 接口，接口使用单一Agent模型模式（直接调用Conversation Node）
3. **Given** 前端配置未设置或使用默认值, **When** 用户进行追问对话, **Then** 前端调用 `/api/conversation/stream` 接口（向后兼容）
4. **Given** 两个接口接收相同的请求参数, **When** 用户发送追问消息, **Then** 两个接口都能正确解析请求参数，包括messageType、message、image、voice、sessionId等字段
5. **Given** 两个接口生成回答, **When** 流式返回响应, **Then** 两个接口返回相同格式的SSE事件流，包括message事件、done事件、error事件等
6. **Given** 前端切换配置, **When** 从单Agent模式切换到多Agent模式, **Then** 前端能够无缝切换接口调用，用户体验一致

---

### Edge Cases

- **Supervisor决策失败**：当Supervisor无法确定合适的Agent时，系统应降级到默认的Conversation Agent，确保对话能够继续
- **Agent调用超时**：当某个Agent调用超时（超过10秒）时，系统应跳过该Agent，选择备用Agent或降级处理
- **多个Agent冲突**：当多个Agent返回冲突的回答时，Supervisor应优先选择置信度更高的Agent，或综合多个Agent的回答
- **认知负载判断错误**：当认知负载判断错误导致回答过难或过易时，系统应允许用户反馈，调整后续回答的难度
- **记忆数据丢失**：当Memory Agent记录的数据丢失时，系统应能够从对话历史中恢复，或重新开始记录
- **工具调用失败**：当领域Agent调用的工具失败时，Agent应能够降级处理，不依赖工具也能生成基本回答
- **意图识别不准确**：当意图识别不准确（置信度<50%）时，系统应使用多个意图的加权结果，或询问用户澄清意图
- **Graph执行中断**：当Graph执行过程中中断时，系统应能够保存中间状态，支持恢复执行
- **新接口调用失败**：当前端调用 `/api/conversation/agent` 接口失败时，前端可以降级到 `/api/conversation/stream` 接口，确保对话能够继续
- **接口切换错误**：当前端切换接口调用时，如果当前请求正在处理，前端应等待当前请求完成后再切换
- **输入输出不一致**：当两个接口的输入输出格式不一致时，系统应记录错误日志，前端应显示友好的错误提示

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: 系统必须实现Follow-up Supervisor作为核心中枢，负责协调多个子Agent协作，不直接生成教学内容
- **FR-002**: Supervisor必须根据孩子的学习意图、认知负载、识别对象和年龄段，选择最合适的子Agent参与回答
- **FR-003**: Supervisor必须遵守单轮对话只解决一个认知目标的规则，以"孩子是否听得懂"为第一优先级
- **FR-004**: 系统必须实现Intent Agent，负责判断孩子的追问意图类型（认知型、探因型、表达型、游戏型、情绪型），只输出意图标签和置信度
- **FR-005**: 系统必须实现Cognitive Load Agent，负责根据孩子年龄、对话轮次、最近输出长度，判断最合适的输出策略（简短讲解、类比讲解、反问引导、暂停探索）
- **FR-006**: 系统必须实现Learning Planner Agent，负责根据意图判断、认知负载建议、识别对象、孩子年龄段，决定是否继续深入、选择哪个领域Agent、是"讲一点"还是"问一个问题"
- **FR-007**: 系统必须实现Science Agent，负责用孩子能理解的方式解释自然现象，可以使用工具（simple_fact_lookup、get_current_time、image_generate_simple），只回答一个知识点，用生活类比，不超过4句话
- **FR-008**: 系统必须实现Language Agent，负责让孩子"说得出口"，不讲语法规则，可以使用工具（simple_dictionary、pronunciation_hint），输出包含可模仿的句子，使用孩子日常语言
- **FR-009**: 系统必须实现Humanities Agent，负责把自然与文化连接起来，提供古诗词、故事、画面感，不要求背诵，必须和眼前看到的事物有关
- **FR-010**: 系统必须实现Interaction Agent，负责把内容说"轻"，给孩子可选动作，不制造学习压力，常用结尾方式（"你想不想试试？"、"我们下一步看什么？"、"要不要换个角度？"）
- **FR-011**: 系统必须实现Reflection Agent，负责判断孩子是否表现出兴趣、出现困惑、需要放松，输出给Memory Agent
- **FR-012**: 系统必须实现Memory Agent，负责记录孩子感兴趣的主题、已理解/未理解的点，为后续对话提供参考
- **FR-013**: 系统必须使用eino Graph结构组织多Agent协作流程，Supervisor控制分支，Domain Agent内部可调用Tool
- **FR-014**: 系统必须支持Graph执行流程：Start → Follow-up Supervisor → Intent Agent → Cognitive Load Agent → Learning Planner Agent → Domain Agent（Science/Language/Humanities）→ Interaction Agent → Reflection Agent → Memory Agent → Back to Supervisor
- **FR-015**: 系统必须支持Domain Agent调用工具增强回答准确性，但工具调用失败时能够降级处理
- **FR-016**: 系统必须支持Supervisor根据多个Agent的返回结果，综合决策最终回答内容
- **FR-017**: 系统必须支持Memory Agent记录的信息在后续对话中被参考，提供个性化回答
- **FR-018**: 系统必须支持单轮回答不超过5句话的规则，避免信息过载
- **FR-019**: 系统必须支持如果内容可能变复杂，优先拆成"一步一步来"，不一次性讲完完整知识体系
- **FR-020**: 系统必须支持永远以"孩子是否听得懂"为第一优先级，避免考试、作业、说教语气
- **FR-021**: 系统必须创建新的 `/api/conversation/agent` 接口，实现多Agent协作流程
- **FR-022**: 系统必须保留 `/api/conversation/stream` 接口，保持单Agent模式不变，确保向后兼容
- **FR-023**: 两个接口必须接收相同的请求参数格式（UnifiedStreamConversationRequest），包括messageType、message、image、voice、sessionId、identificationContext、userAge、maxContextRounds等字段
- **FR-024**: 两个接口必须返回相同格式的SSE流式响应，包括相同的事件类型（message、done、error）和数据结构
- **FR-025**: 前端必须通过配置选择使用哪个接口，配置方式可以是环境变量、配置文件或其他前端配置机制
- **FR-026**: 系统必须支持前端配置切换，允许前端在运行时切换接口调用，无需重启
- **FR-027**: 系统必须保证两个接口的输入输出一致性，确保前端可以无缝切换接口调用

### Key Entities *(include if feature involves data)*

- **Follow-up Supervisor状态（SupervisorState）**：包含当前识别对象、已生成的三张卡片、孩子年龄/年级、最近对话轮数与时长、子Agent的返回结果等属性
- **意图识别结果（IntentResult）**：包含意图类型（认知型、探因型、表达型、游戏型、情绪型）、置信度等属性
- **认知负载建议（CognitiveLoadAdvice）**：包含输出策略（简短讲解、类比讲解、反问引导、暂停探索）、建议理由等属性
- **学习计划决策（LearningPlanDecision）**：包含是否继续深入、选择的领域Agent、教学动作（讲一点/问一个问题）等属性
- **领域Agent回答（DomainAgentResponse）**：包含领域类型（Science/Language/Humanities）、回答内容、使用的工具、工具调用结果等属性
- **交互优化结果（InteractionOptimization）**：包含结尾方式、可选动作等属性
- **反思结果（ReflectionResult）**：包含孩子兴趣状态、困惑状态、放松需求等属性
- **记忆记录（MemoryRecord）**：包含感兴趣的主题、已理解/未理解的点、学习历史等属性
- **Graph执行状态（GraphExecutionState）**：包含当前执行的Agent节点、执行路径、中间结果、错误状态等属性

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 意图识别准确率≥80%（孩子提出的问题能够被正确分类为5种意图类型之一）
- **SC-002**: 认知负载判断准确率≥85%（系统选择的输出策略适合孩子的年龄和认知水平）
- **SC-003**: Supervisor决策准确率≥80%（Supervisor选择的Agent和策略能够生成合适的回答）
- **SC-004**: 单轮回答平均长度≤5句话（符合认知负载控制要求）
- **SC-005**: 领域Agent回答专业性≥90%（Science Agent回答科学问题准确，Language Agent提供可模仿的句子，Humanities Agent提供相关文化内容）
- **SC-006**: 工具调用成功率≥85%（Domain Agent调用的工具能够成功返回结果）
- **SC-007**: Graph执行成功率≥90%（多Agent协作流程能够完整执行，不中断）
- **SC-008**: 回答适龄性≥90%（3-6岁孩子收到简短易懂的回答，7-12岁收到中等深度回答，13-18岁收到深入回答）
- **SC-009**: 交互优化效果≥85%（Interaction Agent的结尾方式能够引导孩子继续探索，不制造学习压力）
- **SC-010**: 记忆功能有效性≥80%（Memory Agent记录的信息能够在后续对话中被正确参考）
- **SC-011**: 回答相关性≥85%（回答内容与孩子的追问问题相关，与识别对象相关）
- **SC-012**: 系统响应时间≤8秒（从孩子发送追问到开始接收回答，包括多Agent协作时间）
- **SC-013**: 用户满意度≥80%（孩子和家长对追问回答的质量和适龄性满意）
- **SC-014**: 接口创建成功率100%（新接口 `/api/conversation/agent` 能够正常创建并响应请求）
- **SC-015**: 向后兼容性100%（旧接口 `/api/conversation/stream` 行为与之前完全一致，不受影响）
- **SC-016**: 输入输出一致性100%（两个接口接收相同的请求参数，返回相同格式的响应）
- **SC-017**: 前端配置切换成功率100%（前端能够正确读取配置并选择对应的接口调用）
- **SC-018**: 接口切换无缝性100%（前端切换接口调用时，用户体验一致，无感知）

## Assumptions

- 系统已经实现基础的对话功能，支持文本输入和流式输出
- 系统已经配置好eino框架，支持Graph结构和Agent节点
- 系统已经实现意图识别功能，能够识别基本的意图类型
- 系统已经实现对话上下文管理，能够获取孩子的年龄、识别对象、对话历史等信息
- 系统已经实现工具调用机制，支持Domain Agent调用外部工具
- 系统已经实现记忆存储机制，支持Memory Agent记录和检索学习状态
- 孩子已经完成拍照识别，识别结果包含对象名称、类别、关键词等信息
- 用户年级信息已经获取，用于内容适配和认知负载判断
- `/api/conversation/stream` 接口已经存在，支持流式输出和统一的多模态输入处理
- 系统配置机制已经支持从环境变量读取配置（如 `USE_AI_MODEL`）
- 前端已经实现配置管理机制，能够读取配置并选择接口调用
- 前端已经实现SSE流式接收机制，能够处理流式响应

## Dependencies

- eino框架必须支持Graph结构和Agent节点，支持Agent之间的协作和状态传递
- 意图识别功能必须可用，能够准确识别5种意图类型（认知型、探因型、表达型、游戏型、情绪型）
- 对话上下文管理机制必须可用，能够提供孩子的年龄、识别对象、对话历史等信息
- 工具调用机制必须可用，支持Domain Agent调用外部工具（simple_fact_lookup、simple_dictionary、pronunciation_hint等）
- 记忆存储机制必须可用，支持Memory Agent记录和检索学习状态
- 流式输出机制必须可用，支持多Agent协作结果的流式返回
- 错误处理和降级机制必须可用，支持Agent调用失败时的降级处理
- `/api/conversation/stream` 接口必须可用，保持单Agent模式不变
- `/api/conversation/agent` 接口必须实现，支持多Agent协作流程
- 前端配置机制必须可用，支持读取配置并选择接口调用
- 两个接口的请求参数格式必须一致（UnifiedStreamConversationRequest）
- 两个接口的响应格式必须一致（SSE事件流格式）

## Configuration

### 接口配置

系统通过接口分离的方式实现多Agent功能：

- **后端接口**：
  - `/api/conversation/stream`：单Agent模式接口（保持不变，向后兼容）
  - `/api/conversation/agent`：多Agent模式接口（新创建）

- **前端配置**：
  - 前端通过配置选择使用哪个接口进行追问的流式回答
  - 配置方式可以是环境变量、配置文件或其他前端配置机制
  - 默认使用 `/api/conversation/stream` 接口（向后兼容）

**接口一致性要求**:
- 两个接口必须接收相同的请求参数格式（UnifiedStreamConversationRequest）
- 两个接口必须返回相同格式的SSE流式响应
- 前端可以无缝切换接口调用，用户体验一致

**实现要求**:
- 后端创建新的 `/api/conversation/agent` 接口，实现多Agent协作流程
- 后端保留 `/api/conversation/stream` 接口，保持单Agent模式不变
- 前端实现配置机制，支持选择接口调用
- 前端实现接口切换逻辑，支持运行时切换

### Agent配置

系统支持通过配置文件控制各个Agent的行为：

- **Supervisor配置**：控制Supervisor的决策策略、Agent选择规则、信息量控制规则
- **Intent Agent配置**：控制意图识别的模型、置信度阈值、意图类型定义
- **Cognitive Load Agent配置**：控制认知负载判断的规则、年龄段对应的策略、输出长度限制
- **Learning Planner Agent配置**：控制教学决策的规则、领域Agent选择规则、教学动作定义
- **Domain Agent配置**：控制各个领域Agent的专业性要求、工具调用规则、回答长度限制
- **Interaction Agent配置**：控制交互优化的策略、结尾方式定义、可选动作定义
- **Reflection Agent配置**：控制反思判断的规则、兴趣/困惑/放松的识别标准
- **Memory Agent配置**：控制记忆记录的规则、存储方式、检索方式

### Graph配置

系统支持通过配置文件控制Graph的执行流程：

- **Graph节点定义**：定义各个Agent节点及其执行顺序
- **Graph分支规则**：定义Supervisor如何选择不同的Agent分支
- **Graph错误处理**：定义Agent调用失败时的降级策略
- **Graph超时控制**：定义各个Agent节点的超时时间

