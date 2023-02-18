package field

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/background"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

const (
	panelReturnAnimCount = 60
	panelHoleCount       = 480
)

var (
	tmpFieldNum      = common.Point{X: 6, Y: 3} // TODO: 要修正
	tmpPanelSize     = common.Point{X: 80, Y: 50}
	tmpDrawPanelTopY = common.ScreenSize.Y - (tmpPanelSize.Y * tmpFieldNum.Y) - 30
)

const (
	PanelTypePlayer int = iota
	PanelTypeEnemy

	panelTypeMax
)

const (
	tmpPanelStatusNormal int = iota
	tmpPanelStatusCrack
	tmpPanelStatusHole

	tmpPanelStatusMax
)

type PanelInfo struct {
	objExists bool

	Type      int
	ObjectID  string
	Status    int
	HoleCount int
}

var (
	imgPanel      [tmpPanelStatusMax][panelTypeMax]int
	blackoutCount = 0
	panels        [][]PanelInfo
)

// Init ...
func Init() error {
	logger.Info("Initialize battle field data")

	panels = make([][]PanelInfo, tmpFieldNum.X)
	for i := 0; i < tmpFieldNum.X; i++ {
		panels[i] = make([]PanelInfo, tmpFieldNum.Y)
	}

	// Initialize images
	files := [tmpPanelStatusMax]string{"normal", "crack", "hole"}
	for i := 0; i < tmpPanelStatusMax; i++ {
		fname := fmt.Sprintf("%sbattle/panel_player_%s.png", common.ImagePath, files[i])
		imgPanel[i][PanelTypePlayer] = dxlib.LoadGraph(fname)
		if imgPanel[i][PanelTypePlayer] < 0 {
			return fmt.Errorf("failed to read player panel image %s", fname)
		}
	}
	for i := 0; i < tmpPanelStatusMax; i++ {
		fname := fmt.Sprintf("%sbattle/panel_enemy_%s.png", common.ImagePath, files[i])
		imgPanel[i][PanelTypeEnemy] = dxlib.LoadGraph(fname)
		if imgPanel[i][PanelTypeEnemy] < 0 {
			return fmt.Errorf("failed to read enemy panel image %s", fname)
		}
	}

	// Initialize panel info
	for x := 0; x < tmpFieldNum.X; x++ {
		t := PanelTypePlayer
		if x >= tmpFieldNum.X/2 {
			t = PanelTypeEnemy
		}
		for y := 0; y < tmpFieldNum.Y; y++ {
			panels[x][y] = PanelInfo{
				Status:    tmpPanelStatusNormal,
				Type:      t,
				HoleCount: 0,
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
		for j := 0; j < panelTypeMax; j++ {
			dxlib.DeleteGraph(imgPanel[i][j])
			imgPanel[i][j] = -1
		}
	}

	background.Unset()
	logger.Info("Successfully cleanuped battle field data")
}

// Draw ...
func Draw() {
	for x := 0; x < tmpFieldNum.X; x++ {
		for y := 0; y < tmpFieldNum.Y; y++ {
			img := imgPanel[panels[x][y].Status][panels[x][y].Type]
			vx := tmpPanelSize.X * x
			vy := tmpDrawPanelTopY + tmpPanelSize.Y*y

			// Note:
			//   panelReturnAnimCount以下の場合StatusはNormalになる
			//   HoleとNormalを点滅させるためCountによってイメージを変える
			if panels[x][y].HoleCount > 0 {
				if panels[x][y].HoleCount < panelReturnAnimCount && (panels[x][y].HoleCount/2)%2 == 0 {
					img = imgPanel[tmpPanelStatusHole][panels[x][y].Type]
				}
			}

			dxlib.DrawGraph(vx, vy, img, true)

			if dm := damage.Get(common.Point{X: x, Y: y}); dm != nil && dm.ShowHitArea {
				x1 := vx
				y1 := vy
				x2 := vx + tmpPanelSize.X
				y2 := vy + tmpPanelSize.Y
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
			panels[x][y].ObjectID = ""
		}
	}

	objs := objanim.GetObjs(objanim.Filter{ObjType: objanim.ObjTypeAll})
	for _, obj := range objs {
		panels[obj.Pos.X][obj.Pos.Y].ObjectID = obj.ObjID
		if panels[obj.Pos.X][obj.Pos.Y].Status == tmpPanelStatusCrack {
			panels[obj.Pos.X][obj.Pos.Y].objExists = true
		}
	}

	if blackoutCount > 0 {
		blackoutCount--
	}

	// Panel status update
	for x := 0; x < len(panels); x++ {
		for y := 0; y < len(panels[x]); y++ {
			if panels[x][y].HoleCount > 0 {
				panels[x][y].HoleCount--
			}

			switch panels[x][y].Status {
			case tmpPanelStatusHole:
				if panels[x][y].HoleCount <= panelReturnAnimCount {
					panels[x][y].Status = tmpPanelStatusNormal
				}
			case tmpPanelStatusCrack:
				// Objectが乗って離れたらHole状態へ
				if panels[x][y].objExists && panels[x][y].ObjectID == "" {
					sound.On(sound.SEPanelBreak)
					panels[x][y].objExists = false
					panels[x][y].Status = tmpPanelStatusHole
					panels[x][y].HoleCount = panelHoleCount
				}
			}
		}
	}
}

func GetPanelInfo(pos common.Point) PanelInfo {
	if pos.X < 0 || pos.X >= tmpFieldNum.X || pos.Y < 0 || pos.Y >= tmpFieldNum.Y {
		return PanelInfo{}
	}

	// Update objectID to latest
	panels[pos.X][pos.Y].ObjectID = objanim.ExistsObject(pos)

	return panels[pos.X][pos.Y]
}

func SetBlackoutCount(cnt int) {
	blackoutCount = cnt
}

func IsBlackout() bool {
	return blackoutCount > 0
}

func ChangePanelType(pos common.Point, pnType int) {
	if pos.X < 0 || pos.X >= tmpFieldNum.X || pos.Y < 0 || pos.Y >= tmpFieldNum.Y {
		return
	}

	panels[pos.X][pos.Y].Type = pnType
}

func PanelBreak(pos common.Point) {
	if pos.X < 0 || pos.X >= tmpFieldNum.X || pos.Y < 0 || pos.Y >= tmpFieldNum.Y {
		return
	}

	if panels[pos.X][pos.Y].Status == tmpPanelStatusHole {
		return
	}

	if panels[pos.X][pos.Y].ObjectID != "" {
		panels[pos.X][pos.Y].Status = tmpPanelStatusCrack
	} else {
		panels[pos.X][pos.Y].Status = tmpPanelStatusHole
		panels[pos.X][pos.Y].HoleCount = panelHoleCount
	}
}

func PanelCrack(pos common.Point) {
	if pos.X < 0 || pos.X >= tmpFieldNum.X || pos.Y < 0 || pos.Y >= tmpFieldNum.Y {
		return
	}

	if panels[pos.X][pos.Y].Status == tmpPanelStatusHole {
		return
	}

	panels[pos.X][pos.Y].Status = tmpPanelStatusCrack
}

func Set4x4Area() {
	tmpFieldNum = common.Point{X: 8, Y: 4}
	common.ScreenSize = common.Point{X: 640, Y: 480}
	tmpDrawPanelTopY = common.ScreenSize.Y - (tmpPanelSize.Y * tmpFieldNum.Y) - 30
	dxlib.SetWindowSize(640, 480)
}

func ResetSet4x4Area() {
	if Is4x4Area() {
		tmpFieldNum = common.Point{X: 6, Y: 3}
		common.ScreenSize = common.Point{X: 480, Y: 320}
		tmpDrawPanelTopY = common.ScreenSize.Y - (tmpPanelSize.Y * tmpFieldNum.Y) - 30
		dxlib.SetWindowSize(480, 320)
	}
}

func Is4x4Area() bool {
	return tmpFieldNum.X == 8 && tmpFieldNum.Y == 4
}
