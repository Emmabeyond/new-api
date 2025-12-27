/*
Copyright (C) 2025 QuantumNous
*/

import React, { useMemo } from 'react';

/**
 * MasonryLayout - 响应式瀑布流布局组件
 * 支持 3/2/1 列响应式布局和入场动画
 */
const MasonryLayout = ({
  items = [],
  renderItem,
  gap = 16,
  className = '',
}) => {
  // 检测用户是否偏好减少动画
  const prefersReducedMotion = 
    typeof window !== 'undefined' && 
    window.matchMedia?.('(prefers-reduced-motion: reduce)').matches;

  // 为每个项目生成动画延迟
  const itemsWithDelay = useMemo(() => {
    return items.map((item, index) => ({
      ...item,
      _animationDelay: prefersReducedMotion ? 0 : index * 80, // 每个卡片延迟 80ms
    }));
  }, [items, prefersReducedMotion]);

  if (items.length === 0) {
    return null;
  }

  return (
    <div
      className={`
        grid gap-4
        grid-cols-1 
        md:grid-cols-2 
        lg:grid-cols-3
        ${className}
      `}
      style={{ gap: `${gap}px` }}
    >
      {itemsWithDelay.map((item, index) => (
        <div
          key={item.id || index}
          className={`
            ${prefersReducedMotion ? '' : 'animate-fade-in-up'}
          `}
          style={{
            animationDelay: `${item._animationDelay}ms`,
            animationFillMode: 'both',
          }}
        >
          {renderItem(item, index)}
        </div>
      ))}
    </div>
  );
};

export default MasonryLayout;
