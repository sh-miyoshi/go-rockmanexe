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
	PanelTypePlayer int = iota
	PanelTypeEnemy

	panelTypeMax
)

// ObjectPosition ...
type ObjectPosition struct {
	ID   string
	X, Y int
}

type PanelInfo struct {
	Type     int
	ObjectID string
	// TODO status(毒とか穴とか)
}

var (
	imgPanel   = [2]int32{-1, -1}
	imgHPFrame int32
	panels     [FieldNumX][FieldNumY]PanelInfo
)

// Init ...
func Init() error {
	logger.Info("Initialize battle field data")

	// Initialize images
	fname := common.ImagePath + "battle/panel_player.png"
	imgPanel[PanelTypePlayer] = dxlib.LoadGraph(fname)
	if imgPanel[PanelTypePlayer] < 0 {
		return fmt.Errorf("Failed to read player panel image %s", fname)
	}
	fname = common.ImagePath + "battle/panel_enemy.png"
	imgPanel[PanelTypeEnemy] = dxlib.LoadGraph(fname)
	if imgPanel[PanelTypeEnemy] < 0 {
		return fmt.Errorf("Failed to read enemy panel image %s", fname)
	}
	fname = common.ImagePath + "battle/hp_frame.png"
	imgHPFrame = dxlib.LoadGraph(fname)
	if imgHPFrame < 0 {
		return fmt.Errorf("Failed to read hp frame image %s", fname)
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
			dxlib.DrawGraph(int32(PanelSizeX*x), int32(DrawPanelTopY+PanelSizeY*y), img, dxlib.TRUE)
		}
	}
}

func DrawFrame(hpFrameX, hpFrameY int32) {
	// TODO Custom Gauge

	// HP Frame
	dxlib.DrawGraph(hpFrameX, hpFrameY, imgHPFrame, dxlib.TRUE)

	// TODO Mind Frame
}

// UpdateObjectPos ...
func UpdateObjectPos(positions []ObjectPosition) {
	// Cleanup at first
	for x := 0; x < FieldNumX; x++ {
		for y := 0; y < FieldNumY; y++ {
			panels[x][y].ObjectID = ""
		}
	}

	for _, pos := range positions {
		panels[pos.X][pos.Y].ObjectID = pos.ID
	}
}

func GetPos(objID string) (x, y int) {
	for x := 0; x < FieldNumX; x++ {
		for y := 0; y < FieldNumY; y++ {
			if panels[x][y].ObjectID == objID {
				return x, y
			}
		}
	}
	return -1, -1
}

func GetPanelInfo(x, y int) PanelInfo {
	return panels[x][y]
}
