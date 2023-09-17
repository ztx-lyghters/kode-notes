package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	USERS_TABLE       = "users"
	NOTES_TABLE       = "notes"
	USERS_NOTES_TABLE = "users_notes"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewDBPostgres(cfg *Config) (*sql.DB, error) {
	data_src := fmt.Sprintf("host=%s port=%s user=%s"+
		" password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := sql.Open("postgres", data_src)
	if err != nil {
		return nil, err
	}

	return db, nil
}
