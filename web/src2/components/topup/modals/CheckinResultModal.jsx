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

import React from 'react';
import { Modal, Button, Typography } from '@douyinfe/semi-ui';
import { Sparkles, Flame } from 'lucide-react';
import '../../checkin/checkin.css';

const { Title } = Typography;

/**
 * 签到结果弹窗组件
 * @param {Object} props
 * @param {boolean} props.visible - 是否显示弹窗
 * @param {Function} props.onClose - 关闭弹窗回调
 * @param {Object} props.result - 签到结果数据
 * @param {Function} props.t - 国际化翻译函数
 * @param {Function} props.renderQuota - 额度渲染函数
 */
const CheckinResultModal = ({ visible, onClose, result, t, renderQuota }) => {
  if (!result) return null;

  return (
    <Modal
      title={null}
      visible={visible}
      onOk={onClose}
      onCancel={onClose}
      footer={null}
      centered
      className='checkin-result-modal'
      width={400}
    >
      <div className='checkin-result-content'>
        {/* 顶部装饰 */}
        <div className='result-decoration'>
          <div className='result-glow' />
          <div className='result-icon'>
            <Sparkles size={32} />
          </div>
        </div>

        <Title heading={4} className='result-title'>
          {t('签到成功')}
        </Title>

        {/* 奖励金额 */}
        <div className='result-reward'>
          <span className='reward-plus'>+</span>
          <span className='reward-amount'>{renderQuota(result.total_reward)}</span>
        </div>

        {/* 奖励明细 */}
        <div className='result-details'>
          <div className='detail-item'>
            <span className='detail-label'>{t('基础奖励')}</span>
            <span className='detail-value'>{renderQuota(result.base_reward)}</span>
          </div>
          {result.bonus_triggered && (
            <div className='detail-item bonus'>
              <span className='detail-label'>
                <Sparkles size={14} className='inline mr-1' />
                {t('惊喜奖励')}
              </span>
              <span className='detail-value bonus'>+{renderQuota(result.bonus_reward)}</span>
            </div>
          )}
        </div>

        {/* 连续签到 */}
        <div className='result-streak'>
          <Flame size={16} className='streak-icon' />
          <span>{t('连续签到')} <strong>{result.consecutive_days}</strong> {t('天')}</span>
        </div>

        <Button
          theme='solid'
          type='primary'
          size='large'
          block
          onClick={onClose}
          className='result-btn'
        >
          {t('太棒了')}
        </Button>
      </div>
    </Modal>
  );
};

export default CheckinResultModal;
