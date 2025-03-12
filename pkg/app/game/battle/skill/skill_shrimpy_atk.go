package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/manager"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	shrimpyAtkNextStepCount = 20
)

type shrimpyAtk struct {
	ID      string
	Arg     skillcore.Argument
	drawer  skilldraw.DrawShrimpyAtk
	pos     point.Point
	count   int
	state   int
	animMgr *manager.Manager
}

func newShrimpyAtk(objID string, arg skillcore.Argument, animMgr *manager.Manager) *shrimpyAtk {
	pos := animMgr.ObjAnimGetObjPos(arg.OwnerID)
	pos.X--
	return &shrimpyAtk{
		ID:      objID,
		Arg:     arg,
		pos:     pos,
		state:   resources.SkillShrimpyAttackStateBegin,
		animMgr: animMgr,
	}
}

func (p *shrimpyAtk) Draw() {
	p.drawer.Draw(p.pos, p.count, shrimpyAtkNextStepCount, p.state)
}

func (p *shrimpyAtk) Update() (bool, error) {
	switch p.state {
	case resources.SkillShrimpyAttackStateBegin:
		if p.drawer.IsDrawEnd(p.count, p.state) {
			p.state = resources.SkillShrimpyAttackStateMove
			p.count = 0
			return false, nil
		}
	case resources.SkillShrimpyAttackStateMove:
		if p.count%shrimpyAtkNextStepCount == 0 {
			p.animMgr.DamageManager().New(damage.Damage{
				DamageType:    damage.TypePosition,
				Pos:           p.pos,
				Power:         int(p.Arg.Power),
				TTL:           shrimpyAtkNextStepCount,
				TargetObjType: p.Arg.TargetType,
				ShowHitArea:   true,
				StrengthType:  damage.StrengthHigh,
				Element:       damage.ElementWater,
			})
		}
		if p.count%shrimpyAtkNextStepCount == shrimpyAtkNextStepCount-1 {
			p.pos.X--
			if p.pos.X < 0 {
				return true, nil
			}
		}
	}
	p.count++
	return false, nil
}

func (p *shrimpyAtk) GetParam() anim.Param {
	return anim.Param{
		ObjID: p.ID,
	}
}

func (p *shrimpyAtk) StopByOwner() {
	p.animMgr.AnimDelete(p.ID)
}
