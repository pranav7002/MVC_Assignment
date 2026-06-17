package controller

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pranav7002/MVC_Assignment/internal/middleware"
	"github.com/pranav7002/MVC_Assignment/internal/models"
	"github.com/pranav7002/MVC_Assignment/internal/simulation"
)

type BattleController struct {
	BattleService  BattleServiceInterface
	VillageService VillageServiceInterface
}

type BattleServiceInterface interface {
	FilterTroop(t []models.TroopDropRequestBody, buildings []models.Building) []models.TroopDropRequestBody
	HydrateTroop(t []models.TroopDropRequestBody) ([]simulation.TroopDrop, error)
	HydrateBuilding(b []models.Building) ([]simulation.BuildingInput, error)
	SaveBattleResult(userID, defendersID string, stars, destructionPct int) error
}

func (c *BattleController) RunSimulation(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.ContextKey("user_id")).(string)
	if !ok {
		WriteError(w, http.StatusUnauthorized, "User ID not found in context")
		return
	}
	defendersID := chi.URLParam(r, "defendersID")

	buildings, err := c.VillageService.GetBuildings(defendersID)
	if err != nil {
		WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	var troopDrop []models.TroopDropRequestBody
	if err := json.NewDecoder(r.Body).Decode(&troopDrop); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	troopInput, err := c.BattleService.HydrateTroop(c.BattleService.FilterTroop(troopDrop, buildings))
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var destructionPct, stars int

	if len(troopInput) >= 1 {
		buildingInput, err := c.BattleService.HydrateBuilding(buildings)
		if err != nil {
			WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
		destructionPct, stars = simulation.NewBattle(buildingInput, troopInput).Simulate()
	}

	if err := c.BattleService.SaveBattleResult(userID, defendersID, stars, destructionPct); err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to save battle result")
		return
	}

	WriteJSON(w, http.StatusOK, map[string]int{
		"stars":              stars,
		"destruction_percent": destructionPct,
	})
} 