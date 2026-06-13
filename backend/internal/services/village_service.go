package services

import (
	"time"

	"github.com/pranav7002/MVC_Assignment/internal/models"
)

type VillageRepositoryInterface interface {
	GetUserBuildings(userID string) ([]models.Building, error)
	GetBuilding(buildingID int64) (models.Building, error)
	GetUserBuildingsByName(userID string, buildingName string) ([]models.Building, error)
	InsertBuilding(userID string, buildingReqBody models.BuildingCreationRequestBody, hp int, size int) error
	UpdateBuilding(userID string, buildingID int64, hp int) error
	MoveBuilding(userID string, buildingID int64, posX, posY int) error
	GetVillage(userID string) (models.Village, error)
	GetBuildingCount(userID string, buildingType string, buildingName string) (int, error)
	RemoveResource(userID string, resourceType string, amount int) error
	AddResourceFromColletor(userID string, resourceType string, amount int, collectionTime time.Time) error
}

type ConfigRepositoryInterface interface {
	GetGameProgressionConfig(thLevel int, buildingType string, buildingName string) (models.GameProgressionConfig, error)
	GetTownHallConfig(name string, level int) (models.TownHallConfig, error)
	GetDefenseConfig(name string, level int) (models.DefenseConfig, error)
	GetResourceConfig(name string, level int) (models.ResourceConfig, error)
	GetStorageConfig(name string, level int) (models.StorageConfig, error)
	GetTrainingGroundsConfig(name string, level int) (models.TrainingGroundsConfig, error)
	GetTroopConfig(name string) (models.TroopConfig, error)
	GetAllResourceConfig() ([]models.ResourceConfig, error)
	GetAllTrainingGroundsConfig() ([]models.TrainingGroundsConfig, error)
	GetAllTroopConfig() ([]models.TroopConfig, error)
}

type VillageService struct {
	VillageRepo VillageRepositoryInterface
	ConfigRepo  ConfigRepositoryInterface
}

func (s *VillageService) GetBuildings(userID string) ([]models.Building, error) {
	buildings, err := s.VillageRepo.GetUserBuildings(userID)
	if err != nil {
		return nil, ErrServer
	}
	return buildings, nil
}

func (s *VillageService) CreateBuilding(userID string, buildingReqBody models.BuildingCreationRequestBody) error {
	village, err := s.VillageRepo.GetVillage(userID)
	if err != nil {
		return ErrServer
	}

	buildingCount, err := s.VillageRepo.GetBuildingCount(
		userID,
		buildingReqBody.BuildingType,
		buildingReqBody.BuildingName,
	)
	if err != nil {
		return ErrServer
	}

	gameProgConfig, err := s.ConfigRepo.GetGameProgressionConfig(
		village.TownHallLevel,
		buildingReqBody.BuildingType,
		buildingReqBody.BuildingName,
	)
	if err != nil {
		return ErrServer
	}

	// CHECK
	if gameProgConfig.MaxBuilt == 0 {
		return ErrBuildingNotUnlocked
	}
	if buildingCount == int(gameProgConfig.MaxBuilt) {
		return ErrBuildingLimitReached 
	}

	var upgradeCost int
	var upgradeCostType string
	var size int
	var maxHP int

	switch buildingReqBody.BuildingType {
	case "town_hall":
		config, err := s.ConfigRepo.GetTownHallConfig(buildingReqBody.BuildingName, 1)
		if err != nil {
			return ErrServer
		}
		upgradeCost = config.UpgradeCost
		upgradeCostType = config.UpgradeCostType
		size = config.Size
		maxHP = config.MaxHP
	case "storage":
		config, err := s.ConfigRepo.GetStorageConfig(buildingReqBody.BuildingName, 1)
		if err != nil {
			return ErrServer
		}
		upgradeCost = config.UpgradeCost
		upgradeCostType = config.UpgradeCostType
		size = config.Size
		maxHP = config.MaxHP
	case "resource":
		config, err := s.ConfigRepo.GetResourceConfig(buildingReqBody.BuildingName, 1)
		if err != nil {
			return ErrServer
		}
		upgradeCost = config.UpgradeCost
		upgradeCostType = config.UpgradeCostType
		size = config.Size
		maxHP = config.MaxHP
	case "defense":
		config, err := s.ConfigRepo.GetDefenseConfig(buildingReqBody.BuildingName, 1)
		if err != nil {
			return ErrServer
		}
		upgradeCost = config.UpgradeCost
		upgradeCostType = config.UpgradeCostType
		size = config.Size
		maxHP = config.MaxHP
	case "training_grounds":
		config, err := s.ConfigRepo.GetTrainingGroundsConfig(buildingReqBody.BuildingName, 1)
		if err != nil {
			return ErrServer
		}
		upgradeCost = config.UpgradeCost
		upgradeCostType = config.UpgradeCostType
		size = config.Size
		maxHP = config.MaxHP
	default:
		return ErrInvalidBuildingType
	}

	switch upgradeCostType {
	case "gold":
		if village.Gold < upgradeCost {
			return ErrInsufficientResources
		}
	case "elixir":
		if village.Elixir < upgradeCost {
			return ErrInsufficientResources
		}
	}

	buildings, err := s.VillageRepo.GetUserBuildings(userID)
	if err != nil {
		return ErrServer
	}

	var villageBitmap [44][44]bool
	for _, building := range buildings {
		for i := building.PosX; i < building.PosX + building.Size; i++ {
			for j := building.PosY; j < building.PosY + building.Size; j++ {
				villageBitmap[i][j] = true
			}
		}
	}

	for i := buildingReqBody.PosX; i < buildingReqBody.PosX + size; i++ {
		for j := buildingReqBody.PosY; j < buildingReqBody.PosY + size; j++ {
			if villageBitmap[i][j] {
				return ErrCollisionDetected
			}
		}
	}

	if err := s.VillageRepo.RemoveResource(userID, upgradeCostType, upgradeCost); err != nil {
		return ErrServer
	}

	if err := s.VillageRepo.InsertBuilding(userID, buildingReqBody, maxHP, size); err != nil {
		return ErrServer
	}

	return nil
}

func (s *VillageService) MoveBuilding(userID string, buildingID int64, reqBody models.BuildingPositionRequestBody) error {
	buildings, err := s.VillageRepo.GetUserBuildings(userID)
	if err != nil {
		return ErrServer
	}

	var b models.Building
	var villageBitmap [44][44]bool
	for _, building := range buildings {
		for i := building.PosX; i < building.PosX+building.Size; i++ {
			for j := building.PosY; j < building.PosY+building.Size; j++ {
				if building.ID == buildingID { 
					b = building
					continue 
				}
				villageBitmap[i][j] = true
			}
		}
	}

	for i := reqBody.PosX; i < reqBody.PosX + b.Size; i++ {
		for j := reqBody.PosY; j < reqBody.PosY + b.Size; j++ {
			if i > 43 || j > 43 {
				return ErrOutOfBounds
			}
			if villageBitmap[i][j] {
				return ErrCollisionDetected
			}
		}
	}

	if err := s.VillageRepo.MoveBuilding(userID, buildingID, reqBody.PosX, reqBody.PosY); err != nil {
		return ErrServer
	}

	return nil
}

func (s *VillageService) UpgradeBuilding(userID string, buildingID int64) error {
	village, err := s.VillageRepo.GetVillage(userID)
	if err != nil {
		return ErrServer
	}
	building, err := s.VillageRepo.GetBuilding(buildingID)
	if err != nil {
		return ErrServer
	}
	gameProgConfig, err := s.ConfigRepo.GetGameProgressionConfig(village.TownHallLevel, building.BuildingType, building.BuildingName)
	if err != nil {
		return ErrServer
	}

	if building.Level >= int(gameProgConfig.MaxLevel) {
		return ErrHighestLevelReached
	}

	var upgradeCost int
	var upgradeCostType string
	var maxHP int

	switch building.BuildingType {
	case "town_hall":
		config, err := s.ConfigRepo.GetTownHallConfig(building.BuildingName, building.Level + 1)
		if err != nil {
			return ErrServer
		}
		upgradeCost = config.UpgradeCost
		upgradeCostType = config.UpgradeCostType
		maxHP = config.MaxHP
	case "storage":
		config, err := s.ConfigRepo.GetStorageConfig(building.BuildingName, building.Level + 1)
		if err != nil {
			return ErrServer
		}
		upgradeCost = config.UpgradeCost
		upgradeCostType = config.UpgradeCostType
		maxHP = config.MaxHP
	case "resource":
		config, err := s.ConfigRepo.GetResourceConfig(building.BuildingName, building.Level + 1)
		if err != nil {
			return ErrServer
		}
		upgradeCost = config.UpgradeCost
		upgradeCostType = config.UpgradeCostType
		maxHP = config.MaxHP
	case "defense":
		config, err := s.ConfigRepo.GetDefenseConfig(building.BuildingName, building.Level + 1)
		if err != nil {
			return ErrServer
		}
		upgradeCost = config.UpgradeCost
		upgradeCostType = config.UpgradeCostType
		maxHP = config.MaxHP
	case "training_grounds":
		config, err := s.ConfigRepo.GetTrainingGroundsConfig(building.BuildingName, building.Level + 1)
		if err != nil {
			return ErrServer
		}
		upgradeCost = config.UpgradeCost
		upgradeCostType = config.UpgradeCostType
		maxHP = config.MaxHP
	default:
		return ErrInvalidBuildingType
	}

	switch upgradeCostType {
	case "gold":
		if village.Gold < upgradeCost {
			return ErrInsufficientResources
		}
	case "elixir":
		if village.Elixir < upgradeCost {
			return ErrInsufficientResources
		}
	}

	if err := s.VillageRepo.UpdateBuilding(userID, buildingID, maxHP); err != nil {
		return ErrServer
	}

	if err := s.VillageRepo.RemoveResource(userID, upgradeCostType, upgradeCost); err != nil {
		return ErrServer
	}

	return nil
}

