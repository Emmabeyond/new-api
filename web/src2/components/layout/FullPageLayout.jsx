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

/**
 * 全屏页面布局组件
 * 用于显示完整的 HTML 页面，不包含导航栏和其他布局元素
 */
const FullPageLayout = ({ children }) => {
  return (
    <div className='fixed inset-0 overflow-auto'>
      {children}
    </div>
  );
};

export default FullPageLayout;
