import React from 'react';
import { Tag, Typography } from '@douyinfe/semi-ui';
import { useTranslation } from 'react-i18next';
import {
  formatCurrency,
  formatDiscount,
  formatRateLimit,
  parseBenefits,
  getLevelTheme,
} from '../../utils/levelUtils';

const { Text } = Typography;

/**
 * 权益行组件
 * 在等级对比表格中渲染一行权益数据
 */
const BenefitRow = ({ row, allLevels, currentLevel }) => {
  const { t, i18n } = useTranslation();

  /**
   * 根据权益类型渲染单元格内容
   */
  const renderCell = (level) => {
    const benefits = parseBenefits(level.benefits);

    switch (row.type) {
      case 'currency':
        // 升级条件（累计充值）
        if (level.min_cumulative_recharge === 0) {
          return <Text type="secondary">-</Text>;
        }
        return <Text>{formatCurrency(level.min_cumulative_recharge, '$', i18n.language)}</Text>;

      case 'tags':
        // 可用渠道分组
        if (!benefits.available_channel_groups || benefits.available_channel_groups.length === 0) {
          return <Text type="secondary">-</Text>;
        }
        return (
          <div className="tags-cell">
            {benefits.available_channel_groups.map((group) => (
              <Tag key={group} size="small" color="cyan">
                {group}
              </Tag>
            ))}
          </div>
        );

      case 'percentage':
        // 优惠倍率
        const discount = formatDiscount(benefits.discount_ratio, t);
        const noDiscountText = t('level.format.no_discount');
        if (discount === noDiscountText) {
          return <Text type="secondary">-</Text>;
        }
        return <Text style={{ color: '#52c41a', fontWeight: 500 }}>{discount}</Text>;

      case 'rate':
        // 速率限制
        const rateLimit = formatRateLimit(benefits.rate_limit, t);
        const unlimitedText = t('level.format.unlimited');
        if (rateLimit === unlimitedText) {
          return <Text type="secondary">{unlimitedText}</Text>;
        }
        return <Text>{rateLimit}</Text>;

      default:
        return <Text type="secondary">-</Text>;
    }
  };

  return (
    <tr className={row.key} role="row">
      <td className="benefit-name" role="rowheader" scope="row">{row.label}</td>
      {allLevels.map((level) => {
        const isCurrentLevel = currentLevel && level.id === currentLevel.id;
        const theme = getLevelTheme(level.priority);

        return (
          <td
            key={level.id}
            className={isCurrentLevel ? 'current-level-cell' : ''}
            role="cell"
            aria-label={`${level.name} - ${row.label}`}
            style={
              isCurrentLevel
                ? {
                    '--tier-bg': theme.bg,
                  }
                : {}
            }
          >
            {renderCell(level)}
          </td>
        );
      })}
    </tr>
  );
};

export default BenefitRow;
