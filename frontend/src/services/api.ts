/**
 * API服务封装
 * 当前使用Mock数据，待后端API就绪后替换
 */

import axios from 'axios';
import type {
  IdentifyRequest,
  IdentifyResponse,
  GenerateCardsRequest,
  GenerateCardsResponse,
  CreateShareRequest,
  CreateShareResponse,
  GetShareResponse,
  GenerateReportRequest,
  GenerateReportResponse,
  ErrorResponse,
  IntentRequest,
  IntentResponse,
  ConversationRequest,
  ConversationStreamEvent,
  VoiceRequest,
  VoiceResponse,
} from '../types/api';

// 从环境变量读取API基础地址，如果没有配置则使用默认值
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 
  (import.meta.env.DEV 
    ? `http://${import.meta.env.VITE_BACKEND_HOST || 'localhost'}:${import.meta.env.VITE_BACKEND_PORT || '8877'}`
    : 'http://localhost:8877');

const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// 请求拦截器
apiClient.interceptors.request.use(
  (config) => {
    // 可以在这里添加token等
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 响应拦截器
apiClient.interceptors.response.use(
  (response) => response.data,
  (error) => {
    // 统一错误处理
    const errorResponse: ErrorResponse = {
      code: error.response?.status || 500,
      message: error.response?.data?.message || error.message || '请求失败',
      detail: error.response?.data?.detail,
    };
    return Promise.reject(errorResponse);
  }
);

/**
 * 图像识别API
 */
export async function identifyImage(
  request: IdentifyRequest
): Promise<IdentifyResponse> {
  const response = await apiClient.post<IdentifyResponse>('/api/explore/identify', request);
  return response as unknown as IdentifyResponse;
}

/**
 * 生成知识卡片API
 */
export async function generateCards(
  request: GenerateCardsRequest
): Promise<GenerateCardsResponse> {
  const response = await apiClient.post<GenerateCardsResponse>('/api/explore/generate-cards', request);
  return response as unknown as GenerateCardsResponse;
}

/**
 * 创建分享链接
 */
export async function createShare(
  request: CreateShareRequest
): Promise<CreateShareResponse> {
  const response = await apiClient.post<CreateShareResponse>('/api/share/create', request);
  return response as unknown as CreateShareResponse;
}

/**
 * 获取分享数据
 */
export async function getShare(shareId: string): Promise<GetShareResponse> {
  const response = await apiClient.get<GetShareResponse>(`/api/share/${shareId}`);
  return response as unknown as GetShareResponse;
}

/**
 * 生成学习报告
 */
export async function generateReport(
  request: GenerateReportRequest
): Promise<GenerateReportResponse> {
  const response = await apiClient.post<GenerateReportResponse>('/api/share/report', request);
  return response as unknown as GenerateReportResponse;
}

/**
 * 意图识别API
 */
export async function recognizeIntent(
  request: IntentRequest
): Promise<IntentResponse> {
  try {
    // 转换字段名：前端使用 text，后端使用 message
    const backendRequest = {
      message: request.text,
      sessionId: request.sessionId,
      context: request.context,
    };
    const response = await apiClient.post<IntentResponse>('/api/conversation/intent', backendRequest);
    return response as unknown as IntentResponse;
  } catch (error: any) {
    // 如果后端不可用，降级到Mock数据
    console.warn('后端API调用失败，使用Mock数据:', error.message);
    await new Promise((resolve) => setTimeout(resolve, 500));
    
    // 简单的关键词匹配作为Mock
    const text = request.text.toLowerCase();
    if (text.includes('卡片') || text.includes('card') || text.includes('生成')) {
      return {
        intent: 'generate_cards',
        confidence: 0.9,
        parameters: {},
      };
    }
    return {
      intent: 'text_response',
      confidence: 0.8,
      parameters: {},
    };
  }
}

/**
 * 对话API（非流式）
 */
export async function sendConversationMessage(
  request: ConversationRequest
): Promise<any> {
  try {
    const response = await apiClient.post('/api/conversation/message', request);
    return response;
  } catch (error: any) {
    console.warn('后端API调用失败，使用Mock数据:', error.message);
    await new Promise((resolve) => setTimeout(resolve, 1000));
    
    // Mock响应
    return {
      type: 'text',
      content: '这是一个Mock响应。',
    };
  }
}

/**
 * 语音识别API
 */
export async function recognizeVoice(
  request: VoiceRequest
): Promise<VoiceResponse> {
  try {
    const response = await apiClient.post<VoiceResponse>('/api/conversation/voice', request);
    return response as unknown as VoiceResponse;
  } catch (error: any) {
    // 如果后端不可用，降级到Mock数据
    console.warn('后端API调用失败，使用Mock数据:', error.message);
    await new Promise((resolve) => setTimeout(resolve, 1500));
    
    return {
      text: '这是语音识别的Mock结果',
      intent: 'text_response',
      confidence: 0.8,
    };
  }
}

