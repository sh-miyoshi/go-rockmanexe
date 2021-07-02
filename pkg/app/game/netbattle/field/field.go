package field

import (
	"fmt"
	"time"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
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

func Draw(playerID string) {
	info, err := netconn.GetFieldInfo()
	if err != nil {
		logger.Error("Failed to get field info: %v", err)
		// TODO error handling
		return
	}

	clientID := config.Get().Net.ClientID

	// TODO update
	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			vx := int32(appfield.PanelSizeX * x)
			vy := int32(appfield.DrawPanelTopY + appfield.PanelSizeY*y)
			dxlib.DrawGraph(vx, vy, imgPanel[0], dxlib.TRUE)
			vx = int32(appfield.PanelSizeX * (x + 3))
			dxlib.DrawGraph(vx, vy, imgPanel[1], dxlib.TRUE)
		}
	}

	for _, obj := range info.Objects {
		reverse := false

		if obj.ClientID != clientID {
			// enemy object
			reverse = true
		}

		viewHP := 0
		if obj.ID != playerID {
			viewHP = obj.HP
		}

		tm := info.CurrentTime.Sub(obj.BaseTime)
		cnt := tm * 60 / time.Second
		imgNo := int(cnt) / field.ImageDelays[obj.Type]
		draw.Object(obj.Type, imgNo, obj.X, obj.Y, draw.Option{
			Reverse:  reverse,
			ViewOfsX: obj.ViewOfsX,
			ViewOfsY: obj.ViewOfsY,
			ViewHP:   viewHP,
		})
	}
}

func GetPanelInfo(x, y int) appfield.PanelInfo {
	info, _ := netconn.GetFieldInfo()
	clientID := config.Get().Net.ClientID

	// TODO update
	if x < 3 {
		// player panel
		id := ""
		for _, obj := range info.Objects {
			if obj.ClientID == clientID && obj.X == x && obj.Y == y {
				id = obj.ID
				break
			}
		}

		return appfield.PanelInfo{
			Type:     appfield.PanelTypePlayer,
			ObjectID: id,
		}
	} else {
		// enemy panel
		id := ""
		for _, obj := range info.Objects {
			if obj.ClientID != clientID && obj.X == x-3 && obj.Y == y {
				id = obj.ID
				break
			}
		}

		return appfield.PanelInfo{
			Type:     appfield.PanelTypeEnemy,
			ObjectID: id,
		}
	}
}
