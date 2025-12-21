/**
 * 收藏网格布局组件
 * 基于 stitch_ui/favorites_page/ 设计稿
 */

import React from 'react';
import { useTranslation } from 'react-i18next';
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
  onToggleCollect?: (recordId: string, collected: boolean) => void; // 切换收藏回调
}

export const CollectionGrid: React.FC<CollectionGridProps> = ({
  records,
  cards = [],
  category = 'all',
  onReExplore,
  onExport,
  onToggleCollect,
}) => {
  const { t } = useTranslation();
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
        console.error(t('collection.exportError'));
      }
    }
  };

  if (filteredRecords.length === 0 && filteredCards.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-20">
        <div className="text-6xl mb-4">📚</div>
        <p className="text-text-sub text-lg font-display">
          {t('collection.emptyMessage')}
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
              isCollected={record.collected}
              onReExplore={onReExplore}
              onToggleCollect={onToggleCollect}
            />
          ))}
        </div>
      )}

      {/* 显示卡片 */}
      {filteredCards.length > 0 && (
        <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-8">
          {filteredCards.map((card) => {
            // 知识卡片组件有自己的收藏逻辑，不需要传入onCollect
            // onCollect prop原本是用于导出功能的，但导出功能应该通过单独的导出按钮实现
            const cardElement = card.type === 'science' ? (
              <ScienceCard 
                key={card.id} 
                card={card} 
                id={`card-${card.id}`}
              />
            ) : card.type === 'poetry' ? (
              <PoetryCard 
                key={card.id} 
                card={card} 
                id={`card-${card.id}`}
              />
            ) : (
              <EnglishCard 
                key={card.id} 
                card={card} 
                id={`card-${card.id}`}
              />
            );

            return (
              <div key={card.id} className="relative group">
                {cardElement}
                {/* 导出按钮（PC和移动端都显示） */}
                <button
                  onClick={() => handleExport(card.id)}
                  className="absolute top-4 right-4 z-10 size-10 rounded-full bg-white/90 hover:bg-white active:scale-95 text-gray-600 shadow-lg flex items-center justify-center transition-all"
                  title={t('collection.exportCardTitle')}
                  aria-label={t('collection.exportCardTitle')}
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

