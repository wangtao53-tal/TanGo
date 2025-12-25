/**
 * 对话消息列表组件
 * 显示所有对话消息，支持流式更新
 */

import { useEffect, useRef, useMemo } from 'react';
import { ConversationMessageComponent } from './ConversationMessage';
import { CardCarousel } from '../cards/CardCarousel';
import { exportCardAsImage } from '../../utils/cardExport';
import type { ConversationMessage } from '../../types/conversation';
import type { KnowledgeCard } from '../../types/exploration';

export interface ConversationListProps {
  messages: ConversationMessage[];
  onCollect?: (cardId: string) => void;
  onExport?: (cardId: string) => void;
}

export function ConversationList({ messages, onCollect, onExport }: ConversationListProps) {
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const listRef = useRef<HTMLDivElement>(null);

  // 自动滚动到底部
  useEffect(() => {
    if (messagesEndRef.current) {
      messagesEndRef.current.scrollIntoView({ behavior: 'smooth' });
    }
  }, [messages]);

  // 检测连续的卡片消息并分组
  const processedMessages = useMemo(() => {
    const result: Array<{ type: 'message' | 'cards'; data: ConversationMessage | ConversationMessage[] }> = [];
    let i = 0;

    while (i < messages.length) {
      const message = messages[i];

      // 检测卡片消息
      if (message.type === 'card' && message.sender === 'assistant') {
        // 收集连续的卡片消息（最多3张）
        const cardMessages: ConversationMessage[] = [];
        let j = i;

        while (j < messages.length && messages[j].type === 'card' && messages[j].sender === 'assistant' && cardMessages.length < 3) {
          cardMessages.push(messages[j]);
          j++;
        }

        // 如果有多张卡片，使用CardCarousel
        if (cardMessages.length > 0) {
          result.push({ type: 'cards', data: cardMessages });
          i = j;
          continue;
        }
      }

      // 普通消息
      result.push({ type: 'message', data: message });
      i++;
    }

    return result;
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
    <div ref={listRef} className="flex-1 overflow-y-auto overflow-x-hidden px-4 py-4 space-y-2 w-full">
      {processedMessages.map((item) => {
        if (item.type === 'cards') {
          // 卡片组：使用CardCarousel
          const cardMessages = item.data as ConversationMessage[];
          const cards: KnowledgeCard[] = cardMessages.map(msg => msg.content as KnowledgeCard);

          return (
            <div key={`cards-${cardMessages[0].id}`} className="flex gap-3 mb-4 pt-2">
              {/* 头像 */}
              <div className="size-10 rounded-full flex items-center justify-center shrink-0 bg-gradient-to-br from-sky-blue to-blue-400 shadow-md">
                <span className="material-symbols-outlined text-white text-xl">auto_awesome</span>
              </div>

              {/* 卡片轮播 */}
              <div className="flex-1">
                <CardCarousel 
                  cards={cards} 
                  onCollect={onCollect}
                  onExport={onExport || (async (cardId: string) => {
                    try {
                      await exportCardAsImage(`card-${cardId}`, { filename: `card-${cardId}-${Date.now()}` });
                    } catch (error) {
                      console.error('导出卡片失败:', error);
                    }
                  })}
                />
                <span className="text-xs text-gray-400 mt-1 block">
                  {new Date(cardMessages[0].timestamp).toLocaleTimeString()}
                </span>
              </div>
            </div>
          );
        } else {
          // 普通消息
          const message = item.data as ConversationMessage;
          return (
            <ConversationMessageComponent
              key={message.id}
              message={message}
              onCollect={onCollect}
            />
          );
        }
      })}
      <div ref={messagesEndRef} />
    </div>
  );
}
