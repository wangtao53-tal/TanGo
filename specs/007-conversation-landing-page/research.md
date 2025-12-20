# 研究文档: H5对话落地页

**创建日期**: 2025-12-19  
**功能**: H5对话落地页 - 流式对话、打字机效果、移动端优先设计

## 研究目标

1. Eino框架流式输出实现方案
2. 上下文窗口管理（20轮对话历史）
3. 基于年级的prompt生成策略
4. 前端打字机效果实现
5. 图片loading占位实现
6. 移动端优先的响应式设计

## 1. Eino框架流式输出实现

### 决策: 使用Eino ChatModel的Stream接口实现流式输出

**理由**:
- Eino框架提供了标准的Stream接口，支持流式返回
- 与现有代码架构一致，使用Ark模型实现
- 支持SSE (Server-Sent Events) 传输协议

**实现方式**:
```go
// 使用Eino ChatModel的Stream方法
streamReader, err := chatModel.Stream(ctx, messages)
if err != nil {
    return err
}

// 读取流式数据并发送到SSE连接
for {
    msg, err := streamReader.Read()
    if err != nil {
        if err == io.EOF {
            break
        }
        return err
    }
    
    // 发送到SSE连接
    sseEvent := fmt.Sprintf("event: message\ndata: %s\n\n", msg.Content)
    fmt.Fprintf(w, sseEvent)
    w.(http.Flusher).Flush()
}
```

**关键点**:
- 使用`chatModel.Stream(ctx, messages)`获取流式读取器
- 通过`streamReader.Read()`逐块读取数据
- 使用SSE格式发送数据: `event: message\ndata: {...}\n\n`
- 每次发送后调用`Flush()`确保数据立即传输

**替代方案**:
- 如果Eino Stream接口不可用，可以使用轮询方式模拟流式输出（不推荐）

## 2. 上下文窗口管理（20轮对话历史）

### 决策: 在内存存储中维护最近20轮对话，转换为Eino Message格式

**理由**:
- 20轮对话足够保持上下文连贯性
- 内存存储速度快，适合实时对话场景
- 需要将内部消息格式转换为Eino的schema.Message格式

**实现方式**:
```go
// 获取最近20轮对话（40条消息：用户+助手）
func (l *StreamLogic) getContextMessages(sessionId string, maxRounds int) []*schema.Message {
    allMessages := l.svcCtx.Storage.GetMessages(sessionId)
    
    // 只取最后maxRounds轮（maxRounds * 2条消息）
    start := 0
    if len(allMessages) > maxRounds*2 {
        start = len(allMessages) - maxRounds*2
    }
    
    // 转换为Eino Message格式
    einoMessages := make([]*schema.Message, 0)
    for i := start; i < len(allMessages); i++ {
        msg := allMessages[i].(types.ConversationMessage)
        if msg.Sender == "user" {
            einoMessages = append(einoMessages, schema.UserMessage(msg.Content))
        } else {
            einoMessages = append(einoMessages, schema.AssistantMessage(msg.Content, nil))
        }
    }
    
    return einoMessages
}
```

**关键点**:
- 限制为最近20轮（40条消息）
- 正确转换消息类型（user/assistant）
- 保持消息顺序

**替代方案**:
- 使用数据库存储历史消息（性能较低，不适合实时场景）
- 使用Redis存储（需要额外依赖）

## 3. 基于年级的prompt生成策略

### 决策: 根据用户年级（3-18岁）动态生成不同难度的prompt

**理由**:
- K12学生认知水平差异大，需要适配内容难度
- 通过prompt控制AI生成内容的复杂度
- 结合课外教育内容，拓展素质教育

**实现方式**:
```go
// 根据年级生成系统prompt
func generateSystemPrompt(age int, objectName string) string {
    var difficulty string
    var contentStyle string
    
    // 根据年龄确定难度和风格
    if age <= 6 {
        difficulty = "简单易懂，使用儿童语言"
        contentStyle = "生动有趣，多用比喻和故事"
    } else if age <= 12 {
        difficulty = "中等难度，使用日常语言"
        contentStyle = "结合生活实际，激发探索兴趣"
    } else {
        difficulty = "较高难度，可以使用专业术语"
        contentStyle = "深入浅出，培养科学思维"
    }
    
    return fmt.Sprintf(`你是一个面向%d岁学生的AI助手，专门帮助学生学习课外知识。
要求：
1. 使用%s的语言风格
2. 内容%s
3. 结合%s相关的科学知识、古诗词和英语表达
4. 拓展素质教育，培养探索精神
5. 内容贴合K12课程，但以课外拓展为主`, age, difficulty, contentStyle, objectName)
}
```

**年级分段策略**:
- **3-6岁（幼儿园）**: 简单词汇，故事化表达，多用图片和比喻
- **7-9岁（小学低年级）**: 日常语言，结合生活实际，激发兴趣
- **10-12岁（小学高年级）**: 中等难度，引入科学概念，培养思维
- **13-15岁（初中）**: 较高难度，可以使用专业术语，深入分析
- **16-18岁（高中）**: 高难度，专业术语，培养科学思维

**关键点**:
- Prompt模板化，支持动态参数注入
- 结合识别对象名称生成个性化prompt
- 平衡课内知识和课外拓展

**替代方案**:
- 使用固定的prompt模板（不够灵活）
- 完全依赖AI模型理解年级（不可控）

## 4. 前端打字机效果实现

### 决策: 使用React Hook实现逐字显示效果，支持流式数据更新

**理由**:
- React Hook可以封装状态逻辑，易于复用
- 支持流式数据实时更新
- 性能优化：使用requestAnimationFrame确保流畅度

**实现方式**:
```typescript
// useTypingEffect Hook
export function useTypingEffect(
  text: string,
  speed: number = 30 // 每30ms显示一个字符
): string {
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
  
  // 当text更新时，重置并重新开始
  useEffect(() => {
    setCurrentIndex(0);
    setDisplayedText('');
  }, [text]);
  
  return displayedText;
}
```

**流式更新支持**:
```typescript
// 在流式接收数据时，逐步更新text
const [streamingText, setStreamingText] = useState('');
const displayedText = useTypingEffect(streamingText);

// SSE事件处理
eventSource.onmessage = (event) => {
  const data = JSON.parse(event.data);
  if (data.type === 'text') {
    setStreamingText(prev => prev + data.content); // 追加新内容
  }
};
```

**关键点**:
- 使用useState管理当前显示文本和索引
- 使用useEffect处理定时器
- 支持流式数据追加更新
- 性能优化：避免不必要的重渲染

**替代方案**:
- 使用CSS动画（不够灵活，难以处理流式数据）
- 使用第三方库（增加依赖）

## 5. 图片loading占位实现

### 决策: 使用React组件实现图片loading占位，支持进度显示

**理由**:
- 提供良好的用户体验，避免空白等待
- 支持进度反馈，让用户知道图片正在生成
- 无缝替换为实际图片

**实现方式**:
```typescript
// ImagePlaceholder组件
export function ImagePlaceholder({ 
  imageUrl, 
  onLoad 
}: { 
  imageUrl?: string; 
  onLoad?: () => void;
}) {
  const [isLoading, setIsLoading] = useState(true);
  const [progress, setProgress] = useState(0);
  
  useEffect(() => {
    if (imageUrl) {
      // 模拟进度更新（实际应该从SSE事件获取）
      const interval = setInterval(() => {
        setProgress(prev => {
          if (prev >= 90) return prev;
          return prev + 10;
        });
      }, 200);
      
      // 预加载图片
      const img = new Image();
      img.onload = () => {
        setIsLoading(false);
        setProgress(100);
        onLoad?.();
        clearInterval(interval);
      };
      img.src = imageUrl;
      
      return () => clearInterval(interval);
    }
  }, [imageUrl, onLoad]);
  
  if (!imageUrl || isLoading) {
    return (
      <div className="image-placeholder">
        <div className="loading-spinner" />
        <div className="progress-bar">
          <div style={{ width: `${progress}%` }} />
        </div>
        <p>正在生成图片... {progress}%</p>
      </div>
    );
  }
  
  return <img src={imageUrl} alt="Generated" />;
}
```

**SSE事件处理**:
```typescript
// 接收图片生成进度事件
eventSource.addEventListener('image_progress', (event) => {
  const data = JSON.parse(event.data);
  setImageProgress(data.progress); // 0-100
});

// 接收图片完成事件
eventSource.addEventListener('image_done', (event) => {
  const data = JSON.parse(event.data);
  setImageUrl(data.url);
});
```

**关键点**:
- 显示loading动画和进度条
- 支持从SSE事件获取真实进度
- 图片加载完成后无缝替换
- 错误处理：图片加载失败时显示错误提示

**替代方案**:
- 使用简单的loading spinner（用户体验较差）
- 不显示进度（用户不知道等待时间）

## 6. 移动端优先的响应式设计

### 决策: 使用Tailwind CSS实现移动端优先的响应式设计

**理由**:
- Tailwind CSS提供强大的响应式工具类
- 移动端优先的设计理念符合项目规范
- 支持触摸、滑动、捏合等移动端交互

**实现方式**:
```typescript
// Tailwind配置 - 移动端优先断点
// tailwind.config.js
export default {
  theme: {
    screens: {
      'sm': '640px',   // 小屏设备（大手机）
      'md': '768px',   // 平板
      'lg': '1024px',  // 小桌面
      'xl': '1280px',  // 大桌面
    },
  },
}

// 组件样式 - 移动端优先
<div className="
  flex flex-col          // 移动端：垂直布局
  md:flex-row           // 平板及以上：水平布局
  gap-2                 // 移动端：小间距
  md:gap-4              // 平板及以上：大间距
  p-4                   // 移动端：小内边距
  md:p-6                // 平板及以上：大内边距
">
```

**触摸交互支持**:
```typescript
// 支持触摸事件
<div
  onTouchStart={handleTouchStart}
  onTouchMove={handleTouchMove}
  onTouchEnd={handleTouchEnd}
  className="touch-pan-y" // Tailwind: 支持垂直滑动
>
  {/* 内容 */}
</div>
```

**关键点**:
- 移动端优先：先设计移动端样式，再适配桌面端
- 使用Tailwind响应式前缀（sm:, md:, lg:, xl:）
- 支持触摸事件和手势操作
- 兼容PC端鼠标事件（hover, click等）

**替代方案**:
- 使用CSS Media Queries（代码冗长）
- 使用CSS-in-JS库（增加运行时开销）

## 7. 对话布局：用户消息右侧，AI消息左侧

### 决策: 使用Flexbox布局，根据消息发送者调整对齐方式

**理由**:
- 符合常见聊天应用的设计模式
- 清晰区分用户和AI消息
- 移动端和PC端都能良好显示

**实现方式**:
```typescript
// ConversationMessage组件
<div className={`
  flex
  ${message.sender === 'user' ? 'justify-end' : 'justify-start'}
  mb-4
`}>
  <div className={`
    max-w-[80%]           // 移动端：最大宽度80%
    md:max-w-[60%]       // 桌面端：最大宽度60%
    px-4 py-3
    rounded-2xl
    ${message.sender === 'user' 
      ? 'bg-primary text-white rounded-br-sm'  // 用户：右侧，主色背景
      : 'bg-gray-100 text-gray-800 rounded-bl-sm'  // AI：左侧，灰色背景
    }
  `}>
    {message.content}
  </div>
</div>
```

**关键点**:
- 用户消息：右侧对齐，主色背景
- AI消息：左侧对齐，灰色背景
- 响应式最大宽度：移动端80%，桌面端60%
- 圆角设计：用户消息右下角小圆角，AI消息左下角小圆角

## 总结

### 技术选型

1. **后端流式输出**: Eino ChatModel.Stream接口 + SSE协议
2. **上下文管理**: 内存存储 + 最近20轮限制
3. **Prompt生成**: 基于年级的动态模板生成
4. **前端打字机效果**: React Hook + 流式数据更新
5. **图片loading**: React组件 + SSE进度事件
6. **响应式设计**: Tailwind CSS移动端优先

### 实施优先级

1. **P1**: Eino流式输出实现（核心功能）
2. **P1**: 上下文窗口管理（核心功能）
3. **P2**: 基于年级的prompt生成（体验优化）
4. **P2**: 前端打字机效果（体验优化）
5. **P2**: 图片loading占位（体验优化）
6. **P3**: 移动端优先响应式设计（体验优化）

### 风险与缓解

1. **Eino Stream接口不可用**: 使用轮询方式模拟（不推荐，但可作为降级方案）
2. **流式输出性能问题**: 优化SSE连接管理，支持重连机制
3. **移动端性能**: 使用React.memo优化组件渲染，避免不必要的重渲染
4. **浏览器兼容性**: 检测SSE支持，不支持时降级到轮询

