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
  CardContentResponse,
  CreateShareRequest,
  CreateShareResponse,
  GetShareResponse,
  GenerateReportRequest,
  GenerateReportResponse,
  ErrorResponse,
  IntentRequest,
  IntentResponse,
  ConversationRequest,
  VoiceRequest,
  VoiceResponse,
  UploadRequest,
  UploadResponse,
} from '../types/api';

// 从环境变量读取API基础地址
// 生产环境默认使用相对路径（通过 Nginx 代理），开发环境使用完整 URL
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL !== undefined
  ? import.meta.env.VITE_API_BASE_URL
  : (import.meta.env.DEV 
    ? `http://${import.meta.env.VITE_BACKEND_HOST || 'localhost'}:${import.meta.env.VITE_BACKEND_PORT || '8877'}`
    : ''); // 生产环境默认使用相对路径，由 Nginx 代理

const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 600000,
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
 * 生成知识卡片API（同步模式）
 * 超时时间设置为6秒（6000ms），优化后目标5秒内
 */
export async function generateCards(
  request: GenerateCardsRequest
): Promise<GenerateCardsResponse> {
  const response = await apiClient.post<GenerateCardsResponse>(
    '/api/explore/generate-cards', 
    request,
    {
      timeout: 6000, // 6秒 = 6000毫秒
    }
  );
  return response as unknown as GenerateCardsResponse;
}

/**
 * 流式生成知识卡片API
 * 使用SSE流式返回，每生成完一张卡片立即返回
 */
export function generateCardsStream(
  request: GenerateCardsRequest,
  callbacks: {
    onMessage?: (text: string, fullText: string) => void; // 流式文本消息回调（字符，完整文本）
    onCard?: (card: CardContentResponse, index: number) => void;
    onError?: (error: Error) => void;
    onComplete?: () => void;
  }
): AbortController {
  const abortController = new AbortController();
  const API_BASE_URL =
    import.meta.env.VITE_API_BASE_URL !== undefined
      ? import.meta.env.VITE_API_BASE_URL
      : (import.meta.env.DEV
        ? `http://${import.meta.env.VITE_BACKEND_HOST || 'localhost'}:${import.meta.env.VITE_BACKEND_PORT || '8877'}`
        : ''); // 生产环境默认使用相对路径

  fetch(`${API_BASE_URL}/api/explore/generate-cards?stream=true`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(request),
    signal: abortController.signal,
  })
    .then(async (response) => {
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      if (!response.body) {
        throw new Error('Response body is null');
      }

      const reader = response.body.getReader();
      const decoder = new TextDecoder();
      let buffer = '';
      let fullText = ''; // 累积完整文本，用于流式消息

      while (true) {
        const { done, value } = await reader.read();

        if (done) {
          if (buffer.trim()) {
            processCardSSEBuffer(buffer, callbacks, fullText);
          }
          callbacks.onComplete?.();
          break;
        }

        const chunk = decoder.decode(value, { stream: true });
        buffer += chunk;
        const result = processCardSSEBuffer(buffer, callbacks, fullText);
        buffer = result.remainingBuffer;
        fullText = result.fullText;
      }
    })
    .catch((error) => {
      if (error.name === 'AbortError') {
        return;
      }
      console.error('流式卡片生成错误:', error);
      callbacks.onError?.(error);
    });

  return abortController;
}

/**
 * 处理卡片生成SSE缓冲区
 */
function processCardSSEBuffer(
  buffer: string,
  callbacks: {
    onMessage?: (text: string, fullText: string) => void;
    onCard?: (card: CardContentResponse, index: number) => void;
    onError?: (error: Error) => void;
    onComplete?: () => void;
  },
  fullText: string
): { remainingBuffer: string; fullText: string } {
  const lines = buffer.split('\n');
  const remainingLines: string[] = [];
  let currentEvent = '';
  let currentData = '';

  for (let i = 0; i < lines.length; i++) {
    const line = lines[i];

    if (line.startsWith('event:')) {
      if (currentEvent && currentData) {
        const result = processCardSSEEvent(currentEvent, currentData, callbacks, fullText);
        fullText = result.fullText;
        currentEvent = '';
        currentData = '';
      }
      currentEvent = line.substring(6).trim();
    } else if (line.startsWith('data:')) {
      const dataLine = line.substring(5).trim();
      if (currentData) {
        currentData += '\n' + dataLine;
      } else {
        currentData = dataLine;
      }
    } else if (line === '') {
      if (currentEvent && currentData) {
        const result = processCardSSEEvent(currentEvent, currentData, callbacks, fullText);
        fullText = result.fullText;
        currentEvent = '';
        currentData = '';
      }
    } else if (line.trim() !== '') {
      if (currentData) {
        currentData += '\n' + line;
      }
    }
  }

  // 如果循环结束后还有完整的事件，处理它
  if (currentEvent && currentData) {
    const result = processCardSSEEvent(currentEvent, currentData, callbacks, fullText);
    fullText = result.fullText;
    currentEvent = '';
    currentData = '';
  }

  // 如果还有不完整的事件，保留到下次处理
  if (currentEvent || currentData) {
    if (currentEvent) {
      remainingLines.push(`event: ${currentEvent}`);
    }
    if (currentData) {
      remainingLines.push(`data: ${currentData}`);
    }
  }

  return {
    remainingBuffer: remainingLines.length > 0 ? remainingLines.join('\n') : '',
    fullText,
  };
}

/**
 * 处理卡片生成SSE事件
 */
function processCardSSEEvent(
  eventType: string,
  dataStr: string,
  callbacks: {
    onMessage?: (text: string, fullText: string) => void;
    onCard?: (card: CardContentResponse, index: number) => void;
    onError?: (error: Error) => void;
    onComplete?: () => void;
  },
  fullText: string
): { fullText: string } {
  try {
    const data = JSON.parse(dataStr);

    if (eventType === 'message' && data.content) {
      // 处理流式文本消息（逐字符返回）
      const char = typeof data.content === 'string' ? data.content : String(data.content);
      fullText += char;
      callbacks.onMessage?.(char, fullText);
    } else if (eventType === 'card' && data.content) {
      const card: CardContentResponse = {
        type: data.content.type,
        title: data.content.title,
        content: data.content.content,
      };
      callbacks.onCard?.(card, data.index || 0);
    } else if (eventType === 'done') {
      callbacks.onComplete?.();
    } else if (eventType === 'error') {
      const errorMessage =
        (data.content as any)?.message || data.message || '未知错误';
      callbacks.onError?.(new Error(errorMessage));
    }
  } catch (err) {
    console.error('解析卡片SSE消息失败:', err);
  }
  return { fullText };
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

