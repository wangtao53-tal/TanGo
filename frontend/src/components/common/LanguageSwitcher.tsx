/**
 * 语言切换组件
 * 支持中英文切换
 */

import { useLanguage } from '../../hooks/useLanguage';
import type { Language } from '../../types/settings';

export function LanguageSwitcher() {
  const { currentLanguage, changeLanguage } = useLanguage();

  const handleLanguageChange = (lang: Language) => {
    changeLanguage(lang);
  };

  return (
    <div className="flex items-center gap-2 bg-white rounded-full border-2 border-gray-200 p-1 shadow-sm">
      <button
        onClick={() => handleLanguageChange('zh')}
        className={`px-4 py-2 rounded-full font-bold text-sm transition-all ${
          currentLanguage === 'zh'
            ? 'bg-[var(--color-primary)] text-white shadow-md'
            : 'text-gray-600 hover:bg-gray-100'
        }`}
      >
        中文
      </button>
      <button
        onClick={() => handleLanguageChange('en')}
        className={`px-4 py-2 rounded-full font-bold text-sm transition-all ${
          currentLanguage === 'en'
            ? 'bg-[var(--color-primary)] text-white shadow-md'
            : 'text-gray-600 hover:bg-gray-100'
        }`}
      >
        English
      </button>
    </div>
  );
}
