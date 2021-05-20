package menu

import (
	"fmt"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/enemy"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/inputs"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/sound"
)

type selectEnemyData struct {
	BattleParam enemy.EnemyParam
	ViewPosX    int32
	ViewPosY    int32
}

type selectValue struct {
	Name    string
	Enemies []selectEnemyData
}

const (
	viewCenterX = 350
	viewCenterY = 150
)

var (
	goBattleSelectData []selectValue
	goBattleCursor     int
	goBattleWaitCount  int
	images             = make(map[int]int32)
)

func goBattleInit() error {
	goBattleCursor = 0
	goBattleWaitCount = 0

	goBattleSelectData = []selectValue{
		{
			Name: "千里の道も一歩から",
			Enemies: []selectEnemyData{
				{
					BattleParam: enemy.EnemyParam{
						CharID: enemy.IDMetall,
						PosX:   4,
						PosY:   1,
						HP:     40,
					},
					ViewPosX: viewCenterX,
					ViewPosY: viewCenterY + 10,
				},
			},
		},
		{
			Name: "侮ることなかれ",
			Enemies: []selectEnemyData{
				{
					BattleParam: enemy.EnemyParam{
						CharID: enemy.IDMetall,
						PosX:   3,
						PosY:   0,
						HP:     40,
					},
					ViewPosX: viewCenterX,
					ViewPosY: viewCenterY - 30,
				},
				{
					BattleParam: enemy.EnemyParam{
						CharID: enemy.IDMetall,
						PosX:   4,
						PosY:   1,
						HP:     40,
					},
					ViewPosX: viewCenterX - 30,
					ViewPosY: viewCenterY + 10,
				},
				{
					BattleParam: enemy.EnemyParam{
						CharID: enemy.IDMetall,
						PosX:   5,
						PosY:   2,
						HP:     40,
					},
					ViewPosX: viewCenterX + 30,
					ViewPosY: viewCenterY + 10,
				},
			},
		},
		{
			Name: "跡追いする電気玉",
			Enemies: []selectEnemyData{
				{
					BattleParam: enemy.EnemyParam{
						CharID: enemy.IDBilly,
						PosX:   5,
						PosY:   1,
						HP:     50,
					},
					ViewPosX: viewCenterX + 20,
					ViewPosY: viewCenterY - 10,
				},
				{
					BattleParam: enemy.EnemyParam{
						CharID: enemy.IDMetall,
						PosX:   3,
						PosY:   0,
						HP:     40,
					},
					ViewPosX: viewCenterX - 30,
					ViewPosY: viewCenterY + 10,
				},
			},
		},
		{
			Name: "遊泳するものたち",
			Enemies: []selectEnemyData{
				{
					BattleParam: enemy.EnemyParam{
						CharID: enemy.IDLark,
						PosX:   3,
						PosY:   0,
						HP:     60,
					},
					ViewPosX: viewCenterX + 10,
					ViewPosY: viewCenterY - 25,
				},
				{
					BattleParam: enemy.EnemyParam{
						CharID: enemy.IDLark,
						PosX:   4,
						PosY:   1,
						HP:     60,
					},
					ViewPosX: viewCenterX - 10,
					ViewPosY: viewCenterY + 25,
				},
			},
		},
		{
			Name: "練習",
			Enemies: []selectEnemyData{
				{
					BattleParam: enemy.EnemyParam{
						CharID: enemy.IDTarget,
						PosX:   4,
						PosY:   1,
						HP:     1000,
					},
					ViewPosX: viewCenterX,
					ViewPosY: viewCenterY + 10,
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
	images = make(map[int]int32)
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
	if inputs.CheckKey(inputs.KeyUp) == 1 && goBattleCursor > 0 {
		sound.On(sound.SECursorMove)
		goBattleCursor--
	} else if inputs.CheckKey(inputs.KeyDown) == 1 && goBattleCursor < len(goBattleSelectData)-1 {
		sound.On(sound.SECursorMove)
		goBattleCursor++
	}

	return false
}

func goBattleDraw() {
	dxlib.DrawBox(20, 30, common.ScreenX-20, 300, dxlib.GetColor(168, 192, 216), dxlib.TRUE)
	dxlib.DrawBox(30, 40, 210, int32(len(goBattleSelectData)*35)+50, dxlib.GetColor(16, 80, 104), dxlib.TRUE)

	for i, s := range goBattleSelectData {
		draw.String(65, 50+int32(i)*35, 0xffffff, s.Name)
	}

	const s = 2
	y := int32(50 + goBattleCursor*35)
	dxlib.DrawTriangle(40, y+s, 40+18-s*2, y+10, 40, y+20-s, 0xffffff, dxlib.TRUE)

	// Show images
	const size = 150
	dxlib.DrawBox(viewCenterX-size/2, viewCenterY-size/2, viewCenterX+size/2, viewCenterY+size/2, 0, dxlib.TRUE)
	for _, e := range goBattleSelectData[goBattleCursor].Enemies {
		dxlib.DrawRotaGraph(e.ViewPosX, e.ViewPosY, 1, 0, images[e.BattleParam.CharID], dxlib.TRUE)
	}
}

func battleEnemies() []enemy.EnemyParam {
	if config.Get().Debug.SkipMenu {
		// Start from battle mode for debug
		// return debug data
		return []enemy.EnemyParam{
			{
				CharID: enemy.IDTarget,
				PosX:   4,
				PosY:   1,
				HP:     1000,
			},
		}
	}

	res := []enemy.EnemyParam{}
	for _, e := range goBattleSelectData[goBattleCursor].Enemies {
		res = append(res, e.BattleParam)
	}

	return res
}
