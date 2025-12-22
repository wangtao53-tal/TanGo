/**
 * 流式对话Hook
 * 用于管理流式对话的状态和连接
 */

import { useState, useEffect, useRef, useCallback } from 'react';
import { flushSync } from 'react-dom';
import { closeSSEConnection } from '../services/sse';
import type { ConversationStreamEvent } from '../types/api';

export interface UseStreamConversationOptions {
  sessionId: string;
  userAge?: number;
  onMessage?: (text: string) => void;
  onError?: (error: Error) => void;
  onComplete?: () => void;
}

export interface UseStreamConversationReturn {
  streamingText: string;
  isStreaming: boolean;
  startStream: (message: string) => void;
  stopStream: () => void;
  error: Error | null;
}

export function useStreamConversation(
  options: UseStreamConversationOptions
): UseStreamConversationReturn {
  const { sessionId, userAge, onMessage, onError, onComplete } = options;
  const [streamingText, setStreamingText] = useState('');
  const [isStreaming, setIsStreaming] = useState(false);
  const [error, setError] = useState<Error | null>(null);
  const eventSourceRef = useRef<EventSource | null>(null);
  const currentTextRef = useRef('');

  const startStream = useCallback(
    (message: string) => {
      // 停止之前的流
      if (eventSourceRef.current) {
        closeSSEConnection(eventSourceRef.current);
      }

      setIsStreaming(true);
      setStreamingText('');
      setError(null);
      currentTextRef.current = '';

      // 构建SSE URL
      const API_BASE_URL =
        import.meta.env.VITE_API_BASE_URL ||
        (import.meta.env.DEV
          ? `http://${import.meta.env.VITE_BACKEND_HOST || 'localhost'}:${import.meta.env.VITE_BACKEND_PORT || '8877'}`
          : 'http://localhost:8877');

      const params = new URLSearchParams({
        sessionId,
        message,
      });
      if (userAge) {
        params.append('userAge', userAge.toString());
      }

      const url = `${API_BASE_URL}/api/conversation/stream?${params.toString()}`;
      const eventSource = new EventSource(url);
      eventSourceRef.current = eventSource;

      // 处理连接建立事件
      eventSource.addEventListener('connected', (event: MessageEvent) => {
        try {
          const data: ConversationStreamEvent = JSON.parse(event.data);
          console.log('SSE连接建立:', data.sessionId);
        } catch (err) {
          console.error('解析connected事件失败:', err);
        }
      });

      // 处理消息事件
      eventSource.addEventListener('message', (event: MessageEvent) => {
        try {
          const data: ConversationStreamEvent = JSON.parse(event.data);
          if (data.type === 'message' && data.content) {
            // 立即追加文本内容到ref
            currentTextRef.current += data.content;
            // 使用flushSync强制同步更新，确保实时渲染
            flushSync(() => {
              setStreamingText(currentTextRef.current);
            });
            onMessage?.(currentTextRef.current);
          }
        } catch (err) {
          console.error('解析message事件失败:', err);
        }
      });

      // 处理完成事件
      eventSource.addEventListener('done', () => {
        setIsStreaming(false);
        closeSSEConnection(eventSource);
        eventSourceRef.current = null;
        onComplete?.();
      });

      // 处理错误事件
      eventSource.addEventListener('error', (event: MessageEvent) => {
        try {
          const data: ConversationStreamEvent = JSON.parse(event.data);
          const error = new Error(data.message || '流式对话错误');
          setError(error);
          setIsStreaming(false);
          closeSSEConnection(eventSource);
          eventSourceRef.current = null;
          onError?.(error);
        } catch (err) {
          const error = new Error('解析错误事件失败');
          setError(error);
          setIsStreaming(false);
          closeSSEConnection(eventSource);
          eventSourceRef.current = null;
          onError?.(error);
        }
      });

      // 处理连接错误
      eventSource.onerror = (err) => {
        console.error('SSE连接错误:', err);
        const error = new Error('SSE连接错误');
        setError(error);
        setIsStreaming(false);
        closeSSEConnection(eventSource);
        eventSourceRef.current = null;
        onError?.(error);
      };
    },
    [sessionId, userAge, onMessage, onError, onComplete]
  );

  const stopStream = useCallback(() => {
    if (eventSourceRef.current) {
      closeSSEConnection(eventSourceRef.current);
      eventSourceRef.current = null;
      setIsStreaming(false);
    }
  }, []);

  // 清理函数
  useEffect(() => {
    return () => {
      if (eventSourceRef.current) {
        closeSSEConnection(eventSourceRef.current);
      }
    };
  }, []);

  return {
    streamingText,
    isStreaming,
    startStream,
    stopStream,
    error,
  };
}

