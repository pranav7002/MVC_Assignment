package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pranav7002/MVC_Assignment/internal/models"
)

type TroopRepository struct {
	DB *pgxpool.Pool
}

func (r *TroopRepository) GetUserTrainedTroops(userID string) ([]models.TroopTrained, error) {
	ctx := context.Background()

	query := `SELECT * FROM troops_trained WHERE user_id = $1`
	rows, err := r.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	troops, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.TroopTrained])
	if err != nil {
		return nil, err
	}

	return troops, nil
}

func (r *TroopRepository) TrainTroop(userID string, troopName string, quantity int, trainingCost int) error {
	ctx := context.Background()

	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `
		UPDATE village 
		SET elixir = elixir - $1 
		WHERE user_id = $2
	`

	_, err = tx.Exec(ctx, query, trainingCost, userID)
	if err != nil {
		return err
	}

	query = `
		INSERT INTO troops_trained (user_id, troop_name, quantity)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, troop_name) 
		DO UPDATE SET quantity = troops_trained.quantity + EXCLUDED.quantity
	`

	_, err = tx.Exec(ctx, query, userID, troopName, quantity)
	if err != nil {
		return err
	}

    if err := tx.Commit(ctx); err != nil {
        return err
    }

	return nil
}

func (r *TroopRepository) DeleteTroop(userID, troopName string) error {
	ctx := context.Background()

	query := `
		UPDATE troops_trained 
		SET quantity = quantity - 1 
		WHERE user_id = $1 AND troop_name = $2 AND quantity > 0
	`

	if _, err := r.DB.Exec(ctx, query, userID, troopName); err != nil {
		return err
	}

	return nil
}