package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type PostgresDB struct {
	Conn *sql.DB
}

func NewPostgresDB(user, password, dbname, host string, port int) (*PostgresDB, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:5432/%s", user, password, host, dbname)
	fmt.Println("CONNSTR", connStr)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresDB{Conn: db}, nil
}

func (db *PostgresDB) Close() error {
	return db.Conn.Close()
}
