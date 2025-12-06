import React, { useState, useEffect } from 'react';
import {
  Button,
  Typography,
  Input,
  ScrollList,
  ScrollItem,
} from '@douyinfe/semi-ui';
import { useTranslation } from 'react-i18next';
import {
  IconGithubLogo,
  IconPlay,
  IconFile,
  IconCopy,
} from '@douyinfe/semi-icons';
import { Link } from 'react-router-dom';
import { copy, showSuccess } from '../../../helpers';
import { API_ENDPOINTS } from '../../../constants/common.constant';
import { useIsMobile } from '../../../hooks/common/useIsMobile';
import ParticlesBackground from './ParticlesBackground';

const { Text } = Typography;

const HeroSection = ({
  serverAddress,
  isDemoSiteMode,
  docsLink,
  version,
}) => {
  const { t, i18n } = useTranslation();
  const isMobile = useIsMobile();
  const isChinese = i18n.language.startsWith('zh');
  
  const endpointItems = API_ENDPOINTS.map((e) => ({ value: e }));
  const [endpointIndex, setEndpointIndex] = useState(0);

  useEffect(() => {
    const timer = setInterval(() => {
      setEndpointIndex((prev) => (prev + 1) % endpointItems.length);
    }, 3000);
    return () => clearInterval(timer);
  }, [endpointItems.length]);

  const handleCopyBaseURL = async () => {
    const ok = await copy(serverAddress);
    if (ok) {
      showSuccess(t('已复制到剪切板'));
    }
  };

  return (
    <div className="w-full min-h-[450px] md:min-h-[550px] relative overflow-hidden">
      {/* 粒子背景 */}
      <ParticlesBackground count={isMobile ? 15 : 30} />
      
      {/* 背景模糊晕染球 */}
      <div className="blur-ball blur-ball-indigo" />
      <div className="blur-ball blur-ball-teal" />
      
      <div className="flex items-center justify-center h-full px-4 py-16 md:py-24 relative z-10">
        <div className="flex flex-col items-center justify-center text-center max-w-4xl mx-auto">
          {/* 顶部标签 */}
          <div className="mb-4 md:mb-6">
            <span className="inline-flex items-center gap-2 px-4 py-1.5 rounded-full bg-gradient-to-r from-indigo-500/10 to-purple-500/10 border border-indigo-500/20 text-sm text-semi-color-text-1">
              <span className="w-2 h-2 rounded-full bg-green-500 animate-pulse" />
              {t('服务运行中')}
            </span>
          </div>

          {/* 标题区域 */}
          <div className="flex flex-col items-center justify-center mb-6 md:mb-8">
            <h1
              className={`hero-title-glow text-3xl md:text-5xl lg:text-6xl xl:text-7xl font-bold text-semi-color-text-0 leading-tight ${isChinese ? 'tracking-wide md:tracking-wider' : ''}`}
            >
              {t('统一的')}
              <br />
              <span className="gradient-text-animated">{t('大模型接口网关')}</span>
            </h1>
            
            <p className="text-base md:text-lg lg:text-xl text-semi-color-text-1 mt-4 md:mt-6 max-w-2xl leading-relaxed">
              {t('更好的价格，更好的稳定性，只需要将模型基址替换为：')}
            </p>

            {/* BASE URL 输入框 */}
            <div className="flex flex-col md:flex-row items-center justify-center gap-4 w-full mt-6 md:mt-8 max-w-lg">
              <Input
                readonly
                value={serverAddress}
                className="flex-1 !rounded-full shimmer-border"
                size={isMobile ? 'default' : 'large'}
                suffix={
                  <div className="flex items-center gap-2">
                    <ScrollList
                      bodyHeight={32}
                      style={{ border: 'unset', boxShadow: 'unset' }}
                    >
                      <ScrollItem
                        mode="wheel"
                        cycled={true}
                        list={endpointItems}
                        selectedIndex={endpointIndex}
                        onSelect={({ index }) => setEndpointIndex(index)}
                      />
                    </ScrollList>
                    <Button
                      type="primary"
                      onClick={handleCopyBaseURL}
                      icon={<IconCopy />}
                      className="!rounded-full btn-glow"
                    />
                  </div>
                }
              />
            </div>
          </div>

          {/* CTA 按钮 */}
          <div className="flex flex-row gap-4 justify-center items-center">
            <Link to="/console">
              <Button
                theme="solid"
                type="primary"
                size={isMobile ? 'default' : 'large'}
                className="!rounded-3xl px-8 py-2 btn-glow"
                icon={<IconPlay />}
              >
                {t('获取密钥')}
              </Button>
            </Link>
            {isDemoSiteMode && version ? (
              <Button
                size={isMobile ? 'default' : 'large'}
                className="flex items-center !rounded-3xl px-6 py-2"
                icon={<IconGithubLogo />}
                onClick={() =>
                  window.open('https://github.com/QuantumNous/new-api', '_blank')
                }
              >
                {version}
              </Button>
            ) : (
              docsLink && (
                <Button
                  size={isMobile ? 'default' : 'large'}
                  className="flex items-center !rounded-3xl px-6 py-2"
                  icon={<IconFile />}
                  onClick={() => window.open(docsLink, '_blank')}
                >
                  {t('文档')}
                </Button>
              )
            )}
          </div>

          {/* 底部特性标签 */}
          <div className="flex flex-wrap justify-center gap-3 mt-8 md:mt-12">
            {[
              { icon: '⚡', text: t('低延迟') },
              { icon: '🔒', text: t('安全可靠') },
              { icon: '💰', text: '1¥ = $1' },
            ].map((item, index) => (
              <span
                key={index}
                className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-semi-color-bg-2 text-semi-color-text-2 text-sm"
              >
                <span>{item.icon}</span>
                <span>{item.text}</span>
              </span>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};

export default HeroSection;
