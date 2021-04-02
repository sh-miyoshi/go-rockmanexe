package menu

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/enemy"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/inputs"
)

type selectValue struct {
	Name    string
	Enemies []enemy.EnemyParam
}

var (
	goBattleSelectData []selectValue
	goBattleCursor     int
)

func goBattleInit() error {
	goBattleCursor = 0

	goBattleSelectData = []selectValue{
		{
			Name: "まずはここから",
			Enemies: []enemy.EnemyParam{
				{
					CharID: enemy.IDMetall,
					PosX:   4,
					PosY:   1,
					HP:     40,
				},
			},
		},
		{
			Name: "侮ることなかれ",
			Enemies: []enemy.EnemyParam{
				{
					CharID: enemy.IDMetall,
					PosX:   3,
					PosY:   0,
					HP:     40,
				},
				{
					CharID: enemy.IDMetall,
					PosX:   4,
					PosY:   1,
					HP:     40,
				},
				{
					CharID: enemy.IDMetall,
					PosX:   5,
					PosY:   2,
					HP:     40,
				},
			},
		},
		{
			Name: "練習",
			Enemies: []enemy.EnemyParam{
				{
					CharID: enemy.IDTarget,
					PosX:   4,
					PosY:   1,
					HP:     1000,
				},
			},
		},
	}

	return nil
}

func goBattleEnd() {
}

func goBattleProcess() bool {
	if inputs.CheckKey(inputs.KeyCancel) == 1 {
		stateChange(stateTop)
		return false
	}
	if inputs.CheckKey(inputs.KeyEnter) == 1 {
		return true
	}
	if inputs.CheckKey(inputs.KeyUp) == 1 && goBattleCursor > 0 {
		goBattleCursor--
	} else if inputs.CheckKey(inputs.KeyDown) == 1 && goBattleCursor < len(goBattleSelectData)-1 {
		goBattleCursor++
	}

	return false
}

func goBattleDraw() {
	// TODO show name, preview image
	draw.String(common.ScreenX/2-20, common.ScreenY/2-20, 0, "未実装")
}

func battleEnemies() []enemy.EnemyParam {
	return goBattleSelectData[goBattleCursor].Enemies
}
