/*
Copyright (C) 2025 QuantumNous
*/

import React from 'react';
import { Modal, Descriptions, Tag } from '@douyinfe/semi-ui';

const OrderDetailModal = ({ visible, order, onClose, t }) => {
  if (!order) return null;

  // 状态标签颜色映射
  const statusColorMap = {
    pending: 'orange',
    success: 'green',
    expired: 'grey',
  };

  // 状态文本映射
  const statusTextMap = {
    pending: t('待支付'),
    success: t('成功'),
    expired: t('已过期'),
  };

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
    { key: t('订单ID'), value: order.id },
    { key: t('订单号'), value: order.trade_no },
    {
      key: t('状态'),
      value: (
        <Tag color={statusColorMap[order.status] || 'default'}>
          {statusTextMap[order.status] || order.status}
        </Tag>
      ),
    },
    { key: t('用户ID'), value: order.user_id },
    { key: t('用户名'), value: order.username || '-' },
    { key: t('用户当前余额'), value: order.user_quota ? `${order.user_quota}` : '-' },
    { key: t('充值金额'), value: `$${order.amount}` },
    { key: t('支付金额'), value: `¥${order.money?.toFixed(2) || '0.00'}` },
    { key: t('支付方式'), value: order.payment_method },
    { key: t('创建时间'), value: formatTime(order.create_time) },
    { key: t('完成时间'), value: formatTime(order.complete_time) },
  ];

  return (
    <Modal
      title={t('订单详情')}
      visible={visible}
      onCancel={onClose}
      footer={null}
      width={600}
    >
      <Descriptions data={data} />
    </Modal>
  );
};

export default OrderDetailModal;
