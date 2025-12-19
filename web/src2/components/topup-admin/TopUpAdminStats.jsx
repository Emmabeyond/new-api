/*
Copyright (C) 2025 QuantumNous
*/

import React from 'react';
import { Card, Skeleton } from '@douyinfe/semi-ui';

const TopUpAdminStats = ({ stats, loading, t }) => {
  const statItems = [
    {
      label: t('总订单数'),
      value: stats?.total_count || 0,
      color: 'var(--semi-color-primary)',
    },
    {
      label: t('成功订单'),
      value: stats?.success_count || 0,
      color: 'var(--semi-color-success)',
    },
    {
      label: t('待支付订单'),
      value: stats?.pending_count || 0,
      color: 'var(--semi-color-warning)',
    },
    {
      label: t('已过期订单'),
      value: stats?.expired_count || 0,
      color: 'var(--semi-color-tertiary)',
    },
    {
      label: t('成功金额'),
      value: `¥${(stats?.success_amount || 0).toFixed(2)}`,
      color: 'var(--semi-color-success)',
    },
  ];

  if (loading && !stats) {
    return (
      <div className='grid grid-cols-2 md:grid-cols-5 gap-4 mb-4'>
        {[1, 2, 3, 4, 5].map((i) => (
          <Card key={i} className='text-center'>
            <Skeleton.Title style={{ width: 60, margin: '0 auto' }} />
            <Skeleton.Paragraph rows={1} style={{ marginTop: 8 }} />
          </Card>
        ))}
      </div>
    );
  }

  return (
    <div className='grid grid-cols-2 md:grid-cols-5 gap-4 mb-4'>
      {statItems.map((item, index) => (
        <Card key={index} className='text-center'>
          <div
            className='text-2xl font-bold'
            style={{ color: item.color }}
          >
            {item.value}
          </div>
          <div className='text-sm text-gray-500 mt-1'>{item.label}</div>
        </Card>
      ))}
    </div>
  );
};

export default TopUpAdminStats;
