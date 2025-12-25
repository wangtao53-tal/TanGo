# 配置文件变更说明

## 📋 概述

本次多Agent功能实现新增了一个前端环境变量配置项。

## 🆕 新增配置项

### 前端配置

#### `VITE_USE_MULTI_AGENT`

- **类型**: `boolean`
- **默认值**: `false`
- **说明**: 控制前端是否使用多Agent模式
  - `true`: 使用多Agent模式（调用 `/api/conversation/agent` 接口）
  - `false`: 使用单Agent模式（调用 `/api/conversation/stream` 接口，默认值）
- **优先级**: localStorage > 环境变量 > 默认值
- **运行时切换**: 支持在Settings页面运行时切换，无需重启应用

**配置位置**: `.env` 文件（前端配置部分）

**示例**:
```bash
# 启用多Agent模式
VITE_USE_MULTI_AGENT=true

# 使用单Agent模式（默认）
VITE_USE_MULTI_AGENT=false
```

## 📝 配置文件更新

### `.env.example` 文件

已在根目录的 `.env.example` 文件中添加了新的配置项：

```bash
# ==================== 前端服务配置 ====================
# ... 其他前端配置 ...
# 是否使用多Agent模式（true=多Agent模式，false=单Agent模式，默认false）
# 注意：此配置也可以通过前端Settings页面运行时切换（优先级更高）
VITE_USE_MULTI_AGENT=false
```

## 🔄 配置优先级

前端API配置的优先级顺序：

1. **localStorage**（最高优先级）
   - 通过Settings页面设置
   - 键名: `useMultiAgent`
   - 值: `"true"` 或 `"false"`

2. **环境变量**
   - `VITE_USE_MULTI_AGENT`
   - 在 `.env` 文件中配置

3. **默认值**（最低优先级）
   - `false`（使用单Agent模式，向后兼容）

## 📍 相关文件

- `frontend/src/config/api.ts` - API配置模块实现
- `frontend/src/pages/Settings.tsx` - Settings页面切换功能
- `.env.example` - 环境变量配置示例

## ✅ 更新检查清单

- [x] 在 `.env.example` 中添加 `VITE_USE_MULTI_AGENT` 配置项
- [x] 添加配置说明注释
- [x] 设置默认值为 `false`（向后兼容）

## 🚀 使用说明

### 方式1：环境变量配置（推荐用于生产环境）

在 `.env` 文件中配置：
```bash
VITE_USE_MULTI_AGENT=true
```

### 方式2：运行时切换（推荐用于开发测试）

在Settings页面切换"单Agent模式"或"多Agent模式"，配置会保存到localStorage。

### 方式3：使用默认值

不配置任何值，系统默认使用单Agent模式（`false`），确保向后兼容。

## 📚 相关文档

- `specs/014-multi-agent-followup/ENV_SETUP_EXAMPLE.md` - 环境变量配置示例
- `specs/014-multi-agent-followup/MODEL_SELECTION_GUIDE.md` - 模型选择机制指南

