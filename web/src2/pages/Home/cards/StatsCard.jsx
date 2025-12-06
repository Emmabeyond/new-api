import React from 'react';
import { Typography } from '@douyinfe/semi-ui';
import { IconBolt } from '@douyinfe/semi-icons';
import { useTranslation } from 'react-i18next';
import BentoCard from '../components/BentoCard';
import useCountUp from '../hooks/useCountUp';
import useInView from '../hooks/useInView';

const { Text } = Typography;

const StatsCard = ({ delay = 0 }) => {
  const { t } = useTranslation();
  const [ref, isInView] = useInView({ threshold: 0.3 });

  const apiCalls = useCountUp({
    end: 1000000,
    duration: 2500,
    enabled: isInView,
  });

  const formatNumber = (num) => {
    if (num >= 1000000) {
      return `${(num / 1000000).toFixed(1)}M+`;
    }
    if (num >= 1000) {
      return `${(num / 1000).toFixed(0)}K+`;
    }
    return num.toString();
  };

  return (
    <BentoCard size="medium" delay={delay}>
      <div ref={ref} className="flex flex-col h-full justify-between relative">
        {/* 背景装饰 */}
        <div className="absolute -right-4 -top-4 w-24 h-24 rounded-full bg-gradient-to-br from-blue-500/10 to-cyan-500/10 blur-xl pointer-events-none" />
        
        <div className="flex items-center gap-2 mb-3 relative z-10">
          <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-blue-500/20 to-cyan-500/20 flex items-center justify-center shadow-lg shadow-blue-500/10">
            <IconBolt className="text-blue-500" size="default" />
          </div>
          <Text className="text-semi-color-text-2 text-sm font-medium">
            {t('API 调用')}
          </Text>
        </div>

        <div className="flex-1 flex flex-col items-start justify-center relative z-10">
          <Text className="counter-number text-semi-color-text-0 text-4xl md:text-5xl font-bold">
            {formatNumber(apiCalls)}
          </Text>
          {/* 进度条装饰 */}
          <div className="w-full h-1.5 mt-3 rounded-full bg-semi-color-bg-2 overflow-hidden">
            <div 
              className="h-full rounded-full bg-gradient-to-r from-blue-500 to-cyan-500 transition-all duration-1000"
              style={{ width: isInView ? '85%' : '0%' }}
            />
          </div>
        </div>

        <div className="flex items-center justify-between relative z-10">
          <Text className="text-semi-color-text-2 text-xs">
            {t('累计请求次数')}
          </Text>
          <span className="text-green-500 text-xs font-medium">↑ 12%</span>
        </div>
      </div>
    </BentoCard>
  );
};

export default StatsCard;
