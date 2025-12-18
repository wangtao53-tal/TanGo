/**
 * 通用卡片组件
 */

import React from 'react';

export interface CardProps extends React.HTMLAttributes<HTMLDivElement> {
  children: React.ReactNode;
  hover?: boolean;
}

export const Card: React.FC<CardProps> = ({
  children,
  hover = false,
  className = '',
  ...props
}) => {
  const hoverClasses = hover
    ? 'hover:-translate-y-1 transition-all duration-300'
    : '';

  return (
    <div
      className={`bg-white border border-gray-100 rounded-3xl p-5 shadow-card ${hoverClasses} ${className}`}
      {...props}
    >
      {children}
    </div>
  );
};

