package battle

import (
	"fmt"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

const (
	fieldNumX  = 6
	fieldNumY  = 3
	panelSizeX = 80
	panelSizeY = 50
)

const (
	typePlayer int = iota
	typeEnemy
	typeMax
)

type panelInfo struct {
	typ int
	// TODO status(毒とか穴とか)
}

var (
	imgPanel [2]int32 = [2]int32{-1, -1}
	panels   [fieldNumX][fieldNumY]panelInfo
)

func fieldInit() error {
	logger.Info("Initialize battle field data")

	// Initialize images
	fname := "data/images/battle/panel_player.png"
	imgPanel[typePlayer] = dxlib.LoadGraph(fname)
	if imgPanel[typePlayer] < 0 {
		return fmt.Errorf("Failed to read player panel image %s", fname)
	}
	fname = "data/images/battle/panel_enemy.png"
	imgPanel[typeEnemy] = dxlib.LoadGraph(fname)
	if imgPanel[typeEnemy] < 0 {
		return fmt.Errorf("Failed to read enemy panel image %s", fname)
	}

	// Initialize panel info
	for x := 0; x < fieldNumX; x++ {
		t := typePlayer
		if x > 2 {
			t = typeEnemy
		}
		for y := 0; y < fieldNumY; y++ {
			panels[x][y] = panelInfo{
				typ: t,
			}
		}
	}
	// TODO: special field

	logger.Info("Successfully initialized battle field data")
	return nil
}

func fieldEnd() {
	logger.Info("Cleanup battle field data")
	for i := 0; i < typeMax; i++ {
		dxlib.DeleteGraph(imgPanel[i])
	}
	logger.Info("Successfully cleanuped battle field data")
}

func fieldDraw() {
	const drawPanelTopY = common.ScreenY - (panelSizeY * 3) - 30

	for x := 0; x < fieldNumX; x++ {
		for y := 0; y < fieldNumY; y++ {
			img := imgPanel[panels[x][y].typ]
			dxlib.DrawGraph(int32(panelSizeX*x), int32(drawPanelTopY+panelSizeY*y), img, dxlib.TRUE)
		}
	}
}
