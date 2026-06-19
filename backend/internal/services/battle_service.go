package services

import (
	"errors"
	"log"

	"github.com/pranav7002/MVC_Assignment/internal/models"
	"github.com/pranav7002/MVC_Assignment/internal/simulation"
)

type BattleService struct {
	BattleRepo  BattleRepositoyInterface
	VillageRepo VillageRepositoryInterface
	ConfigRepo  ConfigRepositoryInterface
}

type BattleRepositoyInterface interface {
	StoreBattle(userID, defendersID, result string, stars, destructionPct, goldLooted, elixirLooted int) error
}

func (s *BattleService) HydrateTroop(t models.TroopDropBody, buildings []models.Building) (simulation.TroopDrop, error) {
	cfg, err := s.ConfigRepo.GetTroopConfig(t.Name)
	if err != nil {
		return simulation.TroopDrop{}, ErrServer
	}
	var villageBitmap [gridSize][gridSize]bool
	for _, building := range buildings {
		for i := building.PosX; i < building.PosX+building.Size; i++ {
			for j := building.PosY; j < building.PosY+building.Size; j++ {
				villageBitmap[i][j] = true
			}
		}
	}
	outOfBounds := t.X >= gridSize || t.Y >= gridSize || t.X < 0 || t.Y < 0
	if villageBitmap[t.X][t.Y] || outOfBounds {
		return simulation.TroopDrop{}, errors.New("Invalid drop location")
	}
	return simulation.TroopDrop{
		Name:  cfg.Name,
		Pos:   simulation.Position{X: int(t.X), Y: int(t.Y)},
		HP:    cfg.Health,
		DPS:   cfg.DPS,
		Range: cfg.Range,
	}, nil
}

func (s *BattleService) HydrateBuildings(b []models.Building) ([]simulation.BuildingInput, error) {
	type BuildingKey struct {
		Name  string
		Level int
	}

	var buildingInput []simulation.BuildingInput

	allDefenseCfg, err := s.ConfigRepo.GetAllDefenseConfig()
	if err != nil {
		return nil, ErrServer
	}

	cfg := make(map[BuildingKey]models.DefenseConfig)
	for _, c := range allDefenseCfg {
		cfg[BuildingKey{c.Name, c.Level}] = c
	}

	for _, building := range b {
		if building.BuildingType == "defense" {
			buildingInput = append(buildingInput, simulation.BuildingInput{
				ID:       int(building.ID),
				Name:     building.BuildingName,
				Type:     building.BuildingType,
				Pos:      simulation.Position{X: building.PosX, Y: building.PosY},
				Size:     building.Size,
				HP:       building.HP,
				AOERange: cfg[BuildingKey{building.BuildingName, building.Level}].AOERange,
				MinRange: cfg[BuildingKey{building.BuildingName, building.Level}].MinRange,
				MaxRange: cfg[BuildingKey{building.BuildingName, building.Level}].MaxRange,
				DPS:      cfg[BuildingKey{building.BuildingName, building.Level}].DPS,
			})
			continue
		}
		buildingInput = append(buildingInput, simulation.BuildingInput{
			ID:   int(building.ID),
			Name: building.BuildingName,
			Type: building.BuildingType,
			Pos:  simulation.Position{X: building.PosX, Y: building.PosY},
			Size: building.Size,
			HP:   building.HP,
		})
	}

	return buildingInput, nil
}

func (s *BattleService) SaveBattleResult(userID, defendersID string, stars, destructionPct int) error {
	defendersVillage, err := s.VillageRepo.GetVillage(defendersID)
	if err != nil {
		log.Println("error:", err)
		return ErrServer
	}

	goldLooted := defendersVillage.Gold * destructionPct / 2
	elixirLooted := defendersVillage.Elixir * destructionPct / 2

	var result string
	switch stars {
	case 0:
		result = "LOSS"
	case 1:
		result = "ONE_STAR"
	case 2:
		result = "TWO_STARS"
	case 3:
		result = "THREE_STARS"
	}

	if err := s.BattleRepo.StoreBattle(
		userID,
		defendersID,
		result,
		stars,
		destructionPct,
		goldLooted,
		elixirLooted,
	); err != nil {
		log.Println("error:", err)
		return ErrServer
	}

	return nil
}
