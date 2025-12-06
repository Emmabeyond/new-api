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

import React, { useEffect, useState, useCallback } from 'react';
import { Modal, Table, Tag, Spin } from '@douyinfe/semi-ui';
import { API, showError } from '../../helpers';
import dayjs from 'dayjs';

const CheckinHistory = ({ visible, onClose, t, renderQuota }) => {
  const [loading, setLoading] = useState(false);
  const [data, setData] = useState([]);
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0,
  });

  const fetchHistory = useCallback(async (page = 1, pageSize = 10) => {
    setLoading(true);
    try {
      const res = await API.get('/api/user/checkin/history', {
        params: { page, page_size: pageSize }
      });
      const { success, data: resData, message } = res.data;
      if (success) {
        setData(resData.items || []);
        setPagination({
          current: resData.page || page,
          pageSize: resData.page_size || pageSize,
          total: resData.total || 0,
        });
      } else {
        showError(message);
      }
    } catch (err) {
      showError(t('获取签到历史失败'));
    } finally {
      setLoading(false);
    }
  }, [t]);

  useEffect(() => {
    if (visible) {
      fetchHistory(1, pagination.pageSize);
    }
  }, [visible, fetchHistory]);

  const handlePageChange = (page) => {
    fetchHistory(page, pagination.pageSize);
  };

  const columns = [
    {
      title: t('日期'),
      dataIndex: 'checkin_date',
      key: 'checkin_date',
      render: (text) => dayjs(text).format('YYYY-MM-DD'),
    },
    {
      title: t('类型'),
      dataIndex: 'checkin_type',
      key: 'checkin_type',
      render: (type) => (
        <Tag color={type === 1 ? 'green' : 'orange'}>
          {type === 1 ? t('正常签到') : t('补签')}
        </Tag>
      ),
    },
    {
      title: t('基础奖励'),
      dataIndex: 'base_reward',
      key: 'base_reward',
      render: (value) => renderQuota(value),
    },
    {
      title: t('惊喜奖励'),
      dataIndex: 'bonus_reward',
      key: 'bonus_reward',
      render: (value) => value > 0 ? (
        <Tag color='orange'>+{renderQuota(value)}</Tag>
      ) : '-',
    },
    {
      title: t('总计'),
      key: 'total',
      render: (_, record) => renderQuota(record.base_reward + record.bonus_reward),
    },
  ];

  return (
    <Modal
      title={t('签到记录')}
      visible={visible}
      onCancel={onClose}
      footer={null}
      width={700}
      centered
    >
      <Spin spinning={loading}>
        <Table
          columns={columns}
          dataSource={data}
          rowKey='id'
          pagination={{
            current: pagination.current,
            pageSize: pagination.pageSize,
            total: pagination.total,
            onChange: handlePageChange,
            showSizeChanger: false,
          }}
        />
      </Spin>
    </Modal>
  );
};

export default CheckinHistory;
