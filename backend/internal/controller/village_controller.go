// GET    /api/village                → get my village (all building instances)
// POST   /api/village/buildings      → place a building (validate: TH level allows it,
//                                      max count not exceeded via game_progression_config,
//                                      player has enough currency, grid position valid)
// PUT    /api/village/buildings/:id  → move (update pos_x, pos_y) or upgrade (level++)
// DELETE /api/village/buildings/:id  → sell building (refund partial cost)

package controller

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pranav7002/MVC_Assignment/internal/models"
)

type VillageServiceInterface interface {
	GetBuildings(userID string) ([]models.Building, error)
}

type VillageController struct {
	VillageService VillageServiceInterface
}

func (villageController *VillageController) BuildingHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	buildings, err := villageController.VillageService.GetBuildings(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("Something bad happened on the server :/"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(buildings)

	if err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}
}