package draw

import (
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/net"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

func Draw() {
	ginfo := net.GetInst().GetGameInfo()
	for _, obj := range ginfo.Objects {
		pos := battlecommon.ViewPos(obj.Pos)
		dxlib.DrawGraph(pos.X, pos.Y, 0, true)
	}
}
