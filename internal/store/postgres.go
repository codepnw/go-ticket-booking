package store

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func InitPostgresDB(addr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", addr)
	if err != nil {
		return nil, fmt.Errorf("connect database failed: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("connect database failed: %v", err)
	}

	log.Println("database connected...")

	return db, err
}
