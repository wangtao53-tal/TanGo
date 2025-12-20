/**
 * 卡片文本提取工具
 * 用于从知识卡片中提取纯文本，用于文本转语音
 */

import type { KnowledgeCard } from '../types/exploration';
import type { ScienceCardContent, PoetryCardContent, EnglishCardContent } from '../types/exploration';

/**
 * 从知识卡片中提取用于文本转语音的文本内容
 */
export function extractCardText(card: KnowledgeCard): string {
  const parts: string[] = [];

  // 添加标题
  if (card.title) {
    parts.push(card.title);
  }

  switch (card.type) {
    case 'science': {
      const content = card.content as ScienceCardContent;
      if (content.name) {
        parts.push(content.name);
      }
      if (content.explanation) {
        parts.push(content.explanation);
      }
      if (content.facts && Array.isArray(content.facts)) {
        parts.push(...content.facts);
      }
      if (content.funFact) {
        parts.push('趣味知识点：' + content.funFact);
      }
      break;
    }
    case 'poetry': {
      const content = card.content as PoetryCardContent;
      if (content.poem) {
        parts.push(content.poem);
      }
      if (content.author) {
        parts.push('作者：' + content.author);
      }
      if (content.explanation) {
        parts.push('解释：' + content.explanation);
      }
      if (content.context) {
        parts.push('情境：' + content.context);
      }
      break;
    }
    case 'english': {
      const content = card.content as EnglishCardContent;
      if (content.words && Array.isArray(content.words)) {
        parts.push('核心单词：' + content.words.join('、'));
      }
      if (content.expressions && Array.isArray(content.expressions)) {
        parts.push('口语表达：' + content.expressions.join('。'));
      }
      if (content.pronunciation) {
        parts.push('发音：' + content.pronunciation);
      }
      break;
    }
  }

  return parts.join('。');
}

/**
 * 检测卡片内容的主要语言（中文/英文）
 */
export function detectCardLanguage(card: KnowledgeCard): string {
  const text = extractCardText(card);
  
  // 简单检测：如果包含大量英文字母，判断为英文
  const englishCharCount = (text.match(/[a-zA-Z]/g) || []).length;
  const totalCharCount = text.length;
  
  if (totalCharCount > 0 && englishCharCount / totalCharCount > 0.3) {
    return 'en-US';
  }
  
  return 'zh-CN';
}
