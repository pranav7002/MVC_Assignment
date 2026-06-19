package services

import (
	"context"

	"github.com/pranav7002/MVC_Assignment/internal/models"
	"github.com/pranav7002/MVC_Assignment/internal/repository"
)

type ShopService struct {
	ShopRepo    *repository.ShopRepository
	VillageRepo *repository.VillageRepository
}

func (s *ShopService) GetShopBuildings(ctx context.Context, userID string) ([]models.ShopBuilding, error) {
	village, err := s.VillageRepo.GetVillage(userID)
	if err != nil {
		return nil, err
	}
	return s.ShopRepo.GetShopBuildings(ctx, village.TownHallLevel)
}

func (s *ShopService) GetShopTroops(ctx context.Context) ([]models.ShopTroop, error) {
	return s.ShopRepo.GetShopTroops(ctx)
}
