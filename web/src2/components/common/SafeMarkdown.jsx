/**
 * 安全的 Markdown 渲染组件
 * 使用 react-markdown 库安全渲染 Markdown 内容
 */
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';

/**
 * 安全的 Markdown 渲染组件
 * @param {string} content - Markdown 内容
 * @param {string} className - CSS 类名
 * @param {Object} props - 其他属性
 * @returns {JSX.Element}
 */
export function SafeMarkdown({ content, className, ...props }) {
  if (!content || typeof content !== 'string') {
    return null;
  }

  return (
    <ReactMarkdown
      className={className}
      remarkPlugins={[remarkGfm]}
      components={{
        // 自定义链接渲染，添加 rel="noopener noreferrer"
        a: ({ node, ...linkProps }) => (
          <a {...linkProps} rel='noopener noreferrer' target='_blank' />
        ),
        // 禁用内嵌 HTML 标签
        html: () => null,
      }}
      {...props}
    >
      {content}
    </ReactMarkdown>
  );
}

export default SafeMarkdown;
