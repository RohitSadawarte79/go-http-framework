package service

import (
	"errors"

	"github.com/RohitSadawarte79/go-http-framework/internal/domain"
)

var ErrValidation = errors.New("validation failed")

type UserService struct {
	repo domain.UserRepository // interface , not concreate type!
}

func NewUserService(repo domain.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetByID(id int) (*domain.User, error) {
	return s.repo.FindByID(id)
}

func (s *UserService) List() ([]*domain.User, error) {
	return s.repo.FindAll()
}

func (s *UserService) GetByEmail(email string) (*domain.User, error) {
	return s.repo.FindByEmail(email)
}

func (s *UserService) Create(user *domain.User) error {
	if user.FirstName == "" || user.LastName == "" {
		return ErrValidation
	}

	return s.repo.Create(user)
}
