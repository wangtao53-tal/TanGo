/**
 * 快速拍照按钮组件
 * 固定在页面底部，所有页面可见
 */

import { useNavigate } from 'react-router-dom';

export function QuickCaptureButton() {
  const navigate = useNavigate();

  const handleClick = () => {
    navigate('/capture');
  };

  return (
    <button
      onClick={handleClick}
      className="fixed bottom-6 right-6 z-50 flex items-center justify-center w-16 h-16 rounded-full bg-[var(--color-primary)] hover:bg-[#5aff2b] text-white shadow-lg hover:shadow-xl transition-all duration-200 active:scale-95 border-2 border-white/30"
      aria-label="快速拍照"
    >
      <span className="material-symbols-outlined text-3xl">
        photo_camera
      </span>
    </button>
  );
}
