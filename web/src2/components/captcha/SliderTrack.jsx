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

import React, { useRef, useState, useCallback, useEffect } from 'react';
import { IconArrowRight, IconTick } from '@douyinfe/semi-icons';
import './captcha.css';

/**
 * SliderTrack 滑块轨道组件
 * @param {Object} props
 * @param {Function} props.onDrag - 拖动回调，返回显示坐标
 * @param {Function} props.onDragEnd - 拖动结束回调，返回显示坐标
 * @param {boolean} props.disabled - 是否禁用
 * @param {string} props.status - 状态: idle, dragging, success, error
 * @param {number} props.maxX - 原始坐标系的最大 X 位置
 * @param {number} props.scaleRatio - 缩放比例
 * @param {number} props.displayedWidth - 显示宽度
 */
const SliderTrack = ({
  onDrag,
  onDragEnd,
  disabled = false,
  status = 'idle',
  maxX = 250,
  scaleRatio = 1,
  displayedWidth = 300
}) => {
  const trackRef = useRef(null);
  const [isDragging, setIsDragging] = useState(false);
  const [position, setPosition] = useState(0);
  const startXRef = useRef(0);
  const startPosRef = useRef(0);

  // 计算显示坐标位置
  const calculateDisplayPosition = useCallback((clientX) => {
    if (!trackRef.current) return 0;
    
    const rect = trackRef.current.getBoundingClientRect();
    const handleWidth = 44; // 滑块宽度
    const trackWidth = rect.width - handleWidth;
    
    let newX = clientX - rect.left - handleWidth / 2;
    newX = Math.max(0, Math.min(newX, trackWidth));
    
    // 映射到显示坐标系（基于缩放后的图片宽度）
    const ratio = newX / trackWidth;
    const displayMaxX = maxX * scaleRatio;
    return Math.round(ratio * displayMaxX);
  }, [maxX, scaleRatio]);

  // 鼠标按下
  const handleMouseDown = useCallback((e) => {
    if (disabled) return;
    
    e.preventDefault();
    setIsDragging(true);
    startXRef.current = e.clientX;
    startPosRef.current = position;
  }, [disabled, position]);

  // 鼠标移动
  const handleMouseMove = useCallback((e) => {
    if (!isDragging || disabled) return;
    
    const newPos = calculateDisplayPosition(e.clientX);
    setPosition(newPos);
    onDrag(newPos);
  }, [isDragging, disabled, calculateDisplayPosition, onDrag]);

  // 鼠标释放
  const handleMouseUp = useCallback((e) => {
    if (!isDragging) return;
    
    setIsDragging(false);
    const finalPos = calculateDisplayPosition(e.clientX);
    onDragEnd(finalPos);
  }, [isDragging, calculateDisplayPosition, onDragEnd]);

  // 触摸开始
  const handleTouchStart = useCallback((e) => {
    if (disabled) return;
    
    const touch = e.touches[0];
    setIsDragging(true);
    startXRef.current = touch.clientX;
    startPosRef.current = position;
  }, [disabled, position]);

  // 触摸移动
  const handleTouchMove = useCallback((e) => {
    if (!isDragging || disabled) return;
    
    const touch = e.touches[0];
    const newPos = calculateDisplayPosition(touch.clientX);
    setPosition(newPos);
    onDrag(newPos);
  }, [isDragging, disabled, calculateDisplayPosition, onDrag]);

  // 触摸结束
  const handleTouchEnd = useCallback((e) => {
    if (!isDragging) return;
    
    setIsDragging(false);
    onDragEnd(position);
  }, [isDragging, position, onDragEnd]);

  // 添加全局事件监听
  useEffect(() => {
    if (isDragging) {
      document.addEventListener('mousemove', handleMouseMove);
      document.addEventListener('mouseup', handleMouseUp);
      document.addEventListener('touchmove', handleTouchMove, { passive: false });
      document.addEventListener('touchend', handleTouchEnd);
    }

    return () => {
      document.removeEventListener('mousemove', handleMouseMove);
      document.removeEventListener('mouseup', handleMouseUp);
      document.removeEventListener('touchmove', handleTouchMove);
      document.removeEventListener('touchend', handleTouchEnd);
    };
  }, [isDragging, handleMouseMove, handleMouseUp, handleTouchMove, handleTouchEnd]);

  // 重置位置
  useEffect(() => {
    if (status === 'idle' && position !== 0) {
      setPosition(0);
    }
  }, [status]);

  // 计算滑块在轨道上的位置百分比
  const getHandlePosition = () => {
    if (!trackRef.current) return 0;
    const trackWidth = trackRef.current.offsetWidth - 44;
    const displayMaxX = maxX * scaleRatio;
    return displayMaxX > 0 ? (position / displayMaxX) * trackWidth : 0;
  };

  return (
    <div 
      ref={trackRef}
      className={`slider-track ${status} ${disabled ? 'disabled' : ''}`}
    >
      {/* 已滑动区域 */}
      <div 
        className="slider-track-fill"
        style={{ width: getHandlePosition() + 22 }}
      />
      
      {/* 滑块手柄 */}
      <div
        className={`slider-handle ${isDragging ? 'dragging' : ''} ${status}`}
        style={{ left: getHandlePosition() }}
        onMouseDown={handleMouseDown}
        onTouchStart={handleTouchStart}
      >
        {status === 'success' ? (
          <IconTick className="slider-handle-icon success" />
        ) : (
          <IconArrowRight className="slider-handle-icon" />
        )}
      </div>
    </div>
  );
};

export default SliderTrack;
