package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
)

type dreamSword struct {
	ID  string
	Arg Argument

	count  int
	drawer skilldraw.DrawDreamSword
}

func newDreamSword(objID string, arg Argument) *dreamSword {
	return &dreamSword{
		ID:  objID,
		Arg: arg,
	}
}

func (p *dreamSword) Draw() {
	pos := localanim.ObjAnimGetObjPos(p.Arg.OwnerID)
	view := battlecommon.ViewPos(pos)

	p.drawer.Draw(view, p.count)
}

func (p *dreamSword) Process() (bool, error) {
	p.count++

	if p.count == 1*resources.SkillSwordDelay {
		sound.On(resources.SEDreamSword)

		userPos := localanim.ObjAnimGetObjPos(p.Arg.OwnerID)
		for x := 1; x <= 2; x++ {
			for y := -1; y <= 1; y++ {
				if objID := field.GetPanelInfo(common.Point{X: userPos.X + x, Y: userPos.Y + y}).ObjectID; objID != "" {
					dm := damage.Damage{
						DamageType:    damage.TypeObject,
						Power:         int(p.Arg.Power),
						TargetObjType: p.Arg.TargetType,
						HitEffectType: resources.EffectTypeNone,
						BigDamage:     true,
						Element:       damage.ElementNone,
						TargetObjID:   objID,
					}
					localanim.DamageManager().New(dm)
				}
			}
		}
	}

	if p.count > resources.SkillSwordEndCount {
		return true, nil
	}
	return false, nil
}

func (p *dreamSword) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *dreamSword) StopByOwner() {
	localanim.AnimDelete(p.ID)
}
