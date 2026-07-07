package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BattleRepository struct {
	DB *pgxpool.Pool
}

func (r *BattleRepository) StoreBattle(userID, defendersID, result string, stars, destructionPct, goldLooted, elixirLooted, attackerTrophyChange, defenderTrophyChange int) error {
	ctx := context.Background()

	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO battles 
			(attacker_id, defender_id, gold_looted, elixir_looted, result, destruction_percentage) 
		VALUES 
			( $1, $2, $3, $4, $5, $6 )
	`

	if _, err := tx.Exec(ctx, query, userID, defendersID, goldLooted, elixirLooted, result, destructionPct); err != nil {
		return err
	}
	query = `
	UPDATE village 
	SET 
		gold = LEAST(gold + $1, max_gold), 
		elixir = LEAST(elixir + $2, max_elixir) 
	WHERE user_id = $3
	`
	if _, err := tx.Exec(ctx, query, goldLooted, elixirLooted, userID); err != nil {
		return err
	}
	query = `
	UPDATE village 
	SET 
		gold = GREATEST(0, gold - $1),
		elixir = GREATEST(0, elixir - $2)
	WHERE user_id = $3
	`
	if _, err := tx.Exec(ctx, query, goldLooted, elixirLooted, defendersID); err != nil {
		return err
	}

	query = `
	UPDATE users 
	SET trophies = GREATEST(0, trophies + $1)
	WHERE id = $2
	`
	if _, err := tx.Exec(ctx, query, attackerTrophyChange, userID); err != nil {
		return err
	}

	query = `
	UPDATE users 
	SET trophies = GREATEST(0, trophies + $1)
	WHERE id = $2
	`
	if _, err := tx.Exec(ctx, query, defenderTrophyChange, defendersID); err != nil {
		return err
	}

	query = `
	DELETE FROM 
		troops_trained 
	WHERE 
		user_id = $1;
	`
	if _, err := tx.Exec(ctx, query, userID); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (r *BattleRepository) GetArmyCount(userID string) (int, error) {
	ctx := context.Background()

	var count int
	query := `
		SELECT 
			COALESCE(SUM(quantity), 0)
		FROM 
			troops_trained
		WHERE 
			user_id = $1;
	`
	err := r.DB.QueryRow(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
