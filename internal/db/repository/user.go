package repository

import (
	"context"
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
