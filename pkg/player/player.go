package player

import "github.com/sh-miyoshi/go-rockmanexe/pkg/chip"

const (
	defaultHP        uint = 100
	defaultShotPower uint = 1

	// FolderSize ...
	FolderSize = 10 // debug
)

// Player ...
type Player struct {
	HP         uint
	HPMax      uint
	ShotPower  uint
	Zenny      uint
	ChipFolder [FolderSize]int
}

// New returns player data with default values
func New() *Player {
	return &Player{
		HP:        defaultHP,
		HPMax:     defaultHP,
		ShotPower: defaultShotPower,
		Zenny:     0,
		ChipFolder: [FolderSize]int{
			chip.IDCannon,
			chip.IDCannon,
			chip.IDCannon,
			chip.IDCannon,
			chip.IDCannon,
			chip.IDCannon,
			chip.IDCannon,
			chip.IDCannon,
			chip.IDCannon,
			chip.IDCannon,
		},
	}
}

// TODO NewWithSaveData(fname string) (*Player, error)
