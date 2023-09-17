package services

import (
	"example.com/kode-notes/core"
	"example.com/kode-notes/repository"
)

type Notes struct {
	repository repository.Notes
}

func NewNotesService(repo *repository.Repository) *Notes {
	return &Notes{repository: repo}
}

func (s *Notes) Create(note *core.Note) (int, error) {
	return s.repository.CreateNote(note)
}

func (s *Notes) GetAll(user_id int) ([]*core.Note, error) {
	return s.repository.GetAllNotes(user_id)
}
