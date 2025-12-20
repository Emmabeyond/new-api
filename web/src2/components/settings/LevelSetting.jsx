import React, { useEffect, useState, useMemo } from 'react';
import {
  Card,
  Button,
  Form,
  Space,
  Tag,
  Popconfirm,
  Toast,
  Spin,
  Tabs,
  TabPane,
  Select,
  Typography,
  Empty,
} from '@douyinfe/semi-ui';
import {
  IconPlus,
  IconDelete,
  IconSave,
  IconClose,
  IconUser,
} from '@douyinfe/semi-icons';
import { Users } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import {
  getAllLevels,
  createLevel,
  updateLevel,
  deleteLevel,
  getLevelStats,
} from '../../services/levelService';
import { API } from '../../helpers';
import { getLevelTheme, parseBenefits } from '../../utils/levelUtils';
import { useResponsive } from '../../hooks/common/useResponsive';
import LevelDiscountMatrix from './LevelDiscountMatrix';
import LevelGroupRateLimits from './LevelGroupRateLimits';
import '../../styles/level-setting.css';

const { Title, Text } = Typography;

const LevelSetting = () => {
  const { t } = useTranslation();
  const { isMobile } = useResponsive();
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);
  const [levels, setLevels] = useState([]);
  const [stats, setStats] = useState({});
  const [selectedLevelId, setSelectedLevelId] = useState(null);
  const [isCreating, setIsCreating] = useState(false);
  const [formApi, setFormApi] = useState(null);
  const [channelGroups, setChannelGroups] = useState([]);

  // 按优先级排序的等级列表
  const sortedLevels = useMemo(() => {
    return [...levels].sort((a, b) => a.priority - b.priority);
  }, [levels]);

  // 当前选中的等级
  const selectedLevel = useMemo(() => {
    if (isCreating) return null;
    return levels.find((l) => l.id === selectedLevelId) || null;
  }, [levels, selectedLevelId, isCreating]);

  // 总用户数
  const totalUsers = useMemo(() => {
    return Object.values(stats).reduce((sum, count) => sum + count, 0);
  }, [stats]);

  const fetchData = async () => {
    setLoading(true);
    try {
      const [levelsRes, statsRes] = await Promise.all([
        getAllLevels(),
        getLevelStats(),
      ]);
      if (levelsRes.success) {
        setLevels(levelsRes.data || []);
      }
      if (statsRes.success) {
        setStats(statsRes.data || {});
      }
    } catch (error) {
      Toast.error(t('加载数据失败'));
    } finally {
      setLoading(false);
    }
  };

  const fetchChannelGroups = async () => {
    try {
      const res = await API.get('/api/level/channel-groups');
      if (res.data.success) {
        setChannelGroups(res.data.data || []);
      }
    } catch (error) {
      console.error('Failed to fetch channel groups:', error);
    }
  };

  useEffect(() => {
    fetchData();
    fetchChannelGroups();
  }, []);

  // 处理添加等级
  const handleAdd = () => {
    setSelectedLevelId(null);
    setIsCreating(true);
  };

  // 处理选择等级
  const handleSelectLevel = (levelId) => {
    setIsCreating(false);
    setSelectedLevelId(levelId);
  };

  // 处理取消
  const handleCancel = () => {
    setIsCreating(false);
    setSelectedLevelId(null);
  };

  // 处理删除
  const handleDelete = async (id) => {
    try {
      const res = await deleteLevel(id);
      if (res.success) {
        Toast.success(t('删除成功'));
        if (selectedLevelId === id) {
          setSelectedLevelId(null);
        }
        fetchData();
      } else {
        Toast.error(res.message || t('删除失败'));
      }
    } catch (error) {
      Toast.error(t('删除失败'));
    }
  };

  // 处理保存
  const handleSubmit = async () => {
    if (!formApi) return;

    try {
      await formApi.validate();
      const values = formApi.getValues();

      // 构建权益配置
      const existingBenefits = selectedLevel
        ? parseBenefits(selectedLevel.benefits)
        : {};

      const benefits = {
        available_channel_groups: values.available_groups || ['default'],
        discount_ratio: values.discount_ratio || 1.0,
        group_discount_ratios: existingBenefits.group_discount_ratios || {},
        rate_limit: {
          total_count: values.rate_limit_total || 0,
          success_count: values.rate_limit_success || 0,
        },
      };

      const levelData = {
        id: values.id,
        name: values.name,
        description: values.description || '',
        priority: values.priority || 1,
        is_default: values.is_default || false,
        min_cumulative_recharge: values.min_cumulative_recharge || 0,
        benefits: JSON.stringify(benefits),
      };

      setSaving(true);
      let res;
      if (isCreating) {
        res = await createLevel(levelData);
      } else {
        res = await updateLevel(selectedLevelId, levelData);
      }

      if (res.success) {
        Toast.success(isCreating ? t('创建成功') : t('更新成功'));
        if (isCreating) {
          setIsCreating(false);
          setSelectedLevelId(values.id);
        }
        fetchData();
      } else {
        Toast.error(res.message || t('操作失败'));
      }
    } catch (error) {
      // 表单验证失败
      console.error('Form validation failed:', error);
    } finally {
      setSaving(false);
    }
  };

  // 处理折扣矩阵保存
  const handleSaveDiscountMatrix = async (updatedBenefits) => {
    if (!selectedLevel) return;

    const levelData = {
      id: selectedLevel.id,
      name: selectedLevel.name,
      description: selectedLevel.description || '',
      priority: selectedLevel.priority || 1,
      is_default: selectedLevel.is_default || false,
      min_cumulative_recharge: selectedLevel.min_cumulative_recharge || 0,
      benefits: JSON.stringify(updatedBenefits),
    };

    const res = await updateLevel(selectedLevel.id, levelData);
    if (res.success) {
      fetchData();
    } else {
      throw new Error(res.message || t('保存失败'));
    }
  };

  // 获取表单初始值
  const getInitialValues = () => {
    if (isCreating) {
      return {
        priority: levels.length + 1,
        discount_ratio: 1.0,
        rate_limit_total: 0,
        rate_limit_success: 1000,
        available_groups: ['default'],
      };
    }

    if (!selectedLevel) return {};

    const benefits = parseBenefits(selectedLevel.benefits);

    return {
      id: selectedLevel.id,
      name: selectedLevel.name,
      description: selectedLevel.description,
      priority: selectedLevel.priority,
      is_default: selectedLevel.is_default,
      min_cumulative_recharge: selectedLevel.min_cumulative_recharge,
      discount_ratio: benefits.discount_ratio || 1.0,
      rate_limit_total: benefits.rate_limit?.total_count || 0,
      rate_limit_success: benefits.rate_limit?.success_count || 0,
      available_groups: benefits.available_channel_groups || ['default'],
    };
  };

  if (loading && levels.length === 0) {
    return (
      <div className="level-setting-loading">
        <Spin size="large" />
      </div>
    );
  }

  return (
    <div className="level-setting-page">
      {/* 页面头部 */}
      <div className="level-setting-header">
        <div className="level-setting-header-left">
          <Title heading={3}>{t('用户等级配置')}</Title>
          <div className="level-setting-stats">
            <Tag color="blue" size="large">
              <Users size={14} style={{ marginRight: 4 }} />
              {t('总用户')}: {totalUsers}
            </Tag>
          </div>
        </div>
        <Button icon={<IconPlus />} theme="solid" onClick={handleAdd}>
          {t('添加等级')}
        </Button>
      </div>

      {/* 主内容区 */}
      <Card className="level-setting-main-card">
        <div className={`level-setting-layout ${isMobile ? 'mobile' : ''}`}>
          {/* 左侧等级列表 */}
          <LevelSidebar
            levels={sortedLevels}
            stats={stats}
            selectedLevelId={selectedLevelId}
            isCreating={isCreating}
            onSelect={handleSelectLevel}
          />

          {/* 右侧详情面板 */}
          <div className="level-detail-panel">
            {!selectedLevel && !isCreating ? (
              <div className="level-detail-empty">
                <div className="level-detail-empty-icon">
                  <Users size={48} />
                </div>
                <Text type="tertiary">{t('请选择一个等级进行编辑，或点击添加等级')}</Text>
              </div>
            ) : (
              <LevelDetailPanel
                level={selectedLevel}
                isCreating={isCreating}
                channelGroups={channelGroups}
                stats={stats}
                saving={saving}
                onSave={handleSubmit}
                onDelete={handleDelete}
                onCancel={handleCancel}
                onSaveDiscountMatrix={handleSaveDiscountMatrix}
                getFormApi={setFormApi}
                getInitialValues={getInitialValues}
                t={t}
              />
            )}
          </div>
        </div>
      </Card>
    </div>
  );
};


/**
 * 左侧等级列表组件
 */
const LevelSidebar = ({ levels, stats, selectedLevelId, isCreating, onSelect }) => {
  const { t } = useTranslation();

  return (
    <div className="level-setting-sidebar">
      <div className="level-sidebar-list">
        {/* 创建新等级的占位项 */}
        {isCreating && (
          <div className="level-sidebar-item creating">
            <div
              className="level-sidebar-dot"
              style={{ background: '#52c41a' }}
            />
            <div className="level-sidebar-content">
              <div className="level-sidebar-name">
                <span className="level-sidebar-name-text" style={{ color: '#52c41a' }}>
                  {t('新等级')}
                </span>
                <Tag size="small" color="green">{t('创建中')}</Tag>
              </div>
            </div>
          </div>
        )}

        {/* 等级列表 */}
        {levels.map((level) => {
          const theme = getLevelTheme(level.priority);
          const isSelected = !isCreating && selectedLevelId === level.id;
          const userCount = stats[level.id] || 0;

          return (
            <div
              key={level.id}
              className={`level-sidebar-item ${isSelected ? 'selected' : ''}`}
              onClick={() => onSelect(level.id)}
              style={isSelected ? { borderLeftColor: theme.primary } : {}}
            >
              <div
                className="level-sidebar-dot"
                style={{ background: theme.primary }}
              />
              <div className="level-sidebar-content">
                <div className="level-sidebar-name">
                  <span
                    className="level-sidebar-name-text"
                    style={{ color: isSelected ? theme.primary : undefined }}
                  >
                    {level.name}
                  </span>
                  {level.is_default && (
                    <Tag size="small" color="green">{t('默认')}</Tag>
                  )}
                </div>
                <div className="level-sidebar-meta">
                  <span>P{level.priority}</span>
                  <span>·</span>
                  <span className="level-sidebar-user-count">
                    <IconUser size="small" />
                    {userCount}
                  </span>
                </div>
              </div>
            </div>
          );
        })}

        {levels.length === 0 && !isCreating && (
          <div className="level-sidebar-empty">
            <Text type="tertiary">{t('暂无等级配置')}</Text>
            <Text type="tertiary" size="small">{t('点击右上角添加等级')}</Text>
          </div>
        )}
      </div>
    </div>
  );
};

/**
 * 右侧详情面板组件
 */
const LevelDetailPanel = ({
  level,
  isCreating,
  channelGroups,
  stats,
  saving,
  onSave,
  onDelete,
  onCancel,
  onSaveDiscountMatrix,
  getFormApi,
  getInitialValues,
  t,
}) => {
  const theme = level ? getLevelTheme(level.priority) : null;
  const userCount = level ? (stats[level.id] || 0) : 0;
  const hasUsers = userCount > 0;

  return (
    <>
      {/* 详情头部 */}
      <div className="level-detail-header">
        <div className="level-detail-title">
          {level && (
            <div
              style={{
                width: 12,
                height: 12,
                borderRadius: '50%',
                background: theme?.primary,
              }}
            />
          )}
          <Title heading={4} style={{ margin: 0 }}>
            {isCreating ? t('创建新等级') : level?.name}
          </Title>
          {level?.is_default && (
            <Tag color="green">{t('默认等级')}</Tag>
          )}
        </div>
        <div className="level-detail-actions">
          {isCreating && (
            <Button icon={<IconClose />} onClick={onCancel}>
              {t('取消')}
            </Button>
          )}
          {!isCreating && level && !level.is_default && (
            <Popconfirm
              title={t('确定删除此等级？')}
              content={hasUsers ? t('该等级下有 {{count}} 个用户，删除后用户将被重置为默认等级', { count: userCount }) : undefined}
              onConfirm={() => onDelete(level.id)}
            >
              <Button icon={<IconDelete />} type="danger">
                {t('删除')}
              </Button>
            </Popconfirm>
          )}
          <Button
            icon={<IconSave />}
            theme="solid"
            onClick={onSave}
            loading={saving}
          >
            {isCreating ? t('创建') : t('保存')}
          </Button>
        </div>
      </div>

      {/* 表单内容 */}
      <Tabs type="line" key={isCreating ? 'creating' : level?.id}>
        <TabPane tab={t('基本信息')} itemKey="basic">
          <div className="level-detail-form">
            <Form
              getFormApi={getFormApi}
              initValues={getInitialValues()}
              labelPosition="left"
              labelWidth={140}
            >
              {/* 基本信息 */}
              <div className="form-section">
                <div className="form-section-title">{t('基本信息')}</div>
                <Form.Input
                  field="id"
                  label={t('等级ID')}
                  placeholder={t('如: tier_1, tier_2')}
                  rules={[
                    { required: true, message: t('请输入等级ID') },
                    {
                      pattern: /^[a-z0-9_-]+$/,
                      message: t('只允许小写字母、数字、下划线和连字符'),
                    },
                  ]}
                  disabled={!isCreating}
                />
                <Form.Input
                  field="name"
                  label={t('等级名称')}
                  placeholder={t('如: Tier 1, VIP')}
                  rules={[{ required: true, message: t('请输入等级名称') }]}
                />
                <Form.TextArea
                  field="description"
                  label={t('描述')}
                  placeholder={t('等级描述')}
                  rows={2}
                />
                <div className="form-row">
                  <Form.InputNumber
                    field="priority"
                    label={t('优先级')}
                    min={1}
                    max={100}
                    placeholder={t('数值越大等级越高')}
                  />
                  <Form.InputNumber
                    field="min_cumulative_recharge"
                    label={t('最低累计充值($)')}
                    min={0}
                    step={1}
                    placeholder={t('0表示无条件')}
                  />
                </div>
                <Form.Switch field="is_default" label={t('默认等级')} />
              </div>

              {/* 权益配置 */}
              <div className="form-section">
                <div className="form-section-title">{t('权益配置')}</div>
                <Form.Select
                  field="available_groups"
                  label={t('可用渠道分组')}
                  placeholder={t('选择可用的渠道分组')}
                  multiple
                  filter
                  style={{ width: '100%' }}
                >
                  {channelGroups.map((group) => (
                    <Select.Option key={group.key} value={group.key}>
                      {group.label || group.key}
                    </Select.Option>
                  ))}
                </Form.Select>
                <Form.InputNumber
                  field="discount_ratio"
                  label={t('全局优惠倍率')}
                  min={0}
                  max={1}
                  step={0.01}
                  precision={2}
                  placeholder={t('1.0表示无折扣，0.8表示8折')}
                />
                <div className="form-row">
                  <Form.InputNumber
                    field="rate_limit_total"
                    label={t('总请求数限制/分钟')}
                    min={0}
                    placeholder={t('0表示无限制')}
                  />
                  <Form.InputNumber
                    field="rate_limit_success"
                    label={t('成功请求数限制/分钟')}
                    min={0}
                    placeholder={t('0表示无限制')}
                  />
                </div>
              </div>
            </Form>
          </div>
        </TabPane>

        {!isCreating && level && (
          <TabPane tab={t('渠道折扣矩阵')} itemKey="discount">
            <LevelDiscountMatrix
              levelId={level.id}
              benefits={parseBenefits(level.benefits)}
              onSave={onSaveDiscountMatrix}
            />
          </TabPane>
        )}

        {!isCreating && level && (
          <TabPane tab={t('渠道限流配置')} itemKey="rateLimit">
            <LevelGroupRateLimits
              levelId={level.id}
              benefits={parseBenefits(level.benefits)}
              channelGroups={channelGroups}
              onSave={onSaveDiscountMatrix}
            />
          </TabPane>
        )}
      </Tabs>
    </>
  );
};

export default LevelSetting;
