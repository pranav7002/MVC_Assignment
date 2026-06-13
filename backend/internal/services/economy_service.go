package services

import "time"

type EconomyService struct {
	VillageRepo VillageRepositoryInterface
	UserRepo    UserRepositoryInterface
	ConfigRepo  ConfigRepositoryInterface
}

func (economyService *EconomyService) CollectGold(userID string, reqTime time.Time) error {
	village, err := economyService.VillageRepo.GetVillage(userID) 
	if err != nil {
		return ErrServer
	}
	resourceConfig, err := economyService.ConfigRepo.GetAllResourceConfig()
	if err != nil {
		return ErrServer
	}
	goldMines, err := economyService.VillageRepo.GetUserBuildingsByName(userID, "gold_mine") 
	if err != nil {
		return ErrServer
	}

	goldPerSecPerLevel := make(map[int]int)
	for _, resourceBuilding := range resourceConfig {
		if resourceBuilding.Name == "elixir_collector" {
			continue
		}
		goldPerSecPerLevel[resourceBuilding.Level] = resourceBuilding.ResourcePerSec
	}

	var goldPerSec int
	for _, goldMine := range goldMines {
		goldPerSec = goldPerSec + goldPerSecPerLevel[goldMine.Level]
	}

	timeElapsed := reqTime.Sub(village.GoldLastCollectedAt)
	secondsElapsed := timeElapsed.Seconds()
	goldCollected := int(secondsElapsed) * goldPerSec

	if err := economyService.VillageRepo.AddResourceFromColletor(userID, "gold", goldCollected, reqTime); err != nil {
		return err
	}

	return nil
}

func (economyService *EconomyService) CollectElixir(userID string, reqTime time.Time) error {
	village, err := economyService.VillageRepo.GetVillage(userID) 
	if err != nil {
		return ErrServer
	}
	resourceConfig, err := economyService.ConfigRepo.GetAllResourceConfig()
	if err != nil {
		return ErrServer
	}
	elixirCollectors, err := economyService.VillageRepo.GetUserBuildingsByName(userID, "elixir_collector") 
	if err != nil {
		return ErrServer
	}

	elixirPerSecPerLevel := make(map[int]int)
	for _, resourceBuilding := range resourceConfig {
		if resourceBuilding.Name == "gold_mine" {
			continue
		}
		elixirPerSecPerLevel[resourceBuilding.Level] = resourceBuilding.ResourcePerSec
	}

	var elixirPerSec int
	for _, elixirCollector := range elixirCollectors {
		elixirPerSec = elixirPerSec + elixirPerSecPerLevel[elixirCollector.Level]
	}

	timeElapsed := reqTime.Sub(village.GoldLastCollectedAt)
	secondsElapsed := timeElapsed.Seconds()
	elixirCollected := int(secondsElapsed) * elixirPerSec

	if err := economyService.VillageRepo.AddResourceFromColletor(userID, "elixir", elixirCollected, reqTime); err != nil {
		return err
	}

	return nil
}