package util

import (
	"database/sql"
	_ "github.com/bmizerany/pq"
	"log"
)

type DB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

func GetTestDb() *sql.DB {
	var err error
	db, err := sql.Open("postgres", "user=weightlog dbname=weightlog_test password=weightlog sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	return db
}
