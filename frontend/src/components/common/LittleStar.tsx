/**
 * Little Star对话气泡组件
 * 参考设计稿中的底部对话气泡
 */

import React from 'react';
import { useTranslation } from 'react-i18next';

export interface LittleStarProps {
  message: string;
  className?: string;
}

export const LittleStar: React.FC<LittleStarProps> = ({ message, className = '' }) => {
  const { t } = useTranslation();
  return (
    <div className={`fixed bottom-6 w-full max-w-[600px] px-4 z-20 pointer-events-none ${className}`}>
      <div className="flex items-end gap-4 pointer-events-auto animate-float" style={{ animationDelay: '1s' }}>
        <div className="relative shrink-0">
          <div className="absolute -inset-1 rounded-full bg-peach-pink blur opacity-40"></div>
          <div className="relative size-16 sm:size-20 rounded-full border-4 border-peach-pink bg-white overflow-hidden shadow-lg">
            {/* Little Star头像 - 使用占位符，后续替换为实际图片 */}
            <div className="w-full h-full bg-gradient-to-br from-yellow-200 to-pink-200 flex items-center justify-center">
              <span className="text-2xl">⭐</span>
            </div>
          </div>
        </div>
        <div className="flex flex-col gap-1 mb-4">
          <p className="text-text-sub text-xs sm:text-sm font-bold ml-4">{t('littleStar.name')}</p>
          <div className="relative bg-white text-text-main px-6 py-4 rounded-2xl rounded-bl-none shadow-[0_10px_30px_-5px_rgba(0,0,0,0.1)] border-2 border-peach-pink/50">
            <p className="text-base sm:text-lg font-bold leading-normal">{message}</p>
            <div className="absolute -bottom-[2px] -left-[2px] w-4 h-4 bg-white border-b-2 border-l-2 border-peach-pink/50 transform rotate-45 translate-y-1/2 -translate-x-1/2 rounded-bl-md"></div>
          </div>
        </div>
      </div>
    </div>
  );
};

