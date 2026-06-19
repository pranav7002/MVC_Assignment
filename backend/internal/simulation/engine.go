package simulation

import (
	"math/rand"
	"sync"
)

type Battle struct {
	BattleGrid        *BattleGrid
	Buildings         []*BuildingEntity
	Troops            []*TroopEntity
	TownHallDestroyed bool
	TotalBuildingHP   int
	Tick              int

	Mu *sync.Mutex
}

type Result struct {
	DestructionPct int
	Stars          int
}

func NewBattle(buildingInputs []BuildingInput) *Battle {
	var buildings []*BuildingEntity
	var totalBuildingHP int

	for _, b := range buildingInputs {
		buildings = append(buildings, &BuildingEntity{
			ID:        b.ID,
			Name:      b.Name,
			Pos:       b.Pos,
			Type:      b.Type,
			Size:      b.Size,
			HP:        b.HP,
			Destroyed: false,

			DPS:         b.DPS,
			MaxRange:    b.MaxRange,
			MinRange:    b.MinRange,
			AOERange:    b.AOERange,
			TargetTroop: -1,
			Cooldown:    0,
		})

		totalBuildingHP = totalBuildingHP + b.HP
	}

	g := NewGrid(buildings)

	return &Battle{
		BattleGrid:      g,
		Buildings:       buildings,
		Tick:            0,
		TotalBuildingHP: totalBuildingHP,
		Mu:              new(sync.Mutex),
	}
}

func (b *Battle) Add(t TroopDrop) {
	b.Troops = append(b.Troops, &TroopEntity{
		ID:       rand.Intn(1e9),
		Name:     t.Name,
		Pos:      t.Pos,
		HP:       t.HP,
		DPS:      t.DPS,
		Range:    t.Range,
		Dead:     false,
		TargetID: -1,
		Path:     nil,
	})
}

func (b *Battle) Step() (Result, bool) {
	finalHP := 0
	allTroopsDead := false
	allBuildingsDestroyed := false
	stars := 0
	destructionPct := 0

	for _, troop := range b.Troops {
		troop.Update(b.Buildings, b.BattleGrid)
	}
	for _, building := range b.Buildings {
		building.Update(b.Tick, b.Troops, b.BattleGrid)
		if building.Name == "Town Hall" && building.Destroyed {
			b.TownHallDestroyed = true
		}
	}

	allTroopsDead = true
	for _, troop := range b.Troops {
		if !troop.Dead {
			allTroopsDead = false
			break
		}
	}
	allBuildingsDestroyed = true
	for _, building := range b.Buildings {
		if !building.Destroyed {
			allBuildingsDestroyed = false
			break
		}
	}

	if allBuildingsDestroyed || allTroopsDead || b.Tick == 1800 {
		for _, building := range b.Buildings {
			finalHP += building.HP
		}
		destructionPct = ((b.TotalBuildingHP - finalHP) * 100) / b.TotalBuildingHP
		if destructionPct >= 50 {
			stars++
		}
		if b.TownHallDestroyed {
			stars++
		}
		if allBuildingsDestroyed {
			stars++
		}
		return Result{
			DestructionPct: destructionPct,
			Stars:          stars,
		}, true
	}
	b.Tick++
	return Result{}, false
}

func (b *Battle) GetState() BattleState {
	troops := make([]TroopState, 0, len(b.Troops))
	for _, t := range b.Troops {
		troops = append(troops, TroopState{
			ID:   t.ID,
			Name: t.Name,
			Pos:  t.Pos,
			HP:   t.HP,
			Dead: t.Dead,
		})
	}

	buildings := make([]BuildingState, 0, len(b.Buildings))
	for _, bg := range b.Buildings {
		buildings = append(buildings, BuildingState{
			ID:   bg.ID,
			Name: bg.Name,
			HP:   bg.HP,
			Dead: bg.Destroyed,
		})
	}

	finalHP := 0
	for _, bg := range b.Buildings {
		finalHP += bg.HP
	}
	destructionPct := 0
	if b.TotalBuildingHP > 0 {
		destructionPct = ((b.TotalBuildingHP - finalHP) * 100) / b.TotalBuildingHP
	}

	stars := 0
	if destructionPct >= 50 {
		stars++
	}
	if b.TownHallDestroyed {
		stars++
	}
	allDestroyed := true
	for _, bg := range b.Buildings {
		if !bg.Destroyed {
			allDestroyed = false
			break
		}
	}
	if allDestroyed {
		stars++
	}

	return BattleState{
		DestructionPct: destructionPct,
		Stars:          stars,
		Troops:         troops,
		Buildings:      buildings,
	}
}
