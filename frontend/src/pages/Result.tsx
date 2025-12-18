/**
 * 识别结果页面组件
 * 基于 stitch_ui/recognition_result_page_1/ 设计稿
 * 展示三张知识卡片
 */

import { useNavigate } from 'react-router-dom';
import { Header } from '../components/common/Header';
import { CardCarousel } from '../components/cards/CardCarousel';
import type { KnowledgeCard } from '../types/exploration';
import { useState } from 'react';

// Mock数据 - 待后端API接入后替换
const mockCards: KnowledgeCard[] = [
  {
    id: 'card-1',
    explorationId: 'exp-1',
    type: 'science',
    title: '银杏',
    content: {
      name: '银杏',
      explanation: '银杏是非常古老的植物，已经在地球上生存了2亿多年。',
      facts: [
        '银杏是现存最古老的树种之一',
        '银杏的叶子在秋天会变成金黄色',
        '银杏的果实可以食用，但需要处理',
      ],
      funFact: '银杏被称为"活化石"，因为它在恐龙时代就已经存在了！',
    },
  },
  {
    id: 'card-2',
    explorationId: 'exp-1',
    type: 'poetry',
    title: '古人怎么看银杏',
    content: {
      poem: '满地翻黄银杏叶，忽惊天地告成功。',
      author: '李清照',
      explanation: '这句诗描写了秋天银杏叶变黄的美丽景象。',
      context: '看到银杏，我们可以联想到秋天的美丽，以及时间的流逝。',
    },
  },
  {
    id: 'card-3',
    explorationId: 'exp-1',
    type: 'english',
    title: '用英语说银杏',
    content: {
      words: ['ginkgo', 'tree', 'ancient'],
      expressions: [
        'This is a ginkgo tree.',
        'The ginkgo leaves are golden in autumn.',
      ],
    },
  },
];

export default function Result() {
  const navigate = useNavigate();
  const [cards] = useState<KnowledgeCard[]>(mockCards);
  const [, setCollectedCards] = useState<Set<string>>(new Set());

  const handleCollect = (cardId: string) => {
    setCollectedCards((prev) => {
      const newSet = new Set(prev);
      if (newSet.has(cardId)) {
        newSet.delete(cardId);
      } else {
        newSet.add(cardId);
      }
      return newSet;
    });
  };

  const handleCollectAll = () => {
    const allCardIds = cards.map((c) => c.id);
    setCollectedCards(new Set(allCardIds));
  };

  return (
    <div className="min-h-screen bg-cloud-white font-display">
      <Header />
      
      <main className="flex flex-col items-center px-4 py-6 w-full max-w-7xl mx-auto">
        {/* 对象信息展示区域 */}
        <div className="flex flex-col md:flex-row items-center justify-between gap-6 mb-8 w-full">
          <div className="flex flex-col items-start gap-2 relative">
            <span className="absolute -top-6 -left-2 rotate-[-5deg] bg-yellow-300 text-yellow-800 px-3 py-1 rounded-lg text-xs font-display shadow-sm border-2 border-yellow-100 transform z-10">
              You found a new friend!
            </span>
            <div className="flex items-center gap-3">
              <span className="px-4 py-1.5 rounded-full bg-science-green/20 text-science-green border-2 border-science-green/30 text-sm font-bold flex items-center gap-1.5 font-display shadow-sm">
                <span className="material-symbols-outlined text-lg">psychiatry</span>
                Natural World
              </span>
            </div>
            <h1 className="text-4xl md:text-6xl font-display font-extrabold text-slate-800 leading-tight mt-2 drop-shadow-sm">
              It's a <span className="text-transparent bg-clip-text bg-gradient-to-r from-science-green to-green-500 relative inline-block">
                Dandelion
                <svg
                  className="absolute w-full h-3 -bottom-1 left-0 text-science-green"
                  preserveAspectRatio="none"
                  viewBox="0 0 100 10"
                >
                  <path d="M0 5 Q 50 10 100 5" fill="none" stroke="currentColor" strokeWidth="3" />
                </svg>
              </span>!
            </h1>
          </div>

          {/* AI Companion 提示 */}
          <div className="flex items-center gap-4 bg-white p-4 pr-8 rounded-3xl max-w-md shadow-card border-2 border-slate-100 relative group cursor-pointer hover:scale-105 transition-all">
            <div className="size-14 rounded-full bg-science-green/20 flex items-center justify-center border-2 border-science-green shrink-0 text-science-green animate-bounce shadow-inner">
              <span className="material-symbols-outlined text-3xl">smart_toy</span>
            </div>
            <div>
              <p className="text-sm font-bold text-slate-500 uppercase mb-0.5 font-display">
                AI Companion says:
              </p>
              <p className="text-base font-bold text-slate-800 leading-tight font-display">
                "Wow! A Dandelion! Let's explore its secrets!"
              </p>
            </div>
            <button className="absolute -top-2 -right-2 size-10 rounded-full bg-sky-blue text-white shadow-md border-4 border-white flex items-center justify-center hover:bg-sky-blue-dark transition-colors">
              <span className="material-symbols-outlined text-xl">volume_up</span>
            </button>
          </div>
        </div>

        {/* 卡片轮播区域 */}
        <CardCarousel cards={cards} onCollect={handleCollect} onCollectAll={handleCollectAll} />

        {/* 底部固定栏 */}
        <footer className="fixed bottom-0 left-0 w-full bg-white/90 backdrop-blur-xl border-t-2 border-slate-100 p-4 z-50 rounded-t-[2rem] shadow-[0_-10px_40px_rgba(0,0,0,0.05)]">
          <div className="max-w-7xl mx-auto flex justify-between items-center gap-4">
            <button
              onClick={() => navigate('/')}
              className="flex items-center gap-2 px-6 py-3 rounded-full bg-slate-100 hover:bg-slate-200 text-slate-600 font-bold transition-all border-2 border-slate-200 group"
            >
              <span className="material-symbols-outlined group-hover:-translate-x-1 transition-transform">
                undo
              </span>
              <span className="hidden sm:inline">Back to Map</span>
            </button>

            {/* 进度指示器 */}
            <div className="flex-1 flex justify-center">
              <div className="flex gap-3 bg-slate-100 px-4 py-2 rounded-full">
                <div className="w-10 h-2.5 rounded-full bg-science-green" />
                <div className="w-2.5 h-2.5 rounded-full bg-slate-300" />
                <div className="w-2.5 h-2.5 rounded-full bg-slate-300" />
              </div>
            </div>

            {/* 收藏所有卡片按钮 */}
            <button
              onClick={handleCollectAll}
              className="flex items-center gap-3 px-8 py-4 rounded-full bg-gradient-to-r from-science-green to-primary text-white font-display font-extrabold text-lg transition-all shadow-lg shadow-science-green/40 transform hover:-translate-y-1 hover:scale-105 border-b-4 border-green-600 active:border-b-0 active:translate-y-0.5"
            >
              <span className="material-symbols-outlined text-2xl animate-bounce">style</span>
              COLLECT CARDS!
            </button>
          </div>
        </footer>
      </main>
    </div>
  );
}
