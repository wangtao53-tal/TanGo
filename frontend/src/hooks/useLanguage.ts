/**
 * 语言切换Hook
 * 提供语言切换和当前语言状态
 */

import { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { userSettingsStorage } from '../services/storage';
import type { Language } from '../types/settings';

export function useLanguage() {
  const { i18n } = useTranslation();
  const [currentLanguage, setCurrentLanguage] = useState<Language>('zh');

  useEffect(() => {
    // 初始化时从localStorage读取语言设置
    const settings = userSettingsStorage.get();
    if (settings?.language) {
      setCurrentLanguage(settings.language);
      i18n.changeLanguage(settings.language);
    } else {
      // 使用默认语言
      const defaultSettings = userSettingsStorage.getDefault();
      userSettingsStorage.save(defaultSettings);
      setCurrentLanguage(defaultSettings.language);
    }
  }, [i18n]);

  const changeLanguage = async (lang: Language) => {
    await i18n.changeLanguage(lang);
    setCurrentLanguage(lang);
    
    // 保存到localStorage
    const settings = userSettingsStorage.get() || userSettingsStorage.getDefault();
    userSettingsStorage.save({
      ...settings,
      language: lang,
      lastUpdated: new Date().toISOString(),
    });
  };

  return {
    currentLanguage,
    changeLanguage,
    t: i18n.t,
  };
}
