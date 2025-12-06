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
import { Typography, Collapsible } from '@douyinfe/semi-ui';
import { IconChevronDown } from '@douyinfe/semi-icons';
import { Gift, Sparkles, Calendar, Zap } from 'lucide-react';

const { Text, Title } = Typography;

const CheckinRules = ({ t, renderQuota }) => {
  const rules = [
    {
      icon: <Calendar size={18} />,
      title: t('基础奖励'),
      desc: t('每日签到即可获得基础额度奖励'),
      color: '#8B5CF6',
    },
    {
      icon: <Zap size={18} />,
      title: t('连续签到加成'),
      items: [
        { days: '1-6', label: t('第1-6天'), multiplier: '1x' },
        { days: '7', label: t('第7天'), multiplier: '5x', highlight: true },
        { days: '8-13', label: t('第8-13天'), multiplier: '1.5x' },
        { days: '14-29', label: t('第14-29天'), multiplier: '2x' },
        { days: '30', label: t('第30天'), multiplier: '20x', highlight: true },
        { days: '31+', label: t('第31天起'), multiplier: '2.5x' },
      ],
      color: '#F59E0B',
    },
    {
      icon: <Sparkles size={18} />,
      title: t('惊喜奖励'),
      desc: t('每次签到有10%概率触发额外惊喜奖励'),
      color: '#EC4899',
    },
    {
      icon: <Gift size={18} />,
      title: t('补签规则'),
      desc: t('可补签最近3天，补签奖励为正常奖励的50%'),
      color: '#10B981',
    },
  ];

  return (
    <div className="checkin-rules-card">
      <Collapsible
        defaultOpen={true}
        collapseHeight={0}
        trigger={(isOpen) => (
          <div className="checkin-rules-trigger">
            <div className="flex items-center gap-2">
              <div className="checkin-rules-icon">
                <Sparkles size={16} />
              </div>
              <Title heading={5} style={{ margin: 0, color: 'var(--semi-color-text-0)' }}>
                {t('签到规则')}
              </Title>
            </div>
            <IconChevronDown 
              style={{ 
                transform: isOpen ? 'rotate(180deg)' : 'rotate(0deg)',
                transition: 'transform 0.3s ease'
              }} 
            />
          </div>
        )}
      >
        <div className="checkin-rules-content">
          {rules.map((rule, index) => (
            <div key={index} className="checkin-rule-item">
              <div 
                className="checkin-rule-icon" 
                style={{ backgroundColor: `${rule.color}20`, color: rule.color }}
              >
                {rule.icon}
              </div>
              <div className="checkin-rule-body">
                <Text strong style={{ color: 'var(--semi-color-text-0)' }}>
                  {rule.title}
                </Text>
                {rule.desc && (
                  <Text type="tertiary" size="small" style={{ display: 'block', marginTop: 4 }}>
                    {rule.desc}
                  </Text>
                )}
                {rule.items && (
                  <div className="checkin-rule-items">
                    {rule.items.map((item, idx) => (
                      <div 
                        key={idx} 
                        className={`checkin-rule-tier ${item.highlight ? 'highlight' : ''}`}
                      >
                        <span className="tier-label">{item.label}</span>
                        <span className={`tier-multiplier ${item.highlight ? 'gold' : ''}`}>
                          {item.multiplier}
                        </span>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            </div>
          ))}
        </div>
      </Collapsible>
    </div>
  );
};

export default CheckinRules;
