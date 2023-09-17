package services

import (
	"example.com/kode-notes/core"
	"example.com/kode-notes/repository"
)

type AuthorizationService interface {
	CreateUser(*core.User) (int, error)
	GenerateToken(*core.User) (string, error)
	ValidateToken(string) (int, error)
	CheckUserExistence(int) bool
}

type NotesService interface {
	Create(*core.Note) (int, error)
	GetAll(user_id int) ([]*core.Note, error)
}

type SpellerService struct {
	Enabled bool
	Yandex  Speller
}

type Services struct {
	AuthorizationService
	NotesService
	SpellerService
}

func New(repo *repository.Repository) *Services {
	return &Services{
		AuthorizationService: NewAuthService(repo),
		NotesService:         NewNotesService(repo),
		SpellerService:       NewSpellersService(),
	}
}
