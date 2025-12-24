/**
 * 分享功能工具函数
 */

import { createShare } from '../services/api';
import type { CreateShareRequest } from '../types/api';
import type { ExplorationRecord, KnowledgeCard } from '../types/exploration';

/**
 * 判断图片数据是否是 base64 格式
 */
function isBase64Image(imageData?: string): boolean {
  if (!imageData) return false;
  // base64 图片通常以 data:image/ 开头，或者是很长的字符串（超过1000字符）
  // URL 通常以 http:// 或 https:// 开头
  if (imageData.startsWith('data:image/') || imageData.startsWith('data:image')) {
    return true;
  }
  // 如果字符串很长且不包含 http，可能是纯 base64（不含 data URL 前缀）
  if (imageData.length > 1000 && !imageData.startsWith('http://') && !imageData.startsWith('https://')) {
    return true;
  }
  return false;
}

/**
 * 将探索记录转换为分享格式
 */
function convertRecordToShareFormat(record: ExplorationRecord) {
  // 只包含 URL 格式的图片，不包含 base64
  const imageData = record.imageData && !isBase64Image(record.imageData) 
    ? record.imageData 
    : undefined;

  return {
    id: record.id,
    timestamp: record.timestamp,
    objectName: record.objectName,
    objectCategory: record.objectCategory,
    age: record.age,
    imageData: imageData, // 只包含 URL，不包含 base64
    cards: record.cards.map((card) => {
      // 确保 content 是一个对象
      let content: Record<string, unknown>;
      
      if (typeof card.content === 'object' && card.content !== null && !Array.isArray(card.content)) {
        // 已经是对象，直接使用
        content = card.content as Record<string, unknown>;
      } else {
        // 如果不是对象，尝试通过 JSON 转换
        try {
          const jsonStr = JSON.stringify(card.content);
          const parsed = JSON.parse(jsonStr);
          if (typeof parsed === 'object' && parsed !== null && !Array.isArray(parsed)) {
            content = parsed;
          } else {
            // 如果解析后仍不是对象，包装成对象
            content = { value: parsed };
          }
        } catch {
          // JSON 转换失败，包装成对象
          content = { value: card.content };
        }
      }
      
      return {
        type: card.type,
        title: card.title,
        content: content,
      };
    }),
  };
}

/**
 * 将知识卡片转换为分享格式
 */
function convertCardToShareFormat(card: KnowledgeCard) {
  // 确保 content 是一个对象
  // 使用 JSON 序列化/反序列化来确保类型正确
  let content: Record<string, unknown>;
  
  if (typeof card.content === 'object' && card.content !== null && !Array.isArray(card.content)) {
    // 已经是对象，直接使用
    content = card.content as Record<string, unknown>;
  } else {
    // 如果不是对象，尝试通过 JSON 转换
    try {
      const jsonStr = JSON.stringify(card.content);
      const parsed = JSON.parse(jsonStr);
      if (typeof parsed === 'object' && parsed !== null && !Array.isArray(parsed)) {
        content = parsed;
      } else {
        // 如果解析后仍不是对象，包装成对象
        content = { value: parsed };
      }
    } catch {
      // JSON 转换失败，包装成对象
      content = { value: card.content };
    }
  }

  return {
    id: card.id,
    explorationId: card.explorationId,
    type: card.type,
    title: card.title,
    content: content,
    collectedAt: card.collectedAt,
  };
}

/**
 * 创建分享链接
 * @param records 探索记录列表
 * @param cards 知识卡片列表（可选）
 * @returns 分享链接URL
 */
export async function createShareLink(
  records: ExplorationRecord[],
  cards: KnowledgeCard[] = []
): Promise<string> {
  // 过滤出有 URL 图片的记录（不包含 base64）
  const recordsWithUrl = records.filter((r) => {
    if (!r.imageData) return true; // 没有图片也可以分享
    return !isBase64Image(r.imageData); // 只保留 URL 格式的图片
  });

  // 检查是否有记录被过滤掉
  const base64Records = records.filter((r) => r.imageData && isBase64Image(r.imageData));
  if (recordsWithUrl.length === 0) {
    if (base64Records.length > 0) {
      throw new Error('所有探索记录的图片都是 base64 格式，无法创建分享。请确保图片已上传到服务器。');
    } else {
      throw new Error('没有可分享的探索记录。');
    }
  }

  // 如果有部分记录被过滤（base64），给出提示
  if (base64Records.length > 0) {
    console.warn(`有 ${base64Records.length} 条探索记录的图片是 base64 格式，已自动跳过，只分享 ${recordsWithUrl.length} 条有 URL 的记录`);
  }

  // 转换为分享格式
  const explorationRecords = recordsWithUrl.map(convertRecordToShareFormat);
  const collectedCards = cards.map(convertCardToShareFormat);

  // 调用API创建分享链接
  const request: CreateShareRequest = {
    explorationRecords,
    collectedCards,
  };

  const response = await createShare(request);

  // 生成完整的分享URL
  const baseUrl = window.location.origin;
  const shareUrl = `${baseUrl}/share/${response.shareId}`;

  return shareUrl;
}

/**
 * 复制文本到剪贴板
 */
export async function copyToClipboard(text: string): Promise<boolean> {
  try {
    if (navigator.clipboard && navigator.clipboard.writeText) {
      await navigator.clipboard.writeText(text);
      return true;
    } else {
      // 降级方案：使用传统的复制方法
      const textArea = document.createElement('textarea');
      textArea.value = text;
      textArea.style.position = 'fixed';
      textArea.style.opacity = '0';
      document.body.appendChild(textArea);
      textArea.select();
      const success = document.execCommand('copy');
      document.body.removeChild(textArea);
      return success;
    }
  } catch (error) {
    console.error('复制失败:', error);
    return false;
  }
}

