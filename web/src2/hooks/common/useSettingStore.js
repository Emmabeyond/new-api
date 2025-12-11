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

import { useState, useCallback, useMemo } from 'react';
import { batchGetOptions, batchSaveOptions } from '../../services/optionService';

/**
 * 设置状态管理 Hook
 * 用于管理运营设置页面的状态，支持批量加载和保存
 */
export function useSettingStore() {
  // 原始数据（从服务器加载）
  const [original, setOriginal] = useState({});
  // 当前数据（包含用户修改）
  const [current, setCurrent] = useState({});
  // 加载状态 - 按模块 key 存储
  const [loading, setLoading] = useState({});
  // 错误信息 - 按配置项 key 存储
  const [errors, setErrors] = useState({});
  // 保存中状态
  const [saving, setSaving] = useState(false);

  /**
   * 加载指定 keys 的配置项
   * @param {string[]} keys - 配置项键名数组
   * @param {string} moduleKey - 模块标识，用于跟踪加载状态
   */
  const loadSettings = useCallback(async (keys, moduleKey) => {
    if (!keys || keys.length === 0) return;

    setLoading((prev) => ({ ...prev, [moduleKey]: true }));
    setErrors((prev) => {
      const newErrors = { ...prev };
      keys.forEach((key) => delete newErrors[key]);
      return newErrors;
    });

    try {
      const response = await batchGetOptions(keys);
      if (response.success && response.data) {
        setOriginal((prev) => ({ ...prev, ...response.data }));
        setCurrent((prev) => ({ ...prev, ...response.data }));
      }
    } catch (error) {
      console.error('Failed to load settings:', error);
    } finally {
      setLoading((prev) => ({ ...prev, [moduleKey]: false }));
    }
  }, []);

  /**
   * 更新单个配置项的值
   * @param {string} key - 配置项键名
   * @param {any} value - 新值
   */
  const updateSetting = useCallback((key, value) => {
    setCurrent((prev) => ({ ...prev, [key]: value }));
    // 清除该项的错误
    setErrors((prev) => {
      if (prev[key]) {
        const newErrors = { ...prev };
        delete newErrors[key];
        return newErrors;
      }
      return prev;
    });
  }, []);

  /**
   * 批量更新配置项
   * @param {Record<string, any>} updates - 要更新的配置项
   */
  const updateSettings = useCallback((updates) => {
    setCurrent((prev) => ({ ...prev, ...updates }));
  }, []);

  /**
   * 保存指定 keys 的配置项（只保存已修改的）
   * @param {string[]} keys - 要保存的配置项键名数组，如果不传则保存所有已修改的
   * @returns {Promise<{success: boolean, results: Array}>}
   */
  const saveSettings = useCallback(
    async (keys) => {
      const keysToSave = keys || Object.keys(current);
      const dirtyItems = keysToSave
        .filter((key) => {
          const originalValue = original[key];
          const currentValue = current[key];
          return JSON.stringify(originalValue) !== JSON.stringify(currentValue);
        })
        .map((key) => ({ key, value: current[key] }));

      if (dirtyItems.length === 0) {
        return { success: true, results: [] };
      }

      setSaving(true);
      setErrors({});

      try {
        const response = await batchSaveOptions(dirtyItems);
        if (response.success && response.data) {
          const { results } = response.data;

          // 更新错误状态
          const newErrors = {};
          results.forEach((result) => {
            if (!result.success) {
              newErrors[result.key] = result.error;
            }
          });
          setErrors(newErrors);

          // 更新 original 为成功保存的值
          const successKeys = results
            .filter((r) => r.success)
            .map((r) => r.key);
          if (successKeys.length > 0) {
            setOriginal((prev) => {
              const updated = { ...prev };
              successKeys.forEach((key) => {
                updated[key] = current[key];
              });
              return updated;
            });
          }

          return {
            success: response.data.failureCount === 0,
            results,
            successCount: response.data.successCount,
            failureCount: response.data.failureCount,
          };
        }
        return { success: false, results: [] };
      } catch (error) {
        console.error('Failed to save settings:', error);
        return { success: false, results: [], error: error.message };
      } finally {
        setSaving(false);
      }
    },
    [original, current]
  );

  /**
   * 重置指定 keys 的配置项为原始值
   * @param {string[]} keys - 要重置的配置项键名数组
   */
  const resetSettings = useCallback(
    (keys) => {
      if (!keys || keys.length === 0) return;
      setCurrent((prev) => {
        const updated = { ...prev };
        keys.forEach((key) => {
          if (original[key] !== undefined) {
            updated[key] = original[key];
          }
        });
        return updated;
      });
      // 清除相关错误
      setErrors((prev) => {
        const newErrors = { ...prev };
        keys.forEach((key) => delete newErrors[key]);
        return newErrors;
      });
    },
    [original]
  );

  /**
   * 获取已修改的配置项键名
   * @param {string[]} keys - 可选，限定检查范围
   * @returns {string[]}
   */
  const getDirtyKeys = useCallback(
    (keys) => {
      const keysToCheck = keys || Object.keys(current);
      return keysToCheck.filter((key) => {
        const originalValue = original[key];
        const currentValue = current[key];
        return JSON.stringify(originalValue) !== JSON.stringify(currentValue);
      });
    },
    [original, current]
  );

  /**
   * 检查是否有未保存的更改
   * @param {string[]} keys - 可选，限定检查范围
   * @returns {boolean}
   */
  const hasUnsavedChanges = useCallback(
    (keys) => {
      return getDirtyKeys(keys).length > 0;
    },
    [getDirtyKeys]
  );

  /**
   * 检查指定模块是否正在加载
   * @param {string} moduleKey - 模块标识
   * @returns {boolean}
   */
  const isModuleLoading = useCallback(
    (moduleKey) => {
      return !!loading[moduleKey];
    },
    [loading]
  );

  /**
   * 获取指定配置项的错误信息
   * @param {string} key - 配置项键名
   * @returns {string|undefined}
   */
  const getError = useCallback(
    (key) => {
      return errors[key];
    },
    [errors]
  );

  return {
    // State
    original,
    current,
    loading,
    errors,
    saving,

    // Actions
    loadSettings,
    updateSetting,
    updateSettings,
    saveSettings,
    resetSettings,

    // Computed
    getDirtyKeys,
    hasUnsavedChanges,
    isModuleLoading,
    getError,
  };
}
