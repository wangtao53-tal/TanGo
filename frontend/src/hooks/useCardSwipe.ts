/**
 * 卡片滑动Hook
 * 处理触摸事件和滑动逻辑
 */

import { useState, useCallback } from 'react';

export interface UseCardSwipeOptions {
  onSwipeLeft?: () => void;  // 左滑（切换到下一张）
  onSwipeRight?: () => void; // 右滑（切换到上一张）
  threshold?: number;         // 滑动阈值（默认30%）
}

export interface UseCardSwipeReturn {
  handleSwipe: (deltaX: number, deltaY: number) => void;
  handleSwipeEnd: () => void;
  isSwiping: boolean;
}

export function useCardSwipe({
  onSwipeLeft,
  onSwipeRight,
  threshold = 0.3, // 30%屏幕宽度
}: UseCardSwipeOptions = {}): UseCardSwipeReturn {
  const [isSwiping, setIsSwiping] = useState(false);
  const [startX, setStartX] = useState(0);
  const [currentX, setCurrentX] = useState(0);

  const handleSwipe = useCallback((deltaX: number, deltaY: number) => {
    // 只处理水平滑动（水平滑动距离大于垂直滑动距离）
    if (Math.abs(deltaX) > Math.abs(deltaY)) {
      setIsSwiping(true);
      setCurrentX(deltaX);
    }
  }, []);

  const handleSwipeEnd = useCallback(() => {
    if (!isSwiping) return;

    // 获取屏幕宽度
    const screenWidth = window.innerWidth;
    const swipeDistance = Math.abs(currentX);
    const swipeThreshold = screenWidth * threshold;

    if (swipeDistance > swipeThreshold) {
      if (currentX > 0 && onSwipeRight) {
        // 向右滑动（切换到上一张）
        onSwipeRight();
      } else if (currentX < 0 && onSwipeLeft) {
        // 向左滑动（切换到下一张）
        onSwipeLeft();
      }
    }

    setIsSwiping(false);
    setCurrentX(0);
    setStartX(0);
  }, [isSwiping, currentX, threshold, onSwipeLeft, onSwipeRight]);

  return {
    handleSwipe,
    handleSwipeEnd,
    isSwiping,
  };
}

