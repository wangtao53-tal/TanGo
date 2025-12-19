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
import type { ConversationMessage } from '../types/conversation';
import type { IdentificationContext } from '../types/api';
import { useState, useEffect, useRef } from 'react';
import { cardStorage } from '../services/storage';
import { sendMessage, createStreamConnection, recognizeUserIntent } from '../services/conversation';
import { fileToBase64, extractBase64Data, compressImage } from '../utils/image';
import { uploadImage, generateCards } from '../services/api';
import type { GenerateCardsRequest, GenerateCardsResponse } from '../types/api';
import type { KnowledgeCard } from '../types/exploration';

interface LocationState {
  objectName: string;
  objectCategory: '自然类' | '生活类' | '人文类';
  confidence: number;
  keywords?: string[];
  age?: number;
  imageData?: string;
}

export default function Result() {
  const navigate = useNavigate();
  const location = useLocation();
  const { t } = useTranslation();
  const [objectName, setObjectName] = useState<string>('Unknown');
  const [objectCategory, setObjectCategory] = useState<'自然类' | '生活类' | '人文类'>('自然类');
  const [messages, setMessages] = useState<ConversationMessage[]>([]);
  const [sessionId, setSessionId] = useState<string>('');
  const [isSending, setIsSending] = useState(false);
  const [sseConnection, setSseConnection] = useState<EventSource | null>(null);
  const [identificationContext, setIdentificationContext] = useState<IdentificationContext | null>(null);
  const [isGeneratingCards, setIsGeneratingCards] = useState(false);
  const hasGeneratedCardsRef = useRef(false); // 使用 ref 防止重复调用
  const generateTimeoutRef = useRef<NodeJS.Timeout | null>(null); // 存储定时器引用

  useEffect(() => {
    const state = location.state as LocationState;
    if (state && state.objectName) {
      // 如果已经生成过卡片，不再重复生成
      if (hasGeneratedCardsRef.current) {
        return;
      }

      setObjectName(state.objectName || 'Unknown');
      setObjectCategory(state.objectCategory || '自然类');
      
      // 生成会话ID
      const newSessionId = `session-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
      setSessionId(newSessionId);

      // 保存识别结果上下文
      const context: IdentificationContext = {
        objectName: state.objectName,
        objectCategory: state.objectCategory,
        confidence: state.confidence,
        keywords: state.keywords,
        age: state.age,
      };
      setIdentificationContext(context);

      // 创建初始系统消息，展示识别结果
      const initialMessages: ConversationMessage[] = [
        {
          id: `msg-init-${Date.now()}`,
          type: 'text',
          content: `${t('result.identifiedAs', '识别为')} ${state.objectName}！置信度：${(state.confidence * 100).toFixed(0)}%`,
          timestamp: new Date().toISOString(),
          sender: 'assistant',
          sessionId: newSessionId,
        },
      ];
      setMessages(initialMessages);

      // 标记为已生成，防止重复调用
      hasGeneratedCardsRef.current = true;

      // 自动生成卡片（延迟一下，确保初始消息已显示）
      // 清除之前的定时器（如果存在）
      if (generateTimeoutRef.current) {
        clearTimeout(generateTimeoutRef.current);
      }
      generateTimeoutRef.current = setTimeout(() => {
        generateCardsAutomatically(state, newSessionId);
        generateTimeoutRef.current = null;
      }, 500);
    } else {
      // 如果没有数据，返回首页
      navigate('/');
    }

    // 清理函数：关闭SSE连接和清除定时器
    return () => {
      if (sseConnection) {
        sseConnection.close();
      }
      if (generateTimeoutRef.current) {
        clearTimeout(generateTimeoutRef.current);
        generateTimeoutRef.current = null;
      }
      // 重置标志，以便下次进入页面时可以重新生成
      hasGeneratedCardsRef.current = false;
    };
  }, [location.state]); // 只依赖 location.state，避免不必要的重复执行

  // 自动生成卡片函数
  const generateCardsAutomatically = async (state: LocationState, sessionId: string) => {
    // 如果正在生成或已经生成过，直接返回
    if (isGeneratingCards) {
      return;
    }
    
    setIsGeneratingCards(true);
    const loadingMessageId = `msg-loading-${Date.now()}`;
    
    try {
      // 添加加载提示消息
      const loadingMessage: ConversationMessage = {
        id: loadingMessageId,
        type: 'text',
        content: '正在为您生成知识卡片...',
        timestamp: new Date().toISOString(),
        sender: 'assistant',
        sessionId,
      };
      setMessages((prev) => [...prev, loadingMessage]);

      // 调用生成卡片API
      const request: GenerateCardsRequest = {
        objectName: state.objectName,
        objectCategory: state.objectCategory,
        age: state.age || 8, // 默认8岁
        keywords: state.keywords,
      };

      const cardsResult: GenerateCardsResponse = await generateCards(request);

      // 移除加载消息
      setMessages((prev) => prev.filter(msg => msg.id !== loadingMessageId));

      // 转换为KnowledgeCard格式并添加到消息列表
      const timestamp = Date.now();
      const knowledgeCards: KnowledgeCard[] = cardsResult.cards.map((card, index) => ({
        id: `card-${card.type}-${timestamp}-${index}`,
        explorationId: `exp-${timestamp}`,
        type: card.type as 'science' | 'poetry' | 'english',
        title: card.title,
        content: card.content as any,
      }));

      // 将卡片作为消息添加到对话中
      const cardMessages: ConversationMessage[] = knowledgeCards.map((card) => ({
        id: `msg-card-${card.id}`,
        type: 'card' as const,
        content: card,
        timestamp: new Date().toISOString(),
        sender: 'assistant' as const,
        sessionId,
      }));

      setMessages((prev) => [...prev, ...cardMessages]);
    } catch (error: any) {
      console.error('自动生成卡片失败:', error);
      // 移除加载消息
      setMessages((prev) => prev.filter(msg => msg.id !== loadingMessageId));
      
      // 添加错误提示
      const errorMessage: ConversationMessage = {
        id: `msg-error-${Date.now()}`,
        type: 'text',
        content: `生成卡片失败：${error?.message || '未知错误'}。您可以稍后通过对话重新生成。`,
        timestamp: new Date().toISOString(),
        sender: 'assistant',
        sessionId,
      };
      setMessages((prev) => [...prev, errorMessage]);
    } finally {
      setIsGeneratingCards(false);
    }
  };

  const handleCollect = async (cardId: string) => {
    // 从消息中查找卡片
    const cardMessage = messages.find((m) => m.type === 'card' && (m.content as any)?.id === cardId);
    if (!cardMessage) return;

    const card = cardMessage.content as any;
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
      // 立即显示用户消息（乐观更新）
      const userMessage: ConversationMessage = {
        id: `msg-temp-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        type: 'text',
        content: text,
        timestamp: new Date().toISOString(),
        sender: 'user',
        sessionId,
      };
      setMessages((prev) => [...prev, userMessage]);

      // 识别意图
      const intent = await recognizeUserIntent(text, sessionId);
      
      // 发送消息，传递识别结果上下文（仅在首次发送时传递）
      const result = await sendMessage(
        text, 
        'text', 
        sessionId, 
        identificationContext || undefined,
        (msg) => {
          // 用户消息已通过乐观更新显示，这里可以更新ID
          const index = messages.findIndex(m => m.id === userMessage.id);
          if (index >= 0) {
            setMessages((prev) => {
              const updated = [...prev];
              updated[index] = { ...updated[index], id: msg.id };
              return updated;
            });
          }
        }
      );
      
      // 更新用户消息ID（如果后端返回了新的ID）
      if (result.userMessage) {
        setMessages((prev) => {
          const index = prev.findIndex(m => m.id === userMessage.id);
          if (index >= 0) {
            const updated = [...prev];
            updated[index] = { ...updated[index], id: result.userMessage.id };
            return updated;
          }
          return prev;
        });
      }
      
      // 添加助手消息
      setMessages((prev) => [...prev, result.message]);

      // 如果是生成卡片意图，使用流式返回
      if (intent.intent === 'generate_cards') {
        const connection = createStreamConnection(sessionId, {
          onMessage: (message) => {
            setMessages((prev) => [...prev, message]);
          },
          onError: (error) => {
            console.error('流式返回错误:', error);
            // 友好的错误提示
            const errorMessage: ConversationMessage = {
              id: `msg-error-${Date.now()}`,
              type: 'text',
              content: '抱歉，生成卡片时出现错误，请稍后重试。',
              timestamp: new Date().toISOString(),
              sender: 'assistant',
              sessionId,
            };
            setMessages((prev) => [...prev, errorMessage]);
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
    } catch (error: any) {
      console.error('发送消息失败:', error);
      // 友好的错误提示
      const errorMessage: ConversationMessage = {
        id: `msg-error-${Date.now()}`,
        type: 'text',
        content: `抱歉，发送消息失败：${error?.message || '未知错误'}。请检查网络连接后重试。`,
        timestamp: new Date().toISOString(),
        sender: 'assistant',
        sessionId,
      };
      setMessages((prev) => [...prev, errorMessage]);
      setIsSending(false);
    }
  };

  const handleVoiceResult = async (text: string) => {
    if (!text.trim() || isSending) return;

    setIsSending(true);
    try {
      // 立即显示语音识别结果作为用户消息（乐观更新）
      const userMessage: ConversationMessage = {
        id: `msg-temp-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        type: 'voice',
        content: text,
        timestamp: new Date().toISOString(),
        sender: 'user',
        sessionId,
      };
      setMessages((prev) => [...prev, userMessage]);

      // 发送消息
      const result = await sendMessage(
        text,
        'voice',
        sessionId,
        identificationContext || undefined,
        (msg) => {
          // 更新用户消息ID
          const index = messages.findIndex(m => m.id === userMessage.id);
          if (index >= 0) {
            setMessages((prev) => {
              const updated = [...prev];
              updated[index] = { ...updated[index], id: msg.id };
              return updated;
            });
          }
        }
      );

      // 更新用户消息ID
      if (result.userMessage) {
        setMessages((prev) => {
          const index = prev.findIndex(m => m.id === userMessage.id);
          if (index >= 0) {
            const updated = [...prev];
            updated[index] = { ...updated[index], id: result.userMessage.id };
            return updated;
          }
          return prev;
        });
      }

      // 添加助手消息
      setMessages((prev) => [...prev, result.message]);
      setIsSending(false);
    } catch (error: any) {
      console.error('发送语音消息失败:', error);
      // 友好的错误提示
      const errorMessage: ConversationMessage = {
        id: `msg-error-${Date.now()}`,
        type: 'text',
        content: `抱歉，发送语音消息失败：${error?.message || '未知错误'}。请重试。`,
        timestamp: new Date().toISOString(),
        sender: 'assistant',
        sessionId,
      };
      setMessages((prev) => [...prev, errorMessage]);
      setIsSending(false);
    }
  };

  const handleImageSelect = async (file: File) => {
    if (isSending) return;

    setIsSending(true);
    try {
      // 1. 立即显示图片作为用户消息（乐观更新）
      const tempImageUrl = URL.createObjectURL(file);
      const userMessage: ConversationMessage = {
        id: `msg-temp-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
        type: 'image',
        content: tempImageUrl, // 临时使用本地URL
        timestamp: new Date().toISOString(),
        sender: 'user',
        sessionId,
      };
      setMessages((prev) => [...prev, userMessage]);

      // 2. 压缩图片
      const compressedBlob = await compressImage(file, 1920, 1920, 0.8);
      
      // 3. 转换为 base64
      const compressedFile = new File([compressedBlob], file.name, { type: 'image/jpeg' });
      const base64 = await fileToBase64(compressedFile);
      const imageData = extractBase64Data(base64);

      // 4. 上传图片到 GitHub（如果失败会自动降级到 base64）
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
      
      // 5. 更新用户消息内容为最终URL
      setMessages((prev) => {
        const index = prev.findIndex(m => m.id === userMessage.id);
        if (index >= 0) {
          const updated = [...prev];
          updated[index] = { ...updated[index], content: imageUrl };
          return updated;
        }
        return prev;
      });
      
      // 6. 发送消息（使用 URL 或 base64），传递识别结果上下文
      const result = await sendMessage(
        imageUrl, 
        'image', 
        sessionId, 
        identificationContext || undefined,
        (msg) => {
          // 更新用户消息ID
          const index = messages.findIndex(m => m.id === userMessage.id);
          if (index >= 0) {
            setMessages((prev) => {
              const updated = [...prev];
              updated[index] = { ...updated[index], id: msg.id };
              return updated;
            });
          }
        }
      );

      // 更新用户消息ID
      if (result.userMessage) {
        setMessages((prev) => {
          const index = prev.findIndex(m => m.id === userMessage.id);
          if (index >= 0) {
            const updated = [...prev];
            updated[index] = { ...updated[index], id: result.userMessage.id };
            return updated;
          }
          return prev;
        });
      }

      // 添加助手消息
      setMessages((prev) => [...prev, result.message]);
      setIsSending(false);
    } catch (error: any) {
      console.error('发送图片失败:', error);
      // 友好的错误提示
      const errorMessage: ConversationMessage = {
        id: `msg-error-${Date.now()}`,
        type: 'text',
        content: `抱歉，发送图片失败：${error?.message || '未知错误'}。请检查图片格式和大小后重试。`,
        timestamp: new Date().toISOString(),
        sender: 'assistant',
        sessionId,
      };
      setMessages((prev) => [...prev, errorMessage]);
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
