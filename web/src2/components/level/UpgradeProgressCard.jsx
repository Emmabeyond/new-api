import React from 'react';
import { Card, Typography, Progress } from '@douyinfe/semi-ui';
import { useTranslation } from 'react-i18next';
import { Trophy, TrendingUp, DollarSign, Target, Zap } from './icons';
import LevelBadge from './LevelBadge';
import {
  formatCurrency,
  getProgressColor,
  calculateRemaining,
} from '../../utils/levelUtils';
import '../../styles/level.css';

const { Title, Text } = Typography;

/**
 * 升级进度卡片组件
 * 显示用户升级到下一等级的进度信息
 */
const UpgradeProgressCard = ({ progress }) => {
  const { t, i18n } = useTranslation();

  // 如果没有下一等级，说明已达最高等级
  if (!progress || !progress.next_level) {
    return (
      <Card 
        className="upgrade-progress-card max-level fade-in"
        role="region"
        aria-label="最高等级状态"
      >
        <Trophy size={48} style={{ color: '#faad14', marginBottom: '16px' }} aria-hidden="true" />
        <Title heading={3} style={{ margin: '8px 0' }}>
          {t('level.max_level')}
        </Title>
        <Text type="secondary">{t('level.max_level_desc')}</Text>
      </Card>
    );
  }

  const {
    current_recharge = 0,
    required_recharge = 0,
    progress_percent = 0,
    next_level,
  } = progress;

  const remaining = calculateRemaining(current_recharge, required_recharge);
  const progressColor = getProgressColor(progress_percent);

  return (
    <Card 
      className="upgrade-progress-card fade-in"
      role="region"
      aria-labelledby="upgrade-progress-title"
      aria-describedby="upgrade-progress-desc"
    >
      {/* 标题 */}
      <div className="progress-header">
        <TrendingUp size={24} style={{ color: '#52c41a' }} aria-hidden="true" />
        <Title heading={4} style={{ margin: 0 }} id="upgrade-progress-title">
          {t('level.upgrade_progress')}
        </Title>
      </div>

      {/* 下一等级信息 */}
      <div className="next-level-info" id="upgrade-progress-desc">
        <Text>{t('level.next_level')}:</Text>
        <LevelBadge level={next_level} size="small" />
        <Text strong>{next_level.name}</Text>
      </div>

      {/* 进度条 */}
      <div className="progress-bar-wrapper">
        <Progress
          percent={progress_percent}
          stroke={progressColor}
          showInfo={false}
          style={{ flex: 1 }}
          aria-label={`升级进度 ${progress_percent.toFixed(1)}%`}
          role="progressbar"
          aria-valuenow={progress_percent}
          aria-valuemin={0}
          aria-valuemax={100}
        />
        <Text type="secondary" style={{ minWidth: '50px', textAlign: 'right' }} aria-live="polite">
          {progress_percent.toFixed(1)}%
        </Text>
      </div>

      {/* 充值信息 */}
      <div className="recharge-info">
        <div className="info-item" role="group" aria-label="累计充值金额">
          <DollarSign size={20} style={{ color: '#52c41a', flexShrink: 0 }} aria-hidden="true" />
          <div>
            <Text type="secondary" style={{ fontSize: '12px', display: 'block' }}>
              {t('level.cumulative_recharge')}
            </Text>
            <Text strong>{formatCurrency(current_recharge, '$', i18n.language)}</Text>
          </div>
        </div>

        <div className="info-item" role="group" aria-label="升级目标金额">
          <Target size={20} style={{ color: '#3370ff', flexShrink: 0 }} aria-hidden="true" />
          <div>
            <Text type="secondary" style={{ fontSize: '12px', display: 'block' }}>
              {t('level.required_recharge')}
            </Text>
            <Text strong>{formatCurrency(required_recharge, '$', i18n.language)}</Text>
          </div>
        </div>

        <div className="info-item highlight" role="group" aria-label="还需充值金额">
          <Zap size={20} style={{ color: '#faad14', flexShrink: 0 }} aria-hidden="true" />
          <div>
            <Text type="secondary" style={{ fontSize: '12px', display: 'block' }}>
              {t('level.remaining_recharge')}
            </Text>
            <Text strong style={{ color: '#faad14' }}>
              {formatCurrency(remaining, '$', i18n.language)}
            </Text>
          </div>
        </div>
      </div>
    </Card>
  );
};

export default UpgradeProgressCard;
