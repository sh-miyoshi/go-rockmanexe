package ncparts

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/system"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

// Memo: 色情報は固定で持たせる。同じ種類で色が異なる場合は別IDを割り当てる

const (
	ColorWhite int = iota
	ColorYellow
	ColorPink
	ColorBlue
)

const (
	IDAttack1_Pink int = iota
	IDAttack1_White
	IDCharge1_Yellow
	IDCharge1_White
	IDHP50_White
	IDHP100_Yellow
	IDCustom1
	IDUnderShirt
)

type NaviCustomParts struct {
	ID          int
	Name        string
	Blocks      []point.Point
	IsPlusParts bool
	Color       int
	Description string
}

var (
	allParts = []NaviCustomParts{
		{
			ID:   IDAttack1_Pink,
			Name: "アタック+1",
			Blocks: []point.Point{
				{X: 0, Y: 0},
				{X: 0, Y: 1},
			},
			IsPlusParts: true,
			Color:       ColorPink,
			Description: "ロックバスターの威力を上げる",
		},
		{
			ID:   IDAttack1_White,
			Name: "アタック+1",
			Blocks: []point.Point{
				{X: 0, Y: 0},
				{X: 0, Y: 1},
			},
			IsPlusParts: true,
			Color:       ColorWhite,
			Description: "ロックバスターの威力を上げる",
		},
		{
			ID:   IDCharge1_Yellow,
			Name: "チャージ+1",
			Blocks: []point.Point{
				{X: 0, Y: 0},
			},
			IsPlusParts: true,
			Color:       ColorYellow,
			Description: "ロックバスターのチャージ速度を上げる",
		},
		{
			ID:   IDCharge1_White,
			Name: "チャージ+1",
			Blocks: []point.Point{
				{X: 0, Y: 0},
			},
			IsPlusParts: true,
			Color:       ColorWhite,
			Description: "ロックバスターのチャージ速度を上げる",
		},
		{
			ID:   IDHP50_White,
			Name: "ＨＰ+50",
			Blocks: []point.Point{
				{X: 0, Y: 0},
				{X: 1, Y: 0},
			},
			IsPlusParts: true,
			Color:       ColorWhite,
			Description: "最大ＨＰを＋５０する",
		},
		{
			ID:   IDHP100_Yellow,
			Name: "ＨＰ+100",
			Blocks: []point.Point{
				{X: 0, Y: 0},
				{X: 1, Y: 0},
				{X: 1, Y: 1},
			},
			IsPlusParts: true,
			Color:       ColorYellow,
			Description: "最大ＨＰを＋１００する",
		},
		{
			ID:   IDUnderShirt,
			Name: "アンダーシャツ",
			Blocks: []point.Point{
				{X: 0, Y: 0},
				{X: 0, Y: 1},
			},
			IsPlusParts: false,
			Color:       ColorWhite,
			Description: "HPが0になる前に1で耐える",
		},
		{
			ID:   IDCustom1,
			Name: "カスタム１",
			Blocks: []point.Point{
				{X: 0, Y: 0},
				{X: 1, Y: 0},
				{X: 0, Y: 1},
				{X: 0, Y: 2},
				{X: -1, Y: 1},
			},
			IsPlusParts: false,
			Color:       ColorBlue,
			Description: "カスタム画面のチップが一枚増える",
		},
	}
)

func Get(id int) NaviCustomParts {
	for _, parts := range allParts {
		if parts.ID == id {
			return parts
		}
	}

	system.SetError(fmt.Sprintf("Navi parts %d is not implemented yet", id))
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
	case ColorBlue:
		return 0x2880DC
	}

	system.SetError(fmt.Sprintf("Color code %d is not implemented yet", color))
	return 0
}
