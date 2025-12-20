import React, { useEffect, useState } from 'react';
import { Card, Progress, Tag, Typography, Space, Spin } from '@douyinfe/semi-ui';
import { getUserLevel } from '../services/levelService';

const { Title, Text } = Typography;

const UserLevelCard = () => {
  const [loading, setLoading] = useState(true);
  const [levelInfo, setLevelInfo] = useState(null);

  useEffect(() => {
    const fetchLevelInfo = async () => {
      try {
        const res = await getUserLevel();
        if (res.success) {
          setLevelInfo(res.data);
        }
      } catch (error) {
        console.error('Failed to fetch level info:', error);
      } finally {
        setLoading(false);
      }
    };
    fetchLevelInfo();
  }, []);

  if (loading) {
    return (
      <Card>
        <Spin />
      </Card>
    );
  }

  if (!levelInfo || !levelInfo.level) {
    return null;
  }

  const { level, benefits, upgrade_progress } = levelInfo;

  return (
    <Card
      title={
        <Space>
          <span>我的等级</span>
          <Tag color="blue" size="large">
            {level.name}
          </Tag>
        </Space>
      }
      style={{ marginBottom: 16 }}
    >
      <Space vertical align="start" style={{ width: '100%' }}>
        {level.description && (
          <Text type="tertiary">{level.description}</Text>
        )}

        {benefits && (
          <div style={{ marginTop: 8 }}>
            <Text strong>当前权益：</Text>
            <div style={{ marginTop: 4 }}>
              {benefits.discount_ratio < 1 && (
                <Tag color="green" style={{ marginRight: 8 }}>
                  {(benefits.discount_ratio * 100).toFixed(0)}% 折扣
                </Tag>
              )}
              {benefits.available_channel_groups?.length > 0 && (
                <Tag color="cyan">
                  可用分组: {benefits.available_channel_groups.join(', ')}
                </Tag>
              )}
            </div>
          </div>
        )}

        {upgrade_progress && upgrade_progress.next_level && (
          <div style={{ width: '100%', marginTop: 16 }}>
            <Text>
              升级到 <Text strong>{upgrade_progress.next_level.name}</Text>
            </Text>
            <Progress
              percent={Math.min(upgrade_progress.progress_percent || 0, 100)}
              showInfo
              style={{ marginTop: 8 }}
            />
            <Text type="tertiary" size="small">
              累计充值: ${upgrade_progress.current_recharge?.toFixed(2) || 0} / $
              {upgrade_progress.required_recharge?.toFixed(2) || 0}
            </Text>
          </div>
        )}
      </Space>
    </Card>
  );
};

export default UserLevelCard;
