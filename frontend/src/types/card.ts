/**
 * 知识卡片相关类型定义
 * 基于 data-model.md
 */

import type { KnowledgeCard, CardContent, ScienceCardContent, PoetryCardContent, EnglishCardContent } from './exploration';

export type { KnowledgeCard, CardContent, ScienceCardContent, PoetryCardContent, EnglishCardContent };

// 卡片类型枚举
export const CardType = {
  SCIENCE: 'science',
  POETRY: 'poetry',
  ENGLISH: 'english',
} as const;

export type CardType = typeof CardType[keyof typeof CardType];

// 卡片类型显示名称
export const CardTypeNames: Record<string, string> = {
  [CardType.SCIENCE]: '科学认知',
  [CardType.POETRY]: '人文素养',
  [CardType.ENGLISH]: '语言能力',
};

// 卡片类型颜色映射（基于设计稿）
export const CardTypeColors: Record<string, string> = {
  [CardType.SCIENCE]: 'science-green',
  [CardType.POETRY]: 'sunny-orange',
  [CardType.ENGLISH]: 'sky-blue',
};

