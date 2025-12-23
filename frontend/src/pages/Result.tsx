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
import { flushSync } from 'react-dom';
import { cardStorage, explorationStorage, conversationStorage } from '../services/storage';
import { createStreamConnectionUnified } from '../services/conversation';
import { closePostSSEConnection } from '../services/sse-post';
import type { UnifiedStreamConversationRequest } from '../types/api';
import { fileToBase64, extractBase64Data, compressImage } from '../utils/image';
import { uploadImage, generateCardsStream } from '../services/api';
import type { GenerateCardsRequest } from '../types/api';
import type { KnowledgeCard, ExplorationRecord } from '../types/exploration';
import { AudioPlaybackProvider } from '../components/cards/ScienceCard';
import { getUserAgeFromStorage } from '../utils/age';

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
  const [sseConnection, setSseConnection] = useState<AbortController | null>(null);
  const [identificationContext, setIdentificationContext] = useState<IdentificationContext | null>(null);
  const [isGeneratingCards, setIsGeneratingCards] = useState(false);
  const hasGeneratedCardsRef = useRef(false); // 使用 ref 防止重复调用
  const generateTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null); // 存储定时器引用
  // 用于流式对话的ref（必须在组件顶层）
  const accumulatedTextRef = useRef('');
  const markdownRef = useRef<boolean>(false);
  const isInitializedRef = useRef(false); // 防止重复初始化

  // 从localStorage获取当前会话ID
  const getCurrentSessionId = (): string | null => {
    try {
      return localStorage.getItem('currentSessionId');
    } catch {
      return null;
    }
  };

  // 保存当前会话ID到localStorage
  const saveCurrentSessionId = (sessionId: string): void => {
    try {
      localStorage.setItem('currentSessionId', sessionId);
    } catch (error: any) {
      // 处理localStorage错误
      if (error && error.name === 'QuotaExceededError') {
        console.error('存储空间不足，无法保存会话ID:', error);
        // 可以提示用户清理存储空间
      } else {
        console.error('保存会话ID失败:', error);
      }
      // 不中断用户操作，仅记录错误
    }
  };

  // 保存识别上下文到localStorage
  const saveIdentificationContext = (context: IdentificationContext): void => {
    try {
      localStorage.setItem('identificationContext', JSON.stringify(context));
    } catch (error: any) {
      // 处理localStorage错误
      if (error && error.name === 'QuotaExceededError') {
        console.error('存储空间不足，无法保存识别上下文:', error);
        // 可以提示用户清理存储空间
      } else {
        console.error('保存识别上下文失败:', error);
      }
      // 不中断用户操作，仅记录错误
    }
  };

  // 从localStorage恢复识别上下文
  const restoreIdentificationContext = (): IdentificationContext | null => {
    try {
      const stored = localStorage.getItem('identificationContext');
      if (stored) {
        const context = JSON.parse(stored) as IdentificationContext;
        
        // 数据验证：确保上下文格式正确
        if (context && typeof context === 'object') {
          // 验证必需字段
          if (context.objectName && context.objectCategory) {
            // 确保分类值有效
            if (!['自然类', '生活类', '人文类'].includes(context.objectCategory)) {
              console.warn('无效的分类值，使用默认值"自然类"');
              context.objectCategory = '自然类';
            }
            return context;
          } else {
            console.warn('识别上下文缺少必需字段');
            return null;
          }
        } else {
          console.warn('识别上下文格式不正确');
          return null;
        }
      }
    } catch (error) {
      console.error('恢复识别上下文失败:', error);
      // 如果恢复失败，返回null，使用默认值
    }
    return null;
  };

  // 恢复对话记录
  const restoreConversation = async (sessionId: string) => {
    try {
      console.log('开始恢复对话记录，sessionId:', sessionId);
      // 注意：hasGeneratedCardsRef.current应该在调用此函数之前就已经设置为true
      // 这里只是确保状态正确，防止在恢复过程中误调用生成卡片接口
      // 刷新页面时不应该生成卡片，只有在首次拍照（从Capture页面跳转）时才生成卡片
      if (!hasGeneratedCardsRef.current) {
        hasGeneratedCardsRef.current = true;
        console.log('确保hasGeneratedCardsRef.current为true，防止误生成卡片');
      }
      
      const savedMessages = await conversationStorage.getMessagesBySessionId(sessionId);
      console.log('从IndexedDB获取到的消息数量:', savedMessages.length);
      
      if (savedMessages.length > 0) {
        // 确保所有流式消息标记为已完成（清除isStreaming和streamingText字段）
        const cleanedMessages = savedMessages.map(msg => ({
          ...msg,
          isStreaming: false,
          streamingText: undefined,
        }));
        
        // 按时间顺序排序（确保顺序正确）
        cleanedMessages.sort((a, b) => 
          new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime()
        );
        
        // 数据验证：确保消息格式正确
        const validMessages = cleanedMessages.filter(msg => 
          msg.id && 
          msg.type && 
          msg.sender && 
          msg.timestamp && 
          msg.sessionId === sessionId
        );
        
        if (validMessages.length !== cleanedMessages.length) {
          console.warn(`恢复的消息中有 ${cleanedMessages.length - validMessages.length} 条格式不正确，已过滤`);
        }
        
        setMessages(validMessages);
        console.log(`已恢复 ${validMessages.length} 条对话记录`);
        
        // 检查是否已有卡片消息（用于日志记录）
        // 重要：判断逻辑是"是否需要生成卡片"（即是否已经有卡片了），而不是"是否有需要识别就重新生成卡片"
        const hasCards = validMessages.some(msg => msg.type === 'card' && msg.sender === 'assistant');
        const cardMessages = validMessages.filter(msg => msg.type === 'card' && msg.sender === 'assistant');
        console.log(`检测到 ${cardMessages.length} 条卡片消息`);
        
        if (hasCards) {
          // hasGeneratedCardsRef.current已经在函数开始时设置为true，这里只是记录日志
          console.log('检测到已有卡片消息，不会重新生成卡片');
        } else {
          // 即使没有卡片消息，也不会生成卡片（hasGeneratedCardsRef.current已经在函数开始时设置为true）
          console.log('没有检测到卡片消息，不会生成卡片（刷新页面时不应该生成卡片）');
        }
        
        // 再次确认hasGeneratedCardsRef.current为true，防止任何误调用
        hasGeneratedCardsRef.current = true;
        console.log('恢复完成，hasGeneratedCardsRef.current =', hasGeneratedCardsRef.current);
      } else {
        console.log('没有找到对话记录');
        // hasGeneratedCardsRef.current已经在函数开始时设置为true，这里只是记录日志
      }
    } catch (error) {
      console.error('恢复对话记录失败:', error);
      // 显示友好提示（可选，可以通过toast或状态显示）
      // 这里仅记录错误，不中断用户操作
      // hasGeneratedCardsRef.current已经在函数开始时设置为true，确保不会误生成卡片
    }
  };

  // 初始化：处理新会话或恢复会话
  useEffect(() => {
    // 防止重复初始化
    if (isInitializedRef.current) {
      console.log('已初始化，跳过重复初始化');
      return;
    }
    isInitializedRef.current = true;
    console.log('开始初始化，location.state:', location.state);

    const state = location.state as LocationState;
    
    // ========== 核心逻辑：基于会话ID来源决定是否生成卡片 ==========
    // 规则：新生成会话ID → 生成卡片；从本地读取会话ID → 不生成卡片
    // ============================================================
    
    // 检测是否从Capture页面跳转过来
    // 使用sessionStorage标记，因为刷新页面时sessionStorage会清空，可以可靠地区分跳转和刷新
    const fromCapturePage = sessionStorage.getItem('fromCapturePage') === 'true';
    const hasState = state && state.objectName && typeof state.objectName === 'string' && state.objectName.trim() !== '';
    
    // 关键判断：只有同时满足"从Capture页面跳转"和"state有值"时，才创建新会话
    // 刷新页面时，即使state有值（浏览器可能保留），但sessionStorage会清空，所以不会创建新会话
    if (fromCapturePage && hasState) {
      // 场景1：从Capture页面进入 → 新生成会话ID → 生成卡片
      console.log('从Capture页面进入（sessionStorage标记），创建新会话，objectName:', state.objectName);
      
      // 重要：在创建新会话之前清除标记，避免useEffect再次执行时误判
      sessionStorage.removeItem('fromCapturePage');
      
      // 清空当前消息列表（清空旧会话）
      setMessages([]);

      setObjectName(state.objectName || 'Unknown');
      setObjectCategory(state.objectCategory || '自然类');
      
      // 生成新会话ID（关键：新生成的会话ID需要生成卡片）
      const newSessionId = `session-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
      setSessionId(newSessionId);
      saveCurrentSessionId(newSessionId);
      console.log('新生成会话ID:', newSessionId, '→ 将生成卡片');
      
      // 重要：新生成会话ID时，允许生成卡片
      hasGeneratedCardsRef.current = false;

      // 保存识别结果上下文
      const context: IdentificationContext = {
        objectName: state.objectName,
        objectCategory: state.objectCategory,
        confidence: state.confidence,
        keywords: state.keywords,
        age: state.age,
      };
      setIdentificationContext(context);
      saveIdentificationContext(context);

      // 创建初始消息列表
      const initialMessages: ConversationMessage[] = [];
      
      // 如果有图片数据，先添加用户图片消息
      if (state.imageData) {
        // 确保图片数据格式正确（支持 data URL 或纯 base64）
        let imageContent = state.imageData;
        if (!imageContent.startsWith('data:') && !imageContent.startsWith('http')) {
          // 如果是纯 base64，添加 data URL 前缀
          imageContent = `data:image/jpeg;base64,${imageContent}`;
        }
        
        const userImageMessage: ConversationMessage = {
          id: `msg-user-image-${Date.now()}`,
          type: 'image',
          content: imageContent,
          timestamp: new Date().toISOString(),
          sender: 'user',
          sessionId: newSessionId,
        };
        initialMessages.push(userImageMessage);
        // 保存到IndexedDB
        conversationStorage.saveMessage(userImageMessage).catch(err => 
          console.error('保存图片消息失败:', err)
        );
      }

      // 添加识别结果消息
      const initMessage: ConversationMessage = {
        id: `msg-init-${Date.now()}`,
        type: 'text',
        content: `${t('result.identifiedAs')} ${state.objectName}！置信度：${(state.confidence * 100).toFixed(0)}%`,
        timestamp: new Date().toISOString(),
        sender: 'assistant',
        sessionId: newSessionId,
      };
      initialMessages.push(initMessage);
      // 保存到IndexedDB
      conversationStorage.saveMessage(initMessage).catch(err => 
        console.error('保存初始消息失败:', err)
      );
      
      setMessages(initialMessages);

      // 自动生成卡片（新生成会话ID时，必须生成卡片）
      // 清除之前的定时器（如果存在）
      if (generateTimeoutRef.current) {
        clearTimeout(generateTimeoutRef.current);
      }
      console.log('新生成会话ID，准备生成卡片，sessionId:', newSessionId);
      generateTimeoutRef.current = setTimeout(() => {
        console.log('定时器触发，调用generateCardsAutomatically，sessionId:', newSessionId);
        // 检查定时器是否仍然有效（可能被清理函数清除了）
        if (generateTimeoutRef.current) {
          generateCardsAutomatically(state, newSessionId);
          generateTimeoutRef.current = null;
        } else {
          console.warn('定时器已被清除，跳过卡片生成');
        }
      }, 500);
    } else {
      // 场景2：刷新对话页面 → 从本地读取会话ID → 不生成卡片
      console.log('刷新对话页面，尝试恢复本地会话');
      
      // 重要：从本地读取会话ID时，不允许生成卡片
      hasGeneratedCardsRef.current = true;
      
      const currentSessionId = getCurrentSessionId();
      if (currentSessionId) {
        console.log('从本地读取会话ID:', currentSessionId, '→ 不生成卡片');
        setSessionId(currentSessionId);
        
        // 恢复识别上下文
        const restoredContext = restoreIdentificationContext();
        if (restoredContext) {
          setIdentificationContext(restoredContext);
          setObjectName(restoredContext.objectName || 'Unknown');
          setObjectCategory(restoredContext.objectCategory || '自然类');
          console.log('恢复识别上下文成功:', restoredContext);
        } else {
          // 如果上下文恢复失败，尝试使用默认值
          console.warn('识别上下文恢复失败，使用默认值');
          setObjectName('Unknown');
          setObjectCategory('自然类');
        }
        
        // 恢复对话记录（异步，等待恢复完成）
        // 从本地读取的会话ID，不生成卡片
        restoreConversation(currentSessionId)
          .then(() => {
            console.log('对话记录恢复完成，不会重新生成卡片（从本地读取的会话ID）');
            // 再次确认hasGeneratedCardsRef.current为true，防止任何误调用
            hasGeneratedCardsRef.current = true;
          })
          .catch(err => {
            console.error('恢复对话记录失败:', err);
            // 恢复失败时，可以选择显示友好提示或引导用户重新开始
            // 这里仅记录错误，不中断用户操作
            // 即使恢复失败，也确保不会生成卡片（因为是从本地读取的会话ID）
            hasGeneratedCardsRef.current = true;
          });
      } else {
        // 如果没有会话记录，返回首页
        console.log('没有找到会话记录，返回首页');
        // 确保不会生成卡片
        hasGeneratedCardsRef.current = true;
        navigate('/');
      }
    }

    // 清理函数：关闭SSE连接和清除定时器
    // 注意：这个清理函数会在useEffect重新执行之前执行（当依赖项变化时）
    // 但是，如果isInitializedRef.current已经是true，useEffect不会重新执行
    // 所以清理函数不应该清除正在进行的操作（如定时器）
    return () => {
      // 只在组件卸载时清理SSE连接
      // 注意：不要在依赖项变化时清除定时器，因为可能会中断正在进行的操作
      // 定时器会在组件卸载时自动清理
      if (sseConnection) {
        closePostSSEConnection(sseConnection);
      }
      // 注意：不要在这里清除generateTimeoutRef，因为可能会中断卡片生成
      // 定时器会在组件卸载时自动清理，或者在下次设置新定时器时清除旧的
      // 注意：不要重置 hasGeneratedCardsRef，因为恢复会话时需要保持这个状态
      // 只有在创建新会话时才会重置（在创建新会话的逻辑中处理）
      // hasGeneratedCardsRef.current = false;
      // 重要：不要重置isInitializedRef，避免useEffect重复执行
      // 只有在组件卸载时才应该重置，但这里不需要，因为组件卸载时会自动清理
      // isInitializedRef.current = false;
    };
  }, [location, navigate, t]); // 依赖 location, navigate, t，确保location变化时触发恢复

  // 自动生成卡片函数
  // 规则：只有新生成会话ID时才会调用此函数
  // 从本地读取会话ID时，hasGeneratedCardsRef.current会被设置为true，不会调用此函数
  const generateCardsAutomatically = async (state: LocationState, sessionId: string) => {
    console.log('generateCardsAutomatically被调用，sessionId:', sessionId, 'hasGeneratedCardsRef.current:', hasGeneratedCardsRef.current, 'isGeneratingCards:', isGeneratingCards);
    
    // 如果正在生成或已经生成过，直接返回
    // 这个检查确保从本地读取会话ID时不会误调用生成卡片接口
    if (isGeneratingCards || hasGeneratedCardsRef.current) {
      console.log('跳过卡片生成：正在生成或已生成过，hasGeneratedCardsRef.current =', hasGeneratedCardsRef.current, 'isGeneratingCards =', isGeneratingCards);
      return;
    }
    
    // 在开始生成前，先标记为已生成，防止重复调用
    hasGeneratedCardsRef.current = true;
    setIsGeneratingCards(true);
    console.log('开始生成卡片，设置hasGeneratedCardsRef.current = true, isGeneratingCards = true');
    const loadingMessageId = `msg-loading-${Date.now()}`;
    
    try {
      // 添加加载提示消息
      const loadingMessage: ConversationMessage = {
        id: loadingMessageId,
        type: 'text',
        content: t('conversation.generatingCards'),
        timestamp: new Date().toISOString(),
        sender: 'assistant',
        sessionId,
      };
      setMessages((prev) => [...prev, loadingMessage]);

      // 获取用户年龄（优先使用state中的年龄，否则从存储中获取）
      const userAge = state.age || getUserAgeFromStorage();
      
      // 调用流式生成卡片API
      const request: GenerateCardsRequest = {
        objectName: state.objectName,
        objectCategory: state.objectCategory,
        age: userAge,
        keywords: state.keywords,
      };

      // 使用流式API，每生成完一张卡片立即显示
      const timestamp = Date.now();
      const receivedCards: KnowledgeCard[] = [];

      // 2秒后显示进度提示
      const progressTimeout = setTimeout(() => {
        if (receivedCards.length === 0) {
          setMessages((prev) => {
            const index = prev.findIndex((m) => m.id === loadingMessageId);
            if (index >= 0) {
              const updated = [...prev];
              updated[index] = {
                ...updated[index],
                content: t('conversation.generatingCardsWait'),
              };
              return updated;
            }
            return prev;
          });
        }
      }, 2000);

      // 用于存储流式文本消息的ID和累积文本
      let streamMessageId: string | null = null;
      let streamFullText = '';

      const abortController = generateCardsStream(request, {
        onMessage: (_char: string, fullText: string) => {
          // 处理流式文本消息（逐字符返回）
          streamFullText = fullText;
          
          // 如果是第一次收到流式消息，将加载消息转换为流式消息
          if (!streamMessageId) {
            streamMessageId = loadingMessageId;
            flushSync(() => {
              setMessages((prev) => {
                const index = prev.findIndex((m) => m.id === loadingMessageId);
                if (index >= 0) {
                  const updated = [...prev];
                  updated[index] = {
                    ...updated[index],
                    isStreaming: true,
                    streamingText: fullText,
                    content: fullText,
                  };
                  return updated;
                }
                return prev;
              });
            });
          } else {
            // 实时更新流式文本
            flushSync(() => {
              setMessages((prev) => {
                const index = prev.findIndex((m) => m.id === streamMessageId);
                if (index >= 0) {
                  const updated = [...prev];
                  updated[index] = {
                    ...updated[index],
                    streamingText: fullText,
                    content: fullText,
                    isStreaming: true,
                  };
                  return updated;
                }
                return prev;
              });
            });
          }
        },
        onCard: async (card, index) => {
          // 验证卡片数据格式
          if (!card || !card.type || !card.title || !card.content) {
            console.error('卡片数据格式不正确:', card);
            return;
          }

          // 验证卡片类型
          if (!['science', 'poetry', 'english'].includes(card.type)) {
            console.error('未知的卡片类型:', card.type);
            return;
          }

          // 立即添加卡片到消息列表
          const knowledgeCard: KnowledgeCard = {
            id: `card-${card.type}-${timestamp}-${index}`,
            explorationId: `exp-${timestamp}`,
            type: card.type as 'science' | 'poetry' | 'english',
            title: card.title,
            content: card.content as any,
          };
          receivedCards.push(knowledgeCard);

          const cardMessage: ConversationMessage = {
            id: `msg-card-${knowledgeCard.id}`,
            type: 'card' as const,
            content: knowledgeCard,
            timestamp: new Date().toISOString(),
            sender: 'assistant' as const,
            sessionId,
          };

          // 保存卡片消息到IndexedDB
          await conversationStorage.saveMessage(cardMessage).catch(err => 
            console.error('保存卡片消息失败:', err)
          );

          // 使用flushSync立即更新UI
          flushSync(() => {
            setMessages((prev) => {
              // 移除加载消息（如果还在）
              const filtered = prev.filter(msg => msg.id !== loadingMessageId);
              return [...filtered, cardMessage];
            });
          });
        },
        onError: async (error) => {
          clearTimeout(progressTimeout);
          console.error('流式生成卡片失败:', error);
          
          // 完成流式输出，将流式消息标记为完成或移除
          flushSync(() => {
            setMessages((prev) => {
              const index = prev.findIndex((m) => m.id === loadingMessageId || m.id === streamMessageId);
              if (index >= 0) {
                const updated = [...prev];
                // 如果有流式文本，保留消息并标记为完成；否则移除
                if (streamFullText) {
                  updated[index] = {
                    ...updated[index],
                    isStreaming: false,
                    streamingText: undefined,
                    content: streamFullText,
                  };
                  return updated;
                } else {
                  // 移除加载消息
                  return updated.filter(msg => msg.id !== loadingMessageId && msg.id !== streamMessageId);
                }
              }
              return prev;
            });
          });
          
          const errorMessage: ConversationMessage = {
            id: `msg-error-${Date.now()}`,
            type: 'text',
            content: t('conversation.generateCardsError', { error: error?.message || t('conversation.unknownError') }),
            timestamp: new Date().toISOString(),
            sender: 'assistant',
            sessionId,
          };
          setMessages((prev) => [...prev, errorMessage]);
          // 保存错误消息到IndexedDB
          await conversationStorage.saveMessage(errorMessage).catch(err => 
            console.error('保存错误消息失败:', err)
          );
        },
        onComplete: async () => {
          clearTimeout(progressTimeout);
          // 完成流式输出，将流式消息标记为完成
          flushSync(() => {
            setMessages((prev) => {
              const index = prev.findIndex((m) => m.id === loadingMessageId || m.id === streamMessageId);
              if (index >= 0) {
                const updated = [...prev];
                // 如果有流式文本，保留消息并标记为完成；否则移除
                if (streamFullText) {
                  const completedMessage: ConversationMessage = {
                    ...updated[index],
                    isStreaming: false,
                    streamingText: undefined,
                    content: streamFullText,
                  };
                  updated[index] = completedMessage;
                  // 保存完成的流式消息到IndexedDB
                  conversationStorage.saveMessage(completedMessage).catch(err => 
                    console.error('保存流式消息失败:', err)
                  );
                  return updated;
                } else {
                  // 移除加载消息
                  return updated.filter(msg => msg.id !== loadingMessageId && msg.id !== streamMessageId);
                }
              }
              return prev;
            });
          });
          
          // 确保所有卡片都已添加
          if (receivedCards.length < 3) {
            console.warn(`只收到${receivedCards.length}张卡片，预期3张`);
          }

          // 保存探索记录到IndexedDB
          try {
            // 验证objectCategory字段
            const category = state.objectCategory || '自然类';
            if (!['自然类', '生活类', '人文类'].includes(category)) {
              console.warn('无效的分类值，使用默认值"自然类":', category);
            }

            const explorationRecord: ExplorationRecord = {
              id: `exp-${timestamp}`,
              timestamp: new Date().toISOString(),
              objectName: state.objectName,
              objectCategory: category as '自然类' | '生活类' | '人文类',
              confidence: state.confidence || 0.95,
              age: userAge,
              imageData: state.imageData,
              cards: receivedCards,
              collected: false,
            };

            await explorationStorage.save(explorationRecord);
            console.log('探索记录已保存:', explorationRecord.id, '分类:', explorationRecord.objectCategory);
          } catch (error) {
            console.error('保存探索记录失败:', error);
          }
        },
      });

      // 保存abortController以便可以取消
      setSseConnection(abortController);
    } catch (error: any) {
      console.error('自动生成卡片失败:', error);
      // 移除加载消息
      setMessages((prev) => prev.filter(msg => msg.id !== loadingMessageId));
      
      // 添加错误提示
      const errorMessage: ConversationMessage = {
        id: `msg-error-${Date.now()}`,
        type: 'text',
        content: t('conversation.generateCardsError', { error: error?.message || t('conversation.unknownError') }),
        timestamp: new Date().toISOString(),
        sender: 'assistant',
        sessionId,
      };
      setMessages((prev) => [...prev, errorMessage]);
      // 保存错误消息到IndexedDB
      await conversationStorage.saveMessage(errorMessage).catch(err => 
        console.error('保存错误消息失败:', err)
      );
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
      // 保存用户消息到IndexedDB
      await conversationStorage.saveMessage(userMessage).catch(err => 
        console.error('保存用户消息失败:', err)
      );

      // 直接使用POST + SSE流式接口，从请求体传递参数，避免中文在URL中的编码问题
      // 创建助手消息占位符
      const assistantMessageId = `msg-assistant-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
      const assistantMessage: ConversationMessage = {
        id: assistantMessageId,
        type: 'text',
        content: '',
        timestamp: new Date().toISOString(),
        sender: 'assistant',
        sessionId,
        isStreaming: true,
        streamingText: '',
      };
      setMessages((prev) => [...prev, assistantMessage]);

      // 获取用户年龄（优先使用识别上下文中的年龄，否则从存储中获取）
      const userAge = identificationContext?.age || getUserAgeFromStorage();
      
      // 构建统一流式请求参数
      const streamRequest: UnifiedStreamConversationRequest = {
        messageType: 'text',
        message: text,
        sessionId,
        userAge,
        maxContextRounds: 20,
      };

      // 首次发送时传递识别结果上下文
      if (identificationContext && !hasGeneratedCardsRef.current) {
        streamRequest.identificationContext = identificationContext;
      }

      // 重置累积文本和markdown状态
      accumulatedTextRef.current = '';
      markdownRef.current = false;

      // 使用统一流式接口
      const abortController = createStreamConnectionUnified(streamRequest, {
        onMessage: (message: ConversationMessage) => {
          // 更新助手消息
          flushSync(() => {
            setMessages((prev) => {
              const index = prev.findIndex((m) => m.id === assistantMessageId);
              if (index >= 0) {
                const updated = [...prev];
                updated[index] = {
                  ...updated[index],
                  content: message.content || '',
                  streamingText: message.streamingText || '',
                  markdown: message.markdown,
                  isStreaming: message.isStreaming,
                };
                return updated;
              }
              return prev;
            });
          });
        },
        onError: (error: Error) => {
          console.error('流式返回错误:', error);
          flushSync(() => {
            setMessages((prev) => {
              const index = prev.findIndex((m) => m.id === assistantMessageId);
              if (index >= 0) {
                const updated = [...prev];
                updated[index] = {
                  ...updated[index],
                  isStreaming: false,
                  content: accumulatedTextRef.current || t('conversation.generateAnswerError', { error: error.message }),
                  streamingText: undefined,
                };
                return updated;
              }
              return prev;
            });
          });
          setIsSending(false);
        },
        onClose: async () => {
          // 流式输出完成
          let completedMessage: ConversationMessage | null = null;
          flushSync(() => {
            setMessages((prev) => {
              const index = prev.findIndex((m) => m.id === assistantMessageId);
              if (index >= 0) {
                const updated = [...prev];
                // 创建完成的消息对象，清除isStreaming和streamingText字段
                completedMessage = {
                  ...updated[index],
                  isStreaming: false,
                  streamingText: undefined,
                } as ConversationMessage;
                updated[index] = completedMessage;
                return updated;
              }
              return prev;
            });
          });
          
          // 保存完成的消息到IndexedDB（清除isStreaming和streamingText字段）
          if (completedMessage && completedMessage !== null && 'content' in completedMessage) {
            try {
              await conversationStorage.saveMessage(completedMessage);
            } catch (err) {
              console.error('保存流式消息失败:', err);
              // 不中断用户操作，仅记录错误
            }
          }
          
          setIsSending(false);
        },
      });

      // 保存abortController以便可以取消
      setSseConnection(abortController);
    } catch (error: any) {
      console.error('发送消息失败:', error);
      // 友好的错误提示
      const errorMessage: ConversationMessage = {
        id: `msg-error-${Date.now()}`,
        type: 'text',
        content: t('conversation.sendMessageError', { error: error?.message || t('conversation.unknownError') }),
        timestamp: new Date().toISOString(),
        sender: 'assistant',
        sessionId,
      };
      setMessages((prev) => [...prev, errorMessage]);
      // 保存错误消息到IndexedDB
      await conversationStorage.saveMessage(errorMessage).catch(err => 
        console.error('保存错误消息失败:', err)
      );
      setIsSending(false);
    }
  };

  const handleVoiceResult = async (text: string, audioBase64?: string) => {
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
      // 保存用户消息到IndexedDB
      await conversationStorage.saveMessage(userMessage).catch(err => 
        console.error('保存语音消息失败:', err)
      );

      // 创建助手消息占位符
      const assistantMessageId = `msg-assistant-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
      const assistantMessage: ConversationMessage = {
        id: assistantMessageId,
        type: 'text',
        content: '',
        timestamp: new Date().toISOString(),
        sender: 'assistant',
        sessionId,
        isStreaming: true,
        streamingText: '',
      };
      setMessages((prev) => [...prev, assistantMessage]);

      // 获取用户年龄（优先使用识别上下文中的年龄，否则从存储中获取）
      const userAge = identificationContext?.age || getUserAgeFromStorage();
      
      // 构建统一流式请求参数
      // 注意：如果前端已经识别了语音，可以发送文本；如果需要后端识别，需要发送音频数据
      const streamRequest: UnifiedStreamConversationRequest = {
        messageType: audioBase64 ? 'voice' : 'text', // 如果有音频数据，使用voice类型；否则使用text类型
        audio: audioBase64 || '', // 如果有音频数据，使用音频；否则为空
        message: text, // 语音识别后的文本
        sessionId,
        userAge,
        maxContextRounds: 20,
      };

      // 首次发送时传递识别结果上下文
      if (identificationContext && !hasGeneratedCardsRef.current) {
        streamRequest.identificationContext = identificationContext;
      }

      // 重置累积文本和markdown状态
      accumulatedTextRef.current = '';
      markdownRef.current = false;

      // 使用统一流式接口
      const abortController = createStreamConnectionUnified(streamRequest, {
        onMessage: (message: ConversationMessage) => {
          // 更新助手消息
          // 使用 requestAnimationFrame 而不是 flushSync，避免在渲染期间触发副作用导致无限循环
          requestAnimationFrame(() => {
            setMessages((prev) => {
              const index = prev.findIndex((m) => m.id === assistantMessageId);
              if (index >= 0) {
                const currentMessage = prev[index];
                const newContent = message.content || '';
                const newStreamingText = message.streamingText || '';
                
                // 只有在内容真正变化时才更新，避免不必要的重新渲染
                if (
                  currentMessage.content !== newContent ||
                  currentMessage.streamingText !== newStreamingText ||
                  currentMessage.markdown !== message.markdown ||
                  currentMessage.isStreaming !== message.isStreaming
                ) {
                  const updated = [...prev];
                  updated[index] = {
                    ...currentMessage,
                    content: newContent,
                    streamingText: newStreamingText,
                    markdown: message.markdown,
                    isStreaming: message.isStreaming,
                  };
                  return updated;
                }
              }
              return prev;
            });
          });
        },
        onError: (error: Error) => {
          console.error('流式返回错误:', error);
          // 使用 setTimeout 延迟状态更新，避免在错误处理中触发副作用
          setTimeout(() => {
            setMessages((prev) => {
              const index = prev.findIndex((m) => m.id === assistantMessageId);
              if (index >= 0) {
                const currentMessage = prev[index];
                const errorContent = accumulatedTextRef.current || t('conversation.generateAnswerError', { error: error.message });
                
                // 只有在内容真正变化时才更新
                if (
                  currentMessage.content !== errorContent ||
                  currentMessage.isStreaming !== false ||
                  currentMessage.streamingText !== undefined
                ) {
                  const updated = [...prev];
                  updated[index] = {
                    ...currentMessage,
                    isStreaming: false,
                    content: errorContent,
                    streamingText: undefined,
                  };
                  return updated;
                }
              }
              return prev;
            });
            setIsSending(false);
          }, 0);
        },
        onClose: async () => {
          // 流式输出完成
          // 使用 setTimeout 延迟状态更新，避免在关闭回调中触发副作用
          setTimeout(async () => {
            let completedMessage: ConversationMessage | null = null;
            setMessages((prev) => {
              const index = prev.findIndex((m) => m.id === assistantMessageId);
              if (index >= 0) {
                const currentMessage = prev[index];
                // 创建完成的消息对象，清除isStreaming和streamingText字段
                completedMessage = {
                  ...currentMessage,
                  isStreaming: false,
                  streamingText: undefined,
                } as ConversationMessage;
                
                // 只有在状态真正变化时才更新
                if (
                  currentMessage.isStreaming !== false ||
                  currentMessage.streamingText !== undefined
                ) {
                  const updated = [...prev];
                  updated[index] = completedMessage;
                  return updated;
                }
              }
              return prev;
            });
            
            // 在状态更新后保存消息
            if (completedMessage && completedMessage !== null && 'content' in completedMessage) {
              conversationStorage.saveMessage(completedMessage).catch(err => {
                console.error('保存流式语音消息失败:', err);
              });
            }
            
            setIsSending(false);
          }, 0);
        },
      });

      // 保存abortController以便可以取消
      setSseConnection(abortController);
    } catch (error: any) {
      console.error('发送语音消息失败:', error);
      // 友好的错误提示
      const errorMessage: ConversationMessage = {
        id: `msg-error-${Date.now()}`,
        type: 'text',
        content: t('conversation.sendVoiceError', { error: error?.message || t('conversation.unknownError') }),
        timestamp: new Date().toISOString(),
        sender: 'assistant',
        sessionId,
      };
      setMessages((prev) => [...prev, errorMessage]);
      // 保存错误消息到IndexedDB
      await conversationStorage.saveMessage(errorMessage).catch(err => 
        console.error('保存错误消息失败:', err)
      );
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
      // 注意：图片消息会在更新为最终URL后再次保存

      // 2. 压缩图片
      const compressedBlob = await compressImage(file, 1920, 1920, 0.8);
      
      // 3. 创建压缩后的文件对象
      const compressedFile = new File([compressedBlob], file.name, { type: 'image/jpeg' });
      
      // 4. 上传图片到 GitHub（使用FormData方式，更高效）
      // 如果上传失败，降级到base64
      let imageUrl: string = '';
      try {
        const uploadResult = await uploadImage(compressedFile, file.name);
        imageUrl = uploadResult.url;
        console.log('图片上传成功:', {
          url: uploadResult.url,
          uploadMethod: uploadResult.uploadMethod,
          filename: uploadResult.filename,
        });
        
        // 如果返回的是base64，给出提示
        if (uploadResult.uploadMethod === 'base64') {
          console.warn('⚠️ 注意：图片使用base64方式，未上传到GitHub。请检查后端GitHub配置。');
        }
      } catch (uploadError: any) {
        console.warn('图片上传失败，降级到 base64:', uploadError);
        // 上传失败时转换为base64，继续流程
        const base64 = await fileToBase64(compressedFile);
        const imageData = extractBase64Data(base64);
        imageUrl = imageData; // 使用base64作为降级方案
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
      // 保存更新后的图片消息到IndexedDB
      const updatedImageMessage: ConversationMessage = {
        ...userMessage,
        content: imageUrl,
      };
      await conversationStorage.saveMessage(updatedImageMessage).catch(err => 
        console.error('保存图片消息失败:', err)
      );
      
      // 6. 创建助手消息占位符
      const assistantMessageId = `msg-assistant-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
      const assistantMessage: ConversationMessage = {
        id: assistantMessageId,
        type: 'text',
        content: '',
        timestamp: new Date().toISOString(),
        sender: 'assistant',
        sessionId,
        isStreaming: true,
        streamingText: '',
      };
      setMessages((prev) => [...prev, assistantMessage]);

      // 7. 构建统一流式请求参数
      // 如果imageUrl是base64，提取纯base64数据；如果是URL，直接使用
      let finalImageData = imageUrl;
      if (imageUrl.startsWith('data:image/')) {
        // 提取base64数据（移除data URL前缀）
        const parts = imageUrl.split(',');
        if (parts.length === 2) {
          finalImageData = parts[1];
        }
      }

      // 获取用户年龄（优先使用识别上下文中的年龄，否则从存储中获取）
      const userAge = identificationContext?.age || getUserAgeFromStorage();
      
      const streamRequest: UnifiedStreamConversationRequest = {
        messageType: 'image',
        image: finalImageData, // 使用base64数据或URL
        sessionId,
        userAge,
        maxContextRounds: 20,
      };

      // 注意：在对话页面中上传图片时，不应该传递identificationContext
      // 因为这不是新的识别结果，而是继续当前会话的对话
      // 只有从Capture页面跳转过来时（新会话）才会传递identificationContext
      // 这里不传递identificationContext，确保不会生成新的知识卡片

      // 重置累积文本和markdown状态
      accumulatedTextRef.current = '';
      markdownRef.current = false;

      // 8. 使用统一流式接口
      const abortController = createStreamConnectionUnified(streamRequest, {
        onMessage: (message: ConversationMessage) => {
          // 更新助手消息
          flushSync(() => {
            setMessages((prev) => {
              const index = prev.findIndex((m) => m.id === assistantMessageId);
              if (index >= 0) {
                const updated = [...prev];
                updated[index] = {
                  ...updated[index],
                  content: message.content || '',
                  streamingText: message.streamingText || '',
                  markdown: message.markdown,
                  isStreaming: message.isStreaming,
                };
                return updated;
              }
              return prev;
            });
          });
        },
        onError: (error: Error) => {
          console.error('流式返回错误:', error);
          flushSync(() => {
            setMessages((prev) => {
              const index = prev.findIndex((m) => m.id === assistantMessageId);
              if (index >= 0) {
                const updated = [...prev];
                updated[index] = {
                  ...updated[index],
                  isStreaming: false,
                  content: accumulatedTextRef.current || t('conversation.generateAnswerError', { error: error.message }),
                  streamingText: undefined,
                };
                return updated;
              }
              return prev;
            });
          });
          setIsSending(false);
        },
        onClose: async () => {
          // 流式输出完成
          let completedMessage: ConversationMessage | null = null;
          flushSync(() => {
            setMessages((prev) => {
              const index = prev.findIndex((m) => m.id === assistantMessageId);
              if (index >= 0) {
                const updated = [...prev];
                // 创建完成的消息对象，清除isStreaming和streamingText字段
                completedMessage = {
                  ...updated[index],
                  isStreaming: false,
                  streamingText: undefined,
                } as ConversationMessage;
                updated[index] = completedMessage;
                return updated;
              }
              return prev;
            });
          });
          
          // 保存完成的消息到IndexedDB（清除isStreaming和streamingText字段）
          if (completedMessage && completedMessage !== null && 'content' in completedMessage) {
            try {
              await conversationStorage.saveMessage(completedMessage);
            } catch (err) {
              console.error('保存流式图片消息失败:', err);
              // 不中断用户操作，仅记录错误
            }
          }
          
          setIsSending(false);
        },
      });

      // 保存abortController以便可以取消
      setSseConnection(abortController);
    } catch (error: any) {
      console.error('发送图片失败:', error);
      // 友好的错误提示
      const errorMessage: ConversationMessage = {
        id: `msg-error-${Date.now()}`,
        type: 'text',
        content: t('conversation.sendImageError', { error: error?.message || t('conversation.unknownError') }),
        timestamp: new Date().toISOString(),
        sender: 'assistant',
        sessionId,
      };
      setMessages((prev) => [...prev, errorMessage]);
      // 保存错误消息到IndexedDB
      await conversationStorage.saveMessage(errorMessage).catch(err => 
        console.error('保存错误消息失败:', err)
      );
      setIsSending(false);
    }
  };

  return (
    <AudioPlaybackProvider>
      <div className="min-h-screen bg-cloud-white font-display flex flex-col">
        <Header />
        
        <main className="flex-1 flex flex-col items-center px-4 py-6 w-full max-w-7xl mx-auto overflow-x-hidden overflow-y-auto">
        {/* 对象信息展示区域 */}
        <div className="flex flex-col md:flex-row items-center justify-between gap-6 mb-8 w-full">
          <div className="flex flex-col items-start gap-2 relative">
            <span className="absolute -top-6 -left-2 rotate-[-5deg] bg-yellow-300 text-yellow-800 px-3 py-1 rounded-lg text-xs font-display shadow-sm border-2 border-yellow-100 transform z-10">
              {t('result.foundNewFriend')}
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
              {t('result.itsA')} <span className="text-transparent bg-clip-text bg-gradient-to-r from-science-green to-green-500 relative inline-block">
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
                {t('result.aiCompanionSays')}
              </p>
              <p className="text-base font-bold text-slate-800 leading-tight font-display">
                {t('result.aiCompanionMessage', { objectName })}
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
    </AudioPlaybackProvider>
  );
}
