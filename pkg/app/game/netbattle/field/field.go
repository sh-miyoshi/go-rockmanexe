package field

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/background"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	battlefield "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/net"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
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

	panels := net.GetInst().GetGameInfo().Panels
	for x := 0; x < battlecommon.FieldNum.X; x++ {
		for y := 0; y < battlecommon.FieldNum.Y; y++ {
			typ := battlefield.PanelTypePlayer
			if panels[x][y].OwnerClientID != clientID {
				typ = battlefield.PanelTypeEnemy
			}
			// TODO: statusを反映させる
			img := f.imgPanel[battlecommon.PanelStatusNormal][typ]
			vx := battlecommon.PanelSize.X * x
			vy := battlecommon.DrawPanelTopY + battlecommon.PanelSize.Y*y

			dxlib.DrawGraph(vx, vy, img, true)
		}
	}
}
