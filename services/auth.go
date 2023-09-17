package services

import (
	"github.com/ztx-lyghters/kode-notes/core"
	"github.com/ztx-lyghters/kode-notes/repository"
)

type Auth struct {
	repository repository.Authorization
}

func NewAuthService(repo *repository.Repository) *Auth {
	return &Auth{repository: repo}
}

func (s *Auth) CreateUser(user *core.User) (int, error) {

	user.Password = GeneratePasswordHash(user.Password)
	user_id, err := s.repository.CreateUser(user)
	if err != nil {
		return 0, err
	}

	return user_id, nil
}

func (s *Auth) GenerateToken(user *core.User) (string, error) {

	user.Password = GeneratePasswordHash(user.Password)
	user, err := s.repository.GetUser(user.Username,
		user.Password)
	if err != nil {
		return "", err
	}

	return assembleTokenJWT(user), nil
}

func (s *Auth) CheckUserExistence(user_id int) bool {
	return s.repository.CheckUserExistence(user_id)
}
