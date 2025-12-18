/*
Copyright (C) 2025 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.
*/

import React from 'react';
import CardPro from '../common/ui/CardPro';
import TopUpAdminTable from './TopUpAdminTable';
import TopUpAdminFilters from './TopUpAdminFilters';
import TopUpAdminStats from './TopUpAdminStats';
import OrderDetailModal from './modals/OrderDetailModal';
import CompleteOrderModal from './modals/CompleteOrderModal';
import { useTopUpAdminData } from '../../hooks/topup-admin/useTopUpAdminData';
import { useIsMobile } from '../../hooks/common/useIsMobile';
import { createCardProPagination } from '../../helpers/utils';

const TopUpAdmin = () => {
  const topUpAdminData = useTopUpAdminData();
  const isMobile = useIsMobile();

  const {
    // Data
    topups,
    stats,
    loading,
    
    // Pagination
    activePage,
    pageSize,
    total,
    handlePageChange,
    handlePageSizeChange,
    
    // Filters
    filters,
    setFilters,
    handleSearch,
    handleExport,
    
    // Detail modal
    showDetail,
    selectedOrder,
    closeDetail,
    openDetail,
    
    // Complete modal
    showComplete,
    completeOrder,
    closeComplete,
    openComplete,
    handleComplete,
    completeLoading,
    
    // Translation
    t,
  } = topUpAdminData;

  return (
    <>
      <OrderDetailModal
        visible={showDetail}
        order={selectedOrder}
        onClose={closeDetail}
        t={t}
      />

      <CompleteOrderModal
        visible={showComplete}
        order={completeOrder}
        onConfirm={handleComplete}
        onCancel={closeComplete}
        loading={completeLoading}
        t={t}
      />

      <div className='mt-[60px] px-2'>
        <TopUpAdminStats stats={stats} loading={loading} t={t} />
        
        <CardPro
          type='type1'
          actionsArea={
            <div className='flex flex-col md:flex-row justify-between items-center gap-2 w-full'>
              <div className='w-full md:w-full lg:w-auto'>
                <TopUpAdminFilters
                  filters={filters}
                  setFilters={setFilters}
                  onSearch={handleSearch}
                  onExport={handleExport}
                  loading={loading}
                  t={t}
                />
              </div>
            </div>
          }
          paginationArea={createCardProPagination({
            currentPage: activePage,
            pageSize: pageSize,
            total: total,
            onPageChange: handlePageChange,
            onPageSizeChange: handlePageSizeChange,
            isMobile: isMobile,
            t: t,
          })}
          t={t}
        >
          <TopUpAdminTable
            data={topups}
            loading={loading}
            onViewDetail={openDetail}
            onComplete={openComplete}
            t={t}
          />
        </CardPro>
      </div>
    </>
  );
};

export default TopUpAdmin;
