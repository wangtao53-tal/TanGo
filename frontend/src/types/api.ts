/**
 * API相关类型定义
 * 基于 contracts/explore.api
 */

// 图像识别请求
export interface IdentifyRequest {
  image: string; // base64编码的图片数据
  age?: number; // 可选：孩子年龄
}

// 图像识别响应
export interface IdentifyResponse {
  objectName: string; // 对象名称（中文）
  objectCategory: '自然类' | '生活类' | '人文类';
  confidence: number; // 识别置信度 0-1
  keywords?: string[]; // 相关关键词
}

// 知识卡片生成请求
export interface GenerateCardsRequest {
  objectName: string;
  objectCategory: '自然类' | '生活类' | '人文类';
  age: number; // 必填
  keywords?: string[];
}

// 知识卡片内容（API响应）
export interface CardContentResponse {
  type: 'science' | 'poetry' | 'english';
  title: string;
  content: Record<string, unknown>; // 根据类型不同结构不同
}

// 知识卡片生成响应
export interface GenerateCardsResponse {
  cards: CardContentResponse[];
}

// 创建分享链接请求
export interface CreateShareRequest {
  explorationRecords: ExplorationRecordForShare[];
  collectedCards: KnowledgeCardForShare[];
}

// 探索记录（分享用）
export interface ExplorationRecordForShare {
  id: string;
  timestamp: string;
  objectName: string;
  objectCategory: '自然类' | '生活类' | '人文类';
  age: number;
  cards: CardContentResponse[];
}

// 知识卡片（分享用）
export interface KnowledgeCardForShare {
  id: string;
  explorationId: string;
  type: 'science' | 'poetry' | 'english';
  title: string;
  content: Record<string, unknown>;
  collectedAt?: string;
}

// 创建分享链接响应
export interface CreateShareResponse {
  shareId: string;
  shareUrl: string;
  expiresAt: string;
}

// 获取分享数据响应
export interface GetShareResponse {
  explorationRecords: ExplorationRecordForShare[];
  collectedCards: KnowledgeCardForShare[];
  createdAt: string;
  expiresAt: string;
}

// 生成学习报告请求
export interface GenerateReportRequest {
  shareId: string;
}

// 学习报告响应
export interface GenerateReportResponse {
  totalExplorations: number;
  totalCollectedCards: number;
  categoryDistribution: Record<string, number>;
  recentCards: KnowledgeCardForShare[];
  generatedAt: string;
}

// 错误响应
export interface ErrorResponse {
  code: number;
  message: string;
  detail?: string;
}

// 意图识别请求
export interface IntentRequest {
  text: string; // 用户输入的文本（前端使用，会转换为 message 发送给后端）
  sessionId?: string; // 对话会话ID
  context?: any[]; // 上下文消息（可选）
}

// 意图识别响应（对应后端的 IntentResult）
export interface IntentResponse {
  intent: 'generate_cards' | 'text_response' | 'image_recognition'; // 意图类型
  confidence: number; // 置信度 0-1
  reason?: string; // 识别原因（对应后端的 reason 字段）
  parameters?: Record<string, any>; // 意图参数（可选，前端扩展字段）
}

// 对话请求
export interface ConversationRequest {
  sessionId?: string; // 对话会话ID（可选，首次请求不传）
  type: 'text' | 'voice' | 'image'; // 输入类型
  content: string; // 内容（文本、base64音频、base64图片）
  inputType: 'text' | 'voice' | 'image'; // 输入类型
}

// 对话响应（流式返回）
export interface ConversationStreamEvent {
  type: 'text' | 'card' | 'image' | 'error' | 'done'; // 事件类型
  content?: any; // 内容（根据类型不同）
  message?: string; // 错误消息
}

// 语音识别请求
export interface VoiceRequest {
  audio: string; // base64编码的音频数据
  sessionId?: string; // 对话会话ID
}

// 语音识别响应
export interface VoiceResponse {
  text: string; // 识别出的文本
  intent?: 'generate_cards' | 'text_response' | 'image_recognition'; // 识别出的意图（可选）
  confidence?: number; // 置信度（可选）
}

// 图片上传请求
export interface UploadRequest {
  imageData: string; // base64编码的图片数据（不含 data URL 前缀）
  filename?: string; // 可选：文件名
}

// 图片上传响应
export interface UploadResponse {
  url: string; // 图片的访问URL（GitHub raw URL 或 base64 data URL）
  filename: string; // 实际存储的文件名
  size?: number; // 图片大小（字节）
  uploadMethod?: 'github' | 'base64'; // 上传方式
}

