import React, { useState, useEffect } from 'react';
import { Typography } from '@douyinfe/semi-ui';
import { IconClock } from '@douyinfe/semi-icons';
import { useTranslation } from 'react-i18next';
import BentoCard from '../components/BentoCard';
import useInView from '../hooks/useInView';

const { Text } = Typography;

const SpeedCard = ({ delay = 0 }) => {
  const { t } = useTranslation();
  const [ref, isInView] = useInView({ threshold: 0.3 });
  const [latency, setLatency] = useState(0);

  useEffect(() => {
    if (!isInView) return;
    
    const targetLatency = 45;
    const duration = 1500;
    const steps = 30;
    const increment = targetLatency / steps;
    let current = 0;
    
    const timer = setInterval(() => {
      current += increment;
      if (current >= targetLatency) {
        setLatency(targetLatency);
        clearInterval(timer);
      } else {
        setLatency(Math.round(current));
      }
    }, duration / steps);

    return () => clearInterval(timer);
  }, [isInView]);

  return (
    <BentoCard size="medium" delay={delay}>
      <div ref={ref} className="flex flex-col h-full justify-between relative">
        {/* 背景装饰 */}
        <div className="absolute -right-4 -bottom-4 w-20 h-20 rounded-full bg-gradient-to-br from-cyan-500/10 to-teal-500/10 blur-xl pointer-events-none" />
        
        <div className="flex items-center gap-2 mb-3 relative z-10">
          <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-cyan-500/20 to-teal-500/20 flex items-center justify-center shadow-lg shadow-cyan-500/10">
            <IconClock className="text-cyan-500" size="default" />
          </div>
          <Text className="text-semi-color-text-2 text-sm font-medium">
            {t('响应速度')}
          </Text>
        </div>

        <div className="flex-1 flex flex-col items-center justify-center relative z-10">
          <div className="flex items-baseline gap-1">
            <Text className="counter-number text-semi-color-text-0 text-4xl md:text-5xl font-bold">
              {latency}
            </Text>
            <Text className="text-cyan-500 text-lg font-medium">ms</Text>
          </div>
          
          {/* 速度指示器 */}
          <div className="flex items-center gap-1 mt-2">
            {[1, 2, 3, 4, 5].map((i) => (
              <div
                key={i}
                className={`w-1.5 rounded-full transition-all duration-300 ${
                  i <= 4 ? 'bg-cyan-500' : 'bg-semi-color-bg-3'
                }`}
                style={{ height: `${8 + i * 4}px` }}
              />
            ))}
          </div>
        </div>

        <Text className="text-semi-color-text-2 text-xs text-center relative z-10">
          {t('平均首字延迟')}
        </Text>
      </div>
    </BentoCard>
  );
};

export default SpeedCard;
