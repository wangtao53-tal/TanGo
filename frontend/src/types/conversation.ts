/**
 * 对话相关类型定义
 */

// 对话消息类型
export type MessageType = 'text' | 'card' | 'image' | 'voice';

// 消息发送者
export type MessageSender = 'user' | 'assistant';

// 对话消息
export interface ConversationMessage {
  id: string; // 消息ID
  type: MessageType; // 消息类型
  content: any; // 消息内容（根据类型不同）
  timestamp: string; // ISO 8601时间戳
  sender: MessageSender; // 发送者
  sessionId?: string; // 对话会话ID
  isStreaming?: boolean; // 是否正在流式返回（仅系统消息）
}

// 对话会话
export interface ConversationSession {
  sessionId: string; // 会话ID
  messages: ConversationMessage[]; // 消息列表
  createdAt: string; // 创建时间
  lastActive: string; // 最后活跃时间
}
