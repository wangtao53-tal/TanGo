/**
 * 收藏网格布局组件
 * 基于 stitch_ui/favorites_page/ 设计稿
 */

import React from 'react';
import type { ExplorationRecord } from '@/types/exploration';
import { CollectionCard } from './CollectionCard';

export interface CollectionGridProps {
  records: ExplorationRecord[];
  category?: 'all' | '自然类' | '生活类' | '人文类';
  onReExplore?: (recordId: string) => void;
}

export const CollectionGrid: React.FC<CollectionGridProps> = ({
  records,
  category = 'all',
  onReExplore,
}) => {
  const filteredRecords =
    category === 'all'
      ? records
      : records.filter((r) => r.objectCategory === category);

  if (filteredRecords.length === 0) {
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
    <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-8">
      {filteredRecords.map((record) => (
        <CollectionCard
          key={record.id}
          record={record}
          onReExplore={onReExplore}
        />
      ))}
    </div>
  );
};

