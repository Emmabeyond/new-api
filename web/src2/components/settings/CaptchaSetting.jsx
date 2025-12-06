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
} from '@douyinfe/semi-ui';
import { API, showError, showSuccess } from '../../helpers';
import { useTranslation } from 'react-i18next';

const CaptchaSetting = () => {
  const { t } = useTranslation();
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [settings, setSettings] = useState({
    enabled: false,
    tolerance_range: 5,
    require_on_login: true,
    require_on_register: true,
    require_on_checkin: true,
    max_attempts: 5,
    block_duration: 300,
  });

  // 获取验证码设置
  const fetchSettings = async () => {
    setLoading(true);
    try {
      const res = await API.get('/api/captcha/settings');
      if (res.data.success && res.data.data) {
        setSettings(res.data.data);
      }
    } catch (err) {
      showError(t('获取验证码设置失败'));
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchSettings();
  }, []);

  // 保存设置
  const handleSave = async () => {
    setSaving(true);
    try {
      const res = await API.put('/api/captcha/settings', settings);
      if (res.data.success) {
        showSuccess(t('验证码设置已保存'));
      } else {
        showError(res.data.message || t('保存失败'));
      }
    } catch (err) {
      showError(t('保存验证码设置失败'));
    } finally {
      setSaving(false);
    }
  };

  // 处理表单变化
  const handleChange = (key, value) => {
    setSettings(prev => ({ ...prev, [key]: value }));
  };

  if (loading) {
    return (
      <div className='flex justify-center items-center min-h-[200px]'>
        <Spin size='large' />
      </div>
    );
  }

  return (
    <Card>
      <Form.Section text={t('滑块验证码设置')}>
        <Row gutter={16}>
          <Col span={24}>
            <Checkbox
              checked={settings.enabled}
              onChange={(e) => handleChange('enabled', e.target.checked)}
            >
              {t('启用滑块验证码')}
            </Checkbox>
          </Col>
        </Row>

        <Row gutter={16} style={{ marginTop: 16 }}>
          <Col xs={24} sm={12} md={8}>
            <Form.Label>{t('容差范围（像素）')}</Form.Label>
            <InputNumber
              value={settings.tolerance_range}
              onChange={(value) => handleChange('tolerance_range', value)}
              min={1}
              max={20}
              style={{ width: '100%' }}
              disabled={!settings.enabled}
            />
          </Col>
          <Col xs={24} sm={12} md={8}>
            <Form.Label>{t('最大尝试次数')}</Form.Label>
            <InputNumber
              value={settings.max_attempts}
              onChange={(value) => handleChange('max_attempts', value)}
              min={1}
              max={20}
              style={{ width: '100%' }}
              disabled={!settings.enabled}
            />
          </Col>
          <Col xs={24} sm={12} md={8}>
            <Form.Label>{t('封禁时长（秒）')}</Form.Label>
            <InputNumber
              value={settings.block_duration}
              onChange={(value) => handleChange('block_duration', value)}
              min={60}
              max={3600}
              style={{ width: '100%' }}
              disabled={!settings.enabled}
            />
          </Col>
        </Row>

        <Row gutter={16} style={{ marginTop: 16 }}>
          <Col span={24}>
            <Form.Label>{t('需要验证的操作')}</Form.Label>
          </Col>
          <Col span={24}>
            <Checkbox
              checked={settings.require_on_login}
              onChange={(e) => handleChange('require_on_login', e.target.checked)}
              disabled={!settings.enabled}
            >
              {t('登录')}
            </Checkbox>
          </Col>
          <Col span={24}>
            <Checkbox
              checked={settings.require_on_register}
              onChange={(e) => handleChange('require_on_register', e.target.checked)}
              disabled={!settings.enabled}
            >
              {t('注册')}
            </Checkbox>
          </Col>
          <Col span={24}>
            <Checkbox
              checked={settings.require_on_checkin}
              onChange={(e) => handleChange('require_on_checkin', e.target.checked)}
              disabled={!settings.enabled}
            >
              {t('签到')}
            </Checkbox>
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
  );
};

export default CaptchaSetting;
