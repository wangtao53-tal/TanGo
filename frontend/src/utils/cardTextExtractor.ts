/**
 * 卡片文本提取工具
 * 用于从知识卡片中提取纯文本，用于文本转语音
 */

import type { KnowledgeCard } from '../types/exploration';
import type { ScienceCardContent, PoetryCardContent, EnglishCardContent } from '../types/exploration';

/**
 * 移除文本中的所有emoji表情符号
 * 用于语音播放时，避免TTS引擎尝试朗读emoji
 */
function removeEmojis(text: string): string {
  if (!text) return text;
  
  // 更全面的emoji Unicode范围正则表达式
  // 包括：表情符号、符号、图标、旗帜、箭头等
  const emojiRegex = /[\u{1F600}-\u{1F64F}]|[\u{1F300}-\u{1F5FF}]|[\u{1F680}-\u{1F6FF}]|[\u{1F900}-\u{1F9FF}]|[\u{2600}-\u{26FF}]|[\u{2700}-\u{27BF}]|[\u{1F1E0}-\u{1F1FF}]|[\u{1F191}-\u{1F251}]|[\u{2934}\u{2935}]|[\u{2190}-\u{21FF}]|[\u{2B00}-\u{2BFF}]|[\u{FE00}-\u{FE0F}]|[\u{200D}]/gu;
  
  // 移除emoji并清理多余的空格
  return text.replace(emojiRegex, '').replace(/\s+/g, ' ').trim();
}

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
      // 处理字段映射：后端可能返回 keywords，前端使用 words
      const rawContent = card.content as any;
      const content: EnglishCardContent = {
        words: rawContent.words || rawContent.keywords || [],
        expressions: rawContent.expressions || [],
        pronunciation: rawContent.pronunciation,
      };
      // 优化：单词之间用逗号和空格分隔，让语音播放更清晰
      if (content.words && Array.isArray(content.words) && content.words.length > 0) {
        // 使用英文逗号和空格，让TTS引擎能更好地识别单词边界
        parts.push('核心单词：' + content.words.join(', '));
      }
      if (content.expressions && Array.isArray(content.expressions) && content.expressions.length > 0) {
        // 表达式之间用句号分隔，每个表达式单独播放
        parts.push('口语表达：' + content.expressions.join('. '));
      }
      if (content.pronunciation) {
        parts.push('发音：' + content.pronunciation);
      }
      break;
    }
  }

  const fullText = parts.join('。');
  // 移除所有emoji表情符号，避免TTS引擎尝试朗读
  return removeEmojis(fullText);
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
