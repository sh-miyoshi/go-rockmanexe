package field

import (
	"fmt"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	appfield "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/field"
)

var (
	info     *field.Info
	imgPanel = [2]int32{-1, -1}
)

func Init(fieldInfo *field.Info) error {
	info = fieldInfo

	// Initialize images
	fname := common.ImagePath + "battle/panel_player.png"
	imgPanel[appfield.PanelTypePlayer] = dxlib.LoadGraph(fname)
	if imgPanel[appfield.PanelTypePlayer] < 0 {
		return fmt.Errorf("failed to read player panel image %s", fname)
	}
	fname = common.ImagePath + "battle/panel_enemy.png"
	imgPanel[appfield.PanelTypeEnemy] = dxlib.LoadGraph(fname)
	if imgPanel[appfield.PanelTypeEnemy] < 0 {
		return fmt.Errorf("failed to read enemy panel image %s", fname)
	}

	return nil
}

func Draw() {
	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			// My Panel
			vx := int32(appfield.PanelSizeX * x)
			vy := int32(appfield.DrawPanelTopY + appfield.PanelSizeY*y)
			dxlib.DrawGraph(vx, vy, imgPanel[0], dxlib.TRUE)

			// Enemy Panel
			vx = int32(appfield.PanelSizeX * (x + 3))
			vy = int32(appfield.DrawPanelTopY + appfield.PanelSizeY*y)
			dxlib.DrawGraph(vx, vy, imgPanel[1], dxlib.TRUE)
		}
	}
}
