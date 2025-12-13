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

import { useState, useEffect, useCallback, useRef } from 'react';

// 后端生成的原始尺寸常量
export const ORIGINAL_IMAGE_WIDTH = 300;
export const ORIGINAL_IMAGE_HEIGHT = 150;
export const ORIGINAL_PUZZLE_SIZE = 50;

// 最小触摸目标尺寸
export const MIN_TOUCH_TARGET_SIZE = 44;

// Debounce 延迟时间 (ms) - 仅用于 resize 事件，初始化不使用
const RESIZE_DEBOUNCE_DELAY = 100;

/**
 * useResponsiveCaptcha - 管理验证码响应式缩放逻辑的 Hook
 * 
 * @param {React.RefObject<HTMLElement>} containerRef - 图片容器的 ref
 * @returns {Object} 缩放状态和坐标转换函数
 */
const useResponsiveCaptcha = (containerRef) => {
  // 缩放状态
  const [scaleRatio, setScaleRatio] = useState(1);
  const [displayedWidth, setDisplayedWidth] = useState(ORIGINAL_IMAGE_WIDTH);
  const [displayedHeight, setDisplayedHeight] = useState(ORIGINAL_IMAGE_HEIGHT);
  
  // Debounce 定时器 ref
  const debounceTimerRef = useRef(null);

  /**
   * 更新尺寸和缩放比例
   */
  const updateDimensions = useCallback((width) => {
    if (width <= 0) {
      console.warn('[useResponsiveCaptcha] Invalid width:', width);
      return;
    }
    
    const newScaleRatio = width / ORIGINAL_IMAGE_WIDTH;
    const newHeight = width / 2; // 保持 2:1 宽高比
    
    setScaleRatio(newScaleRatio);
    setDisplayedWidth(width);
    setDisplayedHeight(newHeight);
    
    if (process.env.NODE_ENV === 'development') {
      console.log('[useResponsiveCaptcha] Dimensions updated:', {
        containerWidth: width,
        scaleRatio: newScaleRatio,
        displayedWidth: width,
        displayedHeight: newHeight
      });
    }
  }, []);

  /**
   * 原始坐标 → 显示坐标
   * @param {number} original - 原始坐标值 (基于 300px 宽度)
   * @returns {number} 显示坐标值
   */
  const toDisplayCoord = useCallback((original) => {
    if (scaleRatio <= 0) {
      console.warn('[useResponsiveCaptcha] Invalid scaleRatio:', scaleRatio);
      return original;
    }
    return original * scaleRatio;
  }, [scaleRatio]);

  /**
   * 显示坐标 → 原始坐标
   * @param {number} display - 显示坐标值
   * @returns {number} 原始坐标值 (基于 300px 宽度)
   */
  const toOriginalCoord = useCallback((display) => {
    if (scaleRatio <= 0) {
      console.warn('[useResponsiveCaptcha] Invalid scaleRatio:', scaleRatio);
      return display;
    }
    return Math.round(display / scaleRatio);
  }, [scaleRatio]);

  /**
   * 计算拼图块显示尺寸
   */
  const puzzleDisplaySize = Math.max(
    ORIGINAL_PUZZLE_SIZE * scaleRatio,
    MIN_TOUCH_TARGET_SIZE
  );

  // 标记是否已完成初始化
  const initializedRef = useRef(false);

  /**
   * 监听容器尺寸变化
   */
  useEffect(() => {
    const container = containerRef?.current;
    if (!container) return;

    // 初始化尺寸 - 立即执行，不使用 debounce
    const initialWidth = container.offsetWidth;
    if (initialWidth > 0) {
      updateDimensions(initialWidth);
      initializedRef.current = true;
    }

    // 创建 ResizeObserver
    const resizeObserver = new ResizeObserver((entries) => {
      for (const entry of entries) {
        const newWidth = entry.contentRect.width;
        
        // 首次回调也立即执行（确保初始化）
        if (!initializedRef.current) {
          updateDimensions(newWidth);
          initializedRef.current = true;
          return;
        }
        
        // 后续 resize 事件使用 debounce
        if (debounceTimerRef.current) {
          clearTimeout(debounceTimerRef.current);
        }
        
        debounceTimerRef.current = setTimeout(() => {
          updateDimensions(newWidth);
        }, RESIZE_DEBOUNCE_DELAY);
      }
    });

    resizeObserver.observe(container);

    // 清理
    return () => {
      resizeObserver.disconnect();
      initializedRef.current = false;
      if (debounceTimerRef.current) {
        clearTimeout(debounceTimerRef.current);
      }
    };
  }, [containerRef, updateDimensions]);

  return {
    // 状态
    scaleRatio,
    displayedWidth,
    displayedHeight,
    puzzleDisplaySize,
    
    // 坐标转换函数
    toDisplayCoord,
    toOriginalCoord,
    
    // 手动更新尺寸
    updateDimensions
  };
};

export default useResponsiveCaptcha;
