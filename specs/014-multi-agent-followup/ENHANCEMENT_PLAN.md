# Multi-Agent 系统智能增强计划

## 📋 当前状态分析

### 已实现的功能
1. ✅ **Supervisor协调机制**：Supervisor节点协调Intent、Cognitive Load、Learning Planner三个Agent
2. ✅ **领域Agent分工**：Science、Language、Humanities三个领域Agent
3. ✅ **交互优化**：Interaction Agent优化回答结尾
4. ✅ **反思记忆**：Reflection Agent和Memory Agent记录学习状态
5. ✅ **消息清理**：所有Agent节点已实现消息清理逻辑

### 缺失的功能
1. ❌ **工具调用（Tool Calling）**：Domain Agent虽然定义了工具，但未实际实现工具调用
2. ❌ **MCP集成**：未集成MCP（Model Context Protocol）资源
3. ❌ **动态工具选择**：Supervisor无法动态决定是否需要调用工具
4. ❌ **工具结果整合**：工具调用结果未整合到回答中

## 🎯 增强目标

### 核心目标
让multi-agent系统更加智能，通过工具调用和MCP集成增强回答的准确性和丰富性。

### 具体目标
1. **实现工具调用机制**：让Domain Agent能够调用外部工具获取准确信息
2. **集成MCP资源**：利用MCP提供的丰富资源（地图、天气、搜索等）
3. **智能工具选择**：Supervisor根据问题类型智能选择是否需要工具
4. **工具结果整合**：将工具调用结果自然整合到回答中

## 🔧 技术方案

### Phase 1: 工具调用基础架构

#### 1.1 工具注册机制
**目标**：建立统一的工具注册和管理机制

**实现方案**：
```go
// 定义工具接口
type Tool interface {
    Name() string
    Description() string
    Execute(ctx context.Context, params map[string]interface{}) (interface{}, error)
}

// 工具注册表
type ToolRegistry struct {
    tools map[string]Tool
}

// 注册工具
func (r *ToolRegistry) Register(tool Tool) {
    r.tools[tool.Name()] = tool
}

// 获取工具
func (r *ToolRegistry) GetTool(name string) (Tool, bool) {
    tool, ok := r.tools[name]
    return tool, ok
}
```

**工具列表**：
- `simple_fact_lookup`: 查找简单事实（Science Agent）
- `simple_dictionary`: 查找单词（Language Agent）
- `pronunciation_hint`: 发音提示（Language Agent）
- `image_generate_simple`: 生成示意图（Science Agent）
- `get_current_time`: 获取当前时间（Science Agent）

#### 1.2 eino工具调用集成
**目标**：集成eino框架的工具调用能力

**实现方案**：
```go
// 在ChatModel中注册工具
func (n *ScienceAgentNode) initChatModelWithTools(ctx context.Context) error {
    // 创建工具定义
    tools := []schema.Tool{
        {
            Type: schema.ToolTypeFunction,
            Function: &schema.FunctionDefinition{
                Name:        "simple_fact_lookup",
                Description: "查找简单事实，用于科学知识查询",
                Parameters: schema.FunctionParameters{
                    Type: schema.FunctionParametersTypeObject,
                    Properties: map[string]interface{}{
                        "query": map[string]interface{}{
                            "type":        "string",
                            "description": "查询关键词",
                        },
                    },
                    Required: []string{"query"},
                },
            },
        },
        // ... 其他工具
    }
    
    // 创建ChatModel配置，包含工具
    cfg := &ark.ChatModelConfig{
        Model: modelName,
        Tools: tools, // 注册工具
    }
    
    // 创建ChatModel
    chatModel, err := ark.NewChatModel(ctx, cfg)
    // ...
}
```

#### 1.3 工具调用处理
**目标**：处理ChatModel返回的工具调用请求

**实现方案**：
```go
// 在Domain Agent中处理工具调用
func (n *ScienceAgentNode) executeReal(ctx context.Context, ...) (*types.DomainAgentResponse, error) {
    messages, err := n.template.Format(ctx, map[string]any{...})
    
    // 调用ChatModel，可能返回工具调用请求
    result, err := n.chatModel.Generate(ctx, cleanMessages)
    if err != nil {
        return n.executeMock(...)
    }
    
    // 检查是否有工具调用请求
    if result.ToolCalls != nil && len(result.ToolCalls) > 0 {
        // 执行工具调用
        toolResults := make(map[string]interface{})
        toolsUsed := []string{}
        
        for _, toolCall := range result.ToolCalls {
            tool, ok := n.toolRegistry.GetTool(toolCall.Function.Name)
            if !ok {
                continue
            }
            
            // 解析参数
            params := make(map[string]interface{})
            json.Unmarshal([]byte(toolCall.Function.Arguments), &params)
            
            // 执行工具
            result, err := tool.Execute(ctx, params)
            if err != nil {
                n.logger.Errorw("工具调用失败", logx.Field("tool", toolCall.Function.Name), logx.Field("error", err))
                continue
            }
            
            toolResults[toolCall.Function.Name] = result
            toolsUsed = append(toolsUsed, toolCall.Function.Name)
        }
        
        // 将工具结果添加到消息中，重新调用ChatModel
        toolMessages := []*schema.Message{}
        for _, toolCall := range result.ToolCalls {
            if result, ok := toolResults[toolCall.Function.Name]; ok {
                toolMessages = append(toolMessages, schema.ToolMessage(
                    toolCall.ID,
                    fmt.Sprintf("%v", result),
                ))
            }
        }
        
        // 重新调用ChatModel，包含工具结果
        messages = append(messages, toolMessages...)
        finalResult, err := n.chatModel.Generate(ctx, messages)
        
        return &types.DomainAgentResponse{
            DomainType:  "Science",
            Content:     finalResult.Content,
            ToolsUsed:   toolsUsed,
            ToolResults: toolResults,
        }, nil
    }
    
    // 没有工具调用，直接返回
    return &types.DomainAgentResponse{
        DomainType:  "Science",
        Content:     result.Content,
        ToolsUsed:   []string{},
        ToolResults: make(map[string]interface{}),
    }, nil
}
```

### Phase 2: MCP资源集成

#### 2.1 MCP资源发现
**目标**：发现并集成可用的MCP资源

**可用MCP资源**（基于当前配置）：
- **地图服务**：`mcp_amap-amap-sse`
  - `maps_geo`: 地址转经纬度
  - `maps_regeocode`: 经纬度转地址
  - `maps_text_search`: 关键词搜索POI
  - `maps_around_search`: 周边搜索
  - `maps_weather`: 天气查询
- **搜索服务**：`mcp_mcpify-google-search`
  - `search_google_scholar`: 学术搜索
- **其他服务**：`mcp_tal_dify_MCP`、`mcp_jmeter`等

#### 2.2 MCP工具包装
**目标**：将MCP资源包装为Agent可用的工具

**实现方案**：
```go
// MCP工具包装器
type MCPToolWrapper struct {
    mcpServer string
    resource  string
    client    MCPClient
}

func (w *MCPToolWrapper) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
    // 调用MCP资源
    result, err := w.client.FetchResource(ctx, w.mcpServer, w.resource, params)
    if err != nil {
        return nil, err
    }
    
    // 格式化结果
    return w.formatResult(result), nil
}

// 注册MCP工具
func registerMCPTools(registry *ToolRegistry) {
    // 地图相关工具
    registry.Register(&MCPToolWrapper{
        mcpServer: "amap-amap-sse",
        resource:  "maps_geo",
        name:      "geo_lookup",
        description: "根据地址查找经纬度坐标",
    })
    
    registry.Register(&MCPToolWrapper{
        mcpServer: "amap-amap-sse",
        resource:  "maps_weather",
        name:      "weather_query",
        description: "查询指定城市的天气信息",
    })
    
    // 搜索相关工具
    registry.Register(&MCPToolWrapper{
        mcpServer: "mcpify-google-search",
        resource:  "search_google_scholar",
        name:      "scholar_search",
        description: "搜索学术论文和研究成果",
    })
}
```

#### 2.3 智能工具选择
**目标**：Supervisor根据问题类型智能选择工具

**实现方案**：
```go
// 在Supervisor中增加工具选择逻辑
func (n *SupervisorNode) SelectTools(ctx context.Context, intent string, message string) []string {
    tools := []string{}
    
    switch intent {
    case "探因型":
        // 科学问题可能需要查找事实
        if strings.Contains(message, "为什么") || strings.Contains(message, "怎么形成") {
            tools = append(tools, "simple_fact_lookup")
        }
        // 地理相关问题可能需要地图服务
        if strings.Contains(message, "哪里") || strings.Contains(message, "位置") {
            tools = append(tools, "geo_lookup")
        }
        // 天气相关问题
        if strings.Contains(message, "天气") || strings.Contains(message, "温度") {
            tools = append(tools, "weather_query")
        }
    case "表达型":
        // 语言问题需要字典和发音
        tools = append(tools, "simple_dictionary", "pronunciation_hint")
    case "认知型":
        // 认知问题可能需要查找事实
        if strings.Contains(message, "是什么") || strings.Contains(message, "特点") {
            tools = append(tools, "simple_fact_lookup")
        }
    }
    
    return tools
}
```

### Phase 3: 工具调用增强

#### 3.1 工具调用链
**目标**：支持多轮工具调用，形成调用链

**实现方案**：
```go
// 工具调用链
type ToolCallChain struct {
    calls []ToolCall
    maxDepth int
}

func (c *ToolCallChain) Execute(ctx context.Context, initialMessage string) (string, error) {
    currentMessage := initialMessage
    depth := 0
    
    for depth < c.maxDepth {
        // 调用ChatModel，可能返回工具调用请求
        result, err := c.chatModel.Generate(ctx, messages)
        if err != nil {
            return "", err
        }
        
        // 如果没有工具调用，返回结果
        if result.ToolCalls == nil || len(result.ToolCalls) == 0 {
            return result.Content, nil
        }
        
        // 执行工具调用
        toolResults := c.executeTools(ctx, result.ToolCalls)
        
        // 将工具结果添加到消息中
        messages = append(messages, toolResults...)
        depth++
    }
    
    return currentMessage, nil
}
```

#### 3.2 工具结果整合
**目标**：将工具调用结果自然整合到回答中

**实现方案**：
```go
// 在Domain Agent中整合工具结果
func (n *ScienceAgentNode) integrateToolResults(content string, toolResults map[string]interface{}) string {
    // 如果工具结果为空，直接返回内容
    if len(toolResults) == 0 {
        return content
    }
    
    // 使用ChatModel整合工具结果
    integrationPrompt := fmt.Sprintf(`请将以下工具调用结果自然整合到回答中：

原始回答：%s

工具调用结果：
%s

要求：
1. 工具结果要自然融入回答，不要显生硬
2. 保持回答的简洁性，不要过度引用工具结果
3. 如果工具结果与回答无关，可以忽略`, content, formatToolResults(toolResults))
    
    // 调用ChatModel整合
    integratedContent, err := n.integrationModel.Generate(ctx, integrationPrompt)
    if err != nil {
        return content // 整合失败，返回原始内容
    }
    
    return integratedContent
}
```

### Phase 4: Supervisor智能协调增强

#### 4.1 动态工具分配
**目标**：Supervisor根据问题动态决定是否需要工具

**实现方案**：
```go
// 在Supervisor中增加工具分配逻辑
func (n *SupervisorNode) CoordinateWithTools(ctx context.Context, state *types.SupervisorState, message string, chatHistory []*schema.Message) (*types.LearningPlanDecision, error) {
    // 1. 识别意图
    intentResult, err := n.intentAgent.RecognizeIntent(ctx, message, chatHistory)
    
    // 2. 判断是否需要工具
    needsTools := n.shouldUseTools(intentResult, message)
    
    // 3. 选择工具
    selectedTools := []string{}
    if needsTools {
        selectedTools = n.SelectTools(ctx, intentResult.Intent, message)
    }
    
    // 4. 制定学习计划（包含工具信息）
    decision, err := n.learningPlannerAgent.PlanLearningWithTools(ctx, intentResult, cognitiveLoadAdvice, selectedTools, ...)
    
    // 5. 将工具信息传递给Domain Agent
    decision.Tools = selectedTools
    
    return decision, nil
}
```

#### 4.2 工具使用策略
**目标**：定义不同场景下的工具使用策略

**策略定义**：
- **高置信度问题**：直接使用工具获取准确信息
- **探索性问题**：先回答，再提供工具增强信息
- **简单问题**：不使用工具，直接回答
- **复杂问题**：使用多个工具，整合结果

## 📝 实施计划

### Phase 1: 工具调用基础架构（1-2周）
1. ✅ 定义工具接口和注册机制
2. ✅ 实现基础工具（simple_fact_lookup、simple_dictionary等）
3. ✅ 集成eino工具调用能力
4. ✅ 在Domain Agent中实现工具调用处理
5. ✅ 测试工具调用流程

### Phase 2: MCP资源集成（1周）
1. ✅ 发现可用MCP资源
2. ✅ 实现MCP工具包装器
3. ✅ 注册MCP工具到工具注册表
4. ✅ 测试MCP工具调用

### Phase 3: 工具调用增强（1周）
1. ✅ 实现工具调用链
2. ✅ 实现工具结果整合
3. ✅ 优化工具调用性能
4. ✅ 测试工具调用链

### Phase 4: Supervisor智能协调增强（1周）
1. ✅ 实现动态工具分配
2. ✅ 定义工具使用策略
3. ✅ 优化Supervisor协调逻辑
4. ✅ 测试智能协调功能

## 🎯 预期效果

### 功能增强
1. **回答准确性提升**：通过工具调用获取准确信息
2. **回答丰富性提升**：通过MCP资源提供更多信息
3. **智能性提升**：Supervisor能够智能选择工具
4. **用户体验提升**：回答更加自然、准确、丰富

### 技术指标
- 工具调用成功率：≥90%
- 工具调用响应时间：≤2秒
- 工具结果整合质量：自然度≥80%
- 智能工具选择准确率：≥85%

## 🔍 技术细节

### 工具调用流程
```
1. Domain Agent接收问题
2. ChatModel生成回答，可能返回工具调用请求
3. 执行工具调用，获取结果
4. 将工具结果添加到消息中
5. 重新调用ChatModel，生成最终回答
6. 整合工具结果到回答中
```

### MCP集成流程
```
1. 发现MCP资源
2. 包装MCP资源为工具
3. 注册到工具注册表
4. Domain Agent调用工具
5. MCP工具包装器调用MCP资源
6. 格式化MCP结果
7. 返回给Domain Agent
```

## 📚 参考文档

- [Eino Tool Calling Documentation](https://www.cloudwego.io/zh/docs/eino/)
- [MCP Protocol Specification](https://modelcontextprotocol.io/)
- [Multi-Agent System Design](./research.md)

## ✅ 检查清单

### Phase 1
- [ ] 工具接口定义
- [ ] 工具注册机制
- [ ] eino工具调用集成
- [ ] Domain Agent工具调用处理
- [ ] 基础工具实现

### Phase 2
- [ ] MCP资源发现
- [ ] MCP工具包装器
- [ ] MCP工具注册
- [ ] MCP工具测试

### Phase 3
- [ ] 工具调用链实现
- [ ] 工具结果整合
- [ ] 性能优化
- [ ] 测试验证

### Phase 4
- [ ] 动态工具分配
- [ ] 工具使用策略
- [ ] Supervisor协调优化
- [ ] 集成测试

