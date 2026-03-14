package repository

import (
	"database/sql"
	"errors"

	"github.com/RohitSadawarte79/go-http-framework/internal/domain"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{
		db: db,
	}
}

func (r *PostgresUserRepository) FindByID(id int) (*domain.User, error) {
	var user domain.User
	row := r.db.QueryRow("SELECT id, first_name, last_name, email, age, created_at FROM users WHERE id=$1;", id)
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Age,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (r *PostgresUserRepository) FindAll() ([]*domain.User, error) {
	var userList []*domain.User
	rows, err := r.db.Query("SELECT id, first_name, last_name, email, age, created_at FROM users ORDER By id;")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var user domain.User
		err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Age,
			&user.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		userList = append(userList, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userList, nil
}

func (r *PostgresUserRepository) Create(user *domain.User) error {

	row := r.db.QueryRow("INSERT INTO users (first_name, last_name, email, age) VALUES ($1, $2, $3, $4) RETURNING id, created_at;",
		user.FirstName, user.LastName, user.Email, user.Age)

	err := row.Scan(
		&user.ID,
		&user.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresUserRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User

	row := r.db.QueryRow("SELECT id, first_name, last_name, email, age, created_at FROM users WHERE email=$1;", email)

	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Age,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}
