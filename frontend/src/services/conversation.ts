/**
 * 对话服务
 * 处理消息发送和接收，管理对话会话
 */

import { sendConversationMessage, recognizeIntent } from './api';
import { createSSEConnection, closeSSEConnection, type SSECallbacks } from './sse';
import { conversationStorage } from './storage';
import type { ConversationMessage } from '../types/conversation';
import type { ConversationStreamEvent } from '../types/api';

/**
 * 生成会话ID
 */
function generateSessionId(): string {
  return `session-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
}

/**
 * 发送消息并处理响应
 */
export async function sendMessage(
  content: string,
  type: 'text' | 'voice' | 'image' = 'text',
  sessionId?: string
): Promise<{ sessionId: string; message: ConversationMessage }> {
  const currentSessionId = sessionId || generateSessionId();

  // 创建用户消息
  const userMessage: ConversationMessage = {
    id: `msg-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
    type: type === 'text' ? 'text' : type === 'voice' ? 'voice' : 'image',
    content,
    timestamp: new Date().toISOString(),
    sender: 'user',
    sessionId: currentSessionId,
  };

  // 保存用户消息到本地
  await conversationStorage.saveMessage(userMessage);

  // 发送到后端
  try {
    const response = await sendConversationMessage({
      sessionId: currentSessionId,
      type,
      content,
      inputType: type,
    });

    // 创建助手消息
    const assistantMessage: ConversationMessage = {
      id: `msg-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
      type: response.type || 'text',
      content: response.content || response,
      timestamp: new Date().toISOString(),
      sender: 'assistant',
      sessionId: currentSessionId,
    };

    // 保存助手消息到本地
    await conversationStorage.saveMessage(assistantMessage);

    return {
      sessionId: currentSessionId,
      message: assistantMessage,
    };
  } catch (error) {
    console.error('发送消息失败:', error);
    throw error;
  }
}

/**
 * 创建SSE连接接收流式返回
 */
export function createStreamConnection(
  sessionId: string,
  callbacks: {
    onMessage?: (message: ConversationMessage) => void;
    onError?: (error: Error) => void;
    onClose?: () => void;
  }
): EventSource {
  const sseCallbacks: SSECallbacks = {
    onMessage: async (event: ConversationStreamEvent) => {
      if (event.type === 'done') {
        callbacks.onClose?.();
        return;
      }

      if (event.type === 'error') {
        callbacks.onError?.(new Error(event.message || '未知错误'));
        return;
      }

      // 创建消息对象
      const message: ConversationMessage = {
        id: `msg-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        type: event.type === 'card' ? 'card' : event.type === 'image' ? 'image' : 'text',
        content: event.content,
        timestamp: new Date().toISOString(),
        sender: 'assistant',
        sessionId,
      };

      // 保存到本地
      await conversationStorage.saveMessage(message);

      // 调用回调
      callbacks.onMessage?.(message);
    },
    onError: callbacks.onError,
    onClose: callbacks.onClose,
  };

  return createSSEConnection(sessionId, sseCallbacks);
}

/**
 * 识别用户意图
 */
export async function recognizeUserIntent(
  text: string,
  sessionId?: string
): Promise<{ intent: string; confidence: number }> {
  try {
    const response = await recognizeIntent({
      text,
      sessionId,
    });
    return {
      intent: response.intent,
      confidence: response.confidence,
    };
  } catch (error) {
    console.error('意图识别失败:', error);
    // 降级处理：默认返回文本回答意图
    return {
      intent: 'text_response',
      confidence: 0.5,
    };
  }
}

/**
 * 获取会话的所有消息
 */
export async function getSessionMessages(sessionId: string): Promise<ConversationMessage[]> {
  return await conversationStorage.getMessagesBySessionId(sessionId);
}

/**
 * 获取所有会话
 */
export async function getAllSessions() {
  return await conversationStorage.getAllSessions();
}
