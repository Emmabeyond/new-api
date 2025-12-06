import React from 'react';
import { Typography } from '@douyinfe/semi-ui';
import { IconCheckCircleStroked } from '@douyinfe/semi-icons';
import { useTranslation } from 'react-i18next';
import BentoCard from '../components/BentoCard';

const { Text } = Typography;

const StabilityCard = ({ delay = 0 }) => {
  const { t } = useTranslation();

  return (
    <BentoCard size="medium" delay={delay}>
      <div className="flex flex-col h-full justify-between">
        <div className="flex items-center gap-2 mb-3">
          <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-green-500/20 to-emerald-500/20 flex items-center justify-center">
            <IconCheckCircleStroked className="text-green-500" size="small" />
          </div>
          <Text className="text-semi-color-text-2 text-sm font-medium">
            {t('服务稳定性')}
          </Text>
        </div>

        <div className="flex-1 flex items-center justify-center relative">
          {/* 脉冲动画圈 */}
          <div className="absolute inset-0 flex items-center justify-center">
            <div className="w-16 h-16 rounded-full bg-green-500/10 pulse-animation" />
          </div>
          <div className="absolute inset-0 flex items-center justify-center">
            <div className="w-12 h-12 rounded-full bg-green-500/20 pulse-animation" style={{ animationDelay: '0.5s' }} />
          </div>
          
          <Text className="counter-number text-semi-color-text-0 text-3xl md:text-4xl font-bold relative z-10">
            99.9%
          </Text>
        </div>

        <Text className="text-semi-color-text-2 text-xs text-center">
          {t('可用性保障')}
        </Text>
      </div>
    </BentoCard>
  );
};

export default StabilityCard;
