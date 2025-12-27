/*
Copyright (C) 2025 QuantumNous
*/

import React, { useRef, useState, useCallback } from 'react';
import { Star } from 'lucide-react';

/**
 * Card3D - 3D 悬浮卡片组件
 * 支持鼠标跟随倾斜效果、玻璃拟态、动态阴影
 */
const Card3D = ({
  children,
  className = '',
  isFeatured = false,
  glowColor = 'rgba(139, 92, 246, 0.3)', // 默认紫色光晕
  disabled = false, // 禁用 3D 效果（用于 reduced-motion）
}) => {
  const cardRef = useRef(null);
  const [transform, setTransform] = useState({ rotateX: 0, rotateY: 0 });
  const [isHovering, setIsHovering] = useState(false);

  // 检测用户是否偏好减少动画
  const prefersReducedMotion = 
    typeof window !== 'undefined' && 
    window.matchMedia?.('(prefers-reduced-motion: reduce)').matches;

  const shouldAnimate = !disabled && !prefersReducedMotion;

  // 处理鼠标移动，计算 3D 倾斜角度
  const handleMouseMove = useCallback((e) => {
    if (!shouldAnimate || !cardRef.current) return;

    const card = cardRef.current;
    const rect = card.getBoundingClientRect();
    
    // 计算鼠标相对于卡片中心的位置
    const centerX = rect.left + rect.width / 2;
    const centerY = rect.top + rect.height / 2;
    const mouseX = e.clientX - centerX;
    const mouseY = e.clientY - centerY;

    // 计算旋转角度（最大 ±10 度）
    const maxRotate = 10;
    const rotateY = (mouseX / (rect.width / 2)) * maxRotate;
    const rotateX = -(mouseY / (rect.height / 2)) * maxRotate;

    setTransform({ rotateX, rotateY });
  }, [shouldAnimate]);

  const handleMouseEnter = useCallback(() => {
    setIsHovering(true);
  }, []);

  const handleMouseLeave = useCallback(() => {
    setIsHovering(false);
    setTransform({ rotateX: 0, rotateY: 0 });
  }, []);

  // 精选卡片的光晕颜色
  const featuredGlow = 'rgba(251, 191, 36, 0.4)';
  const activeGlow = isFeatured ? featuredGlow : glowColor;

  return (
    <div
      ref={cardRef}
      className={`
        relative rounded-2xl overflow-hidden
        transition-all duration-300 ease-out
        ${shouldAnimate ? 'transform-gpu' : ''}
        ${className}
      `}
      style={{
        transform: shouldAnimate && isHovering
          ? `perspective(1000px) rotateX(${transform.rotateX}deg) rotateY(${transform.rotateY}deg) scale(1.02)`
          : 'perspective(1000px) rotateX(0deg) rotateY(0deg) scale(1)',
        transformStyle: 'preserve-3d',
        willChange: shouldAnimate ? 'transform' : 'auto',
      }}
      onMouseMove={handleMouseMove}
      onMouseEnter={handleMouseEnter}
      onMouseLeave={handleMouseLeave}
    >
      {/* 玻璃拟态背景 */}
      <div
        className={`
          absolute inset-0 rounded-2xl
          backdrop-blur-xl
          ${isFeatured 
            ? 'bg-gradient-to-br from-amber-50/90 to-orange-50/90 dark:from-amber-900/30 dark:to-orange-900/30' 
            : 'bg-white/80 dark:bg-gray-800/80'
          }
          border
          ${isFeatured 
            ? 'border-amber-300/50 dark:border-amber-600/50' 
            : 'border-white/20 dark:border-gray-700/50'
          }
          transition-all duration-300
        `}
        style={{
          boxShadow: isHovering
            ? `0 25px 50px -12px rgba(0, 0, 0, 0.25), 0 0 ${isFeatured ? '40px' : '30px'} ${activeGlow}`
            : isFeatured
              ? '0 10px 30px -10px rgba(0, 0, 0, 0.15), 0 0 20px rgba(251, 191, 36, 0.2)'
              : '0 4px 20px -5px rgba(0, 0, 0, 0.1)',
        }}
      />

      {/* 精选标记 - 高级玻璃拟态风格 */}
      {isFeatured && (
        <div className='absolute -top-2 -right-2 z-20'>
          <div className='relative group'>
            {/* 外层光晕动画 */}
            <div className='absolute -inset-1 rounded-xl bg-gradient-to-r from-amber-400 via-yellow-300 to-orange-400 opacity-75 blur-sm animate-pulse' />
            
            {/* 玻璃拟态容器 */}
            <div className='relative flex items-center gap-1.5 px-3 py-1.5 rounded-xl backdrop-blur-md bg-gradient-to-br from-amber-400/90 via-yellow-400/85 to-orange-400/90 border border-amber-200/50 shadow-[0_4px_20px_rgba(251,191,36,0.4),inset_0_1px_0_rgba(255,255,255,0.3)]'>
              {/* 星星图标 - 带闪烁效果 */}
              <div className='relative'>
                <Star size={13} fill='currentColor' className='text-white drop-shadow-[0_0_3px_rgba(255,255,255,0.8)]' />
                {/* 星星闪光点 */}
                <div className='absolute -top-0.5 -right-0.5 w-1.5 h-1.5 bg-white rounded-full animate-ping opacity-75' />
              </div>
              
              {/* 文字 */}
              <span className='text-white text-xs font-semibold tracking-wide drop-shadow-[0_1px_2px_rgba(0,0,0,0.2)]'>
                精选
              </span>
            </div>
            
            {/* 底部微光反射 */}
            <div className='absolute -bottom-1 left-1/2 -translate-x-1/2 w-3/4 h-1 bg-gradient-to-r from-transparent via-amber-300/40 to-transparent blur-sm' />
          </div>
        </div>
      )}

      {/* 内容区域 */}
      <div className='relative z-10 p-5'>
        {children}
      </div>

      {/* 底部渐变装饰线 */}
      <div 
        className={`
          absolute bottom-0 left-0 right-0 h-1
          ${isFeatured 
            ? 'bg-gradient-to-r from-amber-400 via-orange-400 to-amber-400' 
            : 'bg-gradient-to-r from-purple-400 via-pink-400 to-purple-400'
          }
          opacity-0 transition-opacity duration-300
          ${isHovering ? 'opacity-100' : ''}
        `}
      />
    </div>
  );
};

export default Card3D;
