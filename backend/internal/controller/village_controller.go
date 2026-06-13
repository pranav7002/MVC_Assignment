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
	UpgradeBuilding(userID string, buildingID int64) error
}

type VillageController struct {
	VillageService VillageServiceInterface
}

func (villageController *VillageController) BuildingHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	buildings, err := villageController.VillageService.GetBuildings(userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Something bad happened on the server :/")
		return
	}

	writeJSON(w, http.StatusOK, buildings)
}

func (villageController *VillageController) BuildingCreationHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	reqBody := new(models.BuildingCreationRequestBody)

	err := json.NewDecoder(r.Body).Decode(reqBody)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Please provide the correct input!!")
		return
	}

	if err := villageController.VillageService.CreateBuilding(userID, *reqBody); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, "Building Created successfully!")
}

func (villageController *VillageController) BuildingPositionHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	buildingID, err := strconv.ParseInt(chi.URLParam(r, "buildingID"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid building ID. It must be an integer.")
		return
	}

	reqBody := new(models.BuildingPositionRequestBody)

	if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
		writeError(w, http.StatusBadRequest, "Please provide the correct input!!")
		return
	}

	if err := villageController.VillageService.MoveBuilding(userID, buildingID, *reqBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, "Building moved successfully!")
}

func (villageController *VillageController) BuildingUpgradeHandler(w http.ResponseWriter, r *http.Request) {	
	userID := chi.URLParam(r, "userID")
	buildingID, err := strconv.ParseInt(chi.URLParam(r, "BuildingID"), 10, 64) 
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid building ID. It must be an integer.")
		return
	}

	if err := villageController.VillageService.UpgradeBuilding(userID, buildingID); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return 
	} 

	writeJSON(w, http.StatusOK, "Building upgraded successfully!")
}