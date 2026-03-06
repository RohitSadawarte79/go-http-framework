package repository

import (
	"sort"

	"github.com/RohitSadawarte79/go-http-framework/internal/domain"
)

type MemoryUserRepository struct {
	users map[int]*domain.User
}

func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		users: make(map[int]*domain.User, 0),
	}
}

func (r *MemoryUserRepository) FindByID(id int) (*domain.User, error) {
	user, ok := r.users[id]

	if !ok {
		return nil, domain.ErrUserNotFound
	}

	return user, nil
}

func (r *MemoryUserRepository) FindAll() ([]*domain.User, error) {
	var usersList []*domain.User

	for _, user := range r.users {
		usersList = append(usersList, user)
	}

	sort.Slice(usersList, func(i, j int) bool {
		return usersList[i].ID < usersList[j].ID
	})

	return usersList, nil
}

func (r *MemoryUserRepository) Create(user *domain.User) error {
	r.users[user.ID] = user
	return nil
}
