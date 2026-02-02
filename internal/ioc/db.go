package ioc

import (
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConfig struct {
	Driver       string
	Host         string
	Port         int
	Database     string
	Username     string
	Password     string
	MaxIdleConns int
	MaxOpenConns int
}

func InitDB() *gorm.DB {
	cfg := &DBConfig{
		Driver:       viper.GetString("db.driver"),
		Host:         viper.GetString("db.host"),
		Port:         viper.GetInt("db.port"),
		Database:     viper.GetString("db.database"),
		Username:     viper.GetString("db.username"),
		Password:     viper.GetString("db.password"),
		MaxIdleConns: viper.GetInt("db.max_idle_conns"),
		MaxOpenConns: viper.GetInt("db.max_open_conns"),
	}

	var dsn string
	var dialector gorm.Dialector

	switch cfg.Driver {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.Username,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.Database,
		)
		dialector = mysql.Open(dsn)
	case "postgres":
		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			cfg.Host,
			cfg.Port,
			cfg.Username,
			cfg.Password,
			cfg.Database,
		)
		dialector = postgres.Open(dsn)
	default:
		panic("unsupported database driver: " + cfg.Driver)
	}

	// 初始化数据库连接
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("failed to connect to database: " + err.Error())
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to get underlying sql.DB: " + err.Error())
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)

	return db
}
