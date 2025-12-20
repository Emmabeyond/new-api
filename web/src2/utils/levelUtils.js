import { Shield, Star, Crown, Gem } from 'lucide-react';

/**
 * 获取等级主题配置
 * @param {number} priority - 等级优先级
 * @returns {object} 主题配置对象
 */
export const getLevelTheme = (priority) => {
  const themes = {
    1: {
      color: 'grey',
      className: 'tier-basic',
      icon: Shield,
      name: 'Basic',
      primary: '#8c8c8c',
      light: '#d9d9d9',
      bg: '#fafafa',
    },
    2: {
      color: 'blue',
      className: 'tier-standard',
      icon: Star,
      name: 'Standard',
      primary: '#3370ff',
      light: '#91c5ff',
      bg: '#f0f7ff',
    },
    3: {
      color: 'purple',
      className: 'tier-premium',
      icon: Crown,
      name: 'Premium',
      primary: '#722ed1',
      light: '#b37feb',
      bg: '#f9f0ff',
    },
    4: {
      color: 'gold',
      className: 'tier-elite',
      icon: Gem,
      name: 'Elite',
      primary: '#faad14',
      light: '#ffd666',
      bg: '#fffbe6',
    },
  };

  return themes[priority] || themes[1];
};

/**
 * 格式化货币金额
 * @param {number} amount - 金额
 * @param {string} currency - 货币符号，默认 '$'
 * @param {string} locale - 语言环境，默认 'en-US'
 * @returns {string} 格式化后的货币字符串
 */
export const formatCurrency = (amount, currency = '$', locale = 'en-US') => {
  if (amount === null || amount === undefined) {
    return `${currency}0.00`;
  }

  // 规范化 locale：将下划线格式转换为连字符格式（如 zh_CN -> zh-CN）
  const normalizedLocale = locale.replace(/_/g, '-');

  const formatted = new Intl.NumberFormat(normalizedLocale, {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(amount);

  return `${currency}${formatted}`;
};

/**
 * 格式化折扣百分比
 * @param {number} ratio - 优惠倍率（0-1之间，1表示无折扣）
 * @param {function} t - i18n翻译函数
 * @returns {string} 格式化后的折扣字符串
 */
export const formatDiscount = (ratio, t = null) => {
  if (ratio === null || ratio === undefined || ratio >= 1) {
    return t ? t('level.format.no_discount') : '无折扣';
  }

  const discount = ((1 - ratio) * 100).toFixed(0);
  const offText = t ? t('level.format.off') : 'OFF';
  return `${discount}% ${offText}`;
};

/**
 * 格式化速率限制
 * @param {object} rateLimit - 速率限制对象 { total_count, success_count }
 * @param {function} t - i18n翻译函数
 * @returns {string} 格式化后的速率限制字符串
 */
export const formatRateLimit = (rateLimit, t = null) => {
  const unlimitedText = t ? t('level.format.unlimited') : '无限制';
  const perMinuteText = t ? t('level.format.per_minute') : '次/分钟';

  if (!rateLimit) {
    return unlimitedText;
  }

  const { total_count, success_count } = rateLimit;

  // 如果两个都是 0 或不存在，表示无限制
  if ((!total_count || total_count === 0) && (!success_count || success_count === 0)) {
    return unlimitedText;
  }

  // 优先显示成功请求数限制
  if (success_count && success_count > 0) {
    return `${formatNumber(success_count)} ${perMinuteText}`;
  }

  // 否则显示总请求数限制
  if (total_count && total_count > 0) {
    return `${formatNumber(total_count)} ${perMinuteText}`;
  }

  return unlimitedText;
};

/**
 * 格式化数字（添加千分位分隔符）
 * @param {number} num - 数字
 * @param {string} locale - 语言环境，默认 'en-US'
 * @returns {string} 格式化后的数字字符串
 */
export const formatNumber = (num, locale = 'en-US') => {
  if (num === null || num === undefined) {
    return '0';
  }

  // 规范化 locale：将下划线格式转换为连字符格式（如 zh_CN -> zh-CN）
  const normalizedLocale = locale.replace(/_/g, '-');

  return new Intl.NumberFormat(normalizedLocale).format(num);
};

/**
 * 解析权益 JSON 字符串
 * @param {string} benefitsJson - 权益 JSON 字符串
 * @returns {object} 解析后的权益对象
 */
export const parseBenefits = (benefitsJson) => {
  if (!benefitsJson) {
    return {
      available_channel_groups: ['default'],
      discount_ratio: 1.0,
      group_discount_ratios: {},
      rate_limit: {
        total_count: 0,
        success_count: 0,
      },
      model_rate_limits: {},
    };
  }

  try {
    if (typeof benefitsJson === 'string') {
      return JSON.parse(benefitsJson);
    }
    return benefitsJson;
  } catch (error) {
    console.error('Failed to parse benefits JSON:', error);
    return {
      available_channel_groups: ['default'],
      discount_ratio: 1.0,
      group_discount_ratios: {},
      rate_limit: {
        total_count: 0,
        success_count: 0,
      },
      model_rate_limits: {},
    };
  }
};

/**
 * 获取进度条颜色
 * @param {number} percent - 进度百分比
 * @returns {string} 颜色值
 */
export const getProgressColor = (percent) => {
  if (percent < 30) {
    return '#ff4d4f'; // 红色
  } else if (percent < 60) {
    return '#faad14'; // 橙色
  } else if (percent < 90) {
    return '#52c41a'; // 绿色
  } else {
    return '#3370ff'; // 蓝色
  }
};

/**
 * 计算还需充值的金额
 * @param {number} current - 当前充值金额
 * @param {number} required - 需要的充值金额
 * @returns {number} 还需充值的金额
 */
export const calculateRemaining = (current, required) => {
  const remaining = required - current;
  return remaining > 0 ? remaining : 0;
};

/**
 * 格式化日期时间
 * @param {Date|string|number} date - 日期对象、ISO字符串或时间戳
 * @param {string} locale - 语言环境，默认 'en-US'
 * @param {object} options - Intl.DateTimeFormat 选项
 * @returns {string} 格式化后的日期字符串
 */
export const formatDateTime = (date, locale = 'en-US', options = {}) => {
  if (!date) {
    return '';
  }

  const dateObj = date instanceof Date ? date : new Date(date);

  if (isNaN(dateObj.getTime())) {
    return '';
  }

  const defaultOptions = {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    ...options,
  };

  // 规范化 locale：将下划线格式转换为连字符格式（如 zh_CN -> zh-CN）
  const normalizedLocale = locale.replace(/_/g, '-');

  return new Intl.DateTimeFormat(normalizedLocale, defaultOptions).format(dateObj);
};

/**
 * 格式化相对时间（如：2天前）
 * @param {Date|string|number} date - 日期对象、ISO字符串或时间戳
 * @param {string} locale - 语言环境，默认 'en-US'
 * @returns {string} 格式化后的相对时间字符串
 */
export const formatRelativeTime = (date, locale = 'en-US') => {
  if (!date) {
    return '';
  }

  const dateObj = date instanceof Date ? date : new Date(date);

  if (isNaN(dateObj.getTime())) {
    return '';
  }

  const now = new Date();
  const diffInSeconds = Math.floor((now - dateObj) / 1000);

  // 规范化 locale：将下划线格式转换为连字符格式（如 zh_CN -> zh-CN）
  const normalizedLocale = locale.replace(/_/g, '-');

  const rtf = new Intl.RelativeTimeFormat(normalizedLocale, { numeric: 'auto' });

  if (diffInSeconds < 60) {
    return rtf.format(-diffInSeconds, 'second');
  } else if (diffInSeconds < 3600) {
    return rtf.format(-Math.floor(diffInSeconds / 60), 'minute');
  } else if (diffInSeconds < 86400) {
    return rtf.format(-Math.floor(diffInSeconds / 3600), 'hour');
  } else if (diffInSeconds < 2592000) {
    return rtf.format(-Math.floor(diffInSeconds / 86400), 'day');
  } else if (diffInSeconds < 31536000) {
    return rtf.format(-Math.floor(diffInSeconds / 2592000), 'month');
  } else {
    return rtf.format(-Math.floor(diffInSeconds / 31536000), 'year');
  }
};

/**
 * 获取指定渠道的折扣倍率
 * 
 * 向后兼容性说明：
 * - 旧配置：只有 discount_ratio 字段（全局折扣），无 group_discount_ratios
 *   行为：所有渠道都使用 discount_ratio 的值
 * - 新配置：有 group_discount_ratios 字段（渠道专属折扣）
 *   行为：优先使用渠道专属折扣，未配置的渠道回退到 discount_ratio
 * - 混合配置：同时存在两个字段
 *   行为：已配置渠道使用专属折扣，未配置渠道使用全局折扣
 * 
 * @param {object} benefits - 等级权益配置
 * @param {string} channelGroup - 渠道分组名称
 * @returns {number} 折扣倍率（0-1）
 */
export const getDiscountRatio = (benefits, channelGroup) => {
  if (!benefits) {
    return 1.0; // 默认无折扣
  }

  // 优先使用渠道专属折扣（新配置）
  if (benefits.group_discount_ratios && benefits.group_discount_ratios[channelGroup]) {
    return benefits.group_discount_ratios[channelGroup];
  }

  // 回退到全局折扣（旧配置兼容）
  if (benefits.discount_ratio && benefits.discount_ratio > 0) {
    return benefits.discount_ratio;
  }

  return 1.0; // 默认无折扣
};

/**
 * 格式化渠道折扣显示
 * @param {number} ratio - 折扣倍率（0-1）
 * @param {function} t - i18n翻译函数（可选）
 * @returns {string} 格式化后的折扣字符串
 */
export const formatDiscountForChannel = (ratio, t = null) => {
  if (ratio === null || ratio === undefined || ratio >= 1.0) {
    return t ? t('level.format.no_discount') : '-';
  }

  // 显示为 "折扣：*0.80" 格式
  const displayValue = ratio.toFixed(2);
  return `折扣：*${displayValue}`;
};


/**
 * 解析渠道分组限流配置
 * @param {object} benefits - 等级权益配置
 * @returns {object} 渠道分组限流配置 { groupName: { total_count, success_count } }
 */
export const parseGroupRateLimits = (benefits) => {
  if (!benefits || !benefits.group_rate_limits) {
    return {};
  }
  return benefits.group_rate_limits;
};

/**
 * 获取指定渠道分组的限流配置
 * 优先级：group_rate_limits > rate_limit > 无限制
 * 
 * @param {object} benefits - 等级权益配置
 * @param {string} channelGroup - 渠道分组名称
 * @returns {object|null} 限流配置 { total_count, success_count } 或 null（无限制）
 */
export const getRateLimitForGroup = (benefits, channelGroup) => {
  if (!benefits) {
    return null;
  }

  // 1. 优先检查渠道分组限流
  if (benefits.group_rate_limits && channelGroup) {
    const groupLimit = benefits.group_rate_limits[channelGroup];
    if (groupLimit && (groupLimit.total_count > 0 || groupLimit.success_count > 0)) {
      return groupLimit;
    }
  }

  // 2. 回退到全局限流
  if (benefits.rate_limit && (benefits.rate_limit.total_count > 0 || benefits.rate_limit.success_count > 0)) {
    return benefits.rate_limit;
  }

  // 3. 无限制
  return null;
};

/**
 * 格式化限流值显示
 * @param {number} value - 限流值
 * @param {function} t - i18n翻译函数（可选）
 * @returns {string} 格式化后的限流字符串
 */
export const formatRateLimitValue = (value, t = null) => {
  const unlimitedText = t ? t('level.format.unlimited') : '无限制';
  
  if (value === null || value === undefined || value === 0) {
    return unlimitedText;
  }
  
  return formatNumber(value);
};

/**
 * 格式化限流配置显示（包含总请求数和成功请求数）
 * @param {object} rateLimit - 限流配置 { total_count, success_count }
 * @param {function} t - i18n翻译函数（可选）
 * @returns {object} 格式化后的限流显示 { total, success }
 */
export const formatRateLimitConfig = (rateLimit, t = null) => {
  const unlimitedText = t ? t('level.format.unlimited') : '无限制';
  
  if (!rateLimit) {
    return { total: unlimitedText, success: unlimitedText };
  }
  
  return {
    total: formatRateLimitValue(rateLimit.total_count, t),
    success: formatRateLimitValue(rateLimit.success_count, t),
  };
};

/**
 * 验证限流值是否有效（非负整数）
 * @param {any} value - 要验证的值
 * @returns {boolean} 是否有效
 */
export const isValidRateLimitValue = (value) => {
  if (value === null || value === undefined || value === '') {
    return true; // 空值视为有效（表示无限制）
  }
  
  const num = Number(value);
  return Number.isInteger(num) && num >= 0;
};
