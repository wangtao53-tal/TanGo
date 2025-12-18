/**
 * 卡片轮播组件
 * 三张知识卡片横向滑动展示（snap scroll）
 * 基于 stitch_ui/recognition_result_page_1/ 设计稿
 */

import React, { useRef } from 'react';
import { ScienceCard } from './ScienceCard';
import { PoetryCard } from './PoetryCard';
import { EnglishCard } from './EnglishCard';
import type { KnowledgeCard } from '@/types/exploration';

export interface CardCarouselProps {
  cards: KnowledgeCard[];
  onCollect?: (cardId: string) => void;
  onCollectAll?: () => void;
}

export const CardCarousel: React.FC<CardCarouselProps> = ({
  cards,
  onCollect,
}) => {
  const carouselRef = useRef<HTMLDivElement>(null);

  const scrollLeft = () => {
    if (carouselRef.current) {
      carouselRef.current.scrollBy({ left: -400, behavior: 'smooth' });
    }
  };

  const scrollRight = () => {
    if (carouselRef.current) {
      carouselRef.current.scrollBy({ left: 400, behavior: 'smooth' });
    }
  };

  const scienceCard = cards.find(c => c.type === 'science');
  const poetryCard = cards.find(c => c.type === 'poetry');
  const englishCard = cards.find(c => c.type === 'english');

  return (
    <section className="relative flex-1 flex flex-col justify-center w-full min-h-[500px] mb-12 group/carousel">
      {/* 左箭头导航（PC端） */}
      <button
        onClick={scrollLeft}
        className="absolute left-0 top-1/2 -translate-y-1/2 z-20 size-14 bg-white hover:bg-science-green hover:text-white text-slate-400 border-4 border-slate-100 rounded-full flex items-center justify-center transition-all shadow-card hover:shadow-card-hover-green hidden md:flex"
      >
        <span className="material-symbols-outlined text-3xl font-bold">arrow_back</span>
      </button>

      {/* 卡片容器 */}
      <div
        ref={carouselRef}
        className="flex overflow-x-auto gap-8 pb-12 px-4 md:px-12 snap-x snap-mandatory scrollbar-thin h-full items-stretch pt-4"
      >
        {/* 科学认知卡 */}
        {scienceCard && (
          <ScienceCard
            card={scienceCard}
            onCollect={onCollect}
            className="snap-center shrink-0 w-[85vw] md:w-[380px] lg:w-[420px]"
          />
        )}

        {/* 古诗词/人文卡 */}
        {poetryCard && (
          <PoetryCard
            card={poetryCard}
            onCollect={onCollect}
            className="snap-center shrink-0 w-[85vw] md:w-[380px] lg:w-[420px]"
          />
        )}

        {/* 英语表达卡 */}
        {englishCard && (
          <EnglishCard
            card={englishCard}
            onCollect={onCollect}
            className="snap-center shrink-0 w-[85vw] md:w-[380px] lg:w-[420px]"
          />
        )}
      </div>

      {/* 右箭头导航（PC端） */}
      <button
        onClick={scrollRight}
        className="absolute right-0 top-1/2 -translate-y-1/2 z-20 size-14 bg-white hover:bg-science-green hover:text-white text-slate-400 border-4 border-slate-100 rounded-full flex items-center justify-center transition-all shadow-card hover:shadow-card-hover-green hidden md:flex"
      >
        <span className="material-symbols-outlined text-3xl font-bold">arrow_forward</span>
      </button>
    </section>
  );
};

