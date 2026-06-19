package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pranav7002/MVC_Assignment/internal/models"
)

type VillageRepository struct {
	DB *pgxpool.Pool
}

func (r *VillageRepository) GetUserBuildings(userID string) ([]models.Building, error) {
	ctx := context.Background()

	query := `SELECT * FROM building_instance WHERE user_id = $1`
	rows, err := r.DB.Query(ctx, query, userID)
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

func (r *VillageRepository) GetVillage(userID string) (models.Village, error) {
	ctx := context.Background()

	query := `SELECT * FROM village WHERE user_id = $1`
	rows, err := r.DB.Query(ctx, query, userID)
	if err != nil {
		return models.Village{}, err
	}
	defer rows.Close()

	village, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.Village])
	if err != nil {
		return models.Village{}, err
	}
	return village, nil
}

func (r *VillageRepository) GetBuildingCount(userID string, buildingType string, buildingName string) (int, error) {
	ctx := context.Background()

	var count int
	query := `SELECT COUNT(*) FROM building_instance WHERE user_id = $1 AND building_type = $2 AND building_name = $3`
	err := r.DB.QueryRow(ctx, query, userID, buildingType, buildingName).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *VillageRepository) InsertBuilding(userID string, buildingReqBody models.BuildingCreationRequestBody, hp int, size int) error {
	ctx := context.Background()

	query := `
		INSERT INTO building_instance 
		(user_id, building_type, building_name, level, pos_x, pos_y, size, is_upgrading, hp) 
		VALUES ($1, $2, $3, 1, $4, $5, $6, false, $7)
	`

	_, err := r.DB.Exec(ctx, query,
		userID,
		buildingReqBody.BuildingType,
		buildingReqBody.BuildingName,
		buildingReqBody.PosX,
		buildingReqBody.PosY,
		size,
		hp,
	)

	return err
}

func (r *VillageRepository) MoveBuilding(userID string, buildingID int64, posX, posY int) error {
	ctx := context.Background()

	query := `
		UPDATE building_instance 
		SET pos_x = $1, pos_y = $2 
		WHERE id = $3 AND user_id = $4
	`

	result, err := r.DB.Exec(ctx, query, posX, posY, buildingID, userID)
	if err != nil {
		return err
	}

	// Check if the building belongs to the user
	if result.RowsAffected() == 0 {
		return fmt.Errorf("building not found or does not belong to user")
	}

	return nil
}

func (r *VillageRepository) RemoveResource(userID string, resourceType string, amount int) error {
	ctx := context.Background()

	query := fmt.Sprintf(`
		UPDATE village 
		SET %s = %s - $1 
		WHERE user_id = $2
	`, resourceType, resourceType)

	_, err := r.DB.Exec(ctx, query, amount, userID)
	if err != nil {
		return err
	}

	return nil
}

func (r *VillageRepository) AddResourceFromColletor(userID string, resourceType string, amount int, collectionTime time.Time) error {
	ctx := context.Background()

	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := fmt.Sprintf(`
		UPDATE village 
		SET %s = %s + $1 
		WHERE user_id = $2
	`, resourceType, resourceType)

	_, err = tx.Exec(ctx, query, amount, userID)
	if err != nil {
		return err
	}

	query = fmt.Sprintf(`
		UPDATE village 
		SET %s_last_collected_at = $1
		WHERE user_id = $2	
	`, resourceType)

	_, err = tx.Exec(ctx, query, collectionTime, userID)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (r *VillageRepository) GetBuilding(buildingID int64) (models.Building, error) {
	ctx := context.Background()

	query := `
	SELECT *
	FROM 
		building_instance
	WHERE 
		id = $1
	`

	rows, err := r.DB.Query(ctx, query, buildingID)
	if err != nil {
		return models.Building{}, err
	}
	defer rows.Close()

	building, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.Building])
	if err != nil {
		return models.Building{}, err
	}

	return building, nil
}

func (r *VillageRepository) UpgradeBuilding(userID string, buildingID int64, hp int, upgradeCostType string, upgradeCost int) error {
	ctx := context.Background()

	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `
	UPDATE building_instance
	SET 
		hp = $1, level = level + 1
	WHERE 
		id = $2
		AND user_id = $3
	`
	result, err := tx.Exec(ctx, query, hp, buildingID, userID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("building not found or does not belong to user")
	}

	query = fmt.Sprintf(`
		UPDATE village 
		SET %s = %s - $1 
		WHERE user_id = $2
	`, upgradeCostType, upgradeCostType)

	_, err = tx.Exec(ctx, query, upgradeCost, userID)
	if err != nil {
		return err
	}
	tx.Commit(ctx)

	return nil
}

func (r *VillageRepository) GetUserBuildingsByName(userID string, buildingName string) ([]models.Building, error) {
	ctx := context.Background()
	query := `
	SELECT * 
	FROM 
		building_instance 
	WHERE 
		user_id = $1 
		AND building_name = $2
	`

	rows, err := r.DB.Query(ctx, query, userID, buildingName)
	if err != nil {
		return []models.Building{}, err
	}
	defer rows.Close()

	buildings, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Building])
	if err != nil {
		return []models.Building{}, err
	}

	return buildings, nil
}

func (r *VillageRepository) UpgradeTownHall(userID string, upgradeCost int, upgradeCostType string, maxHP int) error {
	ctx := context.Background()

	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `
	UPDATE building_instance
	SET 
		hp = $1, level = level + 1
	WHERE 
		user_id = $2 AND building_name = 'Town Hall'
	`
	if _, err = tx.Exec(ctx, query, maxHP, userID); err != nil {
		return err
	}

	query = `
		UPDATE village 
		SET gold = gold - $1, town_hall_level = town_hall_level + 1
		WHERE user_id = $2
	`
	if _, err = tx.Exec(ctx, query, upgradeCost, userID); err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (r *VillageRepository) GetRandomVillage(attackerID string, attackerTHLevel int, attackerTrophies int) (models.Village, error) {
	ctx := context.Background()

	query := `
	SELECT v.id, v.user_id, v.town_hall_level, v.gold, v.elixir, v.gold_last_collected_at, v.elixir_last_collected_at
	FROM village v
	JOIN users u ON v.user_id = u.id
	WHERE v.user_id != $1
	ORDER BY ABS(v.town_hall_level - $2) ASC, ABS(u.trophies - $3) ASC, RANDOM()
	LIMIT 1
	`

	rows, err := r.DB.Query(ctx, query, attackerID, attackerTHLevel, attackerTrophies)
	if err != nil {
		return models.Village{}, err
	}
	defer rows.Close()

	village, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.Village])
	if err != nil {
		return models.Village{}, err
	}
	return village, nil
}
