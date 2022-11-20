package menu

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/enemy"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/inputs"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
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

func (t *menuDevFeature) Process() error {
	if inputs.CheckKey(inputs.KeyEnter) == 1 {
		sound.On(sound.SEMenuEnter)
		switch t.pointer {
		case devFeatureSelectMapMove:
			return ErrGoMap
		case devFeatureSelectWideArea:
			field.Set4x4Area()

			// Set enemy info
			specificEnemy = []enemy.EnemyParam{
				{
					CharID: enemy.IDTarget,
					Pos:    common.Point{X: 6, Y: 2},
					HP:     1000,
				},
			}

			return ErrGoBattle
		case devFeatureSelectSupportNPC:
			// Set enemy info
			specificEnemy = []enemy.EnemyParam{
				{
					CharID: enemy.IDTarget,
					Pos:    common.Point{X: 4, Y: 1},
					HP:     1000,
				},
			}

			// TODO(味方キャラの追加)
		}
		return nil
	}
	if inputs.CheckKey(inputs.KeyUp) == 1 {
		if t.pointer > 0 {
			sound.On(sound.SECursorMove)
			t.pointer--
		}
	} else if inputs.CheckKey(inputs.KeyDown) == 1 {
		if t.pointer < devFeatureSelectMax-1 {
			sound.On(sound.SECursorMove)
			t.pointer++
		}
	}
	return nil
}
