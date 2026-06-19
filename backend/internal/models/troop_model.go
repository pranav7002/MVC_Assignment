package models

type TroopTrainingReqBody struct {
	TroopName string `json:"troop_name"`
	Quantity  int    `json:"quantity"`
}

type TroopTrained struct {
	UserID    string `json:"user_id" db:"user_id"`
	TroopName string `json:"troop_name" db:"troop_name"`
	Quantity  int    `json:"quantity" db:"quantity"`
}
