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

import {
  createContext,
  useCallback,
  useContext,
  useState,
  useEffect,
  useRef,
} from 'react';
import {
  presetThemes,
  DEFAULT_THEME_ID,
  getThemeById,
  getThemeColors,
} from '../../config/themes';

// Context definitions
const ThemeContext = createContext(null);
export const useTheme = () => useContext(ThemeContext);

const ActualThemeContext = createContext(null);
export const useActualTheme = () => useContext(ActualThemeContext);

const SetThemeContext = createContext(null);
export const useSetTheme = () => useContext(SetThemeContext);

// New contexts for multi-theme support
const CurrentThemeIdContext = createContext(null);
export const useCurrentThemeId = () => useContext(CurrentThemeIdContext);

const PresetThemesContext = createContext(null);
export const usePresetThemes = () => useContext(PresetThemesContext);

const ThemeActionsContext = createContext(null);
export const useThemeActions = () => useContext(ThemeActionsContext);

// 检测系统主题偏好
const getSystemTheme = () => {
  if (typeof window !== 'undefined' && window.matchMedia) {
    return window.matchMedia('(prefers-color-scheme: dark)').matches
      ? 'dark'
      : 'light';
  }
  return 'light';
};

// LocalStorage key for theme config
const THEME_CONFIG_KEY = 'theme-config';
const THEME_MODE_KEY = 'theme-mode';

// Load theme config from localStorage
const loadThemeConfig = () => {
  try {
    const config = localStorage.getItem(THEME_CONFIG_KEY);
    if (config) {
      const parsed = JSON.parse(config);
      return {
        currentThemeId: parsed.currentThemeId || DEFAULT_THEME_ID,
        currentMode: parsed.currentMode || 'auto',
      };
    }
  } catch (error) {
    console.error('Failed to load theme configuration:', error);
  }
  // Fallback: try to read old theme-mode key for backward compatibility
  try {
    const oldMode = localStorage.getItem(THEME_MODE_KEY);
    return {
      currentThemeId: DEFAULT_THEME_ID,
      currentMode: oldMode || 'auto',
    };
  } catch {
    return {
      currentThemeId: DEFAULT_THEME_ID,
      currentMode: 'auto',
    };
  }
};

// Save theme config to localStorage
const saveThemeConfig = (themeId, mode) => {
  try {
    localStorage.setItem(
      THEME_CONFIG_KEY,
      JSON.stringify({
        version: '1.0',
        currentThemeId: themeId,
        currentMode: mode,
      })
    );
    // Also save to old key for backward compatibility
    localStorage.setItem(THEME_MODE_KEY, mode);
  } catch (error) {
    console.error('Failed to save theme configuration:', error);
  }
};

// Style element ID for theme injection
const THEME_STYLE_ID = 'theme-custom-vars';

// Apply theme colors to CSS variables via style injection
const applyThemeToDOM = (themeId, mode) => {
  const colors = getThemeColors(themeId, mode);
  if (!colors) return;

  const root = document.documentElement;

  // Apply custom theme variables to :root (for custom components)
  root.style.setProperty('--theme-primary', colors.primary);
  root.style.setProperty('--theme-primary-hover', colors.primaryHover);
  root.style.setProperty('--theme-primary-active', colors.primaryActive);
  root.style.setProperty('--theme-primary-light', colors.primaryLight);
  root.style.setProperty('--theme-secondary', colors.secondary);
  root.style.setProperty('--theme-success', colors.success);
  root.style.setProperty('--theme-warning', colors.warning);
  root.style.setProperty('--theme-danger', colors.danger);
  root.style.setProperty('--theme-info', colors.info);
  root.style.setProperty('--theme-bg0', colors.bg0);
  root.style.setProperty('--theme-bg1', colors.bg1);
  root.style.setProperty('--theme-bg2', colors.bg2);
  root.style.setProperty('--theme-bg3', colors.bg3);
  root.style.setProperty('--theme-text0', colors.text0);
  root.style.setProperty('--theme-text1', colors.text1);
  root.style.setProperty('--theme-text2', colors.text2);
  root.style.setProperty('--theme-border', colors.border);
  root.style.setProperty('--theme-fill0', colors.fill0);
  root.style.setProperty('--theme-fill1', colors.fill1);
  root.style.setProperty('--theme-fill2', colors.fill2);

  // Inject CSS to override Semi Design variables with high specificity
  // Semi uses body and body[theme-mode=dark] selectors, we need to match or exceed that
  const cssContent = `
    body,
    body[theme-mode=dark],
    body[theme-mode=light],
    body .semi-always-dark,
    body .semi-always-light {
      --semi-color-primary: ${colors.primary} !important;
      --semi-color-primary-hover: ${colors.primaryHover} !important;
      --semi-color-primary-active: ${colors.primaryActive} !important;
      --semi-color-primary-light-default: ${colors.primaryLight} !important;
      --semi-color-secondary: ${colors.secondary} !important;
      --semi-color-success: ${colors.success} !important;
      --semi-color-warning: ${colors.warning} !important;
      --semi-color-danger: ${colors.danger} !important;
      --semi-color-info: ${colors.info} !important;
      --semi-color-bg-0: ${colors.bg0} !important;
      --semi-color-bg-1: ${colors.bg1} !important;
      --semi-color-bg-2: ${colors.bg2} !important;
      --semi-color-bg-3: ${colors.bg3} !important;
      --semi-color-text-0: ${colors.text0} !important;
      --semi-color-text-1: ${colors.text1} !important;
      --semi-color-text-2: ${colors.text2} !important;
      --semi-color-border: ${colors.border} !important;
      --semi-color-fill-0: ${colors.fill0} !important;
      --semi-color-fill-1: ${colors.fill1} !important;
      --semi-color-fill-2: ${colors.fill2} !important;
    }
  `;

  // Remove existing theme style if present
  let styleEl = document.getElementById(THEME_STYLE_ID);
  if (!styleEl) {
    styleEl = document.createElement('style');
    styleEl.id = THEME_STYLE_ID;
    document.head.appendChild(styleEl);
  }
  styleEl.textContent = cssContent;
};

export const ThemeProvider = ({ children }) => {
  // Load initial config
  const initialConfig = loadThemeConfig();

  // Mode state: 'light' | 'dark' | 'auto'
  const [theme, _setTheme] = useState(initialConfig.currentMode);

  // Current theme ID state
  const [currentThemeId, setCurrentThemeId] = useState(initialConfig.currentThemeId);

  // System theme state
  const [systemTheme, setSystemTheme] = useState(getSystemTheme());

  // Preview state
  const [previewThemeId, setPreviewThemeId] = useState(null);
  const previousThemeIdRef = useRef(null);

  // Compute actual theme mode
  const actualTheme = theme === 'auto' ? systemTheme : theme;

  // Get the effective theme ID (preview or current)
  const effectiveThemeId = previewThemeId || currentThemeId;

  // Monitor system theme changes
  useEffect(() => {
    if (typeof window !== 'undefined' && window.matchMedia) {
      const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');

      const handleSystemThemeChange = (e) => {
        setSystemTheme(e.matches ? 'dark' : 'light');
      };

      mediaQuery.addEventListener('change', handleSystemThemeChange);

      return () => {
        mediaQuery.removeEventListener('change', handleSystemThemeChange);
      };
    }
  }, []);

  // Apply theme mode to DOM
  useEffect(() => {
    const body = document.body;
    if (actualTheme === 'dark') {
      body.setAttribute('theme-mode', 'dark');
      document.documentElement.classList.add('dark');
    } else {
      body.removeAttribute('theme-mode');
      document.documentElement.classList.remove('dark');
    }
  }, [actualTheme]);

  // Apply theme colors to DOM
  useEffect(() => {
    applyThemeToDOM(effectiveThemeId, actualTheme);
  }, [effectiveThemeId, actualTheme]);

  // Set theme mode (light/dark/auto)
  const setTheme = useCallback(
    (newTheme) => {
      let themeValue;

      if (typeof newTheme === 'boolean') {
        // Backward compatibility for boolean parameter
        themeValue = newTheme ? 'dark' : 'light';
      } else if (typeof newTheme === 'string') {
        themeValue = newTheme;
      } else {
        themeValue = 'auto';
      }

      _setTheme(themeValue);
      saveThemeConfig(currentThemeId, themeValue);
    },
    [currentThemeId]
  );

  // Apply a theme by ID
  const applyTheme = useCallback(
    (themeId) => {
      // Validate theme ID
      const themeExists = presetThemes.some((t) => t.id === themeId);
      if (!themeExists) {
        console.warn(`Theme ${themeId} not found, falling back to default`);
        themeId = DEFAULT_THEME_ID;
      }

      // Clear preview if active
      if (previewThemeId) {
        setPreviewThemeId(null);
      }

      setCurrentThemeId(themeId);
      saveThemeConfig(themeId, theme);
    },
    [theme, previewThemeId]
  );

  // Preview a theme temporarily
  const previewTheme = useCallback(
    (themeId) => {
      if (!previewThemeId) {
        previousThemeIdRef.current = currentThemeId;
      }
      setPreviewThemeId(themeId);
    },
    [currentThemeId, previewThemeId]
  );

  // Cancel preview and restore previous theme
  const cancelPreview = useCallback(() => {
    setPreviewThemeId(null);
  }, []);

  // Theme actions object
  const themeActions = {
    applyTheme,
    previewTheme,
    cancelPreview,
  };

  return (
    <SetThemeContext.Provider value={setTheme}>
      <ActualThemeContext.Provider value={actualTheme}>
        <ThemeContext.Provider value={theme}>
          <CurrentThemeIdContext.Provider value={currentThemeId}>
            <PresetThemesContext.Provider value={presetThemes}>
              <ThemeActionsContext.Provider value={themeActions}>
                {children}
              </ThemeActionsContext.Provider>
            </PresetThemesContext.Provider>
          </CurrentThemeIdContext.Provider>
        </ThemeContext.Provider>
      </ActualThemeContext.Provider>
    </SetThemeContext.Provider>
  );
};
