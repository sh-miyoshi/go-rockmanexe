package field

import (
	"fmt"
	"sort"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	appfield "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/netconn"
	netconfig "github.com/sh-miyoshi/go-rockmanexe/pkg/net/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/object"
)

var (
	imgPanel = [2]int32{-1, -1}
)

func Init() error {
	// Initialize images
	fname := common.ImagePath + "battle/panel_player_normal.png"
	imgPanel[appfield.PanelTypePlayer] = dxlib.LoadGraph(fname)
	if imgPanel[appfield.PanelTypePlayer] < 0 {
		return fmt.Errorf("failed to read player panel image %s", fname)
	}
	fname = common.ImagePath + "battle/panel_enemy_normal.png"
	imgPanel[appfield.PanelTypeEnemy] = dxlib.LoadGraph(fname)
	if imgPanel[appfield.PanelTypeEnemy] < 0 {
		return fmt.Errorf("failed to read enemy panel image %s", fname)
	}

	return nil
}

func Draw(playerID string) {
	finfo := netconn.GetFieldInfo()
	clientID := config.Get().Net.ClientID

	for x := 0; x < netconfig.FieldNumX; x++ {
		for y := 0; y < netconfig.FieldNumY; y++ {
			vx := int32(appfield.PanelSizeX * x)
			vy := int32(appfield.DrawPanelTopY + appfield.PanelSizeY*y)
			pn := imgPanel[0]
			if finfo.Panels[x][y].OwnerClientID != clientID {
				pn = imgPanel[1]
			}

			dxlib.DrawGraph(vx, vy, pn, dxlib.TRUE)

			if finfo.Panels[x][y].ShowHitArea {
				x1 := vx
				y1 := vy
				x2 := vx + appfield.PanelSizeX
				y2 := vy + appfield.PanelSizeY
				const s = 5
				dxlib.DrawBox(x1+s, y1+s, x2-s, y2-s, 0xffff00, dxlib.TRUE)
			}
		}
	}

	objects := append([]object.Object{}, finfo.Objects...)
	sort.Slice(objects, func(i, j int) bool {
		ii := objects[i].Y*appfield.FieldNumX + objects[i].X
		ij := objects[j].Y*appfield.FieldNumX + objects[j].X
		return ii < ij
	})
	for _, obj := range objects {
		reverse := false

		if obj.ClientID != clientID {
			// enemy object
			reverse = true
		}

		viewHP := 0
		if obj.ID != playerID {
			viewHP = obj.HP
		}

		draw.Object(obj, draw.Option{
			Reverse:  reverse,
			ViewHP:   viewHP,
			ViewChip: obj.ID != playerID,
		})
	}

	effect.Draw()
}

func GetPanelInfo(x, y int) appfield.PanelInfo {
	finfo := netconn.GetFieldInfo()
	clientID := config.Get().Net.ClientID

	id := ""
	for _, obj := range finfo.Objects {
		if obj.Hittable && obj.X == x && obj.Y == y {
			id = obj.ID
			break
		}
	}

	pnType := appfield.PanelTypePlayer
	if finfo.Panels[x][y].OwnerClientID != clientID {
		pnType = appfield.PanelTypeEnemy
	}

	return appfield.PanelInfo{
		Type:     pnType,
		ObjectID: id,
	}
}
