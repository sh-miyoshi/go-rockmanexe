package menu

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/enemy"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/inputs"
)

type selectEnemyData struct {
	BattleParam enemy.EnemyParam
	View        common.Point
}

type selectValue struct {
	Name    string
	Enemies []selectEnemyData
}

const (
	goBattleListShowMax = 7
)

var (
	viewCenter = common.Point{X: 350, Y: 150}

	goBattleSelectData []selectValue
	goBattleCursor     int
	goBattleScroll     int
	goBattleWaitCount  int
	images             = make(map[int]int)
)

func goBattleInit() error {
	goBattleCursor = 0
	goBattleScroll = 0
	goBattleWaitCount = 0

	goBattleSelectData = []selectValue{
		{
			Name: "千里の道も一歩から",
			Enemies: []selectEnemyData{
				{
					BattleParam: enemy.EnemyParam{
						CharID: enemy.IDMetall,
						Pos:    common.Point{X: 4, Y: 1},
						HP:     40,
					},
					View: common.Point{X: viewCenter.X, Y: viewCenter.Y + 10},
				},
			},
		},
		{
			Name: "侮ることなかれ",
			Enemies: []selectEnemyData{
				{
					BattleParam: enemy.EnemyParam{
						CharID: enemy.IDMetall,
						Pos:    common.Point{X: 3, Y: 0},
						HP:     40,
					},
					View: common.Point{X: viewCenter.X, Y: viewCenter.Y - 30},
				},
				{
					BattleParam: enemy.EnemyParam{
						CharID: enemy.IDMetall,
						Pos:    common.Point{X: 4, Y: 1},
						HP:     40,
					},
					View: common.Point{X: viewCenter.X - 30, Y: viewCenter.Y + 10},
				},
				{
					BattleParam: enemy.EnemyParam{
						CharID: enemy.IDMetall,
						Pos:    common.Point{X: 5, Y: 2},
						HP:     40,
					},
					View: common.Point{X: viewCenter.X + 30, Y: viewCenter.Y + 10},
				},
			},
		},
		{
			Name: "跡追いする電気玉",
			Enemies: []selectEnemyData{
				{
					BattleParam: enemy.EnemyParam{
						CharID: enemy.IDBilly,
						Pos:    common.Point{X: 5, Y: 1},
						HP:     50,
					},
					View: common.Point{X: viewCenter.X + 20, Y: viewCenter.Y - 10},
				},
				{
					BattleParam: enemy.EnemyParam{
						CharID: enemy.IDMetall,
						Pos:    common.Point{X: 3, Y: 0},
						HP:     40,
					},
					View: common.Point{X: viewCenter.X - 30, Y: viewCenter.Y + 10},
				},
			},
		},
		{
			Name: "遊泳するものたち",
			Enemies: []selectEnemyData{
				{
					BattleParam: enemy.EnemyParam{
						CharID: enemy.IDLark,
						Pos:    common.Point{X: 3, Y: 0},
						HP:     60,
					},
					View: common.Point{X: viewCenter.X + 10, Y: viewCenter.Y - 25},
				},
				{
					BattleParam: enemy.EnemyParam{
						CharID: enemy.IDLark,
						Pos:    common.Point{X: 4, Y: 1},
						HP:     60,
					},
					View: common.Point{X: viewCenter.X - 10, Y: viewCenter.Y + 25},
				},
			},
		},
		{
			Name: "舞戻るやいば",
			Enemies: []selectEnemyData{
				{
					BattleParam: enemy.EnemyParam{
						CharID: enemy.IDBoomer,
						Pos:    common.Point{X: 5, Y: 1},
						HP:     60,
					},
					View: common.Point{X: viewCenter.X + 20, Y: viewCenter.Y - 5},
				},
				{
					BattleParam: enemy.EnemyParam{
						CharID: enemy.IDMetall,
						Pos:    common.Point{X: 4, Y: 0},
						HP:     40,
					},
					View: common.Point{X: viewCenter.X - 30, Y: viewCenter.Y + 5},
				},
			},
		},
		{
			Name: "火玉に注意！",
			Enemies: []selectEnemyData{
				{
					BattleParam: enemy.EnemyParam{
						CharID: enemy.IDGaroo,
						Pos:    common.Point{X: 5, Y: 0},
						HP:     60,
					},
					View: common.Point{X: viewCenter.X + 20, Y: viewCenter.Y - 20},
				},
				{
					BattleParam: enemy.EnemyParam{
						CharID: enemy.IDGaroo,
						Pos:    common.Point{X: 5, Y: 2},
						HP:     60,
					},
					View: common.Point{X: viewCenter.X - 25, Y: viewCenter.Y + 10},
				},
				{
					BattleParam: enemy.EnemyParam{
						CharID: enemy.IDMetall,
						Pos:    common.Point{X: 3, Y: 0},
						HP:     40,
					},
					View: common.Point{X: viewCenter.X + 25, Y: viewCenter.Y + 10},
				},
			},
		},
		{
			Name: "灼熱の息吹",
			Enemies: []selectEnemyData{
				{
					BattleParam: enemy.EnemyParam{
						CharID: enemy.IDVolgear,
						Pos:    common.Point{X: 5, Y: 0},
						HP:     80,
					},
					View: common.Point{X: viewCenter.X + 30, Y: viewCenter.Y - 5},
				},
				{
					BattleParam: enemy.EnemyParam{
						CharID: enemy.IDGaroo,
						Pos:    common.Point{X: 3, Y: 2},
						HP:     60,
					},
					View: common.Point{X: viewCenter.X - 20, Y: viewCenter.Y + 25},
				},
			},
		},
		{
			Name: "水を操りし者",
			Enemies: []selectEnemyData{
				{
					BattleParam: enemy.EnemyParam{
						CharID: enemy.IDAquaman,
						Pos:    common.Point{X: 4, Y: 1},
						HP:     500,
					},
					View: common.Point{X: viewCenter.X, Y: viewCenter.Y},
				},
			},
		},
		{
			Name: "練習",
			Enemies: []selectEnemyData{
				{
					BattleParam: enemy.EnemyParam{
						CharID: enemy.IDTarget,
						Pos:    common.Point{X: 4, Y: 1},
						HP:     1000,
					},
					View: common.Point{X: viewCenter.X, Y: viewCenter.Y + 10},
				},
			},
		},
	}

	for _, s := range goBattleSelectData {
		for _, e := range s.Enemies {
			name, ext := enemy.GetStandImageFile(e.BattleParam.CharID)
			fname := name + ext
			images[e.BattleParam.CharID] = dxlib.LoadGraph(fname)
		}
	}
	for id, img := range images {
		if img == -1 {
			return fmt.Errorf("failed to load enemy %d image", id)
		}
	}

	return nil
}

func goBattleEnd() {
	for _, img := range images {
		dxlib.DeleteGraph(img)
	}
	images = make(map[int]int)
}

func goBattleProcess() bool {
	if goBattleWaitCount > 0 {
		goBattleWaitCount++
		return goBattleWaitCount > 30
	}

	if inputs.CheckKey(inputs.KeyCancel) == 1 {
		sound.On(sound.SECancel)
		stateChange(stateTop)
		return false
	}
	if inputs.CheckKey(inputs.KeyEnter) == 1 {
		sound.On(sound.SEGoBattle)
		goBattleWaitCount++
		return false
	}

	if inputs.CheckKey(inputs.KeyUp)%10 == 1 {
		if goBattleCursor > 0 {
			sound.On(sound.SECursorMove)
			goBattleCursor--
		} else if goBattleScroll > 0 {
			sound.On(sound.SECursorMove)
			goBattleScroll--
		}
	} else if inputs.CheckKey(inputs.KeyDown)%10 == 1 {
		n := goBattleListShowMax - 1
		if len(goBattleSelectData) < goBattleListShowMax {
			n = len(goBattleSelectData) - 1
		}

		if goBattleCursor < n {
			sound.On(sound.SECursorMove)
			goBattleCursor++
		} else if goBattleScroll < len(goBattleSelectData)-goBattleListShowMax {
			sound.On(sound.SECursorMove)
			goBattleScroll++
		}
	}

	return false
}

func goBattleDraw() {
	dxlib.DrawBox(20, 30, common.ScreenSize.X-20, 300, dxlib.GetColor(168, 192, 216), true)
	dxlib.DrawBox(30, 40, 210, goBattleListShowMax*35+50, dxlib.GetColor(16, 80, 104), true)

	for i := 0; i < goBattleListShowMax; i++ {
		c := i + goBattleScroll
		draw.String(65, 50+i*35, 0xffffff, goBattleSelectData[c].Name)
	}

	const s = 2
	y := 50 + goBattleCursor*35
	dxlib.DrawTriangle(40, y+s, 40+18-s*2, y+10, 40, y+20-s, 0xffffff, true)

	// Show images
	c := goBattleCursor + goBattleScroll
	const size = 150
	dxlib.DrawBox(viewCenter.X-size/2, viewCenter.Y-size/2, viewCenter.X+size/2, viewCenter.Y+size/2, 0, true)
	for _, e := range goBattleSelectData[c].Enemies {
		dxlib.DrawRotaGraph(e.View.X, e.View.Y, 1, 0, images[e.BattleParam.CharID], true)
	}
}

func battleEnemies() []enemy.EnemyParam {
	if config.Get().Debug.SkipMenu {
		// Start from battle mode for debug
		// return debug data
		return []enemy.EnemyParam{
			{
				CharID: enemy.IDTarget,
				Pos:    common.Point{X: 4, Y: 1},
				HP:     1000,
			},
		}
	}

	res := []enemy.EnemyParam{}
	c := goBattleCursor + goBattleScroll
	for _, e := range goBattleSelectData[c].Enemies {
		res = append(res, e.BattleParam)
	}

	return res
}
