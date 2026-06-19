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
			MaxHP:     b.HP,
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

func (b *Battle) Step() {
	for _, troop := range b.Troops {
		troop.Update(b.Buildings, b.BattleGrid)
	}
	for _, building := range b.Buildings {
		building.Update(b.Tick, b.Troops, b.BattleGrid)
		if building.Name == "Town Hall" && building.Destroyed {
			b.TownHallDestroyed = true
		}
	}
	b.Tick++
}

func (b *Battle) GetState() (BattleState, bool) {
	finalHP := 0
	stars := 0
	destructionPct := 0

	troops := make([]TroopState, 0, len(b.Troops))
	allTroopsDead := true
	for _, t := range b.Troops {
		troops = append(troops, TroopState{
			ID:   t.ID,
			Name: t.Name,
			Pos:  t.Pos,
			HP:   t.HP,
			Dead: t.Dead,
		})
		if !t.Dead {
			allTroopsDead = false
		}
	}
	if len(b.Troops) == 0 {
		allTroopsDead = false
	}

	buildings := make([]BuildingState, 0, len(b.Buildings))
	allBuildingsDestroyed := true
	for _, bg := range b.Buildings {
		buildings = append(buildings, BuildingState{
			ID:    bg.ID,
			Name:  bg.Name,
			HP:    bg.HP,
			MaxHP: bg.MaxHP,
			Dead:  bg.Destroyed,
		})
		if !bg.Destroyed {
			allBuildingsDestroyed = false
		}
	}

	for _, building := range b.Buildings {
		finalHP += building.HP
	}
	if b.TotalBuildingHP > 0 {
		destructionPct = ((b.TotalBuildingHP - finalHP) * 100) / b.TotalBuildingHP
	}

	if allBuildingsDestroyed || allTroopsDead || b.Tick >= 1800 {
		if destructionPct >= 50 {
			stars++
		}
		if b.TownHallDestroyed {
			stars++
		}
		if allBuildingsDestroyed {
			stars++
		}
		return BattleState{
			DestructionPct: destructionPct,
			Stars:          stars,
			Troops:         troops,
			Buildings:      buildings,
		}, true
	}

	return BattleState{
		DestructionPct: destructionPct,
		Stars:          stars,
		Troops:         troops,
		Buildings:      buildings,
	}, false
}
