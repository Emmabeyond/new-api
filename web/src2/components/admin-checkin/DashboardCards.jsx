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

import React from 'react';
import { Card, Spin, Typography, Button } from '@douyinfe/semi-ui';
import { IconRefresh, IconCalendarClock, IconCoinMoneyStroked, IconUser } from '@douyinfe/semi-icons';
import { Flame } from 'lucide-react';

const { Title, Text } = Typography;

const DashboardCards = ({ dashboard, loading, t, renderQuota, onRefresh }) => {
  const cards = [
    {
      icon: <IconCalendarClock size='extra-large' style={{ color: 'var(--semi-color-primary)' }} />,
      label: t('今日签到'),
      value: dashboard?.today_checkins || 0,
      suffix: t('次'),
    },
    {
      icon: <IconCoinMoneyStroked size='extra-large' style={{ color: 'var(--semi-color-success)' }} />,
      label: t('今日发放额度'),
      value: renderQuota(dashboard?.today_quota_distributed || 0),
      suffix: '',
    },
    {
      icon: <IconUser size='extra-large' style={{ color: 'var(--semi-color-tertiary)' }} />,
      label: t('活跃用户'),
      value: dashboard?.active_users || 0,
      suffix: t('人'),
    },
    {
      icon: <Flame size={28} style={{ color: 'var(--semi-color-warning)' }} />,
      label: t('平均连续天数'),
      value: (dashboard?.avg_consecutive_days || 0).toFixed(1),
      suffix: t('天'),
    },
  ];

  return (
    <Card
      title={t('签到统计')}
      headerExtraContent={
        <Button
          theme='borderless'
          icon={<IconRefresh />}
          onClick={onRefresh}
          loading={loading}
        >
          {t('刷新')}
        </Button>
      }
    >
      <Spin spinning={loading}>
        <div className='grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4'>
          {cards.map((card, index) => (
            <div 
              key={index} 
              className='p-4 rounded-lg text-center'
              style={{ backgroundColor: 'var(--semi-color-fill-0)' }}
            >
              <div className='flex flex-col items-center gap-2'>
                {card.icon}
                <Text type='secondary' size='small'>{card.label}</Text>
                <Title heading={4}>
                  {card.value}{card.suffix}
                </Title>
              </div>
            </div>
          ))}
        </div>
      </Spin>
    </Card>
  );
};

export default DashboardCards;
