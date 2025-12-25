/**
 * 图片loading占位组件
 * 显示图片生成进度和loading动画
 */

import { useState, useEffect } from 'react';

export interface ImagePlaceholderProps {
  imageUrl?: string;
  progress?: number; // 0-100
  onLoad?: () => void;
  className?: string;
}

export function ImagePlaceholder({
  imageUrl,
  progress = 0,
  onLoad,
  className = '',
}: ImagePlaceholderProps) {
  const [isLoading, setIsLoading] = useState(true);
  const [displayProgress, setDisplayProgress] = useState(0);

  // 更新进度显示
  useEffect(() => {
    setDisplayProgress(Math.min(100, Math.max(0, progress)));
  }, [progress]);

  // 预加载图片
  useEffect(() => {
    if (imageUrl) {
      const img = new Image();
      img.onload = () => {
        setIsLoading(false);
        setDisplayProgress(100);
        onLoad?.();
      };
      img.onerror = () => {
        setIsLoading(false);
        setDisplayProgress(0);
      };
      img.src = imageUrl;
    } else {
      setIsLoading(true);
    }
  }, [imageUrl, onLoad]);

  // 如果图片已加载完成，显示图片
  if (imageUrl && !isLoading) {
    return (
      <img
        src={imageUrl}
        alt="Generated"
        className={`rounded-lg ${className}`}
      />
    );
  }

  // 显示loading占位符
  return (
    <div
      className={`image-placeholder bg-gray-100 rounded-lg p-8 flex flex-col items-center justify-center min-h-[200px] ${className}`}
    >
      {/* Loading spinner */}
      <div className="loading-spinner animate-spin rounded-full h-12 w-12 border-b-2 border-[var(--color-primary)] mb-4" />

      {/* Progress bar */}
      <div className="w-full max-w-xs bg-gray-200 rounded-full h-2 mb-2">
        <div
          className="bg-[var(--color-primary)] h-2 rounded-full transition-all duration-300"
          style={{ width: `${displayProgress}%` }}
        />
      </div>

      {/* Progress text */}
      <p className="text-sm text-gray-600">
        正在生成图片... {displayProgress}%
      </p>
    </div>
  );
}

