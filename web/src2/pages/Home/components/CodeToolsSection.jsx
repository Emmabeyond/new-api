import React, { useState, useContext } from 'react';
import { Typography } from '@douyinfe/semi-ui';
import { IconCopy, IconTerminal, IconInfoCircle, IconSetting, IconCode, IconTick, IconHelpCircle } from '@douyinfe/semi-icons';
import { useTranslation } from 'react-i18next';
import { StatusContext } from '../../../context/Status';
import { copy, showSuccess } from '../../../helpers';
import { useIsMobile } from '../../../hooks/common/useIsMobile';
import { Claude, OpenAI } from '@lobehub/icons';

const { Text } = Typography;

// 可复制的代码块组件
const CopyableCode = ({ code, className = '' }) => {
  const { t } = useTranslation();
  const [copied, setCopied] = useState(false);
  
  const handleCopy = async () => {
    const ok = await copy(code);
    if (ok) {
      setCopied(true);
      showSuccess(t('已复制到剪切板'));
      setTimeout(() => setCopied(false), 2000);
    }
  };
  
  return (
    <div className={`relative group ${className}`}>
      <pre className="bg-semi-color-bg-2 rounded-lg p-3 text-xs md:text-sm font-mono text-semi-color-text-1 overflow-x-auto">
        {code}
      </pre>
      <button
        onClick={handleCopy}
        className="absolute top-2 right-2 p-1.5 rounded-md bg-semi-color-bg-3 hover:bg-semi-color-fill-1 transition-colors opacity-0 group-hover:opacity-100"
      >
        {copied ? <IconTick size="small" className="text-green-500" /> : <IconCopy size="small" />}
      </button>
    </div>
  );
};

// 配置项组件
const ConfigItem = ({ label, value, copyable = true }) => {
  const { t } = useTranslation();
  const [copied, setCopied] = useState(false);
  
  const handleCopy = async () => {
    const ok = await copy(value);
    if (ok) {
      setCopied(true);
      showSuccess(t('已复制到剪切板'));
      setTimeout(() => setCopied(false), 2000);
    }
  };
  
  return (
    <div className="flex items-center justify-between py-2 px-3 bg-semi-color-bg-2 rounded-lg mb-2">
      <div className="flex items-center gap-2 min-w-0 flex-1">
        <span className="text-semi-color-text-2 text-sm shrink-0">{label}</span>
        <code className="text-cyan-500 text-sm font-mono truncate">{value}</code>
      </div>
      {copyable && (
        <button
          onClick={handleCopy}
          className="p-1.5 rounded-md hover:bg-semi-color-fill-1 transition-colors shrink-0 ml-2"
        >
          {copied ? <IconTick size="small" className="text-green-500" /> : <IconCopy size="small" />}
        </button>
      )}
    </div>
  );
};

// 步骤列表组件
const StepList = ({ steps }) => {
  return (
    <div className="space-y-2">
      {steps.map((step, index) => (
        <div key={index} className="flex items-start gap-3">
          <span className="flex items-center justify-center w-6 h-6 rounded-full bg-cyan-500/20 text-cyan-500 text-xs font-medium shrink-0">
            {index + 1}
          </span>
          <span className="text-semi-color-text-1 text-sm pt-0.5">{step}</span>
        </div>
      ))}
    </div>
  );
};

// 注意事项组件
const NoteList = ({ notes }) => {
  return (
    <div className="space-y-2">
      {notes.map((note, index) => (
        <div key={index} className="flex items-start gap-2 text-sm">
          <span className="text-yellow-500 shrink-0">•</span>
          <span className="text-semi-color-text-2">{note}</span>
        </div>
      ))}
    </div>
  );
};

const CodeToolsSection = () => {
  const { t } = useTranslation();
  const [statusState] = useContext(StatusContext);
  const [activeTab, setActiveTab] = useState('claude-code');
  const isMobile = useIsMobile();
  
  const serverAddress = statusState?.status?.server_address || window.location.origin;
  const docsLink = statusState?.status?.docs_link || '';

  // 工具配置数据
  const toolConfigs = {
    'claude-code': {
      key: 'claude-code',
      name: 'Claude Code',
      icon: <Claude.Color size={20} />,
      description: t('Anthropic 官方 AI 编程助手'),
      color: 'purple',
      groupNote: t('请确保在 "claude code" 专用分组创建 API Key'),
      sections: [
        {
          title: t('使用前准备'),
          icon: <IconInfoCircle />,
          content: (
            <div className="space-y-3">
              <StepList steps={[
                t('请确保在 "claude code" 专用分组创建 API Key'),
                t('推荐使用 cc-switch 工具快速切换环境'),
              ]} />
              <div className="p-3 bg-purple-500/10 border border-purple-500/20 rounded-lg">
                <Text className="text-purple-400 text-sm">
                  💡 <a href="https://github.com/farion1231/cc-switch/releases" target="_blank" rel="noopener noreferrer" className="underline hover:text-purple-300">cc-switch</a> {t('是一个图形化工具，可以方便地管理多个 Claude Code 配置')}
                </Text>
              </div>
            </div>
          ),
        },
        {
          title: t('终端配置指南'),
          icon: <IconTerminal />,
          content: (
            <div className="space-y-4">
              <div>
                <Text className="text-semi-color-text-2 text-sm mb-2 block">{t('临时设置（当前终端会话有效）')}</Text>
                <CopyableCode code={`export ANTHROPIC_BASE_URL=${serverAddress}\nexport ANTHROPIC_AUTH_TOKEN=your-api-key`} />
              </div>
              <div>
                <Text className="text-semi-color-text-2 text-sm mb-2 block">{t('永久设置（需要重启终端生效）')}</Text>
                <CopyableCode code={`# 添加到 ~/.zshrc 或 ~/.bash_profile\necho 'export ANTHROPIC_BASE_URL=${serverAddress}' >> ~/.zshrc\necho 'export ANTHROPIC_AUTH_TOKEN=your-api-key' >> ~/.zshrc\nsource ~/.zshrc`} />
              </div>
            </div>
          ),
        },
        {
          title: t('注意事项'),
          icon: <IconInfoCircle />,
          content: (
            <NoteList notes={[
              t('请将 your-api-key 替换为您在 claude code 分组创建的实际 API Key'),
              t('永久设置后需要重新打开终端或执行 source 命令才能生效'),
              t('Windows 用户使用 setx 命令后需要重新打开命令提示符'),
            ]} />
          ),
        },
      ],
    },
    'codex': {
      key: 'codex',
      name: 'Codex CLI',
      icon: <OpenAI size={20} />,
      description: t('OpenAI 官方命令行工具'),
      color: 'green',
      groupNote: t('支持 OpenAI 格式调用，推荐使用 codex 分组'),
      sections: [
        {
          title: t('基础配置'),
          icon: <IconSetting />,
          content: (
            <div className="space-y-2">
              <ConfigItem label={t('站点地址')} value={serverAddress} />
              <ConfigItem label="OPENAI_API_KEY" value="sk-xxxxxx" />
              <ConfigItem label="OPENAI_BASE_URL" value={`${serverAddress}/v1`} />
            </div>
          ),
        },
        {
          title: t('使用说明'),
          icon: <IconInfoCircle />,
          content: (
            <StepList steps={[
              t('登录后进入控制台，创建一个新的 API Key'),
              t('推荐在 codex 分组创建 API Key'),
              t('将 Base URL 设置为上述站点地址'),
              t('使用创建的 API Key 作为 OPENAI_API_KEY'),
            ]} />
          ),
        },
        {
          title: t('终端配置'),
          icon: <IconTerminal />,
          content: (
            <div className="space-y-4">
              <div>
                <Text className="text-semi-color-text-2 text-sm mb-2 block">{t('临时设置')}</Text>
                <CopyableCode code={`export OPENAI_BASE_URL="${serverAddress}/v1"\nexport OPENAI_API_KEY="sk-xxx"`} />
              </div>
              <div>
                <Text className="text-semi-color-text-2 text-sm mb-2 block">{t('或运行配置命令')}</Text>
                <CopyableCode code="codex configure" />
              </div>
            </div>
          ),
        },
      ],
    },
    'cursor': {
      key: 'cursor',
      name: 'Cursor',
      icon: <span className="text-base">⌘</span>,
      description: t('AI 驱动的代码编辑器'),
      color: 'cyan',
      groupNote: t('支持 OpenAI 格式调用，可使用 codex 或 claude code 分组'),
      sections: [
        {
          title: t('基础配置'),
          icon: <IconSetting />,
          content: (
            <div className="space-y-2">
              <ConfigItem label="Override OpenAI Base URL" value={`${serverAddress}/v1`} />
              <ConfigItem label="API Key" value="sk-xxxxxx" />
            </div>
          ),
        },
        {
          title: t('配置步骤'),
          icon: <IconCode />,
          content: (
            <StepList steps={[
              t('打开 Cursor'),
              t('进入 Settings → Models → OpenAI API Key'),
              t('填写上述 Base URL 和 API Key'),
              t('保存后即可使用'),
            ]} />
          ),
        },
        {
          title: t('注意事项'),
          icon: <IconInfoCircle />,
          content: (
            <NoteList notes={[
              t('请将 sk-xxxxxx 替换为您创建的实际 API Key'),
              t('推荐在 codex 或 claude code 分组创建 API Key'),
              t('配置完成后可能需要重启 Cursor 才能生效'),
            ]} />
          ),
        },
      ],
    },
  };

  const tabs = [
    { key: 'claude-code', label: 'Claude Code', color: 'purple' },
    { key: 'codex', label: 'Codex', color: 'green' },
    { key: 'cursor', label: 'Cursor', color: 'cyan' },
  ];

  const currentConfig = toolConfigs[activeTab];

  // 获取标签激活样式
  const getTabStyle = (tab) => {
    const isActive = activeTab === tab.key;
    const colorMap = {
      purple: isActive ? 'bg-purple-500/20 text-purple-400 border-purple-500/30' : '',
      green: isActive ? 'bg-green-500/20 text-green-400 border-green-500/30' : '',
      cyan: isActive ? 'bg-cyan-500/20 text-cyan-400 border-cyan-500/30' : '',
    };
    return isActive 
      ? colorMap[tab.color] 
      : 'text-semi-color-text-2 hover:bg-semi-color-bg-2 border-transparent';
  };

  return (
    <div className="w-full py-8 md:py-16 px-4">
      <div className="max-w-4xl mx-auto">
        {/* 标题区域 */}
        <div className="text-center mb-6 md:mb-10">
          <div className="inline-flex items-center gap-2 px-3 py-1.5 rounded-full bg-cyan-500/10 border border-cyan-500/20 text-xs md:text-sm text-cyan-500 mb-4">
            <IconTerminal size="small" />
            {t('快速接入')}
          </div>
          <h2 className="text-xl md:text-3xl font-bold text-semi-color-text-0 mb-2">
            {t('开发者工具配置')}
          </h2>
          <p className="text-semi-color-text-2 text-sm md:text-base">
            {t('选择你的工具，复制配置，立即开始')}
          </p>
        </div>

        {/* 主容器 */}
        <div className="rounded-2xl bg-semi-color-bg-1 border border-semi-color-border overflow-hidden shadow-xl">
          {/* 工具标签切换 */}
          <div className="flex gap-1 p-3 bg-semi-color-bg-2 border-b border-semi-color-border overflow-x-auto">
            {tabs.map((tab) => (
              <button
                key={tab.key}
                onClick={() => setActiveTab(tab.key)}
                className={`
                  flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium 
                  transition-all duration-200 border whitespace-nowrap
                  ${getTabStyle(tab)}
                `}
              >
                {toolConfigs[tab.key].icon}
                <span>{tab.label}</span>
              </button>
            ))}
          </div>

          {/* 工具描述 + 分组提示 */}
          <div className="px-4 py-3 border-b border-semi-color-border bg-semi-color-bg-1/50">
            <div className="flex items-center gap-2 mb-2">
              {currentConfig.icon}
              <Text className="text-semi-color-text-1 font-medium">
                {currentConfig.name}
              </Text>
              <span className="text-semi-color-text-3">—</span>
              <Text className="text-semi-color-text-2 text-sm">
                {currentConfig.description}
              </Text>
            </div>
            {currentConfig.groupNote && (
              <div className="flex items-center gap-2 px-3 py-2 bg-yellow-500/10 border border-yellow-500/20 rounded-lg">
                <IconInfoCircle className="text-yellow-500 shrink-0" size="small" />
                <Text className="text-yellow-600 dark:text-yellow-400 text-xs">
                  {t('注意')}：{currentConfig.groupNote}
                </Text>
              </div>
            )}
          </div>

          {/* 配置内容区域 - 分区卡片展示 */}
          <div className={`p-4 space-y-4 overflow-y-auto ${isMobile ? 'max-h-[500px]' : 'max-h-[600px]'}`}>
            {currentConfig.sections.map((section, index) => (
              <div key={index} className="rounded-xl bg-semi-color-bg-2/50 border border-semi-color-border overflow-hidden">
                {/* 区块标题 */}
                <div className="flex items-center gap-2 px-4 py-3 bg-semi-color-bg-2 border-b border-semi-color-border">
                  <span className="text-cyan-500">{section.icon}</span>
                  <Text className="text-semi-color-text-1 font-medium text-sm">
                    {section.title}
                  </Text>
                </div>
                {/* 区块内容 */}
                <div className="p-4">
                  {section.content}
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* 底部提示 - 遇到问题查看文档 */}
        <div className="mt-6 p-4 rounded-xl bg-semi-color-bg-1 border border-semi-color-border">
          <div className="flex items-start gap-3">
            <div className="p-2 rounded-lg bg-blue-500/10 shrink-0">
              <IconHelpCircle className="text-blue-500" />
            </div>
            <div className="flex-1">
              <Text className="text-semi-color-text-1 font-medium text-sm block mb-1">
                {t('遇到问题？')}
              </Text>
              <Text className="text-semi-color-text-2 text-xs block mb-2">
                {t('如果按照上述步骤配置后仍有问题，请查看我们的详细文档获取更多帮助。')}
              </Text>
              {docsLink && (
                <a
                  href={docsLink}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-blue-500/10 hover:bg-blue-500/20 border border-blue-500/20 text-blue-500 text-xs font-medium transition-colors"
                >
                  📖 {t('查看完整文档')}
                </a>
              )}
            </div>
          </div>
        </div>

        {/* 底部小提示 */}
        <div className="text-center mt-4">
          <Text className="text-semi-color-text-3 text-xs">
            {t('在控制台获取令牌')}
          </Text>
        </div>
      </div>
    </div>
  );
};

export default CodeToolsSection;
