/**
 * 收藏页面组件
 * 基于 stitch_ui/favorites_page/ 设计稿
 */

import { useState, useEffect } from 'react';
import { Header } from '../components/common/Header';
import { CollectionGrid } from '../components/collection/CollectionGrid';
import { CategoryFilter, type Category } from '../components/collection/CategoryFilter';
import { explorationStorage, cardStorage } from '../services/storage';
import { exportCardAsImage } from '../utils/export';
import type { ExplorationRecord } from '../types/exploration';
import type { KnowledgeCard } from '../types/exploration';

export default function Collection() {
  const [records, setRecords] = useState<ExplorationRecord[]>([]);
  const [cards, setCards] = useState<KnowledgeCard[]>([]);
  const [selectedCategory, setSelectedCategory] = useState<Category>('all');
  const [loading, setLoading] = useState(true);

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

  const handleClearAll = () => {
    // TODO: 实现清空所有收藏功能（需要家长模式）
    console.log('清空所有收藏');
  };

  const handleExportCard = async (cardId: string) => {
    try {
      await exportCardAsImage(`card-${cardId}`, `card-${cardId}`);
    } catch (error) {
      console.error('导出卡片失败:', error);
      alert('导出失败，请重试');
    }
  };

  const handleExportAll = async () => {
    try {
      for (const card of cards) {
        await exportCardAsImage(`card-${card.id}`, `card-${card.id}`);
        // 添加延迟，避免浏览器阻止多个下载
        await new Promise((resolve) => setTimeout(resolve, 500));
      }
    } catch (error) {
      console.error('批量导出失败:', error);
      alert('导出失败，请重试');
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-cloud-white flex items-center justify-center">
        <div className="text-text-sub">加载中...</div>
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
                  My Favorites
                </span>
              </h1>
              <p className="text-text-sub text-sm md:text-base font-medium pl-1 mt-2">
                Keep exploring your collection of wonders!
              </p>
            </div>

            {/* 操作按钮组 */}
            <div className="flex items-center gap-2">
              {/* 导出所有按钮 */}
              {cards.length > 0 && (
                <button
                  onClick={handleExportAll}
                  className="group flex items-center justify-center gap-2 bg-[var(--color-primary)] hover:bg-[#5aff2b] text-white px-4 py-2 rounded-full transition-all shadow-md hover:shadow-lg"
                  title="导出所有卡片"
                >
                  <span className="material-symbols-outlined text-lg">download</span>
                  <span className="text-sm font-bold font-display hidden sm:inline">导出全部</span>
                </button>
              )}
              
              {/* 家长模式控制 */}
              <div className="flex items-center gap-2 bg-white p-1.5 pr-3 rounded-full border border-gray-100 shadow-sm">
                <div className="flex items-center gap-2 pl-3 pr-3 py-1.5 bg-gray-100 rounded-full">
                  <span className="material-symbols-outlined text-gray-500 text-sm">lock_open</span>
                  <span className="text-xs font-bold text-gray-500 uppercase tracking-wider font-display">
                    Parent Mode
                  </span>
                </div>
                <button
                  onClick={handleClearAll}
                  className="group flex items-center justify-center gap-1.5 text-red-400 hover:text-red-500 hover:bg-red-50 px-3 py-1.5 rounded-full transition-all cursor-pointer"
                  title="Only available in Parent Mode"
                >
                  <span className="material-symbols-outlined text-lg">delete</span>
                  <span className="text-sm font-bold font-display">Clear All</span>
                </button>
              </div>
            </div>
          </header>

          {/* 分类筛选 */}
          <CategoryFilter selected={selectedCategory} onSelect={setSelectedCategory} />

          {/* 收藏卡片网格 */}
          <CollectionGrid
            records={records}
            cards={cards}
            category={selectedCategory}
            onReExplore={handleReExplore}
            onExport={handleExportCard}
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
                          Little Star Says:
                        </h4>
                        <p className="text-text-main font-medium leading-relaxed">
                          Go explore interesting knowledge and collect more favorite cards! I'm waiting for your discoveries! ✨
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
