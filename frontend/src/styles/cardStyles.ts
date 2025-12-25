/**
 * 卡片专用样式配置
 * 儿童友好设计：字体、色彩、固定比例
 */

export const cardStyles = {
  // 儿童友好字体配置
  fonts: {
    childFriendly: {
      chinese: '"Comfortaa", "PingFang SC", "Microsoft YaHei UI", sans-serif',
      english: '"Comfortaa", "Nunito", sans-serif',
    },
    sizes: {
      title: 'clamp(20px, 5vw, 28px)',      // 标题：20-28px（增大最小值）
      body: 'clamp(16px, 4vw, 18px)',      // 正文：16-18px（增大最小值）
      small: 'clamp(14px, 3vw, 16px)',   // 小字：14-16px（增大最小值）
    },
    lineHeight: {
      title: 1.4,
      body: 1.6,
      small: 1.5,
    },
  },
  // 固定比例配置
  aspectRatio: '16/9',
  // 色彩配置（增强对比度）
  colors: {
    // 保持现有主题色，但增强对比度
    scienceGreen: '#76FF7A',
    sunnyOrange: '#FF9E64',
    skyBlue: '#40C4FF',
    // 文本颜色（确保对比度≥4.5:1）
    textDark: '#1F2937',
    textLight: '#FFFFFF',
  },
} as const;

