package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pranav7002/MVC_Assignment/internal/models"
)

type VillageRepository struct {
	DB *pgxpool.Pool
}

func (villageRepo *VillageRepository) GetUserBuildings(userID string) ([]models.Building, error) {
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


func (villageRepo *VillageRepository) GetVillage(userID string) (models.Village, error) {
	ctx := context.Background()

	query := `SELECT * FROM village WHERE user_id = $1`
	rows, err := villageRepo.DB.Query(ctx, query, userID)
	if err != nil {
		return models.Village{}, err
	}
	defer rows.Close()

	village, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.Village])
	return village, nil
}

func (villageRepo *VillageRepository) GetBuildingCount(userID string, buildingType string, buildingName string) (int, error) {
	ctx := context.Background()

	var count int
	query := `SELECT COUNT(*) FROM building_instance WHERE user_id = $1 AND building_type = $2 AND building_name = $3`
	err := villageRepo.DB.QueryRow(ctx, query, userID, buildingType, buildingName).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}


func (villageRepo *VillageRepository) InsertBuilding(userID string, buildingReqBody models.BuildingCreationRequestBody, hp int, size int) error {
	ctx := context.Background()

	query := `
		INSERT INTO building_instance 
		(user_id, building_type, building_name, level, pos_x, pos_y, size, is_upgrading, hp) 
		VALUES ($1, $2, $3, 1, $4, $5, $6, false, $7)
	`

	_, err := villageRepo.DB.Exec(ctx, query,
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

func (villageRepo *VillageRepository) MoveBuilding(userID string, buildingID int64, posX, posY int) error {
	ctx := context.Background()

	query := `
		UPDATE building_instance 
		SET pos_x = $1, pos_y = $2 
		WHERE id = $3 AND user_id = $4
	`

	result, err := villageRepo.DB.Exec(ctx, query, posX, posY, buildingID, userID)
	if err != nil {
		return err
	}

	// Check if the building belongs to the user
	if result.RowsAffected() == 0 {
		return fmt.Errorf("building not found or does not belong to user")
	}

	return nil
}

func (villageRepo *VillageRepository) RemoveResource(userID string, resourceType string, amount int) error {
	ctx := context.Background()

	query := fmt.Sprintf(`
		UPDATE village 
		SET %s = %s - $1 
		WHERE user_id = $2
	`, resourceType, resourceType)

	_, err := villageRepo.DB.Exec(ctx, query, amount, userID)
	if err != nil {
		return err
	}

	return nil
}

func (villageRepo *VillageRepository) GetBuilding(buildingID int64) (models.Building, error) {
	ctx := context.Background()

	query := `
	SELECT *
	FROM 
		building_instance
	WHERE 
		id == $1
	`

	rows, err := villageRepo.DB.Query(ctx, query, buildingID) 
	if err != nil {
		return models.Building{}, err
	}
	defer rows.Close()

	building, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.Building]);
	if err != nil {
		return models.Building{}, err
	}

	return building, nil
}

func (villageRepo *VillageRepository) UpdateBuilding(userID string, buildingID int64, hp int) error {
	ctx := context.Background()

	query := `
	UPDATE building_instance
	SET 
		hp = $1 
		AND level = level + 1
	WHERE 
		id == $2
		AND user_id == $3
	`

	result, err := villageRepo.DB.Exec(ctx, query, hp, buildingID, userID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("building not found or does not belong to user")
	}

	return nil
}