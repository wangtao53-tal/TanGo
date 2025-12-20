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
    setIsCollected(newCollectedState);
    
    if (onCollect) {
      onCollect(card.id);
    }

    // 立即保存到本地存储
    if (newCollectedState) {
      const cardToSave = {
        ...card,
        collectedAt: new Date().toISOString(),
      };
      await cardStorage.save(cardToSave);
    } else {
      await cardStorage.delete(card.id);
    }
  };

  return (
    <article
      id={id || `card-${card.id}`}
      className={`flex flex-col rounded-[2.5rem] bg-white border-4 border-science-green shadow-card relative transition-transform hover:-translate-y-2 duration-300 group overflow-hidden w-full max-w-md mx-auto ${className}`}
    >
      {/* 内容区域（纯文本显示，不显示图片） */}
      <div className="bg-white p-6 flex flex-col justify-between relative rounded-[2.2rem]">
        <div className="flex-1 overflow-y-auto pr-2 scrollbar-thin">
          <h3 className="text-3xl font-display font-bold text-science-green mb-2">
            {content.name || card.title}
          </h3>
          <p className="text-base font-semibold text-slate-500 mb-4 bg-slate-50 p-3 rounded-xl border border-slate-100 italic">
            "{content.explanation}"
          </p>

          {/* 关键事实 */}
          <div className="space-y-3 mb-4">
            {content.facts?.map((fact, index) => (
              <div key={index} className="flex gap-3 items-start group/item">
                <div className="mt-1 size-2 rounded-full bg-science-green group-hover/item:scale-150 transition-transform" />
                <p className="text-sm font-bold text-slate-600 leading-snug">{fact}</p>
              </div>
            ))}
          </div>

          {/* 趣味知识点 */}
          {content.funFact && (
            <div className="bg-gradient-to-br from-science-green/20 to-green-100 rounded-2xl p-4 border-2 border-science-green/30 relative mt-2">
              <span className="absolute -top-3 -right-3 bg-science-green text-white rounded-full p-1 border-2 border-white shadow-sm">
                <span className="material-symbols-outlined text-sm">lightbulb</span>
              </span>
              <span className="text-xs font-black text-science-green uppercase tracking-wider block mb-1">
                Fun Fact
              </span>
              <p className="text-sm font-bold text-slate-700">{content.funFact}</p>
            </div>
          )}
        </div>

        {/* 操作按钮 */}
        <div className="flex justify-between items-center mt-4 pt-2">
          <div className="flex items-center gap-2">
            <button
              onClick={handleListen}
              disabled={!isSupported}
              className="flex items-center gap-2 text-sm font-bold text-science-green bg-science-green/20 hover:bg-science-green hover:text-white px-5 py-3 rounded-full transition-all shadow-sm hover:shadow-md disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <span className="material-symbols-outlined text-xl">
                {isPlaying && !isPaused ? 'pause' : 'play_circle'}
              </span>
              {isPlaying && !isPaused ? '暂停' : isPaused ? '继续' : '听'}
            </button>
            {isPlaying && (
              <button
                onClick={handleStop}
                className="flex items-center gap-1 text-xs font-bold text-slate-500 hover:text-slate-700 px-3 py-2 rounded-full transition-all"
              >
                <span className="material-symbols-outlined text-lg">stop</span>
                停止
              </button>
            )}
          </div>
          <button
            onClick={handleCollect}
            className={`size-12 rounded-full flex items-center justify-center transition-all group-active:scale-90 shadow-sm border-2 ${
              isCollected
                ? 'bg-yellow-300 text-yellow-800 border-yellow-400'
                : 'bg-slate-100 text-slate-300 border-slate-200 hover:border-yellow-400'
            }`}
          >
            <span className="material-symbols-outlined text-2xl fill-current">star</span>
          </button>
        </div>
      </div>
    </article>
  );
};

