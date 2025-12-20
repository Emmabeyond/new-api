import React from 'react';
import { Card, Typography, Tag, Divider } from '@douyinfe/semi-ui';
import { useTranslation } from 'react-i18next';
import LevelBadge from './LevelBadge';
import {
  formatCurrency,
  formatDiscount,
  formatRateLimit,
  parseBenefits,
  getLevelTheme,
} from '../../utils/levelUtils';

const { Title, Text } = Typography;

/**
 * 移动端等级对比组件
 * 以卡片列表形式展示等级对比
 */
const LevelComparisonMobile = ({ allLevels, currentLevel }) => {
  const { t } = useTranslation();

  if (!allLevels || allLevels.length === 0) {
    return null;
  }

  // 按优先级排序等级
  const sortedLevels = [...allLevels].sort((a, b) => a.priority - b.priority);

  return (
    <div className="level-comparison-mobile">
      <Title heading={3} style={{ marginBottom: '16px' }}>
        {t('level.comparison.title')}
      </Title>

      {sortedLevels.map((level) => {
        const isCurrentLevel = currentLevel && level.id === currentLevel.id;
        const theme = getLevelTheme(level.priority);
        const benefits = parseBenefits(level.benefits);

        return (
          <Card
            key={level.id}
            className={`level-comparison-card ${isCurrentLevel ? 'current-level' : ''}`}
            style={{
              marginBottom: '16px',
              border: isCurrentLevel ? `2px solid ${theme.primary}` : undefined,
              background: isCurrentLevel ? theme.bg : undefined,
            }}
          >
            {/* 等级标题 */}
            <div style={{ display: 'flex', alignItems: 'center', gap: '12px', marginBottom: '16px' }}>
              <LevelBadge level={level} size="medium" />
              <div style={{ flex: 1 }}>
                <Title heading={4} style={{ margin: 0 }}>
                  {level.name}
                </Title>
                {level.description && (
                  <Text type="secondary" size="small">
                    {level.description}
                  </Text>
                )}
              </div>
              {isCurrentLevel && (
                <Tag color={theme.color} size="large">
                  {t('level.current')}
                </Tag>
              )}
            </div>

            <Divider margin="12px" />

            {/* 权益详情 */}
            <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
              {/* 升级条件 */}
              <div>
                <Text type="secondary" size="small">
                  {t('level.comparison.upgrade_condition')}
                </Text>
                <div style={{ marginTop: '4px' }}>
                  {level.min_cumulative_recharge === 0 ? (
                    <Text>-</Text>
                  ) : (
                    <Text strong>{formatCurrency(level.min_cumulative_recharge)}</Text>
                  )}
                </div>
              </div>

              {/* 可用渠道分组 */}
              <div>
                <Text type="secondary" size="small">
                  {t('level.benefits.channel_groups')}
                </Text>
                <div style={{ marginTop: '4px', display: 'flex', flexWrap: 'wrap', gap: '4px' }}>
                  {benefits.available_channel_groups && benefits.available_channel_groups.length > 0 ? (
                    benefits.available_channel_groups.map((group) => (
                      <Tag key={group} size="small" color="cyan">
                        {group}
                      </Tag>
                    ))
                  ) : (
                    <Text>-</Text>
                  )}
                </div>
              </div>

              {/* 优惠倍率 */}
              <div>
                <Text type="secondary" size="small">
                  {t('level.benefits.discount')}
                </Text>
                <div style={{ marginTop: '4px' }}>
                  {formatDiscount(benefits.discount_ratio) === '无折扣' ? (
                    <Text>-</Text>
                  ) : (
                    <Text strong style={{ color: '#52c41a' }}>
                      {formatDiscount(benefits.discount_ratio)}
                    </Text>
                  )}
                </div>
              </div>

              {/* 速率限制 */}
              <div>
                <Text type="secondary" size="small">
                  {t('level.benefits.rate_limit')}
                </Text>
                <div style={{ marginTop: '4px' }}>
                  <Text>{formatRateLimit(benefits.rate_limit)}</Text>
                </div>
              </div>
            </div>
          </Card>
        );
      })}
    </div>
  );
};

export default LevelComparisonMobile;
