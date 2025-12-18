/*
Copyright (C) 2025 QuantumNous
*/

import React from 'react';
import { Modal, Descriptions, Typography } from '@douyinfe/semi-ui';

const { Text } = Typography;

const CompleteOrderModal = ({ visible, order, onConfirm, onCancel, loading, t }) => {
  if (!order) return null;

  // 格式化时间
  const formatTime = (timestamp) => {
    if (!timestamp) return '-';
    const date = new Date(timestamp * 1000);
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
    });
  };

  const data = [
    { key: t('订单号'), value: order.trade_no },
    { key: t('用户ID'), value: order.user_id },
    { key: t('充值金额'), value: `$${order.amount}` },
    { key: t('支付金额'), value: `¥${order.money?.toFixed(2) || '0.00'}` },
    { key: t('支付方式'), value: order.payment_method },
    { key: t('创建时间'), value: formatTime(order.create_time) },
  ];

  return (
    <Modal
      title={t('确认补单')}
      visible={visible}
      onOk={onConfirm}
      onCancel={onCancel}
      okText={t('确认补单')}
      cancelText={t('取消')}
      confirmLoading={loading}
      width={500}
    >
      <div className='mb-4'>
        <Text type='warning'>
          {t('确认要为该订单进行补单操作吗？补单后将为用户增加相应额度。')}
        </Text>
      </div>
      <Descriptions data={data} />
    </Modal>
  );
};

export default CompleteOrderModal;
