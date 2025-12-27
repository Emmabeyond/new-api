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
import { Typography, Switch, Button } from '@douyinfe/semi-ui';
import { MessageSquareHeart, ExternalLink } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import { API, showError, showSuccess } from '../../../helpers';

const { Title, Text } = Typography;

export default function SettingsGuestbook({ options, refresh }) {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const [panelEnabled, setPanelEnabled] = useState(true);

  useEffect(() => {
    const enabledStr = options['console_setting.guestbook_enabled'];
    setPanelEnabled(
      enabledStr === undefined
        ? true
        : enabledStr === 'true' || enabledStr === true,
    );
  }, [options['console_setting.guestbook_enabled']]);

  const handleToggleEnabled = async (checked) => {
    const newValue = String(checked);
    try {
      const res = await API.put('/api/option/', {
        key: 'console_setting.guestbook_enabled',
        value: newValue,
      });
      if (res.data.success) {
        setPanelEnabled(checked);
        showSuccess(t('设置已保存'));
        refresh?.();
      } else {
        showError(res.data.message || t('保存失败'));
      }
    } catch (error) {
      showError(t('保存失败'));
    }
  };

  const handleManageGuestbook = () => {
    navigate('/console/guestbook-admin');
  };

  return (
    <div>
      {/* 标题栏 */}
      <div className='flex flex-col md:flex-row md:items-center md:justify-between gap-4 mb-4'>
        <div className='order-2 md:order-1 flex items-center gap-2'>
          <MessageSquareHeart size={20} className='text-pink-500' />
          <Title heading={5} style={{ margin: 0 }}>
            {t('精选留言设置')}
          </Title>
        </div>
        {/* 启用开关 */}
        <div className='order-1 md:order-2 flex items-center gap-2'>
          <Switch checked={panelEnabled} onChange={handleToggleEnabled} />
          <Text>{panelEnabled ? t('已启用') : t('已禁用')}</Text>
        </div>
      </div>

      {/* 说明文字 */}
      <div className='mb-4'>
        <Text type='tertiary'>
          {t('启用后，Dashboard 将显示精选留言面板，展示管理员标记的优质用户留言。')}
        </Text>
      </div>

      {/* 管理按钮 */}
      <div className='flex gap-2'>
        <Button
          icon={<ExternalLink size={14} />}
          onClick={handleManageGuestbook}
        >
          {t('管理留言')}
        </Button>
      </div>
    </div>
  );
}
