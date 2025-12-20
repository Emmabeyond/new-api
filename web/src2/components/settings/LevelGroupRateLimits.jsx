import React, { useState, useEffect, useMemo } from 'react';
import {
  Card,
  Table,
  InputNumber,
  Button,
  Toast,
  Typography,
  Empty,
  Spin,
} from '@douyinfe/semi-ui';
import { IconSave } from '@douyinfe/semi-icons';
import { useTranslation } from 'react-i18next';
import { parseBenefits, isValidRateLimitValue } from '../../utils/levelUtils';

const { Text } = Typography;

/**
 * 渠道分组限流配置组件
 * 用于在等级配置页面中配置每个渠道分组的限流规则
 */
const LevelGroupRateLimits = ({
  levelId,
  benefits,
  channelGroups,
  onSave,
}) => {
  const { t } = useTranslation();
  const [loading, setLoading] = useState(false);
  const [editedLimits, setEditedLimits] = useState({});
  const [hasChanges, setHasChanges] = useState(false);

  // 解析当前的限流配置
  const parsedBenefits = useMemo(() => {
    return typeof benefits === 'string' ? parseBenefits(benefits) : benefits;
  }, [benefits]);

  // 初始化编辑状态
  useEffect(() => {
    const groupRateLimits = parsedBenefits?.group_rate_limits || {};
    setEditedLimits(groupRateLimits);
    setHasChanges(false);
  }, [parsedBenefits, levelId]);

  // 处理限流值变更
  const handleLimitChange = (groupKey, field, value) => {
    if (!isValidRateLimitValue(value)) {
      return;
    }

    setEditedLimits((prev) => {
      const newLimits = { ...prev };
      if (!newLimits[groupKey]) {
        newLimits[groupKey] = { total_count: 0, success_count: 0 };
      }
      newLimits[groupKey] = {
        ...newLimits[groupKey],
        [field]: value || 0,
      };
      return newLimits;
    });
    setHasChanges(true);
  };

  // 保存配置
  const handleSave = async () => {
    setLoading(true);
    try {
      // 清理空配置（两个值都为0的）
      const cleanedLimits = {};
      Object.entries(editedLimits).forEach(([key, limit]) => {
        if (limit.total_count > 0 || limit.success_count > 0) {
          cleanedLimits[key] = limit;
        }
      });

      const updatedBenefits = {
        ...parsedBenefits,
        group_rate_limits: cleanedLimits,
      };

      await onSave(updatedBenefits);
      Toast.success(t('保存成功'));
      setHasChanges(false);
    } catch (error) {
      Toast.error(error.message || t('保存失败'));
    } finally {
      setLoading(false);
    }
  };

  // 获取全局限流配置（用于显示默认值）
  const globalRateLimit = parsedBenefits?.rate_limit || { total_count: 0, success_count: 0 };

  // 表格列定义
  const columns = [
    {
      title: t('渠道分组'),
      dataIndex: 'groupKey',
      width: 150,
      render: (text, record) => (
        <Text strong>{record.label || record.groupKey}</Text>
      ),
    },
    {
      title: t('总请求数/分钟'),
      dataIndex: 'total_count',
      width: 180,
      render: (_, record) => {
        const limit = editedLimits[record.groupKey] || {};
        const value = limit.total_count;
        const isDefault = value === undefined || value === 0;
        
        return (
          <div className="rate-limit-input-cell">
            <InputNumber
              value={value || null}
              min={0}
              placeholder={isDefault ? (globalRateLimit.total_count || t('无限制')) : ''}
              onChange={(val) => handleLimitChange(record.groupKey, 'total_count', val)}
              style={{ width: 120 }}
            />
            {isDefault && globalRateLimit.total_count > 0 && (
              <Text type="tertiary" size="small" style={{ marginLeft: 4 }}>
                ({t('继承全局')})
              </Text>
            )}
          </div>
        );
      },
    },
    {
      title: t('成功请求数/分钟'),
      dataIndex: 'success_count',
      width: 180,
      render: (_, record) => {
        const limit = editedLimits[record.groupKey] || {};
        const value = limit.success_count;
        const isDefault = value === undefined || value === 0;
        
        return (
          <div className="rate-limit-input-cell">
            <InputNumber
              value={value || null}
              min={0}
              placeholder={isDefault ? (globalRateLimit.success_count || t('无限制')) : ''}
              onChange={(val) => handleLimitChange(record.groupKey, 'success_count', val)}
              style={{ width: 120 }}
            />
            {isDefault && globalRateLimit.success_count > 0 && (
              <Text type="tertiary" size="small" style={{ marginLeft: 4 }}>
                ({t('继承全局')})
              </Text>
            )}
          </div>
        );
      },
    },
  ];

  // 表格数据
  const dataSource = useMemo(() => {
    return channelGroups.map((group) => ({
      key: group.key,
      groupKey: group.key,
      label: group.label,
    }));
  }, [channelGroups]);

  if (!channelGroups || channelGroups.length === 0) {
    return (
      <div className="level-group-rate-limits-empty">
        <Empty description={t('暂无渠道分组配置')} />
      </div>
    );
  }

  return (
    <div className="level-group-rate-limits">
      <div className="level-group-rate-limits-header">
        <div className="level-group-rate-limits-title">
          <Text strong>{t('渠道分组限流配置')}</Text>
          <Text type="tertiary" size="small" style={{ marginLeft: 8 }}>
            {t('为每个渠道分组设置独立的限流规则，留空表示使用全局限流')}
          </Text>
        </div>
        <Button
          icon={<IconSave />}
          theme="solid"
          onClick={handleSave}
          loading={loading}
          disabled={!hasChanges}
        >
          {t('保存')}
        </Button>
      </div>

      <div className="level-group-rate-limits-info">
        <Text type="tertiary" size="small">
          {t('全局限流')}: {t('总请求数')} {globalRateLimit.total_count || t('无限制')}, {t('成功请求数')} {globalRateLimit.success_count || t('无限制')}
        </Text>
      </div>

      <Table
        columns={columns}
        dataSource={dataSource}
        pagination={false}
        size="small"
        className="level-group-rate-limits-table"
      />
    </div>
  );
};

export default LevelGroupRateLimits;
