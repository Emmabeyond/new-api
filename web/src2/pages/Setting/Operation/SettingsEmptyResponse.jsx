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
import { Button, Col, Form, Row, Spin } from '@douyinfe/semi-ui';
import {
  compareObjects,
  API,
  showError,
  showSuccess,
  showWarning,
} from '../../../helpers';
import { useTranslation } from 'react-i18next';

export default function SettingsEmptyResponse(props) {
  const { t } = useTranslation();
  const [loading, setLoading] = useState(false);
  const [inputs, setInputs] = useState({
    'empty_response_setting.enabled': true,
    'empty_response_setting.max_retry_count': 2,
    'empty_response_setting.excluded_models': '',
    'empty_response_setting.alert_threshold': 10,
  });
  const refForm = useRef();
  const [inputsRow, setInputsRow] = useState(inputs);

  function onSubmit() {
    const updateArray = compareObjects(inputs, inputsRow);
    if (!updateArray.length) return showWarning(t('你似乎并没有修改什么'));
    const requestQueue = updateArray.map((item) => {
      let value = '';
      if (typeof inputs[item.key] === 'boolean') {
        value = String(inputs[item.key]);
      } else {
        value = String(inputs[item.key]);
      }
      return API.put('/api/option/', {
        key: item.key,
        value,
      });
    });
    setLoading(true);
    Promise.all(requestQueue)
      .then((res) => {
        if (requestQueue.length === 1) {
          if (res.includes(undefined)) return;
        } else if (requestQueue.length > 1) {
          if (res.includes(undefined))
            return showError(t('部分保存失败，请重试'));
        }
        showSuccess(t('保存成功'));
        props.refresh();
      })
      .catch(() => {
        showError(t('保存失败，请重试'));
      })
      .finally(() => {
        setLoading(false);
      });
  }

  useEffect(() => {
    const currentInputs = {};
    const booleanFields = ['empty_response_setting.enabled'];
    const numberFields = [
      'empty_response_setting.max_retry_count',
      'empty_response_setting.alert_threshold',
    ];

    for (let key in props.options) {
      if (Object.keys(inputs).includes(key)) {
        if (booleanFields.includes(key)) {
          currentInputs[key] =
            props.options[key] === 'true' || props.options[key] === true;
        } else if (numberFields.includes(key)) {
          currentInputs[key] = Number(props.options[key]) || inputs[key];
        } else {
          currentInputs[key] = props.options[key];
        }
      }
    }
    setInputs((prev) => ({ ...prev, ...currentInputs }));
    setInputsRow((prev) => structuredClone({ ...prev, ...currentInputs }));
    refForm.current.setValues(currentInputs);
  }, [props.options]);

  return (
    <>
      <Spin spinning={loading}>
        <Form
          values={inputs}
          getFormApi={(formAPI) => (refForm.current = formAPI)}
          style={{ marginBottom: 15 }}
        >
          <Form.Section text={t('空回复处理设置')}>
            <Row gutter={16}>
              <Col xs={24} sm={12} md={8} lg={8} xl={8}>
                <Form.Switch
                  field={'empty_response_setting.enabled'}
                  label={t('启用空回复检测和重试')}
                  extraText={t('检测上游渠道返回的空内容并自动重试')}
                  size='default'
                  checkedText='｜'
                  uncheckedText='〇'
                  onChange={(value) =>
                    setInputs({
                      ...inputs,
                      'empty_response_setting.enabled': value,
                    })
                  }
                />
              </Col>
              <Col xs={24} sm={12} md={8} lg={8} xl={8}>
                <Form.InputNumber
                  label={t('最大重试次数')}
                  step={1}
                  min={0}
                  max={10}
                  extraText={t('空回复时自动重试的最大次数')}
                  placeholder={'2'}
                  field={'empty_response_setting.max_retry_count'}
                  onChange={(value) =>
                    setInputs({
                      ...inputs,
                      'empty_response_setting.max_retry_count': parseInt(value) || 0,
                    })
                  }
                />
              </Col>
              <Col xs={24} sm={12} md={8} lg={8} xl={8}>
                <Form.InputNumber
                  label={t('告警阈值')}
                  step={1}
                  min={0}
                  max={100}
                  suffix={'%'}
                  extraText={t('空回复率超过此阈值时记录警告日志')}
                  placeholder={'10'}
                  field={'empty_response_setting.alert_threshold'}
                  onChange={(value) =>
                    setInputs({
                      ...inputs,
                      'empty_response_setting.alert_threshold': parseFloat(value) || 0,
                    })
                  }
                />
              </Col>
            </Row>
            <Row gutter={16}>
              <Col xs={24} sm={16}>
                <Form.TextArea
                  label={t('排除模型列表')}
                  placeholder={t('一行一个模型名称，支持前缀匹配')}
                  extraText={t(
                    '这些模型不会触发空回复检测和重试，例如：o1-preview、o1-mini',
                  )}
                  field={'empty_response_setting.excluded_models'}
                  autosize={{ minRows: 4, maxRows: 8 }}
                  onChange={(value) =>
                    setInputs({
                      ...inputs,
                      'empty_response_setting.excluded_models': value,
                    })
                  }
                />
              </Col>
            </Row>
            <Row>
              <Button size='default' onClick={onSubmit}>
                {t('保存空回复处理设置')}
              </Button>
            </Row>
          </Form.Section>
        </Form>
      </Spin>
    </>
  );
}
