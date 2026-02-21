package psql

import (
	"database/sql"
	"fmt"
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

	fmt.Printf("Новый пользователь успешно зареган c id: %d \n", id)

	return nil
}

func NewUsers(db *sql.DB) *Users {
	return &Users{db}
}
