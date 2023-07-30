package common

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
)

// MoveObject ...
func MoveObject(pos *common.Point, direct int, objPanelType int, isMove bool, GetPanelInfo func(pos common.Point) PanelInfo) bool {
	next := *pos

	// Check field out
	switch direct {
	case common.DirectUp:
		if next.Y <= 0 {
			return false
		}
		next.Y--
	case common.DirectDown:
		if next.Y >= FieldNum.Y-1 {
			return false
		}
		next.Y++
	case common.DirectLeft:
		if next.X <= 0 {
			return false
		}
		next.X--
	case common.DirectRight:
		if next.X >= FieldNum.X-1 {
			return false
		}
		next.X++
	}

	return MoveObjectDirect(pos, next, objPanelType, isMove, GetPanelInfo)
}

func MoveObjectDirect(pos *common.Point, target common.Point, objPanelType int, isMove bool, GetPanelInfo func(pos common.Point) PanelInfo) bool {
	pn := GetPanelInfo(target)
	// Object exists?
	if pn.ObjectID != "" {
		return false
	}
	// Panel owner?
	if objPanelType >= 0 && objPanelType != pn.Type {
		return false
	}
	// Panel Status
	if pn.Status == PanelStatusHole {
		return false
	}

	if isMove {
		*pos = target
	}

	return true
}

func ViewPos(pos common.Point) common.Point {
	return common.Point{
		X: PanelSize.X*pos.X + PanelSize.X/2,
		Y: DrawPanelTopY + PanelSize.Y*pos.Y - 10,
	}
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
