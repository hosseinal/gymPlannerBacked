package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	*sqlx.DB
}

func NewDB(host, port, user, password, dbname string) (*DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &DB{db}, nil
} 