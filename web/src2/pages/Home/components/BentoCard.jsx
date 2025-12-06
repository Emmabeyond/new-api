import React, { useRef, useEffect, useState } from 'react';
import { useActualTheme } from '../../../context/Theme';

const CARD_SIZES = {
  small: 'col-span-1 row-span-1',
  medium: 'col-span-1 row-span-1',
  large: 'col-span-1 md:col-span-2 row-span-1 md:row-span-2',
  wide: 'col-span-1 md:col-span-2 row-span-1',
};

const BentoCard = ({ 
  size = 'medium', 
  children, 
  className = '', 
  delay = 0,
  onClick,
  floating = false,
}) => {
  const actualTheme = useActualTheme();
  const isDark = actualTheme === 'dark';
  const cardRef = useRef(null);
  const [isVisible, setIsVisible] = useState(false);
  const [mousePosition, setMousePosition] = useState({ x: 0, y: 0 });

  useEffect(() => {
    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          setIsVisible(true);
          observer.disconnect();
        }
      },
      { threshold: 0.1 }
    );

    if (cardRef.current) {
      observer.observe(cardRef.current);
    }

    return () => observer.disconnect();
  }, []);

  // 鼠标跟随光效
  const handleMouseMove = (e) => {
    if (!cardRef.current) return;
    const rect = cardRef.current.getBoundingClientRect();
    setMousePosition({
      x: e.clientX - rect.left,
      y: e.clientY - rect.top,
    });
  };

  const sizeClass = CARD_SIZES[size] || CARD_SIZES.medium;

  const baseStyles = `
    relative overflow-hidden rounded-2xl p-5 md:p-6
    transition-all duration-500 ease-out
    hover:scale-[1.02] hover:-translate-y-1
    ${sizeClass}
    ${floating ? 'float-card' : ''}
  `;

  const darkStyles = `
    bg-white/[0.03] backdrop-blur-xl
    border border-white/[0.08]
    hover:border-white/[0.2]
    hover:bg-white/[0.06]
    bento-card-glow
    hover:shadow-[0_8px_32px_rgba(99,102,241,0.15)]
  `;

  const lightStyles = `
    bg-white/90 backdrop-blur-sm border border-gray-200/80
    shadow-sm hover:shadow-xl
    hover:border-indigo-200
  `;

  const animationStyles = isVisible
    ? 'opacity-100 translate-y-0 scale-100'
    : 'opacity-0 translate-y-8 scale-95';

  return (
    <div
      ref={cardRef}
      className={`
        ${baseStyles}
        ${isDark ? darkStyles : lightStyles}
        ${animationStyles}
        ${className}
      `}
      style={{ 
        transitionDelay: `${delay}ms`,
        cursor: onClick ? 'pointer' : 'default'
      }}
      onClick={onClick}
      onMouseMove={handleMouseMove}
    >
      {/* 鼠标跟随光效 */}
      {isDark && (
        <div
          className="absolute pointer-events-none transition-opacity duration-300"
          style={{
            width: '200px',
            height: '200px',
            left: mousePosition.x - 100,
            top: mousePosition.y - 100,
            background: 'radial-gradient(circle, rgba(99, 102, 241, 0.15) 0%, transparent 70%)',
            opacity: isVisible ? 1 : 0,
          }}
        />
      )}
      
      {isDark && <div className="bento-card-glow-effect" />}
      
      {/* 角落装饰 */}
      <div className="absolute top-0 right-0 w-20 h-20 opacity-30 pointer-events-none">
        <div className={`absolute top-2 right-2 w-8 h-[1px] ${isDark ? 'bg-white/20' : 'bg-gray-300'}`} />
        <div className={`absolute top-2 right-2 w-[1px] h-8 ${isDark ? 'bg-white/20' : 'bg-gray-300'}`} />
      </div>
      
      <div className="relative z-10 h-full">
        {children}
      </div>
    </div>
  );
};

export default BentoCard;
