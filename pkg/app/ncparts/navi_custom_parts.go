package ncparts

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
)

// Memo: 色情報は固定で持たせる。同じ種類で色が異なる場合は別IDを割り当てる

const (
	ColorWhite int = iota
	ColorYellow
	ColorPink
)

const (
	IDAttack1 int = iota
	IDCharge1
	IDHP50
	IDHP100
	IDCustom1
	IDUnderShirt
	IDQuickGauge
)

type NaviCustomParts struct {
	ID          int
	Name        string
	Blocks      []common.Point
	IsPlusParts bool
	Color       int
	// TODO: 効果
}

var (
	allParts = []NaviCustomParts{
		{
			ID:   IDAttack1,
			Name: "アタック+1",
			Blocks: []common.Point{
				{X: 0, Y: 0},
				{X: 1, Y: 0},
			},
			IsPlusParts: true,
			Color:       ColorPink,
		},
		{
			ID:   IDUnderShirt,
			Name: "アンダーシャツ",
			Blocks: []common.Point{
				{X: 0, Y: 0},
				{X: 0, Y: 1},
			},
			IsPlusParts: false,
			Color:       ColorWhite,
		},
	}
)

func Get(id int) NaviCustomParts {
	for _, parts := range allParts {
		if parts.ID == id {
			return parts
		}
	}

	common.SetError(fmt.Sprintf("Navi parts %d is not implemented yet", id))
	return NaviCustomParts{}
}

func GetColorCode(color int) uint {
	switch color {
	case ColorWhite:
		return 0xDCD8DC
	case ColorYellow:
		return 0xDCD800
	case ColorPink:
		return 0xDC88C4
	}

	common.SetError(fmt.Sprintf("Color code %d is not implemented yet", color))
	return 0
}
