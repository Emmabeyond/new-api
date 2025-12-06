/*
Copyright (C) 2025 QuantumNous
AGPL-3.0 License
*/

import React, { useEffect, useState } from 'react';
import { Card, Form, InputNumber, Switch, Button, Spin, Typography } from '@douyinfe/semi-ui';
import { IconSetting } from '@douyinfe/semi-icons';
import { API, showError, showSuccess, renderQuota } from '../../helpers';

const { Title, Text } = Typography;

const CheckinSettings = ({ t }) => {
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [settings, setSettings] = useState(null);

  const fetchSettings = async () => {
    setLoading(true);
    try {
      const res = await API.get('/api/admin/checkin/settings');
      if (res.data.success) {
        setSettings(res.data.data);
      } else {
        showError(res.data.message);
      }
    } catch (err) {
      showError(t('获取设置失败'));
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchSettings();
  }, []);

  const handleSave = async (values) => {
    setSaving(true);
    try {
      const res = await API.put('/api/admin/checkin/settings', values);
      if (res.data.success) {
        showSuccess(t('保存成功'));
        setSettings(res.data.data);
      } else {
        showError(res.data.message);
      }
    } catch (err) {
      showError(t('保存失败'));
    } finally {
      setSaving(false);
    }
  };


  if (loading) {
    return (
      <Card>
        <div className='flex justify-center items-center py-8'>
          <Spin size='large' />
        </div>
      </Card>
    );
  }

  return (
    <Card
      title={
        <span>
          <IconSetting className='mr-2' />
          {t('签到奖励配置')}
        </span>
      }
    >
      <Form
        initValues={settings}
        onSubmit={handleSave}
        labelPosition='left'
        labelWidth={180}
      >
        <Form.Switch field='enabled' label={t('启用签到功能')} />

        <Title heading={6} className='mt-6 mb-4'>{t('基础奖励配置')}</Title>
        
        <Form.InputNumber
          field='reward_day_1_6'
          label={t('第1-6天奖励')}
          min={0}
          suffix={t('额度')}
          extraText={t('连续签到第1-6天每天获得的奖励')}
        />
        <Form.InputNumber
          field='reward_day_7'
          label={t('第7天奖励（周奖励）')}
          min={0}
          suffix={t('额度')}
          extraText={t('连续签到满7天的额外奖励')}
        />
        <Form.InputNumber
          field='reward_day_8_13'
          label={t('第8-13天奖励')}
          min={0}
          suffix={t('额度')}
        />
        <Form.InputNumber
          field='reward_day_14_29'
          label={t('第14-29天奖励')}
          min={0}
          suffix={t('额度')}
        />
        <Form.InputNumber
          field='reward_day_30'
          label={t('第30天奖励（月奖励）')}
          min={0}
          suffix={t('额度')}
          extraText={t('连续签到满30天的额外奖励')}
        />
        <Form.InputNumber
          field='reward_day_31_plus'
          label={t('第31天+奖励')}
          min={0}
          suffix={t('额度')}
        />

        <Title heading={6} className='mt-6 mb-4'>{t('惊喜奖励配置')}</Title>
        
        <Form.InputNumber
          field='bonus_min'
          label={t('惊喜奖励最小值')}
          min={0}
          suffix={t('额度')}
        />
        <Form.InputNumber
          field='bonus_max'
          label={t('惊喜奖励最大值')}
          min={0}
          suffix={t('额度')}
        />
        <div className='ml-[180px] mb-4'>
          <Text type='secondary'>{t('惊喜奖励触发概率')}: 10%</Text>
        </div>

        <Title heading={6} className='mt-6 mb-4'>{t('补签配置')}</Title>
        
        <Form.InputNumber
          field='makeup_cost'
          label={t('补签消耗额度')}
          min={0}
          suffix={t('额度')}
          extraText={t('用户补签一天需要消耗的额度')}
        />

        <div className='mt-6'>
          <Button theme='solid' type='primary' htmlType='submit' loading={saving}>
            {t('保存设置')}
          </Button>
        </div>
      </Form>
    </Card>
  );
};

export default CheckinSettings;
