package services

type TroopRepositoryInterface interface {
}

type TroopService struct {
	TroopRepo   TroopRepositoryInterface
	VillageRepo VillageRepositoryInterface
	ConfigRepo  ConfigRepositoryInterface
}

func (s *TroopService) TrainTroop(userID string, troopName string, quantity int) error {
	return nil
}