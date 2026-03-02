package psql

import (
	"database/sql"
)

type Users struct {
	db *sql.DB
}

func (r *Users) Create(email, name, hashPassword string) error {
	var id int64
	// res, err := r.db.Exec("INSERT INTO users (email, name, password) VALUES ($1, $2, $3)", a, b, c)
	err := r.db.QueryRow("INSERT INTO users (email, name, password) VALUES ($1, $2, $3) RETURNING id", email, name, hashPassword).Scan(&id)
	if err != nil {
		return err
	}
	return nil
}

func (r *Users) GetByCredentials(email, hashPassword string) (int64, error) {
	var id int64
	err := r.db.QueryRow("SELECT id FROM users WHERE email=$1 AND password=$2", email, hashPassword).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

func NewUsers(db *sql.DB) *Users {
	return &Users{db}
}
