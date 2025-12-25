# 快速开始：H5对话落地页

**功能**: 007-conversation-landing-page  
**创建日期**: 2025-12-19

## 概述

本文档提供H5对话落地页功能的快速开始指南，包括环境配置、开发流程、测试方法等。

## 前置条件

### 后端环境

- Go 1.21+
- go-zero框架
- Eino框架（github.com/cloudwego/eino）
- Eino-ext（github.com/cloudwego/eino-ext）
- 配置Eino服务地址和认证信息

### 前端环境

- Node.js 18+
- React 19.2+
- Tailwind CSS 4.1+
- 现代浏览器（支持SSE）

## 环境配置

### 后端配置

1. **Eino配置**（在`backend/etc/explore.yaml`中）:
```yaml
AI:
  EinoBaseURL: "https://your-eino-service.com"
  AppID: "your-app-id"
  AppKey: "your-app-key"
  TextGenerationModel: "gpt-5-nano"  # 文本生成模型
  MaxContextRounds: 20                # 最大上下文轮次
```

2. **启动后端服务**:
```bash
cd backend
go run explore.go
```

### 前端配置

1. **环境变量**（在`frontend/.env`中）:
```env
VITE_API_BASE_URL=http://localhost:8877
VITE_BACKEND_HOST=localhost
VITE_BACKEND_PORT=8877
```

2. **安装依赖**:
```bash
cd frontend
npm install
```

3. **启动开发服务器**:
```bash
npm run dev
```

## 开发流程

### 1. 后端开发

#### 步骤1: 创建对话节点（基于Eino Graph）

创建`backend/internal/agent/nodes/conversation_node.go`:

```go
package nodes

import (
    "context"
    "github.com/cloudwego/eino/components/model"
    "github.com/cloudwego/eino/schema"
    "github.com/cloudwego/eino/components/prompt"
)

type ConversationNode struct {
    chatModel model.ChatModel
    template  prompt.ChatTemplate
}

func (n *ConversationNode) StreamConversation(
    ctx context.Context,
    message string,
    contextMessages []*schema.Message,
    userAge int,
) (*schema.StreamReader[*schema.Message], error) {
    // 根据用户年级生成系统prompt
    systemPrompt := generateSystemPrompt(userAge)
    
    // 构建消息列表
    messages := []*schema.Message{
        schema.SystemMessage(systemPrompt),
    }
    messages = append(messages, contextMessages...)
    messages = append(messages, schema.UserMessage(message))
    
    // 调用流式接口
    return n.chatModel.Stream(ctx, messages)
}
```

#### 步骤2: 实现流式逻辑

更新`backend/internal/logic/streamlogic.go`:

```go
func (l *StreamLogic) StreamConversation(
    w http.ResponseWriter,
    req *StreamConversationRequest,
) error {
    // 设置SSE响应头
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")
    
    // 获取最近20轮对话历史
    contextMessages := l.getContextMessages(req.SessionId, 20)
    
    // 调用对话节点
    streamReader, err := l.conversationNode.StreamConversation(
        l.ctx,
        req.Message,
        contextMessages,
        req.UserAge,
    )
    if err != nil {
        return err
    }
    
    // 读取流式数据并发送
    for {
        msg, err := streamReader.Read()
        if err != nil {
            if err == io.EOF {
                break
            }
            return err
        }
        
        // 发送SSE事件
        event := fmt.Sprintf("event: message\ndata: %s\n\n", msg.Content)
        fmt.Fprintf(w, event)
        w.(http.Flusher).Flush()
    }
    
    return nil
}
```

#### 步骤3: 更新Handler

更新`backend/internal/handler/streamhandler.go`:

```go
func StreamConversationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // 解析请求参数
        sessionId := r.URL.Query().Get("sessionId")
        message := r.URL.Query().Get("message")
        userAge, _ := strconv.Atoi(r.URL.Query().Get("userAge"))
        
        // 调用流式逻辑
        logic := logic.NewStreamLogic(r.Context(), svcCtx)
        err := logic.StreamConversation(w, &types.StreamConversationRequest{
            SessionId: sessionId,
            Message:   message,
            UserAge:   userAge,
        })
        
        if err != nil {
            // 发送错误事件
            fmt.Fprintf(w, "event: error\ndata: %s\n\n", err.Error())
        }
    }
}
```

### 2. 前端开发

#### 步骤1: 创建流式对话Hook

创建`frontend/src/hooks/useStreamConversation.ts`:

```typescript
import { useState, useEffect, useRef } from 'react';
import { createSSEConnection } from '../services/sse';

export function useStreamConversation(sessionId: string) {
  const [streamingText, setStreamingText] = useState('');
  const [isStreaming, setIsStreaming] = useState(false);
  const eventSourceRef = useRef<EventSource | null>(null);
  
  const startStream = (message: string, userAge?: number) => {
    setIsStreaming(true);
    setStreamingText('');
    
    const url = `/api/conversation/stream?sessionId=${sessionId}&message=${encodeURIComponent(message)}${userAge ? `&userAge=${userAge}` : ''}`;
    const eventSource = new EventSource(url);
    eventSourceRef.current = eventSource;
    
    eventSource.onmessage = (event) => {
      const data = JSON.parse(event.data);
      if (data.type === 'message') {
        setStreamingText(prev => prev + data.content);
      } else if (data.type === 'done') {
        setIsStreaming(false);
        eventSource.close();
      }
    };
    
    eventSource.onerror = () => {
      setIsStreaming(false);
      eventSource.close();
    };
  };
  
  const stopStream = () => {
    if (eventSourceRef.current) {
      eventSourceRef.current.close();
      setIsStreaming(false);
    }
  };
  
  useEffect(() => {
    return () => {
      stopStream();
    };
  }, []);
  
  return { streamingText, isStreaming, startStream, stopStream };
}
```

#### 步骤2: 创建打字机效果Hook

创建`frontend/src/hooks/useTypingEffect.ts`:

```typescript
import { useState, useEffect } from 'react';

export function useTypingEffect(text: string, speed: number = 30) {
  const [displayedText, setDisplayedText] = useState('');
  const [currentIndex, setCurrentIndex] = useState(0);
  
  useEffect(() => {
    if (currentIndex < text.length) {
      const timer = setTimeout(() => {
        setDisplayedText(text.slice(0, currentIndex + 1));
        setCurrentIndex(currentIndex + 1);
      }, speed);
      
      return () => clearTimeout(timer);
    }
  }, [currentIndex, text, speed]);
  
  useEffect(() => {
    setCurrentIndex(0);
    setDisplayedText('');
  }, [text]);
  
  return displayedText;
}
```

#### 步骤3: 创建图片loading占位组件

创建`frontend/src/components/common/ImagePlaceholder.tsx`:

```typescript
import { useState, useEffect } from 'react';

export function ImagePlaceholder({ 
  imageUrl, 
  progress = 0,
  onLoad 
}: { 
  imageUrl?: string; 
  progress?: number;
  onLoad?: () => void;
}) {
  const [isLoading, setIsLoading] = useState(true);
  
  useEffect(() => {
    if (imageUrl) {
      const img = new Image();
      img.onload = () => {
        setIsLoading(false);
        onLoad?.();
      };
      img.src = imageUrl;
    }
  }, [imageUrl, onLoad]);
  
  if (!imageUrl || isLoading) {
    return (
      <div className="image-placeholder bg-gray-100 rounded-lg p-8 flex flex-col items-center justify-center">
        <div className="loading-spinner animate-spin rounded-full h-12 w-12 border-b-2 border-primary mb-4" />
        <div className="w-full bg-gray-200 rounded-full h-2 mb-2">
          <div 
            className="bg-primary h-2 rounded-full transition-all duration-300"
            style={{ width: `${progress}%` }}
          />
        </div>
        <p className="text-sm text-gray-600">正在生成图片... {progress}%</p>
      </div>
    );
  }
  
  return <img src={imageUrl} alt="Generated" className="rounded-lg" />;
}
```

#### 步骤4: 更新对话页面组件

更新`frontend/src/pages/Result.tsx`:

```typescript
import { useStreamConversation } from '../hooks/useStreamConversation';
import { useTypingEffect } from '../hooks/useTypingEffect';
import { ImagePlaceholder } from '../components/common/ImagePlaceholder';

export default function Result() {
  const { streamingText, isStreaming, startStream } = useStreamConversation(sessionId);
  const displayedText = useTypingEffect(streamingText);
  
  const handleSendMessage = async (text: string) => {
    // 立即显示用户消息
    setMessages(prev => [...prev, {
      id: `msg-${Date.now()}`,
      type: 'text',
      sender: 'user',
      content: text,
      timestamp: new Date().toISOString(),
      sessionId,
    }]);
    
    // 开始流式对话
    startStream(text, userAge);
  };
  
  // 渲染AI消息（带打字机效果）
  const renderAssistantMessage = (message: ConversationMessage) => {
    if (message.isStreaming) {
      return (
        <div className="assistant-message">
          {displayedText}
          <span className="typing-cursor">|</span>
        </div>
      );
    }
    return <div className="assistant-message">{message.content}</div>;
  };
  
  // ...
}
```

## 测试方法

### 1. 单元测试

**后端测试**:
```bash
cd backend
go test ./internal/logic/... -v
go test ./internal/agent/nodes/... -v
```

**前端测试**:
```bash
cd frontend
npm test
```

### 2. 集成测试

**测试流式对话**:
```bash
# 启动后端服务
cd backend && go run explore.go

# 在另一个终端测试SSE连接
curl -N "http://localhost:8877/api/conversation/stream?sessionId=test-123&message=你好&userAge=8"
```

**预期输出**:
```
event: connected
data: {"type":"connected","sessionId":"test-123"}

event: message
data: {"type":"message","content":"你","index":0}

event: message
data: {"type":"message","content":"好","index":1}

...

event: done
data: {"type":"done"}
```

### 3. 端到端测试

1. **启动服务**:
   - 后端: `cd backend && go run explore.go`
   - 前端: `cd frontend && npm run dev`

2. **测试流程**:
   - 访问 `http://localhost:5173`
   - 完成拍照识别
   - 进入对话页面
   - 发送消息，观察打字机效果
   - 测试图片生成loading占位

3. **验证点**:
   - ✅ 用户消息显示在右侧
   - ✅ AI消息显示在左侧
   - ✅ 打字机效果流畅
   - ✅ 图片loading占位正确显示
   - ✅ 上下文关联正确（多轮对话）

## 常见问题

### Q1: SSE连接失败

**原因**: 浏览器不支持SSE或网络问题

**解决**: 
- 检查浏览器版本（Chrome 90+, Safari 14+, Firefox 88+）
- 检查网络连接
- 使用非流式接口作为降级方案

### Q2: 打字机效果卡顿

**原因**: 渲染性能问题

**解决**:
- 使用React.memo优化组件
- 减少不必要的重渲染
- 使用requestAnimationFrame优化动画

### Q3: 上下文窗口超过20轮

**原因**: 消息数量未正确限制

**解决**:
- 检查`getContextMessages`函数的实现
- 确保只返回最近20轮消息
- 添加日志验证消息数量

### Q4: 图片loading占位不显示

**原因**: SSE事件未正确接收或处理

**解决**:
- 检查SSE事件类型是否为`image_progress`
- 验证事件数据格式
- 检查前端事件监听器

## 下一步

1. **性能优化**: 
   - 优化流式输出性能
   - 优化前端渲染性能
   - 添加缓存机制

2. **功能扩展**:
   - 支持语音输入
   - 支持图片上传
   - 支持多模态输入

3. **用户体验优化**:
   - 优化移动端交互
   - 添加错误重试机制
   - 优化加载状态显示

