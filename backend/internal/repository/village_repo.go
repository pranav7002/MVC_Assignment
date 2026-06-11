package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pranav7002/MVC_Assignment/internal/models"
)

type VillageRepository struct {
	DB *pgxpool.Pool
}

func (villageRepo *VillageRepository) FetchUserBuildings(userID string) ([]models.Building, error) {
	ctx := context.Background()

	query := `SELECT * FROM building_instance WHERE user_id = $1`
	rows, err := villageRepo.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	buildings, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Building])
	if err != nil {
		return nil, err
	}

	return buildings, nil
}