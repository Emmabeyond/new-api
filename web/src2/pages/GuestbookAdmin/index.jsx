/*
Copyright (C) 2025 QuantumNous
*/

import React, { useEffect, useState, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import {
  Card,
  Table,
  Button,
  Input,
  Select,
  Tag,
  Typography,
  Space,
  Popconfirm,
  Modal,
  TextArea,
} from '@douyinfe/semi-ui';
import {
  IconSearch,
  IconRefresh,
  IconDelete,
  IconTick,
  IconClose,
  IconStar,
  IconStarStroked,
  IconEdit,
  IconComment,
} from '@douyinfe/semi-icons';
import { MessageCircle, Shield } from 'lucide-react';
import { API, showError, showSuccess } from '../../helpers';
import dayjs from 'dayjs';

const { Text } = Typography;

const MAX_REPLY_LENGTH = 300;

const GuestbookAdmin = () => {
  const { t } = useTranslation();

  // 留言列表
  const [messages, setMessages] = useState([]);
  const [loading, setLoading] = useState(false);
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 20,
    total: 0,
  });

  // 筛选条件
  const [filters, setFilters] = useState({
    status: '',
    keyword: '',
  });

  // 精选数量
  const [featuredCount, setFeaturedCount] = useState(0);
  const MAX_FEATURED = 5;

  // 回复相关状态
  const [replyModalVisible, setReplyModalVisible] = useState(false);
  const [currentMessage, setCurrentMessage] = useState(null);
  const [replyContent, setReplyContent] = useState('');
  const [replySubmitting, setReplySubmitting] = useState(false);

  // 获取留言列表
  const fetchMessages = useCallback(
    async (page = 1, pageSize = 20) => {
      setLoading(true);
      try {
        const params = { page, page_size: pageSize };
        if (filters.status) params.status = filters.status;
        if (filters.keyword) params.keyword = filters.keyword;

        const res = await API.get('/api/guestbook/admin/messages', { params });
        const { success, data, message } = res.data;
        if (success) {
          const messageList = data.messages || [];
          setMessages(messageList);
          setPagination({
            current: data.page || page,
            pageSize: pageSize,
            total: data.total || 0,
          });
          const featured = messageList.filter((m) => m.is_featured).length;
          setFeaturedCount(featured);
        } else {
          showError(message);
        }
      } catch (err) {
        showError(t('获取留言列表失败'));
      } finally {
        setLoading(false);
      }
    },
    [filters, t]
  );

  useEffect(() => {
    fetchMessages();
  }, []);

  const handleSearch = () => fetchMessages(1, pagination.pageSize);
  const handleReset = () => {
    setFilters({ status: '', keyword: '' });
    fetchMessages(1, pagination.pageSize);
  };
  const handlePageChange = (page) => fetchMessages(page, pagination.pageSize);

  // 审核留言
  const handleReview = async (id, status) => {
    try {
      const res = await API.put(`/api/guestbook/admin/messages/${id}/review`, { status });
      const { success, message } = res.data;
      if (success) {
        showSuccess(message || t('审核成功'));
        fetchMessages(pagination.current, pagination.pageSize);
      } else {
        showError(message);
      }
    } catch (err) {
      showError(t('审核失败'));
    }
  };

  // 设置/取消精选
  const handleFeature = async (id, isFeatured) => {
    if (isFeatured && featuredCount >= MAX_FEATURED) {
      showError(t('精选留言数量已达上限') + ` (${MAX_FEATURED}${t('条')})`);
      return;
    }
    try {
      const res = await API.put(`/api/guestbook/admin/messages/${id}/feature`, { is_featured: isFeatured });
      const { success, message } = res.data;
      if (success) {
        showSuccess(message || t('操作成功'));
        fetchMessages(pagination.current, pagination.pageSize);
      } else {
        showError(message);
      }
    } catch (err) {
      showError(t('操作失败'));
    }
  };

  // 删除留言
  const handleDelete = async (id) => {
    try {
      const res = await API.delete(`/api/guestbook/admin/messages/${id}`);
      const { success, message } = res.data;
      if (success) {
        showSuccess(message || t('删除成功'));
        fetchMessages(pagination.current, pagination.pageSize);
      } else {
        showError(message);
      }
    } catch (err) {
      showError(t('删除失败'));
    }
  };

  // 打开回复弹窗
  const openReplyModal = (record) => {
    setCurrentMessage(record);
    setReplyContent(record.admin_reply || '');
    setReplyModalVisible(true);
  };

  // 提交回复
  const handleSubmitReply = async () => {
    if (!replyContent.trim()) {
      showError(t('回复内容不能为空'));
      return;
    }
    if (replyContent.length > MAX_REPLY_LENGTH) {
      showError(t('回复内容不能超过300字'));
      return;
    }

    setReplySubmitting(true);
    try {
      const res = await API.put(`/api/guestbook/admin/messages/${currentMessage.id}/reply`, {
        reply: replyContent.trim(),
      });
      const { success, message } = res.data;
      if (success) {
        showSuccess(message || t('回复成功'));
        setReplyModalVisible(false);
        setReplyContent('');
        setCurrentMessage(null);
        fetchMessages(pagination.current, pagination.pageSize);
      } else {
        showError(message);
      }
    } catch (err) {
      showError(t('回复失败'));
    } finally {
      setReplySubmitting(false);
    }
  };

  // 删除回复
  const handleDeleteReply = async (id) => {
    try {
      const res = await API.delete(`/api/guestbook/admin/messages/${id}/reply`);
      const { success, message } = res.data;
      if (success) {
        showSuccess(message || t('回复已删除'));
        fetchMessages(pagination.current, pagination.pageSize);
      } else {
        showError(message);
      }
    } catch (err) {
      showError(t('删除回复失败'));
    }
  };

  const renderStatus = (status) => {
    const statusMap = {
      pending: { color: 'orange', text: t('待审核') },
      approved: { color: 'green', text: t('已通过') },
      rejected: { color: 'red', text: t('已拒绝') },
    };
    const config = statusMap[status] || { color: 'grey', text: status };
    return <Tag color={config.color}>{config.text}</Tag>;
  };

  const columns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
    {
      title: t('用户'),
      dataIndex: 'username',
      key: 'username',
      width: 100,
      render: (text, record) => (
        <div>
          <Text strong>{text}</Text>
          <br />
          <Text type="tertiary" size="small">ID: {record.user_id}</Text>
        </div>
      ),
    },
    {
      title: t('留言内容'),
      dataIndex: 'content',
      key: 'content',
      render: (text, record) => (
        <div style={{ maxWidth: 350 }}>
          <div style={{ wordBreak: 'break-word' }}>{text}</div>
          {record.admin_reply && (
            <div className="mt-2 p-2 rounded-lg bg-purple-50 dark:bg-purple-900/20 border-l-2 border-purple-500">
              <div className="flex items-center gap-1 mb-1">
                <Shield size={12} className="text-purple-500" />
                <Text size="small" type="tertiary">{t('官方回复')}</Text>
              </div>
              <Text size="small" className="text-purple-700 dark:text-purple-300">
                {record.admin_reply}
              </Text>
            </div>
          )}
        </div>
      ),
    },
    { title: t('状态'), dataIndex: 'status', key: 'status', width: 80, render: renderStatus },
    {
      title: t('回复'),
      dataIndex: 'admin_reply',
      key: 'admin_reply',
      width: 80,
      render: (reply) => reply ? (
        <Tag color="purple" prefixIcon={<MessageCircle size={12} />}>{t('已回复')}</Tag>
      ) : (
        <Text type="tertiary">-</Text>
      ),
    },
    {
      title: t('精选'),
      dataIndex: 'is_featured',
      key: 'is_featured',
      width: 70,
      render: (isFeatured) => isFeatured ? (
        <Tag color="yellow" prefixIcon={<IconStar />}>{t('精选')}</Tag>
      ) : '-',
    },
    {
      title: t('时间'),
      dataIndex: 'created_at',
      key: 'created_at',
      width: 140,
      render: (timestamp) => timestamp ? dayjs.unix(timestamp).format('YYYY-MM-DD HH:mm') : '-',
    },
    {
      title: t('操作'),
      key: 'action',
      width: 320,
      render: (_, record) => (
        <Space wrap>
          {record.status === 'pending' && (
            <>
              <Button theme="solid" type="primary" size="small" icon={<IconTick />}
                onClick={() => handleReview(record.id, 'approved')}>{t('通过')}</Button>
              <Button theme="solid" type="warning" size="small" icon={<IconClose />}
                onClick={() => handleReview(record.id, 'rejected')}>{t('拒绝')}</Button>
            </>
          )}
          {record.status === 'approved' && (
            <>
              <Button theme={record.is_featured ? 'solid' : 'light'}
                type={record.is_featured ? 'warning' : 'tertiary'} size="small"
                icon={record.is_featured ? <IconStar /> : <IconStarStroked />}
                onClick={() => handleFeature(record.id, !record.is_featured)}>
                {record.is_featured ? t('取消精选') : t('精选')}
              </Button>
              <Button theme="light" type="primary" size="small"
                icon={record.admin_reply ? <IconEdit /> : <IconComment />}
                onClick={() => openReplyModal(record)}>
                {record.admin_reply ? t('编辑回复') : t('回复')}
              </Button>
              {record.admin_reply && (
                <Popconfirm title={t('确认删除回复')} content={t('确定要删除这条回复吗？')}
                  onConfirm={() => handleDeleteReply(record.id)}>
                  <Button theme="borderless" type="tertiary" size="small">{t('删除回复')}</Button>
                </Popconfirm>
              )}
            </>
          )}
          <Popconfirm title={t('确认删除')} content={t('确定要删除这条留言吗？此操作不可恢复。')}
            onConfirm={() => handleDelete(record.id)}>
            <Button theme="borderless" type="danger" size="small" icon={<IconDelete />}>{t('删除')}</Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div className="mt-[60px] px-2">
      <div className="w-full max-w-7xl mx-auto space-y-6">
        <Card title={t('留言管理')} headerExtraContent={
          <Space>
            <Text type="tertiary">{t('精选')}: {featuredCount}/{MAX_FEATURED}</Text>
            <Button theme="borderless" icon={<IconRefresh />}
              onClick={() => fetchMessages(pagination.current, pagination.pageSize)}>{t('刷新')}</Button>
          </Space>
        }>
          <div className="flex flex-wrap gap-4 mb-4">
            <Input placeholder={t('搜索内容或用户名')} value={filters.keyword}
              onChange={(value) => setFilters({ ...filters, keyword: value })}
              style={{ width: 200 }} prefix={<IconSearch />} />
            <Select placeholder={t('审核状态')} value={filters.status}
              onChange={(value) => setFilters({ ...filters, status: value })} style={{ width: 140 }}>
              <Select.Option value="">{t('全部状态')}</Select.Option>
              <Select.Option value="pending">{t('待审核')}</Select.Option>
              <Select.Option value="approved">{t('已通过')}</Select.Option>
              <Select.Option value="rejected">{t('已拒绝')}</Select.Option>
            </Select>
            <Button theme="solid" icon={<IconSearch />} onClick={handleSearch}>{t('搜索')}</Button>
            <Button theme="light" onClick={handleReset}>{t('重置')}</Button>
          </div>
          <Table columns={columns} dataSource={messages} rowKey="id" loading={loading}
            pagination={{ current: pagination.current, pageSize: pagination.pageSize,
              total: pagination.total, onChange: handlePageChange, showSizeChanger: false }} />
        </Card>
      </div>

      {/* 回复弹窗 */}
      <Modal title={currentMessage?.admin_reply ? t('编辑回复') : t('回复留言')}
        visible={replyModalVisible} onCancel={() => { setReplyModalVisible(false); setReplyContent(''); }}
        onOk={handleSubmitReply} okText={t('提交')} cancelText={t('取消')}
        confirmLoading={replySubmitting} width={500}>
        {currentMessage && (
          <div className="space-y-4">
            <div className="p-3 rounded-lg bg-gray-50 dark:bg-gray-800">
              <div className="flex items-center gap-2 mb-2">
                <Text strong>{currentMessage.username}</Text>
                <Text type="tertiary" size="small">
                  {dayjs.unix(currentMessage.created_at).format('YYYY-MM-DD HH:mm')}
                </Text>
              </div>
              <Text>{currentMessage.content}</Text>
            </div>
            <div>
              <div className="flex items-center justify-between mb-2">
                <Text strong>{t('官方回复')}</Text>
                <Text type="tertiary" size="small">{replyContent.length}/{MAX_REPLY_LENGTH}</Text>
              </div>
              <TextArea value={replyContent} onChange={setReplyContent}
                placeholder={t('输入回复内容...')} maxCount={MAX_REPLY_LENGTH}
                autosize={{ minRows: 3, maxRows: 6 }} />
            </div>
          </div>
        )}
      </Modal>
    </div>
  );
};

export default GuestbookAdmin;
