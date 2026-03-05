package handlers

import (
	"crypto/sha256"
	"encoding/hex"
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

// MicroHandler 微语处理器
type MicroHandler struct{}

// NewMicroHandler 创建微语处理器
func NewMicroHandler() *MicroHandler {
	return &MicroHandler{}
}

// CreateMicroRequest 创建微语请求
type CreateMicroRequest struct {
	Content string   `json:"content" binding:"required,min=1,max=2000"`
	Images  []string `json:"images"`
	Tags    []string `json:"tags"`
}

// CreateMicro 创建微语
func (h *MicroHandler) CreateMicro(c *gin.Context) {
	userID, isLoggedIn := middleware.GetUserID(c)
	if !isLoggedIn || userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请登录后再发布"})
		return
	}

	var req CreateMicroRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 将图片数组转为 JSON 字符串
	var imagesJSON string
	if len(req.Images) > 0 {
		imagesBytes, _ := json.Marshal(req.Images)
		imagesJSON = string(imagesBytes)
	}

	// 将标签数组转为 JSON 字符串
	var tagsJSON string
	if len(req.Tags) > 0 {
		tagsBytes, _ := json.Marshal(req.Tags)
		tagsJSON = string(tagsBytes)
	}

	micro := models.Micro{
		UserID:  userID,
		Content: req.Content,
		Images:  imagesJSON,
		Tags:    tagsJSON,
	}

	if err := database.DB.Create(&micro).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发布失败"})
		return
	}

	// 加载用户信息
	database.DB.Preload("User").First(&micro, micro.ID)

	c.JSON(http.StatusCreated, micro)
}

// GetMicros 获取微语列表
func (h *MicroHandler) GetMicros(c *gin.Context) {
	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}

	// 标签过滤
	tag := c.Query("tag")

	var micros []models.Micro
	var total int64

	query := database.DB.Model(&models.Micro{}).Where("deleted_at IS NULL")
	
	// 如果有标签过滤
	if tag != "" {
		query = query.Where("tags LIKE ?", "%\""+tag+"\"%")
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

	// 获取当前用户的点赞状态
	userID, isLoggedIn := middleware.GetUserID(c)
	userLikes := make(map[uint]bool)
	if isLoggedIn && userID > 0 {
		var likes []models.MicroLike
		microIDs := make([]uint, len(micros))
		for i, m := range micros {
			microIDs[i] = m.ID
		}
		database.DB.Where("micro_id IN ? AND user_id = ?", microIDs, userID).Find(&likes)
		for _, like := range likes {
			userLikes[like.MicroID] = true
		}
	}

	// 构建响应
	type MicroResponse struct {
		models.Micro
		IsLiked bool `json:"is_liked"`
	}

	responses := make([]MicroResponse, len(micros))
	for i, m := range micros {
		responses[i] = MicroResponse{
			Micro:   m,
			IsLiked: userLikes[m.ID],
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

// GetMicro 获取单条微语
func (h *MicroHandler) GetMicro(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	var micro models.Micro
	if err := database.DB.Preload("User").First(&micro, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "微语不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取失败"})
		return
	}

	// 获取点赞状态
	userID, isLoggedIn := middleware.GetUserID(c)
	isLiked := false
	if isLoggedIn && userID > 0 {
		var like models.MicroLike
		if err := database.DB.Where("micro_id = ? AND user_id = ?", micro.ID, userID).First(&like).Error; err == nil {
			isLiked = true
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    micro,
		"is_liked": isLiked,
	})
}

// DeleteMicro 删除微语
func (h *MicroHandler) DeleteMicro(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	role, _ := c.Get("role")

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	var micro models.Micro
	if err := database.DB.First(&micro, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "微语不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取失败"})
		return
	}

	// 检查权限：本人或管理员可删除
	if micro.UserID != userID && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限删除"})
		return
	}

	// 软删除
	if err := database.DB.Delete(&micro).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// ToggleMicroLike 切换点赞状态
func (h *MicroHandler) ToggleMicroLike(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	// 检查微语是否存在
	var micro models.Micro
	if err := database.DB.First(&micro, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "微语不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取失败"})
		return
	}

	userID, isLoggedIn := middleware.GetUserID(c)

	var existingLike models.MicroLike
	var likeErr error

	if isLoggedIn && userID > 0 {
		// 登录用户：按 user_id 查找
		likeErr = database.DB.Where("micro_id = ? AND user_id = ?", id, userID).First(&existingLike).Error
	} else {
		// 访客：按 IP Hash 查找
		ipHash := getIPHash(c)
		likeErr = database.DB.Where("micro_id = ? AND ip_hash = ?", id, ipHash).First(&existingLike).Error
	}

	if likeErr == gorm.ErrRecordNotFound {
		// 创建点赞
		newLike := models.MicroLike{
			MicroID: uint(id),
		}
		if isLoggedIn && userID > 0 {
			newLike.UserID = &userID
		} else {
			newLike.IPHash = getIPHash(c)
		}

		if err := database.DB.Create(&newLike).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "点赞失败"})
			return
		}

		// 更新点赞计数
		database.DB.Model(&micro).Update("like_count", gorm.Expr("like_count + 1"))

		c.JSON(http.StatusOK, gin.H{
			"liked":      true,
			"like_count": micro.LikeCount + 1,
		})
	} else if likeErr == nil {
		// 取消点赞
		if err := database.DB.Delete(&existingLike).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "取消点赞失败"})
			return
		}

		// 更新点赞计数
		newCount := micro.LikeCount - 1
		if newCount < 0 {
			newCount = 0
		}
		database.DB.Model(&micro).Update("like_count", newCount)

		c.JSON(http.StatusOK, gin.H{
			"liked":      false,
			"like_count": newCount,
		})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
}

// GetMicroHeatmap 获取热力图数据
func (h *MicroHandler) GetMicroHeatmap(c *gin.Context) {
	yearStr := c.DefaultQuery("year", strconv.Itoa(time.Now().Year()))
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		year = time.Now().Year()
	}

	// 计算年份的起止日期
	startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(year, 12, 31, 23, 59, 59, 0, time.UTC)

	// 按日期分组统计
	type DayCount struct {
		Date  string `json:"date"`
		Count int    `json:"count"`
	}

	var results []DayCount

	database.DB.Model(&models.Micro{}).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("created_at >= ? AND created_at <= ? AND deleted_at IS NULL", startDate, endDate).
		Group("DATE(created_at)").
		Find(&results)

	// 转换为 map
	data := make(map[string]int)
	for _, r := range results {
		data[r.Date] = r.Count
	}

	c.JSON(http.StatusOK, gin.H{
		"year": year,
		"data": data,
	})
}

// GetMicroStats 获取统计数据
func (h *MicroHandler) GetMicroStats(c *gin.Context) {
	var totalMicros int64
	var totalLikes int64
	var totalComments int64
	var monthMicros int64

	// 总微语数
	database.DB.Model(&models.Micro{}).Count(&totalMicros)

	// 总点赞数
	database.DB.Model(&models.MicroLike{}).Count(&totalLikes)

	// 总评论数
	database.DB.Model(&models.MicroComment{}).Count(&totalComments)

	// 本月发布数
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	database.DB.Model(&models.Micro{}).Where("created_at >= ? AND deleted_at IS NULL", monthStart).Count(&monthMicros)

	c.JSON(http.StatusOK, gin.H{
		"total_micros":   totalMicros,
		"total_likes":    totalLikes,
		"total_comments": totalComments,
		"month_micros":   monthMicros,
	})
}

// GetMicroTags 获取热门标签
func (h *MicroHandler) GetMicroTags(c *gin.Context) {
	// 获取所有微语的标签
	var micros []models.Micro
	database.DB.Where("tags IS NOT NULL AND tags != '' AND deleted_at IS NULL").Select("tags").Find(&micros)

	// 统计标签使用次数
	tagCount := make(map[string]int)
	for _, micro := range micros {
		if micro.Tags == "" {
			continue
		}
		var tags []string
		if err := json.Unmarshal([]byte(micro.Tags), &tags); err == nil {
			for _, tag := range tags {
				tagCount[tag]++
			}
		}
	}

	// 转换为排序后的列表
	type TagItem struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}

	var tags []TagItem
	for name, count := range tagCount {
		tags = append(tags, TagItem{Name: name, Count: count})
	}

	// 按使用次数排序
	for i := 0; i < len(tags); i++ {
		for j := i + 1; j < len(tags); j++ {
			if tags[j].Count > tags[i].Count {
				tags[i], tags[j] = tags[j], tags[i]
			}
		}
	}

	// 只返回前20个
	if len(tags) > 20 {
		tags = tags[:20]
	}

	c.JSON(http.StatusOK, gin.H{
		"tags": tags,
	})
}

// getIPHash 获取 IP 的哈希值
func getIPHash(c *gin.Context) string {
	ip := c.ClientIP()
	hash := sha256.Sum256([]byte(ip + "micro-salt"))
	return hex.EncodeToString(hash[:])
}

// CreateMicroCommentRequest 创建微语评论请求
type CreateMicroCommentRequest struct {
	MicroID  uint   `json:"micro_id" binding:"required"`
	ParentID *uint  `json:"parent_id"`
	Content  string `json:"content" binding:"required,min=1,max=2000"`
}

// CreateMicroComment 创建微语评论
func (h *MicroHandler) CreateMicroComment(c *gin.Context) {
	userID, isLoggedIn := middleware.GetUserID(c)
	if !isLoggedIn || userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请登录后再评论"})
		return
	}

	var req CreateMicroCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查微语是否存在
	var micro models.Micro
	if err := database.DB.First(&micro, req.MicroID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "微语不存在"})
		return
	}

	comment := models.MicroComment{
		MicroID:  req.MicroID,
		UserID:   userID,
		ParentID: req.ParentID,
		Content:  req.Content,
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "评论失败"})
		return
	}

	// 更新评论计数
	database.DB.Model(&micro).Update("comment_count", gorm.Expr("comment_count + 1"))

	// 加载用户信息
	database.DB.Preload("User").First(&comment, comment.ID)

	c.JSON(http.StatusCreated, comment)
}

// GetMicroComments 获取微语评论列表
func (h *MicroHandler) GetMicroComments(c *gin.Context) {
	microIDStr := c.Query("micro_id")
	if microIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "micro_id is required"})
		return
	}

	microID, err := strconv.ParseUint(microIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid micro_id"})
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
	var comments []models.MicroComment
	var total int64

	query := database.DB.Model(&models.MicroComment{}).
		Where("micro_id = ? AND parent_id IS NULL", microID)

	query.Count(&total)

	if err := query.
		Preload("User").
		Preload("Replies.User").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取评论失败"})
		return
	}

	// 手动加载每个评论的用户信息
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

// DeleteMicroComment 删除微语评论
func (h *MicroHandler) DeleteMicroComment(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	role, _ := c.Get("role")

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment id"})
		return
	}

	var comment models.MicroComment
	if err := database.DB.First(&comment, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	// 更新评论计数
	database.DB.Model(&models.Micro{}).Where("id = ?", comment.MicroID).
		Update("comment_count", gorm.Expr("comment_count - 1"))

	c.JSON(http.StatusOK, gin.H{"message": "评论已删除"})
}