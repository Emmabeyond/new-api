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

import { useState, useEffect, useCallback } from 'react';
import { API } from '../../helpers';

// 缓存键名
const CACHE_KEYS = {
  LEVEL_INFO: 'level_data_info',
  PROGRESS: 'level_data_progress',
  ALL_LEVELS: 'level_data_all',
  TIMESTAMP: 'level_data_timestamp',
};

// 缓存过期时间（5分钟）
const CACHE_EXPIRY_MS = 5 * 60 * 1000;

/**
 * 检查缓存是否过期
 * @returns {boolean} 是否过期
 */
const isCacheExpired = () => {
  try {
    const timestamp = sessionStorage.getItem(CACHE_KEYS.TIMESTAMP);
    if (!timestamp) return true;
    
    const cacheTime = parseInt(timestamp, 10);
    const now = Date.now();
    return now - cacheTime > CACHE_EXPIRY_MS;
  } catch (error) {
    console.error('检查缓存过期时间失败:', error);
    return true;
  }
};

/**
 * 从缓存中获取数据
 * @returns {Object|null} 缓存的数据或 null
 */
const getCachedData = () => {
  try {
    if (isCacheExpired()) {
      clearCache();
      return null;
    }

    const levelInfo = sessionStorage.getItem(CACHE_KEYS.LEVEL_INFO);
    const progress = sessionStorage.getItem(CACHE_KEYS.PROGRESS);
    const allLevels = sessionStorage.getItem(CACHE_KEYS.ALL_LEVELS);

    if (!levelInfo || !progress || !allLevels) {
      return null;
    }

    return {
      levelInfo: JSON.parse(levelInfo),
      progress: JSON.parse(progress),
      allLevels: JSON.parse(allLevels),
    };
  } catch (error) {
    console.error('读取缓存数据失败:', error);
    clearCache();
    return null;
  }
};

/**
 * 将数据保存到缓存
 * @param {Object} data - 要缓存的数据
 */
const setCachedData = (data) => {
  try {
    if (data.levelInfo) {
      sessionStorage.setItem(CACHE_KEYS.LEVEL_INFO, JSON.stringify(data.levelInfo));
    }
    if (data.progress) {
      sessionStorage.setItem(CACHE_KEYS.PROGRESS, JSON.stringify(data.progress));
    }
    if (data.allLevels) {
      sessionStorage.setItem(CACHE_KEYS.ALL_LEVELS, JSON.stringify(data.allLevels));
    }
    sessionStorage.setItem(CACHE_KEYS.TIMESTAMP, Date.now().toString());
  } catch (error) {
    console.error('保存缓存数据失败:', error);
  }
};

/**
 * 清除缓存
 */
const clearCache = () => {
  try {
    Object.values(CACHE_KEYS).forEach((key) => {
      sessionStorage.removeItem(key);
    });
  } catch (error) {
    console.error('清除缓存失败:', error);
  }
};

/**
 * 用户等级数据管理 Hook
 * 提供数据加载、缓存管理和刷新功能
 * 
 * @returns {Object} 包含数据、加载状态、错误信息和刷新函数
 */
export const useLevelData = () => {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [levelInfo, setLevelInfo] = useState(null);
  const [progress, setProgress] = useState(null);
  const [allLevels, setAllLevels] = useState([]);

  /**
   * 从 API 加载数据
   */
  const loadFromAPI = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);

      // 并行加载三个 API
      const [levelRes, progressRes, allLevelsRes] = await Promise.all([
        API.get('/api/level/self'),
        API.get('/api/level/self/progress'),
        API.get('/api/level/compare'),
      ]);

      const data = {
        levelInfo: levelRes.data.success ? levelRes.data.data : null,
        progress: progressRes.data.success ? progressRes.data.data : null,
        allLevels: allLevelsRes.data.success ? allLevelsRes.data.data || [] : [],
      };

      // 更新状态
      setLevelInfo(data.levelInfo);
      setProgress(data.progress);
      setAllLevels(data.allLevels);

      // 保存到缓存
      setCachedData(data);

      return data;
    } catch (err) {
      console.error('加载等级数据失败:', err);
      setError(err);
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  /**
   * 刷新数据（强制从 API 加载）
   */
  const refresh = useCallback(async () => {
    clearCache();
    return loadFromAPI();
  }, [loadFromAPI]);

  /**
   * 初始化加载数据
   */
  useEffect(() => {
    const initializeData = async () => {
      // 先尝试从缓存加载
      const cachedData = getCachedData();
      
      if (cachedData) {
        // 使用缓存数据
        setLevelInfo(cachedData.levelInfo);
        setProgress(cachedData.progress);
        setAllLevels(cachedData.allLevels);
        setLoading(false);
      } else {
        // 缓存不存在或已过期，从 API 加载
        await loadFromAPI();
      }
    };

    initializeData();
  }, [loadFromAPI]);

  return {
    loading,
    error,
    levelInfo,
    progress,
    allLevels,
    refresh,
    clearCache,
  };
};

export default useLevelData;
