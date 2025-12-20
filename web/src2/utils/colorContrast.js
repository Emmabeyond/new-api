/**
 * 颜色对比度工具
 * 用于确保符合 WCAG 2.1 AA 级标准
 */

/**
 * 计算相对亮度
 * @param {number} r - 红色值 (0-255)
 * @param {number} g - 绿色值 (0-255)
 * @param {number} b - 蓝色值 (0-255)
 * @returns {number} 相对亮度 (0-1)
 */
export const getRelativeLuminance = (r, g, b) => {
  const [rs, gs, bs] = [r, g, b].map((c) => {
    const sRGB = c / 255;
    return sRGB <= 0.03928 ? sRGB / 12.92 : Math.pow((sRGB + 0.055) / 1.055, 2.4);
  });
  return 0.2126 * rs + 0.7152 * gs + 0.0722 * bs;
};

/**
 * 计算对比度
 * @param {string} color1 - 颜色1 (hex格式)
 * @param {string} color2 - 颜色2 (hex格式)
 * @returns {number} 对比度 (1-21)
 */
export const getContrastRatio = (color1, color2) => {
  const l1 = getRelativeLuminance(...hexToRgb(color1));
  const l2 = getRelativeLuminance(...hexToRgb(color2));
  const lighter = Math.max(l1, l2);
  const darker = Math.min(l1, l2);
  return (lighter + 0.05) / (darker + 0.05);
};

/**
 * 将 hex 颜色转换为 RGB
 * @param {string} hex - hex 颜色值
 * @returns {number[]} [r, g, b]
 */
export const hexToRgb = (hex) => {
  const result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex);
  return result
    ? [parseInt(result[1], 16), parseInt(result[2], 16), parseInt(result[3], 16)]
    : [0, 0, 0];
};

/**
 * 检查对比度是否符合 WCAG AA 标准
 * @param {string} foreground - 前景色
 * @param {string} background - 背景色
 * @param {boolean} isLargeText - 是否为大文本 (>= 18pt 或 14pt 粗体)
 * @returns {object} { passes: boolean, ratio: number, level: string }
 */
export const checkContrastCompliance = (foreground, background, isLargeText = false) => {
  const ratio = getContrastRatio(foreground, background);
  const requiredRatio = isLargeText ? 3 : 4.5;
  const requiredRatioAAA = isLargeText ? 4.5 : 7;

  return {
    passes: ratio >= requiredRatio,
    ratio: ratio.toFixed(2),
    level: ratio >= requiredRatioAAA ? 'AAA' : ratio >= requiredRatio ? 'AA' : 'Fail',
    required: requiredRatio,
  };
};

/**
 * 等级主题色对比度验证
 * 确保所有主题色在白色背景上都有足够的对比度
 */
export const THEME_COLORS = {
  basic: {
    primary: '#595959',
    background: '#fafafa',
    text: '#262626',
  },
  standard: {
    primary: '#1890ff',
    background: '#e6f7ff',
    text: '#003a8c',
  },
  premium: {
    primary: '#531dab',
    background: '#f9f0ff',
    text: '#22075e',
  },
  elite: {
    primary: '#d48806',
    background: '#fffbe6',
    text: '#874d00',
  },
};

/**
 * 验证所有主题色的对比度
 * @returns {object} 验证结果
 */
export const validateThemeColors = () => {
  const results = {};

  Object.entries(THEME_COLORS).forEach(([theme, colors]) => {
    results[theme] = {
      primaryOnWhite: checkContrastCompliance(colors.primary, '#ffffff'),
      textOnBackground: checkContrastCompliance(colors.text, colors.background),
      primaryOnBackground: checkContrastCompliance(colors.primary, colors.background),
    };
  });

  return results;
};

/**
 * 获取建议的文本颜色（黑色或白色）
 * @param {string} backgroundColor - 背景色
 * @returns {string} '#000000' 或 '#ffffff'
 */
export const getRecommendedTextColor = (backgroundColor) => {
  const luminance = getRelativeLuminance(...hexToRgb(backgroundColor));
  return luminance > 0.5 ? '#000000' : '#ffffff';
};

/**
 * 颜色对比度测试报告
 * 用于开发时验证颜色选择
 */
export const generateContrastReport = () => {
  console.group('🎨 颜色对比度测试报告');

  const themeResults = validateThemeColors();

  Object.entries(themeResults).forEach(([theme, tests]) => {
    console.group(`${theme.toUpperCase()} 主题`);
    Object.entries(tests).forEach(([testName, result]) => {
      const icon = result.passes ? '✅' : '❌';
      console.log(
        `${icon} ${testName}: ${result.ratio}:1 (${result.level}) - 要求: ${result.required}:1`
      );
    });
    console.groupEnd();
  });

  console.groupEnd();
};

// 在开发环境中自动运行测试
if (process.env.NODE_ENV === 'development') {
  // generateContrastReport(); // 取消注释以在开发时查看报告
}
