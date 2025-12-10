/*
Copyright (C) 2025 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/

import React, { useEffect, useState } from 'react';
import {
  Button,
  Card,
  Form,
  Row,
  Col,
  InputNumber,
  Spin,
  Checkbox,
  Typography,
  Divider,
  Table,
  Tag,
  Popconfirm,
} from '@douyinfe/semi-ui';
import { API, showError, showSuccess, timestamp2string } from '../../helpers';
import { useTranslation } from 'react-i18next';
import { Shield, AlertTriangle } from 'lucide-react';

const { Text, Title } = Typography;

const SecuritySetting = () => {
  const { t } = useTranslation();
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [penaltiesLoading, setPenaltiesLoading] = useState(false);
  const [settings, setSettings] = useState({
    enable_channel_masking: true,
    mask_channel_names: true,
    mask_channel_ids: true,
    mask_channel_types: true,
    enable_anti_abuse: false,
    model_switch_window_minutes: 5,
    model_switch_threshold: 10,
    test_content_patterns: '',
    min_content_length: 10,
    test_content_threshold: 20,
    test_content_window_minutes: 5,
    abuse_score_warning_threshold: 50,
    abuse_score_action_threshold: 80,
    penalty_type: 'rate_limit',
    temp_ban_duration_minutes: 30,
    rate_limit_requests: 5,
    whitelist_user_ids: '',
    whitelist_groups: '',
  });
  const [penalties, setPenalties] = useState([]);
  const [penaltiesTotal, setPenaltiesTotal] = useState(0);
  const [penaltiesPage, setPenaltiesPage] = useState(1);

  // 获取安全设置
  const fetchSettings = async () => {
    setLoading(true);
    try {
      const res = await API.get('/api/security/settings');
      if (res.data.success && res.data.data) {
        setSettings(res.data.data);
      }
    } catch (err) {
      showError(t('获取安全设置失败'));
    } finally {
      setLoading(false);
    }
  };

  // 获取处罚列表
  const fetchPenalties = async (page = 1) => {
    setPenaltiesLoading(true);
    try {
      const res = await API.get(`/api/security/penalties?page=${page}&page_size=10`);
      if (res.data.success && res.data.data) {
        setPenalties(res.data.data.penalties || []);
        setPenaltiesTotal(res.data.data.total || 0);
      }
    } catch (err) {
      showError(t('获取处罚列表失败'));
    } finally {
      setPenaltiesLoading(false);
    }
  };

  useEffect(() => {
    fetchSettings();
    fetchPenalties();
  }, []);

  // 保存设置
  const handleSave = async () => {
    setSaving(true);
    try {
      const res = await API.put('/api/security/settings', settings);
      if (res.data.success) {
        showSuccess(t('安全设置已保存'));
      } else {
        showError(res.data.message || t('保存失败'));
      }
    } catch (err) {
      showError(t('保存安全设置失败'));
    } finally {
      setSaving(false);
    }
  };

  // 解除处罚
  const handleLiftPenalty = async (tokenId) => {
    try {
      const res = await API.delete(`/api/security/penalties/${tokenId}`);
      if (res.data.success) {
        showSuccess(t('处罚已解除'));
        fetchPenalties(penaltiesPage);
      } else {
        showError(res.data.message || t('解除处罚失败'));
      }
    } catch (err) {
      showError(t('解除处罚失败'));
    }
  };

  // 处理表单变化
  const handleChange = (key, value) => {
    setSettings(prev => ({ ...prev, [key]: value }));
  };

  // 处罚列表列定义
  const penaltyColumns = [
    {
      title: t('Token ID'),
      dataIndex: 'token_id',
      key: 'token_id',
    },
    {
      title: t('Token 名称'),
      dataIndex: 'token_name',
      key: 'token_name',
    },
    {
      title: t('处罚类型'),
      dataIndex: 'penalty_type',
      key: 'penalty_type',
      render: (type) => {
        const typeMap = {
          'temp_ban': t('临时封禁'),
          'permanent_ban': t('永久封禁'),
          'warning': t('警告'),
        };
        return <Tag color={type === 'permanent_ban' ? 'red' : 'orange'}>{typeMap[type] || type}</Tag>;
      },
    },
    {
      title: t('原因'),
      dataIndex: 'reason',
      key: 'reason',
    },
    {
      title: t('开始时间'),
      dataIndex: 'start_time',
      key: 'start_time',
      render: (time) => timestamp2string(time),
    },
    {
      title: t('结束时间'),
      dataIndex: 'end_time',
      key: 'end_time',
      render: (time) => time ? timestamp2string(time) : t('永久'),
    },
    {
      title: t('操作'),
      key: 'action',
      render: (_, record) => (
        <Popconfirm
          title={t('确认解除处罚？')}
          onConfirm={() => handleLiftPenalty(record.token_id)}
        >
          <Button size='small' type='danger'>
            {t('解除')}
          </Button>
        </Popconfirm>
      ),
    },
  ];

  if (loading) {
    return (
      <div className='flex justify-center items-center min-h-[200px]'>
        <Spin size='large' />
      </div>
    );
  }

  return (
    <div className='space-y-4'>
      {/* 基础设置 */}
      <Card>
        <Form.Section text={t('反滥用保护设置')}>
          {/* 渠道信息脱敏 */}
          <Row gutter={16}>
            <Col span={24}>
              <Title heading={6}>{t('渠道信息脱敏')}</Title>
              <Checkbox
                checked={settings.enable_channel_masking}
                onChange={(e) => handleChange('enable_channel_masking', e.target.checked)}
              >
                {t('启用渠道信息脱敏')}
              </Checkbox>
            </Col>
          </Row>
          <Row gutter={16} style={{ marginTop: 16 }}>
            <Col span={24}>
              <Checkbox
                checked={settings.mask_channel_names}
                onChange={(e) => handleChange('mask_channel_names', e.target.checked)}
                disabled={!settings.enable_channel_masking}
              >
                {t('脱敏渠道名称')}
              </Checkbox>
            </Col>
            <Col span={24}>
              <Checkbox
                checked={settings.mask_channel_ids}
                onChange={(e) => handleChange('mask_channel_ids', e.target.checked)}
                disabled={!settings.enable_channel_masking}
              >
                {t('脱敏渠道 ID')}
              </Checkbox>
            </Col>
            <Col span={24}>
              <Checkbox
                checked={settings.mask_channel_types}
                onChange={(e) => handleChange('mask_channel_types', e.target.checked)}
                disabled={!settings.enable_channel_masking}
              >
                {t('脱敏渠道类型')}
              </Checkbox>
            </Col>
          </Row>

          <Divider margin='24px' />

          {/* 反滥用检测 */}
          <Row gutter={16}>
            <Col span={24}>
              <div className='flex items-center gap-2 mb-4'>
                <Shield size={20} className='text-blue-500' />
                <Text>{t('启用反滥用检测系统，保护 API 资源不被恶意使用')}</Text>
              </div>
              <Checkbox
                checked={settings.enable_anti_abuse}
                onChange={(e) => handleChange('enable_anti_abuse', e.target.checked)}
              >
                {t('启用反滥用保护')}
              </Checkbox>
            </Col>
          </Row>

          <Divider margin='24px' />

          {/* 模型切换检测 */}
          <Row gutter={16}>
            <Col span={24}>
              <Title heading={6}>{t('模型频繁切换检测')}</Title>
            </Col>
          </Row>
          <Row gutter={16} style={{ marginTop: 16 }}>
            <Col xs={24} sm={12} md={8}>
              <Form.Label>{t('检测时间窗口（分钟）')}</Form.Label>
              <InputNumber
                value={settings.model_switch_window_minutes}
                onChange={(value) => handleChange('model_switch_window_minutes', value)}
                min={1}
                max={60}
                style={{ width: '100%' }}
                disabled={!settings.enable_anti_abuse}
              />
              <Text size='small' type='tertiary'>{t('在此时间内检测模型切换次数')}</Text>
            </Col>
            <Col xs={24} sm={12} md={8}>
              <Form.Label>{t('切换次数阈值')}</Form.Label>
              <InputNumber
                value={settings.model_switch_threshold}
                onChange={(value) => handleChange('model_switch_threshold', value)}
                min={1}
                max={100}
                style={{ width: '100%' }}
                disabled={!settings.enable_anti_abuse}
              />
              <Text size='small' type='tertiary'>{t('超过此次数将触发警告')}</Text>
            </Col>
          </Row>

          <Divider margin='24px' />

          {/* 测试内容检测 */}
          <Row gutter={16}>
            <Col span={24}>
              <Title heading={6}>{t('测试内容检测')}</Title>
            </Col>
          </Row>
          <Row gutter={16} style={{ marginTop: 16 }}>
            <Col xs={24} sm={12} md={8}>
              <Form.Label>{t('最小内容长度')}</Form.Label>
              <InputNumber
                value={settings.min_content_length}
                onChange={(value) => handleChange('min_content_length', value)}
                min={1}
                max={100}
                style={{ width: '100%' }}
                disabled={!settings.enable_anti_abuse}
              />
              <Text size='small' type='tertiary'>{t('低于此长度将被标记为短内容')}</Text>
            </Col>
            <Col xs={24} sm={12} md={8}>
              <Form.Label>{t('测试内容阈值')}</Form.Label>
              <InputNumber
                value={settings.test_content_threshold}
                onChange={(value) => handleChange('test_content_threshold', value)}
                min={1}
                max={100}
                style={{ width: '100%' }}
                disabled={!settings.enable_anti_abuse}
              />
              <Text size='small' type='tertiary'>{t('超过此次数将触发警告')}</Text>
            </Col>
            <Col xs={24} sm={12} md={8}>
              <Form.Label>{t('测试内容时间窗口（分钟）')}</Form.Label>
              <InputNumber
                value={settings.test_content_window_minutes}
                onChange={(value) => handleChange('test_content_window_minutes', value)}
                min={1}
                max={60}
                style={{ width: '100%' }}
                disabled={!settings.enable_anti_abuse}
              />
              <Text size='small' type='tertiary'>{t('在此时间内检测测试内容次数')}</Text>
            </Col>
          </Row>
          <Row gutter={16} style={{ marginTop: 16 }}>
            <Col span={24}>
              <Form.Label>{t('测试内容模式（每行一个）')}</Form.Label>
              <Form.TextArea
                value={settings.test_content_patterns}
                onChange={(value) => handleChange('test_content_patterns', value)}
                rows={4}
                placeholder='hi&#10;hello&#10;test&#10;ping&#10;你好&#10;测试'
                disabled={!settings.enable_anti_abuse}
              />
              <Text size='small' type='tertiary'>{t('匹配这些内容将被标记为测试请求')}</Text>
            </Col>
          </Row>

          <Divider margin='24px' />

          {/* 滥用评分阈值 */}
          <Row gutter={16}>
            <Col span={24}>
              <Title heading={6}>{t('滥用评分阈值')}</Title>
            </Col>
          </Row>
          <Row gutter={16} style={{ marginTop: 16 }}>
            <Col xs={24} sm={12} md={8}>
              <Form.Label>{t('警告阈值')}</Form.Label>
              <InputNumber
                value={settings.abuse_score_warning_threshold}
                onChange={(value) => handleChange('abuse_score_warning_threshold', value)}
                min={1}
                max={100}
                style={{ width: '100%' }}
                disabled={!settings.enable_anti_abuse}
              />
              <Text size='small' type='tertiary'>{t('达到此分数将记录警告')}</Text>
            </Col>
            <Col xs={24} sm={12} md={8}>
              <Form.Label>{t('处罚阈值')}</Form.Label>
              <InputNumber
                value={settings.abuse_score_action_threshold}
                onChange={(value) => handleChange('abuse_score_action_threshold', value)}
                min={1}
                max={100}
                style={{ width: '100%' }}
                disabled={!settings.enable_anti_abuse}
              />
              <Text size='small' type='tertiary'>{t('达到此分数将触发处罚')}</Text>
            </Col>
          </Row>

          <Divider margin='24px' />

          {/* 处罚设置 */}
          <Row gutter={16}>
            <Col span={24}>
              <Title heading={6}>{t('处罚设置')}</Title>
            </Col>
          </Row>
          <Row gutter={16} style={{ marginTop: 16 }}>
            <Col xs={24} sm={12} md={8}>
              <Form.Label>{t('处罚类型')}</Form.Label>
              <Form.Select
                value={settings.penalty_type}
                onChange={(value) => handleChange('penalty_type', value)}
                style={{ width: '100%' }}
                disabled={!settings.enable_anti_abuse}
              >
                <Form.Select.Option value='rate_limit'>{t('速率限制')}</Form.Select.Option>
                <Form.Select.Option value='temp_ban'>{t('临时封禁')}</Form.Select.Option>
                <Form.Select.Option value='perm_ban'>{t('永久封禁')}</Form.Select.Option>
              </Form.Select>
              <Text size='small' type='tertiary'>{t('触发处罚时的处理方式')}</Text>
            </Col>
            <Col xs={24} sm={12} md={8}>
              <Form.Label>{t('临时封禁时长（分钟）')}</Form.Label>
              <InputNumber
                value={settings.temp_ban_duration_minutes}
                onChange={(value) => handleChange('temp_ban_duration_minutes', value)}
                min={1}
                max={1440}
                style={{ width: '100%' }}
                disabled={!settings.enable_anti_abuse || settings.penalty_type !== 'temp_ban'}
              />
              <Text size='small' type='tertiary'>{t('临时封禁的时长')}</Text>
            </Col>
            <Col xs={24} sm={12} md={8}>
              <Form.Label>{t('速率限制（请求/分钟）')}</Form.Label>
              <InputNumber
                value={settings.rate_limit_requests}
                onChange={(value) => handleChange('rate_limit_requests', value)}
                min={1}
                max={1000}
                style={{ width: '100%' }}
                disabled={!settings.enable_anti_abuse || settings.penalty_type !== 'rate_limit'}
              />
              <Text size='small' type='tertiary'>{t('速率限制时的请求数')}</Text>
            </Col>
          </Row>

          <Divider margin='24px' />

          {/* 白名单设置 */}
          <Row gutter={16}>
            <Col span={24}>
              <Title heading={6}>{t('白名单设置')}</Title>
            </Col>
          </Row>
          <Row gutter={16} style={{ marginTop: 16 }}>
            <Col xs={24} sm={12}>
              <Form.Label>{t('白名单用户 ID（逗号分隔）')}</Form.Label>
              <Form.Input
                value={settings.whitelist_user_ids}
                onChange={(value) => handleChange('whitelist_user_ids', value)}
                placeholder='1,2,3'
                disabled={!settings.enable_anti_abuse}
              />
              <Text size='small' type='tertiary'>{t('这些用户不受反滥用检测限制')}</Text>
            </Col>
            <Col xs={24} sm={12}>
              <Form.Label>{t('白名单用户组（逗号分隔）')}</Form.Label>
              <Form.Input
                value={settings.whitelist_groups}
                onChange={(value) => handleChange('whitelist_groups', value)}
                placeholder='vip,admin'
                disabled={!settings.enable_anti_abuse}
              />
              <Text size='small' type='tertiary'>{t('这些用户组不受反滥用检测限制')}</Text>
            </Col>
          </Row>

          <Row style={{ marginTop: 24 }}>
            <Col span={24}>
              <Button
                theme='solid'
                type='primary'
                onClick={handleSave}
                loading={saving}
              >
                {t('保存设置')}
              </Button>
            </Col>
          </Row>
        </Form.Section>
      </Card>

      {/* 活跃处罚列表 */}
      <Card>
        <div className='flex items-center gap-2 mb-4'>
          <AlertTriangle size={20} className='text-orange-500' />
          <Title heading={5}>{t('活跃处罚列表')}</Title>
        </div>
        <Table
          columns={penaltyColumns}
          dataSource={penalties}
          loading={penaltiesLoading}
          pagination={{
            currentPage: penaltiesPage,
            pageSize: 10,
            total: penaltiesTotal,
            onPageChange: (page) => {
              setPenaltiesPage(page);
              fetchPenalties(page);
            },
          }}
          empty={t('暂无活跃处罚')}
        />
      </Card>
    </div>
  );
};

export default SecuritySetting;
