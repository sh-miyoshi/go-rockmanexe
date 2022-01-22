package effect

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/effect"
)

var (
	effects []effect.Effect
)

func Draw() {
	for _, eff := range effects {
		p := common.Point{X: int32(eff.X), Y: int32(eff.Y)}
		draw.Effect(eff.Type, eff.Count, p, eff.ViewOfsX, eff.ViewOfsY)
	}
}

func Process() {
	finfo := netconn.GetFieldInfo()
	effects = append(effects, finfo.Effects...)
	netconn.RemoveEffects()

	newList := []effect.Effect{}
	for _, eff := range effects {
		eff.Count++
		n, delay := draw.GetEffectImageInfo(eff.Type)
		if eff.Count < n*delay {
			newList = append(newList, eff)
		}
	}
	effects = newList
}
