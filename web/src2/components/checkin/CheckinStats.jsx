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
import { Typography } from '@douyinfe/semi-ui';
import { IconCalendarClock, IconCoinMoneyStroked } from '@douyinfe/semi-icons';
import { Flame, Trophy, Coins } from 'lucide-react';

const { Title, Text } = Typography;

const CheckinStats = ({ stats, t, renderQuota }) => {
  if (!stats) return null;

  const statItems = [
    {
      icon: <Flame size={28} />,
      label: t('连续签到'),
      value: `${stats.consecutive_days}`,
      unit: t('天'),
      highlight: stats.consecutive_days >= 7,
      gradient: 'linear-gradient(135deg, #f97316 0%, #fb923c 100%)',
      bgGradient: 'linear-gradient(135deg, rgba(249, 115, 22, 0.1) 0%, rgba(251, 146, 60, 0.05) 100%)',
    },
    {
      icon: <Trophy size={28} />,
      label: t('累计签到'),
      value: `${stats.total_checkins}`,
      unit: t('次'),
      gradient: 'linear-gradient(135deg, #8b5cf6 0%, #a855f7 100%)',
      bgGradient: 'linear-gradient(135deg, rgba(139, 92, 246, 0.1) 0%, rgba(168, 85, 247, 0.05) 100%)',
    },
    {
      icon: <Coins size={28} />,
      label: t('累计获得'),
      value: renderQuota(stats.total_quota),
      unit: '',
      gradient: 'linear-gradient(135deg, #10b981 0%, #34d399 100%)',
      bgGradient: 'linear-gradient(135deg, rgba(16, 185, 129, 0.1) 0%, rgba(52, 211, 153, 0.05) 100%)',
    },
  ];

  return (
    <div className='grid grid-cols-1 md:grid-cols-3 gap-4'>
      {statItems.map((item, index) => (
        <div 
          key={index} 
          className={`checkin-stats-card ${item.highlight ? 'highlight' : ''}`}
          style={{ background: item.bgGradient }}
        >
          <div className='flex flex-col items-center gap-3'>
            <div 
              className='stats-icon-wrapper'
              style={{ 
                background: item.gradient,
                width: 56,
                height: 56,
                borderRadius: 16,
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                color: '#fff',
                boxShadow: `0 8px 16px ${item.gradient.includes('f97316') ? 'rgba(249, 115, 22, 0.3)' : 
                  item.gradient.includes('8b5cf6') ? 'rgba(139, 92, 246, 0.3)' : 'rgba(16, 185, 129, 0.3)'}`
              }}
            >
              {item.icon}
            </div>
            <Text type='secondary' size='small'>{item.label}</Text>
            <div className='flex items-baseline gap-1'>
              <Title 
                heading={3} 
                style={{ 
                  margin: 0,
                  background: item.gradient,
                  WebkitBackgroundClip: 'text',
                  WebkitTextFillColor: 'transparent',
                  backgroundClip: 'text',
                }}
              >
                {item.value}
              </Title>
              {item.unit && (
                <Text type='secondary' size='small'>{item.unit}</Text>
              )}
            </div>
          </div>
        </div>
      ))}
    </div>
  );
};

export default CheckinStats;
