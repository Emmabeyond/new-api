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
import { Card, Typography } from '@douyinfe/semi-ui';
import { IconCalendarClock, IconCoinMoneyStroked } from '@douyinfe/semi-icons';
import { Flame } from 'lucide-react';

const { Title, Text } = Typography;

const CheckinStats = ({ stats, t, renderQuota }) => {
  if (!stats) return null;

  const statItems = [
    {
      icon: <Flame size={28} style={{ color: 'var(--semi-color-warning)' }} />,
      label: t('连续签到'),
      value: `${stats.consecutive_days} ${t('天')}`,
      highlight: stats.consecutive_days >= 7,
    },
    {
      icon: <IconCalendarClock size='extra-large' style={{ color: 'var(--semi-color-primary)' }} />,
      label: t('累计签到'),
      value: `${stats.total_checkins} ${t('次')}`,
    },
    {
      icon: <IconCoinMoneyStroked size='extra-large' style={{ color: 'var(--semi-color-success)' }} />,
      label: t('累计获得'),
      value: renderQuota(stats.total_quota),
    },
  ];

  return (
    <div className='grid grid-cols-1 md:grid-cols-3 gap-4'>
      {statItems.map((item, index) => (
        <Card key={index} className='text-center'>
          <div className='flex flex-col items-center gap-2'>
            {item.icon}
            <Text type='secondary' size='small'>{item.label}</Text>
            <Title 
              heading={4} 
              style={item.highlight ? { color: 'var(--semi-color-warning)' } : {}}
            >
              {item.value}
            </Title>
          </div>
        </Card>
      ))}
    </div>
  );
};

export default CheckinStats;
