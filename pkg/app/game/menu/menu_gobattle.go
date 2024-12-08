package menu

import (
	"github.com/cockroachdb/errors"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/enemy"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/list"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/inputs"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type selectEnemyData struct {
	BattleParam enemy.EnemyParam
	View        point.Point
}

type selectValue struct {
	Name    string
	Enemies []selectEnemyData
}

const (
	goBattleListShowMax = 7
)

var (
	viewCenter = point.Point{X: 350, Y: 150}
)

type menuGoBattle struct {
	result     Result
	selectData []selectValue
	waitCount  int
	images     map[int]int
	itemList   list.ItemList
}

func goBattleNew() (*menuGoBattle, error) {
	res := &menuGoBattle{
		images: make(map[int]int),
		result: ResultContinue,
	}
	res.selectData = []selectValue{
		{
			Name: "千里の道も一歩から",
			Enemies: []selectEnemyData{
				{
					BattleParam: enemy.EnemyParam{
						CharID: enemy.IDMetall,
						Pos:    point.Point{X: 4, Y: 1},
						HP:     40,
					},
					View: point.Point{X: viewCenter.X, Y: viewCenter.Y + 10},
				},
			},
		},
		{
			Name: "侮ることなかれ",
			Enemies: []selectEnemyData{
				{
					BattleParam: enemy.EnemyParam{
						CharID: enemy.IDMetall,
						Pos:    point.Point{X: 3, Y: 0},
						HP:     40,
					},
					View: point.Point{X: viewCenter.X, Y: viewCenter.Y - 30},
				},
				{
					BattleParam: enemy.EnemyParam{
						CharID: enemy.IDMetall,
						Pos:    point.Point{X: 4, Y: 1},
						HP:     40,
					},
					View: point.Point{X: viewCenter.X - 30, Y: viewCenter.Y + 10},
				},
				{
					BattleParam: enemy.EnemyParam{
						CharID: enemy.IDMetall,
						Pos:    point.Point{X: 5, Y: 2},
						HP:     40,
					},
					View: point.Point{X: viewCenter.X + 30, Y: viewCenter.Y + 10},
				},
			},
		},
		{
			Name: "跡追いする電気玉",
			Enemies: []selectEnemyData{
				{
					BattleParam: enemy.EnemyParam{
						CharID: enemy.IDBilly,
						Pos:    point.Point{X: 5, Y: 1},
						HP:     50,
					},
					View: point.Point{X: viewCenter.X + 20, Y: viewCenter.Y - 10},
				},
				{
					BattleParam: enemy.EnemyParam{
						CharID: enemy.IDMetall,
						Pos:    point.Point{X: 3, Y: 0},
						HP:     40,
					},
					View: point.Point{X: viewCenter.X - 30, Y: viewCenter.Y + 10},
				},
			},
		},
		{
			Name: "遊泳するものたち",
			Enemies: []selectEnemyData{
				{
					BattleParam: enemy.EnemyParam{
						CharID: enemy.IDLark,
						Pos:    point.Point{X: 3, Y: 0},
						HP:     60,
					},
					View: point.Point{X: viewCenter.X + 10, Y: viewCenter.Y - 25},
				},
				{
					BattleParam: enemy.EnemyParam{
						CharID: enemy.IDLark,
						Pos:    point.Point{X: 4, Y: 1},
						HP:     60,
					},
					View: point.Point{X: viewCenter.X - 10, Y: viewCenter.Y + 25},
				},
			},
		},
		{
			Name: "舞戻るやいば",
			Enemies: []selectEnemyData{
				{
					BattleParam: enemy.EnemyParam{
						CharID: enemy.IDBoomer,
						Pos:    point.Point{X: 5, Y: 1},
						HP:     60,
					},
					View: point.Point{X: viewCenter.X + 20, Y: viewCenter.Y - 5},
				},
				{
					BattleParam: enemy.EnemyParam{
						CharID: enemy.IDMetall,
						Pos:    point.Point{X: 4, Y: 0},
						HP:     40,
					},
					View: point.Point{X: viewCenter.X - 30, Y: viewCenter.Y + 5},
				},
			},
		},
		{
			Name: "火玉に注意！",
			Enemies: []selectEnemyData{
				{
					BattleParam: enemy.EnemyParam{
						CharID: enemy.IDGaroo,
						Pos:    point.Point{X: 5, Y: 0},
						HP:     60,
					},
					View: point.Point{X: viewCenter.X + 20, Y: viewCenter.Y - 20},
				},
				{
					BattleParam: enemy.EnemyParam{
						CharID: enemy.IDGaroo,
						Pos:    point.Point{X: 5, Y: 2},
						HP:     60,
					},
					View: point.Point{X: viewCenter.X - 25, Y: viewCenter.Y + 10},
				},
				{
					BattleParam: enemy.EnemyParam{
						CharID: enemy.IDMetall,
						Pos:    point.Point{X: 3, Y: 0},
						HP:     40,
					},
					View: point.Point{X: viewCenter.X + 25, Y: viewCenter.Y + 10},
				},
			},
		},
		{
			Name: "灼熱の息吹",
			Enemies: []selectEnemyData{
				{
					BattleParam: enemy.EnemyParam{
						CharID: enemy.IDVolgear,
						Pos:    point.Point{X: 5, Y: 0},
						HP:     80,
					},
					View: point.Point{X: viewCenter.X + 30, Y: viewCenter.Y - 5},
				},
				{
					BattleParam: enemy.EnemyParam{
						CharID: enemy.IDGaroo,
						Pos:    point.Point{X: 3, Y: 2},
						HP:     60,
					},
					View: point.Point{X: viewCenter.X - 20, Y: viewCenter.Y + 25},
				},
			},
		},
		{
			Name: "回る殺戮者",
			Enemies: []selectEnemyData{
				{
					BattleParam: enemy.EnemyParam{
						CharID: enemy.IDCirKill,
						Pos:    point.Point{X: 5, Y: 0},
						HP:     150,
					},
					View: point.Point{X: viewCenter.X + 30, Y: viewCenter.Y - 5},
				},
				{
					BattleParam: enemy.EnemyParam{
						CharID: enemy.IDCirKill,
						Pos:    point.Point{X: 3, Y: 2},
						HP:     150,
					},
					View: point.Point{X: viewCenter.X - 20, Y: viewCenter.Y + 25},
				},
			},
		},
		{
			Name: "水玉バブル",
			Enemies: []selectEnemyData{
				{
					BattleParam: enemy.EnemyParam{
						CharID: enemy.IDShrimpy,
						Pos:    point.Point{X: 5, Y: 1},
						HP:     80,
					},
					View: point.Point{X: viewCenter.X, Y: viewCenter.Y},
				},
			},
		},
		{
			Name: "水を操りし者",
			Enemies: []selectEnemyData{
				{
					BattleParam: enemy.EnemyParam{
						CharID: enemy.IDAquaman,
						Pos:    point.Point{X: 4, Y: 1},
						HP:     500,
					},
					View: point.Point{X: viewCenter.X, Y: viewCenter.Y},
				},
			},
		},
		{
			Name: "極寒より訪れし者",
			Enemies: []selectEnemyData{
				{
					BattleParam: enemy.EnemyParam{
						CharID: enemy.IDColdman,
						Pos:    point.Point{X: 4, Y: 1},
						HP:     700,
					},
					View: point.Point{X: viewCenter.X, Y: viewCenter.Y},
				},
			},
		},
		{
			Name: "永遠の好敵手",
			Enemies: []selectEnemyData{
				{
					BattleParam: enemy.EnemyParam{
						CharID: enemy.IDBlues,
						Pos:    point.Point{X: 4, Y: 1},
						HP:     800,
					},
					View: point.Point{X: viewCenter.X, Y: viewCenter.Y},
				},
			},
		},
		{
			Name: "闇よりの使徒",
			Enemies: []selectEnemyData{
				{
					BattleParam: enemy.EnemyParam{
						CharID: enemy.IDForte,
						Pos:    point.Point{X: 4, Y: 1},
						HP:     1500,
					},
					View: point.Point{X: viewCenter.X, Y: viewCenter.Y},
				},
			},
		},
		{
			Name: "練習",
			Enemies: []selectEnemyData{
				{
					BattleParam: enemy.EnemyParam{
						CharID: enemy.IDTarget,
						Pos:    point.Point{X: 4, Y: 1},
						HP:     1000,
					},
					View: point.Point{X: viewCenter.X, Y: viewCenter.Y + 10},
				},
			},
		},
	}

	names := []string{}
	for _, s := range res.selectData {
		for _, e := range s.Enemies {
			name, ext := enemy.GetStandImageFile(e.BattleParam.CharID)
			fname := name + ext
			res.images[e.BattleParam.CharID] = dxlib.LoadGraph(fname)
		}
		names = append(names, s.Name)
	}
	res.itemList.SetList(names, goBattleListShowMax)

	for id, img := range res.images {
		if img == -1 {
			return nil, errors.Newf("failed to load enemy %d image", id)
		}
	}

	return res, nil
}

func (p *menuGoBattle) End() {
	for _, img := range p.images {
		dxlib.DeleteGraph(img)
	}
}

func (p *menuGoBattle) Update() bool {
	if p.waitCount > 0 {
		// 敵を決定後少し待ってからstateを変更する
		p.waitCount++
		return p.waitCount > 30
	}

	if inputs.CheckKey(inputs.KeyCancel) == 1 {
		sound.On(resources.SECancel)
		return true
	}
	if p.itemList.Update() != -1 {
		sound.On(resources.SEGoBattle)
		p.waitCount++
		p.result = ResultGoBattle
		// グローバルデータを編集する
		battleEnemies = []enemy.EnemyParam{}
		c := p.itemList.GetPointer() + p.itemList.GetScroll()
		for _, e := range p.selectData[c].Enemies {
			battleEnemies = append(battleEnemies, e.BattleParam)
		}

		return false
	}

	return false
}

func (p *menuGoBattle) Draw() {
	dxlib.DrawBox(20, 30, config.ScreenSize.X-20, 300, dxlib.GetColor(168, 192, 216), true)
	dxlib.DrawBox(30, 40, 210, goBattleListShowMax*35+50, dxlib.GetColor(16, 80, 104), true)

	for i := 0; i < goBattleListShowMax; i++ {
		c := i + p.itemList.GetScroll()
		draw.String(65, 50+i*35, 0xffffff, p.selectData[c].Name)
	}

	const s = 2
	y := 50 + p.itemList.GetPointer()*35
	dxlib.DrawTriangle(40, y+s, 40+18-s*2, y+10, 40, y+20-s, 0xffffff, true)

	// Show images
	c := p.itemList.GetPointer() + p.itemList.GetScroll()
	const size = 150
	dxlib.DrawBox(viewCenter.X-size/2, viewCenter.Y-size/2, viewCenter.X+size/2, viewCenter.Y+size/2, 0, true)
	for _, e := range p.selectData[c].Enemies {
		dxlib.DrawRotaGraph(e.View.X, e.View.Y, 1, 0, p.images[e.BattleParam.CharID], true)
	}
}

func (p *menuGoBattle) GetResult() Result {
	return p.result
}
