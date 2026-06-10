package repository

import "github.com/jackc/pgx/v5/pgxpool"

type UserRepository struct {
	DB *pgxpool.Pool
}

func (userRepo *UserRepository) InsertUser(username, hash string) error {
	return nil
}