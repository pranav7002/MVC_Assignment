package simulation

type Position struct {
    X int `json:"x"`
    Y int `json:"y"`
}

type BuildingInput struct {
    ID           int      `json:"id"`
    Name         string   `json:"name"`          
    Type         string   `json:"type"`         
    Pos          Position `json:"pos"`
    Size         int      `json:"size"`          
    HP           int      `json:"hp"`           
    DPS          int      `json:"dps"`          
    Range        int      `json:"range"`        
}

type TroopDrop struct {
    Name         string   `json:"name"`        
    Pos          Position `json:"pos"`           
    HP           int      `json:"hp"`           
    DPS          int      `json:"dps"`
    Range        int      `json:"range"`
    Speed        int      `json:"speed"`
}

type BattleGrid struct {
	OccupiedGrid [GridSize][GridSize]bool 
	BuildingGrid [GridSize][GridSize]int
}