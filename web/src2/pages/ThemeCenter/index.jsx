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

import React from 'react';
import { Typography, Card } from '@douyinfe/semi-ui';
import { useTranslation } from 'react-i18next';
import { Palette, Check } from 'lucide-react';
import {
  useCurrentThemeId,
  usePresetThemes,
  useThemeActions,
  useActualTheme,
} from '../../context/Theme';

const { Title, Text } = Typography;

// Theme color preview bar component
const ColorPreviewBar = ({ colors, mode }) => {
  const palette = colors[mode];
  const previewColors = [
    palette.primary,
    palette.primaryHover,
    palette.success,
    palette.warning,
    palette.danger,
  ];

  return (
    <div className="flex h-3 w-full rounded-t-md overflow-hidden">
      {previewColors.map((color, index) => (
        <div
          key={index}
          className="flex-1"
          style={{ backgroundColor: color }}
        />
      ))}
    </div>
  );
};

// Theme card component
const ThemeCard = ({ theme, isActive, onApply, onPreview, onCancelPreview, mode }) => {
  const { t } = useTranslation();

  return (
    <Card
      className={`cursor-pointer transition-all duration-200 hover:shadow-lg ${
        isActive ? 'ring-2 ring-[var(--semi-color-primary)]' : ''
      }`}
      bodyStyle={{ padding: 0 }}
      onClick={() => onApply(theme.id)}
      onMouseEnter={() => onPreview(theme.id)}
      onMouseLeave={onCancelPreview}
    >
      <ColorPreviewBar colors={theme.colors} mode={mode} />
      <div className="p-4">
        <div className="flex items-center justify-between mb-2">
          <Text strong className="text-base">
            {t(theme.nameKey, { defaultValue: theme.name })}
          </Text>
          {isActive && (
            <div className="flex items-center gap-1 text-[var(--semi-color-primary)]">
              <Check size={16} />
              <Text type="primary" size="small">
                {t('theme.current')}
              </Text>
            </div>
          )}
        </div>
        <Text type="tertiary" size="small">
          {t(theme.descriptionKey, { defaultValue: theme.description })}
        </Text>
        
        {/* Color palette preview */}
        <div className="flex gap-1 mt-3">
          {['primary', 'success', 'warning', 'danger', 'info'].map((colorKey) => (
            <div
              key={colorKey}
              className="w-5 h-5 rounded-full border border-[var(--semi-color-border)]"
              style={{ backgroundColor: theme.colors[mode][colorKey] }}
              title={colorKey}
            />
          ))}
        </div>
      </div>
    </Card>
  );
};

// Main ThemeCenter page component
const ThemeCenter = () => {
  const { t } = useTranslation();
  const currentThemeId = useCurrentThemeId();
  const presetThemes = usePresetThemes();
  const { applyTheme, previewTheme, cancelPreview } = useThemeActions();
  const actualTheme = useActualTheme();

  return (
    <div className="mt-[60px] px-4 pb-8">
      {/* Page header */}
      <div className="mb-6">
        <div className="flex items-center gap-3 mb-2">
          <Palette size={28} className="text-[var(--semi-color-primary)]" />
          <Title heading={3} style={{ margin: 0 }}>
            {t('theme.center')}
          </Title>
        </div>
        <Text type="tertiary">
          {t('theme.centerDescription')}
        </Text>
      </div>

      {/* Theme grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5 gap-4">
        {presetThemes.map((theme) => (
          <ThemeCard
            key={theme.id}
            theme={theme}
            isActive={currentThemeId === theme.id}
            onApply={applyTheme}
            onPreview={previewTheme}
            onCancelPreview={cancelPreview}
            mode={actualTheme}
          />
        ))}
      </div>
    </div>
  );
};

export default ThemeCenter;
