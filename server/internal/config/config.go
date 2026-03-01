package config

import (
	"os"
	"strconv"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	CORS     CORSConfig
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string
	Mode string // debug, release, test
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Path string
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret     string
	ExpireTime int // 过期时间（小时）
}

// CORSConfig CORS 配置
type CORSConfig struct {
	AllowOrigins []string
}

// Load 从环境变量加载配置
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Mode: getEnv("GIN_MODE", "release"),
		},
		Database: DatabaseConfig{
			Path: getEnv("DB_PATH", "./data/novablog.db"),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "novablog-secret-key-change-in-production"),
			ExpireTime: getEnvAsInt("JWT_EXPIRE_HOURS", 24*7), // 默认 7 天
		},
		CORS: CORSConfig{
			// 开发环境允许所有 localhost 端口
			AllowOrigins: []string{
				"http://localhost:4321",
				"http://localhost:4322",
				"http://localhost:4323",
				"http://localhost:4324",
				"http://localhost:4325",
				"http://localhost:4326",
				"http://localhost:3000",
				"http://127.0.0.1:4321",
				"http://127.0.0.1:4322",
				"http://127.0.0.1:4323",
				"http://127.0.0.1:4324",
				"http://127.0.0.1:4325",
				"http://127.0.0.1:4326",
				"http://127.0.0.1:3000",
			},
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}