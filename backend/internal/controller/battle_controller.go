package controller

import (
	"encoding/json"
	"math/rand"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/pranav7002/MVC_Assignment/internal/middleware"
	"github.com/pranav7002/MVC_Assignment/internal/models"
	"github.com/pranav7002/MVC_Assignment/internal/simulation"
)
type BattleController struct {
	BattleService  BattleServiceInterface
	VillageService VillageServiceInterface
	WSUpgrader     websocket.Upgrader
	BattleManager  *models.BattleManager
}
type BattleServiceInterface interface {
	FilterTroop(t []models.TroopDropBody, buildings []models.Building) []models.TroopDropBody
	HydrateTroop(t models.TroopDropBody) (simulation.TroopDrop, error)
	HydrateBuilding(b []models.Building) ([]simulation.BuildingInput, error)
 	SaveBattleResult(userID, defendersID string, stars, destructionPct int) error
}

func (c *BattleController) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
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
	buildingInput, err := c.BattleService.HydrateBuilding(buildings)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	conn, err := c.WSUpgrader.Upgrade(w, r, nil)
    if err != nil {
        return  
    }
	defer conn.Close()

	attacker := &models.Client{
		ID: userID,
		Conn: conn,
		Send: make(chan []byte),
	}

	battleID := rand.Intn(1e9)

    c.BattleManager.Mu.Lock()
    c.BattleManager.Battles[battleID] = append(c.BattleManager.Battles[battleID], attacker)
    c.BattleManager.Mu.Unlock()

	go attacker.Read()
	go attacker.Write()

	destructionPct, stars := simulation.NewBattle(buildingInput, nil).Simulate()
}

func (c *BattleController) HandleTroopDrop(client *models.Client, b *simulation.Battle) {
	var troopDrop models.TroopDropBody
	for {
		select {
		case msg := <- client.Incoming:
			if err := json.Unmarshal(msg, &troopDrop); err != nil {
				return 
			}
			t, err := c.BattleService.HydrateTroop(troopDrop) 
			if err != nil {
				return
			}
			b.Add(t)
		}
	}
}