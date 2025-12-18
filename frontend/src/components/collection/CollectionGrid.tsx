/**
 * 收藏网格布局组件
 * 基于 stitch_ui/favorites_page/ 设计稿
 */

import React from 'react';
import type { ExplorationRecord } from '../../types/exploration';
import type { KnowledgeCard } from '../../types/exploration';
import { CollectionCard } from './CollectionCard';
import { ScienceCard } from '../cards/ScienceCard';
import { PoetryCard } from '../cards/PoetryCard';
import { EnglishCard } from '../cards/EnglishCard';
import { exportCardAsImage } from '../../utils/export';

export interface CollectionGridProps {
  records: ExplorationRecord[];
  cards?: KnowledgeCard[];
  category?: 'all' | '自然类' | '生活类' | '人文类';
  onReExplore?: (recordId: string) => void;
  onExport?: (cardId: string) => void;
}

export const CollectionGrid: React.FC<CollectionGridProps> = ({
  records,
  cards = [],
  category = 'all',
  onReExplore,
  onExport,
}) => {
  const filteredRecords =
    category === 'all'
      ? records
      : records.filter((r) => r.objectCategory === category);

  // 从records中提取所有卡片
  const allCardsFromRecords: KnowledgeCard[] = [];
  filteredRecords.forEach((record) => {
    allCardsFromRecords.push(...record.cards);
  });

  // 合并所有卡片
  const allCards = [...allCardsFromRecords, ...cards];
  const filteredCards = category === 'all' 
    ? allCards 
    : allCards.filter((card) => {
        const record = records.find((r) => r.cards.some((c) => c.id === card.id));
        return record?.objectCategory === category;
      });

  const handleExport = async (cardId: string) => {
    if (onExport) {
      onExport(cardId);
    } else {
      try {
        await exportCardAsImage(`card-${cardId}`, `card-${cardId}`);
      } catch (error) {
        console.error('导出失败:', error);
      }
    }
  };

  if (filteredRecords.length === 0 && filteredCards.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-20">
        <div className="text-6xl mb-4">📚</div>
        <p className="text-text-sub text-lg font-display">
          还没有收藏任何卡片，快去探索吧！
        </p>
      </div>
    );
  }

  return (
    <div className="space-y-8">
      {/* 显示探索记录 */}
      {filteredRecords.length > 0 && (
        <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-8">
          {filteredRecords.map((record) => (
            <CollectionCard
              key={record.id}
              record={record}
              onReExplore={onReExplore}
            />
          ))}
        </div>
      )}

      {/* 显示卡片 */}
      {filteredCards.length > 0 && (
        <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-8">
          {filteredCards.map((card) => {
            const cardElement = card.type === 'science' ? (
              <ScienceCard 
                key={card.id} 
                card={card} 
                id={`card-${card.id}`}
                onCollect={handleExport}
              />
            ) : card.type === 'poetry' ? (
              <PoetryCard 
                key={card.id} 
                card={card} 
                id={`card-${card.id}`}
                onCollect={handleExport}
              />
            ) : (
              <EnglishCard 
                key={card.id} 
                card={card} 
                id={`card-${card.id}`}
                onCollect={handleExport}
              />
            );

            return (
              <div key={card.id} className="relative group">
                {cardElement}
                {/* 导出按钮 */}
                <button
                  onClick={() => handleExport(card.id)}
                  className="absolute top-4 right-4 z-10 size-10 rounded-full bg-white/90 hover:bg-white text-gray-600 shadow-lg flex items-center justify-center transition-all opacity-0 group-hover:opacity-100"
                  title="导出卡片"
                >
                  <span className="material-symbols-outlined">download</span>
                </button>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
};

