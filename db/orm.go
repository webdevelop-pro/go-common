package db

import (
	"fmt"

	baseLogger "github.com/webdevelop-pro/go-common/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Storage ...
type Storage interface {
	GetOrm() *gorm.DB
}

type orm struct {
	db *gorm.DB
}

// NewOrm ...
func NewOrm(cfg Config, log baseLogger.Logger) (Storage, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Database,
		cfg.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &orm{
		db: db,
	}, nil
}

// GetOrm ...
func (o orm) GetOrm() *gorm.DB {
	return o.db
}
