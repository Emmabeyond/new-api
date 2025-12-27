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

import React, { useEffect, useState, useRef } from 'react';
import {
  Banner,
  Button,
  Form,
  Row,
  Col,
  Typography,
  Spin,
} from '@douyinfe/semi-ui';
const { Text } = Typography;
import { API, removeTrailingSlash, showError, showSuccess } from '../../../helpers';
import { useTranslation } from 'react-i18next';

export default function SettingsPaymentGatewayLinuxDo(props) {
  const { t } = useTranslation();
  const [loading, setLoading] = useState(false);
  const [inputs, setInputs] = useState({
    LinuxDoPayAddress: 'https://credit.linux.do/epay',
    LinuxDoClientId: '',
    LinuxDoClientSecret: '',
    LinuxDoMinTopUp: 1,
    LinuxDoUnitPrice: 1.0,
  });
  const [originInputs, setOriginInputs] = useState({});
  const formApiRef = useRef(null);

  useEffect(() => {
    if (props.options && formApiRef.current) {
      const currentInputs = {
        LinuxDoPayAddress: props.options.LinuxDoPayAddress || 'https://credit.linux.do/epay',
        LinuxDoClientId: props.options.LinuxDoClientId || '',
        LinuxDoClientSecret: props.options.LinuxDoClientSecret || '',
        LinuxDoMinTopUp:
          props.options.LinuxDoMinTopUp !== undefined
            ? parseInt(props.options.LinuxDoMinTopUp)
            : 1,
        LinuxDoUnitPrice:
          props.options.LinuxDoUnitPrice !== undefined
            ? parseFloat(props.options.LinuxDoUnitPrice)
            : 1.0,
      };
      setInputs(currentInputs);
      setOriginInputs({ ...currentInputs });
      formApiRef.current.setValues(currentInputs);
    }
  }, [props.options]);

  const handleFormChange = (values) => {
    setInputs(values);
  };

  const submitLinuxDoSetting = async () => {
    setLoading(true);
    try {
      const options = [];

      // 支付地址
      options.push({
        key: 'LinuxDoPayAddress',
        value: removeTrailingSlash(inputs.LinuxDoPayAddress),
      });

      // Client ID
      if (inputs.LinuxDoClientId && inputs.LinuxDoClientId !== '') {
        options.push({ key: 'LinuxDoClientId', value: inputs.LinuxDoClientId });
      }

      // Client Secret
      if (inputs.LinuxDoClientSecret && inputs.LinuxDoClientSecret !== '') {
        options.push({
          key: 'LinuxDoClientSecret',
          value: inputs.LinuxDoClientSecret,
        });
      }

      // 最小充值金额
      options.push({
        key: 'LinuxDoMinTopUp',
        value: inputs.LinuxDoMinTopUp.toString(),
      });

      // 充值比例
      if (inputs.LinuxDoUnitPrice !== undefined && inputs.LinuxDoUnitPrice !== null) {
        options.push({
          key: 'LinuxDoUnitPrice',
          value: inputs.LinuxDoUnitPrice.toString(),
        });
      }

      // 发送请求
      const requestQueue = options.map((opt) =>
        API.put('/api/option/', {
          key: opt.key,
          value: opt.value,
        }),
      );

      const results = await Promise.all(requestQueue);

      // 检查所有请求是否成功
      const errorResults = results.filter((res) => !res.data.success);
      if (errorResults.length > 0) {
        errorResults.forEach((res) => {
          showError(res.data.message);
        });
      } else {
        showSuccess(t('更新成功'));
        // 更新本地存储的原始值
        setOriginInputs({ ...inputs });
        props.refresh?.();
      }
    } catch (error) {
      showError(t('更新失败'));
    }
    setLoading(false);
  };

  return (
    <Spin spinning={loading}>
      <Form
        initValues={inputs}
        onValueChange={handleFormChange}
        getFormApi={(api) => (formApiRef.current = api)}
      >
        <Form.Section text={t('LINUX DO Credit 设置')}>
          <Text>
            {t('LINUX DO Credit 是专为 LINUX DO 社区打造的积分流通基础设施，兼容易支付协议。')}
            <a
              href='https://credit.linux.do/'
              target='_blank'
              rel='noreferrer'
            >
              {' '}LINUX DO Credit 官网
            </a>
            <br />
          </Text>
          <Banner
            type='info'
            description={
              <>
                {t('配置步骤：')}
                <ol style={{ margin: '8px 0', paddingLeft: '20px' }}>
                  <li>{t('前往 LINUX DO Credit 集市中心创建应用')}</li>
                  <li>{t('填写应用名称、主页和回调地址（回调地址为：{ServerAddress}/api/user/linuxdo/notify）')}</li>
                  <li>{t('在 API 配置面板获取 Client ID 和 Client Secret')}</li>
                  <li>{t('将获取的凭证填入下方配置')}</li>
                </ol>
              </>
            }
          />

          <Row gutter={{ xs: 8, sm: 16, md: 24, lg: 24, xl: 24, xxl: 24 }}>
            <Col xs={24} sm={24} md={8} lg={8} xl={8}>
              <Form.Input
                field='LinuxDoPayAddress'
                label={t('支付地址')}
                placeholder='https://credit.linux.do/epay'
                extraText={t('LINUX DO Credit 支付接口地址，通常无需修改')}
              />
            </Col>
            <Col xs={24} sm={24} md={8} lg={8} xl={8}>
              <Form.InputNumber
                field='LinuxDoUnitPrice'
                label={t('充值比例（积分/美元）')}
                placeholder='1'
                min={0.01}
                precision={2}
                extraText={t('1美元=多少积分，例如：7.3 表示充值1美元需支付7.3积分')}
              />
            </Col>
            <Col xs={24} sm={24} md={8} lg={8} xl={8}>
              <Form.InputNumber
                field='LinuxDoMinTopUp'
                label={t('最小充值金额')}
                placeholder='1'
                min={1}
                extraText={t('用户最小充值金额（美元）')}
              />
            </Col>
          </Row>

          <Row
            gutter={{ xs: 8, sm: 16, md: 24, lg: 24, xl: 24, xxl: 24 }}
            style={{ marginTop: 16 }}
          >
            <Col xs={24} sm={24} md={12} lg={12} xl={12}>
              <Form.Input
                field='LinuxDoClientId'
                label={t('Client ID')}
                placeholder={t('从 LINUX DO Credit 获取的 Client ID')}
              />
            </Col>
            <Col xs={24} sm={24} md={12} lg={12} xl={12}>
              <Form.Input
                field='LinuxDoClientSecret'
                label={t('Client Secret')}
                placeholder={t('敏感信息不会发送到前端显示')}
                type='password'
              />
            </Col>
          </Row>

          <Button onClick={submitLinuxDoSetting} style={{ marginTop: 16 }}>
            {t('更新 LINUX DO Credit 设置')}
          </Button>
        </Form.Section>
      </Form>
    </Spin>
  );
}

