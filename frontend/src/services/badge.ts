/**
 * 勋章服务
 */

import type { BadgeDetailResponse, UserStats } from '../types/badge';
import { explorationStorage } from './storage';
import { cardStorage } from './storage';
import { conversationStorage } from './storage';

// 从环境变量读取API基础地址
// 生产环境默认使用相对路径（通过 Nginx 代理），开发环境使用完整 URL
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL !== undefined
  ? import.meta.env.VITE_API_BASE_URL
  : (import.meta.env.DEV 
    ? `http://${import.meta.env.VITE_BACKEND_HOST || 'localhost'}:${import.meta.env.VITE_BACKEND_PORT || '8877'}`
    : ''); // 生产环境默认使用相对路径，由 Nginx 代理

/**
 * 获取用户统计数据
 */
export async function getUserStats(): Promise<UserStats> {
  // 从本地存储获取统计数据
  const explorations = await explorationStorage.getAll();
  const cards = await cardStorage.getAll();
  const sessions = await conversationStorage.getAllSessions();

  const explorationCount = explorations.length;
  const collectionCount = cards.length;
  const conversationCount = sessions.length;

  // 调用后端API计算勋章等级
  const response = await fetch(`${API_BASE_URL}/api/badge/stats`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      explorationCount,
      collectionCount,
      conversationCount,
    }),
  });

  if (!response.ok) {
    throw new Error('获取勋章统计失败');
  }

  const data: BadgeDetailResponse = await response.json();
  return data.stats;
}

/**
 * 获取勋章详情（包括所有等级信息）
 */
export async function getBadgeDetail(): Promise<BadgeDetailResponse> {
  // 从本地存储获取统计数据
  const explorations = await explorationStorage.getAll();
  const cards = await cardStorage.getAll();
  const sessions = await conversationStorage.getAllSessions();

  const explorationCount = explorations.length;
  const collectionCount = cards.length;
  const conversationCount = sessions.length;

  // 调用后端API计算勋章等级
  const response = await fetch(`${API_BASE_URL}/api/badge/stats`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      explorationCount,
      collectionCount,
      conversationCount,
    }),
  });

  if (!response.ok) {
    throw new Error('获取勋章详情失败');
  }

  return await response.json();
}

