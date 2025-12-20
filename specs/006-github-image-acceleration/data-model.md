# 数据模型：GitHub 图片加速

**功能**: 006-github-image-acceleration  
**创建日期**: 2025-12-20

## 概述

本文档描述 GitHub 图片加速功能涉及的数据结构和配置模型。

## 配置模型

### CDNConfig

CDN 配置结构，用于控制 CDN 加速行为。

```go
type CDNConfig struct {
    // 是否启用CDN加速
    EnableGitHubCDN bool `json:"enableGitHubCDN" default:"true"`
    
    // CDN提供商：jsdelivr, fastgit
    CDNProvider string `json:"cdnProvider" default:"jsdelivr"`
    
    // 是否启用CDN失败重试
    CDNRetryEnabled bool `json:"cdnRetryEnabled" default:"true"`
    
    // CDN请求超时（秒）
    CDNTimeout int `json:"cdnTimeout" default:"10"`
    
    // 原始URL重试超时（秒）
    OriginalURLTimeout int `json:"originalURLTimeout" default:"15"`
}
```

## URL 转换规则

### GitHub Raw URL 格式

```
https://raw.githubusercontent.com/{owner}/{repo}/{branch}/{path}
```

### jsDelivr CDN URL 格式

```
https://cdn.jsdelivr.net/gh/{owner}/{repo}@{branch}/{path}
```

### FastGit URL 格式

```
https://raw.fastgit.org/{owner}/{repo}/{branch}/{path}
```

## URL 解析结构

### GitHubURLParts

GitHub URL 解析结果结构。

```go
type GitHubURLParts struct {
    Owner  string // 仓库所有者
    Repo   string // 仓库名
    Branch string // 分支名
    Path   string // 文件路径
    IsValid bool  // 是否为有效的GitHub raw URL
}
```

## 监控指标

### CDNMetrics

CDN 使用监控指标。

```go
type CDNMetrics struct {
    TotalRequests      int64 // 总请求数
    CDNSuccessCount    int64 // CDN成功次数
    CDNFailureCount    int64 // CDN失败次数
    RetryCount         int64 // 重试次数
    FallbackToOriginal int64 // 降级到原始URL次数
    FallbackToBase64   int64 // 降级到base64次数
}
```

## 状态流转

### URL 处理流程

```
GitHub Raw URL
    ↓
检测是否为GitHub raw URL
    ↓ (是)
转换为CDN URL (jsDelivr)
    ↓
使用CDN URL调用模型
    ↓ (成功)
完成
    ↓ (失败)
重试原始URL
    ↓ (成功)
完成
    ↓ (失败)
下载并转换为base64
    ↓
使用base64 data URL调用模型
```

## 验证规则

### URL 验证

- GitHub raw URL 必须匹配模式：`https://raw.githubusercontent.com/{owner}/{repo}/{branch}/{path}`
- owner、repo、branch、path 不能为空
- path 必须包含文件扩展名

### CDN URL 验证

- jsDelivr URL 格式：`https://cdn.jsdelivr.net/gh/{owner}/{repo}@{branch}/{path}`
- FastGit URL 格式：`https://raw.fastgit.org/{owner}/{repo}/{branch}/{path}`
