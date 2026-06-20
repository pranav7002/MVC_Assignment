package services

import (
	"log"
	"time"
)

type EconomyService struct {
	VillageRepo VillageRepositoryInterface
	UserRepo    UserRepositoryInterface
	ConfigRepo  ConfigRepositoryInterface
}

func (s *EconomyService) CollectGold(userID string, reqTime time.Time) error {
	village, err := s.VillageRepo.GetVillage(userID)
	if err != nil {
		log.Println("error:", err)
		return ErrServer
	}
	resourceConfig, err := s.ConfigRepo.GetAllResourceConfig()
	if err != nil {
		log.Println("error:", err)
		return ErrServer
	}
	goldMines, err := s.VillageRepo.GetUserBuildingsByName(userID, "Gold Mine")
	if err != nil {
		log.Println("error:", err)
		return ErrServer
	}

	goldPerSecPerLevel := make(map[int]float64)
	for _, cfg := range resourceConfig {
		if cfg.Name == "Elixir Collector" {
			continue
		}
		goldPerSecPerLevel[cfg.Level] = cfg.ResourcePerSec
	}

	var goldPerSec float64
	for _, goldMine := range goldMines {
		goldPerSec = goldPerSec + goldPerSecPerLevel[goldMine.Level]
	}

	timeElapsed := reqTime.Sub(village.GoldLastCollectedAt)
	secondsElapsed := timeElapsed.Seconds()
	goldCollected := int(secondsElapsed * goldPerSec)

	if err := s.VillageRepo.AddResourceFromColletor(userID, "gold", goldCollected, reqTime); err != nil {
		return err
	}

	return nil
}

func (s *EconomyService) CollectElixir(userID string, reqTime time.Time) error {
	village, err := s.VillageRepo.GetVillage(userID)
	if err != nil {
		log.Println("error:", err)
		return ErrServer
	}
	resourceConfig, err := s.ConfigRepo.GetAllResourceConfig()
	if err != nil {
		log.Println("error:", err)
		return ErrServer
	}
	elixirCollectors, err := s.VillageRepo.GetUserBuildingsByName(userID, "Elixir Collector")
	if err != nil {
		log.Println("error:", err)
		return ErrServer
	}

	elixirPerSecPerLevel := make(map[int]float64)
	for _, cfg := range resourceConfig {
		if cfg.Name == "Gold Mine" {
			continue
		}
		elixirPerSecPerLevel[cfg.Level] = cfg.ResourcePerSec
	}

	var elixirPerSec float64
	for _, elixirCollector := range elixirCollectors {
		elixirPerSec = elixirPerSec + elixirPerSecPerLevel[elixirCollector.Level]
	}

	timeElapsed := reqTime.Sub(village.ElixirLastCollectedAt)
	secondsElapsed := timeElapsed.Seconds()
	elixirCollected := int(secondsElapsed * elixirPerSec)

	if err := s.VillageRepo.AddResourceFromColletor(userID, "elixir", elixirCollected, reqTime); err != nil {
		return err
	}

	return nil
}
