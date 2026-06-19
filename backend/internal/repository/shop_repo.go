package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pranav7002/MVC_Assignment/internal/models"
)

type ShopRepository struct {
	DB *pgxpool.Pool
}

func (r *ShopRepository) GetShopBuildings(ctx context.Context, thLevel int) ([]models.ShopBuilding, error) {

	query := `
		SELECT 
			g.building_type, 
			g.building_name, 
			g.max_level, 
			g.max_built,
			COALESCE(d.upgrade_cost, res.upgrade_cost, s.upgrade_cost, t.upgrade_cost, th.upgrade_cost, 0) as cost,
			COALESCE(d.upgrade_cost_type, res.upgrade_cost_type, s.upgrade_cost_type, t.upgrade_cost_type, th.upgrade_cost_type, 'gold') as cost_type,
			COALESCE(d.size, res.size, s.size, t.size, th.size, 3) as size
		FROM game_progression_config g
		LEFT JOIN defense_config d ON g.building_name = d.name AND d.level = 1
		LEFT JOIN resource_config res ON g.building_name = res.name AND res.level = 1
		LEFT JOIN storage_config s ON g.building_name = s.name AND s.level = 1
		LEFT JOIN training_grounds_config t ON g.building_name = t.name AND t.level = 1
		LEFT JOIN town_hall_config th ON g.building_name = th.name AND th.level = 1
		WHERE g.town_hall_level = $1
	`
	rows, err := r.DB.Query(ctx, query, thLevel)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var buildings []models.ShopBuilding
	for rows.Next() {
		var b models.ShopBuilding
		err := rows.Scan(&b.BuildingType, &b.BuildingName, &b.MaxLevel, &b.MaxBuilt, &b.Cost, &b.CostType, &b.Size)
		if err != nil {
			return nil, err
		}
		buildings = append(buildings, b)
	}
	return buildings, nil
}

func (r *ShopRepository) GetShopTroops(ctx context.Context) ([]models.ShopTroop, error) {
	query := `SELECT name, dps, health, range, housing_space, training_cost FROM troop_config`
	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var troops []models.ShopTroop
	for rows.Next() {
		var t models.ShopTroop
		err := rows.Scan(&t.Name, &t.DPS, &t.Health, &t.Range, &t.HousingSpace, &t.TrainingCost)
		if err != nil {
			return nil, err
		}
		troops = append(troops, t)
	}
	return troops, nil
}
