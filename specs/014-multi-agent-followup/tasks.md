# Tasks: Multi-Agent系统工具调用和MCP集成增强

**Feature**: Multi-Agent系统工具调用和MCP集成增强  
**Created**: 2025-12-24  
**Status**: Planning  
**Branch**: `dev-mvp-20251218`

## Implementation Strategy

### MVP Scope
- **Phase 1-2**: 工具调用基础架构和MCP配置
- **Phase 3**: 基础工具实现（simple_fact_lookup、simple_dictionary）
- **Phase 4**: MCP工具集成（tal_time作为示例）
- **Phase 5**: Supervisor智能工具选择

### Incremental Delivery
1. **Week 1**: 工具调用基础架构 + 基础工具
2. **Week 2**: MCP配置和集成
3. **Week 3**: Supervisor智能协调增强
4. **Week 4**: 测试和优化

## Dependencies

### Story Completion Order
1. **US-Tool-1**: 工具调用基础架构（P0，阻塞所有工具相关功能）
2. **US-Tool-2**: MCP配置和集成（P1，依赖US-Tool-1）
3. **US-Tool-3**: Supervisor智能工具选择（P1，依赖US-Tool-1和US-Tool-2）
4. **US-Tool-4**: Domain Agent工具调用集成（P2，依赖US-Tool-1）

### Parallel Opportunities
- 基础工具实现可以并行（simple_fact_lookup、simple_dictionary等）
- MCP工具包装器可以并行实现
- Domain Agent工具调用集成可以并行（Science、Language、Humanities）

## Phase 1: Setup & Configuration

### T001: 创建工具调用基础目录结构
- [X] T001 创建 `backend/internal/tools/` 目录结构
  - `backend/internal/tools/tool.go` - Tool接口定义
  - `backend/internal/tools/registry.go` - ToolRegistry实现
  - `backend/internal/tools/base/` - 基础工具实现目录
  - `backend/internal/tools/mcp/` - MCP工具包装器目录

### T002: 扩展配置结构支持MCP配置
- [ ] T002 在 `backend/internal/config/config.go` 中添加MCP配置结构
  - 添加 `MCPConfig` 类型，包含MCP服务器配置
  - 支持从环境变量或配置文件读取MCP服务器URL
  - 添加 `MCPEnabled` 标志控制MCP功能开关

### T003: 创建MCP配置文件解析
- [ ] T003 创建 `backend/internal/config/mcp.go` 实现MCP配置解析
  - 支持读取MCP服务器配置（URL、认证信息等）
  - 支持从环境变量 `MCP_SERVERS` 读取JSON配置
  - 支持从配置文件读取MCP配置

## Phase 2: Tool Interface & Registry Foundation

### T004: 定义Tool接口
- [X] T004 [P] 在 `backend/internal/tools/tool.go` 中定义Tool接口
  - `Name() string` - 工具名称
  - `Description() string` - 工具描述
  - `Execute(ctx context.Context, params map[string]interface{}) (interface{}, error)` - 执行工具
  - `Parameters() map[string]interface{}` - 工具参数定义（用于Eino工具注册）

### T005: 实现ToolRegistry
- [X] T005 [P] 在 `backend/internal/tools/registry.go` 中实现ToolRegistry
  - `Register(tool Tool)` - 注册工具
  - `GetTool(name string) (Tool, bool)` - 获取工具
  - `ListTools() []Tool` - 列出所有工具
  - `GetToolsForAgent(agentType string) []Tool` - 根据Agent类型获取工具列表

### T006: 创建全局工具注册表实例
- [X] T006 在 `backend/internal/tools/registry.go` 中创建全局工具注册表
  - 创建 `DefaultRegistry` 全局实例
  - 实现工具注册初始化函数
  - 支持从配置加载工具

## Phase 3: Base Tools Implementation [US-Tool-1]

### Story Goal
实现基础工具调用机制，让Domain Agent能够调用简单工具获取信息。

### Independent Test Criteria
- 可以独立测试：Domain Agent调用simple_fact_lookup工具，能够获取事实信息并整合到回答中
- 工具调用失败时能够降级处理，不影响Agent正常回答
- 工具调用响应时间≤2秒

### T007: 实现simple_fact_lookup工具
- [X] T007 [P] [US-Tool-1] 创建 `backend/internal/tools/base/fact_lookup.go` 实现simple_fact_lookup工具
  - 实现Tool接口
  - 接收查询关键词参数
  - 返回简单事实信息（Mock实现，后续可接入真实API）
  - 实现错误处理和超时机制

### T008: 实现simple_dictionary工具
- [X] T008 [P] [US-Tool-1] 创建 `backend/internal/tools/base/dictionary.go` 实现simple_dictionary工具
  - 实现Tool接口
  - 接收单词参数
  - 返回单词释义和例句（Mock实现，后续可接入真实API）
  - 实现错误处理和超时机制

### T009: 实现pronunciation_hint工具
- [X] T009 [P] [US-Tool-1] 创建 `backend/internal/tools/base/pronunciation.go` 实现pronunciation_hint工具
  - 实现Tool接口
  - 接收单词参数
  - 返回发音提示（Mock实现，后续可接入真实API）
  - 实现错误处理和超时机制

### T010: 实现get_current_time工具
- [X] T010 [P] [US-Tool-1] 创建 `backend/internal/tools/base/time.go` 实现get_current_time工具
  - 实现Tool接口
  - 不需要参数
  - 返回当前时间信息
  - 实现错误处理

### T011: 实现image_generate_simple工具
- [X] T011 [P] [US-Tool-1] 创建 `backend/internal/tools/base/image_gen.go` 实现image_generate_simple工具
  - 实现Tool接口
  - 接收描述参数
  - 返回图片URL（Mock实现，后续可接入真实API）
  - 实现错误处理和超时机制

### T012: 注册基础工具到注册表
- [X] T012 [US-Tool-1] 在 `backend/internal/tools/init.go` 中注册所有基础工具
  - 在初始化函数中注册simple_fact_lookup
  - 在初始化函数中注册simple_dictionary
  - 在初始化函数中注册pronunciation_hint
  - 在初始化函数中注册get_current_time
  - 在初始化函数中注册image_generate_simple

## Phase 4: Eino Tool Calling Integration [US-Tool-1]

### T013: 研究Eino工具调用API
- [ ] T013 [US-Tool-1] 研究Eino框架工具调用机制
  - 查看Eino文档了解工具注册方式
  - 查看Eino文档了解工具调用请求处理
  - 创建POC验证工具调用流程
  - 记录研究结果到 `specs/014-multi-agent-followup/research.md`

### T014: 实现Eino工具定义转换
- [X] T014 [US-Tool-1] 创建 `backend/internal/tools/eino_adapter.go` 实现工具定义转换
  - `ToolToEinoSchema(tool Tool) schema.Tool` - 将Tool转换为Eino Schema
  - 支持Function类型工具定义
  - 支持参数定义转换
  - 支持描述信息转换

### T015: 在Science Agent中集成工具调用
- [X] T015 [US-Tool-1] 修改 `backend/internal/agent/nodes/science_agent_node.go` 支持工具调用
  - 在initChatModel中注册工具（simple_fact_lookup、get_current_time、image_generate_simple）
  - 在executeReal中处理工具调用请求
  - 执行工具调用并获取结果
  - 将工具结果添加到消息中，重新调用ChatModel
  - 整合工具结果到最终回答

### T016: 在Language Agent中集成工具调用
- [X] T016 [US-Tool-1] 修改 `backend/internal/agent/nodes/language_agent_node.go` 支持工具调用
  - 在initChatModel中注册工具（simple_dictionary、pronunciation_hint）
  - 在executeReal中处理工具调用请求
  - 执行工具调用并获取结果
  - 将工具结果添加到消息中，重新调用ChatModel
  - 整合工具结果到最终回答

### T017: 实现工具调用错误处理和降级
- [X] T017 [US-Tool-1] 在所有Domain Agent中实现工具调用错误处理
  - 工具调用失败时记录错误日志（已实现）
  - 工具调用失败时降级到不使用工具的回答（已实现）
  - 工具调用超时时自动取消并降级（已实现，通过context timeout）
  - 确保工具调用失败不影响Agent正常回答（已实现）

## Phase 5: MCP Configuration & Integration [US-Tool-2]

### Story Goal
集成MCP资源到Agent系统，支持配置MCP服务器和工具。

### Independent Test Criteria
- 可以独立测试：系统能够读取MCP配置并连接到MCP服务器
- MCP工具可以包装为Agent可用的工具
- MCP工具调用失败时能够降级处理

### T018: 实现MCP客户端基础库
- [X] T018 [US-Tool-2] 创建 `backend/internal/tools/mcp/client.go` 实现MCP客户端
  - 实现MCP协议基础通信（HTTP连接）
  - 支持连接到MCP服务器
  - 支持调用MCP资源
  - 实现错误处理和超时机制

### T019: 实现MCP工具包装器
- [X] T019 [US-Tool-2] 创建 `backend/internal/tools/mcp/wrapper.go` 实现MCP工具包装器
  - 实现Tool接口
  - 包装MCP资源为工具
  - 处理MCP资源调用
  - 格式化MCP结果

### T020: 实现tal_time MCP工具包装
- [X] T020 [P] [US-Tool-2] 创建 `backend/internal/tools/mcp/tal_time.go` 实现tal_time工具包装
  - 包装tal_time MCP资源
  - 实现Tool接口
  - 调用MCP资源获取时间信息
  - 格式化时间结果为Agent可用格式（含降级处理）

### T021: 实现MCP资源发现机制
- [X] T021 [US-Tool-2] 创建 `backend/internal/tools/mcp/discovery.go` 实现MCP资源发现
  - 连接到MCP服务器
  - 发现可用资源
  - 自动包装资源为工具
  - 注册到工具注册表

### T022: 在配置中支持MCP服务器配置
- [X] T022 [US-Tool-2] 创建 `backend/internal/config/mcp.go` 添加MCP配置
  - 添加 `MCPConfig` 结构体
  - 支持从环境变量读取MCP服务器配置
  - 支持从JSON配置文件读取MCP配置
  - 添加MCP功能开关

### T023: 实现MCP配置加载和初始化
- [X] T023 [US-Tool-2] 在 `backend/internal/agent/multiagent_graph.go` 中实现MCP初始化
  - 从配置加载MCP服务器列表
  - 连接到MCP服务器
  - 发现并包装MCP资源
  - 注册MCP工具到工具注册表

### T024: 注册MCP工具到注册表
- [X] T024 [US-Tool-2] 在 `backend/internal/tools/mcp/discovery.go` 中支持MCP工具注册
  - 实现MCP工具注册方法
  - 支持按MCP服务器分组工具
  - 支持工具启用/禁用

## Phase 6: Supervisor Intelligent Tool Selection [US-Tool-3]

### Story Goal
让Supervisor能够根据问题类型智能选择是否需要工具，以及选择哪些工具。

### Independent Test Criteria
- 可以独立测试：Supervisor根据意图和问题内容选择工具
- 工具选择准确率≥85%
- 工具选择失败时能够降级处理

### T025: 实现工具选择决策逻辑
- [X] T025 [US-Tool-3] 创建 `backend/internal/agent/nodes/tool_strategy.go` 实现工具选择逻辑
  - `SelectTools(intent, message, domainAgent, confidence)` - 根据意图和消息选择工具
  - 实现意图到工具的映射规则
  - 实现关键词匹配逻辑
  - 返回推荐的工具列表和策略

### T026: 定义工具使用策略
- [X] T026 [US-Tool-3] 在 `backend/internal/agent/nodes/tool_strategy.go` 中定义工具使用策略
  - 高置信度问题策略：直接使用工具
  - 探索性问题策略：先回答，再提供工具增强
  - 简单问题策略：不使用工具
  - 复杂问题策略：使用多个工具

### T027: 在Supervisor协调中集成工具选择
- [X] T027 [US-Tool-3] 修改 `backend/internal/agent/nodes/supervisor_node.go` 的Coordinate方法
  - 在协调流程中添加工具选择步骤
  - 将选择的工具添加到LearningPlanDecision
  - 将工具信息传递给Domain Agent

### T028: 在LearningPlanDecision中添加工具信息
- [X] T028 [US-Tool-3] 修改 `backend/internal/types/multiagent_types.go` 的LearningPlanDecision结构
  - 添加 `Tools []string` 字段
  - 添加 `ToolStrategy string` 字段
  - 更新相关序列化逻辑

### T029: 在Domain Agent中应用工具选择
- [X] T029 [US-Tool-3] 修改Domain Agent节点支持工具选择
  - Science Agent接收recommendedTools参数（已实现）
  - Language Agent接收recommendedTools参数（已实现）
  - Supervisor工具选择结果传递给Domain Agent（已实现）

## Phase 7: Domain Agent Tool Integration [US-Tool-4]

### Story Goal
让Domain Agent能够实际调用工具，并将工具结果整合到回答中。

### Independent Test Criteria
- 可以独立测试：Domain Agent调用工具获取信息并整合到回答中
- 工具调用成功率≥90%
- 工具结果整合自然度≥80%

### T030: 实现工具调用链处理
- [X] T030 [US-Tool-4] 创建 `backend/internal/agent/nodes/tool_chain.go` 实现工具调用链
  - 支持多轮工具调用（已实现，最大深度3层）
  - 支持工具调用结果链式传递（已实现）
  - 实现调用深度限制（已实现，maxDepth=3）
  - 实现超时控制（已实现，timeout=10秒）

### T031: 实现工具结果整合逻辑
- [X] T031 [US-Tool-4] 工具结果整合逻辑已集成到工具调用链中
  - 工具调用链自动处理工具结果整合（通过ChatModel重新生成）
  - 保持回答简洁性（已实现）
  - 避免过度引用工具结果（已实现）

### T032: 在Science Agent中实现完整工具调用流程
- [X] T032 [US-Tool-4] 完善 `backend/internal/agent/nodes/science_agent_node.go` 的工具调用
  - 实现工具调用链处理（已集成ToolChain）
  - 实现工具结果整合（已集成到ToolChain）
  - 优化工具调用性能（已实现）
  - 添加工具调用日志（已实现）

### T033: 在Language Agent中实现完整工具调用流程
- [X] T033 [US-Tool-4] 完善 `backend/internal/agent/nodes/language_agent_node.go` 的工具调用
  - 实现工具调用链处理（已集成ToolChain）
  - 实现工具结果整合（已集成到ToolChain）
  - 优化工具调用性能（已实现）
  - 添加工具调用日志（已实现）

## Phase 8: Testing & Validation

### T034: 编写工具接口单元测试
- [ ] T034 [P] 创建 `backend/internal/tools/tool_test.go` 测试Tool接口
  - 测试Tool接口方法
  - 测试工具参数验证
  - 测试工具执行错误处理

### T035: 编写ToolRegistry单元测试
- [ ] T035 [P] 创建 `backend/internal/tools/registry_test.go` 测试ToolRegistry
  - 测试工具注册
  - 测试工具获取
  - 测试工具列表
  - 测试按Agent类型获取工具

### T036: 编写基础工具单元测试
- [ ] T036 [P] 创建 `backend/internal/tools/base/*_test.go` 测试基础工具
  - 测试simple_fact_lookup工具
  - 测试simple_dictionary工具
  - 测试pronunciation_hint工具
  - 测试get_current_time工具
  - 测试image_generate_simple工具

### T037: 编写MCP工具包装器单元测试
- [ ] T037 [P] 创建 `backend/internal/tools/mcp/*_test.go` 测试MCP工具
  - 测试MCP客户端连接
  - 测试MCP资源调用
  - 测试MCP工具包装器
  - 测试tal_time工具

### T038: 编写Domain Agent工具调用集成测试
- [ ] T038 创建 `backend/internal/agent/nodes/tool_integration_test.go` 测试工具调用集成
  - 测试Science Agent工具调用
  - 测试Language Agent工具调用
  - 测试工具调用错误处理
  - 测试工具结果整合

### T039: 编写Supervisor工具选择集成测试
- [ ] T039 创建 `backend/internal/agent/nodes/supervisor_tool_test.go` 测试工具选择
  - 测试工具选择逻辑
  - 测试工具使用策略
  - 测试工具选择准确率

### T040: 编写端到端测试
- [ ] T040 创建 `backend/internal/logic/agentlogic_tool_test.go` 测试端到端工具调用
  - 测试完整工具调用流程
  - 测试MCP工具集成
  - 测试工具调用性能
  - 测试工具调用降级机制

## Phase 9: Documentation & Configuration

### T041: 更新环境变量配置文档
- [ ] T041 更新 `.env.example` 添加MCP配置示例
  - 添加 `MCP_ENABLED` 环境变量
  - 添加 `MCP_SERVERS` 环境变量（JSON格式）
  - 添加MCP服务器配置示例

### T042: 创建工具调用使用文档
- [ ] T042 创建 `specs/014-multi-agent-followup/TOOL_USAGE.md` 文档
  - 工具调用使用指南
  - MCP工具配置指南
  - 工具开发指南
  - 故障排查指南

### T043: 创建MCP集成文档
- [ ] T043 创建 `specs/014-multi-agent-followup/MCP_INTEGRATION.md` 文档
  - MCP协议说明
  - MCP服务器配置
  - MCP工具包装指南
  - MCP资源发现机制

### T044: 更新API文档
- [ ] T044 更新API文档说明工具调用功能
  - 说明Domain Agent支持工具调用
  - 说明工具调用结果如何整合
  - 说明工具调用错误处理

## Phase 10: Polish & Cross-Cutting Concerns

### T045: 优化工具调用性能
- [ ] T045 优化工具调用性能
  - 实现工具调用缓存
  - 优化工具调用并发
  - 减少工具调用延迟
  - 确保工具调用响应时间≤2秒

### T046: 增强工具调用错误处理
- [ ] T046 增强工具调用错误处理
  - 完善错误日志
  - 实现错误重试机制
  - 实现错误降级策略
  - 添加错误监控

### T047: 添加工具调用监控和指标
- [ ] T047 添加工具调用监控和指标
  - 工具调用成功率指标
  - 工具调用响应时间指标
  - 工具调用错误率指标
  - 工具使用统计

### T048: 优化Supervisor工具选择算法
- [ ] T048 优化Supervisor工具选择算法
  - 提高工具选择准确率
  - 优化工具选择性能
  - 添加工具选择日志
  - 实现工具选择反馈机制

## Summary

### Task Statistics
- **Total Tasks**: 48
- **Setup Tasks**: 3 (T001-T003)
- **Foundation Tasks**: 3 (T004-T006)
- **Base Tools**: 6 (T007-T012)
- **Eino Integration**: 5 (T013-T017)
- **MCP Integration**: 7 (T018-T024)
- **Supervisor Enhancement**: 5 (T025-T029)
- **Domain Agent Integration**: 4 (T030-T033)
- **Testing**: 7 (T034-T040)
- **Documentation**: 4 (T041-T044)
- **Polish**: 4 (T045-T048)

### Parallel Execution Examples

#### Phase 3: Base Tools (可以并行)
- T007, T008, T009, T010, T011 可以并行实现（不同文件）

#### Phase 4: Eino Integration (部分并行)
- T014 和 T015 可以并行（T014完成后）
- T015 和 T016 可以并行（不同Agent）

#### Phase 5: MCP Integration (部分并行)
- T020 可以独立实现（tal_time工具）
- T021 和 T022 可以并行（不同功能）

### MVP Scope Recommendation
- **Week 1**: T001-T012 (Setup + Base Tools)
- **Week 2**: T013-T017 (Eino Integration)
- **Week 3**: T018-T024 (MCP Integration - tal_time)
- **Week 4**: T025-T029 (Supervisor Enhancement)

### Independent Test Criteria Summary
- **US-Tool-1**: Domain Agent可以调用基础工具，工具调用失败时降级处理
- **US-Tool-2**: 系统可以读取MCP配置并连接到MCP服务器，MCP工具可以包装为Agent可用工具
- **US-Tool-3**: Supervisor可以根据意图和问题选择工具，工具选择准确率≥85%
- **US-Tool-4**: Domain Agent可以调用工具并整合结果，工具调用成功率≥90%，结果整合自然度≥80%
