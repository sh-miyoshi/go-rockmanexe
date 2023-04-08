package draw

import (
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/net"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/router/anim"
)

type animDraw struct {
	images    [anim.TypeMax][]int
	imgDelays [anim.TypeMax]int
}

func (d *animDraw) Init() error {
	return nil
}

func (d *animDraw) End() {

}

func (d *animDraw) Draw() {
	ginfo := net.GetInst().GetGameInfo()
	for _, a := range ginfo.Anims {
		pos := battlecommon.ViewPos(a.Pos)
		ino := a.ActCount / d.imgDelays[a.AnimType]
		if ino >= len(d.images[a.AnimType]) {
			ino = len(d.images[a.AnimType]) - 1
		}
		// TODO offset

		dxlib.DrawRotaGraph(pos.X, pos.Y, 1, 0, d.images[a.AnimType][ino], true)

	}
}
