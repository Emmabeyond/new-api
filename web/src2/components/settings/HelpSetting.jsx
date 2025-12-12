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
import { Tabs, TabPane, Typography } from '@douyinfe/semi-ui';
import { IconFile, IconFolder } from '@douyinfe/semi-icons';
import HelpDocumentSetting from './HelpDocumentSetting';
import HelpCategorySetting from './HelpCategorySetting';

const { Title } = Typography;

const HelpSetting = () => {
  const { t } = useTranslation();
  const [activeTab, setActiveTab] = useState('documents');

  return (
    <div>
      <Tabs
        activeKey={activeTab}
        onChange={setActiveTab}
        type='line'
      >
        <TabPane
          tab={
            <span className='flex items-center gap-2'>
              <IconFile />
              {t('文档管理')}
            </span>
          }
          itemKey='documents'
        >
          <HelpDocumentSetting />
        </TabPane>
        <TabPane
          tab={
            <span className='flex items-center gap-2'>
              <IconFolder />
              {t('分类管理')}
            </span>
          }
          itemKey='categories'
        >
          <HelpCategorySetting />
        </TabPane>
      </Tabs>
    </div>
  );
};

export default HelpSetting;
