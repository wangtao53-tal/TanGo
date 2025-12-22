/**
 * 拍照页面组件
 * 基于 stitch_ui/capture_/_scan_interface/ 设计稿
 */

import { useNavigate } from 'react-router-dom';
import { useState, useRef } from 'react';
import { useTranslation } from 'react-i18next';
import { identifyImage, uploadImage } from '../services/api';
import { fileToBase64, extractBase64Data, compressImage } from '../utils/image';
import { getUserAgeFromStorage } from '../utils/age';
import type { IdentifyResponse } from '../types/api';

export default function Capture() {
  const navigate = useNavigate();
  const { t } = useTranslation();
  const [isProcessing, setIsProcessing] = useState(false);
  const [isListening, setIsListening] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleCaptureClick = () => {
    fileInputRef.current?.click();
  };

  const handleImageSelect = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;

    setIsProcessing(true);
    try {
      // 获取用户年龄（从存储中获取，优先从年级转换）
      const age = getUserAgeFromStorage();

      // 1. 压缩图片
      const compressedBlob = await compressImage(file, 1920, 1920, 0.8);
      
      // 2. 转换为 base64
      const compressedFile = new File([compressedBlob], file.name, { type: 'image/jpeg' });
      const base64 = await fileToBase64(compressedFile);
      const imageData = extractBase64Data(base64);

      // 3. 上传图片到 GitHub（如果失败会自动降级到 base64）
      let imageUrl = imageData; // 默认使用 base64
      try {
        const uploadResult = await uploadImage({
          imageData: imageData,
          filename: file.name,
        });
        imageUrl = uploadResult.url;
        console.log('图片上传成功:', uploadResult.url, '方式:', uploadResult.uploadMethod);
      } catch (uploadError: any) {
        console.warn('图片上传失败，使用 base64:', uploadError);
        // 上传失败时使用 base64，继续流程
      }

      // 4. 调用识别API（使用 URL 或 base64）
      const identifyResult: IdentifyResponse = await identifyImage({
        image: imageUrl, // 使用上传后的 URL 或 base64
        age,
      });

      // 跳转到问答页面，只传递识别结果（不生成卡片）
      // 使用sessionStorage标记从Capture页面跳转，刷新页面时sessionStorage会清空
      sessionStorage.setItem('fromCapturePage', 'true');
      
      navigate('/result', {
        state: {
          objectName: identifyResult.objectName,
          objectCategory: identifyResult.objectCategory,
          confidence: identifyResult.confidence,
          keywords: identifyResult.keywords,
          age,
          imageData: base64, // 保存原始base64用于显示
        },
      });
    } catch (error: any) {
      console.error('处理图片失败:', error);
      const errorMessage = error?.message || error?.detail || t('capture.identifyError');
      // 友好的错误提示
      alert(t('capture.identifyErrorDetail', { error: errorMessage }));
    } finally {
      setIsProcessing(false);
      // 清空input，允许重复选择同一文件
      if (fileInputRef.current) {
        fileInputRef.current.value = '';
      }
    }
  };

  const handleVoiceInput = () => {
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

      recognition.onresult = async (event: any) => {
        const transcript = event.results[0][0].transcript;
        console.log('语音识别结果:', transcript);
        setIsListening(false);
        
        // TODO: 将语音识别结果发送到后端进行意图识别和处理
        // 这里可以导航到对话页面或直接处理
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
      alert(t('capture.voiceNotSupported'));
    }
  };

  return (
    <div className="font-display antialiased overflow-hidden h-screen w-full bg-cloud-white text-text-main select-none flex flex-col">
      {/* 顶部栏 */}
      <div className="relative z-30 w-full px-6 py-4 flex justify-between items-center">
        <div className="flex items-center gap-2 bg-white px-5 py-2 rounded-full border border-gray-100 shadow-soft">
          <span className="material-symbols-outlined text-warm-yellow text-2xl fill-1">auto_awesome</span>
          <span className="text-sm font-bold tracking-wide text-slate-600">{t('capture.aiAutoDetect')}</span>
        </div>
        <button className="size-12 flex items-center justify-center rounded-full bg-white text-slate-400 hover:text-warm-yellow hover:bg-yellow-50 transition-colors border border-gray-100 shadow-soft">
          <span className="material-symbols-outlined">settings</span>
        </button>
      </div>

      {/* 主要内容区域 */}
      <div className="flex-1 flex flex-col items-center justify-center w-full px-4 relative z-10">
        <div className="mb-6 text-center">
          <h2 className="text-2xl md:text-3xl font-extrabold text-slate-800 tracking-tight drop-shadow-sm font-display">
            {t('capture.title')}
          </h2>
        </div>

        {/* 相机取景框 */}
        <div className="relative w-full max-w-3xl aspect-[4/3] flex items-center justify-center">
          <div className="relative w-full h-full border-[8px] border-warm-yellow rounded-[2.5rem] shadow-glow-yellow overflow-hidden bg-slate-100 z-20 group">
            {/* 相机预览占位符 */}
            <div className="w-full h-full bg-gradient-to-br from-gray-200 to-gray-300 flex items-center justify-center">
              <span className="text-6xl text-gray-400">📷</span>
            </div>
            
            {/* 取景框装饰 */}
            <div className="absolute top-6 left-6 w-10 h-10 border-t-[6px] border-l-[6px] border-white/90 rounded-tl-2xl shadow-sm"></div>
            <div className="absolute top-6 right-6 w-10 h-10 border-t-[6px] border-r-[6px] border-white/90 rounded-tr-2xl shadow-sm"></div>
            <div className="absolute bottom-6 left-6 w-10 h-10 border-b-[6px] border-l-[6px] border-white/90 rounded-bl-2xl shadow-sm"></div>
            <div className="absolute bottom-6 right-6 w-10 h-10 border-b-[6px] border-r-[6px] border-white/90 rounded-br-2xl shadow-sm"></div>
            <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-16 h-16 opacity-60">
              <div className="absolute top-1/2 left-0 w-full h-[2px] bg-white shadow-sm"></div>
              <div className="absolute left-1/2 top-0 h-full w-[2px] bg-white shadow-sm"></div>
            </div>
            
            {/* 扫描线动画 */}
            <div className="absolute w-full h-1 bg-warm-yellow/90 shadow-[0_0_20px_rgba(255,215,0,0.8)] animate-scan"></div>
          </div>

          {/* 语音模式按钮 */}
          <button
            onClick={handleVoiceInput}
            className="absolute -right-2 md:-right-20 top-1/2 -translate-y-1/2 z-30 flex flex-col items-center gap-2 group"
          >
            <div className={`size-16 flex items-center justify-center bg-white rounded-full border-2 border-gray-100 text-slate-400 hover:text-warm-yellow hover:border-warm-yellow transition-all duration-300 shadow-soft group-hover:scale-110 group-hover:shadow-glow-yellow ${isListening ? 'animate-pulse border-warm-yellow text-warm-yellow' : ''}`}>
              <span className="material-symbols-outlined text-[36px]">mic</span>
            </div>
            <span className="text-xs font-bold text-slate-500 bg-white px-2 py-1 rounded-md shadow-sm opacity-0 group-hover:opacity-100 transition transform -translate-x-2">
              {t('capture.voiceInput')}
            </span>
          </button>
        </div>

        {/* AI识别提示 */}
        {isProcessing && (
          <div className="mt-8 flex items-center gap-3 bg-white px-6 py-3 rounded-2xl border border-gray-100 shadow-soft animate-bounce-slow">
            <div className="size-8 rounded-full bg-gradient-to-tr from-yellow-300 to-orange-400 flex items-center justify-center shadow-inner ring-2 ring-white">
              <span className="material-symbols-outlined text-white text-lg fill-1">star</span>
            </div>
            <p className="text-sm font-bold text-slate-600">{t('capture.processing')}</p>
          </div>
        )}
      </div>

      {/* 底部操作栏 */}
      <div className="relative z-20 w-full h-auto min-h-[140px] flex items-center justify-center px-10 pb-8 pt-4">
        <div className="flex items-center justify-between w-full max-w-4xl gap-8">
          {/* 相册按钮 */}
          <div className="flex-1 flex justify-end">
            <button
              onClick={() => fileInputRef.current?.click()}
              className="flex flex-col items-center gap-2 group"
            >
              <div className="size-16 rounded-2xl overflow-hidden border-4 border-white shadow-soft group-hover:shadow-md transition-all relative bg-gray-100 group-hover:scale-105 duration-200">
                <div className="w-full h-full bg-gray-200 flex items-center justify-center">
                  <span className="material-symbols-outlined text-gray-400">photo_library</span>
                </div>
              </div>
              <span className="text-xs font-bold text-slate-500 group-hover:text-warm-yellow transition-colors">{t('capture.selectFromAlbum')}</span>
            </button>
          </div>

          {/* 快门按钮 */}
          <div className="shrink-0 mx-6">
            <input
              ref={fileInputRef}
              type="file"
              accept="image/*"
              capture="environment"
              className="hidden"
              onChange={handleImageSelect}
            />
            <button
              onClick={handleCaptureClick}
              disabled={isProcessing}
              className="relative size-28 rounded-full bg-white border border-gray-100 flex items-center justify-center shadow-button transition-transform cursor-pointer group hover:shadow-lg active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <div className="absolute inset-1 rounded-full border-[6px] border-warm-yellow opacity-30 group-hover:opacity-100 transition-opacity"></div>
              <div className="size-[84px] rounded-full bg-warm-yellow border-[4px] border-white shadow-inner flex items-center justify-center group-hover:scale-95 transition-all">
                {isProcessing ? (
                  <span className="material-symbols-outlined text-white text-4xl opacity-90 animate-spin">refresh</span>
                ) : (
                  <span className="material-symbols-outlined text-white text-4xl opacity-90">photo_camera</span>
                )}
              </div>
            </button>
          </div>

          {/* 返回按钮 */}
          <div className="flex-1 flex justify-start">
            <button
              onClick={() => navigate('/')}
              className="flex flex-col items-center gap-2 group"
            >
              <div className="size-16 flex items-center justify-center rounded-full bg-white border-4 border-white shadow-soft group-hover:shadow-md transition-all group-hover:scale-105 duration-200">
                <span className="material-symbols-outlined text-slate-400 text-3xl group-hover:text-slate-600 transition-colors">arrow_back</span>
              </div>
              <span className="text-xs font-bold text-slate-500 group-hover:text-slate-600 transition-colors">{t('common.back')}</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}

