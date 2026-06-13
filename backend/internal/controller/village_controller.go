package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/pranav7002/MVC_Assignment/internal/middleware"
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
	userID, ok := r.Context().Value(middleware.ContextKey("user_id")).(string)
	if !ok {
		WriteError(w, http.StatusUnauthorized, "User ID not found in context")
		return
	}

	buildings, err := villageController.VillageService.GetBuildings(userID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Something bad happened on the server :/")
		return
	}

	WriteJSON(w, http.StatusOK, buildings)
}

func (villageController *VillageController) BuildingCreationHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.ContextKey("user_id")).(string)
	if !ok {
		WriteError(w, http.StatusUnauthorized, "User ID not found in context")
		return
	}

	reqBody := new(models.BuildingCreationRequestBody)

	err := json.NewDecoder(r.Body).Decode(reqBody)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Please provide the correct input!!")
		return
	}

	if err := villageController.VillageService.CreateBuilding(userID, *reqBody); err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, "Building Created successfully!")
}

func (villageController *VillageController) BuildingPositionHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.ContextKey("user_id")).(string)
	if !ok {
		WriteError(w, http.StatusUnauthorized, "User ID not found in context")
		return
	}
	
	buildingID, err := strconv.ParseInt(chi.URLParam(r, "buildingID"), 10, 64)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid building ID. It must be an integer.")
		return
	}

	reqBody := new(models.BuildingPositionRequestBody)

	if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
		WriteError(w, http.StatusBadRequest, "Please provide the correct input!!")
		return
	}

	if err := villageController.VillageService.MoveBuilding(userID, buildingID, *reqBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	WriteJSON(w, http.StatusOK, "Building moved successfully!")
}

func (villageController *VillageController) BuildingUpgradeHandler(w http.ResponseWriter, r *http.Request) {	
	userID, ok := r.Context().Value(middleware.ContextKey("user_id")).(string)
	if !ok {
		WriteError(w, http.StatusUnauthorized, "User ID not found in context")
		return
	}
	
	buildingID, err := strconv.ParseInt(chi.URLParam(r, "BuildingID"), 10, 64) 
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid building ID. It must be an integer.")
		return
	}

	if err := villageController.VillageService.UpgradeBuilding(userID, buildingID); err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return 
	} 

	WriteJSON(w, http.StatusOK, "Building upgraded successfully!")
}