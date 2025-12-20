import React, { useEffect, useState } from 'react';
import { Modal, Form, Select, Input, Toast, Spin } from '@douyinfe/semi-ui';
import {
  getAllLevels,
  setUserLevel,
  syncUserRecharge,
} from '../../../../services/levelService';

const UserLevelModal = ({ visible, onCancel, user, onSuccess, t }) => {
  const [loading, setLoading] = useState(false);
  const [levels, setLevels] = useState([]);
  const [formApi, setFormApi] = useState(null);

  useEffect(() => {
    if (visible) {
      fetchLevels();
    }
  }, [visible]);

  const fetchLevels = async () => {
    try {
      const res = await getAllLevels();
      if (res.success) {
        setLevels(res.data || []);
      }
    } catch (error) {
      console.error('Failed to fetch levels:', error);
    }
  };

  const handleSubmit = async () => {
    if (!formApi) return;
    const values = formApi.getValues();
    
    setLoading(true);
    try {
      const res = await setUserLevel(user.id, values.level_id, values.reason);
      if (res.success) {
        Toast.success('等级调整成功');
        onCancel();
        onSuccess && onSuccess();
      } else {
        Toast.error(res.message || '操作失败');
      }
    } catch (error) {
      Toast.error('操作失败');
    } finally {
      setLoading(false);
    }
  };

  const handleSyncRecharge = async () => {
    setLoading(true);
    try {
      const res = await syncUserRecharge(user.id);
      if (res.success) {
        Toast.success(
          `同步成功，累计充值: $${res.data.cumulative_recharge?.toFixed(2) || 0}` +
            (res.data.level_changed ? `，等级已升级到 ${res.data.new_level}` : '')
        );
        onSuccess && onSuccess();
      } else {
        Toast.error(res.message || '同步失败');
      }
    } catch (error) {
      Toast.error('同步失败');
    } finally {
      setLoading(false);
    }
  };

  const levelOptions = levels.map((level) => ({
    value: level.id,
    label: `${level.name} (${level.id})`,
  }));

  return (
    <Modal
      title={`调整用户等级 - ${user?.username || ''}`}
      visible={visible}
      onOk={handleSubmit}
      onCancel={onCancel}
      okText="确认调整"
      cancelText="取消"
      confirmLoading={loading}
      footer={
        <div style={{ display: 'flex', justifyContent: 'space-between' }}>
          <button
            className="semi-button semi-button-tertiary"
            onClick={handleSyncRecharge}
            disabled={loading}
          >
            同步累计充值
          </button>
          <div>
            <button
              className="semi-button semi-button-tertiary"
              onClick={onCancel}
            >
              取消
            </button>
            <button
              className="semi-button semi-button-primary"
              onClick={handleSubmit}
              disabled={loading}
              style={{ marginLeft: 8 }}
            >
              确认调整
            </button>
          </div>
        </div>
      }
    >
      <Spin spinning={loading}>
        <div style={{ marginBottom: 16 }}>
          <p>当前等级: <strong>{user?.level || 'tier_1'}</strong></p>
          <p>累计充值: <strong>${user?.cumulative_recharge?.toFixed(2) || 0}</strong></p>
        </div>
        <Form
          getFormApi={(api) => setFormApi(api)}
          initValues={{ level_id: user?.level || 'tier_1' }}
          labelPosition="left"
          labelWidth={80}
        >
          <Form.Select
            field="level_id"
            label="目标等级"
            optionList={levelOptions}
            style={{ width: '100%' }}
          />
          <Form.TextArea
            field="reason"
            label="调整原因"
            placeholder="请输入调整原因（可选）"
          />
        </Form>
      </Spin>
    </Modal>
  );
};

export default UserLevelModal;
