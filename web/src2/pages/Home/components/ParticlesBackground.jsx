import React, { useEffect, useRef, useMemo } from 'react';
import { useActualTheme } from '../../../context/Theme';

const ParticlesBackground = ({ count = 30 }) => {
  const actualTheme = useActualTheme();
  const isDark = actualTheme === 'dark';
  const containerRef = useRef(null);

  const particles = useMemo(() => {
    return Array.from({ length: count }, (_, i) => ({
      id: i,
      left: `${Math.random() * 100}%`,
      top: `${Math.random() * 100}%`,
      size: Math.random() * 3 + 2,
      delay: Math.random() * 10,
      duration: Math.random() * 10 + 15,
      opacity: Math.random() * 0.3 + 0.1,
    }));
  }, [count]);

  return (
    <div ref={containerRef} className="particles-container">
      {/* 网格背景 */}
      <div className="grid-background" />
      
      {/* 光晕球 */}
      <div className="glow-orb glow-orb-primary" />
      <div className="glow-orb glow-orb-secondary" />
      <div className="glow-orb glow-orb-accent" />
      
      {/* 粒子 */}
      {particles.map((particle) => (
        <div
          key={particle.id}
          className="particle"
          style={{
            left: particle.left,
            top: particle.top,
            width: `${particle.size}px`,
            height: `${particle.size}px`,
            animationDelay: `${particle.delay}s`,
            animationDuration: `${particle.duration}s`,
            opacity: particle.opacity,
            color: isDark ? 'rgba(139, 92, 246, 0.6)' : 'rgba(99, 102, 241, 0.4)',
          }}
        />
      ))}
    </div>
  );
};

export default ParticlesBackground;
