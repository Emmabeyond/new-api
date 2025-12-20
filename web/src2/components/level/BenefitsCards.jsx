/*
Copyright (C) 2025 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.
*/

import React, { useState } from 'react';
import { Card, Typography, Tag, Modal, Tooltip } from '@douyinfe/semi-ui';
import { IconChevronRight } from '@douyinfe/semi-icons';
import { useTranslation } from 'react-i18next';
import { useResponsive } from '../../hooks/common/useResponsive';
import {
  formatDiscountForChannel,
  formatRateLimitValue,
} from '../../utils/levelUtils';
import './BenefitsCards.css';

const { Title, Text } = Typography;

/**
 * 折扣权益卡片
 */
const DiscountCard = ({ discountRatio, groupDiscountRatios, t, onClick }) => {
  const hasDiscount = discountRatio < 1.0;
  const hasGroupDiscounts = groupDiscountRatios && Object.keys(groupDiscountRatios).length > 0;
  const discountPercent = hasDiscount ? Math.round((1 - discountRatio) * 100) : 0;

  return (
    <div 
      className={`benefit-card discount-card ${hasGroupDiscounts ? 'clickable' : ''}`}
      onClick={hasGroupDiscounts ? onClick : undefined}
      role={hasGroupDiscounts ? 'button' : undefined}
      tabIndex={hasGroupDiscounts ? 0 : undefined}
    >
      <div className="benefit-card-icon">
        <span className="icon-emoji">💰</span>
      </div>
      <div className="benefit-card-content">
        <Text className="benefit-card-label">{t('level.benefits.discount')}</Text>
        <div className="benefit-card-value">
          {hasDiscount ? (
            <span className="discount-value">{discountPercent}% OFF</span>
          ) : (
            <span className="no-benefit">{t('level.no_discount')}</span>
          )}
        </div>
        {hasGroupDiscounts && (
          <div className="benefit-card-hint">
            <Text type="tertiary" size="small">
              {Object.keys(groupDiscountRatios).length} {t('level.group_discount_count')}
            </Text>
            <IconChevronRight size="small" />
          </div>
        )}
      </div>
    </div>
  );
};

/**
 * 并发限制卡片
 */
const RateLimitCard = ({ rateLimit, modelRateLimits, t, onClick }) => {
  const hasRateLimit = rateLimit && (rateLimit.total_count > 0 || rateLimit.success_count > 0);
  const hasModelLimits = modelRateLimits && Object.keys(modelRateLimits).length > 0;
  const displayValue = hasRateLimit 
    ? (rateLimit.success_count > 0 ? rateLimit.success_count : rateLimit.total_count)
    : null;

  return (
    <div 
      className={`benefit-card ratelimit-card ${hasModelLimits ? 'clickable' : ''}`}
      onClick={hasModelLimits ? onClick : undefined}
      role={hasModelLimits ? 'button' : undefined}
      tabIndex={hasModelLimits ? 0 : undefined}
    >
      <div className="benefit-card-icon">
        <span className="icon-emoji">⚡</span>
      </div>
      <div className="benefit-card-content">
        <Text className="benefit-card-label">{t('level.concurrency_limit')}</Text>
        <div className="benefit-card-value">
          {hasRateLimit ? (
            <span className="ratelimit-value">
              {formatRateLimitValue(displayValue, t)}
              <span className="unit">/min</span>
            </span>
          ) : (
            <span className="unlimited">{t('level.unlimited')}</span>
          )}
        </div>
        {hasModelLimits && (
          <div className="benefit-card-hint">
            <Text type="tertiary" size="small">
              {Object.keys(modelRateLimits).length} {t('level.model_limit_count')}
            </Text>
            <IconChevronRight size="small" />
          </div>
        )}
      </div>
    </div>
  );
};

/**
 * 可用渠道卡片
 */
const ChannelGroupsCard = ({ availableGroups, t }) => {
  const hasGroups = availableGroups && availableGroups.length > 0;
  const displayGroups = hasGroups ? availableGroups.slice(0, 3) : [];
  const moreCount = hasGroups ? Math.max(0, availableGroups.length - 3) : 0;

  return (
    <div className="benefit-card channels-card">
      <div className="benefit-card-icon">
        <span className="icon-emoji">🔗</span>
      </div>
      <div className="benefit-card-content">
        <Text className="benefit-card-label">{t('level.available_channels')}</Text>
        <div className="benefit-card-value channels-value">
          {!hasGroups ? (
            <span className="all-channels">{t('level.all_channels')}</span>
          ) : (
            <div className="channel-tags">
              {displayGroups.map((group) => (
                <Tag key={group} color="cyan" size="small">
                  {group}
                </Tag>
              ))}
              {moreCount > 0 && (
                <Tooltip content={availableGroups.slice(3).join(', ')}>
                  <Tag color="grey" size="small">+{moreCount}</Tag>
                </Tooltip>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

/**
 * 分组折扣详情弹窗
 */
const GroupDiscountModal = ({ visible, onClose, groupDiscountRatios, t }) => {
  if (!groupDiscountRatios) return null;

  return (
    <Modal
      title={t('level.group_discount_detail')}
      visible={visible}
      onCancel={onClose}
      footer={null}
      width={480}
    >
      <div className="group-discount-list">
        {Object.entries(groupDiscountRatios).map(([group, ratio]) => (
          <div key={group} className="group-discount-item">
            <Tag color="cyan" size="large">{group}</Tag>
            <Text style={{ color: 'var(--theme-success, #52c41a)', fontWeight: 600 }}>
              {formatDiscountForChannel(ratio, t)}
            </Text>
          </div>
        ))}
      </div>
    </Modal>
  );
};

/**
 * 模型限流详情弹窗
 */
const ModelLimitsModal = ({ visible, onClose, modelRateLimits, t }) => {
  if (!modelRateLimits) return null;

  return (
    <Modal
      title={t('level.model_concurrency')}
      visible={visible}
      onCancel={onClose}
      footer={null}
      width={560}
    >
      <div className="model-limits-list-modal">
        {Object.entries(modelRateLimits).map(([model, limit]) => (
          <div key={model} className="model-limit-item-modal">
            <Tag color="blue">{model}</Tag>
            <Text strong>{limit} /min</Text>
          </div>
        ))}
      </div>
    </Modal>
  );
};

/**
 * 专属权益卡片组件
 */
const BenefitsCards = ({ benefits }) => {
  const { t } = useTranslation();
  const { isMobile } = useResponsive();
  const [discountModalVisible, setDiscountModalVisible] = useState(false);
  const [modelLimitsModalVisible, setModelLimitsModalVisible] = useState(false);

  if (!benefits) return null;

  const {
    discount_ratio = 1.0,
    group_discount_ratios = {},
    rate_limit = {},
    model_rate_limits = {},
    available_channel_groups = [],
  } = benefits;

  return (
    <Card className="benefits-cards-container">
      <div className="benefits-cards-header">
        <Title heading={4} style={{ margin: 0 }}>{t('level.exclusive_benefits')}</Title>
        <Text type="tertiary" size="small">{t('level.benefits_hint')}</Text>
      </div>

      <div className={`benefits-cards-grid ${isMobile ? 'mobile' : ''}`}>
        <DiscountCard
          discountRatio={discount_ratio}
          groupDiscountRatios={group_discount_ratios}
          t={t}
          onClick={() => setDiscountModalVisible(true)}
        />
        <RateLimitCard
          rateLimit={rate_limit}
          modelRateLimits={model_rate_limits}
          t={t}
          onClick={() => setModelLimitsModalVisible(true)}
        />
        <ChannelGroupsCard
          availableGroups={available_channel_groups}
          t={t}
        />
      </div>

      <GroupDiscountModal
        visible={discountModalVisible}
        onClose={() => setDiscountModalVisible(false)}
        groupDiscountRatios={group_discount_ratios}
        t={t}
      />
      <ModelLimitsModal
        visible={modelLimitsModalVisible}
        onClose={() => setModelLimitsModalVisible(false)}
        modelRateLimits={model_rate_limits}
        t={t}
      />
    </Card>
  );
};

export default BenefitsCards;
