package player

// Player ...
type Player struct {
	HP        uint
	HPMax     uint
	ShotPower uint
	Zenny     uint
}

const (
	defaultHP        uint = 100
	defaultShotPower uint = 1
)

// New returns player data with default values
func New() *Player {
	return &Player{
		HP:        defaultHP,
		HPMax:     defaultHP,
		ShotPower: defaultShotPower,
		Zenny:     0,
	}
}

// TODO NewWithSaveData(fname string) (*Player, error)
