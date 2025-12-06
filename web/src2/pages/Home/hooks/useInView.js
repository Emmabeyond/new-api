import { useState, useEffect, useRef } from 'react';

/**
 * 视口检测 Hook
 * @param {Object} options
 * @param {number} [options.threshold=0.1] - 触发阈值（0-1）
 * @param {boolean} [options.triggerOnce=true] - 是否只触发一次
 * @param {string} [options.rootMargin='0px'] - 根元素边距
 * @returns {[React.RefObject, boolean]} [ref, isInView]
 */
const useInView = ({ 
  threshold = 0.1, 
  triggerOnce = true,
  rootMargin = '0px'
} = {}) => {
  const ref = useRef(null);
  const [isInView, setIsInView] = useState(false);

  useEffect(() => {
    const element = ref.current;
    if (!element) return;

    const observer = new IntersectionObserver(
      ([entry]) => {
        const inView = entry.isIntersecting;
        setIsInView(inView);

        if (inView && triggerOnce) {
          observer.disconnect();
        }
      },
      {
        threshold,
        rootMargin,
      }
    );

    observer.observe(element);

    return () => {
      observer.disconnect();
    };
  }, [threshold, triggerOnce, rootMargin]);

  return [ref, isInView];
};

export default useInView;
