/**
 * 学习报告页面组件
 * 基于 stitch_ui/learning_report_page/ 设计稿
 */

import { useState, useEffect } from 'react';
import { Header } from '../components/common/Header';
import { explorationStorage, cardStorage } from '../services/storage';

export default function LearningReport() {
  const [totalExplorations, setTotalExplorations] = useState(0);
  const [totalCollectedCards, setTotalCollectedCards] = useState(0);
  const [categoryDistribution, setCategoryDistribution] = useState<Record<string, number>>({
    自然类: 0,
    生活类: 0,
    人文类: 0,
  });

  useEffect(() => {
    loadReportData();
  }, []);

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
        distribution[r.objectCategory] = (distribution[r.objectCategory] || 0) + 1;
      });
      setCategoryDistribution(distribution);
    } catch (error) {
      console.error('加载报告数据失败:', error);
    }
  };

  const totalCategories = Object.values(categoryDistribution).reduce((a, b) => a + b, 0);
  const naturalPercent = totalCategories > 0 ? (categoryDistribution['自然类'] / totalCategories) * 100 : 0;
  const lifePercent = totalCategories > 0 ? (categoryDistribution['生活类'] / totalCategories) * 100 : 0;
  const humanitiesPercent = totalCategories > 0 ? (categoryDistribution['人文类'] / totalCategories) * 100 : 0;

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
                  Weekly Report
                </span>
              </div>
              <h1 className="text-4xl font-black leading-tight tracking-tight text-text-main md:text-5xl">
                Hi, <span className="text-sky-blue inline-block hover:scale-105 transition-transform cursor-default">Little Explorer!</span> 🌟
              </h1>
              <p className="mt-2 text-lg font-medium text-text-sub">
                You're doing great! Look at your growth this week.
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
            <div className="group relative overflow-hidden rounded-3xl bg-white p-6 transition-all hover:-translate-y-1 shadow-card border-2 border-warm-yellow/20 hover:border-warm-yellow">
              <div className="absolute -right-6 -top-6 h-32 w-32 rounded-full bg-warm-yellow/10 transition-all group-hover:scale-110" />
              <div className="flex items-start justify-between relative z-10">
                <div className="flex flex-col gap-2">
                  <p className="text-xs font-bold uppercase tracking-wider text-text-sub">Exploration Stars</p>
                  <p className="text-5xl font-black text-text-main">{totalExplorations}</p>
                  <div className="inline-flex items-center gap-1 rounded-full bg-warm-yellow/20 px-3 py-1 text-xs font-bold text-text-main w-fit">
                    <span className="material-symbols-outlined text-sm font-bold">arrow_upward</span>
                    Keep exploring!
                  </div>
                </div>
                <div className="flex h-16 w-16 items-center justify-center rounded-2xl bg-warm-yellow text-white shadow-md rotate-3 group-hover:rotate-12 transition-transform">
                  <span className="material-symbols-outlined text-4xl fill-1">star</span>
                </div>
              </div>
            </div>

            {/* 收藏卡片数 */}
            <div className="group relative overflow-hidden rounded-3xl bg-white p-6 transition-all hover:-translate-y-1 shadow-card border-2 border-peach-pink/20 hover:border-peach-pink">
              <div className="absolute -right-6 -top-6 h-32 w-32 rounded-full bg-peach-pink/10 transition-all group-hover:scale-110" />
              <div className="flex items-start justify-between relative z-10">
                <div className="flex flex-col gap-2">
                  <p className="text-xs font-bold uppercase tracking-wider text-text-sub">Total Favorites</p>
                  <p className="text-5xl font-black text-text-main">{totalCollectedCards}</p>
                  <div className="inline-flex items-center gap-1 rounded-full bg-peach-pink/20 px-3 py-1 text-xs font-bold text-peach-pink w-fit">
                    <span className="material-symbols-outlined text-sm font-bold">favorite</span>
                    Great collection!
                  </div>
                </div>
                <div className="flex h-16 w-16 items-center justify-center rounded-2xl bg-peach-pink text-white shadow-md -rotate-3 group-hover:-rotate-12 transition-transform">
                  <span className="material-symbols-outlined text-4xl fill-1">emoji_events</span>
                </div>
              </div>
            </div>

            {/* 专家等级 */}
            <div className="group relative overflow-hidden rounded-3xl bg-white p-6 transition-all hover:-translate-y-1 shadow-card border-2 border-science-green/20 hover:border-science-green">
              <div className="absolute -right-6 -top-6 h-32 w-32 rounded-full bg-science-green/10 transition-all group-hover:scale-110" />
              <div className="flex items-start justify-between relative z-10">
                <div className="flex flex-col gap-2">
                  <p className="text-xs font-bold uppercase tracking-wider text-text-sub">Little Expert</p>
                  <p className="text-3xl font-black leading-tight text-text-main mt-1">
                    Nature<br />Master
                  </p>
                  <div className="mt-1 inline-flex items-center gap-1 rounded-full bg-science-green/20 px-3 py-1 text-xs font-bold text-science-green w-fit">
                    Level Up! 🚀
                  </div>
                </div>
                <div className="flex h-16 w-16 items-center justify-center rounded-2xl bg-science-green text-white shadow-md rotate-3 group-hover:rotate-6 transition-transform">
                  <span className="material-symbols-outlined text-4xl fill-1">forest</span>
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
                  <h3 className="text-xl font-extrabold text-text-main">Knowledge Map</h3>
                </div>
                <button className="rounded-full bg-gray-100 p-2 text-text-sub hover:bg-gray-200 hover:text-text-main transition-colors">
                  <span className="material-symbols-outlined text-xl">info</span>
                </button>
              </div>
              <div className="flex flex-col items-center sm:flex-row sm:justify-center gap-10">
                <div
                  className="relative flex h-52 w-52 shrink-0 items-center justify-center rounded-full shadow-lg ring-8 ring-gray-50"
                  style={{
                    background: `conic-gradient(${naturalPercent > 0 ? '#76FF7A' : '#e5e7eb'} 0% ${naturalPercent}%, ${lifePercent > 0 ? '#FF9E64' : '#e5e7eb'} ${naturalPercent}% ${naturalPercent + lifePercent}%, ${humanitiesPercent > 0 ? '#40C4FF' : '#e5e7eb'} ${naturalPercent + lifePercent}% 100%)`,
                  }}
                >
                  <div className="absolute h-36 w-36 rounded-full bg-white flex items-center justify-center flex-col shadow-inner">
                    <span className="text-xs text-text-sub font-bold uppercase tracking-widest">Total</span>
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
                      <span className="text-sm font-bold text-text-main">Natural</span>
                      <span className="text-xs text-text-sub">{categoryDistribution['自然类']} items</span>
                    </div>
                  </div>
                  <div className="flex items-center gap-3">
                    <div className="h-4 w-4 rounded-full bg-sunny-orange" />
                    <div className="flex flex-col">
                      <span className="text-sm font-bold text-text-main">Life</span>
                      <span className="text-xs text-text-sub">{categoryDistribution['生活类']} items</span>
                    </div>
                  </div>
                  <div className="flex items-center gap-3">
                    <div className="h-4 w-4 rounded-full bg-sky-blue" />
                    <div className="flex flex-col">
                      <span className="text-sm font-bold text-text-main">Humanities</span>
                      <span className="text-xs text-text-sub">{categoryDistribution['人文类']} items</span>
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
                  <h3 className="text-xl font-extrabold text-text-main">Recent Favorites</h3>
                </div>
              </div>
              <div className="flex flex-col gap-4">
                {totalCollectedCards > 0 ? (
                  <p className="text-text-sub">最近收藏了 {totalCollectedCards} 张卡片</p>
                ) : (
                  <div className="flex flex-col items-center justify-center py-10">
                    <span className="text-6xl mb-4">📚</span>
                    <p className="text-text-sub">还没有收藏任何卡片</p>
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}
