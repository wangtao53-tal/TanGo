/**
 * 首页组件
 * 基于 stitch_ui/homepage_/_main_interface/ 设计稿
 */

import { useNavigate } from 'react-router-dom';
import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Header } from '../components/common/Header';
import { LittleStar } from '../components/common/LittleStar';

export default function Home() {
  const navigate = useNavigate();
  const { t } = useTranslation();
  const [isListening, setIsListening] = useState(false);

  const handlePhotoClick = () => {
    navigate('/capture');
  };

  const handleVoiceClick = () => {
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
        console.log('语音识别结果:', transcript);
        // TODO: 处理语音识别结果，可以导航到对话页面或触发搜索
        setIsListening(false);
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
      // 如果不支持Web Speech API，导航到拍照页面使用语音输入功能
      navigate('/capture');
    }
  };

  return (
    <div className="relative flex min-h-screen w-full flex-col">
      <Header />
      
      <main className="flex flex-1 w-full max-w-[800px] flex-col items-center justify-center px-4 pb-24 gap-8 mx-auto">
        <div className="flex flex-col items-center justify-center gap-8 w-full py-8">
          {/* 大圆形拍照按钮 */}
          <div className="relative group cursor-pointer">
            <div className="absolute -inset-1 rounded-full bg-[var(--color-primary)] opacity-30 blur-xl group-hover:opacity-50 transition-opacity duration-500 animate-pulse"></div>
            <button
              onClick={handlePhotoClick}
              className="relative flex flex-col items-center justify-center size-48 sm:size-64 rounded-full bg-[var(--color-primary)] hover:bg-[#5aff2b] text-[#0a3f00] shadow-[var(--shadow-glow)] animate-[var(--animate-pulse-glow)] transition-transform active:scale-95 border-4 border-white/30"
            >
              <span className="material-symbols-outlined text-6xl sm:text-8xl mb-2 drop-shadow-sm">
                photo_camera
              </span>
              <span className="text-lg sm:text-2xl font-bold tracking-tight drop-shadow-sm font-display">
                {t('home.photoButton')}
              </span>
            </button>
          </div>

          {/* 语音触发按钮 */}
          <div className="relative z-10 -mt-8 sm:-mt-12 ml-32 sm:ml-48 animate-float">
            <button
              onClick={handleVoiceClick}
              disabled={isListening}
              className={`flex items-center gap-3 rounded-full bg-[var(--color-sky-blue)] hover:bg-[#60d0ff] text-white py-3 px-6 sm:py-4 sm:px-8 shadow-[var(--shadow-glow-blue)] transition-transform hover:scale-105 active:scale-95 border-2 border-white/20 ${isListening ? 'animate-pulse' : ''}`}
            >
              <span className="material-symbols-outlined text-2xl">{isListening ? 'mic' : 'mic'}</span>
              <span className="text-base sm:text-lg font-bold whitespace-nowrap font-display">
                {isListening ? t('conversation.thinking', 'Listening...') : t('home.voiceButton')}
              </span>
            </button>
          </div>
        </div>

        {/* 三个功能卡片展示区域 */}
        <div className="w-full grid grid-cols-1 sm:grid-cols-3 gap-4 mt-4">
          {/* 科学认知卡片 */}
          <div className="group relative cursor-pointer">
            <div className="absolute inset-0 bg-science-green/30 rounded-3xl blur-md opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
            <div className="relative flex flex-col items-center gap-3 bg-white border border-gray-100 group-hover:border-[var(--color-science-green)]/50 p-5 rounded-3xl hover:-translate-y-1 transition-all duration-300 h-full shadow-[var(--shadow-card)] group-hover:shadow-[var(--shadow-card-hover-green)]">
              <div className="size-16 rounded-full bg-[var(--color-science-green)]/10 flex items-center justify-center text-[var(--color-science-green)] group-hover:scale-110 transition-transform duration-300">
                <span className="material-symbols-outlined text-3xl">science</span>
              </div>
              <h3 className="text-[var(--color-text-main)] text-lg font-bold text-center group-hover:text-[var(--color-science-green)] transition-colors font-[var(--font-display)]">
                科学认知
              </h3>
            </div>
          </div>

          {/* 人文素养卡片 */}
          <div className="group relative cursor-pointer">
            <div className="absolute inset-0 bg-sunny-orange/30 rounded-3xl blur-md opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
            <div className="relative flex flex-col items-center gap-3 bg-white border border-gray-100 group-hover:border-sunny-orange/50 p-5 rounded-3xl hover:-translate-y-1 transition-all duration-300 h-full shadow-card group-hover:shadow-card-hover-orange">
              <div className="size-16 rounded-full bg-sunny-orange/10 flex items-center justify-center text-sunny-orange group-hover:scale-110 transition-transform duration-300">
                <span className="material-symbols-outlined text-3xl">auto_stories</span>
              </div>
              <h3 className="text-text-main text-lg font-bold text-center group-hover:text-sunny-orange transition-colors font-display">
                人文素养
              </h3>
            </div>
          </div>

          {/* 语言能力卡片 */}
          <div className="group relative cursor-pointer">
            <div className="absolute inset-0 bg-sky-blue/30 rounded-3xl blur-md opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
            <div className="relative flex flex-col items-center gap-3 bg-white border border-gray-100 group-hover:border-sky-blue/50 p-5 rounded-3xl hover:-translate-y-1 transition-all duration-300 h-full shadow-card group-hover:shadow-card-hover-blue">
              <div className="size-16 rounded-full bg-sky-blue/10 flex items-center justify-center text-sky-blue group-hover:scale-110 transition-transform duration-300">
                <span className="material-symbols-outlined text-3xl">forum</span>
              </div>
              <h3 className="text-text-main text-lg font-bold text-center group-hover:text-sky-blue transition-colors font-display">
                语言能力
              </h3>
            </div>
          </div>
        </div>
      </main>

      {/* Little Star对话气泡 */}
      <LittleStar message="拍一拍，发现有趣的知识吧～" />

      {/* 背景装饰元素 */}
      <div className="fixed inset-0 pointer-events-none -z-10 overflow-hidden">
        <div className="absolute top-20 left-10 w-72 h-72 bg-science-green/10 rounded-full blur-3xl mix-blend-multiply animate-float" style={{ animationDuration: '6s' }}></div>
        <div className="absolute bottom-20 right-10 w-96 h-96 bg-sky-blue/10 rounded-full blur-3xl mix-blend-multiply animate-float" style={{ animationDuration: '8s', animationDelay: '1s' }}></div>
        <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[500px] h-[500px] bg-sunny-orange/5 rounded-full blur-3xl mix-blend-multiply"></div>
      </div>
    </div>
  );
}

