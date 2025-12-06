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

import React, { useState, useMemo } from 'react';
import { Modal, Button, Typography, Toast } from '@douyinfe/semi-ui';
import { IconAlertTriangle } from '@douyinfe/semi-icons';
import { API, showError, showSuccess } from '../../helpers';
import dayjs from 'dayjs';

const { Text, Title } = Typography;

const MAKEUP_COST = 2000; // 补签消耗额度
const MAKEUP_MAX_DAYS = 3; // 最多可补签天数

const MakeupModal = ({ 
  visible, 
  onClose, 
  onSuccess, 
  checkedDays, 
  currentYear, 
  currentMonth,
  t, 
  renderQuota,
  userDispatch,
  userState
}) => {
  const [selectedDate, setSelectedDate] = useState(null);
  const [loading, setLoading] = useState(false);

  // 计算可补签的日期（最近3天内未签到的日期）
  const availableDates = useMemo(() => {
    const today = dayjs().startOf('day');
    const dates = [];
    
    for (let i = 1; i <= MAKEUP_MAX_DAYS; i++) {
      const date = today.subtract(i, 'day');
      const day = date.date();
      const month = date.month() + 1;
      const year = date.year();
      
      // 检查是否在当前显示的月份内且未签到
      const isInCurrentMonth = year === currentYear && month === currentMonth;
      const isChecked = isInCurrentMonth && checkedDays.includes(day);
      
      // 如果不在当前月份，需要检查该日期是否已签到（这里简化处理，假设未签到）
      if (!isChecked) {
        dates.push({
          date: date.format('YYYY-MM-DD'),
          display: date.format('MM月DD日'),
          dayOfWeek: ['日', '一', '二', '三', '四', '五', '六'][date.day()],
        });
      }
    }
    
    return dates;
  }, [checkedDays, currentYear, currentMonth]);

  // 执行补签
  const handleMakeup = async () => {
    if (!selectedDate) {
      Toast.warning(t('请选择要补签的日期'));
      return;
    }

    setLoading(true);
    try {
      const res = await API.post('/api/user/checkin/makeup', {
        target_date: selectedDate,
      });
      const { success, data, message } = res.data;
      if (success) {
        showSuccess(t('补签成功') + `，${t('获得')} ${renderQuota(data.total_reward)}`);
        // 更新用户额度（扣除补签费用，加上奖励）
        if (userState?.user) {
          const updatedUser = {
            ...userState.user,
            quota: userState.user.quota - MAKEUP_COST + data.total_reward,
          };
          userDispatch({ type: 'login', payload: updatedUser });
        }
        setSelectedDate(null);
        onSuccess();
      } else {
        showError(message);
      }
    } catch (err) {
      showError(t('补签失败'));
    } finally {
      setLoading(false);
    }
  };

  const handleClose = () => {
    setSelectedDate(null);
    onClose();
  };

  return (
    <Modal
      title={t('补签')}
      visible={visible}
      onCancel={handleClose}
      footer={
        <div className='flex justify-end gap-2'>
          <Button onClick={handleClose}>{t('取消')}</Button>
          <Button 
            theme='solid' 
            type='primary' 
            onClick={handleMakeup}
            loading={loading}
            disabled={!selectedDate || availableDates.length === 0}
          >
            {t('确认补签')}
          </Button>
        </div>
      }
      centered
    >
      <div className='space-y-4'>
        {/* 补签说明 */}
        <div 
          className='p-3 rounded-lg flex items-start gap-2'
          style={{ backgroundColor: 'var(--semi-color-warning-light-default)' }}
        >
          <IconAlertTriangle style={{ color: 'var(--semi-color-warning)' }} />
          <div>
            <Text>{t('补签规则')}:</Text>
            <ul className='list-disc list-inside text-sm mt-1' style={{ color: 'var(--semi-color-text-2)' }}>
              <li>{t('只能补签最近3天内的日期')}</li>
              <li>{t('补签消耗')}: {renderQuota(MAKEUP_COST)}</li>
              <li>{t('补签奖励为正常签到的50%')}</li>
              <li>{t('补签不触发惊喜奖励')}</li>
            </ul>
          </div>
        </div>

        {/* 可补签日期列表 */}
        {availableDates.length > 0 ? (
          <div className='space-y-2'>
            <Text strong>{t('选择补签日期')}:</Text>
            <div className='grid grid-cols-1 gap-2'>
              {availableDates.map((item) => (
                <div
                  key={item.date}
                  className={`
                    p-3 rounded-lg border cursor-pointer transition-all
                    ${selectedDate === item.date 
                      ? 'border-primary bg-primary-light' 
                      : 'border-gray-200 hover:border-primary'}
                  `}
                  style={selectedDate === item.date ? {
                    borderColor: 'var(--semi-color-primary)',
                    backgroundColor: 'var(--semi-color-primary-light-default)',
                  } : {}}
                  onClick={() => setSelectedDate(item.date)}
                >
                  <div className='flex items-center justify-between'>
                    <div>
                      <Text strong>{item.display}</Text>
                      <Text type='secondary' className='ml-2'>
                        {t('星期')}{item.dayOfWeek}
                      </Text>
                    </div>
                    {selectedDate === item.date && (
                      <Text style={{ color: 'var(--semi-color-primary)' }}>✓</Text>
                    )}
                  </div>
                </div>
              ))}
            </div>
          </div>
        ) : (
          <div className='text-center py-8'>
            <Text type='secondary'>{t('最近3天内没有可补签的日期')}</Text>
          </div>
        )}

        {/* 费用提示 */}
        {selectedDate && (
          <div className='text-center pt-2'>
            <Text type='secondary'>
              {t('补签将消耗')} <Text strong style={{ color: 'var(--semi-color-warning)' }}>{renderQuota(MAKEUP_COST)}</Text>
            </Text>
          </div>
        )}
      </div>
    </Modal>
  );
};

export default MakeupModal;
