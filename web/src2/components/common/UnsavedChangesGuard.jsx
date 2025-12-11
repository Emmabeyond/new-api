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

import React, { useEffect } from 'react';

/**
 * 未保存更改导航守卫组件
 * 当用户有未保存的更改并尝试刷新或关闭页面时，显示浏览器原生确认对话框
 *
 * 注意：由于当前应用使用 BrowserRouter 而非 data router，
 * 无法使用 useBlocker 拦截路由导航。仅支持拦截浏览器刷新/关闭操作。
 *
 * @param {Object} props
 * @param {boolean} props.hasUnsavedChanges - 是否有未保存的更改
 * @param {React.ReactNode} props.children - 子组件
 */
export function UnsavedChangesGuard({
  hasUnsavedChanges,
  children,
}) {
  // 处理浏览器刷新/关闭
  useEffect(() => {
    const handleBeforeUnload = (e) => {
      if (hasUnsavedChanges) {
        e.preventDefault();
        e.returnValue = '';
        return '';
      }
    };

    window.addEventListener('beforeunload', handleBeforeUnload);
    return () => {
      window.removeEventListener('beforeunload', handleBeforeUnload);
    };
  }, [hasUnsavedChanges]);

  return <>{children}</>;
}

export default UnsavedChangesGuard;
