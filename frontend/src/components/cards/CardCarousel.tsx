/**
 * 卡片轮播容器组件
 * 实现渐进式展示和滑动切换
 */

import React, { useState, useEffect, useMemo } from 'react';
import { useSwipeable } from 'react-swipeable';
import type { KnowledgeCard } from '../../types/exploration';
import { ScienceCard } from './ScienceCard';
import { PoetryCard } from './PoetryCard';
import { EnglishCard } from './EnglishCard';

export interface CardCarouselProps {
  cards: KnowledgeCard[];  // 三张卡片数据（可能部分未生成）
  onCollect?: (cardId: string) => void;   // 收藏回调
  onExport?: (cardId: string) => void;     // 导出回调
}

export const CardCarousel: React.FC<CardCarouselProps> = ({
  cards,
  onCollect,
  onExport,
}) => {
  const [currentIndex, setCurrentIndex] = useState(0);
  const [isMobile, setIsMobile] = useState(false);

  // 检测是否为移动端
  useEffect(() => {
    const checkMobile = () => {
      setIsMobile(window.innerWidth < 1024);
    };
    checkMobile();
    window.addEventListener('resize', checkMobile);
    return () => window.removeEventListener('resize', checkMobile);
  }, []);

  // 按类型排序卡片：science -> poetry -> english（第一张优先显示）
  const sortedCards = useMemo(() => {
    const cardMap = new Map<string, KnowledgeCard>();
    cards.forEach(card => {
      // 如果已存在同类型卡片，保留第一个（通常是先生成的）
      if (!cardMap.has(card.type)) {
        cardMap.set(card.type, card);
      }
    });
    
    const result: (KnowledgeCard | null)[] = [
      cardMap.get('science') || null,
      cardMap.get('poetry') || null,
      cardMap.get('english') || null,
    ];
    
    return result;
  }, [cards]);

  // 只显示已生成的卡片
  const availableCards = useMemo(() => {
    return sortedCards.filter(card => card !== null) as KnowledgeCard[];
  }, [sortedCards]);

  // 确保currentIndex在有效范围内
  useEffect(() => {
    if (availableCards.length > 0 && currentIndex >= availableCards.length) {
      setCurrentIndex(0);
    }
  }, [availableCards.length, currentIndex]);

  const [isTransitioning, setIsTransitioning] = useState(false);

  // 滑动处理
  const handlers = useSwipeable({
    onSwipedLeft: () => {
      // 左滑：切换到下一张
      if (currentIndex < availableCards.length - 1 && !isTransitioning) {
        setIsTransitioning(true);
        setCurrentIndex(prev => prev + 1);
        setTimeout(() => setIsTransitioning(false), 300);
      }
    },
    onSwipedRight: () => {
      // 右滑：切换到上一张
      if (currentIndex > 0 && !isTransitioning) {
        setIsTransitioning(true);
        setCurrentIndex(prev => prev - 1);
        setTimeout(() => setIsTransitioning(false), 300);
      }
    },
    trackMouse: !isMobile, // PC端也支持鼠标拖动
    preventScrollOnSwipe: true,
    delta: 50, // 最小滑动距离
  });

  // PC端切换函数
  const handlePrev = () => {
    if (currentIndex > 0 && !isTransitioning) {
      setIsTransitioning(true);
      setCurrentIndex(prev => prev - 1);
      setTimeout(() => setIsTransitioning(false), 300);
    }
  };

  const handleNext = () => {
    if (currentIndex < availableCards.length - 1 && !isTransitioning) {
      setIsTransitioning(true);
      setCurrentIndex(prev => prev + 1);
      setTimeout(() => setIsTransitioning(false), 300);
    }
  };

  // 导出当前卡片
  const handleExport = () => {
    if (currentCard && onExport) {
      onExport(currentCard.id);
    }
  };

  // 如果没有卡片，显示加载状态
  if (availableCards.length === 0) {
    return (
      <div className="w-full max-w-md mx-auto min-h-[400px] flex items-center justify-center">
        <div className="text-center text-gray-400">
          <span className="material-symbols-outlined text-6xl mb-4 animate-pulse">
            auto_awesome
          </span>
          <p className="text-lg">正在生成卡片...</p>
        </div>
      </div>
    );
  }

  const currentCard = availableCards[currentIndex];

  // 渲染卡片
  const renderCard = (card: KnowledgeCard) => {
    const commonProps = {
      card,
      onCollect,
      className: 'w-full',
      id: `card-${card.id}`, // 确保ID正确传递，用于导出功能
    };

    switch (card.type) {
      case 'science':
        return <ScienceCard {...commonProps} />;
      case 'poetry':
        return <PoetryCard {...commonProps} />;
      case 'english':
        return <EnglishCard {...commonProps} />;
      default:
        return null;
    }
  };

  return (
    <div className="w-full max-w-md mx-auto relative" {...handlers}>
      {/* 卡片容器 */}
      <div className="relative overflow-hidden">
        {availableCards.map((card, index) => (
          <div
            key={card.id}
            className={`w-full transition-all duration-300 ease-out ${
              index === currentIndex 
                ? 'opacity-100 block transform scale-100' 
                : 'opacity-0 hidden transform scale-95'
            }`}
          >
            {renderCard(card)}
          </div>
        ))}
      </div>

      {/* PC端切换按钮 */}
      {!isMobile && availableCards.length > 1 && (
        <>
          <button
            onClick={handlePrev}
            disabled={currentIndex === 0}
            className="absolute left-0 top-1/2 -translate-y-1/2 -translate-x-12 bg-white rounded-full p-2 shadow-lg hover:bg-gray-50 disabled:opacity-30 disabled:cursor-not-allowed transition-all"
            aria-label="上一张"
          >
            <span className="material-symbols-outlined text-2xl text-gray-600">
              chevron_left
            </span>
          </button>
          <button
            onClick={handleNext}
            disabled={currentIndex === availableCards.length - 1}
            className="absolute right-0 top-1/2 -translate-y-1/2 translate-x-12 bg-white rounded-full p-2 shadow-lg hover:bg-gray-50 disabled:opacity-30 disabled:cursor-not-allowed transition-all"
            aria-label="下一张"
          >
            <span className="material-symbols-outlined text-2xl text-gray-600">
              chevron_right
            </span>
          </button>
        </>
      )}

      {/* 导出按钮（PC和移动端都显示） */}
      {onExport && (
        <button
          onClick={handleExport}
          className={`absolute top-4 right-4 bg-white rounded-full shadow-lg hover:bg-gray-50 active:scale-95 transition-all z-10 ${
            isMobile ? 'p-2.5' : 'p-3'
          }`}
          aria-label="导出卡片"
          title="导出卡片"
        >
          <span className={`material-symbols-outlined text-gray-600 ${isMobile ? 'text-lg' : 'text-xl'}`}>
            download
          </span>
        </button>
      )}

      {/* 指示器 */}
      {availableCards.length > 1 && (
        <div className="flex justify-center gap-2 mt-4">
          {availableCards.map((_, index) => (
            <button
              key={index}
              onClick={() => {
                if (!isTransitioning) {
                  setIsTransitioning(true);
                  setCurrentIndex(index);
                  setTimeout(() => setIsTransitioning(false), 300);
                }
              }}
              className={`h-2 rounded-full transition-all duration-300 ${
                index === currentIndex
                  ? 'bg-science-green w-6'
                  : 'bg-gray-300 hover:bg-gray-400 w-2'
              }`}
              aria-label={`切换到第${index + 1}张卡片`}
            />
          ))}
        </div>
      )}
    </div>
  );
};
