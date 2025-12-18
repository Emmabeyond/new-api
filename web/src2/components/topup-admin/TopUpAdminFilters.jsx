/*
Copyright (C) 2025 QuantumNous
*/

import React from 'react';
import { Form, Button, Space, DatePicker } from '@douyinfe/semi-ui';
import { IconSearch, IconDownload, IconRefresh } from '@douyinfe/semi-icons';

const TopUpAdminFilters = ({ filters, setFilters, onSearch, onExport, loading, t }) => {
  const statusOptions = [
    { value: '', label: t('全部状态') },
    { value: 'pending', label: t('待支付') },
    { value: 'success', label: t('成功') },
    { value: 'expired', label: t('已过期') },
  ];

  const paymentMethodOptions = [
    { value: '', label: t('全部方式') },
    { value: 'alipay', label: t('支付宝') },
    { value: 'wxpay', label: t('微信支付') },
    { value: 'stripe', label: 'Stripe' },
    { value: 'creem', label: 'Creem' },
  ];

  const handleFieldChange = (field) => (value) => {
    setFilters((prev) => ({ ...prev, [field]: value }));
  };

  const handleDateRangeChange = (dates) => {
    if (dates && dates.length === 2) {
      setFilters((prev) => ({
        ...prev,
        start_time: dates[0]?.getTime() || null,
        end_time: dates[1]?.getTime() || null,
      }));
    } else {
      setFilters((prev) => ({
        ...prev,
        start_time: null,
        end_time: null,
      }));
    }
  };

  const handleReset = () => {
    setFilters({
      keyword: '',
      status: '',
      payment_method: '',
      user_id: '',
      start_time: null,
      end_time: null,
    });
    onSearch();
  };

  return (
    <Form layout='horizontal' className='flex flex-wrap gap-2 items-end'>
      <Form.Input
        noLabel
        field='keyword_input'
        placeholder={t('搜索订单号')}
        value={filters.keyword}
        onChange={handleFieldChange('keyword')}
        style={{ width: 180 }}
        showClear
      />
      <Form.Select
        noLabel
        field='status_select'
        placeholder={t('订单状态')}
        value={filters.status}
        onChange={handleFieldChange('status')}
        optionList={statusOptions}
        style={{ width: 120 }}
      />
      <Form.Select
        noLabel
        field='payment_method_select'
        placeholder={t('支付方式')}
        value={filters.payment_method}
        onChange={handleFieldChange('payment_method')}
        optionList={paymentMethodOptions}
        style={{ width: 120 }}
      />
      <Form.Input
        noLabel
        field='user_id_input'
        placeholder={t('用户ID')}
        value={filters.user_id}
        onChange={handleFieldChange('user_id')}
        style={{ width: 100 }}
        showClear
      />
      <DatePicker
        type='dateRange'
        placeholder={[t('开始日期'), t('结束日期')]}
        onChange={handleDateRangeChange}
        style={{ width: 260 }}
      />
      <Space>
        <Button
          icon={<IconSearch />}
          theme='solid'
          onClick={onSearch}
          loading={loading}
        >
          {t('搜索')}
        </Button>
        <Button
          icon={<IconRefresh />}
          onClick={handleReset}
        >
          {t('重置')}
        </Button>
        <Button
          icon={<IconDownload />}
          onClick={onExport}
        >
          {t('导出')}
        </Button>
      </Space>
    </Form>
  );
};

export default TopUpAdminFilters;
