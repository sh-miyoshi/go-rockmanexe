package common

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
)

// MoveObject ...
func MoveObject(x, y *int, direct int, objPanelType int, isMove bool, GetPanelInfo func(x, y int) field.PanelInfo) bool {
	nx := *x
	ny := *y

	// Check field out
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

	pn := GetPanelInfo(nx, ny)
	// Object exists?
	if pn.ObjectID != "" {
		return false
	}
	// Panel owner?
	if objPanelType != pn.Type {
		return false
	}
	// Panel Status
	if pn.Status == field.PanelStatusHole {
		return false
	}

	if isMove {
		*x = nx
		*y = ny
	}

	return true
}

func MoveObjectDirect(x, y *int, targetX, targetY int, objPanelType int, isMove bool, GetPanelInfo func(x, y int) field.PanelInfo) bool {
	pn := GetPanelInfo(targetX, targetY)
	if pn.ObjectID != "" {
		return false
	}
	// Check panel type
	if objPanelType != pn.Type {
		return false
	}

	if isMove {
		*x = targetX
		*y = targetY
	}

	return true
}

func ViewPos(x, y int) (viewX, viewY int32) {
	viewX = int32(field.PanelSizeX*x + field.PanelSizeX/2)
	viewY = int32(field.DrawPanelTopY + field.PanelSizeY*y - 10)
	return viewX, viewY
}

func GetOffset(nextPos, nowPos, beforePos int, cnt, totalCnt int, size int) int {
	// if cnt < total_count/2
	//   init_offset = (before - now) * size / 2
	//   offset = init_offset - (before - now)*(count*size/total_count))

	var res int
	if cnt < totalCnt/2 {
		res = (beforePos - nowPos)
	} else {
		res = (nowPos - nextPos)
	}

	return res * size * (totalCnt - 2*cnt) / (totalCnt * 2)
}
