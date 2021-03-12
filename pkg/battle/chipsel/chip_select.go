package chipsel

import (
	"fmt"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/player"
)

const (
	selectMax = 5 // TODO should variable length
	sendBtnNo = -1
)

var (
	count = 0

	selectList []player.ChipInfo
	selected   []int

	imgFrame   int32 = -1
	imgPointer       = []int32{-1, -1}
	pointer          = sendBtnNo
)

// Init ...
func Init(folder []player.ChipInfo) error {
	if imgFrame == -1 {
		fname := common.ImagePath + "battle/chip_select_frame.png"
		imgFrame = dxlib.LoadGraph(fname)
		if imgFrame == -1 {
			return fmt.Errorf("Failed to read frame image: %s", fname)
		}

		fname = common.ImagePath + "battle/pointer.png"
		res := dxlib.LoadDivGraph(fname, 2, 2, 1, 44, 44, imgPointer)
		if res == -1 {
			return fmt.Errorf("Failed to read pointer image: %s", fname)
		}
	}

	count = 0
	selectList = []player.ChipInfo{}
	selected = []int{}

	num := len(folder)
	if num > selectMax {
		num = selectMax
	}
	for i := 0; i < num; i++ {
		selectList = append(selectList, folder[i])
	}

	pointer = sendBtnNo
	if num > 0 {
		pointer = 0
	}

	return nil
}

// Draw ...
func Draw() {
	dxlib.DrawGraph(0, 0, imgFrame, dxlib.TRUE)

	// Show chip data
	for i, s := range selectList {
		// Show Icon
		// TODO selectable()
		x := i*32 + 17
		dxlib.DrawGraph(int32(x), 210, chip.GetIcon(s.ID, true), dxlib.TRUE)

		// Show Detail Data
		if i == pointer {
			c := chip.Get(s.ID)
			// TODO font
			dxlib.DrawGraph(31, 64, c.Image, dxlib.TRUE)
			dxlib.DrawGraph(52, 161, chip.GetTypeImage(c.Type), dxlib.TRUE)
			dxlib.DrawString(20, 25, c.Name, 0x000000)
			dxlib.DrawFormatString(30, 163, 0xffffff, "%s", s.Code)
		}
	}

	// Show pointer
	n := count / 20
	if n%3 != 0 {
		if pointer == sendBtnNo {
			dxlib.DrawGraph(180, 225, imgPointer[1], dxlib.TRUE)
		} else {
			x := (pointer%5)*32 + 8
			y := (pointer/5)*20 + 202
			dxlib.DrawGraph(int32(x), int32(y), imgPointer[0], dxlib.TRUE)
		}
	}
}

// Process ...
func Process() {
	count++
}

// GetSelected ...
func GetSelected() []int {
	return selected
}
