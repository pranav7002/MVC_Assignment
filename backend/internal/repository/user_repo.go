package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	DB *pgxpool.Pool
}

func (r *UserRepository) InsertUser(username, hash string) error {
	ctx := context.Background()

	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	userID := uuid.New().String()
	query := `INSERT INTO users (id, username, password_hash) VALUES ( $3, $1, $2 )`

	_, err = tx.Exec(ctx, query, username, hash, userID)
	if err != nil {
		return err
	}

	query = `
	INSERT INTO village (user_id, town_hall_level, gold, elixir)
	VALUES ($1, $2, $3, $4)
	`
	_, err = tx.Exec(ctx, query, userID, 1, 100000, 100000)
	if err != nil {
		return err
	}

	query = `
		INSERT INTO building_instance (user_id, building_type, building_name, level, pos_x, pos_y, size, is_upgrading, hp)
		VALUES ($1, 'town_hall', 'Town Hall', 1, 8, 8, 4, false, 1500)
	`
	_, err = tx.Exec(ctx, query, userID)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
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

func (r *UserRepository) GetTrophies(userID string) (int, error) {
	ctx := context.Background()
	var trophies int
	query := `SELECT trophies FROM users WHERE id = $1`
	err := r.DB.QueryRow(ctx, query, userID).Scan(&trophies)
	if err != nil {
		return 0, err
	}
	return trophies, nil
}
