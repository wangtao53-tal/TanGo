/**
 * POST + SSE 服务
 * 使用 fetch API 发送 POST 请求，然后手动解析 SSE 流
 * 解决 EventSource 只支持 GET 请求的限制
 */

import type { ConversationStreamEvent, StreamConversationRequest } from '../types/api';

export interface PostSSECallbacks {
  onMessage?: (event: ConversationStreamEvent) => void;
  onError?: (error: Error) => void;
  onClose?: () => void;
}

// 从环境变量读取API基础地址
const API_BASE_URL =
  import.meta.env.VITE_API_BASE_URL ||
  (import.meta.env.DEV
    ? `http://${import.meta.env.VITE_BACKEND_HOST || 'localhost'}:${import.meta.env.VITE_BACKEND_PORT || '8877'}`
    : 'http://localhost:8877');

/**
 * 创建POST + SSE连接
 * 使用 fetch API 发送 POST 请求，然后手动解析 SSE 流
 */
export function createPostSSEConnection(
  request: StreamConversationRequest,
  callbacks: PostSSECallbacks
): AbortController {
  const abortController = new AbortController();

  // 发送POST请求
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

      // 读取流式数据
      while (true) {
        const { done, value } = await reader.read();

        if (done) {
          // 处理缓冲区中剩余的数据
          if (buffer.trim()) {
            processSSEBuffer(buffer, callbacks);
          }
          callbacks.onClose?.();
          break;
        }

        // 解码数据（立即处理，不等待缓冲区满）
        const chunk = decoder.decode(value, { stream: true });
        buffer += chunk;

        // 立即处理完整的事件（不等待所有数据到达）
        buffer = processSSEBuffer(buffer, callbacks);
      }
    })
    .catch((error) => {
      if (error.name === 'AbortError') {
        // 用户主动取消，不调用错误回调
        return;
      }
      console.error('POST SSE连接错误:', error);
      callbacks.onError?.(error);
    });

  return abortController;
}

/**
 * 处理SSE缓冲区，解析并触发回调
 * 返回剩余的不完整数据
 */
function processSSEBuffer(
  buffer: string,
  callbacks: PostSSECallbacks
): string {
  // 解析SSE格式的数据
  // SSE格式: "event: eventType\ndata: {...}\n\n"
  const lines = buffer.split('\n');
  const remainingLines: string[] = [];
  let currentEvent = '';
  let currentData = '';

  for (let i = 0; i < lines.length; i++) {
    const line = lines[i];

    if (line.startsWith('event:')) {
      // 如果之前有未完成的事件，先处理它
      if (currentEvent && currentData) {
        processSSEEvent(currentEvent, currentData, callbacks);
        currentEvent = '';
        currentData = '';
      }
      currentEvent = line.substring(6).trim();
    } else if (line.startsWith('data:')) {
      // 处理多行data（累加）
      const dataLine = line.substring(5).trim();
      if (currentData) {
        currentData += '\n' + dataLine;
      } else {
        currentData = dataLine;
      }
    } else if (line === '') {
      // 空行表示一个完整的事件
      if (currentEvent && currentData) {
        processSSEEvent(currentEvent, currentData, callbacks);
        currentEvent = '';
        currentData = '';
      }
    } else if (line.trim() !== '') {
      // 非空行且不是event或data，可能是数据的一部分
      if (currentData) {
        currentData += '\n' + line;
      }
    }
  }

  // 如果有未完成的事件，保留在缓冲区中
  if (currentEvent || currentData) {
    if (currentEvent) {
      remainingLines.push(`event: ${currentEvent}`);
    }
    if (currentData) {
      remainingLines.push(`data: ${currentData}`);
    }
  }

  return remainingLines.length > 0 ? remainingLines.join('\n') : '';
}

/**
 * 处理单个SSE事件
 */
function processSSEEvent(
  eventType: string,
  dataStr: string,
  callbacks: PostSSECallbacks
): void {
  try {
    const data: ConversationStreamEvent = JSON.parse(dataStr);
    callbacks.onMessage?.(data);

    // 如果收到done事件，关闭连接
    if (data.type === 'done') {
      callbacks.onClose?.();
      return;
    }

    // 如果收到error事件，调用错误回调
    if (data.type === 'error') {
      const errorMessage =
        (data.content as any)?.message || data.message || '未知错误';
      callbacks.onError?.(new Error(errorMessage));
      return;
    }
  } catch (err) {
    console.error('解析SSE消息失败:', err, 'eventType:', eventType, 'data:', dataStr);
    // 不因为解析错误而中断连接，继续处理后续事件
  }
}

/**
 * 关闭POST + SSE连接
 */
export function closePostSSEConnection(abortController: AbortController): void {
  abortController.abort();
}

