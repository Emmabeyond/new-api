/*
Copyright (C) 2025 QuantumNous
*/

import React from 'react';
import { useTranslation } from 'react-i18next';
import { Spin, Empty, Pagination, Typography, Avatar } from '@douyinfe/semi-ui';
import { IconRefresh } from '@douyinfe/semi-icons';
import {
  IllustrationNoContent,
  IllustrationNoContentDark,
} from '@douyinfe/semi-illustrations';
import { getRelativeTime } from '../../helpers';
import Card3D from './Card3D';
import MasonryLayout from './MasonryLayout';
import AdminReply from './AdminReply';

const { Text } = Typography;

// 头像颜色列表
const AVATAR_COLORS = [
  'amber', 'blue', 'cyan', 'green', 'grey', 
  'indigo', 'light-blue', 'light-green', 'lime', 
  'orange', 'pink', 'purple', 'red', 'teal', 'violet', 'yellow'
];

// 生成头像颜色
const getAvatarColor = (username) => {
  const index = username?.charCodeAt(0) % AVATAR_COLORS.length || 0;
  return AVATAR_COLORS[index];
};

// 生成头像光环颜色
const getGlowColor = (username) => {
  const colors = [
    'rgba(251, 191, 36, 0.4)', // amber
    'rgba(59, 130, 246, 0.4)', // blue
    'rgba(6, 182, 212, 0.4)',  // cyan
    'rgba(34, 197, 94, 0.4)',  // green
    'rgba(99, 102, 241, 0.4)', // indigo
    'rgba(236, 72, 153, 0.4)', // pink
    'rgba(139, 92, 246, 0.4)', // purple
    'rgba(239, 68, 68, 0.4)',  // red
    'rgba(20, 184, 166, 0.4)', // teal
  ];
  const index = username?.charCodeAt(0) % colors.length || 0;
  return colors[index];
};

// 留言卡片内容组件
const MessageCardContent = ({ message, t }) => {
  const avatarColor = getAvatarColor(message.username);
  const glowColor = getGlowColor(message.username);

  return (
    <>
      {/* 用户信息 */}
      <div className='flex items-start gap-3'>
        {/* 头像带光环效果 */}
        <div className='relative flex-shrink-0'>
          {/* 光环 */}
          <div 
            className='absolute -inset-1 rounded-full opacity-60 blur-sm'
            style={{ background: glowColor }}
          />
          <Avatar 
            size='small' 
            color={avatarColor}
            className='relative'
          >
            {message.username?.charAt(0)?.toUpperCase()}
          </Avatar>
        </div>
        
        <div className='flex-1 min-w-0'>
          <div className='flex items-center gap-2 mb-1.5'>
            <Text strong className='truncate text-gray-800 dark:text-gray-100'>
              {message.username}
            </Text>
            <Text type='tertiary' size='small'>
              {getRelativeTime(new Date(message.created_at * 1000))}
            </Text>
          </div>
          
          {/* 留言内容 */}
          <Text className='whitespace-pre-wrap break-words text-gray-600 dark:text-gray-300'>
            {message.content}
          </Text>
        </div>
      </div>

      {/* 管理员回复 */}
      {message.admin_reply && (
        <AdminReply 
          reply={message.admin_reply} 
          replyAt={message.admin_reply_at}
        />
      )}
    </>
  );
};

const MessageWall = ({
  messages = [],
  loading = false,
  total = 0,
  page = 1,
  pageSize = 20,
  onPageChange,
  onRefresh,
}) => {
  const { t } = useTranslation();

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
          description={t('成为第一个留言的人吧！')}
        />
      </div>
    );
  }

  // 渲染单个卡片
  const renderMessageCard = (message, index) => (
    <Card3D 
      isFeatured={message.is_featured}
      glowColor={getGlowColor(message.username)}
    >
      <MessageCardContent message={message} t={t} />
    </Card3D>
  );

  return (
    <div className='space-y-6'>
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

      {/* 瀑布流留言列表 */}
      <MasonryLayout
        items={messages}
        renderItem={renderMessageCard}
        gap={20}
      />

      {/* 分页 */}
      {total > pageSize && (
        <div className='flex justify-center pt-6'>
          <Pagination
            currentPage={page}
            pageSize={pageSize}
            total={total}
            onPageChange={onPageChange}
            showSizeChanger={false}
          />
        </div>
      )}
    </div>
  );
};

export default MessageWall;
