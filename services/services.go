package services

import (
	"github.com/ztx-lyghters/kode-notes/core"
	"github.com/ztx-lyghters/kode-notes/repository"
)

type AuthorizationService interface {
	CreateUser(*core.User) (uint, error)
	GenerateToken(*core.User) (string, error)
	ValidateToken(string) (uint, error)
	CheckUserExistence(uint) bool
}

type NotesService interface {
	Create(*core.Note) (uint, error)
	GetAll(user_id uint) ([]*core.Note, error)
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
