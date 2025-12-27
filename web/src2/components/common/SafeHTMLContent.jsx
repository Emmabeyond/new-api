/**
 * 安全的 HTML 内容渲染组件
 * 使用 DOMPurify 净化 HTML 内容后再渲染
 */
import {
  sanitizeHTML,
  sanitizeMarkdown,
  sanitizeStrict,
  sanitizeAboutPage,
} from '../../utils/htmlSanitizer';

/**
 * 安全的 HTML 内容渲染组件（使用 DOMPurify）
 * @param {string} htmlContent - HTML 内容
 * @param {string} mode - 净化模式：'default' | 'markdown' | 'strict' | 'about'
 * @param {string} className - CSS 类名
 * @param {Object} props - 其他属性
 * @returns {JSX.Element}
 */
export function SafeHTMLContent({
  htmlContent,
  mode = 'default',
  className,
  ...props
}) {
  if (!htmlContent || typeof htmlContent !== 'string') {
    return null;
  }

  let sanitizedContent;
  switch (mode) {
    case 'markdown':
      sanitizedContent = sanitizeMarkdown(htmlContent);
      break;
    case 'strict':
      sanitizedContent = sanitizeStrict(htmlContent);
      break;
    case 'about':
      sanitizedContent = sanitizeAboutPage(htmlContent);
      break;
    case 'default':
    default:
      sanitizedContent = sanitizeHTML(htmlContent);
      break;
  }

  return (
    <div
      className={className}
      dangerouslySetInnerHTML={{ __html: sanitizedContent }}
      {...props}
    />
  );
}

export default SafeHTMLContent;
