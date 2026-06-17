package controller

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"sync"
	"time"

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
	BattleManager  *BattleManager
}



type BattleManager struct {
    mu      sync.Mutex
    battles map[int][]*Client
}

type Client struct {
	ID string
	Conn *websocket.Conn

	Send chan []byte
}

func (c *BattleController) handleWebSocket(w http.ResponseWriter, r *http.Request) {
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
	destructionPct, stars = simulation.NewBattle(buildingInput).Simulate()

	conn, err := c.WSUpgrader.Upgrade(w, r, nil)
    if err != nil {
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode(APIResponse{Error: "upgrade failed"})
        return  
    }
	defer conn.Close()

	client := &Client{
		ID: userID,
		Conn: conn,
		Send: make(chan []byte),
	}
	battleID := rand.Intn(1e9)
	c.BattleManager.Add(battleID, client)
	
	go client.read()
	go client.write()
}

func (bm *BattleManager) Add(battleID int, client *Client) {
    bm.mu.Lock()
    defer bm.mu.Unlock()
    bm.battles[battleID] = append(bm.battles[battleID], client)
}

func (c *Client) write() {
	for data := range c.Send {
		err := c.Conn.WriteMessage(
			websocket.TextMessage,
			data,
		)

		if err != nil {
			return
		}
	}
}

func (c *Client) read() {
	for {
		messageType, message, err := c.Conn.ReadMessage()
		if err != nil {
			return
		}

		if messageType != websocket.TextMessage {
			continue
		}

		msg := new(Message)
		if err := json.Unmarshal(message, msg); err != nil {
			return
		}
	}
}

type Message struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// type BattleServiceInterface interface {
// 	FilterTroop(t []models.TroopDropRequestBody, buildings []models.Building) []models.TroopDropRequestBody
// 	HydrateTroop(t []models.TroopDropRequestBody) ([]simulation.TroopDrop, error)
// 	HydrateBuilding(b []models.Building) ([]simulation.BuildingInput, error)
//  SaveBattleResult(userID, defendersID string, stars, destructionPct int) error
// }

// func (c *BattleController) RunSimulation(w http.ResponseWriter, r *http.Request) {
// 	userID, ok := r.Context().Value(middleware.ContextKey("user_id")).(string)
// 	if !ok {
// 		WriteError(w, http.StatusUnauthorized, "User ID not found in context")
// 		return
// 	}
// 	defendersID := chi.URLParam(r, "defendersID")

// 	buildings, err := c.VillageService.GetBuildings(defendersID)
// 	if err != nil {
// 		WriteError(w, http.StatusNotFound, err.Error())
// 		return
// 	}

// 	var troopDrop []models.TroopDropRequestBody
// 	if err := json.NewDecoder(r.Body).Decode(&troopDrop); err != nil {
// 		WriteError(w, http.StatusBadRequest, "Invalid request body")
// 		return
// 	}

// 	troopInput, err := c.BattleService.HydrateTroop(c.BattleService.FilterTroop(troopDrop, buildings))
// 	if err != nil {
// 		WriteError(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	var destructionPct, stars int

// 	if len(troopInput) >= 1 {
// 		buildingInput, err := c.BattleService.HydrateBuilding(buildings)
// 		if err != nil {
// 			WriteError(w, http.StatusInternalServerError, err.Error())
// 			return
// 		}
// 		destructionPct, stars = simulation.NewBattle(buildingInput, troopInput).Simulate()
// 	}

// 	if err := c.BattleService.SaveBattleResult(userID, defendersID, stars, destructionPct); err != nil {
// 		WriteError(w, http.StatusInternalServerError, "Failed to save battle result")
// 		return
// 	}

// 	WriteJSON(w, http.StatusOK, map[string]int{
// 		"stars":              stars,
// 		"destruction_percent": destructionPct,
// 	})
// } 