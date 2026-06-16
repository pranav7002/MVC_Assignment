package simulation

// type BuildingInput struct {
//     ID           int      `json:"id"`
//     Name         string   `json:"name"`          
//     Type         string   `json:"type"`  // "defense", "resource", "storage", "town_hall", "training_grounds"        
//     Pos          Position `json:"pos"`
//     Size         int      `json:"size"`          
//     HP           int      `json:"hp"`           
//     DPS          int      `json:"dps"`          
//     Range        int      `json:"range"`        
// }

const GridSize int = 20

func NewGrid(buildings []BuildingInput) *BattleGrid {
	var og [GridSize][GridSize]bool 
	var tg [GridSize][GridSize]string
	var ig [GridSize][GridSize]int


	for _, b := range buildings {
		for i := b.Pos.X; i < b.Size + b.Pos.X; i++ {
			for j := b.Pos.Y; j < b.Size + b.Pos.Y; j++ {
				og[i][j] = true
				tg[i][j] = b.Type
				ig[i][j] = b.ID
			}
		}
	}

	g := BattleGrid{
		OccupiedGrid: og,
		TypeGrid: tg,
		IDGrid: ig,
	}

	return &g
}

