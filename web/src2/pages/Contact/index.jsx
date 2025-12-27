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

import React, { useEffect, useState } from 'react';
import { API, showError } from '../../helpers';
import { Empty, Spin } from '@douyinfe/semi-ui';
import {
  IllustrationConstruction,
  IllustrationConstructionDark,
} from '@douyinfe/semi-illustrations';
import { useTranslation } from 'react-i18next';
import SafeHTMLContent from '../../components/common/SafeHTMLContent';

// 检查是否为 HTML 内容
const isHtmlContent = (content) => {
  if (!content || typeof content !== 'string') return false;
  const htmlTagRegex = /<\/?[a-z][\s\S]*>/i;
  return htmlTagRegex.test(content);
};

const Contact = () => {
  const { t } = useTranslation();
  const [contact, setContact] = useState('');
  const [loading, setLoading] = useState(true);

  // 判断是否为完整的 HTML 页面
  const isFullHtmlPage = contact && (
    contact.includes('<html') || 
    contact.includes('<body') || 
    contact.includes('<!DOCTYPE')
  );

  const loadContact = async () => {
    // 先从缓存加载
    const cached = localStorage.getItem('contact') || '';
    if (cached) {
      setContact(cached);
      setLoading(false);
    }

    try {
      const res = await API.get('/api/contact');
      const { success, message, data } = res.data;
      if (success && data) {
        setContact(data);
        localStorage.setItem('contact', data);
      } else {
        if (!cached) {
          showError(message || t('加载联系页面失败'));
          setContact('');
        }
      }
    } catch (error) {
      if (!cached) {
        showError(t('加载联系页面失败'));
        setContact('');
      }
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadContact();
  }, []);

  if (loading) {
    return (
      <div className='flex justify-center items-center min-h-screen'>
        <Spin size='large' />
      </div>
    );
  }

  // 如果是完整 HTML 页面，使用 iframe 完全隔离
  if (isFullHtmlPage && !loading) {
    return (
      <div className='mt-[60px] px-2'>
        <iframe
          srcDoc={contact}
          style={{ 
            width: '100%', 
            height: 'calc(100vh - 60px)', 
            border: 'none',
            display: 'block'
          }}
          title={t('联系我们')}
        />
      </div>
    );
  }

  if (!contact || contact.trim() === '') {
    return (
      <div className='mt-[60px] px-2'>
        <div className='flex justify-center items-center h-screen p-8'>
          <Empty
            image={
              <IllustrationConstruction style={{ width: 150, height: 150 }} />
            }
            darkModeImage={
              <IllustrationConstructionDark
                style={{ width: 150, height: 150 }}
              />
            }
            description={t('管理员暂时未设置联系页面内容')}
          />
        </div>
      </div>
    );
  }

  return (
    <div className='mt-[60px] px-2'>
      {contact.startsWith('https://') ? (
        <iframe
          src={contact}
          style={{ width: '100%', height: '100vh', border: 'none' }}
          title={t('联系我们')}
        />
      ) : isHtmlContent(contact) ? (
        <SafeHTMLContent
          htmlContent={contact}
          mode='about'
          className='prose prose-lg max-w-none'
          style={{ fontSize: 'larger' }}
        />
      ) : (
        <div className='prose prose-lg max-w-none' style={{ fontSize: 'larger' }}>
          <pre>{contact}</pre>
        </div>
      )}
    </div>
  );
};

export default Contact;
