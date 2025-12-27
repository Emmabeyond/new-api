import React, { useContext, useState } from 'react';
import { Typography, Button } from '@douyinfe/semi-ui';
import { IconCode, IconCopy, IconTerminal } from '@douyinfe/semi-icons';
import { useTranslation } from 'react-i18next';
import BentoCard from '../components/BentoCard';
import { StatusContext } from '../../../context/Status';
import { copy, showSuccess } from '../../../helpers';
import SafeHTMLContent from '../../../components/common/SafeHTMLContent';

const { Text } = Typography;

const CodeToolsCard = ({ delay = 0 }) => {
  const { t } = useTranslation();
  const [statusState] = useContext(StatusContext);
  const [activeTab, setActiveTab] = useState('claude-code');
  const serverAddress = statusState?.status?.server_address || window.location.origin;

  const toolConfigs = {
    'claude-code': {
      key: 'claude-code',
      label: 'Claude Code',
      icon: <IconCode className="text-purple-500" size="small" />,
      description: t('Anthropic Claude Code 专属渠道'),
      endpoint: `${serverAddress}/v1`,
      models: [
        'claude-3-5-sonnet-20241022',
        'claude-sonnet-4-20250514',
        'claude-opus-4-20250514',
        'claude-3-5-haiku-20241022',
      ],
      configExample: `# Claude Code 配置
# 在终端中设置环境变量：

export ANTHROPIC_BASE_URL="${serverAddress}"
export ANTHROPIC_API_KEY="sk-xxx"

# 或在 Claude Code 设置中配置：
# API Endpoint: ${serverAddress}/v1
# API Key: sk-xxx (你的令牌)`,
    },
    'codex': {
      key: 'codex',
      label: 'Codex',
      icon: <IconTerminal className="text-green-500" size="small" />,
      description: t('OpenAI Codex CLI 专属渠道'),
      endpoint: `${serverAddress}/v1`,
      models: [
        'codex-mini-latest',
        'gpt-4o',
        'gpt-4o-mini',
        'o1',
        'o1-mini',
        'o3-mini',
      ],
      configExample: `# Codex CLI 配置
# 运行配置命令：

codex configure

# 或设置环境变量：

export OPENAI_BASE_URL="${serverAddress}/v1"
export OPENAI_API_KEY="sk-xxx"`,
    },
    'code': {
      key: 'code',
      label: t('通用 Code'),
      icon: <IconCode className="text-blue-500" size="small" />,
      description: t('通用代码工具渠道，兼容 OpenAI 格式'),
      endpoint: `${serverAddress}/v1`,
      models: [t('支持所有 OpenAI 兼容模型')],
      configExample: `# 通用 Code 渠道配置
# 兼容所有 OpenAI SDK 的代码工具

import openai

client = openai.OpenAI(
    base_url="${serverAddress}/v1",
    api_key="sk-xxx"
)

response = client.chat.completions.create(
    model="gpt-4o",
    messages=[{"role": "user", "content": "写一个快速排序"}]
)`,
    },
  };

  const tabs = [
    { key: 'claude-code', label: 'Claude Code', color: 'purple' },
    { key: 'codex', label: 'Codex', color: 'green' },
    { key: 'code', label: t('通用 Code'), color: 'blue' },
  ];

  const currentConfig = toolConfigs[activeTab];

  const handleCopy = async () => {
    const ok = await copy(currentConfig.configExample);
    if (ok) {
      showSuccess(t('已复制到剪切板'));
    }
  };


  const renderCode = (code) => {
    return code.split('\n').map((line, index) => {
      let formattedLine = line
        // Comments
        .replace(/(#.*)$/g, '<span class="code-comment">$1</span>')
        // Keywords
        .replace(/\b(import|from|export|const|await|client)\b/g, '<span class="code-keyword">$1</span>')
        // Strings
        .replace(/(["'])((?:(?!\1)[^\\]|\\.)*)(\1)/g, '<span class="code-string">$1$2$3</span>')
        // Properties/methods
        .replace(/\.([a-zA-Z_][a-zA-Z0-9_]*)/g, '.<span class="code-property">$1</span>');

      return (
        <div key={index} className="flex hover:bg-white/5 transition-colors">
          <span className="text-semi-color-text-3 w-6 text-right mr-3 select-none opacity-40">
            {index + 1}
          </span>
          <SafeHTMLContent htmlContent={formattedLine || '&nbsp;'} mode="default" />
        </div>
      );
    });
  };

  return (
    <BentoCard size="wide" delay={delay}>
      <div className="flex flex-col h-full">
        {/* 头部 */}
        <div className="flex items-center justify-between mb-3">
          <div className="flex items-center gap-2">
            <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-cyan-500/20 to-blue-500/20 flex items-center justify-center shadow-lg shadow-cyan-500/10">
              <IconTerminal className="text-cyan-500" size="default" />
            </div>
            <div>
              <Text className="text-semi-color-text-0 text-sm font-semibold block">
                {t('开发者工具接入')}
              </Text>
              <Text className="text-semi-color-text-2 text-xs">
                {t('Claude Code / Codex 配置教程')}
              </Text>
            </div>
          </div>
          
          {/* 复制按钮 */}
          <Button
            size="small"
            icon={<IconCopy />}
            onClick={handleCopy}
            className="!rounded-lg"
          >
            {t('复制')}
          </Button>
        </div>

        {/* 工具切换标签 */}
        <div className="flex gap-1 mb-2">
          {tabs.map((tab) => (
            <button
              key={tab.key}
              onClick={() => setActiveTab(tab.key)}
              className={`
                px-3 py-1 rounded-lg text-xs font-medium transition-all duration-200
                ${activeTab === tab.key 
                  ? `bg-${tab.color}-500/20 text-${tab.color}-500` 
                  : 'text-semi-color-text-2 hover:bg-semi-color-bg-2'
                }
              `}
            >
              {tab.label}
            </button>
          ))}
        </div>

        {/* 配置信息 */}
        <div className="flex items-center gap-2 mb-2 px-1">
          {currentConfig.icon}
          <Text className="text-semi-color-text-2 text-xs">
            {currentConfig.description}
          </Text>
        </div>

        {/* 支持的模型 */}
        <div className="flex flex-wrap gap-1 mb-2">
          {currentConfig.models.slice(0, 4).map((model, index) => (
            <span
              key={index}
              className="px-2 py-0.5 rounded-md bg-semi-color-bg-2 text-semi-color-text-2 text-xs"
            >
              {model}
            </span>
          ))}
          {currentConfig.models.length > 4 && (
            <span className="px-2 py-0.5 rounded-md bg-semi-color-bg-2 text-semi-color-text-3 text-xs">
              +{currentConfig.models.length - 4}
            </span>
          )}
        </div>

        {/* 代码块 */}
        <div className="flex-1 bento-code-block rounded-xl bg-semi-color-bg-2 p-4 overflow-hidden relative">
          {/* 窗口装饰 */}
          <div className="absolute top-3 left-3 flex gap-1.5">
            <div className="w-2.5 h-2.5 rounded-full bg-red-500/60" />
            <div className="w-2.5 h-2.5 rounded-full bg-yellow-500/60" />
            <div className="w-2.5 h-2.5 rounded-full bg-green-500/60" />
          </div>
          
          <pre className="text-semi-color-text-1 whitespace-pre mt-4 text-xs overflow-x-auto">
            {renderCode(currentConfig.configExample)}
          </pre>
        </div>
      </div>
    </BentoCard>
  );
};

export default CodeToolsCard;
