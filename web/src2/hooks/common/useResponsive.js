/*
Copyright (C) 2025 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/

import { useState, useEffect, useMemo } from 'react';

// 响应式断点定义
export const BREAKPOINTS = {
  mobile: 768,
  tablet: 1200,
};

// 防抖函数
const debounce = (func, wait) => {
  let timeout;
  return function executedFunction(...args) {
    const later = () => {
      clearTimeout(timeout);
      func(...args);
    };
    clearTimeout(timeout);
    timeout = setTimeout(later, wait);
  };
};

/**
 * 响应式 Hook - 检测屏幕尺寸并返回断点信息
 * @returns {Object} 包含 isMobile, isTablet, isDesktop 的对象
 */
export const useResponsive = () => {
  const [windowWidth, setWindowWidth] = useState(
    typeof window !== 'undefined' ? window.innerWidth : 1200
  );

  useEffect(() => {
    // 初始化时设置窗口宽度
    if (typeof window !== 'undefined') {
      setWindowWidth(window.innerWidth);
    }

    // 创建防抖的 resize 处理函数
    const handleResize = debounce(() => {
      setWindowWidth(window.innerWidth);
    }, 200);

    // 添加事件监听
    window.addEventListener('resize', handleResize);

    // 清理函数
    return () => {
      window.removeEventListener('resize', handleResize);
    };
  }, []);

  // 使用 useMemo 缓存计算结果
  const breakpointInfo = useMemo(() => {
    const isMobile = windowWidth < BREAKPOINTS.mobile;
    const isTablet = windowWidth >= BREAKPOINTS.mobile && windowWidth < BREAKPOINTS.tablet;
    const isDesktop = windowWidth >= BREAKPOINTS.tablet;

    return {
      isMobile,
      isTablet,
      isDesktop,
      windowWidth,
      breakpoint: isMobile ? 'mobile' : isTablet ? 'tablet' : 'desktop',
    };
  }, [windowWidth]);

  return breakpointInfo;
};
