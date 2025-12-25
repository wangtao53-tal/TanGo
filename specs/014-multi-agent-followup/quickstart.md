# Quick Start: 多Agent追问功能优化

**Date**: 2025-01-27  
**Feature**: 多Agent追问功能优化

## 概述

多Agent追问功能通过Follow-up Supervisor协调8个子Agent协作，实现智能追问回答。创建新的 `/api/conversation/agent` 接口作为多Agent模式接口，保留 `/api/conversation/stream` 接口作为单Agent模式接口。

## 架构设计

### 多Agent系统架构

```
[Start]
   ↓
[Follow-up Supervisor] - 核心中枢，协调各个Agent
   ↓
[Intent Agent] - 识别意图类型（认知型、探因型、表达型、游戏型、情绪型）
   ↓
[Cognitive Load Agent] - 判断认知负载（简短讲解、类比讲解、反问引导、暂停探索）
   ↓
[Learning Planner Agent] - 决定教学动作（选择领域Agent、讲一点/问一个问题）
   ↓
 ┌───────────────┬───────────────┬───────────────┐
 ▼               ▼               ▼
[Science Agent] [Language Agent] [Humanities Agent]
   │               │               │
 (Tool?)         (Tool?)         (Tool?)
   └──────┬────────┴────────┬──────┘
          ▼
   [Interaction Agent] - 优化交互方式，添加可选动作
          ↓
   [Reflection Agent] - 反思学习状态（兴趣、困惑、放松）
          ↓
   [Memory Agent] - 记录学习状态（感兴趣的主题、已理解/未理解的点）
          ↓
   [End] - 返回最终回答
```

### 接口设计

- **`/api/conversation/stream`**: 单Agent模式接口（保持不变，向后兼容）
- **`/api/conversation/agent`**: 多Agent模式接口（新创建）

两个接口的输入输出格式完全一致，前端可以无缝切换。

## 后端实现

### 1. 创建Agent节点

在 `backend/internal/agent/nodes/` 目录下创建以下Agent节点：

- `supervisor_node.go`: Supervisor节点
- `intent_agent_node.go`: Intent Agent节点
- `cognitive_load_node.go`: Cognitive Load Agent节点
- `learning_planner_node.go`: Learning Planner Agent节点
- `science_agent_node.go`: Science Agent节点
- `language_agent_node.go`: Language Agent节点
- `humanities_agent_node.go`: Humanities Agent节点
- `interaction_agent_node.go`: Interaction Agent节点
- `reflection_agent_node.go`: Reflection Agent节点
- `memory_agent_node.go`: Memory Agent节点

### 2. 创建MultiAgentGraph

在 `backend/internal/agent/` 目录下创建 `multiagent_graph.go`，实现多Agent Graph结构：

```go
type MultiAgentGraph struct {
    ctx    context.Context
    config config.AIConfig
    logger logx.Logger
    
    // Agent节点实例
    supervisorNode      *nodes.SupervisorNode
    intentAgentNode     *nodes.IntentAgentNode
    cognitiveLoadNode   *nodes.CognitiveLoadNode
    learningPlannerNode *nodes.LearningPlannerNode
    scienceAgentNode    *nodes.ScienceAgentNode
    languageAgentNode   *nodes.LanguageAgentNode
    humanitiesAgentNode *nodes.HumanitiesAgentNode
    interactionAgentNode *nodes.InteractionAgentNode
    reflectionAgentNode *nodes.ReflectionAgentNode
    memoryAgentNode    *nodes.MemoryAgentNode
}

// ExecuteMultiAgentConversation 执行多Agent对话流程
func (g *MultiAgentGraph) ExecuteMultiAgentConversation(
    ctx context.Context,
    req *types.UnifiedStreamConversationRequest,
) (*schema.StreamReader[*schema.Message], error) {
    // 1. Supervisor分析上下文
    // 2. Intent Agent识别意图
    // 3. Cognitive Load Agent判断认知负载
    // 4. Learning Planner Agent决定教学动作
    // 5. Domain Agent生成回答
    // 6. Interaction Agent优化交互
    // 7. Reflection Agent反思状态
    // 8. Memory Agent记录状态
    // 9. 返回流式回答
}
```

### 3. 创建新接口Handler和Logic

- `backend/internal/handler/agenthandler.go`: 多Agent接口处理器
- `backend/internal/logic/agentlogic.go`: 多Agent逻辑层

### 4. 更新路由配置

在 `backend/internal/handler/routes.go` 中添加新路由：

```go
{
    Method:  http.MethodPost,
    Path:    "/api/conversation/agent",
    Handler: AgentConversationHandler(serverCtx),
},
```

## 前端实现

### 1. 创建API配置模块

创建 `frontend/src/config/api.ts`:

```typescript
// API配置
export const API_CONFIG = {
  // 是否使用多Agent模式
  useMultiAgent: import.meta.env.VITE_USE_MULTI_AGENT === 'true' || 
                 localStorage.getItem('useMultiAgent') === 'true',
  
  // 根据配置选择接口
  getConversationEndpoint(): string {
    return this.useMultiAgent 
      ? '/api/conversation/agent' 
      : '/api/conversation/stream';
  }
};
```

### 2. 更新对话服务

更新 `frontend/src/services/conversation.ts`:

```typescript
import { API_CONFIG } from '@/config/api';

export async function streamConversation(
  request: UnifiedStreamConversationRequest
): Promise<EventSource> {
  const endpoint = API_CONFIG.getConversationEndpoint();
  // 使用endpoint调用接口
}
```

### 3. 支持运行时切换

在设置页面添加接口切换选项：

```typescript
function toggleMultiAgentMode(enabled: boolean) {
  localStorage.setItem('useMultiAgent', enabled.toString());
  // 重新加载页面或更新配置
}
```

## 配置说明

### 后端配置

无需额外配置，使用现有的Eino配置即可。

### 前端配置

**环境变量配置**（`.env`）:
```bash
VITE_USE_MULTI_AGENT=true  # 启用多Agent模式
```

**localStorage配置**（运行时）:
```javascript
localStorage.setItem('useMultiAgent', 'true');  // 启用多Agent模式
localStorage.setItem('useMultiAgent', 'false'); // 使用单Agent模式
```

**默认值**: `false`（使用单Agent模式，向后兼容）

## 测试

### 后端测试

1. 测试Supervisor决策逻辑
2. 测试Intent Agent意图识别
3. 测试Cognitive Load Agent认知负载判断
4. 测试Learning Planner Agent教学决策
5. 测试Domain Agent专业回答
6. 测试Interaction Agent交互优化
7. 测试Reflection Agent反思判断
8. 测试Memory Agent记忆记录
9. 测试Graph执行流程完整性
10. 测试新接口功能

### 前端测试

1. 测试API配置读取
2. 测试接口选择逻辑
3. 测试接口切换功能
4. 测试错误处理和降级机制

## 部署

### 开发环境

1. 前端设置 `VITE_USE_MULTI_AGENT=false`（默认值，使用单Agent模式）
2. 后端启动服务
3. 前端启动开发服务器

### 生产环境

1. 前端设置 `VITE_USE_MULTI_AGENT=true`（启用多Agent模式）
2. 后端部署服务
3. 前端构建并部署

## 故障排查

### 多Agent模式失败

1. 检查Eino配置是否正确
2. 检查Agent节点是否正常初始化
3. 检查Graph执行流程是否完整
4. 查看错误日志
5. 前端自动降级到单Agent模式

### 接口切换问题

1. 检查前端配置是否正确读取
2. 检查接口路径是否正确
3. 检查请求参数格式是否一致
4. 检查响应格式是否一致

## 下一步

1. 实现各个Agent节点
2. 实现MultiAgentGraph
3. 实现新接口Handler和Logic
4. 更新前端配置和对话服务
5. 编写测试用例
6. 部署和验证

