package simulation

// Backwards implementation of https://www.youtube.com/watch?v=KiCBXu4P-2Y&t=44s

// type BuildingInput struct {
//     ID           int      `json:"id"`
//     Name         string   `json:"name"`          
//     Type         string   `json:"type"`         
//     Pos          Position `json:"pos"`
//     Size         int      `json:"size"`          
//     HP           int      `json:"hp"`           
//     DPS          int      `json:"dps"`          
//     Range        int      `json:"range"`        
// }

// type TroopDrop struct {
//     Name         string   `json:"name"`        
//     Pos          Position `json:"pos"`           
//     HP           int      `json:"hp"`           
//     DPS          int      `json:"dps"`
//     Range        int      `json:"range"`
//     Speed        int      `json:"speed"`
// }



func findPath(t TroopDrop, g *BattleGrid) []Position {
	var path []Position
	switch t.Name {
	case "Barbarian", "Archer", "Wizard": 
		path = findNearestBuilding(Position{X: t.Pos.X, Y: t.Pos.Y}, g)
	// case "Goblin":
	// 	path = findNearestStorage(Position{X: t.Pos.X, Y: t.Pos.Y}, g)
	// case "Giant": 
	// path = findNearestDefence(Position{X: t.Pos.X, Y: t.Pos.Y}, g)
	}

	return path
}

func findNearestBuilding(p Position, g *BattleGrid) []Position {
	queue := make([]Position, 0)
	parent := make(map[Position]Position)

	reached := false

	var visited [GridSize][GridSize]bool 

	dirs := [4]Position{
		{1, 0},
		{-1, 0},
		{0, 1},
		{0, -1},
	}

	queue = append(queue, Position{p.X, p.Y})
	visited[p.X][p.Y] = true

	var current Position
	for len(queue) > 0 {
		current = queue[0]
		queue = queue[1:]

		if isAdjacentToBuilding(current, g) {
			reached = true
			break
		}

		for _, d := range dirs {
			n := Position{X: current.X + d.X, Y: current.Y + d.Y}

			// out of bounds
			if n.X < 0 || n.Y < 0 || n.X > GridSize - 1 || n.Y > GridSize - 1 {
				continue
			}

			// skip visited locations and blocked cells 
			if visited[n.X][n.Y] || g.OccupiedGrid[n.X][n.Y] {
				continue
			}

			queue = append(queue, n)
			parent[Position{n.X, n.Y}] = current
			visited[n.X][n.Y] = true
		}
	}

	if !reached {
		return nil
	}
	
	start := Position{p.X, p.Y}
	var path []Position

	for current != start {
		path = append(path, current)
		current = parent[current]
	}
	path = append(path, start)

	// reverse 
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return path
}

func isAdjacentToBuilding(p Position, g *BattleGrid) bool {
	dirs := [4]Position{
		{1, 0},
		{-1, 0},
		{0, 1},
		{0, -1},
	}

	for _, d := range dirs {
		n := Position{X:p.X + d.X , Y:p.Y + d.Y}

		if n.X < 0 || n.Y < 0 || n.X > GridSize - 1 || n.Y > GridSize - 1 {
			continue
		}

		if g.OccupiedGrid[n.X][n.Y] {
			return true
		}
	}	
	return false 
}

// func findNearestStorage(p Position, g *BattleGrid) []Position {

// }

// func findNearestDefence(p Position, g *BattleGrid) []Position {

// }