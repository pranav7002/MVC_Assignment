package simulation

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
	var troops []*TroopEntity
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

	for i, t := range troopInputs {
		troops = append(troops, &TroopEntity{
			ID: i,   
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

	g := NewGrid(buildings)

	return &Battle{
		BattleGrid: g,
		Buildings: buildings,
		Troops: troops, 
		Tick: 0, 
		TotalBuildingHP: totalBuildingHP,
	}
}

func (b *Battle) Simulate() (int, int) {
	finalHP := 0
	allTroopsDead := false
	allBuildingsDestroyed := false
	stars := 0
	destructionPct := 0

	for b.Tick != 1800 {
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
			break
		}
		b.Tick++
	}
	return destructionPct, stars
}