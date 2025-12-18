/**
 * 语音输入组件
 * 支持语音识别
 */

import { useState } from 'react';
import { useTranslation } from 'react-i18next';

export interface VoiceInputProps {
  onResult: (text: string) => void;
  disabled?: boolean;
}

export function VoiceInput({ onResult, disabled = false }: VoiceInputProps) {
  const { t } = useTranslation();
  const [isListening, setIsListening] = useState(false);

  const handleVoiceClick = () => {
    if (disabled || isListening) return;

    // 启动语音识别
    if ('webkitSpeechRecognition' in window || 'SpeechRecognition' in window) {
      const SpeechRecognition = (window as any).webkitSpeechRecognition || (window as any).SpeechRecognition;
      const recognition = new SpeechRecognition();
      recognition.lang = 'zh-CN';
      recognition.continuous = false;
      recognition.interimResults = false;

      recognition.onstart = () => {
        setIsListening(true);
      };

      recognition.onresult = (event: any) => {
        const transcript = event.results[0][0].transcript;
        setIsListening(false);
        onResult(transcript);
      };

      recognition.onerror = (event: any) => {
        console.error('语音识别错误:', event.error);
        setIsListening(false);
      };

      recognition.onend = () => {
        setIsListening(false);
      };

      recognition.start();
    } else {
      alert('您的浏览器不支持语音识别功能');
    }
  };

  return (
    <button
      onClick={handleVoiceClick}
      disabled={disabled || isListening}
      className={`p-3 rounded-full transition-all ${
        isListening
          ? 'bg-red-500 text-white animate-pulse'
          : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
      } disabled:opacity-50 disabled:cursor-not-allowed`}
      title={t('conversation.voiceInput')}
    >
      <span className="material-symbols-outlined text-2xl">
        {isListening ? 'mic' : 'mic_none'}
      </span>
    </button>
  );
}
