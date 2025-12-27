import React from 'react';
import { useTranslation } from 'react-i18next';
import { IconBolt, IconShield, IconCoinMoneyStroked, IconApps } from '@douyinfe/semi-icons';

const FeaturesSection = () => {
  const { t } = useTranslation();

  // 核心特性数据 - Requirements 3.1, 3.2
  const features = [
    {
      icon: <IconBolt size="large" />,
      title: t('低延迟'),
      value: '<50ms',
      description: t('首字响应'),
      color: 'yellow',
    },
    {
      icon: <IconShield size="large" />,
      title: t('高稳定'),
      value: '99.9%',
      description: t('可用性'),
      color: 'green',
    },
    {
      icon: <IconCoinMoneyStroked size="large" />,
      title: t('价格透明'),
      value: '¥1=$1',
      description: t('超值汇率'),
      color: 'purple',
    },
    {
      icon: <IconApps size="large" />,
      title: t('多模型'),
      value: '30+',
      description: t('供应商支持'),
      color: 'cyan',
    },
  ];

  // 获取颜色样式
  const getColorStyle = (color) => {
    const colorMap = {
      yellow: {
        bg: 'bg-yellow-500/10',
        border: 'border-yellow-500/20',
        text: 'text-yellow-500',
        glow: 'shadow-yellow-500/10',
      },
      green: {
        bg: 'bg-green-500/10',
        border: 'border-green-500/20',
        text: 'text-green-500',
        glow: 'shadow-green-500/10',
      },
      purple: {
        bg: 'bg-purple-500/10',
        border: 'border-purple-500/20',
        text: 'text-purple-500',
        glow: 'shadow-purple-500/10',
      },
      cyan: {
        bg: 'bg-cyan-500/10',
        border: 'border-cyan-500/20',
        text: 'text-cyan-500',
        glow: 'shadow-cyan-500/10',
      },
    };
    return colorMap[color] || colorMap.cyan;
  };

  return (
    <div className="w-full py-8 md:py-16 px-4">
      <div className="max-w-4xl mx-auto">
        {/* 标题区域 */}
        <div className="text-center mb-6 md:mb-10">
          <h2 className="text-xl md:text-3xl font-bold text-semi-color-text-0 mb-2">
            {t('为什么选择我们')}
          </h2>
          <p className="text-semi-color-text-2 text-sm md:text-base">
            {t('开发者关心的核心指标')}
          </p>
        </div>

        {/* 特性卡片网格 - Requirements 3.3 简洁的图标和数据展示 */}
        <div className="grid grid-cols-2 md:grid-cols-4 gap-3 md:gap-4">
          {features.map((feature, index) => {
            const colorStyle = getColorStyle(feature.color);
            return (
              <div
                key={index}
                className={`
                  relative p-4 md:p-6 rounded-2xl 
                  bg-semi-color-bg-1 border border-semi-color-border
                  hover:border-semi-color-border-hover
                  transition-all duration-300 group
                  hover:shadow-lg ${colorStyle.glow}
                `}
              >
                {/* 图标 */}
                <div
                  className={`
                    inline-flex items-center justify-center 
                    w-10 h-10 md:w-12 md:h-12 rounded-xl mb-3 md:mb-4
                    ${colorStyle.bg} ${colorStyle.border} border
                    ${colorStyle.text}
                    group-hover:scale-110 transition-transform duration-300
                  `}
                >
                  {feature.icon}
                </div>

                {/* 数值 */}
                <div className={`text-2xl md:text-3xl font-bold ${colorStyle.text} mb-1`}>
                  {feature.value}
                </div>

                {/* 标题 */}
                <div className="text-semi-color-text-0 font-medium text-sm md:text-base mb-0.5">
                  {feature.title}
                </div>

                {/* 描述 */}
                <div className="text-semi-color-text-3 text-xs md:text-sm">
                  {feature.description}
                </div>
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
};

export default FeaturesSection;
