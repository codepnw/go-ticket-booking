package database

import (
	"context"
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

type TxManager interface {
	WithTx(ctx context.Context, fn func(tx *sql.Tx) error) error
}

type sqlTxManager struct {
	db *sql.DB
}

func NewSqlTxManager(db *sql.DB) TxManager {
	return &sqlTxManager{db: db}
}

func (m *sqlTxManager) WithTx(ctx context.Context, fn func(tx *sql.Tx) error) (err error) {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}
