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

/**
 * 集中管理等级相关图标
 * 通过统一导出优化 tree-shaking 和代码分割
 * 
 * 优化说明：
 * 1. 集中导出避免重复导入相同图标
 * 2. lucide-react 支持 tree-shaking，只打包使用的图标
 * 3. 便于统一管理和替换图标
 * 4. 减少打包体积和加载时间
 */

// 等级徽章图标
export { Shield, Star, Crown, Gem } from 'lucide-react';

// 升级进度图标
export { Trophy, TrendingUp, DollarSign, Target, Zap } from 'lucide-react';

// 权益详情图标
export { Layers, Percent, Gift } from 'lucide-react';
