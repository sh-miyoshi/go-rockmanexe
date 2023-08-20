package menu

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/enemy"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/inputs"
)

const (
	devFeatureSelectMapMove = iota
	devFeatureSelectWideArea
	devFeatureSelectSupportNPC

	devFeatureSelectMax
)

type menuDevFeature struct {
	pointer int
}

func devFeatureNew() (*menuDevFeature, error) {
	return &menuDevFeature{
		pointer: 0,
	}, nil
}

func (t *menuDevFeature) End() {
}

func (t *menuDevFeature) Draw() {
	msgs := []string{
		"マップ移動",
		"4x4 対戦",
		"味方NPC",
	}

	dxlib.DrawBox(20, 30, 230, 300, dxlib.GetColor(168, 192, 216), true)
	dxlib.DrawBox(30, 40, 210, len(msgs)*35+50, dxlib.GetColor(16, 80, 104), true)

	for i, msg := range msgs {
		draw.String(65, 50+i*35, 0xffffff, msg)
	}

	const s = 2
	y := 50 + t.pointer*35
	dxlib.DrawTriangle(40, y+s, 40+18-s*2, y+10, 40, y+20-s, 0xffffff, true)
}

func (t *menuDevFeature) Process() (Result, error) {
	// 隠しコマンド
	if inputs.CheckKey(inputs.KeyDebug) == 1 {
		return ResultGoScratch, nil
	}

	if inputs.CheckKey(inputs.KeyEnter) == 1 {
		sound.On(resources.SEMenuEnter)
		switch t.pointer {
		case devFeatureSelectMapMove:
			return ResultGoMap, nil
		case devFeatureSelectWideArea:
			field.Set4x4Area()

			// Set enemy info
			specificEnemy = []enemy.EnemyParam{
				{
					CharID: enemy.IDAquaman,
					Pos:    common.Point{X: 6, Y: 2},
					HP:     1000,
				},
			}

			return ResultGoBattle, nil
		case devFeatureSelectSupportNPC:
			field.Set4x4Area()

			// Set enemy info
			specificEnemy = []enemy.EnemyParam{
				{
					CharID: enemy.IDAquaman,
					Pos:    common.Point{X: 6, Y: 2},
					HP:     1000,
				},
				{
					CharID: enemy.IDSupportNPC,
					Pos:    common.Point{X: 0, Y: 0},
					HP:     100,
				},
			}

			return ResultGoBattle, nil
		}
		return ResultContinue, nil
	}
	if inputs.CheckKey(inputs.KeyUp) == 1 {
		if t.pointer > 0 {
			sound.On(resources.SECursorMove)
			t.pointer--
		}
	} else if inputs.CheckKey(inputs.KeyDown) == 1 {
		if t.pointer < devFeatureSelectMax-1 {
			sound.On(resources.SECursorMove)
			t.pointer++
		}
	}
	return ResultContinue, nil
}
