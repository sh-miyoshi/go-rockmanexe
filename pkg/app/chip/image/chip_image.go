package chipimage

import (
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
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
		if id >= chip.IDPAIndex {
			continue
		}

		fname := fmt.Sprintf("%schipInfo/detail/%d.png", config.ImagePath, id)
		imgDetails[id] = dxlib.LoadGraph(fname)
		if imgDetails[id] == -1 {
			return errors.Newf("failed to read chip detail image %s", fname)
		}
	}

	// Load Type Image
	tmp := make([]int, 14)
	fname := config.ImagePath + "chipInfo/chip_type.png"
	if res := dxlib.LoadDivGraph(fname, 14, 7, 2, 28, 28, tmp); res == -1 {
		return errors.Newf("failed to read chip type image %s", fname)
	}
	imgTypes = make([]int, chip.TypeMax)
	for i := 0; i < chip.TypeMax; i++ {
		imgTypes[i] = tmp[i]
	}
	for i := chip.TypeMax; i < 14; i++ {
		dxlib.DeleteGraph(tmp[i])
	}

	// Load Icon Image
	imgIcons = make(map[int]int)
	imgMonoIcons = make(map[int]int)

	tmp = make([]int, 230)
	tmp2 := make([]int, 230)
	fname = config.ImagePath + "chipInfo/chip_icon.png"
	if res := dxlib.LoadDivGraph(fname, 230, 30, 8, 28, 28, tmp); res == -1 {
		return errors.Newf("failed to read chip icon image %s", fname)
	}

	fname = config.ImagePath + "chipInfo/chip_icon_mono.png"
	if res := dxlib.LoadDivGraph(fname, 230, 30, 8, 28, 28, tmp2); res == -1 {
		return errors.Newf("failed to read chip monochro icon image %s", fname)
	}
	for _, c := range chip.GetList() {
		index := c.ID - 1
		if c.IconIndex > 0 {
			index = c.IconIndex
		}

		imgIcons[c.ID] = tmp[index]
		imgMonoIcons[c.ID] = tmp2[index]
	}

	fname = config.ImagePath + "chipInfo/pa_icon.png"
	imgIcons[chip.IDPAIndex] = dxlib.LoadGraph(fname)
	if imgIcons[chip.IDPAIndex] == -1 {
		return errors.Newf("failed to load image %s", fname)
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
