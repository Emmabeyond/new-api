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

import React, { useState, useEffect, useRef, useCallback } from 'react';
import { Modal, Button, Spin, Typography } from '@douyinfe/semi-ui';
import { IconRefresh } from '@douyinfe/semi-icons';
import { useTranslation } from 'react-i18next';
import { QRCodeSVG } from 'qrcode.react';
import { API, showError, showSuccess, setUserData, updateAPI } from '../../helpers';

const { Text, Title } = Typography;

// 轮询间隔（毫秒）
const POLL_INTERVAL = 2000;
// 最大轮询次数（5分钟 / 2秒 = 150次）
const MAX_POLL_COUNT = 150;
// 重试次数
const MAX_RETRY_COUNT = 3;

const WeChatThirdPartyLoginModal = ({
  visible,
  onCancel,
  onSuccess,
  action = 'login', // 'login' 或 'bind'
}) => {
  const { t } = useTranslation();
  const [loading, setLoading] = useState(false);
  const [qrCodeUrl, setQrCodeUrl] = useState('');
  const [sessionId, setSessionId] = useState('');
  const [expiresAt, setExpiresAt] = useState(null);
  const [status, setStatus] = useState('pending'); // pending, scanned, confirmed, expired
  const [countdown, setCountdown] = useState(0);
  const [error, setError] = useState('');
  
  const pollTimerRef = useRef(null);
  const countdownTimerRef = useRef(null);
  const pollCountRef = useRef(0);
  const retryCountRef = useRef(0);


  // 生成二维码
  const generateQRCode = useCallback(async () => {
    setLoading(true);
    setError('');
    setStatus('pending');
    pollCountRef.current = 0;
    retryCountRef.current = 0;

    try {
      const res = await API.post('/api/wechat-third-party/generate', { action });
      const { success, message, data } = res.data;
      
      if (success && data) {
        setQrCodeUrl(data.qrCodeUrl);
        setSessionId(data.sessionId);
        setExpiresAt(new Date(data.expiresAt));
        
        // 计算倒计时
        const expiresTime = new Date(data.expiresAt).getTime();
        const now = Date.now();
        const remainingSeconds = Math.max(0, Math.floor((expiresTime - now) / 1000));
        setCountdown(remainingSeconds);
        
        // 开始轮询
        startPolling(data.sessionId);
        // 开始倒计时
        startCountdown();
      } else {
        setError(message || t('生成二维码失败'));
      }
    } catch (err) {
      setError(t('生成二维码失败，请重试'));
    } finally {
      setLoading(false);
    }
  }, [action, t]);

  // 开始轮询状态
  const startPolling = useCallback((sid) => {
    if (pollTimerRef.current) {
      clearInterval(pollTimerRef.current);
    }

    pollTimerRef.current = setInterval(async () => {
      pollCountRef.current += 1;
      
      // 超过最大轮询次数
      if (pollCountRef.current > MAX_POLL_COUNT) {
        stopPolling();
        setStatus('expired');
        return;
      }

      try {
        const res = await API.get(`/api/wechat-third-party/status/${sid}`);
        const { success, message, data } = res.data;
        
        if (success && data) {
          setStatus(data.status);
          retryCountRef.current = 0; // 重置重试计数
          
          if (data.status === 'confirmed') {
            stopPolling();
            // 登录成功
            if (message === 'bind') {
              showSuccess(t('绑定成功'));
            } else {
              showSuccess(t('登录成功'));
            }
            onSuccess && onSuccess(data);
          } else if (data.status === 'expired') {
            stopPolling();
          }
        } else {
          // 请求失败，重试
          retryCountRef.current += 1;
          if (retryCountRef.current >= MAX_RETRY_COUNT) {
            stopPolling();
            setError(message || t('查询状态失败'));
          }
        }
      } catch (err) {
        // 网络错误，重试
        retryCountRef.current += 1;
        if (retryCountRef.current >= MAX_RETRY_COUNT) {
          stopPolling();
          setError(t('网络错误，请重试'));
        }
      }
    }, POLL_INTERVAL);
  }, [onSuccess, t]);

  // 停止轮询
  const stopPolling = useCallback(() => {
    if (pollTimerRef.current) {
      clearInterval(pollTimerRef.current);
      pollTimerRef.current = null;
    }
  }, []);

  // 开始倒计时
  const startCountdown = useCallback(() => {
    if (countdownTimerRef.current) {
      clearInterval(countdownTimerRef.current);
    }

    countdownTimerRef.current = setInterval(() => {
      setCountdown((prev) => {
        if (prev <= 1) {
          clearInterval(countdownTimerRef.current);
          countdownTimerRef.current = null;
          setStatus('expired');
          stopPolling();
          return 0;
        }
        return prev - 1;
      });
    }, 1000);
  }, [stopPolling]);

  // 停止倒计时
  const stopCountdown = useCallback(() => {
    if (countdownTimerRef.current) {
      clearInterval(countdownTimerRef.current);
      countdownTimerRef.current = null;
    }
  }, []);

  // 刷新二维码
  const handleRefresh = useCallback(() => {
    stopPolling();
    stopCountdown();
    generateQRCode();
  }, [generateQRCode, stopPolling, stopCountdown]);

  // 格式化倒计时
  const formatCountdown = (seconds) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  };

  // 获取状态文本
  const getStatusText = () => {
    switch (status) {
      case 'pending':
        return t('请使用微信扫描二维码');
      case 'scanned':
        return t('扫描成功，请在手机上确认');
      case 'confirmed':
        return t('登录成功');
      case 'expired':
        return t('二维码已过期，请刷新');
      default:
        return '';
    }
  };

  // 获取状态颜色
  const getStatusColor = () => {
    switch (status) {
      case 'scanned':
        return 'var(--semi-color-success)';
      case 'expired':
        return 'var(--semi-color-danger)';
      default:
        return 'var(--semi-color-text-2)';
    }
  };

  // Modal 打开时生成二维码
  useEffect(() => {
    if (visible) {
      generateQRCode();
    } else {
      // Modal 关闭时清理
      stopPolling();
      stopCountdown();
      setQrCodeUrl('');
      setSessionId('');
      setStatus('pending');
      setError('');
    }
  }, [visible, generateQRCode, stopPolling, stopCountdown]);

  // 组件卸载时清理
  useEffect(() => {
    return () => {
      stopPolling();
      stopCountdown();
    };
  }, [stopPolling, stopCountdown]);


  return (
    <Modal
      title={action === 'bind' ? t('绑定微信账号') : t('微信扫码登录')}
      visible={visible}
      onCancel={onCancel}
      footer={null}
      centered
      width={400}
      maskClosable={status !== 'scanned'}
    >
      <div className="flex flex-col items-center py-4">
        {loading ? (
          <div className="flex flex-col items-center justify-center h-64">
            <Spin size="large" />
            <Text className="mt-4">{t('正在生成二维码...')}</Text>
          </div>
        ) : error ? (
          <div className="flex flex-col items-center justify-center h-64">
            <Text type="danger" className="mb-4">{error}</Text>
            <Button
              icon={<IconRefresh />}
              onClick={handleRefresh}
            >
              {t('重试')}
            </Button>
          </div>
        ) : (
          <>
            {/* 二维码区域 */}
            <div className="relative mb-4">
              <div
                className="w-48 h-48 rounded-lg bg-white p-2 flex items-center justify-center"
                style={{
                  opacity: status === 'expired' ? 0.3 : 1,
                  transition: 'opacity 0.3s',
                }}
              >
                {qrCodeUrl && (
                  <QRCodeSVG
                    value={qrCodeUrl}
                    size={176}
                    level="M"
                    includeMargin={false}
                  />
                )}
              </div>
              
              {/* 过期遮罩 */}
              {status === 'expired' && (
                <div className="absolute inset-0 flex flex-col items-center justify-center bg-white/80 dark:bg-gray-800/80 rounded-lg">
                  <Text type="danger" className="mb-2">{t('二维码已过期')}</Text>
                  <Button
                    icon={<IconRefresh />}
                    onClick={handleRefresh}
                    size="small"
                  >
                    {t('刷新')}
                  </Button>
                </div>
              )}
              
              {/* 扫描成功遮罩 */}
              {status === 'scanned' && (
                <div className="absolute inset-0 flex flex-col items-center justify-center bg-green-50/90 dark:bg-green-900/50 rounded-lg">
                  <div className="w-12 h-12 rounded-full bg-green-500 flex items-center justify-center mb-2">
                    <svg className="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                    </svg>
                  </div>
                  <Text className="text-green-600 dark:text-green-400">{t('扫描成功')}</Text>
                </div>
              )}
            </div>

            {/* 状态文本 */}
            <Text
              className="mb-2 text-center"
              style={{ color: getStatusColor() }}
            >
              {getStatusText()}
            </Text>

            {/* 倒计时 */}
            {status !== 'expired' && status !== 'confirmed' && countdown > 0 && (
              <Text type="tertiary" size="small">
                {t('二维码有效期')}: {formatCountdown(countdown)}
              </Text>
            )}

            {/* 刷新按钮 */}
            {status !== 'scanned' && status !== 'confirmed' && (
              <Button
                icon={<IconRefresh />}
                type="tertiary"
                size="small"
                className="mt-4"
                onClick={handleRefresh}
              >
                {t('刷新二维码')}
              </Button>
            )}
          </>
        )}
      </div>
    </Modal>
  );
};

export default WeChatThirdPartyLoginModal;
