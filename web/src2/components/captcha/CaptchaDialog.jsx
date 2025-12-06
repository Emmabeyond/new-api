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

import React, { useState, useCallback, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { Modal, Spin, Button, Toast } from '@douyinfe/semi-ui';
import { IconRefresh, IconClose } from '@douyinfe/semi-icons';
import SliderTrack from './SliderTrack';
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
  const [puzzleX, setPuzzleX] = useState(0);
  const [status, setStatus] = useState('idle'); // idle, dragging, success, error
  const [imageLoaded, setImageLoaded] = useState(false);
  const [verifying, setVerifying] = useState(false);

  // 重置状态
  useEffect(() => {
    if (open) {
      setPuzzleX(0);
      setStatus('idle');
      setImageLoaded(false);
    }
  }, [open, challenge?.session_id]);

  // 处理拖动
  const handleDrag = useCallback((x) => {
    setPuzzleX(x);
    setStatus('dragging');
  }, []);

  // 处理拖动结束
  const handleDragEnd = useCallback(async (x) => {
    if (!challenge || verifying) return;

    setVerifying(true);
    try {
      const result = await onVerify(x);
      
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
  }, [challenge, verifying, onVerify, t]);

  // 处理图片加载完成
  const handleImageLoad = useCallback(() => {
    setImageLoaded(true);
  }, []);

  // 处理刷新
  const handleRefresh = useCallback(() => {
    setPuzzleX(0);
    setStatus('idle');
    setImageLoaded(false);
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
        <div className="captcha-image-container">
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
                onLoad={handleImageLoad}
                style={{ display: imageLoaded ? 'block' : 'none' }}
              />
              
              {/* 拼图块 */}
              {imageLoaded && (
                <img
                  src={challenge.puzzle_image}
                  alt="puzzle piece"
                  className={`captcha-puzzle-piece ${status} ${status === 'idle' ? 'idle' : ''}`}
                  style={{
                    left: puzzleX,
                    top: challenge.puzzle_y
                  }}
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
          maxX={challenge ? challenge.width - 50 : 250}
        />
      </div>
    </Modal>
  );
};

export default CaptchaDialog;
