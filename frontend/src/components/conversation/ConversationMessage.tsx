/**
 * 单条消息组件
 * 支持文本、卡片、图片、语音消息，支持打字机效果
 */

import { memo } from 'react';
import { useTranslation } from 'react-i18next';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
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

function ConversationMessageComponentInner({ message, onCollect }: ConversationMessageProps) {
  const { t } = useTranslation();
  const isUser = message.sender === 'user';

  // 如果是AI消息且正在流式返回，优先使用streamingText（实时渲染）
  // 否则使用content
  const displayText = !isUser && message.isStreaming && message.streamingText !== undefined
    ? message.streamingText
    : typeof message.content === 'string'
    ? message.content
    : '';

  // 使用打字机效果（仅在流式返回时，且streamingText存在）
  // 注意：由于我们已经使用flushSync实时更新streamingText，打字机效果可能不需要
  // 但保留作为可选效果
  const typingText = useTypingEffect({
    text: displayText,
    speed: 30,
    enabled: false, // 禁用打字机效果，直接显示streamingText以实现真正的实时渲染
  });

  const renderContent = () => {
    switch (message.type) {
      case 'text':
        // 流式返回时直接显示streamingText，不使用打字机效果
        const textToShow = !isUser && message.isStreaming && message.streamingText !== undefined
          ? message.streamingText
          : displayText;
        
        // 检测是否为Markdown格式（通过markdown字段或内容检测）
        const isMarkdown = message.markdown || 
          (typeof textToShow === 'string' && (
            textToShow.includes('```') || // 代码块
            textToShow.includes('##') || // 标题
            textToShow.includes('- ') || // 列表
            textToShow.includes('* ') || // 列表
            textToShow.includes('[') && textToShow.includes('](') // 链接
          ));
        
        return (
          <div className={`px-4 py-3 rounded-2xl ${
            isUser
              ? 'bg-[var(--color-primary)] text-white rounded-br-sm'
              : 'bg-gray-100 text-gray-800 rounded-bl-sm'
          }`}>
            {isMarkdown ? (
              <div className="text-sm font-medium prose prose-sm max-w-none dark:prose-invert">
                <ReactMarkdown
                  remarkPlugins={[remarkGfm]}
                  components={{
                    // 自定义代码块样式
                    code: ({ node, inline, className, children, ...props }) => {
                      return inline ? (
                        <code className="bg-gray-200 px-1 py-0.5 rounded text-xs" {...props}>
                          {children}
                        </code>
                      ) : (
                        <pre className="bg-gray-200 p-3 rounded-lg overflow-x-auto my-2">
                          <code className="text-xs" {...props}>
                            {children}
                          </code>
                        </pre>
                      );
                    },
                    // 自定义链接样式
                    a: ({ node, ...props }) => (
                      <a className="text-blue-600 hover:text-blue-800 underline" {...props} />
                    ),
                    // 自定义列表样式
                    ul: ({ node, ...props }) => (
                      <ul className="list-disc list-inside my-2 space-y-1" {...props} />
                    ),
                    ol: ({ node, ...props }) => (
                      <ol className="list-decimal list-inside my-2 space-y-1" {...props} />
                    ),
                    // 自定义标题样式
                    h1: ({ node, ...props }) => (
                      <h1 className="text-lg font-bold mt-3 mb-2" {...props} />
                    ),
                    h2: ({ node, ...props }) => (
                      <h2 className="text-base font-bold mt-2 mb-1" {...props} />
                    ),
                    h3: ({ node, ...props }) => (
                      <h3 className="text-sm font-bold mt-2 mb-1" {...props} />
                    ),
                  }}
                >
                  {textToShow}
                </ReactMarkdown>
                {!isUser && message.isStreaming && (
                  <span className="inline-block w-2 h-4 bg-gray-600 ml-1 animate-pulse" />
                )}
              </div>
            ) : (
              <p className="text-sm font-medium whitespace-pre-wrap">
                {textToShow}
                {!isUser && message.isStreaming && (
                  <span className="inline-block w-2 h-4 bg-gray-600 ml-1 animate-pulse" />
                )}
              </p>
            )}
          </div>
        );

      case 'card':
        // 将content转换为KnowledgeCard格式
        const card = message.content as KnowledgeCard;
        
        // 验证卡片数据
        if (!card || !card.type || !card.content) {
          console.error('卡片数据无效:', card);
          return (
            <div className="px-4 py-3 rounded-2xl bg-red-50 text-red-600 border border-red-200">
              <p className="text-sm">卡片数据格式错误，无法显示</p>
            </div>
          );
        }

        if (card.type === 'science') {
          return <ScienceCard card={card} onCollect={onCollect} />;
        } else if (card.type === 'poetry') {
          return <PoetryCard card={card} onCollect={onCollect} />;
        } else if (card.type === 'english') {
          return <EnglishCard card={card} onCollect={onCollect} />;
        }
        
        // 未知卡片类型
        console.warn('未知的卡片类型:', card.type);
        return (
          <div className="px-4 py-3 rounded-2xl bg-yellow-50 text-yellow-600 border border-yellow-200">
            <p className="text-sm">未知的卡片类型: {card.type}</p>
          </div>
        );

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
          ? 'bg-gradient-to-br from-primary to-[#5aff2b] shadow-md'
          : 'bg-gradient-to-br from-sky-blue to-blue-400 shadow-md'
      }`}>
        {isUser ? (
          <span className="material-symbols-outlined text-white text-xl">face</span>
        ) : (
          <span className="material-symbols-outlined text-white text-xl">auto_awesome</span>
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

// 使用React.memo优化性能，避免不必要的重新渲染
export const ConversationMessageComponent = memo(ConversationMessageComponentInner, (prevProps, nextProps) => {
  // 自定义比较函数：只有当消息内容真正改变时才重新渲染
  if (prevProps.message.id !== nextProps.message.id) {
    return false; // 不同消息，需要重新渲染
  }
  
  // 相同消息，比较关键字段
  if (prevProps.message.isStreaming !== nextProps.message.isStreaming) {
    return false; // 流式状态改变，需要重新渲染
  }
  
  if (prevProps.message.streamingText !== nextProps.message.streamingText) {
    return false; // 流式文本改变，需要重新渲染
  }
  
  if (prevProps.message.content !== nextProps.message.content) {
    return false; // 内容改变，需要重新渲染
  }
  
  if (prevProps.message.markdown !== nextProps.message.markdown) {
    return false; // Markdown标记改变，需要重新渲染
  }
  
  // 其他字段相同，可以跳过渲染
  return true;
});
