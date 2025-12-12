package model

import (
	"errors"
)

const (
	HelpCategoryStatusEnabled  = 1
	HelpCategoryStatusDisabled = 2
)

// HelpCategory 帮助文档分类
type HelpCategory struct {
	Id        int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string `json:"name" gorm:"type:varchar(100);not null"`
	SortOrder int    `json:"sort_order" gorm:"default:0"`
	Status    int    `json:"status" gorm:"default:1"` // 1=启用, 2=禁用
	CreatedAt int64  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt int64  `json:"updated_at" gorm:"autoUpdateTime"`
}

// HelpCategoryWithCount 带文档数量的分类
type HelpCategoryWithCount struct {
	HelpCategory
	DocumentCount int64 `json:"document_count"`
}

func (HelpCategory) TableName() string {
	return "help_categories"
}

// GetAllHelpCategories 获取所有启用的分类（按排序权重）
func GetAllHelpCategories() ([]HelpCategory, error) {
	var categories []HelpCategory
	err := DB.Where("status = ?", HelpCategoryStatusEnabled).
		Order("sort_order ASC, id ASC").
		Find(&categories).Error
	return categories, err
}

// GetAllHelpCategoriesAdmin 管理员获取所有分类（包含禁用的）
func GetAllHelpCategoriesAdmin() ([]HelpCategoryWithCount, error) {
	var categories []HelpCategoryWithCount
	err := DB.Model(&HelpCategory{}).
		Select("help_categories.*, COUNT(help_documents.id) as document_count").
		Joins("LEFT JOIN help_documents ON help_documents.category_id = help_categories.id").
		Group("help_categories.id").
		Order("sort_order ASC, id ASC").
		Find(&categories).Error
	return categories, err
}

// GetHelpCategoryById 根据ID获取分类
func GetHelpCategoryById(id int) (*HelpCategory, error) {
	var category HelpCategory
	err := DB.First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// CreateHelpCategory 创建分类
func CreateHelpCategory(category *HelpCategory) error {
	if category.Name == "" {
		return errors.New("分类名称不能为空")
	}
	return DB.Create(category).Error
}

// UpdateHelpCategory 更新分类
func UpdateHelpCategory(category *HelpCategory) error {
	if category.Name == "" {
		return errors.New("分类名称不能为空")
	}
	return DB.Save(category).Error
}

// DeleteHelpCategory 删除分类
func DeleteHelpCategory(id int) error {
	// 检查分类下是否有文档
	var count int64
	DB.Model(&HelpDocument{}).Where("category_id = ?", id).Count(&count)
	if count > 0 {
		return errors.New("该分类下存在文档，无法删除")
	}
	return DB.Delete(&HelpCategory{}, id).Error
}

// GetHelpCategoryDocumentCount 获取分类下的文档数量
func GetHelpCategoryDocumentCount(categoryId int) (int64, error) {
	var count int64
	err := DB.Model(&HelpDocument{}).Where("category_id = ?", categoryId).Count(&count).Error
	return count, err
}
