/**
 * 通用按钮组件
 * 支持多种样式，参考设计稿
 */

import React from 'react';

export interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'outline' | 'ghost';
  size?: 'sm' | 'md' | 'lg';
  children: React.ReactNode;
}

export const Button: React.FC<ButtonProps> = ({
  variant = 'primary',
  size = 'md',
  className = '',
  children,
  ...props
}) => {
  const baseClasses = 'font-display font-bold rounded-full transition-all active:scale-95';
  
  const variantClasses = {
    primary: 'bg-primary hover:bg-[#5aff2b] text-[#0a3f00] shadow-glow',
    secondary: 'bg-sky-blue hover:bg-[#60d0ff] text-white shadow-glow-blue',
    outline: 'bg-white border-2 border-gray-100 hover:border-primary text-text-main',
    ghost: 'bg-transparent hover:bg-gray-50 text-text-main',
  };

  const sizeClasses = {
    sm: 'px-4 py-2 text-sm',
    md: 'px-6 py-3 text-base',
    lg: 'px-8 py-4 text-lg',
  };

  return (
    <button
      className={`${baseClasses} ${variantClasses[variant]} ${sizeClasses[size]} ${className}`}
      {...props}
    >
      {children}
    </button>
  );
};

