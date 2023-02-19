package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

var (
	ErrNoRecord = fmt.Errorf("no matching record found")
	insertOp    = "insert"
	deleteOp    = "delete"
	updateOp    = "update"
)

type Database struct {
	Conn   *sql.DB
	Logger zerolog.Logger
}

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	DbName   string
	Logger   zerolog.Logger
}

func Init(cfg Config) (Database, error) {
	db := Database{}
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DbName)
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return db, err
	}

	db.Conn = conn
	db.Logger = cfg.Logger
	err = db.Conn.Ping()
	if err != nil {
		return db, err
	}
	return db, nil
}
