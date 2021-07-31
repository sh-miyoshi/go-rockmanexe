package effect

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/effect"
)

var (
	effects []effect.Effect
)

func Draw() {
	for _, eff := range effects {
		draw.Effect(eff.Type, eff.Count, eff.X, eff.Y, eff.ViewOfsX, eff.ViewOfsY)
	}
}

func Process() {
	finfo, _ := netconn.GetFieldInfo()
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
