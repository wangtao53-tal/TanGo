/**
 * 学习报告页面组件
 * 基于 stitch_ui/learning_report_page/ 设计稿
 */

import { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import { Header } from '../components/common/Header';
import { explorationStorage, cardStorage } from '../services/storage';
import { createShareLink, copyToClipboard } from '../utils/share';
import { isInCurrentWeek } from '../utils/week';
import { getUserStats } from '../services/badge';
import type { UserStats } from '../types/badge';

export default function LearningReport() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const [totalExplorations, setTotalExplorations] = useState(0);
  const [totalCollectedCards, setTotalCollectedCards] = useState(0);
  const [categoryDistribution, setCategoryDistribution] = useState<Record<string, number>>({
    自然类: 0,
    生活类: 0,
    人文类: 0,
  });
  const [isSharing, setIsSharing] = useState(false);
  const [shareSuccess, setShareSuccess] = useState(false);
  const [badgeStats, setBadgeStats] = useState<UserStats | null>(null);
  const [badgeLoading, setBadgeLoading] = useState(true);

  useEffect(() => {
    loadReportData();
    loadBadgeStats();
  }, []);

  const loadBadgeStats = async () => {
    try {
      const stats = await getUserStats();
      setBadgeStats(stats);
    } catch (error) {
      console.error('加载勋章数据失败:', error);
    } finally {
      setBadgeLoading(false);
    }
  };

  const loadReportData = async () => {
    try {
      const records = await explorationStorage.getAll();
      const cards = await cardStorage.getAll();

      setTotalExplorations(records.length);
      setTotalCollectedCards(cards.length);

      // 计算类别分布
      const distribution: Record<string, number> = {
        自然类: 0,
        生活类: 0,
        人文类: 0,
      };
      
      records.forEach((r) => {
        // 确保objectCategory有效
        const category = r.objectCategory || '自然类'; // 默认值
        if (['自然类', '生活类', '人文类'].includes(category)) {
          distribution[category] = (distribution[category] || 0) + 1;
        } else {
          // 如果分类无效，使用默认值
          console.warn('无效的分类值，使用默认值"自然类":', category, '记录ID:', r.id);
          distribution['自然类'] = (distribution['自然类'] || 0) + 1;
        }
      });
      
      setCategoryDistribution(distribution);
      
      // 验证数据一致性
      const totalCategories = Object.values(distribution).reduce((a, b) => a + b, 0);
      if (totalCategories !== records.length) {
        console.warn('数据不一致：知识地图总数与探索次数不匹配', {
          totalCategories,
          totalExplorations: records.length,
        });
      }
    } catch (error) {
      console.error('加载报告数据失败:', error);
    }
  };

  const totalCategories = Object.values(categoryDistribution).reduce((a, b) => a + b, 0);
  const naturalPercent = totalCategories > 0 ? (categoryDistribution['自然类'] / totalCategories) * 100 : 0;
  const lifePercent = totalCategories > 0 ? (categoryDistribution['生活类'] / totalCategories) * 100 : 0;
  const humanitiesPercent = totalCategories > 0 ? (categoryDistribution['人文类'] / totalCategories) * 100 : 0;

  const handleShare = async () => {
    if (isSharing) return;
    
    setIsSharing(true);
    setShareSuccess(false);

    try {
      // 获取所有探索记录和卡片
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

  return (
    <div className="min-h-screen bg-cloud-white font-display">
      <Header />

      <main className="flex-1 px-4 py-8 md:px-10 lg:px-20">
        <div className="mx-auto flex max-w-[1024px] flex-col gap-8">
          {/* 报告头部 */}
          <div className="relative flex flex-col gap-4 md:flex-row md:items-end md:justify-between p-6 bg-white rounded-3xl border-2 border-gray-100 shadow-card">
            <div className="absolute top-0 right-0 p-10 opacity-5 pointer-events-none">
              <span className="material-symbols-outlined text-9xl rotate-12">sunny</span>
            </div>
            <div className="relative z-10">
              <div className="flex items-center gap-2 mb-2">
                <span className="inline-flex items-center justify-center rounded-full bg-sky-blue/10 px-3 py-1 text-xs font-extrabold text-sky-blue uppercase tracking-wide">
                  {t('report.weeklyReport')}
                </span>
              </div>
              <h1 className="text-4xl font-black leading-tight tracking-tight text-text-main md:text-5xl">
                {t('report.greeting')} 🌟
              </h1>
              <p className="mt-2 text-lg font-medium text-text-sub">
                {t('report.subtitle')}
              </p>
            </div>
            <div className="relative z-10 mt-4 md:mt-0">
              <span className="inline-flex items-center gap-2 rounded-2xl bg-warm-yellow/20 px-5 py-3 text-sm font-bold text-text-main border-2 border-warm-yellow/30">
                <span className="material-symbols-outlined text-[20px] text-primary">calendar_month</span>
                {new Date().toLocaleDateString('zh-CN', { month: 'long', day: 'numeric' })}
              </span>
            </div>
          </div>

          {/* 统计卡片 */}
          <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
            {/* 探索次数 */}
            <div 
              onClick={() => navigate('/capture')}
              className="group relative overflow-hidden rounded-3xl bg-white p-6 transition-all hover:-translate-y-1 shadow-card border-2 border-warm-yellow/20 hover:border-warm-yellow cursor-pointer"
            >
              <div className="absolute -right-6 -top-6 h-32 w-32 rounded-full bg-warm-yellow/10 transition-all group-hover:scale-110" />
              <div className="flex items-start justify-between relative z-10">
                <div className="flex flex-col gap-2">
                  <p className="text-xs font-bold uppercase tracking-wider text-text-sub">{t('report.explorationStars')}</p>
                  <p className="text-5xl font-black text-text-main">{totalExplorations}</p>
                  <div className="inline-flex items-center gap-1 rounded-full bg-warm-yellow/20 px-3 py-1 text-xs font-bold text-text-main w-fit">
                    <span className="material-symbols-outlined text-sm font-bold">arrow_upward</span>
                    {t('report.keepExploring')}
                  </div>
                </div>
                <div className="flex h-16 w-16 items-center justify-center rounded-2xl bg-warm-yellow text-white shadow-md rotate-3 group-hover:rotate-12 transition-transform">
                  <span className="material-symbols-outlined text-4xl fill-1">star</span>
                </div>
              </div>
            </div>

            {/* 收藏卡片数 */}
            <div 
              onClick={() => navigate('/collection')}
              className="group relative overflow-hidden rounded-3xl bg-white p-6 transition-all hover:-translate-y-1 shadow-card border-2 border-peach-pink/20 hover:border-peach-pink cursor-pointer"
            >
              <div className="absolute -right-6 -top-6 h-32 w-32 rounded-full bg-peach-pink/10 transition-all group-hover:scale-110" />
              <div className="flex items-start justify-between relative z-10">
                <div className="flex flex-col gap-2">
                  <p className="text-xs font-bold uppercase tracking-wider text-text-sub">{t('report.totalFavorites')}</p>
                  <p className="text-5xl font-black text-text-main">{totalCollectedCards}</p>
                  <div className="inline-flex items-center gap-1 rounded-full bg-peach-pink/20 px-3 py-1 text-xs font-bold text-peach-pink w-fit">
                    <span className="material-symbols-outlined text-sm font-bold">favorite</span>
                    {t('report.greatCollection')}
                  </div>
                </div>
                <div className="flex h-16 w-16 items-center justify-center rounded-2xl bg-peach-pink text-white shadow-md -rotate-3 group-hover:-rotate-12 transition-transform">
                  <span className="material-symbols-outlined text-4xl fill-1">emoji_events</span>
                </div>
              </div>
            </div>

            {/* 专家等级 */}
            <div 
              onClick={() => navigate('/badge')}
              className="group relative overflow-hidden rounded-3xl bg-white p-6 transition-all hover:-translate-y-1 shadow-card border-2 border-science-green/20 hover:border-science-green cursor-pointer"
            >
              <div className="absolute -right-6 -top-6 h-32 w-32 rounded-full bg-science-green/10 transition-all group-hover:scale-110" />
              <div className="flex items-start justify-between relative z-10">
                <div className="flex flex-col gap-2">
                  <p className="text-xs font-bold uppercase tracking-wider text-text-sub">{t('report.littleExpert')}</p>
                  {badgeLoading ? (
                    <p className="text-3xl font-black leading-tight text-text-main mt-1">
                      {t('report.loading', '加载中...')}
                    </p>
                  ) : badgeStats ? (
                    <>
                      <p 
                        className="text-3xl font-black leading-tight text-text-main mt-1"
                        style={{ color: badgeStats.currentLevelInfo.color }}
                      >
                        {badgeStats.currentLevelInfo.title}
                      </p>
                      {badgeStats.nextLevelInfo && (
                        <div className="mt-1 inline-flex items-center gap-1 rounded-full bg-science-green/20 px-3 py-1 text-xs font-bold text-science-green w-fit">
                          {t('report.levelUp')} 🚀
                        </div>
                      )}
                    </>
                  ) : (
                    <p className="text-3xl font-black leading-tight text-text-main mt-1">
                      {t('report.natureMaster')}
                    </p>
                  )}
                </div>
                <div 
                  className="flex h-16 w-16 items-center justify-center rounded-2xl text-white shadow-md rotate-3 group-hover:rotate-6 transition-transform"
                  style={{ 
                    backgroundColor: badgeStats?.currentLevelInfo.color || '#76FF7A'
                  }}
                >
                  {badgeStats?.currentLevelInfo.icon ? (
                    <span className="text-4xl">{badgeStats.currentLevelInfo.icon}</span>
                  ) : (
                    <span className="material-symbols-outlined text-4xl fill-1">forest</span>
                  )}
                </div>
              </div>
            </div>
          </div>

          {/* 知识地图 */}
          <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
            <div className="flex flex-col justify-between rounded-3xl bg-white p-8 shadow-card border-2 border-gray-100">
              <div className="mb-8 flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <div className="h-8 w-2 rounded-full bg-primary" />
                  <h3 className="text-xl font-extrabold text-text-main">{t('report.knowledgeMap')}</h3>
                </div>
              </div>
              <div className="flex flex-col items-center sm:flex-row sm:justify-center gap-10">
                <div
                  className="relative flex h-52 w-52 shrink-0 items-center justify-center rounded-full shadow-lg ring-8 ring-gray-50"
                  style={{
                    background: `conic-gradient(${naturalPercent > 0 ? '#76FF7A' : '#e5e7eb'} 0% ${naturalPercent}%, ${lifePercent > 0 ? '#FF9E64' : '#e5e7eb'} ${naturalPercent}% ${naturalPercent + lifePercent}%, ${humanitiesPercent > 0 ? '#40C4FF' : '#e5e7eb'} ${naturalPercent + lifePercent}% 100%)`,
                  }}
                >
                  <div className="absolute h-36 w-36 rounded-full bg-white flex items-center justify-center flex-col shadow-inner">
                    <span className="text-xs text-text-sub font-bold uppercase tracking-widest">{t('report.total')}</span>
                    <span className="text-4xl font-black text-text-main">{totalCategories}</span>
                  </div>
                  <div
                    className="absolute -left-4 top-8 flex h-10 w-10 items-center justify-center rounded-full bg-science-green text-white shadow-md border-4 border-white animate-bounce-slow"
                    style={{ animationDelay: '0s' }}
                  >
                    <span className="material-symbols-outlined text-lg">eco</span>
                  </div>
                  <div
                    className="absolute -right-4 top-8 flex h-10 w-10 items-center justify-center rounded-full bg-peach-pink text-white shadow-md border-4 border-white animate-bounce-slow"
                    style={{ animationDelay: '1s' }}
                  >
                    <span className="material-symbols-outlined text-lg">pets</span>
                  </div>
                  <div
                    className="absolute bottom-0 left-1/2 -translate-x-1/2 translate-y-1/2 flex h-10 w-10 items-center justify-center rounded-full bg-sky-blue text-white shadow-md border-4 border-white animate-bounce-slow"
                    style={{ animationDelay: '2s' }}
                  >
                    <span className="material-symbols-outlined text-lg">menu_book</span>
                  </div>
                </div>
                <div className="flex flex-col gap-4">
                  <div className="flex items-center gap-3">
                    <div className="h-4 w-4 rounded-full bg-science-green" />
                    <div className="flex flex-col">
                      <span className="text-sm font-bold text-text-main">{t('report.categoryNatural')}</span>
                      <span className="text-xs text-text-sub">{categoryDistribution['自然类']} {t('report.items')}</span>
                    </div>
                  </div>
                  <div className="flex items-center gap-3">
                    <div className="h-4 w-4 rounded-full bg-sunny-orange" />
                    <div className="flex flex-col">
                      <span className="text-sm font-bold text-text-main">{t('report.categoryLife')}</span>
                      <span className="text-xs text-text-sub">{categoryDistribution['生活类']} {t('report.items')}</span>
                    </div>
                  </div>
                  <div className="flex items-center gap-3">
                    <div className="h-4 w-4 rounded-full bg-sky-blue" />
                    <div className="flex flex-col">
                      <span className="text-sm font-bold text-text-main">{t('report.categoryHumanities')}</span>
                      <span className="text-xs text-text-sub">{categoryDistribution['人文类']} {t('report.items')}</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            {/* 最近收藏 */}
            <div className="flex flex-col rounded-3xl bg-white p-8 shadow-card border-2 border-gray-100">
              <div className="mb-6 flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <div className="h-8 w-2 rounded-full bg-primary" />
                  <h3 className="text-xl font-extrabold text-text-main">{t('report.recentFavorites')}</h3>
                </div>
              </div>
              <div className="flex flex-col gap-4">
                {totalCollectedCards > 0 ? (
                  <p className="text-text-sub">{t('report.recentFavoritesMessage', { totalCollectedCards })}</p>
                ) : (
                  <div className="flex flex-col items-center justify-center py-10">
                    <span className="text-6xl mb-4">📚</span>
                    <p className="text-text-sub">{t('report.noCards')}</p>
                  </div>
                )}
              </div>
            </div>
          </div>

          {/* 分享区域 */}
          <div className="flex flex-col items-center justify-center gap-6 rounded-3xl border-2 border-dashed border-primary/30 bg-primary/5 p-8 md:flex-row md:justify-between relative overflow-hidden">
            <div className="absolute -left-10 -bottom-10 h-32 w-32 rounded-full bg-primary/10 blur-2xl"></div>
            <div className="flex items-center gap-5 relative z-10">
              <div className="flex h-14 w-14 items-center justify-center rounded-2xl bg-white text-primary shadow-comic border-2 border-primary/20 rotate-3">
                <span className="material-symbols-outlined text-3xl">auto_fix_high</span>
              </div>
              <div className="flex flex-col">
                <h4 className="text-xl font-black text-text-main">{t('share.readyToShare', '准备好分享了吗？')}</h4>
                <p className="text-sm font-bold text-text-muted">{t('share.shareWithParents', '与家长分享你的探索成果！')}</p>
              </div>
            </div>
            <div className="flex w-full flex-col gap-4 sm:w-auto sm:flex-row relative z-10">
              <button
                onClick={handleShare}
                disabled={isSharing}
                className={`flex items-center justify-center gap-2 rounded-2xl border-2 ${
                  shareSuccess
                    ? 'border-green-300 bg-green-50 text-green-600'
                    : 'border-gray-200 bg-surface hover:bg-gray-50 hover:text-text-main hover:border-gray-300'
                } px-6 py-3.5 text-sm font-extrabold transition-all shadow-sm sm:w-auto disabled:opacity-50 disabled:cursor-not-allowed`}
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
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}
