package draw

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/net"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

type animDraw struct {
}

func (d *animDraw) Init() error {
	return nil
}

func (d *animDraw) End() {

}

func (d *animDraw) Draw() {
	ginfo := net.GetInst().GetGameInfo()
	for _, a := range ginfo.Anims {
		dxlib.DrawFormatString(0, 100, 0xff0000, "(%d, %d): %d", a.Pos.X, a.Pos.Y, a.DrawType)
	}
}
