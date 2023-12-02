package skill

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/object"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type aquaman struct {
	ID  string
	Arg Argument

	count  int
	state  int
	pos    point.Point
	atkID  string
	drawer skilldraw.DrawAquaman
}

func newAquaman(objID string, arg Argument) *aquaman {
	return &aquaman{
		ID:    objID,
		Arg:   arg,
		state: resources.SkillAquamanStateInit,
		pos:   localanim.ObjAnimGetObjPos(arg.OwnerID),
	}
}

func (p *aquaman) Draw() {
	view := battlecommon.ViewPos(p.pos)
	p.drawer.Draw(view, p.count, p.state)
}

func (p *aquaman) Process() (bool, error) {
	switch p.state {
	case resources.SkillAquamanStateInit:
		field.SetBlackoutCount(300)
		SetChipNameDraw("アクアマン", true)
		p.setState(resources.SkillAquamanStateAppear)
		return false, nil
	case resources.SkillAquamanStateAppear:
		if p.count == 70 {
			p.setState(resources.SkillAquamanStateCreatePipe)
			return false, nil
		}
	case resources.SkillAquamanStateCreatePipe:
		if p.count == 10 {
			obj := &object.WaterPipe{}
			pm := object.ObjectParam{
				Pos:           point.Point{X: p.pos.X + 1, Y: p.pos.Y},
				HP:            500,
				OnwerCharType: objanim.ObjTypePlayer,
				AttackNum:     1,
				Interval:      50,
				Power:         int(p.Arg.Power),
			}
			if err := obj.Init(p.ID, pm); err != nil {
				return false, fmt.Errorf("water pipe create failed: %w", err)
			}
			p.atkID = localanim.ObjAnimNew(obj)
			localanim.ObjAnimAddActiveAnim(p.atkID)

			p.setState(resources.SkillAquamanStateAttack)
			return false, nil
		}
	case resources.SkillAquamanStateAttack:
		if !localanim.ObjAnimIsProcessing(p.atkID) {
			field.SetBlackoutCount(0)
			return true, nil
		}
	}

	p.count++
	return false, nil
}

func (p *aquaman) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *aquaman) StopByOwner() {
	// Nothing to do after throwing
}

func (p *aquaman) setState(nextState int) {
	p.count = 0
	p.state = nextState
}
