package field

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/background"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

type extendPanelInfo struct {
	info      battlecommon.PanelInfo
	objExists bool
}

const (
	panelReturnAnimCount = 60
	panelHoleCount       = 480
)

const (
	tmpPanelStatusNormal int = iota
	tmpPanelStatusCrack
	tmpPanelStatusHole

	tmpPanelStatusMax
)

var (
	imgPanel      [tmpPanelStatusMax][battlecommon.PanelTypeMax]int
	blackoutCount = 0
	panels        [][]extendPanelInfo
)

// Init ...
func Init() error {
	logger.Info("Initialize battle field data")

	panels = make([][]extendPanelInfo, battlecommon.FieldNum.X)
	for i := 0; i < battlecommon.FieldNum.X; i++ {
		panels[i] = make([]extendPanelInfo, battlecommon.FieldNum.Y)
	}

	// Initialize images
	files := [tmpPanelStatusMax]string{"normal", "crack", "hole"}
	for i := 0; i < tmpPanelStatusMax; i++ {
		fname := fmt.Sprintf("%sbattle/panel_player_%s.png", common.ImagePath, files[i])
		imgPanel[i][battlecommon.PanelTypePlayer] = dxlib.LoadGraph(fname)
		if imgPanel[i][battlecommon.PanelTypePlayer] < 0 {
			return fmt.Errorf("failed to read player panel image %s", fname)
		}
	}
	for i := 0; i < tmpPanelStatusMax; i++ {
		fname := fmt.Sprintf("%sbattle/panel_enemy_%s.png", common.ImagePath, files[i])
		imgPanel[i][battlecommon.PanelTypeEnemy] = dxlib.LoadGraph(fname)
		if imgPanel[i][battlecommon.PanelTypeEnemy] < 0 {
			return fmt.Errorf("failed to read enemy panel image %s", fname)
		}
	}

	// Initialize panel info
	for x := 0; x < battlecommon.FieldNum.X; x++ {
		t := battlecommon.PanelTypePlayer
		if x >= battlecommon.FieldNum.X/2 {
			t = battlecommon.PanelTypeEnemy
		}
		for y := 0; y < battlecommon.FieldNum.Y; y++ {
			panels[x][y] = extendPanelInfo{
				info: battlecommon.PanelInfo{
					Status:    tmpPanelStatusNormal,
					Type:      t,
					HoleCount: 0,
				},
				objExists: false,
			}
		}
	}

	// TODO: Map情報から取得する
	if err := background.Set(background.Type秋原町); err != nil {
		return fmt.Errorf("failed to load background: %w", err)
	}

	logger.Info("Successfully initialized battle field data")
	return nil
}

// End ...
func End() {
	logger.Info("Cleanup battle field data")
	for i := 0; i < tmpPanelStatusMax; i++ {
		for j := 0; j < battlecommon.PanelTypeMax; j++ {
			dxlib.DeleteGraph(imgPanel[i][j])
			imgPanel[i][j] = -1
		}
	}

	background.Unset()
	logger.Info("Successfully cleanuped battle field data")
}

// Draw ...
func Draw() {
	for x := 0; x < battlecommon.FieldNum.X; x++ {
		for y := 0; y < battlecommon.FieldNum.Y; y++ {
			img := imgPanel[panels[x][y].info.Status][panels[x][y].info.Type]
			vx := battlecommon.PanelSize.X * x
			vy := battlecommon.DrawPanelTopY + battlecommon.PanelSize.Y*y

			// Note:
			//   panelReturnAnimCount以下の場合StatusはNormalになる
			//   HoleとNormalを点滅させるためCountによってイメージを変える
			if panels[x][y].info.HoleCount > 0 {
				if panels[x][y].info.HoleCount < panelReturnAnimCount && (panels[x][y].info.HoleCount/2)%2 == 0 {
					img = imgPanel[tmpPanelStatusHole][panels[x][y].info.Type]
				}
			}

			dxlib.DrawGraph(vx, vy, img, true)

			if dm := damage.Get(common.Point{X: x, Y: y}); dm != nil && dm.ShowHitArea {
				x1 := vx
				y1 := vy
				x2 := vx + battlecommon.PanelSize.X
				y2 := vy + battlecommon.PanelSize.Y
				const s = 5
				dxlib.DrawBox(x1+s, y1+s, x2-s, y2-s, 0xffff00, true)
			}
		}
	}
}

func DrawBlackout() {
	if blackoutCount > 0 {
		dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_ALPHA, 128)
		dxlib.DrawBox(0, 0, common.ScreenSize.X, common.ScreenSize.Y, 0x000000, true)
		dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_NOBLEND, 255)
	}
}

func Update() {
	// Cleanup at first
	for x := 0; x < len(panels); x++ {
		for y := 0; y < len(panels[x]); y++ {
			panels[x][y].info.ObjectID = ""
		}
	}

	objs := localanim.ObjAnimGetObjs(objanim.FilterAll)
	for _, obj := range objs {
		panels[obj.Pos.X][obj.Pos.Y].info.ObjectID = obj.ObjID
		if panels[obj.Pos.X][obj.Pos.Y].info.Status == tmpPanelStatusCrack {
			panels[obj.Pos.X][obj.Pos.Y].objExists = true
		}
	}

	if blackoutCount > 0 {
		blackoutCount--
	}

	// Panel status update
	for x := 0; x < len(panels); x++ {
		for y := 0; y < len(panels[x]); y++ {
			if panels[x][y].info.HoleCount > 0 {
				panels[x][y].info.HoleCount--
			}

			switch panels[x][y].info.Status {
			case tmpPanelStatusHole:
				if panels[x][y].info.HoleCount <= panelReturnAnimCount {
					panels[x][y].info.Status = tmpPanelStatusNormal
				}
			case tmpPanelStatusCrack:
				// Objectが乗って離れたらHole状態へ
				if panels[x][y].objExists && panels[x][y].info.ObjectID == "" {
					sound.On(sound.SEPanelBreak)
					panels[x][y].objExists = false
					panels[x][y].info.Status = tmpPanelStatusHole
					panels[x][y].info.HoleCount = panelHoleCount
				}
			}
		}
	}
}

func GetPanelInfo(pos common.Point) battlecommon.PanelInfo {
	if pos.X < 0 || pos.X >= battlecommon.FieldNum.X || pos.Y < 0 || pos.Y >= battlecommon.FieldNum.Y {
		return battlecommon.PanelInfo{}
	}

	// Update objectID to latest
	panels[pos.X][pos.Y].info.ObjectID = localanim.ObjAnimExistsObject(pos)

	return panels[pos.X][pos.Y].info
}

func SetBlackoutCount(cnt int) {
	blackoutCount = cnt
}

func IsBlackout() bool {
	return blackoutCount > 0
}

func ChangePanelType(pos common.Point, pnType int) {
	if pos.X < 0 || pos.X >= battlecommon.FieldNum.X || pos.Y < 0 || pos.Y >= battlecommon.FieldNum.Y {
		return
	}

	panels[pos.X][pos.Y].info.Type = pnType
}

func PanelBreak(pos common.Point) {
	if pos.X < 0 || pos.X >= battlecommon.FieldNum.X || pos.Y < 0 || pos.Y >= battlecommon.FieldNum.Y {
		return
	}

	if panels[pos.X][pos.Y].info.Status == tmpPanelStatusHole {
		return
	}

	if panels[pos.X][pos.Y].info.ObjectID != "" {
		panels[pos.X][pos.Y].info.Status = tmpPanelStatusCrack
	} else {
		panels[pos.X][pos.Y].info.Status = tmpPanelStatusHole
		panels[pos.X][pos.Y].info.HoleCount = panelHoleCount
	}
}

func PanelCrack(pos common.Point) {
	if pos.X < 0 || pos.X >= battlecommon.FieldNum.X || pos.Y < 0 || pos.Y >= battlecommon.FieldNum.Y {
		return
	}

	if panels[pos.X][pos.Y].info.Status == tmpPanelStatusHole {
		return
	}

	panels[pos.X][pos.Y].info.Status = tmpPanelStatusCrack
}

func Set4x4Area() {
	battlecommon.FieldNum = common.Point{X: 8, Y: 4}
	common.ScreenSize = common.Point{X: 640, Y: 480}
	battlecommon.DrawPanelTopY = common.ScreenSize.Y - (battlecommon.PanelSize.Y * battlecommon.FieldNum.Y) - 30
	dxlib.SetWindowSize(640, 480)
}

func ResetSet4x4Area() {
	if Is4x4Area() {
		battlecommon.FieldNum = common.Point{X: 6, Y: 3}
		common.ScreenSize = common.Point{X: 480, Y: 320}
		battlecommon.DrawPanelTopY = common.ScreenSize.Y - (battlecommon.PanelSize.Y * battlecommon.FieldNum.Y) - 30
		dxlib.SetWindowSize(480, 320)
	}
}

func Is4x4Area() bool {
	return battlecommon.FieldNum.X == 8 && battlecommon.FieldNum.Y == 4
}
