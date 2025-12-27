import React, { useContext } from 'react';
import { Button } from '@douyinfe/semi-ui';
import { useTranslation } from 'react-i18next';
import { IconCopy, IconArrowRight } from '@douyinfe/semi-icons';
import { Link } from 'react-router-dom';
import { copy, showSuccess } from '../../../helpers';
import { useIsMobile } from '../../../hooks/common/useIsMobile';
import { StatusContext } from '../../../context/Status';
import ParticlesBackground from './ParticlesBackground';
import { Claude, OpenAI } from '@lobehub/icons';

const DevHeroSection = ({ serverAddress }) => {
  const [statusState] = useContext(StatusContext);
  const docsLink = statusState?.status?.docs_link || '';
  const { t } = useTranslation();
  const isMobile = useIsMobile();

  const handleCopyBaseURL = async () => {
    const ok = await copy(serverAddress);
    if (ok) {
      showSuccess(t('已复制到剪切板'));
    }
  };

  // 支持的工具图标 - Requirements 1.4
  const tools = [
    { name: 'Claude Code', icon: <Claude.Color size={24} /> },
    { name: 'Codex', icon: <OpenAI size={24} /> },
    { name: 'Cursor', icon: <span className="text-lg">⌘</span> },
    { name: 'Windsurf', icon: <span className="text-lg">🏄</span> },
  ];

  return (
    <div className="w-full min-h-[350px] md:min-h-[450px] relative overflow-hidden">
      {/* 粒子背景 */}
      <ParticlesBackground count={isMobile ? 10 : 20} />
      
      {/* 背景模糊晕染 - 终端/代码风格 Requirements 1.1 */}
      <div className="absolute top-1/4 left-1/4 w-64 h-64 bg-purple-500/10 rounded-full blur-3xl pointer-events-none" />
      <div className="absolute bottom-1/4 right-1/4 w-48 h-48 bg-cyan-500/10 rounded-full blur-3xl pointer-events-none" />
      
      <div className="flex items-center justify-center h-full px-4 py-12 md:py-20 relative z-10">
        <div className="flex flex-col items-center justify-center text-center max-w-3xl mx-auto">
          
          {/* 顶部标签 */}
          <div className="mb-4 md:mb-6">
            <span className="inline-flex items-center gap-2 px-3 py-1.5 rounded-full bg-green-500/10 border border-green-500/20 text-xs md:text-sm text-green-500">
              <span className="w-2 h-2 rounded-full bg-green-500 animate-pulse" />
              {t('服务运行中')}
            </span>
          </div>

          {/* 主标题 - Requirements 1.2 简洁标语「为开发者打造」 */}
          <h1 className="text-2xl md:text-4xl lg:text-5xl font-bold text-semi-color-text-0 mb-3 md:mb-4">
            {t('为开发者打造的')}
            <br />
            <span className="bg-gradient-to-r from-purple-500 to-cyan-500 bg-clip-text text-transparent">
              {t('AI 模型网关')}
            </span>
          </h1>

          {/* 副标题 - 工具列表 Requirements 1.4 */}
          <div className="flex items-center justify-center gap-3 md:gap-4 mb-6 md:mb-8 flex-wrap">
            {tools.map((tool) => (
              <div
                key={tool.name}
                className="flex items-center gap-1.5 px-2.5 py-1 rounded-lg bg-semi-color-bg-2 border border-semi-color-border text-semi-color-text-2 text-xs md:text-sm"
              >
                {tool.icon}
                <span>{tool.name}</span>
              </div>
            ))}
          </div>

          {/* Base URL 复制框 - Requirements 1.3 */}
          <div className="w-full max-w-md mb-6 md:mb-8">
            <div className="flex items-center gap-2 p-1 rounded-xl bg-semi-color-bg-1 border border-semi-color-border">
              <code className="flex-1 px-3 py-2 text-sm md:text-base text-semi-color-text-1 font-mono truncate">
                {serverAddress}
              </code>
              <Button
                theme="solid"
                type="primary"
                icon={<IconCopy />}
                onClick={handleCopyBaseURL}
                className="!rounded-lg"
              >
                {t('复制')}
              </Button>
            </div>
            <p className="text-semi-color-text-3 text-xs mt-2">
              {t('将此地址设置为你的 API Base URL')}
            </p>
          </div>

          {/* CTA 按钮 - 获取 API Key / 文档 */}
          <div className="flex flex-row gap-3 justify-center items-center">
            <Link to="/console/token">
              <Button
                theme="solid"
                type="primary"
                size={isMobile ? 'default' : 'large'}
                className="!rounded-xl px-6"
                icon={<IconArrowRight />}
                iconPosition="right"
              >
                {t('获取 API Key')}
              </Button>
            </Link>
            {docsLink && (
              <Button
                size={isMobile ? 'default' : 'large'}
                className="!rounded-xl px-6"
                onClick={() => window.open(docsLink, '_blank')}
              >
                {t('查看文档')}
              </Button>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default DevHeroSection;
