package gorm

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Config struct {
	Type         string `json:"type" yaml:"type"`
	Host         string `json:"host" yaml:"host"`
	Port         int    `json:"port" yaml:"port"`
	Username     string `json:"username" yaml:"username"`
	Password     string `json:"password" yaml:"password"`
	Database     string `json:"database" yaml:"database"`
	SSLMode      string `json:"ssl_mode" yaml:"ssl_mode"`
	MaxOpenConns int    `json:"max_open_conns" yaml:"max_open_conns"`
	MaxIdleConns int    `json:"max_idle_conns" yaml:"max_idle_conns"`
	TablePrefix  string `json:"table_prefix" yaml:"table_prefix"`
	Debug        bool   `json:"debug" yaml:"debug"`
}

type Database struct {
	DB     *gorm.DB
	Config *Config
}

func NewDatabase(config *Config) (*Database, error) {
	var db *gorm.DB
	var err error

	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   config.TablePrefix,
			SingularTable: true,
		},
	}

	if config.Debug {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	} else {
		gormConfig.Logger = logger.Default.LogMode(logger.Error)
	}

	switch config.Type {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.Username, config.Password, config.Host, config.Port, config.Database)
		db, err = gorm.Open(mysql.Open(dsn), gormConfig)
	case "postgres":
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			config.Host, config.Port, config.Username, config.Password, config.Database, config.SSLMode)
		db, err = gorm.Open(postgres.Open(dsn), gormConfig)
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(config.Database), gormConfig)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", config.Type)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	if config.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	}

	if config.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	}

	sqlDB.SetConnMaxLifetime(time.Hour)

	return &Database{
		DB:     db,
		Config: config,
	}, nil
}

func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	return sqlDB.Close()
}

func (d *Database) Ping() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	return sqlDB.Ping()
}

func (d *Database) Migrate(models ...interface{}) error {
	return d.DB.AutoMigrate(models...)
}

func (d *Database) Transaction(fn func(tx *gorm.DB) error) error {
	return d.DB.Transaction(fn)
}
