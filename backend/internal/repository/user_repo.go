package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	DB *pgxpool.Pool
}

func (userRepo *UserRepository) InsertUser(username, hash string) error {
    ctx := context.Background()	

    query := `INSERT INTO users (id, username, password_hash) VALUES ( gen_random_uuid(), $1, $2 )`
    _, err := userRepo.DB.Exec(ctx, query, username, hash)
    
    return err 
}

