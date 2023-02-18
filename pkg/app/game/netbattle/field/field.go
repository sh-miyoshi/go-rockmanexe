package field

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/background"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	battlefield "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	netconn "github.com/sh-miyoshi/go-rockmanexe/pkg/app/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	netconfig "github.com/sh-miyoshi/go-rockmanexe/pkg/net/config"
)

type Field struct {
	imgPanel [battlecommon.PanelStatusMax][2]int
}

func New() (*Field, error) {
	logger.Info("Initialize battle field data")

	res := &Field{}

	// TODO: Serverから取得する
	if err := background.Set(background.Type秋原町); err != nil {
		return nil, fmt.Errorf("failed to load background: %w", err)
	}

	// Initialize images
	files := [battlecommon.PanelStatusMax]string{"normal", "crack", "hole"}
	for i := 0; i < battlecommon.PanelStatusMax; i++ {
		fname := fmt.Sprintf("%sbattle/panel_player_%s.png", common.ImagePath, files[i])
		res.imgPanel[i][battlefield.PanelTypePlayer] = dxlib.LoadGraph(fname)
		if res.imgPanel[i][battlefield.PanelTypePlayer] < 0 {
			return nil, fmt.Errorf("failed to read player panel image %s", fname)
		}
	}
	for i := 0; i < battlecommon.PanelStatusMax; i++ {
		fname := fmt.Sprintf("%sbattle/panel_enemy_%s.png", common.ImagePath, files[i])
		res.imgPanel[i][battlefield.PanelTypeEnemy] = dxlib.LoadGraph(fname)
		if res.imgPanel[i][battlefield.PanelTypeEnemy] < 0 {
			return nil, fmt.Errorf("failed to read enemy panel image %s", fname)
		}
	}

	logger.Info("Successfully initialized battle field data")
	return res, nil
}

func (f *Field) End() {
	background.Unset()

	for i := 0; i < battlecommon.PanelStatusMax; i++ {
		for j := 0; j < 2; j++ {
			dxlib.DeleteGraph(f.imgPanel[i][j])
			f.imgPanel[i][j] = -1
		}
	}
}

func (f *Field) Draw() {
	clientID := config.Get().Net.ClientID

	panels := netconn.GetInst().GetGameInfo().Panels
	for x := 0; x < netconfig.FieldNumX; x++ {
		for y := 0; y < netconfig.FieldNumY; y++ {
			typ := battlefield.PanelTypePlayer
			if panels[x][y].OwnerClientID != clientID {
				typ = battlefield.PanelTypeEnemy
			}
			img := f.imgPanel[panels[x][y].Status][typ]
			vx := battlecommon.PanelSize.X * x
			vy := battlefield.DrawPanelTopY + battlecommon.PanelSize.Y*y

			dxlib.DrawGraph(vx, vy, img, true)
		}
	}
}

func GetPanelInfo(pos common.Point) battlefield.PanelInfo {
	ginfo := netconn.GetInst().GetGameInfo()
	clientID := config.Get().Net.ClientID

	id := ""
	for _, obj := range ginfo.Objects {
		if obj.Hittable && obj.X == pos.X && obj.Y == pos.Y {
			id = obj.ID
			break
		}
	}

	pnType := battlefield.PanelTypePlayer
	if ginfo.Panels[pos.X][pos.Y].OwnerClientID != clientID {
		pnType = battlefield.PanelTypeEnemy
	}

	return battlefield.PanelInfo{
		Type:     pnType,
		ObjectID: id,

		// TODO 未実装
		// Status    int
		// HoleCount int
	}
}
