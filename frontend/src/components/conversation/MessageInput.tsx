/**
 * 消息输入组件
 * 支持文本输入
 */

import { useState, KeyboardEvent } from 'react';
import { useTranslation } from 'react-i18next';

export interface MessageInputProps {
  onSend: (text: string) => void;
  disabled?: boolean;
}

export function MessageInput({ onSend, disabled = false }: MessageInputProps) {
  const { t } = useTranslation();
  const [input, setInput] = useState('');

  const handleSend = () => {
    if (input.trim() && !disabled) {
      onSend(input.trim());
      setInput('');
    }
  };

  const handleKeyPress = (e: KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  return (
    <div className="flex items-center gap-2 px-4 py-3 bg-white border-t border-gray-200">
      <input
        type="text"
        value={input}
        onChange={(e) => setInput(e.target.value)}
        onKeyPress={handleKeyPress}
        placeholder={t('conversation.placeholder')}
        disabled={disabled}
        className="flex-1 px-4 py-2 rounded-full border-2 border-gray-200 focus:border-[var(--color-primary)] focus:outline-none disabled:opacity-50"
      />
      <button
        onClick={handleSend}
        disabled={disabled || !input.trim()}
        className="px-6 py-2 rounded-full bg-[var(--color-primary)] text-white font-bold disabled:opacity-50 disabled:cursor-not-allowed hover:bg-[#5aff2b] transition-colors"
      >
        {t('conversation.send')}
      </button>
    </div>
  );
}
