/*
Copyright (C) 2025 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.
*/

import { useState, useCallback, useRef } from 'react';
import { API } from '../../helpers';

const MAX_RETRY_COUNT = 3;
const RETRY_DELAY = 1000; // 1秒

/**
 * useCaptcha 验证码自定义 Hook
 * 管理验证码挑战的获取和验证
 * Requirements: 8.4, 8.5
 */
const useCaptcha = () => {
  const [challenge, setChallenge] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [token, setToken] = useState(null);
  const retryCountRef = useRef(0);

  // 延迟函数
  const delay = (ms) => new Promise(resolve => setTimeout(resolve, ms));

  // 获取验证挑战（带重试逻辑）
  const fetchChallenge = useCallback(async (retryOnError = true) => {
    setLoading(true);
    setError(null);
    
    const attemptFetch = async (attempt) => {
      try {
        const res = await API.get('/api/captcha/challenge');
        const { success, data, error: apiError, message } = res.data;
        
        if (success) {
          setChallenge(data);
          retryCountRef.current = 0; // 重置重试计数
          return data;
        } else {
          const errorMsg = apiError?.message || message || '获取验证码失败';
          throw new Error(errorMsg);
        }
      } catch (err) {
        const errorMsg = err.response?.data?.error?.message || err.message || '网络错误，请检查网络连接';
        
        // 如果还有重试次数且允许重试
        if (retryOnError && attempt < MAX_RETRY_COUNT) {
          await delay(RETRY_DELAY * (attempt + 1)); // 递增延迟
          return attemptFetch(attempt + 1);
        }
        
        setError(errorMsg);
        return null;
      }
    };

    try {
      return await attemptFetch(0);
    } finally {
      setLoading(false);
    }
  }, []);

  // 验证滑块位置
  const verifyChallenge = useCallback(async (x) => {
    if (!challenge) {
      return { success: false, message: '验证会话无效' };
    }

    setLoading(true);
    setError(null);

    try {
      const res = await API.post('/api/captcha/verify', {
        session_id: challenge.session_id,
        x: x
      });
      
      const { success, data, error: apiError, message } = res.data;
      
      if (success) {
        setToken(data.token);
        return { success: true, token: data.token };
      } else {
        const errorMsg = apiError?.message || message || '验证失败';
        return { success: false, message: errorMsg };
      }
    } catch (err) {
      const errorMsg = err.response?.data?.error?.message || '验证失败，请重试';
      return { success: false, message: errorMsg };
    } finally {
      setLoading(false);
    }
  }, [challenge]);

  // 重置状态
  const reset = useCallback(() => {
    setChallenge(null);
    setLoading(false);
    setError(null);
    setToken(null);
  }, []);

  // 检查验证码是否启用
  const checkCaptchaStatus = useCallback(async () => {
    try {
      const res = await API.get('/api/captcha/status');
      const { success, data } = res.data;
      
      if (success) {
        return data;
      }
      return null;
    } catch (err) {
      return null;
    }
  }, []);

  return {
    challenge,
    loading,
    error,
    token,
    fetchChallenge,
    verifyChallenge,
    reset,
    checkCaptchaStatus
  };
};

export default useCaptcha;
