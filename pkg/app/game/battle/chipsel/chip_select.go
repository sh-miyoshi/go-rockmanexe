package chipsel

import (
	"github.com/cockroachdb/errors"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
	chipimage "github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip/image"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/inputs"
	"github.com/stretchr/stew/slice"
)

const (
	sendBtnNo   = -1
	unisonBtnNo = -2
	rowMax      = 5
)

const (
	selectStateNormal = iota
	selectStateUnison
)

// ChipSelect holds chip selection state.
type ChipSelect struct {
	count       int
	selectList  []player.ChipInfo
	selected    []int
	state       int
	imgFrame    int
	imgPointer  []int
	imgSoulIcon int
	pointer     int
}

// NewChipSelect creates a new ChipSelect instance.
func NewChipSelect() *ChipSelect {
	return &ChipSelect{
		count:       0,
		state:       selectStateNormal,
		imgFrame:    -1,
		imgPointer:  []int{-1, -1},
		imgSoulIcon: -1,
		pointer:     sendBtnNo,
	}
}

// Init initializes the chip selection.
func (c *ChipSelect) Init(folder []player.ChipInfo, chipSelectMax int) error {
	if c.imgFrame == -1 {
		fname := config.ImagePath + "battle/chip_select_frame.png"
		c.imgFrame = dxlib.LoadGraph(fname)
		if c.imgFrame == -1 {
			return errors.Newf("failed to read frame image: %s", fname)
		}
	}

	if c.imgPointer[0] == -1 && c.imgPointer[1] == -1 {
		fname := config.ImagePath + "battle/pointer.png"
		res := dxlib.LoadDivGraph(fname, 2, 2, 1, 44, 44, c.imgPointer)
		if res == -1 {
			return errors.Newf("failed to read pointer image: %s", fname)
		}
	}

	if c.imgSoulIcon == -1 {
		fname := config.ImagePath + "battle/soul_icon.png"
		c.imgSoulIcon = dxlib.LoadGraph(fname)
		if c.imgSoulIcon == -1 {
			return errors.Newf("failed to read soul icon image: %s", fname)
		}
	}

	c.count = 0
	c.selectList = []player.ChipInfo{}
	c.selected = []int{}
	c.state = selectStateNormal

	num := len(folder)
	if num > chipSelectMax {
		num = chipSelectMax
	}
	for i := 0; i < num; i++ {
		c.selectList = append(c.selectList, folder[i])
	}

	c.pointer = sendBtnNo
	if num > 0 {
		c.pointer = 0
	}

	sound.On(resources.SEChipSelectOpen)

	return nil
}

func (c *ChipSelect) End() {
	if c.imgFrame != -1 {
		dxlib.DeleteGraph(c.imgFrame)
		c.imgFrame = -1
	}
	if c.imgPointer[0] != -1 {
		dxlib.DeleteGraph(c.imgPointer[0])
		c.imgPointer[0] = -1
	}
	if c.imgPointer[1] != -1 {
		dxlib.DeleteGraph(c.imgPointer[1])
		c.imgPointer[1] = -1
	}
	if c.imgSoulIcon != -1 {
		dxlib.DeleteGraph(c.imgSoulIcon)
		c.imgSoulIcon = -1
	}
}

// Draw renders the chip selection UI.
func (c *ChipSelect) Draw() {
	if c.imgFrame == -1 {
		// Waiting for initialization.
		return
	}

	baseY := 0
	if field.Is4x4Area() {
		baseY = 20
	}

	dxlib.DrawGraph(0, baseY, c.imgFrame, true)

	// Show chip data.
	for i, s := range c.selectList {
		x := (i%rowMax)*32 + 17
		y := (i / rowMax) * 48
		draw.ChipCode(x+10, y+240+baseY, s.Code, 50)
		if !slice.Contains(c.selected, i) {
			// Show Icon.
			dxlib.DrawGraph(x, y+210+baseY, chipimage.GetIcon(s.ID, c.selectable(i)), true)
		}

		// Show Detail Data.
		if i == c.pointer && c.state == selectStateNormal {
			ch := chip.Get(s.ID)
			dxlib.DrawGraph(31, 64+baseY, chipimage.GetDetail(ch.ID), true)
			dxlib.DrawGraph(52, 161+baseY, chipimage.GetType(ch.Type), true)
			draw.String(20, 25+baseY, 0x000000, "%s", ch.Name)
			draw.ChipCode(30, 163+baseY, s.Code, 100)
			if ch.Power != 0 {
				draw.Number(95, 163+baseY, int(ch.Power), draw.NumberOption{
					Color:        draw.NumberColorWhite,
					Length:       3,
					RightAligned: true,
				})
			}
		}
	}

	// Show pointer.
	if c.state == selectStateNormal {
		n := c.count / 20
		if n%3 != 0 {
			if c.pointer == unisonBtnNo {
				dxlib.DrawGraph(180, 285+baseY, c.imgPointer[0], true)
			} else if c.pointer == sendBtnNo {
				dxlib.DrawGraph(180, 225+baseY, c.imgPointer[1], true)
			} else {
				x := (c.pointer%rowMax)*32 + 8
				y := (c.pointer/rowMax)*48 + 202 + baseY
				dxlib.DrawGraph(x, y, c.imgPointer[0], true)
			}
		}
	} else {
		// WIP
	}

	// Show Selected Chips.
	for i, s := range c.selected {
		y := i*32 + 50
		dxlib.DrawGraph(193, y+baseY, chipimage.GetIcon(c.selectList[s].ID, true), true)
	}
}

// Update handles input and updates the chip selection state.
func (c *ChipSelect) Update() bool {
	c.count++
	max := len(c.selectList)

	switch c.state {
	case selectStateNormal:
		if inputs.CheckKey(inputs.KeyEnter) == 1 {
			if c.pointer == sendBtnNo {
				sound.On(resources.SEChipSelectEnd)
				return true
			}
			if c.pointer == unisonBtnNo {
				sound.On(resources.SEMenuEnter)
				c.state = selectStateUnison
				return false
			} else if c.selectable(c.pointer) {
				sound.On(resources.SESelect)
				c.selected = append(c.selected, c.pointer)
			} else {
				sound.On(resources.SEDenied)
			}
		} else {
			if max == 0 {
				return false
			}

			if inputs.CheckKey(inputs.KeyCancel) == 1 {
				if len(c.selected) > 0 {
					sound.On(resources.SECancel)
					c.selected = c.selected[:len(c.selected)-1]
				}
			} else if inputs.CheckKey(inputs.KeyRight) == 1 {
				sound.On(resources.SECursorMove)
				if c.pointer == rowMax-1 || c.pointer == max-1 {
					c.pointer = sendBtnNo
				} else if c.pointer == sendBtnNo || c.pointer == unisonBtnNo {
					c.pointer = 0
				} else {
					c.pointer++
				}
			} else if inputs.CheckKey(inputs.KeyLeft) == 1 {
				sound.On(resources.SECursorMove)
				if c.pointer == sendBtnNo || c.pointer == unisonBtnNo {
					c.pointer = max - 1
				} else if c.pointer == 0 {
					c.pointer = sendBtnNo
				} else {
					c.pointer--
				}
			} else if inputs.CheckKey(inputs.KeyUp) == 1 {
				if c.pointer == unisonBtnNo {
					sound.On(resources.SECursorMove)
					c.pointer = sendBtnNo
				} else if c.pointer >= rowMax {
					sound.On(resources.SECursorMove)
					c.pointer -= rowMax
				}
			} else if inputs.CheckKey(inputs.KeyDown) == 1 {
				if max > rowMax && c.pointer >= 0 && c.pointer < rowMax {
					sound.On(resources.SECursorMove)
					c.pointer += rowMax
					if c.pointer >= max {
						c.pointer = max - 1
					}
				} else if c.pointer == sendBtnNo {
					sound.On(resources.SECursorMove)
					c.pointer = unisonBtnNo
				}
			}
		}
	case selectStateUnison:
		// WIP
	}

	return false
}

// GetSelected returns the indices of the selected chips.
func (c *ChipSelect) GetSelected() []int {
	return c.selected
}

// selectable determines if a chip at the index can be selected.
func (c *ChipSelect) selectable(no int) bool {
	if slice.Contains(c.selected, no) {
		// already selected
		return false
	}

	ch := chip.Get(c.selectList[no].ID)
	target := chip.SelectParam{
		Name: ch.Name,
		Code: c.selectList[no].Code,
	}
	list := []chip.SelectParam{}
	for _, s := range c.selected {
		ch2 := chip.Get(c.selectList[s].ID)
		list = append(list, chip.SelectParam{
			Name: ch2.Name,
			Code: c.selectList[s].Code,
		})
	}
	return chip.Selectable(target, list)
}
