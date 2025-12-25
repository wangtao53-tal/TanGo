/**
 * 勋章系统类型定义
 */

// 勋章等级信息
export interface BadgeLevel {
  level: number; // 1-10
  title: string; // 称号，如"小小专家"、"自然大师"
  minScore: number; // 最低分数要求
  icon: string; // 图标标识
  color: string; // 主题颜色
  description: string; // 等级描述
}

// 用户统计数据
export interface UserStats {
  explorationCount: number; // 探索次数
  collectionCount: number; // 收藏次数
  conversationCount: number; // 对话次数（会话数）
  totalScore: number; // 总分
  currentLevel: number; // 当前等级 1-10
  currentLevelInfo: BadgeLevel; // 当前等级信息
  nextLevelInfo?: BadgeLevel; // 下一等级信息
  progress: number; // 当前等级进度 0-100
}

// 勋章详情响应
export interface BadgeDetailResponse {
  stats: UserStats;
  allLevels: BadgeLevel[]; // 所有等级信息
  recentUpgrade?: {
    fromLevel: number;
    toLevel: number;
    upgradedAt: string;
  }; // 最近升级信息
}

