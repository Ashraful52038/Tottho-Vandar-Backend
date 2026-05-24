package postgres

import (
	"fmt"
	"log"

	"github.com/Ashraful52038/tottho-vandar-backend/internal/config"
	"github.com/Ashraful52038/tottho-vandar-backend/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Dhaka",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	// Auto migrate
	err = db.AutoMigrate(
		&domain.User{},
		&domain.Post{},
		&domain.Comment{},
		&domain.Like{},
		&domain.Tag{},
		&domain.PostTag{},
		&domain.Follow{},
		&domain.Notification{},
	)
	if err != nil {
		return nil, err
	}

	log.Println("Database connected successfully")
	return db, nil
}
