/**
 * 探索相关类型定义
 * 基于 data-model.md
 */

// 用户档案
export interface UserProfile {
  age: number; // 3-18
  grade?: string; // K1-K12格式，如 "K1", "G3"
  lastUpdated: string; // ISO 8601
}

// 探索记录
export interface ExplorationRecord {
  id: string; // UUID
  timestamp: string; // ISO 8601
  objectName: string; // 识别出的对象名称（中文）
  objectCategory: '自然类' | '生活类' | '人文类';
  confidence: number; // 0-1
  age: number; // 探索时的年龄
  imageData?: string; // base64，仅本地存储
  cards: KnowledgeCard[];
  collected: boolean;
}

// 知识卡片
export interface KnowledgeCard {
  id: string; // UUID
  explorationId: string;
  type: 'science' | 'poetry' | 'english';
  title: string;
  content: CardContent;
  collectedAt?: string; // ISO 8601
}

// 卡片内容（根据类型不同）
export type CardContent = 
  | ScienceCardContent 
  | PoetryCardContent 
  | EnglishCardContent;

// 类型守卫函数
export function isScienceCardContent(content: CardContent): content is ScienceCardContent {
  return 'name' in content && 'explanation' in content && 'facts' in content;
}

export function isPoetryCardContent(content: CardContent): content is PoetryCardContent {
  return 'poem' in content && 'explanation' in content;
}

export function isEnglishCardContent(content: CardContent): content is EnglishCardContent {
  return 'words' in content && 'expressions' in content;
}

// 科学认知卡内容
export interface ScienceCardContent {
  name: string;
  explanation: string;
  facts: string[];
  funFact: string;
}

// 古诗词/人文卡内容
export interface PoetryCardContent {
  poem: string;
  poemTitle?: string;
  author?: string;
  explanation: string;
  context: string;
}

// 英语表达卡内容
export interface EnglishCardContent {
  words: string[];
  expressions: string[];
  pronunciation?: string;
}

