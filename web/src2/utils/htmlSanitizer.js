/**
 * HTML 净化模块
 * 使用 DOMPurify 净化 HTML 内容，防止 XSS 攻击
 */
import DOMPurify from 'dompurify';

// 默认配置：允许常见格式化标签
const DEFAULT_CONFIG = {
  ALLOWED_TAGS: [
    'p',
    'br',
    'strong',
    'em',
    'u',
    'a',
    'ul',
    'ol',
    'li',
    'h1',
    'h2',
    'h3',
    'h4',
    'h5',
    'h6',
    'blockquote',
    'code',
    'pre',
  ],
  ALLOWED_ATTR: ['href', 'target', 'rel', 'class'],
  ALLOW_DATA_ATTR: false,
  FORBID_TAGS: ['script', 'style', 'iframe', 'object', 'embed'],
  FORBID_ATTR: ['onerror', 'onload', 'onclick', 'onmouseover'],
  KEEP_CONTENT: true,
  RETURN_DOM: false,
  RETURN_DOM_FRAGMENT: false,
  RETURN_DOM_IMPORT: false,
  SAFE_FOR_TEMPLATES: false,
};

// Markdown 配置：允许更多标签
const MARKDOWN_CONFIG = {
  ...DEFAULT_CONFIG,
  ALLOWED_TAGS: [
    ...DEFAULT_CONFIG.ALLOWED_TAGS,
    'table',
    'thead',
    'tbody',
    'tr',
    'th',
    'td',
    'img',
  ],
  ALLOWED_ATTR: [...DEFAULT_CONFIG.ALLOWED_ATTR, 'src', 'alt', 'title'],
};

// 严格配置：仅允许基本格式化
const STRICT_CONFIG = {
  ALLOWED_TAGS: ['p', 'br', 'strong', 'em'],
  ALLOWED_ATTR: [],
  KEEP_CONTENT: true,
};

/**
 * 净化 HTML 内容，移除危险标签和属性
 * @param {string} dirtyHTML - 待净化的 HTML 字符串
 * @param {Object} options - 可选配置
 * @returns {string} 净化后的 HTML 字符串
 */
export function sanitizeHTML(dirtyHTML, options = {}) {
  try {
    if (!dirtyHTML || typeof dirtyHTML !== 'string') {
      return '';
    }
    const config = { ...DEFAULT_CONFIG, ...options };
    return DOMPurify.sanitize(dirtyHTML, config);
  } catch (error) {
    console.error('HTML sanitization failed:', error);
    // 返回纯文本版本（移除所有 HTML 标签）
    return dirtyHTML.replace(/<[^>]*>/g, '');
  }
}

/**
 * 净化 Markdown 解析后的 HTML
 * @param {string} markdownHTML - Markdown 转换的 HTML
 * @returns {string} 净化后的 HTML
 */
export function sanitizeMarkdown(markdownHTML) {
  try {
    if (!markdownHTML || typeof markdownHTML !== 'string') {
      return '';
    }
    return DOMPurify.sanitize(markdownHTML, MARKDOWN_CONFIG);
  } catch (error) {
    console.error('Markdown HTML sanitization failed:', error);
    return markdownHTML.replace(/<[^>]*>/g, '');
  }
}

/**
 * 严格模式净化（仅允许纯文本格式化标签）
 * @param {string} dirtyHTML - 待净化的 HTML
 * @returns {string} 净化后的 HTML
 */
export function sanitizeStrict(dirtyHTML) {
  try {
    if (!dirtyHTML || typeof dirtyHTML !== 'string') {
      return '';
    }
    return DOMPurify.sanitize(dirtyHTML, STRICT_CONFIG);
  } catch (error) {
    console.error('Strict HTML sanitization failed:', error);
    return dirtyHTML.replace(/<[^>]*>/g, '');
  }
}

export default {
  sanitizeHTML,
  sanitizeMarkdown,
  sanitizeStrict,
  DEFAULT_CONFIG,
  MARKDOWN_CONFIG,
  STRICT_CONFIG,
};
