package common

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

func MoveObject(pos *point.Point, direct int, objPanelType int, isMove bool, GetPanelInfo func(pos point.Point) PanelInfo) bool {
	next := *pos

	// Check field out
	switch direct {
	case config.DirectUp:
		next.Y--
	case config.DirectDown:
		next.Y++
	case config.DirectLeft:
		next.X--
	case config.DirectRight:
		next.X++
	}

	return MoveObjectDirect(pos, next, objPanelType, isMove, GetPanelInfo)
}

func MoveObjectDirect(pos *point.Point, target point.Point, objPanelType int, isMove bool, GetPanelInfo func(pos point.Point) PanelInfo) bool {
	if target.X < 0 || target.Y < 0 || target.X >= FieldNum.X || target.Y >= FieldNum.Y {
		return false
	}

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

func ReverseDirect(direct int) int {
	switch direct {
	case config.DirectLeft:
		return config.DirectRight
	case config.DirectRight:
		return config.DirectLeft
	}
	return direct
}

var (
	playerImageNums   = [PlayerActMax]int{4, 6, 6, 6, 7, 7, 6, 6, 4, 4}
	playerImageDelays = [PlayerActMax]int{1, 2, 2, 6, 3, 4, 1, 4, 4, 2}

	playerBShotNumDelays = map[resources.SoulUnison][2]int{
		resources.SoulUnisonNone: {6, 1},
		resources.SoulUnisonAqua: {7, 2},
	}
)

func GetPlayerActCount(soulUnison resources.SoulUnison, actType int, keepCount int) int {
	if actType < 0 || actType >= PlayerActMax {
		return 1
	}
	if actType == PlayerActBShot {
		return playerBShotNumDelays[soulUnison][0] * (playerBShotNumDelays[soulUnison][1] + keepCount)
	}
	return playerImageDelays[actType] * (playerImageNums[actType] + keepCount)
}

func GetPlayerImageInfo(soulUnison resources.SoulUnison, actType int) (num, delay int) {
	if actType == PlayerActBShot {
		return playerBShotNumDelays[soulUnison][0], playerBShotNumDelays[soulUnison][1]
	}

	return playerImageNums[actType], playerImageDelays[actType]
}
