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

import React, { useState, useEffect, useCallback } from 'react';
import { Modal, Typography, Card, Button, Spin } from '@douyinfe/semi-ui';
import { IconRefresh } from '@douyinfe/semi-icons';
import { QRCodeSVG } from 'qrcode.react';
import { SiAlipay, SiWechat } from 'react-icons/si';
import { useOrderPolling } from '../../../hooks/topup/useOrderPolling';

const { Text, Title } = Typography;

// 二维码有效期（5分钟）
const QR_EXPIRE_TIME = 5 * 60 * 1000;

/**
 * 二维码支付弹窗组件
 * 
 * @param {Object} props
 * @param {boolean} props.visible - 弹窗是否可见
 * @param {Function} props.onClose - 关闭回调
 * @param {Function} props.onSuccess - 支付成功回调
 * @param {Function} props.onRefresh - 刷新二维码回调
 * @param {string} props.qrCodeUrl - 二维码链接
 * @param {string} props.tradeNo - 订单号
 * @param {number} props.amount - 充值数量
 * @param {string} props.money - 支付金额
 * @param {string} props.paymentMethod - 支付方式 (alipay/wxpay)
 * @param {Function} props.t - 国际化函数
 */
const QRCodePaymentModal = ({
  visible,
  onClose,
  onSuccess,
  onRefresh,
  qrCodeUrl,
  tradeNo,
  amount,
  money,
  paymentMethod,
  t,
}) => {
  const [countdown, setCountdown] = useState(QR_EXPIRE_TIME / 1000);
  const [isExpired, setIsExpired] = useState(false);
  const [refreshing, setRefreshing] = useState(false);

  // 使用订单轮询 Hook
  const { status, isPolling, startPolling, stopPolling } = useOrderPolling(
    tradeNo,
    {
      expireTime: QR_EXPIRE_TIME,
      onSuccess: (rechargedAmount) => {
        onSuccess?.(rechargedAmount);
      },
      onExpire: () => {
        setIsExpired(true);
      },
    }
  );

  // 弹窗打开时开始轮询和倒计时
  useEffect(() => {
    if (visible && tradeNo && qrCodeUrl) {
      setCountdown(QR_EXPIRE_TIME / 1000);
      setIsExpired(false);
      startPolling();
    }
    return () => {
      stopPolling();
    };
  }, [visible, tradeNo, qrCodeUrl, startPolling, stopPolling]);

  // 倒计时逻辑
  useEffect(() => {
    if (!visible || isExpired) return;

    const timer = setInterval(() => {
      setCountdown((prev) => {
        if (prev <= 1) {
          setIsExpired(true);
          return 0;
        }
        return prev - 1;
      });
    }, 1000);

    return () => clearInterval(timer);
  }, [visible, isExpired]);

  // 格式化倒计时显示
  const formatCountdown = useCallback((seconds) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  }, []);

  // 处理关闭
  const handleClose = () => {
    stopPolling();
    onClose?.();
  };

  // 处理刷新
  const handleRefresh = async () => {
    setRefreshing(true);
    stopPolling();
    try {
      await onRefresh?.();
      setCountdown(QR_EXPIRE_TIME / 1000);
      setIsExpired(false);
    } finally {
      setRefreshing(false);
    }
  };

  // 获取支付方式图标和名称
  const getPaymentInfo = () => {
    if (paymentMethod === 'alipay') {
      return {
        icon: <SiAlipay size={20} color='#1677FF' />,
        name: t('支付宝'),
        color: '#1677FF',
      };
    }
    return {
      icon: <SiWechat size={20} color='#07C160' />,
      name: t('微信'),
      color: '#07C160',
    };
  };

  const paymentInfo = getPaymentInfo();

  return (
    <Modal
      title={
        <div className='flex items-center'>
          {paymentInfo.icon}
          <span className='ml-2'>{paymentInfo.name}{t('扫码支付')}</span>
        </div>
      }
      visible={visible}
      onCancel={handleClose}
      footer={null}
      maskClosable={false}
      centered
      width={360}
    >
      <div className='flex flex-col items-center space-y-4 py-2'>
        {/* 订单信息 */}
        <Card className='w-full !rounded-xl !border-0 bg-slate-50 dark:bg-slate-800'>
          <div className='space-y-2'>
            <div className='flex justify-between items-center'>
              <Text className='text-slate-600 dark:text-slate-400'>
                {t('订单号')}：
              </Text>
              <Text className='text-slate-900 dark:text-slate-100 font-mono text-sm'>
                {tradeNo}
              </Text>
            </div>
            <div className='flex justify-between items-center'>
              <Text className='text-slate-600 dark:text-slate-400'>
                {t('充值数量')}：
              </Text>
              <Text className='text-slate-900 dark:text-slate-100'>
                {amount}
              </Text>
            </div>
            <div className='flex justify-between items-center'>
              <Text className='text-slate-600 dark:text-slate-400'>
                {t('支付金额')}：
              </Text>
              <Text strong style={{ color: 'red', fontSize: '16px' }}>
                ¥{money}
              </Text>
            </div>
          </div>
        </Card>

        {/* 二维码区域 */}
        <div className='relative'>
          {isExpired ? (
            <div className='w-[200px] h-[200px] flex flex-col items-center justify-center bg-slate-100 dark:bg-slate-700 rounded-lg'>
              <Text className='text-slate-500 dark:text-slate-400 mb-3'>
                {t('二维码已过期')}
              </Text>
              <Button
                icon={<IconRefresh />}
                onClick={handleRefresh}
                loading={refreshing}
              >
                {t('刷新二维码')}
              </Button>
            </div>
          ) : qrCodeUrl ? (
            <div className='p-3 bg-white rounded-lg shadow-sm'>
              <QRCodeSVG
                value={qrCodeUrl}
                size={180}
                level='M'
                includeMargin={false}
              />
            </div>
          ) : (
            <div className='w-[200px] h-[200px] flex items-center justify-center'>
              <Spin size='large' />
            </div>
          )}
        </div>

        {/* 倒计时和提示 */}
        {!isExpired && (
          <div className='text-center space-y-1'>
            <Text className='text-slate-500 dark:text-slate-400'>
              {t('请使用')}{paymentInfo.name}{t('扫描二维码完成支付')}
            </Text>
            <div className='flex items-center justify-center space-x-1'>
              <Text className='text-slate-400 dark:text-slate-500'>
                {t('二维码有效期')}：
              </Text>
              <Text
                strong
                style={{
                  color: countdown <= 60 ? '#f5222d' : paymentInfo.color,
                }}
              >
                {formatCountdown(countdown)}
              </Text>
            </div>
          </div>
        )}

        {/* 轮询状态提示 */}
        {isPolling && !isExpired && (
          <div className='flex items-center space-x-2 text-slate-400'>
            <Spin size='small' />
            <Text type='tertiary' size='small'>
              {t('等待支付中...')}
            </Text>
          </div>
        )}
      </div>
    </Modal>
  );
};

export default QRCodePaymentModal;
