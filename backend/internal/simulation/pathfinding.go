package simulation

// Backwards implementation of https://www.youtube.com/watch?v=KiCBXu4P-2Y&t=44s

func FindTarget(t *TroopEntity, g *BattleGrid) BuildingTarget {
	var path []Position
	var id int
	switch t.Name {
	case "Goblin":
		path, id = bfs(t.Pos, g, "storage")
		if path == nil {
			path, id = bfs(t.Pos, g, "")
		}
	case "Giant":
		path, id = bfs(t.Pos, g, "defense")
		if path == nil {
			path, id = bfs(t.Pos, g, "")
		}
	default:
		path, id = bfs(t.Pos, g, "")
	}

	if path == nil {
		return BuildingTarget{}
	}

	return BuildingTarget{
		ID:   id,
		Path: path,
	}
}

// filter = "" for any building
func bfs(p Position, g *BattleGrid, filter string) ([]Position, int) {
	queue := make([]Position, 0)
	parent := make(map[Position]Position)

	reached := false

	var visited [GridSize][GridSize]bool

	dirs := [8]Position{
		{1, 0}, {-1, 0},
		{0, 1}, {0, -1},
		{1, 1}, {1, -1},
		{-1, 1}, {-1, -1},
	}

	queue = append(queue, Position{p.X, p.Y})
	visited[p.X][p.Y] = true

	var ID int
	var current Position
	for len(queue) > 0 {
		current = queue[0]
		queue = queue[1:]

		var isAdj bool
		isAdj, ID = isAdjacent(current, g, filter)
		if isAdj {
			reached = true
			break
		}

		for _, d := range dirs {
			n := Position{X: current.X + d.X, Y: current.Y + d.Y}

			// out of bounds
			if n.X < 0 || n.Y < 0 || n.X > GridSize-1 || n.Y > GridSize-1 {
				continue
			}

			// skip if any adjacent to diagonal is blocked
			if d.X != 0 && d.Y != 0 {
				if g.OccupiedGrid[current.X][n.Y] || g.OccupiedGrid[n.X][current.Y] {
					continue
				}
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
		return nil, 0
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

	return path, ID
}

func isAdjacent(p Position, g *BattleGrid, filter string) (bool, int) {
	dirs := [8]Position{
		{1, 0}, {-1, 0},
		{0, 1}, {0, -1},
		{1, 1}, {1, -1},
		{-1, 1}, {-1, -1},
	}

	for _, d := range dirs {
		n := Position{X: p.X + d.X, Y: p.Y + d.Y}

		if n.X < 0 || n.Y < 0 || n.X > GridSize-1 || n.Y > GridSize-1 {
			continue
		}

		if !g.OccupiedGrid[n.X][n.Y] {
			continue
		}

		if filter == "" || filter == g.TypeGrid[n.X][n.Y] {
			return true, g.IDGrid[n.X][n.Y]
		}
	}
	return false, 0
}
