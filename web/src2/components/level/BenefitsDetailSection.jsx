import React from 'react';
import { Card, Typography, Tag, Divider } from '@douyinfe/semi-ui';
import { useTranslation } from 'react-i18next';
import { Layers, Percent, Zap, Gift } from './icons';
import { formatDiscount, formatRateLimit } from '../../utils/levelUtils';
import '../../styles/level.css';

const { Title, Text } = Typography;

/**
 * 权益详情区域组件
 * 显示当前等级的详细权益信息
 */
const BenefitsDetailSection = ({ benefits }) => {
  const { t } = useTranslation();

  if (!benefits) {
    return null;
  }

  const {
    available_channel_groups = [],
    discount_ratio = 1.0,
    group_discount_ratios = {},
    rate_limit = {},
  } = benefits || {};

  const hasSpecialDiscounts = group_discount_ratios && Object.keys(group_discount_ratios).length > 0;

  return (
    <div className="benefits-detail-section fade-in" role="region" aria-labelledby="benefits-title">
      <Title heading={3} style={{ marginBottom: '16px' }} id="benefits-title">
        {t('level.benefits.title')}
      </Title>

      <div className="benefits-grid">
        {/* 可用渠道分组卡片 */}
        <Card className="benefit-card" role="article" aria-label="可用渠道分组">
          <div className="benefit-card-header">
            <Layers size={20} style={{ color: '#3370ff' }} aria-hidden="true" />
            <Text strong>{t('level.benefits.channel_groups')}</Text>
          </div>
          <div className="benefit-card-content">
            {available_channel_groups.length > 0 ? (
              <div style={{ display: 'flex', flexWrap: 'wrap', gap: '8px' }} role="list" aria-label="渠道分组列表">
                {available_channel_groups.map((group) => (
                  <Tag key={group} color="cyan" size="large" role="listitem">
                    {group}
                  </Tag>
                ))}
              </div>
            ) : (
              <Text type="secondary">{t('level.benefits.no_groups')}</Text>
            )}
          </div>
        </Card>

        {/* 优惠倍率卡片 */}
        <Card className="benefit-card" role="article" aria-label="优惠倍率">
          <div className="benefit-card-header">
            <Percent size={20} style={{ color: '#52c41a' }} aria-hidden="true" />
            <Text strong>{t('level.benefits.discount')}</Text>
          </div>
          <div className="benefit-card-content">
            <Text
              style={{
                fontSize: '32px',
                fontWeight: 'bold',
                color: '#52c41a',
                display: 'block',
                marginTop: '8px',
              }}
              aria-label={`优惠倍率 ${formatDiscount(discount_ratio, t)}`}
            >
              {formatDiscount(discount_ratio, t)}
            </Text>
          </div>
        </Card>

        {/* 速率限制卡片 */}
        <Card className="benefit-card" role="article" aria-label="速率限制">
          <div className="benefit-card-header">
            <Zap size={20} style={{ color: '#faad14' }} aria-hidden="true" />
            <Text strong>{t('level.benefits.rate_limit')}</Text>
          </div>
          <div className="benefit-card-content">
            <Text
              style={{
                fontSize: '20px',
                fontWeight: '500',
                display: 'block',
                marginTop: '8px',
              }}
              aria-label={`速率限制 ${formatRateLimit(rate_limit, t)}`}
            >
              {formatRateLimit(rate_limit, t)}
            </Text>
          </div>
        </Card>

        {/* 特殊分组优惠卡片 - 仅在有特殊优惠时显示 */}
        {hasSpecialDiscounts && (
          <Card className="benefit-card" role="article" aria-label="特殊分组优惠">
            <div className="benefit-card-header">
              <Gift size={20} style={{ color: '#722ed1' }} aria-hidden="true" />
              <Text strong>{t('level.benefits.special_discounts')}</Text>
            </div>
            <div className="benefit-card-content">
              <div style={{ display: 'flex', flexWrap: 'wrap', gap: '8px', marginTop: '8px' }} role="list" aria-label="特殊优惠列表">
                {Object.entries(group_discount_ratios).map(([group, ratio]) => (
                  <Tag key={group} color="orange" size="large" role="listitem">
                    {group}: {formatDiscount(ratio, t)}
                  </Tag>
                ))}
              </div>
            </div>
          </Card>
        )}
      </div>
    </div>
  );
};

export default BenefitsDetailSection;
