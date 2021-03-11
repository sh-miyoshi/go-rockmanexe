package field

import (
	"fmt"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/common"
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
	typePlayer int = iota
	typeEnemy
	typeMax
)

// ObjectPosition ...
type ObjectPosition struct {
	ID   string
	X, Y int
}

type panelInfo struct {
	typ      int
	objectID string
	// TODO status(毒とか穴とか)
}

var (
	imgPanel = [2]int32{-1, -1}
	panels   [FieldNumX][FieldNumY]panelInfo
)

// Init ...
func Init() error {
	logger.Info("Initialize battle field data")

	// Initialize images
	fname := common.ImagePath + "battle/panel_player.png"
	imgPanel[typePlayer] = dxlib.LoadGraph(fname)
	if imgPanel[typePlayer] < 0 {
		return fmt.Errorf("Failed to read player panel image %s", fname)
	}
	fname = common.ImagePath + "battle/panel_enemy.png"
	imgPanel[typeEnemy] = dxlib.LoadGraph(fname)
	if imgPanel[typeEnemy] < 0 {
		return fmt.Errorf("Failed to read enemy panel image %s", fname)
	}

	// Initialize panel info
	for x := 0; x < FieldNumX; x++ {
		t := typePlayer
		if x > 2 {
			t = typeEnemy
		}
		for y := 0; y < FieldNumY; y++ {
			panels[x][y] = panelInfo{
				typ: t,
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
	for i := 0; i < typeMax; i++ {
		dxlib.DeleteGraph(imgPanel[i])
		imgPanel[i] = -1
	}
	logger.Info("Successfully cleanuped battle field data")
}

// Draw ...
func Draw() {
	for x := 0; x < FieldNumX; x++ {
		for y := 0; y < FieldNumY; y++ {
			img := imgPanel[panels[x][y].typ]
			dxlib.DrawGraph(int32(PanelSizeX*x), int32(DrawPanelTopY+PanelSizeY*y), img, dxlib.TRUE)
		}
	}
}

// UpdateObjectPos ...
func UpdateObjectPos(positions []ObjectPosition) {
	// Cleanup at first
	for x := 0; x < FieldNumX; x++ {
		for y := 0; y < FieldNumY; y++ {
			panels[x][y].objectID = ""
		}
	}

	for _, pos := range positions {
		panels[pos.X][pos.Y].objectID = pos.ID
	}
}
