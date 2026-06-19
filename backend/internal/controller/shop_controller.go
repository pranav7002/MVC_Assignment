package controller

import (
	"net/http"

	"github.com/pranav7002/MVC_Assignment/internal/middleware"
	"github.com/pranav7002/MVC_Assignment/internal/services"
)

type ShopController struct {
	ShopService *services.ShopService
}

func (c *ShopController) ShopBuildingsHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.ContextKey("user_id")).(string)
	if !ok {
		WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	buildings, err := c.ShopService.GetShopBuildings(r.Context(), userID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, buildings)
}

func (c *ShopController) ShopTroopsHandler(w http.ResponseWriter, r *http.Request) {
	troops, err := c.ShopService.GetShopTroops(r.Context())
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, troops)
}
