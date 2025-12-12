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
import { Link } from 'react-router-dom';
import { Tag } from '@douyinfe/semi-ui';
import SkeletonWrapper from '../components/SkeletonWrapper';

const HeaderLogo = ({
  isMobile,
  isConsoleRoute,
  logo,
  logoLoaded,
  isLoading,
  systemName,
  isSelfUseMode,
  isDemoSiteMode,
  t,
}) => {
  if (isMobile && isConsoleRoute) {
    return null;
  }

  return (
    <Link to='/' className='group flex items-center gap-3'>
      {/* Logo 容器 - 带光晕效果 */}
      <div className='relative'>
        {/* 光晕背景 */}
        <div 
          className='absolute inset-0 rounded-xl bg-gradient-to-r from-violet-500 via-purple-500 to-pink-500 opacity-0 group-hover:opacity-70 blur-lg transition-all duration-500 scale-150'
          style={{ animation: 'pulse 2s ease-in-out infinite' }}
        />
        {/* Logo 外框 */}
        <div className='relative w-9 h-9 md:w-10 md:h-10 rounded-xl bg-gradient-to-br from-violet-500 via-purple-500 to-pink-500 p-[2px] shadow-lg group-hover:shadow-purple-500/50 transition-all duration-300'>
          <div className='w-full h-full rounded-[10px] bg-white dark:bg-zinc-900 flex items-center justify-center overflow-hidden'>
            <SkeletonWrapper loading={isLoading || !logoLoaded} type='image' />
            <img
              src={logo}
              alt='logo'
              className={`w-7 h-7 md:w-8 md:h-8 object-contain transition-all duration-300 group-hover:scale-110 ${!isLoading && logoLoaded ? 'opacity-100' : 'opacity-0'}`}
            />
          </div>
        </div>
      </div>
      
      {/* 品牌名称 */}
      <div className='hidden md:flex items-center gap-2'>
        <div className='flex items-center gap-2'>
          <SkeletonWrapper
            loading={isLoading}
            type='title'
            width={120}
            height={24}
          >
            {/* 渐变文字效果 */}
            <span 
              className='text-xl font-bold transition-all duration-300 group-hover:tracking-wide'
              style={{
                background: 'linear-gradient(135deg, #8b5cf6 0%, #a855f7 25%, #ec4899 50%, #8b5cf6 75%, #a855f7 100%)',
                backgroundSize: '200% auto',
                WebkitBackgroundClip: 'text',
                WebkitTextFillColor: 'transparent',
                backgroundClip: 'text',
                animation: 'shimmer 3s linear infinite',
              }}
            >
              {systemName}
            </span>
          </SkeletonWrapper>
          {(isSelfUseMode || isDemoSiteMode) && !isLoading && (
            <Tag
              color={isSelfUseMode ? 'purple' : 'blue'}
              className='text-xs px-2 py-0.5 rounded-full whitespace-nowrap shadow-sm border-0'
              size='small'
              style={{
                background: isSelfUseMode 
                  ? 'linear-gradient(135deg, #8b5cf6 0%, #a855f7 100%)' 
                  : 'linear-gradient(135deg, #3b82f6 0%, #6366f1 100%)',
                color: 'white',
              }}
            >
              {isSelfUseMode ? t('自用模式') : t('演示站点')}
            </Tag>
          )}
        </div>
      </div>
      
      {/* CSS 动画 */}
      <style>{`
        @keyframes shimmer {
          0% { background-position: 0% center; }
          50% { background-position: 100% center; }
          100% { background-position: 0% center; }
        }
        @keyframes pulse {
          0%, 100% { opacity: 0; transform: scale(1.5); }
          50% { opacity: 0.5; transform: scale(1.8); }
        }
      `}</style>
    </Link>
  );
};

export default HeaderLogo;
