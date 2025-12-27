/*
Copyright (C) 2025 QuantumNous
*/

import React from 'react';
import { useTranslation } from 'react-i18next';
import { Typography } from '@douyinfe/semi-ui';
import { Shield, MessageCircle } from 'lucide-react';
import { getRelativeTime } from '../../helpers';

const { Text } = Typography;

/**
 * AdminReply - 管理员回复展示组件
 * 渐变背景、官方徽章、回复内容和时间
 */
const AdminReply = ({ reply, replyAt }) => {
  const { t } = useTranslation();

  if (!reply) return null;

  return (
    <div className='mt-4 relative'>
      {/* 连接线 */}
      <div className='absolute -top-2 left-6 w-0.5 h-2 bg-gradient-to-b from-transparent to-purple-400/50' />
      
      {/* 回复卡片 */}
      <div className='relative rounded-xl overflow-hidden'>
        {/* 渐变背景 */}
        <div className='absolute inset-0 bg-gradient-to-br from-purple-500/10 via-indigo-500/10 to-purple-600/10 dark:from-purple-500/20 dark:via-indigo-500/20 dark:to-purple-600/20' />
        
        {/* 左侧装饰条 */}
        <div className='absolute left-0 top-0 bottom-0 w-1 bg-gradient-to-b from-purple-500 to-indigo-500' />
        
        {/* 内容区域 */}
        <div className='relative p-4 pl-5'>
          {/* 头部：官方徽章 + 时间 */}
          <div className='flex items-center justify-between mb-2'>
            <div className='flex items-center gap-2'>
              {/* 官方徽章 */}
              <div className='flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-gradient-to-r from-purple-500 to-indigo-500 text-white text-xs font-medium shadow-sm'>
                <Shield size={12} fill='currentColor' />
                <span>{t('官方回复')}</span>
              </div>
              
              {/* 回复图标 */}
              <MessageCircle size={14} className='text-purple-400' />
            </div>
            
            {/* 回复时间 */}
            {replyAt && (
              <Text type='tertiary' size='small'>
                {getRelativeTime(new Date(replyAt * 1000))}
              </Text>
            )}
          </div>
          
          {/* 回复内容 */}
          <Text className='whitespace-pre-wrap break-words text-gray-700 dark:text-gray-200'>
            {reply}
          </Text>
        </div>
      </div>
    </div>
  );
};

export default AdminReply;
