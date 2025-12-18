/**
 * 主题配置
 * 基于设计稿中的颜色和字体系统
 */

export const theme = {
  colors: {
    primary: '#4cdf20',
    'cloud-white': '#F8F8F8',
    'science-green': '#76FF7A',
    'sunny-orange': '#FF9E64',
    'sky-blue': '#40C4FF',
    'warm-yellow': '#FFE580',
    'peach-pink': '#FFB7C5',
    'text-main': '#1F2937',
    'text-sub': '#6B7280',
    'white-card': '#FFFFFF',
  },
  fonts: {
    display: ['Manrope', 'Noto Sans SC', 'sans-serif'],
    body: ['Noto Sans', 'Quicksand', 'sans-serif'],
  },
  borderRadius: {
    DEFAULT: '1rem',
    lg: '2rem',
    xl: '3rem',
    '2xl': '4rem',
    full: '9999px',
  },
} as const;

