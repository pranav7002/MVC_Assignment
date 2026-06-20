package simulation

type TroopEntity struct {
	ID              int
	Name            string
	Pos             Position
	HP              int
	DPS             int
	Range           int
	Dead            bool
	TargetID        int // ID of the building troop is targeting
	Path            []Position
	Speed           float64
	Steps 			float64
}

const buffer float64 = 0.5

func (t *TroopEntity) Update(buildings []*BuildingEntity, g *BattleGrid) {
	if t.Dead {
		return
	}

	// 1. Troop has no target
	hasTarget, b := buildingExists(t.TargetID, buildings)
	if !hasTarget {
		target := FindTarget(t, g)
		if target.Path == nil {
			return
		}

		t.TargetID = target.ID
		t.Path = target.Path
		return
	}

	// 2. Target in Attack Range
	if inRange(t, b) {
		b.HP = b.HP - t.DPS
		if b.HP <= 0 {
			b.Destroyed = true
		}
		return
	}

	// 3. Troop has target, not in range
	if len(t.Path) > 1 {
		t.Steps += t.Speed
		if t.Steps >= 1.0 && len(t.Path) > 1 {
			t.Pos = t.Path[1]
			t.Path = t.Path[1:]
			t.Steps -= 1.0
		}
	}
}

func buildingExists(id int, buildings []*BuildingEntity) (bool, *BuildingEntity) {
	for _, b := range buildings {
		if b.ID == id && !b.Destroyed {
			return true, b
		}
	}
	return false, nil
}

func inRange(t *TroopEntity, b *BuildingEntity) bool {
	cx := max(b.Pos.X, min(t.Pos.X, b.Pos.X+b.Size-1))
	cy := max(b.Pos.Y, min(t.Pos.Y, b.Pos.Y+b.Size-1))

	if Dist(Position{cx, cy}, t.Pos) <= float64(t.Range)+buffer {
		return true
	}

	return false
}
