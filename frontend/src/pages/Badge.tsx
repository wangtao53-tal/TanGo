/**
 * 勋章页面
 */

import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Header } from '../components/common/Header';
import { getBadgeDetail, getUserStats } from '../services/badge';
import type { BadgeDetailResponse, BadgeLevel, UserStats } from '../types/badge';

export default function Badge() {
  const navigate = useNavigate();
  const { t } = useTranslation();
  const [loading, setLoading] = useState(true);
  const [stats, setStats] = useState<UserStats | null>(null);
  const [allLevels, setAllLevels] = useState<BadgeLevel[]>([]);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadBadgeData();
  }, []);

  const loadBadgeData = async () => {
    try {
      setLoading(true);
      const data: BadgeDetailResponse = await getBadgeDetail();
      setStats(data.stats);
      setAllLevels(data.allLevels);
      setError(null);
    } catch (err) {
      console.error('加载勋章数据失败:', err);
      setError('加载勋章数据失败，请稍后重试');
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gradient-to-b from-green-50 to-white">
        <Header title="我的勋章" />
        <div className="flex items-center justify-center h-64">
          <div className="text-green-600">加载中...</div>
        </div>
      </div>
    );
  }

  if (error || !stats) {
    return (
      <div className="min-h-screen bg-gradient-to-b from-green-50 to-white">
        <Header title="我的勋章" />
        <div className="flex items-center justify-center h-64">
          <div className="text-red-600">{error || '数据加载失败'}</div>
        </div>
      </div>
    );
  }

  const currentLevelInfo = stats.currentLevelInfo;
  const nextLevelInfo = stats.nextLevelInfo;

  return (
    <div className="min-h-screen bg-gradient-to-b from-green-50 to-white">
      <Header title="我的勋章" />
      
      <main className="flex-1 px-4 py-6 md:px-6 lg:px-8">
        <div className="mx-auto flex max-w-[1024px] flex-col gap-6 lg:gap-8">
          {/* PC端：两列布局，移动端：单列 */}
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            {/* 当前勋章卡片 */}
            <div 
              className="lg:col-span-1 bg-white rounded-2xl shadow-lg p-4 sm:p-5 lg:p-6 border-2"
              style={{ borderColor: currentLevelInfo.color }}
            >
              <div className="flex items-center justify-between mb-3 sm:mb-4">
                <div className="flex-1">
                  <div className="text-xs sm:text-sm text-gray-500 mb-1">当前等级</div>
                  <div className="text-xl sm:text-2xl font-bold mb-1 sm:mb-2" style={{ color: currentLevelInfo.color }}>
                    {currentLevelInfo.title}
                  </div>
                  <div className="text-xs text-gray-400">{currentLevelInfo.description}</div>
                </div>
                <div className="text-4xl sm:text-5xl lg:text-6xl ml-2">{currentLevelInfo.icon}</div>
              </div>

              {/* 进度条 */}
              {nextLevelInfo && (
                <div className="mt-3 sm:mt-4">
                  <div className="flex justify-between text-xs text-gray-500 mb-2">
                    <span>距离下一级</span>
                    <span>{stats.progress}%</span>
                  </div>
                  <div className="w-full bg-gray-200 rounded-full h-2">
                    <div
                      className="h-2 rounded-full transition-all duration-300"
                      style={{
                        width: `${stats.progress}%`,
                        backgroundColor: currentLevelInfo.color,
                      }}
                    />
                  </div>
                  <div className="text-xs text-gray-400 mt-2">
                    下一级: {nextLevelInfo.title} (需要 {nextLevelInfo.minScore} 分)
                  </div>
                </div>
              )}

              {/* 统计数据 */}
              <div className="mt-4 sm:mt-5 pt-4 sm:pt-5 border-t border-gray-100 grid grid-cols-3 gap-3 sm:gap-4">
                <div className="text-center">
                  <div className="text-xl sm:text-2xl font-bold text-green-600">{stats.explorationCount}</div>
                  <div className="text-xs text-gray-500 mt-1">探索次数</div>
                </div>
                <div className="text-center">
                  <div className="text-xl sm:text-2xl font-bold text-green-600">{stats.collectionCount}</div>
                  <div className="text-xs text-gray-500 mt-1">收藏次数</div>
                </div>
                <div className="text-center">
                  <div className="text-xl sm:text-2xl font-bold text-green-600">{stats.conversationCount}</div>
                  <div className="text-xs text-gray-500 mt-1">对话次数</div>
                </div>
              </div>

              {/* 总分 */}
              <div className="mt-3 sm:mt-4 text-center">
                <div className="text-xs sm:text-sm text-gray-500">总分</div>
                <div className="text-2xl sm:text-3xl font-bold mt-1" style={{ color: currentLevelInfo.color }}>
                  {stats.totalScore}
                </div>
              </div>
            </div>

            {/* 所有等级列表 */}
            <div className="lg:col-span-2 bg-white rounded-2xl shadow-lg p-4 sm:p-5 lg:p-6">
              <div className="text-base sm:text-lg font-bold mb-3 sm:mb-4 text-gray-800">所有等级</div>
              <div className="space-y-3 sm:space-y-4 max-h-[600px] overflow-y-auto">
                {allLevels.map((level) => {
                  const isCurrentLevel = level.level === stats.currentLevel;
                  const isUnlocked = stats.totalScore >= level.minScore;

                  return (
                    <div
                      key={level.level}
                      className={`flex items-center p-3 sm:p-4 rounded-xl border-2 transition-all ${
                        isCurrentLevel
                          ? 'border-green-500 bg-green-50'
                          : isUnlocked
                          ? 'border-gray-200 bg-gray-50'
                          : 'border-gray-100 bg-gray-50 opacity-50'
                      }`}
                    >
                      <div className="text-3xl sm:text-4xl mr-3 sm:mr-4 flex-shrink-0">{level.icon}</div>
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center gap-2 flex-wrap">
                          <div className={`text-sm sm:text-base font-bold ${isCurrentLevel ? 'text-green-600' : 'text-gray-600'}`}>
                            {level.title}
                          </div>
                          {isCurrentLevel && (
                            <span className="px-2 py-1 bg-green-500 text-white text-xs rounded-full flex-shrink-0">
                              当前
                            </span>
                          )}
                        </div>
                        <div className="text-xs text-gray-500 mt-1">{level.description}</div>
                        <div className="text-xs text-gray-400 mt-1">需要 {level.minScore} 分</div>
                      </div>
                      <div className="text-right flex-shrink-0 ml-2">
                        <div className="text-xs sm:text-sm font-bold text-gray-600">Lv.{level.level}</div>
                      </div>
                    </div>
                  );
                })}
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}

