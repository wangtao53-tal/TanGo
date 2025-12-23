/**
 * SSE (Server-Sent Events) 服务
 * 用于接收流式返回数据
 */

import type { ConversationStreamEvent, UnifiedStreamConversationRequest } from '../types/api';

export interface SSECallbacks {
  onMessage?: (event: ConversationStreamEvent) => void;
  onError?: (error: Error) => void;
  onClose?: () => void;
}

// 从环境变量读取API基础地址
// 生产环境默认使用相对路径（通过 Nginx 代理），开发环境使用完整 URL
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL !== undefined
  ? import.meta.env.VITE_API_BASE_URL
  : (import.meta.env.DEV 
    ? `http://${import.meta.env.VITE_BACKEND_HOST || 'localhost'}:${import.meta.env.VITE_BACKEND_PORT || '8877'}`
    : ''); // 生产环境默认使用相对路径，由 Nginx 代理

/**
 * 创建SSE连接（统一接口，支持POST请求）
 */
export function createSSEConnectionUnified(
  request: UnifiedStreamConversationRequest,
  callbacks: SSECallbacks
): AbortController {
  const abortController = new AbortController();
  
  fetch(`${API_BASE_URL}/api/conversation/stream`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(request),
    signal: abortController.signal,
  })
    .then(async (response) => {
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      if (!response.body) {
        throw new Error('Response body is null');
      }

      const reader = response.body.getReader();
      const decoder = new TextDecoder();
      let buffer = '';

      while (true) {
        const { done, value } = await reader.read();

        if (done) {
          callbacks.onClose?.();
          break;
        }

        buffer += decoder.decode(value, { stream: true });
        const lines = buffer.split('\n');
        buffer = lines.pop() || '';

        for (const line of lines) {
          if (line.startsWith('event: ')) {
            // 事件类型已解析，继续处理下一行
            continue;
          }
          if (line.startsWith('data: ')) {
            const dataStr = line.substring(6).trim();
            if (dataStr) {
              try {
                const data: ConversationStreamEvent = JSON.parse(dataStr);
                callbacks.onMessage?.(data);
                
                // 如果收到done事件，关闭连接
                if (data.type === 'done') {
                  callbacks.onClose?.();
                  return;
                }
              } catch (error) {
                console.error('解析SSE消息失败:', error);
                callbacks.onError?.(error as Error);
              }
            }
          }
        }
      }
    })
    .catch((error) => {
      if (error.name === 'AbortError') {
        return; // 用户主动取消，不触发错误回调
      }
      console.error('SSE连接错误:', error);
      callbacks.onError?.(error as Error);
    });

  return abortController;
}

/**
 * 创建SSE连接（兼容旧版本，使用GET请求）
 * @deprecated 请使用 createSSEConnectionUnified
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
export function closeSSEConnection(eventSource: EventSource | AbortController): void {
  if (eventSource instanceof AbortController) {
    eventSource.abort();
  } else {
    eventSource.close();
  }
}
