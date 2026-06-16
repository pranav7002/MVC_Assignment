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
    Range        int
    AOERange     int           // mortar
    MinRange     int           // radius where mortar cant attack
    TargetTroop  int           
    CooldownTick int           // ticks bw each attack
}

