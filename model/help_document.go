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


// InitDefaultAIToolsHelpDocuments 初始化 AI 工具相关帮助文档
// 创建 AI 工具指南分类和相关文档
func InitDefaultAIToolsHelpDocuments() error {
	// 检查是否已有 AI 工具指南分类
	var existingCategory HelpCategory
	err := DB.Where("name = ?", "AI 工具指南").First(&existingCategory).Error
	if err == nil {
		return nil // 分类已存在，不需要初始化
	}

	// 创建 AI 工具指南分类
	aiCategory := HelpCategory{
		Name:      "AI 工具指南",
		SortOrder: 100, // 较高的排序值，显示在后面
		Status:    HelpCategoryStatusEnabled,
	}
	if err := DB.Create(&aiCategory).Error; err != nil {
		return err
	}

	// 创建帮助文档
	documents := []HelpDocument{
		{
			Title:      "API 调用最佳实践",
			CategoryId: aiCategory.Id,
			SortOrder:  1,
			Status:     HelpDocumentStatusEnabled,
			Content: `# API 调用最佳实践

## 1. 选择合适的模型

根据任务类型选择最适合的模型：

| 任务类型 | 推荐模型 | 说明 |
|---------|---------|------|
| 日常对话 | GPT-5.2 Instant | 响应快速，成本低 |
| 代码编写 | Claude Sonnet 4.5 | Agent 能力强，编码效率高 |
| 深度推理 | GPT-5.2 Thinking | 400K 上下文，推理能力强 |
| 复杂任务 | Claude Opus 4.5 | 综合能力最强 |

## 2. 优化 Prompt 设计

- **明确任务目标**：清晰描述期望的输出格式和内容
- **提供上下文**：给出必要的背景信息
- **分步骤执行**：复杂任务拆分为多个步骤
- **使用示例**：提供输入输出示例帮助模型理解

## 3. 控制成本

- 使用流式响应减少等待时间
- 合理设置 max_tokens 限制
- 缓存常用查询结果
- 使用更经济的模型处理简单任务

## 4. 错误处理

- 实现重试机制（指数退避）
- 处理速率限制（429 错误）
- 记录请求日志便于排查问题
- 设置合理的超时时间`,
		},
		{
			Title:      "模型选择建议",
			CategoryId: aiCategory.Id,
			SortOrder:  2,
			Status:     HelpDocumentStatusEnabled,
			Content: `# 模型选择建议

## GPT-5.2 系列

### GPT-5.2 Instant
- **特点**：响应速度快，成本低
- **适用场景**：日常对话、简单问答、内容生成
- **上下文窗口**：128K tokens

### GPT-5.2 Thinking
- **特点**：深度推理能力强，支持复杂逻辑
- **适用场景**：数学推理、代码调试、复杂分析
- **上下文窗口**：400K tokens

### GPT-5.2 Pro
- **特点**：研究级能力，最强综合性能
- **适用场景**：学术研究、专业分析、高精度任务
- **上下文窗口**：400K tokens

## Claude 4.5 系列

### Claude Sonnet 4.5
- **特点**：编码能力最强，Agent 能力出色
- **适用场景**：代码编写、自动化任务、工具调用
- **定价**：$3/$15 per M tokens
- **Token 效率**：比上代提升 76%

### Claude Opus 4.5
- **特点**：复杂推理能力最强
- **适用场景**：复杂任务、长文档处理、深度分析
- **定价**：$5/$25 per M tokens（降价 2/3）

## 选择建议

| 需求 | 推荐模型 |
|------|---------|
| 快速响应 | GPT-5.2 Instant |
| 写代码 | Claude Sonnet 4.5 |
| 推理任务 | GPT-5.2 Thinking |
| 复杂分析 | Claude Opus 4.5 |
| 成本敏感 | GPT-5.2 Instant / Claude Sonnet 4.5 |`,
		},
		{
			Title:      "常见错误排查",
			CategoryId: aiCategory.Id,
			SortOrder:  3,
			Status:     HelpDocumentStatusEnabled,
			Content: `# 常见错误排查

## 1. 认证错误 (401)

**问题**：API Key 无效或过期

**解决方案**：
- 检查 API Key 是否正确复制
- 确认 API Key 未被禁用
- 检查 Key 的权限范围

## 2. 速率限制 (429)

**问题**：请求频率超过限制

**解决方案**：
- 实现指数退避重试
- 使用请求队列控制并发
- 升级账户等级获取更高限额

## 3. 上下文超限 (400)

**问题**：输入 token 数超过模型限制

**解决方案**：
- 精简 Prompt 内容
- 使用支持更大上下文的模型
- 分批处理长文档

## 4. 模型幻觉

**问题**：模型生成不准确或虚假信息

**解决方案**：
- 要求模型引用来源
- 使用 RAG 技术提供真实数据
- 分步验证关键信息
- 设置较低的 temperature

## 5. 响应截断

**问题**：输出内容不完整

**解决方案**：
- 增加 max_tokens 设置
- 检查 finish_reason 字段
- 使用流式响应获取完整内容

## 6. 超时错误

**问题**：请求超时未响应

**解决方案**：
- 增加客户端超时时间
- 使用流式响应减少等待
- 检查网络连接状态`,
		},
	}

	for _, doc := range documents {
		if err := DB.Create(&doc).Error; err != nil {
			return err
		}
	}

	return nil
}


// InitDefaultFAQ 初始化 AI 工具相关 FAQ
// 在系统首次部署时创建示例 FAQ
func InitDefaultFAQ() error {
	// 检查是否已有 FAQ 配置
	var option Option
	err := DB.Where("`key` = ?", "console_setting.faq").First(&option).Error
	if err == nil && option.Value != "" && option.Value != "[]" {
		return nil // 已有 FAQ 配置，不需要初始化
	}

	// 创建默认 FAQ 数据
	defaultFAQ := `[
		{"question": "GPT-5.2 有哪些版本？", "answer": "GPT-5.2 有三个版本：Instant（快速日常对话）、Thinking（深度推理，400K上下文）、Pro（研究级任务）"},
		{"question": "Claude Sonnet 4.5 vs Opus 4.5 怎么选？", "answer": "Sonnet 4.5 编码和 Agent 能力最强，适合代码编写；Opus 4.5 复杂推理更好，适合深度分析任务"},
		{"question": "GPT-5.2 上下文窗口多大？", "answer": "GPT-5.2 Thinking 和 Pro 版本支持 400K tokens 上下文，可同时处理数百文档"},
		{"question": "Claude Opus 4.5 定价是多少？", "answer": "$5/$25 per M tokens（输入/输出），比之前版本降价约 2/3"},
		{"question": "推理任务选什么模型？", "answer": "推荐 GPT-5.2 Thinking 或 Claude Opus 4.5，两者都有出色的深度推理能力"},
		{"question": "如何降低 API 成本？", "answer": "使用 Instant/Haiku 版本处理简单任务、缓存常用查询结果、精简 Prompt 内容"},
		{"question": "模型幻觉怎么处理？", "answer": "要求模型引用来源、分步验证关键信息、使用 RAG 技术提供真实数据"},
		{"question": "Token 效率如何优化？", "answer": "Claude 4.5 系列 Token 效率比上代提升 76%，同样内容消耗更少 Token"}
	]`

	// 保存到数据库
	return UpdateOption("console_setting.faq", defaultFAQ)
}
