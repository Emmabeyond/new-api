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

import React, { useMemo } from 'react';
import { Button, Dropdown } from '@douyinfe/semi-ui';
import { Sun, Moon, Monitor, Palette, ChevronRight } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import {
  useActualTheme,
  useCurrentThemeId,
  usePresetThemes,
  useThemeActions,
} from '../../../context/Theme';

// Color preview dot component
const ColorDot = ({ color, size = 12 }) => (
  <div
    className="rounded-full border border-[var(--semi-color-border)]"
    style={{
      backgroundColor: color,
      width: size,
      height: size,
      minWidth: size,
    }}
  />
);

const ThemeToggle = ({ theme, onThemeToggle, t }) => {
  const navigate = useNavigate();
  const actualTheme = useActualTheme();
  const currentThemeId = useCurrentThemeId();
  const presetThemes = usePresetThemes();
  const { applyTheme } = useThemeActions();

  // Get current theme object
  const currentTheme = useMemo(() => {
    return presetThemes?.find((t) => t.id === currentThemeId) || presetThemes?.[0];
  }, [currentThemeId, presetThemes]);

  // Get quick access themes (first 5)
  const quickThemes = useMemo(() => {
    return presetThemes?.slice(0, 5) || [];
  }, [presetThemes]);

  const themeOptions = useMemo(
    () => [
      {
        key: 'light',
        icon: <Sun size={18} />,
        buttonIcon: <Sun size={18} />,
        label: t('浅色模式'),
        description: t('始终使用浅色主题'),
      },
      {
        key: 'dark',
        icon: <Moon size={18} />,
        buttonIcon: <Moon size={18} />,
        label: t('深色模式'),
        description: t('始终使用深色主题'),
      },
      {
        key: 'auto',
        icon: <Monitor size={18} />,
        buttonIcon: <Monitor size={18} />,
        label: t('自动模式'),
        description: t('跟随系统主题设置'),
      },
    ],
    [t],
  );

  const getItemClassName = (isSelected) =>
    isSelected
      ? '!bg-semi-color-primary-light-default !font-semibold'
      : 'hover:!bg-semi-color-fill-1';

  const currentButtonIcon = useMemo(() => {
    const currentOption = themeOptions.find((option) => option.key === theme);
    return currentOption?.buttonIcon || themeOptions[2].buttonIcon;
  }, [theme, themeOptions]);

  const handleThemeSelect = (themeId) => {
    applyTheme(themeId);
  };

  const handleMoreThemes = () => {
    navigate('/console/theme-center');
  };

  return (
    <Dropdown
      position='bottomRight'
      render={
        <Dropdown.Menu>
          {/* Mode selection section */}
          {themeOptions.map((option) => (
            <Dropdown.Item
              key={option.key}
              icon={option.icon}
              onClick={() => onThemeToggle(option.key)}
              className={getItemClassName(theme === option.key)}
            >
              <div className='flex flex-col'>
                <span>{option.label}</span>
                <span className='text-xs text-semi-color-text-2'>
                  {option.description}
                </span>
              </div>
            </Dropdown.Item>
          ))}

          {theme === 'auto' && (
            <>
              <Dropdown.Divider />
              <div className='px-3 py-2 text-xs text-semi-color-text-2'>
                {t('当前跟随系统')}：
                {actualTheme === 'dark' ? t('深色') : t('浅色')}
              </div>
            </>
          )}

          {/* Theme selection section */}
          <Dropdown.Divider />
          <div className='px-3 py-2 text-xs text-semi-color-text-2 font-medium'>
            {t('theme.center')}
          </div>
          
          {quickThemes.map((themeItem) => (
            <Dropdown.Item
              key={themeItem.id}
              onClick={() => handleThemeSelect(themeItem.id)}
              className={getItemClassName(currentThemeId === themeItem.id)}
            >
              <div className='flex items-center gap-2 w-full'>
                <ColorDot color={themeItem.colors[actualTheme]?.primary} />
                <span className='flex-1'>
                  {t(themeItem.nameKey, { defaultValue: themeItem.name })}
                </span>
                {currentThemeId === themeItem.id && (
                  <span className='text-xs text-semi-color-primary'>
                    {t('theme.current')}
                  </span>
                )}
              </div>
            </Dropdown.Item>
          ))}

          <Dropdown.Divider />
          <Dropdown.Item
            icon={<Palette size={18} />}
            onClick={handleMoreThemes}
            className='hover:!bg-semi-color-fill-1'
          >
            <div className='flex items-center justify-between w-full'>
              <span>{t('theme.moreThemes')}</span>
              <ChevronRight size={16} className='text-semi-color-text-2' />
            </div>
          </Dropdown.Item>
        </Dropdown.Menu>
      }
    >
      <Button
        icon={currentButtonIcon}
        aria-label={t('切换主题')}
        theme='borderless'
        type='tertiary'
        className='!p-1.5 !text-current focus:!bg-semi-color-fill-1 !rounded-full !bg-semi-color-fill-0 hover:!bg-semi-color-fill-1'
      />
    </Dropdown>
  );
};

export default ThemeToggle;
