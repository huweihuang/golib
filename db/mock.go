package db

import (
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetDBMock() (*gorm.DB, sqlmock.Sqlmock, error) {
	// mock一个*sql.DB对象，不需要连接真实的数据库
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		return nil, nil, err
	}

	//测试时不需要真正连接数据库
	c := &gorm.Config{
		DisableAutomaticPing: true,
		Logger:               logger.Default.LogMode(logger.Info),
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      mockDB,
		SkipInitializeWithVersion: true}), c)
	if err != nil {
		return nil, nil, err
	}

	return db, mock, nil
}
