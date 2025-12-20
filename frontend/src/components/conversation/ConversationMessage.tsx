/**
 * 单条消息组件
 * 支持文本、卡片、图片、语音消息，支持打字机效果
 */

import { useTranslation } from 'react-i18next';
import type { ConversationMessage } from '../../types/conversation';
import { ScienceCard } from '../cards/ScienceCard';
import { PoetryCard } from '../cards/PoetryCard';
import { EnglishCard } from '../cards/EnglishCard';
import type { KnowledgeCard } from '../../types/exploration';
import { useTypingEffect } from '../../hooks/useTypingEffect';

export interface ConversationMessageProps {
  message: ConversationMessage;
  onCollect?: (cardId: string) => void;
}

export function ConversationMessageComponent({ message, onCollect }: ConversationMessageProps) {
  const { t } = useTranslation();
  const isUser = message.sender === 'user';

  // 如果是AI消息且正在流式返回，使用打字机效果
  const displayText = !isUser && message.isStreaming && message.streamingText
    ? message.streamingText
    : typeof message.content === 'string'
    ? message.content
    : '';

  // 使用打字机效果（仅在流式返回时）
  const typingText = useTypingEffect({
    text: displayText,
    speed: 30,
    enabled: !isUser && message.isStreaming === true && !!message.streamingText,
  });

  const renderContent = () => {
    switch (message.type) {
      case 'text':
        const textToShow = !isUser && message.isStreaming && message.streamingText
          ? typingText
          : displayText;
        return (
          <div className={`px-4 py-3 rounded-2xl ${
            isUser
              ? 'bg-[var(--color-primary)] text-white rounded-br-sm'
              : 'bg-gray-100 text-gray-800 rounded-bl-sm'
          }`}>
            <p className="text-sm font-medium whitespace-pre-wrap">
              {textToShow}
              {!isUser && message.isStreaming && (
                <span className="inline-block w-2 h-4 bg-gray-600 ml-1 animate-pulse" />
              )}
            </p>
          </div>
        );

      case 'card':
        // 将content转换为KnowledgeCard格式
        const card = message.content as KnowledgeCard;
        if (card.type === 'science') {
          return <ScienceCard card={card} onCollect={onCollect} />;
        } else if (card.type === 'poetry') {
          return <PoetryCard card={card} onCollect={onCollect} />;
        } else if (card.type === 'english') {
          return <EnglishCard card={card} onCollect={onCollect} />;
        }
        return null;

      case 'image':
        return (
          <div className="rounded-2xl overflow-hidden max-w-sm">
            <img
              src={message.content}
              alt={t('conversation.imageInput')}
              className="w-full h-auto"
            />
          </div>
        );

      case 'voice':
        return (
          <div className="flex items-center gap-2 px-4 py-3 rounded-2xl bg-gray-100">
            <span className="material-symbols-outlined text-[var(--color-primary)]">
              mic
            </span>
            <p className="text-sm font-medium text-gray-800">{message.content}</p>
          </div>
        );

      default:
        return null;
    }
  };

  return (
    <div className={`flex gap-3 mb-4 ${isUser ? 'flex-row-reverse' : 'flex-row'}`}>
      {/* 头像 */}
      <div className={`size-10 rounded-full flex items-center justify-center shrink-0 ${
        isUser
          ? 'bg-[var(--color-primary)]'
          : 'bg-gray-200'
      }`}>
        {isUser ? (
          <span className="material-symbols-outlined text-white text-xl">person</span>
        ) : (
          <span className="material-symbols-outlined text-gray-600 text-xl">smart_toy</span>
        )}
      </div>

      {/* 消息内容 */}
      <div className={`flex-1 ${isUser ? 'items-end' : 'items-start'} flex flex-col`}>
        {renderContent()}
        <span className="text-xs text-gray-400 mt-1">
          {new Date(message.timestamp).toLocaleTimeString()}
        </span>
      </div>
    </div>
  );
}
