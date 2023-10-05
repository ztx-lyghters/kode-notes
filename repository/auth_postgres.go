package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/ztx-lyghters/kode-notes/core"
)

type AuthPostgres struct {
	db *sql.DB
}

func NewAuthPostgres(db *sql.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user *core.User) (uint, error) {

	var id uint

	query := fmt.Sprintf("INSERT INTO %s(username, "+
		"password_hash) VALUES ($1, $2) RETURNING id",
		USERS_TABLE)

	row := r.db.QueryRow(query, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *AuthPostgres) CheckUserExistence(user_id uint) bool {
	var username string
	query := fmt.Sprintf("SELECT username FROM %s WHERE id=$1", USERS_TABLE)
	row := r.db.QueryRow(query, user_id)
	err := row.Scan(&username)
	if err != nil || strings.TrimSpace(username) == "" {
		return false
	}

	return true
}

func (r *AuthPostgres) GetUser(username string,
	password_hash string) (*core.User,
	error) {
	var user core.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE username=$1 AND password_hash=$2", USERS_TABLE)

	row := r.db.QueryRow(query, username, password_hash)
	err := row.Scan(&user.Id)
	if err != nil {
		return nil, err
	}
	user.Username = username
	user.Password = password_hash

	return &user, nil
}
