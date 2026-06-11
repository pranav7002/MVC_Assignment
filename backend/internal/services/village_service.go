package services

import (
	"errors"
	"fmt"

	"github.com/pranav7002/MVC_Assignment/internal/models"
)

type VillageRepositoryInterface interface {
	FetchUserBuildings(userID string) ([]models.Building, error)
	InsertBuilding(userID string, buildingReqBody models.BuildingCreationRequestBody, hp int, size int) error
	MoveBuilding(userID string, buildingID int64, posX, posY int) error
	GetGameProgressionConfig(thLevel int, buildingType string, buildingName string) (models.GameProgressionConfig, error)
	GetVillage(userID string) (models.Village, error)
	GetBuildingCount(userID string, buildingType string, buildingName string) (int, error)
	GetTownHallConfig(name string, level int) (models.TownHallConfig, error)
	GetDefenceConfig(name string, level int) (models.DefenceConfig, error)
	GetResourceConfig(name string, level int) (models.ResourceConfig, error)
	GetStorageConfig(name string, level int) (models.StorageConfig, error)
	GetTrainingGroundsConfig(name string, level int) (models.TrainingGroundsConfig, error)
}

type VillageService struct {
	VillageRepo VillageRepositoryInterface
}

func (villageService *VillageService) GetBuildings(userID string) ([]models.Building, error) {
	buildings, err := villageService.VillageRepo.FetchUserBuildings(userID)
	if err != nil {
		return nil, err
	}
	return buildings, nil
}

func (villageService *VillageService) CreateBuilding(userID string, buildingReqBody models.BuildingCreationRequestBody) error {
	village, err := villageService.VillageRepo.GetVillage(userID)
	if err != nil { 
		return err 
	}

	buildingCount, err := villageService.VillageRepo.GetBuildingCount(
		userID, 
		buildingReqBody.BuildingType, 
		buildingReqBody.BuildingName,
	)
	if err != nil { 
		return err 
	}

	gameProgConfig, err := villageService.VillageRepo.GetGameProgressionConfig(
		village.TownHallLevel, 
		buildingReqBody.BuildingType, 
		buildingReqBody.BuildingName,
	)	
	if err != nil { 
		return err 
	}

	// CHECK
	if gameProgConfig.MaxBuilt == 0 {
    	return errors.New("building not available at this town hall level")
	}
	if buildingCount == int(gameProgConfig.MaxBuilt) {
    	return errors.New("more buildings not allowed at this town hall level")
	}

	var upgradeCost int
	var upgradeCostType string
	var size int
	var maxHP int

	switch buildingReqBody.BuildingType {
	case "town_hall":
		config, err := villageService.VillageRepo.GetTownHallConfig(buildingReqBody.BuildingName, 1)
		if err != nil {
			return err
		}
		upgradeCost = config.UpgradeCost
		upgradeCostType = config.UpgradeCostType
		size = config.Size
		maxHP = config.MaxHP
	case "storage":
		config, err := villageService.VillageRepo.GetStorageConfig(buildingReqBody.BuildingName, 1)
		if err != nil {
			return err
		}
		upgradeCost = config.UpgradeCost
		upgradeCostType = config.UpgradeCostType
		size = config.Size
		maxHP = config.MaxHP
	case "resource":
		config, err := villageService.VillageRepo.GetResourceConfig(buildingReqBody.BuildingName, 1)
		if err != nil {
			return err
		}
		upgradeCost = config.UpgradeCost
		upgradeCostType = config.UpgradeCostType
		size = config.Size
		maxHP = config.MaxHP
	case "defence":
		config, err := villageService.VillageRepo.GetDefenceConfig(buildingReqBody.BuildingName, 1)
		if err != nil {
			return err
		}
		upgradeCost = config.UpgradeCost
		upgradeCostType = config.UpgradeCostType
		size = config.Size
		maxHP = config.MaxHP
	case "barracks":
		config, err := villageService.VillageRepo.GetTrainingGroundsConfig(buildingReqBody.BuildingName, 1)
		if err != nil {
			return err
		}
		upgradeCost = config.UpgradeCost
		upgradeCostType = config.UpgradeCostType
		size = config.Size
		maxHP = config.MaxHP
	default:
		return fmt.Errorf("invalid building type: %s", buildingReqBody.BuildingType)
	}

	
	switch upgradeCostType {
	case "gold":
		if village.Gold < upgradeCost {
			return errors.New("Insufficient Gold!!")
		}
	case "elixir":
		if village.Elixir < upgradeCost {
			return errors.New("Insufficient Elixir!!")
		}
	}

	buildings, err := villageService.VillageRepo.FetchUserBuildings(userID)
	if err != nil { 
		return err 
	}

	var villageBitmap [44][44]bool  
	for _, building := range buildings {
		for i := building.PosX; i <= building.Size; i++ {
			for j := building.PosY; j <= building.Size; j++ {
				villageBitmap[i][j] = true;
			}
		}
	}

	for i := buildingReqBody.PosX; i <= size; i++ {
		for j := buildingReqBody.PosY; j <= size; j++ {
			if villageBitmap[i][j] == true {
				return errors.New("Collision Detected!!")
			}
		}
	}

 	if err := villageService.VillageRepo.InsertBuilding(userID, buildingReqBody, maxHP, size); err != nil {
		return err
	}

	return nil
}

func (villageService *VillageService) MoveBuilding(userID string, buildingID int64, reqBody models.BuildingPositionRequestBody) error {
	buildings, err := villageService.VillageRepo.FetchUserBuildings(userID)
	if err != nil { 
		return err 
	}

	var villageBitmap [44][44]bool  
	for _, building := range buildings {
		for i := building.PosX; i <= building.PosX + building.Size; i++ {
			for j := building.PosY; j <= building.PosY + building.Size; j++ {
				villageBitmap[i][j] = true;
			}
		}
	}

	for i := reqBody.PosX; i <= reqBody.PosX + reqBody.Size; i++ {
		for j := reqBody.PosY; j <= reqBody.PosY + reqBody.Size; j++ {
			if villageBitmap[i][j] == true {
				return errors.New("Collision Detected!!")
			}
		}
	}

	if err := villageService.VillageRepo.MoveBuilding(userID, buildingID, reqBody.PosX, reqBody.PosY); err != nil {
		return err
	} 

	return nil
}