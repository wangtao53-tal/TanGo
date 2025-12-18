/**
 * 国际化配置
 * 使用 react-i18next
 */

import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import zh from './locales/zh';
import en from './locales/en';

// 从localStorage获取保存的语言设置，默认中文
const savedLanguage = localStorage.getItem('tango_language') as 'zh' | 'en' | null;
const defaultLanguage = savedLanguage || 'zh';

i18n
  .use(initReactI18next)
  .init({
    resources: {
      zh: {
        translation: zh,
      },
      en: {
        translation: en,
      },
    },
    lng: defaultLanguage,
    fallbackLng: 'zh',
    interpolation: {
      escapeValue: false, // React已经转义了
    },
  });

// 监听语言变化，保存到localStorage
i18n.on('languageChanged', (lng) => {
  localStorage.setItem('tango_language', lng);
});

export default i18n;
