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
  UploadRequest,
  UploadResponse,
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
    console.error('API请求错误:', {
      url: error.config?.url,
      method: error.config?.method,
      status: error.response?.status,
      statusText: error.response?.statusText,
      data: error.response?.data,
      message: error.message,
      baseURL: error.config?.baseURL,
    });
    
    const errorResponse: ErrorResponse = {
      code: error.response?.status || 500,
      message: error.response?.data?.message || error.message || '请求失败',
      detail: error.response?.data?.detail || error.response?.statusText,
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
 * 超时时间设置为3分钟（180000ms），因为卡片生成可能需要较长时间
 */
export async function generateCards(
  request: GenerateCardsRequest
): Promise<GenerateCardsResponse> {
  const response = await apiClient.post<GenerateCardsResponse>(
    '/api/explore/generate-cards', 
    request,
    {
      timeout: 180000, // 3分钟 = 180000毫秒
    }
  );
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
    // 转换请求格式以匹配后端API
    const backendRequest: any = {
      message: request.content,
      sessionId: request.sessionId,
      identificationContext: request.identificationContext,
    };
    
    // 根据类型设置相应的字段
    if (request.type === 'image') {
      backendRequest.image = request.content;
    } else if (request.type === 'voice') {
      backendRequest.voice = request.content;
    }
    
    const response = await apiClient.post('/api/conversation/message', backendRequest);
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

/**
 * 图片上传API
 */
export async function uploadImage(
  request: UploadRequest
): Promise<UploadResponse> {
  try {
    console.log('开始上传图片，请求数据:', {
      imageDataLength: request.imageData?.length || 0,
      filename: request.filename,
      baseURL: API_BASE_URL,
    });
    
    const response = await apiClient.post<UploadResponse>('/api/upload/image', request);
    
    console.log('图片上传成功:', response);
    return response as unknown as UploadResponse;
  } catch (error: any) {
    console.error('图片上传失败:', {
      error,
      code: error?.code,
      message: error?.message,
      detail: error?.detail,
      response: error?.response,
      request: {
        imageDataLength: request.imageData?.length || 0,
        filename: request.filename,
      },
    });
    throw error;
  }
}

