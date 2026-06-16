package simulation

// Backwards implementation of https://www.youtube.com/watch?v=KiCBXu4P-2Y&t=44s

func findPath(t TroopDrop, g *BattleGrid) []Position {
	var path []Position
	switch t.Name {
	case "Goblin":
		path = bfs(t.Pos, g, "storage")
		if path == nil {
			path = bfs(t.Pos, g, "")
		}
	case "Giant": 
		path = bfs(t.Pos, g, "defense")
		if path == nil {
			path = bfs(t.Pos, g, "")
		}
	default: 
		path = bfs(t.Pos, g, "")
	}

	return path
}

// filter = "" for any building 
func bfs(p Position, g *BattleGrid, filter string) []Position {
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

	var current Position
	for len(queue) > 0 {
		current = queue[0]
		queue = queue[1:]

		if isAdjacent(current, g, filter) {
			reached = true
			break
		}

		for _, d := range dirs {
			n := Position{X: current.X + d.X, Y: current.Y + d.Y}

			// out of bounds
			if n.X < 0 || n.Y < 0 || n.X > GridSize - 1 || n.Y > GridSize - 1 {
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

func isAdjacent(p Position, g *BattleGrid, filter string) bool {
	dirs := [8]Position{
		{1, 0}, {-1, 0},
		{0, 1}, {0, -1},
		{1, 1}, {1, -1},
		{-1, 1}, {-1, -1},
	}

	for _, d := range dirs {
		n := Position{X:p.X + d.X , Y:p.Y + d.Y}

		if n.X < 0 || n.Y < 0 || n.X > GridSize - 1 || n.Y > GridSize - 1 {
			continue
		}

		if !g.OccupiedGrid[n.X][n.Y] {
			continue
		}

		if filter == "" || filter == g.TypeGrid[n.X][n.Y] {
			return true
		}
	}	
	return false 
}