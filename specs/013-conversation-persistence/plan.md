# Implementation Plan: 对话页面会话记录持久化

**Branch**: `dev-mvp-20251218` | **Date**: 2025-12-20 | **Spec**: [spec.md](./spec.md)

**Note**: MVP版本阶段，所有开发工作统一在 `dev-mvp-20251218` 分支进行，不采用一个功能一个分支的策略。

## Summary

实现对话页面会话记录的完整持久化功能，确保用户在刷新页面、切换页面或重新打开应用后能够恢复之前的对话记录，只有在用户主动从拍照页面（Capture页面）选择图片并识别后跳转到对话页面时才创建新会话。

**重要区分**：
- **拍照上传**：从Capture页面选择图片并识别，通过`location.state`传递识别结果跳转到Result页面 → **创建新会话**
- **上传图片**：在对话页面中通过ImageInput组件上传图片 → **继续当前会话**，不创建新会话

核心功能：
1. **刷新页面后恢复对话记录**：自动恢复所有消息、卡片和识别结果上下文，**不会重新生成卡片**（如果已有卡片消息）
2. **切换页面后保持对话记录**：在不同页面间切换后返回时保持对话连续性
3. **重新拍照上传时创建新会话**：只有从拍照页面（Capture页面）选择图片并识别后跳转到对话页面时才创建新会话并清空之前的对话记录

技术方案：优化现有的持久化逻辑，确保所有消息和上下文信息正确保存到IndexedDB和localStorage，改进恢复逻辑以处理所有场景，包括流式消息状态的处理。确保刷新页面时判断"是否需要生成卡片"（即是否已经有卡片了），而不是判断"是否有需要识别就重新生成卡片"。

## Technical Context

**Language/Version**: TypeScript 5.x, React 18.x  
**Primary Dependencies**: React, React Router, IndexedDB (通过现有的storage服务)  
**Storage**: IndexedDB (conversationStorage), localStorage (sessionId和identificationContext)  
**Testing**: Jest + React Testing Library (前端单元测试)  
**Target Platform**: Web (移动端优先，支持PC端)  
**Project Type**: Web application (frontend)  
**Performance Goals**: 恢复对话记录时间≤2秒，切换页面恢复时间≤1秒，消息保存成功率≥99%  
**Constraints**: 必须兼容现有的对话功能，不影响流式消息和卡片生成功能，支持优雅降级（IndexedDB失败时使用localStorage）  
**Scale/Scope**: 单页面优化（Result.tsx），涉及存储服务（storage.ts）和对话消息处理逻辑

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**规范检查项**（基于项目规范）：

- [x] **原则一：中文优先规范** - 所有文档和生成内容必须使用中文（除非技术限制）
- [x] **原则二：K12 教育游戏化设计规范** - 持久化功能不影响现有的游戏化元素和儿童友好设计
- [x] **原则三：可发布应用规范** - 实现必须达到生产级标准，确保数据持久化的可靠性和性能
- [x] **原则四：多语言和年级设置规范** - 持久化功能不影响多语言和年级设置
- [x] **原则五：AI优先（模型优先）规范** - 持久化功能不影响AI模型调用和流式输出
- [x] **原则六：移动端优先规范** - 持久化功能在移动端和PC端都能正常工作
- [x] **原则七：用户体验流程规范** - 持久化功能确保用户体验的连续性，不影响识别和对话流程
- [x] **原则八：对话Agent技术规范** - 持久化功能不影响Agent的流式输出和Markdown格式支持

**合规性说明**：所有实现均符合项目规范要求，无违反项。持久化功能是基础功能，确保用户体验的连续性，符合所有规范要求。

## Project Structure

### Documentation (this feature)

```text
specs/013-conversation-persistence/
├── plan.md              # This file (/speckit.plan command output)
├── spec.md              # Feature specification
├── checklists/
│   └── requirements.md  # Quality checklist
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
frontend/
├── src/
│   ├── pages/
│   │   └── Result.tsx                    # 对话页面（需要优化：改进恢复逻辑、确保所有消息保存）
│   ├── services/
│   │   └── storage.ts                     # 存储服务（已存在，可能需要优化错误处理）
│   └── types/
│       ├── conversation.ts                # 对话类型定义（已存在，可能需要扩展）
│       └── api.ts                         # API类型定义（已存在，包含IdentificationContext）
```

**Structure Decision**: 主要修改Result.tsx中的持久化和恢复逻辑，确保所有消息和上下文信息正确保存和恢复。利用现有的storage.ts服务，优化错误处理和降级策略。

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

无违反项，无需填写。

## Phase 0: Research & Analysis

### 现有实现分析

1. **持久化现状**：
   - ✅ 已有IndexedDB存储对话消息（conversationStorage.saveMessage）
   - ✅ 已有localStorage存储sessionId和identificationContext
   - ✅ 已有恢复对话记录的逻辑（restoreConversation）
   - ✅ 已有检测新识别结果的逻辑（通过location.state）

2. **问题分析**：
   - ⚠️ **问题1**：流式消息保存时可能包含`isStreaming: true`状态，恢复时应该标记为已完成
   - ⚠️ **问题2**：从Capture页面传入新识别结果时，虽然创建了新会话，但可能没有正确清空旧会话的IndexedDB数据（只是清空了内存状态）
   - ⚠️ **问题3**：切换页面后返回时，如果location.state为空，会尝试恢复会话，但可能在某些边界情况下没有正确恢复
   - ⚠️ **问题4**：错误处理不够完善，如果IndexedDB操作失败，没有优雅降级到localStorage
   - ⚠️ **问题5**：恢复对话记录时，没有正确处理流式消息的`streamingText`字段（应该清空）
   - ⚠️ **问题6**：需要明确区分"拍照上传"（从Capture页面）和"上传图片"（在对话页面中），只有拍照上传才创建新会话
   - ⚠️ **问题7**：刷新页面时，应该判断"是否需要生成卡片"（即是否已经有卡片了），而不是判断"是否有需要识别就重新生成卡片"

3. **代码审查发现**：
   - `Result.tsx`第98-115行：`restoreConversation`函数已实现，但需要确保正确处理流式消息状态
   - `Result.tsx`第127-211行：创建新会话的逻辑已实现，但需要确保清空旧会话数据
   - `Result.tsx`第655-673行：流式消息完成时，需要确保保存的消息不包含`isStreaming: true`
   - `storage.ts`第353-504行：conversationStorage服务已实现，但需要优化错误处理

### 技术调研

1. **IndexedDB最佳实践**：
   - 使用事务确保数据一致性
   - 处理存储配额错误（QuotaExceededError）
   - 实现优雅降级到localStorage

2. **React Router状态管理**：
   - `location.state`在页面刷新后会丢失，必须依赖持久化存储
   - 使用`useEffect`监听路由变化，但需要正确处理依赖项

3. **流式消息状态处理**：
   - 流式消息完成时，必须清除`isStreaming`和`streamingText`字段
   - 恢复时，所有消息的`isStreaming`应该为`false`或`undefined`

### 依赖关系

- **依赖现有功能**：
  - conversationStorage服务（已存在）
  - localStorage API（浏览器原生）
  - IndexedDB API（浏览器原生）
  - React Router的location.state（用于检测新识别结果）

- **不影响的功能**：
  - 流式消息生成和显示
  - 卡片生成和展示
  - 语音输入和图片上传
  - 对话Agent的流式输出

## Phase 1: Design

### 架构设计

1. **持久化策略**：
   - **消息存储**：所有消息（包括文本、语音、图片、卡片）保存到IndexedDB
   - **上下文存储**：识别结果上下文（IdentificationContext）保存到localStorage
   - **会话ID存储**：当前会话ID保存到localStorage（key: `currentSessionId`）
   - **降级策略**：如果IndexedDB失败，至少保存关键信息到localStorage

2. **恢复策略**：
   - **检测新会话**：通过`location.state`检测是否有新识别结果（`location.state`包含`objectName`且为有效字符串）
   - **恢复旧会话**：如果没有新识别结果（刷新页面或直接访问），从localStorage读取sessionId，从IndexedDB恢复消息
   - **状态恢复**：恢复识别结果上下文、消息列表、会话ID，确保所有流式消息标记为已完成
   - **卡片生成判断**：恢复时检查是否已有卡片消息，如果已有卡片，标记为已生成，**不会重新生成卡片**。判断逻辑：判断"是否需要生成卡片"（即是否已经有卡片了），而不是判断"是否有需要识别就重新生成卡片"

3. **新会话创建策略**：
   - **触发条件**：从Capture页面选择图片并识别后跳转到Result页面（`location.state`包含`objectName`且为有效字符串）
   - **清理操作**：清空当前消息列表（内存），生成新sessionId，保存新识别结果上下文
   - **数据隔离**：新会话使用新的sessionId，旧会话数据保留在IndexedDB中（不删除，支持历史记录）
   - **重要区分**：
     - **拍照上传**（从Capture页面）：创建新会话，清空消息列表，生成新sessionId
     - **上传图片**（在对话页面中通过ImageInput组件）：继续当前会话，不清空消息列表，使用当前sessionId

### 数据流设计

1. **消息保存流程**：
   ```
   用户发送消息 → 创建消息对象 → 保存到IndexedDB → 更新UI
   流式消息完成 → 清除isStreaming字段 → 保存到IndexedDB → 更新UI
   ```

2. **会话恢复流程**：
   ```
   页面加载 → 检查location.state → 
   有新识别结果（location.state包含objectName且为有效字符串）？ 
   → 是：创建新会话（从Capture页面跳转）
   → 否：从localStorage读取sessionId → 从IndexedDB恢复消息 → 恢复识别结果上下文 → 
        检查是否已有卡片消息 → 如果有，标记为已生成，不重新生成卡片 → 更新UI
   ```

3. **新会话创建流程**：
   ```
   从Capture页面跳转 → location.state包含识别结果（objectName且为有效字符串） → 
   清空当前消息列表 → 生成新sessionId → 保存新识别结果上下文 → 
   创建初始消息（图片+识别结果） → 保存到IndexedDB → 自动生成卡片
   ```

4. **上传图片流程**（在对话页面中）：
   ```
   用户在对话页面中通过ImageInput组件上传图片 → 
   使用当前sessionId → 添加图片消息到当前会话 → 
   保存到IndexedDB → 发送流式请求（不传递identificationContext，不生成卡片） → 
   更新UI（继续当前会话，不清空消息列表）
   ```

### 错误处理设计

1. **IndexedDB错误处理**：
   - 捕获`QuotaExceededError`，提示用户清理存储空间
   - 捕获其他IndexedDB错误，降级到localStorage保存关键信息
   - 记录错误日志，但不中断用户操作

2. **localStorage错误处理**：
   - 捕获`QuotaExceededError`，提示用户清理存储空间
   - 捕获其他localStorage错误，记录日志，但不中断用户操作

3. **恢复失败处理**：
   - 如果恢复失败，显示友好提示，引导用户重新开始对话
   - 如果部分数据恢复失败，尽可能恢复可用数据

### 性能优化

1. **恢复性能**：
   - 使用异步加载，不阻塞UI渲染
   - 按时间顺序恢复消息，确保顺序正确
   - 限制恢复的消息数量（如果需要，最多恢复最近1000条）

2. **保存性能**：
   - 使用批量保存（如果需要）
   - 异步保存，不阻塞用户操作
   - 错误处理不阻塞主流程

## Phase 2: Implementation

> **Note**: Detailed implementation tasks will be created in `tasks.md` by `/speckit.tasks` command.

### 实现概览

1. **优化消息保存逻辑**：
   - 确保所有消息（包括流式消息完成时）正确保存到IndexedDB
   - 流式消息完成时，清除`isStreaming`和`streamingText`字段再保存
   - 添加错误处理和降级策略

2. **优化恢复逻辑**：
   - 改进`restoreConversation`函数，确保正确处理所有消息类型
   - 确保恢复时所有流式消息标记为已完成
   - 添加恢复失败的错误处理

3. **优化新会话创建逻辑**：
   - 确保从Capture页面传入新识别结果时（`location.state`包含`objectName`且为有效字符串），正确创建新会话
   - 确保新会话ID和识别结果上下文正确保存
   - 确保旧会话数据不干扰新会话（通过sessionId隔离）
   - **重要**：确保在对话页面中通过ImageInput组件上传图片时，继续当前会话，不创建新会话，不清空消息列表

4. **优化错误处理**：
   - 添加IndexedDB错误处理和降级策略
   - 添加localStorage错误处理
   - 添加用户友好的错误提示

5. **测试和验证**：
   - 测试刷新页面后恢复对话记录
   - 测试切换页面后保持对话记录
   - 测试重新拍照上传时创建新会话
   - 测试错误处理和降级策略

### 关键实现点

1. **流式消息状态处理**：
   - 流式消息完成时（`onClose`回调），必须清除`isStreaming`和`streamingText`字段
   - 保存到IndexedDB的消息不应该包含`isStreaming: true`

2. **拍照上传 vs 上传图片的处理**：
   - **拍照上传**（从Capture页面）：检测`location.state`包含`objectName`且为有效字符串 → 创建新会话
   - **上传图片**（在对话页面中）：使用当前`sessionId`，不清空消息列表，继续当前会话

3. **会话隔离**：
   - 使用sessionId区分不同会话
   - 新会话创建时（从Capture页面跳转），使用新的sessionId，旧会话数据保留在IndexedDB中
   - 在对话页面中上传图片时，使用当前sessionId，不创建新会话

4. **恢复时机和卡片生成判断**：
   - 页面加载时（`useEffect`）检查是否需要恢复
   - 切换页面后返回时，通过React Router的location变化触发恢复
   - **重要**：恢复时检查是否已有卡片消息，如果已有卡片，标记为已生成，**不会重新生成卡片**
   - 判断逻辑：判断"是否需要生成卡片"（即是否已经有卡片了），而不是判断"是否有需要识别就重新生成卡片"

4. **数据完整性**：
   - 确保所有消息按时间顺序恢复
   - 确保识别结果上下文完整恢复
   - 确保会话ID和消息的sessionId一致

## Risk Assessment

### 技术风险

1. **IndexedDB兼容性**：
   - **风险**：某些旧浏览器可能不支持IndexedDB
   - **缓解**：使用现有的storage服务，已有兼容性处理

2. **存储配额限制**：
   - **风险**：用户设备存储空间不足
   - **缓解**：实现优雅降级，优先保存关键信息，提示用户清理空间

3. **数据一致性**：
   - **风险**：恢复时数据不完整或不一致
   - **缓解**：添加数据验证，确保恢复的数据完整且一致

### 业务风险

1. **用户体验**：
   - **风险**：恢复时间过长影响用户体验
   - **缓解**：优化恢复逻辑，使用异步加载，目标恢复时间≤2秒

2. **数据丢失**：
   - **风险**：某些情况下数据可能丢失
   - **缓解**：实现多重备份（IndexedDB + localStorage），添加错误处理和用户提示

## Success Metrics

- **SC-001**: 刷新页面后恢复对话记录时间≤2秒，恢复成功率100%
- **SC-002**: 切换页面后恢复对话记录时间≤1秒，恢复成功率100%
- **SC-003**: 新会话创建成功率100%
- **SC-004**: 消息保存成功率≥99%
- **SC-005**: 识别结果上下文保存成功率100%
- **SC-006**: 对话连续性保持率100%
- **SC-007**: 数据完整性100%，无消息丢失或顺序错误

## Next Steps

1. 运行 `/speckit.tasks` 创建详细的实现任务列表
2. 开始实现Phase 2中的任务
3. 测试和验证所有功能
4. 确保符合所有成功标准

