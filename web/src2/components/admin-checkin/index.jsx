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

import React, { useEffect, useState, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { 
  Card, 
  Table, 
  Button, 
  Input, 
  DatePicker, 
  Select, 
  Tag, 
  Modal, 
  InputNumber,
  Spin,
  Typography,
  Tabs,
  TabPane
} from '@douyinfe/semi-ui';
import { IconSearch, IconRefresh, IconEdit, IconSetting } from '@douyinfe/semi-icons';
import { API, showError, showSuccess, renderQuota } from '../../helpers';
import dayjs from 'dayjs';
import DashboardCards from './DashboardCards';
import CheckinSettings from './CheckinSettings';

const { Title, Text } = Typography;

const AdminCheckin = () => {
  const { t } = useTranslation();
  
  // 仪表盘数据
  const [dashboard, setDashboard] = useState(null);
  const [dashboardLoading, setDashboardLoading] = useState(true);
  
  // 签到记录
  const [records, setRecords] = useState([]);
  const [recordsLoading, setRecordsLoading] = useState(false);
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 20,
    total: 0,
  });
  
  // 筛选条件
  const [filters, setFilters] = useState({
    userId: '',
    startDate: null,
    endDate: null,
    checkinType: 0,
  });
  
  // 调整连续天数弹窗
  const [adjustVisible, setAdjustVisible] = useState(false);
  const [adjustLoading, setAdjustLoading] = useState(false);
  const [adjustData, setAdjustData] = useState({
    userId: null,
    currentDays: 0,
    newDays: 0,
  });

  // 获取仪表盘数据
  const fetchDashboard = useCallback(async () => {
    setDashboardLoading(true);
    try {
      const res = await API.get('/api/admin/checkin/dashboard');
      const { success, data, message } = res.data;
      if (success) {
        setDashboard(data);
      } else {
        showError(message);
      }
    } catch (err) {
      showError(t('获取仪表盘数据失败'));
    } finally {
      setDashboardLoading(false);
    }
  }, [t]);

  // 获取签到记录
  const fetchRecords = useCallback(async (page = 1, pageSize = 20) => {
    setRecordsLoading(true);
    try {
      const params = {
        page,
        page_size: pageSize,
      };
      
      if (filters.userId) {
        params.user_id = filters.userId;
      }
      if (filters.startDate) {
        params.start_date = dayjs(filters.startDate).format('YYYY-MM-DD');
      }
      if (filters.endDate) {
        params.end_date = dayjs(filters.endDate).format('YYYY-MM-DD');
      }
      if (filters.checkinType > 0) {
        params.checkin_type = filters.checkinType;
      }
      
      const res = await API.get('/api/admin/checkin/records', { params });
      const { success, data, message } = res.data;
      if (success) {
        setRecords(data.items || []);
        setPagination({
          current: data.page || page,
          pageSize: data.page_size || pageSize,
          total: data.total || 0,
        });
      } else {
        showError(message);
      }
    } catch (err) {
      showError(t('获取签到记录失败'));
    } finally {
      setRecordsLoading(false);
    }
  }, [filters, t]);

  // 初始化加载
  useEffect(() => {
    fetchDashboard();
    fetchRecords();
  }, []);

  // 搜索
  const handleSearch = () => {
    fetchRecords(1, pagination.pageSize);
  };

  // 重置筛选
  const handleReset = () => {
    setFilters({
      userId: '',
      startDate: null,
      endDate: null,
      checkinType: 0,
    });
    fetchRecords(1, pagination.pageSize);
  };

  // 分页变化
  const handlePageChange = (page) => {
    fetchRecords(page, pagination.pageSize);
  };

  // 打开调整弹窗
  const handleOpenAdjust = async (userId) => {
    try {
      const res = await API.get(`/api/admin/checkin/user/${userId}`);
      const { success, data, message } = res.data;
      if (success) {
        setAdjustData({
          userId,
          currentDays: data.checkin?.consecutive_days || 0,
          newDays: data.checkin?.consecutive_days || 0,
        });
        setAdjustVisible(true);
      } else {
        showError(message);
      }
    } catch (err) {
      showError(t('获取用户签到信息失败'));
    }
  };

  // 提交调整
  const handleAdjust = async () => {
    if (adjustData.newDays < 0) {
      showError(t('连续天数不能为负数'));
      return;
    }
    
    setAdjustLoading(true);
    try {
      const res = await API.put(`/api/admin/checkin/user/${adjustData.userId}/consecutive`, {
        consecutive_days: adjustData.newDays,
      });
      const { success, message } = res.data;
      if (success) {
        showSuccess(t('调整成功'));
        setAdjustVisible(false);
        fetchRecords(pagination.current, pagination.pageSize);
      } else {
        showError(message);
      }
    } catch (err) {
      showError(t('调整失败'));
    } finally {
      setAdjustLoading(false);
    }
  };

  // 表格列定义
  const columns = [
    {
      title: t('用户ID'),
      dataIndex: 'user_id',
      key: 'user_id',
      width: 100,
    },
    {
      title: t('签到日期'),
      dataIndex: 'checkin_date',
      key: 'checkin_date',
      render: (text) => dayjs(text).format('YYYY-MM-DD'),
    },
    {
      title: t('类型'),
      dataIndex: 'checkin_type',
      key: 'checkin_type',
      render: (type) => (
        <Tag color={type === 1 ? 'green' : 'orange'}>
          {type === 1 ? t('正常签到') : t('补签')}
        </Tag>
      ),
    },
    {
      title: t('基础奖励'),
      dataIndex: 'base_reward',
      key: 'base_reward',
      render: (value) => renderQuota(value),
    },
    {
      title: t('惊喜奖励'),
      dataIndex: 'bonus_reward',
      key: 'bonus_reward',
      render: (value) => value > 0 ? (
        <Tag color='orange'>+{renderQuota(value)}</Tag>
      ) : '-',
    },
    {
      title: t('总计'),
      key: 'total',
      render: (_, record) => renderQuota(record.base_reward + record.bonus_reward),
    },
    {
      title: t('操作'),
      key: 'action',
      width: 100,
      render: (_, record) => (
        <Button
          theme='borderless'
          type='tertiary'
          icon={<IconEdit />}
          onClick={() => handleOpenAdjust(record.user_id)}
        >
          {t('调整')}
        </Button>
      ),
    },
  ];

  return (
    <div className='w-full max-w-7xl mx-auto space-y-6'>
      <Tabs type='line'>
        <TabPane tab={t('签到管理')} itemKey='records'>
          {/* 仪表盘 */}
          <DashboardCards 
            dashboard={dashboard} 
            loading={dashboardLoading} 
            t={t} 
            renderQuota={renderQuota}
            onRefresh={fetchDashboard}
          />

          {/* 签到记录 */}
          <Card
        title={t('签到记录')}
        headerExtraContent={
          <Button
            theme='borderless'
            icon={<IconRefresh />}
            onClick={() => fetchRecords(pagination.current, pagination.pageSize)}
          >
            {t('刷新')}
          </Button>
        }
      >
        {/* 筛选条件 */}
        <div className='flex flex-wrap gap-4 mb-4'>
          <Input
            placeholder={t('用户ID')}
            value={filters.userId}
            onChange={(value) => setFilters({ ...filters, userId: value })}
            style={{ width: 120 }}
          />
          <DatePicker
            placeholder={t('开始日期')}
            value={filters.startDate}
            onChange={(date) => setFilters({ ...filters, startDate: date })}
            style={{ width: 150 }}
          />
          <DatePicker
            placeholder={t('结束日期')}
            value={filters.endDate}
            onChange={(date) => setFilters({ ...filters, endDate: date })}
            style={{ width: 150 }}
          />
          <Select
            placeholder={t('签到类型')}
            value={filters.checkinType}
            onChange={(value) => setFilters({ ...filters, checkinType: value })}
            style={{ width: 120 }}
          >
            <Select.Option value={0}>{t('全部')}</Select.Option>
            <Select.Option value={1}>{t('正常签到')}</Select.Option>
            <Select.Option value={2}>{t('补签')}</Select.Option>
          </Select>
          <Button theme='solid' icon={<IconSearch />} onClick={handleSearch}>
            {t('搜索')}
          </Button>
          <Button theme='light' onClick={handleReset}>
            {t('重置')}
          </Button>
        </div>

        {/* 表格 */}
        <Table
          columns={columns}
          dataSource={records}
          rowKey='id'
          loading={recordsLoading}
          pagination={{
            current: pagination.current,
            pageSize: pagination.pageSize,
            total: pagination.total,
            onChange: handlePageChange,
            showSizeChanger: false,
          }}
        />
      </Card>

          {/* 调整连续天数弹窗 */}
          <Modal
            title={t('调整连续签到天数')}
            visible={adjustVisible}
            onOk={handleAdjust}
            onCancel={() => setAdjustVisible(false)}
            confirmLoading={adjustLoading}
            centered
          >
            <div className='space-y-4'>
              <div>
                <Text type='secondary'>{t('用户ID')}: </Text>
                <Text strong>{adjustData.userId}</Text>
              </div>
              <div>
                <Text type='secondary'>{t('当前连续天数')}: </Text>
                <Text strong>{adjustData.currentDays}</Text>
              </div>
              <div>
                <Text>{t('新的连续天数')}: </Text>
                <InputNumber
                  value={adjustData.newDays}
                  onChange={(value) => setAdjustData({ ...adjustData, newDays: value })}
                  min={0}
                  style={{ width: 120 }}
                />
              </div>
            </div>
          </Modal>
        </TabPane>

        <TabPane tab={<span><IconSetting className='mr-1' />{t('奖励配置')}</span>} itemKey='settings'>
          <div className='mt-4'>
            <CheckinSettings t={t} />
          </div>
        </TabPane>
      </Tabs>
    </div>
  );
};

export default AdminCheckin;
