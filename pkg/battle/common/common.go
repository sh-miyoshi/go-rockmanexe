package common

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/common"
)

// MoveObject ...
func MoveObject(x, y *int, direct int, isMove bool) bool {
	nx := *x
	ny := *y

	switch direct {
	case common.DirectUp:
		if ny <= 0 {
			return false
		}
		ny--
	case common.DirectDown:
		if ny >= field.FieldNumY-1 {
			return false
		}
		ny++
	case common.DirectLeft:
		if nx <= 0 {
			return false
		}
		nx--
	case common.DirectRight:
		if nx >= field.FieldNumX-1 {
			return false
		}
		nx++
	}

	// TODO field panel is player?

	if isMove {
		*x = nx
		*y = ny
	}

	return true
}

func ViewPos(x, y int) (viewX, viewY int32) {
	viewX = int32(field.PanelSizeX*x + field.PanelSizeX/2)
	viewY = int32(field.DrawPanelTopY + field.PanelSizeY*y - 10)
	return viewX, viewY
}
