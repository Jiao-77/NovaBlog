package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/novablog/server/internal/database"
	"github.com/novablog/server/internal/middleware"
	"github.com/novablog/server/internal/models"
	"gorm.io/gorm"
)

// LikeHandler 点赞处理器
type LikeHandler struct{}

// NewLikeHandler 创建点赞处理器
func NewLikeHandler() *LikeHandler {
	return &LikeHandler{}
}

// LikeRequest 点赞请求
type LikeRequest struct {
	PostID string `json:"post_id" binding:"required"`
}

// LikeResponse 点赞响应
type LikeResponse struct {
	Liked    bool `json:"liked"`
	LikeCount int  `json:"like_count"`
}

// ToggleLike 切换点赞状态（点赞/取消点赞）
func (h *LikeHandler) ToggleLike(c *gin.Context) {
	var req LikeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取用户 ID（可选，支持未登录用户）
	userID, isLoggedIn := middleware.GetUserID(c)

	// 获取访客 IP Hash（用于未登录用户的防刷）
	ipHash := ""
	if !isLoggedIn {
		ip := c.ClientIP()
		hash := sha256.Sum256([]byte(ip + "novablog-salt")) // 加盐防止反向推导
		ipHash = hex.EncodeToString(hash[:])[:64]
	}

	// 检查是否已点赞
	var existingLike models.Like
	var err error

	if isLoggedIn {
		err = database.DB.Where("post_id = ? AND user_id = ?", req.PostID, userID).First(&existingLike).Error
	} else {
		err = database.DB.Where("post_id = ? AND ip_hash = ?", req.PostID, ipHash).First(&existingLike).Error
	}

	liked := false
	if err == gorm.ErrRecordNotFound {
		// 未点赞，创建点赞记录
		like := models.Like{
			PostID: req.PostID,
			IPHash: ipHash,
		}
		if isLoggedIn {
			like.UserID = &userID
		}

		if err := database.DB.Create(&like).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to like"})
			return
		}
		liked = true

		// 更新点赞计数
		h.updateLikeCount(req.PostID, 1)
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	} else {
		// 已点赞，取消点赞
		if err := database.DB.Delete(&existingLike).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unlike"})
			return
		}
		liked = false

		// 更新点赞计数
		h.updateLikeCount(req.PostID, -1)
	}

	// 获取当前点赞数
	likeCount := h.getLikeCount(req.PostID)

	c.JSON(http.StatusOK, LikeResponse{
		Liked:     liked,
		LikeCount: likeCount,
	})
}

// GetLikeStatus 获取点赞状态
func (h *LikeHandler) GetLikeStatus(c *gin.Context) {
	postID := c.Query("post_id")
	if postID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "post_id is required"})
		return
	}

	// 获取用户 ID（可选）
	userID, isLoggedIn := middleware.GetUserID(c)

	// 获取访客 IP Hash
	ipHash := ""
	if !isLoggedIn {
		ip := c.ClientIP()
		hash := sha256.Sum256([]byte(ip + "novablog-salt"))
		ipHash = hex.EncodeToString(hash[:])[:64]
	}

	// 检查是否已点赞
	var existingLike models.Like
	var err error

	if isLoggedIn {
		err = database.DB.Where("post_id = ? AND user_id = ?", postID, userID).First(&existingLike).Error
	} else {
		err = database.DB.Where("post_id = ? AND ip_hash = ?", postID, ipHash).First(&existingLike).Error
	}

	liked := err == nil

	// 获取点赞数
	likeCount := h.getLikeCount(postID)

	c.JSON(http.StatusOK, LikeResponse{
		Liked:     liked,
		LikeCount: likeCount,
	})
}

// updateLikeCount 更新点赞计数
func (h *LikeHandler) updateLikeCount(postID string, delta int) {
	var likeCount models.LikeCount
	result := database.DB.FirstOrCreate(&likeCount, models.LikeCount{PostID: postID})
	if result.Error != nil {
		return
	}

	likeCount.Count += delta
	if likeCount.Count < 0 {
		likeCount.Count = 0
	}

	database.DB.Save(&likeCount)
}

// getLikeCount 获取点赞数
func (h *LikeHandler) getLikeCount(postID string) int {
	var likeCount models.LikeCount
	if err := database.DB.Where("post_id = ?", postID).First(&likeCount).Error; err != nil {
		return 0
	}
	return likeCount.Count
}