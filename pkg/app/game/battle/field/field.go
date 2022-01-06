package field

import (
	"fmt"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

const (
	FieldNumX     = 6
	FieldNumY     = 3
	PanelSizeX    = 80
	PanelSizeY    = 50
	DrawPanelTopY = common.ScreenY - (PanelSizeY * 3) - 30

	panelReturnAnimCount = 60
	panelHoleCount       = 480
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

	panelStatusMax
)

type PanelInfo struct {
	Type      int
	ObjectID  string
	Status    int
	HoleCount int
	ObjExists bool
}

var (
	imgPanel      [panelStatusMax][panelTypeMax]int32
	blackoutCount = 0
	panels        [FieldNumX][FieldNumY]PanelInfo
)

// Init ...
func Init() error {
	logger.Info("Initialize battle field data")

	// Initialize images
	files := [panelStatusMax]string{"normal", "crack", "hole"}
	for i := 0; i < panelStatusMax; i++ {
		fname := fmt.Sprintf("%sbattle/panel_player_%s.png", common.ImagePath, files[i])
		imgPanel[i][PanelTypePlayer] = dxlib.LoadGraph(fname)
		if imgPanel[i][PanelTypePlayer] < 0 {
			return fmt.Errorf("failed to read player panel image %s", fname)
		}
	}
	for i := 0; i < panelStatusMax; i++ {
		fname := fmt.Sprintf("%sbattle/panel_enemy_%s.png", common.ImagePath, files[i])
		imgPanel[i][PanelTypeEnemy] = dxlib.LoadGraph(fname)
		if imgPanel[i][PanelTypeEnemy] < 0 {
			return fmt.Errorf("failed to read enemy panel image %s", fname)
		}
	}

	// Initialize panel info
	for x := 0; x < FieldNumX; x++ {
		t := PanelTypePlayer
		if x > 2 {
			t = PanelTypeEnemy
		}
		for y := 0; y < FieldNumY; y++ {
			panels[x][y] = PanelInfo{
				Status:    PanelStatusNormal,
				Type:      t,
				HoleCount: 0,
				ObjExists: false,
			}
		}
	}

	logger.Info("Successfully initialized battle field data")
	return nil
}

// End ...
func End() {
	logger.Info("Cleanup battle field data")
	for i := 0; i < panelStatusMax; i++ {
		for j := 0; j < panelTypeMax; j++ {
			dxlib.DeleteGraph(imgPanel[i][j])
			imgPanel[i][j] = -1
		}
	}
	logger.Info("Successfully cleanuped battle field data")
}

// Draw ...
func Draw() {
	for x := 0; x < FieldNumX; x++ {
		for y := 0; y < FieldNumY; y++ {
			img := imgPanel[panels[x][y].Status][panels[x][y].Type]
			vx := int32(PanelSizeX * x)
			vy := int32(DrawPanelTopY + PanelSizeY*y)

			// Note:
			//   panelReturnAnimCount以下の場合StatusはNormalになる
			//   HoleとNormalを点滅させるためCountによってイメージを変える
			if panels[x][y].HoleCount > 0 {
				if panels[x][y].HoleCount < panelReturnAnimCount && (panels[x][y].HoleCount/2)%2 == 0 {
					img = imgPanel[PanelStatusHole][panels[x][y].Type]
				}
			}

			dxlib.DrawGraph(vx, vy, img, dxlib.TRUE)

			if dm := damage.Get(x, y); dm != nil && dm.ShowHitArea {
				x1 := vx
				y1 := vy
				x2 := vx + PanelSizeX
				y2 := vy + PanelSizeY
				const s = 5
				dxlib.DrawBox(x1+s, y1+s, x2-s, y2-s, 0xffff00, dxlib.TRUE)
			}
		}
	}
}

func DrawBlackout() {
	if blackoutCount > 0 {
		dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_ALPHA, 128)
		dxlib.DrawBox(0, 0, common.ScreenX, common.ScreenY, 0x000000, dxlib.TRUE)
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
		panels[obj.PosX][obj.PosY].ObjectID = obj.ObjID
		if panels[obj.PosX][obj.PosY].Status == PanelStatusCrack {
			panels[obj.PosX][obj.PosY].ObjExists = true
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
				if panels[x][y].ObjExists && panels[x][y].ObjectID == "" {
					sound.On(sound.SEPanelBreak)
					panels[x][y].ObjExists = false
					panels[x][y].Status = PanelStatusHole
					panels[x][y].HoleCount = panelHoleCount
				}
			}
		}
	}
}

func GetPanelInfo(x, y int) PanelInfo {
	return panels[x][y]
}

func SetBlackoutCount(cnt int) {
	blackoutCount = cnt
}

func IsBlackout() bool {
	return blackoutCount > 0
}

func PanelBreak(x, y int) {
	if panels[x][y].Status == PanelStatusHole {
		return
	}

	if panels[x][y].ObjectID != "" {
		panels[x][y].Status = PanelStatusCrack
	} else {
		panels[x][y].Status = PanelStatusHole
		panels[x][y].HoleCount = panelHoleCount
	}
}
