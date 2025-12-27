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
import { useNavigate } from 'react-router-dom';
import { TextArea, Button, Typography } from '@douyinfe/semi-ui';
import { IconSend } from '@douyinfe/semi-icons';
import { LogIn } from 'lucide-react';
import { showError } from '../../helpers';

const { Text } = Typography;

const MAX_LENGTH = 500;

const SubmitForm = ({
  isLoggedIn = false,
  submitting = false,
  onSubmit,
}) => {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const [content, setContent] = useState('');

  const handleSubmit = async () => {
    const trimmedContent = content.trim();
    
    if (!trimmedContent) {
      showError(t('留言内容不能为空'));
      return;
    }

    if (trimmedContent.length > MAX_LENGTH) {
      showError(t('留言内容不能超过500字'));
      return;
    }

    const success = await onSubmit(trimmedContent);
    if (success) {
      setContent('');
    }
  };

  const handleKeyDown = (e) => {
    // Ctrl/Cmd + Enter 提交
    if ((e.ctrlKey || e.metaKey) && e.key === 'Enter') {
      e.preventDefault();
      if (isLoggedIn && !submitting && content.trim()) {
        handleSubmit();
      }
    }
  };

  // 未登录状态
  if (!isLoggedIn) {
    return (
      <div className='p-6 rounded-xl bg-gradient-to-br from-semi-color-fill-0 to-semi-color-fill-1 border border-semi-color-border'>
        <div className='text-center'>
          <div className='inline-flex items-center justify-center w-12 h-12 rounded-full bg-semi-color-primary-light-default mb-3'>
            <LogIn size={24} className='text-semi-color-primary' />
          </div>
          <Text className='block mb-4'>
            {t('登录后即可发表留言')}
          </Text>
          <Button
            theme='solid'
            type='primary'
            onClick={() => navigate('/login')}
          >
            {t('立即登录')}
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className='p-4 rounded-xl bg-semi-color-bg-1 border border-semi-color-border'>
      <TextArea
        value={content}
        onChange={setContent}
        onKeyDown={handleKeyDown}
        placeholder={t('分享你的想法...')}
        maxCount={MAX_LENGTH}
        showClear
        autosize={{ minRows: 3, maxRows: 6 }}
        disabled={submitting}
      />
      
      <div className='flex items-center justify-between mt-3'>
        <Text type='tertiary' size='small'>
          {t('Ctrl + Enter 快捷提交')}
        </Text>
        
        <Button
          theme='solid'
          type='primary'
          icon={<IconSend />}
          loading={submitting}
          disabled={!content.trim() || submitting}
          onClick={handleSubmit}
        >
          {t('发表留言')}
        </Button>
      </div>
    </div>
  );
};

export default SubmitForm;
