/*
Copyright (C) 2025 QuantumNous
*/

import React, { useContext } from 'react';
import { useTranslation } from 'react-i18next';
import { Typography, Tabs, TabPane } from '@douyinfe/semi-ui';
import { IconComment, IconUser } from '@douyinfe/semi-icons';
import { MessageSquare, Sparkles } from 'lucide-react';
import { UserContext } from '../../context/User';
import MessageWall from '../../components/guestbook/MessageWall';
import SubmitForm from '../../components/guestbook/SubmitForm';
import MyMessages from '../../components/guestbook/MyMessages';
import useGuestbook from '../../hooks/guestbook/useGuestbook';

const { Title, Text } = Typography;

const GuestbookPage = () => {
  const { t } = useTranslation();
  const [userState] = useContext(UserContext);
  const isLoggedIn = !!userState?.user;

  const {
    messages,
    myMessages,
    loading,
    myMessagesLoading,
    total,
    page,
    pageSize,
    submitting,
    fetchMessages,
    fetchMyMessages,
    submitMessage,
    deleteMyMessage,
    setPage,
  } = useGuestbook();

  return (
    <div className='mt-[60px] px-4 py-8 max-w-6xl mx-auto'>
      {/* 页面头部 - 动画渐变背景 */}
      <div className='relative mb-10 text-center overflow-hidden rounded-3xl py-12 px-6'>
        {/* 动画渐变背景 */}
        <div className='absolute inset-0 bg-gradient-to-br from-purple-500/20 via-pink-500/20 to-orange-500/20 dark:from-purple-600/30 dark:via-pink-600/30 dark:to-orange-600/30' />
        <div className='absolute inset-0 bg-gradient-to-tr from-blue-500/10 via-transparent to-amber-500/10 animate-gradient-shift' />
        
        {/* 装饰元素 */}
        <div className='absolute top-4 left-8 opacity-20'>
          <Sparkles size={24} className='text-purple-500 animate-pulse' />
        </div>
        <div className='absolute bottom-6 right-12 opacity-20'>
          <Sparkles size={20} className='text-pink-500 animate-pulse delay-300' />
        </div>
        <div className='absolute top-1/2 left-4 opacity-10'>
          <div className='w-32 h-32 rounded-full bg-gradient-to-br from-purple-400 to-pink-400 blur-3xl' />
        </div>
        <div className='absolute top-1/3 right-8 opacity-10'>
          <div className='w-24 h-24 rounded-full bg-gradient-to-br from-orange-400 to-amber-400 blur-3xl' />
        </div>

        {/* 内容 */}
        <div className='relative z-10'>
          {/* 图标 */}
          <div className='inline-flex items-center justify-center w-20 h-20 rounded-2xl bg-gradient-to-br from-purple-500 to-pink-500 mb-5 shadow-lg shadow-purple-500/30'>
            <MessageSquare size={40} className='text-white' />
          </div>
          
          {/* 标题 */}
          <Title 
            heading={1} 
            className='!mb-3 !text-3xl md:!text-4xl bg-gradient-to-r from-purple-600 via-pink-600 to-orange-500 bg-clip-text text-transparent'
          >
            {t('留言板')}
          </Title>
          
          {/* 描述 */}
          <Text type='tertiary' className='text-base md:text-lg'>
            {t('分享你的想法和反馈，与社区互动')}
          </Text>
        </div>
      </div>

      {/* 留言提交表单 */}
      <div className='mb-8'>
        <SubmitForm
          isLoggedIn={isLoggedIn}
          submitting={submitting}
          onSubmit={submitMessage}
        />
      </div>

      {/* 标签页：留言墙 / 我的留言 */}
      <Tabs type='line' defaultActiveKey='wall'>
        <TabPane
          tab={
            <span className='flex items-center gap-2'>
              <IconComment />
              {t('留言墙')}
            </span>
          }
          itemKey='wall'
        >
          <MessageWall
            messages={messages}
            loading={loading}
            total={total}
            page={page}
            pageSize={pageSize}
            onPageChange={setPage}
            onRefresh={fetchMessages}
          />
        </TabPane>

        {isLoggedIn && (
          <TabPane
            tab={
              <span className='flex items-center gap-2'>
                <IconUser />
                {t('我的留言')}
              </span>
            }
            itemKey='my'
          >
            <MyMessages
              messages={myMessages}
              loading={myMessagesLoading}
              onDelete={deleteMyMessage}
              onRefresh={fetchMyMessages}
            />
          </TabPane>
        )}
      </Tabs>
    </div>
  );
};

export default GuestbookPage;
