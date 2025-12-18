/**
 * 对话消息列表组件
 * 显示所有对话消息，支持流式更新
 */

import { useEffect, useRef } from 'react';
import { ConversationMessageComponent } from './ConversationMessage';
import type { ConversationMessage } from '../../types/conversation';

export interface ConversationListProps {
  messages: ConversationMessage[];
  onCollect?: (cardId: string) => void;
}

export function ConversationList({ messages, onCollect }: ConversationListProps) {
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const listRef = useRef<HTMLDivElement>(null);

  // 自动滚动到底部
  useEffect(() => {
    if (messagesEndRef.current) {
      messagesEndRef.current.scrollIntoView({ behavior: 'smooth' });
    }
  }, [messages]);

  if (messages.length === 0) {
    return (
      <div className="flex-1 flex items-center justify-center text-gray-400">
        <div className="text-center">
          <span className="material-symbols-outlined text-6xl mb-4">chat_bubble_outline</span>
          <p className="text-lg">开始对话吧～</p>
        </div>
      </div>
    );
  }

  return (
    <div ref={listRef} className="flex-1 overflow-y-auto px-4 py-4 space-y-2">
      {messages.map((message) => (
        <ConversationMessageComponent
          key={message.id}
          message={message}
          onCollect={onCollect}
        />
      ))}
      <div ref={messagesEndRef} />
    </div>
  );
}
