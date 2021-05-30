package field

import (
	"fmt"
	"time"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	appfield "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/field"
)

var (
	imgPanel = [2]int32{-1, -1}
)

func Init() error {
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
	info, err := netconn.GetFieldInfo()
	if err != nil {
		logger.Error("Failed to get field info: %v", err)
		// TODO error handling
		return
	}

	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			// My Panel
			vx := int32(appfield.PanelSizeX * x)
			vy := int32(appfield.DrawPanelTopY + appfield.PanelSizeY*y)
			dxlib.DrawGraph(vx, vy, imgPanel[0], dxlib.TRUE)

			// Show objects in my panel
			if info.MyArea[x][y].ID != "" {
				obj := info.MyArea[x][y]
				tm := info.CurrentTime.Sub(obj.BaseTime)
				cnt := tm * 60 / time.Second
				imgNo := int(cnt) / field.ImageDelays[obj.Type]
				draw.Object(obj.Type, imgNo, x, y, false)
			}

			// Enemy Panel
			vx = int32(appfield.PanelSizeX * (x + 3))
			vy = int32(appfield.DrawPanelTopY + appfield.PanelSizeY*y)
			dxlib.DrawGraph(vx, vy, imgPanel[1], dxlib.TRUE)

			// Show objects in enemy panel
			if info.EnemyArea[x][y].ID != "" {
				obj := info.EnemyArea[x][y]
				tm := info.CurrentTime.Sub(obj.BaseTime)
				cnt := tm * 60 / time.Second
				imgNo := int(cnt) / field.ImageDelays[obj.Type]
				draw.Object(info.EnemyArea[x][y].Type, imgNo, x+3, y, true)
			}
		}
	}
}

func GetPanelInfo(x, y int) appfield.PanelInfo {
	info, _ := netconn.GetFieldInfo()

	if x < 3 {
		// player panel
		return appfield.PanelInfo{
			Type:     appfield.PanelTypePlayer,
			ObjectID: info.MyArea[x][y].ID,
		}
	} else {
		// enemy panel
		return appfield.PanelInfo{
			Type:     appfield.PanelTypeEnemy,
			ObjectID: info.EnemyArea[x-3][y].ID,
		}
	}
}
