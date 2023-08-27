package naviparts

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
	IDRapid1
	IDCharge1
	IDHP50
	IDHP100
	IDCustom1
	IDUnderShirt
	IDQuickGauge
)

type NaviParts struct {
	ID          int
	Name        string
	Blocks      []common.Point
	IsPlusParts bool
	Color       int
	// TODO: 効果
}

var (
	allParts = []NaviParts{
		{
			ID:   IDAttack1,
			Name: "アタック+1",
			Blocks: []common.Point{
				{X: 0, Y: 0},
				{X: 0, Y: 1},
			},
			IsPlusParts: true,
			Color:       ColorPink,
		},
	}
)

func Get(id int) NaviParts {
	for _, parts := range allParts {
		if parts.ID == id {
			return parts
		}
	}

	common.SetError(fmt.Sprintf("Navi parts %d is not implemented yet", id))
	return NaviParts{}
}
