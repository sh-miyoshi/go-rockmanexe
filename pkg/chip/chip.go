package chip

import (
	"fmt"
	"io/ioutil"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/common"
	"github.com/stretchr/stew/slice"
	yaml "gopkg.in/yaml.v2"
)

// Chip ...
type Chip struct {
	ID        int    `yaml:"id"`
	Name      string `yaml:"name"`
	Power     uint   `yaml:"power"`
	Type      int    `yaml:"type"`
	Code      string `yaml:"code"`
	PlayerAct int    `yaml:"player_act"`

	Image int32
}

const (
	// IDCannon ...
	IDCannon = iota

	idMax
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
	imgIcons     []int32
	imgMonoIcons []int32
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

	for i := 0; i < idMax; i++ {
		fname := fmt.Sprintf("%schipInfo/detail/%d.png", common.ImagePath, i)
		chipData[i].Image = dxlib.LoadGraph(fname)
		if chipData[i].Image == -1 {
			return fmt.Errorf("Failed to read chip detail image %s", fname)
		}
	}

	// Load Type Image
	tmp := make([]int32, 14)
	fname = common.ImagePath + "chipInfo/chip_type.png"
	if res := dxlib.LoadDivGraph(fname, 14, 7, 2, 28, 28, tmp); res == -1 {
		return fmt.Errorf("Failed to read chip type image %s", fname)
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
		return fmt.Errorf("Failed to read chip icon image %s", fname)
	}
	fname = common.ImagePath + "chipInfo/chip_icon_mono.png"
	if res := dxlib.LoadDivGraph(fname, 240, 30, 8, 28, 28, tmp2); res == -1 {
		return fmt.Errorf("Failed to read chip monochro icon image %s", fname)
	}

	imgIcons = make([]int32, idMax)
	imgMonoIcons = make([]int32, idMax)
	used := []int{}

	// Set icons by manual
	used = append(used, 0)
	imgIcons[IDCannon] = tmp[0]
	imgMonoIcons[IDCannon] = tmp2[0]

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
	return chipData[id]
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
