/**
 * 卡片导出Hook
 * 处理长按和右键菜单触发导出
 */

import { useRef, useCallback, useEffect, useState } from 'react';
import { exportCardAsImage } from '../utils/cardExport';

export interface UseCardExportOptions {
  cardId: string; // 卡片ID，用于生成elementId和filename
  onSuccess?: () => void; // 导出成功回调
  onError?: (error: Error) => void; // 导出失败回调
  longPressDelay?: number; // 长按延迟时间（毫秒，默认500ms）
}

export interface UseCardExportReturn {
  onTouchStart: (e: React.TouchEvent) => void;
  onTouchEnd: () => void;
  onTouchCancel: () => void;
  onContextMenu: (e: React.MouseEvent) => void;
  isExporting: boolean;
}

export function useCardExport({
  cardId,
  onSuccess,
  onError,
  longPressDelay = 500,
}: UseCardExportOptions): UseCardExportReturn {
  const longPressTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const isExportingRef = useRef(false);
  const [isExporting, setIsExporting] = useState(false);

  // 清理定时器
  const clearLongPressTimer = useCallback(() => {
    if (longPressTimerRef.current) {
      clearTimeout(longPressTimerRef.current);
      longPressTimerRef.current = null;
    }
  }, []);

  // 导出卡片
  const handleExport = useCallback(async () => {
    // 防止重复触发
    if (isExportingRef.current) {
      return;
    }

    isExportingRef.current = true;
    setIsExporting(true);
    
    const elementId = `card-${cardId}`;
    const filename = `card-${cardId}-${Date.now()}`;

    try {
      await exportCardAsImage(elementId, { filename });
      if (onSuccess) {
        onSuccess();
      }
    } catch (error) {
      console.error('导出卡片失败:', error);
      if (onError) {
        onError(error as Error);
      }
    } finally {
      isExportingRef.current = false;
      setIsExporting(false);
    }
  }, [cardId, onSuccess, onError]);

  // 移动端长按处理
  const onTouchStart = useCallback((_e: React.TouchEvent) => {
    clearLongPressTimer();
    
    longPressTimerRef.current = setTimeout(() => {
      handleExport();
    }, longPressDelay);
  }, [handleExport, longPressDelay, clearLongPressTimer]);

  const onTouchEnd = useCallback(() => {
    clearLongPressTimer();
  }, [clearLongPressTimer]);

  const onTouchCancel = useCallback(() => {
    clearLongPressTimer();
  }, [clearLongPressTimer]);

  // PC端右键菜单处理
  const onContextMenu = useCallback((e: React.MouseEvent) => {
    e.preventDefault(); // 阻止默认右键菜单
    handleExport();
  }, [handleExport]);

  // 组件卸载时清理定时器
  useEffect(() => {
    return () => {
      clearLongPressTimer();
    };
  }, [clearLongPressTimer]);

  return {
    onTouchStart,
    onTouchEnd,
    onTouchCancel,
    onContextMenu,
    isExporting,
  };
}
