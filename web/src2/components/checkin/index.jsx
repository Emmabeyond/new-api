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

import React, { useEffect, useState, useContext, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { Card, Button, Modal, Spin, Toast, Tag, Typography } from '@douyinfe/semi-ui';
import { IconGift, IconCalendar, IconHistory, IconTick } from '@douyinfe/semi-icons';
import { API, showError, showSuccess, renderQuota } from '../../helpers';
import { UserContext } from '../../context/User';
import CheckinCalendar from './CheckinCalendar';
import CheckinStats from './CheckinStats';
import CheckinHistory from './CheckinHistory';
import MakeupModal from './MakeupModal';

const { Title, Text } = Typography;

const Checkin = () => {
  const { t } = useTranslation();
  const [userState, userDispatch] = useContext(UserContext);
  
  // 签到状态
  const [stats, setStats] = useState(null);
  const [calendar, setCalendar] = useState([]);
  const [loading, setLoading] = useState(true);
  const [checkinLoading, setCheckinLoading] = useState(false);
  
  // 弹窗状态
  const [historyVisible, setHistoryVisible] = useState(false);
  const [makeupVisible, setMakeupVisible] = useState(false);
  const [resultVisible, setResultVisible] = useState(false);
  const [checkinResult, setCheckinResult] = useState(null);
  
  // 当前日期
  const [currentYear, setCurrentYear] = useState(new Date().getFullYear());
  const [currentMonth, setCurrentMonth] = useState(new Date().getMonth() + 1);

  // 获取签到统计
  const fetchStats = useCallback(async () => {
    try {
      const res = await API.get('/api/user/checkin/stats');
      const { success, data, message } = res.data;
      if (success) {
        setStats(data);
      } else {
        showError(message);
      }
    } catch (err) {
      showError(t('获取签到统计失败'));
    }
  }, [t]);

  // 获取签到日历
  const fetchCalendar = useCallback(async (year, month) => {
    try {
      const res = await API.get('/api/user/checkin/calendar', {
        params: { year, month }
      });
      const { success, data, message } = res.data;
      if (success) {
        setCalendar(data.days || []);
      } else {
        showError(message);
      }
    } catch (err) {
      showError(t('获取签到日历失败'));
    }
  }, [t]);

  // 初始化加载
  useEffect(() => {
    const init = async () => {
      setLoading(true);
      await Promise.all([
        fetchStats(),
        fetchCalendar(currentYear, currentMonth)
      ]);
      setLoading(false);
    };
    init();
  }, [fetchStats, fetchCalendar, currentYear, currentMonth]);

  // 执行签到
  const handleCheckin = async () => {
    if (stats?.checked_in_today) {
      Toast.warning(t('今日已签到'));
      return;
    }
    
    setCheckinLoading(true);
    try {
      const res = await API.post('/api/user/checkin');
      const { success, data, message } = res.data;
      if (success) {
        setCheckinResult(data);
        setResultVisible(true);
        // 刷新数据
        await fetchStats();
        await fetchCalendar(currentYear, currentMonth);
        // 更新用户额度
        if (userState.user) {
          const updatedUser = {
            ...userState.user,
            quota: userState.user.quota + data.total_reward,
          };
          userDispatch({ type: 'login', payload: updatedUser });
        }
      } else {
        showError(message);
      }
    } catch (err) {
      showError(t('签到失败'));
    } finally {
      setCheckinLoading(false);
    }
  };

  // 补签成功回调
  const handleMakeupSuccess = async () => {
    setMakeupVisible(false);
    await fetchStats();
    await fetchCalendar(currentYear, currentMonth);
  };

  // 月份切换
  const handleMonthChange = (year, month) => {
    setCurrentYear(year);
    setCurrentMonth(month);
  };

  if (loading) {
    return (
      <div className='flex justify-center items-center min-h-[400px]'>
        <Spin size='large' />
      </div>
    );
  }

  return (
    <div className='w-full max-w-4xl mx-auto space-y-6'>
      {/* 签到卡片 */}
      <Card className='checkin-main-card'>
        <div className='flex flex-col md:flex-row md:items-center md:justify-between gap-4'>
          <div className='flex-1'>
            <Title heading={4} className='mb-2'>
              <IconGift className='mr-2' />
              {t('每日签到')}
            </Title>
            <Text type='secondary'>
              {stats?.checked_in_today 
                ? t('今日已签到，明天再来吧！')
                : t('签到领取额度奖励，连续签到奖励更多！')}
            </Text>
          </div>
          <div className='flex gap-3'>
            <Button
              theme='solid'
              type='primary'
              size='large'
              icon={stats?.checked_in_today ? <IconTick /> : <IconGift />}
              loading={checkinLoading}
              disabled={stats?.checked_in_today}
              onClick={handleCheckin}
            >
              {stats?.checked_in_today ? t('已签到') : t('立即签到')}
            </Button>
            <Button
              theme='light'
              type='tertiary'
              size='large'
              icon={<IconCalendar />}
              onClick={() => setMakeupVisible(true)}
            >
              {t('补签')}
            </Button>
          </div>
        </div>
      </Card>

      {/* 统计信息 */}
      <CheckinStats stats={stats} t={t} renderQuota={renderQuota} />

      {/* 签到日历 */}
      <Card
        title={
          <div className='flex items-center justify-between w-full'>
            <span>
              <IconCalendar className='mr-2' />
              {t('签到日历')}
            </span>
            <Button
              theme='borderless'
              type='tertiary'
              icon={<IconHistory />}
              onClick={() => setHistoryVisible(true)}
            >
              {t('签到记录')}
            </Button>
          </div>
        }
      >
        <CheckinCalendar
          year={currentYear}
          month={currentMonth}
          checkedDays={calendar}
          onMonthChange={handleMonthChange}
          t={t}
        />
      </Card>

      {/* 签到结果弹窗 */}
      <Modal
        title={t('签到成功')}
        visible={resultVisible}
        onOk={() => setResultVisible(false)}
        onCancel={() => setResultVisible(false)}
        footer={
          <Button theme='solid' type='primary' onClick={() => setResultVisible(false)}>
            {t('太棒了')}
          </Button>
        }
        centered
      >
        {checkinResult && (
          <div className='text-center py-4'>
            <div className='text-6xl mb-4'>🎉</div>
            <Title heading={3} className='mb-4'>
              +{renderQuota(checkinResult.total_reward)}
            </Title>
            <div className='space-y-2'>
              <Text>{t('基础奖励')}: {renderQuota(checkinResult.base_reward)}</Text>
              {checkinResult.bonus_triggered && (
                <div>
                  <Tag color='orange' size='large'>
                    🎁 {t('惊喜奖励')}: +{renderQuota(checkinResult.bonus_reward)}
                  </Tag>
                </div>
              )}
              <div className='mt-4'>
                <Text type='secondary'>
                  {t('连续签到')}: {checkinResult.consecutive_days} {t('天')}
                </Text>
              </div>
            </div>
          </div>
        )}
      </Modal>

      {/* 签到历史弹窗 */}
      <CheckinHistory
        visible={historyVisible}
        onClose={() => setHistoryVisible(false)}
        t={t}
        renderQuota={renderQuota}
      />

      {/* 补签弹窗 */}
      <MakeupModal
        visible={makeupVisible}
        onClose={() => setMakeupVisible(false)}
        onSuccess={handleMakeupSuccess}
        checkedDays={calendar}
        currentYear={currentYear}
        currentMonth={currentMonth}
        t={t}
        renderQuota={renderQuota}
        userDispatch={userDispatch}
        userState={userState}
      />
    </div>
  );
};

export default Checkin;
