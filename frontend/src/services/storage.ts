/**
 * 本地存储服务
 * 使用 localStorage 和 IndexedDB
 */

import type { UserProfile, ExplorationRecord, KnowledgeCard } from '../types/exploration';

const STORAGE_KEYS = {
  USER_PROFILE: 'tango_user_profile',
} as const;

/**
 * 用户档案存储（localStorage）
 */
export const userProfileStorage = {
  /**
   * 获取用户档案
   */
  get(): UserProfile | null {
    const data = localStorage.getItem(STORAGE_KEYS.USER_PROFILE);
    return data ? JSON.parse(data) : null;
  },

  /**
   * 保存用户档案
   */
  save(profile: UserProfile): void {
    localStorage.setItem(STORAGE_KEYS.USER_PROFILE, JSON.stringify(profile));
  },

  /**
   * 清除用户档案
   */
  clear(): void {
    localStorage.removeItem(STORAGE_KEYS.USER_PROFILE);
  },
};

/**
 * IndexedDB 数据库名称和版本
 */
const DB_NAME = 'TanGoDB';
const DB_VERSION = 1;

/**
 * 初始化IndexedDB
 */
function initDB(): Promise<IDBDatabase> {
  return new Promise((resolve, reject) => {
    const request = indexedDB.open(DB_NAME, DB_VERSION);

    request.onerror = () => reject(request.error);
    request.onsuccess = () => resolve(request.result);

    request.onupgradeneeded = (event) => {
      const db = (event.target as IDBOpenDBRequest).result;

      // 探索记录表
      if (!db.objectStoreNames.contains('explorations')) {
        const explorationStore = db.createObjectStore('explorations', {
          keyPath: 'id',
        });
        explorationStore.createIndex('timestamp', 'timestamp', { unique: false });
        explorationStore.createIndex('category', 'objectCategory', { unique: false });
      }

      // 知识卡片表
      if (!db.objectStoreNames.contains('cards')) {
        const cardStore = db.createObjectStore('cards', {
          keyPath: 'id',
        });
        cardStore.createIndex('explorationId', 'explorationId', { unique: false });
        cardStore.createIndex('type', 'type', { unique: false });
      }
    };
  });
}

/**
 * 探索记录存储（IndexedDB）
 */
export const explorationStorage = {
  /**
   * 保存探索记录
   */
  async save(record: ExplorationRecord): Promise<void> {
    const db = await initDB();
    const transaction = db.transaction('explorations', 'readwrite');
    const store = transaction.objectStore('explorations');
    await new Promise<void>((resolve, reject) => {
      const request = store.put(record);
      request.onsuccess = () => resolve();
      request.onerror = () => reject(request.error);
    });
  },

  /**
   * 获取所有探索记录
   */
  async getAll(): Promise<ExplorationRecord[]> {
    const db = await initDB();
    const transaction = db.transaction('explorations', 'readonly');
    const store = transaction.objectStore('explorations');
    const index = store.index('timestamp');
    
    return new Promise((resolve, reject) => {
      const request = index.getAll();
      request.onsuccess = () => resolve(request.result);
      request.onerror = () => reject(request.error);
    });
  },

  /**
   * 根据ID获取探索记录
   */
  async getById(id: string): Promise<ExplorationRecord | null> {
    const db = await initDB();
    const transaction = db.transaction('explorations', 'readonly');
    const store = transaction.objectStore('explorations');
    
    return new Promise((resolve, reject) => {
      const request = store.get(id);
      request.onsuccess = () => resolve(request.result || null);
      request.onerror = () => reject(request.error);
    });
  },

  /**
   * 删除探索记录
   */
  async delete(id: string): Promise<void> {
    const db = await initDB();
    const transaction = db.transaction('explorations', 'readwrite');
    const store = transaction.objectStore('explorations');
    
    await new Promise<void>((resolve, reject) => {
      const request = store.delete(id);
      request.onsuccess = () => resolve();
      request.onerror = () => reject(request.error);
    });
  },
};

/**
 * 知识卡片存储（IndexedDB）
 */
export const cardStorage = {
  /**
   * 保存知识卡片
   */
  async save(card: KnowledgeCard): Promise<void> {
    const db = await initDB();
    const transaction = db.transaction('cards', 'readwrite');
    const store = transaction.objectStore('cards');
    await new Promise<void>((resolve, reject) => {
      const request = store.put(card);
      request.onsuccess = () => resolve();
      request.onerror = () => reject(request.error);
    });
  },

  /**
   * 获取所有收藏的卡片
   */
  async getAll(): Promise<KnowledgeCard[]> {
    const db = await initDB();
    const transaction = db.transaction('cards', 'readonly');
    const store = transaction.objectStore('cards');
    
    return new Promise((resolve, reject) => {
      const request = store.getAll();
      request.onsuccess = () => resolve(request.result);
      request.onerror = () => reject(request.error);
    });
  },

  /**
   * 根据类型获取卡片
   */
  async getByType(type: 'science' | 'poetry' | 'english'): Promise<KnowledgeCard[]> {
    const db = await initDB();
    const transaction = db.transaction('cards', 'readonly');
    const store = transaction.objectStore('cards');
    const index = store.index('type');
    
    return new Promise((resolve, reject) => {
      const request = index.getAll(type);
      request.onsuccess = () => resolve(request.result);
      request.onerror = () => reject(request.error);
    });
  },

  /**
   * 删除卡片
   */
  async delete(id: string): Promise<void> {
    const db = await initDB();
    const transaction = db.transaction('cards', 'readwrite');
    const store = transaction.objectStore('cards');
    
    await new Promise<void>((resolve, reject) => {
      const request = store.delete(id);
      request.onsuccess = () => resolve();
      request.onerror = () => reject(request.error);
    });
  },
};

