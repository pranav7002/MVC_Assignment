package models

type ShopBuilding struct {
	BuildingType string `json:"building_type"`
	BuildingName string `json:"building_name"`
	MaxLevel     int    `json:"max_level"`
	MaxBuilt     int    `json:"max_built"`
	Cost         int    `json:"cost"`
	CostType     string `json:"cost_type"`
	Size         int    `json:"size"`
}

type ShopTroop struct {
	Name         string `json:"name"`
	DPS          int    `json:"dps"`
	Health       int    `json:"health"`
	Range        int    `json:"range"`
	HousingSpace int    `json:"housing_space"`
	TrainingCost int    `json:"training_cost"`
}
