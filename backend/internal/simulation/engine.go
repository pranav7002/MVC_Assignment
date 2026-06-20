package simulation

import (
	"math/rand"
	"sync"
)

const LAST_TICK int = 1200

type Battle struct {
	BattleGrid        *BattleGrid
	Buildings         []*BuildingEntity
	Troops     		  []*TroopEntity
	TotalTroops		  int
	TownHallDestroyed bool
	TotalBuildingHP   int
	Tick              int

	Mu *sync.Mutex
}

func NewBattle(buildingInputs []BuildingInput, totalTroops int) *Battle {
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
			TargetTroop: 0,
			Cooldown:    0,
		})

		totalBuildingHP = totalBuildingHP + b.HP
	}

	g := NewGrid(buildings)

	return &Battle{
		BattleGrid:      g,
		Buildings:       buildings,
		TotalTroops: 	 totalTroops,
		Tick:            0,
		TotalBuildingHP: totalBuildingHP,
		Mu:              new(sync.Mutex),
	}
}

func (b *Battle) Add(t TroopDrop) {
	b.Troops = append(b.Troops, &TroopEntity{
		ID:              rand.Intn(1e9),
		Name:            t.Name,
		Pos:             t.Pos,
		HP:              t.HP,
		DPS:             t.DPS,
		Range:           t.Range,
		Dead:            false,
		TargetID:        0,
		Path:            nil,
		Speed:           t.Speed,
		Steps: 0,
	})
	b.TotalTroops--
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
	destructionPct := 0
	stars := 0

	troops := make([]TroopState, 0, len(b.Troops))
	allTroopsDead := true
	for _, t := range b.Troops {
		troops = append(troops, TroopState{
			ID:   t.ID,
			Name: t.Name,
			Pos:  t.Pos,
			HP:   max(0, t.HP),
			Dead: t.Dead,
		})
		if !t.Dead {
			allTroopsDead = false
		}
	}
	buildings := make([]BuildingState, 0, len(b.Buildings))
	allBuildingsDestroyed := true
	buildingHpRemaining := 0
	for _, bg := range b.Buildings {
		buildings = append(buildings, BuildingState{
			ID:    bg.ID,
			Name:  bg.Name,
			HP:    max(0, bg.HP),
			MaxHP: bg.MaxHP,
			Dead:  bg.Destroyed,
		})
		buildingHpRemaining += max(0, bg.HP)
		if !bg.Destroyed {
			allBuildingsDestroyed = false
		}
	}

	destructionPct = ((b.TotalBuildingHP - buildingHpRemaining) * 100) / b.TotalBuildingHP
	if b.TownHallDestroyed {
		stars++
	}
	if destructionPct >= 50 {
		stars++
	}
	if allBuildingsDestroyed {
		stars++ 
	}

	done := (b.TotalTroops == 0 && allTroopsDead) || allBuildingsDestroyed || b.Tick == LAST_TICK

	if done {
		return BattleState{
		DestructionPct: destructionPct,
		Stars:          stars,
		Troops:         troops,
		Buildings:      buildings,
	}, done
	}

	return BattleState{
		DestructionPct: destructionPct,
		Stars:          stars,
		Troops:         troops,
		Buildings:      buildings,
	}, false
}
