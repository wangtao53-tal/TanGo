/**
 * 卡片轮播容器组件
 * 实现渐进式展示和滑动切换
 */

import React, { useState, useEffect, useMemo } from 'react';
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
  const [currentIndex, setCurrentIndex] = useState(1); // 默认显示中间一张（主体介绍）
  const [isMobile, setIsMobile] = useState(false);
  const [showSwipeHint, setShowSwipeHint] = useState(true); // 显示滑动提示
  const [isInitialized, setIsInitialized] = useState(false); // 是否已初始化

  // 检测是否为移动端
  useEffect(() => {
    const checkMobile = () => {
      setIsMobile(window.innerWidth < 1024);
    };
    checkMobile();
    window.addEventListener('resize', checkMobile);
    return () => window.removeEventListener('resize', checkMobile);
  }, []);

  // 按类型排序卡片：poetry（左） -> science（中，默认显示） -> english（右）
  const sortedCards = useMemo(() => {
    const cardMap = new Map<string, KnowledgeCard>();
    cards.forEach(card => {
      // 如果已存在同类型卡片，保留第一个（通常是先生成的）
      if (!cardMap.has(card.type)) {
        cardMap.set(card.type, card);
      }
    });
    
    const result: (KnowledgeCard | null)[] = [
      cardMap.get('poetry') || null,   // 索引0：古诗文（左）
      cardMap.get('science') || null,   // 索引1：主体介绍（中，默认显示）
      cardMap.get('english') || null,   // 索引2：英语学习（右）
    ];
    
    return result;
  }, [cards]);

  // 只显示已生成的卡片，保持排序顺序：poetry（左） -> science（中） -> english（右）
  const availableCards = useMemo(() => {
    // 按照固定顺序构建数组，确保顺序正确
    const result: KnowledgeCard[] = [];
    
    // 按顺序添加：poetry -> science -> english
    sortedCards.forEach(card => {
      if (card !== null) {
        result.push(card);
      }
    });
    
    return result;
  }, [sortedCards]);

  // 确保currentIndex在有效范围内，默认显示第二张卡片（主体介绍，索引1）
  useEffect(() => {
    if (availableCards.length > 0) {
      // 找到 science 卡片在 availableCards 中的索引
      const scienceIndex = availableCards.findIndex(card => card.type === 'science');
      
      // 如果还没有初始化
      if (!isInitialized) {
        // 如果三张卡片都存在，默认显示索引1（第二张，主体介绍）
        if (availableCards.length === 3 && scienceIndex === 1) {
          setCurrentIndex(1);
          setIsInitialized(true);
          return;
        }
        
        // 如果只有1-2张卡片，优先显示 science 卡片，否则显示中间位置
        if (availableCards.length <= 2) {
          if (scienceIndex >= 0) {
            setCurrentIndex(scienceIndex);
          } else {
            const middleIndex = Math.floor(availableCards.length / 2);
            setCurrentIndex(middleIndex);
          }
          setIsInitialized(true);
          return;
        }
      }
      
      // 如果已经初始化，但 currentIndex 超出范围，重置为有效索引
      if (isInitialized && currentIndex >= availableCards.length) {
        // 优先显示 science 卡片，否则显示中间位置
        if (scienceIndex >= 0 && scienceIndex < availableCards.length) {
          setCurrentIndex(scienceIndex);
        } else {
          const middleIndex = Math.floor(availableCards.length / 2);
          setCurrentIndex(middleIndex);
        }
      }
      
      // 如果索引无效，重置为0
      if (currentIndex < 0) {
        setCurrentIndex(0);
      }
    }
  }, [availableCards.length, isInitialized]);

  // 自动隐藏提示（5秒后）
  useEffect(() => {
    if (showSwipeHint && isMobile && availableCards.length > 1 && currentIndex === 1) {
      const timer = setTimeout(() => {
        setShowSwipeHint(false);
      }, 5000);
      return () => clearTimeout(timer);
    }
  }, [showSwipeHint, isMobile, availableCards.length, currentIndex]);

  const [isTransitioning, setIsTransitioning] = useState(false);

  // 切换函数（PC和移动端都使用）
  const handlePrev = () => {
    // 左按钮：向左切换（显示索引更小的卡片，即古诗文）
    if (!isTransitioning && availableCards.length > 1) {
      setIsTransitioning(true);
      setShowSwipeHint(false); // 隐藏提示
      
      const prevIndex = currentIndex > 0 
        ? currentIndex - 1 
        : availableCards.length - 1; // 循环：第一张时，点击上一张跳到最后一张
      
      setCurrentIndex(prevIndex);
      setTimeout(() => setIsTransitioning(false), 300);
    }
  };

  const handleNext = () => {
    // 右按钮：向右切换（显示索引更大的卡片，即英语学习）
    if (!isTransitioning && availableCards.length > 1) {
      setIsTransitioning(true);
      setShowSwipeHint(false); // 隐藏提示
      
      const nextIndex = currentIndex < availableCards.length - 1 
        ? currentIndex + 1 
        : 0; // 循环：最后一张时，点击下一张回到第一张
      
      setCurrentIndex(nextIndex);
      setTimeout(() => setIsTransitioning(false), 300);
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

  // 确保 currentIndex 在有效范围内
  const safeCurrentIndex = Math.max(0, Math.min(currentIndex, availableCards.length - 1));
  const currentCard = availableCards[safeCurrentIndex];

  // 导出当前卡片
  const handleExport = () => {
    if (currentCard && onExport) {
      onExport(currentCard.id);
    }
  };

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
    <div className="w-full max-w-md mx-auto relative overflow-x-hidden px-2 md:px-0">
      {/* 卡片容器 */}
      <div className="relative overflow-hidden w-full">
        {availableCards.map((card, index) => (
          <div
            key={card.id}
            className={`w-full transition-all duration-300 ease-out ${
              index === safeCurrentIndex 
                ? 'opacity-100 block transform scale-100' 
                : 'opacity-0 hidden transform scale-95'
            }`}
          >
            {renderCard(card)}
          </div>
        ))}
      </div>

      {/* 切换按钮（PC和移动端都显示） */}
      {availableCards.length > 1 && (
        <>
          <button
            onClick={handlePrev}
            className={`absolute top-1/2 -translate-y-1/2 bg-white rounded-full p-2 shadow-lg hover:bg-gray-50 active:scale-95 transition-all z-10 ${
              isMobile 
                ? 'left-2' 
                : 'left-0 -translate-x-12'
            }`}
            aria-label="上一张（古诗文）"
            title="古诗文"
          >
            <span className={`material-symbols-outlined text-gray-600 ${isMobile ? 'text-xl' : 'text-2xl'}`}>
              chevron_left
            </span>
          </button>
          <button
            onClick={handleNext}
            className={`absolute top-1/2 -translate-y-1/2 bg-white rounded-full p-2 shadow-lg hover:bg-gray-50 active:scale-95 transition-all z-10 ${
              isMobile 
                ? 'right-2' 
                : 'right-0 translate-x-12'
            }`}
            aria-label="下一张（英语学习）"
            title="英语学习"
          >
            <span className={`material-symbols-outlined text-gray-600 ${isMobile ? 'text-xl' : 'text-2xl'}`}>
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
                index === safeCurrentIndex
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
