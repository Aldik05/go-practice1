package main

import (
    "database/sql"
    "log"

    _ "github.com/lib/pq"
)

var db *sql.DB

func initDB() {
	var err error
	connStr := "host=localhost port=5432 user=postgres password=Aldik2005 dbname=shop sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
}
