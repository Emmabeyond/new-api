import React from 'react';
import { Card, Skeleton } from '@douyinfe/semi-ui';
import '../../styles/level.css';

/**
 * 等级页面骨架屏组件
 * 在数据加载时显示占位符
 */
const LevelSkeleton = () => {
  return (
    <div className="user-level-page">
      {/* 标题骨架 */}
      <Skeleton.Title style={{ width: '200px', marginBottom: '24px' }} />

      {/* 当前等级卡片骨架 */}
      <Card style={{ marginBottom: '24px' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '16px' }}>
          <Skeleton.Avatar size="large" />
          <div style={{ flex: 1 }}>
            <Skeleton.Title style={{ width: '150px', marginBottom: '8px' }} />
            <Skeleton.Paragraph rows={1} style={{ width: '250px' }} />
          </div>
          <Skeleton.Button style={{ width: '80px' }} />
        </div>
      </Card>

      {/* 升级进度卡片骨架 */}
      <Card style={{ marginBottom: '24px' }}>
        <Skeleton.Title style={{ width: '120px', marginBottom: '16px' }} />
        <Skeleton.Paragraph rows={1} style={{ marginBottom: '16px' }} />
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: '16px' }}>
          <Skeleton.Paragraph rows={2} />
          <Skeleton.Paragraph rows={2} />
          <Skeleton.Paragraph rows={2} />
        </div>
      </Card>

      {/* 等级对比表格骨架 */}
      <Card style={{ marginBottom: '24px' }}>
        <Skeleton.Title style={{ width: '120px', marginBottom: '16px' }} />
        <Skeleton.Paragraph rows={5} />
      </Card>

      {/* 权益详情骨架 */}
      <div>
        <Skeleton.Title style={{ width: '120px', marginBottom: '16px' }} />
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(280px, 1fr))', gap: '16px' }}>
          <Card>
            <Skeleton.Paragraph rows={3} />
          </Card>
          <Card>
            <Skeleton.Paragraph rows={3} />
          </Card>
          <Card>
            <Skeleton.Paragraph rows={3} />
          </Card>
        </div>
      </div>
    </div>
  );
};

export default LevelSkeleton;
