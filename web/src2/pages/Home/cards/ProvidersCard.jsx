import React, { useState, useEffect } from 'react';
import { Typography } from '@douyinfe/semi-ui';
import { useTranslation } from 'react-i18next';
import BentoCard from '../components/BentoCard';
import {
  Moonshot,
  OpenAI,
  XAI,
  Zhipu,
  Volcengine,
  Cohere,
  Claude,
  Gemini,
  Minimax,
  Wenxin,
  Spark,
  DeepSeek,
  Qwen,
  Grok,
  AzureAI,
  Hunyuan,
} from '@lobehub/icons';

const { Text } = Typography;

const PROVIDERS = [
  { Icon: OpenAI, name: 'OpenAI' },
  { Icon: Claude, Color: Claude.Color, name: 'Claude' },
  { Icon: Gemini, Color: Gemini.Color, name: 'Gemini' },
  { Icon: DeepSeek, Color: DeepSeek.Color, name: 'DeepSeek' },
  { Icon: Qwen, Color: Qwen.Color, name: 'Qwen' },
  { Icon: Zhipu, Color: Zhipu.Color, name: 'Zhipu' },
  { Icon: Grok, name: 'Grok' },
  { Icon: XAI, name: 'xAI' },
  { Icon: Moonshot, name: 'Moonshot' },
  { Icon: Volcengine, Color: Volcengine.Color, name: 'Volcengine' },
  { Icon: Minimax, Color: Minimax.Color, name: 'Minimax' },
  { Icon: Wenxin, Color: Wenxin.Color, name: 'Wenxin' },
  { Icon: Spark, Color: Spark.Color, name: 'Spark' },
  { Icon: Cohere, Color: Cohere.Color, name: 'Cohere' },
  { Icon: AzureAI, Color: AzureAI.Color, name: 'Azure' },
  { Icon: Hunyuan, Color: Hunyuan.Color, name: 'Hunyuan' },
];

const ProvidersCard = ({ delay = 0 }) => {
  const { t } = useTranslation();
  const [hoveredIndex, setHoveredIndex] = useState(null);
  const [animatedCount, setAnimatedCount] = useState(0);

  // 数字动画
  useEffect(() => {
    const timer = setTimeout(() => {
      const interval = setInterval(() => {
        setAnimatedCount((prev) => {
          if (prev >= 30) {
            clearInterval(interval);
            return 30;
          }
          return prev + 1;
        });
      }, 50);
    }, delay + 300);
    return () => clearTimeout(timer);
  }, [delay]);

  return (
    <BentoCard size="large" delay={delay}>
      <div className="flex flex-col h-full">
        {/* 标题区域 */}
        <div className="mb-4 flex items-center justify-between">
          <Text className="text-semi-color-text-2 text-sm font-medium">
            {t('支持的供应商')}
          </Text>
          <span className="px-2 py-0.5 rounded-full bg-indigo-500/10 text-indigo-500 text-xs font-medium">
            {t('持续更新')}
          </span>
        </div>
        
        {/* 图标网格 */}
        <div className="flex-1 flex items-center justify-center relative">
          {/* 背景装饰 */}
          <div className="absolute inset-0 flex items-center justify-center opacity-20 pointer-events-none">
            <div className="w-32 h-32 rounded-full bg-gradient-to-br from-indigo-500/30 to-purple-500/30 blur-2xl" />
          </div>
          
          <div className="grid grid-cols-4 gap-3 md:gap-4 relative z-10">
            {PROVIDERS.map(({ Icon, Color, name }, index) => {
              const IconComponent = Color || Icon;
              const isHovered = hoveredIndex === index;
              return (
                <div
                  key={name}
                  className={`
                    provider-icon-float w-10 h-10 md:w-12 md:h-12 
                    flex items-center justify-center 
                    rounded-xl transition-all duration-300
                    ${isHovered ? 'scale-125 z-20' : 'opacity-70 hover:opacity-100'}
                  `}
                  style={{
                    animationDelay: `${index * 100}ms`,
                    transform: isHovered ? 'scale(1.25) translateY(-4px)' : undefined,
                  }}
                  title={name}
                  onMouseEnter={() => setHoveredIndex(index)}
                  onMouseLeave={() => setHoveredIndex(null)}
                >
                  <IconComponent size={isHovered ? 36 : 32} />
                </div>
              );
            })}
          </div>
        </div>

        {/* 底部统计 */}
        <div className="mt-4 flex items-center justify-between">
          <div className="flex items-baseline gap-1">
            <Text className="text-semi-color-text-0 text-3xl md:text-4xl font-bold counter-number">
              {animatedCount}
            </Text>
            <Text className="text-indigo-500 text-xl font-bold">+</Text>
          </div>
          <Text className="text-semi-color-text-2 text-sm">
            {t('大模型供应商')}
          </Text>
        </div>
      </div>
    </BentoCard>
  );
};

export default ProvidersCard;
