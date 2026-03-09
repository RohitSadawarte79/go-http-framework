package domain

import (
	"errors"
	"time"
)

type ContextKey string

const ParamKey ContextKey = "params"

type User struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Age       int
	CreatedAt time.Time
}

type UserRepository interface {
	FindByID(id int) (*User, error)
	FindAll() ([]*User, error)
	Create(user *User) error
	FindByEmail(email string) (*User, error)
}

var ErrUserNotFound = errors.New("user not found")
var ErrDuplicateUser = errors.New("user already exists")
