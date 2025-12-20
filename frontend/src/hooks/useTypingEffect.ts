/**
 * 打字机效果Hook
 * 实现逐字显示文本的效果
 */

import { useState, useEffect, useRef } from 'react';

export interface UseTypingEffectOptions {
  text: string;
  speed?: number; // 每字符显示间隔（毫秒），默认30ms
  enabled?: boolean; // 是否启用打字机效果，默认true
}

export function useTypingEffect(options: UseTypingEffectOptions): string {
  const { text, speed = 30, enabled = true } = options;
  const [displayedText, setDisplayedText] = useState('');
  const [currentIndex, setCurrentIndex] = useState(0);
  const timerRef = useRef<NodeJS.Timeout | null>(null);

  // 当text更新时，重置并重新开始
  useEffect(() => {
    setCurrentIndex(0);
    setDisplayedText('');
    
    // 清除之前的定时器
    if (timerRef.current) {
      clearTimeout(timerRef.current);
    }
  }, [text]);

  // 打字机效果
  useEffect(() => {
    if (!enabled || currentIndex >= text.length) {
      return;
    }

    timerRef.current = setTimeout(() => {
      setDisplayedText(text.slice(0, currentIndex + 1));
      setCurrentIndex(currentIndex + 1);
    }, speed);

    return () => {
      if (timerRef.current) {
        clearTimeout(timerRef.current);
      }
    };
  }, [currentIndex, text, speed, enabled]);

  // 如果禁用打字机效果，直接显示完整文本
  if (!enabled) {
    return text;
  }

  return displayedText;
}

