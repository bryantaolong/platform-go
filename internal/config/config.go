package config

import (
	"log"
	"os"

	"github.com/bryantaolong/platform/internal/model"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Port      string
	JWTSecret string
	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
	SSLMode   string
}

func Load() *Config {
	// 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		log.Println("未找到 .env 文件，使用默认配置")
	}

	return &Config{
		Port:      getEnv("PORT", "8081"),
		JWTSecret: getEnv("JWT_SECRET", "ThisIsAVerySecretKeyForYourJWTAuthenticationAndItShouldBeLongEnough"),
		DBHost:    getEnv("DB_HOST", "localhost"),
		DBPort:    getEnv("DB_PORT", "5432"),
		DBUser:    getEnv("DB_USER", "platform_user"),
		DBPass:    getEnv("DB_PASS", "123456"),
		DBName:    getEnv("DB_NAME", "platform"),
		SSLMode:   getEnv("SSL_MODE", "disable"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func InitDB(cfg *Config) *gorm.DB {
	dsn := "host=" + cfg.DBHost +
		" user=" + cfg.DBUser +
		" password=" + cfg.DBPass +
		" dbname=" + cfg.DBName +
		" port=" + cfg.DBPort +
		" sslmode=" + cfg.SSLMode +
		" TimeZone=Asia/Shanghai"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 安全迁移模式
	if err := db.Migrator().AutoMigrate(&model.User{}); err != nil {
		log.Fatal("数据库迁移失败:", err)
	}

	// 确保索引存在
	if !db.Migrator().HasIndex(&model.User{}, "idx_username") {
		if err := db.Migrator().CreateIndex(&model.User{}, "Username"); err != nil {
			log.Fatal("创建索引失败:", err)
		}
	}

	if !db.Migrator().HasIndex(&model.User{}, "idx_deleted") {
		if err := db.Migrator().CreateIndex(&model.User{}, "Deleted"); err != nil {
			log.Fatal("创建索引失败:", err)
		}
	}

	return db
}
