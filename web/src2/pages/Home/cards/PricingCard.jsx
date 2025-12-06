import React from 'react';
import { Typography } from '@douyinfe/semi-ui';
import { IconPriceTag } from '@douyinfe/semi-icons';
import { useTranslation } from 'react-i18next';
import BentoCard from '../components/BentoCard';

const { Text } = Typography;

const PricingCard = ({ delay = 0 }) => {
  const { t } = useTranslation();

  return (
    <BentoCard size="medium" delay={delay}>
      <div className="flex flex-col h-full justify-between">
        <div className="flex items-center gap-2 mb-3">
          <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-amber-500/20 to-orange-500/20 flex items-center justify-center">
            <IconPriceTag className="text-amber-500" size="small" />
          </div>
          <Text className="text-semi-color-text-2 text-sm font-medium">
            {t('充值比例')}
          </Text>
        </div>

        <div className="flex-1 flex flex-col items-center justify-center gap-2">
          {/* 汇率展示 */}
          <div className="flex items-center gap-3">
            <div className="flex flex-col items-center">
              <Text className="text-semi-color-text-0 text-3xl md:text-4xl font-bold">
                ¥1
              </Text>
              <Text className="text-semi-color-text-2 text-xs mt-1">CNY</Text>
            </div>
            <div className="flex items-center justify-center w-8 h-8 rounded-full bg-gradient-to-r from-amber-500/30 to-orange-500/30">
              <Text className="text-amber-500 text-lg font-bold">=</Text>
            </div>
            <div className="flex flex-col items-center">
              <Text className="text-semi-color-text-0 text-3xl md:text-4xl font-bold">
                $1
              </Text>
              <Text className="text-semi-color-text-2 text-xs mt-1">USD</Text>
            </div>
          </div>
          
          {/* 标签 */}
          <div className="px-3 py-1 rounded-full bg-gradient-to-r from-amber-500/20 to-orange-500/20 animate-pulse">
            <Text className="text-amber-500 text-xs font-semibold">
              {t('超值汇率')}
            </Text>
          </div>
        </div>

        <Text className="text-semi-color-text-2 text-xs text-center">
          {t('充值即享优惠')}
        </Text>
      </div>
    </BentoCard>
  );
};

export default PricingCard;
