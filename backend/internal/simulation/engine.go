package simulation

type Battle struct {
	BattleGrid *BattleGrid
	Buildings []*BuildingEntity
	Troops []*TroopEntity
	Tick int 
}

func StartEngine(buildingInputs []BuildingInput, troopInputs []TroopDrop) *Battle {
	var buildings []*BuildingEntity
	var troops []*TroopEntity

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
			TargetTroop: 0,        
			CooldownTick: 10,
		})
	}

	for i, t := range troopInputs {
		troops = append(troops, &TroopEntity{
			ID: i,   
			Name: t.Name,   
			Pos: t.Pos,    
			HP: t.HP,        
			DPS: t.DPS,   
			Range: t.Range,    
			Speed: t.Speed,             
			Dead: false,   
			TargetID: 0,          
			Path: nil,    
		})
	}

	g := NewGrid(buildings)

	return &Battle{
		BattleGrid: g,
		Buildings: buildings,
		Troops: troops, 
		Tick: 0, 
	}
}