package menu

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/enemy"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/list"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/inputs"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	devFeatureSelectMapMove = iota
	devFeatureSelectWideArea
	devFeatureSelectSupportNPC
	devFeatureSelectTalkAI
)

type menuDevFeature struct {
	itemList list.ItemList
	result   Result
}

func devFeatureNew() (*menuDevFeature, error) {
	res := &menuDevFeature{
		result: ResultContinue,
	}
	res.itemList.SetList([]string{
		"マップ移動",
		"4x4 対戦",
		"味方NPC",
		"AIと会話",
	}, -1)

	return res, nil
}

func (t *menuDevFeature) End() {
}

func (t *menuDevFeature) Draw() {
	dxlib.DrawBox(20, 30, 230, 300, dxlib.GetColor(168, 192, 216), true)
	dxlib.DrawBox(30, 40, 210, len(t.itemList.GetList())*35+50, dxlib.GetColor(16, 80, 104), true)

	for i, msg := range t.itemList.GetList() {
		draw.String(65, 50+i*35, 0xffffff, msg)
	}

	const s = 2
	y := 50 + t.itemList.GetPointer()*35
	dxlib.DrawTriangle(40, y+s, 40+18-s*2, y+10, 40, y+20-s, 0xffffff, true)
}

func (t *menuDevFeature) Update() bool {
	// 隠しコマンド
	if inputs.CheckKey(inputs.KeyDebug) == 1 {
		t.result = ResultGoScratch
		return true
	}

	sel := t.itemList.Update()
	if sel != -1 {
		sound.On(resources.SEMenuEnter)
		switch t.itemList.GetPointer() {
		case devFeatureSelectMapMove:
			t.result = ResultGoMap
			return true
		case devFeatureSelectWideArea:
			field.Set4x4Area()

			// Set enemy info
			battleEnemies = []enemy.EnemyParam{
				{
					CharID: enemy.IDAquaman,
					Pos:    point.Point{X: 6, Y: 2},
					HP:     1000,
				},
			}

			t.result = ResultGoBattle
			return true
		case devFeatureSelectSupportNPC:
			field.Set4x4Area()

			// Set enemy info
			battleEnemies = []enemy.EnemyParam{
				{
					CharID: enemy.IDAquaman,
					Pos:    point.Point{X: 6, Y: 2},
					HP:     1000,
				},
				{
					CharID: enemy.IDSupportNPC,
					Pos:    point.Point{X: 0, Y: 0},
					HP:     100,
				},
			}

			t.result = ResultGoBattle
			return true
		case devFeatureSelectTalkAI:
			t.result = ResultGoTalkAI
			return true
		}
	}
	return false
}

func (t *menuDevFeature) GetResult() Result {
	return t.result
}
