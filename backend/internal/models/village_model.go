package models

import "time"

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

type Village struct {
	ID                    int64     `json:"id" db:"id"`
	UserID                string    `json:"user_id" db:"user_id"`
	TownHallLevel         int       `json:"town_hall_level" db:"town_hall_level"`
	Gold                  int       `json:"gold" db:"gold"`
	Elixir                int       `json:"elixir" db:"elixir"`
	GoldLastCollectedAt   time.Time `json:"gold_last_collected_at" db:"gold_last_collected_at"`
	ElixirLastCollectedAt time.Time `json:"elixir_last_collected_at" db:"elixir_last_collected_at"`
}

type BuildingCreationRequestBody struct {
	BuildingType string `json:"building_type"`
	BuildingName string `json:"building_name"`
	PosX         int    `json:"pos_x"`
	PosY         int    `json:"pos_y"`
}

type BuildingPositionRequestBody struct {
	PosX int `json:"pos_x"`
	PosY int `json:"pos_y"`
}

type VillageResBody struct {
	TownHallLevel int `json:"town_hall_level" db:"town_hall_level"`
	Gold          int `json:"gold" db:"gold"`
	Elixir        int `json:"elixir" db:"elixir"`
}
