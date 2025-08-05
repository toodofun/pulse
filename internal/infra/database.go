package infra

import (
	"fmt"
	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"pulse/internal/config"
)

type Database struct {
	*gorm.DB
}

func NewDatabase(cfg config.Database) (*Database, error) {
	var driver gorm.Dialector

	switch cfg.Driver {
	case "sqlite":
		driver = sqlite.Open(cfg.DSN)
	case "mysql":
		driver = mysql.Open(cfg.DSN)
	case "postgres":
		driver = postgres.Open(cfg.DSN)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}

	dbLogLevel := logger.Silent

	client, err := gorm.Open(driver, &gorm.Config{
		Logger: logger.Default.LogMode(dbLogLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	db, err := client.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	db.SetMaxIdleConns(cfg.MaxIdleConn)
	db.SetMaxOpenConns(cfg.MaxOpenConn)
	db.SetConnMaxLifetime(cfg.ConnMaxLift)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdle)

	return &Database{DB: client}, nil
}
