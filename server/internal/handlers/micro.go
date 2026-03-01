package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/novablog/server/internal/database"
	"github.com/novablog/server/internal/middleware"
	"github.com/novablog/server/internal/models"
	"gorm.io/gorm"
)

func mustMarshal(v interface{}) []byte {
	data, _ := json.Marshal(v)
	return data
}

type MicroHandler struct{}

func NewMicroHandler() *MicroHandler {
	return &MicroHandler{}
}

type CreateMicroRequest struct {
	Content  string   `json:"content" binding:"required,max=2000"`
	Images   []string `json:"images"`
	Tags     []string `json:"tags"`
	IsPublic bool     `json:"is_public"`
}

type UpdateMicroRequest struct {
	Content  string   `json:"content" binding:"required,max=2000"`
	Images   []string `json:"images"`
	Tags     []string `json:"tags"`
	IsPublic bool     `json:"is_public"`
}

type MicroResponse struct {
	models.MicroPost
	LikeCount int  `json:"like_count"`
	IsLiked   bool `json:"is_liked"`
}

func (h *MicroHandler) CreateMicro(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok || userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请登录后再发布"})
		return
	}

	var req CreateMicroRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	imagesJSON := "[]"
	if len(req.Images) > 0 {
		imagesJSON = string(mustMarshal(req.Images))
	}

	tagsJSON := "[]"
	if len(req.Tags) > 0 {
		tagsJSON = string(mustMarshal(req.Tags))
	}

	micro := models.MicroPost{
		UserID:   userID,
		Content:  req.Content,
		Images:   imagesJSON,
		Tags:     tagsJSON,
		IsPublic: req.IsPublic,
	}

	if err := database.DB.Create(&micro).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发布失败"})
		return
	}

	database.DB.Preload("User").First(&micro, micro.ID)

	c.JSON(http.StatusCreated, MicroResponse{
		MicroPost: micro,
		LikeCount: 0,
		IsLiked:   false,
	})
}

func (h *MicroHandler) GetMicros(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	userIDQuery := c.Query("user_id")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}

	var micros []models.MicroPost
	var total int64

	query := database.DB.Model(&models.MicroPost{}).Where("is_public = ?", true)

	if userIDQuery != "" {
		query = query.Where("user_id = ?", userIDQuery)
	}

	query.Count(&total)

	if err := query.
		Preload("User").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&micros).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取失败"})
		return
	}

	currentUserID, _ := middleware.GetUserID(c)

	responses := make([]MicroResponse, len(micros))
	for i, micro := range micros {
		var likeCount int64
		database.DB.Model(&models.MicroPostLike{}).Where("micro_post_id = ?", micro.ID).Count(&likeCount)

		isLiked := false
		if currentUserID > 0 {
			var count int64
			database.DB.Model(&models.MicroPostLike{}).
				Where("micro_post_id = ? AND user_id = ?", micro.ID, currentUserID).
				Count(&count)
			isLiked = count > 0
		}

		responses[i] = MicroResponse{
			MicroPost: micro,
			LikeCount: int(likeCount),
			IsLiked:   isLiked,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": responses,
		"pagination": gin.H{
			"page":       page,
			"page_size":  pageSize,
			"total":      total,
			"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

func (h *MicroHandler) GetMicro(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var micro models.MicroPost
	if err := database.DB.Preload("User").First(&micro, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	var likeCount int64
	database.DB.Model(&models.MicroPostLike{}).Where("micro_post_id = ?", micro.ID).Count(&likeCount)

	currentUserID, _ := middleware.GetUserID(c)
	isLiked := false
	if currentUserID > 0 {
		var count int64
		database.DB.Model(&models.MicroPostLike{}).
			Where("micro_post_id = ? AND user_id = ?", micro.ID, currentUserID).
			Count(&count)
		isLiked = count > 0
	}

	c.JSON(http.StatusOK, MicroResponse{
		MicroPost: micro,
		LikeCount: int(likeCount),
		IsLiked:   isLiked,
	})
}

func (h *MicroHandler) UpdateMicro(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var micro models.MicroPost
	if err := database.DB.First(&micro, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	if micro.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权修改"})
		return
	}

	var req UpdateMicroRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	imagesJSON := "[]"
	if len(req.Images) > 0 {
		imagesJSON = string(mustMarshal(req.Images))
	}

	tagsJSON := "[]"
	if len(req.Tags) > 0 {
		tagsJSON = string(mustMarshal(req.Tags))
	}

	updates := map[string]interface{}{
		"content":   req.Content,
		"images":    imagesJSON,
		"tags":      tagsJSON,
		"is_public": req.IsPublic,
	}

	if err := database.DB.Model(&micro).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	database.DB.Preload("User").First(&micro, micro.ID)

	c.JSON(http.StatusOK, micro)
}

func (h *MicroHandler) DeleteMicro(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	role, _ := c.Get("role")

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var micro models.MicroPost
	if err := database.DB.First(&micro, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	if micro.UserID != userID && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权删除"})
		return
	}

	if err := database.DB.Delete(&micro).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func (h *MicroHandler) ToggleLike(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok || userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请登录后再点赞"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var micro models.MicroPost
	if err := database.DB.First(&micro, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	var existingLike models.MicroPostLike
	result := database.DB.Where("micro_post_id = ? AND user_id = ?", id, userID).First(&existingLike)

	if result.Error == nil {
		database.DB.Delete(&existingLike)
		c.JSON(http.StatusOK, gin.H{"liked": false, "message": "取消点赞"})
		return
	}

	like := models.MicroPostLike{
		MicroPostID: uint(id),
		UserID:      userID,
	}

	if err := database.DB.Create(&like).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "点赞失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"liked": true, "message": "点赞成功"})
}

func (h *MicroHandler) GetHeatmap(c *gin.Context) {
	userIDQuery := c.Query("user_id")
	yearStr := c.DefaultQuery("year", strconv.Itoa(time.Now().Year()))
	year, _ := strconv.Atoi(yearStr)

	query := database.DB.Model(&models.MicroPost{}).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("is_public = ?", true).
		Group("DATE(created_at)")

	if userIDQuery != "" {
		query = query.Where("user_id = ?", userIDQuery)
	}

	if year > 0 {
		startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)
		query = query.Where("created_at >= ? AND created_at < ?", startDate, endDate)
	}

	var results []struct {
		Date  string `json:"date"`
		Count int    `json:"count"`
	}

	if err := query.Scan(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取失败"})
		return
	}

	c.JSON(http.StatusOK, results)
}

func (h *MicroHandler) GetStats(c *gin.Context) {
	userIDQuery := c.Query("user_id")

	var totalMicros int64
	var totalUsers int64

	query := database.DB.Model(&models.MicroPost{}).Where("is_public = ?", true)
	if userIDQuery != "" {
		query = query.Where("user_id = ?", userIDQuery)
	}
	query.Count(&totalMicros)

	database.DB.Model(&models.User{}).Count(&totalUsers)

	var topUsers []struct {
		UserID    uint   `json:"user_id"`
		Username  string `json:"username"`
		Nickname  string `json:"nickname"`
		Avatar    string `json:"avatar"`
		PostCount int    `json:"post_count"`
	}

	database.DB.Model(&models.MicroPost{}).
		Select("user_id, COUNT(*) as post_count").
		Where("is_public = ?", true).
		Group("user_id").
		Order("post_count DESC").
		Limit(10).
		Scan(&topUsers)

	for i := range topUsers {
		var user models.User
		if err := database.DB.First(&user, topUsers[i].UserID).Error; err == nil {
			topUsers[i].Username = user.Username
			topUsers[i].Nickname = user.Nickname
			topUsers[i].Avatar = user.Avatar
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"total_micros": totalMicros,
		"total_users":  totalUsers,
		"top_users":    topUsers,
	})
}
