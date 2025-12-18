# 数据模型设计

**创建日期**: 2025-12-18  
**功能**: TanGo 多模态探索核心功能

## 实体定义

### 1. UserProfile（用户档案）

**存储位置**: 前端 localStorage

**用途**: 存储孩子的基本信息，用于知识卡片内容分级

**字段**:
- `age` (number, required): 孩子年龄，范围 3-18
- `grade` (string, optional): 年级，格式如 "K1", "K2", "G1", "G2" 等
- `lastUpdated` (string, ISO 8601): 最后更新时间

**验证规则**:
- age 必须在 3-18 之间
- grade 格式必须符合 K1-K12 规范

**示例**:
```json
{
  "age": 8,
  "grade": "G3",
  "lastUpdated": "2025-12-18T10:00:00Z"
}
```

### 2. ExplorationRecord（探索记录）

**存储位置**: 前端 IndexedDB，分享时上传到后端

**用途**: 记录每次探索的详细信息

**字段**:
- `id` (string, UUID): 唯一标识符
- `timestamp` (string, ISO 8601): 探索时间
- `objectName` (string): 识别出的对象名称（中文）
- `objectCategory` (string): 对象类别（"自然类" | "生活类" | "人文类"）
- `confidence` (number, 0-1): 识别置信度
- `age` (number): 探索时的年龄（用于内容分级）
- `imageData` (string, base64, optional): 原始图片数据（仅本地存储，不分享）
- `cards` (Card[]): 生成的三张知识卡片
- `collected` (boolean): 是否已收藏

**索引**:
- 主键：id
- 索引：timestamp（用于按时间排序）
- 索引：objectCategory（用于分类筛选）

**示例**:
```json
{
  "id": "exp-123e4567-e89b-12d3-a456-426614174000",
  "timestamp": "2025-12-18T10:00:00Z",
  "objectName": "银杏",
  "objectCategory": "自然类",
  "confidence": 0.95,
  "age": 8,
  "imageData": "data:image/jpeg;base64,...",
  "cards": [
    {
      "type": "science",
      "title": "银杏",
      "content": {...}
    },
    ...
  ],
  "collected": true
}
```

### 3. KnowledgeCard（知识卡片）

**存储位置**: 前端 IndexedDB（作为探索记录的一部分），分享时上传到后端

**用途**: 存储生成的三张知识卡片内容

**字段**:
- `id` (string, UUID): 唯一标识符
- `explorationId` (string): 关联的探索记录ID
- `type` (string): 卡片类型（"science" | "poetry" | "english"）
- `title` (string): 卡片标题
- `content` (object): 卡片内容（根据类型不同结构不同）
- `collectedAt` (string, ISO 8601, optional): 收藏时间

**卡片内容结构**:

**科学认知卡 (science)**:
```json
{
  "type": "science",
  "title": "银杏",
  "content": {
    "name": "银杏",
    "explanation": "银杏是非常古老的植物，已经在地球上生存了2亿多年。",
    "facts": [
      "银杏是现存最古老的树种之一",
      "银杏的叶子在秋天会变成金黄色",
      "银杏的果实可以食用，但需要处理"
    ],
    "funFact": "银杏被称为'活化石'，因为它在恐龙时代就已经存在了！"
  }
}
```

**古诗词/人文卡 (poetry)**:
```json
{
  "type": "poetry",
  "title": "古人怎么看银杏",
  "content": {
    "poem": "满地翻黄银杏叶，忽惊天地告成功。",
    "poemSource": "《夜坐》- 李清照",
    "explanation": "这句诗描写了秋天银杏叶变黄的美丽景象。",
    "context": "看到银杏，我们可以联想到秋天的美丽，以及时间的流逝。"
  }
}
```

**英语表达卡 (english)**:
```json
{
  "type": "english",
  "title": "用英语说银杏",
  "content": {
    "keywords": ["ginkgo", "tree", "ancient"],
    "expressions": [
      "This is a ginkgo tree.",
      "The ginkgo leaves are golden in autumn."
    ],
    "pronunciation": "ginkgo: /ˈɡɪŋkoʊ/"
  }
}
```

**索引**:
- 主键：id
- 索引：explorationId（用于关联探索记录）
- 索引：type（用于按类型筛选）
- 索引：collectedAt（用于按收藏时间排序）

### 4. ShareLink（分享链接）

**存储位置**: 后端内存/Redis（临时存储）

**用途**: 存储分享给家长的数据

**字段**:
- `shareId` (string, UUID): 分享链接唯一标识符
- `explorationRecords` (ExplorationRecord[]): 探索记录列表
- `collectedCards` (KnowledgeCard[]): 收藏的卡片列表
- `createdAt` (string, ISO 8601): 创建时间
- `expiresAt` (string, ISO 8601): 过期时间（默认7天后）

**TTL**: 7天（604800秒）

**示例**:
```json
{
  "shareId": "share-123e4567-e89b-12d3-a456-426614174000",
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

### 5. LearningReport（学习报告）

**存储位置**: 后端临时生成，不持久化

**用途**: 家长端一键生成的学习报告

**字段**:
- `totalExplorations` (number): 总探索次数
- `totalCollectedCards` (number): 总收藏卡片数
- `categoryDistribution` (object): 类别分布
  - `natural` (number): 自然类探索次数
  - `life` (number): 生活类探索次数
  - `humanity` (number): 人文类探索次数
- `recentCards` (KnowledgeCard[]): 最近收藏的卡片（最多10张）
- `generatedAt` (string, ISO 8601): 生成时间

**示例**:
```json
{
  "totalExplorations": 15,
  "totalCollectedCards": 45,
  "categoryDistribution": {
    "natural": 8,
    "life": 5,
    "humanity": 2
  },
  "recentCards": [...],
  "generatedAt": "2025-12-18T10:00:00Z"
}
```

## 数据关系

```
UserProfile (1)
    ↓
ExplorationRecord (N) - 每次探索时使用UserProfile中的age
    ↓
KnowledgeCard (3 per ExplorationRecord) - 每张探索记录生成3张卡片
    ↓
ShareLink (1) - 包含多个ExplorationRecord和KnowledgeCard
    ↓
LearningReport (1) - 从ShareLink数据生成
```

## 数据验证规则

### UserProfile
- age: 必须为整数，范围 3-18
- grade: 必须符合格式 "K1" | "K2" | "G1" | "G2" | ... | "G12"
- lastUpdated: 必须为有效的ISO 8601格式

### ExplorationRecord
- id: 必须为有效的UUID v4
- timestamp: 必须为有效的ISO 8601格式
- objectName: 不能为空，最大长度100字符
- objectCategory: 必须是 "自然类" | "生活类" | "人文类" 之一
- confidence: 必须在 0-1 之间
- age: 必须在 3-18 之间
- cards: 必须是包含3个元素的数组

### KnowledgeCard
- id: 必须为有效的UUID v4
- explorationId: 必须关联到存在的ExplorationRecord
- type: 必须是 "science" | "poetry" | "english" 之一
- title: 不能为空，最大长度200字符
- content: 必须符合对应类型的结构

### ShareLink
- shareId: 必须为有效的UUID v4
- createdAt: 必须为有效的ISO 8601格式
- expiresAt: 必须晚于createdAt

## 数据迁移策略

**版本管理**:
- 使用版本号标记数据结构（如 v1.0）
- 在UserProfile中存储当前数据版本

**迁移机制**:
- 检测到版本不匹配时，执行迁移脚本
- 迁移失败时，提供数据重置选项

## 数据清理策略

**前端清理**:
- 探索记录：保留最近100条，超出部分提示用户清理
- 图片数据：仅保留最近30天的原始图片，超出部分删除imageData字段
- 收藏卡片：永久保留，除非用户手动删除

**后端清理**:
- ShareLink：TTL到期自动删除
- 不存储用户原始图片数据（隐私保护）

## 隐私保护

1. **本地存储**: 所有用户数据仅存储在浏览器本地，不上传到服务器（除非用户主动分享）
2. **分享数据**: 分享时仅上传探索记录和卡片内容，不包含原始图片
3. **数据加密**: 分享链接使用UUID v4，不可预测
4. **自动过期**: 分享链接7天后自动过期
5. **无账户系统**: 完全匿名使用，不收集用户身份信息

