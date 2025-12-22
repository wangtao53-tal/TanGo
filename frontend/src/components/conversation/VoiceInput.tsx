/**
 * 语音输入组件
 * 支持语音识别
 */

import { useState, useRef, useEffect, useCallback } from 'react';
import { useTranslation } from 'react-i18next';

export interface VoiceInputProps {
  onResult: (text: string) => void;
  disabled?: boolean;
}

export function VoiceInput({ onResult, disabled = false }: VoiceInputProps) {
  const { t } = useTranslation();
  const [isListening, setIsListening] = useState(false);
  const recognitionRef = useRef<any>(null);
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);
  const onResultRef = useRef(onResult);

  // 保持 onResult 引用最新
  useEffect(() => {
    onResultRef.current = onResult;
  }, [onResult]);

  // 停止识别的辅助函数
  const stopRecognition = useCallback(() => {
    if (recognitionRef.current) {
      try {
        // 先尝试 abort，如果失败再尝试 stop
        if (typeof recognitionRef.current.abort === 'function') {
          recognitionRef.current.abort();
        } else {
          recognitionRef.current.stop();
        }
      } catch (error) {
        console.warn('停止语音识别时出错:', error);
        // 即使出错也继续清理
      }
      recognitionRef.current = null;
    }
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
      timeoutRef.current = null;
    }
    // 使用 setTimeout 延迟状态更新，避免在事件处理中直接更新
    setTimeout(() => {
      setIsListening(false);
    }, 0);
  }, []);

  const handleVoiceClick = () => {
    if (disabled) return;

    // 如果已有识别实例在运行，点击按钮停止识别
    if (recognitionRef.current) {
      console.log('用户手动停止语音识别');
      stopRecognition();
      return;
    }

    // 启动语音识别
    if ('webkitSpeechRecognition' in window || 'SpeechRecognition' in window) {
      const SpeechRecognition = (window as any).webkitSpeechRecognition || (window as any).SpeechRecognition;
      const recognition = new SpeechRecognition();
      recognition.lang = 'zh-CN';
      recognition.continuous = false;
      recognition.interimResults = false;

      recognitionRef.current = recognition;

      recognition.onstart = () => {
        // 使用 setTimeout 延迟状态更新
        setTimeout(() => {
          setIsListening(true);
        }, 0);
        // 设置超时，30秒后自动停止
        timeoutRef.current = setTimeout(() => {
          console.warn('语音识别超时，自动停止');
          stopRecognition();
          setTimeout(() => {
            alert(t('conversation.voiceTimeout', '语音识别超时，请重试'));
          }, 0);
        }, 30000);
      };

      recognition.onresult = (event: any) => {
        try {
          // 先停止识别和超时
          if (timeoutRef.current) {
            clearTimeout(timeoutRef.current);
            timeoutRef.current = null;
          }
          
          // 检查是否有识别结果
          if (!event.results || event.results.length === 0) {
            console.warn('语音识别结果为空');
            stopRecognition();
            return;
          }
          
          const transcript = event.results[0][0].transcript;
          console.log('语音识别成功:', transcript);
          
          // 停止识别
          if (recognitionRef.current) {
            try {
              if (typeof recognitionRef.current.abort === 'function') {
                recognitionRef.current.abort();
              } else {
                recognitionRef.current.stop();
              }
            } catch (error) {
              console.warn('停止识别时出错:', error);
            }
            recognitionRef.current = null;
          }
          
          // 使用 setTimeout 延迟状态更新和回调
          setTimeout(() => {
            setIsListening(false);
            // 确保识别结果自动发送到对话框
            if (transcript && transcript.trim()) {
              onResultRef.current(transcript.trim());
            }
          }, 0);
        } catch (error) {
          console.error('处理语音识别结果失败:', error);
          stopRecognition();
        }
      };

      recognition.onerror = (event: any) => {
        // 如果是用户主动取消（aborted），不处理错误
        if (event.error === 'aborted') {
          console.log('语音识别已取消（用户主动停止）');
          // 清理资源，但不显示错误
          if (timeoutRef.current) {
            clearTimeout(timeoutRef.current);
            timeoutRef.current = null;
          }
          // 确保状态更新
          setTimeout(() => {
            setIsListening(false);
          }, 0);
          recognitionRef.current = null;
          return;
        }
        
        // 其他错误才打印错误日志
        console.error('语音识别错误:', event.error);
        
        // 先停止识别
        stopRecognition();
        
        // 根据错误类型提供友好的错误提示
        let errorMessage = t('conversation.voiceError', '语音识别失败');
        switch (event.error) {
          case 'network':
            errorMessage = t('conversation.voiceNetworkError', '网络连接失败，请检查网络后重试');
            break;
          case 'no-speech':
            errorMessage = t('conversation.voiceNoSpeech', '未检测到语音，请重试');
            break;
          case 'audio-capture':
            errorMessage = t('conversation.voiceAudioError', '无法访问麦克风，请检查权限设置');
            break;
          case 'not-allowed':
            errorMessage = t('conversation.voiceNotAllowed', '麦克风权限被拒绝，请在浏览器设置中允许访问');
            break;
          default:
            errorMessage = t('conversation.voiceError', '语音识别失败，请重试');
        }
        
        // 使用 setTimeout 延迟显示错误，避免在渲染期间触发
        setTimeout(() => {
          alert(errorMessage);
        }, 0);
      };

      recognition.onend = () => {
        console.log('语音识别结束');
        // 清理超时
        if (timeoutRef.current) {
          clearTimeout(timeoutRef.current);
          timeoutRef.current = null;
        }
        // 使用 setTimeout 延迟状态更新
        setTimeout(() => {
          setIsListening(false);
        }, 0);
        recognitionRef.current = null;
      };

      try {
        recognition.start();
      } catch (error) {
        console.error('启动语音识别失败:', error);
        stopRecognition();
        setTimeout(() => {
          alert(t('conversation.voiceStartError', '启动语音识别失败，请重试'));
        }, 0);
      }
    } else {
      alert(t('conversation.voiceNotSupported', '您的浏览器不支持语音识别功能'));
    }
  };

  // 组件卸载时清理语音识别
  useEffect(() => {
    return () => {
      stopRecognition();
    };
  }, [stopRecognition]);

  return (
    <button
      onClick={handleVoiceClick}
      disabled={disabled}
      className={`p-3 rounded-full transition-all ${
        isListening
          ? 'bg-red-500 text-white animate-pulse hover:bg-red-600'
          : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
      } disabled:opacity-50 disabled:cursor-not-allowed`}
      title={isListening ? t('conversation.voiceStop', '点击停止语音识别') : t('conversation.voiceInput')}
    >
      <span className="material-symbols-outlined text-2xl">
        {isListening ? 'mic' : 'mic_none'}
      </span>
    </button>
  );
}
