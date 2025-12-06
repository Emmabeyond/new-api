import { useState, useEffect, useRef } from 'react';

/**
 * 数字计数动画 Hook
 * @param {Object} options
 * @param {number} options.end - 目标数值
 * @param {number} [options.duration=2000] - 动画持续时间（毫秒）
 * @param {number} [options.start=0] - 起始数值
 * @param {boolean} [options.enabled=true] - 是否启用动画
 * @returns {number} 当前显示的数值
 */
const useCountUp = ({ end, duration = 2000, start = 0, enabled = true }) => {
  const [count, setCount] = useState(start);
  const frameRef = useRef(null);
  const startTimeRef = useRef(null);

  useEffect(() => {
    // 检测用户是否偏好减少动画
    const prefersReducedMotion = window.matchMedia(
      '(prefers-reduced-motion: reduce)'
    ).matches;

    if (!enabled || prefersReducedMotion) {
      setCount(end);
      return;
    }

    const animate = (timestamp) => {
      if (!startTimeRef.current) {
        startTimeRef.current = timestamp;
      }

      const progress = Math.min(
        (timestamp - startTimeRef.current) / duration,
        1
      );

      // 使用 easeOutQuart 缓动函数
      const easeOutQuart = 1 - Math.pow(1 - progress, 4);
      const currentValue = Math.floor(start + (end - start) * easeOutQuart);

      setCount(currentValue);

      if (progress < 1) {
        frameRef.current = requestAnimationFrame(animate);
      }
    };

    frameRef.current = requestAnimationFrame(animate);

    return () => {
      if (frameRef.current) {
        cancelAnimationFrame(frameRef.current);
      }
    };
  }, [end, duration, start, enabled]);

  return count;
};

export default useCountUp;
