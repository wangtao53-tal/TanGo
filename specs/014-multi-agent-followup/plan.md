# Implementation Plan: Multi-Agent系统智能增强

**Feature**: Multi-Agent系统工具调用和MCP集成增强  
**Created**: 2025-12-24  
**Status**: Planning  
**Branch**: `dev-mvp-20251218`

## Technical Context

### Current State
- ✅ Multi-Agent系统基础架构已实现
- ✅ Supervisor协调机制已实现
- ✅ Domain Agent（Science/Language/Humanities）已实现
- ✅ Interaction/Reflection/Memory Agent已实现
- ✅ **工具调用基础架构已实现**：Tool接口、ToolRegistry、5个基础工具
- ✅ **工具调用处理逻辑已实现**：Science Agent和Language Agent已支持工具调用处理
- ⏳ **工具注册到ChatModel待验证**：需要在实际环境中测试工具调用流程
- ❌ **MCP集成未实现**：未集成MCP资源
- ❌ **智能工具选择未实现**：Supervisor无法动态决定是否需要工具

### Technology Stack
- **Backend**: Go + go-zero框架
- **AI Framework**: Eino框架（CloudWeGo）
- **Tool Calling**: Eino ChatModel工具调用能力
- **MCP**: Model Context Protocol（MCP）资源集成
- **Storage**: 内存存储（MemoryStorage）

### Dependencies
- `github.com/cloudwego/eino`: Eino框架核心
- `github.com/cloudwego/eino-ext/components/model/ark`: Ark模型实现
- `github.com/cloudwego/eino/schema`: Eino Schema定义
- MCP客户端库（需要确认）

### Integration Points
- **Eino Tool Calling**: 需要在ChatModel中注册工具，处理工具调用请求
- **MCP Integration**: 需要包装MCP资源为工具，集成到工具注册表
- **Supervisor Enhancement**: 需要增加工具选择逻辑

### Constraints
- 必须保持向后兼容，不影响现有功能
- 工具调用失败时必须降级处理
- 工具调用响应时间必须≤2秒
- 必须支持Mock模式（工具不可用时）

## Constitution Check

### Code Quality
- ✅ 遵循Go代码规范
- ✅ 使用中文注释和文档
- ✅ 实现错误处理和降级机制

### Architecture
- ✅ 遵循go-zero架构模式
- ✅ 使用Eino框架最佳实践
- ✅ 保持模块化和可扩展性

### Testing
- ⏳ 需要添加工具调用单元测试
- ⏳ 需要添加MCP集成测试
- ⏳ 需要添加集成测试

## Phase 0: Research & Design

### Research Tasks

#### 1. Eino Tool Calling机制研究 ✅
**问题**: 如何在Eino框架中实现工具调用？

**研究发现**:
- ✅ Eino ChatModel通过`Generate`方法返回`ToolCalls`字段
- ✅ 工具调用处理流程已实现：检查ToolCalls → 执行工具 → 创建ToolMessage → 重新调用ChatModel
- ✅ 工具结果整合逻辑已实现
- ✅ 工具调用错误处理和降级机制已实现
- ⏳ 工具注册到ChatModel的方式需要在实际环境中验证

**输出**: `research.md` 中已增加工具调用研究章节（章节13）

#### 2. MCP集成方式研究 ✅
**问题**: 如何集成MCP资源到Agent系统？

**研究发现**:
- ✅ MCP协议通过SSE或HTTP接口访问
- ✅ MCP工具包装模式：包装资源为Tool接口实现
- ✅ MCP资源发现机制设计
- ✅ MCP错误处理策略

**输出**: `research.md` 中已增加MCP集成研究章节（章节14）

#### 3. 智能工具选择策略研究 ✅
**问题**: Supervisor如何智能选择工具？

**研究发现**:
- ✅ 工具选择决策树：根据意图类型和问题关键词
- ✅ 工具使用场景分析：Science Agent和Language Agent的工具映射
- ✅ 工具调用链设计：支持多轮调用，限制深度
- ✅ 工具结果整合策略：使用ChatModel自然整合

**输出**: `research.md` 中已增加智能工具选择研究章节（章节15）

### Design Artifacts

#### 1. 工具接口设计
**文件**: `data-model.md`
**内容**:
- Tool接口定义
- ToolRegistry设计
- 工具调用流程设计

#### 2. MCP集成设计
**文件**: `data-model.md`
**内容**:
- MCP工具包装器设计
- MCP资源映射设计
- MCP调用流程设计

#### 3. Supervisor增强设计
**文件**: `data-model.md`
**内容**:
- 工具选择逻辑设计
- 工具分配策略设计
- 工具结果整合设计

## Phase 1: Core Implementation

### Task 1: 工具调用基础架构
**优先级**: P0  
**估计时间**: 3-5天

**子任务**:
1. 定义Tool接口和ToolRegistry
2. 实现基础工具（simple_fact_lookup、simple_dictionary等）
3. 集成Eino工具调用能力
4. 在Domain Agent中实现工具调用处理
5. 实现工具调用错误处理和降级机制

**验收标准**:
- ✅ Tool接口定义清晰
- ✅ ToolRegistry可以注册和获取工具
- ✅ Domain Agent可以调用工具
- ✅ 工具调用失败时降级处理

### Task 2: MCP资源集成
**优先级**: P1  
**估计时间**: 2-3天

**子任务**:
1. 发现可用MCP资源
2. 实现MCP工具包装器
3. 注册MCP工具到工具注册表
4. 测试MCP工具调用

**验收标准**:
- ✅ MCP资源可以包装为工具
- ✅ MCP工具可以注册到工具注册表
- ✅ Domain Agent可以调用MCP工具
- ✅ MCP工具调用失败时降级处理

### Task 3: Supervisor智能协调增强
**优先级**: P1  
**估计时间**: 2-3天

**子任务**:
1. 实现动态工具分配逻辑
2. 定义工具使用策略
3. 优化Supervisor协调逻辑
4. 实现工具结果整合

**验收标准**:
- ✅ Supervisor可以根据问题类型选择工具
- ✅ 工具使用策略清晰明确
- ✅ 工具结果可以整合到回答中
- ✅ 工具选择准确率≥85%

## Phase 2: Testing & Validation

### Unit Tests
- Tool接口和ToolRegistry单元测试
- Domain Agent工具调用单元测试
- MCP工具包装器单元测试
- Supervisor工具选择单元测试

### Integration Tests
- 工具调用流程集成测试
- MCP集成测试
- Supervisor协调增强集成测试
- 端到端测试

### Performance Tests
- 工具调用响应时间测试（目标≤2秒）
- 工具调用成功率测试（目标≥90%）
- 工具结果整合质量测试（自然度≥80%）

## Phase 3: Documentation & Deployment

### Documentation
- 工具调用使用文档
- MCP集成文档
- Supervisor增强文档
- API文档更新

### Deployment
- 环境变量配置更新
- 部署脚本更新
- 监控和日志配置

## Risks & Mitigation

### Risk 1: Eino工具调用API不熟悉
**影响**: 高  
**概率**: 中  
**缓解措施**: 
- 深入研究Eino文档
- 创建POC验证工具调用流程
- 咨询Eino社区

### Risk 2: MCP集成复杂度高
**影响**: 中  
**概率**: 中  
**缓解措施**:
- 先实现简单的MCP工具包装器
- 逐步增加MCP资源支持
- 实现完善的错误处理

### Risk 3: 工具调用性能问题
**影响**: 中  
**概率**: 低  
**缓解措施**:
- 实现工具调用超时机制
- 实现工具调用缓存
- 优化工具调用流程

## Success Criteria

### Functional
- ✅ Domain Agent可以调用工具获取准确信息
- ✅ MCP资源可以集成到Agent系统
- ✅ Supervisor可以智能选择工具
- ✅ 工具结果可以整合到回答中

### Non-Functional
- ✅ 工具调用成功率≥90%
- ✅ 工具调用响应时间≤2秒
- ✅ 工具结果整合自然度≥80%
- ✅ 智能工具选择准确率≥85%

## Next Steps

1. **立即开始**: Phase 0研究阶段
2. **研究重点**: Eino工具调用机制和MCP集成方式
3. **设计重点**: 工具接口设计和Supervisor增强设计
4. **实现重点**: 工具调用基础架构和MCP资源集成
