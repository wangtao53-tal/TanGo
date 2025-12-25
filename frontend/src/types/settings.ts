/**
 * 用户设置类型定义
 */

// 语言设置
export type Language = 'zh' | 'en';

// 用户设置
export interface UserSettings {
  language: Language; // 语言设置，默认中文
  grade?: string; // 年级设置（K1-K12）
  lastUpdated: string; // 最后更新时间（ISO 8601）
}
