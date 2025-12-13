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

import React, { useState, useCallback, useEffect, useRef } from 'react';
import { useTranslation } from 'react-i18next';
import { Modal, Spin, Button, Toast } from '@douyinfe/semi-ui';
import { IconRefresh, IconClose } from '@douyinfe/semi-icons';
import SliderTrack from './SliderTrack';
import useResponsiveCaptcha, { ORIGINAL_PUZZLE_SIZE } from './useResponsiveCaptcha';
import './captcha.css';

/**
 * CaptchaDialog 验证码弹窗组件
 * @param {Object} props
 * @param {boolean} props.open - 是否打开
 * @param {Object} props.challenge - 验证挑战数据
 * @param {boolean} props.loading - 是否加载中
 * @param {string} props.error - 错误信息
 * @param {Function} props.onVerify - 验证回调
 * @param {Function} props.onClose - 关闭回调
 * @param {Function} props.onRefresh - 刷新回调
 */
const CaptchaDialog = ({
  open,
  challenge,
  loading,
  error,
  onVerify,
  onClose,
  onRefresh
}) => {
  const { t } = useTranslation();
  const [puzzleX, setPuzzleX] = useState(0); // 显示坐标
  const [status, setStatus] = useState('idle'); // idle, dragging, success, error
  const [bgImageLoaded, setBgImageLoaded] = useState(false);
  const [puzzleImageLoaded, setPuzzleImageLoaded] = useState(false);
  const [verifying, setVerifying] = useState(false);
  
  // 图片容器 ref
  const imageContainerRef = useRef(null);
  
  // 响应式缩放 Hook
  const {
    scaleRatio,
    displayedWidth,
    displayedHeight,
    puzzleDisplaySize,
    toDisplayCoord,
    toOriginalCoord
  } = useResponsiveCaptcha(imageContainerRef);
  
  // 计算是否两张图片都已加载
  const imageLoaded = bgImageLoaded && puzzleImageLoaded;

  // 重置状态
  useEffect(() => {
    if (open) {
      setPuzzleX(0);
      setStatus('idle');
      setBgImageLoaded(false);
      setPuzzleImageLoaded(false);
    }
  }, [open, challenge?.session_id]);

  // 处理拖动 - x 是显示坐标
  const handleDrag = useCallback((displayX) => {
    setPuzzleX(displayX);
    setStatus('dragging');
  }, []);

  // 处理拖动结束 - 转换为原始坐标后发送到后端
  const handleDragEnd = useCallback(async (displayX) => {
    if (!challenge || verifying) return;

    // 获取容器实际宽度，确保坐标转换正确
    // 这是一个安全措施，防止 scaleRatio 还未更新的情况
    let actualScaleRatio = scaleRatio;
    if (imageContainerRef.current) {
      const containerWidth = imageContainerRef.current.offsetWidth;
      if (containerWidth > 0) {
        actualScaleRatio = containerWidth / 300; // ORIGINAL_IMAGE_WIDTH = 300
      }
    }

    // 将显示坐标转换为原始坐标
    const originalX = actualScaleRatio > 0 
      ? Math.round(displayX / actualScaleRatio)
      : toOriginalCoord(displayX);
    
    if (process.env.NODE_ENV === 'development') {
      console.log('[CaptchaDialog] Verify coordinates:', {
        displayX,
        originalX,
        scaleRatio,
        actualScaleRatio
      });
    }

    setVerifying(true);
    try {
      const result = await onVerify(originalX);
      
      if (result.success) {
        setStatus('success');
        Toast.success(t('验证成功'));
      } else {
        setStatus('error');
        Toast.error(result.message || t('验证失败，请重试'));
        // 重置滑块位置
        setTimeout(() => {
          setPuzzleX(0);
          setStatus('idle');
        }, 500);
      }
    } catch (err) {
      setStatus('error');
      Toast.error(t('验证失败，请重试'));
      setTimeout(() => {
        setPuzzleX(0);
        setStatus('idle');
      }, 500);
    } finally {
      setVerifying(false);
    }
  }, [challenge, verifying, onVerify, toOriginalCoord, scaleRatio, t]);

  // 处理背景图加载完成
  const handleBgImageLoad = useCallback(() => {
    setBgImageLoaded(true);
  }, []);

  // 处理拼图块加载完成
  const handlePuzzleImageLoad = useCallback(() => {
    setPuzzleImageLoaded(true);
  }, []);

  // 处理刷新
  const handleRefresh = useCallback(() => {
    setPuzzleX(0);
    setStatus('idle');
    setBgImageLoaded(false);
    setPuzzleImageLoaded(false);
    onRefresh();
  }, [onRefresh]);

  return (
    <Modal
      title={null}
      visible={open}
      onCancel={onClose}
      footer={null}
      centered
      className="captcha-dialog"
      width={340}
      closeIcon={<IconClose />}
    >
      <div className="captcha-dialog-content">
        {/* 标题栏 */}
        <div className="captcha-dialog-header">
          <span className="captcha-dialog-title">{t('安全验证')}</span>
          <Button
            theme="borderless"
            type="tertiary"
            icon={<IconRefresh />}
            onClick={handleRefresh}
            disabled={loading || verifying}
            className="captcha-refresh-btn"
          />
        </div>

        {/* 图片区域 */}
        <div 
          ref={imageContainerRef}
          className="captcha-image-container"
          style={{
            width: '100%',
            height: displayedHeight,
            maxWidth: 300
          }}
        >
          {(loading || !imageLoaded) && (
            <div className="captcha-loading-skeleton">
              <Spin />
            </div>
          )}
          
          {error && (
            <div className="captcha-error">
              <span>{error}</span>
              <Button size="small" onClick={handleRefresh}>
                {t('重试')}
              </Button>
            </div>
          )}

          {challenge && !error && (
            <>
              {/* 背景图（含缺口） */}
              <img
                src={challenge.bg_image}
                alt="captcha background"
                className="captcha-bg-image"
                onLoad={handleBgImageLoad}
                style={{ 
                  display: bgImageLoaded ? 'block' : 'none',
                  width: '100%',
                  height: '100%'
                }}
              />
              
              {/* 拼图块 - 使用缩放后的坐标和尺寸 */}
              {imageLoaded && (
                <img
                  src={challenge.puzzle_image}
                  alt="puzzle piece"
                  className={`captcha-puzzle-piece ${status} ${status === 'idle' ? 'idle' : ''}`}
                  onLoad={handlePuzzleImageLoad}
                  style={{
                    left: puzzleX,
                    top: toDisplayCoord(challenge.puzzle_y),
                    width: puzzleDisplaySize,
                    height: puzzleDisplaySize
                  }}
                />
              )}
              
              {/* 隐藏的拼图块用于预加载 */}
              {!puzzleImageLoaded && (
                <img
                  src={challenge.puzzle_image}
                  alt=""
                  onLoad={handlePuzzleImageLoad}
                  style={{ display: 'none' }}
                />
              )}
            </>
          )}
        </div>

        {/* 提示文字 */}
        <div className="captcha-hint">
          {status === 'success' ? (
            <span className="captcha-hint-success">{t('验证成功')}</span>
          ) : status === 'error' ? (
            <span className="captcha-hint-error">{t('验证失败，请重试')}</span>
          ) : (
            <span>{t('向右拖动滑块完成验证')}</span>
          )}
        </div>

        {/* 滑块轨道 */}
        <SliderTrack
          onDrag={handleDrag}
          onDragEnd={handleDragEnd}
          disabled={loading || !imageLoaded || verifying || status === 'success'}
          status={status}
          maxX={challenge ? challenge.width - ORIGINAL_PUZZLE_SIZE : 250}
          scaleRatio={scaleRatio}
          displayedWidth={displayedWidth}
        />
      </div>
    </Modal>
  );
};

export default CaptchaDialog;
