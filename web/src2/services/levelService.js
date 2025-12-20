import { API } from '../helpers';

/**
 * 获取所有等级配置（管理员）
 */
export async function getAllLevels() {
  const response = await API.get('/api/level/admin/');
  return response.data;
}

/**
 * 创建等级配置（管理员）
 */
export async function createLevel(level) {
  const response = await API.post('/api/level/admin/', level);
  return response.data;
}

/**
 * 更新等级配置（管理员）
 */
export async function updateLevel(id, level) {
  const response = await API.put(`/api/level/admin/${id}`, level);
  return response.data;
}

/**
 * 删除等级配置（管理员）
 */
export async function deleteLevel(id) {
  const response = await API.delete(`/api/level/admin/${id}`);
  return response.data;
}

/**
 * 获取等级用户统计（管理员）
 */
export async function getLevelStats() {
  const response = await API.get('/api/level/admin/stats');
  return response.data;
}

/**
 * 获取当前用户等级信息
 */
export async function getUserLevel() {
  const response = await API.get('/api/level/self');
  return response.data;
}

/**
 * 获取用户升级进度
 */
export async function getUserLevelProgress() {
  const response = await API.get('/api/level/self/progress');
  return response.data;
}

/**
 * 获取用户等级权益详情
 */
export async function getUserLevelBenefits() {
  const response = await API.get('/api/level/self/benefits');
  return response.data;
}

/**
 * 获取等级对比信息
 */
export async function getLevelComparison() {
  const response = await API.get('/api/level/compare');
  return response.data;
}

/**
 * 调整用户等级（管理员）
 */
export async function setUserLevel(userId, levelId, reason) {
  const response = await API.put(`/api/level/admin/user/${userId}`, {
    level_id: levelId,
    reason: reason,
  });
  return response.data;
}

/**
 * 获取用户等级变更历史（管理员）
 */
export async function getUserLevelHistory(userId, page = 1, pageSize = 10) {
  const response = await API.get(`/api/level/admin/user/${userId}/history`, {
    params: { page, page_size: pageSize },
  });
  return response.data;
}

/**
 * 同步用户累计充值（管理员）
 */
export async function syncUserRecharge(userId) {
  const response = await API.post(`/api/level/admin/user/${userId}/sync-recharge`);
  return response.data;
}
