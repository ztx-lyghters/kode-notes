package services

import (
	"github.com/ztx-lyghters/kode-notes/core"
	"github.com/ztx-lyghters/kode-notes/repository"
)

type Notes struct {
	repository repository.Notes
}

func NewNotesService(repo *repository.Repository) *Notes {
	return &Notes{repository: repo}
}

func (s *Notes) Create(note *core.Note) (uint, error) {
	return s.repository.CreateNote(note)
}

func (s *Notes) GetAll(user_id uint) ([]*core.Note, error) {
	return s.repository.GetAllNotes(user_id)
}
