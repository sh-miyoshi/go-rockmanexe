package draw

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/oldnet/effect"
)

func drawEffect(images [effect.TypeMax][]int, imgNo int, eff effect.Effect) {
	view := battlecommon.ViewPos(common.Point{X: eff.X, Y: eff.Y})
	view.X += eff.ViewOfsX
	view.Y += eff.ViewOfsY

	dxlib.DrawRotaGraph(view.X, view.Y, 1, 0, images[eff.Type][imgNo], true)
}
