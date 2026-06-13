package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pranav7002/MVC_Assignment/internal/models"
)

type ConfigRepository struct {
	DB *pgxpool.Pool
}

func (configRepo *ConfigRepository) GetGameProgressionConfig(thLevel int, buildingType string, buildingName string) (models.GameProgressionConfig, error) {
	ctx := context.Background()

	query := `SELECT * FROM game_progression_config WHERE building_type = $1 AND building_name = $2 AND town_hall_level = $3`

	rows, err := configRepo.DB.Query(ctx, query, buildingType, buildingName, thLevel)
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

func (configRepo *ConfigRepository) GetTownHallConfig(name string, level int) (models.TownHallConfig, error) {
	ctx := context.Background()

	query := `SELECT * FROM town_hall_config WHERE name = $1 AND level = $2`
	rows, err := configRepo.DB.Query(ctx, query, name, level)
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

func (configRepo *ConfigRepository) GetDefenseConfig(name string, level int) (models.DefenseConfig, error) {
	ctx := context.Background()

	query := `SELECT * FROM defense_config WHERE name = $1 AND level = $2`
	rows, err := configRepo.DB.Query(ctx, query, name, level)
	if err != nil {
		return models.DefenseConfig{}, err
	}
	defer rows.Close()

	config, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[models.DefenseConfig])
	if err != nil {
		return models.DefenseConfig{}, err
	}

	return config, nil
}

func (configRepo *ConfigRepository) GetResourceConfig(name string, level int) (models.ResourceConfig, error) {
	ctx := context.Background()

	query := `SELECT * FROM resource_config WHERE name = $1 AND level = $2`
	rows, err := configRepo.DB.Query(ctx, query, name, level)
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

func (configRepo *ConfigRepository) GetStorageConfig(name string, level int) (models.StorageConfig, error) {
	ctx := context.Background()

	query := `SELECT * FROM storage_config WHERE name = $1 AND level = $2`
	rows, err := configRepo.DB.Query(ctx, query, name, level)
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

func (configRepo *ConfigRepository) GetTrainingGroundsConfig(name string, level int) (models.TrainingGroundsConfig, error) {
	ctx := context.Background()

	query := `SELECT * FROM training_grounds_config WHERE name = $1 AND level = $2`
	rows, err := configRepo.DB.Query(ctx, query, name, level)
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

func (configRepo *ConfigRepository) GetAllResourceConfig() ([]models.ResourceConfig, error) {
	ctx := context.Background()

	query := `SELECT * FROM resource_config`
	rows, err := configRepo.DB.Query(ctx, query)
	if err != nil {
		return []models.ResourceConfig{}, err
	}
	defer rows.Close()

	config, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.ResourceConfig])
	if err != nil {
		return []models.ResourceConfig{}, err
	}

	return config, nil
}