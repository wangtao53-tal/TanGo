/**
 * 文本转语音Hook
 * 使用Web Speech API实现文本转语音功能
 */

import { useState, useRef, useCallback } from 'react';

export interface UseTextToSpeechOptions {
  rate?: number; // 语速 (0.1-10, 默认0.9)
  pitch?: number; // 音调 (0-2, 默认1.0)
  volume?: number; // 音量 (0-1, 默认1.0)
  onStart?: () => void;
  onEnd?: () => void;
  onError?: (error: Error) => void;
}

export interface UseTextToSpeechReturn {
  isPlaying: boolean;
  isPaused: boolean;
  play: (text: string, language?: string) => void;
  pause: () => void;
  resume: () => void;
  stop: () => void;
  isSupported: boolean;
}

export function useTextToSpeech(
  options: UseTextToSpeechOptions = {}
): UseTextToSpeechReturn {
  const {
    rate = 0.9,
    pitch = 1.0,
    volume = 1.0,
    onStart,
    onEnd,
    onError,
  } = options;

  const [isPlaying, setIsPlaying] = useState(false);
  const [isPaused, setIsPaused] = useState(false);
  const utteranceRef = useRef<SpeechSynthesisUtterance | null>(null);
  const synthRef = useRef<SpeechSynthesis | null>(null);

  // 检查浏览器支持
  const isSupported =
    typeof window !== 'undefined' && 'speechSynthesis' in window;

  // 初始化
  if (isSupported && !synthRef.current) {
    synthRef.current = window.speechSynthesis;
  }

  const play = useCallback(
    (text: string, language: string = 'zh-CN') => {
      if (!isSupported || !synthRef.current) {
        const error = new Error('浏览器不支持Web Speech API');
        onError?.(error);
        return;
      }

      // 停止当前播放
      stop();

      try {
        const utterance = new SpeechSynthesisUtterance(text);
        utterance.lang = language;
        utterance.rate = rate;
        utterance.pitch = pitch;
        utterance.volume = volume;

        utterance.onstart = () => {
          setIsPlaying(true);
          setIsPaused(false);
          onStart?.();
        };

        utterance.onend = () => {
          setIsPlaying(false);
          setIsPaused(false);
          utteranceRef.current = null;
          onEnd?.();
        };

        utterance.onerror = (event) => {
          setIsPlaying(false);
          setIsPaused(false);
          utteranceRef.current = null;
          const error = new Error(
            `语音播放错误: ${event.error || '未知错误'}`
          );
          onError?.(error);
        };

        utteranceRef.current = utterance;
        synthRef.current.speak(utterance);
      } catch (error) {
        setIsPlaying(false);
        setIsPaused(false);
        onError?.(error as Error);
      }
    },
    [isSupported, rate, pitch, volume, onStart, onEnd, onError]
  );

  const pause = useCallback(() => {
    if (synthRef.current && isPlaying && !isPaused) {
      synthRef.current.pause();
      setIsPaused(true);
    }
  }, [isPlaying, isPaused]);

  const resume = useCallback(() => {
    if (synthRef.current && isPaused) {
      synthRef.current.resume();
      setIsPaused(false);
    }
  }, [isPaused]);

  const stop = useCallback(() => {
    if (synthRef.current) {
      synthRef.current.cancel();
      setIsPlaying(false);
      setIsPaused(false);
      utteranceRef.current = null;
    }
  }, []);

  return {
    isPlaying,
    isPaused,
    play,
    pause,
    resume,
    stop,
    isSupported,
  };
}
