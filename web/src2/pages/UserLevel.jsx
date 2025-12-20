/*
Copyright (C) 2025 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.
*/

import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { Typography, Spin, Card, Tag, Collapsible } from '@douyinfe/semi-ui';
import { IconChevronDown } from '@douyinfe/semi-icons';
import { showError, API } from '../helpers';
import LevelBadge from '../components/level/LevelBadge';
import LevelSkeleton from '../components/level/LevelSkeleton';
import { useResponsive } from '../hooks/common/useResponsive';
import { useLevelData } from '../hooks/common/useLevelData';
import {
  getLevelTheme,
  parseBenefits,
  formatCurrency,
  getDiscountRatio,
  formatDiscountForChannel,
  getRateLimitForGroup,
  formatRateLimitValue,
} from '../utils/levelUtils';
import '../styles/level.css';

const { Title, Text } = Typography;

const UserLevel = () => {
  const { t, i18n } = useTranslation();
  const { isMobile } = useResponsive();
  const { loading, error, levelInfo, allLevels } = useLevelData();
  const [channelGroups, setChannelGroups] = useState([]);
  const [channelLoading, setChannelLoading] = useState(true);

  useEffect(() => {
    fetchChannelGroups();
  }, []);

  const fetchChannelGroups = async () => {
    try {
      const res = await API.get('/api/level/channel-groups');
      if (res.data.success) {
        const validGroups = (res.data.data || []).filter(
          (group) => group.key && group.key.trim() !== ''
        );
        setChannelGroups(validGroups);
      }
    } catch (err) {
      console.error('Failed to fetch channel groups:', err);
    } finally {
      setChannelLoading(false);
    }
  };

  if (error) {
    showError(t('加载等级信息失败'));
  }

  if (loading) {
    return <LevelSkeleton />;
  }

  if (!levelInfo) {
    return (
      <div className="flex justify-center items-center h-screen">
        <Text>{t('暂无等级信息')}</Text>
      </div>
    );
  }

  const { level, benefits } = levelInfo;
  const sortedLevels = [...(allLevels || [])].sort((a, b) => a.priority - b.priority);

  return (
    <div className="limits-page">
      {/* 页面标题 */}
      <div className="limits-header">
        <Title heading={2} style={{ margin: 0 }}>Limits</Title>
        <Tag color="blue" size="large">{level.name}</Tag>
      </div>

      <Card className="limits-main-card">
        <div className={`limits-layout ${isMobile ? 'mobile' : ''}`}>
          {/* 左侧等级列表 */}
          <div className="levels-sidebar">
            {sortedLevels.map((lvl) => {
              const isCurrentLevel = level && lvl.id === level.id;
              const theme = getLevelTheme(lvl.priority);
              
              return (
                <div
                  key={lvl.id}
                  className={`level-item ${isCurrentLevel ? 'current' : ''}`}
                  style={isCurrentLevel ? { borderLeftColor: theme.primary } : {}}
                >
                  <div className="level-item-dot" style={{ background: theme.primary }} />
                  <div className="level-item-content">
                    <Text strong style={{ color: theme.primary }}>{lvl.name}</Text>
                    <Text type="tertiary" size="small">
                      {lvl.min_cumulative_recharge === 0
                        ? t('level.no_requirement')
                        : `${t('level.cumulative_recharge')}: ${formatCurrency(lvl.min_cumulative_recharge, '$', i18n.language)}`}
                    </Text>
                  </div>
                </div>
              );
            })}

            {/* 当前等级提示 */}
            <div className="current-level-hint">
              <Text type="tertiary" size="small">{t('level.hint.recharge_upgrade')}</Text>
            </div>
          </div>

          {/* 右侧矩阵表格 */}
          <div className="limits-matrix">
            {channelLoading ? (
              <div className="matrix-loading">
                <Spin />
              </div>
            ) : (
              <div className="matrix-table-wrapper">
                <table className="matrix-table">
                  <thead>
                    <tr>
                      <th className="level-col">{t('level.comparison.level')}</th>
                      {channelGroups.map((group) => (
                        <th key={group.key}>
                          {group.label || group.key}
                        </th>
                      ))}
                    </tr>
                  </thead>
                  <tbody>
                    {sortedLevels.map((lvl) => {
                      const isCurrentLevel = level && lvl.id === level.id;
                      const theme = getLevelTheme(lvl.priority);
                      const lvlBenefits = parseBenefits(lvl.benefits);

                      return (
                        <tr
                          key={lvl.id}
                          className={isCurrentLevel ? 'current-level-row' : ''}
                          style={isCurrentLevel ? { background: theme.bg } : {}}
                        >
                          <td className="level-name">{lvl.name}</td>
                          {channelGroups.map((group) => {
                            const ratio = getDiscountRatio(lvlBenefits, group.key);
                            const modelLimits = lvlBenefits.model_rate_limits || {};
                            const groupRateLimit = getRateLimitForGroup(lvlBenefits, group.key);

                            return (
                              <td key={group.key}>
                                <LevelCellContent
                                  ratio={ratio}
                                  modelLimits={modelLimits}
                                  groupRateLimit={groupRateLimit}
                                  t={t}
                                />
                              </td>
                            );
                          })}
                        </tr>
                      );
                    })}
                  </tbody>
                </table>
              </div>
            )}
          </div>
        </div>
      </Card>

      {/* 底部专属权益 */}
      <BenefitsCollapsible benefits={benefits} t={t} />
    </div>
  );
};

/**
 * 单元格内容组件
 */
const LevelCellContent = ({ ratio, modelLimits, groupRateLimit, t }) => {
  const hasModelLimits = modelLimits && Object.keys(modelLimits).length > 0;
  const discountText = formatDiscountForChannel(ratio, t);
  const hasDiscount = ratio < 1.0;
  const hasRateLimit = groupRateLimit && (groupRateLimit.total_count > 0 || groupRateLimit.success_count > 0);

  if (!hasDiscount && !hasModelLimits && !hasRateLimit) {
    return <Text type="tertiary">-</Text>;
  }

  return (
    <div className="cell-content">
      {hasDiscount && (
        <Text style={{ color: '#52c41a', fontWeight: 500 }}>{discountText}</Text>
      )}
      {hasRateLimit && (
        <div className="cell-rate-limit">
          <Text size="small" type="tertiary">{t('level.concurrency_limit')}:</Text>
          <div className="rate-limit-values">
            {groupRateLimit.success_count > 0 && (
              <Text size="small">{formatRateLimitValue(groupRateLimit.success_count, t)}/min</Text>
            )}
          </div>
        </div>
      )}
      {hasModelLimits && (
        <div className="cell-limits">
          <Text size="small" type="tertiary">{t('level.concurrency')}:</Text>
          <div className="model-limits-list">
            {Object.entries(modelLimits).slice(0, 8).map(([model, limit]) => (
              <div key={model} className="model-limit-item">
                <Tag size="small" color="grey">{model}</Tag>
                <Text size="small">{limit}</Text>
              </div>
            ))}
            {Object.keys(modelLimits).length > 8 && (
              <Text type="tertiary" size="small">...</Text>
            )}
          </div>
          <div className="other-limit">
            <Text size="small" type="tertiary">{t('level.other')}</Text>
            <Text size="small">{t('level.unlimited')}</Text>
          </div>
        </div>
      )}
    </div>
  );
};

/**
 * 专属权益折叠面板
 */
const BenefitsCollapsible = ({ benefits, t }) => {
  if (!benefits) return null;

  const discount_ratio = benefits.discount_ratio ?? 1.0;
  const group_discount_ratios = benefits.group_discount_ratios || {};
  const rate_limit = benefits.rate_limit || {};
  const model_rate_limits = benefits.model_rate_limits || {};

  const hasGroupDiscounts = group_discount_ratios && Object.keys(group_discount_ratios).length > 0;
  const hasModelLimits = model_rate_limits && Object.keys(model_rate_limits).length > 0;

  return (
    <Card className="benefits-card">
      <div className="benefits-header">
        <Title heading={4} style={{ margin: 0 }}>{t('level.exclusive_benefits')}</Title>
        <Text type="tertiary" size="small">{t('level.benefits_hint')}</Text>
      </div>

      <div className="benefits-list">
        {/* 并发限制 */}
        <Collapsible
          collapseHeight={48}
          className="benefit-collapsible"
        >
          <div className="benefit-row">
            <Text className="benefit-label">{t('level.concurrency_limit')}</Text>
            <Text>{rate_limit?.total_count > 0 ? rate_limit.total_count : t('level.unlimited')}</Text>
          </div>
          {hasModelLimits && (
            <div className="benefit-detail">
              <Text type="tertiary" size="small">{t('level.model_concurrency')}</Text>
              <div className="model-limits-grid">
                {Object.entries(model_rate_limits).map(([model, limit]) => (
                  <div key={model} className="model-limit-row">
                    <Tag size="small">{model}</Tag>
                    <Text size="small">{limit}</Text>
                  </div>
                ))}
              </div>
            </div>
          )}
        </Collapsible>

        {/* 分组折扣 */}
        <Collapsible
          collapseHeight={48}
          className="benefit-collapsible"
        >
          <div className="benefit-row">
            <Text className="benefit-label">{t('level.group_discount')}</Text>
            <Text>{discount_ratio < 1.0 ? formatDiscountForChannel(discount_ratio, t) : t('level.no_discount')}</Text>
          </div>
          {hasGroupDiscounts && (
            <div className="benefit-detail">
              <Text type="tertiary" size="small">{t('level.group_discount_detail')}</Text>
              <div className="group-discounts-grid">
                {Object.entries(group_discount_ratios).map(([group, ratio]) => (
                  <div key={group} className="group-discount-row">
                    <Tag size="small" color="cyan">{group}</Tag>
                    <Text size="small" style={{ color: '#52c41a' }}>
                      {formatDiscountForChannel(ratio, t)}
                    </Text>
                  </div>
                ))}
              </div>
            </div>
          )}
        </Collapsible>

        {/* 模型折扣 */}
        <Collapsible
          collapseHeight={48}
          className="benefit-collapsible"
        >
          <div className="benefit-row">
            <Text className="benefit-label">{t('level.model_discount')}</Text>
            <Text>{t('level.no_discount')}</Text>
          </div>
          <div className="benefit-detail">
            <Text type="tertiary" size="small">{t('level.model_discount_detail')}</Text>
          </div>
        </Collapsible>
      </div>
    </Card>
  );
};

export default UserLevel;
