package services

import "github.com/pranav7002/MVC_Assignment/internal/models"

type VillageRepositoryInterface interface {
	FetchUserBuildings(userID string) ([]models.Building, error)
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