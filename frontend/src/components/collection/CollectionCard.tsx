/**
 * 收藏卡片组件
 * 基于设计稿中的卡片样式
 */

import React from 'react';
import { useNavigate } from 'react-router-dom';
import type { ExplorationRecord } from '@/types/exploration';

export interface CollectionCardProps {
  record: ExplorationRecord;
  onReExplore?: (recordId: string) => void;
}

const categoryConfig: Record<string, { icon: string; colorClasses: { bg: string; text: string; border: string }; label: string }> = {
  自然类: {
    icon: 'forest',
    colorClasses: {
      bg: 'bg-green-100',
      text: 'text-green-600',
      border: 'border-green-100',
    },
    label: 'Natural',
  },
  生活类: {
    icon: 'house',
    colorClasses: {
      bg: 'bg-orange-100',
      text: 'text-orange-600',
      border: 'border-orange-100',
    },
    label: 'Life',
  },
  人文类: {
    icon: 'palette',
    colorClasses: {
      bg: 'bg-purple-100',
      text: 'text-purple-600',
      border: 'border-purple-100',
    },
    label: 'Humanities',
  },
};

export const CollectionCard: React.FC<CollectionCardProps> = ({
  record,
  onReExplore,
}) => {
  const navigate = useNavigate();
  const config = categoryConfig[record.objectCategory] || categoryConfig['自然类'];
  const timeAgo = getTimeAgo(record.timestamp);

  const handleReExplore = () => {
    if (onReExplore) {
      onReExplore(record.id);
    } else {
      // 默认行为：导航到结果页面
      navigate(`/result?explorationId=${record.id}`);
    }
  };

  return (
    <div className="group relative flex flex-col bg-white rounded-3xl p-4 border border-white shadow-card hover:shadow-card-hover-green transition-all duration-300 hover:-translate-y-2">
      {/* 类别标签 */}
      <div
        className={`absolute top-6 right-6 z-10 bg-white/90 backdrop-blur-sm px-3 py-1.5 rounded-full shadow-sm flex items-center gap-1.5 border ${config.colorClasses.border}`}
      >
        <span className={`material-symbols-outlined ${config.colorClasses.text} text-sm`}>
          {config.icon}
        </span>
        <span className={`text-xs font-bold ${config.colorClasses.text} font-display`}>
          {config.label}
        </span>
      </div>

      {/* 缩略图 */}
      <div className="w-full aspect-[4/3] rounded-2xl bg-gray-100 mb-4 overflow-hidden relative shadow-inner">
        {record.imageData ? (
          <img
            src={record.imageData}
            alt={record.objectName}
            className="w-full h-full object-cover transition-transform duration-700 group-hover:scale-110"
          />
        ) : (
          <div className="w-full h-full bg-gradient-to-br from-gray-200 to-gray-300 flex items-center justify-center">
            <span className="text-6xl">📷</span>
          </div>
        )}
        <div className="absolute inset-0 bg-gradient-to-t from-black/10 to-transparent opacity-0 group-hover:opacity-100 transition-opacity" />
      </div>

      {/* 卡片信息 */}
      <div className="flex flex-col gap-2 px-2 pb-2">
        <h3 className="text-text-main text-xl font-bold font-display tracking-tight group-hover:text-primary transition-colors">
          {record.objectName}
        </h3>
        <div className="flex items-center justify-between mt-2">
          <span className="text-gray-400 text-xs font-semibold bg-gray-50 px-2 py-1 rounded-lg">
            {timeAgo}
          </span>
          <button
            onClick={handleReExplore}
            className="flex items-center gap-2 bg-primary hover:bg-[#5aff2b] text-white px-5 py-2.5 rounded-2xl font-bold text-sm transition-all shadow-md shadow-primary/20 active:scale-95"
          >
            <span>Re-explore</span>
            <span className="material-symbols-outlined text-lg">rocket_launch</span>
          </button>
        </div>
      </div>
    </div>
  );
};

function getTimeAgo(timestamp: string): string {
  const now = new Date();
  const time = new Date(timestamp);
  const diffMs = now.getTime() - time.getTime();
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));

  if (diffDays === 0) return '今天';
  if (diffDays === 1) return '1天前';
  if (diffDays < 7) return `${diffDays}天前`;
  if (diffDays < 30) return `${Math.floor(diffDays / 7)}周前`;
  if (diffDays < 365) return `${Math.floor(diffDays / 30)}个月前`;
  return `${Math.floor(diffDays / 365)}年前`;
}

