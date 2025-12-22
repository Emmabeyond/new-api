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

import React from 'react';
import { Button, Select } from '@douyinfe/semi-ui';
import CardPro from '../../common/ui/CardPro';
import LogsTable from './UsageLogsTable';
import LogsActions from './UsageLogsActions';
import LogsFilters from './UsageLogsFilters';
import ColumnSelectorModal from './modals/ColumnSelectorModal';
import UserInfoModal from './modals/UserInfoModal';
import { useLogsData } from '../../../hooks/usage-logs/useUsageLogsData';
import { useIsMobile } from '../../../hooks/common/useIsMobile';

const LogsPage = () => {
  const logsData = useLogsData();
  const isMobile = useIsMobile();
  const { t, hasMore, isLoadingMore, loadMoreLogs, pageSize, handlePageSizeChange, logs } = logsData;

  // Custom pagination area for cursor-based pagination
  const paginationArea = (
    <div className="flex items-center justify-between w-full gap-4">
      <span
        className='text-sm select-none'
        style={{ color: 'var(--semi-color-text-2)' }}
      >
        {t('已加载')} {logs.length} {t('条记录')}
      </span>
      <div className="flex items-center gap-3">
        <Select
          value={pageSize}
          onChange={handlePageSizeChange}
          style={{ width: 100 }}
          size={isMobile ? 'small' : 'default'}
        >
          <Select.Option value={10}>10 {t('条/页')}</Select.Option>
          <Select.Option value={20}>20 {t('条/页')}</Select.Option>
          <Select.Option value={50}>50 {t('条/页')}</Select.Option>
          <Select.Option value={100}>100 {t('条/页')}</Select.Option>
        </Select>
        {hasMore && (
          <Button
            onClick={loadMoreLogs}
            loading={isLoadingMore}
            size={isMobile ? 'small' : 'default'}
          >
            {t('加载更多')}
          </Button>
        )}
      </div>
    </div>
  );

  return (
    <>
      {/* Modals */}
      <ColumnSelectorModal {...logsData} />
      <UserInfoModal {...logsData} />

      {/* Main Content */}
      <CardPro
        type='type2'
        statsArea={<LogsActions {...logsData} />}
        searchArea={<LogsFilters {...logsData} />}
        paginationArea={paginationArea}
        t={logsData.t}
      >
        <LogsTable {...logsData} />
      </CardPro>
    </>
  );
};

export default LogsPage;
