package field

import (
	"fmt"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

const (
	// FieldNumX ...
	FieldNumX = 6
	// FieldNumY ...
	FieldNumY = 3
	// PanelSizeX ...
	PanelSizeX = 80
	// PanelSizeY ...
	PanelSizeY = 50
	// DrawPanelTopY ...
	DrawPanelTopY = common.ScreenY - (PanelSizeY * 3) - 30
)

const (
	PanelTypePlayer int = iota
	PanelTypeEnemy

	panelTypeMax
)

type PanelInfo struct {
	Type     int
	ObjectID string
	// TODO status(毒とか穴とか)
}

var (
	imgPanel = [2]int32{-1, -1}

	panels [FieldNumX][FieldNumY]PanelInfo
)

// Init ...
func Init() error {
	logger.Info("Initialize battle field data")

	// Initialize images
	fname := common.ImagePath + "battle/panel_player.png"
	imgPanel[PanelTypePlayer] = dxlib.LoadGraph(fname)
	if imgPanel[PanelTypePlayer] < 0 {
		return fmt.Errorf("failed to read player panel image %s", fname)
	}
	fname = common.ImagePath + "battle/panel_enemy.png"
	imgPanel[PanelTypeEnemy] = dxlib.LoadGraph(fname)
	if imgPanel[PanelTypeEnemy] < 0 {
		return fmt.Errorf("failed to read enemy panel image %s", fname)
	}

	// Initialize panel info
	for x := 0; x < FieldNumX; x++ {
		t := PanelTypePlayer
		if x > 2 {
			t = PanelTypeEnemy
		}
		for y := 0; y < FieldNumY; y++ {
			panels[x][y] = PanelInfo{
				Type: t,
			}
		}
	}
	// TODO: special field

	logger.Info("Successfully initialized battle field data")
	return nil
}

// End ...
func End() {
	logger.Info("Cleanup battle field data")
	for i := 0; i < panelTypeMax; i++ {
		dxlib.DeleteGraph(imgPanel[i])
		imgPanel[i] = -1
	}
	logger.Info("Successfully cleanuped battle field data")
}

// Draw ...
func Draw() {
	for x := 0; x < FieldNumX; x++ {
		for y := 0; y < FieldNumY; y++ {
			img := imgPanel[panels[x][y].Type]
			vx := int32(PanelSizeX * x)
			vy := int32(DrawPanelTopY + PanelSizeY*y)
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

func Update() {
	// Cleanup at first
	for x := 0; x < len(panels); x++ {
		for y := 0; y < len(panels[x]); y++ {
			panels[x][y].ObjectID = ""
		}
	}

	objs := anim.GetObjs(anim.Filter{ObjType: anim.ObjTypePlayer | anim.ObjTypeEnemy})
	for _, obj := range objs {
		panels[obj.PosX][obj.PosY].ObjectID = obj.ObjID
	}
}

func GetPanelInfo(x, y int) PanelInfo {
	return panels[x][y]
}
