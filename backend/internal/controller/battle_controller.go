package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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
	HydrateTroop(t models.TroopDropBody, buildings []models.Building) (simulation.TroopDrop, error)
	HydrateBuildings(b []models.Building) ([]simulation.BuildingInput, error)
	SaveBattleResult(userID, defendersID string, stars, destructionPct int) error
	FindMatch(userID string) (models.MatchmakingResBody, error)
	GetTotalTroops(userID string) (int, error)
}

func (c *BattleController) MatchmakingHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.ContextKey("user_id")).(string)
	if !ok {
		WriteError(w, http.StatusUnauthorized, "User ID not found in context")
		return
	}

	match, err := c.BattleService.FindMatch(userID)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, match)
}

func (c *BattleController) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		WriteError(w, http.StatusUnauthorized, "user_id query param required")
		return
	}
	defendersID := chi.URLParam(r, "defendersID")

	buildings, err := c.VillageService.GetBuildings(defendersID)
	if err != nil {
		WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	buildingInput, err := c.BattleService.HydrateBuildings(buildings)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	totalTroops, err := c.BattleService.GetTotalTroops(userID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
	}

	conn, err := c.WSUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	attacker := &models.Client{
		ID:       userID,
		Conn:     conn,
		Send:     make(chan []byte),
		Incoming: make(chan []byte),
		Done:     make(chan struct{}),
	}

	battleID := uuid.New().String()

	c.BattleManager.Mu.Lock()
	c.BattleManager.Battles[battleID] = append(c.BattleManager.Battles[battleID], attacker)
	c.BattleManager.Mu.Unlock()

	go attacker.Read()
	go attacker.Write()

	battle := simulation.NewBattle(buildingInput, totalTroops)
	go c.HandleTroopDrop(attacker, battle, buildings)

	ticker := time.NewTicker(100 * time.Millisecond)

	defer func() {
		conn.Close()

		c.BattleManager.Mu.Lock()
		delete(c.BattleManager.Battles, battleID)
		c.BattleManager.Mu.Unlock()

		ticker.Stop()
	}()

	var finalState simulation.BattleState
	for range ticker.C {
		battle.Mu.Lock()
		battle.Step()
		state, done := battle.GetState()
		battle.Mu.Unlock()
		msg, err := json.Marshal(&state)
		if err != nil {
			return
		}
		select {
		case attacker.Send <- msg:
			if done {
				finalState = state
				close(attacker.Send)
				goto saveBattle
			}
		case <-attacker.Done:
			return
		}
	}
	saveBattle:
	c.BattleService.SaveBattleResult(userID, defendersID, finalState.Stars, finalState.DestructionPct)
}

func (c *BattleController) HandleTroopDrop(client *models.Client, b *simulation.Battle, buildings []models.Building) {
	var troopDrop models.TroopDropBody
	for msg := range client.Incoming {
		if err := json.Unmarshal(msg, &troopDrop); err != nil {
			return
		}
		t, err := c.BattleService.HydrateTroop(troopDrop, buildings)
		if err != nil {
			return
		}
		b.Mu.Lock()
		b.Add(t)
		b.Mu.Unlock()
	}
}

