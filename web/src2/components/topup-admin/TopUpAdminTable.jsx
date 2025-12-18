/*
Copyright (C) 2025 QuantumNous
*/

import React from 'react';
import { Table, Tag, Button, Space } from '@douyinfe/semi-ui';
import { IconEyeOpened, IconCheckCircleStroked } from '@douyinfe/semi-icons';

const TopUpAdminTable = ({ data, loading, onViewDetail, onComplete, t }) => {
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

  const columns = [
    {
      title: t('订单ID'),
      dataIndex: 'id',
      width: 80,
    },
    {
      title: t('用户ID'),
      dataIndex: 'user_id',
      width: 80,
    },
    {
      title: t('充值金额'),
      dataIndex: 'amount',
      width: 100,
      render: (amount) => `$${amount}`,
    },
    {
      title: t('支付金额'),
      dataIndex: 'money',
      width: 100,
      render: (money) => `¥${money?.toFixed(2) || '0.00'}`,
    },
    {
      title: t('订单号'),
      dataIndex: 'trade_no',
      width: 200,
      ellipsis: true,
    },
    {
      title: t('支付方式'),
      dataIndex: 'payment_method',
      width: 100,
    },
    {
      title: t('创建时间'),
      dataIndex: 'create_time',
      width: 170,
      render: formatTime,
    },
    {
      title: t('完成时间'),
      dataIndex: 'complete_time',
      width: 170,
      render: formatTime,
    },
    {
      title: t('状态'),
      dataIndex: 'status',
      width: 90,
      render: (status) => (
        <Tag color={statusColorMap[status] || 'default'}>
          {statusTextMap[status] || status}
        </Tag>
      ),
    },
    {
      title: t('操作'),
      dataIndex: 'action',
      width: 150,
      fixed: 'right',
      render: (_, record) => (
        <Space>
          <Button
            icon={<IconEyeOpened />}
            size='small'
            onClick={() => onViewDetail(record)}
          >
            {t('详情')}
          </Button>
          {record.status === 'pending' && (
            <Button
              icon={<IconCheckCircleStroked />}
              size='small'
              type='warning'
              onClick={() => onComplete(record)}
            >
              {t('补单')}
            </Button>
          )}
        </Space>
      ),
    },
  ];

  return (
    <Table
      columns={columns}
      dataSource={data}
      loading={loading}
      rowKey='id'
      pagination={false}
      scroll={{ x: 1200 }}
      size='small'
    />
  );
};

export default TopUpAdminTable;
