package controller

import (
	"encoding/json"
	"net/http"

	"github.com/pranav7002/MVC_Assignment/internal/middleware"
	"github.com/pranav7002/MVC_Assignment/internal/models"
)

type TroopServiceInterface interface {
	TrainTroop(userID string, troopName string, quantity int) error
	GetTrainedTroops(userID string) ([]models.TroopTrained, error)
}

type TroopController struct {
	TroopService TroopServiceInterface
}

func (c *TroopController) TrainTroopHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.ContextKey("user_id")).(string)
	if !ok {
		WriteError(w, http.StatusUnauthorized, "User ID not found in context")
		return
	}

	reqBody := new(models.TroopTrainingReqBody)
	if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
		WriteError(w, http.StatusBadRequest, "Please provide the correct input!!")
		return
	}

	if err := c.TroopService.TrainTroop(userID, reqBody.TroopName, reqBody.Quantity); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, "Troops trained successfully!")
}

func (c *TroopController) TroopHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.ContextKey("user_id")).(string)
	if !ok {
		WriteError(w, http.StatusUnauthorized, "User ID not found in context")
		return
	}

	troops, err := c.TroopService.GetTrainedTroops(userID) 
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, troops)
}