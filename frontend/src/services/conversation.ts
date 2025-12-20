/**
 * 对话服务
 * 处理消息发送和接收，管理对话会话
 */

import { sendConversationMessage, recognizeIntent } from './api';
import { createSSEConnectionUnified, closeSSEConnection, type SSECallbacks } from './sse';
import { conversationStorage } from './storage';
import type { ConversationMessage } from '../types/conversation';
import type { ConversationStreamEvent, IdentificationContext, UnifiedStreamConversationRequest } from '../types/api';

/**
 * 生成会话ID
 */
function generateSessionId(): string {
  return `session-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
}

/**
 * 发送消息并处理响应（支持乐观更新）
 * @param onUserMessage 用户消息创建后的回调，用于立即显示用户消息（乐观更新）
 */
export async function sendMessage(
  content: string,
  type: 'text' | 'voice' | 'image' = 'text',
  sessionId?: string,
  identificationContext?: IdentificationContext,
  onUserMessage?: (message: ConversationMessage) => void
): Promise<{ sessionId: string; message: ConversationMessage; userMessage: ConversationMessage }> {
  const currentSessionId = sessionId || generateSessionId();

  // 创建用户消息（使用临时ID，乐观更新）
  const tempId = `msg-temp-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
  const userMessage: ConversationMessage = {
    id: tempId,
    type: type === 'text' ? 'text' : type === 'voice' ? 'voice' : 'image',
    content,
    timestamp: new Date().toISOString(),
    sender: 'user',
    sessionId: currentSessionId,
  };

  // 立即调用回调，显示用户消息（乐观更新）
  if (onUserMessage) {
    onUserMessage(userMessage);
  }

  // 保存用户消息到本地
  await conversationStorage.saveMessage(userMessage);

  // 发送到后端
  try {
    const response = await sendConversationMessage({
      sessionId: currentSessionId,
      type,
      content,
      inputType: type,
      identificationContext,
    });

    // 更新用户消息ID（如果后端返回了新的ID）
    if (response.userMessageId) {
      userMessage.id = response.userMessageId;
    }

    // 创建助手消息
    const assistantMessage: ConversationMessage = {
      id: response.message?.id || `msg-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
      type: response.message?.type || response.type || 'text',
      sender: 'assistant',
      content: response.message?.content || response.content || response,
      timestamp: response.message?.timestamp || new Date().toISOString(),
      sessionId: currentSessionId,
      isStreaming: response.message?.isStreaming,
    };

    // 保存助手消息到本地
    await conversationStorage.saveMessage(assistantMessage);

    return {
      sessionId: currentSessionId,
      message: assistantMessage,
      userMessage,
    };
  } catch (error) {
    console.error('发送消息失败:', error);
    // 如果失败，可以标记用户消息为错误状态
    throw error;
  }
}

/**
 * 创建统一流式连接（支持文本、语音、图片三种输入方式）
 */
export function createStreamConnectionUnified(
  request: UnifiedStreamConversationRequest,
  callbacks: {
    onMessage?: (message: ConversationMessage) => void;
    onError?: (error: Error) => void;
    onClose?: () => void;
  }
): AbortController {
  let assistantMessageId: string | undefined;
  let fullText = '';
  let isMarkdown = false;

  const sseCallbacks: SSECallbacks = {
    onMessage: async (event: ConversationStreamEvent) => {
      if (event.type === 'connected') {
        assistantMessageId = event.messageId;
        return;
      }

      if (event.type === 'done') {
        // 保存完整的助手消息
        if (assistantMessageId && fullText) {
          const assistantMessage: ConversationMessage = {
            id: assistantMessageId,
            type: 'text',
            content: fullText,
            timestamp: new Date().toISOString(),
            sender: 'assistant',
            sessionId: request.sessionId || '',
            markdown: isMarkdown,
          };
          await conversationStorage.saveMessage(assistantMessage);
        }
        callbacks.onClose?.();
        return;
      }

      if (event.type === 'error') {
        callbacks.onError?.(new Error(event.message || event.content?.message || '未知错误'));
        return;
      }

      if (event.type === 'voice_recognized') {
        // 语音识别完成事件，不需要创建消息
        return;
      }

      if (event.type === 'image_uploaded') {
        // 图片上传完成事件，不需要创建消息
        return;
      }

      if (event.type === 'message') {
        // 流式文本消息
        assistantMessageId = event.messageId || assistantMessageId;
        const char = event.content as string;
        if (char) {
          fullText += char;
          // 检测Markdown格式
          if (!isMarkdown && fullText.length > 10) {
            isMarkdown = /[#*_`\[\]()]/.test(fullText);
          }

          // 创建临时消息对象（用于实时显示）
          const message: ConversationMessage = {
            id: assistantMessageId || `msg-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
            type: 'text',
            content: fullText,
            timestamp: new Date().toISOString(),
            sender: 'assistant',
            sessionId: request.sessionId || '',
            isStreaming: true,
            markdown: isMarkdown,
          };

          // 调用回调（实时更新UI）
          callbacks.onMessage?.(message);
        }
        return;
      }

      // 其他类型的事件（card等）
      const message: ConversationMessage = {
        id: event.messageId || `msg-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        type: event.type === 'card' ? 'card' : event.type === 'image' ? 'image' : 'text',
        content: event.content,
        timestamp: new Date().toISOString(),
        sender: 'assistant',
        sessionId: request.sessionId || '',
      };

      // 保存到本地
      await conversationStorage.saveMessage(message);

      // 调用回调
      callbacks.onMessage?.(message);
    },
    onError: callbacks.onError,
    onClose: callbacks.onClose,
  };

  return createSSEConnectionUnified(request, sseCallbacks);
}

/**
 * 创建SSE连接接收流式返回（兼容旧版本）
 * @deprecated 请使用 createStreamConnectionUnified
 */
export function createStreamConnection(
  sessionId: string,
  callbacks: {
    onMessage?: (message: ConversationMessage) => void;
    onError?: (error: Error) => void;
    onClose?: () => void;
  }
): EventSource {
  // 转换为统一接口调用
  const request: UnifiedStreamConversationRequest = {
    messageType: 'text',
    message: '',
    sessionId,
  };
  const abortController = createStreamConnectionUnified(request, callbacks);
  // 返回一个模拟的EventSource对象（为了兼容性）
  return {
    close: () => abortController.abort(),
  } as EventSource;
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
