package db

import (
	"time"

	sql "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func SetupDB(addr, dbName, user, passwd, logLevel string) (*gorm.DB, error) {
	dsn := FormatDSN(addr, dbName, user, passwd)
	level := formatLogLevel(logLevel)
	engine, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(level),
	})
	if err != nil {
		return nil, err
	}
	db, err := engine.DB()
	if err != nil {
		return nil, err
	}
	// Set the maximum lifetime of the connection (less than server setting)
	db.SetConnMaxLifetime(30 * time.Minute)

	return engine, nil
}

func GetDB() *gorm.DB {
	return DB
}

func Close() error {
	db, err := DB.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

// FormatDSN formats the given Config into a DSN string which can be passed to the driver.
func FormatDSN(addr, dbName, user, passwd string) string {
	cfg := sql.Config{
		User:                 user,
		Passwd:               passwd,
		Net:                  "tcp",
		Addr:                 addr,
		DBName:               dbName,
		ParseTime:            true,
		AllowNativePasswords: true,
	}
	return cfg.FormatDSN()
}

func formatLogLevel(level string) logger.LogLevel {
	switch level {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	}
	return logger.Silent
}
