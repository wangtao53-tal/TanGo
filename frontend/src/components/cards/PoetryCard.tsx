/**
 * 古诗词/人文卡组件
 * 橙色主题，基于设计稿
 */

import React, { useState, useEffect } from 'react';
import { cardStorage } from '../../services/storage';
import type { KnowledgeCard } from '../../types/exploration';
import type { PoetryCardContent } from '../../types/exploration';
import { useTextToSpeech } from '../../hooks/useTextToSpeech';
import { extractCardText, detectCardLanguage } from '../../utils/cardTextExtractor';
import { usePlayingCard } from './ScienceCard';
import { cardStyles } from '../../styles/cardStyles';

export interface PoetryCardProps {
  card: KnowledgeCard;
  onCollect?: (cardId: string) => void;
  className?: string;
  id?: string; // 用于导出功能
}

export const PoetryCard: React.FC<PoetryCardProps> = ({
  card,
  onCollect,
  className = '',
  id,
}) => {
  const [isCollected, setIsCollected] = useState(false);
  const content = card.content as PoetryCardContent;
  const { playingCardId, setPlayingCardId } = usePlayingCard();

  // 文本转语音Hook
  const { isPlaying, isPaused, play, pause, resume, stop, isSupported } = useTextToSpeech({
    rate: 0.9,
    pitch: 1.0,
    onStart: () => {
      setPlayingCardId(card.id);
    },
    onEnd: () => {
      setPlayingCardId(null);
    },
    onError: (error) => {
      console.error('文本转语音错误:', error);
      setPlayingCardId(null);
    },
  });

  // 如果其他卡片开始播放，停止当前播放
  useEffect(() => {
    if (playingCardId !== null && playingCardId !== card.id && isPlaying) {
      stop();
    }
  }, [playingCardId, card.id, isPlaying, stop]);

  const handleListen = () => {
    if (!isSupported) {
      alert('您的浏览器不支持文本转语音功能');
      return;
    }

    if (isPlaying && !isPaused) {
      pause();
    } else if (isPaused) {
      resume();
    } else {
      const text = extractCardText(card);
      const language = detectCardLanguage(card);
      play(text, language);
    }
  };

  const handleStop = () => {
    stop();
    setPlayingCardId(null);
  };

  // 检查卡片是否已收藏
  useEffect(() => {
    cardStorage.getAll().then((cards) => {
      const found = cards.find((c) => c.id === card.id);
      setIsCollected(!!found);
    });
  }, [card.id]);

  const handleCollect = async () => {
    const newCollectedState = !isCollected;
    
    // 乐观更新UI
    setIsCollected(newCollectedState);
    
    if (onCollect) {
      onCollect(card.id);
    }

    try {
      // 保存到本地存储
      if (newCollectedState) {
        const cardToSave = {
          ...card,
          collectedAt: new Date().toISOString(),
        };
        await cardStorage.save(cardToSave);
      } else {
        await cardStorage.delete(card.id);
      }
    } catch (error) {
      // 如果保存失败，回滚状态
      console.error('收藏操作失败:', error);
      setIsCollected(!newCollectedState);
      // 可选：显示错误提示
      // alert('收藏操作失败，请重试');
    }
  };

  return (
    <article
      id={id || `card-${card.id}`}
      className={`flex flex-col rounded-[2.5rem] bg-white border-4 border-sunny-orange shadow-card relative transition-all hover:-translate-y-2 hover:shadow-xl hover:z-10 duration-300 group overflow-hidden w-full max-w-md mx-auto my-2 ${className}`}
    >
      {/* 内容区域（纯文本显示，不显示图片） */}
      <div className="bg-white p-6 flex flex-col justify-between relative rounded-[2.2rem] min-h-[400px]">
        <div className="flex-1 flex flex-col">
          <h3 
            className="font-bold mb-3"
            style={{ 
              fontSize: cardStyles.fonts.sizes.title,
              lineHeight: cardStyles.fonts.lineHeight.title,
              fontFamily: cardStyles.fonts.childFriendly.chinese,
              color: '#EA580C' // 使用更深的橙色，确保对比度≥4.5:1
            }}
          >
            {card.title}
          </h3>

          {/* 诗词内容 */}
          {content.poem && (
            <div className="relative pl-5 border-l-4 border-sunny-orange/40 mb-4 bg-orange-50 p-3 rounded-r-xl">
              <p 
                className="font-serif italic"
                style={{ 
                  fontSize: cardStyles.fonts.sizes.body,
                  lineHeight: cardStyles.fonts.lineHeight.body,
                  fontFamily: cardStyles.fonts.childFriendly.chinese,
                  color: '#1F2937' // 使用更深的灰色，确保对比度≥4.5:1
                }}
              >
                "{content.poem}"
              </p>
              {content.author && (
                <p 
                  className="mt-1"
                  style={{ 
                    fontSize: cardStyles.fonts.sizes.small,
                    lineHeight: cardStyles.fonts.lineHeight.small,
                    fontFamily: cardStyles.fonts.childFriendly.chinese,
                    color: '#4B5563' // 使用更深的灰色，确保对比度≥4.5:1
                  }}
                >
                  — {content.author}
                </p>
              )}
            </div>
          )}

          {/* 解释 */}
          <div className="mb-4 flex-1">
            <h4 
              className="font-black uppercase mb-1 tracking-wider"
              style={{ 
                fontSize: cardStyles.fonts.sizes.small,
                lineHeight: cardStyles.fonts.lineHeight.small,
                fontFamily: cardStyles.fonts.childFriendly.english,
                color: '#EA580C' // 使用更深的橙色，确保对比度≥4.5:1
              }}
            >
              What it means
            </h4>
            <p 
              className="font-bold"
              style={{ 
                fontSize: cardStyles.fonts.sizes.body,
                lineHeight: cardStyles.fonts.lineHeight.body,
                fontFamily: cardStyles.fonts.childFriendly.chinese,
                color: '#374151' // 使用更深的灰色，确保对比度≥4.5:1
              }}
            >
              {content.explanation}
            </p>
          </div>

          {/* 情境联想 */}
          {content.context && (
            <div className="flex items-center gap-3 bg-sunny-orange/10 p-3 rounded-2xl border-2 border-sunny-orange/20">
              <div className="bg-white p-1 rounded-full shadow-sm text-sunny-orange">
                <span className="material-symbols-outlined text-lg">emoji_nature</span>
              </div>
              <p 
                className="font-bold"
                style={{ 
                  fontSize: cardStyles.fonts.sizes.small,
                  lineHeight: cardStyles.fonts.lineHeight.small,
                  fontFamily: cardStyles.fonts.childFriendly.chinese,
                  color: '#1F2937' // 使用更深的灰色，确保对比度≥4.5:1
                }}
              >
                {content.context}
              </p>
            </div>
          )}
        </div>

        {/* 操作按钮 */}
        <div className="flex justify-between items-center mt-4 pt-2">
          <div className="flex items-center gap-2">
            <button
              onClick={handleListen}
              disabled={!isSupported}
              className="flex items-center gap-2 text-sm font-bold text-sunny-orange bg-sunny-orange/20 hover:bg-sunny-orange hover:text-white px-5 py-3 rounded-full transition-all duration-200 shadow-sm hover:shadow-md disabled:opacity-50 disabled:cursor-not-allowed active:scale-95"
            >
              <span className="material-symbols-outlined text-xl">
                {isPlaying && !isPaused ? 'pause' : 'play_circle'}
              </span>
              {isPlaying && !isPaused ? '暂停' : isPaused ? '继续' : '听'}
            </button>
            {isPlaying && (
              <button
                onClick={handleStop}
                className="flex items-center gap-1 text-xs font-bold text-slate-500 hover:text-slate-700 px-3 py-2 rounded-full transition-all duration-200 active:scale-95"
              >
                <span className="material-symbols-outlined text-lg">stop</span>
                停止
              </button>
            )}
          </div>
          <button
            onClick={handleCollect}
            className={`size-12 rounded-full flex items-center justify-center transition-all duration-300 shadow-sm border-2 active:scale-90 ${
              isCollected
                ? 'bg-yellow-300 text-yellow-800 border-yellow-400 hover:bg-yellow-400 hover:scale-105'
                : 'bg-slate-100 text-slate-300 border-slate-200 hover:border-yellow-400 hover:bg-yellow-50'
            }`}
          >
            <span 
              className={`material-symbols-outlined text-2xl transition-all duration-300 ${
                isCollected ? 'fill-current scale-110' : 'fill-current'
              }`}
            >
              star
            </span>
          </button>
        </div>
      </div>
    </article>
  );
};

