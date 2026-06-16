package models

type GameProgressionConfig struct {
	ID            int64  `json:"id" db:"id"`
	TownHallLevel int    `json:"town_hall_level" db:"town_hall_level"`
	BuildingType  string `json:"building_type" db:"building_type"`
	BuildingName  string `json:"building_name" db:"building_name"`
	MaxLevel      int64  `json:"max_level" db:"max_level"`
	MaxBuilt      int64  `json:"max_built" db:"max_built"`
}

type StorageConfig struct {
	Name               string `json:"name" db:"name"`
	Level              int    `json:"level" db:"level"`
	ResourceType       string `json:"resource_type" db:"resource_type"`
	MaxCapacity        int    `json:"max_capacity" db:"max_capacity"`
	UpgradeCost        int    `json:"upgrade_cost" db:"upgrade_cost"`
	UpgradeCostType    string `json:"upgrade_cost_type" db:"upgrade_cost_type"`
	UpgradeDurationSec int    `json:"upgrade_duration_sec" db:"upgrade_duration_sec"`
	MaxHP              int    `json:"max_hp" db:"max_hp"`
	Size               int    `json:"size" db:"size"`
}

type ResourceConfig struct {
	Name               string `json:"name" db:"name"`
	Level              int    `json:"level" db:"level"`
	ResourceType       string `json:"resource_type" db:"resource_type"`
	MaxCapacity        int    `json:"max_capacity" db:"max_capacity"`
	ResourcePerSec     int    `json:"resource_per_sec" db:"resource_per_sec"`
	UpgradeCost        int    `json:"upgrade_cost" db:"upgrade_cost"`
	UpgradeCostType    string `json:"upgrade_cost_type" db:"upgrade_cost_type"`
	UpgradeDurationSec int    `json:"upgrade_duration_sec" db:"upgrade_duration_sec"`
	MaxHP              int    `json:"max_hp" db:"max_hp"`
	Size               int    `json:"size" db:"size"`
}

type TownHallConfig struct {
	Name               string `json:"name" db:"name"`
	Level              int    `json:"level" db:"level"`
	UpgradeCost        int    `json:"upgrade_cost" db:"upgrade_cost"`
	UpgradeCostType    string `json:"upgrade_cost_type" db:"upgrade_cost_type"`
	UpgradeDurationSec int    `json:"upgrade_duration_sec" db:"upgrade_duration_sec"`
	MaxHP              int    `json:"max_hp" db:"max_hp"`
	Size               int    `json:"size" db:"size"`
}

type DefenseConfig struct {
	Name               string `json:"name" db:"name"`
	Level              int    `json:"level" db:"level"`
	UpgradeCost        int    `json:"upgrade_cost" db:"upgrade_cost"`
	UpgradeCostType    string `json:"upgrade_cost_type" db:"upgrade_cost_type"`
	UpgradeDurationSec int    `json:"upgrade_duration_sec" db:"upgrade_duration_sec"`
	DPS                int    `json:"dps" db:"dps"`
	MaxHP              int    `json:"max_hp" db:"max_hp"`
	Range              int    `json:"range" db:"range"`
	AOERange           int    `json:"aoe_range" db:"aoe_range"`
	Size               int    `json:"size" db:"size"`
}

type TrainingGroundsConfig struct {
	Name               string `json:"name" db:"name"`
	Level              int    `json:"level" db:"level"`
	HousingSpace       int    `json:"housing_space" db:"housing_space"`
	UpgradeCost        int    `json:"upgrade_cost" db:"upgrade_cost"`
	UpgradeCostType    string `json:"upgrade_cost_type" db:"upgrade_cost_type"`
	UpgradeDurationSec int    `json:"upgrade_duration_sec" db:"upgrade_duration_sec"`
	MaxHP              int    `json:"max_hp" db:"max_hp"`
	Size               int    `json:"size" db:"size"`
}

type TroopConfig struct {
	Name         string `json:"name" db:"name"`
	DPS          int    `json:"dps" db:"dps"`
	Health       int    `json:"health" db:"health"`
	Range        int    `json:"range" db:"range"`
	Speed        int    `json:"speed" db:"speed"`
	HousingSpace int    `json:"housing_space" db:"housing_space"`
	TrainingCost int    `json:"training_cost" db:"training_cost"`
}
