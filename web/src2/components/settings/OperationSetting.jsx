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

import React, { useState, useEffect, useCallback, Suspense, useMemo } from 'react';
import { Collapse, Spin, Button, Toast, Banner, List, Typography } from '@douyinfe/semi-ui';
import { IconRefresh } from '@douyinfe/semi-icons';
import { useTranslation } from 'react-i18next';
import { CollapsibleSettingPanel } from './CollapsibleSettingPanel';
import { settingModules } from '../../config/settingModules';
import { useSettingStore } from '../../hooks/common/useSettingStore';
import { UnsavedChangesGuard } from '../common/UnsavedChangesGuard';

const { Text } = Typography;

const OperationSetting = () => {
  const { t } = useTranslation();
  const [activeKeys, setActiveKeys] = useState(['general']);
  const [loadedModules, setLoadedModules] = useState(new Set(['general']));

  const {
    current,
    original,
    errors,
    saving,
    loadSettings,
    updateSetting,
    saveSettings,
    resetSettings,
    getDirtyKeys,
    hasUnsavedChanges,
    isModuleLoading,
  } = useSettingStore();

  // 初始加载第一个模块的数据
  useEffect(() => {
    const firstModule = settingModules[0];
    if (firstModule) {
      loadSettings(firstModule.settingKeys, firstModule.key);
    }
  }, [loadSettings]);

  // 处理面板展开/折叠
  const handleCollapseChange = useCallback(
    (keys) => {
      setActiveKeys(keys);

      // 加载新展开的模块数据
      keys.forEach((key) => {
        if (!loadedModules.has(key)) {
          const module = settingModules.find((m) => m.key === key);
          if (module) {
            loadSettings(module.settingKeys, module.key);
            setLoadedModules((prev) => new Set([...prev, key]));
          }
        }
      });
    },
    [loadedModules, loadSettings]
  );

  // 检查模块是否有未保存的更改
  const isModuleDirty = useCallback(
    (moduleKey) => {
      const module = settingModules.find((m) => m.key === moduleKey);
      if (!module) return false;
      return hasUnsavedChanges(module.settingKeys);
    },
    [hasUnsavedChanges]
  );

  // 检查模块是否有错误
  const hasModuleError = useCallback(
    (moduleKey) => {
      const module = settingModules.find((m) => m.key === moduleKey);
      if (!module) return false;
      return module.settingKeys.some((key) => errors[key]);
    },
    [errors]
  );

  // 保存所有更改
  const handleSaveAll = useCallback(async () => {
    const result = await saveSettings();
    if (result.success) {
      Toast.success(t('保存成功'));
    } else if (result.failureCount > 0) {
      Toast.error(t('部分保存失败，请检查错误信息'));
    } else {
      Toast.error(result.error || t('保存失败'));
    }
  }, [saveSettings, t]);

  // 刷新数据
  const handleRefresh = useCallback(() => {
    loadedModules.forEach((moduleKey) => {
      const module = settingModules.find((m) => m.key === moduleKey);
      if (module) {
        loadSettings(module.settingKeys, module.key);
      }
    });
  }, [loadedModules, loadSettings]);

  // 为子组件提供的 props
  const getModuleProps = useCallback(
    (moduleKey) => {
      const module = settingModules.find((m) => m.key === moduleKey);
      if (!module) return {};

      // 构建 options 对象（兼容旧组件接口）
      const options = {};
      module.settingKeys.forEach((key) => {
        if (current[key] !== undefined) {
          options[key] = current[key];
        }
      });

      return {
        options,
        refresh: handleRefresh,
        // 新接口
        store: {
          current,
          original,
          updateSetting,
          saveSettings: () => saveSettings(module.settingKeys),
          resetSettings: () => resetSettings(module.settingKeys),
          getDirtyKeys: () => getDirtyKeys(module.settingKeys),
          hasUnsavedChanges: () => hasUnsavedChanges(module.settingKeys),
          errors,
        },
      };
    },
    [
      current,
      original,
      errors,
      handleRefresh,
      updateSetting,
      saveSettings,
      resetSettings,
      getDirtyKeys,
      hasUnsavedChanges,
    ]
  );

  const hasAnyUnsavedChanges = hasUnsavedChanges();

  // 获取失败的配置项列表
  const failedItems = useMemo(() => {
    return Object.entries(errors).map(([key, error]) => ({ key, error }));
  }, [errors]);

  // 重试失败的配置项
  const handleRetryFailed = useCallback(async () => {
    const failedKeys = Object.keys(errors);
    if (failedKeys.length === 0) return;

    const result = await saveSettings(failedKeys);
    if (result.success) {
      Toast.success(t('重试成功'));
    } else if (result.failureCount > 0) {
      Toast.error(t('部分保存仍然失败'));
    }
  }, [errors, saveSettings, t]);

  return (
    <UnsavedChangesGuard hasUnsavedChanges={hasAnyUnsavedChanges}>
    <div style={{ padding: '10px 0' }}>
      {/* 错误显示区域 */}
      {failedItems.length > 0 && (
        <Banner
          type="danger"
          description={
            <div>
              <div style={{ marginBottom: '8px' }}>
                {t('以下配置项保存失败：')}
              </div>
              <List
                size="small"
                dataSource={failedItems}
                renderItem={(item) => (
                  <List.Item style={{ padding: '4px 0' }}>
                    <Text strong>{item.key}</Text>
                    <Text type="danger" style={{ marginLeft: '8px' }}>
                      {item.error}
                    </Text>
                  </List.Item>
                )}
              />
              <Button
                icon={<IconRefresh />}
                size="small"
                style={{ marginTop: '8px' }}
                onClick={handleRetryFailed}
                loading={saving}
              >
                {t('重试失败项')}
              </Button>
            </div>
          }
          style={{ marginBottom: '16px' }}
          closeIcon={null}
        />
      )}

      {/* 全局操作栏 */}
      {hasAnyUnsavedChanges && (
        <div
          style={{
            marginBottom: '16px',
            padding: '12px 16px',
            background: 'var(--semi-color-warning-light-default)',
            borderRadius: '6px',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
          }}
        >
          <span>{t('您有未保存的更改')}</span>
          <Button
            theme="solid"
            type="warning"
            loading={saving}
            onClick={handleSaveAll}
          >
            {t('保存所有更改')}
          </Button>
        </div>
      )}

      {/* 折叠面板 */}
      <Collapse
        activeKey={activeKeys}
        onChange={handleCollapseChange}
        style={{ background: 'transparent' }}
      >
        {settingModules.map((module) => {
          const ModuleComponent = module.Component;
          const moduleProps = getModuleProps(module.key);
          const isLoaded = loadedModules.has(module.key);

          return (
            <CollapsibleSettingPanel
              key={module.key}
              itemKey={module.key}
              title={t(module.title)}
              icon={module.icon}
              isDirty={isModuleDirty(module.key)}
              isLoading={isModuleLoading(module.key)}
              hasError={hasModuleError(module.key)}
            >
              {isLoaded && (
                <Suspense
                  fallback={
                    <Spin style={{ display: 'block', margin: '20px auto' }} />
                  }
                >
                  <ModuleComponent {...moduleProps} />
                </Suspense>
              )}
            </CollapsibleSettingPanel>
          );
        })}
      </Collapse>
    </div>
    </UnsavedChangesGuard>
  );
};

export default OperationSetting;
