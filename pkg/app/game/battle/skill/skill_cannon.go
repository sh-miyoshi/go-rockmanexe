package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

const (
	TypeNormalCannon int = iota
	TypeHighCannon
	TypeMegaCannon

	TypeCannonMax
)

const (
	delayCannonAtk  = 2
	delayCannonBody = 6
)

type cannon struct {
	ID   string
	Type int
	Arg  Argument

	count int
}

func newCannon(objID string, cannonType int, arg Argument) *cannon {
	return &cannon{
		ID:   objID,
		Type: cannonType,
		Arg:  arg,
	}
}

func (p *cannon) Draw() {
	pos := objanim.GetObjPos(p.Arg.OwnerID)
	view := battlecommon.ViewPos(pos)

	n := p.count / delayCannonBody
	if n < len(imgCannonBody[p.Type]) {
		if n >= 3 {
			view.X -= 15
		}

		dxlib.DrawRotaGraph(view.X+48, view.Y-12, 1, 0, imgCannonBody[p.Type][n], true)
	}

	n = (p.count - 15) / delayCannonAtk
	if n >= 0 && n < len(imgCannonAtk[p.Type]) {
		dxlib.DrawRotaGraph(view.X+90, view.Y-10, 1, 0, imgCannonAtk[p.Type][n], true)
	}
}

func (p *cannon) Process() (bool, error) {
	p.count++

	if p.count == 20 {
		sound.On(sound.SECannon)
		pos := objanim.GetObjPos(p.Arg.OwnerID)
		dm := damage.Damage{
			Pos:           pos,
			Power:         int(p.Arg.Power),
			TTL:           1,
			TargetType:    p.Arg.TargetType,
			HitEffectType: effect.TypeCannonHit,
			BigDamage:     true,
		}

		if p.Arg.TargetType == damage.TargetEnemy {
			for x := pos.X + 1; x < field.FieldNum.X; x++ {
				dm.Pos.X = x
				if field.GetPanelInfo(common.Point{X: x, Y: dm.Pos.Y}).ObjectID != "" {
					damage.New(dm)
					break
				}
			}
		} else {
			for x := pos.X - 1; x >= 0; x-- {
				dm.Pos.X = x
				if field.GetPanelInfo(common.Point{X: x, Y: dm.Pos.Y}).ObjectID != "" {
					damage.New(dm)
					break
				}
			}
		}
	}

	max := len(imgCannonBody[p.Type]) * delayCannonBody
	if max < len(imgCannonAtk[p.Type])*delayCannonAtk+15 {
		max = len(imgCannonAtk[p.Type])*delayCannonAtk + 15
	}

	if p.count > max {
		return true, nil
	}
	return false, nil
}

func (p *cannon) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		AnimType: anim.AnimTypeSkill,
	}
}

func (p *cannon) AtDelete() {
	if p.Arg.RemoveObject != nil {
		p.Arg.RemoveObject(p.ID)
	}
}
