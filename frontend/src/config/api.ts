/**
 * API配置模块
 * 支持环境变量和localStorage配置，用于选择使用单Agent模式或多Agent模式
 */

/**
 * API配置对象
 */
export const API_CONFIG = {
  /**
   * 是否使用多Agent模式
   * 优先级：localStorage > 环境变量 > 默认值(false)
   */
  get useMultiAgent(): boolean {
    // 1. 优先从localStorage读取
    const localStorageValue = localStorage.getItem('useMultiAgent');
    if (localStorageValue !== null) {
      return localStorageValue === 'true';
    }

    // 2. 从环境变量读取
    const envValue = import.meta.env.VITE_USE_MULTI_AGENT;
    if (envValue !== undefined) {
      return envValue === 'true';
    }

    // 3. 默认值：false（使用单Agent模式，向后兼容）
    return false;
  },

  /**
   * 设置是否使用多Agent模式
   * @param enabled 是否启用多Agent模式
   */
  setUseMultiAgent(enabled: boolean): void {
    localStorage.setItem('useMultiAgent', enabled.toString());
  },

  /**
   * 根据配置获取对话接口路径
   * @returns 接口路径
   */
  getConversationEndpoint(): string {
    return this.useMultiAgent 
      ? '/api/conversation/agent' 
      : '/api/conversation/stream';
  },

  /**
   * 获取当前使用的模式名称
   * @returns 模式名称
   */
  getCurrentMode(): 'single-agent' | 'multi-agent' {
    return this.useMultiAgent ? 'multi-agent' : 'single-agent';
  },
};

