/**
 * 收藏页面组件
 * 基于 stitch_ui/favorites_page/ 设计稿
 */

import { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { Header } from '../components/common/Header';
import { CollectionGrid } from '../components/collection/CollectionGrid';
import { explorationStorage, cardStorage } from '../services/storage';
import { exportCardAsImage } from '../utils/export';
import { createShareLink, copyToClipboard } from '../utils/share';
import { isInCurrentWeek } from '../utils/week';
import type { ExplorationRecord } from '../types/exploration';
import type { KnowledgeCard } from '../types/exploration';

export default function Collection() {
  const { t } = useTranslation();
  const [records, setRecords] = useState<ExplorationRecord[]>([]);
  const [cards, setCards] = useState<KnowledgeCard[]>([]);
  const [loading, setLoading] = useState(true);
  const [isSharing, setIsSharing] = useState(false);
  const [shareSuccess, setShareSuccess] = useState(false);

  useEffect(() => {
    loadRecords();
  }, []);

  const loadRecords = async () => {
    try {
      const allRecords = await explorationStorage.getAll();
      // 只显示已收藏的记录
      const collectedRecords = allRecords.filter((r) => r.collected);
      setRecords(collectedRecords);
      
      // 加载所有收藏的卡片
      const allCards = await cardStorage.getAll();
      setCards(allCards);
    } catch (error) {
      console.error('加载收藏记录失败:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleReExplore = (recordId: string) => {
    // TODO: 导航到结果页面，显示该探索记录
    console.log('重新探索:', recordId);
  };

  const handleToggleCollect = async (recordId: string, collected: boolean) => {
    try {
      const record = records.find(r => r.id === recordId);
      if (record) {
        // 更新探索记录的collected字段
        record.collected = collected;
        await explorationStorage.save(record);
        
        // 如果取消收藏，从列表中移除该记录
        if (!collected) {
          setRecords(records.filter(r => r.id !== recordId));
        } else {
          // 如果收藏，添加到列表（如果不在列表中）
          if (!records.find(r => r.id === recordId)) {
            setRecords([...records, record]);
          }
        }
        
        // 重新加载数据以更新UI
        await loadRecords();
      }
    } catch (error) {
      console.error('切换收藏状态失败:', error);
    }
  };

  const handleExportCard = async (cardId: string) => {
    try {
      await exportCardAsImage(`card-${cardId}`, `card-${cardId}`);
    } catch (error) {
      console.error('导出卡片失败:', error);
      alert(t('collection.exportError'));
    }
  };

  const handleShare = async () => {
    if (isSharing) return;
    
    setIsSharing(true);
    setShareSuccess(false);

    try {
      // 获取所有探索记录和卡片（不限制收藏状态）
      const allRecords = await explorationStorage.getAll();
      const allCards = await cardStorage.getAll();

      // 过滤出当前周（周一到周日）的记录和卡片
      const currentWeekRecords = allRecords.filter((r) => isInCurrentWeek(r.timestamp));
      const currentWeekCards = allCards.filter((c) => {
        // 如果卡片有关联的探索记录，检查探索记录的时间
        const relatedRecord = allRecords.find((r) => r.id === c.explorationId);
        if (relatedRecord) {
          return isInCurrentWeek(relatedRecord.timestamp);
        }
        // 如果卡片有收藏时间，检查收藏时间
        if (c.collectedAt) {
          return isInCurrentWeek(c.collectedAt);
        }
        // 如果没有时间信息，不包含
        return false;
      });

      if (currentWeekRecords.length === 0 && currentWeekCards.length === 0) {
        alert(t('share.noDataToShare', '本周还没有探索记录，先去探索一些内容吧！'));
        return;
      }

      // 限制最多10条探索记录（按时间倒序，最新的10条）
      const sortedRecords = currentWeekRecords.sort((a, b) => 
        new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime()
      );
      const limitedRecords = sortedRecords.slice(0, 10);

      // 只包含这10条记录相关的卡片
      const recordIds = new Set(limitedRecords.map(r => r.id));
      const limitedCards = currentWeekCards.filter(c => 
        recordIds.has(c.explorationId)
      );

      // 创建分享链接（分享最多10条记录和相关的卡片）
      const shareUrl = await createShareLink(limitedRecords, limitedCards);

      // 复制到剪贴板
      const success = await copyToClipboard(shareUrl);
      
      if (success) {
        setShareSuccess(true);
        setTimeout(() => setShareSuccess(false), 3000);
      } else {
        // 如果复制失败，显示链接让用户手动复制
        const userConfirmed = confirm(
          `${t('share.linkCreated', '分享链接已创建')}:\n${shareUrl}\n\n${t('share.copyManually', '请手动复制链接')}`
        );
        if (userConfirmed) {
          setShareSuccess(true);
          setTimeout(() => setShareSuccess(false), 3000);
        }
      }
    } catch (error: any) {
      console.error('分享失败:', error);
      const errorMessage = error.message || t('share.shareError', '分享失败，请稍后重试');
      alert(errorMessage);
    } finally {
      setIsSharing(false);
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-cloud-white flex items-center justify-center">
        <div className="text-text-sub">{t('common.loading')}</div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-cloud-white font-display">
      <Header showFavorites={false} showReport={true} />

      <main className="flex-1 overflow-y-auto p-4 md:p-6 lg:px-10 lg:py-8 scroll-smooth z-10">
        <div className="max-w-6xl mx-auto flex flex-col gap-8 pb-10">
          {/* 页面头部 */}
          <header className="flex flex-col md:flex-row md:items-end justify-between gap-6 bg-white/60 backdrop-blur-sm p-6 rounded-3xl border border-white shadow-sm">
            <div className="flex flex-col gap-1">
              <h1 className="text-text-main text-3xl md:text-4xl font-extrabold tracking-tight flex items-center gap-3 font-display">
                <div className="bg-gradient-to-br from-orange-400 to-primary text-white p-3 rounded-2xl shadow-lg rotate-3">
                  <span className="material-symbols-outlined text-3xl">stars</span>
                </div>
                <span className="bg-clip-text text-transparent bg-gradient-to-r from-text-main to-orange-800">
                  {t('collection.title')}
                </span>
              </h1>
              <p className="text-text-sub text-sm md:text-base font-medium pl-1 mt-2">
                {t('collection.subtitle')}
              </p>
            </div>
            <button
              onClick={handleShare}
              disabled={isSharing || (records.length === 0 && cards.length === 0)}
              className={`flex items-center justify-center gap-2 rounded-2xl border-2 ${
                shareSuccess
                  ? 'border-green-300 bg-green-50 text-green-600'
                  : 'border-gray-200 bg-surface hover:bg-gray-50 hover:text-text-main hover:border-gray-300'
              } px-6 py-3 text-sm font-extrabold transition-all shadow-sm disabled:opacity-50 disabled:cursor-not-allowed`}
            >
              {shareSuccess ? (
                <>
                  <span className="material-symbols-outlined">check_circle</span>
                  {t('share.copied', '已复制！')}
                </>
              ) : (
                <>
                  <span className="material-symbols-outlined">{isSharing ? 'hourglass_empty' : 'share'}</span>
                  {isSharing ? t('share.creating', '创建中...') : t('share.shareWithParents', '分享给家长')}
                </>
              )}
            </button>
          </header>

          {/* 收藏卡片网格 */}
          <CollectionGrid
            records={records}
            cards={cards}
            category="all"
            onReExplore={handleReExplore}
            onExport={handleExportCard}
            onToggleCollect={handleToggleCollect}
          />

          {/* Little Star 鼓励消息 */}
          {records.length > 0 && (
            <div className="col-span-1 md:col-span-2 xl:col-span-3 mt-8">
              <div className="flex flex-row items-end gap-4 max-w-2xl mx-auto">
                <div className="bg-center bg-contain bg-no-repeat size-24 shrink-0 drop-shadow-lg">
                  <div className="w-full h-full bg-gradient-to-br from-yellow-200 to-pink-200 rounded-full flex items-center justify-center">
                    <span className="text-4xl">⭐</span>
                  </div>
                </div>
                <div className="relative mb-6">
                  <div className="bg-white p-6 rounded-3xl rounded-bl-none border-2 border-primary/20 shadow-lg relative z-10 max-w-lg">
                    <div className="flex items-start gap-3">
                      <span className="material-symbols-outlined text-primary text-3xl">lightbulb</span>
                      <div>
                        <h4 className="text-primary font-bold font-display text-lg mb-1">
                          {t('collection.littleStarSays')}
                        </h4>
                        <p className="text-text-main font-medium leading-relaxed">
                          {t('collection.littleStarMessage')}
                        </p>
                      </div>
                    </div>
                  </div>
                  <svg
                    className="absolute -bottom-2 -left-2 w-8 h-8 text-white transform rotate-0 z-20"
                    fill="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path d="M0 0 L24 0 L24 24 Z" transform="scale(1, -1) translate(0, -24)" />
                  </svg>
                </div>
              </div>
            </div>
          )}
        </div>
      </main>
    </div>
  );
}
