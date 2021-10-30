package dao

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
	"fmt"
)

func InitDB() *sql.DB {
	db, err := sql.Open("sqlite3", "db/corner_data.db") // Open the created SQLite File
	
	if err != nil {
		panic(err)
	}
	
	fmt.Println("Database opened")
	
	return db
}

func CloseDB(db *sql.DB) {
	db.Close()
	fmt.Println("Database closed")
}