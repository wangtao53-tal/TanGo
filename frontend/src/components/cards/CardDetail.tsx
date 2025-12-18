/**
 * 卡片详情组件
 * 展示卡片完整细节
 */

import { ScienceCard } from './ScienceCard';
import { PoetryCard } from './PoetryCard';
import { EnglishCard } from './EnglishCard';
import type { KnowledgeCard } from '../../types/exploration';

export interface CardDetailProps {
  card: KnowledgeCard;
  onCollect?: (cardId: string) => void;
  onClose?: () => void;
}

export function CardDetail({ card, onCollect, onClose }: CardDetailProps) {
  const renderCard = () => {
    switch (card.type) {
      case 'science':
        return <ScienceCard card={card} onCollect={onCollect} id={`card-detail-${card.id}`} />;
      case 'poetry':
        return <PoetryCard card={card} onCollect={onCollect} id={`card-detail-${card.id}`} />;
      case 'english':
        return <EnglishCard card={card} onCollect={onCollect} id={`card-detail-${card.id}`} />;
      default:
        return null;
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm p-4">
      <div className="relative w-full max-w-md max-h-[90vh] overflow-y-auto">
        {renderCard()}
        {onClose && (
          <button
            onClick={onClose}
            className="absolute top-4 right-4 size-10 rounded-full bg-white/90 hover:bg-white text-gray-600 shadow-lg flex items-center justify-center transition-all"
          >
            <span className="material-symbols-outlined">close</span>
          </button>
        )}
      </div>
    </div>
  );
}
