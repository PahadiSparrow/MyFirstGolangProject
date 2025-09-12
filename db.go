package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func init() {
	var err error
	// Update with your DB connection info
	DB, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/myapp")
	if err != nil {
		log.Fatal("DB Connection Error: ", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("DB Ping Failed: ", err)
	}
}
