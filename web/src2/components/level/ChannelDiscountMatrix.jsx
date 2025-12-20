import React, { useState, useEffect } from 'react';
import { Card, Typography, Spin, Toast } from '@douyinfe/semi-ui';
import { useTranslation } from 'react-i18next';
import { API } from '../../helpers';
import LevelBadge from './LevelBadge';
import { 
  getLevelTheme, 
  parseBenefits, 
  formatCurrency,
  getDiscountRatio,
  formatDiscountForChannel
} from '../../utils/levelUtils';
import '../../styles/level.css';

const { Title, Text } = Typography;

/**
 * 渠道折扣矩阵组件
 * 展示等级-渠道折扣的二维表格
 * 行：等级，列：渠道分组
 */
const ChannelDiscountMatrix = ({ allLevels, currentLevel, userCumulativeUSD }) => {
  const { t, i18n } = useTranslation();
  const [channelGroups, setChannelGroups] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchChannelGroups();
  }, []);

  const fetchChannelGroups = async () => {
    try {
      const res = await API.get('/api/level/channel-groups');
      if (res.data.success) {
        // 过滤掉无效的分组（key 为空的）
        const validGroups = (res.data.data || []).filter(
          (group) => group.key && group.key.trim() !== ''
        );
        setChannelGroups(validGroups);
      } else {
        Toast.error(t('level.error.fetch_channel_groups_failed'));
      }
    } catch (error) {
      console.error('Failed to fetch channel groups:', error);
      Toast.error(t('level.error.fetch_channel_groups_failed'));
    } finally {
      setLoading(false);
    }
  };

  /**
   * 获取指定等级和渠道的折扣倍率
   */
  const getChannelDiscount = (level, channelGroupKey) => {
    const benefits = parseBenefits(level.benefits);
    return getDiscountRatio(benefits, channelGroupKey);
  };

  if (loading) {
    return (
      <Card className="channel-discount-matrix">
        <Spin size="large" />
      </Card>
    );
  }

  if (!allLevels || allLevels.length === 0) {
    return null;
  }

  if (!channelGroups || channelGroups.length === 0) {
    return (
      <Card className="channel-discount-matrix">
        <Text type="secondary">{t('level.error.no_channel_groups')}</Text>
      </Card>
    );
  }

  // 按优先级排序等级
  const sortedLevels = [...allLevels].sort((a, b) => a.priority - b.priority);

  return (
    <Card className="channel-discount-matrix fade-in" role="region" aria-labelledby="matrix-title">
      <Title heading={3} style={{ marginBottom: '16px' }} id="matrix-title">
        {t('level.comparison.channel_discount_matrix')}
      </Title>

      <div className="table-wrapper">
        <table role="table" aria-label="渠道折扣矩阵">
          <thead>
            <tr role="row">
              <th className="level-name-col" role="columnheader" scope="col">
                {t('level.comparison.level')}
              </th>
              <th className="upgrade-condition-col" role="columnheader" scope="col">
                {t('level.comparison.upgrade_condition')}
              </th>
              {channelGroups.map((group) => (
                <th key={group.key} role="columnheader" scope="col">
                  {group.label || group.key}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {sortedLevels.map((level) => {
              const isCurrentLevel = currentLevel && level.id === currentLevel.id;
              const theme = getLevelTheme(level.priority);

              return (
                <tr
                  key={level.id}
                  className={isCurrentLevel ? 'current-level-row' : ''}
                  role="row"
                  style={
                    isCurrentLevel
                      ? {
                          '--tier-bg': theme.bg,
                          '--tier-primary': theme.primary,
                        }
                      : {}
                  }
                >
                  {/* 等级名称列 */}
                  <td role="cell" className="level-name-cell">
                    <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                      <LevelBadge level={level} size="small" />
                      <Text strong>{level.name}</Text>
                      {isCurrentLevel && (
                        <Text type="tertiary" size="small">
                          ({t('level.current')})
                        </Text>
                      )}
                    </div>
                  </td>

                  {/* 升级条件列 */}
                  <td role="cell" className="upgrade-condition-cell">
                    {level.min_cumulative_recharge === 0 ? (
                      <Text type="secondary">-</Text>
                    ) : (
                      <Text>{formatCurrency(level.min_cumulative_recharge, '$', i18n.language)}</Text>
                    )}
                  </td>

                  {/* 渠道折扣列 */}
                  {channelGroups.map((group) => {
                    const ratio = getChannelDiscount(level, group.key);
                    const discountText = formatDiscountForChannel(ratio, t);

                    return (
                      <td key={group.key} role="cell" className="discount-cell">
                        {discountText === '-' || discountText === t('level.format.no_discount') ? (
                          <Text type="secondary">{discountText}</Text>
                        ) : (
                          <Text style={{ color: '#52c41a', fontWeight: 500 }}>{discountText}</Text>
                        )}
                      </td>
                    );
                  })}
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>
    </Card>
  );
};

export default ChannelDiscountMatrix;
