package common

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

// MoveObject ...
func MoveObject(pos *point.Point, direct int, objPanelType int, isMove bool, GetPanelInfo func(pos point.Point) PanelInfo) bool {
	next := *pos

	// Check field out
	switch direct {
	case config.DirectUp:
		if next.Y <= 0 {
			return false
		}
		next.Y--
	case config.DirectDown:
		if next.Y >= FieldNum.Y-1 {
			return false
		}
		next.Y++
	case config.DirectLeft:
		if next.X <= 0 {
			return false
		}
		next.X--
	case config.DirectRight:
		if next.X >= FieldNum.X-1 {
			return false
		}
		next.X++
	}

	return MoveObjectDirect(pos, next, objPanelType, isMove, GetPanelInfo)
}

func MoveObjectDirect(pos *point.Point, target point.Point, objPanelType int, isMove bool, GetPanelInfo func(pos point.Point) PanelInfo) bool {
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

func ViewPos(pos point.Point) point.Point {
	return point.Point{
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
