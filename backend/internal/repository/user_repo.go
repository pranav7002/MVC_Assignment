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

func (userRepo *UserRepository) GetAttributeFromUsername(username, column string) (string, error) {
    ctx := context.Background()	

	var attribute string

    query := `SELECT $1 FROM users WHERE username = "$2"`
	err := userRepo.DB.QueryRow(ctx, query, column, username).Scan(&attribute)
	
	if err != nil {
		return "", err
	}

	return attribute, nil
}  
