package models

type Building struct {
	ID           int64  `json:"id"`
	UserID       string `json:"user_id"`
	BuildingType string `json:"building_type"`
	BuildingName string `json:"building_name"`
	Level        int    `json:"level"`
	PosX         int    `json:"pos_x"`
	PosY         int    `json:"pos_y"`
	Size         int    `json:"size"`
	IsUpgrading  bool   `json:"is_upgrading"`
	HP           int    `json:"hp"`
}

