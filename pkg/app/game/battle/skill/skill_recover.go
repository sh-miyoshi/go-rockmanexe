package skill

import (
	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
)

type recover struct {
	ID         string
	OwnerID    string
	Power      uint
	TargetType int

	count int
}

func (p *recover) Draw() {
	n := (p.count / delayRecover) % len(imgRecover)
	if n >= 0 {
		px, py := objanim.GetObjPos(p.OwnerID)
		x, y := battlecommon.ViewPos(px, py)
		dxlib.DrawRotaGraph(x, y, 1, 0, imgRecover[n], dxlib.TRUE)
	}
}

func (p *recover) Process() (bool, error) {
	if p.count == 0 {
		sound.On(sound.SERecover)
		px, py := objanim.GetObjPos(p.OwnerID)
		damage.New(damage.Damage{
			PosX:          px,
			PosY:          py,
			Power:         -int(p.Power),
			TTL:           1,
			TargetType:    p.TargetType,
			HitEffectType: effect.TypeNone,
		})
	}

	p.count++

	if p.count > len(imgRecover)*delayRecover {
		return true, nil
	}
	return false, nil
}

func (p *recover) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		AnimType: anim.AnimTypeEffect,
	}
}
