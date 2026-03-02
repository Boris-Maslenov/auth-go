package psql

import (
	"auth-test/internal/service"
	"database/sql"
	"fmt"
)

type Tokens struct {
	db *sql.DB
}

func (r *Tokens) Create(td service.RefreshData) error {
	_, err := r.db.Exec("INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)", td.UserId, td.Token, td.ExpiresAt)
	return err
}

func (r *Tokens) Get(token string) (service.RefreshData, error) {
	var t service.RefreshData
	err := r.db.QueryRow("SELECT id, user_id, token, expires_at FROM refresh_tokens WHERE token=$1", token).Scan(&t.Id, &t.UserId, &t.Token, &t.ExpiresAt)
	if err != nil {
		return t, err
	}

	fmt.Println("GET t.UserId", t.UserId)

	_, err = r.db.Exec("DELETE FROM refresh_tokens WHERE user_id=$1", t.UserId)

	return t, err
}

func NewTokens(db *sql.DB) *Tokens {
	return &Tokens{db}
}
