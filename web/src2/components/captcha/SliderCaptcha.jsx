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

import React, { useState, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { Checkbox, Spin } from '@douyinfe/semi-ui';
import { IconTick } from '@douyinfe/semi-icons';
import CaptchaDialog from './CaptchaDialog';
import useCaptcha from './useCaptcha';
import './captcha.css';

/**
 * SliderCaptcha 滑块验证码组件
 * @param {Object} props
 * @param {Function} props.onSuccess - 验证成功回调，接收 token 参数
 * @param {Function} props.onError - 错误回调
 * @param {string} props.triggerText - 触发器文本
 * @param {boolean} props.disabled - 是否禁用
 */
const SliderCaptcha = ({ 
  onSuccess, 
  onError, 
  triggerText,
  disabled = false 
}) => {
  const { t } = useTranslation();
  const [isOpen, setIsOpen] = useState(false);
  const [isVerified, setIsVerified] = useState(false);
  
  const {
    challenge,
    loading,
    error,
    fetchChallenge,
    verifyChallenge,
    reset
  } = useCaptcha();

  // 打开验证弹窗
  const handleOpen = useCallback(async () => {
    if (disabled || isVerified) return;
    
    setIsOpen(true);
    await fetchChallenge();
  }, [disabled, isVerified, fetchChallenge]);

  // 关闭验证弹窗
  const handleClose = useCallback(() => {
    setIsOpen(false);
  }, []);

  // 刷新验证码
  const handleRefresh = useCallback(async () => {
    await fetchChallenge();
  }, [fetchChallenge]);

  // 验证滑块位置
  const handleVerify = useCallback(async (x) => {
    const result = await verifyChallenge(x);
    
    if (result.success) {
      setIsVerified(true);
      setIsOpen(false);
      if (onSuccess) {
        onSuccess(result.token);
      }
      return result;
    } else {
      // 验证失败，获取新的验证码
      await fetchChallenge();
      return result;
    }
  }, [verifyChallenge, fetchChallenge, onSuccess]);

  // 重置状态
  const handleReset = useCallback(() => {
    setIsVerified(false);
    reset();
  }, [reset]);

  return (
    <div className="slider-captcha-container">
      {/* 触发器 - 复选框样式 */}
      <div 
        className={`captcha-trigger ${isVerified ? 'verified' : ''} ${disabled ? 'disabled' : ''}`}
        onClick={handleOpen}
      >
        {isVerified ? (
          <div className="captcha-trigger-verified">
            <IconTick className="captcha-check-icon" />
            <span>{t('验证成功')}</span>
          </div>
        ) : (
          <div className="captcha-trigger-content">
            <Checkbox 
              checked={false} 
              disabled={disabled}
              onChange={() => {}}
            />
            <span className="captcha-trigger-text">
              {triggerText || t('请完成人机验证后继续')}
            </span>
          </div>
        )}
      </div>

      {/* 验证弹窗 */}
      <CaptchaDialog
        open={isOpen}
        challenge={challenge}
        loading={loading}
        error={error}
        onVerify={handleVerify}
        onClose={handleClose}
        onRefresh={handleRefresh}
      />
    </div>
  );
};

export default SliderCaptcha;
