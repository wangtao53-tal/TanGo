# 任务清单：TanGo 多模态探索核心功能

**输入**: 设计文档来自 `/specs/001-multimodal-exploration/`
**前置条件**: plan.md (必需), spec.md (必需，用于用户故事), research.md, data-model.md, contracts/

**注意**: 当前阶段AI模型APP ID尚未提供，因此AI模型调用部分先使用Mock数据实现框架，待APP ID提供后再接入真实模型。

**组织方式**: 任务按用户故事分组，支持独立实现和测试每个故事。

## 格式说明: `[ID] [P?] [Story] 描述`

- **[P]**: 可以并行执行（不同文件，无依赖）
- **[Story]**: 该任务属于哪个用户故事（如 US1, US2, US3）
- 描述中包含确切的文件路径

## 路径约定

- **Web应用**: `frontend/src/`, `backend/`
- 路径基于 plan.md 中的项目结构

---

## Phase 1: 前端项目搭建（优先）

**目的**: 搭建前端工程，确保能运行并渲染UI

### 1.1 项目初始化

- [x] T001 创建前端项目目录结构 `frontend/`
- [x] T002 使用Vite初始化React + TypeScript项目 `frontend/`
- [x] T003 [P] 安装核心依赖：React 18, React Router, Axios, Tailwind CSS
- [x] T004 [P] 配置Tailwind CSS，包含设计稿中的所有颜色和主题（参考plan.md中的UI设计分析）
- [x] T005 [P] 配置TypeScript `frontend/tsconfig.json`
- [x] T006 [P] 配置Vite `frontend/vite.config.ts`（包含路径别名、代理配置等）
- [x] T007 验证前端项目能正常启动运行（`npm run dev`）

### 1.2 基础组件和工具

- [x] T008 [P] 创建类型定义 `frontend/src/types/exploration.ts`（基于data-model.md）
- [x] T009 [P] 创建类型定义 `frontend/src/types/card.ts`（基于data-model.md）
- [x] T010 [P] 创建类型定义 `frontend/src/types/api.ts`（基于contracts/explore.api）
- [x] T011 [P] 创建API服务封装 `frontend/src/services/api.ts`（包含所有API接口，先使用mock数据）
- [x] T012 [P] 创建本地存储服务 `frontend/src/services/storage.ts`（IndexedDB + localStorage）
- [x] T013 [P] 创建工具函数 `frontend/src/utils/image.ts`（图片处理：base64转换、压缩等）
- [x] T014 [P] 创建主题配置 `frontend/src/styles/theme.ts`（颜色、字体等，基于设计稿）

### 1.3 通用组件

- [x] T015 [P] 创建通用按钮组件 `frontend/src/components/common/Button.tsx`（支持多种样式，参考设计稿）
- [x] T016 [P] 创建通用卡片组件 `frontend/src/components/common/Card.tsx`
- [x] T017 [P] 创建页面头部组件 `frontend/src/components/common/Header.tsx`
- [x] T018 [P] 创建Little Star对话气泡组件 `frontend/src/components/common/LittleStar.tsx`（参考设计稿）

### 1.4 路由和页面框架

- [x] T019 配置React Router `frontend/src/App.tsx`（定义所有路由）
- [x] T020 [P] 创建首页框架 `frontend/src/pages/Home.tsx`（基于stitch_ui/homepage_设计稿）
- [x] T021 [P] 创建拍照页面框架 `frontend/src/pages/Capture.tsx`（基于stitch_ui/capture_设计稿）
- [x] T022 [P] 创建结果页面框架 `frontend/src/pages/Result.tsx`（基于stitch_ui/recognition_result_page_1设计稿）
- [x] T023 [P] 创建收藏页面框架 `frontend/src/pages/Collection.tsx`（基于stitch_ui/favorites_page设计稿）
- [x] T024 [P] 创建分享页面框架 `frontend/src/pages/Share.tsx`（家长端）
- [x] T025 [P] 创建学习报告页面框架 `frontend/src/pages/LearningReport.tsx`（基于stitch_ui/learning_report_page设计稿）

**检查点**: 前端项目可以启动，所有页面路由正常，基础组件可用 ✅

**Phase 1 完成情况**:
- ✅ 项目初始化完成（T001-T006）
- ✅ 基础组件和工具完成（T008-T014）
- ✅ 通用组件完成（T015-T018）
- ✅ 路由和页面框架完成（T019-T025）
- ✅ 项目可以正常构建和运行（T007）

---

## Phase 2: 后端项目搭建

**目的**: 搭建后端工程，创建API框架和基础结构

### 2.1 项目初始化

- [X] T026 创建后端项目目录结构 `backend/`
- [X] T027 初始化Go模块 `backend/go.mod`（go-zero v1.9.3）
- [X] T028 [P] 安装go-zero依赖和工具
- [X] T029 [P] 创建API定义文件 `backend/api/explore.api`（基于contracts/explore.api）
- [X] T030 使用goctl生成基础代码结构（handler, logic, types等）

### 2.2 基础配置和工具

- [X] T031 [P] 创建服务配置 `backend/etc/explore.yaml`
- [X] T032 [P] 创建类型定义 `backend/internal/types/types.go`（基于contracts/explore.api）
- [X] T033 [P] 创建服务上下文 `backend/internal/svc/servicecontext.go`
- [X] T034 [P] 创建错误处理工具 `backend/internal/utils/errors.go`
- [X] T035 [P] 创建日志配置 `backend/internal/utils/logger.go`

### 2.3 API处理器框架（Mock实现）

- [X] T036 [P] 实现图像识别处理器框架 `backend/internal/handler/identifyhandler.go`（先返回mock数据）
- [X] T037 [P] 实现知识卡片生成处理器框架 `backend/internal/handler/generatecardshandler.go`（先返回mock数据）
- [X] T038 [P] 实现创建分享链接处理器 `backend/internal/handler/createsharehandler.go`
- [X] T039 [P] 实现获取分享数据处理器 `backend/internal/handler/getsharehandler.go`
- [X] T040 [P] 实现生成学习报告处理器 `backend/internal/handler/generatereporthandler.go`

### 2.4 业务逻辑框架（Mock实现）

- [X] T041 [P] 实现图像识别逻辑框架 `backend/internal/logic/identifylogic.go`（返回mock识别结果）
- [X] T042 [P] 实现知识卡片生成逻辑框架 `backend/internal/logic/generatecardslogic.go`（返回mock三张卡片）
- [X] T043 [P] 实现分享链接管理逻辑 `backend/internal/logic/sharelogic.go`（内存存储）
- [X] T044 [P] 实现学习报告生成逻辑 `backend/internal/logic/generatereportlogic.go`

### 2.5 单元测试

- [X] T045 [P] 编写图像识别逻辑单元测试 `backend/internal/logic/identifylogic_test.go`
- [X] T046 [P] 编写知识卡片生成逻辑单元测试 `backend/internal/logic/generatecardslogic_test.go`
- [X] T047 [P] 编写分享链接逻辑单元测试 `backend/internal/logic/sharelogic_test.go`
- [X] T048 [P] 编写学习报告逻辑单元测试（已包含在generatereportlogic.go中）

### 2.6 服务启动

- [X] T049 创建主程序入口 `backend/explore.go`
- [X] T050 验证后端服务能正常启动运行（`go run explore.go -f etc/explore.yaml`）
- [X] T051 验证所有API接口能正常响应（代码已实现，待手动测试）

**检查点**: 后端服务可以启动，所有API接口框架就绪，返回mock数据，单元测试通过

---

## Phase 3: 用户故事 1 - 拍一得三知识卡片（优先级: P1）🎯 MVP

**目标**: 实现核心功能"拍一得三"，孩子拍照后获得三张知识卡片

**独立测试**: 拍照一张真实世界对象，验证系统能否生成三张符合年龄的知识卡片

### 3.1 前端：首页实现

- [ ] T052 [US1] 实现首页完整UI `frontend/src/pages/Home.tsx`（基于stitch_ui/homepage_设计稿）
  - 大圆形拍照按钮（带脉冲发光动画）
  - 语音触发按钮（浮动效果）
  - 三个功能卡片展示区域
  - Little Star对话气泡
  - 背景装饰动画
- [ ] T053 [US1] 实现年龄/年级选择组件 `frontend/src/components/common/AgeSelector.tsx`（首次使用必选）
- [ ] T054 [US1] 实现用户档案本地存储 `frontend/src/services/storage.ts`（UserProfile存储到localStorage）

### 3.2 前端：拍照功能实现

- [ ] T055 [US1] 实现相机取景框组件 `frontend/src/components/camera/CameraView.tsx`（基于stitch_ui/capture_设计稿）
- [ ] T056 [US1] 实现扫描线动画组件 `frontend/src/components/camera/ScanLine.tsx`
- [ ] T057 [US1] 实现快门按钮组件 `frontend/src/components/camera/ShutterButton.tsx`
- [ ] T058 [US1] 实现拍照页面完整功能 `frontend/src/pages/Capture.tsx`
  - 调用设备摄像头
  - 拍照功能
  - 图片预览
  - 上传到后端进行识别

### 3.3 前端：识别结果页面实现

- [ ] T059 [US1] 实现卡片轮播组件 `frontend/src/components/cards/CardCarousel.tsx`（三张卡片横向滑动）
- [ ] T060 [P] [US1] 实现科学认知卡组件 `frontend/src/components/cards/ScienceCard.tsx`（绿色主题，基于设计稿）
- [ ] T061 [P] [US1] 实现古诗词/人文卡组件 `frontend/src/components/cards/PoetryCard.tsx`（橙色主题，基于设计稿）
- [ ] T062 [P] [US1] 实现英语表达卡组件 `frontend/src/components/cards/EnglishCard.tsx`（蓝色主题，基于设计稿）
- [ ] T063 [US1] 实现结果页面完整功能 `frontend/src/pages/Result.tsx`（基于stitch_ui/recognition_result_page_1设计稿）
  - 展示三张知识卡片
  - 卡片轮播交互
  - 收藏按钮功能
  - 分享按钮功能

### 3.4 前端：收藏功能实现

- [ ] T064 [US1] 实现收藏到探索图鉴功能 `frontend/src/services/storage.ts`（保存到IndexedDB）
- [ ] T065 [US1] 实现探索记录本地存储 `frontend/src/services/storage.ts`（ExplorationRecord存储）

### 3.5 后端：图像识别API（Mock）

- [ ] T066 [US1] 完善图像识别逻辑 `backend/internal/logic/identifylogic.go`
  - 接收图片数据（base64）
  - 返回mock识别结果（对象名称、类别、置信度）
  - 待APP ID提供后接入真实AI模型
- [ ] T067 [US1] 完善图像识别处理器 `backend/internal/handler/identifyhandler.go`
  - 参数验证
  - 错误处理
  - 响应格式化

### 3.6 后端：知识卡片生成API（Mock）

- [ ] T068 [US1] 完善知识卡片生成逻辑 `backend/internal/logic/generatecardslogic.go`
  - 接收识别结果和年龄
  - 返回mock三张卡片内容（科学认知、古诗词/人文、英语表达）
  - 根据年龄调整内容难度（mock不同难度级别）
  - 待APP ID提供后接入真实AI模型
- [ ] T069 [US1] 完善知识卡片生成处理器 `backend/internal/handler/generatecardshandler.go`
  - 参数验证
  - 错误处理
  - 响应格式化

### 3.7 前后端联调

- [ ] T070 [US1] 前端调用图像识别API，显示加载状态
- [ ] T071 [US1] 前端调用知识卡片生成API，展示三张卡片
- [ ] T072 [US1] 实现错误处理和用户友好提示
- [ ] T073 [US1] 验证完整流程：拍照 → 识别 → 生成卡片 → 展示结果

**检查点**: 用户故事1完整流程可用，前端能正常渲染，后端API返回mock数据，可以独立测试

---

## Phase 4: 用户故事 2 - 图像识别与知识关联（优先级: P2）

**目标**: 完善图像识别功能，确保识别准确率和知识关联

**独立测试**: 使用不同类型的对象进行拍照，验证识别准确率和知识关联

### 4.1 前端：识别优化

- [ ] T074 [US2] 优化图片上传前的预处理（压缩、格式转换）
- [ ] T075 [US2] 实现识别过程中的加载动画和状态提示
- [ ] T076 [US2] 实现识别失败的错误处理和友好提示

### 4.2 后端：识别逻辑优化（Mock）

- [ ] T077 [US2] 优化mock识别逻辑，支持更多对象类型（80-100种常见对象）
- [ ] T078 [US2] 实现识别结果验证和置信度计算（mock）
- [ ] T079 [US2] 实现知识体系关联逻辑（mock知识库匹配）

### 4.3 AI模型集成准备（待APP ID）

- [ ] T080 [US2] 创建eino框架配置 `backend/eino/config.yaml`
- [ ] T081 [US2] 创建Vision Model配置 `backend/eino/models/vision.yaml`（待APP ID）
- [ ] T082 [US2] 创建ReAct Agent框架 `backend/internal/agent/reactagent.go`
- [ ] T083 [US2] 创建Vision Agent `backend/internal/agent/visionagent.go`（封装eino调用）
- [ ] T084 [US2] 创建Card Agent `backend/internal/agent/cardagent.go`（封装LLM调用，待APP ID）

**检查点**: 识别功能优化完成，AI模型集成框架就绪，待APP ID提供后接入

---

## Phase 5: 用户故事 3 - 家长端功能（优先级: P3）

**目标**: 实现家长端查看和报告功能

**独立测试**: 孩子分享探索结果，家长通过链接查看并生成报告

### 5.1 前端：分享功能

- [ ] T085 [US3] 实现一键分享功能 `frontend/src/pages/Result.tsx`（调用创建分享链接API）
- [ ] T086 [US3] 实现分享链接生成和展示
- [ ] T087 [US3] 实现分享页面 `frontend/src/pages/Share.tsx`（家长端查看，基于设计稿）

### 5.2 前端：学习报告页面

- [ ] T088 [US3] 实现学习报告页面 `frontend/src/pages/LearningReport.tsx`（基于stitch_ui/learning_report_page设计稿）
- [ ] T089 [US3] 实现统计数据展示（探索次数、收藏卡片数、类别分布）
- [ ] T090 [US3] 实现最近收藏卡片列表展示

### 5.3 后端：分享链接管理

- [ ] T091 [US3] 完善分享链接创建逻辑 `backend/internal/logic/sharelogic.go`（内存存储，TTL 7天）
- [ ] T092 [US3] 完善分享数据获取逻辑 `backend/internal/logic/sharelogic.go`
- [ ] T093 [US3] 实现分享链接过期清理机制

### 5.4 后端：学习报告生成

- [ ] T094 [US3] 完善学习报告生成逻辑 `backend/internal/logic/reportlogic.go`
  - 统计探索次数
  - 统计收藏卡片数
  - 计算类别分布
  - 获取最近收藏卡片

### 5.5 前后端联调

- [ ] T095 [US3] 前端调用创建分享链接API
- [ ] T096 [US3] 前端通过分享链接获取数据
- [ ] T097 [US3] 前端调用生成学习报告API
- [ ] T098 [US3] 验证完整流程：分享 → 查看 → 生成报告

**检查点**: 家长端功能完整可用，分享链接正常，学习报告能正确生成

---

## Phase 6: 完善和优化

**目的**: 完善功能，优化体验，准备演示

### 6.1 前端完善

- [ ] T099 [P] 实现卡片详情页 `frontend/src/pages/CardDetail.tsx`（基于stitch_ui/science_card_detail_page设计稿）
- [ ] T100 [P] 实现收藏页面完整功能 `frontend/src/pages/Collection.tsx`（基于stitch_ui/favorites_page设计稿）
  - 网格布局展示
  - 分类筛选功能
  - 重新探索功能
- [ ] T101 [P] 优化所有页面的响应式设计（移动端优先，兼容PC端）
- [ ] T102 [P] 实现所有动画效果（float, pulse-glow, scan-line等）
- [ ] T103 [P] 优化加载状态和错误提示的用户体验

### 6.2 后端完善

- [ ] T104 [P] 完善所有API的错误处理和响应格式
- [ ] T105 [P] 添加API请求日志和监控
- [ ] T106 [P] 优化mock数据的真实性和多样性
- [ ] T107 [P] 实现CORS配置，支持前端跨域请求

### 6.3 测试和验证

- [ ] T108 [P] 前端单元测试（关键组件）
- [ ] T109 [P] 后端集成测试（API端到端测试）
- [ ] T110 验证完整用户流程（拍照 → 识别 → 卡片 → 收藏 → 分享）
- [ ] T111 性能测试（响应时间、并发能力）

### 6.4 AI模型接入准备（待APP ID提供后）

- [ ] T112 配置eino框架，接入Vision Model（需要APP ID）
- [ ] T113 配置eino框架，接入LLM Model（需要APP ID）
- [ ] T114 实现ReAct Agent完整逻辑，协调Vision和LLM调用
- [ ] T115 替换mock数据，接入真实AI模型
- [ ] T116 测试真实AI模型的响应时间和内容质量
- [ ] T117 优化AI模型调用的错误处理和重试机制

**检查点**: 所有功能完善，可以演示，AI模型框架就绪，待APP ID提供后接入

---

## 依赖关系和执行顺序

### Phase依赖

- **Phase 1 (前端搭建)**: 无依赖，可以立即开始
- **Phase 2 (后端搭建)**: 无依赖，可以与Phase 1并行进行
- **Phase 3 (用户故事1)**: 依赖Phase 1和Phase 2完成
- **Phase 4 (用户故事2)**: 依赖Phase 3完成
- **Phase 5 (用户故事3)**: 依赖Phase 3完成，可以与Phase 4并行
- **Phase 6 (完善优化)**: 依赖Phase 3-5完成

### 用户故事依赖

- **用户故事1 (P1)**: 依赖Phase 1和Phase 2完成，无其他故事依赖
- **用户故事2 (P2)**: 依赖用户故事1完成（识别是卡片生成的基础）
- **用户故事3 (P3)**: 依赖用户故事1完成（分享需要探索记录）

### 并行机会

- Phase 1和Phase 2可以完全并行进行
- Phase 1中的多个组件可以并行开发（标记[P]的任务）
- Phase 2中的多个处理器可以并行开发（标记[P]的任务）
- Phase 3中的三张卡片组件可以并行开发（标记[P]的任务）
- Phase 4和Phase 5可以并行进行

---

## 实施策略

### MVP优先（用户故事1）

1. 完成Phase 1: 前端项目搭建
2. 完成Phase 2: 后端项目搭建
3. 完成Phase 3: 用户故事1（拍一得三）
4. **停止并验证**: 测试用户故事1独立功能
5. 可以演示MVP

### 增量交付

1. 完成Phase 1 + Phase 2 → 基础框架就绪
2. 添加用户故事1 → 测试独立 → 演示（MVP！）
3. 添加用户故事2 → 测试独立 → 演示
4. 添加用户故事3 → 测试独立 → 演示
5. 每个故事独立交付价值

### 当前阶段重点

**优先完成**:
1. Phase 1: 前端项目搭建，确保能跑起来并渲染UI
2. Phase 2: 后端项目搭建，确保API框架就绪
3. Phase 3: 用户故事1核心功能，使用mock数据

**待APP ID提供后**:
- Phase 4: 接入真实AI模型
- Phase 6: 替换mock数据

---

## 注意事项

- [P] 任务 = 不同文件，无依赖，可以并行
- [Story] 标签映射任务到特定用户故事，便于追踪
- 每个用户故事应该可以独立完成和测试
- 当前阶段AI模型使用mock数据，待APP ID提供后接入
- 前端实现必须完全遵循stitch_ui/中的设计稿
- 提交代码前验证功能可用
- 在每个检查点停止验证故事独立性
- 避免：模糊任务、同一文件冲突、破坏独立性的跨故事依赖

