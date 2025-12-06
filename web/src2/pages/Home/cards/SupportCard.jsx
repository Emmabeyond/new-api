import React from 'react';
import { Typography } from '@douyinfe/semi-ui';
import { IconCustomerSupport } from '@douyinfe/semi-icons';
import { useTranslation } from 'react-i18next';
import BentoCard from '../components/BentoCard';

const { Text } = Typography;

const SupportCard = ({ delay = 0 }) => {
  const { t } = useTranslation();

  const features = [
    { icon: '📚', text: t('完整文档') },
    { icon: '💬', text: t('在线客服') },
    { icon: '🔧', text: t('技术支持') },
  ];

  return (
    <BentoCard size="medium" delay={delay}>
      <div className="flex flex-col h-full justify-between">
        <div className="flex items-center gap-2 mb-3">
          <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-indigo-500/20 to-blue-500/20 flex items-center justify-center shadow-lg shadow-indigo-500/10">
            <IconCustomerSupport className="text-indigo-500" size="default" />
          </div>
          <Text className="text-semi-color-text-2 text-sm font-medium">
            {t('服务支持')}
          </Text>
        </div>

        <div className="flex-1 flex flex-col justify-center gap-2">
          {features.map((feature, index) => (
            <div
              key={index}
              className="flex items-center gap-2 px-3 py-2 rounded-lg bg-semi-color-bg-2 hover:bg-semi-color-bg-3 transition-colors"
            >
              <span className="text-base">{feature.icon}</span>
              <Text className="text-semi-color-text-1 text-sm">
                {feature.text}
              </Text>
            </div>
          ))}
        </div>

        <div className="flex items-center justify-center gap-1 mt-2">
          <span className="w-2 h-2 rounded-full bg-green-500 animate-pulse" />
          <Text className="text-green-500 text-xs font-medium">
            {t('7×24 在线')}
          </Text>
        </div>
      </div>
    </BentoCard>
  );
};

export default SupportCard;
