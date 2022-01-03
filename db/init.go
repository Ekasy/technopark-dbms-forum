package db

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/stdlib"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(connString string) (*Database, error) {
	db, err := sql.Open("pgx", connString)
	if err != nil {
		log.Fatalf("could not connect to database: %s", err.Error())
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("error with ping: %s", err.Error())
		return nil, err
	}

	db.SetMaxOpenConns(10)

	return &Database{db: db}, nil
}

func (db *Database) Close() {
	db.db.Close()
}

func (db *Database) GetPool() *sql.DB {
	return db.db
}
