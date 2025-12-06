import React, { useContext, useState } from 'react';
import { Typography, Button } from '@douyinfe/semi-ui';
import { IconCode, IconCopy } from '@douyinfe/semi-icons';
import { useTranslation } from 'react-i18next';
import BentoCard from '../components/BentoCard';
import { StatusContext } from '../../../context/Status';
import { copy, showSuccess } from '../../../helpers';

const { Text } = Typography;

const CodeExampleCard = ({ delay = 0 }) => {
  const { t } = useTranslation();
  const [statusState] = useContext(StatusContext);
  const [activeTab, setActiveTab] = useState('python');
  const serverAddress = statusState?.status?.server_address || window.location.origin;

  const codeExamples = {
    python: `import openai

client = openai.OpenAI(
    base_url="${serverAddress}/v1",
    api_key="sk-xxx"
)

response = client.chat.completions.create(
    model="gpt-4o",
    messages=[{"role": "user", "content": "Hello!"}]
)`,
    curl: `curl ${serverAddress}/v1/chat/completions \\
  -H "Authorization: Bearer sk-xxx" \\
  -H "Content-Type: application/json" \\
  -d '{
    "model": "gpt-4o",
    "messages": [{"role": "user", "content": "Hello!"}]
  }'`,
    node: `import OpenAI from 'openai';

const client = new OpenAI({
  baseURL: '${serverAddress}/v1',
  apiKey: 'sk-xxx'
});

const response = await client.chat.completions.create({
  model: 'gpt-4o',
  messages: [{ role: 'user', content: 'Hello!' }]
});`,
  };

  const tabs = [
    { key: 'python', label: 'Python' },
    { key: 'node', label: 'Node.js' },
    { key: 'curl', label: 'cURL' },
  ];

  const handleCopy = async () => {
    const ok = await copy(codeExamples[activeTab]);
    if (ok) {
      showSuccess(t('已复制到剪切板'));
    }
  };

  const renderCode = (code) => {
    return code.split('\n').map((line, index) => {
      let formattedLine = line
        // Keywords
        .replace(/\b(import|from|const|await|new)\b/g, '<span class="code-keyword">$1</span>')
        // Strings
        .replace(/(["'])((?:(?!\1)[^\\]|\\.)*)(\1)/g, '<span class="code-string">$1$2$3</span>')
        // Properties/methods
        .replace(/\.([a-zA-Z_][a-zA-Z0-9_]*)/g, '.<span class="code-property">$1</span>')
        // Comments
        .replace(/(#.*)$/g, '<span class="code-comment">$1</span>')
        // curl flags
        .replace(/(-H|-d|curl)/g, '<span class="code-keyword">$1</span>');

      return (
        <div key={index} className="flex hover:bg-white/5 transition-colors">
          <span className="text-semi-color-text-3 w-6 text-right mr-3 select-none opacity-40">
            {index + 1}
          </span>
          <span dangerouslySetInnerHTML={{ __html: formattedLine || '&nbsp;' }} />
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
            <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-violet-500/20 to-purple-500/20 flex items-center justify-center shadow-lg shadow-violet-500/10">
              <IconCode className="text-violet-500" size="default" />
            </div>
            <div>
              <Text className="text-semi-color-text-0 text-sm font-semibold block">
                {t('快速接入')}
              </Text>
              <Text className="text-semi-color-text-2 text-xs">
                {t('兼容 OpenAI SDK')}
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

        {/* 语言切换标签 */}
        <div className="flex gap-1 mb-2">
          {tabs.map((tab) => (
            <button
              key={tab.key}
              onClick={() => setActiveTab(tab.key)}
              className={`
                px-3 py-1 rounded-lg text-xs font-medium transition-all duration-200
                ${activeTab === tab.key 
                  ? 'bg-violet-500/20 text-violet-500' 
                  : 'text-semi-color-text-2 hover:bg-semi-color-bg-2'
                }
              `}
            >
              {tab.label}
            </button>
          ))}
        </div>

        {/* 代码块 */}
        <div className="flex-1 bento-code-block rounded-xl bg-semi-color-bg-2 p-4 overflow-hidden relative">
          {/* 窗口装饰 */}
          <div className="absolute top-3 left-3 flex gap-1.5">
            <div className="w-2.5 h-2.5 rounded-full bg-red-500/60" />
            <div className="w-2.5 h-2.5 rounded-full bg-yellow-500/60" />
            <div className="w-2.5 h-2.5 rounded-full bg-green-500/60" />
          </div>
          
          <pre className="text-semi-color-text-1 whitespace-pre mt-4">
            {renderCode(codeExamples[activeTab])}
          </pre>
        </div>
      </div>
    </BentoCard>
  );
};

export default CodeExampleCard;
