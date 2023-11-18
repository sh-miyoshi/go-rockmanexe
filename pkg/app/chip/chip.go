package chip

import (
	"fmt"
	"io/ioutil"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	yaml "gopkg.in/yaml.v2"
)

// Chip ...
type Chip struct {
	ID          int    `yaml:"id"`
	Name        string `yaml:"name"`
	Power       uint   `yaml:"power"`
	Type        int    `yaml:"type"`
	PlayerAct   int    `yaml:"player_act"`
	ForMe       bool   `yaml:"for_me"`
	KeepCount   int    `yaml:"keep_cnt"`
	Description string `yaml:"description"`

	IsProgramAdvance bool
}

type SelectParam struct {
	ID   int
	Name string
	Code string
}

const (
	// ID must be same as in chipList.yaml

	IDCannon      = 1
	IDHighCannon  = 2
	IDMegaCannon  = 3
	IDVulcan1     = 5
	IDSpreadGun   = 8
	IDHeatShot    = 9
	IDHeatV       = 10
	IDHeatSide    = 11
	IDThunderBall = 15
	IDWideShot    = 18
	IDFlameLine1  = 21
	IDFlameLine2  = 22
	IDFlameLine3  = 23
	IDTornado     = 42
	IDMiniBomb    = 44
	IDSword       = 54
	IDWideSword   = 55
	IDLongSword   = 56
	IDBoomerang1  = 69
	IDBambooLance = 75
	IDCountBomb   = 93
	IDCrackout    = 106
	IDDoubleCrack = 107
	IDTripleCrack = 108
	IDRecover10   = 109
	IDRecover30   = 110
	IDPanelSteal  = 118
	IDAreaSteal   = 119
	IDQuickGauge  = 128
	IDInvisible   = 133
	IDAttack10    = 148
	IDAquaman     = 211
	IDShockWave   = 229

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
	buf, err := ioutil.ReadFile(fname)
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
		res = append(res, c.ID)
	}
	return res
}

func Get(id int) Chip {
	for _, c := range chipData {
		if c.ID == id {
			return c
		}
	}

	common.SetError(fmt.Sprintf("No such chip ID %d in list %+v", id, chipData))
	return Chip{}
}

func GetByName(name string) Chip {
	for _, c := range chipData {
		if c.Name == name {
			return c
		}
	}

	common.SetError(fmt.Sprintf("No such chip %s in list %+v", name, chipData))
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
