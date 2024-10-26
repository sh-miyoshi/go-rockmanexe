package chip

import (
	"fmt"
	"os"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/system"
	yaml "gopkg.in/yaml.v2"
)

// Chip ...
type Chip struct {
	ID            int    `yaml:"id"`
	Name          string `yaml:"name"`
	Power         uint   `yaml:"power"`
	Type          int    `yaml:"type"`
	PlayerAct     int    `yaml:"player_act"`
	ForMe         bool   `yaml:"for_me"`
	KeepCount     int    `yaml:"keep_cnt"`
	Description   string `yaml:"description"`
	IsImplemented bool   `yaml:"is_implemented"`
	IconIndex     int    `yaml:"icon_index"`

	IsProgramAdvance bool
}

type SelectParam struct {
	ID   int
	Name string
	Code string
}

const (
	// ID must be same as in chipList.yaml

	IDCannon       = 1
	IDHighCannon   = 2
	IDMegaCannon   = 3
	IDVulcan1      = 5
	IDVulcan2      = 6
	IDVulcan3      = 7
	IDSpreadGun    = 8
	IDHeatShot     = 9
	IDHeatV        = 10
	IDHeatSide     = 11
	IDBubbleShot   = 12
	IDBubbleV      = 13
	IDBubbleSide   = 14
	IDThunderBall1 = 15
	IDThunderBall2 = 16
	IDThunderBall3 = 17
	IDWideShot1    = 18
	IDWideShot2    = 19
	IDWideShot3    = 20
	IDFlameLine1   = 21
	IDFlameLine2   = 22
	IDFlameLine3   = 23
	IDTornado      = 42
	IDMiniBomb     = 44
	IDSword        = 54
	IDWideSword    = 55
	IDLongSword    = 56
	IDBoomerang1   = 69
	IDBambooLance  = 75
	IDCountBomb    = 93
	IDCrackout     = 106
	IDDoubleCrack  = 107
	IDTripleCrack  = 108
	IDRecover10    = 109
	IDRecover30    = 110
	IDRecover50    = 111
	IDRecover80    = 112
	IDRecover120   = 113
	IDRecover150   = 114
	IDRecover200   = 115
	IDRecover300   = 116
	IDPanelSteal   = 118
	IDAreaSteal    = 119
	IDPanelReturn  = 123
	IDDeathMatch1  = 124
	IDDeathMatch2  = 125
	IDQuickGauge   = 128
	IDInvisible    = 133
	IDAttack10     = 148
	IDAquaman      = 181
	IDForteAnother = 236
	IDShockWave    = 247

	// Program Advance
	IDPAIndex    = 1000
	IDDreamSword = IDPAIndex + 1
)

const (
	TypeNone = iota
	TypeWind
	TypeBreaking
	TypeSword
	TypeCracking
	TypeObstacle
	TypeRecovery
	TypeInvis
	TypePlus
	TypeFire
	TypeWater
	TypeElec
	TypeWood

	TypeMax
)

var (
	chipData []Chip
)

func Init(fname string) error {
	// Load chip data
	buf, err := os.ReadFile(fname)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(buf, &chipData); err != nil {
		return err
	}

	// Load Program Advance chips
	setPAData()

	return nil
}

func GetIDList() []int {
	res := []int{}
	for _, c := range chipData {
		if c.IsImplemented {
			res = append(res, c.ID)
		}
	}
	return res
}

func GetList() []Chip {
	res := []Chip{}
	for _, c := range chipData {
		if c.IsImplemented {
			res = append(res, c)
		}
	}
	return res
}

func Get(id int) Chip {
	for _, c := range chipData {
		if c.ID == id {
			if !c.IsImplemented {
				system.SetError(fmt.Sprintf("Chip ID %d is not implemented", id))
			}
			return c
		}
	}

	system.SetError(fmt.Sprintf("No such chip ID %d in list %+v", id, chipData))
	return Chip{}
}

func GetByName(name string) Chip {
	for _, c := range chipData {
		if c.Name == name {
			if !c.IsImplemented {
				system.SetError(fmt.Sprintf("Chip ID %d is not implemented", c.ID))
			}
			return c
		}
	}

	system.SetError(fmt.Sprintf("No such chip %s in list %+v", name, chipData))
	return Chip{}
}

func Selectable(target SelectParam, currentList []SelectParam) bool {
	name := target.Name
	code := target.Code
	for _, c := range currentList {
		if c.Name != name {
			name = "-"
		}
		if c.Code != code {
			// If either c.Code or code is "*", set the other code
			// Both of c.Code and code is not "*", set "-"
			if code == "*" {
				code = c.Code
			} else if c.Code != "*" {
				code = "-"
			}
		}
	}

	if name != "-" || code != "-" {
		return true
	}

	return false
}
