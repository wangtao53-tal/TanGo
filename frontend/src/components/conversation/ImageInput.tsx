/**
 * 图片输入组件
 * 支持图片上传
 */

import { useRef } from 'react';
import { useTranslation } from 'react-i18next';

export interface ImageInputProps {
  onImageSelect: (file: File) => void;
  disabled?: boolean;
}

export function ImageInput({ onImageSelect, disabled = false }: ImageInputProps) {
  const { t } = useTranslation();
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleImageClick = () => {
    if (disabled) return;
    fileInputRef.current?.click();
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file && file.type.startsWith('image/')) {
      onImageSelect(file);
    }
    // 清空input，允许重复选择同一文件
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }
  };

  return (
    <>
      <button
        onClick={handleImageClick}
        disabled={disabled}
        className="p-3 rounded-full bg-gray-100 text-gray-600 hover:bg-gray-200 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
        title={t('conversation.imageInput')}
      >
        <span className="material-symbols-outlined text-2xl">image</span>
      </button>
      <input
        ref={fileInputRef}
        type="file"
        accept="image/*"
        className="hidden"
        onChange={handleFileChange}
        disabled={disabled}
      />
    </>
  );
}
