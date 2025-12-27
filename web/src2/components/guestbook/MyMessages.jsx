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

import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Spin, Empty, Typography, Tag, Popconfirm, Button } from '@douyinfe/semi-ui';
import { IconRefresh, IconDelete } from '@douyinfe/semi-icons';
import { Clock, CheckCircle, XCircle, Star } from 'lucide-react';
import {
  IllustrationNoContent,
  IllustrationNoContentDark,
} from '@douyinfe/semi-illustrations';
import { getRelativeTime } from '../../helpers';

const { Text } = Typography;

// 状态配置
const STATUS_CONFIG = {
  pending: {
    color: 'amber',
    icon: Clock,
    label: '待审核',
  },
  approved: {
    color: 'green',
    icon: CheckCircle,
    label: '已通过',
  },
  rejected: {
    color: 'red',
    icon: XCircle,
    label: '已拒绝',
  },
};

// 我的留言卡片
const MyMessageCard = ({ message, onDelete, deleting }) => {
  const { t } = useTranslation();
  const statusConfig = STATUS_CONFIG[message.status] || STATUS_CONFIG.pending;
  const StatusIcon = statusConfig.icon;

  return (
    <div className='p-4 rounded-xl bg-semi-color-bg-1 border border-semi-color-border'>
      {/* 头部：状态和时间 */}
      <div className='flex items-center justify-between mb-3'>
        <div className='flex items-center gap-2'>
          <Tag
            color={statusConfig.color}
            size='small'
            className='flex items-center gap-1'
          >
            <StatusIcon size={12} />
            {t(statusConfig.label)}
          </Tag>
          
          {message.is_featured && (
            <Tag color='orange' size='small' className='flex items-center gap-1'>
              <Star size={12} fill='currentColor' />
              {t('精选')}
            </Tag>
          )}
        </div>
        
        <Text type='tertiary' size='small'>
          {getRelativeTime(new Date(message.created_at * 1000))}
        </Text>
      </div>

      {/* 留言内容 */}
      <Text className='whitespace-pre-wrap break-words block mb-3'>
        {message.content}
      </Text>

      {/* 操作按钮 */}
      <div className='flex justify-end'>
        <Popconfirm
          title={t('确认删除')}
          content={t('删除后无法恢复，确定要删除这条留言吗？')}
          onConfirm={() => onDelete(message.id)}
          okText={t('删除')}
          cancelText={t('取消')}
          okType='danger'
        >
          <Button
            type='danger'
            theme='borderless'
            size='small'
            icon={<IconDelete />}
            loading={deleting}
          >
            {t('删除')}
          </Button>
        </Popconfirm>
      </div>
    </div>
  );
};

const MyMessages = ({
  messages = [],
  loading = false,
  onDelete,
  onRefresh,
}) => {
  const { t } = useTranslation();
  const [deletingId, setDeletingId] = useState(null);

  const handleDelete = async (messageId) => {
    setDeletingId(messageId);
    await onDelete(messageId);
    setDeletingId(null);
  };

  // 加载状态
  if (loading && messages.length === 0) {
    return (
      <div className='flex justify-center items-center py-16'>
        <Spin size='large' />
      </div>
    );
  }

  // 空状态
  if (!loading && messages.length === 0) {
    return (
      <div className='py-16'>
        <Empty
          image={<IllustrationNoContent style={{ width: 150, height: 150 }} />}
          darkModeImage={<IllustrationNoContentDark style={{ width: 150, height: 150 }} />}
          title={t('暂无留言')}
          description={t('你还没有发表过留言')}
        />
      </div>
    );
  }

  return (
    <div className='space-y-4'>
      {/* 刷新按钮 */}
      <div className='flex justify-end'>
        <button
          onClick={onRefresh}
          disabled={loading}
          className='flex items-center gap-1 px-3 py-1.5 text-sm text-semi-color-text-2 hover:text-semi-color-primary rounded-lg hover:bg-semi-color-fill-0 transition-colors disabled:opacity-50'
        >
          <IconRefresh spin={loading} />
          {t('刷新')}
        </button>
      </div>

      {/* 留言列表 */}
      <div className='grid gap-4'>
        {messages.map((message) => (
          <MyMessageCard
            key={message.id}
            message={message}
            onDelete={handleDelete}
            deleting={deletingId === message.id}
          />
        ))}
      </div>
    </div>
  );
};

export default MyMessages;
