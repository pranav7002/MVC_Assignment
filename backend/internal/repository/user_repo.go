package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	DB *pgxpool.Pool
}

func (r *UserRepository) InsertUser(username, hash string) error {
	ctx := context.Background()

	query := `INSERT INTO users (id, username, password_hash) VALUES ( gen_random_uuid(), $1, $2 )`
	_, err := r.DB.Exec(ctx, query, username, hash)

	return err
}

func (r *UserRepository) GetAttributeFromUsername(username, column string) (string, error) {
	ctx := context.Background()

	var attribute string

	query := fmt.Sprintf(`SELECT %s FROM users WHERE username = $1`, column)
	err := r.DB.QueryRow(ctx, query, username).Scan(&attribute)

	if err != nil {
		return "", err
	}

	return attribute, nil
}
