package db

import (
	"database/sql/driver"

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

type TestCase struct {
	Name       string
	Query      MockQuery
	Args       interface{}
	WantErr    error
	WantResult interface{}
}

// MockQuery contains the necessary data required to mock a SQL query, from the
// query string, to the arguments passed into any given query.
type MockQuery struct {
	SQL  string
	Args []driver.Value
	// Rows are rows created from SQLmock
	// Depracated: This field is used in the older ExpectQueries function.
	// It shouldnt be used anymore as it requires a reference back to the original
	// mock controller. Use MockRows instead.
	Rows     *sqlmock.Rows
	MockRows *MockRows
	Results  driver.Result
	Err      error
}

type MockRows struct {
	Columns []string
	Rows    [][]driver.Value
}

func ExpectExec(mock sqlmock.Sqlmock, q MockQuery) {
	mock.ExpectBegin()
	mock.ExpectExec(q.SQL).
		WithArgs(q.Args...).
		WillReturnResult(q.Results)
	mock.ExpectCommit()
}

func ExpectQuery(mock sqlmock.Sqlmock, q MockQuery) {
	mock.ExpectQuery(q.SQL).
		WithArgs(q.Args...).
		WillReturnRows(q.Rows)
}
