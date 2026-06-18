package simulation

import "math/rand"

type Battle struct {
	BattleGrid *BattleGrid
	Buildings []*BuildingEntity
	Troops []*TroopEntity
	TownHallDestroyed bool
	TotalBuildingHP int
	Tick int 
}

func NewBattle(buildingInputs []BuildingInput, troopInputs []TroopDrop) *Battle {
	var buildings []*BuildingEntity
	var totalBuildingHP int

	for _, b := range buildingInputs {
		buildings = append(buildings, &BuildingEntity{
			ID: b.ID,       
			Name: b.Name,  
			Pos: b.Pos,    
			Type: b.Type,                 
			Size: b.Size,     
			HP: b.HP,       
			Destroyed: false,
			
			DPS: b.DPS,      
			MaxRange: b.MaxRange,   
			MinRange: b.MinRange,               
			AOERange: b.AOERange,             
			TargetTroop: -1,        
			Cooldown: 0,
		})

		totalBuildingHP = totalBuildingHP + b.HP
	}

	g := NewGrid(buildings)

	return &Battle{
		BattleGrid: g,
		Buildings: buildings,
		Tick: 0, 
		TotalBuildingHP: totalBuildingHP,
	}
}

func (b *Battle) Add(t TroopDrop) {
		b.Troops = append(b.Troops, &TroopEntity{
			ID: rand.Intn(1e9),   
			Name: t.Name,   
			Pos: t.Pos,    
			HP: t.HP,        
			DPS: t.DPS,   
			Range: t.Range,              
			Dead: false,   
			TargetID: -1,          
			Path: nil,    
		})
}

func (b *Battle) Step() (int, int) {
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
	}
	b.Tick++

	return destructionPct, stars
}