package chipsel

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/inputs"
	"github.com/stretchr/stew/slice"
)

const (
	selectMax = 5 // TODO should variable length
	sendBtnNo = -1
)

var (
	count = 0

	selectList []player.ChipInfo
	selected   []int

	imgFrame   int = -1
	imgPointer     = []int{-1, -1}
	pointer        = sendBtnNo
)

// Init ...
func Init(folder []player.ChipInfo) error {
	if imgFrame == -1 {
		fname := common.ImagePath + "battle/chip_select_frame.png"
		imgFrame = dxlib.LoadGraph(fname)
		if imgFrame == -1 {
			return fmt.Errorf("failed to read frame image: %s", fname)
		}

		fname = common.ImagePath + "battle/pointer.png"
		res := dxlib.LoadDivGraph(fname, 2, 2, 1, 44, 44, imgPointer)
		if res == -1 {
			return fmt.Errorf("failed to read pointer image: %s", fname)
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

	sound.On(sound.SEChipSelectOpen)

	return nil
}

// Draw ...
func Draw() {
	if imgFrame == -1 {
		// Waiting initialize
		return
	}

	baseY := 0
	if field.Is4x4Area() {
		baseY = 20
	}

	dxlib.DrawGraph(0, baseY, imgFrame, true)

	// Show chip data
	for i, s := range selectList {
		x := i*32 + 17
		draw.ChipCode(x+10, 240+baseY, s.Code, 50)
		if !slice.Contains(selected, i) {
			// Show Icon
			dxlib.DrawGraph(x, 210+baseY, chip.GetIcon(s.ID, selectable(i)), true)
		}

		// Show Detail Data
		if i == pointer {
			c := chip.Get(s.ID)
			dxlib.DrawGraph(31, 64+baseY, c.Image, true)
			dxlib.DrawGraph(52, 161+baseY, chip.GetTypeImage(c.Type), true)
			draw.String(20, 25+baseY, 0x000000, "%s", c.Name)
			draw.ChipCode(30, 163+baseY, s.Code, 100)
			if c.Power != 0 {
				draw.Number(95, 163+baseY, int(c.Power), draw.NumberOption{
					Color:        draw.NumberColorWhite,
					Length:       3,
					RightAligned: true,
				})
			}
		}
	}

	// Show pointer
	n := count / 20
	if n%3 != 0 {
		if pointer == sendBtnNo {
			dxlib.DrawGraph(180, 225+baseY, imgPointer[1], true)
		} else {
			x := (pointer%5)*32 + 8
			y := (pointer/5)*20 + 202 + baseY
			dxlib.DrawGraph(x, y, imgPointer[0], true)
		}
	}

	// Show Selected Chips
	for i, s := range selected {
		y := i*32 + 50
		dxlib.DrawGraph(193, y+baseY, chip.GetIcon(selectList[s].ID, true), true)
	}
}

// Process ...
func Process() bool {
	count++
	max := len(selectList)

	if inputs.CheckKey(inputs.KeyEnter) == 1 {
		if pointer == sendBtnNo {
			sound.On(sound.SEChipSelectEnd)
			return true
		}
		if selectable(pointer) {
			sound.On(sound.SESelect)
			selected = append(selected, pointer)
		} else {
			sound.On(sound.SEDenied)
		}
	} else {
		if max == 0 {
			return false
		}

		if inputs.CheckKey(inputs.KeyCancel) == 1 {
			if len(selected) > 0 {
				sound.On(sound.SECancel)
				selected = selected[:len(selected)-1]
			}
		} else if inputs.CheckKey(inputs.KeyRight) == 1 {
			sound.On(sound.SECursorMove)
			if pointer == max-1 {
				pointer = sendBtnNo
			} else if pointer == sendBtnNo {
				pointer = 0
			} else {
				pointer++
			}
		} else if inputs.CheckKey(inputs.KeyLeft) == 1 {
			sound.On(sound.SECursorMove)
			if pointer == sendBtnNo {
				pointer = max - 1
			} else if pointer == 0 {
				pointer = sendBtnNo
			} else {
				pointer--
			}
		}
	}

	return false
}

// GetSelected ...
func GetSelected() []int {
	return selected
}

func selectable(no int) bool {
	if slice.Contains(selected, no) {
		// already selected
		return false
	}

	c := chip.Get(selectList[no].ID)
	target := chip.SelectParam{Name: c.Name, Code: selectList[no].Code}
	list := []chip.SelectParam{}
	for _, s := range selected {
		c := chip.Get(selectList[s].ID)
		list = append(list, chip.SelectParam{Name: c.Name, Code: selectList[s].Code})
	}
	return chip.Selectable(target, list)
}
