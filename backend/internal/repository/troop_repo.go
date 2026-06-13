package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type TroopRepository struct {
	DB *pgxpool.Pool
}
