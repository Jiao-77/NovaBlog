package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/novablog/server/internal/database"
	"github.com/novablog/server/internal/middleware"
	"github.com/novablog/server/internal/models"
	"gorm.io/gorm"
)

// CommentHandler 评论处理器
type CommentHandler struct{}

// NewCommentHandler 创建评论处理器
func NewCommentHandler() *CommentHandler {
	return &CommentHandler{}
}

// CreateCommentRequest 创建评论请求
type CreateCommentRequest struct {
	PostID   string `json:"post_id" binding:"required"`
	ParentID *uint  `json:"parent_id"`
	Content  string `json:"content" binding:"required,min=1,max=2000"`
}

// CreateComment 创建评论
func (h *CommentHandler) CreateComment(c *gin.Context) {
	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 尝试从上下文获取用户 ID（如果有认证中间件）
	userID, isLoggedIn := middleware.GetUserID(c)
	
	// 如果没有登录，返回错误（评论需要登录）
	if !isLoggedIn || userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请登录后再评论"})
		return
	}

	comment := models.Comment{
		PostID:   req.PostID,
		UserID:   userID,
		ParentID: req.ParentID,
		Content:  req.Content,
		Status:   "approved", // 默认直接通过，可改为 pending 需审核
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create comment"})
		return
	}

	// 加载用户信息
	database.DB.Preload("User").First(&comment, comment.ID)

	c.JSON(http.StatusCreated, comment)
}

// GetComments 获取文章评论列表
func (h *CommentHandler) GetComments(c *gin.Context) {
	postID := c.Query("post_id")
	if postID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "post_id is required"})
		return
	}

	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 获取顶级评论（非回复）
	var comments []models.Comment
	var total int64

	query := database.DB.Model(&models.Comment{}).
		Where("post_id = ? AND status = ? AND parent_id IS NULL", postID, "approved")

	query.Count(&total)

	// 使用 Joins 显式加载用户信息，确保数据完整
	if err := query.
		Joins("LEFT JOIN users ON users.id = comments.user_id").
		Preload("Replies.User").
		Order("comments.created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get comments"})
		return
	}

	// 手动加载每个评论的用户信息（确保 Preload 正确工作）
	for i := range comments {
		if comments[i].UserID > 0 {
			var user models.User
			if err := database.DB.First(&user, comments[i].UserID).Error; err == nil {
				comments[i].User = user
			}
		}
		// 加载回复的用户信息
		for j := range comments[i].Replies {
			if comments[i].Replies[j].UserID > 0 {
				var replyUser models.User
				if err := database.DB.First(&replyUser, comments[i].Replies[j].UserID).Error; err == nil {
					comments[i].Replies[j].User = replyUser
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": comments,
		"pagination": gin.H{
			"page":       page,
			"page_size":  pageSize,
			"total":      total,
			"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// DeleteComment 删除评论（仅限本人或管理员）
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	role, _ := c.Get("role")

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment id"})
		return
	}

	var comment models.Comment
	if err := database.DB.First(&comment, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	// 检查权限：本人或管理员可删除
	if comment.UserID != userID && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
		return
	}

	// 软删除
	if err := database.DB.Delete(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "comment deleted"})
}