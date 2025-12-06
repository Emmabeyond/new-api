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
import { Sparkles, Crown, Flame } from 'lucide-react';
import { API, showError, showSuccess, renderQuota } from '../../helpers';
import { UserContext } from '../../context/User';
import CheckinCalendar from './CheckinCalendar';
import CheckinStats from './CheckinStats';
import CheckinHistory from './CheckinHistory';
import CheckinRules from './CheckinRules';
import MakeupModal from './MakeupModal';
import './checkin.css';

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
      {/* 签到主卡片 - 高级渐变风格 */}
      <div className='checkin-hero-card'>
        {/* 粒子背景效果 */}
        <div className='checkin-particles'>
          {[...Array(20)].map((_, i) => (
            <div key={i} className='particle' style={{
              left: `${Math.random() * 100}%`,
              animationDelay: `${Math.random() * 5}s`,
              animationDuration: `${3 + Math.random() * 4}s`
            }} />
          ))}
        </div>
        
        {/* 装饰光晕 */}
        <div className='checkin-glow checkin-glow-1' />
        <div className='checkin-glow checkin-glow-2' />
        
        <div className='checkin-hero-content'>
          {/* 左侧信息区 */}
          <div className='checkin-hero-info'>
            <div className='checkin-hero-badge'>
              <Crown size={14} />
              <span>{t('每日福利')}</span>
            </div>
            <Title heading={2} className='checkin-hero-title'>
              {t('每日签到')}
            </Title>
            <Text className='checkin-hero-desc'>
              {stats?.checked_in_today 
                ? t('今日已签到，明天再来吧！')
                : t('签到领取额度奖励，连续签到奖励更多！')}
            </Text>
            
            {/* 连续签到展示 */}
            {stats?.consecutive_days > 0 && (
              <div className='checkin-streak-badge'>
                <Flame size={16} className='streak-icon' />
                <span>{t('连续签到')} <strong>{stats.consecutive_days}</strong> {t('天')}</span>
              </div>
            )}
          </div>
          
          {/* 右侧按钮区 */}
          <div className='checkin-hero-actions'>
            <button
              className={`checkin-main-btn ${stats?.checked_in_today ? 'checked' : ''}`}
              disabled={stats?.checked_in_today || checkinLoading}
              onClick={handleCheckin}
            >
              {checkinLoading ? (
                <Spin size='small' />
              ) : stats?.checked_in_today ? (
                <>
                  <IconTick size='large' />
                  <span>{t('已签到')}</span>
                </>
              ) : (
                <>
                  <Sparkles size={24} />
                  <span>{t('立即签到')}</span>
                </>
              )}
            </button>
            <Button
              theme='light'
              type='tertiary'
              size='large'
              icon={<IconCalendar />}
              onClick={() => setMakeupVisible(true)}
              className='checkin-makeup-btn'
            >
              {t('补签')}
            </Button>
          </div>
        </div>
      </div>

      {/* 签到规则 */}
      <CheckinRules t={t} renderQuota={renderQuota} />

      {/* 统计信息 */}
      <CheckinStats stats={stats} t={t} renderQuota={renderQuota} />

      {/* 签到日历 */}
      <div className='checkin-calendar-card'>
        <div className='calendar-card-header'>
          <div className='calendar-card-title'>
            <IconCalendar className='mr-2' style={{ color: '#8b5cf6' }} />
            <span>{t('签到日历')}</span>
          </div>
          <Button
            theme='borderless'
            type='tertiary'
            icon={<IconHistory />}
            onClick={() => setHistoryVisible(true)}
          >
            {t('签到记录')}
          </Button>
        </div>
        <div className='calendar-card-body'>
          <CheckinCalendar
            year={currentYear}
            month={currentMonth}
            checkedDays={calendar}
            onMonthChange={handleMonthChange}
            t={t}
          />
        </div>
      </div>

      {/* 签到结果弹窗 - 高级风格 */}
      <Modal
        title={null}
        visible={resultVisible}
        onOk={() => setResultVisible(false)}
        onCancel={() => setResultVisible(false)}
        footer={null}
        centered
        className='checkin-result-modal'
        width={400}
      >
        {checkinResult && (
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
              <span className='reward-amount'>{renderQuota(checkinResult.total_reward)}</span>
            </div>
            
            {/* 奖励明细 */}
            <div className='result-details'>
              <div className='detail-item'>
                <span className='detail-label'>{t('基础奖励')}</span>
                <span className='detail-value'>{renderQuota(checkinResult.base_reward)}</span>
              </div>
              {checkinResult.bonus_triggered && (
                <div className='detail-item bonus'>
                  <span className='detail-label'>
                    <Sparkles size={14} className='inline mr-1' />
                    {t('惊喜奖励')}
                  </span>
                  <span className='detail-value bonus'>+{renderQuota(checkinResult.bonus_reward)}</span>
                </div>
              )}
            </div>
            
            {/* 连续签到 */}
            <div className='result-streak'>
              <Flame size={16} className='streak-icon' />
              <span>{t('连续签到')} <strong>{checkinResult.consecutive_days}</strong> {t('天')}</span>
            </div>
            
            <Button 
              theme='solid' 
              type='primary' 
              size='large'
              block
              onClick={() => setResultVisible(false)}
              className='result-btn'
            >
              {t('太棒了')}
            </Button>
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
