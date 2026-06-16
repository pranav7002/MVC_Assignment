package simulation

type BuildingEntity struct {
    ID        int
    Name      string
    Type      string           
    Pos       Position
    Size      int
    HP        int
    Destroyed bool
    // Only for defense buildings:
    DPS          int
    MaxRange     int
    MinRange     int           // Mortar dead zone (e.g., 4 cells)
    AOERange     int           // only for Mortar (e.g., 2 cells)
    TargetTroop  int           // ID of troop currently being targeted
    CooldownTick int           // tick when this defense can fire again
}

const cooldownTicks int = 10

func (b *BuildingEntity) UpdateBuilding(tick int, troops []*TroopEntity, g *BattleGrid) {
	if b.Destroyed {
		removeBuilding(b, g)
		return
	}

	if tick != b.CooldownTick {
		return 
	}

	t := findTargetTroop(troops, b)
	if t == nil {
		return
	}

	switch b.Name {
	case "mortar":
		targets := findTargetsInAOERange(troops, t.Pos, b.AOERange)
		for _, target := range targets {
			target.HP = target.HP - b.DPS
			if (t.HP <= 0) {
				t.Dead = true
			}
		}
	default:
		t.HP = t.HP - b.DPS
		if (t.HP <= 0) {
			t.Dead = true
		}
	}
	b.CooldownTick = b.CooldownTick + cooldownTicks
}

func removeBuilding(b *BuildingEntity, g *BattleGrid) {
	for i := b.Pos.X; i < b.Size + b.Pos.X; i++ {
		for j := b.Pos.Y; j < b.Size + b.Pos.Y; j++ {
			g.OccupiedGrid[i][j] = false
			g.TypeGrid[i][j] = ""
			g.IDGrid[i][j] = 0 
		}
	}
}

func findTargetTroop(troops []*TroopEntity, b *BuildingEntity) *TroopEntity {
	cx := b.Pos.X + (b.Size-1)/2
	cy := b.Pos.Y + (b.Size-1)/2
	for _, t := range troops {
        if t.Dead { 
			continue 
		}
		if Dist(Position{cx, cy}, t.Pos) <= float64(b.MaxRange) && Dist(Position{cx, cy}, t.Pos) >= float64(b.MinRange) {
			return t
		}
	}
	return nil
}

func findTargetsInAOERange(troops []*TroopEntity ,p Position, AOERange int) []*TroopEntity {
	var targets []*TroopEntity
	for _, t := range troops {
		if Dist(t.Pos, p) <= float64(AOERange) {
			targets = append(targets, t)
		}
	}
	return targets
}