import React, { useState, useEffect } from 'react';
import { Card, Typography, Spin, Toast, Tag } from '@douyinfe/semi-ui';
import { useTranslation } from 'react-i18next';
import { API } from '../../helpers';
import LevelBadge from './LevelBadge';
import { 
  getLevelTheme, 
  parseBenefits, 
  getDiscountRatio, 
  formatDiscountForChannel,
  formatCurrency 
} from '../../utils/levelUtils';
import '../../styles/level.css';

const { Title, Text } = Typography;

/**
 * 等级对比表格组件
 * 展示渠道分组×等级的二维折扣矩阵
 */
const LevelComparisonTable = ({ allLevels, currentLevel }) => {
  const { t, i18n } = useTranslation();
  const [channelGroups, setChannelGroups] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    fetchChannelGroups();
  }, []);

  /**
   * 获取所有渠道分组列表
   */
  const fetchChannelGroups = async () => {
    try {
      setLoading(true);
      setError(null);
      const res = await API.get('/api/level/channel-groups');
      
      if (res.data.success) {
        const groups = res.data.data || [];
        // 排序逻辑：default 排第一，其他按字母顺序
        const sortedGroups = groups.sort((a, b) => {
          // default 分组始终排在第一位
          if (a.key === 'default') return -1;
          if (b.key === 'default') return 1;
          // 其他分组按 key 字母顺序排列
          return a.key.localeCompare(b.key);
        });
        setChannelGroups(sortedGroups);
      } else {
        const errorMsg = res.data.message || t('level.error.fetch_channel_groups_failed');
        setError(errorMsg);
        Toast.error(errorMsg);
      }
    } catch (err) {
      console.error('Failed to fetch channel groups:', err);
      const errorMsg = t('level.error.fetch_channel_groups_failed');
      setError(errorMsg);
      Toast.error(errorMsg);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <Card className="level-comparison-table">
        <div style={{ 
          display: 'flex', 
          justifyContent: 'center', 
          alignItems: 'center', 
          minHeight: '200px',
          padding: '24px'
        }}>
          <Spin size="large" tip={t('level.loading')} />
        </div>
      </Card>
    );
  }

  if (error) {
    return (
      <Card className="level-comparison-table">
        <div style={{ 
          display: 'flex', 
          justifyContent: 'center', 
          alignItems: 'center', 
          minHeight: '200px',
          padding: '24px'
        }}>
          <Typography.Text type="danger">{error}</Typography.Text>
        </div>
      </Card>
    );
  }

  if (!allLevels || allLevels.length === 0) {
    return null;
  }

  // 按优先级排序等级
  const sortedLevels = [...allLevels].sort((a, b) => a.priority - b.priority);

  /**
   * 渲染折扣单元格
   */
  const renderDiscountCell = (level, channelGroup) => {
    const benefits = parseBenefits(level.benefits);
    const ratio = getDiscountRatio(benefits, channelGroup);
    const formattedDiscount = formatDiscountForChannel(ratio, t);
    
    if (formattedDiscount === '-' || formattedDiscount === t('level.format.no_discount')) {
      return <Text type="secondary">-</Text>;
    }
    
    return <Text style={{ color: '#52c41a', fontWeight: 500 }}>{formattedDiscount}</Text>;
  };

  return (
    <Card className="level-comparison-table fade-in" role="region" aria-labelledby="comparison-title">
      <Title heading={3} style={{ marginBottom: '16px' }} id="comparison-title">
        {t('level.comparison.discount_matrix_title') || '等级折扣对比'}
      </Title>

      <div className="table-wrapper">
        <table role="table" aria-label="渠道分组折扣矩阵">
          <thead>
            <tr role="row">
              <th className="benefit-name-col" role="columnheader" scope="col">
                {t('level.comparison.level_name') || '等级'}
              </th>
              <th role="columnheader" scope="col">
                {t('level.comparison.upgrade_condition') || '升级条件'}
              </th>
              {channelGroups.map((group) => (
                <th key={group.key || group} role="columnheader" scope="col" style={{ textAlign: 'center' }}>
                  <Tag color="cyan" size="large">{group.label || group.key || group}</Tag>
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
                  role="row"
                  className={isCurrentLevel ? 'current-level-row' : ''}
                  style={
                    isCurrentLevel
                      ? {
                          '--tier-bg': theme.bg,
                          '--tier-primary': theme.primary,
                        }
                      : {}
                  }
                >
                  <td 
                    role="rowheader" 
                    scope="row"
                    className={isCurrentLevel ? 'current-level-cell' : ''}
                  >
                    <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                      <LevelBadge level={level} size="small" />
                      <span style={{ fontWeight: isCurrentLevel ? 600 : 400 }}>
                        {level.name}
                      </span>
                    </div>
                  </td>
                  <td 
                    role="cell"
                    className={isCurrentLevel ? 'current-level-cell' : ''}
                  >
                    {level.min_cumulative_recharge === 0 ? (
                      <Text type="secondary">-</Text>
                    ) : (
                      <Text>{formatCurrency(level.min_cumulative_recharge, '$', i18n.language)}</Text>
                    )}
                  </td>
                  {channelGroups.map((group) => (
                    <td 
                      key={group.key || group} 
                      role="cell"
                      style={{ textAlign: 'center' }}
                      className={isCurrentLevel ? 'current-level-cell' : ''}
                    >
                      {renderDiscountCell(level, group.key || group)}
                    </td>
                  ))}
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>
    </Card>
  );
};

export default LevelComparisonTable;
