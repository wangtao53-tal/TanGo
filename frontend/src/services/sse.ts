/**
 * SSE (Server-Sent Events) 服务
 * 用于接收流式返回数据
 */

import type { ConversationStreamEvent } from '../types/api';

export interface SSECallbacks {
  onMessage?: (event: ConversationStreamEvent) => void;
  onError?: (error: Error) => void;
  onClose?: () => void;
}

// 从环境变量读取API基础地址，如果没有配置则使用默认值
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 
  (import.meta.env.DEV 
    ? `http://${import.meta.env.VITE_BACKEND_HOST || 'localhost'}:${import.meta.env.VITE_BACKEND_PORT || '8877'}`
    : 'http://localhost:8877');

/**
 * 创建SSE连接
 */
export function createSSEConnection(
  sessionId: string,
  callbacks: SSECallbacks
): EventSource {
  const url = `${API_BASE_URL}/api/conversation/stream?sessionId=${sessionId}`;
  const eventSource = new EventSource(url);

  eventSource.onmessage = (event) => {
    try {
      const data: ConversationStreamEvent = JSON.parse(event.data);
      callbacks.onMessage?.(data);
      
      // 如果收到done事件，关闭连接
      if (data.type === 'done') {
        eventSource.close();
        callbacks.onClose?.();
      }
    } catch (error) {
      console.error('解析SSE消息失败:', error);
      callbacks.onError?.(error as Error);
    }
  };

  eventSource.onerror = (error) => {
    console.error('SSE连接错误:', error);
    callbacks.onError?.(new Error('SSE连接错误'));
  };

  eventSource.addEventListener('error', () => {
    eventSource.close();
    callbacks.onClose?.();
  });

  return eventSource;
}

/**
 * 关闭SSE连接
 */
export function closeSSEConnection(eventSource: EventSource): void {
  eventSource.close();
}
