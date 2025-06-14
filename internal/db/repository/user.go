package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/prajwalbharadwajbm/broker/internal/db"
)

func AddUser(ctx context.Context, email string, hashedPassword []byte) (string, error) {
	db := db.GetClient()

	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id`
	var userId string
	err := db.QueryRowContext(dbCtx, query, email, hashedPassword).Scan(&userId)
	if err != nil {
		return "", err
	}
	return userId, nil
}

func GetUserByEmail(ctx context.Context, email string) (string, []byte, error) {
	db := db.GetClient()

	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var userId string
	var hashedPassword []byte

	query := `SELECT id, password_hash FROM users WHERE email = $1`
	err := db.QueryRowContext(dbCtx, query, email).Scan(&userId, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil, errors.New("user not found")
		}
		return "", nil, err
	}
	return userId, hashedPassword, nil
}
