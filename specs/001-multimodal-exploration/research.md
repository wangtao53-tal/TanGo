# 技术研究与决策文档

**创建日期**: 2025-12-18  
**功能**: TanGo 多模态探索核心功能

## 技术选型决策

### 1. eino框架集成

**决策**: 使用字节的eino框架统一管理AI模型调用

**理由**:
- eino是字节开源的云原生AI框架，专为AI模型集成设计
- 提供统一的模型调用接口，简化多模型管理
- 支持模型配置管理和监控，便于运维
- 与go-zero框架兼容，可以无缝集成

**替代方案考虑**:
- 直接调用AI模型API（如OpenAI API、Claude API）
  - 优点：实现简单，无需额外框架
  - 缺点：需要自行管理配置、错误处理、重试机制
  - 决策：不采用，因为eino提供更好的统一管理

**实施要点**:
- 需要在eino配置文件中定义Vision Model和LLM Model
- 通过APP ID配置模型访问权限
- 实现统一的错误处理和重试机制

### 2. 图像识别模型选择

**决策**: 待定，根据APP ID申请情况选择

**候选方案**:

**方案A: GPT-4V**
- 优点：通用性强，识别准确率高，支持中文
- 缺点：成本较高，响应时间可能较长
- 适用场景：需要高准确率的通用识别

**方案B: Claude 3.5 Sonnet with Vision**
- 优点：识别准确率高，对中文支持好，响应速度快
- 缺点：需要API访问权限
- 适用场景：优先考虑（如果可用）

**方案C: 专用图像分类模型（通过eino集成）**
- 优点：针对性强，可能成本更低
- 缺点：需要训练或找到合适的预训练模型
- 适用场景：如果字节内部有合适的模型

**选择标准**:
1. 识别准确率≥90%（目标95%+）
2. 支持中文输出
3. 响应时间≤3秒
4. 支持80-100种常见对象识别
5. 成本可控（黑客马拉松阶段）

**实施建议**:
- 优先尝试Claude 3.5 Sonnet with Vision
- 如果不可用，使用GPT-4V
- 通过eino框架统一调用，便于后续切换

### 3. 知识卡片生成模型选择

**决策**: 使用GPT-4或Claude 3.5 Sonnet

**理由**:
- 内容生成质量高，符合K12教育要求
- 支持长文本生成（三张卡片内容）
- 支持结构化输出（JSON格式）
- 支持中文内容生成
- 可以通过prompt控制内容难度和风格

**Prompt设计要点**:
1. **科学认知卡Prompt**:
   - 输入：对象名称、类别、年龄
   - 输出：名称、一句话解释、2-3个关键事实、1个趣味知识点
   - 难度控制：根据年龄调整语言复杂度和内容深度

2. **古诗词/人文卡Prompt**:
   - 输入：对象名称、类别、年龄
   - 输出：相关诗词（节选）、白话解释、情境联想
   - 重点：不是背诗，而是"看到→联想到"

3. **英语表达卡Prompt**:
   - 输入：对象名称、类别、年龄
   - 输出：核心单词、1-2句口语级表达
   - 重点：强调"会说"，不考语法

**内容质量控制**:
- 实现内容验证机制（检查是否包含必要字段）
- 实现内容过滤（确保内容适合儿童）
- 实现重试机制（如果生成内容不符合要求）

### 4. AI Agent框架选择

**决策**: 使用ReAct Agent（单Agent模式）

**理由**:
1. **简单直接**: 我们的场景主要是顺序调用两个AI模型，不需要复杂的多Agent协作
2. **易于实现**: ReAct Agent模式适合"推理-行动-观察"的循环，我们的流程是线性的
3. **快速开发**: 适合黑客马拉松快速实现，降低复杂度
4. **易于调试**: 单Agent模式更容易追踪和调试

**架构设计**:
```
用户拍照
  ↓
ReAct Agent启动
  ↓
Action 1: 调用Vision Model（图像识别）
  ↓
Observation: 获取识别结果
  ↓
Reasoning: 判断识别是否成功
  ↓
Action 2: 调用LLM Model（知识卡片生成）
  ↓
Observation: 获取三张卡片内容
  ↓
Reasoning: 验证内容质量
  ↓
Final Action: 返回结果
```

**实现方式**:
- 在go-zero的logic层实现ReAct Agent逻辑
- 使用eino框架封装AI模型调用
- Agent负责协调Vision Model和LLM Model的调用顺序
- 实现简单的错误重试机制

**不采用Multi-Agent的原因**:
- 过度设计：我们的场景不需要多个Agent协作
- 复杂度高：需要Agent间通信、任务分配等机制
- 开发时间长：不适合黑客马拉松快速实现

### 5. 前端本地存储方案

**决策**: 使用IndexedDB + localStorage混合方案

**理由**:
- **IndexedDB**: 
  - 存储大量探索记录和知识卡片（支持结构化数据）
  - 支持索引查询，性能好
  - 存储容量大（通常几GB）
- **localStorage**: 
  - 存储用户设置（年龄、年级等简单配置）
  - 访问速度快
  - 适合小数据量

**数据结构设计**:

**IndexedDB数据库: "TanGoDB"**
- **ObjectStore: "explorations"**
  - 字段：id, timestamp, objectName, objectCategory, age, cards, imageData
- **ObjectStore: "cards"**
  - 字段：id, explorationId, cardType, content, collectedAt

**localStorage**:
- key: "userProfile"
- value: { age: number, grade: string }

**实施要点**:
- 使用idb库（IndexedDB的Promise封装）简化操作
- 实现数据迁移机制（如果数据结构变更）
- 实现数据清理机制（避免存储过多数据）

### 6. 分享链接实现方案

**决策**: 使用临时分享链接，数据存储在内存或Redis

**理由**:
- 无需账户系统，简化实现
- 临时存储，自动过期，保护隐私
- 实现简单，适合MVP版本

**实现方式**:
1. 前端生成分享请求（包含探索记录和卡片数据）
2. 后端创建临时存储（内存或Redis，TTL 7天）
3. 返回分享链接（包含唯一shareId）
4. 家长端通过链接ID获取数据

**数据结构**:
```json
{
  "shareId": "uuid-v4",
  "explorationRecords": [
    {
      "id": "exp-1",
      "timestamp": "2025-12-18T10:00:00Z",
      "objectName": "银杏",
      "objectCategory": "自然类",
      "age": 8,
      "cards": [...]
    }
  ],
  "collectedCards": [...],
  "createdAt": "2025-12-18T10:00:00Z",
  "expiresAt": "2025-12-25T10:00:00Z"
}
```

**存储选择**:
- MVP版本：使用内存存储（map[string]ShareData）
- 生产版本：使用Redis（支持分布式和持久化）

**安全考虑**:
- shareId使用UUID v4，不可预测
- 设置TTL，自动过期
- 不存储敏感信息（如原始图片）

### 7. 前端技术栈确认

**决策**: React 18 + Vite + Tailwind CSS

**理由**:
- **React 18**: 最新稳定版本，支持并发特性
- **Vite**: 快速开发服务器，构建速度快
- **Tailwind CSS**: 快速构建UI，移动端优先设计友好

**移动端优先实现**:
- 使用Tailwind的响应式设计（sm:, md:, lg:断点）
- 默认移动端样式，PC端通过断点扩展
- 使用触摸友好的交互（大按钮、手势支持）

**设计稿集成**:
- 参考：https://stitch.withgoogle.com/projects/3345948492361031874
- 根据HTML和设计图进行组件化实现
- 保持设计风格一致（儿童友好、简洁有趣）

## 待确认事项

1. **AI模型APP ID**: 需要申请Vision Model和LLM Model的APP ID
2. **eino配置**: 需要确认eino与go-zero的具体集成方式
3. **模型选择**: 根据APP ID申请情况确定具体使用的模型
4. ~~**设计稿**: 需要获取完整的HTML和设计图文件~~ ✅ 已完成：设计稿已提供在 `stitch_ui/` 文件夹

## 风险评估

1. **AI模型可用性**: 如果APP ID申请失败，需要备用方案
   - 风险等级：中
   - 缓解措施：准备直接调用公开API的备用方案

2. **响应时间**: AI模型调用可能超过5秒
   - 风险等级：中
   - 缓解措施：实现异步处理，前端显示加载状态

3. **内容质量**: AI生成的内容可能不符合K12要求
   - 风险等级：中
   - 缓解措施：优化prompt设计，实现内容验证和过滤

4. **识别准确率**: 可能无法达到90%的目标
   - 风险等级：低
   - 缓解措施：选择高质量模型，实现识别结果验证

### 7. 前端UI实现方案

**决策**: 完全基于提供的UI设计稿（`stitch_ui/` 文件夹）实现，使用React + Tailwind CSS

**理由**:
- 设计稿已经完成，包含所有页面的HTML和样式
- 使用Tailwind CSS可以快速实现设计稿中的样式
- React组件化可以复用设计稿中的结构

**设计稿分析结果**:
1. **页面结构**: 已识别所有页面（首页、拍照、结果、详情、收藏、报告）
2. **颜色系统**: 已提取所有颜色值（primary green, sunny orange, sky blue等）
3. **组件模式**: 已识别可复用组件（按钮、卡片、对话框等）
4. **动画效果**: 已识别所有动画（float, pulse-glow, scan-line等）

**实现策略**:
- 将设计稿中的HTML结构转换为React组件
- 使用Tailwind CSS配置匹配设计稿的颜色和主题
- 实现所有动画效果（CSS动画 + Tailwind动画类）
- 确保响应式设计（移动端优先）

**Tailwind配置要点**:
```javascript
// tailwind.config.js 需要包含的设计稿颜色
colors: {
  primary: "#4cdf20",
  "science-green": "#76FF7A",
  "sunny-orange": "#FF9E64",
  "sky-blue": "#40C4FF",
  "warm-yellow": "#FFE580",
  "peach-pink": "#FFB7C5",
  "cloud-white": "#F8F8F8",
  // ... 其他颜色
}
```

**组件化策略**:
- 通用组件：Button, Card, Header, LittleStar
- 功能组件：CameraView, CardCarousel, CollectionGrid
- 页面组件：Home, Capture, Result, Collection, Share, LearningReport

## 最佳实践参考

1. **ReAct Agent实现**: 参考LangChain的ReAct实现模式
2. **eino集成**: 参考eino官方文档和示例
3. **go-zero最佳实践**: 参考go-zero官方文档
4. **React移动端优化**: 参考React Native Web的最佳实践
5. **儿童友好UI设计**: 参考Duolingo、Khan Academy Kids等教育应用
6. **Tailwind CSS设计系统**: 参考设计稿中的颜色和样式模式

