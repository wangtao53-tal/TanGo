/**
 * 科学认知卡组件
 * 绿色主题，基于设计稿
 */

import React, { useState, useEffect, useContext, createContext } from 'react';
import { cardStorage } from '../../services/storage';
import type { KnowledgeCard } from '../../types/exploration';
import type { ScienceCardContent } from '../../types/exploration';
import { useTextToSpeech } from '../../hooks/useTextToSpeech';
import { extractCardText, detectCardLanguage } from '../../utils/cardTextExtractor';
import { cardStyles } from '../../styles/cardStyles';

// 全局播放状态上下文（确保同时只播放一张卡片）
interface AudioPlaybackContextType {
  currentPlayingCardId: string | null;
  setCurrentPlayingCardId: (id: string | null) => void;
}

const AudioPlaybackContext = createContext<AudioPlaybackContextType>({
  currentPlayingCardId: null,
  setCurrentPlayingCardId: () => {},
});

export const AudioPlaybackProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [currentPlayingCardId, setCurrentPlayingCardId] = useState<string | null>(null);
  return (
    <AudioPlaybackContext.Provider value={{ currentPlayingCardId, setCurrentPlayingCardId }}>
      {children}
    </AudioPlaybackContext.Provider>
  );
};

export const usePlayingCard = () => {
  const context = useContext(AudioPlaybackContext);
  return {
    playingCardId: context.currentPlayingCardId,
    setPlayingCardId: context.setCurrentPlayingCardId,
  };
};

export interface ScienceCardProps {
  card: KnowledgeCard;
  onCollect?: (cardId: string) => void;
  className?: string;
  id?: string; // 用于导出功能
}

export const ScienceCard: React.FC<ScienceCardProps> = ({
  card,
  onCollect,
  className = '',
  id,
}) => {
  const [isCollected, setIsCollected] = useState(false);
  const content = card.content as ScienceCardContent;
  const { playingCardId, setPlayingCardId } = usePlayingCard();
  const isCurrentlyPlaying = playingCardId === card.id;

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
      className={`flex flex-col rounded-[2.5rem] bg-white border-4 border-science-green shadow-card relative transition-all hover:-translate-y-2 hover:shadow-xl hover:z-10 duration-300 group overflow-hidden w-full max-w-md mx-auto my-2 ${className}`}
    >
      {/* 内容区域（纯文本显示，不显示图片） */}
      <div className="bg-white p-6 flex flex-col justify-between relative rounded-[2.2rem] min-h-[400px]">
        <div className="flex-1 flex flex-col">
          <h3 
            className="font-bold mb-2"
            style={{ 
              fontSize: cardStyles.fonts.sizes.title,
              lineHeight: cardStyles.fonts.lineHeight.title,
              fontFamily: cardStyles.fonts.childFriendly.chinese,
              color: '#059669' // 使用更深的绿色，确保对比度≥4.5:1
            }}
          >
            {content.name || card.title}
          </h3>
          <p 
            className="font-semibold mb-4 bg-slate-50 p-3 rounded-xl border border-slate-100 italic"
            style={{ 
              fontSize: cardStyles.fonts.sizes.body,
              lineHeight: cardStyles.fonts.lineHeight.body,
              fontFamily: cardStyles.fonts.childFriendly.chinese,
              color: '#374151' // 使用更深的灰色，确保对比度≥4.5:1
            }}
          >
            "{content.explanation}"
          </p>

          {/* 关键事实 */}
          <div className="space-y-3 mb-4 flex-1">
            {content.facts?.map((fact, index) => (
              <div key={index} className="flex gap-3 items-start group/item">
                <div className="mt-1 size-2 rounded-full bg-science-green group-hover/item:scale-150 transition-transform" />
                <p 
                  className="font-bold"
                  style={{ 
                    fontSize: cardStyles.fonts.sizes.body,
                    lineHeight: cardStyles.fonts.lineHeight.body,
                    fontFamily: cardStyles.fonts.childFriendly.chinese,
                    color: '#374151' // 使用更深的灰色，确保对比度≥4.5:1
                  }}
                >
                  {fact}
                </p>
              </div>
            ))}
          </div>

          {/* 趣味知识点 */}
          {content.funFact && (
            <div className="bg-gradient-to-br from-science-green/20 to-green-100 rounded-2xl p-4 border-2 border-science-green/30 relative mt-2">
              <span className="absolute -top-3 -right-3 bg-science-green text-white rounded-full p-1 border-2 border-white shadow-sm">
                <span className="material-symbols-outlined text-sm">lightbulb</span>
              </span>
              <span 
                className="font-black uppercase tracking-wider block mb-1"
                style={{ 
                  fontSize: cardStyles.fonts.sizes.small,
                  lineHeight: cardStyles.fonts.lineHeight.small,
                  fontFamily: cardStyles.fonts.childFriendly.english,
                  color: '#059669' // 使用更深的绿色，确保对比度≥4.5:1
                }}
              >
                Fun Fact
              </span>
              <p 
                className="font-bold"
                style={{ 
                  fontSize: cardStyles.fonts.sizes.body,
                  lineHeight: cardStyles.fonts.lineHeight.body,
                  fontFamily: cardStyles.fonts.childFriendly.chinese,
                  color: '#1F2937' // 使用更深的灰色，确保对比度≥4.5:1
                }}
              >
                {content.funFact}
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
              className="flex items-center gap-2 text-sm font-bold text-science-green bg-science-green/20 hover:bg-science-green hover:text-white px-5 py-3 rounded-full transition-all duration-200 shadow-sm hover:shadow-md disabled:opacity-50 disabled:cursor-not-allowed active:scale-95"
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

