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

// About 页面配置：完全宽松配置
// 因为 About 页面只有管理员可以设置，完全信任管理员输入的内容
const ABOUT_PAGE_CONFIG = {
  FORBID_TAGS: [], // 不禁止任何标签（包括 script）
  FORBID_ATTR: [], // 不禁止任何属性
  ALLOW_DATA_ATTR: true, // 允许 data-* 属性
  ALLOW_UNKNOWN_PROTOCOLS: false, // 不允许未知协议（安全考虑）
  ADD_TAGS: ['style', 'script'], // 明确添加 style 和 script 标签
  ADD_ATTR: ['style', 'onclick', 'onload', 'onerror'], // 明确添加事件属性
  KEEP_CONTENT: true,
  RETURN_DOM: false,
  RETURN_DOM_FRAGMENT: false,
  RETURN_DOM_IMPORT: false,
  SAFE_FOR_TEMPLATES: false,
  WHOLE_DOCUMENT: false,
  FORCE_BODY: false,
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

/**
 * 净化 About 页面的 HTML 内容
 * 管理员内容完全可信，直接返回原始内容
 * @param {string} aboutHTML - About 页面的 HTML
 * @returns {string} 原始 HTML（不净化）
 */
export function sanitizeAboutPage(aboutHTML) {
  if (!aboutHTML || typeof aboutHTML !== 'string') {
    return '';
  }
  // 管理员内容完全可信，直接返回原始内容
  return aboutHTML;
}

export default {
  sanitizeHTML,
  sanitizeMarkdown,
  sanitizeStrict,
  sanitizeAboutPage,
  DEFAULT_CONFIG,
  MARKDOWN_CONFIG,
  STRICT_CONFIG,
  ABOUT_PAGE_CONFIG,
};
