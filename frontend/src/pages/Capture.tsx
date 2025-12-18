/**
 * 拍照页面组件
 * 基于 stitch_ui/capture_/_scan_interface/ 设计稿
 */

import { useNavigate } from 'react-router-dom';

export default function Capture() {
  const navigate = useNavigate();

  return (
    <div className="font-display antialiased overflow-hidden h-screen w-full bg-cloud-white text-text-main select-none flex flex-col">
      {/* 顶部栏 */}
      <div className="relative z-30 w-full px-6 py-4 flex justify-between items-center">
        <div className="flex items-center gap-2 bg-white px-5 py-2 rounded-full border border-gray-100 shadow-soft">
          <span className="material-symbols-outlined text-warm-yellow text-2xl fill-1">auto_awesome</span>
          <span className="text-sm font-bold tracking-wide text-slate-600">AI Auto-Detect</span>
        </div>
        <button className="size-12 flex items-center justify-center rounded-full bg-white text-slate-400 hover:text-warm-yellow hover:bg-yellow-50 transition-colors border border-gray-100 shadow-soft">
          <span className="material-symbols-outlined">settings</span>
        </button>
      </div>

      {/* 主要内容区域 */}
      <div className="flex-1 flex flex-col items-center justify-center w-full px-4 relative z-10">
        <div className="mb-6 text-center">
          <h2 className="text-2xl md:text-3xl font-extrabold text-slate-800 tracking-tight drop-shadow-sm font-display">
            Align the target inside the frame!
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
          <button className="absolute -right-2 md:-right-20 top-1/2 -translate-y-1/2 z-30 flex flex-col items-center gap-2 group">
            <div className="size-16 flex items-center justify-center bg-white rounded-full border-2 border-gray-100 text-slate-400 hover:text-warm-yellow hover:border-warm-yellow transition-all duration-300 shadow-soft group-hover:scale-110 group-hover:shadow-glow-yellow">
              <span className="material-symbols-outlined text-[36px]">mic</span>
            </div>
            <span className="text-xs font-bold text-slate-500 bg-white px-2 py-1 rounded-md shadow-sm opacity-0 group-hover:opacity-100 transition transform -translate-x-2">
              Voice Mode
            </span>
          </button>
        </div>

        {/* AI识别提示 */}
        <div className="mt-8 flex items-center gap-3 bg-white px-6 py-3 rounded-2xl border border-gray-100 shadow-soft animate-bounce-slow">
          <div className="size-8 rounded-full bg-gradient-to-tr from-yellow-300 to-orange-400 flex items-center justify-center shadow-inner ring-2 ring-white">
            <span className="material-symbols-outlined text-white text-lg fill-1">star</span>
          </div>
          <p className="text-sm font-bold text-slate-600">Little Star is identifying...</p>
        </div>
      </div>

      {/* 底部操作栏 */}
      <div className="relative z-20 w-full h-auto min-h-[140px] flex items-center justify-center px-10 pb-8 pt-4">
        <div className="flex items-center justify-between w-full max-w-4xl gap-8">
          {/* 相册按钮 */}
          <div className="flex-1 flex justify-end">
            <button className="flex flex-col items-center gap-2 group">
              <div className="size-16 rounded-2xl overflow-hidden border-4 border-white shadow-soft group-hover:shadow-md transition-all relative bg-gray-100 group-hover:scale-105 duration-200">
                <div className="w-full h-full bg-gray-200 flex items-center justify-center">
                  <span className="material-symbols-outlined text-gray-400">photo_library</span>
                </div>
              </div>
              <span className="text-xs font-bold text-slate-500 group-hover:text-warm-yellow transition-colors">Album</span>
            </button>
          </div>

          {/* 快门按钮 */}
          <div className="shrink-0 mx-6">
            <button
              onClick={() => navigate('/result')}
              className="relative size-28 rounded-full bg-white border border-gray-100 flex items-center justify-center shadow-button transition-transform cursor-pointer group hover:shadow-lg active:scale-95"
            >
              <div className="absolute inset-1 rounded-full border-[6px] border-warm-yellow opacity-30 group-hover:opacity-100 transition-opacity"></div>
              <div className="size-[84px] rounded-full bg-warm-yellow border-[4px] border-white shadow-inner flex items-center justify-center group-hover:scale-95 transition-all">
                <span className="material-symbols-outlined text-white text-4xl opacity-90">photo_camera</span>
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
              <span className="text-xs font-bold text-slate-500 group-hover:text-slate-600 transition-colors">Return</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}

