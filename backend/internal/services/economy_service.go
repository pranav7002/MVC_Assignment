package services

import "time"

type EconomyService struct {
	VillageRepo VillageRepositoryInterface
	UserRepo    UserRepositoryInterface
	ConfigRepo  ConfigRepositoryInterface
}

func (s *EconomyService) CollectGold(userID string, reqTime time.Time) error {
	village, err := s.VillageRepo.GetVillage(userID) 
	if err != nil {
		return ErrServer
	}
	resourceConfig, err := s.ConfigRepo.GetAllResourceConfig()
	if err != nil {
		return ErrServer
	}
	goldMines, err := s.VillageRepo.GetUserBuildingsByName(userID, "gold_mine") 
	if err != nil {
		return ErrServer
	}

	goldPerSecPerLevel := make(map[int]int)
	for _, cfg := range resourceConfig {
		if cfg.Name == "elixir_collector" {
			continue
		}
		goldPerSecPerLevel[cfg.Level] = cfg.ResourcePerSec
	}

	var goldPerSec int
	for _, goldMine := range goldMines {
		goldPerSec = goldPerSec + goldPerSecPerLevel[goldMine.Level]
	}

	timeElapsed := reqTime.Sub(village.GoldLastCollectedAt)
	secondsElapsed := timeElapsed.Seconds()
	goldCollected := int(secondsElapsed) * goldPerSec

	if err := s.VillageRepo.AddResourceFromColletor(userID, "gold", goldCollected, reqTime); err != nil {
		return err
	}

	return nil
}

func (s *EconomyService) CollectElixir(userID string, reqTime time.Time) error {
	village, err := s.VillageRepo.GetVillage(userID) 
	if err != nil {
		return ErrServer
	}
	resourceConfig, err := s.ConfigRepo.GetAllResourceConfig()
	if err != nil {
		return ErrServer
	}
	elixirCollectors, err := s.VillageRepo.GetUserBuildingsByName(userID, "elixir_collector") 
	if err != nil {
		return ErrServer
	}

	elixirPerSecPerLevel := make(map[int]int)
	for _, cfg := range resourceConfig {
		if cfg.Name == "gold_mine" {
			continue
		}
		elixirPerSecPerLevel[cfg.Level] = cfg.ResourcePerSec
	}

	var elixirPerSec int
	for _, elixirCollector := range elixirCollectors {
		elixirPerSec = elixirPerSec + elixirPerSecPerLevel[elixirCollector.Level]
	}

	timeElapsed := reqTime.Sub(village.GoldLastCollectedAt)
	secondsElapsed := timeElapsed.Seconds()
	elixirCollected := int(secondsElapsed) * elixirPerSec

	if err := s.VillageRepo.AddResourceFromColletor(userID, "elixir", elixirCollected, reqTime); err != nil {
		return err
	}

	return nil
}