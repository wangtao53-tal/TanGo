/**
 * 识别结果页面组件
 * 基于 stitch_ui/recognition_result_page_1/ 设计稿
 * 展示三张知识卡片
 */

import { useNavigate, useLocation } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Header } from '../components/common/Header';
import { ConversationList } from '../components/conversation/ConversationList';
import { MessageInput } from '../components/conversation/MessageInput';
import { VoiceInput } from '../components/conversation/VoiceInput';
import { ImageInput } from '../components/conversation/ImageInput';
import type { KnowledgeCard } from '../types/exploration';
import type { ConversationMessage } from '../types/conversation';
import { useState, useEffect } from 'react';
import { cardStorage } from '../services/storage';
import { sendMessage, createStreamConnection, recognizeUserIntent } from '../services/conversation';
import { fileToBase64, extractBase64Data, compressImage } from '../utils/image';
import { uploadImage } from '../services/api';

interface LocationState {
  objectName: string;
  objectCategory: '自然类' | '生活类' | '人文类';
  confidence: number;
  cards: KnowledgeCard[];
  imageData?: string;
}

export default function Result() {
  const navigate = useNavigate();
  const location = useLocation();
  const { t } = useTranslation();
  const [cards, setCards] = useState<KnowledgeCard[]>([]);
  const [objectName, setObjectName] = useState<string>('Unknown');
  const [objectCategory, setObjectCategory] = useState<'自然类' | '生活类' | '人文类'>('自然类');
  const [messages, setMessages] = useState<ConversationMessage[]>([]);
  const [sessionId, setSessionId] = useState<string>('');
  const [isSending, setIsSending] = useState(false);
  const [sseConnection, setSseConnection] = useState<EventSource | null>(null);

  useEffect(() => {
    const state = location.state as LocationState;
    if (state && state.cards) {
      setCards(state.cards);
      setObjectName(state.objectName || 'Unknown');
      setObjectCategory(state.objectCategory || '自然类');
      
      // 生成会话ID
      const newSessionId = `session-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
      setSessionId(newSessionId);

      // 创建初始消息：将卡片作为对话消息
      const initialMessages: ConversationMessage[] = [
        {
          id: `msg-init-${Date.now()}`,
          type: 'text',
          content: `${t('result.identifiedAs')} ${state.objectName}!`,
          timestamp: new Date().toISOString(),
          sender: 'assistant',
          sessionId: newSessionId,
        },
        ...state.cards.map((card, index) => ({
          id: `msg-card-${index}`,
          type: 'card' as const,
          content: card,
          timestamp: new Date().toISOString(),
          sender: 'assistant' as const,
          sessionId: newSessionId,
        })),
      ];
      setMessages(initialMessages);
    } else {
      // 如果没有数据，返回首页
      navigate('/');
    }

    // 清理函数：关闭SSE连接
    return () => {
      if (sseConnection) {
        sseConnection.close();
      }
    };
  }, [location, navigate, t]);

  const handleCollect = async (cardId: string) => {
    const card = cards.find((c) => c.id === cardId) || 
                 messages.find((m) => m.type === 'card' && (m.content as KnowledgeCard).id === cardId)?.content as KnowledgeCard;
    if (!card) return;

    // 收藏或取消收藏
    const existingCard = await cardStorage.getAll().then(cards => cards.find(c => c.id === cardId));
    if (existingCard) {
      await cardStorage.delete(cardId);
    } else {
      const cardToSave = {
        ...card,
        collectedAt: new Date().toISOString(),
      };
      await cardStorage.save(cardToSave);
    }
  };

  const handleSendMessage = async (text: string) => {
    if (!text.trim() || isSending) return;

    setIsSending(true);
    try {
      // 识别意图
      const intent = await recognizeUserIntent(text, sessionId);
      
      // 发送消息
      const result = await sendMessage(text, 'text', sessionId);
      setMessages((prev) => [...prev, result.message]);

      // 如果是生成卡片意图，使用流式返回
      if (intent.intent === 'generate_cards') {
        const connection = createStreamConnection(sessionId, {
          onMessage: (message) => {
            setMessages((prev) => [...prev, message]);
          },
          onError: (error) => {
            console.error('流式返回错误:', error);
            setIsSending(false);
          },
          onClose: () => {
            setIsSending(false);
          },
        });
        setSseConnection(connection);
      } else {
        setIsSending(false);
      }
    } catch (error) {
      console.error('发送消息失败:', error);
      setIsSending(false);
    }
  };

  const handleVoiceResult = async (text: string) => {
    await handleSendMessage(text);
  };

  const handleImageSelect = async (file: File) => {
    if (isSending) return;

    setIsSending(true);
    try {
      // 1. 压缩图片
      const compressedBlob = await compressImage(file, 1920, 1920, 0.8);
      
      // 2. 转换为 base64
      const compressedFile = new File([compressedBlob], file.name, { type: 'image/jpeg' });
      const base64 = await fileToBase64(compressedFile);
      const imageData = extractBase64Data(base64);

      // 3. 上传图片到 GitHub（如果失败会自动降级到 base64）
      let imageUrl = imageData; // 默认使用 base64
      try {
        const uploadResult = await uploadImage({
          imageData: imageData,
          filename: file.name,
        });
        imageUrl = uploadResult.url;
        console.log('图片上传成功:', uploadResult.url, '方式:', uploadResult.uploadMethod);
      } catch (uploadError: any) {
        console.warn('图片上传失败，使用 base64:', uploadError);
        // 上传失败时使用 base64，继续流程
      }
      
      // 4. 发送消息（使用 URL 或 base64）
      const result = await sendMessage(imageUrl, 'image', sessionId);
      setMessages((prev) => [...prev, result.message]);
      setIsSending(false);
    } catch (error) {
      console.error('发送图片失败:', error);
      setIsSending(false);
    }
  };

  return (
    <div className="min-h-screen bg-cloud-white font-display flex flex-col">
      <Header />
      
      <main className="flex-1 flex flex-col items-center px-4 py-6 w-full max-w-7xl mx-auto overflow-hidden">
        {/* 对象信息展示区域 */}
        <div className="flex flex-col md:flex-row items-center justify-between gap-6 mb-8 w-full">
          <div className="flex flex-col items-start gap-2 relative">
            <span className="absolute -top-6 -left-2 rotate-[-5deg] bg-yellow-300 text-yellow-800 px-3 py-1 rounded-lg text-xs font-display shadow-sm border-2 border-yellow-100 transform z-10">
              You found a new friend!
            </span>
            <div className="flex items-center gap-3">
              <span className={`px-4 py-1.5 rounded-full ${
                objectCategory === '自然类' ? 'bg-science-green/20 text-science-green border-science-green/30' :
                objectCategory === '生活类' ? 'bg-sunny-orange/20 text-sunny-orange border-sunny-orange/30' :
                'bg-sky-blue/20 text-sky-blue border-sky-blue/30'
              } border-2 text-sm font-bold flex items-center gap-1.5 font-display shadow-sm`}>
                <span className="material-symbols-outlined text-lg">psychiatry</span>
                {objectCategory}
              </span>
            </div>
            <h1 className="text-4xl md:text-6xl font-display font-extrabold text-slate-800 leading-tight mt-2 drop-shadow-sm">
              It's a <span className="text-transparent bg-clip-text bg-gradient-to-r from-science-green to-green-500 relative inline-block">
                {objectName}
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
                {t('result.continueChat', `"Wow! A ${objectName}! Let's explore its secrets!"`)}
              </p>
            </div>
            <button className="absolute -top-2 -right-2 size-10 rounded-full bg-sky-blue text-white shadow-md border-4 border-white flex items-center justify-center hover:bg-sky-blue-dark transition-colors">
              <span className="material-symbols-outlined text-xl">volume_up</span>
            </button>
          </div>
        </div>

        {/* 对话消息列表 */}
        <div className="flex-1 w-full max-w-4xl mx-auto mb-24">
          <ConversationList messages={messages} onCollect={handleCollect} />
        </div>

        {/* 底部输入栏 */}
        <footer className="fixed bottom-0 left-0 w-full bg-white/90 backdrop-blur-xl border-t-2 border-slate-100 z-50">
          <div className="max-w-4xl mx-auto px-4 py-3">
            <div className="flex items-center gap-2">
              <VoiceInput onResult={handleVoiceResult} disabled={isSending} />
              <ImageInput onImageSelect={handleImageSelect} disabled={isSending} />
              <div className="flex-1">
                <MessageInput onSend={handleSendMessage} disabled={isSending} />
              </div>
            </div>
          </div>
        </footer>
      </main>
    </div>
  );
}
