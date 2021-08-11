package mysql

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"log"
)

var db *sql.DB
var mock sqlmock.Sqlmock

func NewMock() {
	var err error
	db, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
}

func GetSqlDB() *sql.DB {
	return db
}

func GetSqlMock() sqlmock.Sqlmock {
	return mock
}