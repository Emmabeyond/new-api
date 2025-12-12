/*
Copyright (C) 2025 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/

import React, { useState, useRef, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import {
  Button,
  ButtonGroup,
  Tooltip,
  Typography,
  Divider,
} from '@douyinfe/semi-ui';
import {
  IconBold,
  IconItalic,
  IconLink,
  IconImage,
  IconCode,
  IconList,
  IconH1,
  IconH2,
  IconQuote,
  IconMinus,
  IconEyeOpened,
  IconEdit2,
} from '@douyinfe/semi-icons';
import { Columns } from 'lucide-react';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import rehypeHighlight from 'rehype-highlight';
import 'highlight.js/styles/github-dark.css';

const { Text } = Typography;

const MarkdownEditor = ({ value = '', onChange, height = 400 }) => {
  const { t } = useTranslation();
  const textareaRef = useRef(null);
  const [viewMode, setViewMode] = useState('split'); // 'edit', 'preview', 'split'

  const insertText = useCallback((before, after = '', placeholder = '') => {
    const textarea = textareaRef.current;
    if (!textarea) return;

    const start = textarea.selectionStart;
    const end = textarea.selectionEnd;
    const selectedText = value.substring(start, end) || placeholder;
    const newText = value.substring(0, start) + before + selectedText + after + value.substring(end);
    
    onChange(newText);
    
    // 设置光标位置
    setTimeout(() => {
      textarea.focus();
      const newCursorPos = start + before.length + selectedText.length;
      textarea.setSelectionRange(newCursorPos, newCursorPos);
    }, 0);
  }, [value, onChange]);

  const toolbarItems = [
    {
      icon: <IconH1 />,
      tooltip: t('一级标题'),
      action: () => insertText('# ', '', t('标题')),
    },
    {
      icon: <IconH2 />,
      tooltip: t('二级标题'),
      action: () => insertText('## ', '', t('标题')),
    },
    { type: 'divider' },
    {
      icon: <IconBold />,
      tooltip: t('粗体'),
      action: () => insertText('**', '**', t('粗体文本')),
    },
    {
      icon: <IconItalic />,
      tooltip: t('斜体'),
      action: () => insertText('*', '*', t('斜体文本')),
    },
    { type: 'divider' },
    {
      icon: <IconLink />,
      tooltip: t('链接'),
      action: () => insertText('[', '](url)', t('链接文本')),
    },
    {
      icon: <IconImage />,
      tooltip: t('图片'),
      action: () => insertText('![', '](url)', t('图片描述')),
    },
    { type: 'divider' },
    {
      icon: <IconCode />,
      tooltip: t('代码块'),
      action: () => insertText('```\n', '\n```', t('代码')),
    },
    {
      icon: <IconQuote />,
      tooltip: t('引用'),
      action: () => insertText('> ', '', t('引用文本')),
    },
    { type: 'divider' },
    {
      icon: <IconList />,
      tooltip: t('无序列表'),
      action: () => insertText('- ', '', t('列表项')),
    },
    {
      icon: <IconMinus />,
      tooltip: t('分隔线'),
      action: () => insertText('\n---\n', ''),
    },
  ];

  const renderToolbar = () => (
    <div className='flex items-center justify-between p-2 border-b' style={{ borderColor: 'var(--semi-color-border)' }}>
      <div className='flex items-center gap-1'>
        {toolbarItems.map((item, index) => {
          if (item.type === 'divider') {
            return <Divider layout='vertical' key={index} className='mx-1 h-4' />;
          }
          return (
            <Tooltip content={item.tooltip} key={index}>
              <Button
                theme='borderless'
                type='tertiary'
                size='small'
                icon={item.icon}
                onClick={item.action}
              />
            </Tooltip>
          );
        })}
      </div>
      
      <ButtonGroup size='small'>
        <Tooltip content={t('仅编辑')}>
          <Button
            theme={viewMode === 'edit' ? 'solid' : 'borderless'}
            type={viewMode === 'edit' ? 'primary' : 'tertiary'}
            icon={<IconEdit2 />}
            onClick={() => setViewMode('edit')}
          />
        </Tooltip>
        <Tooltip content={t('分栏视图')}>
          <Button
            theme={viewMode === 'split' ? 'solid' : 'borderless'}
            type={viewMode === 'split' ? 'primary' : 'tertiary'}
            icon={<Columns size={14} />}
            onClick={() => setViewMode('split')}
          />
        </Tooltip>
        <Tooltip content={t('仅预览')}>
          <Button
            theme={viewMode === 'preview' ? 'solid' : 'borderless'}
            type={viewMode === 'preview' ? 'primary' : 'tertiary'}
            icon={<IconEyeOpened />}
            onClick={() => setViewMode('preview')}
          />
        </Tooltip>
      </ButtonGroup>
    </div>
  );

  const renderEditor = () => (
    <textarea
      ref={textareaRef}
      value={value}
      onChange={(e) => onChange(e.target.value)}
      placeholder={t('在此输入 Markdown 内容...')}
      className='w-full h-full p-4 resize-none outline-none font-mono text-sm'
      style={{
        backgroundColor: 'var(--semi-color-bg-0)',
        color: 'var(--semi-color-text-0)',
        border: 'none',
      }}
    />
  );

  const renderPreview = () => (
    <div
      className='w-full h-full p-4 overflow-auto prose prose-semi dark:prose-invert max-w-none'
      style={{ backgroundColor: 'var(--semi-color-bg-1)' }}
    >
      {value ? (
        <ReactMarkdown
          remarkPlugins={[remarkGfm]}
          rehypePlugins={[rehypeHighlight]}
          components={{
            a: ({ node, ...props }) => {
              const isExternal = props.href?.startsWith('http');
              return (
                <a
                  {...props}
                  target={isExternal ? '_blank' : undefined}
                  rel={isExternal ? 'noopener noreferrer' : undefined}
                  className='text-semi-color-primary hover:underline'
                />
              );
            },
            code: ({ node, inline, className, children, ...props }) => {
              if (inline) {
                return (
                  <code className='bg-semi-color-fill-0 px-1.5 py-0.5 rounded text-sm' {...props}>
                    {children}
                  </code>
                );
              }
              return (
                <code className={className} {...props}>
                  {children}
                </code>
              );
            },
            pre: ({ node, ...props }) => (
              <pre className='bg-semi-color-fill-0 rounded-lg p-4 overflow-x-auto' {...props} />
            ),
            table: ({ node, ...props }) => (
              <div className='overflow-x-auto'>
                <table className='w-full border-collapse border border-semi-color-border' {...props} />
              </div>
            ),
            th: ({ node, ...props }) => (
              <th className='border border-semi-color-border bg-semi-color-fill-0 px-4 py-2 text-left' {...props} />
            ),
            td: ({ node, ...props }) => (
              <td className='border border-semi-color-border px-4 py-2' {...props} />
            ),
            img: ({ node, ...props }) => (
              <img className='max-w-full h-auto rounded-lg' {...props} />
            ),
          }}
        >
          {value}
        </ReactMarkdown>
      ) : (
        <Text type='tertiary'>{t('预览区域')}</Text>
      )}
    </div>
  );

  return (
    <div
      className='border rounded-lg overflow-hidden'
      style={{
        borderColor: 'var(--semi-color-border)',
        height: height + 48, // 48px for toolbar
      }}
    >
      {renderToolbar()}
      
      <div className='flex' style={{ height }}>
        {viewMode === 'edit' && (
          <div className='w-full h-full'>{renderEditor()}</div>
        )}
        
        {viewMode === 'preview' && (
          <div className='w-full h-full'>{renderPreview()}</div>
        )}
        
        {viewMode === 'split' && (
          <>
            <div className='w-1/2 h-full border-r' style={{ borderColor: 'var(--semi-color-border)' }}>
              {renderEditor()}
            </div>
            <div className='w-1/2 h-full'>{renderPreview()}</div>
          </>
        )}
      </div>
    </div>
  );
};

export default MarkdownEditor;
