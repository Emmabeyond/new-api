/*
Copyright (C) 2025 QuantumNous
*/

import { useState, useEffect, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { API, showError, showSuccess } from '../../helpers';

export const useTopUpAdminData = () => {
  const { t } = useTranslation();
  
  // Data state
  const [topups, setTopups] = useState([]);
  const [stats, setStats] = useState(null);
  const [loading, setLoading] = useState(false);
  
  // Pagination state
  const [activePage, setActivePage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const [total, setTotal] = useState(0);
  
  // Filter state
  const [filters, setFilters] = useState({
    keyword: '',
    status: '',
    payment_method: '',
    user_id: '',
    start_time: null,
    end_time: null,
  });
  
  // Detail modal state
  const [showDetail, setShowDetail] = useState(false);
  const [selectedOrder, setSelectedOrder] = useState(null);
  
  // Complete modal state
  const [showComplete, setShowComplete] = useState(false);
  const [completeOrder, setCompleteOrder] = useState(null);
  const [completeLoading, setCompleteLoading] = useState(false);

  // Build query params
  const buildQueryParams = useCallback(() => {
    const params = new URLSearchParams();
    params.append('page', activePage);
    params.append('page_size', pageSize);
    
    if (filters.keyword) params.append('keyword', filters.keyword);
    if (filters.status) params.append('status', filters.status);
    if (filters.payment_method) params.append('payment_method', filters.payment_method);
    if (filters.user_id) params.append('user_id', filters.user_id);
    if (filters.start_time) params.append('start_time', Math.floor(filters.start_time / 1000));
    if (filters.end_time) params.append('end_time', Math.floor(filters.end_time / 1000));
    
    return params.toString();
  }, [activePage, pageSize, filters]);

  // Fetch topups
  const fetchTopups = useCallback(async () => {
    setLoading(true);
    try {
      const queryString = buildQueryParams();
      const res = await API.get(`/api/user/topup?${queryString}`);
      const { success, data, message } = res.data;
      if (success) {
        setTopups(data.items || []);
        setTotal(data.total || 0);
      } else {
        showError(message);
      }
    } catch (err) {
      showError(t('获取充值记录失败'));
    } finally {
      setLoading(false);
    }
  }, [buildQueryParams, t]);

  // Fetch stats
  const fetchStats = useCallback(async () => {
    try {
      const params = new URLSearchParams();
      if (filters.keyword) params.append('keyword', filters.keyword);
      if (filters.status) params.append('status', filters.status);
      if (filters.payment_method) params.append('payment_method', filters.payment_method);
      if (filters.user_id) params.append('user_id', filters.user_id);
      if (filters.start_time) params.append('start_time', Math.floor(filters.start_time / 1000));
      if (filters.end_time) params.append('end_time', Math.floor(filters.end_time / 1000));
      
      const res = await API.get(`/api/user/topup/stats?${params.toString()}`);
      const { success, data, message } = res.data;
      if (success) {
        setStats(data);
      }
    } catch (err) {
      console.error('获取统计数据失败:', err);
    }
  }, [filters]);

  // Load data on mount and when dependencies change
  useEffect(() => {
    fetchTopups();
    fetchStats();
  }, [fetchTopups, fetchStats]);

  // Handlers
  const handlePageChange = (page) => {
    setActivePage(page);
  };

  const handlePageSizeChange = (size) => {
    setPageSize(size);
    setActivePage(1);
  };

  const handleSearch = () => {
    setActivePage(1);
    fetchTopups();
    fetchStats();
  };

  const handleExport = async () => {
    try {
      const params = new URLSearchParams();
      if (filters.keyword) params.append('keyword', filters.keyword);
      if (filters.status) params.append('status', filters.status);
      if (filters.payment_method) params.append('payment_method', filters.payment_method);
      if (filters.user_id) params.append('user_id', filters.user_id);
      if (filters.start_time) params.append('start_time', Math.floor(filters.start_time / 1000));
      if (filters.end_time) params.append('end_time', Math.floor(filters.end_time / 1000));
      
      window.open(`/api/user/topup/export?${params.toString()}`, '_blank');
    } catch (err) {
      showError(t('导出失败'));
    }
  };

  const openDetail = async (order) => {
    try {
      const res = await API.get(`/api/user/topup/${order.id}`);
      const { success, data, message } = res.data;
      if (success) {
        setSelectedOrder(data);
        setShowDetail(true);
      } else {
        showError(message);
      }
    } catch (err) {
      showError(t('获取订单详情失败'));
    }
  };

  const closeDetail = () => {
    setShowDetail(false);
    setSelectedOrder(null);
  };

  const openComplete = (order) => {
    setCompleteOrder(order);
    setShowComplete(true);
  };

  const closeComplete = () => {
    setShowComplete(false);
    setCompleteOrder(null);
  };

  const handleComplete = async () => {
    if (!completeOrder) return;
    
    setCompleteLoading(true);
    try {
      const res = await API.post('/api/user/topup/complete', {
        trade_no: completeOrder.trade_no,
      });
      const { success, message } = res.data;
      if (success) {
        showSuccess(t('补单成功'));
        closeComplete();
        fetchTopups();
        fetchStats();
      } else {
        showError(message);
      }
    } catch (err) {
      showError(t('补单失败'));
    } finally {
      setCompleteLoading(false);
    }
  };

  return {
    // Data
    topups,
    stats,
    loading,
    
    // Pagination
    activePage,
    pageSize,
    total,
    handlePageChange,
    handlePageSizeChange,
    
    // Filters
    filters,
    setFilters,
    handleSearch,
    handleExport,
    
    // Detail modal
    showDetail,
    selectedOrder,
    closeDetail,
    openDetail,
    
    // Complete modal
    showComplete,
    completeOrder,
    closeComplete,
    openComplete,
    handleComplete,
    completeLoading,
    
    // Translation
    t,
    
    // Refresh
    refresh: fetchTopups,
  };
};
