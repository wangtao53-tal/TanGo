/**
 * 页面头部组件
 * 参考设计稿中的header样式
 */

import React from 'react';
import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

export interface HeaderProps {
  title?: string;
  showFavorites?: boolean;
  showReport?: boolean;
}

export const Header: React.FC<HeaderProps> = ({
  title,
  showFavorites = true,
  showReport = true,
}) => {
  const { t } = useTranslation();
  
  return (
    <header className="w-full px-4 py-6 z-10">
      <div className="flex items-center justify-between rounded-full bg-white/80 backdrop-blur-md border border-white shadow-lg px-4 py-2 max-w-[1024px] mx-auto">
        <div className="flex items-center gap-3">
          <div className="size-12 rounded-full border-2 border-warm-yellow overflow-hidden bg-white shadow-sm">
            {/* 头像图标 */}
            <img 
              src="/icon.png" 
              alt="TanGo" 
              className="w-full h-full object-cover"
            />
          </div>
          <h2 className="hidden sm:block text-text-main text-lg font-bold tracking-tight font-display">
            {title || t('header.title')}
          </h2>
        </div>
        <div className="flex items-center gap-3">
          {showFavorites && (
            <Link
              to="/collection"
              className="group flex items-center gap-2 rounded-full bg-white hover:bg-gray-50 border border-gray-100 px-4 py-2 transition-all shadow-sm"
            >
              <span className="material-symbols-outlined text-warm-yellow group-hover:scale-110 transition-transform">
                star
              </span>
              <span className="text-sm font-medium text-text-main hidden sm:block">{t('header.favorites')}</span>
            </Link>
          )}
          {showReport && (
            <Link
              to="/report"
              className="group flex items-center gap-2 rounded-full bg-white hover:bg-gray-50 border border-gray-100 px-4 py-2 transition-all shadow-sm"
            >
              <span className="material-symbols-outlined text-sky-blue group-hover:scale-110 transition-transform">
                menu_book
              </span>
              <span className="text-sm font-medium text-text-main hidden sm:block">{t('common.report')}</span>
            </Link>
          )}
          <Link
            to="/settings"
            className="group flex items-center gap-2 rounded-full bg-white hover:bg-gray-50 border border-gray-100 px-4 py-2 transition-all shadow-sm"
          >
            <span className="material-symbols-outlined text-slate-600 group-hover:scale-110 transition-transform">
              settings
            </span>
            <span className="text-sm font-medium text-text-main hidden sm:block">{t('settings.title')}</span>
          </Link>
        </div>
      </div>
    </header>
  );
};

