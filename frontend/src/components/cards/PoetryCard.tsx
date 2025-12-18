/**
 * 古诗词/人文卡组件
 * 橙色主题，基于设计稿
 */

import React, { useState } from 'react';
import type { KnowledgeCard } from '@/types/exploration';
import type { PoetryCardContent } from '@/types/exploration';

export interface PoetryCardProps {
  card: KnowledgeCard;
  onCollect?: (cardId: string) => void;
  className?: string;
}

export const PoetryCard: React.FC<PoetryCardProps> = ({
  card,
  onCollect,
  className = '',
}) => {
  const [isCollected, setIsCollected] = useState(false);
  const content = card.content as PoetryCardContent;

  const handleCollect = () => {
    setIsCollected(!isCollected);
    if (onCollect) {
      onCollect(card.id);
    }
  };

  return (
    <article
      className={`flex flex-col rounded-[2.5rem] bg-white border-4 border-sunny-orange shadow-card relative transition-transform hover:-translate-y-2 duration-300 group overflow-hidden ${className}`}
    >
      {/* 顶部图片区域（45%高度） */}
      <div className="h-[45%] w-full relative overflow-hidden rounded-t-[2.2rem]">
        <div
          className="absolute inset-0 bg-cover bg-center transform hover:scale-110 transition-transform duration-700"
          style={{
            backgroundImage: `url(https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=800)`,
          }}
        />
        <div className="absolute top-4 left-4 bg-white/90 backdrop-blur-sm px-4 py-2 rounded-2xl text-sunny-orange font-display font-bold text-sm border-2 border-sunny-orange shadow-sm flex items-center gap-2">
          <span className="material-symbols-outlined text-lg">history_edu</span>
          Story Time
        </div>
        <div
          className="absolute bottom-0 left-0 w-full h-8 bg-white"
          style={{ clipPath: 'polygon(0 100%, 100% 100%, 100% 0, 0 100%)' }}
        />
      </div>

      {/* 底部内容区域（55%高度） */}
      <div className="h-[55%] bg-white p-6 flex flex-col justify-between relative rounded-b-[2.2rem]">
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

