package skill

import (
	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
)

type aquamanShot struct {
	ID         string
	OwnerID    string
	Power      uint
	TargetType int

	x       int32
	y       int32
	ofsX    int32
	ofsY    int32
	targetX int
	targetY int
	count   int
}

func newAquamanShot(objID string, arg Argument) *aquamanShot {
	px, py := objanim.GetObjPos(arg.OwnerID)
	x, y := battlecommon.ViewPos(px, py)

	return &aquamanShot{
		ID:         objID,
		OwnerID:    arg.OwnerID,
		Power:      arg.Power,
		TargetType: arg.TargetType,
		x:          x - 40,
		y:          y + 10,
		targetX:    px - 2,
		targetY:    py,
	}
}

func (p *aquamanShot) Draw() {
	dxlib.DrawRotaGraph(p.x+p.ofsX, p.y+p.ofsY, 1, 0, imgAquamanShot[0], dxlib.TRUE)
}

func (p *aquamanShot) Process() (bool, error) {
	p.count++

	if p.count == 1 {
		sound.On(sound.SEBombThrow)
	}

	const size = 120
	p.ofsX -= 6
	p.ofsY = 10*p.ofsX*p.ofsX/(size*size) - 20*p.ofsX/size

	if p.ofsX < -size {
		pn := field.GetPanelInfo(p.targetX, p.targetY)
		if pn.Status == field.PanelStatusHole {
			return true, nil
		}

		sound.On(sound.SEWaterLanding)
		anim.New(effect.Get(effect.TypeWaterBomb, p.targetX, p.targetY, 0))
		anim.New(effect.Get(effect.TypeWaterBomb, p.targetX-1, p.targetY, 0))
		damage.New(damage.Damage{
			PosX:          p.targetX,
			PosY:          p.targetY,
			Power:         int(p.Power),
			TTL:           20,
			TargetType:    p.TargetType,
			HitEffectType: effect.TypeNone,
			BigDamage:     true,
		})
		damage.New(damage.Damage{
			PosX:          p.targetX - 1,
			PosY:          p.targetY,
			Power:         int(p.Power),
			TTL:           20,
			TargetType:    p.TargetType,
			HitEffectType: effect.TypeNone,
			BigDamage:     true,
		})

		return true, nil
	}
	return false, nil
}

func (p *aquamanShot) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		AnimType: anim.AnimTypeSkill,
	}
}
