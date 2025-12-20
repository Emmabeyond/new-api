package model

import (
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

// LevelConfig 等级配置
type LevelConfig struct {
	Id                    string `json:"id" gorm:"primaryKey;type:varchar(32)"`           // 等级标识符
	Name                  string `json:"name" gorm:"type:varchar(64)"`                    // 等级名称
	Description           string `json:"description" gorm:"type:varchar(255)"`           // 等级描述
	Priority              int    `json:"priority" gorm:"index"`                           // 优先级（越大越高）
	IsDefault             bool   `json:"is_default" gorm:"default:false;index"`          // 是否默认等级
	MinCumulativeRecharge float64 `json:"min_cumulative_recharge" gorm:"default:0"`      // 最低累计充值（美元）
	Benefits              string `json:"benefits" gorm:"type:text"`                       // JSON 格式的权益配置
	CreatedAt             int64  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt             int64  `json:"updated_at" gorm:"autoUpdateTime"`
}

// LevelBenefits 权益配置结构
type LevelBenefits struct {
	AvailableChannelGroups []string                    `json:"available_channel_groups"` // 可用渠道分组列表
	DiscountRatio          float64                     `json:"discount_ratio"`           // 全局优惠倍率（0-1，1表示无折扣）
	GroupDiscountRatios    map[string]float64          `json:"group_discount_ratios"`    // 针对特定渠道分组的优惠倍率
	RateLimit              *RateLimitConfig            `json:"rate_limit"`               // 全局速率限制配置
	GroupRateLimits        map[string]*RateLimitConfig `json:"group_rate_limits"`        // 针对特定渠道分组的速率限制
	ModelRateLimits        map[string]*RateLimitConfig `json:"model_rate_limits"`        // 针对特定模型的速率限制
}

// RateLimitConfig 速率限制配置
type RateLimitConfig struct {
	TotalCount   int `json:"total_count"`   // 总请求数限制（每分钟），0表示无限制
	SuccessCount int `json:"success_count"` // 成功请求数限制（每分钟）
}

func (LevelConfig) TableName() string {
	return "level_configs"
}

// GetBenefits 解析权益配置
func (l *LevelConfig) GetBenefits() (*LevelBenefits, error) {
	if l.Benefits == "" {
		return &LevelBenefits{
			DiscountRatio: 1.0,
		}, nil
	}
	var benefits LevelBenefits
	err := json.Unmarshal([]byte(l.Benefits), &benefits)
	if err != nil {
		return nil, err
	}
	return &benefits, nil
}

// SetBenefits 设置权益配置
func (l *LevelConfig) SetBenefits(benefits *LevelBenefits) error {
	data, err := json.Marshal(benefits)
	if err != nil {
		return err
	}
	l.Benefits = string(data)
	return nil
}

// GetAllLevelConfigs 获取所有等级配置
func GetAllLevelConfigs() ([]*LevelConfig, error) {
	var levels []*LevelConfig
	err := DB.Order("priority asc").Find(&levels).Error
	return levels, err
}

// GetLevelConfigById 根据ID获取等级配置
func GetLevelConfigById(id string) (*LevelConfig, error) {
	var level LevelConfig
	err := DB.Where("id = ?", id).First(&level).Error
	if err != nil {
		return nil, err
	}
	return &level, nil
}

// GetDefaultLevelConfig 获取默认等级配置
func GetDefaultLevelConfig() (*LevelConfig, error) {
	var level LevelConfig
	err := DB.Where("is_default = ?", true).First(&level).Error
	if err != nil {
		return nil, err
	}
	return &level, nil
}


// GetLevelConfigByPriority 根据优先级获取等级配置
func GetLevelConfigByPriority(priority int) (*LevelConfig, error) {
	var level LevelConfig
	err := DB.Where("priority = ?", priority).First(&level).Error
	if err != nil {
		return nil, err
	}
	return &level, nil
}

// CreateLevelConfig 创建等级配置
func CreateLevelConfig(level *LevelConfig) error {
	return DB.Create(level).Error
}

// UpdateLevelConfig 更新等级配置
func UpdateLevelConfig(level *LevelConfig) error {
	return DB.Save(level).Error
}

// DeleteLevelConfig 删除等级配置
func DeleteLevelConfig(id string) error {
	return DB.Delete(&LevelConfig{}, "id = ?", id).Error
}

// GetLevelUserCount 获取某等级的用户数量
func GetLevelUserCount(levelId string) (int64, error) {
	var count int64
	err := DB.Model(&User{}).Where("level = ?", levelId).Count(&count).Error
	return count, err
}

// GetAllLevelUserStats 获取所有等级的用户统计
func GetAllLevelUserStats() (map[string]int64, error) {
	type Result struct {
		Level string
		Count int64
	}
	var results []Result
	err := DB.Model(&User{}).Select("level, count(*) as count").Group("level").Find(&results).Error
	if err != nil {
		return nil, err
	}
	stats := make(map[string]int64)
	for _, r := range results {
		stats[r.Level] = r.Count
	}
	return stats, nil
}

// ValidateLevelConfig 验证等级配置
func ValidateLevelConfig(level *LevelConfig) error {
	if level.Id == "" {
		return errors.New("等级ID不能为空")
	}
	if level.Name == "" {
		return errors.New("等级名称不能为空")
	}
	if level.Priority < 0 {
		return errors.New("等级优先级不能为负数")
	}
	if level.MinCumulativeRecharge < 0 {
		return errors.New("最低累计充值不能为负数")
	}
	
	// 验证权益配置
	if level.Benefits != "" {
		benefits, err := level.GetBenefits()
		if err != nil {
			return errors.New("权益配置格式错误: " + err.Error())
		}
		if benefits.DiscountRatio < 0 || benefits.DiscountRatio > 10 {
			return errors.New("优惠倍率必须在0-10之间")
		}
		for group, ratio := range benefits.GroupDiscountRatios {
			if ratio < 0 || ratio > 10 {
				return errors.New("分组 " + group + " 的优惠倍率必须在0-10之间")
			}
		}
	}
	return nil
}

// CheckLevelIdExists 检查等级ID是否已存在
func CheckLevelIdExists(id string) (bool, error) {
	var count int64
	err := DB.Model(&LevelConfig{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

// CheckLevelPriorityExists 检查等级优先级是否已存在（排除指定ID）
func CheckLevelPriorityExists(priority int, excludeId string) (bool, error) {
	var count int64
	query := DB.Model(&LevelConfig{}).Where("priority = ?", priority)
	if excludeId != "" {
		query = query.Where("id != ?", excludeId)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

// GetHigherLevelConfig 获取比当前累计充值更高的等级配置
func GetHigherLevelConfig(cumulativeRecharge float64, currentPriority int) (*LevelConfig, error) {
	var level LevelConfig
	err := DB.Where("min_cumulative_recharge <= ? AND priority > ?", cumulativeRecharge, currentPriority).
		Order("priority desc").First(&level).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &level, nil
}

// InitDefaultLevelConfigs 初始化默认等级配置
func InitDefaultLevelConfigs() error {
	// 检查是否已有等级配置
	var count int64
	if err := DB.Model(&LevelConfig{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil // 已有配置，不需要初始化
	}

	// 创建默认等级
	defaultLevels := []*LevelConfig{
		{
			Id:                    "tier_1",
			Name:                  "Tier 1",
			Description:           "基础等级，无需充值",
			Priority:              1,
			IsDefault:             true,
			MinCumulativeRecharge: 0,
			Benefits:              `{"available_channel_groups":["default"],"discount_ratio":1.0,"rate_limit":{"total_count":0,"success_count":0}}`,
		},
		{
			Id:                    "tier_2",
			Name:                  "Tier 2",
			Description:           "累计充值 $50 解锁",
			Priority:              2,
			IsDefault:             false,
			MinCumulativeRecharge: 50,
			Benefits:              `{"available_channel_groups":["default","vip"],"discount_ratio":1.0,"rate_limit":{"total_count":0,"success_count":0}}`,
		},
		{
			Id:                    "tier_3",
			Name:                  "Tier 3",
			Description:           "累计充值 $1000 解锁",
			Priority:              3,
			IsDefault:             false,
			MinCumulativeRecharge: 1000,
			Benefits:              `{"available_channel_groups":["default","vip","premium"],"discount_ratio":0.8,"rate_limit":{"total_count":0,"success_count":0}}`,
		},
	}

	for _, level := range defaultLevels {
		if err := DB.Create(level).Error; err != nil {
			return err
		}
	}
	return nil
}
