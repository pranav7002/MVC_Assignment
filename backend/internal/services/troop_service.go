package services

import "github.com/pranav7002/MVC_Assignment/internal/models"

type TroopRepositoryInterface interface {
	GetUserTrainedTroops(userID string) ([]models.TroopTrained, error)
	TrainTroop(userID string, troopName string, quantity int, trainingCost int) error
}

type TroopService struct {
	TroopRepo   TroopRepositoryInterface
	VillageRepo VillageRepositoryInterface
	ConfigRepo  ConfigRepositoryInterface
}

func (s *TroopService) TrainTroop(userID string, troopName string, quantity int) error {
	village, err := s.VillageRepo.GetVillage(userID)
	if err != nil {
		return ErrServer 
	}
	troopConfig, err := s.ConfigRepo.GetAllTroopConfig()
	if err != nil {
		return ErrServer
	}

	TroopHousingSpace := make(map[string]int)
	var troop models.TroopConfig
	var found bool

	for _,t := range troopConfig {
		TroopHousingSpace[t.Name] = t.HousingSpace

		if t.Name == troopName {
			troop = t
			found = true
		}
	}
	if !found {
		return ErrInvalidTroop
	}

	trainingCost := troop.TrainingCost * quantity
	housingSpaceRequired := troop.HousingSpace * quantity

	if village.Elixir < trainingCost {
		return ErrInsufficientElixir 
	}

	trainingGrounds, err := s.VillageRepo.GetUserBuildingsByName(userID, "training_grounds")
	if err != nil {
		return ErrServer
	}
	trainingGroundsConfig, err := s.ConfigRepo.GetAllTrainingGroundsConfig()
	if err != nil {
		return ErrServer
	}
	
	troopTrained, err := s.TroopRepo.GetUserTrainedTroops(userID)
	if err != nil {
		return ErrServer
	}

	var housingSpaceUsed int 
	for _, t := range troopTrained {
		housingSpaceUsed = housingSpaceUsed + t.Quantity * TroopHousingSpace[t.TroopName]
	}

	housingSpacePerTrainingGroundPerLevel := make(map[int]int)
	for _, cfg := range trainingGroundsConfig {
		housingSpacePerTrainingGroundPerLevel[cfg.Level] = cfg.HousingSpace
	}

	var housingSpace int 
	for _, trainingGround := range trainingGrounds {
		housingSpace = housingSpace + housingSpacePerTrainingGroundPerLevel[trainingGround.Level]
	}

	if housingSpaceUsed + housingSpaceRequired > housingSpace {
		return ErrInsufficientHousingSpace
	}

	if err := s.TroopRepo.TrainTroop(userID, troopName, quantity, trainingCost); err != nil {
		return ErrServer
	}

	return nil
}

func (s *TroopService) GetTrainedTroops(userID string) ([]models.TroopTrained, error) {
	troopTrained, err := s.TroopRepo.GetUserTrainedTroops(userID)
	if err != nil {
		return []models.TroopTrained{}, ErrServer
	}

	return troopTrained, ErrServer
}