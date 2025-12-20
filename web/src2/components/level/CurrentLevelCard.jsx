import React from 'react';
import { Card, Typography, Tag } from '@douyinfe/semi-ui';
import { useTranslation } from 'react-i18next';
import LevelBadge from './LevelBadge';
import { getLevelTheme } from '../../utils/levelUtils';
import '../../styles/level.css';

const { Title, Text } = Typography;

/**
 * 当前等级卡片组件
 * 显示用户当前等级的详细信息
 */
const CurrentLevelCard = ({ level }) => {
  const { t } = useTranslation();

  if (!level) {
    return null;
  }

  const levelTheme = getLevelTheme(level.priority);

  return (
    <Card
      className={`current-level-card ${levelTheme.className} fade-in`}
      bodyStyle={{ padding: 0 }}
      aria-labelledby="current-level-title"
      role="region"
      aria-label="当前等级信息"
    >
      <div className="card-content">
        <LevelBadge level={level} size="large" theme={levelTheme} />

        <div className="level-info">
          <Title heading={2} style={{ margin: '0 0 8px 0' }} id="current-level-title">
            {level.name || t('level.unknown')}
          </Title>
          <Text type="secondary" style={{ fontSize: '14px' }}>
            {level.description || t('level.no_description')}
          </Text>
        </div>

        <Tag
          size="large"
          color={levelTheme.color}
          style={{
            padding: '8px 16px',
            fontSize: '14px',
            fontWeight: 500,
            color: 'var(--text-primary)', /* 确保文本对比度 */
          }}
          aria-label="当前等级标识"
        >
          {t('level.current')}
        </Tag>
      </div>
    </Card>
  );
};

export default CurrentLevelCard;
