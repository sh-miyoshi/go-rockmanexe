package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type bambooLance struct {
	ID  string
	Arg skillcore.Argument

	count  int
	drawer skilldraw.DrawBamboolance
}

func newBambooLance(objID string, arg skillcore.Argument) *bambooLance {
	res := &bambooLance{
		ID:  objID,
		Arg: arg,
	}
	res.drawer.Init()

	return res
}

func (p *bambooLance) Draw() {
	p.drawer.Draw(p.count)
}

func (p *bambooLance) Process() (bool, error) {
	p.count++

	if p.count == 5 {
		dm := damage.Damage{
			DamageType:    damage.TypePosition,
			Pos:           point.Point{X: battlecommon.FieldNum.X - 1},
			Power:         int(p.Arg.Power),
			TTL:           5,
			TargetObjType: p.Arg.TargetType,
			HitEffectType: resources.EffectTypeBambooHit,
			ShowHitArea:   false,
			BigDamage:     true,
			PushLeft:      1,
			Element:       damage.ElementWood,
		}
		for y := 0; y < battlecommon.FieldNum.Y; y++ {
			dm.Pos.Y = y
			localanim.DamageManager().New(dm)
		}
	}

	if p.count > 10 {
		return true, nil
	}

	return false, nil
}

func (p *bambooLance) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *bambooLance) StopByOwner() {
	// Nothing to do after throwing
}
