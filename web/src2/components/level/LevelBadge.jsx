import React from 'react';
import { Shield, Star, Crown, Gem } from './icons';
import './LevelBadge.css';

/**
 * 等级徽章组件
 * 根据等级优先级显示不同的图标和主题色
 */
const LevelBadge = ({ level, size = 'medium', theme }) => {
  if (!level) return null;

  const sizeMap = {
    small: { icon: 16, badge: 24 },
    medium: { icon: 24, badge: 40 },
    large: { icon: 48, badge: 80 },
  };

  const levelTheme = theme || getLevelTheme(level.priority);
  const IconComponent = levelTheme.icon;
  const sizes = sizeMap[size];

  return (
    <div
      className={`level-badge level-badge-${size} ${levelTheme.className}`}
      style={{
        width: sizes.badge,
        height: sizes.badge,
      }}
      role="img"
      aria-label={`${level.name} 等级徽章`}
    >
      <div className="badge-circle">
        <IconComponent size={sizes.icon} aria-hidden="true" />
      </div>
    </div>
  );
};

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
    },
    2: {
      color: 'blue',
      className: 'tier-standard',
      icon: Star,
      name: 'Standard',
    },
    3: {
      color: 'purple',
      className: 'tier-premium',
      icon: Crown,
      name: 'Premium',
    },
    4: {
      color: 'gold',
      className: 'tier-elite',
      icon: Gem,
      name: 'Elite',
    },
  };

  return themes[priority] || themes[1];
};

export default LevelBadge;
