/**
 * 英语表达卡组件
 * 蓝色主题，基于设计稿
 */

import React, { useState, useEffect, useRef, useCallback } from 'react';
import { cardStorage } from '../../services/storage';
import type { KnowledgeCard } from '../../types/exploration';
import type { EnglishCardContent } from '../../types/exploration';
import { useTextToSpeech } from '../../hooks/useTextToSpeech';
import { detectCardLanguage } from '../../utils/cardTextExtractor';

/**
 * 移除文本中的所有emoji表情符号
 */
function removeEmojis(text: string): string {
  if (!text) return text;
  
  // 更全面的emoji Unicode范围正则表达式
  const emojiRegex = /[\u{1F600}-\u{1F64F}]|[\u{1F300}-\u{1F5FF}]|[\u{1F680}-\u{1F6FF}]|[\u{1F900}-\u{1F9FF}]|[\u{2600}-\u{26FF}]|[\u{2700}-\u{27BF}]|[\u{1F1E0}-\u{1F1FF}]|[\u{1F191}-\u{1F251}]|[\u{2934}\u{2935}]|[\u{2190}-\u{21FF}]|[\u{2B00}-\u{2BFF}]|[\u{FE00}-\u{FE0F}]|[\u{200D}]/gu;
  
  // 移除emoji并清理多余的空格
  return text.replace(emojiRegex, '').replace(/\s+/g, ' ').trim();
}
import { usePlayingCard } from './ScienceCard';
import { cardStyles } from '../../styles/cardStyles';

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
  // 处理字段映射：后端可能返回 keywords，前端使用 words
  const rawContent = card.content as any;
  const content: EnglishCardContent = {
    words: rawContent.words || rawContent.keywords || [],
    expressions: rawContent.expressions || [],
    pronunciation: rawContent.pronunciation,
  };
  const { playingCardId, setPlayingCardId } = usePlayingCard();
  const isCurrentlyPlaying = playingCardId === card.id;
  // 用于跟踪正在高亮的单词
  const [highlightedWords, setHighlightedWords] = useState<Set<string>>(new Set());
  // 用于跟踪是否正在播放自定义序列
  const [isPlayingSequence, setIsPlayingSequence] = useState(false);

  // 文本转语音Hook - 英语卡片使用更慢的语速（0.7）以便学习
  const { isPlaying, stop, isSupported } = useTextToSpeech({
    rate: 0.7,
    pitch: 1.0,
    onStart: () => {
      setPlayingCardId(card.id);
    },
    onEnd: () => {
      setPlayingCardId(null);
      setHighlightedWords(new Set());
    },
    onError: (error) => {
      console.error('文本转语音错误:', error);
      setPlayingCardId(null);
      setHighlightedWords(new Set());
    },
  });

  // 如果其他卡片开始播放，停止当前播放
  useEffect(() => {
    if (playingCardId !== null && playingCardId !== card.id) {
      // 停止自定义播放序列
      if (synthRef.current && isPlayingSequence) {
        synthRef.current.cancel();
        setIsPlayingSequence(false);
        setHighlightedWords(new Set());
        currentIndexRef.current = 0;
      }
      // 停止 useTextToSpeech hook 的播放
      if (isPlaying) {
        stop();
      }
    }
  }, [playingCardId, card.id, isPlaying, isPlayingSequence, stop]);

  // 用于存储播放序列的引用
  const playSequenceRef = useRef<Array<{ text: string; highlightWords?: boolean }>>([]);
  const currentIndexRef = useRef(0);
  const synthRef = useRef<SpeechSynthesis | null>(null);

  useEffect(() => {
    if (typeof window !== 'undefined' && 'speechSynthesis' in window) {
      synthRef.current = window.speechSynthesis;
    }
  }, []);

  // 顺序播放函数
  const playSequence = useCallback((sequence: Array<{ text: string; highlightWords?: boolean }>, language: string) => {
    if (!synthRef.current || sequence.length === 0) return;

    // 停止当前播放
    synthRef.current.cancel();
    
    playSequenceRef.current = sequence;
    currentIndexRef.current = 0;
    setIsPlayingSequence(true);
    setPlayingCardId(card.id);

    const playNext = () => {
      if (currentIndexRef.current >= playSequenceRef.current.length) {
        // 播放完成
        setIsPlayingSequence(false);
        setPlayingCardId(null);
        setHighlightedWords(new Set());
        return;
      }

      const item = playSequenceRef.current[currentIndexRef.current];
      
      // 如果需要高亮单词，设置高亮状态
      if (item.highlightWords && content.words) {
        setHighlightedWords(new Set(content.words));
      } else {
        // 其他部分播放时，清除单词高亮
        setHighlightedWords(new Set());
      }

      // 播放当前文本
      const utterance = new SpeechSynthesisUtterance(item.text);
      utterance.lang = language;
      utterance.rate = 0.7;
      utterance.pitch = 1.0;
      utterance.volume = 1.0;

      utterance.onend = () => {
        currentIndexRef.current++;
        // 短暂延迟后播放下一段
        setTimeout(playNext, 200);
      };

      utterance.onerror = () => {
        currentIndexRef.current++;
        setTimeout(playNext, 200);
      };

      if (synthRef.current) {
        synthRef.current.speak(utterance);
      }
    };

    playNext();
  }, [card.id, content.words, setPlayingCardId]);

  const handleListen = () => {
    if (!isSupported) {
      alert('您的浏览器不支持文本转语音功能');
      return;
    }

    // 如果正在播放，停止播放
    if (synthRef.current && (isPlayingSequence || currentIndexRef.current > 0)) {
      synthRef.current.cancel();
      setIsPlayingSequence(false);
      setPlayingCardId(null);
      setHighlightedWords(new Set());
      currentIndexRef.current = 0;
      return;
    }

    // 构建播放文本序列
    const playSequenceItems: Array<{ text: string; highlightWords?: boolean }> = [];
    
    // 1. 播放标题（移除emoji）
    if (card.title) {
      playSequenceItems.push({ text: removeEmojis(card.title) });
    }
    
    // 2. 播放单词部分（需要高亮，移除emoji）
    if (content.words && content.words.length > 0) {
      const wordsText = '核心单词：' + content.words.map(w => removeEmojis(w)).join(', ');
      playSequenceItems.push({ text: wordsText, highlightWords: true });
    }
    
    // 3. 播放表达式部分（移除emoji）
    if (content.expressions && content.expressions.length > 0) {
      const expressionsText = '口语表达：' + content.expressions.map(e => removeEmojis(e)).join('. ');
      playSequenceItems.push({ text: expressionsText });
    }
    
    // 4. 播放发音部分（移除emoji）
    if (content.pronunciation) {
      playSequenceItems.push({ text: '发音：' + removeEmojis(content.pronunciation) });
    }
    
    // 开始顺序播放
    const language = detectCardLanguage(card);
    playSequence(playSequenceItems, language);
  };

  const handleStop = () => {
    // 停止自定义播放序列
    if (synthRef.current) {
      synthRef.current.cancel();
    }
    setIsPlayingSequence(false);
    currentIndexRef.current = 0;
    setPlayingCardId(null);
    setHighlightedWords(new Set());
    // 也停止 useTextToSpeech hook 的播放（如果正在使用）
    stop();
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
    
    // 乐观更新UI
    setIsCollected(newCollectedState);
    
    if (onCollect) {
      onCollect(card.id);
    }

    try {
      // 保存到本地存储
      if (newCollectedState) {
        const cardToSave = {
          ...card,
          collectedAt: new Date().toISOString(),
        };
        await cardStorage.save(cardToSave);
      } else {
        await cardStorage.delete(card.id);
      }
    } catch (error) {
      // 如果保存失败，回滚状态
      console.error('收藏操作失败:', error);
      setIsCollected(!newCollectedState);
      // 可选：显示错误提示
      // alert('收藏操作失败，请重试');
    }
  };


  return (
    <article
      id={id || `card-${card.id}`}
      className={`flex flex-col rounded-[2.5rem] bg-white border-4 border-sky-blue shadow-card relative transition-all hover:-translate-y-2 hover:shadow-xl hover:z-10 duration-300 group overflow-hidden w-full max-w-md mx-auto my-2 ${className}`}
    >
      {/* 内容区域（纯文本显示，不显示图片） */}
      <div className="bg-white p-6 flex flex-col justify-between relative rounded-[2.2rem] min-h-[400px]">
        <div className="flex-1 flex flex-col">
          <div className="flex items-center justify-between mb-4">
            <h3 
              className="font-bold"
              style={{ 
                fontSize: cardStyles.fonts.sizes.title,
                lineHeight: cardStyles.fonts.lineHeight.title,
                fontFamily: cardStyles.fonts.childFriendly.english,
                color: '#0284C7' // 使用更深的蓝色，确保对比度≥4.5:1
              }}
            >
              {card.title}
            </h3>
            <button
              onClick={handleListen}
              disabled={!isSupported}
              className={`size-10 rounded-full flex items-center justify-center transition-all duration-200 shadow-sm disabled:opacity-50 disabled:cursor-not-allowed active:scale-95 ${
                isCurrentlyPlaying && isPlayingSequence
                  ? 'bg-sky-blue text-white'
                  : 'bg-sky-blue/20 hover:bg-sky-blue hover:text-white text-sky-blue'
              }`}
            >
              <span className="material-symbols-outlined text-xl">
                {isCurrentlyPlaying && isPlayingSequence ? 'pause' : 'volume_up'}
              </span>
            </button>
          </div>

          <div className="space-y-4 flex-1">
            {/* 核心单词 */}
            {content.words && content.words.length > 0 && (
              <div>
                <h4 
                  className="font-black uppercase mb-2 tracking-wider"
                  style={{ 
                    fontSize: cardStyles.fonts.sizes.small,
                    lineHeight: cardStyles.fonts.lineHeight.small,
                    fontFamily: cardStyles.fonts.childFriendly.english,
                    color: '#0284C7' // 使用更深的蓝色，确保对比度≥4.5:1
                  }}
                >
                  Magic Words
                </h4>
                <div className="flex gap-2 flex-wrap">
                  {content.words.map((word, index) => {
                    const isHighlighted = highlightedWords.has(word);
                    return (
                      <span
                        key={index}
                        className={`px-3 py-1.5 rounded-xl font-bold border transition-all duration-300 ${
                          isHighlighted
                            ? 'bg-sky-blue text-white border-sky-blue scale-110 shadow-lg'
                            : 'bg-sky-blue/10 border-sky-blue/20 hover:bg-sky-blue hover:text-white cursor-pointer'
                        }`}
                        style={{ 
                          fontSize: cardStyles.fonts.sizes.small,
                          lineHeight: cardStyles.fonts.lineHeight.small,
                          fontFamily: cardStyles.fonts.childFriendly.english,
                          color: isHighlighted ? '#FFFFFF' : '#0284C7' // 高亮时使用白色，否则使用蓝色
                        }}
                      >
                        {word}
                      </span>
                    );
                  })}
                </div>
              </div>
            )}

            {/* 口语表达 */}
            {content.expressions && content.expressions.length > 0 && (
              <div>
                <h4 
                  className="font-black uppercase mb-2 tracking-wider"
                  style={{ 
                    fontSize: cardStyles.fonts.sizes.small,
                    lineHeight: cardStyles.fonts.lineHeight.small,
                    fontFamily: cardStyles.fonts.childFriendly.english,
                    color: '#0284C7' // 使用更深的蓝色，确保对比度≥4.5:1
                  }}
                >
                  Let's Talk!
                </h4>
                <div className="bg-sky-50 p-4 rounded-2xl border-2 border-sky-blue/20 relative">
                  <div 
                    className="absolute -top-3 right-4 bg-sky-200 text-sky-800 font-bold px-2 py-0.5 rounded-full"
                    style={{ 
                      fontSize: cardStyles.fonts.sizes.small,
                      fontFamily: cardStyles.fonts.childFriendly.english
                    }}
                  >
                    Try saying this
                  </div>
                  {content.expressions.map((expr, index) => (
                    <p 
                      key={index} 
                      className="mb-2 font-bold"
                      style={{ 
                        fontSize: cardStyles.fonts.sizes.body,
                        lineHeight: cardStyles.fonts.lineHeight.body,
                        fontFamily: cardStyles.fonts.childFriendly.english,
                        color: '#1F2937' // 使用更深的灰色，确保对比度≥4.5:1
                      }}
                    >
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
          <div className="flex items-center gap-2">
            <button
              onClick={handleListen}
              disabled={!isSupported}
              className="flex items-center gap-2 text-sm font-bold text-sky-blue bg-sky-blue/20 hover:bg-sky-blue hover:text-white px-5 py-3 rounded-full transition-all duration-200 shadow-sm hover:shadow-md disabled:opacity-50 disabled:cursor-not-allowed active:scale-95"
            >
              <span className="material-symbols-outlined text-xl">
                {isCurrentlyPlaying && isPlayingSequence ? 'pause' : 'play_circle'}
              </span>
              {isCurrentlyPlaying && isPlayingSequence ? '暂停' : '听'}
            </button>
            {isCurrentlyPlaying && isPlayingSequence && (
              <button
                onClick={handleStop}
                className="flex items-center gap-1 text-xs font-bold text-slate-500 hover:text-slate-700 px-3 py-2 rounded-full transition-all duration-200 active:scale-95"
              >
                <span className="material-symbols-outlined text-lg">stop</span>
                停止
              </button>
            )}
          </div>
          <button
            onClick={handleCollect}
            className={`size-12 rounded-full flex items-center justify-center transition-all duration-300 shadow-sm border-2 active:scale-90 ${
              isCollected
                ? 'bg-yellow-300 text-yellow-800 border-yellow-400 hover:bg-yellow-400 hover:scale-105'
                : 'bg-slate-100 text-slate-300 border-slate-200 hover:border-yellow-400 hover:bg-yellow-50'
            }`}
          >
            <span 
              className={`material-symbols-outlined text-2xl transition-all duration-300 ${
                isCollected ? 'fill-current scale-110' : 'fill-current'
              }`}
            >
              star
            </span>
          </button>
        </div>
      </div>
    </article>
  );
};

