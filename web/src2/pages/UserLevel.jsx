/*
Copyright (C) 2025 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.
*/

import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { Typography, Spin, Card, Tag } from '@douyinfe/semi-ui';
import { showError, API } from '../helpers';
import LevelBadge from '../components/level/LevelBadge';
import LevelSkeleton from '../components/level/LevelSkeleton';
import BenefitsCards from '../components/level/BenefitsCards';
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
        // 排序逻辑：default 排第一，其他按字母顺序
        const sortedGroups = validGroups.sort((a, b) => {
          // default 分组始终排在第一位
          if (a.key === 'default') return -1;
          if (b.key === 'default') return 1;
          // 其他分组按 key 字母顺序排列
          return a.key.localeCompare(b.key);
        });
        setChannelGroups(sortedGroups);
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

  // 计算每个分组在所有等级中的可用性（只要有一个等级可用就算可用）
  const getGroupAvailabilityScore = (groupKey) => {
    let hasAvailable = false;
    let hasUnavailable = false;
    
    for (const lvl of sortedLevels) {
      const lvlBenefits = parseBenefits(lvl.benefits);
      const availableGroups = lvlBenefits.available_channel_groups || [];
      const isAvailable = availableGroups.length === 0 || availableGroups.includes(groupKey);
      
      if (isAvailable) {
        hasAvailable = true;
      } else {
        hasUnavailable = true;
      }
    }
    
    // 如果所有等级都不可用，返回 0（排到最后）
    if (!hasAvailable) return 0;
    // 如果所有等级都可用，返回 2（排在前面）
    if (!hasUnavailable) return 2;
    // 如果部分可用，返回 1（排在中间）
    return 1;
  };

  // 对分组进行最终排序：default 第一，全部可用的在前，部分可用的在中间，全部不可用的在后
  const finalSortedGroups = [...channelGroups].sort((a, b) => {
    // default 分组始终排在第一位
    if (a.key === 'default') return -1;
    if (b.key === 'default') return 1;
    
    // 根据所有等级的可用性排序
    const aScore = getGroupAvailabilityScore(a.key);
    const bScore = getGroupAvailabilityScore(b.key);
    
    if (aScore !== bScore) {
      return bScore - aScore; // 分数高的排在前面
    }
    
    // 分数相同时，按字母顺序
    return a.key.localeCompare(b.key);
  });

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
                      {finalSortedGroups.map((group) => (
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
                          {finalSortedGroups.map((group) => {
                            // 检查该等级是否允许使用此渠道分组
                            const availableGroups = lvlBenefits.available_channel_groups || [];
                            const isGroupAvailable = availableGroups.length === 0 || availableGroups.includes(group.key);
                            
                            const ratio = getDiscountRatio(lvlBenefits, group.key);
                            const modelLimits = lvlBenefits.model_rate_limits || {};
                            const groupRateLimit = getRateLimitForGroup(lvlBenefits, group.key);

                            return (
                              <td key={group.key}>
                                <LevelCellContent
                                  ratio={ratio}
                                  modelLimits={modelLimits}
                                  groupRateLimit={groupRateLimit}
                                  isAvailable={isGroupAvailable}
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
      <BenefitsCards benefits={benefits} />
    </div>
  );
};

/**
 * 单元格内容组件
 */
const LevelCellContent = ({ ratio, modelLimits, groupRateLimit, isAvailable, t }) => {
  // 如果该等级不允许使用此渠道，显示 -
  if (!isAvailable) {
    return (
      <div className="cell-content">
        <Text type="tertiary" style={{ fontSize: '18px', fontWeight: 500 }}>-</Text>
      </div>
    );
  }

  const hasModelLimits = modelLimits && Object.keys(modelLimits).length > 0;
  const discountText = formatDiscountForChannel(ratio, t);
  const hasDiscount = ratio < 1.0;
  const hasRateLimit = groupRateLimit && (groupRateLimit.total_count > 0 || groupRateLimit.success_count > 0);

  return (
    <div className="cell-content">
      {hasDiscount && (
        <Text style={{ color: 'var(--theme-success, #52c41a)', fontWeight: 500 }}>{discountText}</Text>
      )}
      {/* 始终显示并发限制信息 */}
      <div className="cell-rate-limit">
        <Text size="small" type="tertiary">
          {t('level.concurrency_limit')}: {' '}
          {hasRateLimit ? (
            groupRateLimit.success_count > 0 ? (
              `${formatRateLimitValue(groupRateLimit.success_count, t)}/min`
            ) : null
          ) : (
            t('level.unlimited')
          )}
        </Text>
      </div>
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

export default UserLevel;
