/**
 * 分类筛选组件
 * 基于设计稿中的筛选按钮
 */

import React from 'react';
import { useTranslation } from 'react-i18next';

export type Category = 'all' | '自然类' | '生活类' | '人文类';

export interface CategoryFilterProps {
  selected: Category;
  onSelect: (category: Category) => void;
}

const categories: Array<{ 
  value: Category; 
  icon: string; 
  labelKey: string; 
  colorClasses: { bg: string; text: string; hoverBg: string; hoverBorder: string; icon: string }
}> = [
  { 
    value: 'all', 
    icon: 'grid_view', 
    labelKey: 'collection.category.all', 
    colorClasses: { bg: 'bg-yellow-100', text: 'text-yellow-300', hoverBg: 'bg-yellow-200', hoverBorder: 'border-yellow-200', icon: 'text-yellow-300' }
  },
  { 
    value: '自然类', 
    icon: 'forest', 
    labelKey: 'collection.category.natural', 
    colorClasses: { bg: 'bg-green-100', text: 'text-green-600', hoverBg: 'bg-green-200', hoverBorder: 'border-green-200', icon: 'text-green-600' }
  },
  { 
    value: '生活类', 
    icon: 'house', 
    labelKey: 'collection.category.life', 
    colorClasses: { bg: 'bg-orange-100', text: 'text-orange-500', hoverBg: 'bg-orange-200', hoverBorder: 'border-orange-200', icon: 'text-orange-500' }
  },
  { 
    value: '人文类', 
    icon: 'palette', 
    labelKey: 'collection.category.humanities', 
    colorClasses: { bg: 'bg-purple-100', text: 'text-purple-500', hoverBg: 'bg-purple-200', hoverBorder: 'border-purple-200', icon: 'text-purple-500' }
  },
];

export const CategoryFilter: React.FC<CategoryFilterProps> = ({
  selected,
  onSelect,
}) => {
  const { t } = useTranslation();
  return (
    <div className="flex flex-wrap items-center gap-4 sticky top-0 md:relative z-20 py-2">
      {categories.map((category) => {
        const isSelected = selected === category.value;
        return (
          <button
            key={category.value}
            onClick={() => onSelect(category.value)}
            className={`flex items-center gap-2 px-6 py-3 rounded-full transition-all border-2 ${
              isSelected
                ? 'bg-text-main text-white shadow-lg shadow-text-main/20 hover:scale-105 border-transparent'
                : 'bg-white text-text-main hover:bg-gray-50 border-transparent hover:border-gray-200 shadow-sm hover:shadow-md group'
            }`}
          >
            {isSelected ? (
              <span className={`material-symbols-outlined ${category.colorClasses.icon}`}>
                {category.icon}
              </span>
            ) : (
              <div className={`${category.colorClasses.bg} p-1 rounded-full group-hover:${category.colorClasses.hoverBg} transition-colors`}>
                <span className={`material-symbols-outlined ${category.colorClasses.text} text-sm`}>
                  {category.icon}
                </span>
              </div>
            )}
            <span className="text-sm font-bold font-display">
              {t(category.labelKey as any)}
            </span>
          </button>
        );
      })}
    </div>
  );
};

