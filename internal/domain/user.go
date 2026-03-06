package domain

import "errors"

type ContextKey string

const ParamKey ContextKey = "params"

type User struct {
	ID        int
	FirstName string
	LastName  string
	Age       int
}

type UserRepository interface {
	FindByID(id int) (*User, error)
	FindAll() ([]*User, error)
	Create(user *User) error
}

var ErrUserNotFound = errors.New("user not found")
var ErrDuplicateUser = errors.New("user already exists")
