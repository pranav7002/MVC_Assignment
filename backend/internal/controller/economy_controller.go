package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/pranav7002/MVC_Assignment/internal/middleware"
	"github.com/pranav7002/MVC_Assignment/internal/models"
)

type EconomyController struct {
	EconomyService EconomyServiceInterface
	VillageService VillageServiceInterface
}

type EconomyServiceInterface interface {
	CollectGold(userID string, reqTime time.Time) error
	CollectElixir(userID string, reqTime time.Time) error
}

func (economyController *EconomyController) ResourceCollectionHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.ContextKey("user_id")).(string)
	if !ok {
		WriteError(w, http.StatusUnauthorized, "User ID not found in context")
		return
	}

	reqTime := time.Now()

	reqBody := new(models.ResourceCollectionReqBody)
	if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
		WriteError(w, http.StatusBadRequest, "Please provide the correct input!!")
		return
	}

	switch reqBody.ResourceType {
	case "gold":
		if err := economyController.EconomyService.CollectGold(userID, reqTime); err != nil {
			WriteError(w, http.StatusInternalServerError, "Something bad happened on the server :/")
			return
		}
	case "elixir":
		if err := economyController.EconomyService.CollectElixir(userID, reqTime); err != nil {
			WriteError(w, http.StatusInternalServerError, "Something bad happened on the server :/")
			return
		}
	default:
		WriteError(w, http.StatusBadRequest, "Invalid resource type. Must be 'gold' or 'elixir'")
		return 
	}

	WriteJSON(w, http.StatusOK, "Resource collected successfully!")
}