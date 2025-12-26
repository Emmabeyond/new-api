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

import { useState, useEffect, useRef, useCallback } from 'react';
import { API } from '../../helpers';

/**
 * 订单状态轮询 Hook
 * 
 * @param {string} tradeNo - 订单号
 * @param {Object} options - 配置选项
 * @param {number} options.interval - 轮询间隔（毫秒），默认 2000
 * @param {number} options.maxRetries - 最大重试次数，默认 3
 * @param {number} options.maxBackoff - 最大退避时间（毫秒），默认 30000
 * @param {number} options.expireTime - 二维码过期时间（毫秒），默认 300000 (5分钟)
 * @param {Function} options.onSuccess - 支付成功回调
 * @param {Function} options.onExpire - 二维码过期回调
 * @param {Function} options.onError - 错误回调
 * @returns {Object} { status, isPolling, error, startPolling, stopPolling }
 */
export const useOrderPolling = (tradeNo, options = {}) => {
  const {
    interval = 2000,
    maxRetries = 3,
    maxBackoff = 30000,
    expireTime = 300000,
    onSuccess,
    onExpire,
    onError,
  } = options;

  const [status, setStatus] = useState('pending'); // pending | success | failed | expired
  const [isPolling, setIsPolling] = useState(false);
  const [error, setError] = useState(null);

  const timerRef = useRef(null);
  const retryCountRef = useRef(0);
  const startTimeRef = useRef(null);
  const isPollingRef = useRef(false);

  // 计算指数退避延迟
  const getBackoffDelay = useCallback((retryCount) => {
    const delay = interval * Math.pow(2, retryCount);
    return Math.min(delay, maxBackoff);
  }, [interval, maxBackoff]);

  // 查询订单状态
  const queryOrderStatus = useCallback(async () => {
    if (!tradeNo || !isPollingRef.current) return;

    // 检查是否过期
    if (startTimeRef.current && Date.now() - startTimeRef.current >= expireTime) {
      setStatus('expired');
      setIsPolling(false);
      isPollingRef.current = false;
      onExpire?.();
      return;
    }

    try {
      const res = await API.get(`/api/user/topup/status?trade_no=${tradeNo}`);
      const { message, data } = res.data;

      if (message === 'success') {
        retryCountRef.current = 0; // 重置重试计数

        if (data.status === 'success') {
          setStatus('success');
          setIsPolling(false);
          isPollingRef.current = false;
          onSuccess?.(data.amount);
          return;
        } else if (data.status === 'failed') {
          setStatus('failed');
          setIsPolling(false);
          isPollingRef.current = false;
          return;
        }

        // 继续轮询
        if (isPollingRef.current) {
          timerRef.current = setTimeout(queryOrderStatus, interval);
        }
      } else {
        throw new Error(data || '查询订单状态失败');
      }
    } catch (err) {
      retryCountRef.current += 1;
      setError(err.message || '网络错误');

      if (retryCountRef.current >= maxRetries) {
        onError?.(err);
      }

      // 使用指数退避继续轮询
      if (isPollingRef.current) {
        const backoffDelay = getBackoffDelay(retryCountRef.current);
        timerRef.current = setTimeout(queryOrderStatus, backoffDelay);
      }
    }
  }, [tradeNo, interval, expireTime, maxRetries, getBackoffDelay, onSuccess, onExpire, onError]);

  // 开始轮询
  const startPolling = useCallback(() => {
    if (!tradeNo || isPollingRef.current) return;

    setIsPolling(true);
    setError(null);
    setStatus('pending');
    isPollingRef.current = true;
    retryCountRef.current = 0;
    startTimeRef.current = Date.now();

    // 立即执行第一次查询
    queryOrderStatus();
  }, [tradeNo, queryOrderStatus]);

  // 停止轮询
  const stopPolling = useCallback(() => {
    setIsPolling(false);
    isPollingRef.current = false;
    if (timerRef.current) {
      clearTimeout(timerRef.current);
      timerRef.current = null;
    }
  }, []);

  // 组件卸载时清理
  useEffect(() => {
    return () => {
      isPollingRef.current = false;
      if (timerRef.current) {
        clearTimeout(timerRef.current);
      }
    };
  }, []);

  // tradeNo 变化时重置状态
  useEffect(() => {
    setStatus('pending');
    setError(null);
    retryCountRef.current = 0;
  }, [tradeNo]);

  return {
    status,
    isPolling,
    error,
    startPolling,
    stopPolling,
  };
};

export default useOrderPolling;
