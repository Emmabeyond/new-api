package controller

import (
	"net/http"
	"strconv"

	"github.com/QuantumNous/new-api/model"
	"github.com/gin-gonic/gin"
)

// ==================== 公开接口 ====================

// GetHelpCategories 获取所有启用的帮助分类（带文档列表）
func GetHelpCategories(c *gin.Context) {
	data, err := model.GetCachedHelpData()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "获取帮助分类失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    data,
	})
}

// GetHelpDocument 获取单个帮助文档详情
func GetHelpDocument(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的文档ID",
		})
		return
	}

	doc, err := model.GetCachedHelpDocument(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "文档不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    doc,
	})
}

// SearchHelpDocuments 搜索帮助文档
func SearchHelpDocuments(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "搜索关键词不能为空",
		})
		return
	}

	docs, err := model.SearchHelpDocuments(query)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "搜索失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    docs,
	})
}

// ==================== 管理员接口 ====================

// AdminGetHelpCategories 管理员获取所有帮助分类
func AdminGetHelpCategories(c *gin.Context) {
	categories, err := model.GetAllHelpCategoriesAdmin()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "获取分类失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    categories,
	})
}

// AdminCreateHelpCategory 管理员创建帮助分类
func AdminCreateHelpCategory(c *gin.Context) {
	var category model.HelpCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的请求参数",
		})
		return
	}

	if err := model.CreateHelpCategory(&category); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	model.InvalidateHelpCache()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "分类创建成功",
		"data":    category,
	})
}

// AdminUpdateHelpCategory 管理员更新帮助分类
func AdminUpdateHelpCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的分类ID",
		})
		return
	}

	var category model.HelpCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的请求参数",
		})
		return
	}

	category.Id = id
	if err := model.UpdateHelpCategory(&category); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	model.InvalidateHelpCache()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "分类更新成功",
	})
}

// AdminDeleteHelpCategory 管理员删除帮助分类
func AdminDeleteHelpCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的分类ID",
		})
		return
	}

	if err := model.DeleteHelpCategory(id); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	model.InvalidateHelpCache()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "分类删除成功",
	})
}

// AdminGetHelpDocuments 管理员获取所有帮助文档
func AdminGetHelpDocuments(c *gin.Context) {
	docs, err := model.GetAllHelpDocumentsAdmin()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "获取文档失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    docs,
	})
}

// AdminCreateHelpDocument 管理员创建帮助文档
func AdminCreateHelpDocument(c *gin.Context) {
	var doc model.HelpDocument
	if err := c.ShouldBindJSON(&doc); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的请求参数",
		})
		return
	}

	if err := model.CreateHelpDocument(&doc); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	model.InvalidateHelpCache()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "文档创建成功",
		"data":    doc,
	})
}

// AdminUpdateHelpDocument 管理员更新帮助文档
func AdminUpdateHelpDocument(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的文档ID",
		})
		return
	}

	var doc model.HelpDocument
	if err := c.ShouldBindJSON(&doc); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的请求参数",
		})
		return
	}

	doc.Id = id
	if err := model.UpdateHelpDocument(&doc); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	model.InvalidateHelpDocumentCache(id)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "文档更新成功",
	})
}

// AdminDeleteHelpDocument 管理员删除帮助文档
func AdminDeleteHelpDocument(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的文档ID",
		})
		return
	}

	if err := model.DeleteHelpDocument(id); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	model.InvalidateHelpDocumentCache(id)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "文档删除成功",
	})
}
