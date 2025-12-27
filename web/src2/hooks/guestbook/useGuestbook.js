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

import { useState, useEffect, useCallback, useContext } from 'react';
import { API, showError, showSuccess } from '../../helpers';
import { UserContext } from '../../context/User';
import { useTranslation } from 'react-i18next';

const useGuestbook = () => {
  const { t } = useTranslation();
  const [userState] = useContext(UserContext);
  const isLoggedIn = !!userState?.user;

  // 留言墙状态
  const [messages, setMessages] = useState([]);
  const [loading, setLoading] = useState(true);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize] = useState(20);

  // 我的留言状态
  const [myMessages, setMyMessages] = useState([]);
  const [myMessagesLoading, setMyMessagesLoading] = useState(false);

  // 提交状态
  const [submitting, setSubmitting] = useState(false);

  // 获取已审核留言列表
  const fetchMessages = useCallback(async (pageNum = page) => {
    setLoading(true);
    try {
      const res = await API.get('/api/guestbook/messages', {
        params: { page: pageNum, page_size: pageSize },
      });
      const { success, data, message } = res.data;
      if (success) {
        setMessages(data.messages || []);
        setTotal(data.total || 0);
      } else {
        showError(message || t('获取留言失败'));
      }
    } catch (err) {
      showError(t('获取留言失败'));
    } finally {
      setLoading(false);
    }
  }, [page, pageSize, t]);

  // 获取我的留言
  const fetchMyMessages = useCallback(async () => {
    if (!isLoggedIn) return;
    
    setMyMessagesLoading(true);
    try {
      const res = await API.get('/api/guestbook/my-messages');
      const { success, data, message } = res.data;
      if (success) {
        setMyMessages(data || []);
      } else {
        showError(message || t('获取我的留言失败'));
      }
    } catch (err) {
      showError(t('获取我的留言失败'));
    } finally {
      setMyMessagesLoading(false);
    }
  }, [isLoggedIn, t]);

  // 提交留言
  const submitMessage = useCallback(async (content) => {
    if (!isLoggedIn) {
      showError(t('请先登录'));
      return false;
    }

    setSubmitting(true);
    try {
      const res = await API.post('/api/guestbook/messages', { content });
      const { success, message } = res.data;
      if (success) {
        showSuccess(message || t('留言提交成功，等待审核'));
        // 刷新我的留言列表
        await fetchMyMessages();
        return true;
      } else {
        showError(message || t('留言提交失败'));
        return false;
      }
    } catch (err) {
      showError(t('留言提交失败'));
      return false;
    } finally {
      setSubmitting(false);
    }
  }, [isLoggedIn, fetchMyMessages, t]);

  // 删除我的留言
  const deleteMyMessage = useCallback(async (messageId) => {
    if (!isLoggedIn) {
      showError(t('请先登录'));
      return false;
    }

    try {
      const res = await API.delete(`/api/guestbook/messages/${messageId}`);
      const { success, message } = res.data;
      if (success) {
        showSuccess(message || t('留言删除成功'));
        // 刷新我的留言列表
        await fetchMyMessages();
        // 刷新留言墙
        await fetchMessages();
        return true;
      } else {
        showError(message || t('删除失败'));
        return false;
      }
    } catch (err) {
      showError(t('删除失败'));
      return false;
    }
  }, [isLoggedIn, fetchMyMessages, fetchMessages, t]);

  // 页码变化
  const handlePageChange = useCallback((newPage) => {
    setPage(newPage);
    fetchMessages(newPage);
  }, [fetchMessages]);

  // 初始化加载
  useEffect(() => {
    fetchMessages();
  }, []);

  // 登录状态变化时加载我的留言
  useEffect(() => {
    if (isLoggedIn) {
      fetchMyMessages();
    } else {
      setMyMessages([]);
    }
  }, [isLoggedIn]);

  return {
    // 留言墙
    messages,
    loading,
    total,
    page,
    pageSize,
    fetchMessages,
    setPage: handlePageChange,
    
    // 我的留言
    myMessages,
    myMessagesLoading,
    fetchMyMessages,
    
    // 操作
    submitting,
    submitMessage,
    deleteMyMessage,
  };
};

export default useGuestbook;
