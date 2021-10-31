package chip

import (
	"fmt"
	"io/ioutil"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/stretchr/stew/slice"
	yaml "gopkg.in/yaml.v2"
)

// Chip ...
type Chip struct {
	ID        int    `yaml:"id"`
	Name      string `yaml:"name"`
	Power     uint   `yaml:"power"`
	Type      int    `yaml:"type"`
	PlayerAct int    `yaml:"player_act"`
	ForMe     bool   `yaml:"for_me"`
	KeepCount int    `yaml:"keep_cnt"`

	Image int32
}

type SelectParam struct {
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
	IDThunderBall = 15
	IDWideShot    = 18
	IDMiniBomb    = 44
	IDSword       = 54
	IDWideSword   = 55
	IDLongSword   = 56
	IDBoomerang1  = 69
	IDRecover10   = 109
	IDRecover30   = 110
	IDShockWave   = 229
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

	typeMax
)

var (
	imgIcons     map[int]int32
	imgMonoIcons map[int]int32
	imgTypes     []int32
	chipData     []Chip
)

// Init ...
func Init(fname string) error {
	// Load chip data
	buf, err := ioutil.ReadFile(fname)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(buf, &chipData); err != nil {
		return err
	}

	for i, c := range chipData {
		fname := fmt.Sprintf("%schipInfo/detail/%d.png", common.ImagePath, c.ID)
		chipData[i].Image = dxlib.LoadGraph(fname)
		if chipData[i].Image == -1 {
			return fmt.Errorf("failed to read chip detail image %s", fname)
		}
	}

	// Load Type Image
	tmp := make([]int32, 14)
	fname = common.ImagePath + "chipInfo/chip_type.png"
	if res := dxlib.LoadDivGraph(fname, 14, 7, 2, 28, 28, tmp); res == -1 {
		return fmt.Errorf("failed to read chip type image %s", fname)
	}
	imgTypes = make([]int32, typeMax)
	for i := 0; i < typeMax; i++ {
		imgTypes[i] = tmp[i]
	}
	for i := typeMax; i < 14; i++ {
		dxlib.DeleteGraph(tmp[i])
	}

	// Load Icon Image
	tmp = make([]int32, 240)
	tmp2 := make([]int32, 240)
	fname = common.ImagePath + "chipInfo/chip_icon.png"
	if res := dxlib.LoadDivGraph(fname, 240, 30, 8, 28, 28, tmp); res == -1 {
		return fmt.Errorf("failed to read chip icon image %s", fname)
	}
	fname = common.ImagePath + "chipInfo/chip_icon_mono.png"
	if res := dxlib.LoadDivGraph(fname, 240, 30, 8, 28, 28, tmp2); res == -1 {
		return fmt.Errorf("failed to read chip monochro icon image %s", fname)
	}

	imgIcons = make(map[int]int32)
	imgMonoIcons = make(map[int]int32)
	used := []int{}

	// Set icons by manual
	for _, c := range chipData {
		// tmp and tmp2 start with 0, but chip id start with 1
		imgIcons[c.ID] = tmp[c.ID-1]
		imgMonoIcons[c.ID] = tmp2[c.ID-1]
		used = append(used, c.ID-1)
	}

	// Release unused images
	for i := 0; i < 240; i++ {
		if !slice.Contains(used, i) {
			dxlib.DeleteGraph(tmp[i])
			dxlib.DeleteGraph(tmp2[i])
		}
	}

	return nil
}

// Get ...
func Get(id int) Chip {
	for _, c := range chipData {
		if c.ID == id {
			return c
		}
	}

	panic(fmt.Sprintf("No such chip ID %d in list %+v", id, chipData))
}

func GetByName(name string) Chip {
	for _, c := range chipData {
		if c.Name == name {
			return c
		}
	}

	panic(fmt.Sprintf("No such chip %s in list %+v", name, chipData))
}

// GetIcon ...
func GetIcon(id int, colored bool) int32 {
	if colored {
		return imgIcons[id]
	}
	return imgMonoIcons[id]
}

// GetTypeImage ...
func GetTypeImage(typ int) int32 {
	return imgTypes[typ]
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
