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

import React from 'react';
import {
  IconSetting,
  IconMonitorStroked,
  IconCreditCard,
  IconShield,
  IconList,
  IconApps,
  IconMenu,
} from '@douyinfe/semi-icons';

/**
 * 运营设置模块配置
 * 每个模块定义了其标识、标题、图标、懒加载组件和所需的配置项键名
 */
export const settingModules = [
  {
    key: 'general',
    title: '通用设置',
    icon: <IconSetting />,
    Component: React.lazy(() => import('../pages/Setting/Operation/SettingsGeneral')),
    settingKeys: [
      'TopUpLink',
      'general_setting.docs_link',
      'general_setting.quota_display_type',
      'general_setting.custom_currency_symbol',
      'general_setting.custom_currency_exchange_rate',
      'QuotaPerUnit',
      'RetryTimes',
      'USDExchangeRate',
      'DisplayTokenStatEnabled',
      'DefaultCollapseSidebar',
      'DemoSiteEnabled',
      'SelfUseModeEnabled',
    ],
  },
  {
    key: 'headerNav',
    title: '顶栏模块管理',
    icon: <IconMenu />,
    Component: React.lazy(() => import('../pages/Setting/Operation/SettingsHeaderNavModules')),
    settingKeys: ['HeaderNavModules'],
  },
  {
    key: 'sidebarAdmin',
    title: '左侧边栏模块管理（管理员）',
    icon: <IconApps />,
    Component: React.lazy(() => import('../pages/Setting/Operation/SettingsSidebarModulesAdmin')),
    settingKeys: ['SidebarModulesAdmin'],
  },
  {
    key: 'sensitiveWords',
    title: '屏蔽词过滤设置',
    icon: <IconShield />,
    Component: React.lazy(() => import('../pages/Setting/Operation/SettingsSensitiveWords')),
    settingKeys: [
      'CheckSensitiveEnabled',
      'CheckSensitiveOnPromptEnabled',
      'SensitiveWords',
    ],
  },
  {
    key: 'log',
    title: '日志设置',
    icon: <IconList />,
    Component: React.lazy(() => import('../pages/Setting/Operation/SettingsLog')),
    settingKeys: ['LogConsumeEnabled'],
  },
  {
    key: 'monitoring',
    title: '监控设置',
    icon: <IconMonitorStroked />,
    Component: React.lazy(() => import('../pages/Setting/Operation/SettingsMonitoring')),
    settingKeys: [
      'ChannelDisableThreshold',
      'QuotaRemindThreshold',
      'AutomaticDisableChannelEnabled',
      'AutomaticEnableChannelEnabled',
      'AutomaticDisableKeywords',
      'monitor_setting.auto_test_channel_enabled',
      'monitor_setting.auto_test_channel_minutes',
    ],
  },
  {
    key: 'creditLimit',
    title: '额度设置',
    icon: <IconCreditCard />,
    Component: React.lazy(() => import('../pages/Setting/Operation/SettingsCreditLimit')),
    settingKeys: [
      'QuotaForNewUser',
      'PreConsumedQuota',
      'QuotaForInviter',
      'QuotaForInvitee',
      'quota_setting.enable_free_model_pre_consume',
    ],
  },
];

/**
 * 获取所有模块的配置项键名
 * @returns {string[]}
 */
export function getAllSettingKeys() {
  return settingModules.flatMap((module) => module.settingKeys);
}

/**
 * 根据模块 key 获取模块配置
 * @param {string} key - 模块标识
 * @returns {Object|undefined}
 */
export function getModuleByKey(key) {
  return settingModules.find((module) => module.key === key);
}
