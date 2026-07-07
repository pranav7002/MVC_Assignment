package simulation

import "math"

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type BuildingInput struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Pos      Position `json:"pos"`
	Size     int      `json:"size"`
	HP       int      `json:"hp"`
	DPS      int      `json:"dps"`
	MaxRange int      `json:"max_range"`
	MinRange int      `json:"min_range"`
	AOERange int      `json:"aoe_range"`
}

type TroopDrop struct {
	Name  string   `json:"name"`
	Pos   Position `json:"pos"`
	HP    int      `json:"hp"`
	DPS   int      `json:"dps"`
	Range int      `json:"range"`
	Speed float64  `json:"speed"`
}

type BattleGrid struct {
	OccupiedGrid [GridSize][GridSize]bool
	TypeGrid     [GridSize][GridSize]string
	IDGrid       [GridSize][GridSize]int
}

type BuildingTarget struct {
	ID   int
	Path []Position
}

func Dist(n, m Position) float64 {
	dist := math.Sqrt(math.Pow(float64(m.X-n.X), 2) + math.Pow(float64(m.Y-n.Y), 2))
	return dist
}

// BATTLE EVENTS

type TroopState struct {
	ID   int      `json:"id"`
	Name string   `json:"name"`
	Pos  Position `json:"pos"`
	HP    int      `json:"hp"`
	MaxHP int      `json:"max_hp"`
	Dead  bool     `json:"dead"`
}

type BuildingState struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	HP    int    `json:"hp"`
	MaxHP int    `json:"max_hp"`
	Dead  bool   `json:"dead"`
}

type BattleState struct {
	DestructionPct int             `json:"destruction_pct"`
	Stars          int             `json:"stars"`
	Troops         []TroopState    `json:"troops"`
	Buildings      []BuildingState `json:"buildings"`
}