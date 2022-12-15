package field

import (
	"fmt"

	originaldxlib "github.com/sh-miyoshi/dxlib"
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
	FieldNum      = common.Point{X: 6, Y: 3}
	PanelSize     = common.Point{X: 80, Y: 50}
	DrawPanelTopY = common.ScreenSize.Y - (PanelSize.Y * FieldNum.Y) - 30
)

const (
	PanelTypePlayer int = iota
	PanelTypeEnemy

	panelTypeMax
)

const (
	PanelStatusNormal int = iota
	PanelStatusCrack
	PanelStatusHole

	PanelStatusMax
)

type PanelInfo struct {
	objExists bool

	Type      int
	ObjectID  string
	Status    int
	HoleCount int
}

var (
	imgPanel      [PanelStatusMax][panelTypeMax]int
	blackoutCount = 0
	panels        [][]PanelInfo
)

// Init ...
func Init() error {
	logger.Info("Initialize battle field data")

	panels = make([][]PanelInfo, FieldNum.X)
	for i := 0; i < FieldNum.X; i++ {
		panels[i] = make([]PanelInfo, FieldNum.Y)
	}

	// Initialize images
	files := [PanelStatusMax]string{"normal", "crack", "hole"}
	for i := 0; i < PanelStatusMax; i++ {
		fname := fmt.Sprintf("%sbattle/panel_player_%s.png", common.ImagePath, files[i])
		imgPanel[i][PanelTypePlayer] = dxlib.LoadGraph(fname)
		if imgPanel[i][PanelTypePlayer] < 0 {
			return fmt.Errorf("failed to read player panel image %s", fname)
		}
	}
	for i := 0; i < PanelStatusMax; i++ {
		fname := fmt.Sprintf("%sbattle/panel_enemy_%s.png", common.ImagePath, files[i])
		imgPanel[i][PanelTypeEnemy] = dxlib.LoadGraph(fname)
		if imgPanel[i][PanelTypeEnemy] < 0 {
			return fmt.Errorf("failed to read enemy panel image %s", fname)
		}
	}

	// Initialize panel info
	for x := 0; x < FieldNum.X; x++ {
		t := PanelTypePlayer
		if x >= FieldNum.X/2 {
			t = PanelTypeEnemy
		}
		for y := 0; y < FieldNum.Y; y++ {
			panels[x][y] = PanelInfo{
				Status:    PanelStatusNormal,
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
	for i := 0; i < PanelStatusMax; i++ {
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
	for x := 0; x < FieldNum.X; x++ {
		for y := 0; y < FieldNum.Y; y++ {
			img := imgPanel[panels[x][y].Status][panels[x][y].Type]
			vx := PanelSize.X * x
			vy := DrawPanelTopY + PanelSize.Y*y

			// Note:
			//   panelReturnAnimCount以下の場合StatusはNormalになる
			//   HoleとNormalを点滅させるためCountによってイメージを変える
			if panels[x][y].HoleCount > 0 {
				if panels[x][y].HoleCount < panelReturnAnimCount && (panels[x][y].HoleCount/2)%2 == 0 {
					img = imgPanel[PanelStatusHole][panels[x][y].Type]
				}
			}

			dxlib.DrawGraph(vx, vy, img, true)

			if dm := damage.Get(common.Point{X: x, Y: y}); dm != nil && dm.ShowHitArea {
				x1 := vx
				y1 := vy
				x2 := vx + PanelSize.X
				y2 := vy + PanelSize.Y
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
		if panels[obj.Pos.X][obj.Pos.Y].Status == PanelStatusCrack {
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
			case PanelStatusHole:
				if panels[x][y].HoleCount <= panelReturnAnimCount {
					panels[x][y].Status = PanelStatusNormal
				}
			case PanelStatusCrack:
				// Objectが乗って離れたらHole状態へ
				if panels[x][y].objExists && panels[x][y].ObjectID == "" {
					sound.On(sound.SEPanelBreak)
					panels[x][y].objExists = false
					panels[x][y].Status = PanelStatusHole
					panels[x][y].HoleCount = panelHoleCount
				}
			}
		}
	}
}

func GetPanelInfo(pos common.Point) PanelInfo {
	if pos.X < 0 || pos.X >= FieldNum.X || pos.Y < 0 || pos.Y >= FieldNum.Y {
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
	if pos.X < 0 || pos.X >= FieldNum.X || pos.Y < 0 || pos.Y >= FieldNum.Y {
		return
	}

	panels[pos.X][pos.Y].Type = pnType
}

func PanelBreak(pos common.Point) {
	if pos.X < 0 || pos.X >= FieldNum.X || pos.Y < 0 || pos.Y >= FieldNum.Y {
		return
	}

	if panels[pos.X][pos.Y].Status == PanelStatusHole {
		return
	}

	if panels[pos.X][pos.Y].ObjectID != "" {
		panels[pos.X][pos.Y].Status = PanelStatusCrack
	} else {
		panels[pos.X][pos.Y].Status = PanelStatusHole
		panels[pos.X][pos.Y].HoleCount = panelHoleCount
	}
}

func PanelCrack(pos common.Point) {
	if pos.X < 0 || pos.X >= FieldNum.X || pos.Y < 0 || pos.Y >= FieldNum.Y {
		return
	}

	if panels[pos.X][pos.Y].Status == PanelStatusHole {
		return
	}

	panels[pos.X][pos.Y].Status = PanelStatusCrack
}

func Set4x4Area() {
	FieldNum = common.Point{X: 8, Y: 4}
	common.ScreenSize = common.Point{X: 640, Y: 480}
	DrawPanelTopY = common.ScreenSize.Y - (PanelSize.Y * FieldNum.Y) - 30
	originaldxlib.SetWindowSize(640, 480)
}

func ResetSet4x4Area() {
	if Is4x4Area() {
		FieldNum = common.Point{X: 6, Y: 3}
		common.ScreenSize = common.Point{X: 480, Y: 320}
		DrawPanelTopY = common.ScreenSize.Y - (PanelSize.Y * FieldNum.Y) - 30
		originaldxlib.SetWindowSize(480, 320)
	}
}

func Is4x4Area() bool {
	return FieldNum.X == 8 && FieldNum.Y == 4
}
