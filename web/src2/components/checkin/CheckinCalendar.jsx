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

import React, { useMemo } from 'react';
import { Button, Typography } from '@douyinfe/semi-ui';
import { IconChevronLeft, IconChevronRight } from '@douyinfe/semi-icons';
import { Check } from 'lucide-react';

const { Text } = Typography;

const CheckinCalendar = ({ year, month, checkedDays, onMonthChange, t }) => {
  // 获取月份的天数
  const getDaysInMonth = (year, month) => {
    return new Date(year, month, 0).getDate();
  };

  // 获取月份第一天是星期几 (0-6, 0是周日)
  const getFirstDayOfMonth = (year, month) => {
    return new Date(year, month - 1, 1).getDay();
  };

  // 生成日历数据
  const calendarData = useMemo(() => {
    const daysInMonth = getDaysInMonth(year, month);
    const firstDay = getFirstDayOfMonth(year, month);
    const today = new Date();
    const isCurrentMonth = today.getFullYear() === year && today.getMonth() + 1 === month;
    const currentDay = today.getDate();

    const days = [];
    
    // 填充月初空白
    for (let i = 0; i < firstDay; i++) {
      days.push({ day: null, isChecked: false, isToday: false, isPast: true });
    }
    
    // 填充日期
    for (let day = 1; day <= daysInMonth; day++) {
      const isChecked = checkedDays.includes(day);
      const isToday = isCurrentMonth && day === currentDay;
      const isPast = isCurrentMonth ? day < currentDay : 
        (year < today.getFullYear() || (year === today.getFullYear() && month < today.getMonth() + 1));
      const isFuture = isCurrentMonth ? day > currentDay :
        (year > today.getFullYear() || (year === today.getFullYear() && month > today.getMonth() + 1));
      
      days.push({ day, isChecked, isToday, isPast, isFuture });
    }

    return days;
  }, [year, month, checkedDays]);

  // 星期标题
  const weekDays = [
    t('日'), t('一'), t('二'), t('三'), t('四'), t('五'), t('六')
  ];

  // 上一个月
  const handlePrevMonth = () => {
    if (month === 1) {
      onMonthChange(year - 1, 12);
    } else {
      onMonthChange(year, month - 1);
    }
  };

  // 下一个月
  const handleNextMonth = () => {
    const today = new Date();
    const nextMonth = month === 12 ? 1 : month + 1;
    const nextYear = month === 12 ? year + 1 : year;
    
    // 不允许查看未来月份
    if (nextYear > today.getFullYear() || 
        (nextYear === today.getFullYear() && nextMonth > today.getMonth() + 1)) {
      return;
    }
    
    onMonthChange(nextYear, nextMonth);
  };

  // 判断是否可以查看下一个月
  const canGoNext = () => {
    const today = new Date();
    const nextMonth = month === 12 ? 1 : month + 1;
    const nextYear = month === 12 ? year + 1 : year;
    return !(nextYear > today.getFullYear() || 
             (nextYear === today.getFullYear() && nextMonth > today.getMonth() + 1));
  };

  return (
    <div className='checkin-calendar'>
      {/* 月份导航 */}
      <div className='flex items-center justify-between mb-4'>
        <Button
          theme='borderless'
          icon={<IconChevronLeft />}
          onClick={handlePrevMonth}
        />
        <Text strong size='large'>
          {year}{t('年')}{month}{t('月')}
        </Text>
        <Button
          theme='borderless'
          icon={<IconChevronRight />}
          onClick={handleNextMonth}
          disabled={!canGoNext()}
        />
      </div>

      {/* 星期标题 */}
      <div className='grid grid-cols-7 gap-1 mb-2'>
        {weekDays.map((day, index) => (
          <div 
            key={index} 
            className='text-center py-2 text-sm font-medium'
            style={{ color: 'var(--semi-color-text-2)' }}
          >
            {day}
          </div>
        ))}
      </div>

      {/* 日期格子 */}
      <div className='grid grid-cols-7 gap-2'>
        {calendarData.map((item, index) => (
          <div
            key={index}
            className={`calendar-day-cell ${item.day === null ? 'empty' : ''}`}
          >
            {item.day && (
              <div 
                className={`calendar-day ${item.isChecked ? 'checked' : ''} ${item.isToday ? 'today' : ''} ${item.isFuture ? 'future' : ''} ${!item.isChecked && item.isPast ? 'past' : ''}`}
              >
                <span className='day-number'>{item.day}</span>
                {item.isChecked && (
                  <span className='day-check'>
                    <Check size={12} />
                  </span>
                )}
              </div>
            )}
          </div>
        ))}
      </div>

      {/* 图例 */}
      <div className='calendar-legend'>
        <div className='legend-item'>
          <div className='legend-dot checked' />
          <Text type='secondary' size='small'>{t('已签到')}</Text>
        </div>
        <div className='legend-item'>
          <div className='legend-dot today' />
          <Text type='secondary' size='small'>{t('今天')}</Text>
        </div>
      </div>
    </div>
  );
};

export default CheckinCalendar;
