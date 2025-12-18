/**
 * 本地存储服务
 * 使用 localStorage 和 IndexedDB
 */

import type { UserProfile, ExplorationRecord, KnowledgeCard } from '../types/exploration';
import type { ConversationMessage, ConversationSession } from '../types/conversation';
import type { UserSettings } from '../types/settings';

const STORAGE_KEYS = {
  USER_PROFILE: 'tango_user_profile',
  USER_SETTINGS: 'tango_user_settings',
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
const DB_VERSION = 2; // 升级版本以添加新表

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

      // 对话消息表
      if (!db.objectStoreNames.contains('conversations')) {
        const conversationStore = db.createObjectStore('conversations', {
          keyPath: 'id',
        });
        conversationStore.createIndex('sessionId', 'sessionId', { unique: false });
        conversationStore.createIndex('timestamp', 'timestamp', { unique: false });
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

  /**
   * 批量保存卡片
   */
  async saveBatch(cards: KnowledgeCard[]): Promise<void> {
    const db = await initDB();
    const transaction = db.transaction('cards', 'readwrite');
    const store = transaction.objectStore('cards');
    
    await Promise.all(
      cards.map(
        (card) =>
          new Promise<void>((resolve, reject) => {
            const request = store.put(card);
            request.onsuccess = () => resolve();
            request.onerror = () => reject(request.error);
          })
      )
    );
  },

  /**
   * 批量删除卡片
   */
  async deleteBatch(ids: string[]): Promise<void> {
    const db = await initDB();
    const transaction = db.transaction('cards', 'readwrite');
    const store = transaction.objectStore('cards');
    
    await Promise.all(
      ids.map(
        (id) =>
          new Promise<void>((resolve, reject) => {
            const request = store.delete(id);
            request.onsuccess = () => resolve();
            request.onerror = () => reject(request.error);
          })
      )
    );
  },
};

/**
 * 用户设置存储（localStorage）
 */
export const userSettingsStorage = {
  /**
   * 获取用户设置
   */
  get(): UserSettings | null {
    const data = localStorage.getItem(STORAGE_KEYS.USER_SETTINGS);
    return data ? JSON.parse(data) : null;
  },

  /**
   * 保存用户设置
   */
  save(settings: UserSettings): void {
    localStorage.setItem(STORAGE_KEYS.USER_SETTINGS, JSON.stringify(settings));
  },

  /**
   * 清除用户设置
   */
  clear(): void {
    localStorage.removeItem(STORAGE_KEYS.USER_SETTINGS);
  },

  /**
   * 获取默认设置
   */
  getDefault(): UserSettings {
    return {
      language: 'zh',
      lastUpdated: new Date().toISOString(),
    };
  },
};

/**
 * 对话消息存储（IndexedDB）
 */
export const conversationStorage = {
  /**
   * 保存对话消息
   */
  async saveMessage(message: ConversationMessage): Promise<void> {
    const db = await initDB();
    const transaction = db.transaction('conversations', 'readwrite');
    const store = transaction.objectStore('conversations');
    await new Promise<void>((resolve, reject) => {
      const request = store.put(message);
      request.onsuccess = () => resolve();
      request.onerror = () => reject(request.error);
    });
  },

  /**
   * 根据会话ID获取所有消息
   */
  async getMessagesBySessionId(sessionId: string): Promise<ConversationMessage[]> {
    const db = await initDB();
    const transaction = db.transaction('conversations', 'readonly');
    const store = transaction.objectStore('conversations');
    const index = store.index('sessionId');
    
    return new Promise((resolve, reject) => {
      const request = index.getAll(sessionId);
      request.onsuccess = () => {
        const messages = request.result || [];
        // 按时间戳排序
        messages.sort((a, b) => new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime());
        resolve(messages);
      };
      request.onerror = () => reject(request.error);
    });
  },

  /**
   * 获取所有会话
   */
  async getAllSessions(): Promise<ConversationSession[]> {
    const db = await initDB();
    const transaction = db.transaction('conversations', 'readonly');
    const store = transaction.objectStore('conversations');
    const index = store.index('sessionId');
    
    return new Promise((resolve, reject) => {
      const request = index.getAll();
      request.onsuccess = () => {
        const messages = request.result || [];
        // 按sessionId分组
        const sessionMap = new Map<string, ConversationMessage[]>();
        messages.forEach((msg) => {
          if (msg.sessionId) {
            if (!sessionMap.has(msg.sessionId)) {
              sessionMap.set(msg.sessionId, []);
            }
            sessionMap.get(msg.sessionId)!.push(msg);
          }
        });
        
        // 转换为会话对象
        const sessions: ConversationSession[] = Array.from(sessionMap.entries()).map(([sessionId, msgs]) => {
          msgs.sort((a, b) => new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime());
          return {
            sessionId,
            messages: msgs,
            createdAt: msgs[0]?.timestamp || new Date().toISOString(),
            lastActive: msgs[msgs.length - 1]?.timestamp || new Date().toISOString(),
          };
        });
        
        // 按最后活跃时间排序
        sessions.sort((a, b) => new Date(b.lastActive).getTime() - new Date(a.lastActive).getTime());
        resolve(sessions);
      };
      request.onerror = () => reject(request.error);
    });
  },

  /**
   * 删除会话的所有消息
   */
  async deleteSession(sessionId: string): Promise<void> {
    const db = await initDB();
    const transaction = db.transaction('conversations', 'readwrite');
    const store = transaction.objectStore('conversations');
    const index = store.index('sessionId');
    
    return new Promise((resolve, reject) => {
      const request = index.openCursor(IDBKeyRange.only(sessionId));
      request.onsuccess = (event) => {
        const cursor = (event.target as IDBRequest<IDBCursorWithValue>).result;
        if (cursor) {
          cursor.delete();
          cursor.continue();
        } else {
          resolve();
        }
      };
      request.onerror = () => reject(request.error);
    });
  },

  /**
   * 删除单条消息
   */
  async deleteMessage(id: string): Promise<void> {
    const db = await initDB();
    const transaction = db.transaction('conversations', 'readwrite');
    const store = transaction.objectStore('conversations');
    
    await new Promise<void>((resolve, reject) => {
      const request = store.delete(id);
      request.onsuccess = () => resolve();
      request.onerror = () => reject(request.error);
    });
  },
};

