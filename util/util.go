package util

import (
	"database/sql"
	_ "github.com/bmizerany/pq"
	"log"
	"os"
)

var db_url string

func init() {
	db_url = os.Getenv("DATABASE_URL")
	if db_url == "" {
		db_url = "user=weightlog dbname=weightlog password=weightlog sslmode=disable"
	}
}

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

func GetDb() *sql.DB {
	db, err := sql.Open("postgres", db_url)
	if err != nil {
		panic(err)
	}
	return db
}
func GetLastId(db DB, tablename string) (id int64, err error) {
	row := db.QueryRow("SELECT currval(pg_get_serial_sequence($1, 'id'));", tablename)
	err = row.Scan(&id)
	return
}
