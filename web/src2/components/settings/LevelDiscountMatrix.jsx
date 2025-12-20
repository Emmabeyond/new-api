import React, { useState, useEffect } from 'react';
import {
  Card,
  Table,
  InputNumber,
  Button,
  Toast,
  Spin,
  Typography,
  Space,
  Popconfirm,
} from '@douyinfe/semi-ui';
import { IconSave, IconDelete } from '@douyinfe/semi-icons';
import { useTranslation } from 'react-i18next';
import { API } from '../../helpers';

const { Title, Text } = Typography;

/**
 * LevelDiscountMatrix 组件
 * 管理员用于编辑等级-渠道折扣矩阵的组件
 * 
 * @param {Object} props
 * @param {string} props.levelId - 当前编辑的等级 ID
 * @param {Object} props.benefits - 当前等级的权益配置
 * @param {Function} props.onSave - 保存回调函数
 */
const LevelDiscountMatrix = ({ levelId, benefits, onSave }) => {
  const { t } = useTranslation();
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [channelGroups, setChannelGroups] = useState([]);
  const [discountRatios, setDiscountRatios] = useState({});
  const [errors, setErrors] = useState({});

  useEffect(() => {
    fetchChannelGroups();
  }, []);

  useEffect(() => {
    if (benefits) {
      // 初始化折扣倍率
      const groupDiscounts = benefits.group_discount_ratios || {};
      setDiscountRatios(groupDiscounts);
    }
  }, [benefits]);

  const fetchChannelGroups = async () => {
    try {
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
        Toast.error('获取渠道分组失败');
      }
    } catch (error) {
      console.error('Failed to fetch channel groups:', error);
      Toast.error('获取渠道分组失败');
    } finally {
      setLoading(false);
    }
  };

  /**
   * 处理折扣倍率变更
   */
  const handleDiscountChange = (channelGroup, value) => {
    setDiscountRatios((prev) => ({
      ...prev,
      [channelGroup]: value,
    }));

    // 清除该渠道的错误
    if (errors[channelGroup]) {
      setErrors((prev) => {
        const newErrors = { ...prev };
        delete newErrors[channelGroup];
        return newErrors;
      });
    }
  };

  /**
   * 验证折扣倍率
   */
  const validateDiscount = (channelGroup, value) => {
    if (value === null || value === undefined || value === '') {
      return true; // 空值表示使用全局折扣
    }

    if (value < 0 || value > 1) {
      setErrors((prev) => ({
        ...prev,
        [channelGroup]: '折扣倍率必须在 0 到 1 之间',
      }));
      return false;
    }

    return true;
  };

  /**
   * 清空指定渠道的折扣配置
   */
  const handleClearDiscount = (channelGroup) => {
    setDiscountRatios((prev) => {
      const newRatios = { ...prev };
      delete newRatios[channelGroup];
      return newRatios;
    });

    // 清除错误
    if (errors[channelGroup]) {
      setErrors((prev) => {
        const newErrors = { ...prev };
        delete newErrors[channelGroup];
        return newErrors;
      });
    }

    Toast.success(`已清空 ${channelGroup} 的折扣配置，将使用全局折扣`);
  };

  /**
   * 保存折扣矩阵
   */
  const handleSave = async () => {
    // 验证所有折扣倍率
    let hasError = false;
    const newErrors = {};

    Object.entries(discountRatios).forEach(([channelGroup, value]) => {
      if (value !== null && value !== undefined && value !== '') {
        if (value < 0 || value > 1) {
          newErrors[channelGroup] = '折扣倍率必须在 0 到 1 之间';
          hasError = true;
        }
      }
    });

    if (hasError) {
      setErrors(newErrors);
      Toast.error('请修正输入错误');
      return;
    }

    // 构建新的 benefits 对象
    const updatedBenefits = {
      ...benefits,
      group_discount_ratios: discountRatios,
    };

    setSaving(true);
    try {
      await onSave(updatedBenefits);
      Toast.success('保存成功');
    } catch (error) {
      Toast.error('保存失败');
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return (
      <Card>
        <Spin size="large" />
      </Card>
    );
  }

  if (!channelGroups || channelGroups.length === 0) {
    return (
      <Card>
        <Text type="secondary">暂无渠道分组</Text>
      </Card>
    );
  }

  const columns = [
    {
      title: '渠道分组',
      dataIndex: 'channelGroupLabel',
      width: 200,
      render: (text) => <Text strong>{text}</Text>,
    },
    {
      title: '折扣倍率',
      dataIndex: 'discount',
      width: 200,
      render: (_, record) => {
        const channelGroup = record.channelGroup;
        const value = discountRatios[channelGroup];
        const hasError = !!errors[channelGroup];

        return (
          <div>
            <InputNumber
              value={value}
              onChange={(val) => handleDiscountChange(channelGroup, val)}
              onBlur={() => validateDiscount(channelGroup, value)}
              min={0}
              max={1}
              step={0.01}
              precision={2}
              placeholder="留空使用全局折扣"
              style={{ width: '100%' }}
              validateStatus={hasError ? 'error' : undefined}
            />
            {hasError && (
              <Text type="danger" size="small">
                {errors[channelGroup]}
              </Text>
            )}
          </div>
        );
      },
    },
    {
      title: '显示折扣',
      dataIndex: 'display',
      width: 150,
      render: (_, record) => {
        const channelGroup = record.channelGroup;
        const value = discountRatios[channelGroup];
        const globalDiscount = benefits?.discount_ratio || 1.0;

        if (value === null || value === undefined || value === '') {
          return (
            <Text type="secondary">
              使用全局: {(globalDiscount * 100).toFixed(0)}%
            </Text>
          );
        }

        if (value === 1.0) {
          return <Text type="secondary">无折扣</Text>;
        }

        return (
          <Text style={{ color: 'var(--theme-success, #52c41a)', fontWeight: 500 }}>
            {(value * 100).toFixed(0)}% 折扣
          </Text>
        );
      },
    },
    {
      title: '操作',
      width: 100,
      render: (_, record) => {
        const channelGroup = record.channelGroup;
        const hasDiscount = discountRatios[channelGroup] !== null &&
                           discountRatios[channelGroup] !== undefined &&
                           discountRatios[channelGroup] !== '';

        return (
          <Popconfirm
            title="确定清空此渠道的折扣配置？"
            content="清空后将使用全局折扣倍率"
            onConfirm={() => handleClearDiscount(channelGroup)}
            disabled={!hasDiscount}
          >
            <Button
              icon={<IconDelete />}
              size="small"
              type="tertiary"
              disabled={!hasDiscount}
            >
              清空
            </Button>
          </Popconfirm>
        );
      },
    },
  ];

  const dataSource = channelGroups.map((group) => ({
    key: group.key || group,
    channelGroup: group.key || group,
    channelGroupLabel: group.label || group.key || group,
  }));

  return (
    <Card
      title={
        <Space>
          <Title heading={5} style={{ margin: 0 }}>
            渠道折扣矩阵配置
          </Title>
          <Text type="tertiary" size="small">
            (等级: {levelId})
          </Text>
        </Space>
      }
      headerExtraContent={
        <Button
          icon={<IconSave />}
          theme="solid"
          onClick={handleSave}
          loading={saving}
          disabled={Object.keys(errors).length > 0}
        >
          保存配置
        </Button>
      }
    >
      <div style={{ marginBottom: '16px' }}>
        <Text type="secondary">
          为每个渠道分组配置独立的折扣倍率。留空表示使用全局折扣倍率 (
          {benefits?.discount_ratio ? (benefits.discount_ratio * 100).toFixed(0) : 100}%)。
          折扣倍率范围: 0-1，其中 1 表示无折扣，0.8 表示 8 折。
        </Text>
      </div>

      <Table
        columns={columns}
        dataSource={dataSource}
        pagination={false}
        bordered
      />
    </Card>
  );
};

export default LevelDiscountMatrix;
