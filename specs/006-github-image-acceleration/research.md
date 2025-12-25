# 研究文档：GitHub 图片加速方案

**功能**: 006-github-image-acceleration  
**创建日期**: 2025-12-20

## 研究目标

解决 GitHub raw.githubusercontent.com URL 偶发访问超时问题，找到稳定可靠的图片加速方案。

## 研究问题与决策

### 1. GitHub 图片 CDN 加速方案选择

**问题**: 如何解决 GitHub raw URL 访问不稳定和超时问题？

**决策**: 使用 jsDelivr CDN 作为主要加速方案，FastGit 作为备选

**理由**: 
- jsDelivr 是免费的开源 CDN，稳定可靠，全球加速
- 支持 GitHub 仓库文件直接加速
- 访问速度快，成功率高
- FastGit 作为备选方案，提供额外保障

**URL 转换规则**:
- 原始URL: `https://raw.githubusercontent.com/{owner}/{repo}/{branch}/{path}`
- jsDelivr: `https://cdn.jsdelivr.net/gh/{owner}/{repo}@{branch}/{path}`
- FastGit: `https://raw.fastgit.org/{owner}/{repo}/{branch}/{path}`

**示例**:
```
原始: https://raw.githubusercontent.com/wangtao53-tal/image/main/tango/IMG_5829.JPG
jsDelivr: https://cdn.jsdelivr.net/gh/wangtao53-tal/image@main/tango/IMG_5829.JPG
FastGit: https://raw.fastgit.org/wangtao53-tal/image/main/tango/IMG_5829.JPG
```

**替代方案考虑**: 
- **ghproxy**: 需要代理整个URL，可能增加延迟 - 作为最后备选
- **Static CDN**: 稳定性不如 jsDelivr - 不采用
- **自建CDN**: 成本高，维护复杂 - 不采用

### 2. URL 转换实现位置

**问题**: URL 转换应该在哪个层级实现？

**决策**: 在 ImageRecognitionNode 中实现 URL 转换

**理由**: 
- 识别节点是使用图片URL的地方，在这里转换最合适
- 不影响存储层的URL生成逻辑，保持向后兼容
- 可以灵活控制是否使用CDN

**实现方式**:
- 检测 GitHub raw URL 模式
- 自动转换为 jsDelivr CDN URL
- 如果CDN失败，重试原始URL
- 如果仍失败，下载并转换为base64

**替代方案考虑**: 
- **存储层转换**: 在 GitHubStorage 中生成CDN URL - 被拒绝，因为需要保持原始URL的兼容性
- **前端转换**: 前端转换URL - 被拒绝，因为后端需要控制CDN使用策略

### 3. 重试和降级策略

**问题**: CDN失败时如何处理？

**决策**: 实现三级降级策略

**策略**:
1. **第一级**: 使用 jsDelivr CDN URL
2. **第二级**: CDN失败时重试原始 GitHub raw URL
3. **第三级**: 如果仍失败，下载图片并转换为base64 data URL

**理由**: 
- 保证最大兼容性
- 逐步降级，不影响正常流程
- 提供多重保障

**超时设置**:
- CDN请求超时：10秒
- 原始URL重试超时：15秒
- 总超时：不超过模型调用超时（45秒）

### 4. 配置和监控

**问题**: 如何控制CDN使用和监控效果？

**决策**: 
- 添加配置选项控制CDN使用（默认启用）
- 记录CDN使用情况和失败率
- 提供监控指标

**配置选项**:
- `EnableGitHubCDN`: 是否启用CDN加速（默认true）
- `CDNProvider`: CDN提供商（jsdelivr/fastgit，默认jsdelivr）
- `CDNRetryEnabled`: 是否启用CDN失败重试（默认true）

**监控指标**:
- CDN使用次数
- CDN成功率
- CDN失败重试次数
- 降级到原始URL次数
- 降级到base64次数

## 技术选型总结

| 技术点 | 选型 | 理由 |
|--------|------|------|
| 主要CDN | jsDelivr | 稳定可靠，全球加速，免费 |
| 备选CDN | FastGit | 国内访问快，作为备选 |
| 实现位置 | ImageRecognitionNode | 不影响存储层，灵活可控 |
| 降级策略 | 三级降级 | CDN → 原始URL → base64 |
| 配置方式 | 环境变量/配置文件 | 灵活配置，易于调整 |

## 风险评估

**低风险**:
- URL转换逻辑简单，不改变核心功能
- jsDelivr CDN稳定可靠

**中风险**:
- CDN服务可能偶尔不可用（有降级策略）
- FastGit可能不稳定（有备选方案）

**高风险**:
- 无（有完善的降级策略）

## 实施建议

1. **第一阶段**: 实现 jsDelivr CDN 转换和基本重试机制
2. **第二阶段**: 添加 FastGit 备选支持
3. **第三阶段**: 完善配置和监控

## 参考资源

- jsDelivr 文档: https://www.jsdelivr.com/
- FastGit 文档: https://github.com/FastGitORG
- GitHub raw URL 格式: https://raw.githubusercontent.com/
