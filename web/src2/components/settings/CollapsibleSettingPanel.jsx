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

import React, { Suspense } from 'react';
import { Collapse, Spin, Tag, Typography } from '@douyinfe/semi-ui';
import { IconAlertCircle } from '@douyinfe/semi-icons';

const { Text } = Typography;

/**
 * 可折叠设置面板组件
 * @param {Object} props
 * @param {string} props.itemKey - 面板唯一标识
 * @param {string} props.title - 面板标题
 * @param {React.ReactNode} props.icon - 面板图标
 * @param {boolean} props.isDirty - 是否有未保存的更改
 * @param {boolean} props.isLoading - 是否正在加载
 * @param {boolean} props.hasError - 是否有错误
 * @param {React.ReactNode} props.children - 面板内容
 */
export function CollapsibleSettingPanel({
  itemKey,
  title,
  icon,
  isDirty = false,
  isLoading = false,
  hasError = false,
  children,
}) {
  const header = (
    <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
      {icon}
      <Text strong>{title}</Text>
      {isDirty && (
        <Tag color="orange" size="small" style={{ marginLeft: '8px' }}>
          未保存
        </Tag>
      )}
      {hasError && (
        <IconAlertCircle style={{ color: 'var(--semi-color-danger)', marginLeft: '4px' }} />
      )}
    </div>
  );

  return (
    <Collapse.Panel itemKey={itemKey} header={header}>
      <Suspense fallback={<Spin style={{ display: 'block', margin: '20px auto' }} />}>
        {isLoading ? (
          <Spin style={{ display: 'block', margin: '20px auto' }} />
        ) : (
          children
        )}
      </Suspense>
    </Collapse.Panel>
  );
}

export default CollapsibleSettingPanel;
