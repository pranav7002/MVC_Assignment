package models

type TroopDropRequestBody struct {
	Name string `json:"name"`
	X uint8 `json:"x"`
	Y uint8 `json:"y"`
}