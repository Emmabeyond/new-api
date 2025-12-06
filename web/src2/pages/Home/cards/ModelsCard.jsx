import React from 'react';
import { Typography } from '@douyinfe/semi-ui';
import { IconLayers } from '@douyinfe/semi-icons';
import { useTranslation } from 'react-i18next';
import BentoCard from '../components/BentoCard';

const { Text } = Typography;

const MODELS = [
  { name: 'GPT-5.1', tag: 'OpenAI' },
  { name: 'Claude Opus 4.5', tag: 'Anthropic' },
  { name: 'Gemini 3 Pro', tag: 'Google' },
  { name: 'DeepSeek V3.2', tag: 'DeepSeek' },
];

const ModelsCard = ({ delay = 0 }) => {
  const { t } = useTranslation();

  return (
    <BentoCard size="medium" delay={delay}>
      <div className="flex flex-col h-full justify-between">
        <div className="flex items-center gap-2 mb-3">
          <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-pink-500/20 to-rose-500/20 flex items-center justify-center shadow-lg shadow-pink-500/10">
            <IconLayers className="text-pink-500" size="default" />
          </div>
          <Text className="text-semi-color-text-2 text-sm font-medium">
            {t('热门模型')}
          </Text>
        </div>

        <div className="flex-1 flex flex-col gap-1.5">
          {MODELS.map((model, index) => (
            <div
              key={model.name}
              className="flex items-center justify-between px-2 py-1.5 rounded-lg bg-semi-color-bg-2 hover:bg-semi-color-bg-3 transition-colors"
              style={{ animationDelay: `${index * 100}ms` }}
            >
              <Text className="text-semi-color-text-0 text-xs font-medium">
                {model.name}
              </Text>
              <span className="text-[10px] px-1.5 py-0.5 rounded bg-pink-500/10 text-pink-500">
                {model.tag}
              </span>
            </div>
          ))}
        </div>

        <Text className="text-semi-color-text-2 text-xs text-center mt-2">
          {t('支持 200+ 模型')}
        </Text>
      </div>
    </BentoCard>
  );
};

export default ModelsCard;
