package repository

import (
	"database/sql"

	"example.com/kode-notes/core"
)

type Authorization interface {
	CreateUser(*core.User) (int, error)
	GetUser(string, string) (*core.User, error)
	CheckUserExistence(int) bool
}

type Notes interface {
	CreateNote(*core.Note) (int, error)
	GetAllNotes(int) ([]*core.Note, error)
}

type Repository struct {
	Authorization
	Notes
}

func New(db *sql.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Notes:         NewNotesPostgres(db),
	}
}
