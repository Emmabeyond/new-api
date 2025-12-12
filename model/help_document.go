package model

import (
	"errors"
)

const (
	HelpDocumentStatusEnabled  = 1
	HelpDocumentStatusDisabled = 2
)

// HelpDocument 帮助文档
type HelpDocument struct {
	Id         int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Title      string `json:"title" gorm:"type:varchar(255);not null"`
	Content    string `json:"content" gorm:"type:text"`
	CategoryId int    `json:"category_id" gorm:"index"`
	SortOrder  int    `json:"sort_order" gorm:"default:0"`
	Status     int    `json:"status" gorm:"default:1"` // 1=启用, 2=禁用
	CreatedAt  int64  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  int64  `json:"updated_at" gorm:"autoUpdateTime"`
}

// HelpDocumentBrief 文档简要信息（不含内容）
type HelpDocumentBrief struct {
	Id         int    `json:"id"`
	Title      string `json:"title"`
	CategoryId int    `json:"category_id"`
	SortOrder  int    `json:"sort_order"`
}

// HelpCategoryWithDocuments 带文档列表的分类
type HelpCategoryWithDocuments struct {
	Id        int                 `json:"id"`
	Name      string              `json:"name"`
	SortOrder int                 `json:"sort_order"`
	Documents []HelpDocumentBrief `json:"documents"`
}

func (HelpDocument) TableName() string {
	return "help_documents"
}

// GetAllHelpDocuments 获取所有启用的文档（按分类和排序权重）
func GetAllHelpDocuments() ([]HelpDocument, error) {
	var documents []HelpDocument
	err := DB.Where("status = ?", HelpDocumentStatusEnabled).
		Order("category_id ASC, sort_order ASC, id ASC").
		Find(&documents).Error
	return documents, err
}

// GetHelpDocumentsGroupedByCategory 获取按分类分组的文档列表
func GetHelpDocumentsGroupedByCategory() ([]HelpCategoryWithDocuments, error) {
	// 获取所有启用的分类
	categories, err := GetAllHelpCategories()
	if err != nil {
		return nil, err
	}

	// 获取所有启用的文档
	var documents []HelpDocumentBrief
	err = DB.Model(&HelpDocument{}).
		Select("id, title, category_id, sort_order").
		Where("status = ?", HelpDocumentStatusEnabled).
		Order("sort_order ASC, id ASC").
		Find(&documents).Error
	if err != nil {
		return nil, err
	}

	// 按分类分组
	docMap := make(map[int][]HelpDocumentBrief)
	for _, doc := range documents {
		docMap[doc.CategoryId] = append(docMap[doc.CategoryId], doc)
	}

	// 组装结果
	result := make([]HelpCategoryWithDocuments, 0, len(categories))
	for _, cat := range categories {
		docs := docMap[cat.Id]
		if docs == nil {
			docs = []HelpDocumentBrief{}
		}
		result = append(result, HelpCategoryWithDocuments{
			Id:        cat.Id,
			Name:      cat.Name,
			SortOrder: cat.SortOrder,
			Documents: docs,
		})
	}

	return result, nil
}

// GetAllHelpDocumentsAdmin 管理员获取所有文档（包含禁用的）
func GetAllHelpDocumentsAdmin() ([]HelpDocument, error) {
	var documents []HelpDocument
	err := DB.Order("category_id ASC, sort_order ASC, id ASC").
		Find(&documents).Error
	return documents, err
}

// GetHelpDocumentById 根据ID获取文档
func GetHelpDocumentById(id int) (*HelpDocument, error) {
	var document HelpDocument
	err := DB.First(&document, id).Error
	if err != nil {
		return nil, err
	}
	return &document, nil
}

// GetHelpDocumentByIdPublic 获取启用的文档（公开接口）
func GetHelpDocumentByIdPublic(id int) (*HelpDocument, error) {
	var document HelpDocument
	err := DB.Where("id = ? AND status = ?", id, HelpDocumentStatusEnabled).
		First(&document).Error
	if err != nil {
		return nil, err
	}
	return &document, nil
}

// CreateHelpDocument 创建文档
func CreateHelpDocument(document *HelpDocument) error {
	if document.Title == "" {
		return errors.New("文档标题不能为空")
	}
	if document.Content == "" {
		return errors.New("文档内容不能为空")
	}
	return DB.Create(document).Error
}

// UpdateHelpDocument 更新文档
func UpdateHelpDocument(document *HelpDocument) error {
	if document.Title == "" {
		return errors.New("文档标题不能为空")
	}
	if document.Content == "" {
		return errors.New("文档内容不能为空")
	}
	return DB.Save(document).Error
}

// DeleteHelpDocument 删除文档
func DeleteHelpDocument(id int) error {
	return DB.Delete(&HelpDocument{}, id).Error
}

// SearchHelpDocuments 搜索文档（按标题）
func SearchHelpDocuments(query string) ([]HelpDocumentBrief, error) {
	var documents []HelpDocumentBrief
	err := DB.Model(&HelpDocument{}).
		Select("id, title, category_id, sort_order").
		Where("status = ? AND title LIKE ?", HelpDocumentStatusEnabled, "%"+query+"%").
		Order("sort_order ASC, id ASC").
		Find(&documents).Error
	return documents, err
}
