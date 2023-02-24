package chipimage

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/stretchr/stew/slice"
)

var (
	imgDetails   map[int]int
	imgIcons     map[int]int
	imgMonoIcons map[int]int
	imgTypes     []int
)

func Init() error {
	// Load Chip Image
	imgDetails = make(map[int]int)
	for _, id := range chip.GetIDList() {
		fname := fmt.Sprintf("%schipInfo/detail/%d.png", common.ImagePath, id)
		imgDetails[id] = dxlib.LoadGraph(fname)
		if imgDetails[id] == -1 {
			return fmt.Errorf("failed to read chip detail image %s", fname)
		}
	}

	// Load Type Image
	tmp := make([]int, 14)
	fname := common.ImagePath + "chipInfo/chip_type.png"
	if res := dxlib.LoadDivGraph(fname, 14, 7, 2, 28, 28, tmp); res == -1 {
		return fmt.Errorf("failed to read chip type image %s", fname)
	}
	imgTypes = make([]int, chip.TypeMax)
	for i := 0; i < chip.TypeMax; i++ {
		imgTypes[i] = tmp[i]
	}
	for i := chip.TypeMax; i < 14; i++ {
		dxlib.DeleteGraph(tmp[i])
	}

	// Load Icon Image
	tmp = make([]int, 240)
	tmp2 := make([]int, 240)
	fname = common.ImagePath + "chipInfo/chip_icon.png"
	if res := dxlib.LoadDivGraph(fname, 240, 30, 8, 28, 28, tmp); res == -1 {
		return fmt.Errorf("failed to read chip icon image %s", fname)
	}
	fname = common.ImagePath + "chipInfo/chip_icon_mono.png"
	if res := dxlib.LoadDivGraph(fname, 240, 30, 8, 28, 28, tmp2); res == -1 {
		return fmt.Errorf("failed to read chip monochro icon image %s", fname)
	}

	imgIcons = make(map[int]int)
	imgMonoIcons = make(map[int]int)
	used := []int{}

	// Set icons by manual
	for _, id := range chip.GetIDList() {
		// tmp and tmp2 start with 0, but chip id start with 1
		imgIcons[id] = tmp[id-1]
		imgMonoIcons[id] = tmp2[id-1]
		used = append(used, id-1)
	}
	fname = common.ImagePath + "chipInfo/pa_icon.png"
	imgIcons[chip.IDPAIndex] = dxlib.LoadGraph(fname)
	if imgIcons[chip.IDPAIndex] == -1 {
		return fmt.Errorf("failed to load image %s", fname)
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

func GetDetail(id int) int {
	img, ok := imgDetails[id]
	if ok {
		return img
	}
	return -1
}

func GetIcon(id int, colored bool) int {
	if id > chip.IDPAIndex {
		return imgIcons[chip.IDPAIndex]
	}

	if colored {
		return imgIcons[id]
	}
	return imgMonoIcons[id]
}

func GetType(typ int) int {
	return imgTypes[typ]
}
