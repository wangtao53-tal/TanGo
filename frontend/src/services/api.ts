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
} from '../types/api';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8888';

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
 * 图像识别API（当前使用Mock数据）
 */
export async function identifyImage(
  _request: IdentifyRequest
): Promise<IdentifyResponse> {
  // TODO: 待后端API就绪后替换
  // return apiClient.post<IdentifyResponse>('/api/explore/identify', request);
  
  // Mock数据
  await new Promise((resolve) => setTimeout(resolve, 1000)); // 模拟网络延迟
  
  return {
    objectName: '银杏',
    objectCategory: '自然类',
    confidence: 0.95,
    keywords: ['植物', '古老', '叶子'],
  };
}

/**
 * 生成知识卡片API（当前使用Mock数据）
 */
export async function generateCards(
  request: GenerateCardsRequest
): Promise<GenerateCardsResponse> {
  // TODO: 待后端API就绪后替换
  // return apiClient.post<GenerateCardsResponse>('/api/explore/generate-cards', request);
  
  // Mock数据
  await new Promise((resolve) => setTimeout(resolve, 2000)); // 模拟网络延迟
  
  return {
    cards: [
      {
        type: 'science',
        title: request.objectName,
        content: {
          name: request.objectName,
          explanation: `${request.objectName}是非常古老的植物，已经在地球上生存了2亿多年。`,
          facts: [
            '生长在阳光充足的地方',
            '叶子像扇子一样',
          ],
          funFact: '银杏的每一部分都可以食用！',
        },
      },
      {
        type: 'poetry',
        title: '相关诗词',
        content: {
          poem: '轻飘飘地随风飘，飞向远方...',
          explanation: '古人认为银杏是勇敢的旅行者！它们离开家，只靠一点风就能飞得很远。',
          context: '就像小小的降落伞冒险！',
        },
      },
      {
        type: 'english',
        title: request.objectName,
        content: {
          words: ['Flower', 'Seed', 'Blow'],
          expressions: [
            `Look at the yellow ${request.objectName}!`,
            'Make a wish and blow!',
          ],
        },
      },
    ],
  };
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

