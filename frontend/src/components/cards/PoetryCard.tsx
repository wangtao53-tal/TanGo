/**
 * 古诗词/人文卡组件
 * 橙色主题，基于设计稿
 */

import React, { useState, useEffect } from 'react';
import { cardStorage } from '../../services/storage';
import type { KnowledgeCard } from '../../types/exploration';
import type { PoetryCardContent } from '../../types/exploration';

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
      className={`flex flex-col rounded-[2.5rem] bg-white border-4 border-sunny-orange shadow-card relative transition-transform hover:-translate-y-2 duration-300 group overflow-hidden w-full max-w-md mx-auto ${className}`}
    >
      {/* 内容区域（纯文本显示，不显示图片） */}
      <div className="bg-white p-6 flex flex-col justify-between relative rounded-[2.2rem]">
        <div className="flex-1 overflow-y-auto pr-2 scrollbar-thin">
          <h3 className="text-3xl font-display font-bold text-sunny-orange mb-3">
            {card.title}
          </h3>

          {/* 诗词内容 */}
          {content.poem && (
            <div className="relative pl-5 border-l-4 border-sunny-orange/40 mb-4 bg-orange-50 p-3 rounded-r-xl">
              <p className="text-slate-700 text-lg font-serif italic">"{content.poem}"</p>
              {content.author && (
                <p className="text-xs text-slate-500 mt-1">— {content.author}</p>
              )}
            </div>
          )}

          {/* 解释 */}
          <div className="mb-4">
            <h4 className="text-xs font-black text-sunny-orange uppercase mb-1 tracking-wider">
              What it means
            </h4>
            <p className="text-sm font-bold text-slate-600 leading-relaxed">
              {content.explanation}
            </p>
          </div>

          {/* 情境联想 */}
          {content.context && (
            <div className="flex items-center gap-3 bg-sunny-orange/10 p-3 rounded-2xl border-2 border-sunny-orange/20">
              <div className="bg-white p-1 rounded-full shadow-sm text-sunny-orange">
                <span className="material-symbols-outlined text-lg">emoji_nature</span>
              </div>
              <p className="text-xs font-bold text-slate-700">{content.context}</p>
            </div>
          )}
        </div>

        {/* 操作按钮 */}
        <div className="flex justify-between items-center mt-4 pt-2">
          <button className="flex items-center gap-2 text-sm font-bold text-sunny-orange bg-sunny-orange/20 hover:bg-sunny-orange hover:text-white px-5 py-3 rounded-full transition-all shadow-sm hover:shadow-md">
            <span className="material-symbols-outlined text-xl">play_circle</span>
            Listen
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

