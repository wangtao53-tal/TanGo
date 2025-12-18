/**
 * 英语表达卡组件
 * 蓝色主题，基于设计稿
 */

import React, { useState, useEffect } from 'react';
import { cardStorage } from '../../services/storage';
import type { KnowledgeCard } from '../../types/exploration';
import type { EnglishCardContent } from '../../types/exploration';

export interface EnglishCardProps {
  card: KnowledgeCard;
  onCollect?: (cardId: string) => void;
  className?: string;
  id?: string; // 用于导出功能
}

export const EnglishCard: React.FC<EnglishCardProps> = ({
  card,
  onCollect,
  className = '',
  id,
}) => {
  const [isCollected, setIsCollected] = useState(false);
  const [isPlaying, setIsPlaying] = useState(false);
  const content = card.content as EnglishCardContent;

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

  const handlePlay = () => {
    setIsPlaying(!isPlaying);
    // TODO: 实现音频播放功能
  };

  return (
    <article
      id={id || `card-${card.id}`}
      className={`flex flex-col rounded-[2.5rem] bg-white border-4 border-sky-blue shadow-card relative transition-transform hover:-translate-y-2 duration-300 group overflow-hidden w-full max-w-md mx-auto ${className}`}
      style={{ minHeight: '600px', maxHeight: '800px' }}
    >
      {/* 顶部图片区域（45%高度） */}
      <div className="h-[45%] w-full relative overflow-hidden rounded-t-[2.2rem]">
        <div
          className="absolute inset-0 bg-cover bg-center transform hover:scale-110 transition-transform duration-700"
          style={{
            backgroundImage: `url(https://images.unsplash.com/photo-1456513080510-7bf3a84b82f8?w=800)`,
          }}
        />
        <div className="absolute top-4 left-4 bg-white/90 backdrop-blur-sm px-4 py-2 rounded-2xl text-sky-blue font-display font-bold text-sm border-2 border-sky-blue shadow-sm flex items-center gap-2">
          <span className="material-symbols-outlined text-lg">translate</span>
          Say it in English
        </div>
        <div
          className="absolute bottom-0 left-0 w-full h-8 bg-white"
          style={{ clipPath: 'polygon(0 100%, 100% 100%, 100% 0, 0 100%)' }}
        />
      </div>

      {/* 底部内容区域（55%高度） */}
      <div className="h-[55%] bg-white p-6 flex flex-col justify-between relative rounded-b-[2.2rem]">
        <div className="flex-1 overflow-y-auto pr-2 scrollbar-thin">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-3xl font-display font-bold text-sky-blue">
              {card.title}
            </h3>
            <button
              onClick={handlePlay}
              className={`size-10 rounded-full flex items-center justify-center transition-colors shadow-sm ${
                isPlaying
                  ? 'bg-sky-blue text-white'
                  : 'bg-sky-blue/20 hover:bg-sky-blue hover:text-white text-sky-blue'
              }`}
            >
              <span className="material-symbols-outlined text-xl">volume_up</span>
            </button>
          </div>

          <div className="space-y-4">
            {/* 核心单词 */}
            {content.words && content.words.length > 0 && (
              <div>
                <h4 className="text-xs font-black text-sky-blue uppercase mb-2 tracking-wider">
                  Magic Words
                </h4>
                <div className="flex gap-2 flex-wrap">
                  {content.words.map((word, index) => (
                    <span
                      key={index}
                      className="px-3 py-1.5 bg-sky-blue/10 rounded-xl text-xs font-bold text-sky-blue border border-sky-blue/20 hover:bg-sky-blue hover:text-white cursor-pointer transition-colors"
                    >
                      {word}
                    </span>
                  ))}
                </div>
              </div>
            )}

            {/* 口语表达 */}
            {content.expressions && content.expressions.length > 0 && (
              <div>
                <h4 className="text-xs font-black text-sky-blue uppercase mb-2 tracking-wider">
                  Let's Talk!
                </h4>
                <div className="bg-sky-50 p-4 rounded-2xl border-2 border-sky-blue/20 relative">
                  <div className="absolute -top-3 right-4 bg-sky-200 text-sky-800 text-[10px] font-bold px-2 py-0.5 rounded-full">
                    Try saying this
                  </div>
                  {content.expressions.map((expr, index) => (
                    <p key={index} className="text-sm text-slate-700 mb-2 font-bold">
                      "{expr}"
                    </p>
                  ))}
                </div>
              </div>
            )}
          </div>
        </div>

        {/* 操作按钮 */}
        <div className="flex justify-between items-center mt-4 pt-2">
          <button
            onClick={handlePlay}
            className="flex items-center gap-2 text-sm font-bold text-sky-blue bg-sky-blue/20 hover:bg-sky-blue hover:text-white px-5 py-3 rounded-full transition-all shadow-sm hover:shadow-md"
          >
            <span className="material-symbols-outlined text-xl">mic</span>
            Practice
          </button>
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

