package repository

import (
	"database/sql"

	"github.com/ztx-lyghters/kode-notes/core"
)

type Authorization interface {
	CreateUser(*core.User) (uint, error)
	GetUser(string, string) (*core.User, error)
	CheckUserExistence(uint) bool
}

type Notes interface {
	CreateNote(*core.Note) (uint, error)
	GetAllNotes(uint) ([]*core.Note, error)
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
