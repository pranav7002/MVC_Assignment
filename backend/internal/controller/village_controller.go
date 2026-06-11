// PUT    /api/village/buildings/:id  → move (update pos_x, pos_y) or upgrade (level++)
// DELETE /api/village/buildings/:id  → sell building (refund partial cost)

package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/pranav7002/MVC_Assignment/internal/models"
)

type VillageServiceInterface interface {
	GetBuildings(userID string) ([]models.Building, error)
	CreateBuilding(userID string, reqBody models.BuildingCreationRequestBody) error
	MoveBuilding(userID string, buildingID int64, reqBody models.BuildingPositionRequestBody) error
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

func (villageController *VillageController) BuildingCreationHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	reqBody := new(models.BuildingCreationRequestBody)

	err := json.NewDecoder(r.Body).Decode(reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Please provide the correct input!!"))
		return
	}

	if err := villageController.VillageService.CreateBuilding(userID, *reqBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Building Created successfully!"))
}

func (villageController *VillageController) BuildingPositionHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	buildingID, err := strconv.ParseInt(chi.URLParam(r, "buildingID"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid building ID. It must be an integer.", http.StatusBadRequest)
		return
	}

	reqBody := new(models.BuildingPositionRequestBody)

	if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Please provide the correct input!!"))
		return
	}

	if err := villageController.VillageService.MoveBuilding(userID, buildingID, *reqBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Building Moved successfully!"))
}
