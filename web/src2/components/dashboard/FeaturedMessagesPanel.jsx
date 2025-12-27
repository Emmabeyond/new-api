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

import React, { useState, useEffect } from 'react';
import { Card, Avatar, Spin, Typography } from '@douyinfe/semi-ui';
import { MessageSquareHeart, Star, ArrowRight } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { API, showError } from '../../helpers';
import { getRelativeTime } from '../../helpers';

const { Text, Paragraph } = Typography;

const FeaturedMessagesPanel = ({ CARD_PROPS }) => {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const [messages, setMessages] = useState([]);
  const [loading, setLoading] = useState(true);

  // 获取精选留言
  const fetchFeaturedMessages = async () => {
    setLoading(true);
    try {
      const res = await API.get('/api/guestbook/featured', {
        params: { limit: 3 },
      });
      const { success, data, message } = res.data;
      if (success) {
        setMessages(data || []);
      } else {
        console.error('Failed to fetch featured messages:', message);
      }
    } catch (err) {
      console.error('Error fetching featured messages:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchFeaturedMessages();
  }, []);

  // 如果没有精选留言，不渲染面板
  if (!loading && messages.length === 0) {
    return null;
  }

  // 点击跳转到留言板
  const handleNavigate = () => {
    navigate('/guestbook');
  };

  // 格式化时间
  const formatTime = (timestamp) => {
    if (!timestamp) return '';
    const date = new Date(timestamp * 1000);
    return getRelativeTime(date.toISOString());
  };

  return (
    <Card
      {...CARD_PROPS}
      className='shadow-sm !rounded-2xl'
      title={
        <div className='flex items-center justify-between w-full'>
          <div className='flex items-center gap-2'>
            <MessageSquareHeart size={16} className='text-pink-500' />
            <span>{t('精选留言')}</span>
          </div>
          <div
            className='flex items-center gap-1 text-xs text-gray-500 cursor-pointer hover:text-blue-500 transition-colors'
            onClick={handleNavigate}
          >
            <span>{t('查看全部')}</span>
            <ArrowRight size={14} />
          </div>
        </div>
      }
      bodyStyle={{ padding: '12px 16px' }}
    >
      {loading ? (
        <div className='flex justify-center items-center py-8'>
          <Spin />
        </div>
      ) : (
        <div className='space-y-3'>
          {messages.map((msg) => (
            <div
              key={msg.id}
              className='p-3 bg-gray-50 dark:bg-gray-800 rounded-lg cursor-pointer hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors'
              onClick={handleNavigate}
            >
              <div className='flex items-start gap-3'>
                <Avatar
                  size='small'
                  style={{ backgroundColor: '#3b82f6' }}
                >
                  {msg.username?.charAt(0)?.toUpperCase() || 'U'}
                </Avatar>
                <div className='flex-1 min-w-0'>
                  <div className='flex items-center gap-2 mb-1'>
                    <Text strong className='text-sm truncate'>
                      {msg.username}
                    </Text>
                    <Star size={12} className='text-amber-500 fill-amber-400 flex-shrink-0 drop-shadow-[0_0_3px_rgba(251,191,36,0.6)]' />
                    <Text type='tertiary' className='text-xs flex-shrink-0'>
                      {formatTime(msg.created_at)}
                    </Text>
                  </div>
                  <Paragraph
                    ellipsis={{ rows: 2 }}
                    className='text-sm text-gray-600 dark:text-gray-300 !mb-0'
                  >
                    {msg.content}
                  </Paragraph>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </Card>
  );
};

export default FeaturedMessagesPanel;
