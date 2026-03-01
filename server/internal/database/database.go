package database

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/novablog/server/internal/config"
	"github.com/novablog/server/internal/models"
	"github.com/novablog/server/internal/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Initialize 初始化数据库连接
func Initialize(cfg *config.Config) error {
	var err error

	// 确保数据目录存在
	dbDir := filepath.Dir(cfg.Database.Path)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return fmt.Errorf("failed to create database directory: %w", err)
	}

	// 连接 SQLite 数据库
	DB, err = gorm.Open(sqlite.Open(cfg.Database.Path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	// 自动迁移数据库表
	if err := autoMigrate(); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	// 创建默认管理员用户
	if err := createDefaultAdmin(); err != nil {
		return fmt.Errorf("failed to create default admin: %w", err)
	}

	return nil
}

// autoMigrate 自动迁移数据库表结构
func autoMigrate() error {
	return DB.AutoMigrate(
		&models.User{},
		&models.Comment{},
		&models.Like{},
		&models.LikeCount{},
		&models.PostMeta{},
	)
}

// Close 关闭数据库连接
func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// createDefaultAdmin 创建默认管理员用户
func createDefaultAdmin() error {
	// 检查是否已存在 admin 用户
	var count int64
	DB.Model(&models.User{}).Where("username = ?", "admin").Count(&count)
	if count > 0 {
		return nil // 已存在，跳过
	}

	// 创建 admin 用户
	hashedPassword, err := utils.HashPassword("admin")
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	admin := models.User{
		Username: "admin",
		Email:    "admin@novablog.local",
		Password: hashedPassword,
		Role:     "admin",
		Nickname: "Administrator",
	}

	if err := DB.Create(&admin).Error; err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}

	fmt.Println("✅ Default admin user created: admin / admin")
	return nil
}
