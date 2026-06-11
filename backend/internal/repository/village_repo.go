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

func (villageRepo *VillageRepository) GetGameProgressionConfig(thLevel int, buildingType string, buildingName string) (models.GameProgressionConfig, error) {
	ctx := context.Background()

	query := `SELECT * FROM game_progression_config WHERE building_type = $1 AND building_name = $2 AND town_hall_level = $3`

	rows, err := villageRepo.DB.Query(ctx, query, buildingType, buildingName, thLevel)
	if err != nil {
		return models.GameProgressionConfig{}, err
	}
	defer rows.Close()

	gameProgConfig, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.GameProgressionConfig])
	if err != nil {
		return models.GameProgressionConfig{}, err
	}
	return gameProgConfig, nil
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

func (villageRepo *VillageRepository) GetTownHallConfig(name string, level int) (models.TownHallConfig, error) {
	ctx := context.Background()

	query := `SELECT * FROM town_hall_config WHERE name = $1 AND level = $2`
	rows, err := villageRepo.DB.Query(ctx, query, name, level)
	if err != nil {
		return models.TownHallConfig{}, err
	}
	defer rows.Close()

	config, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.TownHallConfig])
	if err != nil {
		return models.TownHallConfig{}, err
	}

	return config, nil
}

func (villageRepo *VillageRepository) GetDefenceConfig(name string, level int) (models.DefenceConfig, error) {
	ctx := context.Background()

	query := `SELECT * FROM defence_config WHERE name = $1 AND level = $2`
	rows, err := villageRepo.DB.Query(ctx, query, name, level)
	if err != nil {
		return models.DefenceConfig{}, err
	}
	defer rows.Close()

	config, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.DefenceConfig])
	if err != nil {
		return models.DefenceConfig{}, err
	}

	return config, nil
}

func (villageRepo *VillageRepository) GetResourceConfig(name string, level int) (models.ResourceConfig, error) {
	ctx := context.Background()

	query := `SELECT * FROM resource_config WHERE name = $1 AND level = $2`
	rows, err := villageRepo.DB.Query(ctx, query, name, level)
	if err != nil {
		return models.ResourceConfig{}, err
	}
	defer rows.Close()

	config, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.ResourceConfig])
	if err != nil {
		return models.ResourceConfig{}, err
	}

	return config, nil
}

func (villageRepo *VillageRepository) GetStorageConfig(name string, level int) (models.StorageConfig, error) {
	ctx := context.Background()

	query := `SELECT * FROM storage_config WHERE name = $1 AND level = $2`
	rows, err := villageRepo.DB.Query(ctx, query, name, level)
	if err != nil {
		return models.StorageConfig{}, err
	}
	defer rows.Close()

	config, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.StorageConfig])
	if err != nil {
		return models.StorageConfig{}, err
	}

	return config, nil
}

func (villageRepo *VillageRepository) GetTrainingGroundsConfig(name string, level int) (models.TrainingGroundsConfig, error) {
	ctx := context.Background()

	query := `SELECT * FROM training_grounds_config WHERE name = $1 AND level = $2`
	rows, err := villageRepo.DB.Query(ctx, query, name, level)
	if err != nil {
		return models.TrainingGroundsConfig{}, err
	}
	defer rows.Close()

	config, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.TrainingGroundsConfig])
	if err != nil {
		return models.TrainingGroundsConfig{}, err
	}

	return config, nil
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