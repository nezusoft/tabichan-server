package user

import (
	"database/sql"
	"errors"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) CreateUser(user *User) error {
	query := "INSERT INTO users (username, password, email) VALUES ($1, $2, $3) RETURNING id"
	return r.DB.QueryRow(query, user.Username, user.Password, user.Email).Scan(&user.ID)
}

func (r *Repository) GetUserByUsername(username string) (*User, error) {
	query := "SELECT id, username, password, email FROM users WHERE username=$1"
	user := &User{}
	err := r.DB.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password, &user.Email)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	return user, err
}
