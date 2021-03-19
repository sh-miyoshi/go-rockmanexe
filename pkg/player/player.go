package player

import "github.com/sh-miyoshi/go-rockmanexe/pkg/chip"

const (
	defaultHP        uint = 100
	defaultShotPower uint = 1

	// FolderSize ...
	FolderSize = 10 // debug
)

// ChipInfo ...
type ChipInfo struct {
	ID   int
	Code string
}

// Player ...
type Player struct {
	HP         uint
	HPMax      uint
	ShotPower  uint
	Zenny      uint
	ChipFolder [FolderSize]ChipInfo
}

// New returns player data with default values
func New() *Player {
	return &Player{
		HP:        defaultHP,
		HPMax:     defaultHP,
		ShotPower: defaultShotPower,
		Zenny:     0,
		ChipFolder: [FolderSize]ChipInfo{
			{ID: chip.IDSword, Code: "a"},
			{ID: chip.IDWideSword, Code: "a"},
			{ID: chip.IDLongSword, Code: "a"},
			{ID: chip.IDCannon, Code: "b"},
			{ID: chip.IDCannon, Code: "b"},
			{ID: chip.IDCannon, Code: "b"},
			{ID: chip.IDCannon, Code: "c"},
			{ID: chip.IDCannon, Code: "c"},
			{ID: chip.IDCannon, Code: "*"},
			{ID: chip.IDCannon, Code: "*"},
		},
	}
}

// TODO NewWithSaveData(fname string) (*Player, error)
