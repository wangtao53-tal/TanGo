/**
 * 分享页面组件（家长端）
 * 基于设计稿，展示孩子分享的探索结果
 */

import { useParams } from 'react-router-dom';
import { useState, useEffect } from 'react';
import { getShare } from '../services/api';
import type { GetShareResponse } from '../types/api';
import { CollectionGrid } from '../components/collection/CollectionGrid';
import type { ExplorationRecord } from '../types/exploration';

export default function Share() {
  const { shareId } = useParams<{ shareId: string }>();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [shareData, setShareData] = useState<GetShareResponse | null>(null);
  const [records, setRecords] = useState<ExplorationRecord[]>([]);

  useEffect(() => {
    if (shareId) {
      loadShareData(shareId);
    } else {
      setError('分享链接无效');
      setLoading(false);
    }
  }, [shareId]);

  const loadShareData = async (id: string) => {
    try {
      const data = await getShare(id);
      setShareData(data);

      // 转换为 ExplorationRecord 格式
      const convertedRecords: ExplorationRecord[] = data.explorationRecords.map((r) => ({
        id: r.id,
        timestamp: r.timestamp,
        objectName: r.objectName,
        objectCategory: r.objectCategory,
        confidence: 0.95, // 分享数据中没有，使用默认值
        age: r.age,
        cards: r.cards.map((c) => ({
          id: `card-${c.type}-${r.id}`,
          explorationId: r.id,
          type: c.type,
          title: c.title,
          content: c.content as any,
        })),
        collected: true,
      }));

      setRecords(convertedRecords);
    } catch (err: any) {
      setError(err.message || '加载分享数据失败');
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-cloud-white flex items-center justify-center">
        <div className="text-text-sub">加载中...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-cloud-white flex items-center justify-center">
        <div className="text-center">
          <div className="text-6xl mb-4">😢</div>
          <p className="text-text-main text-lg font-display mb-2">加载失败</p>
          <p className="text-text-sub">{error}</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-cloud-white font-display">
      <main className="flex-1 px-4 py-8 md:px-10 lg:px-20">
        <div className="max-w-6xl mx-auto flex flex-col gap-8">
          {/* 页面头部 */}
          <header className="flex flex-col gap-4 bg-white/60 backdrop-blur-sm p-6 rounded-3xl border border-white shadow-sm">
            <h1 className="text-3xl md:text-4xl font-extrabold tracking-tight text-text-main font-display">
              孩子的探索成果
            </h1>
            {shareData && (
              <div className="flex items-center gap-4 text-text-sub text-sm">
                <span>创建时间: {new Date(shareData.createdAt).toLocaleString('zh-CN')}</span>
                <span>•</span>
                <span>过期时间: {new Date(shareData.expiresAt).toLocaleString('zh-CN')}</span>
              </div>
            )}
          </header>

          {/* 探索记录列表 */}
          {records.length > 0 ? (
            <CollectionGrid records={records} />
          ) : (
            <div className="flex flex-col items-center justify-center py-20">
              <div className="text-6xl mb-4">📚</div>
              <p className="text-text-sub text-lg font-display">暂无探索记录</p>
            </div>
          )}
        </div>
      </main>
    </div>
  );
}
