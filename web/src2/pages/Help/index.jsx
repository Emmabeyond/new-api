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

import React, { useEffect, useState, useMemo, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { API, showError, copy } from '../../helpers';
import { Input, Empty, Spin, Typography, Tooltip } from '@douyinfe/semi-ui';
import { IconSearch, IconCopy, IconTick, IconMenu } from '@douyinfe/semi-icons';
import {
  IllustrationNoContent,
  IllustrationNoContentDark,
} from '@douyinfe/semi-illustrations';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import rehypeHighlight from 'rehype-highlight';
import 'highlight.js/styles/github-dark.css';

const { Title, Text } = Typography;

// 递归提取 React 元素中的纯文本
const extractTextFromChildren = (children) => {
  if (typeof children === 'string') {
    return children;
  }
  if (Array.isArray(children)) {
    return children.map(extractTextFromChildren).join('');
  }
  if (children?.props?.children) {
    return extractTextFromChildren(children.props.children);
  }
  return '';
};

// 代码块组件，带复制功能
const CodeBlock = ({ children, className, ...props }) => {
  const { t } = useTranslation();
  const [copied, setCopied] = useState(false);
  
  const handleCopy = useCallback(async () => {
    // 递归提取纯文本内容
    const code = extractTextFromChildren(children).replace(/\n$/, '');
    await copy(code, t('代码已复制'));
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  }, [children, t]);

  // 获取语言类型
  const language = className?.replace('language-', '') || '';

  return (
    <div className='relative group'>
      {/* 语言标签和复制按钮 */}
      <div className='absolute top-0 right-0 flex items-center gap-2 px-3 py-2 opacity-0 group-hover:opacity-100 transition-opacity z-10'>
        {language && (
          <span className='text-xs text-gray-400 uppercase'>{language}</span>
        )}
        <Tooltip content={copied ? t('已复制') : t('复制代码')}>
          <button
            onClick={handleCopy}
            className='p-1.5 rounded-md bg-gray-700 hover:bg-gray-600 text-gray-300 hover:text-white transition-colors'
          >
            {copied ? <IconTick size='small' /> : <IconCopy size='small' />}
          </button>
        </Tooltip>
      </div>
      <code className={className} {...props}>
        {children}
      </code>
    </div>
  );
};

const HelpPage = () => {
  const { t } = useTranslation();
  const [categories, setCategories] = useState([]);
  const [selectedDocId, setSelectedDocId] = useState(null);
  const [selectedDoc, setSelectedDoc] = useState(null);
  const [loading, setLoading] = useState(true);
  const [docLoading, setDocLoading] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  const [isMobile, setIsMobile] = useState(window.innerWidth < 768);
  const [showSidebar, setShowSidebar] = useState(true);


  // 响应式处理
  useEffect(() => {
    const handleResize = () => {
      const mobile = window.innerWidth < 768;
      setIsMobile(mobile);
      if (!mobile) setShowSidebar(true);
    };
    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, []);

  // 加载分类和文档列表
  useEffect(() => {
    loadCategories();
  }, []);

  const loadCategories = async () => {
    try {
      const res = await API.get('/api/help/categories');
      const { success, data, message } = res.data;
      if (success) {
        setCategories(data || []);
        // 默认选中第一个文档
        if (data && data.length > 0 && data[0].documents?.length > 0) {
          const firstDocId = data[0].documents[0].id;
          setSelectedDocId(firstDocId);
          loadDocument(firstDocId);
        }
      } else {
        showError(message || t('加载帮助文档失败'));
      }
    } catch (error) {
      showError(t('加载帮助文档失败'));
    } finally {
      setLoading(false);
    }
  };

  const loadDocument = async (id) => {
    setDocLoading(true);
    try {
      const res = await API.get(`/api/help/documents/${id}`);
      const { success, data, message } = res.data;
      if (success) {
        setSelectedDoc(data);
        if (isMobile) setShowSidebar(false);
      } else {
        showError(message || t('加载文档失败'));
      }
    } catch (error) {
      showError(t('加载文档失败'));
    } finally {
      setDocLoading(false);
    }
  };

  const handleDocSelect = (docId) => {
    setSelectedDocId(docId);
    loadDocument(docId);
  };

  // 搜索过滤
  const filteredCategories = useMemo(() => {
    if (!searchQuery.trim()) return categories;
    
    const query = searchQuery.toLowerCase();
    return categories.map(cat => ({
      ...cat,
      documents: cat.documents.filter(doc => 
        doc.title.toLowerCase().includes(query)
      )
    })).filter(cat => cat.documents.length > 0);
  }, [categories, searchQuery]);

  // 空状态
  if (loading) {
    return (
      <div className='flex justify-center items-center min-h-screen bg-gradient-to-br from-semi-color-bg-0 to-semi-color-bg-1'>
        <Spin size='large' />
      </div>
    );
  }

  if (categories.length === 0) {
    return (
      <div className='flex justify-center items-center min-h-screen p-8 bg-gradient-to-br from-semi-color-bg-0 to-semi-color-bg-1'>
        <Empty
          image={<IllustrationNoContent style={{ width: 150, height: 150 }} />}
          darkModeImage={<IllustrationNoContentDark style={{ width: 150, height: 150 }} />}
          description={t('暂无帮助文档')}
        />
      </div>
    );
  }


  return (
    <div className='flex min-h-screen pt-16 bg-gradient-to-br from-semi-color-bg-0 to-semi-color-bg-1'>
      {/* 侧边栏 */}
      <aside 
        className={`
          ${isMobile ? (showSidebar ? 'fixed inset-0 z-50 pt-16' : 'hidden') : 'sticky top-16 h-[calc(100vh-4rem)]'}
          ${isMobile ? 'w-full' : 'w-72'} 
          flex-shrink-0 
          border-r border-semi-color-border 
          bg-semi-color-bg-0/80 backdrop-blur-sm
          overflow-y-auto
          transition-all duration-300
        `}
      >
        <div className='p-5'>
          {/* 移动端关闭按钮 */}
          {isMobile && (
            <div className='flex justify-between items-center mb-4'>
              <Text strong className='text-lg'>{t('文档目录')}</Text>
              <button
                onClick={() => setShowSidebar(false)}
                className='p-2 rounded-lg hover:bg-semi-color-fill-0 transition-colors'
              >
                ✕
              </button>
            </div>
          )}
          
          <Input
            prefix={<IconSearch className='text-semi-color-text-2' />}
            placeholder={t('搜索文档...')}
            value={searchQuery}
            onChange={setSearchQuery}
            showClear
            className='mb-5'
            style={{ borderRadius: '10px' }}
          />
          
          {filteredCategories.map(category => (
            <div key={category.id} className='mb-6'>
              <Text 
                strong 
                className='text-semi-color-text-2 text-xs uppercase tracking-wider mb-3 block px-2'
                style={{ letterSpacing: '0.1em' }}
              >
                {category.name}
              </Text>
              <div className='space-y-1'>
                {category.documents.map(doc => (
                  <div
                    key={doc.id}
                    onClick={() => handleDocSelect(doc.id)}
                    className={`
                      px-3 py-2.5 rounded-lg cursor-pointer transition-all duration-200
                      ${selectedDocId === doc.id
                        ? 'bg-gradient-to-r from-semi-color-primary-light-default to-semi-color-primary-light-hover text-semi-color-primary shadow-sm'
                        : 'hover:bg-semi-color-fill-0 hover:translate-x-1'
                      }
                    `}
                  >
                    <Text 
                      className={`text-sm ${selectedDocId === doc.id ? 'font-medium text-semi-color-primary' : ''}`}
                    >
                      {doc.title}
                    </Text>
                  </div>
                ))}
              </div>
            </div>
          ))}

          {filteredCategories.length === 0 && searchQuery && (
            <div className='text-center py-12'>
              <Text type='tertiary'>{t('未找到匹配的文档')}</Text>
            </div>
          )}
        </div>
      </aside>


      {/* 内容区 */}
      <main className='flex-1 overflow-y-auto'>
        {/* 移动端菜单按钮 */}
        {isMobile && !showSidebar && (
          <div className='sticky top-16 z-40 p-3 bg-semi-color-bg-0/80 backdrop-blur-sm border-b border-semi-color-border'>
            <button
              onClick={() => setShowSidebar(true)}
              className='flex items-center gap-2 px-4 py-2 rounded-lg bg-semi-color-fill-0 hover:bg-semi-color-fill-1 transition-colors'
            >
              <IconMenu />
              <span>{t('文档目录')}</span>
            </button>
          </div>
        )}
        
        {docLoading ? (
          <div className='flex justify-center items-center h-64'>
            <Spin />
          </div>
        ) : selectedDoc ? (
          <article className='max-w-4xl mx-auto px-6 py-8 md:px-10 md:py-12'>
            {/* 文档标题 */}
            <div className='mb-8 pb-6 border-b border-semi-color-border'>
              <Title heading={1} className='!text-3xl md:!text-4xl !font-bold !mb-3'>
                {selectedDoc.title}
              </Title>
              {selectedDoc.updated_at && (
                <Text type='tertiary' className='text-sm'>
                  {t('最后更新')}: {new Date(selectedDoc.updated_at * 1000).toLocaleDateString()}
                </Text>
              )}
            </div>
            
            {/* Markdown 内容 */}
            <div className='prose prose-semi dark:prose-invert max-w-none prose-headings:scroll-mt-20'>
              <ReactMarkdown
                remarkPlugins={[remarkGfm]}
                rehypePlugins={[rehypeHighlight]}
                components={{
                  // 链接
                  a: ({ node, ...props }) => {
                    const isExternal = props.href?.startsWith('http');
                    return (
                      <a
                        {...props}
                        target={isExternal ? '_blank' : undefined}
                        rel={isExternal ? 'noopener noreferrer' : undefined}
                        className='text-semi-color-primary hover:text-semi-color-primary-hover underline decoration-semi-color-primary/30 hover:decoration-semi-color-primary transition-colors'
                      />
                    );
                  },
                  // 行内代码
                  code: ({ node, inline, className, children, ...props }) => {
                    if (inline) {
                      return (
                        <code 
                          className='bg-semi-color-fill-1 text-semi-color-danger px-1.5 py-0.5 rounded text-sm font-mono' 
                          {...props}
                        >
                          {children}
                        </code>
                      );
                    }
                    return (
                      <CodeBlock className={className} {...props}>
                        {children}
                      </CodeBlock>
                    );
                  },
                  // 代码块容器
                  pre: ({ node, children, ...props }) => (
                    <pre 
                      className='!bg-gray-900 !rounded-xl !p-0 overflow-hidden shadow-lg border border-gray-700/50' 
                      {...props}
                    >
                      <div className='p-4 overflow-x-auto'>
                        {children}
                      </div>
                    </pre>
                  ),
                  // 标题
                  h1: ({ node, ...props }) => (
                    <h1 className='!text-2xl !font-bold !mt-10 !mb-4 !text-semi-color-text-0' {...props} />
                  ),
                  h2: ({ node, ...props }) => (
                    <h2 className='!text-xl !font-semibold !mt-8 !mb-3 !text-semi-color-text-0 !border-b !border-semi-color-border !pb-2' {...props} />
                  ),
                  h3: ({ node, ...props }) => (
                    <h3 className='!text-lg !font-semibold !mt-6 !mb-2 !text-semi-color-text-0' {...props} />
                  ),
                  // 段落
                  p: ({ node, ...props }) => (
                    <p className='!my-4 !leading-7 !text-semi-color-text-1' {...props} />
                  ),
                  // 列表
                  ul: ({ node, ...props }) => (
                    <ul className='!my-4 !pl-6 !list-disc !text-semi-color-text-1' {...props} />
                  ),
                  ol: ({ node, ...props }) => (
                    <ol className='!my-4 !pl-6 !list-decimal !text-semi-color-text-1' {...props} />
                  ),
                  li: ({ node, ...props }) => (
                    <li className='!my-1 !leading-7' {...props} />
                  ),
                  // 引用
                  blockquote: ({ node, ...props }) => (
                    <blockquote 
                      className='!border-l-4 !border-semi-color-primary !bg-semi-color-fill-0 !pl-4 !py-3 !pr-4 !my-4 !rounded-r-lg !italic !text-semi-color-text-2' 
                      {...props} 
                    />
                  ),
                  // 表格
                  table: ({ node, ...props }) => (
                    <div className='overflow-x-auto my-6 rounded-lg border border-semi-color-border'>
                      <table className='w-full border-collapse' {...props} />
                    </div>
                  ),
                  th: ({ node, ...props }) => (
                    <th className='bg-semi-color-fill-0 px-4 py-3 text-left font-semibold text-semi-color-text-0 border-b border-semi-color-border' {...props} />
                  ),
                  td: ({ node, ...props }) => (
                    <td className='px-4 py-3 border-b border-semi-color-border text-semi-color-text-1' {...props} />
                  ),
                  // 图片
                  img: ({ node, ...props }) => (
                    <img className='max-w-full h-auto rounded-lg shadow-md my-6' {...props} />
                  ),
                  // 分隔线
                  hr: ({ node, ...props }) => (
                    <hr className='!my-8 !border-semi-color-border' {...props} />
                  ),
                }}
              >
                {selectedDoc.content}
              </ReactMarkdown>
            </div>
          </article>
        ) : (
          <div className='flex flex-col justify-center items-center h-64 gap-4'>
            <div className='w-16 h-16 rounded-full bg-semi-color-fill-0 flex items-center justify-center'>
              <IconSearch size='extra-large' className='text-semi-color-text-2' />
            </div>
            <Text type='tertiary'>{t('请从左侧选择一篇文档')}</Text>
          </div>
        )}
      </main>
    </div>
  );
};

export default HelpPage;
