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

import { API } from '../helpers';

/**
 * 批量获取配置项
 * @param {string[]} keys - 配置项键名数组
 * @returns {Promise<{success: boolean, data: Record<string, any>}>}
 */
export async function batchGetOptions(keys) {
  if (!keys || keys.length === 0) {
    return { success: true, data: {} };
  }
  const response = await API.get('/api/option/batch', {
    params: { keys: keys.join(',') },
  });
  return response.data;
}

/**
 * 批量保存配置项
 * @param {Array<{key: string, value: any}>} options - 配置项数组
 * @returns {Promise<{success: boolean, data: {results: Array, successCount: number, failureCount: number}}>}
 */
export async function batchSaveOptions(options) {
  if (!options || options.length === 0) {
    return { success: true, data: { results: [], successCount: 0, failureCount: 0 } };
  }
  const response = await API.put('/api/option/batch', { options });
  return response.data;
}
